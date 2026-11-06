# 01: Read a Zip's tree in mini-fstore

**What to build:** mini-fstore gains an internal endpoint that returns the nested folder and file tree of a stored Zip file. It reads the Zip central directory only. No entry is decompressed, no temporary file is created, and nothing is written.

The reply is one root folder node with an empty name. File nodes carry the entry name and its uncompressed size. Folder nodes carry children and no size. Siblings are ordered folders first, then files, each group case-insensitive by name. Duplicate names keep their archive order.

An explicit directory entry becomes a folder node, so an empty folder appears. A folder inferred from a file path merges into the same node. A repeated directory entry merges. Every file entry becomes its own node, so a repeated file path appears twice. A file and a folder at the same path appear as two nodes.

Both `/` and `\` are separators. Leading separators and empty segments are ignored. Every other character is kept as text, including `..`, drive prefixes, control characters, and long names.

The request carries a maximum depth and a maximum entry count, where zero means the server default of 32 and 10000. Entries deeper than the depth are omitted and the reply says it was truncated by depth. Once the entry cap is reached, the scan stops, later entries are omitted, and the reply says it was truncated by entry count. The reply reports the included entry count and not the total.

A missing or logically deleted stored file returns `FILE_NOT_FOUND`. A file that cannot be read as a Zip returns `ILLEGAL_FORMAT`.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] The endpoint reads the central directory without decompressing an entry, creating a temporary file, or writing anywhere.
- [ ] The reply is one root folder node with an empty name, containing file nodes with uncompressed sizes and folder nodes with children and no size.
- [ ] Siblings are folders first, then files, each group case-insensitive by name. Duplicate names keep archive order.
- [ ] Explicit directory entries become folder nodes. Inferred parents and repeated directory entries merge. Repeated file entries appear once each. A file and folder at one path appear as two nodes.
- [ ] `/` and `\` both act as separators. Leading separators and empty segments are ignored. All other characters are kept and returned as text.
- [ ] Entries deeper than the request maximum are omitted and `truncated` is true with `depth` in `truncatedReasons`.
- [ ] No more than the request maximum entries are represented. The scan stops at the cap, `truncated` is true with `entries` in `truncatedReasons`, and `includedEntries` matches the reply.
- [ ] A zero or absent depth and entry count fall back to 32 and 10000.
- [ ] A missing or logically deleted stored file returns `FILE_NOT_FOUND`. An unreadable Zip returns `ILLEGAL_FORMAT`.
- [ ] The pure tree builder is covered by table tests. They cover nesting, both separators, leading separators, `..` and control characters, explicit folders, empty folders, repeated files, repeated folders, file and folder at one path, ordering, depth boundary, entry boundary, empty archive.
- [ ] One test builds a real Zip and covers a normal tree, an explicit directory entry, a backslash path, a non-Zip file, a truncated file, and a cap small enough to truncate.
