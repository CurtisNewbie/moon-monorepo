# AGENTS.md

## Repository Structure

Manual monorepo with:
- **Backend**: Go microservices (gatekeeper, acct, vfm, user-vault, mini-fstore, logbot)
- **Frontend**: Angular application (`frontend/moon/`)
- **Deployment**: Docker Compose configs in `deploy/`

Each backend service is independent with its own `go.mod/go.sum` and `conf.yml`.

## External Dependencies (Required for Backend Development)

- MySQL >= 5.7
- RabbitMQ
- Consul
- Redis
- [event-pump](https://github.com/CurtisNewbie/event-pump) >= v0.0.18

**Order matters**: Start middleware services before backend services.

### Frontend

```bash
cd frontend/moon/
npm ci
ng serve  # or NODE_OPTIONS=--openssl-legacy-provider ng serve if openssl version incompatible
```

Frontend dev server proxies to `http://127.0.0.1:7070` (gatekeeper) via `src/proxy.conf.json`.

## Database Schema Migrations

**No auto-migration**. DDL changes are maintained as versioned SQL files in each service's `schema/` folder.

- Execute DDL scripts manually based on your version
- Check `changes/changes.md` for migration history
- Example: `backend/user-vault/schema/v0.0.4.sql`

## Configuration Files

- **Development**: `conf.yml` in each backend service
- **Production**: `conf-prod.yml` uses environment variables with `${VAR_NAME}` syntax
- Production env vars are in `deploy/backend.env`

## Backend Framework

All backend services use [miso](https://github.com/CurtisNewbie/miso) framework:

Use miso skill, install if absent:

```sh
npx sklls add https://github.com/curtisnewbie/miso
```

See service-specific `doc/config.md` for custom properties.

## Frontend i18n

**Preferred approach**: Use `trl` pure pipe for static translation keys in templates.

```html
{{'module' | trl:'key'}}
{{'module' | trl:'key':'paramName':paramValue}}
```

**NEVER inject TrlPipe** - pipes are for templates only.

For dynamic keys (constructed at runtime), inject `I18n` service:
```typescript
constructor(private i18n: I18n) {}
this.snackBar.open(this.i18n.trl("module", "key"), "ok");
```

See `frontend/moon/doc/agents/i18n.md` for full guidelines.

## Deployment

Production uses Docker Compose:
- All configs externalized to environment variables
- Use `conf-prod.yml` (not `conf.yml`) for production
- Includes Nginx reverse proxy, Prometheus monitoring, Grafana dashboard

Config files in `deploy/`: `docker-compose.yml`, `backend.env`, `nginx.conf`, `prometheus.yml`, `grafana_dashboard.json`

## Entry Points

- **Frontend**: `frontend/moon/src/main.ts`
- **Backend services**: `main.go` in each service's root directory