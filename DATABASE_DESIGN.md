# Дизайн на база данни

Таблици (ключови полета):
- departments: id, name, description
- categories: id, department_id (FK), name, description
- products: id, category_id (FK), name, short_description, long_description, base_price, base_production_time_days, image_url, base_material, default_width, default_height, default_depth, is_made_to_order
- product_options: id, product_id (FK), option_type, option_name, price_modifier_type, price_modifier_value, production_time_modifier_days, production_time_modifier_percent
- users: id, role, name, email (unique), password_hash, address, phone, created_at
- orders: id, user_id (FK), status, total_price, estimated_production_time_days, payment_method, payment_status, created_at, updated_at
- order_items: id, order_id (FK), product_id (FK), quantity, unit_price, line_total, calculated_production_time_days, selected_options_json
- stock: id, material_name, quantity_available, unit

Връзки: 1:N между отдели→категории, категории→продукти, продукти→опции; поръчки→позиции. GORM миграции и seed в `internal/database`.

