import { Component, OnInit, ViewChild } from "@angular/core";
import { isEnterKey } from "src/common/condition";
import { ConfirmDialog } from "src/common/dialog";
import { HttpClient } from "@angular/common/http";
import { MatSnackBar } from "@angular/material/snack-bar";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";

@Component({
  selector: "app-bookmark-blacklist",
  templateUrl: "./bookmark-blacklist.component.html",
  styleUrls: ["./bookmark-blacklist.component.css"],
})
export class BookmarkBlacklistComponent implements OnInit {
  readonly isEnterKeyPressed = isEnterKey;
  readonly tabcol = ["id", "name", "operation"];

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  tabdat = [];
  isEnter = isEnterKey;
  file = null;

  searchName = null;
  showUploadPanel = false;

  constructor(
    private http: HttpClient,
    private confirmDialog: ConfirmDialog,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {}

  fetchList() {
    this.http
      .post<any>(`vfm/bookmark/blacklist/list`, {
        paging: this.pagingController.paging,
        name: this.searchName,
      })
      .subscribe({
        next: (r) => {
          if (r.error) {
            this.snackBar.open(r.msg, "ok", { duration: 6000 });
            return;
          }

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
    this.http.post<any>(`vfm/bookmark/blacklist/remove`, { id: id }).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
      },
      complete: () => this.fetchList(),
    });
  }

  resetSearchName() {
    this.searchName = null;
    this.fetchList();
  }
}
