# Go Financial Planning

A web application for personal financial planning with:
- user registration and login
- future income/expense tracking
- monthly projection charts
- stateful authentication using cookie-based sessions

This is a **study project** focused on AI-assisted development with **Codex**. It was built for learning purposes and to test previously acquired knowledge in the Go language.

## Technologies

- Backend: Go + Chi + SQLite + Goose
- Frontend: React + Vite + Recharts
- Security: Argon2id (password hashing), `HttpOnly` cookie, database-backed session storage

## Prerequisites

Before running the project, install:

- Go 1.24+ (or a version compatible with `go.mod`)
- Node.js 20+ and npm

Check your versions:

```bash
go version
node -v
npm -v
```

## Project Structure

```text
backend/
  cmd/api/main.go
  internal/
    app/
    auth/
    db/
    handlers/
    repository/
frontend/
  src/
```

## Run the Project (Step by Step)

Open **2 terminals** at the project root.

### 1) Backend (terminal 1)

```bash
cd backend
go mod tidy
go run ./cmd/api/main.go
```

If it starts correctly, you should see something like:

```text
goose: no migrations to run. current version: 1
api listening on :8080
```

Notes:
- Migrations run automatically on startup.
- The SQLite database is created at `backend/finance.db` by default.
- The project uses a pure Go SQLite driver (`modernc.org/sqlite`), so CGO is not required.

### 2) Frontend (terminal 2)

```bash
cd frontend
npm install
npm run dev
```

Then open:

- `http://localhost:5173`

## Environment Variables

### Backend

- `API_ADDR` (default: `:8080`)
- `DB_PATH` (default: `./finance.db`)
- `FRONTEND_ORIGIN` (default: `http://localhost:5173`)
- `APP_PRODUCTION` (default: `false`)
- `SESSION_TTL` (default: `168h`)

PowerShell example:

```powershell
$env:API_ADDR=":8090"
$env:FRONTEND_ORIGIN="http://localhost:5173"
go run ./cmd/api/main.go
```

### Frontend

- `VITE_API_URL` (default: `http://localhost:8080/api/v1`)

If backend runs on a different port (for example, `8090`):

```powershell
$env:VITE_API_URL="http://localhost:8090/api/v1"
npm run dev
```

## API Endpoints

### Healthcheck

- `GET /healthz`

### Authentication

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/me`

### Transactions

- `GET /api/v1/transactions`
- `POST /api/v1/transactions`
- `PUT /api/v1/transactions/{id}`
- `DELETE /api/v1/transactions/{id}`

### Forecast

- `GET /api/v1/forecast/monthly?months=12`

## Payload Rules

- `kind`: `income` or `expense`
- `amount`: positive decimal (example: `1299.90`)
- `dueDate`: `YYYY-MM-DD`
- Login accepts only `email` and `password`
- Register accepts `name`, `email`, and `password`

## Troubleshooting

### 1) Port already in use (`bind: ... 8080`)

Run backend on another port:

```powershell
$env:API_ADDR=":8090"
go run ./cmd/api/main.go
```

Then update frontend:

```powershell
$env:VITE_API_URL="http://localhost:8090/api/v1"
npm run dev
```

### 2) CORS error

Make sure `FRONTEND_ORIGIN` in backend matches the exact frontend URL (including port).

### 3) `unknown field: name` during login

Use only this login payload:

```json
{
  "email": "user@example.com",
  "password": "password1234"
}
```

### 4) Frontend not updating

Stop and run again:

```bash
npm run dev
```

If needed, rebuild:

```bash
cd frontend
npm run build
```

## Implemented Security

- Password hashing with Argon2id
- Stateful sessions with random token and SHA-256 hash stored in database
- `HttpOnly` cookie with `SameSite=Lax`
- Security headers (`X-Frame-Options`, `X-Content-Type-Options`, `CSP`)
