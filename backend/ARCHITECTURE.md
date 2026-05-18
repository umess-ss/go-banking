# Go Banking Backend Architecture

## Purpose

The backend is a Go HTTP API for a banking dashboard. It handles authentication, account management, money movement, transaction history, and health checks.

The application is currently a modular monolith: one deployable Go binary, organized by domain packages. This keeps local development simple while making the codebase easier to split later if any domain needs to become a separate service.

## High-Level Shape

```text
HTTP client
  |
  v
Chi router
  |
  +-- global middleware
  |   +-- request logging
  |   +-- panic recovery
  |   +-- CORS
  |
  +-- public routes
  |   +-- /health
  |   +-- /ready
  |   +-- /auth/register
  |   +-- /auth/login
  |
  +-- protected routes
      +-- AuthMiddleware
      +-- /auth/me
      +-- /accounts
      +-- /accounts/{id}
      +-- /accounts/{id}/deposit
      +-- /accounts/{id}/withdraw
      +-- /accounts/{id}/transactions
      +-- /transactions
      +-- /transfer
```

## Runtime Entry Point

The application starts in:

```text
cmd/api/main.go
```

`main.go` is responsible for:

- Loading environment configuration.
- Creating the structured logger.
- Validating required config.
- Opening the PostgreSQL connection pool.
- Creating repositories, services, and handlers.
- Registering domain routes.
- Installing global middleware.
- Starting the HTTP server.

The entry point should stay thin. Domain rules belong in service packages, database queries belong in repository packages, and HTTP request/response concerns belong in handlers.

## Package Layout

```text
internal
+-- account
|   +-- model.go       # Account domain model
|   +-- repository.go  # Account SQL queries and transfer transaction
|   +-- service.go     # Account business rules
|   +-- handler.go     # Account HTTP handlers
|   +-- routes.go      # Account route registration
+-- auth
|   +-- model.go       # User/auth request and response models
|   +-- repository.go  # User SQL queries
|   +-- service.go     # Registration, login, JWT generation
|   +-- handler.go     # Auth HTTP handlers
|   +-- routers.go     # Auth route registration
+-- transaction
|   +-- model.go       # Transaction request and domain models
|   +-- repository.go  # Transaction SQL queries
|   +-- service.go     # Transaction read operations
|   +-- handler.go     # Transaction HTTP handlers
|   +-- routes.go      # Transaction route registration
+-- health
|   +-- handler.go     # Liveness and readiness handlers
|   +-- routes.go      # Health route registration
+-- config
|   +-- config.go      # Environment loading and validation
+-- database
|   +-- postgres.go    # PostgreSQL pool creation
+-- logger
|   +-- logger.go      # slog logger setup
+-- middleware
|   +-- auth_middleware.go
|   +-- cors.go
|   +-- logger_middleware.go
+-- response
    +-- response.go    # Shared JSON response helpers
```

## Layering Rules

Each domain follows this dependency direction:

```text
routes -> handler -> service -> repository -> database
```

The intended responsibilities are:

- `routes.go`: attach handlers to URL paths and add route-specific middleware.
- `handler.go`: decode requests, read route params/context, call services, write JSON responses.
- `service.go`: validate inputs and enforce business rules.
- `repository.go`: execute SQL and map rows into domain models.
- `model.go`: define request, response, and persisted data shapes for the domain.

Handlers should not contain SQL. Repositories should not know about HTTP. Services should not write responses.

## Domains

### Auth

Package:

```text
internal/auth
```

Responsibilities:

- Register users.
- Normalize email addresses.
- Hash passwords with bcrypt.
- Validate login credentials.
- Generate JWT access tokens.
- Return the current user for `/auth/me`.

Public routes:

```text
POST /auth/register
POST /auth/login
```

Protected routes:

```text
GET /auth/me
```

Key dependencies:

- `auth.UserRepository` for user persistence.
- `bcrypt` for password hashing and verification.
- `jwt/v5` for token signing.

### Account

Package:

```text
internal/account
```

Responsibilities:

- Create user-owned bank accounts.
- List accounts for the authenticated user.
- Fetch one account by ID and owner.
- Deposit into an account.
- Withdraw from an account.
- Transfer between accounts.

Protected routes:

```text
GET  /accounts/
POST /accounts/
GET  /accounts/{id}
POST /accounts/{id}/deposit
POST /accounts/{id}/withdraw
POST /transfer
```

Business rules:

- Account name is required.
- Account type defaults to `savings`.
- Allowed account types are `savings`, `checking`, and `current`.
- Initial balance cannot be negative.
- Deposit amount must be positive.
- Withdrawal amount must be greater than zero and cannot exceed balance.
- Transfer source and destination cannot be the same account.
- Transfer amount must be greater than zero.

Transfer behavior:

- Transfers run inside a database transaction.
- Source and destination account rows are selected `FOR UPDATE`.
- The source account is debited and the destination account is credited.
- A `transfer` transaction row is inserted before commit.

### Transaction

Package:

```text
internal/transaction
```

Responsibilities:

- Store deposit, withdrawal, and transfer records.
- Generate transaction reference numbers.
- List all transactions visible to the authenticated user.
- List transactions for a specific account owned by the authenticated user.

Protected routes:

```text
GET /transactions
GET /accounts/{id}/transactions
```

Transaction records use:

- `type`: `deposit`, `withdraw`, or `transfer`
- `status`: currently defaults to `success`
- `reference_number`: generated when missing

### Health

Package:

```text
internal/health
```

Responsibilities:

- Lightweight liveness check.
- Readiness check that verifies database connectivity.

Routes:

```text
GET /health
GET /ready
```

`/health` confirms the HTTP process is alive. `/ready` confirms the API can reach PostgreSQL.

## Authentication Flow

1. A user logs in through `POST /auth/login`.
2. The auth service verifies the bcrypt password hash.
3. The auth service signs a JWT containing:

```json
{
  "user_id": 1,
  "email": "user@email.com",
  "exp": 1234567890
}
```

4. The frontend sends the token with protected requests:

```http
Authorization: Bearer <token>
```

5. `middleware.AuthMiddleware` validates the token and stores `userID` in request context.
6. Protected handlers read the user ID with `middleware.GetUserIDFromContext`.

## Configuration

Configuration is loaded by:

```text
internal/config/config.go
```

Supported environment variables:

```text
APP_ENV       # development by default
PORT          # 8080 by default
DATABASE_URL  # required
JWT_SECRET    # required
```

`Config.Validate()` prevents the server from starting without `DATABASE_URL` and `JWT_SECRET`.

## Database

Database access uses `pgxpool`.

Connection setup lives in:

```text
internal/database/postgres.go
```

Schema changes live in:

```text
migrations/
```

Current core tables:

- `users`
- `accounts`
- `transactions`

Repositories own SQL. Services should call repository methods rather than embedding queries.

## Response Format

JSON response helpers live in:

```text
internal/response/response.go
```

Successful responses follow:

```json
{
  "success": true,
  "message": "operation completed",
  "data": {}
}
```

Error responses follow:

```json
{
  "success": false,
  "error": "error message"
}
```

## Middleware

Global middleware:

- `middleware.Logger`: logs method, path, status, duration, remote address, and user agent.
- `middleware.Recovery`: catches panics and returns a standard 500 response.
- `middleware.CORS`: sets CORS headers.

Route-specific middleware:

- `middleware.AuthMiddleware`: protects authenticated routes.

## Logging

Logging uses the standard library `log/slog`.

Logger setup lives in:

```text
internal/logger/logger.go
```

Behavior:

- `APP_ENV=production`: JSON logs at info level.
- Other environments: text logs at debug level.

## Dependency Construction

Dependencies are manually wired in `cmd/api/main.go`.

Current construction order:

```text
config
logger
database pool
repositories
services
handlers
routes
server
```

Manual wiring is enough for the current project size. If constructors become noisy, introduce a small application bootstrap package before considering a dependency injection framework.

## Adding a New Domain

Use this package shape:

```text
internal/<domain>
+-- model.go
+-- repository.go
+-- service.go
+-- handler.go
+-- routes.go
```

Recommended steps:

1. Define request, response, and domain structs in `model.go`.
2. Add database queries in `repository.go`.
3. Add validation and business rules in `service.go`.
4. Add HTTP handlers in `handler.go`.
5. Add route registration in `routes.go`.
6. Wire the repository, service, handler, and routes in `cmd/api/main.go`.
7. Add migrations for schema changes.

## Current Boundaries and Tradeoffs

The backend is not a microservice system. It is one process with clear domain packages.

Intentional choices:

- Keep auth, accounts, and transactions in one binary for now.
- Keep PostgreSQL as the single source of truth.
- Use explicit constructors instead of global state.
- Use route registration per domain to keep `main.go` small.
- Share middleware and response helpers across domains.

Known coupling:

- `account` depends on `transaction` to write deposit and withdrawal transaction records.
- `account.TransferTx` currently performs the account balance update and transfer transaction insert in the account repository.
- Auth middleware reads `JWT_SECRET` from the environment directly.

These are acceptable for the current modular monolith phase. If the project grows, the first cleanup targets should be moving token verification behind an auth service/helper and extracting transaction-writing behavior behind an interface.
