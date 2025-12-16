# Архитектура и дизайн

## Backend (Go + Fiber + GORM + PostgreSQL)
- Слоеве: `handlers` (HTTP), `services` (бизнес логика), `database` (връзка, миграции/seed), `models` (GORM модели), `routes` (регистрация на маршрути), `middleware` (JWT).
- REST API под `/api/*`. CORS е активиран за React клиента.
- JWT за удостоверяване; роли: `client`, `admin`.

## Frontend (React + Ant Design)
- Vite + TypeScript, React Router, AntD компоненти.
- Всички текстове са на български. Комуникация с бекенда чрез Axios към `/api`.

## API обзор
- Публично: каталога, търсене, детайли, препоръки, създаване на поръчка, плащане с карта (симулация).
- Потребител: `GET /api/user/orders`, `GET /api/user/orders/:id`.
- Админ: CRUD за отдели/категории/продукти/опции; поръчки и смяна на статус.

