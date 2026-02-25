# Repository Guidelines

## Project Structure & Module Organization
- `internal/` holds Go domain logic, grouped by layer: `handler/`, `service/`, `repository/`, `views/`, etc. UI assets (CSS/JS/templates) also live under `internal/views/...` and are embedded via `go-assets`.
- Browser-delivered static artifacts reside in `static/`. Database migrations are in `internal/migrate/scripts/` and are embedded into `assets/assets.go` (regenerated via `go-assets-builder`).
- Top-level helper scripts: `build.sh`, `start-dev.sh`, `gen-dev-certs.sh`.

## Build, Test, and Development Commands
- `go build -o app` — default binary build; `build.sh` wraps this and also refreshes embedded assets via `go-assets-builder`.
- `./start-dev.sh` — spins up the dev server with live assets.

## Coding Style & Naming Conventions
- Go code must be formatted with `gofmt` before committing. Follow existing naming: long-form nouns such as `ResourceServiceMetric`, and exported structs live in `internal/model`.
- Keep HTML/CSS/JS assets in ASCII; prefer gomponents for server-rendered views under `internal/views`. Use descriptive CSS class names scoped per page file.
- Repository layer owns SQL; services should call repository methods instead of issuing raw queries directly.

## Testing Guidelines
- run `build.sh` to check if app compiles.

## Commit & Pull Request Guidelines
- Commit messages follow the format `type(scope1,scope2): summary`, e.g., `feat(service,resource): add resource service metric edit/archiving`. Use lowercase type verbs (`feat`, `fix`, `refactor`, etc.).

## Agent-Specific Instructions
- Never run destructive git commands (`reset --hard`, `checkout -- <file>`) without explicit approval.

