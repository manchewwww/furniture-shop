# Иновации и интеграции

## Динамично ETA (Estimated Time of Arrival)

- Пер продукт: `ETA_product = base_production_time_days + Σ(option_time_modifiers)`.
- Пер поръчка: `ETA_order = MAX(ETA_product_i × quantity_i)` (паралелно производство).
- Буфер при натоварване: добавя се +X дни, когато броят активни поръчки е над праг (`orders_in_production`).
- Оптимизация: ако всички артикули са „в наличност“ → `ETA_order = 1` ден.

## Stripe интеграция

- Плащания се инициират от бекенда (Stripe Secret: `STRIPE_SECRET`).
- Webhook: `POST /api/webhooks/stripe` валидира подписа (`STRIPE_WEBHOOK_SECRET`) и обработва събития:
  - `payment_intent.succeeded` → запис в `payments`, статус на поръчка = `PAID`.
  - `payment_intent.payment_failed` → запис в `payments`, статус = `CANCELLED`.
- Съхранени процедури се използват за транзакционно обновяване на поръчката и тоталите.

## SendGrid (SMTP) имейли

- Конфигурира се чрез SMTP: `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS`.
- Събития:
  - Регистрация → имейл за приветствие.
  - Създадена поръчка → потвърждение с номер и стойност.
  - Промяна на статус → известие (напр. `PROCESSING`, `PAID`, `SHIPPED`).

## Система за препоръки

- Брояч по продукт в таблица `recommendation_counters` (`product_id`, `count`).
- Увеличава се при преглед/поръчка; препоръките извличат популярни продукти в същата категория (сортиране по `count DESC`).
