# go-rest-api-starter

A production-ready Go REST API template.

**Stack:** Gin · PostgreSQL (pgx/v5) · MongoDB · JWT Auth · Swagger · Docker

---

## Project Structure

```
.
├── cmd/main.go               # Entry point, wires config → packages
├── api/
│   ├── api.go                # Gin router setup + Swagger annotations
│   ├── docs/                 # Generated Swagger docs (swag init)
│   ├── handlers/             # HTTP handlers (auth, profile, user)
│   ├── http/                 # Response types and status codes
│   └── middleware/           # JWT auth + role middleware
├── config/
│   └── config.go             # Environment config loader
├── migrations/
│   └── postgres/             # SQL migration files
├── models/                   # Shared request/response/DB types
├── pkg/
│   ├── email/                # SMTP welcome email
│   ├── jwt/                  # Token generation, validation, blacklist
│   ├── logger/               # Structured logger (zap)
│   └── utils/                # Password hashing, random generation
├── storage/
│   ├── storage.go            # StorageI and UserRepoImpl interfaces
│   ├── mongo/                # MongoDB implementation
│   └── postgres/             # PostgreSQL implementation
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

---

## API Endpoints

| Method | Path | Auth | Role | Description |
|--------|------|------|------|-------------|
| POST | `/api/v1/auth/login` | No | — | Login, receive JWT |
| POST | `/api/v1/auth/logout` | No | — | Invalidate token |
| GET | `/api/v1/profile` | Yes | Any | Get own profile |
| PUT | `/api/v1/profile` | Yes | Any | Update own profile |
| PUT | `/api/v1/profile/password` | Yes | Any | Change password |
| POST | `/api/v1/user` | Yes | SYSTEM_ADMIN | Create user |
| GET | `/api/v1/user` | Yes | Any | List users |
| GET | `/api/v1/user/:user_id` | Yes | Any | Get user by ID |
| PUT | `/api/v1/user/:user_id` | Yes | SYSTEM_ADMIN | Update user |
| DELETE | `/api/v1/user/:user_id` | Yes | SYSTEM_ADMIN | Soft-delete user |
| GET | `/swagger/*any` | No | — | Swagger UI |

Protected routes require `Authorization: Bearer <token>` header.

---

## Quick Start

### 1. Configure environment

```bash
cp .env.example .env
# Fill in SECRET_KEY, POSTGRES_PASSWORD, and SMTP_* values
```

### 2. Run with Docker

```bash
make docker-up
```

Starts the API and a PostgreSQL instance. Open [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html).

### 3. Or run locally

```bash
# Apply migrations (requires golang-migrate)
make migration-up DB_URL="postgres://postgres:pass@localhost:5432/go_rest_api_starter?sslmode=disable"

make run
```

---

## Storage Backends

The app implements a common `StorageI` interface. Switch backends by changing one line in `main.go`:

```go
// PostgreSQL (default)
pgStore, err := postgres.NewPostgres(ctx, cfg)

// MongoDB
mongoStore, err := mongo.NewMongo(ctx, cfg)
```

Both backends implement identical soft-delete behaviour (`deleted_at` field).

---

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make run` | Run the API locally |
| `make build` | Build binary to `bin/` |
| `make test` | Run all tests |
| `make tidy` | Run `go mod tidy` |
| `make lint` | Run golangci-lint |
| `make swag-init` | Regenerate Swagger docs |
| `make docker-up` | Start API + PostgreSQL with Docker |
| `make docker-down` | Stop Docker containers |
| `make docker-logs` | Tail API container logs |
| `make migration-up` | Apply SQL migrations |
| `make migration-down` | Roll back SQL migrations |
| `make rename-module NEW_MODULE=<path>` | Rename the Go module across all files |

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ENVIRONMENT` | `debug` | `debug`, `test`, or `release` |
| `SERVICE_HOST` | `localhost` | Host shown in Swagger UI |
| `HTTP_PORT` | `:8080` | Listen address |
| `SECRET_KEY` | _(required)_ | JWT signing secret |
| `POSTGRES_HOST` | `0.0.0.0` | PostgreSQL host |
| `POSTGRES_PORT` | `5432` | PostgreSQL port |
| `POSTGRES_USER` | `postgres` | PostgreSQL user |
| `POSTGRES_PASSWORD` | _(required)_ | PostgreSQL password |
| `POSTGRES_DATABASE` | _(service name)_ | PostgreSQL database |
| `MONGO_URI` | `mongodb://localhost:27017` | MongoDB connection URI |
| `MONGO_DATABASE` | _(service name)_ | MongoDB database name |
| `SMTP_HOST` | `smtp.gmail.com` | SMTP server host |
| `SMTP_PORT` | `587` | SMTP server port |
| `SMTP_USERNAME` | _(empty)_ | SMTP username |
| `SMTP_PASSWORD` | _(empty)_ | SMTP password |
| `SMTP_FROM` | _(empty)_ | Sender email address |

---

## Using as a Template

### 1. Create your project

```bash
mkdir my-project && cd my-project
git init
git remote add template https://github.com/Ramazon1227/go-rest-api-starter.git
git fetch template
git merge template/main --allow-unrelated-histories
```

### 2. Rename the module

```bash
make rename-module NEW_MODULE=github.com/your-username/my-project
```

This rewrites every import path in all `.go` files and `go.mod`, then runs `go mod tidy`.

### 3. Add a new resource

1. Add types to `models/`
2. Add methods to `storage/storage.go` interface
3. Implement in `storage/postgres/` and `storage/mongo/`
4. Add handlers in `api/handlers/`
5. Register routes in `api/api.go`
6. Run `make swag-init` to regenerate Swagger docs

---

## License

MIT
