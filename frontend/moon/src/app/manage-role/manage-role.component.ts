import { Component, OnDestroy, OnInit, ViewChild } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { animateElementExpanding } from "src/animate/animate-util";
import { MngRoleDialogComponent } from "../mng-role-dialog/mng-role-dialog.component";
import { isEnterKey } from "src/common/condition";
import { dialogSize } from "src/common/dialog-size";
import { HttpClient } from "@angular/common/http";
import { Env } from "src/common/env-util";
import { MatSnackBar } from "@angular/material/snack-bar";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";
import { I18n } from "../i18n.service";

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
export class ManageRoleComponent implements OnInit, OnDestroy {
  private readonly compactColumns = ["name", "updateTime"];
  private readonly wideColumns = [
    "roleNo",
    "name",
    "createBy",
    "createTime",
    "updateBy",
    "updateTime",
  ];
  displayedColumns: string[] = this.wideColumns;
  private readonly modeSub = this.env.mobile$.subscribe((mobile) => {
    this.displayedColumns = mobile ? this.compactColumns : this.wideColumns;
  });

  isEnter = isEnterKey;
  newRoleDialog = false;
  newRoleName = "";
  roles: ERole[] = [];

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  trl(k) {
    return this.i18n.trl("manage-role", k);
  }

  constructor(
    private http: HttpClient,
    private dialog: MatDialog,
    public env: Env,
    private snackBar: MatSnackBar,
    public i18n: I18n
  ) {}

  reset() {
    this.newRoleDialog = false;
    this.pagingController.firstPage();
  }

  ngOnInit(): void {}

  ngOnDestroy(): void {
    this.modeSub.unsubscribe();
  }

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
        ...dialogSize(1000),
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
      this.snackBar.open(this.trl("pleaseEnterNewRoleName"), "ok", {
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
