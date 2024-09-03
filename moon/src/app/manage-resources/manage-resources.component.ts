import { Component, OnInit } from "@angular/core";
import { getExpanded, isIdEqual } from "src/animate/animate-util";
import { environment } from "src/environments/environment";
import { PagingController } from "src/common/paging";
import { Toaster } from "../notification.service";
import { isEnterKey } from "src/common/condition";
import { MngResDialogComponent } from "../mng-res-dialog/mng-res-dialog.component";
import { MatDialog } from "@angular/material/dialog";
import { HttpClient } from "@angular/common/http";

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
  pagingController: PagingController;

  readonly tabcol = [
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
    private toaster: Toaster,
    private dialog: MatDialog
  ) {}

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
      .post<any>(`${environment.uservault}/open/api/resource/list`, {
        paging: this.pagingController.paging,
      })
      .subscribe({
        next: (r) => {
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

  onPagingControllerReady(pc) {
    this.pagingController = pc;
    this.pagingController.onPageChanged = () => this.fetchList();
    this.fetchList();
  }

  createNewRes() {
    if (!this.newResName) {
      this.toaster.toast("Please enter new resource name");
      return;
    }
    if (!this.newResCode) {
      this.toaster.toast("Please enter new resource code");
      return;
    }

    this.http
      .post(`${environment.uservault}/open/api/resource/add`, {
        name: this.newResName,
        code: this.newResCode,
      })
      .subscribe({
        next: (r) => {
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
