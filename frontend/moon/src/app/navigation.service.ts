import { Injectable } from "@angular/core";
import { Router } from "@angular/router";
import { NavType } from "./routes";

@Injectable({
  providedIn: "root",
})
export class NavigationService {
  constructor(private router: Router) { }

  /** Navigate to using Router*/
  public navigateTo(nt: NavType, extra?: any[]): void {
    this.navigateToUrl(nt, extra);
  }

  /** Navigate to using Router*/
  public navigateToUrl(url: string, extra?: any[]): void {
    let arr: any[] = [url];
    if (extra != null) arr = arr.concat(extra);
    this.router.navigate(arr);
  }
}

