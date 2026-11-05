# 02: Selection follows clicks and picking stops collapsing

**What to build:** Clicking a directory in the tree marks it and leaves the tree where the user navigated. Picking a Parent Directory on Manage Drone Tasks, or a target on fape gallery creation, no longer collapses the whole tree, so the chosen node stays visible and a sibling or child is one click away.

**Blocked by:** 01

**Status:** done

- [x] Clicking a node marks it as selected and emits it, with no collapse.
- [x] Manage Drone Tasks no longer collapses the tree after a directory is picked, and the selected node stays highlighted.
- [x] Fape gallery creation no longer collapses the tree after a directory is picked, and the selected node stays highlighted.
- [x] The tree component's public collapse method and the view queries that existed only to call it are gone, with no remaining callers.
- [x] The Manage Drone Tasks page spec covers the last-selected-directory round trip and a click that leaves the tree expanded.
- [x] Existing Manage Drone Tasks search tests still pass.