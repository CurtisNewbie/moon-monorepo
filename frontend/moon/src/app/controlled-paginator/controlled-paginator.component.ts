import {
  AfterViewInit,
  Component,
  EventEmitter,
  OnInit,
  Output,
  ViewChild,
} from "@angular/core";
import { MatPaginator, PageEvent } from "@angular/material/paginator";
import { isEnterKey } from "src/common/condition";
import { Env } from "src/common/env-util";
import { Paging, PagingConst } from "src/common/paging";

@Component({
  selector: "app-controlled-paginator",
  templateUrl: "./controlled-paginator.component.html",
  styleUrls: ["./controlled-paginator.component.css"],
})
export class ControlledPaginatorComponent implements OnInit, AfterViewInit {
  PAGE_LIMIT_OPTIONS: number[] = PagingConst.getPagingLimitOptions();
  paging: Paging = {
    page: 1,
    limit: 10,
    total: 0,
  };

  @ViewChild("paginator", { static: true })
  paginator: MatPaginator;

  goto: string = "1";
  maxPage: number = 1;

  @Output("pageChanged")
  pageChangedEmitter = new EventEmitter<Paging>();

  constructor(public env: Env) {}

  ngAfterViewInit(): void {
    // first page
    this.pageChangedEmitter.emit(this.paging);
  }

  ngOnInit(): void {
    this.paginator.page.subscribe((evt) => {
      this.paging.page = evt.pageIndex + 1;
      this.paging.limit = evt.pageSize;
      this.goto = String(evt.pageIndex + 1);
      this.pageChangedEmitter.emit(this.paging);
    });

    if (this.env.isMobile()) {
      this.PAGE_LIMIT_OPTIONS = [5, 10, 30];
      this.paging.limit = this.PAGE_LIMIT_OPTIONS[0];
    }
  }

  onGoToPageKeyUp(evt) {
    if (!this.goto) {
      this.goto = String(1);
      return;
    }
    let n = parseInt(this.goto);
    if (Number.isNaN(n)) {
      this.goto = "";
      return;
    }
    if (n < 1) {
      this.goto = "1";
      n = 1;
    }

    let maxPage = this.maxPage;
    if (n > maxPage) {
      n = maxPage;
    }
    this.goto = String(n);

    if (!isEnterKey(evt)) {
      return;
    }
    this.goToPage(n);
  }

  goToPage(n) {
    this.paginator.pageIndex = n - 1;
    this.emitPageEvent();
  }

  emitPageEvent() {
    const event: PageEvent = {
      length: this.paginator.length,
      pageIndex: this.paginator.pageIndex,
      pageSize: this.paginator.pageSize,
    };
    this.paginator.page.next(event);
  }

  /** go to last page */
  lastPage() {
    this.paginator.lastPage();
  }

  /** go to first page */
  firstPage(): boolean {
    if (this.atFirstPage()) {
      return false;
    }
    this.paginator.firstPage();
    return true;
  }

  nextPage(): boolean {
    let b = this.paginator.hasNextPage();
    this.paginator.nextPage();
    return b;
  }

  prevPage(): boolean {
    let b = this.paginator.hasPreviousPage();
    this.paginator.previousPage();
    return b;
  }

  /** is at first page */
  atFirstPage(): boolean {
    return this.paginator.pageIndex == 0;
  }

  /** set the paginator controlled by this controller */
  control(paginator: MatPaginator) {
    this.paginator = paginator;
  }

  onTotalChanged(p: Paging): void {
    this._updatePages(p.total);
  }

  private _updatePages(total: number): void {
    this.paging.total = total;
    this.maxPage = Math.ceil(total / this.paging.limit);
  }

  /** set page limit */
  setPageLimit(limit: number): void {
    this.paging.limit = limit;
  }
}
