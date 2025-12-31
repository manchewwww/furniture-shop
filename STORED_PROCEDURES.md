# Stored Procedures (PostgreSQL, PL/pgSQL)

This file contains ready‑to‑run stored procedures that implement core furniture‑shop business logic. They are designed to work with a realistic schema and can be applied via a SQL migration.

Assumed tables (simplified; adjust column names if yours differ):
- users(id, email, role)
- products(id, name, category_id, base_production_days, base_price)
- cart_items(id, user_id, product_id, quantity)
- orders(id, user_id, status, total_price, eta_days, created_at)
- order_items(id, order_id, product_id, quantity, unit_price, line_total)
- payments(id, order_id, status, amount, created_at)
- product_recommendations(product_id PRIMARY KEY, popularity_score)

Status pipeline used: NEW → PROCESSING → PAID → COMPLETED; CANCELLED (terminal).

## SQL

```sql
-- 1) Create an order from a user's cart
-- Creates a new order in status NEW, copies items from cart_items to order_items,
-- computes total_price and ETA (via calculate_order_eta_for_order), clears the cart,
-- and updates product recommendation popularity (by quantity).
CREATE OR REPLACE FUNCTION create_order_from_cart(p_user_id BIGINT)
RETURNS BIGINT
LANGUAGE plpgsql
AS $$
DECLARE
  v_order_id BIGINT;
  v_total_price NUMERIC := 0;
  v_eta INT := 1;
  v_now TIMESTAMP := NOW();
BEGIN
  -- Ensure there are items in the user's cart
  IF NOT EXISTS (
    SELECT 1 FROM cart_items WHERE user_id = p_user_id
  ) THEN
    RAISE EXCEPTION 'Cart is empty for user_id=%', p_user_id;
  END IF;

  -- Insert the ORDER header (status NEW; total and eta set below)
  INSERT INTO orders (user_id, status, total_price, eta_days, created_at)
  VALUES (p_user_id, 'NEW', 0, 0, v_now)
  RETURNING id INTO v_order_id;

  -- Copy lines from cart_items to order_items; compute totals
  INSERT INTO order_items (order_id, product_id, quantity, unit_price, line_total)
  SELECT
    v_order_id,
    p.id,
    ci.quantity,
    p.base_price,
    p.base_price * ci.quantity
  FROM cart_items ci
  JOIN products p ON p.id = ci.product_id
  WHERE ci.user_id = p_user_id;

  -- Sum total price
  SELECT COALESCE(SUM(oi.line_total), 0)
    INTO v_total_price
  FROM order_items oi
  WHERE oi.order_id = v_order_id;

  -- Calculate ETA (max(product.base_production_days * quantity); min 1)
  SELECT calculate_order_eta_for_order(v_order_id) INTO v_eta;

  -- Update order totals
  UPDATE orders
     SET total_price = v_total_price,
         eta_days    = COALESCE(v_eta, 1)
   WHERE id = v_order_id;

  -- Increment popularity by quantities for all items
  PERFORM increment_product_popularity(oi.product_id, oi.quantity)
  FROM order_items oi
  WHERE oi.order_id = v_order_id;

  -- Clear the user's cart
  DELETE FROM cart_items WHERE user_id = p_user_id;

  RETURN v_order_id;
END;
$$;


-- 2) Calculate ETA for an order
-- ETA = MAX(product.base_production_days * order_items.quantity), at least 1 day.
CREATE OR REPLACE FUNCTION calculate_order_eta_for_order(p_order_id BIGINT)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
  v_eta INT;
BEGIN
  SELECT GREATEST(1, COALESCE(MAX(p.base_production_days * oi.quantity), 1))
    INTO v_eta
  FROM order_items oi
  JOIN products p ON p.id = oi.product_id
  WHERE oi.order_id = p_order_id;

  RETURN v_eta;
END;
$$;


-- 3) Update order status with pipeline validation
-- Valid transitions:
--   NEW -> PROCESSING
--   PROCESSING -> PAID
--   PAID -> COMPLETED
--   ANY (if not COMPLETED) -> CANCELLED
CREATE OR REPLACE FUNCTION update_order_status(p_order_id BIGINT, p_new_status TEXT)
RETURNS BOOLEAN
LANGUAGE plpgsql
AS $$
DECLARE
  v_old_status TEXT;
  v_allowed BOOLEAN := FALSE;
BEGIN
  SELECT status INTO v_old_status FROM orders WHERE id = p_order_id FOR UPDATE;
  IF v_old_status IS NULL THEN
    RAISE EXCEPTION 'Order % not found', p_order_id;
  END IF;

  -- Validate transitions
  IF p_new_status = 'CANCELLED' AND v_old_status <> 'COMPLETED' THEN
    v_allowed := TRUE;
  ELSIF v_old_status = 'NEW' AND p_new_status = 'PROCESSING' THEN
    v_allowed := TRUE;
  ELSIF v_old_status = 'PROCESSING' AND p_new_status = 'PAID' THEN
    v_allowed := TRUE;
  ELSIF v_old_status = 'PAID' AND p_new_status = 'COMPLETED' THEN
    v_allowed := TRUE;
  END IF;

  IF NOT v_allowed THEN
    RAISE EXCEPTION 'Invalid status transition: % -> %', v_old_status, p_new_status;
  END IF;

  UPDATE orders SET status = p_new_status WHERE id = p_order_id;
  RETURN TRUE;
END;
$$;


-- 4) Process payment result
-- p_payment_status: 'success'/'paid' or 'failure'/'declined'/'cancelled'
-- On success: upsert payment row and set order status to PAID (if currently PROCESSING or NEW then PROCESSING→PAID).
-- On failure/cancel: set order status to CANCELLED.
CREATE OR REPLACE FUNCTION process_payment_result(
  p_order_id BIGINT,
  p_payment_status TEXT,
  p_amount NUMERIC
)
RETURNS TEXT
LANGUAGE plpgsql
AS $$
DECLARE
  v_effective TEXT := LOWER(p_payment_status);
  v_final_order_status TEXT;
BEGIN
  -- Persist payment record (simple insert; customize as needed)
  INSERT INTO payments (order_id, status, amount, created_at)
  VALUES (p_order_id, v_effective, p_amount, NOW());

  IF v_effective IN ('success', 'paid') THEN
    -- Move order to PROCESSING if needed, then to PAID
    BEGIN
      PERFORM update_order_status(p_order_id, 'PROCESSING');
    EXCEPTION WHEN OTHERS THEN
      -- ignore if already past PROCESSING
    END;
    PERFORM update_order_status(p_order_id, 'PAID');
    v_final_order_status := 'PAID';
  ELSE
    -- Failure or cancel
    PERFORM update_order_status(p_order_id, 'CANCELLED');
    v_final_order_status := 'CANCELLED';
  END IF;

  RETURN v_final_order_status;
END;
$$;


-- 5) Increment product recommendation popularity (UPSERT)
-- Increases popularity_score by p_delta (default 1). Creates the row if missing.
CREATE OR REPLACE FUNCTION increment_product_popularity(
  p_product_id BIGINT,
  p_delta INT DEFAULT 1
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
  v_new_score INT;
BEGIN
  -- Upsert-like update
  LOOP
    UPDATE product_recommendations
       SET popularity_score = popularity_score + p_delta
     WHERE product_id = p_product_id
     RETURNING popularity_score INTO v_new_score;

    EXIT WHEN FOUND;

    -- Not found – try insert
    BEGIN
      INSERT INTO product_recommendations(product_id, popularity_score)
      VALUES (p_product_id, p_delta)
      RETURNING popularity_score INTO v_new_score;
      EXIT;
    EXCEPTION WHEN unique_violation THEN
      -- Someone else inserted concurrently; retry
    END;
  END LOOP;

  RETURN v_new_score;
END;
$$;


-- 6) Utility: recalculate order totals (optional helper)
-- Recomputes order total_price and eta_days from order items/products.
CREATE OR REPLACE FUNCTION recalc_order_totals(p_order_id BIGINT)
RETURNS TABLE(total_price NUMERIC, eta_days INT)
LANGUAGE plpgsql
AS $$
BEGIN
  RETURN QUERY
  WITH totals AS (
    SELECT
      COALESCE(SUM(oi.line_total), 0) AS total_price,
      GREATEST(1, COALESCE(MAX(p.base_production_days * oi.quantity), 1)) AS eta_days
    FROM order_items oi
    JOIN products p ON p.id = oi.product_id
    WHERE oi.order_id = p_order_id
  )
  UPDATE orders o
     SET total_price = t.total_price,
         eta_days = t.eta_days
    FROM totals t
    WHERE o.id = p_order_id
  RETURNING t.total_price, t.eta_days;
END;
$$;
```

## Usage notes
- Execute this file in your PostgreSQL database (psql or migration tool).
- If your actual columns differ (e.g., products.price vs base_price), adjust the SELECT/INSERT/UPDATE parts accordingly.
- These procedures are independent from the ORM and can be called by the REST layer when needed.