# 05: Directory Browse History lists every cursor, with its sort

**What to build:** The Directory Browse History list stops collapsing to one entry per directory. Every stored cursor row — one per user, directory and sort — is its own entry, most recently updated first, still excluding directories marked deleted. Each entry shows the sort order it was recorded under, in its own column, and its Resume restores that entry's sort and file.

**Blocked by:** 04

**Status:** ready-for-agent

- [ ] A directory read under several sorts appears once per sort, newest first, with no grouping by directory.
- [ ] Each entry shows the sort order it was recorded under.
- [ ] Directories marked deleted still do not appear.
- [ ] Resume opens the directory under that entry's sort and lands on that entry's file.
- [ ] The entry cap counts entries, not directories.
