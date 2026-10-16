# 01: Extract the effective-sort resolver

**What to build:** One backend resolver that answers what sort order a directory is actually listed by, given the requested sort field and direction and the directory itself: a comic directory is forced to name-ascending, and an omitted direction resolves to ascending for name and custom ordering and descending for time. Both the directory listing and the file-position calculation call this resolver instead of each carrying its own copy of the rule. Observable behavior is unchanged.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] Directory listing and file-position calculation both resolve their sort through the one resolver; no inline copy of the rule remains.
- [ ] A comic directory lists and calculates under name-ascending regardless of the requested field and direction, exactly as before.
- [ ] An omitted direction resolves to ascending for name and custom, and to descending for time.
- [ ] Existing listing and position behavior is unchanged.
