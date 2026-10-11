import { Injectable } from "@angular/core";

// Mobile is inclusive at this boundary. Keep in step with the CSS tiers in
// src/styles.css (the CSS cannot import this constant): the global mobile tier
// ends at MOBILE_MAX_WIDTH and the desktop tier starts at MOBILE_MAX_WIDTH + 1.
const MOBILE_MAX_WIDTH = 768;

@Injectable({
  providedIn: "root",
})
export class Env {
  isMobile(): boolean {
    return window.innerWidth <= MOBILE_MAX_WIDTH;
  }
}
