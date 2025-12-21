import { Component, OnDestroy, OnInit } from "@angular/core";
import { UserInfo } from "src/common/user-info";
import { UserService } from "../user.service";
import { copyToClipboard } from "src/common/clipboard";
import { PlatformNotificationService } from "../platform-notification.service";
import { HttpClient } from "@angular/common/http";
import { Version } from "../version";
import { MatSnackBar } from "@angular/material/snack-bar";
import { I18n } from "../i18n.service";

@Component({
  selector: "app-nav",
  templateUrl: "./nav.component.html",
  styleUrls: ["./nav.component.css"],
})
export class NavComponent implements OnInit, OnDestroy {
  version = Version;
  userInfo: UserInfo = null;
  copyToClipboard = (s) => {
    this.snackBar.open("Copied to clipboard", "ok", { duration: 3000 });
    copyToClipboard(s);
  };
  unreadCount = 0;
  fetching = false;

  constructor(
    public i18n: I18n,
    private userService: UserService,
    private http: HttpClient,
    private platformNotification: PlatformNotificationService,
    private snackBar: MatSnackBar
  ) {
    platformNotification.subscribeChange().subscribe({
      next: () => {
        if (this.userInfo) {
          this.fetchUnreadNotificationCount();
        }
      },
    });
  }

  ngOnDestroy(): void {}

  hasRes(code) {
    return this.userService.hasResource(code);
  }

  hasAnyRes(...codes: string[]) {
    for (let c of codes) {
      if (this.hasRes(c)) return true;
    }
    return false;
  }

  ngOnInit(): void {
    this.userService.userInfoObservable.subscribe({
      next: (user) => {
        this.userInfo = user;
      },
    });

    if (this.userInfo) {
      this.fetchUnreadNotificationCount();
    }
  }

  /** log out current user and navigate back to login page */
  logout(): void {
    this.userService.logout();
  }

  fetchUnreadNotificationCount() {
    if (this.fetching) {
      return;
    }

    this.fetching = true;
    console.log("fetch unread notification count");
    return this.http
      .get<any>(
        `user-vault/open/api/v2/notification/count?curr=${this.unreadCount}`
      )
      .subscribe({
        next: (resp) => {
          if (resp) {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 });
              return;
            }
            this.unreadCount = resp.data;
          }
        },
        complete: () => {
          this.fetching = false;
        },
        error: () => {
          this.fetching = false;
        },
      });
  }

  openGithub() {
    window.open("https://github.com/curtisnewbie/moon-monorepo", "_blank");
  }
}
