import { Component, OnDestroy, OnInit } from "@angular/core";
import { UserInfo } from "src/common/user-info";
import { UserService } from "../user.service";
import { copyToClipboard } from "src/common/clipboard";
import { environment } from "src/environments/environment";
import { PlatformNotificationService } from "../platform-notification.service";
import { Toaster } from "../notification.service";
import { HttpClient } from "@angular/common/http";

@Component({
  selector: "app-nav",
  templateUrl: "./nav.component.html",
  styleUrls: ["./nav.component.css"],
})
export class NavComponent implements OnInit, OnDestroy {
  userInfo: UserInfo = null;
  copyToClipboard = (s) => {
    this.toaster.toast("Copied to clipboard");
    copyToClipboard(s);
  };
  unreadCount: 0;

  constructor(
    private userService: UserService,
    private http: HttpClient,
    private platformNotification: PlatformNotificationService,
    private toaster: Toaster
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
    return this.http
      .get<any>(`${environment.uservault}/open/api/v1/notification/count`)
      .subscribe({
        next: (res) => (this.unreadCount = res.data),
      });
  }
}
