# Release SOP

How to cut a moon-monorepo release. A release is four artifacts kept in sync:

| Artifact | Location |
| --- | --- |
| version strings | 6 `version.go` + `frontend/moon/src/app/version.ts` |
| schema increment | `docs/schema/vX.Y.Z/<service>.sql` + `docs/schema/current/<service>.sql` |
| schema manifest | `release/db-changes.md` |
| release notes | `release/vX.Y.Z.md` → GitHub release body |

The tag is the last step. Do not tag until the version commit has been reviewed.

## 0. Prerequisites

- On `main`, clean tree, in sync with `origin/main`:
  ```sh
  git status --short          # empty
  git fetch origin && git status -sb
  ```
- Target tag does not exist yet:
  ```sh
  git tag -l vX.Y.Z           # empty
  ```
- `go`, `node`/`npm`, `python3` available.
- Database access for the migration; GitHub access for publishing.
- Agreed version number (`vX.Y.Z`) and closed scope.

## 1. Prepare the release

### 1.1 Close scope

```sh
git log --oneline "$(git describe --tags --abbrev=0)"..HEAD
```

Everything intended for the release is merged; no open ticket under `.scratch/` targets it.

### 1.2 Verify the schema increment

```sh
ls docs/schema/vX.Y.Z/
grep -n "Release vX.Y.Z" release/db-changes.md
git diff --stat "$(git describe --tags --abbrev=0)"..HEAD -- docs/schema
```

Must be true:

- every pending DDL for every service is in `docs/schema/vX.Y.Z/`
- `docs/schema/current/` already contains the same changes
- `release/db-changes.md` has the release section, the SQL links, and the apply order ("before deploying the code", non-idempotent)

Rules: [`docs/agents/db-schema.md`](../docs/agents/db-schema.md). No pending DDL for a service means no file for it.

### 1.3 Write release notes

Create `release/vX.Y.Z.md` from the commit range above. It must cover:

- user-facing features and fixes
- required migration, with the exact `source` command and the deploy order
- config/deploy deltas (env vars, nginx, gatekeeper, compose), explicitly stating when there are none

Add it to the list in [Release notes](#release-notes).

### 1.4 Commit the preparation

```sh
git add release/ README.md backend/AGENTS.md
git commit -m "docs(release): prepare vX.Y.Z"
git status --short          # empty
```

The next commit must contain version changes only.

### 1.5 Build check (recommended, not part of the script)

```sh
for d in acct user-vault mini-fstore gatekeeper logbot vfm; do
  ( cd "backend/$d" && go build ./... ) || echo "FAILED: $d"
done

cd frontend/moon
NODE_OPTIONS=--openssl-legacy-provider npx ng build --prod --output-path "$(mktemp -d)"
```

Build order is documented in [`backend/AGENTS.md`](../backend/AGENTS.md); frontend details in [`frontend/AGENTS.md`](../frontend/AGENTS.md).

## 2. Bump the version

```sh
python3 release/release.py              # peek: prints last tag + guessed next (exit 1 is expected)
python3 release/release.py vX.Y.Z       # rewrites the version constants only
```

Verify the result is version-only:

```sh
git status --short                      # exactly the 7 files below
git diff                                # exactly 7 changed lines
```

Files:

- `backend/acct/internal/server/version.go`
- `backend/gatekeeper/gatekeeper/version.go`
- `backend/logbot/logbot/version.go`
- `backend/mini-fstore/internal/server/version.go`
- `backend/user-vault/internal/server/version.go`
- `backend/vfm/internal/vfm/version.go`
- `frontend/moon/src/app/version.ts`

`frontend/moon/package.json` stays at `0.0.0` by convention — do not bump it.

The script does not commit, tag, push, or verify anything; it also runs `go fmt ./...` per module.

## 3. Review, commit, tag

Review the diff, then:

```sh
git add -u
git commit -m "Release vX.Y.Z"

git show --stat HEAD                    # exactly the 7 version files
git status --short                      # empty
```

Tag the reviewed commit:

```sh
git tag vX.Y.Z                          # lightweight
git tag --points-at HEAD                # prints vX.Y.Z
git describe --tags                     # prints vX.Y.Z
```

Never move or reuse a tag. If the tag is wrong and has not been pushed:

```sh
git tag -d vX.Y.Z
```

## 4. Push

```sh
git push
git push origin vX.Y.Z
git ls-remote --tags origin vX.Y.Z      # confirm the remote tag
```

No force-push. Once pushed, the tag, `docs/schema/vX.Y.Z/`, and its `db-changes.md` section are frozen.

## 5. Publish the GitHub release

`gh` is not installed here; use the web UI:

1. GitHub → **Releases** → **Draft a new release**
2. Tag: `vX.Y.Z`; title: `vX.Y.Z`
3. Body: contents of `release/vX.Y.Z.md`, starting from `## Highlights`
4. Publish

## 6. Deploy

No deploy pipeline lives in the repo: `deploy/docker-compose.yml` builds each service from `deploy/<service>/build`, populated by hand.

1. Back up the database.
2. Apply DDL for every release between the deployed version and the target, ascending:
   ```sh
   source docs/schema/vX.Y.Z/<service>.sql
   ```
   Run each version once — `ADD COLUMN`, `ADD KEY`, and `CREATE TABLE` are not idempotent.
3. Deploy backend services (rebuild binaries into `deploy/<service>/build`, then `docker-compose up -d --build <service>`). Deploy the service only after its DDL.
4. Deploy the frontend:
   ```sh
   cd frontend/moon
   NODE_OPTIONS=--openssl-legacy-provider npx ng build --prod
   # upload dist to the nginx html root
   ```
5. Apply the release's config deltas (gatekeeper filters, nginx locations, env vars).

## 7. After the release

- Open the next `docs/schema/vX.Y.(Z+1)/` directory only after tagging this one.
- A fix on an already-released version ships as `vX.Y.Z-fix.N` from a `dev_vX.Y.Z` branch (see `v0.0.4-fix.1`); do not fold it into an unreleased version silently.

## Release notes

- [v0.0.5](./v0.0.5.md)
