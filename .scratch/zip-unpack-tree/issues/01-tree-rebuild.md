# 01: Rebuild the archive's folder tree when unpacking

**What to build:** Unpacking an archive recreates the folder layers it was packed with, instead of dropping every entry flat into the target directory with its whole path as its name. A nested entry becomes a file under freshly created ancestor folders, and folders that already exist in the target directory are reused.

Entry names that cannot be honoured are dropped rather than rewritten, and each drop is logged with the path and the reason. Name clashes resolve through the same numbered convention the zip download action uses.

Nothing is written outside the target directory, and every row carries the requesting user as uploader.

**Blocked by:** None (can start immediately)

**Status:** done

- [x] An entry at `Documents/doc/jdk-12.0.2_doc-all/docs/api/a.html` produces a chain of folders under the target directory, with `a.html` as a child of the deepest one and named only `a.html`.
- [x] Entries sharing ancestors create each folder once; the second entry reuses it.
- [x] A folder that already exists in the target directory under the same name is reused rather than duplicated.
- [x] A folder name held by an existing file gets a numbered suffix and its subtree still lands under it.
- [x] Two entries resolving to the same file path get ` (n)` before the extension on the later one, counting from one, in entry order.
- [x] An entry whose path escapes the target directory, contains a `..` segment, carries a control character, ends a segment with a dot or a space, or has a segment longer than 255 bytes is dropped and logged with its path and reason.
- [x] A leading separator is dropped and a backslash is read as a separator, so a Windows packed archive unpacks.
- [x] Archive directory entries never produce a row.
- [x] The path logic lives in a pure planner and is covered by table-driven tests that need no database.
- [x] The existing unpack path still writes rows through the current directory and file helpers, with the same ownership and parent rules.
