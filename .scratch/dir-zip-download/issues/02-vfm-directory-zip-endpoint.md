# 02: vfm: directory zip download endpoint

**What to build:** An owner can ask vfm for a download of one directory and receive a bundle key that, once fetched, yields a zip of that directory's direct file children. Children are ordered by name ascending. Sub-directories are skipped entirely, and nothing inside them reaches the archive. Duplicate child names are both kept, the later one suffixed ` (n)`. A directory that is not owned by the caller, does not exist, is not a directory, or is logically deleted is refused. A directory with no direct files still returns a working key that produces an empty zip.

**Blocked by:** 01

**Status:** ready-for-agent

- [ ] Requesting a zip for an owned directory returns a key that downloads the directory's direct files.
- [ ] Only direct file children are included: no sub-directory entries, no files from inside sub-directories.
- [ ] Entries appear in name-ascending order.
- [ ] Two children sharing a name both survive, the later suffixed ` (n)` before its extension, resolved the same way on every request.
- [ ] A directory owned by another user is refused.
- [ ] A missing directory, a file row, and a logically deleted directory are each refused.
- [ ] An empty directory returns a key that downloads a valid empty zip.
- [ ] The archive name follows the directory name, including non-ASCII names.
- [ ] The existing single-file download token endpoint still refuses directories.
