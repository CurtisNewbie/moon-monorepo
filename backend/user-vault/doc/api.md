# API Endpoints

- POST /open/api/user/login
  - Description: User Login using password, a JWT token is generated and returned
  - Expected Access Scope: PUBLIC
  - Header Parameter:
    - "x-forwarded-for": 
    - "user-agent": 
  - JSON Request:
    - "username": (string) 
    - "password": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (string) response data
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/user/login' \
      -H 'x-forwarded-for: ' \
      -H 'user-agent: ' \
      -H 'Content-Type: application/json' \
      -d '{"password":"","username":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiUserLogin(rail miso.Rail, req LoginReq, xForwardedFor string, userAgent string) (string, error) {
    	var res miso.GnResp[string]
    	err := miso.NewDynTClient(rail, "/open/api/user/login", "user-vault").
    		AddHeader("xForwardedFor", xForwardedFor).
    		AddHeader("userAgent", userAgent).
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return "", err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface LoginReq {
      username?: string;
      password?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: string;                 // response data
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    userLogin() {
      let xForwardedFor: any | null = null;
      let userAgent: any | null = null;
      let req: LoginReq | null = null;
      this.http.post<any>(`/user-vault/open/api/user/login`, req,
        {
          headers: {
            "x-forwarded-for": xForwardedFor
            "user-agent": userAgent
          }
        })
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: string = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/user/register/request
  - Description: User request registration, approval needed
  - Expected Access Scope: PUBLIC
  - JSON Request:
    - "username": (string) 
    - "password": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/user/register/request' \
      -H 'Content-Type: application/json' \
      -d '{"password":"","username":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiUserRegister(rail miso.Rail, req RegisterReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/user/register/request", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface RegisterReq {
      username?: string;
      password?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    userRegister() {
      let req: RegisterReq | null = null;
      this.http.post<any>(`/user-vault/open/api/user/register/request`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/user/add
  - Description: Admin create new user
  - Bound to Resource: `"manage-users"`
  - JSON Request:
    - "username": (string) 
    - "password": (string) 
    - "roleNo": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/user/add' \
      -H 'Content-Type: application/json' \
      -d '{"password":"","roleNo":"","username":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiAdminAddUser(rail miso.Rail, req AddUserParam) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/user/add", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface AddUserParam {
      username?: string;
      password?: string;
      roleNo?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    adminAddUser() {
      let req: AddUserParam | null = null;
      this.http.post<any>(`/user-vault/open/api/user/add`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/user/list
  - Description: Admin list users
  - Bound to Resource: `"manage-users"`
  - JSON Request:
    - "username": (*string) 
    - "roleNo": (*string) 
    - "isDisabled": (*int) 
    - "paging": (Paging) 
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (PageRes[github.com/curtisnewbie/user-vault/api.UserInfo]) response data
      - "paging": (Paging) pagination parameters
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
      - "payload": ([]api.UserInfo) payload values in current page
        - "id": (int) 
        - "username": (string) 
        - "roleName": (string) 
        - "roleNo": (string) 
        - "userNo": (string) 
        - "reviewStatus": (string) 
        - "isDisabled": (int) 
        - "createTime": (int64) 
        - "createBy": (string) 
        - "updateTime": (int64) 
        - "updateBy": (string) 
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/user/list' \
      -H 'Content-Type: application/json' \
      -d '{"isDisabled":0,"paging":{"limit":0,"page":0,"total":0},"roleNo":"","username":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiAdminListUsers(rail miso.Rail, req ListUserReq) (PageRes, error) {
    	var res miso.GnResp[PageRes]
    	err := miso.NewDynTClient(rail, "/open/api/user/list", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat PageRes
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ListUserReq {
      username?: string;
      roleNo?: string;
      isDisabled?: number;
      paging?: Paging;
    }

    export interface Paging {
      limit?: number;                // page limit
      page?: number;                 // page number, 1-based
      total?: number;                // total count
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: PageRes;
    }

    export interface PageRes {
      paging?: Paging;
      payload?: UserInfo[];
    }

    export interface Paging {
      limit?: number;                // page limit
      page?: number;                 // page number, 1-based
      total?: number;                // total count
    }

    export interface UserInfo {
      id?: number;
      username?: string;
      roleName?: string;
      roleNo?: string;
      userNo?: string;
      reviewStatus?: string;
      isDisabled?: number;
      createTime?: number;
      createBy?: string;
      updateTime?: number;
      updateBy?: string;
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    adminListUsers() {
      let req: ListUserReq | null = null;
      this.http.post<any>(`/user-vault/open/api/user/list`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: PageRes = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/user/info/update
  - Description: Admin update user info
  - Bound to Resource: `"manage-users"`
  - JSON Request:
    - "userNo": (string) 
    - "roleNo": (string) 
    - "isDisabled": (int) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/user/info/update' \
      -H 'Content-Type: application/json' \
      -d '{"isDisabled":0,"roleNo":"","userNo":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiAdminUpdateUser(rail miso.Rail, req AdminUpdateUserReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/user/info/update", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface AdminUpdateUserReq {
      userNo?: string;
      roleNo?: string;
      isDisabled?: number;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    adminUpdateUser() {
      let req: AdminUpdateUserReq | null = null;
      this.http.post<any>(`/user-vault/open/api/user/info/update`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/user/registration/review
  - Description: Admin review user registration
  - Bound to Resource: `"manage-users"`
  - JSON Request:
    - "userId": (int) 
    - "reviewStatus": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/user/registration/review' \
      -H 'Content-Type: application/json' \
      -d '{"reviewStatus":"","userId":0}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiAdminReviewUser(rail miso.Rail, req AdminReviewUserReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/user/registration/review", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface AdminReviewUserReq {
      userId?: number;
      reviewStatus?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    adminReviewUser() {
      let req: AdminReviewUserReq | null = null;
      this.http.post<any>(`/user-vault/open/api/user/registration/review`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- GET /open/api/user/info
  - Description: User get user info
  - Expected Access Scope: PUBLIC
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (UserInfoRes) response data
      - "id": (int) 
      - "username": (string) 
      - "roleName": (string) 
      - "roleNo": (string) 
      - "userNo": (string) 
      - "registerDate": (string) 
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8089/open/api/user/info'
    ```

  - Miso HTTP Client:
    ```go
    func ApiUserGetUserInfo(rail miso.Rail) (UserInfoRes, error) {
    	var res miso.GnResp[UserInfoRes]
    	err := miso.NewDynTClient(rail, "/open/api/user/info", "user-vault").
    		Get().
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat UserInfoRes
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: UserInfoRes;
    }

    export interface UserInfoRes {
      id?: number;
      username?: string;
      roleName?: string;
      roleNo?: string;
      userNo?: string;
      registerDate?: string;
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    userGetUserInfo() {
      this.http.get<any>(`/user-vault/open/api/user/info`)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: UserInfoRes = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/user/password/update
  - Description: User update password
  - Bound to Resource: `"basic-user"`
  - JSON Request:
    - "prevPassword": (string) 
    - "newPassword": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/user/password/update' \
      -H 'Content-Type: application/json' \
      -d '{"newPassword":"","prevPassword":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiUserUpdatePassword(rail miso.Rail, req UpdatePasswordReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/user/password/update", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface UpdatePasswordReq {
      prevPassword?: string;
      newPassword?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    userUpdatePassword() {
      let req: UpdatePasswordReq | null = null;
      this.http.post<any>(`/user-vault/open/api/user/password/update`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/token/exchange
  - Description: Exchange token
  - Expected Access Scope: PUBLIC
  - JSON Request:
    - "token": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (string) response data
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/token/exchange' \
      -H 'Content-Type: application/json' \
      -d '{"token":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ExchangeTokenEp(rail miso.Rail, req ExchangeTokenReq) (string, error) {
    	var res miso.GnResp[string]
    	err := miso.NewDynTClient(rail, "/open/api/token/exchange", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return "", err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ExchangeTokenReq {
      token?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: string;                 // response data
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    exchangeTokenEp() {
      let req: ExchangeTokenReq | null = null;
      this.http.post<any>(`/user-vault/open/api/token/exchange`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: string = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- GET /open/api/token/user
  - Description: Get user info by token. This endpoint is expected to be accessible publicly
  - Expected Access Scope: PUBLIC
  - Query Parameter:
    - "token": jwt token
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (UserInfoBrief) response data
      - "id": (int) 
      - "username": (string) 
      - "roleName": (string) 
      - "roleNo": (string) 
      - "userNo": (string) 
      - "registerDate": (string) 
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8089/open/api/token/user?token='
    ```

  - Miso HTTP Client:
    ```go
    func ApiGetTokenUserInfo(rail miso.Rail, token string) (UserInfoBrief, error) {
    	var res miso.GnResp[UserInfoBrief]
    	err := miso.NewDynTClient(rail, "/open/api/token/user", "user-vault").
    		AddQueryParams("token", token).
    		Get().
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat UserInfoBrief
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: UserInfoBrief;
    }

    export interface UserInfoBrief {
      id?: number;
      username?: string;
      roleName?: string;
      roleNo?: string;
      userNo?: string;
      registerDate?: string;
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    getTokenUserInfo() {
      let token: any | null = null;
      this.http.get<any>(`/user-vault/open/api/token/user?token=${token}`)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: UserInfoBrief = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/access/history
  - Description: User list access logs
  - Bound to Resource: `"basic-user"`
  - JSON Request:
    - "paging": (Paging) 
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (PageRes[github.com/curtisnewbie/user-vault/internal/vault.ListedAccessLog]) response data
      - "paging": (Paging) pagination parameters
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
      - "payload": ([]vault.ListedAccessLog) payload values in current page
        - "id": (int) 
        - "userAgent": (string) 
        - "ipAddress": (string) 
        - "username": (string) 
        - "url": (string) 
        - "accessTime": (int64) 
        - "success": (bool) 
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/access/history' \
      -H 'Content-Type: application/json' \
      -d '{"paging":{"limit":0,"page":0,"total":0}}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiUserListAccessHistory(rail miso.Rail, req ListAccessLogReq) (PageRes, error) {
    	var res miso.GnResp[PageRes]
    	err := miso.NewDynTClient(rail, "/open/api/access/history", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat PageRes
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ListAccessLogReq {
      paging?: Paging;
    }

    export interface Paging {
      limit?: number;                // page limit
      page?: number;                 // page number, 1-based
      total?: number;                // total count
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: PageRes;
    }

    export interface PageRes {
      paging?: Paging;
      payload?: ListedAccessLog[];
    }

    export interface Paging {
      limit?: number;                // page limit
      page?: number;                 // page number, 1-based
      total?: number;                // total count
    }

    export interface ListedAccessLog {
      id?: number;
      userAgent?: string;
      ipAddress?: string;
      username?: string;
      url?: string;
      accessTime?: number;
      success?: boolean;
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    userListAccessHistory() {
      let req: ListAccessLogReq | null = null;
      this.http.post<any>(`/user-vault/open/api/access/history`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: PageRes = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/user/key/generate
  - Description: User generate user key
  - Bound to Resource: `"basic-user"`
  - JSON Request:
    - "password": (string) 
    - "keyName": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/user/key/generate' \
      -H 'Content-Type: application/json' \
      -d '{"keyName":"","password":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiUserGenUserKey(rail miso.Rail, req GenUserKeyReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/user/key/generate", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface GenUserKeyReq {
      password?: string;
      keyName?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    userGenUserKey() {
      let req: GenUserKeyReq | null = null;
      this.http.post<any>(`/user-vault/open/api/user/key/generate`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/user/key/list
  - Description: User list user keys
  - Bound to Resource: `"basic-user"`
  - JSON Request:
    - "paging": (Paging) 
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
    - "name": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (PageRes[github.com/curtisnewbie/user-vault/internal/vault.ListedUserKey]) response data
      - "paging": (Paging) pagination parameters
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
      - "payload": ([]vault.ListedUserKey) payload values in current page
        - "id": (int) 
        - "secretKey": (string) 
        - "name": (string) 
        - "expirationTime": (int64) 
        - "createTime": (int64) 
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/user/key/list' \
      -H 'Content-Type: application/json' \
      -d '{"name":"","paging":{"limit":0,"page":0,"total":0}}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiUserListUserKeys(rail miso.Rail, req ListUserKeysReq) (PageRes, error) {
    	var res miso.GnResp[PageRes]
    	err := miso.NewDynTClient(rail, "/open/api/user/key/list", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat PageRes
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ListUserKeysReq {
      paging?: Paging;
      name?: string;
    }

    export interface Paging {
      limit?: number;                // page limit
      page?: number;                 // page number, 1-based
      total?: number;                // total count
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: PageRes;
    }

    export interface PageRes {
      paging?: Paging;
      payload?: ListedUserKey[];
    }

    export interface Paging {
      limit?: number;                // page limit
      page?: number;                 // page number, 1-based
      total?: number;                // total count
    }

    export interface ListedUserKey {
      id?: number;
      secretKey?: string;
      name?: string;
      expirationTime?: number;
      createTime?: number;
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    userListUserKeys() {
      let req: ListUserKeysReq | null = null;
      this.http.post<any>(`/user-vault/open/api/user/key/list`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: PageRes = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/user/key/delete
  - Description: User delete user key
  - Bound to Resource: `"basic-user"`
  - JSON Request:
    - "userKeyId": (int) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/user/key/delete' \
      -H 'Content-Type: application/json' \
      -d '{"userKeyId":0}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiUserDeleteUserKey(rail miso.Rail, req DeleteUserKeyReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/user/key/delete", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface DeleteUserKeyReq {
      userKeyId?: number;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    userDeleteUserKey() {
      let req: DeleteUserKeyReq | null = null;
      this.http.post<any>(`/user-vault/open/api/user/key/delete`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/resource/add
  - Description: Admin add resource
  - Bound to Resource: `"manage-resources"`
  - JSON Request:
    - "name": (string) 
    - "code": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/resource/add' \
      -H 'Content-Type: application/json' \
      -d '{"code":"","name":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiAdminAddResource(rail miso.Rail, req CreateResReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/resource/add", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface CreateResReq {
      name?: string;
      code?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    adminAddResource() {
      let req: CreateResReq | null = null;
      this.http.post<any>(`/user-vault/open/api/resource/add`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/resource/remove
  - Description: Admin remove resource
  - Bound to Resource: `"manage-resources"`
  - JSON Request:
    - "resCode": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/resource/remove' \
      -H 'Content-Type: application/json' \
      -d '{"resCode":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiAdminRemoveResource(rail miso.Rail, req DeleteResourceReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/resource/remove", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface DeleteResourceReq {
      resCode?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    adminRemoveResource() {
      let req: DeleteResourceReq | null = null;
      this.http.post<any>(`/user-vault/open/api/resource/remove`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- GET /open/api/resource/brief/candidates
  - Description: List all resource candidates for role
  - Bound to Resource: `"manage-resources"`
  - Query Parameter:
    - "roleNo": Role No
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": ([]vault.ResBrief) response data
      - "code": (string) 
      - "name": (string) 
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8089/open/api/resource/brief/candidates?roleNo='
    ```

  - Miso HTTP Client:
    ```go
    func ApiListResCandidates(rail miso.Rail, roleNo string) ([]ResBrief, error) {
    	var res miso.GnResp[[]ResBrief]
    	err := miso.NewDynTClient(rail, "/open/api/resource/brief/candidates", "user-vault").
    		AddQueryParams("roleNo", roleNo).
    		Get().
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat []ResBrief
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: ResBrief[];
    }

    export interface ResBrief {
      code?: string;
      name?: string;
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    listResCandidates() {
      let roleNo: any | null = null;
      this.http.get<any>(`/user-vault/open/api/resource/brief/candidates?roleNo=${roleNo}`)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: ResBrief[] = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/resource/list
  - Description: Admin list resources
  - Bound to Resource: `"manage-resources"`
  - JSON Request:
    - "paging": (Paging) 
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (ListResResp) response data
      - "paging": (Paging) 
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
      - "payload": ([]vault.WRes) 
        - "id": (int) 
        - "code": (string) 
        - "name": (string) 
        - "createTime": (int64) 
        - "createBy": (string) 
        - "updateTime": (int64) 
        - "updateBy": (string) 
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/resource/list' \
      -H 'Content-Type: application/json' \
      -d '{"paging":{"limit":0,"page":0,"total":0}}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiAdminListRes(rail miso.Rail, req ListResReq) (ListResResp, error) {
    	var res miso.GnResp[ListResResp]
    	err := miso.NewDynTClient(rail, "/open/api/resource/list", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat ListResResp
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ListResReq {
      paging?: Paging;
    }

    export interface Paging {
      limit?: number;                // page limit
      page?: number;                 // page number, 1-based
      total?: number;                // total count
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: ListResResp;
    }

    export interface ListResResp {
      paging?: Paging;
      payload?: WRes[];
    }

    export interface Paging {
      limit?: number;                // page limit
      page?: number;                 // page number, 1-based
      total?: number;                // total count
    }

    export interface WRes {
      id?: number;
      code?: string;
      name?: string;
      createTime?: number;
      createBy?: string;
      updateTime?: number;
      updateBy?: string;
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    adminListRes() {
      let req: ListResReq | null = null;
      this.http.post<any>(`/user-vault/open/api/resource/list`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: ListResResp = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- GET /open/api/resource/brief/user
  - Description: List resources that are accessible to current user
  - Expected Access Scope: PUBLIC
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": ([]vault.ResBrief) response data
      - "code": (string) 
      - "name": (string) 
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8089/open/api/resource/brief/user'
    ```

  - Miso HTTP Client:
    ```go
    func ApiListUserAccessibleRes(rail miso.Rail) ([]ResBrief, error) {
    	var res miso.GnResp[[]ResBrief]
    	err := miso.NewDynTClient(rail, "/open/api/resource/brief/user", "user-vault").
    		Get().
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat []ResBrief
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: ResBrief[];
    }

    export interface ResBrief {
      code?: string;
      name?: string;
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    listUserAccessibleRes() {
      this.http.get<any>(`/user-vault/open/api/resource/brief/user`)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: ResBrief[] = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- GET /open/api/resource/brief/all
  - Description: List all resource brief info
  - Expected Access Scope: PUBLIC
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": ([]vault.ResBrief) response data
      - "code": (string) 
      - "name": (string) 
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8089/open/api/resource/brief/all'
    ```

  - Miso HTTP Client:
    ```go
    func ApiListAllResBrief(rail miso.Rail) ([]ResBrief, error) {
    	var res miso.GnResp[[]ResBrief]
    	err := miso.NewDynTClient(rail, "/open/api/resource/brief/all", "user-vault").
    		Get().
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat []ResBrief
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: ResBrief[];
    }

    export interface ResBrief {
      code?: string;
      name?: string;
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    listAllResBrief() {
      this.http.get<any>(`/user-vault/open/api/resource/brief/all`)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: ResBrief[] = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/role/resource/add
  - Description: Admin add resource to role
  - Bound to Resource: `"manage-resources"`
  - JSON Request:
    - "roleNo": (string) 
    - "resCode": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/role/resource/add' \
      -H 'Content-Type: application/json' \
      -d '{"resCode":"","roleNo":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiAdminBindRoleRes(rail miso.Rail, req AddRoleResReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/role/resource/add", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface AddRoleResReq {
      roleNo?: string;
      resCode?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    adminBindRoleRes() {
      let req: AddRoleResReq | null = null;
      this.http.post<any>(`/user-vault/open/api/role/resource/add`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/role/resource/remove
  - Description: Admin remove resource from role
  - Bound to Resource: `"manage-resources"`
  - JSON Request:
    - "roleNo": (string) 
    - "resCode": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/role/resource/remove' \
      -H 'Content-Type: application/json' \
      -d '{"resCode":"","roleNo":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiAdminUnbindRoleRes(rail miso.Rail, req RemoveRoleResReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/role/resource/remove", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface RemoveRoleResReq {
      roleNo?: string;
      resCode?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    adminUnbindRoleRes() {
      let req: RemoveRoleResReq | null = null;
      this.http.post<any>(`/user-vault/open/api/role/resource/remove`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/role/add
  - Description: Admin add role
  - Bound to Resource: `"manage-resources"`
  - JSON Request:
    - "name": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/role/add' \
      -H 'Content-Type: application/json' \
      -d '{"name":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiAdminAddRole(rail miso.Rail, req AddRoleReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/role/add", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface AddRoleReq {
      name?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    adminAddRole() {
      let req: AddRoleReq | null = null;
      this.http.post<any>(`/user-vault/open/api/role/add`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/role/list
  - Description: Admin list roles
  - Bound to Resource: `"manage-resources"`
  - JSON Request:
    - "paging": (Paging) 
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (ListRoleResp) response data
      - "payload": ([]vault.WRole) 
        - "id": (int) 
        - "roleNo": (string) 
        - "name": (string) 
        - "createTime": (int64) 
        - "createBy": (string) 
        - "updateTime": (int64) 
        - "updateBy": (string) 
      - "paging": (Paging) 
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/role/list' \
      -H 'Content-Type: application/json' \
      -d '{"paging":{"limit":0,"page":0,"total":0}}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiAdminListRoles(rail miso.Rail, req ListRoleReq) (ListRoleResp, error) {
    	var res miso.GnResp[ListRoleResp]
    	err := miso.NewDynTClient(rail, "/open/api/role/list", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat ListRoleResp
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ListRoleReq {
      paging?: Paging;
    }

    export interface Paging {
      limit?: number;                // page limit
      page?: number;                 // page number, 1-based
      total?: number;                // total count
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: ListRoleResp;
    }

    export interface ListRoleResp {
      payload?: WRole[];
      paging?: Paging;
    }

    export interface WRole {
      id?: number;
      roleNo?: string;
      name?: string;
      createTime?: number;
      createBy?: string;
      updateTime?: number;
      updateBy?: string;
    }

    export interface Paging {
      limit?: number;                // page limit
      page?: number;                 // page number, 1-based
      total?: number;                // total count
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    adminListRoles() {
      let req: ListRoleReq | null = null;
      this.http.post<any>(`/user-vault/open/api/role/list`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: ListRoleResp = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- GET /open/api/role/brief/all
  - Description: Admin list role brief info
  - Bound to Resource: `"manage-resources"`
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": ([]vault.RoleBrief) response data
      - "roleNo": (string) 
      - "name": (string) 
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8089/open/api/role/brief/all'
    ```

  - Miso HTTP Client:
    ```go
    func ApiAdminListRoleBriefs(rail miso.Rail) ([]RoleBrief, error) {
    	var res miso.GnResp[[]RoleBrief]
    	err := miso.NewDynTClient(rail, "/open/api/role/brief/all", "user-vault").
    		Get().
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat []RoleBrief
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: RoleBrief[];
    }

    export interface RoleBrief {
      roleNo?: string;
      name?: string;
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    adminListRoleBriefs() {
      this.http.get<any>(`/user-vault/open/api/role/brief/all`)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: RoleBrief[] = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/role/resource/list
  - Description: Admin list resources of role
  - Bound to Resource: `"manage-resources"`
  - JSON Request:
    - "paging": (Paging) 
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
    - "roleNo": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (ListRoleResResp) response data
      - "paging": (Paging) 
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
      - "payload": ([]vault.ListedRoleRes) 
        - "id": (int) 
        - "resCode": (string) 
        - "resName": (string) 
        - "createTime": (int64) 
        - "createBy": (string) 
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/role/resource/list' \
      -H 'Content-Type: application/json' \
      -d '{"paging":{"limit":0,"page":0,"total":0},"roleNo":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiAdminListRoleRes(rail miso.Rail, req ListRoleResReq) (ListRoleResResp, error) {
    	var res miso.GnResp[ListRoleResResp]
    	err := miso.NewDynTClient(rail, "/open/api/role/resource/list", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat ListRoleResResp
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ListRoleResReq {
      paging?: Paging;
      roleNo?: string;
    }

    export interface Paging {
      limit?: number;                // page limit
      page?: number;                 // page number, 1-based
      total?: number;                // total count
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: ListRoleResResp;
    }

    export interface ListRoleResResp {
      paging?: Paging;
      payload?: ListedRoleRes[];
    }

    export interface Paging {
      limit?: number;                // page limit
      page?: number;                 // page number, 1-based
      total?: number;                // total count
    }

    export interface ListedRoleRes {
      id?: number;
      resCode?: string;
      resName?: string;
      createTime?: number;
      createBy?: string;
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    adminListRoleRes() {
      let req: ListRoleResReq | null = null;
      this.http.post<any>(`/user-vault/open/api/role/resource/list`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: ListRoleResResp = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/role/info
  - Description: Get role info
  - Expected Access Scope: PUBLIC
  - JSON Request:
    - "roleNo": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (RoleInfoResp) response data
      - "roleNo": (string) 
      - "name": (string) 
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/role/info' \
      -H 'Content-Type: application/json' \
      -d '{"roleNo":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiGetRoleInfo(rail miso.Rail, req RoleInfoReq) (RoleInfoResp, error) {
    	var res miso.GnResp[RoleInfoResp]
    	err := miso.NewDynTClient(rail, "/open/api/role/info", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat RoleInfoResp
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface RoleInfoReq {
      roleNo?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: RoleInfoResp;
    }

    export interface RoleInfoResp {
      roleNo?: string;
      name?: string;
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    getRoleInfo() {
      let req: RoleInfoReq | null = null;
      this.http.post<any>(`/user-vault/open/api/role/info`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: RoleInfoResp = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/path/list
  - Description: Admin list paths
  - Bound to Resource: `"manage-resources"`
  - JSON Request:
    - "resCode": (string) 
    - "pgroup": (string) 
    - "url": (string) 
    - "ptype": (string) path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible
    - "paging": (Paging) 
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (ListPathResp) response data
      - "paging": (Paging) 
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
      - "payload": ([]vault.WPath) 
        - "id": (int) 
        - "pgroup": (string) 
        - "pathNo": (string) 
        - "method": (string) 
        - "desc": (string) 
        - "url": (string) 
        - "ptype": (string) path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible
        - "createTime": (int64) 
        - "createBy": (string) 
        - "updateTime": (int64) 
        - "updateBy": (string) 
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/path/list' \
      -H 'Content-Type: application/json' \
      -d '{"paging":{"limit":0,"page":0,"total":0},"pgroup":"","ptype":"","resCode":"","url":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiAdminListPaths(rail miso.Rail, req ListPathReq) (ListPathResp, error) {
    	var res miso.GnResp[ListPathResp]
    	err := miso.NewDynTClient(rail, "/open/api/path/list", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat ListPathResp
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ListPathReq {
      resCode?: string;
      pgroup?: string;
      url?: string;
      ptype?: string;                // path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible
      paging?: Paging;
    }

    export interface Paging {
      limit?: number;                // page limit
      page?: number;                 // page number, 1-based
      total?: number;                // total count
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: ListPathResp;
    }

    export interface ListPathResp {
      paging?: Paging;
      payload?: WPath[];
    }

    export interface Paging {
      limit?: number;                // page limit
      page?: number;                 // page number, 1-based
      total?: number;                // total count
    }

    export interface WPath {
      id?: number;
      pgroup?: string;
      pathNo?: string;
      method?: string;
      desc?: string;
      url?: string;
      ptype?: string;                // path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible
      createTime?: number;
      createBy?: string;
      updateTime?: number;
      updateBy?: string;
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    adminListPaths() {
      let req: ListPathReq | null = null;
      this.http.post<any>(`/user-vault/open/api/path/list`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: ListPathResp = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/path/resource/bind
  - Description: Admin bind resource to path
  - Bound to Resource: `"manage-resources"`
  - JSON Request:
    - "pathNo": (string) 
    - "resCode": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/path/resource/bind' \
      -H 'Content-Type: application/json' \
      -d '{"pathNo":"","resCode":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiAdminBindResPath(rail miso.Rail, req BindPathResReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/path/resource/bind", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface BindPathResReq {
      pathNo?: string;
      resCode?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    adminBindResPath() {
      let req: BindPathResReq | null = null;
      this.http.post<any>(`/user-vault/open/api/path/resource/bind`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/path/resource/unbind
  - Description: Admin unbind resource and path
  - Bound to Resource: `"manage-resources"`
  - JSON Request:
    - "pathNo": (string) 
    - "resCode": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/path/resource/unbind' \
      -H 'Content-Type: application/json' \
      -d '{"pathNo":"","resCode":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiAdminUnbindResPath(rail miso.Rail, req UnbindPathResReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/path/resource/unbind", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface UnbindPathResReq {
      pathNo?: string;
      resCode?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    adminUnbindResPath() {
      let req: UnbindPathResReq | null = null;
      this.http.post<any>(`/user-vault/open/api/path/resource/unbind`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/path/delete
  - Description: Admin delete path
  - Bound to Resource: `"manage-resources"`
  - JSON Request:
    - "pathNo": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/path/delete' \
      -H 'Content-Type: application/json' \
      -d '{"pathNo":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiAdminDeletePath(rail miso.Rail, req DeletePathReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/path/delete", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface DeletePathReq {
      pathNo?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    adminDeletePath() {
      let req: DeletePathReq | null = null;
      this.http.post<any>(`/user-vault/open/api/path/delete`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/path/update
  - Description: Admin update path
  - Bound to Resource: `"manage-resources"`
  - JSON Request:
    - "type": (string) path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible
    - "pathNo": (string) 
    - "group": (string) 
    - "resCode": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/path/update' \
      -H 'Content-Type: application/json' \
      -d '{"group":"","pathNo":"","resCode":"","type":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiAdminUpdatePath(rail miso.Rail, req UpdatePathReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/path/update", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface UpdatePathReq {
      type?: string;                 // path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible
      pathNo?: string;
      group?: string;
      resCode?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    adminUpdatePath() {
      let req: UpdatePathReq | null = null;
      this.http.post<any>(`/user-vault/open/api/path/update`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /remote/user/info
  - Description: Fetch user info
  - JSON Request:
    - "userId": (*int) 
    - "userNo": (*string) 
    - "username": (*string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (UserInfo) response data
      - "id": (int) 
      - "username": (string) 
      - "roleName": (string) 
      - "roleNo": (string) 
      - "userNo": (string) 
      - "reviewStatus": (string) 
      - "isDisabled": (int) 
      - "createTime": (int64) 
      - "createBy": (string) 
      - "updateTime": (int64) 
      - "updateBy": (string) 
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/remote/user/info' \
      -H 'Content-Type: application/json' \
      -d '{"userId":0,"userNo":"","username":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiFetchUserInfo(rail miso.Rail, req FindUserReq) (UserInfo, error) {
    	var res miso.GnResp[UserInfo]
    	err := miso.NewDynTClient(rail, "/remote/user/info", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat UserInfo
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface FindUserReq {
      userId?: number;
      userNo?: string;
      username?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: UserInfo;
    }

    export interface UserInfo {
      id?: number;
      username?: string;
      roleName?: string;
      roleNo?: string;
      userNo?: string;
      reviewStatus?: string;
      isDisabled?: number;
      createTime?: number;
      createBy?: string;
      updateTime?: number;
      updateBy?: string;
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    fetchUserInfo() {
      let req: FindUserReq | null = null;
      this.http.post<any>(`/user-vault/remote/user/info`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: UserInfo = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- GET /remote/user/id
  - Description: Fetch id of user with the username
  - Query Parameter:
    - "username": Username
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (int) response data
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8089/remote/user/id?username='
    ```

  - Miso HTTP Client:
    ```go
    func ApiFetchUserIdByName(rail miso.Rail, username string) (int, error) {
    	var res miso.GnResp[int]
    	err := miso.NewDynTClient(rail, "/remote/user/id", "user-vault").
    		AddQueryParams("username", username).
    		Get().
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return 0, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: number;                 // response data
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    fetchUserIdByName() {
      let username: any | null = null;
      this.http.get<any>(`/user-vault/remote/user/id?username=${username}`)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: number = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /remote/user/userno/username
  - Description: Fetch usernames of users with the userNos
  - JSON Request:
    - "userNos": ([]string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (FetchUsernamesRes) response data
      - "userNoToUsername": (map[string]string) 
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/remote/user/userno/username' \
      -H 'Content-Type: application/json' \
      -d '{"userNos":[]}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiFetchUsernamesByNosEp(rail miso.Rail, req FetchNameByUserNoReq) (FetchUsernamesRes, error) {
    	var res miso.GnResp[FetchUsernamesRes]
    	err := miso.NewDynTClient(rail, "/remote/user/userno/username", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat FetchUsernamesRes
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface FetchNameByUserNoReq {
      userNos?: string[];
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: FetchUsernamesRes;
    }

    export interface FetchUsernamesRes {
      userNoToUsername?: map[string]string;
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    fetchUsernamesByNosEp() {
      let req: FetchNameByUserNoReq | null = null;
      this.http.post<any>(`/user-vault/remote/user/userno/username`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: FetchUsernamesRes = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /remote/user/list/with-role
  - Description: Fetch users with the role_no
  - JSON Request:
    - "roleNo": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": ([]api.UserInfo) response data
      - "id": (int) 
      - "username": (string) 
      - "roleName": (string) 
      - "roleNo": (string) 
      - "userNo": (string) 
      - "reviewStatus": (string) 
      - "isDisabled": (int) 
      - "createTime": (int64) 
      - "createBy": (string) 
      - "updateTime": (int64) 
      - "updateBy": (string) 
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/remote/user/list/with-role' \
      -H 'Content-Type: application/json' \
      -d '{"roleNo":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiFindUserWithRoleEp(rail miso.Rail, req FetchUsersWithRoleReq) ([]UserInfo, error) {
    	var res miso.GnResp[[]UserInfo]
    	err := miso.NewDynTClient(rail, "/remote/user/list/with-role", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat []UserInfo
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface FetchUsersWithRoleReq {
      roleNo?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: UserInfo[];
    }

    export interface UserInfo {
      id?: number;
      username?: string;
      roleName?: string;
      roleNo?: string;
      userNo?: string;
      reviewStatus?: string;
      isDisabled?: number;
      createTime?: number;
      createBy?: string;
      updateTime?: number;
      updateBy?: string;
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    findUserWithRoleEp() {
      let req: FetchUsersWithRoleReq | null = null;
      this.http.post<any>(`/user-vault/remote/user/list/with-role`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: UserInfo[] = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /remote/user/list/with-resource
  - Description: Fetch users that have access to the resource
  - JSON Request:
    - "resourceCode": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": ([]api.UserInfo) response data
      - "id": (int) 
      - "username": (string) 
      - "roleName": (string) 
      - "roleNo": (string) 
      - "userNo": (string) 
      - "reviewStatus": (string) 
      - "isDisabled": (int) 
      - "createTime": (int64) 
      - "createBy": (string) 
      - "updateTime": (int64) 
      - "updateBy": (string) 
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/remote/user/list/with-resource' \
      -H 'Content-Type: application/json' \
      -d '{"resourceCode":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiFindUserWithResourceEp(rail miso.Rail, req FetchUserWithResourceReq) ([]UserInfo, error) {
    	var res miso.GnResp[[]UserInfo]
    	err := miso.NewDynTClient(rail, "/remote/user/list/with-resource", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat []UserInfo
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface FetchUserWithResourceReq {
      resourceCode?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: UserInfo[];
    }

    export interface UserInfo {
      id?: number;
      username?: string;
      roleName?: string;
      roleNo?: string;
      userNo?: string;
      reviewStatus?: string;
      isDisabled?: number;
      createTime?: number;
      createBy?: string;
      updateTime?: number;
      updateBy?: string;
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    findUserWithResourceEp() {
      let req: FetchUserWithResourceReq | null = null;
      this.http.post<any>(`/user-vault/remote/user/list/with-resource`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: UserInfo[] = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /remote/resource/add
  - Description: Report resource. This endpoint should be used internally by another backend service.
  - JSON Request:
    - "name": (string) 
    - "code": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/remote/resource/add' \
      -H 'Content-Type: application/json' \
      -d '{"code":"","name":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiReportResourceEp(rail miso.Rail, req CreateResReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/remote/resource/add", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface CreateResReq {
      name?: string;
      code?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    reportResourceEp() {
      let req: CreateResReq | null = null;
      this.http.post<any>(`/user-vault/remote/resource/add`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /remote/path/resource/access-test
  - Description: Validate resource access
  - JSON Request:
    - "roleNo": (string) 
    - "url": (string) 
    - "method": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (CheckResAccessResp) response data
      - "valid": (bool) 
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/remote/path/resource/access-test' \
      -H 'Content-Type: application/json' \
      -d '{"method":"","roleNo":"","url":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiCheckResourceAccessEp(rail miso.Rail, req CheckResAccessReq) (CheckResAccessResp, error) {
    	var res miso.GnResp[CheckResAccessResp]
    	err := miso.NewDynTClient(rail, "/remote/path/resource/access-test", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat CheckResAccessResp
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface CheckResAccessReq {
      roleNo?: string;
      url?: string;
      method?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: CheckResAccessResp;
    }

    export interface CheckResAccessResp {
      valid?: boolean;
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    checkResourceAccessEp() {
      let req: CheckResAccessReq | null = null;
      this.http.post<any>(`/user-vault/remote/path/resource/access-test`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: CheckResAccessResp = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /remote/path/add
  - Description: Report endpoint info
  - Bound to Resource: `"manage-resources"`
  - JSON Request:
    - "type": (string) path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible
    - "url": (string) 
    - "group": (string) 
    - "method": (string) 
    - "desc": (string) 
    - "resCode": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/remote/path/add' \
      -H 'Content-Type: application/json' \
      -d '{"desc":"","group":"","method":"","resCode":"","type":"","url":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiReportPath(rail miso.Rail, req CreatePathReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/remote/path/add", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface CreatePathReq {
      type?: string;                 // path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible
      url?: string;
      group?: string;
      method?: string;
      desc?: string;
      resCode?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    reportPath() {
      let req: CreatePathReq | null = null;
      this.http.post<any>(`/user-vault/remote/path/add`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/password/list-site-passwords
  - Description: List site password records
  - Bound to Resource: `"basic-user"`
  - JSON Request:
    - "alias": (string) 
    - "site": (string) 
    - "username": (string) 
    - "paging": (Paging) 
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (PageRes[github.com/curtisnewbie/user-vault/internal/vault.ListSitePasswordRes]) response data
      - "paging": (Paging) pagination parameters
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
      - "payload": ([]vault.ListSitePasswordRes) payload values in current page
        - "recordId": (string) 
        - "site": (string) 
        - "alias": (string) 
        - "username": (string) 
        - "createTime": (int64) 
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/password/list-site-passwords' \
      -H 'Content-Type: application/json' \
      -d '{"alias":"","paging":{"limit":0,"page":0,"total":0},"site":"","username":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiListSitePasswords(rail miso.Rail, req ListSitePasswordReq) (PageRes, error) {
    	var res miso.GnResp[PageRes]
    	err := miso.NewDynTClient(rail, "/open/api/password/list-site-passwords", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat PageRes
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ListSitePasswordReq {
      alias?: string;
      site?: string;
      username?: string;
      paging?: Paging;
    }

    export interface Paging {
      limit?: number;                // page limit
      page?: number;                 // page number, 1-based
      total?: number;                // total count
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: PageRes;
    }

    export interface PageRes {
      paging?: Paging;
      payload?: ListSitePasswordRes[];
    }

    export interface Paging {
      limit?: number;                // page limit
      page?: number;                 // page number, 1-based
      total?: number;                // total count
    }

    export interface ListSitePasswordRes {
      recordId?: string;
      site?: string;
      alias?: string;
      username?: string;
      createTime?: number;
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    listSitePasswords() {
      let req: ListSitePasswordReq | null = null;
      this.http.post<any>(`/user-vault/open/api/password/list-site-passwords`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: PageRes = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/password/add-site-password
  - Description: Add site password record
  - Bound to Resource: `"basic-user"`
  - JSON Request:
    - "site": (string) 
    - "alias": (string) 
    - "username": (string) 
    - "sitePassword": (string) 
    - "loginPassword": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/password/add-site-password' \
      -H 'Content-Type: application/json' \
      -d '{"alias":"","loginPassword":"","site":"","sitePassword":"","username":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiAddSitePassword(rail miso.Rail, req AddSitePasswordReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/password/add-site-password", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface AddSitePasswordReq {
      site?: string;
      alias?: string;
      username?: string;
      sitePassword?: string;
      loginPassword?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    addSitePassword() {
      let req: AddSitePasswordReq | null = null;
      this.http.post<any>(`/user-vault/open/api/password/add-site-password`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/password/remove-site-password
  - Description: Remove site password record
  - Bound to Resource: `"basic-user"`
  - JSON Request:
    - "recordId": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/password/remove-site-password' \
      -H 'Content-Type: application/json' \
      -d '{"recordId":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiRemoveSitePassword(rail miso.Rail, req RemoveSitePasswordRes) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/password/remove-site-password", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface RemoveSitePasswordRes {
      recordId?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    removeSitePassword() {
      let req: RemoveSitePasswordRes | null = null;
      this.http.post<any>(`/user-vault/open/api/password/remove-site-password`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/password/decrypt-site-password
  - Description: Decrypt site password
  - Bound to Resource: `"basic-user"`
  - JSON Request:
    - "loginPassword": (string) 
    - "recordId": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (DecryptSitePasswordRes) response data
      - "decrypted": (string) 
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/password/decrypt-site-password' \
      -H 'Content-Type: application/json' \
      -d '{"loginPassword":"","recordId":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiDecryptSitePassword(rail miso.Rail, req DecryptSitePasswordReq) (DecryptSitePasswordRes, error) {
    	var res miso.GnResp[DecryptSitePasswordRes]
    	err := miso.NewDynTClient(rail, "/open/api/password/decrypt-site-password", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		var dat DecryptSitePasswordRes
    		return dat, err
    	}
    	dat, err := res.Res()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return dat, err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface DecryptSitePasswordReq {
      loginPassword?: string;
      recordId?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
      data?: DecryptSitePasswordRes;
    }

    export interface DecryptSitePasswordRes {
      decrypted?: string;
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    decryptSitePassword() {
      let req: DecryptSitePasswordReq | null = null;
      this.http.post<any>(`/user-vault/open/api/password/decrypt-site-password`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
            let dat: DecryptSitePasswordRes = resp.data;
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/password/edit-site-password
  - Description: Edit site password
  - Bound to Resource: `"basic-user"`
  - JSON Request:
    - "recordId": (string) 
    - "site": (string) 
    - "alias": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/password/edit-site-password' \
      -H 'Content-Type: application/json' \
      -d '{"alias":"","recordId":"","site":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiEditSitePassword(rail miso.Rail, req EditSitePasswordReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/password/edit-site-password", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface EditSitePasswordReq {
      recordId?: string;
      site?: string;
      alias?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    editSitePassword() {
      let req: EditSitePasswordReq | null = null;
      this.http.post<any>(`/user-vault/open/api/password/edit-site-password`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/user/clear-failed-login-attempts
  - Description: Admin clear user's failed login attempts
  - Bound to Resource: `"manage-users"`
  - JSON Request:
    - "userNo": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/user/clear-failed-login-attempts' \
      -H 'Content-Type: application/json' \
      -d '{"userNo":""}'
    ```

  - Miso HTTP Client:
    ```go
    func ApiClearUserFailedLoginAttempts(rail miso.Rail, req ClearUserFailedLoginAttemptsReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/user/clear-failed-login-attempts", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ClearUserFailedLoginAttemptsReq {
      userNo?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    clearUserFailedLoginAttempts() {
      let req: ClearUserFailedLoginAttemptsReq | null = null;
      this.http.post<any>(`/user-vault/open/api/user/clear-failed-login-attempts`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/v1/notification/create
  - Description: Create platform notification
  - Bound to Resource: `"postbox:notification:create"`
  - JSON Request:
    - "title": (string) 
    - "message": (string) 
    - "receiverUserNos": ([]string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/v1/notification/create' \
      -H 'Content-Type: application/json' \
      -d '{"message":"","receiverUserNos":[],"title":""}'
    ```

  - Miso HTTP Client:
    ```go
    func SendCreateNotificationReq(rail miso.Rail, req CreateNotificationReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/v1/notification/create", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface CreateNotificationReq {
      title?: string;
      message?: string;
      receiverUserNos?: string[];
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    sendCreateNotificationReq() {
      let req: CreateNotificationReq | null = null;
      this.http.post<any>(`/user-vault/open/api/v1/notification/create`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/v1/notification/query
  - Description: Query platform notification
  - Bound to Resource: `"postbox:notification:query"`
  - JSON Request:
    - "page": (Paging) 
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
    - "status": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/v1/notification/query' \
      -H 'Content-Type: application/json' \
      -d '{"page":{"limit":0,"page":0,"total":0},"status":""}'
    ```

  - Miso HTTP Client:
    ```go
    func SendQueryNotificationReq(rail miso.Rail, req QueryNotificationReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/v1/notification/query", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface QueryNotificationReq {
      page?: Paging;
      status?: string;
    }

    export interface Paging {
      limit?: number;                // page limit
      page?: number;                 // page number, 1-based
      total?: number;                // total count
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    sendQueryNotificationReq() {
      let req: QueryNotificationReq | null = null;
      this.http.post<any>(`/user-vault/open/api/v1/notification/query`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- GET /open/api/v1/notification/count
  - Description: Count received platform notification
  - Bound to Resource: `"postbox:notification:query"`
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8089/open/api/v1/notification/count'
    ```

  - Miso HTTP Client:
    ```go
    func SendRequest(rail miso.Rail) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/v1/notification/count", "user-vault").
    		Get().
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    sendRequest() {
      this.http.get<any>(`/user-vault/open/api/v1/notification/count`)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/v1/notification/open
  - Description: Record user opened platform notification
  - Bound to Resource: `"postbox:notification:query"`
  - JSON Request:
    - "notifiNo": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/v1/notification/open' \
      -H 'Content-Type: application/json' \
      -d '{"notifiNo":""}'
    ```

  - Miso HTTP Client:
    ```go
    func SendOpenNotificationReq(rail miso.Rail, req OpenNotificationReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/v1/notification/open", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface OpenNotificationReq {
      notifiNo?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    sendOpenNotificationReq() {
      let req: OpenNotificationReq | null = null;
      this.http.post<any>(`/user-vault/open/api/v1/notification/open`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- POST /open/api/v1/notification/open-all
  - Description: Mark all notifications opened
  - Bound to Resource: `"postbox:notification:query"`
  - JSON Request:
    - "notifiNo": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8089/open/api/v1/notification/open-all' \
      -H 'Content-Type: application/json' \
      -d '{"notifiNo":""}'
    ```

  - Miso HTTP Client:
    ```go
    func SendOpenNotificationReq(rail miso.Rail, req OpenNotificationReq) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/v1/notification/open-all", "user-vault").
    		PostJson(req).
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface OpenNotificationReq {
      notifiNo?: string;
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string;            // error code
      msg?: string;                  // message
      error?: boolean;               // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    sendOpenNotificationReq() {
      let req: OpenNotificationReq | null = null;
      this.http.post<any>(`/user-vault/open/api/v1/notification/open-all`, req)
        .subscribe({
          next: (resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 })
              return;
            }
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- GET /open/api/v2/notification/count
  - Description: Count received platform notification using long polling
  - Bound to Resource: `"postbox:notification:query"`
  - Query Parameter:
    - "curr": Current count (used to implement long polling)
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8089/open/api/v2/notification/count?curr='
    ```

  - Miso HTTP Client:
    ```go
    func SendRequest(rail miso.Rail, curr string) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/open/api/v2/notification/count", "user-vault").
    		AddQueryParams("curr", curr).
    		Get().
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    sendRequest() {
      let curr: any | null = null;
      this.http.get<any>(`/user-vault/open/api/v2/notification/count?curr=${curr}`)
        .subscribe({
          next: () => {
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- GET /metrics
  - Description: Collect prometheus metrics information
  - Header Parameter:
    - "Authorization": Basic authorization if enabled
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8089/metrics' \
      -H 'Authorization: '
    ```

  - Miso HTTP Client:
    ```go
    func SendRequest(rail miso.Rail, authorization string) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/metrics", "user-vault").
    		AddHeader("authorization", authorization).
    		Get().
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    sendRequest() {
      let authorization: any | null = null;
      this.http.get<any>(`/user-vault/metrics`,
        {
          headers: {
            "Authorization": authorization
          }
        })
        .subscribe({
          next: () => {
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- GET /debug/pprof
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8089/debug/pprof'
    ```

  - Miso HTTP Client:
    ```go
    func SendRequest(rail miso.Rail) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/debug/pprof", "user-vault").
    		Get().
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    sendRequest() {
      this.http.get<any>(`/user-vault/debug/pprof`)
        .subscribe({
          next: () => {
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- GET /debug/pprof/:name
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8089/debug/pprof/:name'
    ```

  - Miso HTTP Client:
    ```go
    func SendRequest(rail miso.Rail) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/debug/pprof/:name", "user-vault").
    		Get().
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    sendRequest() {
      this.http.get<any>(`/user-vault/debug/pprof/:name`)
        .subscribe({
          next: () => {
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- GET /debug/pprof/cmdline
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8089/debug/pprof/cmdline'
    ```

  - Miso HTTP Client:
    ```go
    func SendRequest(rail miso.Rail) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/debug/pprof/cmdline", "user-vault").
    		Get().
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    sendRequest() {
      this.http.get<any>(`/user-vault/debug/pprof/cmdline`)
        .subscribe({
          next: () => {
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- GET /debug/pprof/profile
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8089/debug/pprof/profile'
    ```

  - Miso HTTP Client:
    ```go
    func SendRequest(rail miso.Rail) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/debug/pprof/profile", "user-vault").
    		Get().
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    sendRequest() {
      this.http.get<any>(`/user-vault/debug/pprof/profile`)
        .subscribe({
          next: () => {
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- GET /debug/pprof/symbol
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8089/debug/pprof/symbol'
    ```

  - Miso HTTP Client:
    ```go
    func SendRequest(rail miso.Rail) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/debug/pprof/symbol", "user-vault").
    		Get().
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    sendRequest() {
      this.http.get<any>(`/user-vault/debug/pprof/symbol`)
        .subscribe({
          next: () => {
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- GET /debug/pprof/trace
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8089/debug/pprof/trace'
    ```

  - Miso HTTP Client:
    ```go
    func SendRequest(rail miso.Rail) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/debug/pprof/trace", "user-vault").
    		Get().
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    sendRequest() {
      this.http.get<any>(`/user-vault/debug/pprof/trace`)
        .subscribe({
          next: () => {
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

- GET /doc/api
  - Description: Serve the generated API documentation webpage
  - Expected Access Scope: PUBLIC
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8089/doc/api'
    ```

  - Miso HTTP Client:
    ```go
    func SendRequest(rail miso.Rail) error {
    	var res miso.GnResp[any]
    	err := miso.NewDynTClient(rail, "/doc/api", "user-vault").
    		Get().
    		Json(&res)
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    		return err
    	}
    	err = res.Err()
    	if err != nil {
    		rail.Errorf("Request failed, %v", err)
    	}
    	return err
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    sendRequest() {
      this.http.get<any>(`/user-vault/doc/api`)
        .subscribe({
          next: () => {
          },
          error: (err) => {
            console.log(err)
            this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
          }
        });
    }
    ```

# Event Pipelines

- CreateNotifiPipeline
  - Description: Pipeline that creates notifications to the specified list of users
  - RabbitMQ Queue: `pieline.user-vault.create-notifi`
  - RabbitMQ Exchange: `pieline.user-vault.create-notifi`
  - RabbitMQ RoutingKey: `#`
  - Event Payload:
    - "title": (string) notification title
    - "message": (string) notification content
    - "receiverUserNos": ([]string) user_no of receivers

- CreateNotifiByAccessPipeline
  - Description: Pipeline that creates notifications to users who have access to the specified resource
  - RabbitMQ Queue: `pieline.user-vault.create-notifi.by-access`
  - RabbitMQ Exchange: `pieline.user-vault.create-notifi.by-access`
  - RabbitMQ RoutingKey: `#`
  - Event Payload:
    - "title": (string) notification title
    - "message": (string) notification content
    - "resCode": (string) resource code
