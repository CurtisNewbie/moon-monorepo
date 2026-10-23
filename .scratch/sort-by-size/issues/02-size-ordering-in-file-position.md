# 02: Size ordering in file position

**What to build:** The position calculation understands a size sort. Given a file, a directory, and the current size sort, it answers the page that file occupies in the size-ordered listing, under the same File-only coercion as the listing. This is what makes a search-result jump or a deep link land on the page the user is actually looking at, even though a size sort produces no resume cursor.

**Blocked by:** 01

**Status:** ready-for-agent

- [ ] A position request under size returns the page the file occupies in the size-ordered listing, both directions
- [ ] A position request under size applies the File-only type coercion
- [ ] The result matches the page the listing actually places the file on
- [ ] No resume cursor is created or read as part of the position calculation
- [ ] Tests at the position seam cover both directions and the coerced type filter
