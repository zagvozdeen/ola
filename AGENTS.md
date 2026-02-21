# Repository Guidelines

## Project Structure & Module Organization
- Backend entrypoint: `cmd/main.go`.
- Core Go modules live in `internal/`:
  - `internal/api` for HTTP handlers and routing.
  - `internal/store` for data access/services.
  - `internal/db` for PostgreSQL connection and Goose migrations (`internal/db/migrations`).
  - `internal/config` and `internal/logger` for runtime setup.
- Frontend is a Vite multi-page app under `web/`:
  - `web/landing` (public landing page),
  - `web/spa/admin` (admin SPA),
  - `web/spa/tma` (Telegram mini app),
  - `web/shared` (shared TS composables/utilities),
  - `web/public` (static assets).
- Server-rendered HTML template: `templates/index.html`.

## Build, Test, and Development Commands
- `cp .env.example .env` to bootstrap local env vars.
- `docker compose up -d` to start PostgreSQL.
- `make dev` to run the Go API locally (`GOEXPERIMENT=jsonv2 go run cmd/main.go`).
- `GOEXPERIMENT=jsonv2 go test ./...` to run Go tests/check package compilation.
- `npm run dev` to start Vite dev server for frontend work.
- `npm run build` to run type-check (`vue-tsc -b`) and produce `dist/`.
- `npx eslint "web/**/*.{ts,vue}"` to lint TS/Vue sources.

## Coding Style & Naming Conventions
- Go: run `gofmt`; keep package names lowercase; exported symbols use `PascalCase`.
- TypeScript/Vue lint rules are enforced via `eslint.config.ts`:
  - 2-space indentation, single quotes, no semicolons, trailing commas (multiline), strict `===`.
  - Prefer `type` imports (`@typescript-eslint/consistent-type-imports`).
  - Vue component tags should use `PascalCase`.
- Keep shared frontend code in `web/shared/src` and import via `@shared/*`.

## Testing Guidelines
- There are currently no committed `_test.go` files; add tests alongside new Go packages when changing logic.
- Name Go tests `*_test.go` with `TestXxx` functions.
- For frontend changes, at minimum run lint and build before PR.

## Commit & Pull Request Guidelines
- Use Conventional Commit style seen in history: `feat: ...`, `refactor(scope): ...`, `chore: ...`.
- Keep commits focused (one concern per commit).
- PRs should include:
  - clear summary and motivation,
  - linked issue/task (if available),
  - validation steps (commands run),
  - screenshots for UI changes (`web/landing`, `web/spa/*`),
  - migration/rollback notes for DB schema updates.

## Security & Configuration Tips
- Do not commit secrets from `.env`; keep `.env.example` as the source of required keys.
- Review `APP_SECRET`, DB credentials, and Telegram tokens before deploying.
