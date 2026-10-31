# 01: Tree marks and reveals a selected directory

**What to build:** The shared directory tree accepts a selected directory key and shows it. On the Manage Drone Tasks panel, the remembered Parent Directory comes back on screen with its ancestors expanded, its label scrolled into view, and a selected style marking exactly that node. The same key applied twice does not emit twice, and a key that no longer exists in the tree leaves no highlight.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] The tree input that carries the selected directory is named `selectedFileKey` at every call site.
- [ ] Setting a key before the tree data arrives applies the selection once the data lands, and setting it after arrival applies it immediately.
- [ ] Every ancestor between the root and the selected node is expanded, including grouping nodes that carry no file key.
- [ ] The selected node's label carries a selected style class and `aria-selected="true"`, and no other node does.
- [ ] The selected node is scrolled into view inside the tree's scroll box.
- [ ] Applying the same key twice emits `selected` once.
- [ ] A key absent from the loaded tree clears the highlight and emits nothing.
- [ ] Selecting a node programmatically still emits `selected`, so Manage Drone Tasks fills its Parent Directory name.
- [ ] A new component spec covers the cases above through inputs, outputs, expansion state, and the rendered class, and runs under the frontend test command.