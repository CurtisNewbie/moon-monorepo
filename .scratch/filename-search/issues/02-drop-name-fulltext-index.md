# 02: Drop the file name fulltext index

**What to build:** Once the listing matches names by substring, the fulltext index on the file name column of `file_info` serves nothing. Drop it in the in-flight release and record the change, so the schema stops advertising a search path the service no longer uses. The fulltext index on `versioned_file` stays, because versioned-file search keeps its current match.

**Blocked by:** 01

**Status:** done

- [x] The v0.0.6 increment drops the fulltext index from `file_info`
- [x] The current snapshot for the service carries the table without that index
- [x] `release/db-changes.md` lists the increment file under the v0.0.6 section
- [x] The `versioned_file` fulltext index is untouched in both the increment and the snapshot
- [x] The change is applied by hand, with no migration tooling or auto-migration added
