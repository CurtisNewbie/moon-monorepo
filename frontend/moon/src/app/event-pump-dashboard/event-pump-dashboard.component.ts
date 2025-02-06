import { Component, OnInit } from "@angular/core";
import { setEventPumpCookie } from "src/common/api-util";

@Component({
  selector: "app-event-pump-dashboard",
  templateUrl: "./event-pump-dashboard.component.html",
  styleUrls: ["./event-pump-dashboard.component.css"],
})
export class EventPumpDashboardComponent implements OnInit {
  constructor() {}

  ngOnInit(): void {
    setEventPumpCookie();
  }
}
