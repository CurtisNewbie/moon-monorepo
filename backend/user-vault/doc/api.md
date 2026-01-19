# API Endpoints

## Contents

- [POST /open/api/user/login](#post-openapiuserlogin)
- [POST /open/api/user/register/request](#post-openapiuserregisterrequest)
- [POST /open/api/user/add](#post-openapiuseradd)
- [POST /open/api/user/list](#post-openapiuserlist)
- [POST /open/api/user/info/update](#post-openapiuserinfoupdate)
- [POST /open/api/user/registration/review](#post-openapiuserregistrationreview)
- [GET /open/api/user/info](#get-openapiuserinfo)
- [POST /open/api/user/password/update](#post-openapiuserpasswordupdate)
- [POST /open/api/token/exchange](#post-openapitokenexchange)
- [GET /open/api/token/user](#get-openapitokenuser)
- [POST /open/api/access/history](#post-openapiaccesshistory)
- [POST /open/api/user/key/generate](#post-openapiuserkeygenerate)
- [POST /open/api/user/key/list](#post-openapiuserkeylist)
- [POST /open/api/user/key/delete](#post-openapiuserkeydelete)
- [POST /open/api/resource/add](#post-openapiresourceadd)
- [POST /open/api/resource/remove](#post-openapiresourceremove)
- [GET /open/api/resource/brief/candidates](#get-openapiresourcebriefcandidates)
- [POST /open/api/resource/list](#post-openapiresourcelist)
- [GET /open/api/resource/brief/user](#get-openapiresourcebriefuser)
- [GET /open/api/resource/brief/all](#get-openapiresourcebriefall)
- [POST /open/api/role/resource/add](#post-openapiroleresourceadd)
- [POST /open/api/role/resource/remove](#post-openapiroleresourceremove)
- [POST /open/api/role/add](#post-openapiroleadd)
- [POST /open/api/role/list](#post-openapirolelist)
- [GET /open/api/role/brief/all](#get-openapirolebriefall)
- [POST /open/api/role/resource/list](#post-openapiroleresourcelist)
- [POST /open/api/role/info](#post-openapiroleinfo)
- [POST /open/api/path/list](#post-openapipathlist)
- [POST /open/api/path/resource/bind](#post-openapipathresourcebind)
- [POST /open/api/path/resource/unbind](#post-openapipathresourceunbind)
- [POST /open/api/path/delete](#post-openapipathdelete)
- [POST /open/api/path/update](#post-openapipathupdate)
- [POST /remote/user/info](#post-remoteuserinfo)
- [POST /internal/v1/user/info/common](#post-internalv1userinfocommon)
- [GET /remote/user/id](#get-remoteuserid)
- [POST /remote/user/userno/username](#post-remoteuserusernousername)
- [POST /remote/user/list/with-role](#post-remoteuserlistwith-role)
- [POST /remote/user/list/with-resource](#post-remoteuserlistwith-resource)
- [POST /remote/resource/add](#post-remoteresourceadd)
- [POST /remote/path/resource/access-test](#post-remotepathresourceaccess-test)
- [POST /remote/path/add](#post-remotepathadd)
- [POST /open/api/password/list-site-passwords](#post-openapipasswordlist-site-passwords)
- [POST /open/api/password/add-site-password](#post-openapipasswordadd-site-password)
- [POST /open/api/password/remove-site-password](#post-openapipasswordremove-site-password)
- [POST /open/api/password/decrypt-site-password](#post-openapipassworddecrypt-site-password)
- [POST /open/api/password/edit-site-password](#post-openapipasswordedit-site-password)
- [POST /open/api/user/clear-failed-login-attempts](#post-openapiuserclear-failed-login-attempts)
- [POST /open/api/note/list-notes](#post-openapinotelist-notes)
- [POST /open/api/note/save-note](#post-openapinotesave-note)
- [POST /open/api/note/update-note](#post-openapinoteupdate-note)
- [POST /open/api/note/delete-note](#post-openapinotedelete-note)
- [POST /open/api/v1/notification/create](#post-openapiv1notificationcreate)
- [POST /open/api/v1/notification/query](#post-openapiv1notificationquery)
- [GET /open/api/v1/notification/count](#get-openapiv1notificationcount)
- [POST /open/api/v1/notification/open](#post-openapiv1notificationopen)
- [POST /open/api/v1/notification/open-all](#post-openapiv1notificationopen-all)
- [GET /open/api/v2/notification/count](#get-openapiv2notificationcount)
- [GET /debug/trace/recorder/run](#get-debugtracerecorderrun)
- [GET /debug/trace/recorder/snapshot](#get-debugtracerecordersnapshot)
- [GET /debug/trace/recorder/stop](#get-debugtracerecorderstop)
- [POST /debug/task/disable-workers](#post-debugtaskdisable-workers)
- [POST /debug/task/enable-workers](#post-debugtaskenable-workers)

## POST /open/api/user/login

- Description: User Login using password, a JWT token is generated and returned
- Expected Access Scope: PUBLIC
- Header Parameter:
  - "x-forwarded-for": 
  - "user-agent": 
- JSON Request:
    - "username": (string) Required.
    - "password": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type LoginReq struct {
  	Username string `json:"username"` // Required.
  	Password string `json:"password"` // Required.
  }

  // User Login using password, a JWT token is generated and returned
  func ApiUserLogin(rail miso.Rail, req LoginReq, xForwardedFor string, userAgent string) (string, error) {
  	var res miso.GnResp[string]
  	err := miso.NewDynClient(rail, "/open/api/user/login", "user-vault").
  		AddHeader("xForwardedFor", xForwardedFor).
  		AddHeader("userAgent", userAgent).
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return "", err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface LoginReq {
    username?: string;             // Required.
    password?: string;             // Required.
  }

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

## POST /open/api/user/register/request

- Description: User request registration, approval needed
- Expected Access Scope: PUBLIC
- JSON Request:
    - "username": (string) Required.
    - "password": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type RegisterReq struct {
  	Username string `json:"username"` // Required.
  	Password string `json:"password"` // Required.
  }

  // User request registration, approval needed
  func ApiUserRegister(rail miso.Rail, req RegisterReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/user/register/request", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface RegisterReq {
    username?: string;             // Required.
    password?: string;             // Required.
  }

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

## POST /open/api/user/add

- Description: Admin create new user
- Bound to Resource: `"manage-users"`
- JSON Request:
    - "username": (string) Required.
    - "password": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type AddUserParam struct {
  	Username string `json:"username"` // Required.
  	Password string `json:"password"` // Required.
  	RoleNo string `json:"roleNo"`
  }

  // Admin create new user
  func ApiAdminAddUser(rail miso.Rail, req AddUserParam) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/user/add", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface AddUserParam {
    username?: string;             // Required.
    password?: string;             // Required.
    roleNo?: string;
  }

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

## POST /open/api/user/list

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
    - "data": (PageRes[github.com/curtisnewbie/user-vault/internal/vault.UserInfo]) response data
      - "paging": (Paging) pagination parameters
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
      - "payload": ([]vault.UserInfo) payload values in current page
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListUserReq struct {
  	Username *string `json:"username"`
  	RoleNo *string `json:"roleNo"`
  	IsDisabled *int `json:"isDisabled"`
  	Paging miso.Paging `json:"paging"`
  }


  type UserInfo struct {
  	Id int `json:"id"`
  	Username string `json:"username"`
  	RoleName string `json:"roleName"`
  	RoleNo string `json:"roleNo"`
  	UserNo string `json:"userNo"`
  	ReviewStatus string `json:"reviewStatus"`
  	IsDisabled int `json:"isDisabled"`
  	CreateTime atom.Time `json:"createTime"`
  	CreateBy string `json:"createBy"`
  	UpdateTime atom.Time `json:"updateTime"`
  	UpdateBy string `json:"updateBy"`
  }

  // Admin list users
  func ApiAdminListUsers(rail miso.Rail, req ListUserReq) (miso.PageRes[UserInfo], error) {
  	var res miso.GnResp[miso.PageRes[UserInfo]]
  	err := miso.NewDynClient(rail, "/open/api/user/list", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		var dat miso.PageRes[UserInfo]
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
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

## POST /open/api/user/info/update

- Description: Admin update user info
- Bound to Resource: `"manage-users"`
- JSON Request:
    - "userNo": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type AdminUpdateUserReq struct {
  	UserNo string `json:"userNo"`  // Required.
  	RoleNo string `json:"roleNo"`
  	IsDisabled int `json:"isDisabled"`
  }

  // Admin update user info
  func ApiAdminUpdateUser(rail miso.Rail, req AdminUpdateUserReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/user/info/update", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface AdminUpdateUserReq {
    userNo?: string;               // Required.
    roleNo?: string;
    isDisabled?: number;
  }

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

## POST /open/api/user/registration/review

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type AdminReviewUserReq struct {
  	UserId int `json:"userId"`
  	ReviewStatus string `json:"reviewStatus"`
  }

  // Admin review user registration
  func ApiAdminReviewUser(rail miso.Rail, req AdminReviewUserReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/user/registration/review", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface AdminReviewUserReq {
    userId?: number;
    reviewStatus?: string;
  }

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

## GET /open/api/user/info

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type UserInfoRes struct {
  	Id int `json:"id"`
  	Username string `json:"username"`
  	RoleName string `json:"roleName"`
  	RoleNo string `json:"roleNo"`
  	UserNo string `json:"userNo"`
  	RegisterDate string `json:"registerDate"`
  }

  // User get user info
  func ApiUserGetUserInfo(rail miso.Rail) (UserInfoRes, error) {
  	var res miso.GnResp[UserInfoRes]
  	err := miso.NewDynClient(rail, "/open/api/user/info", "user-vault").
  		Get().
  		Json(&res)
  	if err != nil {
  		var dat UserInfoRes
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
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

## POST /open/api/user/password/update

- Description: User update password
- Bound to Resource: `"basic-user"`
- JSON Request:
    - "prevPassword": (string) Required.
    - "newPassword": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type UpdatePasswordReq struct {
  	PrevPassword string `json:"prevPassword"` // Required.
  	NewPassword string `json:"newPassword"` // Required.
  }

  // User update password
  func ApiUserUpdatePassword(rail miso.Rail, req UpdatePasswordReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/user/password/update", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface UpdatePasswordReq {
    prevPassword?: string;         // Required.
    newPassword?: string;          // Required.
  }

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

## POST /open/api/token/exchange

- Description: Exchange token
- Expected Access Scope: PUBLIC
- JSON Request:
    - "token": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ExchangeTokenReq struct {
  	Token string `json:"token"`    // Required.
  }

  // Exchange token
  func ExchangeTokenEp(rail miso.Rail, req ExchangeTokenReq) (string, error) {
  	var res miso.GnResp[string]
  	err := miso.NewDynClient(rail, "/open/api/token/exchange", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return "", err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface ExchangeTokenReq {
    token?: string;                // Required.
  }

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

## GET /open/api/token/user

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type UserInfoBrief struct {
  	Id int `json:"id"`
  	Username string `json:"username"`
  	RoleName string `json:"roleName"`
  	RoleNo string `json:"roleNo"`
  	UserNo string `json:"userNo"`
  	RegisterDate string `json:"registerDate"`
  }

  // Get user info by token. This endpoint is expected to be accessible publicly
  func ApiGetTokenUserInfo(rail miso.Rail, token string) (UserInfoBrief, error) {
  	var res miso.GnResp[UserInfoBrief]
  	err := miso.NewDynClient(rail, "/open/api/token/user", "user-vault").
  		AddQuery("token", token).
  		Get().
  		Json(&res)
  	if err != nil {
  		var dat UserInfoBrief
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
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

## POST /open/api/access/history

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListAccessLogReq struct {
  	Paging miso.Paging `json:"paging"`
  }


  type ListedAccessLog struct {
  	Id int `json:"id"`
  	UserAgent string `json:"userAgent"`
  	IpAddress string `json:"ipAddress"`
  	Username string `json:"username"`
  	Url string `json:"url"`
  	AccessTime atom.Time `json:"accessTime"`
  	Success bool `json:"success"`
  }

  // User list access logs
  func ApiUserListAccessHistory(rail miso.Rail, req ListAccessLogReq) (miso.PageRes[ListedAccessLog], error) {
  	var res miso.GnResp[miso.PageRes[ListedAccessLog]]
  	err := miso.NewDynClient(rail, "/open/api/access/history", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		var dat miso.PageRes[ListedAccessLog]
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface ListAccessLogReq {
    paging?: Paging;
  }

  export interface Paging {
    limit?: number;                // page limit
    page?: number;                 // page number, 1-based
    total?: number;                // total count
  }

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

## POST /open/api/user/key/generate

- Description: User generate user key
- Bound to Resource: `"basic-user"`
- JSON Request:
    - "password": (string) Required.
    - "keyName": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type GenUserKeyReq struct {
  	Password string `json:"password"` // Required.
  	KeyName string `json:"keyName"` // Required.
  }

  // User generate user key
  func ApiUserGenUserKey(rail miso.Rail, req GenUserKeyReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/user/key/generate", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface GenUserKeyReq {
    password?: string;             // Required.
    keyName?: string;              // Required.
  }

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

## POST /open/api/user/key/list

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListUserKeysReq struct {
  	Paging miso.Paging `json:"paging"`
  	Name string `json:"name"`
  }


  type ListedUserKey struct {
  	Id int `json:"id"`
  	SecretKey string `json:"secretKey"`
  	Name string `json:"name"`
  	ExpirationTime atom.Time `json:"expirationTime"`
  	CreateTime atom.Time `json:"createTime"`
  }

  // User list user keys
  func ApiUserListUserKeys(rail miso.Rail, req ListUserKeysReq) (miso.PageRes[ListedUserKey], error) {
  	var res miso.GnResp[miso.PageRes[ListedUserKey]]
  	err := miso.NewDynClient(rail, "/open/api/user/key/list", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		var dat miso.PageRes[ListedUserKey]
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
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

## POST /open/api/user/key/delete

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type DeleteUserKeyReq struct {
  	UserKeyId int `json:"userKeyId"`
  }

  // User delete user key
  func ApiUserDeleteUserKey(rail miso.Rail, req DeleteUserKeyReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/user/key/delete", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface DeleteUserKeyReq {
    userKeyId?: number;
  }

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

## POST /open/api/resource/add

- Description: Admin add resource
- Bound to Resource: `"manage-resources"`
- JSON Request:
    - "name": (string) Required. Max length: 32.
    - "code": (string) Required. Max length: 32.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type CreateResReq struct {
  	Name string `json:"name"`      // Required. Max length: 32.
  	Code string `json:"code"`      // Required. Max length: 32.
  }

  // Admin add resource
  func ApiAdminAddResource(rail miso.Rail, req CreateResReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/resource/add", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface CreateResReq {
    name?: string;                 // Required. Max length: 32.
    code?: string;                 // Required. Max length: 32.
  }

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

## POST /open/api/resource/remove

- Description: Admin remove resource
- Bound to Resource: `"manage-resources"`
- JSON Request:
    - "resCode": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type DeleteResourceReq struct {
  	ResCode string `json:"resCode"` // Required.
  }

  // Admin remove resource
  func ApiAdminRemoveResource(rail miso.Rail, req DeleteResourceReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/resource/remove", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface DeleteResourceReq {
    resCode?: string;              // Required.
  }

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

## GET /open/api/resource/brief/candidates

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ResBrief struct {
  	Code string `json:"code"`
  	Name string `json:"name"`
  }

  // List all resource candidates for role
  func ApiListResCandidates(rail miso.Rail, roleNo string) ([]ResBrief, error) {
  	var res miso.GnResp[[]ResBrief]
  	err := miso.NewDynClient(rail, "/open/api/resource/brief/candidates", "user-vault").
  		AddQuery("roleNo", roleNo).
  		Get().
  		Json(&res)
  	if err != nil {
  		var dat []ResBrief
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
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

## POST /open/api/resource/list

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListResReq struct {
  	Paging miso.Paging `json:"paging"`
  }

  type ListResResp struct {
  	Paging miso.Paging `json:"paging"`
  	Payload []WRes `json:"payload"`
  }

  type WRes struct {
  	Id int `json:"id"`
  	Code string `json:"code"`
  	Name string `json:"name"`
  	CreateTime atom.Time `json:"createTime"`
  	CreateBy string `json:"createBy"`
  	UpdateTime atom.Time `json:"updateTime"`
  	UpdateBy string `json:"updateBy"`
  }

  // Admin list resources
  func ApiAdminListRes(rail miso.Rail, req ListResReq) (ListResResp, error) {
  	var res miso.GnResp[ListResResp]
  	err := miso.NewDynClient(rail, "/open/api/resource/list", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		var dat ListResResp
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface ListResReq {
    paging?: Paging;
  }

  export interface Paging {
    limit?: number;                // page limit
    page?: number;                 // page number, 1-based
    total?: number;                // total count
  }

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

## GET /open/api/resource/brief/user

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ResBrief struct {
  	Code string `json:"code"`
  	Name string `json:"name"`
  }

  // List resources that are accessible to current user
  func ApiListUserAccessibleRes(rail miso.Rail) ([]ResBrief, error) {
  	var res miso.GnResp[[]ResBrief]
  	err := miso.NewDynClient(rail, "/open/api/resource/brief/user", "user-vault").
  		Get().
  		Json(&res)
  	if err != nil {
  		var dat []ResBrief
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
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

## GET /open/api/resource/brief/all

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ResBrief struct {
  	Code string `json:"code"`
  	Name string `json:"name"`
  }

  // List all resource brief info
  func ApiListAllResBrief(rail miso.Rail) ([]ResBrief, error) {
  	var res miso.GnResp[[]ResBrief]
  	err := miso.NewDynClient(rail, "/open/api/resource/brief/all", "user-vault").
  		Get().
  		Json(&res)
  	if err != nil {
  		var dat []ResBrief
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
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

## POST /open/api/role/resource/add

- Description: Admin add resource to role
- Bound to Resource: `"manage-resources"`
- JSON Request:
    - "roleNo": (string) Required.
    - "resCode": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type AddRoleResReq struct {
  	RoleNo string `json:"roleNo"`  // Required.
  	ResCode string `json:"resCode"` // Required.
  }

  // Admin add resource to role
  func ApiAdminBindRoleRes(rail miso.Rail, req AddRoleResReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/role/resource/add", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface AddRoleResReq {
    roleNo?: string;               // Required.
    resCode?: string;              // Required.
  }

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

## POST /open/api/role/resource/remove

- Description: Admin remove resource from role
- Bound to Resource: `"manage-resources"`
- JSON Request:
    - "roleNo": (string) Required.
    - "resCode": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type RemoveRoleResReq struct {
  	RoleNo string `json:"roleNo"`  // Required.
  	ResCode string `json:"resCode"` // Required.
  }

  // Admin remove resource from role
  func ApiAdminUnbindRoleRes(rail miso.Rail, req RemoveRoleResReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/role/resource/remove", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface RemoveRoleResReq {
    roleNo?: string;               // Required.
    resCode?: string;              // Required.
  }

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

## POST /open/api/role/add

- Description: Admin add role
- Bound to Resource: `"manage-resources"`
- JSON Request:
    - "name": (string) Required. Max length: 32.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type AddRoleReq struct {
  	Name string `json:"name"`      // Required. Max length: 32.
  }

  // Admin add role
  func ApiAdminAddRole(rail miso.Rail, req AddRoleReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/role/add", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface AddRoleReq {
    name?: string;                 // Required. Max length: 32.
  }

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

## POST /open/api/role/list

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListRoleReq struct {
  	Paging miso.Paging `json:"paging"`
  }

  type ListRoleResp struct {
  	Payload []WRole `json:"payload"`
  	Paging miso.Paging `json:"paging"`
  }

  type WRole struct {
  	Id int `json:"id"`
  	RoleNo string `json:"roleNo"`
  	Name string `json:"name"`
  	CreateTime atom.Time `json:"createTime"`
  	CreateBy string `json:"createBy"`
  	UpdateTime atom.Time `json:"updateTime"`
  	UpdateBy string `json:"updateBy"`
  }

  // Admin list roles
  func ApiAdminListRoles(rail miso.Rail, req ListRoleReq) (ListRoleResp, error) {
  	var res miso.GnResp[ListRoleResp]
  	err := miso.NewDynClient(rail, "/open/api/role/list", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		var dat ListRoleResp
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface ListRoleReq {
    paging?: Paging;
  }

  export interface Paging {
    limit?: number;                // page limit
    page?: number;                 // page number, 1-based
    total?: number;                // total count
  }

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

## GET /open/api/role/brief/all

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type RoleBrief struct {
  	RoleNo string `json:"roleNo"`
  	Name string `json:"name"`
  }

  // Admin list role brief info
  func ApiAdminListRoleBriefs(rail miso.Rail) ([]RoleBrief, error) {
  	var res miso.GnResp[[]RoleBrief]
  	err := miso.NewDynClient(rail, "/open/api/role/brief/all", "user-vault").
  		Get().
  		Json(&res)
  	if err != nil {
  		var dat []RoleBrief
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
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

## POST /open/api/role/resource/list

- Description: Admin list resources of role
- Bound to Resource: `"manage-resources"`
- JSON Request:
    - "paging": (Paging) 
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
    - "roleNo": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListRoleResReq struct {
  	Paging miso.Paging `json:"paging"`
  	RoleNo string `json:"roleNo"`  // Required.
  }

  type ListRoleResResp struct {
  	Paging miso.Paging `json:"paging"`
  	Payload []ListedRoleRes `json:"payload"`
  }

  type ListedRoleRes struct {
  	Id int `json:"id"`
  	ResCode string `json:"resCode"`
  	ResName string `json:"resName"`
  	CreateTime atom.Time `json:"createTime"`
  	CreateBy string `json:"createBy"`
  }

  // Admin list resources of role
  func ApiAdminListRoleRes(rail miso.Rail, req ListRoleResReq) (ListRoleResResp, error) {
  	var res miso.GnResp[ListRoleResResp]
  	err := miso.NewDynClient(rail, "/open/api/role/resource/list", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		var dat ListRoleResResp
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface ListRoleResReq {
    paging?: Paging;
    roleNo?: string;               // Required.
  }

  export interface Paging {
    limit?: number;                // page limit
    page?: number;                 // page number, 1-based
    total?: number;                // total count
  }

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

## POST /open/api/role/info

- Description: Get role info
- Expected Access Scope: PUBLIC
- JSON Request:
    - "roleNo": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type RoleInfoReq struct {
  	RoleNo string `json:"roleNo"`  // Required.
  }

  type RoleInfoResp struct {
  	RoleNo string `json:"roleNo"`
  	Name string `json:"name"`
  }

  // Get role info
  func ApiGetRoleInfo(rail miso.Rail, req RoleInfoReq) (RoleInfoResp, error) {
  	var res miso.GnResp[RoleInfoResp]
  	err := miso.NewDynClient(rail, "/open/api/role/info", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		var dat RoleInfoResp
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface RoleInfoReq {
    roleNo?: string;               // Required.
  }

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

## POST /open/api/path/list

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListPathReq struct {
  	ResCode string `json:"resCode"`
  	Pgroup string `json:"pgroup"`
  	Url string `json:"url"`
  	Ptype string `json:"ptype"`    // path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible
  	Paging miso.Paging `json:"paging"`
  }

  type ListPathResp struct {
  	Paging miso.Paging `json:"paging"`
  	Payload []WPath `json:"payload"`
  }

  type WPath struct {
  	Id int `json:"id"`
  	Pgroup string `json:"pgroup"`
  	PathNo string `json:"pathNo"`
  	Method string `json:"method"`
  	Desc string `json:"desc"`
  	Url string `json:"url"`
  	Ptype string `json:"ptype"`    // path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible
  	CreateTime atom.Time `json:"createTime"`
  	CreateBy string `json:"createBy"`
  	UpdateTime atom.Time `json:"updateTime"`
  	UpdateBy string `json:"updateBy"`
  }

  // Admin list paths
  func ApiAdminListPaths(rail miso.Rail, req ListPathReq) (ListPathResp, error) {
  	var res miso.GnResp[ListPathResp]
  	err := miso.NewDynClient(rail, "/open/api/path/list", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		var dat ListPathResp
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
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

## POST /open/api/path/resource/bind

- Description: Admin bind resource to path
- Bound to Resource: `"manage-resources"`
- JSON Request:
    - "pathNo": (string) Required.
    - "resCode": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type BindPathResReq struct {
  	PathNo string `json:"pathNo"`  // Required.
  	ResCode string `json:"resCode"` // Required.
  }

  // Admin bind resource to path
  func ApiAdminBindResPath(rail miso.Rail, req BindPathResReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/path/resource/bind", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface BindPathResReq {
    pathNo?: string;               // Required.
    resCode?: string;              // Required.
  }

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

## POST /open/api/path/resource/unbind

- Description: Admin unbind resource and path
- Bound to Resource: `"manage-resources"`
- JSON Request:
    - "pathNo": (string) Required.
    - "resCode": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type UnbindPathResReq struct {
  	PathNo string `json:"pathNo"`  // Required.
  	ResCode string `json:"resCode"` // Required.
  }

  // Admin unbind resource and path
  func ApiAdminUnbindResPath(rail miso.Rail, req UnbindPathResReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/path/resource/unbind", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface UnbindPathResReq {
    pathNo?: string;               // Required.
    resCode?: string;              // Required.
  }

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

## POST /open/api/path/delete

- Description: Admin delete path
- Bound to Resource: `"manage-resources"`
- JSON Request:
    - "pathNo": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type DeletePathReq struct {
  	PathNo string `json:"pathNo"`  // Required.
  }

  // Admin delete path
  func ApiAdminDeletePath(rail miso.Rail, req DeletePathReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/path/delete", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface DeletePathReq {
    pathNo?: string;               // Required.
  }

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

## POST /open/api/path/update

- Description: Admin update path
- Bound to Resource: `"manage-resources"`
- JSON Request:
    - "type": (string) path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible. Enums: ["PROTECTED","PUBLIC"]. Required.
    - "pathNo": (string) Required.
    - "group": (string) Required. Max length: 20.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type UpdatePathReq struct {
  	Type string `json:"type"`      // path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible. Enums: ["PROTECTED","PUBLIC"]. Required.
  	PathNo string `json:"pathNo"`  // Required.
  	Group string `json:"group"`    // Required. Max length: 20.
  	ResCode string `json:"resCode"`
  }

  // Admin update path
  func ApiAdminUpdatePath(rail miso.Rail, req UpdatePathReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/path/update", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface UpdatePathReq {
    type?: string;                 // path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible. Enums: ["PROTECTED","PUBLIC"]. Required.
    pathNo?: string;               // Required.
    group?: string;                // Required. Max length: 20.
    resCode?: string;
  }

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

## POST /remote/user/info

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type FindUserReq struct {
  	UserId *int `json:"userId"`
  	UserNo *string `json:"userNo"`
  	Username *string `json:"username"`
  }

  type UserInfo struct {
  	Id int `json:"id"`
  	Username string `json:"username"`
  	RoleName string `json:"roleName"`
  	RoleNo string `json:"roleNo"`
  	UserNo string `json:"userNo"`
  	ReviewStatus string `json:"reviewStatus"`
  	IsDisabled int `json:"isDisabled"`
  	CreateTime atom.Time `json:"createTime"`
  	CreateBy string `json:"createBy"`
  	UpdateTime atom.Time `json:"updateTime"`
  	UpdateBy string `json:"updateBy"`
  }

  // Fetch user info
  func ApiFetchUserInfo(rail miso.Rail, req FindUserReq) (UserInfo, error) {
  	var res miso.GnResp[UserInfo]
  	err := miso.NewDynClient(rail, "/remote/user/info", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		var dat UserInfo
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface FindUserReq {
    userId?: number;
    userNo?: string;
    username?: string;
  }

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

## POST /internal/v1/user/info/common

- Description: System fetch user info as common.User
- JSON Request:
    - "userId": (*int) 
    - "userNo": (*string) 
    - "username": (*string) 
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (User) response data
      - "userNo": (string) 
      - "username": (string) 
      - "roleNo": (string) 
- cURL:
  ```sh
  curl -X POST 'http://localhost:8089/internal/v1/user/info/common' \
    -H 'Content-Type: application/json' \
    -d '{"userId":0,"userNo":"","username":""}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type FindUserReq struct {
  	UserId *int `json:"userId"`
  	UserNo *string `json:"userNo"`
  	Username *string `json:"username"`
  }

  type User struct {
  	UserNo string `json:"userNo"`
  	Username string `json:"username"`
  	RoleNo string `json:"roleNo"`
  }

  // System fetch user info as common.User
  func ApiSysFetchUserInfo(rail miso.Rail, req FindUserReq) (miso.User, error) {
  	var res miso.GnResp[miso.User]
  	err := miso.NewDynClient(rail, "/internal/v1/user/info/common", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		var dat miso.User
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface FindUserReq {
    userId?: number;
    userNo?: string;
    username?: string;
  }

  export interface Resp {
    errorCode?: string;            // error code
    msg?: string;                  // message
    error?: boolean;               // whether the request was successful
    data?: User;
  }

  export interface User {
    userNo?: string;
    username?: string;
    roleNo?: string;
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

  sysFetchUserInfo() {
    let req: FindUserReq | null = null;
    this.http.post<any>(`/user-vault/internal/v1/user/info/common`, req)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: User = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
  }
  ```

## GET /remote/user/id

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Fetch id of user with the username
  func ApiFetchUserIdByName(rail miso.Rail, username string) (int, error) {
  	var res miso.GnResp[int]
  	err := miso.NewDynClient(rail, "/remote/user/id", "user-vault").
  		AddQuery("username", username).
  		Get().
  		Json(&res)
  	if err != nil {
  		return 0, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
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

## POST /remote/user/userno/username

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type FetchNameByUserNoReq struct {
  	UserNos []string `json:"userNos"`
  }

  type FetchUsernamesRes struct {
  	UserNoToUsername map[string]string `json:"userNoToUsername"`
  }

  // Fetch usernames of users with the userNos
  func ApiFetchUsernamesByNosEp(rail miso.Rail, req FetchNameByUserNoReq) (FetchUsernamesRes, error) {
  	var res miso.GnResp[FetchUsernamesRes]
  	err := miso.NewDynClient(rail, "/remote/user/userno/username", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		var dat FetchUsernamesRes
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface FetchNameByUserNoReq {
    userNos?: string[];
  }

  export interface Resp {
    errorCode?: string;            // error code
    msg?: string;                  // message
    error?: boolean;               // whether the request was successful
    data?: FetchUsernamesRes;
  }

  export interface FetchUsernamesRes {
    userNoToUsername?: Map<string,string>;
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

## POST /remote/user/list/with-role

- Description: Fetch users with the role_no
- JSON Request:
    - "roleNo": (string) Required.
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": ([]vault.UserInfo) response data
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type FetchUsersWithRoleReq struct {
  	RoleNo string `json:"roleNo"`  // Required.
  }

  type UserInfo struct {
  	Id int `json:"id"`
  	Username string `json:"username"`
  	RoleName string `json:"roleName"`
  	RoleNo string `json:"roleNo"`
  	UserNo string `json:"userNo"`
  	ReviewStatus string `json:"reviewStatus"`
  	IsDisabled int `json:"isDisabled"`
  	CreateTime atom.Time `json:"createTime"`
  	CreateBy string `json:"createBy"`
  	UpdateTime atom.Time `json:"updateTime"`
  	UpdateBy string `json:"updateBy"`
  }

  // Fetch users with the role_no
  func ApiFindUserWithRoleEp(rail miso.Rail, req FetchUsersWithRoleReq) ([]UserInfo, error) {
  	var res miso.GnResp[[]UserInfo]
  	err := miso.NewDynClient(rail, "/remote/user/list/with-role", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		var dat []UserInfo
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface FetchUsersWithRoleReq {
    roleNo?: string;               // Required.
  }

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

## POST /remote/user/list/with-resource

- Description: Fetch users that have access to the resource
- JSON Request:
    - "resourceCode": (string) 
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": ([]vault.UserInfo) response data
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type FetchUserWithResourceReq struct {
  	ResourceCode string `json:"resourceCode"`
  }

  type UserInfo struct {
  	Id int `json:"id"`
  	Username string `json:"username"`
  	RoleName string `json:"roleName"`
  	RoleNo string `json:"roleNo"`
  	UserNo string `json:"userNo"`
  	ReviewStatus string `json:"reviewStatus"`
  	IsDisabled int `json:"isDisabled"`
  	CreateTime atom.Time `json:"createTime"`
  	CreateBy string `json:"createBy"`
  	UpdateTime atom.Time `json:"updateTime"`
  	UpdateBy string `json:"updateBy"`
  }

  // Fetch users that have access to the resource
  func ApiFindUserWithResourceEp(rail miso.Rail, req FetchUserWithResourceReq) ([]UserInfo, error) {
  	var res miso.GnResp[[]UserInfo]
  	err := miso.NewDynClient(rail, "/remote/user/list/with-resource", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		var dat []UserInfo
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface FetchUserWithResourceReq {
    resourceCode?: string;
  }

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

## POST /remote/resource/add

- Description: Report resource. This endpoint should be used internally by another backend service.
- JSON Request:
    - "name": (string) Required. Max length: 32.
    - "code": (string) Required. Max length: 32.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type CreateResReq struct {
  	Name string `json:"name"`      // Required. Max length: 32.
  	Code string `json:"code"`      // Required. Max length: 32.
  }

  // Report resource. This endpoint should be used internally by another backend service.
  func ApiReportResourceEp(rail miso.Rail, req CreateResReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/remote/resource/add", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface CreateResReq {
    name?: string;                 // Required. Max length: 32.
    code?: string;                 // Required. Max length: 32.
  }

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

## POST /remote/path/resource/access-test

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type CheckResAccessReq struct {
  	RoleNo string `json:"roleNo"`
  	Url string `json:"url"`
  	Method string `json:"method"`
  }

  type CheckResAccessResp struct {
  	Valid bool `json:"valid"`
  }

  // Validate resource access
  func ApiCheckResourceAccessEp(rail miso.Rail, req CheckResAccessReq) (CheckResAccessResp, error) {
  	var res miso.GnResp[CheckResAccessResp]
  	err := miso.NewDynClient(rail, "/remote/path/resource/access-test", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		var dat CheckResAccessResp
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface CheckResAccessReq {
    roleNo?: string;
    url?: string;
    method?: string;
  }

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

## POST /remote/path/add

- Description: Report endpoint info
- Bound to Resource: `"manage-resources"`
- JSON Request:
    - "type": (string) path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible. Enums: ["PROTECTED","PUBLIC"]. Required.
    - "url": (string) Required. Max length: 128.
    - "group": (string) Required. Max length: 20.
    - "method": (string) Required. Max length: 10.
    - "desc": (string) Max length: 255.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type CreatePathReq struct {
  	Type string `json:"type"`      // path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible. Enums: ["PROTECTED","PUBLIC"]. Required.
  	Url string `json:"url"`        // Required. Max length: 128.
  	Group string `json:"group"`    // Required. Max length: 20.
  	Method string `json:"method"`  // Required. Max length: 10.
  	Desc string `json:"desc"`      // Max length: 255.
  	ResCode string `json:"resCode"`
  }

  // Report endpoint info
  func ApiReportPath(rail miso.Rail, req CreatePathReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/remote/path/add", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface CreatePathReq {
    type?: string;                 // path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible. Enums: ["PROTECTED","PUBLIC"]. Required.
    url?: string;                  // Required. Max length: 128.
    group?: string;                // Required. Max length: 20.
    method?: string;               // Required. Max length: 10.
    desc?: string;                 // Max length: 255.
    resCode?: string;
  }

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

## POST /open/api/password/list-site-passwords

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListSitePasswordReq struct {
  	Alias string `json:"alias"`
  	Site string `json:"site"`
  	Username string `json:"username"`
  	Paging miso.Paging `json:"paging"`
  }


  type ListSitePasswordRes struct {
  	RecordId string `json:"recordId"`
  	Site string `json:"site"`
  	Alias string `json:"alias"`
  	Username string `json:"username"`
  	CreateTime atom.Time `json:"createTime"`
  }

  // List site password records
  func ApiListSitePasswords(rail miso.Rail, req ListSitePasswordReq) (miso.PageRes[ListSitePasswordRes], error) {
  	var res miso.GnResp[miso.PageRes[ListSitePasswordRes]]
  	err := miso.NewDynClient(rail, "/open/api/password/list-site-passwords", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		var dat miso.PageRes[ListSitePasswordRes]
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
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

## POST /open/api/password/add-site-password

- Description: Add site password record
- Bound to Resource: `"basic-user"`
- JSON Request:
    - "site": (string) 
    - "alias": (string) 
    - "username": (string) Required.
    - "sitePassword": (string) Required.
    - "loginPassword": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type AddSitePasswordReq struct {
  	Site string `json:"site"`
  	Alias string `json:"alias"`
  	Username string `json:"username"` // Required.
  	SitePassword string `json:"sitePassword"` // Required.
  	LoginPassword string `json:"loginPassword"` // Required.
  }

  // Add site password record
  func ApiAddSitePassword(rail miso.Rail, req AddSitePasswordReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/password/add-site-password", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface AddSitePasswordReq {
    site?: string;
    alias?: string;
    username?: string;             // Required.
    sitePassword?: string;         // Required.
    loginPassword?: string;        // Required.
  }

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

## POST /open/api/password/remove-site-password

- Description: Remove site password record
- Bound to Resource: `"basic-user"`
- JSON Request:
    - "recordId": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type RemoveSitePasswordRes struct {
  	RecordId string `json:"recordId"` // Required.
  }

  // Remove site password record
  func ApiRemoveSitePassword(rail miso.Rail, req RemoveSitePasswordRes) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/password/remove-site-password", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface RemoveSitePasswordRes {
    recordId?: string;             // Required.
  }

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

## POST /open/api/password/decrypt-site-password

- Description: Decrypt site password
- Bound to Resource: `"basic-user"`
- JSON Request:
    - "loginPassword": (string) Required.
    - "recordId": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type DecryptSitePasswordReq struct {
  	LoginPassword string `json:"loginPassword"` // Required.
  	RecordId string `json:"recordId"` // Required.
  }

  type DecryptSitePasswordRes struct {
  	Decrypted string `json:"decrypted"`
  }

  // Decrypt site password
  func ApiDecryptSitePassword(rail miso.Rail, req DecryptSitePasswordReq) (DecryptSitePasswordRes, error) {
  	var res miso.GnResp[DecryptSitePasswordRes]
  	err := miso.NewDynClient(rail, "/open/api/password/decrypt-site-password", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		var dat DecryptSitePasswordRes
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface DecryptSitePasswordReq {
    loginPassword?: string;        // Required.
    recordId?: string;             // Required.
  }

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

## POST /open/api/password/edit-site-password

- Description: Edit site password
- Bound to Resource: `"basic-user"`
- JSON Request:
    - "recordId": (string) 
    - "site": (string) 
    - "username": (string) 
    - "alias": (string) 
    - "sitePassword": (string) new site password, optional
    - "loginPassword": (string) only used when site password is provided
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
- cURL:
  ```sh
  curl -X POST 'http://localhost:8089/open/api/password/edit-site-password' \
    -H 'Content-Type: application/json' \
    -d '{"alias":"","loginPassword":"","recordId":"","site":"","sitePassword":"","username":""}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type EditSitePasswordReq struct {
  	RecordId string `json:"recordId"`
  	Site string `json:"site"`
  	Username string `json:"username"`
  	Alias string `json:"alias"`
  	SitePassword string `json:"sitePassword"` // new site password, optional
  	LoginPassword string `json:"loginPassword"` // only used when site password is provided
  }

  // Edit site password
  func ApiEditSitePassword(rail miso.Rail, req EditSitePasswordReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/password/edit-site-password", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface EditSitePasswordReq {
    recordId?: string;
    site?: string;
    username?: string;
    alias?: string;
    sitePassword?: string;         // new site password, optional
    loginPassword?: string;        // only used when site password is provided
  }

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

## POST /open/api/user/clear-failed-login-attempts

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ClearUserFailedLoginAttemptsReq struct {
  	UserNo string `json:"userNo"`
  }

  // Admin clear user's failed login attempts
  func ApiClearUserFailedLoginAttempts(rail miso.Rail, req ClearUserFailedLoginAttemptsReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/user/clear-failed-login-attempts", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface ClearUserFailedLoginAttemptsReq {
    userNo?: string;
  }

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

## POST /open/api/note/list-notes

- Description: List User Notes
- Bound to Resource: `"basic-user"`
- JSON Request:
    - "keywords": (string) 
    - "paging": (Paging) 
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (PageRes[github.com/curtisnewbie/user-vault/internal/note.Note]) response data
      - "paging": (Paging) pagination parameters
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
      - "payload": ([]note.Note) payload values in current page
        - "recordId": (string) 
        - "title": (string) 
        - "content": (string) 
        - "userNo": (string) 
        - "createdAt": (int64) 
        - "updatedAt": (int64) 
- cURL:
  ```sh
  curl -X POST 'http://localhost:8089/open/api/note/list-notes' \
    -H 'Content-Type: application/json' \
    -d '{"keywords":"","paging":{"limit":0,"page":0,"total":0}}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListNoteReq struct {
  	Keywords string `json:"keywords"`
  	Paging miso.Paging `json:"paging"`
  }


  type Note struct {
  	RecordId string `json:"recordId"`
  	Title string `json:"title"`
  	Content string `json:"content"`
  	UserNo string `json:"userNo"`
  	CreatedAt atom.Time `json:"createdAt"`
  	UpdatedAt atom.Time `json:"updatedAt"`
  }

  // List User Notes
  func ApiListNotes(rail miso.Rail, req ListNoteReq) (miso.PageRes[Note], error) {
  	var res miso.GnResp[miso.PageRes[Note]]
  	err := miso.NewDynClient(rail, "/open/api/note/list-notes", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		var dat miso.PageRes[Note]
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface ListNoteReq {
    keywords?: string;
    paging?: Paging;
  }

  export interface Paging {
    limit?: number;                // page limit
    page?: number;                 // page number, 1-based
    total?: number;                // total count
  }

  export interface Resp {
    errorCode?: string;            // error code
    msg?: string;                  // message
    error?: boolean;               // whether the request was successful
    data?: PageRes;
  }

  export interface PageRes {
    paging?: Paging;
    payload?: Note[];
  }

  export interface Paging {
    limit?: number;                // page limit
    page?: number;                 // page number, 1-based
    total?: number;                // total count
  }

  export interface Note {
    recordId?: string;
    title?: string;
    content?: string;
    userNo?: string;
    createdAt?: number;
    updatedAt?: number;
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

  listNotes() {
    let req: ListNoteReq | null = null;
    this.http.post<any>(`/user-vault/open/api/note/list-notes`, req)
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

- Angular NgTable Demo:
  ```html
  <table mat-table [dataSource]="tabdata" class="mb-4" style="width: 100%;">
  	<ng-container matColumnDef="recordId">
  		<th mat-header-cell *matHeaderCellDef> RecordId </th>
  		<td mat-cell *matCellDef="let u"> {{u.recordId}} </td>
  	</ng-container>
  	<ng-container matColumnDef="title">
  		<th mat-header-cell *matHeaderCellDef> Title </th>
  		<td mat-cell *matCellDef="let u"> {{u.title}} </td>
  	</ng-container>
  	<ng-container matColumnDef="content">
  		<th mat-header-cell *matHeaderCellDef> Content </th>
  		<td mat-cell *matCellDef="let u"> {{u.content}} </td>
  	</ng-container>
  	<ng-container matColumnDef="userNo">
  		<th mat-header-cell *matHeaderCellDef> UserNo </th>
  		<td mat-cell *matCellDef="let u"> {{u.userNo}} </td>
  	</ng-container>
  	<ng-container matColumnDef="createdAt">
  		<th mat-header-cell *matHeaderCellDef> CreatedAt </th>
  		<td mat-cell *matCellDef="let u"> {{u.createdAt | date: 'yyyy-MM-dd HH:mm:ss'}} </td>
  	</ng-container>
  	<ng-container matColumnDef="updatedAt">
  		<th mat-header-cell *matHeaderCellDef> UpdatedAt </th>
  		<td mat-cell *matCellDef="let u"> {{u.updatedAt | date: 'yyyy-MM-dd HH:mm:ss'}} </td>
  	</ng-container>
  	<tr mat-row *matRowDef="let row; columns: ['recordId','title','content','userNo','createdAt','updatedAt'];"></tr>
  	<tr mat-header-row *matHeaderRowDef="['recordId','title','content','userNo','createdAt','updatedAt']"></tr>
  </table>
  ```

## POST /open/api/note/save-note

- Description: User Save Note
- Bound to Resource: `"basic-user"`
- JSON Request:
    - "title": (string) Required.
    - "content": (string) 
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
- cURL:
  ```sh
  curl -X POST 'http://localhost:8089/open/api/note/save-note' \
    -H 'Content-Type: application/json' \
    -d '{"content":"","title":""}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type SaveNoteReq struct {
  	Title string `json:"title"`    // Required.
  	Content string `json:"content"`
  }

  // User Save Note
  func ApiSaveNote(rail miso.Rail, req SaveNoteReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/note/save-note", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface SaveNoteReq {
    title?: string;                // Required.
    content?: string;
  }

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

  saveNote() {
    let req: SaveNoteReq | null = null;
    this.http.post<any>(`/user-vault/open/api/note/save-note`, req)
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

## POST /open/api/note/update-note

- Description: User Update Note
- Bound to Resource: `"basic-user"`
- JSON Request:
    - "recordId": (string) Required.
    - "title": (string) Required.
    - "content": (string) 
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
- cURL:
  ```sh
  curl -X POST 'http://localhost:8089/open/api/note/update-note' \
    -H 'Content-Type: application/json' \
    -d '{"content":"","recordId":"","title":""}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type UpdateNoteReq struct {
  	RecordId string `json:"recordId"` // Required.
  	Title string `json:"title"`    // Required.
  	Content string `json:"content"`
  }

  // User Update Note
  func ApiUpdateNote(rail miso.Rail, req UpdateNoteReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/note/update-note", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface UpdateNoteReq {
    recordId?: string;             // Required.
    title?: string;                // Required.
    content?: string;
  }

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

  updateNote() {
    let req: UpdateNoteReq | null = null;
    this.http.post<any>(`/user-vault/open/api/note/update-note`, req)
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

## POST /open/api/note/delete-note

- Description: User Delete Note
- Bound to Resource: `"basic-user"`
- JSON Request:
    - "recordId": (string) 
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
- cURL:
  ```sh
  curl -X POST 'http://localhost:8089/open/api/note/delete-note' \
    -H 'Content-Type: application/json' \
    -d '{"recordId":""}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ApiDeleteNoteReq struct {
  	RecordId string `json:"recordId"`
  }

  // User Delete Note
  func ApiDeleteNote(rail miso.Rail, req ApiDeleteNoteReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/note/delete-note", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface ApiDeleteNoteReq {
    recordId?: string;
  }

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

  deleteNote() {
    let req: ApiDeleteNoteReq | null = null;
    this.http.post<any>(`/user-vault/open/api/note/delete-note`, req)
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

## POST /open/api/v1/notification/create

- Description: Create platform notification
- Bound to Resource: `"postbox:notification:create"`
- JSON Request:
    - "title": (string) Max length: 255.
    - "message": (string) Max length: 1000.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type CreateNotificationReq struct {
  	Title string `json:"title"`    // Max length: 255.
  	Message string `json:"message"` // Max length: 1000.
  	ReceiverUserNos []string `json:"receiverUserNos"`
  }

  // Create platform notification
  func SendCreateNotificationReq(rail miso.Rail, req CreateNotificationReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/v1/notification/create", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface CreateNotificationReq {
    title?: string;                // Max length: 255.
    message?: string;              // Max length: 1000.
    receiverUserNos?: string[];
  }

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

## POST /open/api/v1/notification/query

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
    - "data": (PageRes[github.com/curtisnewbie/user-vault/internal/repo.ListedNotification]) response data
      - "paging": (Paging) pagination parameters
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
      - "payload": ([]repo.ListedNotification) payload values in current page
        - "id": (int) 
        - "notifiNo": (string) 
        - "title": (string) 
        - "message": (string) 
        - "status": (string) 
        - "createTime": (int64) 
- cURL:
  ```sh
  curl -X POST 'http://localhost:8089/open/api/v1/notification/query' \
    -H 'Content-Type: application/json' \
    -d '{"page":{"limit":0,"page":0,"total":0},"status":""}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type QueryNotificationReq struct {
  	Page miso.Paging `json:"page"`
  	Status string `json:"status"`
  }


  type ListedNotification struct {
  	Id int `json:"id"`
  	NotifiNo string `json:"notifiNo"`
  	Title string `json:"title"`
  	Message string `json:"message"`
  	Status string `json:"status"`
  	CreateTime atom.Time `json:"createTime"`
  }

  // Query platform notification
  func SendQueryNotificationReq(rail miso.Rail, req QueryNotificationReq) (miso.PageRes[ListedNotification], error) {
  	var res miso.GnResp[miso.PageRes[ListedNotification]]
  	err := miso.NewDynClient(rail, "/open/api/v1/notification/query", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		var dat miso.PageRes[ListedNotification]
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
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

  export interface Resp {
    errorCode?: string;            // error code
    msg?: string;                  // message
    error?: boolean;               // whether the request was successful
    data?: PageRes;
  }

  export interface PageRes {
    paging?: Paging;
    payload?: ListedNotification[];
  }

  export interface Paging {
    limit?: number;                // page limit
    page?: number;                 // page number, 1-based
    total?: number;                // total count
  }

  export interface ListedNotification {
    id?: number;
    notifiNo?: string;
    title?: string;
    message?: string;
    status?: string;
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

  sendQueryNotificationReq() {
    let req: QueryNotificationReq | null = null;
    this.http.post<any>(`/user-vault/open/api/v1/notification/query`, req)
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

## GET /open/api/v1/notification/count

- Description: Count received platform notification
- Bound to Resource: `"postbox:notification:query"`
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (int) response data
- cURL:
  ```sh
  curl -X GET 'http://localhost:8089/open/api/v1/notification/count'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Count received platform notification
  func SendRequest(rail miso.Rail) (int, error) {
  	var res miso.GnResp[int]
  	err := miso.NewDynClient(rail, "/open/api/v1/notification/count", "user-vault").
  		Get().
  		Json(&res)
  	if err != nil {
  		return 0, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
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

  sendRequest() {
    this.http.get<any>(`/user-vault/open/api/v1/notification/count`)
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

## POST /open/api/v1/notification/open

- Description: Record user opened platform notification
- Bound to Resource: `"postbox:notification:query"`
- JSON Request:
    - "notifiNo": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type OpenNotificationReq struct {
  	NotifiNo string `json:"notifiNo"` // Required.
  }

  // Record user opened platform notification
  func SendOpenNotificationReq(rail miso.Rail, req OpenNotificationReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/v1/notification/open", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface OpenNotificationReq {
    notifiNo?: string;             // Required.
  }

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

## POST /open/api/v1/notification/open-all

- Description: Mark all notifications opened
- Bound to Resource: `"postbox:notification:query"`
- JSON Request:
    - "notifiNo": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type OpenNotificationReq struct {
  	NotifiNo string `json:"notifiNo"` // Required.
  }

  // Mark all notifications opened
  func SendOpenNotificationReq(rail miso.Rail, req OpenNotificationReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/v1/notification/open-all", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface OpenNotificationReq {
    notifiNo?: string;             // Required.
  }

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

## GET /open/api/v2/notification/count

- Description: Count received platform notification using long polling
- Bound to Resource: `"postbox:notification:query"`
- Query Parameter:
  - "curr": Current count (used to implement long polling)
- cURL:
  ```sh
  curl -X GET 'http://localhost:8089/open/api/v2/notification/count?curr='
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Count received platform notification using long polling
  func SendRequest(rail miso.Rail, curr string) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/v2/notification/count", "user-vault").
  		AddQuery("curr", curr).
  		Get().
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
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

## GET /debug/trace/recorder/run

- Description: Start FlightRecorder. Recorded result is written to trace.out when it's finished or stopped.
- Query Parameter:
  - "duration": Duration of the flight recording. Required. Duration cannot exceed 30 min.
- cURL:
  ```sh
  curl -X GET 'http://localhost:8089/debug/trace/recorder/run?duration='
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Start FlightRecorder. Recorded result is written to trace.out when it's finished or stopped.
  func SendRequest(rail miso.Rail, duration string) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/debug/trace/recorder/run", "user-vault").
  		AddQuery("duration", duration).
  		Get().
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
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
    let duration: any | null = null;
    this.http.get<any>(`/user-vault/debug/trace/recorder/run?duration=${duration}`)
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

## GET /debug/trace/recorder/snapshot

- Description: FlightRecorder take snapshot. Recorded result is written to trace.out.
- cURL:
  ```sh
  curl -X GET 'http://localhost:8089/debug/trace/recorder/snapshot'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // FlightRecorder take snapshot. Recorded result is written to trace.out.
  func SendRequest(rail miso.Rail) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/debug/trace/recorder/snapshot", "user-vault").
  		Get().
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
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
    this.http.get<any>(`/user-vault/debug/trace/recorder/snapshot`)
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

## GET /debug/trace/recorder/stop

- Description: Stop existing FlightRecorder session.
- cURL:
  ```sh
  curl -X GET 'http://localhost:8089/debug/trace/recorder/stop'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Stop existing FlightRecorder session.
  func SendRequest(rail miso.Rail) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/debug/trace/recorder/stop", "user-vault").
  		Get().
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
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
    this.http.get<any>(`/user-vault/debug/trace/recorder/stop`)
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

## POST /debug/task/disable-workers

- Description: Manually Disable Distributed Task Worker By Name. Use '*' as a special placeholder for all tasks currently registered. For debugging only.
- JSON Request:
    - "tasks": ([]string) 
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
- cURL:
  ```sh
  curl -X POST 'http://localhost:8089/debug/task/disable-workers' \
    -H 'Content-Type: application/json' \
    -d '{"tasks":[]}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type disableTaskWorkerReq struct {
  	Tasks []string `json:"tasks"`
  }

  // Manually Disable Distributed Task Worker By Name. Use '*' as a special placeholder for all tasks currently registered. For debugging only.
  func SendDisableTaskWorkerReq(rail miso.Rail, req disableTaskWorkerReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/debug/task/disable-workers", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface disableTaskWorkerReq {
    tasks?: string[];
  }

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

  sendDisableTaskWorkerReq() {
    let req: disableTaskWorkerReq | null = null;
    this.http.post<any>(`/user-vault/debug/task/disable-workers`, req)
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

## POST /debug/task/enable-workers

- Description: Manually enable previously disabled Distributed Task Worker By Name. Use '*' as a special placeholder for all tasks currently registered. For debugging only.
- JSON Request:
    - "tasks": ([]string) 
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
- cURL:
  ```sh
  curl -X POST 'http://localhost:8089/debug/task/enable-workers' \
    -H 'Content-Type: application/json' \
    -d '{"tasks":[]}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type disableTaskWorkerReq struct {
  	Tasks []string `json:"tasks"`
  }

  // Manually enable previously disabled Distributed Task Worker By Name. Use '*' as a special placeholder for all tasks currently registered. For debugging only.
  func SendDisableTaskWorkerReq(rail miso.Rail, req disableTaskWorkerReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/debug/task/enable-workers", "user-vault").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface disableTaskWorkerReq {
    tasks?: string[];
  }

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

  sendDisableTaskWorkerReq() {
    let req: disableTaskWorkerReq | null = null;
    this.http.post<any>(`/user-vault/debug/task/enable-workers`, req)
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

# Event Pipelines

- CreateNotifiPipeline
  - Description: Pipeline that creates notifications to the specified list of users
  - RabbitMQ Queue: `pieline.user-vault.create-notifi`
  - RabbitMQ Exchange: `pieline.user-vault.create-notifi`
  - RabbitMQ RoutingKey: `#`
  - Event Payload:
    - "title": (string) notification title. Max length: 255.
    - "message": (string) notification content. Max length: 65000.
    - "receiverUserNos": ([]string) user_no of receivers

- CreateNotifiByAccessPipeline
  - Description: Pipeline that creates notifications to users who have access to the specified resource
  - RabbitMQ Queue: `pieline.user-vault.create-notifi.by-access`
  - RabbitMQ Exchange: `pieline.user-vault.create-notifi.by-access`
  - RabbitMQ RoutingKey: `#`
  - Event Payload:
    - "title": (string) notification title. Max length: 255.
    - "message": (string) notification content. Max length: 65000.
    - "resCode": (string) resource code. Required.
