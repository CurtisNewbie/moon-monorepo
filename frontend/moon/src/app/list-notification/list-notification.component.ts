import { Component, Inject, OnInit, ViewChild } from "@angular/core";
import {
  MAT_DIALOG_DATA,
  MatDialog,
  MatDialogRef,
} from "@angular/material/dialog";
import { PlatformNotificationService } from "../platform-notification.service";
import { HttpClient } from "@angular/common/http";
import { Env } from "src/common/env-util";
import { MatSnackBar } from "@angular/material/snack-bar";
import { ConfirmDialogComponent } from "../dialog/confirm/confirm-dialog.component";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";

export interface ShowNatificationDialogData {
  title: string;
  time: string;
  message: string;
}

@Component({
  selector: "show-notification-dialog-component",
  template: `
    <h1 mat-dialog-title>{{ data.title }}</h1>
    <div mat-dialog-content>
      <p>{{ data.time }}</p>
      <pre style="text-wrap: pretty; line-break: anywhere">{{
        data.message
      }}</pre>
    </div>
    <div mat-dialog-actions class="d-flex justify-content-end">
      <button mat-button [mat-dialog-close]="true">Ok</button>
    </div>
  `,
})
export class ShowNotificationDialogComponent {
  constructor(
    public dialogRef: MatDialogRef<
      ShowNotificationDialogComponent,
      ShowNatificationDialogData
    >,
    @Inject(MAT_DIALOG_DATA) public data: ShowNatificationDialogData
  ) {}
}

export interface Notification {
  id: number;
  notifiNo: string;
  title: string;
  brief: string;
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
  readonly columns: string[] = this.env.isMobile()
    ? ["title", "status", "createTime"]
    : ["id", "notifiNo", "title", "brief", "status", "createTime"];
  query = {
    onlyInitMessage: true,
  };
  data: Notification[] = [];

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  constructor(
    private http: HttpClient,
    private dialog: MatDialog,
    private platformNotification: PlatformNotificationService,
    public env: Env,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {}

  fetchList() {
    this.http
      .post<any>(`user-vault/open/api/v1/notification/query`, {
        status: this.query.onlyInitMessage ? "INIT" : "",
        page: this.pagingController.paging,
      })
      .subscribe((resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        if (resp.data) {
          this.data = [];
          if (resp.data.payload) {
            for (let r of resp.data.payload) {
              if (r.createTime) r.createTime = new Date(r.createTime);
              if (r.message == null) {
                r.message = "";
              }

              let th = 100;
              if (r.message.length <= th) {
                r.brief = r.message;
              } else {
                r.brief = "... " + r.message.substring(r.message.length - th);
              }

              this.data.push(r);
            }
          }
          this.pagingController.onTotalChanged(resp.data.paging);
        }
      });
  }

  reset() {
    this.query.onlyInitMessage = true;
    if (!this.pagingController.firstPage()) {
      this.fetchList();
    }
  }

  markOpened(notifiNo: string) {
    this.http
      .post<any>(`user-vault/open/api/v1/notification/open`, {
        notifiNo: notifiNo,
      })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
        },
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

    const dialogRef: MatDialogRef<ShowNotificationDialogComponent, boolean> =
      this.dialog.open(ShowNotificationDialogComponent, {
        width: "900px",
        data: {
          title: n.title,
          time: `Notification Time: ${timeStr}`,
          message: n.message,
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
          .post<any>(`user-vault/open/api/v1/notification/open-all`, {
            notifiNo: last,
          })
          .subscribe({
            next: (resp) => {
              if (resp.error) {
                this.snackBar.open(resp.msg, "ok", { duration: 6000 });
                return;
              }
              this.platformNotification.triggerChange();
              this.fetchList();
            },
          });
      }
    });
  }
}
