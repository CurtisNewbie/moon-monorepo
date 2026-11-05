Status: ready-for-agent

# Responsive layout for phones

## Problem Statement

The client runs on phones and on desktop, and the phone case is handled one screen at a time. One service answers whether the window is narrow, and it answers by reading `window.innerWidth` at the moment of the call. Twenty components call it, eighteen times in code and twenty four times in templates. Where the answer is stored instead of read, it is stored when the component is constructed, so rotating the phone or resizing the window leaves the screen on the layout it started with.

Table columns are hidden twice over, once by choosing a narrower column list in TypeScript and once by CSS `display:none` on the same cells, and the two lists disagree. Dialog widths are fixed pixels, so on a 390px viewport the 800px dialog is clamped to Material's 80vw default and its content is cut off. The 768px boundary itself lives in one TypeScript constant and one CSS block, tied together by a comment that asks the next editor to keep them in step.

The cost is drift. Each new screen picks its own mechanism, its own threshold and its own idea of what a phone user does not need, and no rule stops it.

## Solution

One width signal owned by `Env`, reactive and single-valued. `isMobile()` answers synchronously and `mobile$` emits when the answer changes; both read the same media query list, so the two can never disagree. Everything else is a written rule plus a checker that fails lint when the rule is broken.

Cosmetic differences (a tooltip, a hint, a checkbox) become CSS classes. Structural differences (a different component, a different behaviour) subscribe to `mobile$`. Table columns become a TypeScript array chosen by mode, and CSS never hides a column. No screen hides a capability on phones; it adapts the interaction instead. Breakpoint values exist at 768 and 769 only. Dialogs take their width from one helper, and the viewport meta allows pinch zoom again.

## User Stories

1. As a phone user, I want every screen to open in the phone layout, so that I never have to rotate or reload to get it.
2. As a phone user, I want the layout to follow a rotation or a resize, so that the screen never keeps a layout that no longer fits.
3. As a phone user, I want to reach every feature a desktop user can reach, so that the phone client is not a reduced client.
4. As a phone user, I want a table to show the columns that matter at 390px, so that rows are readable without scrolling sideways.
5. As a phone user, I want every dialog to fit my screen, so that I can read it and reach its buttons.
6. As a phone user, I want pinch zoom to work, so that I can read text the compact type scale made small.
7. As a phone user, I want the header to stay on one line, so that navigation does not push content off the screen.
8. As a desktop user, I want the desktop layout to look as it does today, so that work done for phones costs me nothing.
9. As a desktop user with a narrow window, I want the same adaptation a phone gets, so that the layout matches the space I gave it.
10. As a developer adding a screen, I want one documented way to adapt it, so that I do not invent a threshold of my own.
11. As a developer, I want the wrong way to fail lint, so that the rule survives the next agent as well as me.
12. As a developer, I want the width signal testable without a browser, so that a table's columns can be asserted at both widths.
13. As a developer, I want the per-component breakpoints gone, so that there is one place to look when a width question comes up.
14. As a reviewer, I want the rule doc named from `frontend/AGENTS.md`, so that it is read before a screen is touched.

## Implementation Decisions

- **One width source.** `Env` holds a single media query list for `(max-width: 768px)`. `isMobile()` returns its `matches` and `mobile$` emits on change, starting with the current value. The list is the source rather than `window.innerWidth` because the two diverge once a scrollbar is present, and two sources would let the synchronous answer and the stream disagree. `isMobile()` keeps its name and stays synchronous, so template branches keep working while a component migrates.
- **Compact stays at 768 or below.** The boundary does not move, so no shipped screen changes width class. The repo defines its own display utilities and does not depend on Bootstrap's breakpoints, whose `md` tier starts at 768 on the wide side of this boundary; Bootstrap's responsive display classes are therefore banned for layout use.
- **Two tier blocks in the global stylesheet.** One compact block at `max-width: 768px` and one wide block at `min-width: 769px`, holding the existing type-scale rules plus two utilities, `.only-compact` and `.only-wide`. A component stylesheet may declare a width media query only at those two values.
- **The checker carries the rule.** A node script run from `npm run lint` reads the width from the service and the two widths from the global stylesheet, fails when they disagree, then scans component styles for any other width and source files for a `window.innerWidth` read outside the service. A rule that lives only in a document is a rule the next agent will break.
- **Columns are TypeScript.** A table keeps a desktop list and a compact list as constants, assigns from `isMobile()` at construction, and reassigns on `mobile$`. No CSS rule targets `.mat-column-*`. A cell hidden in CSS still renders and still holds its width under `table-layout: fixed`, so the CSS route fights the layout the component sets in TypeScript.
- **Templates branch only for structure.** A template may test `env.isMobile()` when it renders a different component or changes behaviour, never to drop an element. Cosmetic differences use the utility classes. No tool can judge that intent, so this rule rests on the rule doc and on review; the checker does not scan templates.
- **Dialogs go through one helper.** It returns `width: min(<px>px, 95vw)` together with `maxWidth: 95vw`, and every dialog open call passes it. Material's 80vw default cap is what clips content today.
- **Viewport meta** keeps `width=device-width, initial-scale=1` and drops the scale pins, so pinch zoom works.
- **The rule doc** is `frontend/moon/doc/agents/responsive.md`, named from `frontend/AGENTS.md` the way the HTTP, i18n and pagination docs are. It holds the decision table, the column snippet, the dialog helper, the banned patterns and the checker's name.
- **Batches.** Foundation first, then tables, then dialogs, then the shell, then the outliers. The foundation blocks the rest; the four migration batches run independently.
- **The image viewer** moves onto the canonical boundary and loses its private 400, 500 and 700 tiers, keeping only fullscreen sizing.
- **Not hidden, adapted.** Where a compact screen cannot show an interactive affordance, the interaction changes shape or route. Dropping the affordance is a product decision and needs its own approval, so it never arrives as a layout change.

## Testing Decisions

A good test here asserts what a screen shows at a given viewport, not which field or listener produced it.

- **Seam (new, one).** A test helper fakes `window.matchMedia`, returning a `MediaQueryList` whose `matches` the test sets and whose listeners the test fires on demand. `Env` is built against it, and every width-dependent assertion in the suite goes through this single seam.
- **Service.** `isMobile()` follows the fake's matches. `mobile$` starts with the current value, emits only when the value changes, and stays silent when a listener fires with an unchanged value.
- **Tables.** For each migrated table, assert the compact column list at 390px, the desktop list at 1280px, and that a mode flip mid-life reassigns the columns. Prior art is `mng-files.component.spec.ts` and `watched-gallery.component.spec.ts`.
- **Checker.** Run it against this repo and expect a zero exit, then against a fixture holding a 600px media query and expect a non-zero exit.
- **Manual check.** The user confirms layout appearance at 390px on their phone per batch. This feature runs no browser automation.

## Out of Scope

- A third tier for tablets, and any layout between 769 and 1024px.
- Landscape-specific layout. A landscape phone is 844px wide and takes the wide layout.
- Touch handling beyond today's state. The tooltip `touchGestures: "off"` workaround stays.
- Redesigning any screen. The batches make existing screens fit, nothing more.
- The desktop type scale itself, beyond moving it into the canonical tier blocks.

## Further Notes

- Rejected a new `LayoutService` beside `Env`. One more name for one question that `Env` already answers, and the injection sites stay as they are.
- Rejected `min(width, height)` and device detection. A narrow desktop window needs what a phone needs, and device identity is not what decides layout.
- Rejected CSS-only column hiding, for the fixed-layout reason above.
- Rejected generating the global stylesheet from the TypeScript constant. It buys one number and adds a build step; the checker catches the same failure for a fraction of the cost.
- `isMobile` keeps its name although it measures viewport width. Renaming forty four call sites buys vocabulary, not behaviour, and the rule doc states what the predicate actually answers.
