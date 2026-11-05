# 02: Refuse an archive that is too deep, before extraction starts

**What to build:** Starting an unpack checks how deep the archive's tree goes. An archive deeper than the configured limit is refused in the same call that starts the unpack, so nothing is extracted and no rows are written. An archive that cannot be read is refused the same way.

The limit is a vfm configuration value, sent with the unpack request. mini-fstore holds no policy of its own and falls back to the same default when a caller omits it. The check reads the archive's central directory, so it never decompresses.

**Blocked by:** None (can start immediately)

**Status:** done

- [x] The unpack request in mini-fstore carries the maximum depth, where zero means the server default.
- [x] mini-fstore refuses the request synchronously when any file entry sits deeper than the limit, with its own failure code distinct from an unreadable archive.
- [x] A corrupt or truncated archive is refused synchronously with the existing illegal format code.
- [x] Nothing is enqueued for a refused archive, so no extraction runs.
- [x] Depth counts path segments in file entries, ignores empty segments, ignores directory entries, and counts the archive's own top folder as level one.
- [x] vfm reads its own configuration and sends the value; the effective default is 32 and every change applies from the next unpack.
- [x] The depth check is covered by a pure function tested with table cases, plus one test that builds a real zip and asserts both sides of the limit.
