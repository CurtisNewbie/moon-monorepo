# Code Review — Backend

The backend review contract, named from `backend/AGENTS.md`.

## Membership

- `miso` — mandatory before any backend Go code. Load the skill, then walk
  `references/code-review-checklist.md` entry by entry against the diff.
- `code-review` — the two-axis pass (Standards + Spec) run on every diff.

## Standards sources

- `backend/AGENTS.md` — service layout, build order, framework rules.
- `docs/agents/db-schema.md` — schema snapshot and release manifest rules.
- `docs/adr/` — decisions the diff must respect.

## Reviewing a diff

Walk every entry above against the diff, then report per axis.
