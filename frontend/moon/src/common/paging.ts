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

  paging: Paging = {
    page: 1,
    limit: 10,
    total: 0,
  };
  maxPage: number = 1;

  private paginator: MatPaginator = null;

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

  public nextPage(): boolean {
    let b = this.paginator.hasNextPage();
    this.paginator.nextPage();
    return b;
  }

  public prevPage(): boolean {
    let b = this.paginator.hasPreviousPage();
    this.paginator.previousPage();
    return b;
  }

  /** is at first page */
  public atFirstPage(): boolean {
    return this.paginator.pageIndex == 0;
  }

  /** set the paginator controlled by this controller */
  public control(paginator: MatPaginator) {
    this.paginator = paginator;
  }

  public onTotalChanged(p: Paging): void {
    this._updatePages(p.total);
  }

  private _updatePages(total: number): void {
    this.paging.total = total;
    this.maxPage = Math.ceil(total / this.paging.limit);
  }

  /** set page limit */
  public setPageLimit(limit: number): void {
    this.paging.limit = limit;
  }

  public onPageEvent(e: PageEvent): void {
    this.paging.page = e.pageIndex + 1;
    this.paging.limit = e.pageSize;
  }
}
