import { Component, OnInit } from "@angular/core";
import { environment } from "src/environments/environment";
import { PagingController } from "src/common/paging";
import { HttpClient } from "@angular/common/http";

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
  readonly COLUMNS_TO_BE_DISPLAYED: string[] = [
    "id",
    "user",
    "accessTime",
    "success",
    "ipAddress",
    "userAgent",
    "url",
  ];
  accessLogList: AccessLog[] = [];
  pagingController: PagingController;

  constructor(private http: HttpClient) {}

  ngOnInit() {}

  fetchAccessLogList(): void {
    this.http
      .post<any>(`${environment.uservault}/open/api/access/history`, {
        paging: this.pagingController.paging,
      })
      .subscribe({
        next: (resp) => {
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

  onPagingControllerReady(pc) {
    this.pagingController = pc;
    this.pagingController.onPageChanged = () => this.fetchAccessLogList();
    this.fetchAccessLogList();
  }
}
