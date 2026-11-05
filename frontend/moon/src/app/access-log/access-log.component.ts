import { Component, OnDestroy, OnInit, ViewChild } from "@angular/core";
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
export class AccessLogComponent implements OnInit, OnDestroy {
  private readonly compactColumns = ["accessTime", "success", "ipAddress"];
  private readonly wideColumns = [
    "accessTime",
    "success",
    "ipAddress",
    "userAgent",
    "url",
  ];
  displayedColumns: string[] = this.wideColumns;
  private readonly modeSub = this.env.mobile$.subscribe((mobile) => {
    this.displayedColumns = mobile ? this.compactColumns : this.wideColumns;
  });

  accessLogList: AccessLog[] = [];

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  constructor(
    private env: Env,
    private http: HttpClient,
    private snackBar: MatSnackBar,
  ) {}

  ngOnInit() {}

  ngOnDestroy(): void {
    this.modeSub.unsubscribe();
  }

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
