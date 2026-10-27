# 03: Dialogs fit the screen

**What to build:** Every dialog opens at a width the viewport can hold, and its content stays reachable at 390px. Each open call takes its width from the shared helper, so a dialog is at most its designed width and never wider than the screen.

**Blocked by:** 01.

**Status:** ready-for-agent

- [ ] Every dialog open call in the app passes its width through the helper.
- [ ] No fixed pixel dialog width remains in an open call.
- [ ] At 390px each dialog fits the viewport, and its content scrolls rather than being clipped.
- [ ] At desktop widths each dialog keeps its current designed width.
- [ ] `npm run lint` and the full test suite pass.
- [ ] Confirmed on a 390px viewport on the user's phone.
