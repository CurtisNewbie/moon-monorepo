import { Component, OnInit } from "@angular/core";
import {
  ChangePasswordParam,
  emptyChangePasswordParam,
} from "src/common/user-info";
import { NavigationService } from "../navigation.service";
import { UserService } from "../user.service";
import { hasText } from "src/common/str-util";
import { NavType } from "../routes";
import { MatSnackBar } from "@angular/material/snack-bar";

@Component({
  selector: "app-change-password",
  templateUrl: "./change-password.component.html",
  styleUrls: ["./change-password.component.css"],
})
export class ChangePasswordComponent implements OnInit {
  changePasswordParam: ChangePasswordParam = emptyChangePasswordParam();
  newPasswordConfirm: string = null;

  constructor(
    private nav: NavigationService,
    private userService: UserService,
    private snackBar: MatSnackBar,
  ) { }

  ngOnInit() {
  }

  changePassword() {
    if (
      !hasText(this.changePasswordParam.prevPassword) ||
      !hasText(this.changePasswordParam.newPassword) ||
      !hasText(this.newPasswordConfirm)
    ) {
      this.snackBar.open("Please enter passwords", "ok", { duration: 3000 });;
      return;
    }

    if (this.changePasswordParam.newPassword !== this.newPasswordConfirm) {
      this.snackBar.open("Confirmed password is not matched", "ok", { duration: 3000 });;
      return;
    }

    if (
      this.changePasswordParam.prevPassword ===
      this.changePasswordParam.newPassword
    ) {
      this.snackBar.open("new password must be different", "ok", { duration: 3000 });;
      return;
    }

    this.userService.changePassword(this.changePasswordParam).subscribe({
      next: (result) => {
        this.snackBar.open("Password changed", "ok", { duration: 3000 });;
        this.nav.navigateTo(NavType.MANAGE_USER);
      },
      complete: () => {
        this.changePasswordParam = emptyChangePasswordParam();
        this.newPasswordConfirm = null;
      },
    });
  }
}
