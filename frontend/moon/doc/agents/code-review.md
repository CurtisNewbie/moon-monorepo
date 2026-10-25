# Code Review — Frontend

The frontend review contract, named from `frontend/AGENTS.md`.

## Membership

- `code-review` — the two-axis pass (Standards + Spec) run on every diff.
- `frontend/AGENTS.md` requires no skill of its own, so the membership is the review pass alone.

## Standards sources

- `frontend/AGENTS.md` — Angular 11 / NgModule, TSLint, test commands, i18n, pagination.
- `frontend/moon/doc/agents/` — `http.md`, `i18n.md`, `pagination.md`.
- `frontend/moon/tslint.json` — the lint rules tooling enforces.

## Reviewing a diff

Walk every entry above against the diff, then report per axis.
