# 04: Directory Browse History lists each directory once, with faithful resume

**What to build:** The Directory Browse History list shows each directory once, taken from its most recently updated cursor row, and leaves out directories marked deleted. Each entry carries the sort order and file of that visit, and its Resume opens the directory under that sort and lands on that file, instead of the fixed name sort it forces today.

**Blocked by:** 02

**Status:** done

- [ ] A directory read under several sorts appears once, taken from its most recent visit.
- [ ] Directories marked deleted do not appear in the list.
- [ ] Resume from the list opens the directory under that entry's sort and lands on that entry's file.
- [ ] Resuming from the list no longer forces a fixed name sort.
