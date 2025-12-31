# Архитектура и дизайн

Този проект е пълноценен e-commerce за мебели с Go/Fiber бекенд и React/Vite фронтенд. Доменната йерархия е фиксирана: Отдел → Категория → Продукт.

## Бекенд (Go, Fiber, GORM, Postgres)
- Слоеве и модули
  - HTTP (маршрути/контролери): `internal/server/http/**`
  - Сервизи (доменна логика): `internal/service/**`
  - Хранилища/репозитории (GORM): `internal/storage/**`
  - Модели (ентитети): `internal/entities/**`; DTO: `internal/dtos/**`
  - Конфигурация/БД: `internal/config/**`, `internal/database/**`
- Ключови маршрути (извадка)
  - Каталог (публични): `GET /api/departments`, `GET /api/departments/:id/categories`, `GET /api/categories/:id/products`, `GET /api/products/:id`, `GET /api/products/:id/recommendations`, `GET /api/products/search?query=...`
  - Аутентикация: `POST /api/auth/register`, `POST /api/auth/login`, `GET /api/user/me`
  - Количка (JWT): `GET/PUT/DELETE /api/user/cart`, `POST/PATCH/DELETE /api/user/cart/items[:id]`
  - Поръчки (JWT): `POST /api/orders`, `GET /api/user/orders`, `GET /api/user/orders/:id`, `POST /api/user/orders/:id/pay`
  - Плащания (Stripe – симулация): `POST /api/webhooks/stripe`
  - Админ (JWT + роля admin): `GET/POST/PUT/DELETE /api/admin/departments|categories|products|product_options`, `POST /api/admin/upload`
- Сигурност: JWT (контекст: user_id, user_email, user_role), guard за админ.

## Фронтенд (React, Vite, TypeScript)
- Страници: Каталог (Отдели → Категории → Продукти), Търсене, Детайли на продукт (опции), Количка, Checkout, Мои поръчки, Админ (Отдели/Категории/Продукти/Поръчки).
- Контексти: `AuthContext` (JWT), `CartContext` (гост: localStorage; JWT: синхрон с бекенда).

## Бизнес правила
- Препоръки: броячи се увеличават при отваряне на продукт и при създаване на поръчка (по количество). Връщат топ N от същата категория, без текущия продукт.
- ETA: за артикул – базово време + модификатори от опции; за поръчка – максимум от артикулите + корекция според натоварването. Ако всички артикули са „налични“ на ниво продукт – ETA = 1 ден.
- Плащания: картова авторизация (симулирана чрез Stripe). Уебхук актуализира PaymentStatus и OrderStatus.

