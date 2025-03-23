# API Endpoints

## Contents

- [POST /log/error/list](#post-logerrorlist)
- [GET /auth/resource](#get-authresource)
- [GET /metrics](#get-metrics)

## POST /log/error/list

- Description: List error logs
- Bound to Resource: `"manage-logbot"`
- JSON Request:
    - "app": (string) 
    - "page": (Paging) 
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
- cURL:
  ```sh
  curl -X POST 'http://localhost:8087/log/error/list' \
    -H 'Content-Type: application/json' \
    -d '{"app":"","page":{"limit":0,"page":0,"total":0}}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListErrorLogReq struct {
  	App string `json:"app"`
  	Page miso.Paging `json:"page"`
  }

  func SendListErrorLogReq(rail miso.Rail, req ListErrorLogReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/log/error/list", "logbot").
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
  export interface ListErrorLogReq {
    app?: string;
    page?: Paging;
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

  sendListErrorLogReq() {
    let req: ListErrorLogReq | null = null;
    this.http.post<any>(`/logbot/log/error/list`, req)
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

## GET /auth/resource

- Description: Expose resource and endpoint information to other backend service for authorization.
- Expected Access Scope: PROTECTED
- JSON Response:
    - "resources": ([]auth.Resource) 
      - "name": (string) resource name
      - "code": (string) resource code, unique identifier
    - "paths": ([]auth.Endpoint) 
      - "type": (string) access scope type: PROTECTED/PUBLIC
      - "url": (string) endpoint url
      - "group": (string) app name
      - "desc": (string) description of the endpoint
      - "resCode": (string) resource code
      - "method": (string) http method
- cURL:
  ```sh
  curl -X GET 'http://localhost:8087/auth/resource'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ResourceInfoRes struct {
  	Resources []Resource
  	Paths []Endpoint
  }

  type Resource struct {
  	Name string `json:"name"`      // resource name
  	Code string `json:"code"`      // resource code, unique identifier
  }

  type Endpoint struct {
  	Type string `json:"type"`      // access scope type: PROTECTED/PUBLIC
  	Url string `json:"url"`        // endpoint url
  	Group string `json:"group"`    // app name
  	Desc string `json:"desc"`      // description of the endpoint
  	ResCode string `json:"resCode"` // resource code
  	Method string `json:"method"`  // http method
  }

  func SendRequest(rail miso.Rail) (ResourceInfoRes, error) {
  	var res miso.GnResp[ResourceInfoRes]
  	err := miso.NewDynTClient(rail, "/auth/resource", "logbot").
  		Get().
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat ResourceInfoRes
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
  export interface ResourceInfoRes {
    resources?: Resource[];
    paths?: Endpoint[];
  }

  export interface Resource {
    name?: string;                 // resource name
    code?: string;                 // resource code, unique identifier
  }

  export interface Endpoint {
    type?: string;                 // access scope type: PROTECTED/PUBLIC
    url?: string;                  // endpoint url
    group?: string;                // app name
    desc?: string;                 // description of the endpoint
    resCode?: string;              // resource code
    method?: string;               // http method
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
    this.http.get<ResourceInfoRes>(`/logbot/auth/resource`)
      .subscribe({
        next: (resp) => {
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
  }
  ```

## GET /metrics

- Description: Collect prometheus metrics information
- Header Parameter:
  - "Authorization": Basic authorization if enabled
- cURL:
  ```sh
  curl -X GET 'http://localhost:8087/metrics' \
    -H 'Authorization: '
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  func SendRequest(rail miso.Rail, authorization string) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/metrics", "logbot").
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
    this.http.get<any>(`/logbot/metrics`,
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
