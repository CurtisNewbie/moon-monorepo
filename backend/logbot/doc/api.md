# API Endpoints

## Contents

- [POST /log/error/list](#post-logerrorlist)

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
  		var dat ListErrorLogResp
  		return dat, err
  	}
  	return res.Data, nil
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

# Event Pipelines

- ReportLogPipeline
  - RabbitMQ Queue: `logbot:error-log:report:pipeline`
  - RabbitMQ Exchange: `logbot:error-log:report:pipeline`
  - RabbitMQ RoutingKey: `#`
  - Event Payload:
    - "Node": (string) 
    - "App": (string) 
    - "Time": (int64) 
    - "TraceId": (string) 
    - "SpanId": (string) 
    - "FuncName": (string) 
    - "Message": (string) 
