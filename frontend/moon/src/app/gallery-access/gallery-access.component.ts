import { Component, Inject, OnInit } from "@angular/core";
import { isEnterKey } from "src/common/condition";
import { PagingController } from "src/common/paging";
import { Toaster } from "../notification.service";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { environment } from "src/environments/environment";
import { HttpClient } from "@angular/common/http";

export type GalleryAccessGranted = {
  id: number;
  userNo: string;
  createTime: any;
};

export interface GrantGalleryAccessDialogData {
  galleryNo: string;
  name: string;
}

@Component({
  selector: "app-gallery-access",
  templateUrl: "./gallery-access.component.html",
  styleUrls: ["./gallery-access.component.css"],
})
export class GalleryAccessComponent implements OnInit {
  readonly columns: string[] = ["username", "createTime", "removeButton"];
  grantedTo: string = "";
  grantedAccesses: GalleryAccessGranted[] = [];
  pagingController: PagingController;
  isEnterPressed = isEnterKey;

  constructor(
    private http: HttpClient,
    private toaster: Toaster,
    public dialogRef: MatDialogRef<
      GalleryAccessComponent,
      GrantGalleryAccessDialogData
    >,
    @Inject(MAT_DIALOG_DATA) public data: GrantGalleryAccessDialogData
  ) {}

  ngOnInit() {}

  grantAccess() {
    this.grantFolderAccess();
  }

  grantFolderAccess() {
    if (!this.grantedTo) {
      this.toaster.toast("Enter username first");
      return;
    }

    this.http
      .post<void>(`${environment.vfm}/open/api/gallery/access/grant`, {
        galleryNo: this.data.galleryNo,
        username: this.grantedTo,
      })
      .subscribe({
        next: () => {
          this.toaster.toast("Access granted");
          this.fetchAccessGranted();
        },
      });
  }

  fetchAccessGranted() {
    this.fetchFolderAccessGranted();
  }

  fetchFolderAccessGranted() {
    this.http
      .post<any>(`${environment.vfm}/open/api/gallery/access/list`, {
        galleryNo: this.data.galleryNo,
        paging: this.pagingController.paging,
      })
      .subscribe({
        next: (resp) => {
          this.grantedAccesses = [];
          if (resp.data.payload) {
            for (let g of resp.data.payload) {
              g.createTime = new Date(g.createTime);
              this.grantedAccesses.push(g);
            }
          }
          this.pagingController.onTotalChanged(resp.data.paging);
        },
      });
  }

  removeAccess(access): void {
    this.removeFolderAccess(access.userNo);
  }

  removeFolderAccess(userNo: string): void {
    this.http
      .post<void>(`${environment.vfm}/open/api/gallery/access/remove`, {
        userNo: userNo,
        galleryNo: this.data.galleryNo,
      })
      .subscribe({
        next: () => {
          this.fetchAccessGranted();
        },
      });
  }

  onPagingControllerReady(pc) {
    this.pagingController = pc;
    this.pagingController.onPageChanged = () => this.fetchAccessGranted();
    this.fetchAccessGranted();
  }
}
