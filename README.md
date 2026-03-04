# go-rest-api-starter

A production-ready Go REST API template. Stack: **Gin** + **PostgreSQL (pgx/v5)** + **JWT Auth** + **Docker**.

## Project Structure

```
.
├── cmd/api/main.go          # Entry point
├── internal/
│   ├── config/              # Environment config loader
│   ├── handler/             # HTTP handlers (auth, user)
│   ├── middleware/          # JWT auth middleware
│   ├── model/               # Request/response types
│   ├── repository/          # pgx DB queries
│   ├── service/             # Business logic
│   └── server/              # Gin router setup
├── pkg/
│   ├── database/            # pgx pool connection
│   ├── jwt/                 # Token generation & validation
│   └── response/            # Standardized JSON responses
├── migrations/              # SQL migration files
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── .env.example
```

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/v1/auth/register` | No | Register a new user |
| POST | `/api/v1/auth/login` | No | Login, receive JWT |
| GET | `/api/v1/users/me` | Yes | Get current user profile |
| PUT | `/api/v1/users/me` | Yes | Update current user profile |

Protected routes require `Authorization: Bearer <token>` header.

---

## Using This Repo as a Template for New Projects

This repo is designed to be pulled as a starting point for new Go REST API projects. Follow these steps:

### 1. Create your new project repository

Create a new empty repo on GitHub (or any git host), then locally:

```bash
mkdir my-new-project
cd my-new-project
git init
```

### 2. Add this repo as a remote and pull the template

```bash
git remote add template https://github.com/Ramazon1227/go-rest-api-starter.git
git fetch template
git merge template/main --allow-unrelated-histories
```

### 3. Connect your new project to its own remote origin

```bash
git remote add origin https://github.com/your-username/my-new-project.git
git push -u origin main
```

> You can keep `template` remote to pull future improvements:
> ```bash
> git fetch template
> git merge template/main
> ```

### 4. Rename the Go module

Replace the module path in `go.mod` and all import paths:

```bash
# macOS / Linux
find . -type f -name "*.go" | xargs sed -i '' 's|github.com/Ramazon1227/go-rest-api-starter|github.com/your-username/my-new-project|g'
sed -i '' 's|github.com/Ramazon1227/go-rest-api-starter|github.com/your-username/my-new-project|g' go.mod
```

Then tidy dependencies:

```bash
go mod tidy
```

### 5. Configure environment

```bash
cp .env.example .env
# Edit .env with your DB credentials, JWT secret, etc.
```

### 6. Run with Docker (recommended for first run)

```bash
make docker-up
```

This starts the API and a PostgreSQL instance. The migration file in `migrations/` is automatically applied on first startup.

### 7. Or run locally

Ensure PostgreSQL is running, then apply the migration and start the server:

```bash
make migrate
make run
```

---

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make run` | Run the API locally |
| `make build` | Build binary to `bin/api` |
| `make test` | Run all tests |
| `make tidy` | Run `go mod tidy` |
| `make docker-up` | Start API + PostgreSQL with Docker |
| `make docker-down` | Stop Docker containers |
| `make docker-logs` | Tail API container logs |
| `make migrate` | Apply SQL migrations manually |

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | HTTP server port |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | DB username |
| `DB_PASSWORD` | _(empty)_ | DB password |
| `DB_NAME` | `myapp` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode |
| `JWT_SECRET` | `changeme` | JWT signing secret — **change in production** |
| `JWT_EXPIRY_HOURS` | `24` | Token expiry in hours |

---

## Extending the Template

- **Add a new resource**: create `model/`, `repository/`, `service/`, and `handler/` files following the existing `user` pattern, then register routes in `internal/server/server.go`.
- **Add a migration**: create `migrations/002_your_change.sql` and run `make migrate`.
- **Add middleware**: create a handler in `internal/middleware/` and apply it to a route group in `server.go`.

## License

MIT
