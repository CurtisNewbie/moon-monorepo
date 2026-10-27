# 02: Tables follow the layout mode

**What to build:** Every table shows the columns that fit the current layout, and follows a rotation or resize. Each screen keeps one desktop column list and one compact list, picks between them from the layout mode at construction, and reassigns when the mode changes. No stylesheet hides a column anymore.

**Blocked by:** 01.

**Status:** ready-for-agent

- [ ] Each listed screen shows its compact column list below the boundary and its desktop list above it: gallery, manage-files, access-log, list-notification, manage-user, manage-keys, manage-logs, manage-paths, manage-resources, manage-role, manage-bookmarks, site-password, drone-task, mng-note.
- [ ] A mode change while the screen is open reassigns the columns without a reload.
- [ ] Subscriptions are torn down with the component.
- [ ] No stylesheet under the app hides a table column by cell class.
- [ ] No column list is stored from a width read made anywhere except the layout signal.
- [ ] The tables that already have specs assert both column sets through the shared media query fake, including a mid-life flip.
- [ ] `npm run lint` and the full test suite pass.
- [ ] Confirmed on a 390px viewport on the user's phone.
