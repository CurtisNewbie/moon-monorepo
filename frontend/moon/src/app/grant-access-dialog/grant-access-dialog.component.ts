import { Component, Inject, OnInit } from "@angular/core";
import { MatDialogRef, MAT_DIALOG_DATA } from "@angular/material/dialog";
import { environment } from "src/environments/environment";
import { FileAccessGranted } from "src/common/file-info";
import { PagingController } from "src/common/paging";
import { Toaster } from "../notification.service";
import { isEnterKey } from "src/common/condition";
import { HttpClient } from "@angular/common/http";

export interface GrantAccessDialogData {
  folderNo?: string;
  name: string;
}

@Component({
  selector: "app-grant-access-dialog",
  templateUrl: "./grant-access-dialog.component.html",
  styleUrls: ["./grant-access-dialog.component.css"],
})
export class GrantAccessDialogComponent implements OnInit {
  readonly columns: string[] = ["username", "createDate", "removeButton"];
  grantedTo: string = "";
  grantedAccesses: FileAccessGranted[] = [];
  pagingController: PagingController;
  isEnterPressed = isEnterKey;

  constructor(
    private http: HttpClient,
    private toaster: Toaster,
    public dialogRef: MatDialogRef<
      GrantAccessDialogComponent,
      GrantAccessDialogData
    >,
    @Inject(MAT_DIALOG_DATA) public data: GrantAccessDialogData
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
      .post<void>(`${environment.vfm}/open/api/vfolder/share`, {
        folderNo: this.data.folderNo,
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
      .post<any>(`${environment.vfm}/open/api/vfolder/granted/list`, {
        folderNo: this.data.folderNo,
        paging: this.pagingController.paging,
      })
      .subscribe({
        next: (resp) => {
          this.grantedAccesses = [];
          if (resp.data.payload) {
            for (let g of resp.data.payload) {
              g.createDate = new Date(g.createTime);
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
      .post<void>(`${environment.vfm}/open/api/vfolder/access/remove`, {
        userNo: userNo,
        folderNo: this.data.folderNo,
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
