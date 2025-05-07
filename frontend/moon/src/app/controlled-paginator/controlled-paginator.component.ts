import {
  Component,
  EventEmitter,
  OnInit,
  Output,
  ViewChild,
} from "@angular/core";
import { MatPaginator, PageEvent } from "@angular/material/paginator";
import { isEnterKey } from "src/common/condition";
import { PagingController } from "src/common/paging";

@Component({
  selector: "app-controlled-paginator",
  templateUrl: "./controlled-paginator.component.html",
  styleUrls: ["./controlled-paginator.component.css"],
})
export class ControlledPaginatorComponent implements OnInit {
  @ViewChild("paginator", { static: true })
  paginator: MatPaginator;

  goto: string;
  maxPage: number;

  @Output("controllerReady")
  controllerEmitter = new EventEmitter<PagingController>();

  pagingController = new PagingController();

  constructor() {}

  ngOnInit(): void {
    this.pagingController.control(this.paginator);
    this.controllerEmitter.emit(this.pagingController);
    this.goto = String(1);

    this.paginator.page.subscribe((evt) => {
      this.goto = String(evt.pageIndex + 1);
      this.pagingController.onPageEvent(evt);
    });
    this.maxPage = 1;
  }

  goToPage(evt) {
    if (!this.goto) {
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

    let maxPage = this.pagingController.maxPage;
    if (n > maxPage) {
      n = maxPage;
    }
    this.goto = String(n);

    if (!isEnterKey(evt)) {
      return;
    }

    this.paginator.pageIndex = n - 1;
    const event: PageEvent = {
      length: this.paginator.length,
      pageIndex: this.paginator.pageIndex,
      pageSize: this.paginator.pageSize,
    };
    this.paginator.page.next(event);
  }
}
