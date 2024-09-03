import { Injectable, OnDestroy } from "@angular/core";
import { BehaviorSubject, Subscription, timer } from "rxjs";
import { Observable } from "rxjs";
import { environment } from "src/environments/environment";
import { Resp } from "src/common/resp";
import { ChangePasswordParam, UserInfo } from "src/common/user-info";
import { NavigationService } from "./navigation.service";
import { Toaster } from "./notification.service";
import { NavType } from "./routes";
import { getToken, setToken, onEmptyToken } from "src/common/api-util";
import { HttpClient } from "@angular/common/http";

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

  constructor(
    private http: HttpClient,
    private nav: NavigationService,
    private notifi: Toaster
  ) {
    onEmptyToken(() => this.logout());
  }

  public fetchUserResources() {
    this.http
      .get<any>(`${environment.uservault}/open/api/resource/brief/user`)
      .subscribe({
        next: (res) => {
          this.resources = new Set();
          if (res.data) {
            for (let r of res.data) {
              this.resources.add(r.code);
            }
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
    return this.http.post<Resp<any>>(
      `${environment.uservault}/open/api/user/login`,
      {
        username: username,
        password: password,
      }
    );
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
    return this.http.post<any>(
      `${environment.uservault}/open/api/user/register`,
      {
        username,
        password,
        userRole,
      }
    );
  }

  /**
   * Register user
   * @param username
   * @param password
   * @returns
   */
  public register(username: string, password: string): Observable<Resp<any>> {
    return this.http.post<any>(
      `${environment.uservault}/open/api/user/register/request`,
      { username, password }
    );
  }

  /**
   * Fetch user info
   */
  public fetchUserInfo(callback = null): void {
    this.http
      .get<any>(`${environment.uservault}/open/api/user/info`)
      .subscribe({
        next: (resp) => {
          if (resp.data) {
            this._notifyUserInfo(resp.data);
            if (callback) {
              callback();
            }
          } else {
            this.notifi.toast("Please login first");
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
    return this.http.get<any>(`${environment.uservault}/open/api/user/info`);
  }

  /**
   * Exchange Token
   */
  private exchangeToken(token: string): Observable<Resp<string>> {
    return this.http.post<any>(
      `${environment.uservault}/open/api/token/exchange`,
      {
        token: token,
      }
    );
  }

  /**
   * Change password
   */
  public changePassword(param: ChangePasswordParam): Observable<Resp<any>> {
    return this.http.post<any>(
      `${environment.uservault}/open/api/user/password/update`,
      param
    );
  }

  public fetchRoleBriefs(): Observable<Resp<any>> {
    return this.http.get<any>(`${environment.uservault}/open/api/role/brief/all`);
  }

  public fetchAllResBrief(): Observable<Resp<any>> {
    return this.http.get<any>(
      `${environment.uservault}/open/api/resource/brief/all`
    );
  }
}
