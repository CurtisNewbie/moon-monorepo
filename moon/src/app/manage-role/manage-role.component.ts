import { Component, OnInit } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { animateElementExpanding } from "src/animate/animate-util";
import { environment } from "src/environments/environment";
import { PagingController } from "src/common/paging";
import { MngRoleDialogComponent } from "../mng-role-dialog/mng-role-dialog.component";
import { Toaster } from "../notification.service";
import { isEnterKey } from "src/common/condition";
import { HttpClient } from "@angular/common/http";

export interface ERole {
  id?: number;
  roleNo?: String;
  name?: String;
  createTime?: Date;
  createBy?: String;
  updateTime?: Date;
  updateBy?: String;
}

@Component({
  selector: "app-manage-role",
  templateUrl: "./manage-role.component.html",
  styleUrls: ["./manage-role.component.css"],
  animations: [animateElementExpanding()],
})
export class ManageRoleComponent implements OnInit {
  newRoleDialog = false;
  newRoleName = "";
  pagingController: PagingController;

  readonly tabcol = [
    "id",
    "name",
    "roleNo",
    "createBy",
    "createTime",
    "updateBy",
    "updateTime",
  ];
  roles: ERole[] = [];

  isEnter = isEnterKey;

  constructor(
    private http: HttpClient,
    private dialog: MatDialog,
    private toaster: Toaster
  ) {}

  reset() {
    this.newRoleDialog = false;
    this.pagingController.firstPage();
  }

  ngOnInit(): void {}

  fetchList() {
    this.http
      .post<any>(`${environment.uservault}/open/api/role/list`, {
        paging: this.pagingController.paging,
      })
      .subscribe({
        next: (r) => {
          this.roles = [];
          if (r.data && r.data.payload) {
            for (let ro of r.data.payload) {
              if (ro.createTime) ro.createTime = new Date(ro.createTime);
              if (ro.updateTime) ro.updateTime = new Date(ro.updateTime);
              this.roles.push(ro);
            }
          }
          this.pagingController.onTotalChanged(r.data.paging);
        },
      });
  }

  onPagingControllerReady(pc) {
    this.pagingController = pc;
    this.pagingController.onPageChanged = () => this.fetchList();
    this.fetchList();
  }

  openMngRoleDialog(role: ERole) {
    this.dialog
      .open(MngRoleDialogComponent, {
        width: "1000px",
        data: {
          roleNo: role.roleNo,
        },
      })
      .afterClosed()
      .subscribe({
        complete: () => {
          this.fetchList();
        },
      });
  }

  createNewRole() {
    if (!this.newRoleName) {
      this.toaster.toast("Please enter new role name");
      return;
    }

    this.http
      .post(`${environment.uservault}/open/api/role/add`, {
        name: this.newRoleName,
      })
      .subscribe({
        next: (res) => {
          this.newRoleDialog = false;
          this.newRoleName = null;
          this.fetchList();
        },
      });
  }
}
