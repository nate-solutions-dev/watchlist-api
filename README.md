# Watchlist API

Movie and show personal watchlist tracker. Built in Go with Gin, PostgreSQL (Supabase), and JWT auth. Integrates TMDB for media metadata and Anthropic Claude for AI-generated summaries.

## Status

- ✅ Auth (register, login, JWT-protected /me)
- ✅ Watchlist CRUD
- ✅ Health check endpoint
- ⬜ TMDB integration
- ⬜ Claude AI summary endpoint
- ⬜ Frontend (React + Vite)
- ⬜ Deployed to Railway
- ⬜ Swagger docs

## Tech Stack

- Go 1.22
- Gin
- pgx/v5
- PostgreSQL (Supabase)
- JWT
- bcrypt
- godotenv

## Local Setup

1. Clone the repo.
2. Copy `.env.example` to `.env` and fill in values (never commit `.env`).
3. Run `go mod download`.
4. Run `go run ./cmd/main.go`.

The server listens on port **8080** by default. Set the `PORT` environment variable to use another port.

## API Routes

**Public**

- `GET /health` — Liveness check (no database).
- `POST /public/auth/register` — Register a new account.
- `POST /public/auth/login` — Log in and receive a JWT.

**Protected** (requires `Authorization: Bearer <token>`)

- `GET /private/auth/me` — Current user from JWT.
- `GET /private/watchlist` — List watchlist items for the authenticated user.
- `POST /private/watchlist` — Add a watchlist item.
- `PATCH /private/watchlist/:id` — Update a watchlist item.
- `DELETE /private/watchlist/:id` — Delete a watchlist item.

## Architecture

The API follows clean architecture: HTTP routing hands off to controllers, which call services for business logic and repositories for data access. Shared shapes use DTOs; persisted entities use models. Database access uses pgx/v5 directly (no ORM). Mutating watchlist operations enforce ownership by scoping updates and deletes with `USER_DATA_ID` taken from the JWT context, not from the request body.

## Roadmap

Ideas for v2 include grouped watchlists (`WATCH_LIST_GROUP`), title reviews (`TITLE_REVIEW`), refresh tokens, and password reset flows.

## License

MIT
