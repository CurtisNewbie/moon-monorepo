# 02: Selection follows clicks and picking stops collapsing

**What to build:** Clicking a directory in the tree marks it and leaves the tree where the user navigated. Picking a Parent Directory on Manage Drone Tasks, or a target on fape gallery creation, no longer collapses the whole tree, so the chosen node stays visible and a sibling or child is one click away.

**Blocked by:** 01

**Status:** ready-for-agent

- [ ] Clicking a node marks it as selected and emits it, with no collapse.
- [ ] Manage Drone Tasks no longer collapses the tree after a directory is picked, and the selected node stays highlighted.
- [ ] Fape gallery creation no longer collapses the tree after a directory is picked, and the selected node stays highlighted.
- [ ] The tree component's public collapse method and the view queries that existed only to call it are gone, with no remaining callers.
- [ ] The Manage Drone Tasks page spec covers the last-selected-directory round trip and a click that leaves the tree expanded.
- [ ] Existing Manage Drone Tasks search tests still pass.