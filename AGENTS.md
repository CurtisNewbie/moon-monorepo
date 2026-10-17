# AGENTS.md

## Dev workflow

Load the `dev-flow` skill before starting any non-trivial task — a feature, a bugfix, or workflow setup — and follow its lane.

The skills it dispatches live in the offloaded store, not installed. Reach any of them through the `offloaded-skill` skill; the `skill` tool cannot.

## Repository Structure

Manual monorepo: Go microservices in `backend/<service>/`, Angular app in `frontend/moon/`, Docker Compose configs in `deploy/`.

Each backend service is standalone — its own `go.mod/go.sum` and `conf.yml`, no `go.work`.

## Runtime Dependencies (for Backend)

- MySQL >= 5.7
- RabbitMQ
- Consul
- Redis
- [event-pump](https://github.com/CurtisNewbie/event-pump) >= v0.0.18

**Start order**: middleware services first, then backend services.

## Backend Development

### Build & Test

Each service builds and tests standalone:

```bash
cd backend/<service>
go build ./...
go test ./...
```

### Inter-Service Dependencies

Some `go.mod` files use `replace` directives pointing to sibling directories. Build order matters:

- **acct**, **user-vault**, **mini-fstore**: no dependency on another service
- **gatekeeper**, **logbot**: `replace user-vault => ../user-vault`
- **vfm**: `replace user-vault => ../user-vault` & `replace mini-fstore => ../mini-fstore`

### Configuration Files

- **Development**: `conf.yml` in each backend service
- **Production**: uses `${VAR_NAME}` syntax to read from environment variables
  - Most services: `conf-prod.yml`
  - ⚠️ **user-vault** and **vfm** use `prod-conf.yml` (not `conf-prod.yml`)
- Production env vars in `deploy/backend.env`

### Database Schema Migrations

**No auto-migration.** Versioned DDL lives in `docs/schema/`: `docs/schema/current/` holds the live per-service snapshot, `docs/schema/vX.Y.Z/` holds each release's increments, and the repo-root `changes/changes.md` is the release manifest. See `docs/agents/db-schema.md` for the layout, the update rules, and how the manifest is maintained.

## Backend Framework

All backend services use [miso](https://github.com/CurtisNewbie/miso) framework.

**Mandatory:** load the `miso` skill before writing any backend Go code, and follow its conventions (error handling with `errs` — never `fmt.Errorf`/`errors.New`, bootstrap, rail logging, db, config, etc.). Install if absent:

```sh
npx skills add https://github.com/curtisnewbie/miso
```

See service-specific `doc/config.md` for custom properties.

## Frontend

### Setup & Run

```bash
cd frontend/moon/
npm ci
ng serve  # or NODE_OPTIONS=--openssl-legacy-provider ng serve if openssl version incompatible
```

Dev server proxies all requests to `http://127.0.0.1:7070` (gatekeeper).

### Key Facts

- **Angular 11** with `NgModule`
- **Angular Material** with pink-bluegrey theme
- **TSLint** (not ESLint) — run with `ng lint`
- **Tests**: Karma/Jasmine on ChromeHeadless — `npm test -- --watch=false --browsers=ChromeHeadless`. Specs are written by hand: the component schematic sets `skipTests: true`, so `ng generate` creates none (that setting affects generators only, never `ng build` or `ng test`).
- **HTTP interceptors**: token attach and global 401/403 handling — see `frontend/moon/doc/agents/http.md`
- `patch-package` runs on `postinstall` — patches exist in `patches/`

### i18n

**Preferred**: `trl` pure pipe in templates.

```html
{{'module' | trl:'key'}}
{{'module' | trl:'key':'paramName':paramValue}}
```

**NEVER inject TrlPipe** — pipes are templates-only.

For dynamic keys (TypeScript code), inject `I18n` service:
```typescript
constructor(private i18n: I18n) {}
this.snackBar.open(this.i18n.trl("module", "key"), "ok");
```

See `frontend/moon/doc/agents/i18n.md` for full guidelines.

### Pagination

Paged lists are driven by `ControlledPaginatorComponent`. See `frontend/moon/doc/agents/pagination.md` before wiring a list fetch or configuring a page limit — the limit has to be applied before the paginator's first fetch, which its `pageChanged` output triggers.

## Backend API Documentation

**Always reference backend `doc/api.md` when working on frontend-backend integration.**

Each backend service maintains auto-generated API documentation at `backend/{service}/doc/api.md` including endpoints, request/response schemas, and TypeScript interfaces.

This prevents field name mistakes (e.g., `dirFileKey` vs `fileKey`).

## Deployment

Production runs on Docker Compose: Nginx reverse proxy, Prometheus monitoring, Grafana dashboard. Config files differ from development — see [Configuration Files](#configuration-files).

## Entry Points

- **Frontend**: `frontend/moon/src/main.ts`
- **Backend services**: `main.go` in each service root (none use `cmd/` subdirectory)

## Agent skills

### Issue tracker

Issues and specs live as local markdown under `.scratch/`; `docs/agents/issue-tracker.md` covers the conventions and how skills read and write them.

### Triage labels

The five canonical triage labels (default strings). `docs/agents/triage-labels.md` maps each to this repo's label.

### Domain docs

This is a single-context repo: root `GLOSSARY.md` + `docs/adr/`. Read them before exploring an unfamiliar area; `docs/agents/domain.md` covers how the skills consume them.
