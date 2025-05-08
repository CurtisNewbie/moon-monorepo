import { Component, Inject, OnInit, ViewChild } from "@angular/core";
import { MatDialogRef, MAT_DIALOG_DATA } from "@angular/material/dialog";
import { ResBrief } from "../user.service";
import { HttpClient } from "@angular/common/http";
import { ConfirmDialog } from "src/common/dialog";
import { MatSnackBar } from "@angular/material/snack-bar";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";

export interface DialogDat {
  roleName: string;
  roleNo: string;
}

export interface ListedRoleRes {
  id?: number;
  resCode?: string;
  resName?: string;
  createTime?: Date;
  createBy?: string;
}

@Component({
  selector: "app-mng-role-dialog",
  templateUrl: "./mng-role-dialog.component.html",
  styleUrls: ["./mng-role-dialog.component.css"],
})
export class MngRoleDialogComponent implements OnInit {
  readonly tabcol = [
    "id",
    "resCode",
    "resName",
    "createTime",
    "createBy",
    "operation",
  ];

  roleRes: ListedRoleRes[] = [];
  resBriefs: ResBrief[] = [];
  addResCode: string = null;

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  constructor(
    public dialogRef: MatDialogRef<MngRoleDialogComponent, DialogDat>,
    @Inject(MAT_DIALOG_DATA) public dat: DialogDat,
    private http: HttpClient,
    private confirmDialog: ConfirmDialog,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.fetchResourceCandidates();
  }

  fetchResourceCandidates() {
    this.http
      .get<any>(
        `user-vault/open/api/resource/brief/candidates?roleNo=${this.dat.roleNo}`
      )
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.resBriefs = resp.data;
        },
      });
  }

  addResource() {
    if (!this.addResCode) {
      this.snackBar.open("Please select resource to add", "ok", {
        duration: 3000,
      });
      return;
    }

    this.http
      .post<any>(`user-vault/open/api/role/resource/add`, {
        roleNo: this.dat.roleNo,
        resCode: this.addResCode,
      })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.addResCode = null;
          this.listResources();
          this.fetchResourceCandidates();
        },
      });
  }

  listResources() {
    this.http
      .post<any>(`user-vault/open/api/role/resource/list`, {
        roleNo: this.dat.roleNo,
        paging: this.pagingController.paging,
      })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.roleRes = [];
          if (resp.data && resp.data.payload) {
            for (let r of resp.data.payload) {
              if (r.createTime) r.createTime = new Date(r.createTime);
              this.roleRes.push(r);
            }
            this.pagingController.onTotalChanged(resp.data.paging);
          }
        },
      });
  }

  delRes(roleRes: ListedRoleRes) {
    this.confirmDialog.show(
      "Unbind Resource",
      [
        `Confirm to unbind resource '${roleRes.resCode}' from role '${this.dat.roleName}'?`,
      ],
      () => {
        this.http
          .post<any>(`user-vault/open/api/role/resource/remove`, {
            roleNo: this.dat.roleNo,
            resCode: roleRes.resCode,
          })
          .subscribe({
            next: (resp) => {
              if (resp.error) {
                this.snackBar.open(resp.msg, "ok", { duration: 6000 });
                return;
              }
              this.listResources();
              this.fetchResourceCandidates();
            },
          });
      }
    );
  }
}
