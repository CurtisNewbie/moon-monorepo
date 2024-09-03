import { Component, OnInit } from "@angular/core";
import {
  ChangePasswordParam,
  emptyChangePasswordParam,
} from "src/common/user-info";
import { NavigationService } from "../navigation.service";
import { Toaster } from "../notification.service";
import { UserService } from "../user.service";
import { hasText } from "src/common/str-util";
import { NavType } from "../routes";

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
    private toaster: Toaster
  ) { }

  ngOnInit() {
  }

  changePassword() {
    if (
      !hasText(this.changePasswordParam.prevPassword) ||
      !hasText(this.changePasswordParam.newPassword) ||
      !hasText(this.newPasswordConfirm)
    ) {
      this.toaster.toast("Please enter passwords");
      return;
    }

    if (this.changePasswordParam.newPassword !== this.newPasswordConfirm) {
      this.toaster.toast("Confirmed password is not matched");
      return;
    }

    if (
      this.changePasswordParam.prevPassword ===
      this.changePasswordParam.newPassword
    ) {
      this.toaster.toast("new password must be different");
      return;
    }

    this.userService.changePassword(this.changePasswordParam).subscribe({
      next: (result) => {
        this.toaster.toast("Password changed");
        this.nav.navigateTo(NavType.MANAGE_USER);
      },
      complete: () => {
        this.changePasswordParam = emptyChangePasswordParam();
        this.newPasswordConfirm = null;
      },
    });
  }
}
