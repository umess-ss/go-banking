# Go Banking Backend

Go Banking backend is a Go HTTP API built with Chi, PostgreSQL, pgx, JWT, and bcrypt.

For architecture details, read:

```text
ARCHITECTURE.md
```

## Package Layout

```text
cmd/api
+-- main.go                 # API entrypoint and dependency wiring

internal
+-- account                 # Account model, repository, service, handler, routes
+-- auth                    # User/auth model, repository, service, handler, routes
+-- config                  # Environment loading and validation
+-- database                # PostgreSQL connection pool
+-- health                  # Liveness and readiness checks
+-- logger                  # slog logger setup
+-- middleware              # Auth, CORS, request logging, panic recovery
+-- response                # Shared JSON response helpers
+-- transaction             # Transaction model, repository, service, handler, routes

migrations                  # Goose SQL migrations
```

## Requirements

- Go 1.26 or newer
- PostgreSQL
- Goose migration CLI

Install Goose:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Environment

Create a backend env file:

```bash
cp .env.example .env
```

Required values:

```env
APP_ENV=development
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/go_banking_api?sslmode=disable
JWT_SECRET=replace_with_a_long_random_secret
```

## Database

Create the database:

```bash
createdb go_banking_api
```

Run migrations:

```bash
make migrate-up
```

Other migration commands:

```bash
make migrate-status
make migrate-down
```

## Run

```bash
make run
```

The API starts on:

```text
http://localhost:8080
```

## Test

```bash
go test ./...
```

If the Go build cache is read-only in your environment:

```bash
GOCACHE=/tmp/go-build go test ./...
```

## Routes

Public routes:

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/health` | Liveness check |
| `GET` | `/ready` | Readiness check with database ping |
| `POST` | `/auth/register` | Register a user |
| `POST` | `/auth/login` | Login and receive a JWT |

Protected routes require:

```http
Authorization: Bearer <token>
```

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/auth/me` | Current authenticated user |
| `GET` | `/accounts/` | List accounts |
| `POST` | `/accounts/` | Create account |
| `GET` | `/accounts/{id}` | Get account |
| `POST` | `/accounts/{id}/deposit` | Deposit money |
| `POST` | `/accounts/{id}/withdraw` | Withdraw money |
| `GET` | `/accounts/{id}/transactions` | Account transactions |
| `GET` | `/transactions` | User transaction history |
| `POST` | `/transfer` | Transfer between accounts |

## Response Format

Success:

```json
{
  "success": true,
  "message": "operation completed",
  "data": {}
}
```

Error:

```json
{
  "success": false,
  "error": "error message"
}
```
