import { Component, Inject, OnInit } from "@angular/core";
import {
  MatDialog,
  MatDialogRef,
  MAT_DIALOG_DATA,
} from "@angular/material/dialog";
import { PagingController } from "src/common/paging";
import { environment } from "src/environments/environment";
import { ConfirmDialogComponent } from "../dialog/confirm/confirm-dialog.component";
import { WPath } from "../manage-paths/manage-paths.component";
import { WRes } from "../manage-resources/manage-resources.component";
import { ConfirmDialog } from "src/common/dialog";
import { HttpClient } from "@angular/common/http";
import { MatSnackBar } from "@angular/material/snack-bar";

export interface DialogDat {
  res: WRes;
}

@Component({
  selector: "app-mng-res-dialog",
  templateUrl: "./mng-res-dialog.component.html",
  styleUrls: ["./mng-res-dialog.component.css"],
})
export class MngResDialogComponent implements OnInit {
  readonly tabcol = [
    "id",
    "pgroup",
    "method",
    "url",
    "ptype",
    "desc",
    "option",
  ];
  paths: WPath[] = [];
  pagingController: PagingController = null;

  constructor(
    public dialogRef: MatDialogRef<MngResDialogComponent, DialogDat>,
    @Inject(MAT_DIALOG_DATA) public dat: DialogDat,
    private dialog: MatDialog,
    private confirmDialog: ConfirmDialog,
    private http: HttpClient,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {}

  listPathsBound() {
    this.http
      .post<any>(`user-vault/open/api/path/list`, {
        paging: this.pagingController.paging,
        resCode: this.dat.res.code,
      })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.paths = [];
          if (resp.data && resp.data.payload) {
            for (let ro of resp.data.payload) {
              if (ro.createTime) ro.createTime = new Date(ro.createTime);
              if (ro.updateTime) ro.updateTime = new Date(ro.updateTime);
              this.paths.push(ro);
            }
          }
          this.pagingController.onTotalChanged(resp.data.paging);
        },
      });
  }

  deleteResource() {
    const dialogRef: MatDialogRef<ConfirmDialogComponent, boolean> =
      this.dialog.open(ConfirmDialogComponent, {
        width: "500px",
        data: {
          title: "Remove Resource",
          msg: [`You sure you want to delete resource '${this.dat.res.name}'`],
        },
      });

    dialogRef.afterClosed().subscribe((confirm) => {
      console.log(confirm);
      if (confirm) {
        this.http
          .post<any>(`user-vault/open/api/resource/remove`, {
            resCode: this.dat.res.code,
          })
          .subscribe({
            next: (resp) => {
              if (resp.error) {
                this.snackBar.open(resp.msg, "ok", { duration: 6000 });
                return;
              }
              this.dialogRef.close();
            },
          });
      }
    });
  }

  unbind(pathNo: string, resCode: string, pathUrl: string, resName: string) {
    if (!resCode || !pathNo) {
      return;
    }

    const title = "Unbind Resource Path";
    const msg = [`Resource '${resName}'`, `Path: '${pathUrl}'`];

    this.confirmDialog.show(title, msg, () => {
      this.http
        .post<any>(`user-vault/open/api/path/resource/unbind`, {
          pathNo: pathNo,
          resCode: resCode,
        })
        .subscribe((resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.listPathsBound();
        });
    });
  }

  onPagingControllerReady(pc) {
    this.pagingController = pc;
    this.pagingController.PAGE_LIMIT_OPTIONS = [5];
    this.pagingController.paging.limit = 5;
    this.pagingController.onPageChanged = () => this.listPathsBound();
    this.listPathsBound();
  }
}
