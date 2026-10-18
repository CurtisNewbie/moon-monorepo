# Database Change Log

This change log only contains major changes to database schema.

Applying a release: execute the SQL files for every release between your current version and the target, in ascending release order, **before deploying that release's code**. The scripts are not idempotent — run each version once.

## Release v0.0.3

1. [docs/schema/v0.0.3/user-vault.sql](../docs/schema/v0.0.3/user-vault.sql)

## Release v0.0.4

1. [docs/schema/v0.0.4/user-vault.sql](../docs/schema/v0.0.4/user-vault.sql)
2. [docs/schema/v0.0.4/vfm.sql](../docs/schema/v0.0.4/vfm.sql)

## Release v0.0.5

Apply before deploying `vfm`.

1. [docs/schema/v0.0.5/vfm.sql](../docs/schema/v0.0.5/vfm.sql)

