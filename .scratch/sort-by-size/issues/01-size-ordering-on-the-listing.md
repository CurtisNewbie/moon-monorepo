# 01: Size ordering on the listing

**What to build:** A size sort on the file list. Asking for it returns the directory's files ordered by size, largest first unless ascending is requested. Because directories have no meaningful size, the same request is coerced to the File type filter, so the result never mixes directories in. Equal sizes fall back to row id in the same direction, so pagination stays stable. Keyword search keeps filtering and size ranks the matches. A comic directory still opens name-ascending.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] `orderBy=size` returns files ordered by size, descending when the direction is omitted
- [ ] An explicit ascending direction is honored
- [ ] Equal sizes are ordered by row id in the same direction, stable across pages
- [ ] A size request with the type filter absent, All, or Directory returns files only
- [ ] A size request combined with a name keyword returns the keyword matches ordered by size
- [ ] A comic directory requesting size returns name-ascending
- [ ] No schema change is introduced
- [ ] Tests at the listing seam cover size order, tie-break, coercion, keyword combination, and the comic force
