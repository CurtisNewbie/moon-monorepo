# 03: Search returns to the selected directory

**What to build:** After a user searches the tree and then clears the search box, the tree comes back to the selected directory's path with that directory still highlighted, instead of a flat collapsed tree. While a search is active, the selected node stays marked wherever it appears among the matches.

**Blocked by:** 01

**Status:** done

- [x] Clearing the search box re-expands every ancestor between the root and the selected node.
- [x] The selected node keeps its highlight after the search is cleared and while results are shown.
- [x] Search matching, match highlighting, and result visibility behave as they do today.
- [x] The component spec covers clear-search restoring the selected path and the highlight surviving an active search.