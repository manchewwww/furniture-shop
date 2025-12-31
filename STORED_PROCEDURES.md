```sql

CREATE OR REPLACE FUNCTION create_order_from_cart(p_user_id BIGINT)
RETURNS BIGINT
LANGUAGE plpgsql
AS $$
DECLARE
  v_order_id      BIGINT;
  v_total_price   NUMERIC(12,2) := 0;
  v_eta           INT := 1;
  v_now           TIMESTAMP := NOW();
  v_has_items     BOOLEAN;
BEGIN
  PERFORM 1
  FROM cart_items
  WHERE user_id = p_user_id
  FOR UPDATE;

  SELECT EXISTS(SELECT 1 FROM cart_items WHERE user_id = p_user_id) INTO v_has_items;
  IF NOT v_has_items THEN
    RAISE EXCEPTION 'Cannot create order: cart is empty for user_id=%', p_user_id;
  END IF;

  INSERT INTO orders (user_id, status, total_price, eta_days, created_at)
  VALUES (p_user_id, 'NEW', 0, 0, v_now)
  RETURNING id INTO v_order_id;

  INSERT INTO order_items (order_id, product_id, quantity, unit_price, line_total)
  SELECT
    v_order_id,
    p.id,
    p.base_price,
    ROUND(p.base_price * GREATEST(ci.quantity, 1), 2)
  FROM cart_items ci
  JOIN products p ON p.id = ci.product_id
  WHERE ci.user_id = p_user_id;

  SELECT COALESCE(ROUND(SUM(oi.line_total), 2), 0)
    INTO v_total_price
  FROM order_items oi
  WHERE oi.order_id = v_order_id;

  SELECT calculate_order_eta_for_order(v_order_id) INTO v_eta;

  UPDATE orders
     SET total_price = v_total_price,
         eta_days    = COALESCE(v_eta, 1)
   WHERE id = v_order_id;

  PERFORM increment_product_popularity(oi.product_id, oi.quantity)
  FROM order_items oi
  WHERE oi.order_id = v_order_id;

  DELETE FROM cart_items WHERE user_id = p_user_id;

  RETURN v_order_id;
END;
$$;


CREATE OR REPLACE FUNCTION calculate_order_eta_for_order(p_order_id BIGINT)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
  v_eta INT := 1;
BEGIN
  SELECT GREATEST(1, COALESCE(MAX(p.base_production_days * oi.quantity), 1))
    INTO v_eta
  FROM order_items oi
  JOIN products p ON p.id = oi.product_id
  WHERE oi.order_id = p_order_id;

  RETURN v_eta;
END;
$$;


CREATE OR REPLACE FUNCTION update_order_status(p_order_id BIGINT, p_new_status TEXT)
RETURNS BOOLEAN
LANGUAGE plpgsql
AS $$
DECLARE
  v_old_status TEXT;
  v_allowed    BOOLEAN := FALSE;
BEGIN
  SELECT status
    INTO v_old_status
  FROM orders
  WHERE id = p_order_id
  FOR UPDATE;

  IF v_old_status IS NULL THEN
    RAISE EXCEPTION 'Order % not found', p_order_id;
  END IF;

  IF v_old_status = p_new_status THEN
    RETURN TRUE;
  END IF;

  IF v_old_status = 'COMPLETED' THEN
    RAISE EXCEPTION 'Cannot change status of COMPLETED order %', p_order_id;
  END IF;

  IF p_new_status = 'CANCELLED' THEN
  ELSIF v_old_status = 'NEW' AND p_new_status = 'PROCESSING' THEN
    v_allowed := TRUE;
  ELSIF v_old_status = 'PROCESSING' AND p_new_status = 'PAID' THEN
    v_allowed := TRUE;
  ELSIF v_old_status = 'PAID' AND p_new_status = 'COMPLETED' THEN
    v_allowed := TRUE;
  END IF;

  IF NOT v_allowed THEN
    RAISE EXCEPTION 'Invalid status transition: % -> % for order %', v_old_status, p_new_status, p_order_id;
  END IF;

  UPDATE orders
     SET status = p_new_status
   WHERE id = p_order_id;

  RETURN TRUE;
END;
$$;


CREATE OR REPLACE FUNCTION process_payment_result(
  p_order_id       BIGINT,
  p_payment_status TEXT,
  p_amount         NUMERIC(12,2),
  p_transaction_id TEXT
)
RETURNS TEXT
LANGUAGE plpgsql
AS $$
DECLARE
  v_effective TEXT := LOWER(p_payment_status);
  v_order_status TEXT;
BEGIN
  IF p_transaction_id IS NULL OR length(p_transaction_id) = 0 THEN
    RAISE EXCEPTION 'Missing transaction_id for order %', p_order_id;
  END IF;

  BEGIN
    INSERT INTO payments (order_id, status, amount, created_at, transaction_id)
    VALUES (p_order_id, v_effective, p_amount, NOW(), p_transaction_id);
  EXCEPTION WHEN unique_violation THEN
    SELECT status INTO v_order_status FROM orders WHERE id = p_order_id;
    RETURN COALESCE(v_order_status, 'UNKNOWN');
  END;

  IF v_effective IN ('success','paid') THEN
    BEGIN
      PERFORM update_order_status(p_order_id, 'PROCESSING');
    EXCEPTION WHEN OTHERS THEN
    END;
    PERFORM update_order_status(p_order_id, 'PAID');
  ELSE
    PERFORM update_order_status(p_order_id, 'CANCELLED');
  END IF;

  SELECT status INTO v_order_status FROM orders WHERE id = p_order_id;
  RETURN COALESCE(v_order_status, 'UNKNOWN');
END;
$$;


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
  IF p_delta = 0 THEN
    RETURN NULL;
  END IF;

  LOOP
    UPDATE product_recommendations
       SET popularity_score = popularity_score + p_delta
     WHERE product_id = p_product_id
     RETURNING popularity_score INTO v_new_score;

    EXIT WHEN FOUND;

    BEGIN
      INSERT INTO product_recommendations(product_id, popularity_score)
      VALUES (p_product_id, p_delta)
      RETURNING popularity_score INTO v_new_score;
      EXIT;
    EXCEPTION WHEN unique_violation THEN
    END;
  END LOOP;

  RETURN v_new_score;
END;
$$;

CREATE OR REPLACE FUNCTION recalc_order_totals(p_order_id BIGINT)
RETURNS TABLE(total_price NUMERIC(12,2), eta_days INT)
LANGUAGE plpgsql
AS $$
BEGIN
  RETURN QUERY
  WITH totals AS (
    SELECT
      COALESCE(ROUND(SUM(oi.line_total), 2), 0)::NUMERIC(12,2) AS total_price,
      GREATEST(1, COALESCE(MAX(p.base_production_days * oi.quantity), 1))::INT AS eta_days
    FROM order_items oi
    JOIN products p ON p.id = oi.product_id
    WHERE oi.order_id = p_order_id
  )
  UPDATE orders o
     SET total_price = t.total_price,
         eta_days    = t.eta_days
    FROM totals t
   WHERE o.id = p_order_id
  RETURNING t.total_price, t.eta_days;
END;
$$;
```
