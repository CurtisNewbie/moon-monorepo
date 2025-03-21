import { Component, OnInit } from "@angular/core";
import { environment } from "src/environments/environment";
import { PagingController } from "src/common/paging";
import { HttpClient } from "@angular/common/http";
import { MatSnackBar } from "@angular/material/snack-bar";

export interface OperateLog {
  /** name of operation */
  operateName: string;

  /** description of operation */
  operateDesc: string;

  /** when the operation happens */
  operateTime: Date;

  /** parameters used for the operation */
  operateParam: string;

  /** username */
  username: string;

  /** primary key of user */
  userId: number;
}

@Component({
  selector: "app-operate-history",
  templateUrl: "./operate-history.component.html",
  styleUrls: ["./operate-history.component.css"],
})
export class OperateHistoryComponent implements OnInit {
  operateLogList: OperateLog[] = [];
  pagingController: PagingController;
  COLUMNS_TO_BE_DISPLAYED = [
    "id",
    "user",
    "operateName",
    "operateDesc",
    "operateTime",
    "operateParam",
  ];

  constructor(private http: HttpClient, private snackBar: MatSnackBar) {}

  ngOnInit() {}

  fetchOperateLogList(): void {
    this.http
      .post<any>(
        `user-vault/open/api/operate/history`,
        this.pagingController.paging
      )
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.operateLogList = [];
          if (resp.data.operateLogVoList) {
            for (let r of resp.data.operateLogVoList) {
              if (r.operateTime) r.operateTime = new Date(r.operateTime);
              this.operateLogList.push(r);
            }
          }
          this.pagingController.onTotalChanged(resp.data.paging);
        },
      });
  }

  onPagingControllerReady(pc) {
    this.pagingController = pc;
    this.pagingController.onPageChanged = () => this.fetchOperateLogList();
    this.fetchOperateLogList();
  }
}
