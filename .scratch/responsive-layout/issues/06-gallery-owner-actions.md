# 06: Gallery owner actions

**What to build:** Update writes the name typed in the wide row's panel, and the statement it runs sets only `name` and `update_by`. On a phone, the row's kebab menu offers browse and directory only; update, delete and share leave the compact route by product decision.

**Blocked by:** 05.

**Status:** ready-for-human

- [x] Update in the wide row's expanded panel posts the name typed in that panel.
- [x] `UpdateGallery` runs `UPDATE gallery SET name = ?, update_by = ? WHERE gallery_no = ?`; no other column is written.
- [x] A compact row's kebab menu offers browse and directory only.
- [x] `frontend/moon/doc/agents/responsive.md` records the compact route as it now stands, and names the three dropped actions as a product decision.
- [x] `npm run lint:breakpoints`, the full test suite and the production build (`NODE_OPTIONS=--openssl-legacy-provider ng build`) pass.
- [ ] Confirmed on a 390px viewport on the user's phone.
