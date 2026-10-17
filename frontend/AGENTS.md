# AGENTS.md

The Angular app lives in `frontend/moon/`; entry point `frontend/moon/src/main.ts`.

## Setup & Run

```bash
cd frontend/moon/
npm ci
ng serve  # or NODE_OPTIONS=--openssl-legacy-provider ng serve if openssl version incompatible
```

Dev server proxies all requests to `http://127.0.0.1:7070` (gatekeeper).

## Key Facts

- **Angular 11** with `NgModule`
- **Angular Material** with pink-bluegrey theme
- **TSLint** (not ESLint) — run with `ng lint`
- **Tests**: Karma/Jasmine on ChromeHeadless — `npm test -- --watch=false --browsers=ChromeHeadless`. Specs are written by hand: the component schematic sets `skipTests: true`, so `ng generate` creates none (that setting affects generators only, never `ng build` or `ng test`).
- **HTTP interceptors**: token attach and global 401/403 handling — see `frontend/moon/doc/agents/http.md`
- `patch-package` runs on `postinstall` — patches exist in `patches/`

## i18n

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

## Pagination

Paged lists are driven by `ControlledPaginatorComponent`. See `frontend/moon/doc/agents/pagination.md` before wiring a list fetch or configuring a page limit — the limit has to be applied before the paginator's first fetch, which its `pageChanged` output triggers.
