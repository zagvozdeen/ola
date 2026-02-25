# Repository Guidelines

## Project Structure & Module Organization
- `cmd/main.go`: application entry point.
- `internal/`: backend code by layer:
  - `internal/api`: HTTP handlers, routing, validation, permissions.
  - `internal/store`: data-access layer and domain models.
  - `internal/db`: PostgreSQL connection + embedded Goose migrations (`internal/db/migrations`).
  - `internal/config`, `internal/logger`, `internal/event_bus`, `internal/worker_pool`, `internal/seeder`: supporting infrastructure.
- `web/`: frontend sources:
  - `web/spa`: Vue 3 + TypeScript admin/client SPA.
  - `web/landing`: landing page entry.
- `public/`: static assets. `templates/`: server templates. `dist/`: build output.
- `deploy/`: deployment configs (`nginx`, `systemd`).

## Build, Test, and Development Commands
- `docker compose up -d`: start local PostgreSQL (`localhost:5432`, db/user/pass: `ola`).
- `make dev`: run backend with `config.local.yaml` (`GOEXPERIMENT=jsonv2` enabled).
- `npm run dev`: run Vite dev server for frontend work.
- `make build`: build Linux backend binary (`./ola`) and frontend assets (`dist/`).
- `make test`: run backend tests/check compilation.
- `npx eslint web/spa/src --ext .ts,.vue`: run frontend lint checks.
- `npm run build`: Vue type-check + production frontend build.

## Coding Style & Naming Conventions
- Go: keep code `gofmt`-formatted; package names lowercase; file names use snake_case where practical (for example, `worker_pool.go`).
- Vue/TS: ESLint-enforced style from `eslint.config.ts`:
  - 2-space indentation, single quotes, no semicolons, trailing commas on multiline objects/arrays.
  - Prefer `import type` for types.
  - No `console` except `console.warn`/`console.error`.
- Naming patterns:
  - Vue components in PascalCase (`PageOrders.vue`, `HeaderMenu.vue`).
  - Composables prefixed with `use` (`useAuthState.ts`).

## Testing Guidelines
- There is currently no comprehensive committed test suite; add tests with new features/fixes.
- Backend tests should live beside code as `*_test.go`.
- Frontend changes should at minimum pass `npm run build` and ESLint; include manual verification notes for changed flows.

## Commit & Pull Request Guidelines
- Follow existing Conventional Commit style seen in history: `feat: ...`, `fix: ...`, `refactor(scope): ...`, `chore: ...`.
- Keep commits focused and atomic.
- PRs should include:
  - clear summary and motivation,
  - linked issue/task (if available),
  - commands run (build/lint/test),
  - screenshots or short recordings for UI changes.

## Security & Configuration Tips
- Never commit real secrets or production credentials.
- Use local config files for development (`config.local.yaml`) and review `deploy/` changes carefully before release.
