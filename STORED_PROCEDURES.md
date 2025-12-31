```sql
-- 1) Create an order from the user's cart
CREATE OR REPLACE FUNCTION create_order_from_cart(p_user_id BIGINT)
RETURNS BIGINT
LANGUAGE plpgsql
AS $$
DECLARE
  v_order_id      BIGINT;
  v_total_price   NUMERIC(12,2) := 0;
  v_eta           INT := 1;
  v_now           TIMESTAMP := NOW();
BEGIN
  -- lock cart items and ensure not empty
  PERFORM 1 FROM cart_items WHERE cart_id = (SELECT id FROM carts WHERE user_id = p_user_id) FOR UPDATE;
  IF NOT FOUND THEN
    RAISE EXCEPTION 'Cannot create order: empty cart for user_id=%', p_user_id;
  END IF;

  INSERT INTO orders (user_id, status, total_price, estimated_production_time_days, payment_method, payment_status, created_at)
  VALUES (p_user_id, 'new', 0, 1, 'stripe', 'pending', v_now)
  RETURNING id INTO v_order_id;

  INSERT INTO order_items (order_id, product_id, quantity, unit_price, line_total, selected_options_json, calculated_production_time_days, created_at)
  SELECT
    v_order_id,
    p.id,
    ci.quantity,
    p.base_price,
    ROUND(p.base_price * GREATEST(ci.quantity, 1), 2),
    COALESCE(ci.selected_options_json, '[]'),
    GREATEST(1, p.base_production_time_days),
    v_now
  FROM cart_items ci
  JOIN carts c ON c.id = ci.cart_id AND c.user_id = p_user_id
  JOIN products p ON p.id = ci.product_id;

  SELECT COALESCE(ROUND(SUM(oi.line_total), 2), 0) INTO v_total_price
  FROM order_items oi WHERE oi.order_id = v_order_id;

  SELECT calculate_order_eta_for_order(v_order_id) INTO v_eta;

  UPDATE orders
     SET total_price = v_total_price,
         estimated_production_time_days = COALESCE(v_eta, 1)
   WHERE id = v_order_id;

  -- increment popularity counters
  INSERT INTO recommendation_counters(product_id, count)
  SELECT oi.product_id, oi.quantity FROM order_items oi WHERE oi.order_id = v_order_id
  ON CONFLICT (product_id) DO UPDATE SET count = recommendation_counters.count + EXCLUDED.count;

  -- clear cart
  DELETE FROM cart_items WHERE cart_id = (SELECT id FROM carts WHERE user_id = p_user_id);

  RETURN v_order_id;
END;
$$;

-- 2) Calculate ETA for an order: MAX(base_time_with_modifiers * qty)
CREATE OR REPLACE FUNCTION calculate_order_eta_for_order(p_order_id BIGINT)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
  v_eta INT := 1;
BEGIN
  SELECT GREATEST(1, COALESCE(MAX( (p.base_production_time_days) * oi.quantity ), 1))
    INTO v_eta
  FROM order_items oi
  JOIN products p ON p.id = oi.product_id
  WHERE oi.order_id = p_order_id;

  -- Optional: add buffer based on production load
  -- SELECT v_eta + CASE WHEN (SELECT COUNT(1) FROM orders WHERE status IN ('new','processing')) > 20 THEN 2 ELSE 0 END INTO v_eta;

  RETURN v_eta;
END;
$$;

-- 3) Process payment result from Stripe webhook
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
  IF COALESCE(p_transaction_id,'') = '' THEN
    RAISE EXCEPTION 'Missing transaction_id for order %', p_order_id;
  END IF;

  INSERT INTO payments (order_id, status, amount, created_at, transaction_id)
  VALUES (p_order_id, v_effective, p_amount, NOW(), p_transaction_id)
  ON CONFLICT DO NOTHING;

  IF v_effective IN ('succeeded','success','paid') THEN
    UPDATE orders SET payment_status = 'paid', status = 'paid' WHERE id = p_order_id;
  ELSE
    UPDATE orders SET payment_status = 'cancelled', status = 'cancelled' WHERE id = p_order_id;
  END IF;

  SELECT status INTO v_order_status FROM orders WHERE id = p_order_id;
  RETURN COALESCE(v_order_status, 'unknown');
END;
$$;
```
