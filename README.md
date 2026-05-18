# Go Banking

Go Banking is a full-stack banking dashboard built with a Go REST API, PostgreSQL, and a Next.js frontend. It supports user registration and login, JWT-protected account management, deposits, withdrawals, transfers, transaction history, and a dashboard summary.

This README is the operational source of truth for running, testing, and preparing the project for production.

## Stack

| Layer | Technology |
| --- | --- |
| Backend | Go, Chi router, pgx, JWT, bcrypt |
| Database | PostgreSQL with Goose migrations |
| Frontend | Next.js App Router, React, TypeScript, Tailwind CSS |
| Auth | JWT access token stored by the frontend and sent as `Authorization: Bearer <token>` |

## Repository Layout

```text
.
+-- backend
|   +-- cmd/api                  # API entrypoint
|   +-- internal
|   |   +-- config               # Environment loading
|   |   +-- database             # PostgreSQL connection
|   |   +-- handlers             # HTTP handlers
|   |   +-- middleware           # Auth, CORS, logging, recovery
|   |   +-- models               # Request/response/domain structs
|   |   +-- repository           # Database access
|   |   +-- services             # Business logic
|   +-- migrations               # Goose SQL migrations
|   +-- Makefile
+-- frontend
    +-- src/app                  # Next.js routes
    +-- src/components           # Shared UI and layout
    +-- src/lib                  # Auth helpers and validation
    +-- src/services             # API client modules
    +-- src/types                # TypeScript API/domain types
```

## Prerequisites

- Go 1.26 or newer, matching `backend/go.mod`
- Node.js 20 or newer
- npm
- PostgreSQL 14 or newer
- Goose migration CLI

Install Goose if it is not available:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Make sure your Go bin directory is on `PATH`, commonly:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Environment Variables

### Backend

Create `backend/.env` from the example:

```bash
cp backend/.env.example backend/.env
```

Required values:

```env
APP_ENV=development
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/go_banking_api?sslmode=disable
JWT_SECRET=replace_with_a_long_random_secret
```

Production requirements:

- Use a strong `JWT_SECRET`, at least 32 random bytes.
- Use a least-privilege PostgreSQL user.
- Use SSL for remote PostgreSQL connections.
- Do not commit real `.env` files.

### Frontend

Create `frontend/.env.local`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

For production, set `NEXT_PUBLIC_API_URL` to the public API origin.

## Local Development

### 1. Create the Database

Create a PostgreSQL database that matches `DATABASE_URL`:

```bash
createdb go_banking_api
```

If you use a different database name, user, password, host, or port, update `backend/.env` and `backend/Makefile`.

### 2. Run Migrations

```bash
cd backend
make migrate-up
```

Useful migration commands:

```bash
make migrate-status
make migrate-down
```

### 3. Start the API

```bash
cd backend
make run
```

The API starts on `http://localhost:8080` by default.

Health check:

```bash
curl http://localhost:8080/health
```

### 4. Start the Frontend

In another terminal:

```bash
cd frontend
npm install
npm run dev
```

Open `http://localhost:3000`.

## Test and Quality Commands

Backend:

```bash
cd backend
go test ./...
```

If your environment has a read-only Go cache, use:

```bash
GOCACHE=/tmp/go-build go test ./...
```

Frontend:

```bash
cd frontend
npm run lint
npm run build
```

Run these before opening a pull request or deploying.

## API Overview

All responses use this envelope:

```json
{
  "success": true,
  "message": "operation completed",
  "data": {}
}
```

Errors use:

```json
{
  "success": false,
  "error": "error message"
}
```

Protected routes require:

```http
Authorization: Bearer <access_token>
```

### Public Routes

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/health` | API health check |
| `POST` | `/auth/register` | Register a user |
| `POST` | `/auth/login` | Login and receive an access token |

Register payload:

```json
{
  "name": "Example User",
  "email": "user@email.com",
  "password": "password123"
}
```

Login payload:

```json
{
  "email": "user@email.com",
  "password": "password123"
}
```

Login response data:

```json
{
  "access_token": "jwt_token",
  "user": {
    "id": 1,
    "name": "Example User",
    "email": "user@email.com",
    "created_at": "2026-05-18T00:00:00Z"
  }
}
```

### Protected Routes

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/auth/me` | Return the current authenticated user |
| `GET` | `/accounts/` | List current user's accounts |
| `POST` | `/accounts/` | Create an account |
| `GET` | `/accounts/{id}` | Get one account owned by the user |
| `POST` | `/accounts/{id}/deposit` | Deposit into an account |
| `POST` | `/accounts/{id}/withdraw` | Withdraw from an account |
| `GET` | `/accounts/{id}/transactions` | List transactions for one account |
| `GET` | `/transactions` | List all current user's transactions |
| `POST` | `/transfer` | Transfer between accounts |

Create account payload:

```json
{
  "name": "Main Savings",
  "account_number": "ACC1779083000001",
  "account_type": "savings",
  "balance": 1000,
  "currency": "NPR"
}
```

Deposit or withdraw payload:

```json
{
  "amount": 500
}
```

Transfer payload:

```json
{
  "from_account_id": 1,
  "to_account_id": 2,
  "amount": 250
}
```

## Frontend Routes

| Path | Description |
| --- | --- |
| `/` | Public landing page |
| `/register` | Registration form |
| `/login` | Login form |
| `/dashboard` | Authenticated dashboard |
| `/accounts` | Account list and creation |
| `/accounts/[id]` | Account detail and account transactions |
| `/transactions` | Transaction history and transfer workflow |

Protected frontend pages use `ProtectedRoute`, which redirects unauthenticated users to `/login`.

## Data Model

Core tables:

- `users`: registered users with unique email and bcrypt password hash.
- `accounts`: user-owned bank accounts with unique account numbers, type, balance, and currency.
- `transactions`: deposit, withdrawal, and transfer records with reference numbers.

Migration files live in `backend/migrations` and should be the only way schema changes are introduced.

## License

No license has been declared for this repository.
