import { Component, OnInit, ViewChild } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { MatSnackBar } from "@angular/material/snack-bar";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";
import { Env } from "src/common/env-util";

export interface AccessLog {
  id: number;
  accessTime: Date;
  success: boolean;
  ipAddress: string;
  username: string;
  url: string;
  userAgent: string;
}

@Component({
  selector: "app-access-log",
  templateUrl: "./access-log.component.html",
  styleUrls: ["./access-log.component.css"],
})
export class AccessLogComponent implements OnInit {
  readonly COLUMNS_TO_BE_DISPLAYED: string[] = this.env.isMobile()
    ? ["accessTime", "success", "ipAddress"]
    : ["accessTime", "success", "ipAddress", "userAgent", "url"];
  accessLogList: AccessLog[] = [];

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  constructor(
    private env: Env,
    private http: HttpClient,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit() {}

  fetchAccessLogList(): void {
    this.http
      .post<any>(`user-vault/open/api/access/history`, {
        paging: this.pagingController.paging,
      })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }

          if (resp.data.payload) {
            this.accessLogList = resp.data.payload;
          } else {
            this.accessLogList = [];
          }
          this.pagingController.onTotalChanged(resp.data.paging);
        },
        error: (err) => {
          console.log(err);
        },
      });
  }
}
