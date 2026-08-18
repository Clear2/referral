# Repository Guidelines

## Project Structure & Module Organization

This workspace contains a backend and a pnpm frontend workspace with two independently deployed applications:

- `backend/`: Go modular monolith using Gin, Fx, Ent, PostgreSQL, Redis, gRPC, and RabbitMQ. Business code lives in `internal/modules/<domain>`. Infrastructure composition is in `internal/app`; transports are under `internal/router` and `internal/rpc`. Ent schemas are maintained in `ent/schema`; generated Ent code must not be edited manually. SQL migrations live in `ent/migrate/migrations`.
- `frontend/apps/web/`: React 19 and React Router user application. It owns `/`, user login, referral dashboards, and public invitation registration at `/ref/:code`.
- `frontend/apps/admin/`: React 19 and React Router operations console deployed under `/admin/`. It owns user, referral, Credit, RBAC, and audit administration.
- `frontend/packages/api/`: shared API contracts and browser client behavior. `frontend/packages/ui/` contains presentation primitives shared by both applications.

Keep request types inside their owning module. Reserve `backend/internal/dto` for genuinely shared transport contracts.

## Build, Test, and Development Commands

Run commands from the relevant application directory.

```bash
cd backend
make compose-up          # start PostgreSQL and RabbitMQ
go run ./cmd/migrate up  # apply database migrations
go run ./cmd/app         # run HTTP :8999 and gRPC :8998
make test                # race-enabled internal tests with coverage
make format              # gofumpt and import formatting
make linter-golangci     # Go lint checks
make ent-gen             # regenerate Ent after schema changes
make swag-v1             # regenerate Swagger documentation

cd frontend
pnpm dev:web          # user application on :5173
pnpm dev:admin        # admin console on :5174
pnpm typecheck        # check both applications
pnpm build            # build both applications
```

## Coding Style & Naming Conventions

Format Go with `gofumpt`; package names are lowercase and domain-oriented. Use `RegisterRoutes` for module route registration and constructors named `New<Type>`. Keep Controller → Service → Repository dependencies one-way. TypeScript and TSX use Prettier, two-space indentation, and PascalCase React components.

## Testing Guidelines

Place Go tests beside source files as `*_test.go`; prefer table-driven tests and mocks only at real seams. Run targeted package tests while developing, then `make test`. Frontend changes must pass workspace typecheck/build; manually verify affected login, referral registration, Credit, user management, and analytics flows. Routing changes must preserve `/` for users, `/admin/` for operations, and `/api/` for the backend.

## AI Collaboration

Project-local skills live under `.agents/skills` and are shared with Claude-compatible tools through `.claude/skills`. Use `referral-feature` for end-to-end changes, `referral-investigate` for defects, `referral-code-review` for read-only diff and pull request review, and `referral-release-check` before handoff or release. Record durable decisions using `.agents/notes/README.md`; do not create notes for trivial local edits.

## Commit & Pull Request Guidelines

History is minimal; use concise Conventional Commit subjects such as `feat: add activity analytics` or `fix: rotate refresh token`. Pull requests should describe behavior changes, migrations/configuration, verification commands, linked issues, and screenshots for UI changes.

## Security & Configuration

Never commit `backend/config.yaml`, passwords, JWTs, or raw refresh tokens. Update `config.example.yaml` for new settings. Production cookies require HTTPS and `cookie_secure: true`.
