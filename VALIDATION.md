# VALIDATION (Изискване → Доказателство)

| Изискване | Статус | Доказателство | Бележка |
|---|---|---|---|
| Каталог в БД; Отдел→Категория→Продукт | PASS | Ентитети: `internal/entities/catalog/{department.go,category.go,product.go}`; Роути: `internal/server/http/handler/catalog/routes.go` | Йерархията е константна.
| Търсене в каталог | PASS | `GET /api/products/search`; имплементация: `internal/storage/postgres/catalog/product_repository.go: Search` | ILIKE по име/описание.
| Админ CRUD (Отдели/Категории/Продукти/Опции) | PASS | Роути: `internal/server/http/handler/admin/routes.go`; UI: `frontend/src/pages/Admin{Departments,Categories,Products}.tsx` | Добавени и екрани за редакция.
| Seed: ≥2–3 отдела, ≥2–3 категории/отдел, ≥8–10 продукта/категория | PASS | `internal/database/seed.go` (5 отдела; 5 категории/отдел; 10 продукта/категория) | `SEED_RESET=true` за чист старт.
| Количка се пази в локална БД (сървъра) | PASS | API: `GET/PUT/DELETE /api/user/cart`, `POST/PATCH/DELETE /api/user/cart/items`; модели: `internal/entities/orders/cart.go` | Гост – LocalStorage; JWT – DB.
| Редакция/показване на количка | PASS | UI: `frontend/src/pages/Cart.tsx`; контекст: `frontend/src/store/CartContext.tsx` | Сумиране и ETA визуализация.
| Обработка на заявки/поръчки | PASS | `POST /api/orders`; сервиз: `internal/service/domain/orders/service.go` | Цена и ETA от опции.
| Проста, ефективна препоръчка | PASS | Брояч при преглед/поръчка; `GET /api/products/:id/recommendations` | Сортиране по популярност в категория.
| Билинг/плащания (симулация) | PASS | Stripe webhook: `internal/server/http/handler/payments/handler.go`; сервиз: `internal/service/domain/payments/service.go` | Статуси paid/declined/cancelled и преходи.
| Картова авторизация (симулирана) | PASS | `POST /api/user/orders/:id/pay` + webhook | Е2Е поток документиран.
| Разширен конвейер на поръчка | PASS | Статуси: new→processing→in_production→shipped→delivered/cancelled; енум: `internal/entities/orders/enums.go` | Преходи през уебхук/админ.
