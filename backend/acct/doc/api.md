# API Endpoints

## Contents

- [POST /open/api/v1/cashflow/list](#post-openapiv1cashflowlist)
- [POST /open/api/v1/cashflow/import/wechat](#post-openapiv1cashflowimportwechat)
- [GET /open/api/v1/cashflow/list-currency](#get-openapiv1cashflowlist-currency)
- [POST /open/api/v1/cashflow/list-statistics](#post-openapiv1cashflowlist-statistics)
- [POST /open/api/v1/cashflow/plot-statistics](#post-openapiv1cashflowplot-statistics)
- [GET /auth/resource](#get-authresource)
- [GET /debug/trace/recorder/run](#get-debugtracerecorderrun)
- [GET /debug/trace/recorder/snapshot](#get-debugtracerecordersnapshot)
- [GET /debug/trace/recorder/stop](#get-debugtracerecorderstop)

## POST /open/api/v1/cashflow/list

- Bound to Resource: `"acct:ManageCashflows"`
- JSON Request:
    - "paging": (Paging) Paging
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
    - "direction": (string) Flow Direction: IN / OUT. Enums: ["IN","OUT",""].
    - "transTimeStart": (int64) Transaction Time Range Start
    - "transTimeEnd": (int64) Transaction Time Range End
    - "transId": (string) Transaction ID
    - "category": (string) Category Code
    - "minAmt": (*string) Minimum amount
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (PageRes[github.com/curtisnewbie/acct/internal/flow.ListCashFlowRes]) response data
      - "paging": (Paging) pagination parameters
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
      - "payload": ([]flow.ListCashFlowRes) payload values in current page
        - "direction": (string) Flow Direction: IN / OUT
        - "transTime": (int64) Transaction Time
        - "transId": (string) Transaction ID
        - "counterparty": (string) Counterparty of the transaction
        - "paymentMethod": (string) Payment Method
        - "amount": (string) Amount
        - "currency": (string) Currency
        - "extra": (string) Extra Information
        - "category": (string) Category Code
        - "categoryName": (string) Category Name
        - "remark": (string) Remark
        - "createdAt": (int64) Create Time
- cURL:
  ```sh
  curl -X POST 'http://localhost:8093/open/api/v1/cashflow/list' \
    -H 'Content-Type: application/json' \
    -d '{"category":"","direction":"","minAmt":"","paging":{"limit":0,"page":0,"total":0},"transId":"","transTimeEnd":0,"transTimeStart":0}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListCashFlowReq struct {
  	Paging miso.Paging `json:"paging"`
  	Direction string `json:"direction"` // Flow Direction: IN / OUT. Enums: ["IN","OUT",""].
  	TransTimeStart *atom.Time `json:"transTimeStart"` // Transaction Time Range Start
  	TransTimeEnd *atom.Time `json:"transTimeEnd"` // Transaction Time Range End
  	TransId string `json:"transId"` // Transaction ID
  	Category string `json:"category"` // Category Code
  	MinAmt *money.Amt `json:"minAmt"` // Minimum amount
  }


  type ListCashFlowRes struct {
  	Direction string `json:"direction"` // Flow Direction: IN / OUT
  	TransTime atom.Time `json:"transTime"` // Transaction Time
  	TransId string `json:"transId"` // Transaction ID
  	Counterparty string `json:"counterparty"` // Counterparty of the transaction
  	PaymentMethod string `json:"paymentMethod"` // Payment Method
  	Amount string `json:"amount"`  // Amount
  	Currency string `json:"currency"` // Currency
  	Extra string `json:"extra"`    // Extra Information
  	Category string `json:"category"` // Category Code
  	CategoryName string `json:"categoryName"` // Category Name
  	Remark string `json:"remark"`  // Remark
  	CreatedAt atom.Time `json:"createdAt"` // Create Time
  }

  func ApiListCashFlows(rail miso.Rail, req ListCashFlowReq) (miso.PageRes[ListCashFlowRes], error) {
  	var res miso.GnResp[miso.PageRes[ListCashFlowRes]]
  	err := miso.NewDynClient(rail, "/open/api/v1/cashflow/list", "acct").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat miso.PageRes[ListCashFlowRes]
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
  export interface ListCashFlowReq {
    paging?: Paging;
    direction?: string;            // Flow Direction: IN / OUT. Enums: ["IN","OUT",""].
    transTimeStart?: number;       // Transaction Time Range Start
    transTimeEnd?: number;         // Transaction Time Range End
    transId?: string;              // Transaction ID
    category?: string;             // Category Code
    minAmt?: string;               // Minimum amount
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
    payload?: ListCashFlowRes[];
  }

  export interface Paging {
    limit?: number;                // page limit
    page?: number;                 // page number, 1-based
    total?: number;                // total count
  }

  export interface ListCashFlowRes {
    direction?: string;            // Flow Direction: IN / OUT
    transTime?: number;            // Transaction Time
    transId?: string;              // Transaction ID
    counterparty?: string;         // Counterparty of the transaction
    paymentMethod?: string;        // Payment Method
    amount?: string;               // Amount
    currency?: string;             // Currency
    extra?: string;                // Extra Information
    category?: string;             // Category Code
    categoryName?: string;         // Category Name
    remark?: string;               // Remark
    createdAt?: number;            // Create Time
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

  listCashFlows() {
    let req: ListCashFlowReq | null = null;
    this.http.post<any>(`/acct/open/api/v1/cashflow/list`, req)
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

## POST /open/api/v1/cashflow/import/wechat

- Bound to Resource: `"acct:ManageCashflows"`
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
- cURL:
  ```sh
  curl -X POST 'http://localhost:8093/open/api/v1/cashflow/import/wechat'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  func ApiImportWechatCashflows(rail miso.Rail) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/open/api/v1/cashflow/import/wechat", "acct").
  		Post(nil).
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

- JSON Request / Response Object In TypeScript:
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

  importWechatCashflows() {
    this.http.post<any>(`/acct/open/api/v1/cashflow/import/wechat`, null)
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

## GET /open/api/v1/cashflow/list-currency

- Bound to Resource: `"acct:ManageCashflows"`
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": ([]string) response data
- cURL:
  ```sh
  curl -X GET 'http://localhost:8093/open/api/v1/cashflow/list-currency'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  func ApiListCurrency(rail miso.Rail) ([]string, error) {
  	var res miso.GnResp[[]string]
  	err := miso.NewDynClient(rail, "/open/api/v1/cashflow/list-currency", "acct").
  		Get().
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat []string
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
  export interface Resp {
    errorCode?: string;            // error code
    msg?: string;                  // message
    error?: boolean;               // whether the request was successful
    data?: string[];               // response data
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

  listCurrency() {
    this.http.get<any>(`/acct/open/api/v1/cashflow/list-currency`)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: string[] = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
  }
  ```

## POST /open/api/v1/cashflow/list-statistics

- Bound to Resource: `"acct:ManageCashflows"`
- JSON Request:
    - "paging": (Paging) Paging Info
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
    - "aggType": (string) Aggregation Type. Enums: ["YEARLY","MONTHLY","WEEKLY"].
    - "aggRange": (string) Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD).
    - "currency": (string) Currency
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (PageRes[github.com/curtisnewbie/acct/internal/flow.ApiListStatisticsRes]) response data
      - "paging": (Paging) pagination parameters
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
      - "payload": ([]flow.ApiListStatisticsRes) payload values in current page
        - "aggType": (string) Aggregation Type.
        - "aggRange": (string) Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD).
        - "aggValue": (string) Aggregation Value.
        - "currency": (string) Currency
- cURL:
  ```sh
  curl -X POST 'http://localhost:8093/open/api/v1/cashflow/list-statistics' \
    -H 'Content-Type: application/json' \
    -d '{"aggRange":"","aggType":"","currency":"","paging":{"limit":0,"page":0,"total":0}}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ApiListStatisticsReq struct {
  	Paging miso.Paging `json:"paging"`
  	AggType string `json:"aggType"` // Aggregation Type. Enums: ["YEARLY","MONTHLY","WEEKLY"].
  	AggRange string `json:"aggRange"` // Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD).
  	Currency string `json:"currency"` // Currency
  }


  type ApiListStatisticsRes struct {
  	AggType string `json:"aggType"` // Aggregation Type.
  	AggRange string `json:"aggRange"` // Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD).
  	AggValue string `json:"aggValue"` // Aggregation Value.
  	Currency string `json:"currency"` // Currency
  }

  func ApiListCashflowStatistics(rail miso.Rail, req ApiListStatisticsReq) (miso.PageRes[ApiListStatisticsRes], error) {
  	var res miso.GnResp[miso.PageRes[ApiListStatisticsRes]]
  	err := miso.NewDynClient(rail, "/open/api/v1/cashflow/list-statistics", "acct").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat miso.PageRes[ApiListStatisticsRes]
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
  export interface ApiListStatisticsReq {
    paging?: Paging;
    aggType?: string;              // Aggregation Type. Enums: ["YEARLY","MONTHLY","WEEKLY"].
    aggRange?: string;             // Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD).
    currency?: string;             // Currency
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
    payload?: ApiListStatisticsRes[];
  }

  export interface Paging {
    limit?: number;                // page limit
    page?: number;                 // page number, 1-based
    total?: number;                // total count
  }

  export interface ApiListStatisticsRes {
    aggType?: string;              // Aggregation Type.
    aggRange?: string;             // Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD).
    aggValue?: string;             // Aggregation Value.
    currency?: string;             // Currency
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

  listCashflowStatistics() {
    let req: ApiListStatisticsReq | null = null;
    this.http.post<any>(`/acct/open/api/v1/cashflow/list-statistics`, req)
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

## POST /open/api/v1/cashflow/plot-statistics

- Bound to Resource: `"acct:ManageCashflows"`
- JSON Request:
    - "startTime": (int64) Start time
    - "endTime": (int64) End time
    - "aggType": (string) Aggregation Type. Enums: ["YEARLY","MONTHLY","WEEKLY"].
    - "currency": (string) Currency
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": ([]flow.ApiPlotStatisticsRes) response data
      - "aggRange": (string) Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD).
      - "aggValue": (string) Aggregation Value.
- cURL:
  ```sh
  curl -X POST 'http://localhost:8093/open/api/v1/cashflow/plot-statistics' \
    -H 'Content-Type: application/json' \
    -d '{"aggType":"","currency":"","endTime":0,"startTime":0}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ApiPlotStatisticsReq struct {
  	StartTime atom.Time `json:"startTime"` // Start time
  	EndTime atom.Time `json:"endTime"` // End time
  	AggType string `json:"aggType"` // Aggregation Type. Enums: ["YEARLY","MONTHLY","WEEKLY"].
  	Currency string `json:"currency"` // Currency
  }

  type ApiPlotStatisticsRes struct {
  	AggRange string `json:"aggRange"` // Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD).
  	AggValue string `json:"aggValue"` // Aggregation Value.
  }

  func ApiPlotCashflowStatistics(rail miso.Rail, req ApiPlotStatisticsReq) ([]ApiPlotStatisticsRes, error) {
  	var res miso.GnResp[[]ApiPlotStatisticsRes]
  	err := miso.NewDynClient(rail, "/open/api/v1/cashflow/plot-statistics", "acct").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat []ApiPlotStatisticsRes
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
  export interface ApiPlotStatisticsReq {
    startTime?: number;            // Start time
    endTime?: number;              // End time
    aggType?: string;              // Aggregation Type. Enums: ["YEARLY","MONTHLY","WEEKLY"].
    currency?: string;             // Currency
  }

  export interface Resp {
    errorCode?: string;            // error code
    msg?: string;                  // message
    error?: boolean;               // whether the request was successful
    data?: ApiPlotStatisticsRes[];
  }

  export interface ApiPlotStatisticsRes {
    aggRange?: string;             // Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD).
    aggValue?: string;             // Aggregation Value.
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

  plotCashflowStatistics() {
    let req: ApiPlotStatisticsReq | null = null;
    this.http.post<any>(`/acct/open/api/v1/cashflow/plot-statistics`, req)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: ApiPlotStatisticsRes[] = resp.data;
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
  curl -X GET 'http://localhost:8093/auth/resource'
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
  	err := miso.NewDynClient(rail, "/auth/resource", "acct").
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
    this.http.get<ResourceInfoRes>(`/acct/auth/resource`)
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
  curl -X GET 'http://localhost:8093/debug/trace/recorder/run?duration='
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Start FlightRecorder. Recorded result is written to trace.out when it's finished or stopped.
  func SendRequest(rail miso.Rail, duration string) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/debug/trace/recorder/run", "acct").
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
    this.http.get<any>(`/acct/debug/trace/recorder/run?duration=${duration}`)
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
  curl -X GET 'http://localhost:8093/debug/trace/recorder/snapshot'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // FlightRecorder take snapshot. Recorded result is written to trace.out.
  func SendRequest(rail miso.Rail) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/debug/trace/recorder/snapshot", "acct").
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
    this.http.get<any>(`/acct/debug/trace/recorder/snapshot`)
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
  curl -X GET 'http://localhost:8093/debug/trace/recorder/stop'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Stop existing FlightRecorder session.
  func SendRequest(rail miso.Rail) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/debug/trace/recorder/stop", "acct").
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
    this.http.get<any>(`/acct/debug/trace/recorder/stop`)
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
