# Архитектура и дизайн

Три-слойна архитектура: React (Vite, TS, Ant Design) → Go Fiber (HTTP API, домейн услуги, GORM) → PostgreSQL (модели и съхранени процедури).

## Backend (Go, Fiber, GORM, Postgres)

- Слоеве и директории:
  - HTTP: `internal/server/http/**` (маршрути, middleware, статични файлове `/uploads`).
  - Service: `internal/service/**` (бизнес логика, валидации).
  - Storage: `internal/storage/postgres/**` (репозитории и заявки).
  - Domain модели: `internal/entities/**`, DTO: `internal/dtos/**`.
  - Конфигурация/База: `internal/config/**`, `internal/database/**`.
- Основни публични ендпойнти (каталог):
  - `GET /api/departments`
  - `GET /api/departments/:id/categories`
  - `GET /api/categories/:id/products`
  - `GET /api/products/:id`
  - `GET /api/products/:id/recommendations`
  - `GET /api/products/search?query=...`
- Потребители и количка (JWT):
  - `POST /api/auth/register`, `POST /api/auth/login`, `GET /api/user/me`
  - `GET/PUT/DELETE /api/user/cart`, `POST/PATCH/DELETE /api/user/cart/items[:id]`
- Поръчки и плащания:
  - `POST /api/orders`, `GET /api/user/orders`, `GET /api/user/orders/:id`
  - `POST /api/user/orders/:id/pay` (Stripe)
  - `POST /api/webhooks/stripe` (Webhook)
- Админ (JWT + роля `admin`):
  - `GET/POST/PUT/DELETE /api/admin/departments|categories|products|product_options`
  - `POST /api/admin/upload` (качване на изображения в `uploads/`)

## Frontend (React, Vite, TypeScript)

- Страници: каталог, продукт детайл, количка, checkout, поръчки, админ.
- Състояние: `AuthContext` (JWT), `CartContext` (LocalStorage/DB sync).

## Сигурност

- JWT (подписан токен) в Authorization header. Претенции: `user_id`, `email`, `role`.
- Guard-ове по роля/идентификация за user/admin маршрути.
- Не се логват чувствителни данни; CORS конфигурация и валидация на входа.
