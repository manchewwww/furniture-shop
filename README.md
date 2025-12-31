# Furniture Shop

A full‑stack demo storefront for configurable furniture.
- Backend: Go (Fiber, GORM, PostgreSQL), Stripe payments
- Frontend: React + Vite + TypeScript, Ant Design UI

## Project Structure
- Backend: `cmd/server` (entry), domain in `internal/{server,service,storage,entities,dtos,config,database,validation}`
- Frontend: `frontend/src/{pages,api,store,components}`, entry `frontend/src/main.tsx`, assets `frontend/public/`
- Config: `internal/config/appconfig.json` and `.env` (copy from `.env.example`)
- Docs: `ARCHITECTURE_AND_DESIGN.md`, `DATABASE_DESIGN.md`, `FUNCTIONALITY_AND_USE_CASES.md`

## Prerequisites
- Go 1.21+
- Node.js 18+ and npm
- PostgreSQL 14+
- Stripe account (Secret + Webhook secret)

## Setup
1. Copy environment and set required values:
   - `cp .env.example .env` then set `DB_USER`, `DB_PASSWORD`, `JWT_SECRET`, `STRIPE_SECRET_KEY`, `STRIPE_WEBHOOK_SECRET`.
2. Configure database and CORS in `internal/config/appconfig.json` (or set `APP_CONFIG` to a custom path).
3. Frontend API base (optional): set `VITE_API_URL` (defaults to `http://localhost:8080/api`).
4. SMTP (optional for email notifications): set `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS`, `FROM_EMAIL`. If not set, emails are logged.

## Run (Local)
- Backend API (port 8080): `go run ./cmd/server`
  - Health: `GET /health`
  - API base: `/api` (serves uploads under `/uploads`)
- Frontend dev (port 5173):
  - `cd frontend && npm install`
  - `npm run dev`

### Seed DB
Auto-migration runs on startup. Initial seed inserts 3 departments, 6 categories, and 8–10 products per category, plus an admin user.

Admin credentials: email `admin@example.com`, password `admin123`.

## Build & Test
- Backend build: `go build -o bin/server ./cmd/server`
- Backend tests: `go test -cover ./...`
- Frontend build: `npm --prefix frontend run build`
- Frontend tests: `npm --prefix frontend run test`

## Payments (Stripe)
- Checkout sessions created on order (card method) and for re-pay: `/api/user/orders/:id/pay`
- Webhook: `POST /api/webhooks/stripe` (set your Stripe endpoint to this URL)
- On success, the frontend redirects to `/orders?open=<order_id>` and expands the created order.
 - Local testing via Stripe CLI:
   - `stripe listen --forward-to localhost:8080/api/webhooks/stripe`
   - Use the printed webhook secret as `STRIPE_WEBHOOK_SECRET`.

## Contributing & Conventions
- Follow DTO placement under `internal/dtos/<feature>/` (one DTO per file, lower‑case package names).
- Code style: Go `go fmt`/`go vet`; TS components PascalCase; 2‑space indentation.
- Commits: Conventional Commits; tests near changed code.
- See AGENTS.md for detailed repository guidelines.
