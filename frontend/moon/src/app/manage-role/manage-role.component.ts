import { Component, OnInit, ViewChild } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { animateElementExpanding } from "src/animate/animate-util";
import { MngRoleDialogComponent } from "../mng-role-dialog/mng-role-dialog.component";
import { isEnterKey } from "src/common/condition";
import { HttpClient } from "@angular/common/http";
import { Env } from "src/common/env-util";
import { MatSnackBar } from "@angular/material/snack-bar";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";

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
  isEnter = isEnterKey;
  newRoleDialog = false;
  newRoleName = "";
  readonly tabcol = this.env.isMobile()
    ? ["roleNo", "name", "updateTime"]
    : ["roleNo", "name", "createBy", "createTime", "updateBy", "updateTime"];
  roles: ERole[] = [];

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  constructor(
    private http: HttpClient,
    private dialog: MatDialog,
    public env: Env,
    private snackBar: MatSnackBar
  ) {}

  reset() {
    this.newRoleDialog = false;
    this.pagingController.firstPage();
  }

  ngOnInit(): void {}

  fetchList() {
    this.http
      .post<any>(`user-vault/open/api/role/list`, {
        paging: this.pagingController.paging,
      })
      .subscribe({
        next: (r) => {
          if (r.error) {
            this.snackBar.open(r.msg, "ok", { duration: 6000 });
            return;
          }
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

  openMngRoleDialog(role: ERole) {
    this.dialog
      .open(MngRoleDialogComponent, {
        width: "1000px",
        data: {
          roleName: role.name,
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
      this.snackBar.open("Please enter new role name", "ok", {
        duration: 3000,
      });
      return;
    }

    this.http
      .post<any>(`user-vault/open/api/role/add`, {
        name: this.newRoleName,
      })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.newRoleDialog = false;
          this.newRoleName = null;
          this.fetchList();
        },
      });
  }
}
