Status: ready-for-agent

# Directory tree row style

## Problem Statement

The shared directory tree (`app-dir-tree-nav`) and the Zip preview dialog's tree (`app-zip-tree-preview`) render the same kind of hierarchy in two visual languages. The directory tree keeps Material's 48px node height and padding, prints a literal `/` before every name, lets long names open a horizontal scrollbar inside the tree box, and has no hover or indent guide. The Zip tree uses compact 32px flex rows, one fixed leading icon slot, ellipsized names, a hover tint, and a 19px indent with a faint guide line. In a deep path the directory tree is the harder one to scan.

## Solution

Give the directory tree's rows the Zip tree's anatomy, in place. The panel, search field, selection highlight, expansion and reveal, prefill, and emits keep their current behaviour. Only the row markup and stylesheet change: one flex row per directory, a fixed leading icon slot where the folder glyph doubles as the toggle, an ellipsized name, a hover tint, and an indented branch with a guide line. Names lose the leading slash; the hierarchy already states the path.

## Implementation Decisions

- **Row anatomy.** Copy `app-zip-tree-preview`'s row rules into `dir-tree-nav.component.css`: 32px minimum height, 4px radius, flex row, hover background `rgba(255, 255, 255, 0.06)`.
- **Icon slot.** A 32px leading column for every node. A folder with children uses an icon button that toggles and shows `bi-folder2-open` / `bi-folder-fill`. A folder without children shows `bi-folder-fill` in the same slot, not interactive. The name starts at the same x on every row at a level.
- **Name.** One line, ellipsized, full name in `title`. No `/` prefix. The inline `text-align`, `white-space`, and `overflow` hack is removed.
- **Branch indent.** `.dir-tree ul` gets 19px left padding with a 1px `rgba(255, 255, 255, 0.08)` left border, list style off. Material's node minimum height is cancelled inside `.dir-tree`.
- **Panel scroll.** The tree box keeps its 500px cap and changes `overflow: scroll` to `overflow: auto`, so the vertical bar appears only when the tree is taller than the box and no horizontal bar is reserved.
- **Phones.** At 768px and below, rows and the icon slot grow to 36px, indent tightens to 13px, matching the Zip tree. The old 200px name cap is removed; ellipsis bounds the name at any width.
- **Unchanged behaviour.** Search filtering, bold match styling, selected class and `aria-selected`, ancestor expansion, centring scroll, search prefill, fetch-once, synchronous click emit, and deferred input emit stay exactly as they are. Node order stays the server's order. No size column: directory nodes carry no size.

## Contracts and Invariants

- Public surface unchanged: inputs `title` and `selectedFileKey`; outputs `selected` and `treeLoaded`.
- A click on a name selects and emits; a click on the toggle only expands or collapses.
- At most one label carries `dir-tree-name-selected`, and it still renders `rgba(233, 30, 99, 0.16)`.
- Search still collapses to matches, and clearing it restores the selected path.
- No new HTTP call, no change to `DirTree`, no schema or payload change.

## Acceptance Criteria

- [ ] Every directory row renders as one `.dir-tree-row` with a leading icon slot and an ellipsized name.
- [ ] Folders with children toggle on the folder glyph; leaf folders show the same glyph, non-interactive, in the same slot.
- [ ] Nested levels indent 19px with a faint guide line; at 768px and below, rows are 36px and indent is 13px.
- [ ] Names render without the leading slash and expose the full name through `title`.
- [ ] The tree box has no horizontal scrollbar and scrolls vertically only when needed.
- [ ] Search, bold matches, selection highlight, `aria-selected`, ancestor expansion, centring reveal, prefill, emits, and fetch-once behave as before.
- [ ] The component spec covers plain-name rendering, the tooltip, the row/slot shape, and the unchanged selection/search/scroll behaviour.

## Testing Decisions

The existing `dir-tree-nav.component.spec.ts` is the seam: it already drives the component through HTTP, clicks, search, and computed style. Text expectations change from `/name` to `name`; new assertions cover the tooltip and the row's leading slot. Hover, guide line, and breakpoint sizes are stylesheet-only and stay under the repo's breakpoint lint, as the current mobile rule does. Verification runs `npm test -- --watch=false --browsers=ChromeHeadless`, `ng lint`, and `npm run lint:breakpoints`. A browser check is only for judging the final look, and is asked for separately.

## Out of Scope

- The Zip dialog shell: title, entry count, truncation notice, size column, close action.
- Backend, `DirTree` service, tree ordering, selection logic, or expansion defaults.
- `app-zip-tree-preview` itself.
- New i18n keys or other user-facing copy.
- Keyboard navigation, multi-select, or type-ahead.
