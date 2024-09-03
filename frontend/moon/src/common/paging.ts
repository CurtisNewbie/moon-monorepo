import { MatPaginator, PageEvent } from "@angular/material/paginator";

/** Pagination info */
export interface Paging {
  /** page number */
  page: number;
  /** page size */
  limit: number;
  /** total number of items */
  total: number;
}

/**
 * Constants related paging
 */
export class PagingConst {
  /** Get default paging limit options */
  public static getPagingLimitOptions(): number[] {
    return [10, 30, 50, 100, 500];
  }
}

/**
 * Controller for pagination, internal properties are non-private, thus can be directly bound with directive
 */
export class PagingController {

  PAGE_LIMIT_OPTIONS: number[] = PagingConst.getPagingLimitOptions();
  paging: Paging = {
    page: 1,
    limit: this.PAGE_LIMIT_OPTIONS[0],
    total: 0,
  };
  pages: number[] = [1];

  private paginator: MatPaginator = null;

  /** callback invoked when current page is changed */
  onPageChanged: () => void = null;

  /** go to last page */
  public lastPage() {
    this.paginator.lastPage();
  }

  /** go to first page */
  public firstPage(): boolean {
    if (this.atFirstPage()) {
      return false;
    }
    this.paginator.firstPage();
    return true;
  }

  /** is at first page */
  public atFirstPage(): boolean {
    return this.paginator.pageIndex == 0;
  }

  /** set the paginator controlled by this controller */
  public control(paginator: MatPaginator) {
    this.paginator = paginator;
    if (paginator) {
      paginator.page.subscribe((e) => this.onPageEvent(e));
    }
  }

  /** update the list of pages that it can select based on total */
  public onTotalChanged(p: Paging): void {
    this._updatePages(p.total);
  }

  /** update the list of pages that it can select based on total */
  private _updatePages(total: number): void {
    this.pages = [];
    this.paging.total = total;
    let maxPage = Math.ceil(total / this.paging.limit);
    for (let i = 1; i <= maxPage; i++) {
      this.pages.push(i);
    }
    if (this.pages.length === 0) {
      this.pages.push(1);
    }
  }

  /** set page limit */
  public setPageLimit(limit: number): void {
    this.paging.limit = limit;
  }

  public onPageEvent(e: PageEvent): void {
    // console.log(e);
    this.paging.page = e.pageIndex + 1;
    this.paging.limit = e.pageSize;
    if (this.onPageChanged) this.onPageChanged();
  }
}
