# 01: Responsive foundation

**What to build:** One width signal and the rules around it. After this ticket, `Env` answers "is this the compact layout" from a single media query list, `mobile$` reports when that answer changes, the global stylesheet owns the only two breakpoint blocks, and lint fails when anyone drifts from them. The dialog helper, the pinch-zoom fix and the agent rule doc land here so the migration batches have something to follow.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [ ] `Env` reads one media query list for the 768px boundary; `isMobile()` returns its current match and `mobile$` emits the same value, so the two cannot disagree.
- [ ] `mobile$` starts with the current value, emits on each change, and stays silent when the media query reports an unchanged value.
- [ ] The global stylesheet holds one compact block and one wide block, holding today's type-scale rules, plus the compact-only and wide-only display utilities.
- [ ] A checker run from the lint script fails when the service's boundary and the stylesheet's two block widths disagree, when any component stylesheet declares a width other than 768 or 769, or when anything outside the service reads the window width directly.
- [ ] One dialog-width helper produces a width capped by both its pixel value and the viewport, with a matching maximum width, and is the documented way to size a dialog.
- [ ] The viewport meta keeps device width and initial scale and no longer pins the scale, so pinch zoom works.
- [ ] The rule doc exists, is named from `frontend/AGENTS.md` the way the HTTP, i18n and pagination docs are, and states the decision table (cosmetic is CSS, structural is a subscription, columns are TypeScript), the banned patterns, the column snippet, the dialog helper and the checker.
- [ ] A test helper fakes the media query list so a test can set the current match and fire a change on demand; the service has specs covering the initial value, a change, and a no-op change.
- [ ] `npm run lint` and the full test suite pass.
