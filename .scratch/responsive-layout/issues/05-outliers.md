# 05: Outliers

**What to build:** The five screens whose compact handling does not fit the general pattern. The image viewer loses its private width tiers and sizes itself against the screen. The directory tree nav takes its width from the layout mode instead of computing one inline. The paginator stops removing its page-size control on a phone and keeps every control reachable. The gallery's row click stays a behaviour swap rather than a hidden action, and is written down as one.

**Blocked by:** 01.

**Status:** ready-for-agent

- [ ] The image viewer declares no width breakpoint of its own and fills the screen at 390x844 and at desktop sizes.
- [ ] The directory tree nav takes its width from the layout mode; no width is computed inline from a window read.
- [ ] The paginator offers the page-size control at 390px; every control it renders on desktop is reachable on a phone.
- [ ] The gallery's compact row click is documented in the rule doc as a behaviour swap, and every action reachable from a desktop row is reachable on a phone by some path.
- [ ] `npm run lint` and the full test suite pass.
- [ ] Confirmed on a 390px viewport on the user's phone.
