import { Injectable, OnDestroy } from "@angular/core";
import { BehaviorSubject, Subscription, timer } from "rxjs";
import { Observable } from "rxjs";
import { Resp } from "src/common/resp";
import { ChangePasswordParam, UserInfo } from "src/common/user-info";
import { NavigationService } from "./navigation.service";
import { NavType } from "./routes";
import { getToken, setToken, onEmptyToken } from "src/common/api-util";
import { HttpClient } from "@angular/common/http";
import { Router } from "@angular/router";
import { MatSnackBar } from "@angular/material/snack-bar";

export interface RoleBrief {
  roleNo?: string;
  name?: string;
  code?: string;
}

export interface ResBrief {
  code?: string;
  name?: string;
}

@Injectable({
  providedIn: "root",
})
export class UserService implements OnDestroy {
  private userInfoSubject = new BehaviorSubject<UserInfo>(null);
  private resourceSubject = new BehaviorSubject<any>(null);
  private resources: Set<string> = null;

  // refreshed every 5min
  private tokenRefresher: Subscription = timer(60_000, 360_000).subscribe(
    () => {
      let t = getToken();
      if (t != null) {
        this.exchangeToken(t).subscribe({
          next: (resp) => {
            setToken(resp.data);
          },
        });
      }
    }
  );

  userInfoObservable: Observable<UserInfo> =
    this.userInfoSubject.asObservable();

  resourceObservable: Observable<any> = this.resourceSubject.asObservable();

  constructor(
    private http: HttpClient,
    private nav: NavigationService,
    private router: Router,
    private snackBar: MatSnackBar
  ) {
    onEmptyToken(() => {
      if (this.router.url != "/register") {
        this.logout();
      }
    });
  }

  public fetchUserResources() {
    this.http.get<any>(`user-vault/open/api/resource/brief/user`).subscribe({
      next: (res) => {
        this.resources = new Set();
        if (res.data) {
          for (let r of res.data) {
            this.resources.add(r.code);
          }
          this.resourceSubject.next(this.resources);
        }
      },
    });
  }

  ngOnDestroy(): void {
    this.tokenRefresher.unsubscribe();
  }

  hasResource(code): boolean {
    if (this.resources == null) return false;
    return this.resources.has(code);
  }

  /**
   * Attempt to sign-in
   * @param username
   * @param password
   */
  public login(username: string, password: string): Observable<Resp<any>> {
    return this.http.post<Resp<any>>(`user-vault/open/api/user/login`, {
      username: username,
      password: password,
    });
  }

  /**
   * Logout current user
   */
  public logout(): void {
    setToken(null);
    this.resources = null;
    this._notifyUserInfo(null);
    this.nav.navigateTo(NavType.LOGIN_PAGE);
  }

  /**
   * Add user, only admin is allowed to add user
   * @param username
   * @param password
   * @returns
   */
  public addUser(
    username: string,
    password: string,
    userRole: string
  ): Observable<Resp<any>> {
    return this.http.post<any>(`user-vault/open/api/user/register`, {
      username,
      password,
      userRole,
    });
  }

  /**
   * Register user
   * @param username
   * @param password
   * @returns
   */
  public register(username: string, password: string): Observable<Resp<any>> {
    return this.http.post<any>(`user-vault/open/api/user/register/request`, {
      username,
      password,
    });
  }

  /**
   * Fetch user info
   */
  public fetchUserInfo(callback = null): void {
    this.http.get<any>(`user-vault/open/api/user/info`).subscribe({
      next: (resp) => {
        if (resp.data) {
          this._notifyUserInfo(resp.data);
          if (callback) {
            callback();
          }
        } else {
          this.snackBar.open("Please login first", "ok", { duration: 3000 });
          setToken(null);
          this._notifyUserInfo(null);
          this.nav.navigateTo(NavType.LOGIN_PAGE);
        }
      },
    });
  }

  private _notifyUserInfo(userInfo: UserInfo): void {
    this.userInfoSubject.next(userInfo);
  }

  /**
   * Fetch user details
   */
  public fetchUserDetails(): Observable<
    Resp<{
      id;
      username;
      role;
      registerDate;
    }>
  > {
    return this.http.get<any>(`user-vault/open/api/user/info`);
  }

  /**
   * Exchange Token
   */
  private exchangeToken(token: string): Observable<Resp<string>> {
    return this.http.post<any>(`user-vault/open/api/token/exchange`, {
      token: token,
    });
  }

  /**
   * Change password
   */
  public changePassword(param: ChangePasswordParam): Observable<Resp<any>> {
    return this.http.post<any>(
      `user-vault/open/api/user/password/update`,
      param
    );
  }

  public fetchRoleBriefs(): Observable<Resp<any>> {
    return this.http.get<any>(`user-vault/open/api/role/brief/all`);
  }

  public fetchAllResBrief(): Observable<Resp<any>> {
    return this.http.get<any>(`user-vault/open/api/resource/brief/all`);
  }
}
