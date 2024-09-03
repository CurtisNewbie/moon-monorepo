import { Component } from "@angular/core";
import { UserService } from "./user.service";

@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.css"],
})
export class AppComponent {
  constructor(private userService: UserService) {}

  ngOnInit(): void {
    this.userService.userInfoObservable.subscribe({
      next: (user) => {
        if (user) {
          this.userService.fetchUserResources();
        }
      },
    });
    this.userService.fetchUserInfo();
  }
}
