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

# Furniture Shop — Go + React

University course project (Stages 1–3) for a Furniture Shop with Departments → Categories → Products hierarchy. Backend in Go (Fiber, GORM, PostgreSQL), frontend in React (Vite, TS, Ant Design), payments via Stripe, and email via SendGrid (SMTP).

## Features

- Catalog browsing and search
- Shopping cart (guest in LocalStorage; user persisted in DB)
- Orders and Stripe payments
- Dynamic ETA calculation per order based on product options and production load
- Email notifications via SendGrid (registration, order created, status updates)

## Getting Started

### Prerequisites

- Go 1.21+
- Node 18+
- PostgreSQL 14+

### Environment

Create a `.env` in repo root:

```
# Database (either DSN or individual credentials per internal/config)
DB_DSN=host=localhost user=postgres password=postgres dbname=furniture port=5432 sslmode=disable

# Stripe
STRIPE_SECRET=sk_test_xxx
STRIPE_WEBHOOK_SECRET=whsec_xxx

# SendGrid (SMTP)
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASS=SG.xxxxx

# Optional seed reset
SEED_RESET=true
```

Note: The codebase can also read DB credentials via app config (`internal/config`) and env (`DB_USER`, `DB_PASSWORD`). Using `DB_DSN` is recommended for simplicity.

### Backend

- Run (dev): `go run ./cmd/server`
- Build: `go build -o bin/server ./cmd/server`
- Format: `go fmt ./...`
- Tests: `go test ./...`

Backend serves API under `/api` and static uploads under `/uploads`.

### Frontend

```
cd frontend
npm install
npm run dev     # start Vite dev server
npm run build   # production build
npm run preview # preview build
```

## Project Structure

- Backend entry: `cmd/server/main.go`
- Domain: `internal/entities`, Services: `internal/service`
- HTTP server/routes: `internal/server/http`
- Storage (Postgres): `internal/storage/postgres`
- Config: `internal/config`, seed: `internal/database/seed.go`
- Frontend: `frontend/src/**` (api, pages, components)

## Seeding

Auto-migration and initial seed run on startup. Use `SEED_RESET=true` to truncate and reseed. Images are linked from `./uploads` and served at `/uploads`.

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
