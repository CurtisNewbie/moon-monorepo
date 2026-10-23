# 04: API doc and v0.0.6 release note

**What to build:** The vfm API documentation states that size is a sort value, that it applies only to the File type filter, and that it is not persisted. The v0.0.6 release note records the feature for users. No schema snapshot, schema increment, or change log entry is needed, because no table changes.

**Blocked by:** 01, 02, 03

**Status:** ready-for-agent

- [ ] The vfm API doc lists size among the sort values and states the File-only rule
- [ ] `release/v0.0.6.md` exists and records the size sort from the user's point of view
- [ ] `release/db-changes.md` gains no v0.0.6 entry
- [ ] No schema snapshot or versioned DDL file is added
