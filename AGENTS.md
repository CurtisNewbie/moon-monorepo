# AGENTS.md

## Dev workflow

Load the `dev-flow` skill before starting any task — a feature, a bugfix, or workflow setup — and follow its lane (unless user chooses not to).

## Repository Structure

Manual monorepo: Go microservices in `backend/<service>/`, Angular app in `frontend/moon/`, Docker Compose configs in `deploy/`.

- Go services, backend config, build order, DB schema, deployment → read `backend/AGENTS.md` first.
- Angular app → read `frontend/AGENTS.md` first.

### Sibling project: drone

`drone` is a private Go backend service (miso) that belongs to this system but is not in this monorepo — its code lives at `../drone`, it is not published publicly and whether it will be is undecided. It is optional: nothing here depends on it. It calls this deployment's services at runtime (registered in `user-vault`, logs collected by `logbot`) and is consumed by the Angular app; its `doc/agents/` covers architecture, API and config. When a task touches drone, work in `../drone` and read its `AGENTS.md` first.

**Drone work stays in drone.** Specs and tickets for a drone feature live in `../drone`, never in this repo's `.scratch/`; ADRs about drone go to `../drone/docs/adr/`, not `docs/adr/` here. Conventions: `../drone/docs/agents/issue-tracker.md`. This repo only records the cross-cutting facts: shared API contracts, DB schema shared with these services, and the release/deploy order.

## Release

Cutting a release — version bump, schema manifest, tag, GitHub notes, deploy order — follow `release/README.md`. Schema history lives in `release/db-changes.md`; per-release notes in `release/<version>.md`.

## Agent skills

### Issue tracker

Issues and specs live as local markdown under `.scratch/`; `docs/agents/issue-tracker.md` covers the conventions and how skills read and write them. Drone's specs and tickets are not tracked here — see [Sibling project: drone](#sibling-project-drone).

### Triage labels

The five canonical triage labels (default strings). `docs/agents/triage-labels.md` maps each to this repo's label.

### Domain docs

This is a single-context repo: root `GLOSSARY.md` + `docs/adr/`. Read them before exploring an unfamiliar area; `docs/agents/domain.md` covers how the skills consume them.

### Code review

Repo-wide review contract: skills to load on every run, and their review rules. Backend and frontend code answer to their own contracts. See `docs/agents/code-review.md`.
