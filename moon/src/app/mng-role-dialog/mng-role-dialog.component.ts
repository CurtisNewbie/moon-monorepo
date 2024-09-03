import { Component, Inject, OnInit } from "@angular/core";
import { MatDialogRef, MAT_DIALOG_DATA } from "@angular/material/dialog";
import { environment } from "src/environments/environment";
import { PagingController } from "src/common/paging";
import { Toaster } from "../notification.service";
import { ResBrief } from "../user.service";
import { HttpClient } from "@angular/common/http";

export interface DialogDat {
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
  pagingController: PagingController = null;
  roleRes: ListedRoleRes[] = [];
  resBriefs: ResBrief[] = [];
  addResCode: string = null;

  constructor(
    public dialogRef: MatDialogRef<MngRoleDialogComponent, DialogDat>,
    @Inject(MAT_DIALOG_DATA) public dat: DialogDat,
    private http: HttpClient,
    private toaster: Toaster
  ) {}

  ngOnInit(): void {
    this.fetchResourceCandidates();
  }

  fetchResourceCandidates() {
    this.http
      .get<any>(
        `${environment.uservault}/open/api/resource/brief/candidates?roleNo=${this.dat.roleNo}`
      )
      .subscribe({
        next: (res) => {
          this.resBriefs = res.data;
        },
      });
  }

  addResource() {
    if (!this.addResCode) {
      this.toaster.toast("Please select resource to add");
      return;
    }

    this.http
      .post<any>(`${environment.uservault}/open/api/role/resource/add`, {
        roleNo: this.dat.roleNo,
        resCode: this.addResCode,
      })
      .subscribe({
        next: (res) => {
          this.addResCode = null;
          this.listResources();
          this.fetchResourceCandidates();
        },
      });
  }

  listResources() {
    this.http
      .post<any>(`${environment.uservault}/open/api/role/resource/list`, {
        roleNo: this.dat.roleNo,
        paging: this.pagingController.paging,
      })
      .subscribe({
        next: (res) => {
          this.roleRes = [];
          if (res.data && res.data.payload) {
            for (let r of res.data.payload) {
              if (r.createTime) r.createTime = new Date(r.createTime);
              this.roleRes.push(r);
            }
            this.pagingController.onTotalChanged(res.data.paging);
          }
        },
      });
  }

  onPagingControllerReady(pc) {
    this.pagingController = pc;
    this.pagingController.onPageChanged = () => this.listResources();
    this.listResources();
  }

  delRes(roleRes: ListedRoleRes) {
    this.http
      .post<any>(`${environment.uservault}/open/api/role/resource/remove`, {
        roleNo: this.dat.roleNo,
        resCode: roleRes.resCode,
      })
      .subscribe({
        next: (res) => {
          this.listResources();
          this.fetchResourceCandidates();
        },
      });
  }
}
