# AGENTS.md

## Dev workflow

Load the `dev-flow` skill before starting any task — a feature, a bugfix, or workflow setup — and follow its lane (unless user chooses not to).

## Repository Structure

Manual monorepo: Go microservices in `backend/<service>/`, Angular app in `frontend/moon/`, Docker Compose configs in `deploy/`.

- Go services, backend config, build order, DB schema, deployment → read `backend/AGENTS.md` first.
- Angular app → read `frontend/AGENTS.md` first.

## Agent skills

### Issue tracker

Issues and specs live as local markdown under `.scratch/`; `docs/agents/issue-tracker.md` covers the conventions and how skills read and write them.

### Triage labels

The five canonical triage labels (default strings). `docs/agents/triage-labels.md` maps each to this repo's label.

### Domain docs

This is a single-context repo: root `GLOSSARY.md` + `docs/adr/`. Read them before exploring an unfamiliar area; `docs/agents/domain.md` covers how the skills consume them.
