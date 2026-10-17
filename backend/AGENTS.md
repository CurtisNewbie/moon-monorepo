# AGENTS.md

## Structure

Each service under `backend/<service>/` is standalone: its own `go.mod`/`go.sum` and `conf.yml`, no `go.work`. Entry point is `main.go` at the service root (none use `cmd/`).

## Build & Test

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

## Framework

All backend services use [miso](https://github.com/CurtisNewbie/miso) framework.

**Mandatory:** load the `miso` skill before writing any backend Go code, and follow its conventions (error handling with `errs` — never `fmt.Errorf`/`errors.New`, bootstrap, rail logging, db, config, etc.). Install if absent:

```sh
npx skills add https://github.com/curtisnewbie/miso
```

See service-specific `doc/config.md` for custom properties.

## Configuration

- **Development**: `conf.yml` in each backend service
- **Production**: uses `${VAR_NAME}` syntax to read from environment variables
  - Most services: `conf-prod.yml`
  - ⚠️ **user-vault** and **vfm** use `prod-conf.yml` (not `conf-prod.yml`)
- Production env vars in `deploy/backend.env`

## Runtime Dependencies

- MySQL >= 5.7
- RabbitMQ
- Consul
- Redis
- [event-pump](https://github.com/CurtisNewbie/event-pump) >= v0.0.18

**Start order**: middleware services first, then backend services.

## Database Schema Migrations

**No auto-migration.** Versioned DDL lives in `docs/schema/`: `docs/schema/current/` holds the live per-service snapshot, `docs/schema/vX.Y.Z/` holds each release's increments, and the repo-root `changes/changes.md` is the release manifest. See `docs/agents/db-schema.md` for the layout, the update rules, and how the manifest is maintained.

## Deployment

Production runs on Docker Compose (`deploy/docker-compose.yml`): Nginx reverse proxy, Prometheus monitoring, Grafana dashboard. Config files differ from development — see [Configuration](#configuration).
