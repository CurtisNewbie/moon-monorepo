# 04: Header and page shell on a phone

**What to build:** The header keeps its controls on one line at 390px without pushing content off the screen, and the page container sizes itself to the visible viewport rather than a fixed viewport height. Nothing is dropped from the header; the controls change size or presentation if they do not fit.

**Blocked by:** 01.

**Status:** ready-for-agent

- [ ] At 390px the header renders on one line with no horizontal scrolling on the page.
- [ ] Every header control present on desktop is still reachable on a phone, through the same element or a different one.
- [ ] The page container uses the dynamic viewport height, with the current fixed height kept as the fallback.
- [ ] The desktop header looks as it does today.
- [ ] `npm run lint` and the full test suite pass.
- [ ] Confirmed on a 390px viewport on the user's phone.
