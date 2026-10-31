# 01: Folder-with-cover glyph

**What to build:** The visual for a directory row that has a cover image: the folder silhouette with the cover inset inside it, so the row reads as a directory holding a picture. Drawn and shown to the user for approval before any implementation work starts — this ticket's deliverable is the design, not the wiring.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] The composite shows a folder shape with its tab visible; the cover sits inside the folder body, clipped to the folder's inner area — never a bare thumbnail with a folder decoration on top
- [ ] It reads as a folder, not as an image, at the Preview column's existing row height (60px) and cell width
- [ ] The plain folder glyph for a directory with no cover is unchanged, and the two read as the same family beside each other
- [ ] File rows keep their current plain thumbnail, so a file and a folder-with-cover are visually distinct
- [ ] Works in both the wide and the compact (≤768px) layouts without touching the neighbouring columns
- [ ] Delivered as a reviewable artifact (a rendered mock or an HTML/CSS prototype at real row size) that shows folder-with-cover, plain folder, and an image file side by side
- [ ] The user has approved the artifact; the approved approach is captured in the ticket's comments so the implementer builds that one
