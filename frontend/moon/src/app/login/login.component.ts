import { Component, OnInit } from "@angular/core";
import { NavigationService } from "../navigation.service";
import { NavType } from "../routes";
import { UserService } from "../user.service";
import { setToken } from "src/common/api-util";
import { PlatformNotificationService } from "../platform-notification.service";
import { MatSnackBar } from "@angular/material/snack-bar";

@Component({
  selector: "app-login",
  templateUrl: "./login.component.html",
  styleUrls: ["./login.component.css"],
})
export class LoginComponent implements OnInit {
  usernameInput: string = "";
  passwordInput: string = "";

  constructor(
    private userService: UserService,
    private nav: NavigationService,
    private platformNotification: PlatformNotificationService,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit() {
    this.userService.userInfoObservable.subscribe((user) => {
      if (user) {
        this.nav.navigateTo(NavType.USER_DETAILS);
        this.platformNotification.triggerChange();
      }
    });
  }

  /**
   * login request
   */
  public login(): void {
    if (!this.usernameInput || !this.passwordInput) {
      this.snackBar.open("Please enter username and password", "ok", {
        duration: 3000,
      });
      return;
    }
    this.userService.login(this.usernameInput, this.passwordInput).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        setToken(resp.data);
        this.routeToHomePage();
        this.userService.fetchUserInfo();
        this.userService.fetchUserResources();
      },
      complete: () => {
        this.passwordInput = "";
        this.platformNotification.triggerChange();
      },
    });
  }

  goToRegisterPage(): void {
    this.nav.navigateTo(NavType.REGISTER_PAGE);
  }

  private routeToHomePage(): void {
    this.nav.navigateTo(NavType.USER_DETAILS);
  }

  passwordInputKeyPressed(event: any): void {
    if (event.key === "Enter") {
      this.login();
    }
  }
}
