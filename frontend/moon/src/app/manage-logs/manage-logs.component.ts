import { Component, OnInit, ViewChild } from "@angular/core";
import { Paging } from "src/common/paging";
import { isEnterKey } from "src/common/condition";
import { HttpClient } from "@angular/common/http";
import { Env } from "src/common/env-util";
import { MatSnackBar } from "@angular/material/snack-bar";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";

export interface ListedErrorLog {
  id?: number;
  node?: string;
  app?: string;
  caller?: string;
  traceId?: string;
  spanId?: string;
  errMsg?: string;
  rtime?: any;
}

export interface ListErrorLogReq {
  app?: string;
  page?: Paging;
}

export interface ListErrorLogResp {
  page: Paging;
  payload: ListedErrorLog[];
}

@Component({
  selector: "app-manage-logs",
  templateUrl: "./manage-logs.component.html",
  styleUrls: ["./manage-logs.component.css"],
})
export class ManageLogsComponent implements OnInit {
  readonly tabcol = this.env.isMobile()
    ? ["rtime", "errMsg"]
    : ["rtime", "app", "caller", "errMsg"];

  qryApp = "";
  tabdat = [];
  isEnter = isEnterKey;

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  constructor(
    private http: HttpClient,
    public env: Env,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {}

  reset() {
    this.qryApp = "";
    this.pagingController.firstPage();
  }

  fetchList() {
    this.http
      .post<any>(`logbot/log/error/list`, {
        app: this.qryApp,
        page: this.pagingController.paging,
      })
      .subscribe({
        next: (r) => {
          if (r.error) {
            this.snackBar.open(r.msg, "ok", { duration: 6000 });
            return;
          }
          this.tabdat = [];
          if (r.data && r.data.payload) {
            for (let ro of r.data.payload) {
              if (ro.ctime) ro.createTime = new Date(ro.ctime);
              this.tabdat.push(ro);
            }
          }
          this.pagingController.onTotalChanged(r.data.page);
        },
      });
  }
}
