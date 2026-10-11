# Pagination

Every paged list in the frontend uses `ControlledPaginatorComponent` (`src/app/controlled-paginator/`). It owns page and page-size state and tells you when they change; your component owns the request and the response.

## Wiring a paged list

1. Declare the query `static` — the host reads the paginator during initialization, when a dynamic query would still be `undefined`:

```ts
@ViewChild(ControlledPaginatorComponent, { static: true })
pagingController: ControlledPaginatorComponent;
```

2. Apply the page limit in `ngOnInit`, before anything is fetched:

```ts
ngOnInit(): void {
  this.pagingController.setPageLimitOptions([5, 10, 30, 50, 100]);
}
```

3. Fetch from the paginator's output, not from `ngOnInit`:

```html
<app-controlled-paginator (pageChanged)="listItems()"></app-controlled-paginator>
```

`ngAfterViewInit` emits `pageChanged` once — that emit is the initial fetch — and then `ready`. Pick one of the two and fetch there. A fetch in `ngOnInit` as well sends a second request that races the first, and the loser's response can land last and overwrite the winner.

4. Send the paginator's own `paging` object as the request paging, and report the total back:

```ts
this.listReq.paging = this.pagingController.paging;
// ...
this.pagingController.onTotalChanged(resp.data.paging);
```

## Reference

| Member | Behaviour |
| --- | --- |
| `paging` | `{ page, limit, total }`. Default `limit` is 10. Assign it by reference as the request paging. |
| `setPageLimitOptions(opts, disableMobileOpts = true)` | Sets `paging.limit = opts[0]` — the first option becomes the default page size. |
| `onTotalChanged(paging)` | Stores the total and recomputes the page-number list. Call it from the response handler. |
| `pageChanged` output | Emitted on user navigation, and once from `ngAfterViewInit`. |
| `ready` output | Emitted once from `ngAfterViewInit`, after `pageChanged`. |
| `firstPage()` / `lastPage()` / `syncPage(n)` | Imperative navigation for refresh and search flows. `firstPage()` returns `false` when already on page 1. |

## Two behaviours that surprise people

**Changing `pageSize` or `pageSizeOptions` never emits `pageChanged`.** Material's paginator re-emits only on navigation or on its own page-size selector. A limit applied after the first fetch is displayed but never refetched, so the control shows one page size while the list holds another. Apply the limit in `ngOnInit` and it lands before the initial fetch.

**A dynamic `@ViewChild` is `undefined` while the host is initializing.** Angular materializes a `static` query on the create pass; a dynamic one only after the child's lifecycle has run. So both a synchronous read in `ngOnInit` and a handler triggered by the child's `ngAfterViewInit` emit see `undefined`. Guarding with `this.pagingController ? this.pagingController.paging : { page: 1, limit: 10 }` hides the failure and silently fetches the fallback page size instead. Use the static query; drop the guard.

## Mobile

The breakpoint is **768px, inclusive**: `Env.isMobile()` (`src/common/env-util.ts`) and the CSS tiers in `src/styles.css` share that single value — keep them in step.

On mobile the paginator hides its page-size selector. If the host never calls `setPageLimitOptions`, the paginator applies its own `[5, 10, 30]` list there; a host call replaces that list and disables the mobile fallback. Either way the first option is the page size on every device, so choose it deliberately.
