import { Component, OnInit } from "@angular/core";

import { environment } from "src/environments/environment";
import { PagingController } from "src/common/paging";
import { ConfirmDialogComponent } from "../dialog/confirm/confirm-dialog.component";
import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { PlatformNotificationService } from "../platform-notification.service";
import { HttpClient } from "@angular/common/http";

export interface Notification {
  id: number;
  notifiNo: string;
  title: string;
  message: string;
  status: string;
  createTime: Date;
}

@Component({
  selector: "app-list-notification",
  templateUrl: "./list-notification.component.html",
  styleUrls: ["./list-notification.component.css"],
})
export class ListNotificationComponent implements OnInit {
  readonly columns: string[] = [
    "id",
    "notifiNo",
    "title",
    "status",
    "createTime",
  ];
  query = {
    onlyInitMessage: true,
  };
  pagingController: PagingController;
  data: Notification[] = [];

  constructor(
    private http: HttpClient,
    private dialog: MatDialog,
    private platformNotification: PlatformNotificationService
  ) {}

  ngOnInit(): void {}

  fetchList() {
    this.http
      .post<any>(`${environment.uservault}/open/api/v1/notification/query`, {
        status: this.query.onlyInitMessage ? "INIT" : "",
        page: this.pagingController.paging,
      })
      .subscribe((resp) => {
        if (resp.data) {
          this.data = [];
          if (resp.data.payload) {
            for (let r of resp.data.payload) {
              if (r.createTime) r.createTime = new Date(r.createTime);
              this.data.push(r);
            }
          }
          this.pagingController.onTotalChanged(resp.data.paging);
        }
      });
  }

  reset() {
    this.query = {
      onlyInitMessage: true,
    };
    if (!this.pagingController.firstPage()) {
      this.fetchList();
    }
  }

  onPagingControllerReady(pc) {
    this.pagingController = pc;
    this.pagingController.onPageChanged = () => this.fetchList();
    this.fetchList();
  }

  markOpened(notifiNo: string) {
    this.http
      .post<any>(`${environment.uservault}/open/api/v1/notification/open`, {
        notifiNo: notifiNo,
      })
      .subscribe({
        complete: () => {
          this.platformNotification.triggerChange();
        },
      });
  }

  showNotification(n: Notification) {
    let timeStr = "";
    if (n.createTime) {
      timeStr = n.createTime.toISOString().split(".")[0].replace("T", " ");
    }
    let lines = n.message.split(`\n`);

    const dialogRef: MatDialogRef<ConfirmDialogComponent, boolean> =
      this.dialog.open(ConfirmDialogComponent, {
        width: "700px",
        data: {
          title: n.title,
          msg: [`Notification Time: ${timeStr}`, ...lines],
        },
      });

    dialogRef.afterOpened().subscribe(() => {
      if (n.status != "OPENED") {
        this.markOpened(n.notifiNo);
      }
    });
    dialogRef.afterClosed().subscribe(() => this.fetchList());
  }

  markAllOpened() {
    if (!this.data) {
      return;
    }
    let last = this.data[0].notifiNo;

    const dialogRef: MatDialogRef<ConfirmDialogComponent, boolean> =
      this.dialog.open(ConfirmDialogComponent, {
        width: "700px",
        data: {
          title: "Mark All Notifications Opened?",
          msg: ["Are your sure you want to mark all notifications as opened?"],
        },
      });

    dialogRef.afterClosed().subscribe((res) => {
      if (res) {
        this.http
          .post<any>(
            `${environment.uservault}/open/api/v1/notification/open-all`,
            {
              notifiNo: last,
            }
          )
          .subscribe({
            next: () => {
              this.platformNotification.triggerChange();
              this.fetchList();
            },
          });
      }
    });
  }
}
