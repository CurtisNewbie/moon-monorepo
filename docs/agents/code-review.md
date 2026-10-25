# Code Review

The repo-wide review contract. Root-level work — deploy, release, docs, cross-service contracts —
answers to this file. Files under `backend/` and `frontend/` answer to their own contracts:
`docs/agents/code-review-backend.md` and `frontend/moon/doc/agents/code-review.md`.

## Membership

- `dev-flow` — the workflow. It carries no review rules, so the review pass skips it.
- `code-review` — the two-axis pass (Standards + Spec) run on every diff.

## Standards sources

- `AGENTS.md` — repo structure, the drone boundary, release rules.
- `docs/adr/` — decisions the diff must respect.
- `release/README.md` — release-time changes.

## Reviewing a diff

Walk every entry above against the diff, then report per axis. A diff spanning contexts walks each
context's contract.
