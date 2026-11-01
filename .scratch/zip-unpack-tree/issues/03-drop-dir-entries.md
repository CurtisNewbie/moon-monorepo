# 03: Stop storing archive directory entries as files

**What to build:** Directory entries inside an archive stop becoming stored files. A zip entry that is a folder, such as `Documents/`, no longer creates a stored file with zero bytes, no longer appears in the unpack reply, and no longer becomes a row in either service.

**Blocked by:** None (can start immediately)

**Status:** done

- [x] Extracting an archive that contains directory entries yields no stored file for any of them.
- [x] No reply entry is produced for a directory entry.
- [x] File entries in the same archive are extracted exactly as before, with their paths intact.
- [x] Covered by a test that builds a real zip with both directory and file entries.
