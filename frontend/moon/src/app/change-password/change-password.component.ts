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
import { I18n } from "../i18n.service";

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
    private i18n: I18n
  ) {}

  ngOnInit() {}

  changePassword() {
    if (
      !hasText(this.changePasswordParam.prevPassword) ||
      !hasText(this.changePasswordParam.newPassword) ||
      !hasText(this.newPasswordConfirm)
    ) {
      this.snackBar.open(this.i18n.trl("change-password", "pleaseEnterPasswords"), "ok", { duration: 3000 });
      return;
    }

    if (this.changePasswordParam.newPassword !== this.newPasswordConfirm) {
      this.snackBar.open(this.i18n.trl("change-password", "confirmedPasswordNotMatched"), "ok", {
        duration: 3000,
      });
      return;
    }

    if (
      this.changePasswordParam.prevPassword ===
      this.changePasswordParam.newPassword
    ) {
      this.snackBar.open(this.i18n.trl("change-password", "newPasswordMustBeDifferent"), "ok", {
        duration: 3000,
      });
      return;
    }

    this.userService.changePassword(this.changePasswordParam).subscribe({
      next: (result) => {
        if (result.error) {
          this.snackBar.open(result.msg, "ok", { duration: 6000 });
          return;
        }
        this.snackBar.open(this.i18n.trl("change-password", "passwordChanged"), "ok", { duration: 3000 });
        this.nav.navigateTo(NavType.MANAGE_USER);
      },
      complete: () => {
        this.changePasswordParam = emptyChangePasswordParam();
        this.newPasswordConfirm = null;
      },
    });
  }
}
