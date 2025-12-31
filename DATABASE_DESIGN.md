# Дизайн на базата данни

## Основни таблици (GORM модели)
- departments: id, name, description, image_url, created_at, updated_at
- categories: id, department_id (FK), name, description, created_at, updated_at
- products: id, category_id (FK), name, short_description, long_description, base_price, base_production_time_days, image_url, base_material, default_width, default_height, default_depth, created_at, updated_at
- product_options: id, product_id (FK), option_type, option_name, price_modifier_type, price_modifier_value, production_time_modifier_days, production_time_modifier_percent
- recommendation_counters: id, product_id (UNIQUE), count

## Потребители, количка, поръчки
- users: id, role, name, email, address, phone, password_hash, created_at, updated_at
- carts: id, user_id (UNIQUE), created_at, updated_at
- cart_items: id, cart_id (FK), product_id (FK), quantity, selected_options_json, created_at, updated_at
- orders: id, user_id (FK), status, total_price, estimated_production_time_days, payment_method, payment_status, created_at, updated_at
- order_items: id, order_id (FK), product_id (FK), quantity, unit_price, line_total, selected_options_json, calculated_production_time_days, created_at, updated_at

## Индекси и бележки
- `recommendation_counters.product_id` (unique)
- `carts.user_id` (unique)
- Референции и каскади се управляват на приложно ниво чрез GORM.
- Каталогът е 1→N→N по веригата Department→Category→Product.

