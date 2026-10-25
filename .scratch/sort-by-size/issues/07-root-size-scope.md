# 07: A size order at root spans every directory

**What to build:** When the effective sort is size and the request names no directory, the listing drops the parent filter and returns every file the user owns, in every directory, files only, ordered by size and then by row id in the requested direction. A size order that does name a directory stays scoped to it. Time, name and custom at root keep listing root-level rows only. The position calculation already treats an absent parent as no parent filter, so it needs no change and already lands on the page this listing shows.

**Blocked by:** none

**Status:** done

- [x] A size order with no directory returns files from several directories, ordered by size descending
- [x] The direction control still flips the order, and row id breaks ties across directories so paging is stable
- [x] A size order with a keyword at root ranks only the matches, across every directory
- [x] A size order inside a directory returns only that directory's files
- [x] The listing stays files only under a size order, with no directory rows
- [x] Time, name and custom orders at root still return root-level rows only
- [x] Hidden, logically deleted and sensitive rows stay out
- [x] No resume cursor is recorded or offered under a size order, at root or in a directory
- [x] Tests at the listing seam cover the root listing, the in-directory listing, the keyword case, and the three orders that do not become global
