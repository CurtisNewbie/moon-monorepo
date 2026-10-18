# Releasing moon-monorepo

A release is four things kept in sync:

1. a **git tag** `vX.Y.Z` on `main`
2. the **version string in code** (6 `version.go` + 1 `version.ts`)
3. the **schema increment** under `docs/schema/vX.Y.Z/` and its entry in `release/db-changes.md`
4. **release notes** (`release/vX.Y.Z.md`), published as the GitHub release body

## Pre-release checklist

1. **Scope is closed** — everything intended for the release is merged to `main`; no open ticket under `.scratch/` targets it.
2. **Schema increment is complete** — `docs/schema/vX.Y.Z/<service>.sql` holds every pending DDL for every service, `docs/schema/current/<service>.sql` already contains the same changes, and `release/db-changes.md` has the `## Release vX.Y.Z` section. Rules: [`docs/agents/db-schema.md`](../docs/agents/db-schema.md).
3. **Migration is documented** — `release/db-changes.md` states the apply order and that the DDL runs *before* the matching code is deployed. Scripts are not idempotent.
4. **Release notes are written** — `release/vX.Y.Z.md` covers user-facing changes, migration steps, and config/deploy deltas.
5. **Config/deploy deltas are recorded** — new env vars, nginx locations, gatekeeper filters, compose changes. `deploy/docker-compose.yml` and `deploy/backend.env` are often unchanged; say so explicitly.
6. **Tree is clean and on `main`** — docs committed first.
7. **Bump the version** — `python3 release/release.py vX.Y.Z`, review the 7 changed files.
8. **Commit and tag** — commit the bump as `Release vX.Y.Z`, inspect it, then `git tag vX.Y.Z`.
9. **Push** — `git push && git push origin vX.Y.Z`.
10. **Publish** — create the GitHub release from the pushed tag, body from `release/vX.Y.Z.md`.

Build/test gates are per-service ([`backend/AGENTS.md`](../backend/AGENTS.md), [`frontend/AGENTS.md`](../frontend/AGENTS.md)). The release script does not run them; they are a separate decision, not a release step.

## Release notes

- [v0.0.5](./v0.0.5.md)

## `release/release.py`

```sh
python3 release/release.py            # print last tag + guessed next, exit 1 (dry peek)
python3 release/release.py v0.0.5     # update the version constants only; no commit, no tag
```

What it does:

- rewrites every `version.go` whose body matches `Version = "..."` to `package <pkg>` + `var ( Version = "..." )`
- rewrites every `version.ts` matching `export const Version = "..."`
- runs `go fmt ./...` in every directory containing a `go.mod`

Files it must update in this repo:

- `backend/acct/internal/server/version.go`
- `backend/gatekeeper/gatekeeper/version.go`
- `backend/logbot/logbot/version.go`
- `backend/mini-fstore/internal/server/version.go`
- `backend/user-vault/internal/server/version.go`
- `backend/vfm/internal/vfm/version.go`
- `frontend/moon/src/app/version.ts`

What it does **not** do: commit, tag, push, verify branch or clean tree, verify the schema increment or `db-changes.md`, limit the changed set, or run builds/tests.

## Tags

- `vX.Y.Z`, lightweight, on the release commit on `main`.
- Never move or reuse a tag.
- A fix on an already-released version ships as `vX.Y.Z-fix.N` from a `dev_vX.Y.Z` branch (see `v0.0.4-fix.1`); do not roll it into an unreleased version silently.
- Push the branch before the tag: `git push && git push origin vX.Y.Z`.
- Once the tag is pushed, `docs/schema/vX.Y.Z/` and its `db-changes.md` section are frozen.
- Tag the release **before** opening the next `docs/schema/vX.Y.(Z+1)/` directory.

## Deploying a release

No deploy pipeline lives in the repo: `deploy/docker-compose.yml` builds each service from `deploy/<service>/build`, populated by hand.

1. Back up the database.
2. Apply DDL for every release between the deployed version and the target, ascending:
   `source docs/schema/vX.Y.Z/<service>.sql`
   Run once — `ADD COLUMN`, `ADD KEY`, and `CREATE TABLE` are not idempotent.
3. Deploy backend services (rebuild binaries into `deploy/<service>/build`, then `docker-compose up -d --build <service>`).
4. Deploy the frontend (`NODE_OPTIONS=--openssl-legacy-provider npx ng build --prod` on Node 17+, output to the nginx html root; the OpenSSL flag is required on modern Node).
5. Apply the release's config deltas.
