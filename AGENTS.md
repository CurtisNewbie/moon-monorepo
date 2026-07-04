# AGENTS.md

## Repository Structure

Manual monorepo with:
- **Backend**: Go microservices (gatekeeper, acct, vfm, user-vault, mini-fstore, logbot)
- **Frontend**: Angular application (`frontend/moon/`)
- **Deployment**: Docker Compose configs in `deploy/`

Each backend service is independent with its own `go.mod/go.sum` and `conf.yml`. No `go.work` — each service is built standalone.

## Runtime Dependencies (for Backend)

- MySQL >= 5.7
- RabbitMQ
- Consul
- Redis
- [event-pump](https://github.com/CurtisNewbie/event-pump) >= v0.0.18

**Start order**: middleware services first, then backend services.

## Backend Development

### Build & Test

No Makefiles. Run per service:

```bash
cd backend/<service>
go build ./...
go test ./...
```

### Inter-Service Dependencies

Some `go.mod` files use `replace` directives pointing to sibling directories. Build order matters:

- **acct**: no deps on other services
- **user-vault**: no deps on other services (but others depend on it)
- **mini-fstore**: no deps on other services
- **gatekeeper**: `replace user-vault => ../user-vault`
- **logbot**: `replace user-vault => ../user-vault`
- **vfm**: `replace user-vault => ../user-vault` & `replace mini-fstore => ../mini-fstore`

Go version: 1.24.4 across all services.

### Configuration Files

- **Development**: `conf.yml` in each backend service
- **Production**: uses `${VAR_NAME}` syntax to read from environment variables
  - Most services: `conf-prod.yml`
  - ⚠️ **user-vault** and **vfm** use `prod-conf.yml` (not `conf-prod.yml`)
- Production env vars in `deploy/backend.env`

### Database Schema Migrations

**No auto-migration**. DDL changes are versioned SQL files in the root `schema/` directory, organized by version:

```
schema/
├── v0.0.0/           ← baseline (initial schemas)
│   ├── acct.sql
│   ├── logbot.sql
│   ├── user-vault.sql
│   ├── vfm.sql
│   └── mini-fstore.sql
├── v0.0.3/
│   └── user-vault.sql
├── v0.0.4/
│   ├── user-vault.sql
│   └── vfm.sql
└── v0.0.5/
    └── vfm.sql
```

- `gatekeeper` has no schema (stateless gateway only)
- Each version directory contains all schema changes for that release across all services
- All `CREATE TABLE` statements use `database.table_name` format (no backtick quoting on table names)
- `v0.0.0/` is the baseline — contains only initial `CREATE TABLE`/`CREATE DATABASE` DDL per service
- For new migrations: create a new version directory (e.g., `schema/v0.0.6/`), add per-service SQL files
  - Each service gets its own file: `schema/v0.0.6/acct.sql`, `schema/v0.0.6/vfm.sql`, etc.
  - A version directory only needs files for services that changed
  - Version directory name is determined from git tags — use the next semantic version (check `git tag --sort=-version:refname` for latest)
- Execute DDL scripts manually per your version (`source schema/v0.0.6/acct.sql`, etc.)
- Track migration history via `changes/changes.md`
  - When adding new DDL changes, append a new `## Release vX.Y.Z` section listing the changed SQL files

## Backend Framework

All backend services use [miso](https://github.com/CurtisNewbie/miso) framework.

Use miso skill, install if absent:

```sh
npx sklls add https://github.com/curtisnewbie/miso
```

See service-specific `doc/config.md` for custom properties.

## Frontend

### Setup & Run

```bash
cd frontend/moon/
npm ci
ng serve  # or NODE_OPTIONS=--openssl-legacy-provider ng serve if openssl version incompatible
```

Dev server proxies all requests to `http://127.0.0.1:7070` (gatekeeper) via `src/proxy.conf.json`.

### Key Facts

- **Angular 11** with `NgModule` (not standalone components)
- **Angular Material** with pink-bluegrey theme
- **TSLint** (not ESLint) — run with `ng lint`
- **Testing**: Karma + Jasmine (`ng test`), Protractor E2E (`ng e2e`)
- `angular.json` has `skipTests: true` — `ng generate component` creates no `.spec.ts` by default
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

## Backend API Documentation

**Always reference backend `doc/api.md` when working on frontend-backend integration.**

Each backend service maintains auto-generated API documentation at `backend/{service}/doc/api.md` including endpoints, request/response schemas, and TypeScript interfaces.

```bash
cat backend/vfm/doc/api.md | grep -A 20 "dir-thumbnail"
```

This prevents field name mistakes (e.g., `dirFileKey` vs `fileKey`).

## Deployment

Production uses Docker Compose:
- Configs externalized to environment variables
- Use production config files (not `conf.yml`)
- Includes Nginx reverse proxy, Prometheus monitoring, Grafana dashboard

Config files: `deploy/docker-compose.yml`, `deploy/backend.env`, `deploy/nginx.conf`, `deploy/prometheus.yml`, `deploy/grafana_dashboard.json`

## Entry Points

- **Frontend**: `frontend/moon/src/main.ts`
- **Backend services**: `main.go` in each service root (none use `cmd/` subdirectory)
