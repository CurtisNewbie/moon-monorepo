# 01: Folder-with-cover glyph

**What to build:** The visual for a directory row that has a cover image: the folder silhouette with the cover inset inside it, so the row reads as a directory holding a picture. Drawn and shown to the user for approval before any implementation work starts — this ticket's deliverable is the design, not the wiring.

**Blocked by:** None (can start immediately)

**Status:** done

- [x] The composite shows a folder shape with its tab visible; the cover sits inside the folder body, clipped to the folder's inner area — never a bare thumbnail with a folder decoration on top
- [x] It reads as a folder, not as an image, at the Preview column's existing row height (60px) and cell width
- [x] The plain folder glyph for a directory with no cover is unchanged, and the two read as the same family beside each other
- [x] File rows keep their current plain thumbnail, so a file and a folder-with-cover are visually distinct
- [x] Works in both the wide and the compact (≤768px) layouts without touching the neighbouring columns
- [x] Delivered as a reviewable artifact (a rendered mock or an HTML/CSS prototype at real row size) that shows folder-with-cover, plain folder, and an image file side by side
- [x] The user has approved the artifact; the approved approach is captured in the ticket's comments so the implementer builds that one

## Comments

Approved 2026-10-31 — mock `design/folder-glyph.html` (self-contained; the rendered PNG was a review artifact and is not kept), user approved as-is.

Approved CSS approach, reproduce this:
1. `.dir-cover` = folder body: 44x34, `border-radius: 0 4px 4px 4px` (square top-left, where the tab joins), light gradient `#f4edf1 -> #c9c0c7`, `filter: drop-shadow(0 1px 1.5px rgba(0,0,0,.5))`.
2. `::before` = tab: 20x9 at `left:0; top:-7px` (2px tucked into the body), `border-radius: 4px 0 0 0`, `clip-path: polygon(0 0, calc(100% - 3px) 0, 100% 100%, 0 100%)`.
3. cover `<img>`: absolutely pinned 3px left/right/bottom, 6px under the tab, `object-fit: cover`, `border-radius: 2px`; the frame clips it.
4. Thicker top rail (6px) than side rails (3px) — the knob that makes tab + frame read as one folder rather than a photo mat.
5. `margin: 3px 0` keeps the row at 48px, so glyph -> cover swaps never shift the list. Cover box 38x25, silhouette 44x41 (tab 7px proud).

Deliberate deviation: the composite is ~2x the 20px plain glyph (a cover inside 20px is unreadable). Plain folder glyph and file thumbnail markup unchanged.
