# 04: Remembered directory prefilter

**What to build:** Opening Manage Drone Tasks on a remembered Parent Directory shows only that directory's branch instead of the whole tree. The tree's search box carries the directory's name, the matching directories and the path down to them are rendered, and the remembered one is highlighted. Clicking a different directory never refills the search box, and clearing it brings back the whole tree with the selection still marked.

**Blocked by:** 01, 03

**Status:** ready-for-agent

- [ ] Applying a key through the tree's input fills the search box with that directory's name.
- [ ] The tree then renders only directories whose name matches and the ancestor path down to each match, with the selected one highlighted.
- [ ] A directory name that occurs more than once shows every match with its own ancestor path, and only the applied key's node is highlighted.
- [ ] Clicking a node leaves the search box untouched and the tree keeps whatever it was showing.
- [ ] Clearing the search box restores the whole tree and the selected directory's path, without emitting `selected`.
- [ ] The component spec covers the prefill, the click does not prefill, and the multi-match case.
- [ ] The Manage Drone Tasks page spec keeps passing with the prefilter in effect.