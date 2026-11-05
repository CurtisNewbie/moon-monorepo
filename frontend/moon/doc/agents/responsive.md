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

A screen adapts its interaction on phones. When an action leaves the compact layout, that is a product decision and this doc records it.

The gallery's row click is the worked example. On a wide screen the row expands in place and holds the name input with the update, delete, share, directory and browse actions. On a compact screen the same tap browses into the gallery, and the row's kebab menu carries browse and directory. Update, delete and share are desktop-only by product decision: their work starts from the widened row, so a phone keeps the browse shortcut and leaves owner administration to the wide layout. A template branch that renders a desktop-only control needs that decision behind it, and one that drops an action without it is a bug. When a swap changes the route, check the compact route against the wide one and write down any action it no longer reaches.

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

Why not CSS: the CSS list drifts from the TypeScript one, and a column dropped with `display: none` still renders its cells. That drift is what this rule removes.

### Column widths

Angular Material tags both the header cell and the body cell of a column with
`mat-column-<matColumnDef>`. Size the column once, in the component stylesheet:

```css
.mat-column-notifiNo {
  width: 220px;
  max-width: 220px;
}
```

`max-width` is not optional for anything wider than 200px. `src/styles.css` caps
every cell for wrapping:

```css
td,
th {
  white-space: normal;
  word-wrap: break-word;
  max-width: 200px;
}
```

A column that asks for more keeps rendering at exactly 200px, however the width
was declared — inline on the cell or in the stylesheet — and devtools shows
`width: 200px; max-width: 200px` on the header. Put both properties on the
`.mat-column-*` rule: the class beats the `td, th` element selector on
specificity, and one rule reaches the header and the body cells.

- The width includes cell padding. The first cell of a row has 24px left padding
  and the last cell 24px right padding (Material's own rules; some components
  override `.mat-cell`), so usable content is `width - padding`.
- Size to the widest value the column holds, not to the header label. A long
  unbreakable ID wraps mid-token instead of widening the column, so a width a
  few pixels short shows up as one stray character on a second line.
- Under 200px, a `max-width` alone is enough; that is what the other tables do
  (`.mat-column-name { max-width: 150px }`). Above 200px, both are required.

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
