# API Endpoints

- POST /open/api/v1/cashflow/list
  - Bound to Resource: `"acct:ManageCashflows"`
  - JSON Request:
    - "paging": (Paging) Paging
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
    - "direction": (string) Flow Direction: IN / OUT
    - "transTimeStart": (int64) Transaction Time Range Start
    - "transTimeEnd": (int64) Transaction Time Range End
    - "transId": (string) Transaction ID
    - "category": (string) Category Code
    - "minAmt": (string) Minimum amount
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

  - JSON Request Object In TypeScript:
    ```ts
    export interface ListCashFlowReq {
      paging?: Paging
      direction?: string             // Flow Direction: IN / OUT
      transTimeStart?: number        // Transaction Time Range Start
      transTimeEnd?: number          // Transaction Time Range End
      transId?: string               // Transaction ID
      category?: string              // Category Code
      minAmt?: string                // Minimum amount
    }

    export interface Paging {
      limit?: number                 // page limit
      page?: number                  // page number, 1-based
      total?: number                 // total count
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string             // error code
      msg?: string                   // message
      error?: boolean                // whether the request was successful
      data?: PageRes
    }

    export interface PageRes {
      paging?: Paging
      payload?: ListCashFlowRes[]
    }

    export interface Paging {
      limit?: number                 // page limit
      page?: number                  // page number, 1-based
      total?: number                 // total count
    }

    export interface ListCashFlowRes {
      direction?: string             // Flow Direction: IN / OUT
      transTime?: number             // Transaction Time
      transId?: string               // Transaction ID
      counterparty?: string          // Counterparty of the transaction
      paymentMethod?: string         // Payment Method
      amount?: string                // Amount
      currency?: string              // Currency
      extra?: string                 // Extra Information
      category?: string              // Category Code
      categoryName?: string          // Category Name
      remark?: string                // Remark
      createdAt?: number             // Create Time
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
    ```

- POST /open/api/v1/cashflow/import/wechat
  - Bound to Resource: `"acct:ManageCashflows"`
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8093/open/api/v1/cashflow/import/wechat'
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string             // error code
      msg?: string                   // message
      error?: boolean                // whether the request was successful
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

    this.http.post<any>(`/acct/open/api/v1/cashflow/import/wechat`)
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
    ```

- GET /open/api/v1/cashflow/list-currency
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

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string             // error code
      msg?: string                   // message
      error?: boolean                // whether the request was successful
      data?: string[]                // response data
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
    ```

- POST /open/api/v1/cashflow/list-statistics
  - Bound to Resource: `"acct:ManageCashflows"`
  - JSON Request:
    - "paging": (Paging) Paging Info
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
    - "aggType": (string) Aggregation Type.
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

  - JSON Request Object In TypeScript:
    ```ts
    export interface ApiListStatisticsReq {
      paging?: Paging
      aggType?: string               // Aggregation Type.
      aggRange?: string              // Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD).
      currency?: string              // Currency
    }

    export interface Paging {
      limit?: number                 // page limit
      page?: number                  // page number, 1-based
      total?: number                 // total count
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string             // error code
      msg?: string                   // message
      error?: boolean                // whether the request was successful
      data?: PageRes
    }

    export interface PageRes {
      paging?: Paging
      payload?: ApiListStatisticsRes[]
    }

    export interface Paging {
      limit?: number                 // page limit
      page?: number                  // page number, 1-based
      total?: number                 // total count
    }

    export interface ApiListStatisticsRes {
      aggType?: string               // Aggregation Type.
      aggRange?: string              // Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD).
      aggValue?: string              // Aggregation Value.
      currency?: string              // Currency
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
    ```

- POST /open/api/v1/cashflow/plot-statistics
  - Bound to Resource: `"acct:ManageCashflows"`
  - JSON Request:
    - "startTime": (int64) Start time
    - "endTime": (int64) End time
    - "aggType": (string) Aggregation Type.
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

  - JSON Request Object In TypeScript:
    ```ts
    export interface ApiPlotStatisticsReq {
      startTime?: number             // Start time
      endTime?: number               // End time
      aggType?: string               // Aggregation Type.
      currency?: string              // Currency
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string             // error code
      msg?: string                   // message
      error?: boolean                // whether the request was successful
      data?: ApiPlotStatisticsRes[]
    }

    export interface ApiPlotStatisticsRes {
      aggRange?: string              // Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD).
      aggValue?: string              // Aggregation Value.
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
    ```

- GET /auth/resource
  - Description: Expose resource and endpoint information to other backend service for authorization.
  - Expected Access Scope: PROTECTED
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (ResourceInfoRes) response data
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

  - JSON Response Object In TypeScript:
    ```ts
    export interface GnResp {
      errorCode?: string             // error code
      msg?: string                   // message
      error?: boolean                // whether the request was successful
      data?: ResourceInfoRes
    }

    export interface ResourceInfoRes {
      resources?: Resource[]
      paths?: Endpoint[]
    }

    export interface Resource {
      name?: string                  // resource name
      code?: string                  // resource code, unique identifier
    }

    export interface Endpoint {
      type?: string                  // access scope type: PROTECTED/PUBLIC
      url?: string                   // endpoint url
      group?: string                 // app name
      desc?: string                  // description of the endpoint
      resCode?: string               // resource code
      method?: string                // http method
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

    this.http.get<any>(`/acct/auth/resource`)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: ResourceInfoRes = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```

- GET /metrics
  - Description: Collect prometheus metrics information
  - Header Parameter:
    - "Authorization": Basic authorization if enabled
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8093/metrics' \
      -H 'Authorization: '
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    let authorization: any | null = null;
    this.http.get<any>(`/acct/metrics`,
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
    ```

- GET /debug/pprof
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8093/debug/pprof'
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    this.http.get<any>(`/acct/debug/pprof`)
      .subscribe({
        next: () => {
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```

- GET /debug/pprof/:name
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8093/debug/pprof/:name'
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    this.http.get<any>(`/acct/debug/pprof/:name`)
      .subscribe({
        next: () => {
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```

- GET /debug/pprof/cmdline
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8093/debug/pprof/cmdline'
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    this.http.get<any>(`/acct/debug/pprof/cmdline`)
      .subscribe({
        next: () => {
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```

- GET /debug/pprof/profile
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8093/debug/pprof/profile'
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    this.http.get<any>(`/acct/debug/pprof/profile`)
      .subscribe({
        next: () => {
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```

- GET /debug/pprof/symbol
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8093/debug/pprof/symbol'
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    this.http.get<any>(`/acct/debug/pprof/symbol`)
      .subscribe({
        next: () => {
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```

- GET /debug/pprof/trace
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8093/debug/pprof/trace'
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    this.http.get<any>(`/acct/debug/pprof/trace`)
      .subscribe({
        next: () => {
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```

- GET /doc/api
  - Description: Serve the generated API documentation webpage
  - Expected Access Scope: PUBLIC
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8093/doc/api'
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    this.http.get<any>(`/acct/doc/api`)
      .subscribe({
        next: () => {
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```