import { Injectable } from "@angular/core";

@Injectable({
  providedIn: "root",
})
export class Env {
  isMobile(): boolean {
    return window.innerWidth < 768;
  }
}
