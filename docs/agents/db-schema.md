# Database Schema Migrations

**No auto-migration.** DDL changes are versioned SQL files in the root `schema/` directory, applied by hand.

```
schema/
├── current/          ← live snapshot, one file per service (hand-edited)
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

## Rules

- `gatekeeper` has no schema (stateless gateway only)
- `current/` holds the full current `CREATE TABLE`/`CREATE DATABASE` DDL for each service — the schema as it is now
- Each version directory contains all schema changes for that release across all services
- All `CREATE TABLE` statements use `database.table_name` format (no backtick quoting on table names)
- A schema change updates both, in the same change:
  - the service's snapshot: `schema/current/<service>.sql`
  - the version directory increment: `schema/v<release>/<service>.sql`
- Version directories are append-only history — never edited once the release is tagged. Never replay a version directory on top of `current/`: `current/` already contains those changes.
- Only one version directory may be unreleased at a time. While `schema/vX.Y.Z/` is untagged it is the in-flight release, and every pending schema change for every service goes into it. Tag the release before opening the next version directory.
- Version directory name is determined from git tags — use the next semantic version (`git tag --sort=-version:refname` for the latest)
- A version directory only needs files for services that changed: `schema/v0.0.6/acct.sql`, `schema/v0.0.6/vfm.sql`, etc.
- Execute DDL scripts manually per version: `source schema/v0.0.6/acct.sql`

## `changes/changes.md`

The release manifest for schema changes, and the only place release-level schema history is recorded.

- One `## Release vX.Y.Z` section per version directory, in ascending version order — the file reads oldest to newest, top to bottom.
- Each section lists the SQL files changed in that release, as a numbered list of markdown links relative to the `changes/` directory:
  `1. [schema/v0.0.5/vfm.sql](../schema/v0.0.5/vfm.sql)`
- A new section is added in the same change that opens its version directory, and its entry list grows as further services change within that release.
- A version's section is final once the release is tagged; never reorder, remove, or rewrite a released section. Correcting the ascending order of released sections is the only permitted edit.
- Only schema changes are logged here — not application changes.
