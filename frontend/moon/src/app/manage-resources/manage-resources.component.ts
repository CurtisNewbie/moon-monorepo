import { Component, OnInit, ViewChild } from "@angular/core";
import { getExpanded, isIdEqual } from "src/animate/animate-util";
import { isEnterKey } from "src/common/condition";
import { MngResDialogComponent } from "../mng-res-dialog/mng-res-dialog.component";
import { MatDialog } from "@angular/material/dialog";
import { HttpClient } from "@angular/common/http";
import { Env } from "src/common/env-util";
import { MatSnackBar } from "@angular/material/snack-bar";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";
import { I18n } from "../i18n.service";

export interface WRes {
  id?: number;
  code?: string;
  name?: string;
  createTime?: Date;
  createBy?: string;
  updateTime?: Date;
  updateBy?: string;
}

@Component({
  selector: "app-manage-resources",
  templateUrl: "./manage-resources.component.html",
  styleUrls: ["./manage-resources.component.css"],
})
export class ManageResourcesComponent implements OnInit {
  newResDialog = false;
  newResName = null;
  newResCode = null;
  expandedElement: WRes = null;

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  readonly tabcol = this.env.isMobile()
    ? ["name", "code", "updateTime"]
    : [
        "id",
        "name",
        "code",
        "createBy",
        "createTime",
        "updateBy",
        "updateTime",
      ];
  resources: WRes[] = [];

  idEquals = isIdEqual;
  getExpandedEle = (row) => getExpanded(row, this.expandedElement);
  isEnter = isEnterKey;

  constructor(
    private http: HttpClient,
    private dialog: MatDialog,
    public env: Env,
    private snackBar: MatSnackBar,
    public i18n: I18n
  ) {}

  trl(k) {
    return this.i18n.trl("manage-resources", k);
  }

  reset() {
    this.expandedElement = null;
    this.newResDialog = false;
    this.newResName = null;
    this.newResCode = null;
    this.pagingController.firstPage();
  }

  ngOnInit(): void {}

  fetchList() {
    this.http
      .post<any>(`user-vault/open/api/resource/list`, {
        paging: this.pagingController.paging,
      })
      .subscribe({
        next: (r) => {
          if (r.error) {
            this.snackBar.open(r.msg, "ok", { duration: 6000 });
            return;
          }
          this.resources = [];
          if (r.data && r.data.payload) {
            for (let ro of r.data.payload) {
              if (ro.createTime) ro.createTime = new Date(ro.createTime);
              if (ro.updateTime) ro.updateTime = new Date(ro.updateTime);
              this.resources.push(ro);
            }
          }
          this.pagingController.onTotalChanged(r.data.paging);
        },
      });
  }

  createNewRes() {
    if (!this.newResName) {
      this.snackBar.open(this.trl("pleaseEnterNewResourceName"), "ok", {
        duration: 3000,
      });
      return;
    }
    if (!this.newResCode) {
      this.snackBar.open(this.trl("pleaseEnterNewResourceCode"), "ok", {
        duration: 3000,
      });
      return;
    }

    this.http
      .post<any>(`user-vault/open/api/resource/add`, {
        name: this.newResName,
        code: this.newResCode,
      })
      .subscribe({
        next: (r) => {
          if (r.error) {
            this.snackBar.open(r.msg, "ok", { duration: 6000 });
            return;
          }
          this.newResDialog = false;
          this.newResName = null;
          this.newResCode = null;
          this.fetchList();
        },
      });
  }

  openMngResDialog(r: WRes) {
    this.dialog
      .open(MngResDialogComponent, {
        width: "1000px",
        data: {
          res: { ...r },
        },
      })
      .afterClosed()
      .subscribe({
        complete: () => {
          this.fetchList();
        },
      });
  }
}
