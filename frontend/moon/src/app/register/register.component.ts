import { Component, OnInit } from "@angular/core";
import { NavigationService } from "../navigation.service";
import { UserService } from "../user.service";
import { isEnterKey } from "src/common/condition";
import { NavType } from "../routes";
import { MatSnackBar } from "@angular/material/snack-bar";

@Component({
  selector: "app-register",
  templateUrl: "./register.component.html",
  styleUrls: ["./register.component.css"],
})
export class RegisterComponent implements OnInit {
  usernameInput: string = "";
  passwordInput: string = "";
  isEnter = isEnterKey;

  constructor(
    private userService: UserService,
    private nav: NavigationService,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit() {}

  register(): void {
    if (!this.usernameInput || !this.passwordInput) {
      this.snackBar.open("Please enter username and password", "ok", {
        duration: 3000,
      });
      return;
    }
    this.userService
      .register(this.usernameInput, this.passwordInput)
      .subscribe({
        next: (r) => {
          if (r.error) {
            this.snackBar.open(r.msg, "ok", { duration: 5000 });
            return;
          }

          this.snackBar.open(
            "Registered, please wait for administrator's approval",
            "ok",
            { duration: 5000 }
          );
          this.nav.navigateTo(NavType.LOGIN_PAGE);
        },
        complete: () => {
          this.usernameInput = "";
          this.passwordInput = "";
        },
      });
  }

  gotoLoginPage(): void {
    this.nav.navigateTo(NavType.LOGIN_PAGE);
  }
}
