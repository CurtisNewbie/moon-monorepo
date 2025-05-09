import { Component, ElementRef, OnInit, ViewChild } from "@angular/core";
import { isEnterKey } from "src/common/condition";
import { Observable } from "rxjs";
import { HttpClient, HttpEvent } from "@angular/common/http";
import { ConfirmDialog } from "src/common/dialog";
import { MatSnackBar } from "@angular/material/snack-bar";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";

@Component({
  selector: "app-manage-bookmarks",
  templateUrl: "./manage-bookmarks.component.html",
  styleUrls: ["./manage-bookmarks.component.css"],
})
export class ManageBookmarksComponent implements OnInit {
  readonly isEnterKeyPressed = isEnterKey;
  readonly tabcol = ["id", "name", "operation"];

  tabdat = [];
  isEnter = isEnterKey;
  file = null;

  searchName = null;
  showUploadPanel = false;

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  @ViewChild("uploadFileInput")
  uploadFileInput: ElementRef;

  constructor(
    private http: HttpClient,
    private confirmDialog: ConfirmDialog,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {}

  fetchList() {
    this.http
      .post<any>(`vfm/bookmark/list`, {
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

  upload() {
    if (!this.file) {
      return null;
    }
    this.uploadToTmpFile(this.file).subscribe({
      complete: () => {
        this.file = null;
        this.fetchList();
        if (this.uploadFileInput) {
          this.uploadFileInput.nativeElement.value = null;
        }
        this.snackBar.open("Bookmarks uploaded", "ok", { duration: 3000 });
        this.showUploadPanel = false;
      },
    });
  }

  onFileSelected(files: File[]) {
    if (files == null || files.length < 1) {
      this.snackBar.open("Please select file", "ok", { duration: 3000 });
      return;
    }
    this.file = files[0];
  }

  popToRemove(id, name) {
    this.confirmDialog.show(
      "Remove Bookmark",
      [`Removing Bookmark ${name}`],
      () => {
        this.remove(id);
      }
    );
  }

  remove(id) {
    this.http.post<any>(`vfm/bookmark/remove`, { id: id }).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
      },
      complete: () => this.fetchList(),
    });
  }

  uploadToTmpFile(file: File): Observable<HttpEvent<any>> {
    return this.http.put<HttpEvent<any>>("vfm/bookmark/file/upload", file, {
      observe: "events",
      reportProgress: true,
      withCredentials: true,
    });
  }

  resetSearchName() {
    this.searchName = null;
    this.fetchList();
  }
}
