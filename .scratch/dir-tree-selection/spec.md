# Directory Tree Selection

Status: ready-for-agent

## Problem Statement

On Manage Drone Tasks, a user can only choose the Parent Directory by clicking a node in the embedded directory tree. The page already remembers a default: it asks the backend for the last selected directory and passes that key to the tree, which then re-emits it, so the Parent Directory field is filled in.

The tree itself never shows that choice. There is no selected style, the node's ancestors are never expanded, and both pages that embed the tree collapse every node immediately after a directory is selected. In a deep tree the chosen directory is invisible, and after every pick the user is thrown back to the root. A user reopening the panel cannot see where the task will land without navigating to it again.

## Solution

The directory tree owns its selected node. Handing it a file key selects that node, expands every ancestor between the root and it, scrolls it into view, and marks it with a selected style, so the already chosen Parent Directory is on screen when the panel opens.

Clicking a node does the same: it becomes the highlighted node, its ancestors stay open, and the tree no longer collapses. Clearing a search restores the tree to the selected node's path instead of leaving it flat.

The selection is display state in the shared tree component. What is remembered between visits stays the backend's last selected directory; this feature adds no storage.

## User Stories

1. As a moon user creating a drone task, I want the directory tree to mark the Parent Directory already chosen for me, so that I can see where the task's files will land before submitting.
2. As a moon user opening the New Task panel, I want the tree to expand down to the remembered directory, so that I do not have to click through folders to find it.
3. As a moon user with a deep folder hierarchy, I want the selected directory scrolled into view, so that I see it without hunting inside a scroll box.
4. As a moon user picking a different directory, I want the picked node to stay highlighted after I click it, so that I can confirm the choice before submitting.
5. As a moon user picking a directory, I want the tree to stay where I navigated to, so that choosing a sibling or a child does not send me back to the root.
6. As a moon user searching for a directory, I want the selected node still marked when it appears in the results, so that I can tell the current choice from the matches.
7. As a moon user clearing the search box, I want the tree to return to the selected directory's path, so that I land back on my choice rather than on a collapsed tree.
8. As a moon user creating a fape gallery, I want the same highlighted selection and stable expansion, so that the tree behaves the same way on both creation pages.
9. As a moon user moving a file into a directory, I want the node I click to be marked while the dialog is open, so that I know which target I picked.
10. As a moon user, I want a remembered directory that no longer exists to leave the tree unmarked rather than showing a stale highlight, so that I am not misled about the target.
11. As a moon user whose remembered directory cannot be fetched, I want a clear message and an empty Parent Directory, so that I can pick one myself.
12. As a moon user on a phone, I want the highlighted node readable inside the narrow tree box, so that the marker is visible without horizontal scrolling.
13. As a moon user, I want the tree to load its directory data once per panel open, so that selecting nodes does not trigger repeated requests.
14. As an implementer, I want one component to own the selection rules, so that the three embedding pages cannot drift apart.

## Implementation Decisions

The shared directory tree component, used by Manage Drone Tasks, fape gallery creation, the move-file dialog, and the Manage Files directory dialog, gains the selected-node state. No other component keeps a copy of it.

**Rename the input.** The write-only input `prevSelectedFileKey` becomes `selectedFileKey`, since it now both selects and displays. Every call site is updated in the same change.

**Selection model.** The component stores the selected node itself. Applying a key locates the node in the loaded top-down tree, records it, expands each ancestor between the root and the node with the tree control's `expand`, and applies a selected style class to that node's label. Nodes without a file key are still expanded when they sit on the path.

**Reveal.** After the node is marked, its own scroll box is scrolled so the node sits centred vertically; the page around it never moves. The component has no lazy loading, so the node is always present once the tree has loaded.

**Programmatic selection still emits.** Setting the input emits `selected` for the located node, because the embedding pages use that emission to fill their own fields. A key equal to the one already applied emits nothing, so re-rendering cannot loop. A key that is absent from the tree clears the highlight and emits nothing.

**Clicking.** Clicking a node records and highlights it and emits it. Nothing collapses. The page handlers in Manage Drone Tasks and fape gallery drop their `collapseAll()` calls, and the public `collapseAll()` method is removed from the component along with the view queries that existed only to call it.

**Search.** Search still collapses the tree to show matches. Clearing the search restores the expansion along the selected node's path, so the selection is visible again.

**Highlight style.** The label of the selected node gets a class that renders a distinct background in the theme's pink at low opacity, plus `aria-selected="true"` for screen readers. The component's stylesheet carries the rule.

**Sources of the key.** Manage Drone Tasks keeps resolving the remembered directory from the backend after the tree reports it has loaded. Fape gallery keeps passing no key. No new client-side persistence is introduced.

## Testing Decisions

A good test here drives the component through its inputs and user clicks and asserts what a user would see or receive: the emitted node, which ancestors are expanded, whether the selected class is present, and whether scrolling was requested. Tests do not assert private fields or the BFS implementation.

The primary seam is a new component spec for the directory tree, answering the top-down tree request through `HttpClientTestingModule`. It covers: a key set before load is applied after load and emits once; ancestors are expanded; the selected class lands on exactly one node; a deep key is revealed by scrolling the tree's own box; an unknown key clears the mark and emits nothing; the same input key twice emits once; a click marks and emits even on the node already selected, without collapsing; a pending input emit is dropped when a click supersedes it; clearing a search returns to the selected path without emitting or refetching.

Mobile width is a stylesheet rule with no behaviour of its own, so the repo's breakpoint lint covers it rather than a test.

The existing Manage Drone Tasks spec is extended to cover the page contract: after the last-selected-directory response, the tree receives that key and the Parent Directory field carries the name, and selecting a node leaves the tree expanded. Its current search tests are the prior art for stubbing HTTP and pipes.

Verification runs through the frontend test command, which sets the Node options the toolchain needs, plus `ng lint`. A browser check is not required: the assertions cover DOM class, expansion state, and emission, all reachable in the headless test browser.

## Out of Scope

- Backend changes. The remembered directory endpoint and the top-down tree endpoint keep their contracts.
- Persisting a selected directory in browser storage or in the route.
- Keyboard navigation, multi-select, or type-ahead in the tree.
- Lazy loading of tree branches.
- Changing search matching, matching style, or result ordering.
- The Manage Files breadcrumb path display.

## Further Notes

The top-down tree response is a single root, and the root itself is selectable like any other node, so the root can carry the highlight.

Fape gallery already passes an input that is only ever undefined, so it gains the selected style and stable expansion but no remembered default. That page's remembered directory is a separate concern.

## Contracts

### Data

No schema or payload changes. The tree nodes keep the shape `DirTopDownTreeNode`: a display `name`, an optional `fileKey`, and an optional `child` array. A node without a `fileKey` is a grouping node and is not selectable; it can still be expanded on the path to a selected node.

### Interfaces

The tree component's public surface:

- Input `title`, a translation key, default `defaultTitle`.
- Input `selectedFileKey`, `string | undefined`, default `undefined`. Writing a key applies the selection; there is no getter.
- Output `selected`, emitting a shallow copy of the selected `DirTopDownTreeNode`. It fires on every click, and on an input-applied selection that differs from the current one. A click emits synchronously; an input application emits after the change-detection pass that bound the input, so the page can update its own bound fields without a changed-after-checked error.
- Output `treeLoaded`, emitting once per successful fetch of the tree data.

No HTTP contract changes: the tree keeps reading the top-down tree, and Manage Drone Tasks keeps reading the last selected directory and creating tasks with a directory key.

### Config

N/A. No new configuration keys.

## Invariants

- At most one node carries the selected style at any time.
- The selected style is never shown for a key that is absent from the loaded tree.
- Whenever a selection is applied, by input or by click, every ancestor between the root and that node is expanded, unless a search is active.
- Applying the same key to the input twice emits once; a click always emits, including on the node already selected.
- A click never collapses the tree.
- The tree data is fetched once per component instance; applying a selection never refetches.
- The highlighted node always matches the key last applied by the input or last clicked.

## States and failure

Before the tree data arrives, an input key is pending. When the fetch succeeds, the pending key is applied, the node is revealed, and `treeLoaded` fires.

If the fetch fails, the tree stays empty, no selection is applied, and `treeLoaded` does not fire. Manage Drone Tasks therefore keeps its Parent Directory empty and its remembered-directory request unstarted; the user picks a directory once a reload succeeds.

If the remembered-directory request fails at the transport level, the tree stays unselected and the Parent Directory stays empty, so the user picks a directory by hand. The page logs the failure and shows no message for this one.

An unknown key clears the highlight and emits nothing, so the page keeps whatever field value it already had and no selection is invented.

There are no retries in this feature. A failed fetch is retried by reopening the panel, which creates a fresh component instance.

## Boundaries

The directory tree service read (`GET /vfm/open/api/file/dir/top-down-tree`) and the drone last-selected-directory read (`GET /drone/open/api/task/last-selected-dir`) are the only external reads. This feature adds no writes and holds no server state.

Both calls are read-only and idempotent; ordering between them matters only in that the tree must report loaded before the page asks for the remembered directory. When either service is down, the tree renders empty or unselected rather than failing the page, and task creation remains available with a directory the user picks after recovery.
