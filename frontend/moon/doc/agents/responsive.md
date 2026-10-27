# Responsive Layout

One width signal decides the layout: `Env` (`src/common/env-util.ts`). It owns a single media query list for the boundary, and answers two ways from that one list:

- `env.isMobile()` — synchronous, returns the current match. Read it, never store it.
- `env.mobile$` — starts with the current value, emits when and only when the answer changes.

The compact layout is **768px and below**; the wide layout is **769px and above**. There is no third tier. `src/styles.css` holds the only two width media queries in the app, and `scripts/check-breakpoints.js` fails lint when anything drifts from that.

## Which mechanism

| Difference | Mechanism |
| --- | --- |
| Cosmetic — a hint, a tooltip, padding, icon or type size, an element that looks different | CSS, inside the compact or wide block of `src/styles.css`, or in a component stylesheet at one of the two widths. Use `.only-compact` / `.only-wide` when the element exists at one tier only. |
| Structural — a different component, a different route, a different behaviour | Subscribe to `env.mobile$` in TypeScript and reassign the state. Tear the subscription down in `ngOnDestroy`. |
| Table columns | Two TypeScript arrays on the component, assigned at construction and reassigned from `mobile$`. Never CSS. |
| A template that renders a different component or changes behaviour | `env.isMobile()` in the template is allowed. A template branch that merely drops an element is not. |

A screen does not hide a capability on phones; it adapts the interaction.

An allowed swap changes the shape of the interaction, not the set of things it can reach. The gallery's row click is the worked example: on a wide screen the row expands in place to show the gallery's actions, and on a compact screen the same tap browses into the gallery while the row's kebab menu carries the same Update, Delete, Share and Directory actions (`src/app/gallery/gallery.component.html`). The row is not inert on a phone, and no action is dropped; it takes a different route. A template branch that renders a desktop-only button is not a swap: that is a hidden action, and it needs a product decision rather than a layout one. When a swap changes the route, check that the compact route still reaches every action the wide route offers before calling it one.

## Table columns

Keep the wide list and the compact list as constants, choose at construction, and follow the signal:

```ts
private readonly compactColumns = ['name', 'operation'];
private readonly wideColumns = ['name', 'url', 'createdAt', 'operation'];
displayedColumns: string[] = this.wideColumns;
private readonly modeSub = this.env.mobile$.subscribe((mobile) => {
  this.displayedColumns = mobile ? this.compactColumns : this.wideColumns;
});

ngOnDestroy(): void {
  this.modeSub.unsubscribe();
}
```

The subscription's first emission carries the current value and sets `displayedColumns` before the view renders, so do not assign it from `isMobile()` in `ngOnInit` as well. A mid-life rotation or resize reassigns it with no reload.

Why not CSS: a cell hidden with `display: none` still renders and still holds its width under `table-layout: fixed`, and the CSS list drifts from the TypeScript one. That drift is what this rule removes.

## Dialogs

Size every dialog through `dialogSize` (`src/common/dialog-size.ts`):

```ts
import { dialogSize } from 'src/common/dialog-size';

this.dialog.open(SomeDialog, { ...dialogSize(800), data });
```

It returns `width: min(800px, 95vw)` and `maxWidth: 95vw`. Material's 80vw default cap is what clips content on a phone; passing the helper keeps the designed width on desktop and fits the screen on a phone.

A dialog designed to fill part of the viewport rather than a fixed pixel width takes `dialogSizeVw` from the same module: it takes a percentage, sets that value as the width and applies the same 95vw maximum, so `dialogSizeVw(80)` reads as "80% of the viewport, never edge to edge".

## What is banned

- `window.innerWidth` or `document.documentElement.clientWidth` anywhere outside `src/common/env-util.ts`. Ask `Env`.
- Caching `env.isMobile()` in a field. A field survives the resize that should have changed it.
- `display: none` on `.mat-column-*`, `.mat-cell` or `.mat-header-cell` to drop a table column. Choose the column list in TypeScript.
- Bootstrap's responsive display classes (`d-none d-md-block`, `d-sm-flex`, …) for layout. Bootstrap's `md` tier does not share this boundary.
- Any width in a media query other than `768px` or `769px`. Component stylesheets use the same two values, with no private tiers.

## Utilities

`.only-compact` shows at 768px and below and is hidden above; `.only-wide` shows at 769px and above and is hidden below. Both use `!important` so they beat Material's own display rules. They are cosmetic only: reaching for either to remove a feature is the banned pattern, not the rule.

## The checker

```bash
npm run lint:breakpoints   # node scripts/check-breakpoints.js
```

It runs from `npm run lint`, before `ng lint`. It fails when the boundary in `env-util.ts` and the two blocks in `src/styles.css` disagree, when any stylesheet declares a width other than 768/769, or when a `.ts`/`.html` file reads the viewport width outside `Env` — naming each violation by file and line. There is no allow-list: every file answers to the rule.

The checker does not read templates for intent. A template branch on `env.isMobile()` that drops an element instead of changing structure passes the checker and rests on review.

## Viewport

`src/index.html` keeps `width=device-width, initial-scale=1` and no scale pins, so pinch zoom works. Do not add `minimum-scale` or `maximum-scale` back.
