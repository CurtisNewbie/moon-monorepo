import { Component, OnInit } from "@angular/core";
import { isEnterKey } from "src/common/condition";
import { PagingController } from "src/common/paging";
import { ConfirmDialog } from "src/common/dialog";
import { environment } from "src/environments/environment";
import { HttpClient } from "@angular/common/http";

@Component({
  selector: "app-bookmark-blacklist",
  templateUrl: "./bookmark-blacklist.component.html",
  styleUrls: ["./bookmark-blacklist.component.css"],
})
export class BookmarkBlacklistComponent implements OnInit {
  readonly isEnterKeyPressed = isEnterKey;
  readonly tabcol = ["id", "name", "operation"];

  pagingController: PagingController;
  tabdat = [];
  isEnter = isEnterKey;
  file = null;

  searchName = null;
  showUploadPanel = false;

  constructor(private http: HttpClient, private confirmDialog: ConfirmDialog) {}

  ngOnInit(): void {}

  onPagingControllerReady(pc) {
    this.pagingController = pc;
    this.pagingController.onPageChanged = () => this.fetchList();
    this.fetchList();
  }

  fetchList() {
    this.http
      .post<any>(`${environment.vfm}/bookmark/blacklist/list`, {
        paging: this.pagingController.paging,
        name: this.searchName,
      })
      .subscribe({
        next: (r) => {
          this.tabdat = r.data.payload;
          this.pagingController.onTotalChanged(r.data.paging);
        },
      });
  }

  popToRemove(id, name) {
    this.confirmDialog.show(
      "Remove Bookmark Blacklist",
      [`Removing Bookmark Blacklist ${name}`],
      () => {
        this.remove(id);
      }
    );
  }

  remove(id) {
    this.http
      .post<any>(`${environment.vfm}/bookmark/blacklist/remove`, { id: id })
      .subscribe({
        complete: () => this.fetchList(),
      });
  }

  resetSearchName() {
    this.searchName = null;
    this.fetchList();
  }
}
