import { Component, OnInit } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { getExpanded, isIdEqual } from "src/animate/animate-util";
import { PagingController } from "src/common/paging";
import { MngPathDialogComponent } from "../mng-path-dialog/mng-path-dialog.component";
import { isEnterKey } from "src/common/condition";
import { HttpClient } from "@angular/common/http";
import { MatSnackBar } from "@angular/material/snack-bar";

export interface CreatePathReq {
  type?: string; // path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible
  url?: string;
  group?: string;
  method?: string;
  desc?: string;
  resCode?: string;
}

export interface WPath {
  id?: number;
  pgroup?: string;
  pathNo?: string;
  method?: string;
  desc?: string;
  url?: string;
  ptype?: string;
  createTime?: Date;
  createBy?: string;
  updateTime?: Date;
  updateBy?: string;
}

@Component({
  selector: "app-manage-paths",
  templateUrl: "./manage-paths.component.html",
  styleUrls: ["./manage-paths.component.css"],
})
export class ManagePathsComponent implements OnInit {
  searchPath = null;
  searchGroup = null;
  searchType = null;
  PATH_TYPES = [
    { val: "PROTECTED", name: "Protected" },
    { val: "PUBLIC", name: "Public" },
  ];
  METHOD_OPTIONS = [
    { val: "GET", name: "Get" },
    { val: "PUT", name: "Put" },
    { val: "POST", name: "Post" },
    { val: "DELETE", name: "Delete" },
    { val: "HEAD", name: "Head" },
    { val: "OPTION", name: "Option" },
  ];
  newPathReq: CreatePathReq = { type: "PROTECTED" };
  showNewPath = false;

  expandedElement: WPath = null;
  pagingController: PagingController;

  readonly tabcol = [
    "id",
    "pgroup",
    "url",
    "method",
    "ptype",
    "desc",
    "createBy",
    "createTime",
  ];
  paths: WPath[] = [];

  idEquals = isIdEqual;
  getExpandedEle = (row) => getExpanded(row, this.expandedElement);
  isEnter = isEnterKey;

  constructor(
    private http: HttpClient,
    private dialog: MatDialog,
    private snackBar: MatSnackBar
  ) {}

  reset() {
    this.expandedElement = null;
    this.searchGroup = null;
    this.searchPath = null;
    this.searchType = null;
    this.pagingController.firstPage();
  }

  ngOnInit(): void {}

  openMngPathDialog(p: WPath) {
    this.dialog
      .open(MngPathDialogComponent, {
        width: "700px",
        data: {
          path: { ...p },
        },
      })
      .afterClosed()
      .subscribe({
        complete: () => {
          this.fetchList();
        },
      });
  }

  fetchList() {
    this.http
      .post<any>(`user-vault/open/api/path/list`, {
        paging: this.pagingController.paging,
        pgroup: this.searchGroup,
        url: this.searchPath,
        ptype: this.searchType,
      })
      .subscribe({
        next: (r) => {
          this.paths = [];
          if (r.data && r.data.payload) {
            for (let ro of r.data.payload) {
              if (ro.createTime) ro.createTime = new Date(ro.createTime);
              if (ro.updateTime) ro.updateTime = new Date(ro.updateTime);
              this.paths.push(ro);
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

  createPath() {
    this.http
      .post<any>(`/user-vault/remote/path/add`, this.newPathReq)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.newPathReq = { type: "PROTECTED" };
          this.showNewPath = false;
          this.fetchList();
        },
        error: (err) => {
          console.log(err);
          this.snackBar.open("Request failed, unknown error", "ok", {
            duration: 3000,
          });
        },
      });
  }
}
