import { Injectable } from "@angular/core";

@Injectable({
  providedIn: "root",
})
export class Env {
  private _isMobile: boolean = null;

  isMobile(): boolean {
    if (this._isMobile != null) {
      return this._isMobile;
    }
    this._isMobile = /iPhone|iPad|iPod|Android/i.test(navigator.userAgent);
    return this._isMobile;
  }
}
