import { Component, OnDestroy, OnInit } from "@angular/core";
import { UserInfo } from "src/common/user-info";
import { UserService } from "../user.service";
import { copyToClipboard } from "src/common/clipboard";
import { PlatformNotificationService } from "../platform-notification.service";
import { Version } from "../version";
import { MatSnackBar } from "@angular/material/snack-bar";
import { I18n } from "../i18n.service";
import { WebSocketNotificationService } from "../websocket-notification.service";

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
  window = window;

  constructor(
    public i18n: I18n,
    private userService: UserService,
    private platformNotification: PlatformNotificationService,
    private snackBar: MatSnackBar,
    private wsService: WebSocketNotificationService,
  ) {
    this.wsService.count$.subscribe((count) => {
      this.unreadCount = count;
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
  }

  /** log out current user and navigate back to login page */
  logout(): void {
    this.userService.logout();
  }

  openGithub() {
    window.open("https://github.com/curtisnewbie/moon-monorepo", "_blank");
  }
}
