# 05: Outliers

**What to build:** The five screens whose compact handling does not fit the general pattern. The image viewer loses its private width tiers and sizes itself against the screen. The directory tree nav takes its width from the layout mode instead of computing one inline. The paginator stops removing its page-size control on a phone and keeps every control reachable. The gallery's row click stays a behaviour swap rather than a hidden action, and is written down as one.

**Blocked by:** 01.

**Status:** done

- [x] The image viewer declares no width breakpoint of its own and fills the screen at 390x844 and at desktop sizes, keeping its long-portrait and tiny-image scaling behaviour.
- [x] The gallery's create control reaches a phone user, and every owner action (update, delete, share, directory) plus browse is reachable from a compact row.
- [x] The drone task list's update action is reachable on a phone.
- [x] The directory tree nav takes its width from the layout mode; no width is computed inline from a window read.
- [x] The paginator offers the page-size control at 390px; every control it renders on desktop is reachable on a phone.
- [x] The gallery's compact row click is documented in the rule doc as a behaviour swap, and every action reachable from a desktop row is reachable on a phone by some path.
- [x] The checker's allow-list names no file once this ticket lands.
- [x] The breakpoint checker, the full test suite and the production build (`NODE_OPTIONS=--openssl-legacy-provider ng build`) pass.
- [x] Confirmed on a 390px viewport on the user's phone.
