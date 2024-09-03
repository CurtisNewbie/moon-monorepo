import { Component, ElementRef, OnInit, ViewChild } from "@angular/core";
import { PagingController } from "src/common/paging";
import { isEnterKey } from "src/common/condition";
import { environment } from "src/environments/environment";
import { Observable } from "rxjs";
import { HttpClient, HttpEvent } from "@angular/common/http";
import { Toaster } from "../notification.service";
import { ConfirmDialog } from "src/common/dialog";

@Component({
  selector: "app-manage-bookmarks",
  templateUrl: "./manage-bookmarks.component.html",
  styleUrls: ["./manage-bookmarks.component.css"],
})
export class ManageBookmarksComponent implements OnInit {
  readonly isEnterKeyPressed = isEnterKey;
  readonly tabcol = ["id", "name", "operation"];

  pagingController: PagingController;
  tabdat = [];
  isEnter = isEnterKey;
  file = null;

  searchName = null;
  showUploadPanel = false;

  @ViewChild("uploadFileInput")
  uploadFileInput: ElementRef;

  constructor(
    private http: HttpClient,
    private toaster: Toaster,
    private confirmDialog: ConfirmDialog
  ) {}

  ngOnInit(): void {}

  onPagingControllerReady(pc) {
    this.pagingController = pc;
    this.pagingController.onPageChanged = () => this.fetchList();
    this.fetchList();
  }

  fetchList() {
    this.http
      .post<any>(`${environment.vfm}/bookmark/list`, {
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
        this.toaster.toast("Bookmarks uploaded");
        this.showUploadPanel = false;
      },
    });
  }

  onFileSelected(files: File[]) {
    if (files == null || files.length < 1) {
      this.toaster.toast("Please select file");
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
    this.http
      .post<any>(`${environment.vfm}/bookmark/remove`, { id: id })
      .subscribe({
        complete: () => this.fetchList(),
      });
  }

  uploadToTmpFile(file: File): Observable<HttpEvent<any>> {
    return this.http.put<HttpEvent<any>>(
      environment.vfm + "/bookmark/file/upload",
      file,
      {
        observe: "events",
        reportProgress: true,
        withCredentials: true,
      }
    );
  }

  resetSearchName() {
    this.searchName = null;
    this.fetchList();
  }
}
