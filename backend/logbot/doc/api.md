# API Endpoints

## Contents

- [POST /log/error/list](#post-logerrorlist)
- [GET /auth/resource](#get-authresource)
- [GET /debug/trace/recorder/run](#get-debugtracerecorderrun)
- [GET /debug/trace/recorder/snapshot](#get-debugtracerecordersnapshot)
- [GET /debug/trace/recorder/stop](#get-debugtracerecorderstop)

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
    - "data": (ListErrorLogResp) response data
      - "page": (Paging) 
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
      - "payload": ([]logbot.ListedErrorLog) 
        - "id": (int64) 
        - "node": (string) 
        - "app": (string) 
        - "caller": (string) 
        - "traceId": (string) 
        - "spanId": (string) 
        - "errMsg": (string) 
        - "rtime": (int64) 
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

  type ListErrorLogResp struct {
  	Page miso.Paging `json:"page"`
  	Payload []ListedErrorLog `json:"payload"`
  }

  type ListedErrorLog struct {
  	Id int64 `json:"id"`
  	Node string `json:"node"`
  	App string `json:"app"`
  	Caller string `json:"caller"`
  	TraceId string `json:"traceId"`
  	SpanId string `json:"spanId"`
  	ErrMsg string `json:"errMsg"`
  	RTime atom.Time `json:"rtime"`
  }

  // List error logs
  func SendListErrorLogReq(rail miso.Rail, req ListErrorLogReq) (ListErrorLogResp, error) {
  	var res miso.GnResp[ListErrorLogResp]
  	err := miso.NewDynClient(rail, "/log/error/list", "logbot").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat ListErrorLogResp
  		return dat, err
  	}
  	dat, err := res.Res()
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  	}
  	return dat, err
  }
  ```

- JSON Request / Response Object In TypeScript:
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

  export interface Resp {
    errorCode?: string;            // error code
    msg?: string;                  // message
    error?: boolean;               // whether the request was successful
    data?: ListErrorLogResp;
  }

  export interface ListErrorLogResp {
    page?: Paging;
    payload?: ListedErrorLog[];
  }

  export interface Paging {
    limit?: number;                // page limit
    page?: number;                 // page number, 1-based
    total?: number;                // total count
  }

  export interface ListedErrorLog {
    id?: number;
    node?: string;
    app?: string;
    caller?: string;
    traceId?: string;
    spanId?: string;
    errMsg?: string;
    rtime?: number;
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
          let dat: ListErrorLogResp = resp.data;
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
  	Resources []Resource `json:"resources"`
  	Paths []Endpoint `json:"paths"`
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

  // Expose resource and endpoint information to other backend service for authorization.
  func SendRequest(rail miso.Rail) (ResourceInfoRes, error) {
  	var res miso.GnResp[ResourceInfoRes]
  	err := miso.NewDynClient(rail, "/auth/resource", "logbot").
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

- JSON Request / Response Object In TypeScript:
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

## GET /debug/trace/recorder/run

- Description: Start FlightRecorder. Recorded result is written to trace.out when it's finished or stopped.
- Query Parameter:
  - "duration": Duration of the flight recording. Required. Duration cannot exceed 30 min.
- cURL:
  ```sh
  curl -X GET 'http://localhost:8087/debug/trace/recorder/run?duration='
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Start FlightRecorder. Recorded result is written to trace.out when it's finished or stopped.
  func SendRequest(rail miso.Rail, duration string) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/debug/trace/recorder/run", "logbot").
  		AddQueryParams("duration", duration).
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
    let duration: any | null = null;
    this.http.get<any>(`/logbot/debug/trace/recorder/run?duration=${duration}`)
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
  curl -X GET 'http://localhost:8087/debug/trace/recorder/snapshot'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // FlightRecorder take snapshot. Recorded result is written to trace.out.
  func SendRequest(rail miso.Rail) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/debug/trace/recorder/snapshot", "logbot").
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
    this.http.get<any>(`/logbot/debug/trace/recorder/snapshot`)
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
  curl -X GET 'http://localhost:8087/debug/trace/recorder/stop'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Stop existing FlightRecorder session.
  func SendRequest(rail miso.Rail) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/debug/trace/recorder/stop", "logbot").
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
    this.http.get<any>(`/logbot/debug/trace/recorder/stop`)
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
