import { Component, Inject, OnInit, ViewChild } from "@angular/core";
import { isEnterKey } from "src/common/condition";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { HttpClient } from "@angular/common/http";
import { MatSnackBar } from "@angular/material/snack-bar";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";

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
  isEnterPressed = isEnterKey;

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  constructor(
    private http: HttpClient,
    private snackBar: MatSnackBar,
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
      this.snackBar.open("Enter username first", "ok", { duration: 3000 });
      return;
    }

    this.http
      .post<any>(`vfm/open/api/gallery/access/grant`, {
        galleryNo: this.data.galleryNo,
        username: this.grantedTo,
      })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.snackBar.open("Access granted", "ok", { duration: 3000 });
          this.fetchFolderAccessGranted();
        },
      });
  }

  fetchFolderAccessGranted() {
    this.http
      .post<any>(`vfm/open/api/gallery/access/list`, {
        galleryNo: this.data.galleryNo,
        paging: this.pagingController.paging,
      })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
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
      .post<any>(`vfm/open/api/gallery/access/remove`, {
        userNo: userNo,
        galleryNo: this.data.galleryNo,
      })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.fetchFolderAccessGranted();
        },
      });
  }

}
