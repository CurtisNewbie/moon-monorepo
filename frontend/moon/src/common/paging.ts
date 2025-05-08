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
