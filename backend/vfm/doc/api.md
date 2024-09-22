# API Endpoints

- GET /open/api/file/upload/duplication/preflight
  - Description: Preflight check for duplicate file uploads
  - Bound to Resource: `"manage-files"`
  - Query Parameter:
    - "fileName": 
    - "parentFileKey": 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (bool) response data
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8086/open/api/file/upload/duplication/preflight?fileName=&parentFileKey='
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string             // error code
      msg?: string                   // message
      error?: boolean                // whether the request was successful
      data?: boolean                 // response data
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

    let fileName: any | null = null;
    let parentFileKey: any | null = null;
    this.http.get<any>(`/vfm/open/api/file/upload/duplication/preflight?fileName=${fileName}&parentFileKey=${parentFileKey}`)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: boolean = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```

- GET /open/api/file/parent
  - Description: User fetch parent file info
  - Bound to Resource: `"manage-files"`
  - Query Parameter:
    - "fileKey": 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (*vfm.ParentFileInfo) response data
      - "fileKey": (string) 
      - "fileName": (string) 
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8086/open/api/file/parent?fileKey='
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string             // error code
      msg?: string                   // message
      error?: boolean                // whether the request was successful
      data?: ParentFileInfo
    }

    export interface ParentFileInfo {
      fileKey?: string
      fileName?: string
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

    let fileKey: any | null = null;
    this.http.get<any>(`/vfm/open/api/file/parent?fileKey=${fileKey}`)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: ParentFileInfo = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```

- POST /open/api/file/move-to-dir
  - Description: User move file into directory
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "uuid": (string) 
    - "parentFileUuid": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/file/move-to-dir' \
      -H 'Content-Type: application/json' \
      -d '{"parentFileUuid":"","uuid":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface MoveIntoDirReq {
      uuid?: string
      parentFileUuid?: string
    }
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

    let req: MoveIntoDirReq | null = null;
    this.http.post<any>(`/vfm/open/api/file/move-to-dir`, req)
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

- POST /open/api/file/batch-move-to-dir
  - Description: User move files into directory
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "instructions": ([]vfm.MoveIntoDirReq) 
      - "uuid": (string) 
      - "parentFileUuid": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/file/batch-move-to-dir' \
      -H 'Content-Type: application/json' \
      -d '{"instructions":{"parentFileUuid":"","uuid":""}}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface BatchMoveIntoDirReq {
      instructions?: MoveIntoDirReq[]
    }

    export interface MoveIntoDirReq {
      uuid?: string
      parentFileUuid?: string
    }
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

    let req: BatchMoveIntoDirReq | null = null;
    this.http.post<any>(`/vfm/open/api/file/batch-move-to-dir`, req)
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

- POST /open/api/file/make-dir
  - Description: User make directory
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "parentFile": (string) 
    - "name": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (string) response data
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/file/make-dir' \
      -H 'Content-Type: application/json' \
      -d '{"name":"","parentFile":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface MakeDirReq {
      parentFile?: string
      name?: string
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string             // error code
      msg?: string                   // message
      error?: boolean                // whether the request was successful
      data?: string                  // response data
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

    let req: MakeDirReq | null = null;
    this.http.post<any>(`/vfm/open/api/file/make-dir`, req)
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
    ```

- GET /open/api/file/dir/list
  - Description: User list directories
  - Bound to Resource: `"manage-files"`
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": ([]vfm.ListedDir) response data
      - "id": (int) 
      - "uuid": (string) 
      - "name": (string) 
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8086/open/api/file/dir/list'
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string             // error code
      msg?: string                   // message
      error?: boolean                // whether the request was successful
      data?: ListedDir[]
    }

    export interface ListedDir {
      id?: number
      uuid?: string
      name?: string
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

    this.http.get<any>(`/vfm/open/api/file/dir/list`)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: ListedDir[] = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```

- POST /open/api/file/list
  - Description: User list files
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "paging": (Paging) 
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
    - "filename": (*string) 
    - "folderNo": (*string) 
    - "fileType": (*string) 
    - "parentFile": (*string) 
    - "sensitive": (*bool) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (PageRes[github.com/curtisnewbie/vfm/internal/vfm.ListedFile]) response data
      - "paging": (Paging) pagination parameters
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
      - "payload": ([]vfm.ListedFile) payload values in current page
        - "id": (int) 
        - "uuid": (string) 
        - "name": (string) 
        - "uploadTime": (int64) 
        - "uploaderName": (string) 
        - "sizeInBytes": (int64) 
        - "fileType": (string) 
        - "updateTime": (int64) 
        - "parentFileName": (string) 
        - "sensitiveMode": (string) 
        - "thumbnailToken": (string) 
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/file/list' \
      -H 'Content-Type: application/json' \
      -d '{"fileType":"","filename":"","folderNo":"","paging":{"limit":0,"page":0,"total":0},"parentFile":"","sensitive":false}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ListFileReq {
      paging?: Paging
      filename?: string
      folderNo?: string
      fileType?: string
      parentFile?: string
      sensitive?: boolean
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
      payload?: ListedFile[]
    }

    export interface Paging {
      limit?: number                 // page limit
      page?: number                  // page number, 1-based
      total?: number                 // total count
    }

    export interface ListedFile {
      id?: number
      uuid?: string
      name?: string
      uploadTime?: number
      uploaderName?: string
      sizeInBytes?: number
      fileType?: string
      updateTime?: number
      parentFileName?: string
      sensitiveMode?: string
      thumbnailToken?: string
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

    let req: ListFileReq | null = null;
    this.http.post<any>(`/vfm/open/api/file/list`, req)
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

- POST /open/api/file/delete
  - Description: User delete file
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "uuid": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/file/delete' \
      -H 'Content-Type: application/json' \
      -d '{"uuid":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface DeleteFileReq {
      uuid?: string
    }
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

    let req: DeleteFileReq | null = null;
    this.http.post<any>(`/vfm/open/api/file/delete`, req)
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

- POST /open/api/file/dir/truncate
  - Description: User delete truncate directory recursively
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "uuid": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/file/dir/truncate' \
      -H 'Content-Type: application/json' \
      -d '{"uuid":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface DeleteFileReq {
      uuid?: string
    }
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

    let req: DeleteFileReq | null = null;
    this.http.post<any>(`/vfm/open/api/file/dir/truncate`, req)
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

- POST /open/api/file/delete/batch
  - Description: User delete file in batch
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "fileKeys": ([]string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/file/delete/batch' \
      -H 'Content-Type: application/json' \
      -d '{"fileKeys":[]}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface BatchDeleteFileReq {
      fileKeys?: string[]
    }
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

    let req: BatchDeleteFileReq | null = null;
    this.http.post<any>(`/vfm/open/api/file/delete/batch`, req)
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

- POST /open/api/file/create
  - Description: User create file
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "filename": (string) 
    - "fstoreFileId": (string) 
    - "parentFile": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/file/create' \
      -H 'Content-Type: application/json' \
      -d '{"filename":"","fstoreFileId":"","parentFile":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface CreateFileReq {
      filename?: string
      fstoreFileId?: string
      parentFile?: string
    }
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

    let req: CreateFileReq | null = null;
    this.http.post<any>(`/vfm/open/api/file/create`, req)
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

- POST /open/api/file/info/update
  - Description: User update file
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "id": (int) 
    - "name": (string) 
    - "sensitiveMode": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/file/info/update' \
      -H 'Content-Type: application/json' \
      -d '{"id":0,"name":"","sensitiveMode":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface UpdateFileReq {
      id?: number
      name?: string
      sensitiveMode?: string
    }
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

    let req: UpdateFileReq | null = null;
    this.http.post<any>(`/vfm/open/api/file/info/update`, req)
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

- POST /open/api/file/token/generate
  - Description: User generate temporary token
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "fileKey": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (string) response data
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/file/token/generate' \
      -H 'Content-Type: application/json' \
      -d '{"fileKey":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface GenerateTempTokenReq {
      fileKey?: string
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string             // error code
      msg?: string                   // message
      error?: boolean                // whether the request was successful
      data?: string                  // response data
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

    let req: GenerateTempTokenReq | null = null;
    this.http.post<any>(`/vfm/open/api/file/token/generate`, req)
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
    ```

- POST /open/api/file/unpack
  - Description: User unpack zip
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "fileKey": (string) 
    - "parentFileKey": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/file/unpack' \
      -H 'Content-Type: application/json' \
      -d '{"fileKey":"","parentFileKey":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface UnpackZipReq {
      fileKey?: string
      parentFileKey?: string
    }
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

    let req: UnpackZipReq | null = null;
    this.http.post<any>(`/vfm/open/api/file/unpack`, req)
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

- GET /open/api/file/token/qrcode
  - Description: User generate qrcode image for temporary token
  - Expected Access Scope: PUBLIC
  - Query Parameter:
    - "token": Generated temporary file key
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8086/open/api/file/token/qrcode?token='
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    let token: any | null = null;
    this.http.get<any>(`/vfm/open/api/file/token/qrcode?token=${token}`)
      .subscribe({
        next: () => {
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```

- GET /open/api/vfolder/brief/owned
  - Description: User list virtual folder briefs
  - Bound to Resource: `"manage-files"`
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": ([]vfm.VFolderBrief) response data
      - "folderNo": (string) 
      - "name": (string) 
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8086/open/api/vfolder/brief/owned'
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string             // error code
      msg?: string                   // message
      error?: boolean                // whether the request was successful
      data?: VFolderBrief[]
    }

    export interface VFolderBrief {
      folderNo?: string
      name?: string
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

    this.http.get<any>(`/vfm/open/api/vfolder/brief/owned`)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: VFolderBrief[] = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```

- POST /open/api/vfolder/list
  - Description: User list virtual folders
  - Bound to Resource: `"manage-files"`
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
    - "data": (ListVFolderRes) response data
      - "paging": (Paging) 
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
      - "payload": ([]vfm.ListedVFolder) 
        - "id": (int) 
        - "folderNo": (string) 
        - "name": (string) 
        - "createTime": (int64) 
        - "createBy": (string) 
        - "updateTime": (int64) 
        - "updateBy": (string) 
        - "ownership": (string) 
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/vfolder/list' \
      -H 'Content-Type: application/json' \
      -d '{"name":"","paging":{"limit":0,"page":0,"total":0}}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ListVFolderReq {
      paging?: Paging
      name?: string
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
      data?: ListVFolderRes
    }

    export interface ListVFolderRes {
      paging?: Paging
      payload?: ListedVFolder[]
    }

    export interface Paging {
      limit?: number                 // page limit
      page?: number                  // page number, 1-based
      total?: number                 // total count
    }

    export interface ListedVFolder {
      id?: number
      folderNo?: string
      name?: string
      createTime?: number
      createBy?: string
      updateTime?: number
      updateBy?: string
      ownership?: string
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

    let req: ListVFolderReq | null = null;
    this.http.post<any>(`/vfm/open/api/vfolder/list`, req)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: ListVFolderRes = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```

- POST /open/api/vfolder/create
  - Description: User create virtual folder
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "name": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (string) response data
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/vfolder/create' \
      -H 'Content-Type: application/json' \
      -d '{"name":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface CreateVFolderReq {
      name?: string
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string             // error code
      msg?: string                   // message
      error?: boolean                // whether the request was successful
      data?: string                  // response data
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

    let req: CreateVFolderReq | null = null;
    this.http.post<any>(`/vfm/open/api/vfolder/create`, req)
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
    ```

- POST /open/api/vfolder/file/add
  - Description: User add file to virtual folder
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "folderNo": (string) 
    - "fileKeys": ([]string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/vfolder/file/add' \
      -H 'Content-Type: application/json' \
      -d '{"fileKeys":[],"folderNo":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface AddFileToVfolderReq {
      folderNo?: string
      fileKeys?: string[]
    }
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

    let req: AddFileToVfolderReq | null = null;
    this.http.post<any>(`/vfm/open/api/vfolder/file/add`, req)
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

- POST /open/api/vfolder/file/remove
  - Description: User remove file from virtual folder
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "folderNo": (string) 
    - "fileKeys": ([]string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/vfolder/file/remove' \
      -H 'Content-Type: application/json' \
      -d '{"fileKeys":[],"folderNo":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface RemoveFileFromVfolderReq {
      folderNo?: string
      fileKeys?: string[]
    }
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

    let req: RemoveFileFromVfolderReq | null = null;
    this.http.post<any>(`/vfm/open/api/vfolder/file/remove`, req)
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

- POST /open/api/vfolder/share
  - Description: Share access to virtual folder
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "folderNo": (string) 
    - "username": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/vfolder/share' \
      -H 'Content-Type: application/json' \
      -d '{"folderNo":"","username":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ShareVfolderReq {
      folderNo?: string
      username?: string
    }
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

    let req: ShareVfolderReq | null = null;
    this.http.post<any>(`/vfm/open/api/vfolder/share`, req)
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

- POST /open/api/vfolder/access/remove
  - Description: Remove granted access to virtual folder
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "folderNo": (string) 
    - "userNo": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/vfolder/access/remove' \
      -H 'Content-Type: application/json' \
      -d '{"folderNo":"","userNo":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface RemoveGrantedFolderAccessReq {
      folderNo?: string
      userNo?: string
    }
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

    let req: RemoveGrantedFolderAccessReq | null = null;
    this.http.post<any>(`/vfm/open/api/vfolder/access/remove`, req)
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

- POST /open/api/vfolder/granted/list
  - Description: List granted access to virtual folder
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "paging": (Paging) 
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
    - "folderNo": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (ListGrantedFolderAccessRes) response data
      - "paging": (Paging) 
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
      - "payload": ([]vfm.ListedFolderAccess) 
        - "userNo": (string) 
        - "username": (string) 
        - "createTime": (int64) 
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/vfolder/granted/list' \
      -H 'Content-Type: application/json' \
      -d '{"folderNo":"","paging":{"limit":0,"page":0,"total":0}}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ListGrantedFolderAccessReq {
      paging?: Paging
      folderNo?: string
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
      data?: ListGrantedFolderAccessRes
    }

    export interface ListGrantedFolderAccessRes {
      paging?: Paging
      payload?: ListedFolderAccess[]
    }

    export interface Paging {
      limit?: number                 // page limit
      page?: number                  // page number, 1-based
      total?: number                 // total count
    }

    export interface ListedFolderAccess {
      userNo?: string
      username?: string
      createTime?: number
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

    let req: ListGrantedFolderAccessReq | null = null;
    this.http.post<any>(`/vfm/open/api/vfolder/granted/list`, req)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: ListGrantedFolderAccessRes = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```

- POST /open/api/vfolder/remove
  - Description: Remove virtual folder
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "folderNo": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/vfolder/remove' \
      -H 'Content-Type: application/json' \
      -d '{"folderNo":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface RemoveVFolderReq {
      folderNo?: string
    }
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

    let req: RemoveVFolderReq | null = null;
    this.http.post<any>(`/vfm/open/api/vfolder/remove`, req)
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

- GET /open/api/gallery/brief/owned
  - Description: List owned gallery brief info
  - Bound to Resource: `"manage-files"`
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": ([]vfm.VGalleryBrief) response data
      - "galleryNo": (string) 
      - "name": (string) 
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8086/open/api/gallery/brief/owned'
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string             // error code
      msg?: string                   // message
      error?: boolean                // whether the request was successful
      data?: VGalleryBrief[]
    }

    export interface VGalleryBrief {
      galleryNo?: string
      name?: string
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

    this.http.get<any>(`/vfm/open/api/gallery/brief/owned`)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: VGalleryBrief[] = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```

- POST /open/api/gallery/new
  - Description: Create new gallery
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "name": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (*vfm.Gallery) response data
      - "id": (int64) 
      - "galleryNo": (string) 
      - "userNo": (string) 
      - "name": (string) 
      - "dirFileKey": (string) 
      - "createTime": (int64) 
      - "createBy": (string) 
      - "updateTime": (int64) 
      - "updateBy": (string) 
      - "isDel": (bool) 
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/gallery/new' \
      -H 'Content-Type: application/json' \
      -d '{"name":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface CreateGalleryCmd {
      name?: string
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string             // error code
      msg?: string                   // message
      error?: boolean                // whether the request was successful
      data?: Gallery
    }

    export interface Gallery {
      id?: number
      galleryNo?: string
      userNo?: string
      name?: string
      dirFileKey?: string
      createTime?: number
      createBy?: string
      updateTime?: number
      updateBy?: string
      isDel?: boolean
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

    let req: CreateGalleryCmd | null = null;
    this.http.post<any>(`/vfm/open/api/gallery/new`, req)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: Gallery = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```

- POST /open/api/gallery/update
  - Description: Update gallery
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "galleryNo": (string) 
    - "name": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/gallery/update' \
      -H 'Content-Type: application/json' \
      -d '{"galleryNo":"","name":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface UpdateGalleryCmd {
      galleryNo?: string
      name?: string
    }
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

    let req: UpdateGalleryCmd | null = null;
    this.http.post<any>(`/vfm/open/api/gallery/update`, req)
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

- POST /open/api/gallery/delete
  - Description: Delete gallery
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "galleryNo": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/gallery/delete' \
      -H 'Content-Type: application/json' \
      -d '{"galleryNo":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface DeleteGalleryCmd {
      galleryNo?: string
    }
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

    let req: DeleteGalleryCmd | null = null;
    this.http.post<any>(`/vfm/open/api/gallery/delete`, req)
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

- POST /open/api/gallery/list
  - Description: List galleries
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "paging": (Paging) 
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (PageRes[github.com/curtisnewbie/vfm/internal/vfm.VGallery]) response data
      - "paging": (Paging) pagination parameters
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
      - "payload": ([]vfm.VGallery) payload values in current page
        - "id": (int64) 
        - "galleryNo": (string) 
        - "userNo": (string) 
        - "name": (string) 
        - "createBy": (string) 
        - "updateBy": (string) 
        - "isOwner": (bool) 
        - "createTime": (string) 
        - "updateTime": (string) 
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/gallery/list' \
      -H 'Content-Type: application/json' \
      -d '{"paging":{"limit":0,"page":0,"total":0}}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ListGalleriesCmd {
      paging?: Paging
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
      payload?: VGallery[]
    }

    export interface Paging {
      limit?: number                 // page limit
      page?: number                  // page number, 1-based
      total?: number                 // total count
    }

    export interface VGallery {
      id?: number
      galleryNo?: string
      userNo?: string
      name?: string
      createBy?: string
      updateBy?: string
      isOwner?: boolean
      createTime?: string
      updateTime?: string
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

    let req: ListGalleriesCmd | null = null;
    this.http.post<any>(`/vfm/open/api/gallery/list`, req)
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

- POST /open/api/gallery/access/grant
  - Description: Grant access to the galleries
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "galleryNo": (string) 
    - "username": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/gallery/access/grant' \
      -H 'Content-Type: application/json' \
      -d '{"galleryNo":"","username":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface PermitGalleryAccessCmd {
      galleryNo?: string
      username?: string
    }
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

    let req: PermitGalleryAccessCmd | null = null;
    this.http.post<any>(`/vfm/open/api/gallery/access/grant`, req)
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

- POST /open/api/gallery/access/remove
  - Description: Remove access to the galleries
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "galleryNo": (string) 
    - "userNo": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/gallery/access/remove' \
      -H 'Content-Type: application/json' \
      -d '{"galleryNo":"","userNo":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface RemoveGalleryAccessCmd {
      galleryNo?: string
      userNo?: string
    }
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

    let req: RemoveGalleryAccessCmd | null = null;
    this.http.post<any>(`/vfm/open/api/gallery/access/remove`, req)
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

- POST /open/api/gallery/access/list
  - Description: List granted access to the galleries
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "galleryNo": (string) 
    - "paging": (Paging) 
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (PageRes[github.com/curtisnewbie/vfm/internal/vfm.ListedGalleryAccessRes]) response data
      - "paging": (Paging) pagination parameters
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
      - "payload": ([]vfm.ListedGalleryAccessRes) payload values in current page
        - "id": (int) 
        - "galleryNo": (string) 
        - "userNo": (string) 
        - "username": (string) 
        - "createTime": (int64) 
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/gallery/access/list' \
      -H 'Content-Type: application/json' \
      -d '{"galleryNo":"","paging":{"limit":0,"page":0,"total":0}}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ListGrantedGalleryAccessCmd {
      galleryNo?: string
      paging?: Paging
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
      payload?: ListedGalleryAccessRes[]
    }

    export interface Paging {
      limit?: number                 // page limit
      page?: number                  // page number, 1-based
      total?: number                 // total count
    }

    export interface ListedGalleryAccessRes {
      id?: number
      galleryNo?: string
      userNo?: string
      username?: string
      createTime?: number
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

    let req: ListGrantedGalleryAccessCmd | null = null;
    this.http.post<any>(`/vfm/open/api/gallery/access/list`, req)
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

- POST /open/api/gallery/images
  - Description: List images of gallery
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "galleryNo": (string) 
    - "paging": (Paging) 
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (*vfm.ListGalleryImagesResp) response data
      - "images": ([]vfm.ImageInfo) 
        - "thumbnailToken": (string) 
        - "fileTempToken": (string) 
      - "paging": (Paging) 
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/gallery/images' \
      -H 'Content-Type: application/json' \
      -d '{"galleryNo":"","paging":{"limit":0,"page":0,"total":0}}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ListGalleryImagesCmd {
      galleryNo?: string
      paging?: Paging
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
      data?: ListGalleryImagesResp
    }

    export interface ListGalleryImagesResp {
      images?: ImageInfo[]
      paging?: Paging
    }

    export interface ImageInfo {
      thumbnailToken?: string
      fileTempToken?: string
    }

    export interface Paging {
      limit?: number                 // page limit
      page?: number                  // page number, 1-based
      total?: number                 // total count
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

    let req: ListGalleryImagesCmd | null = null;
    this.http.post<any>(`/vfm/open/api/gallery/images`, req)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: ListGalleryImagesResp = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```

- POST /open/api/gallery/image/transfer
  - Description: Host selected images on gallery
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "images": ([]vfm.CreateGalleryImageCmd) 
      - "galleryNo": (string) 
      - "name": (string) 
      - "fileKey": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/gallery/image/transfer' \
      -H 'Content-Type: application/json' \
      -d '{"images":{"fileKey":"","galleryNo":"","name":""}}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface TransferGalleryImageReq {
      images?: CreateGalleryImageCmd[]
    }

    export interface CreateGalleryImageCmd {
      galleryNo?: string
      name?: string
      fileKey?: string
    }
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

    let req: TransferGalleryImageReq | null = null;
    this.http.post<any>(`/vfm/open/api/gallery/image/transfer`, req)
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

- POST /open/api/versioned-file/list
  - Description: List versioned files
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "paging": (Paging) paging params
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
    - "name": (*string) file name
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (PageRes[github.com/curtisnewbie/vfm/internal/vfm.ApiListVerFileRes]) response data
      - "paging": (Paging) pagination parameters
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
      - "payload": ([]vfm.ApiListVerFileRes) payload values in current page
        - "verFileId": (string) versioned file id
        - "name": (string) file name
        - "fileKey": (string) file key
        - "sizeInBytes": (int64) size in bytes
        - "uploadTime": (int64) last upload time
        - "createTime": (int64) create time of the versioned file record
        - "updateTime": (int64) Update time of the versioned file record
        - "thumbnail": (string) thumbnail token
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/versioned-file/list' \
      -H 'Content-Type: application/json' \
      -d '{"name":"","paging":{"limit":0,"page":0,"total":0}}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ApiListVerFileReq {
      paging?: Paging
      name?: string                  // file name
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
      payload?: ApiListVerFileRes[]
    }

    export interface Paging {
      limit?: number                 // page limit
      page?: number                  // page number, 1-based
      total?: number                 // total count
    }

    export interface ApiListVerFileRes {
      verFileId?: string             // versioned file id
      name?: string                  // file name
      fileKey?: string               // file key
      sizeInBytes?: number           // size in bytes
      uploadTime?: number            // last upload time
      createTime?: number            // create time of the versioned file record
      updateTime?: number            // Update time of the versioned file record
      thumbnail?: string             // thumbnail token
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

    let req: ApiListVerFileReq | null = null;
    this.http.post<any>(`/vfm/open/api/versioned-file/list`, req)
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

- POST /open/api/versioned-file/history
  - Description: List versioned file history
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "paging": (Paging) paging params
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
    - "verFileId": (string) versioned file id
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (PageRes[github.com/curtisnewbie/vfm/internal/vfm.ApiListVerFileHistoryRes]) response data
      - "paging": (Paging) pagination parameters
        - "limit": (int) page limit
        - "page": (int) page number, 1-based
        - "total": (int) total count
      - "payload": ([]vfm.ApiListVerFileHistoryRes) payload values in current page
        - "name": (string) file name
        - "fileKey": (string) file key
        - "sizeInBytes": (int64) size in bytes
        - "uploadTime": (int64) last upload time
        - "thumbnail": (string) thumbnail token
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/versioned-file/history' \
      -H 'Content-Type: application/json' \
      -d '{"paging":{"limit":0,"page":0,"total":0},"verFileId":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ApiListVerFileHistoryReq {
      paging?: Paging
      verFileId?: string             // versioned file id
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
      payload?: ApiListVerFileHistoryRes[]
    }

    export interface Paging {
      limit?: number                 // page limit
      page?: number                  // page number, 1-based
      total?: number                 // total count
    }

    export interface ApiListVerFileHistoryRes {
      name?: string                  // file name
      fileKey?: string               // file key
      sizeInBytes?: number           // size in bytes
      uploadTime?: number            // last upload time
      thumbnail?: string             // thumbnail token
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

    let req: ApiListVerFileHistoryReq | null = null;
    this.http.post<any>(`/vfm/open/api/versioned-file/history`, req)
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

- POST /open/api/versioned-file/accumulated-size
  - Description: Query versioned file log accumulated size
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "verFileId": (string) versioned file id
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (ApiQryVerFileAccuSizeRes) response data
      - "sizeInBytes": (int64) total size in bytes
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/versioned-file/accumulated-size' \
      -H 'Content-Type: application/json' \
      -d '{"verFileId":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ApiQryVerFileAccuSizeReq {
      verFileId?: string             // versioned file id
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string             // error code
      msg?: string                   // message
      error?: boolean                // whether the request was successful
      data?: ApiQryVerFileAccuSizeRes
    }

    export interface ApiQryVerFileAccuSizeRes {
      sizeInBytes?: number           // total size in bytes
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

    let req: ApiQryVerFileAccuSizeReq | null = null;
    this.http.post<any>(`/vfm/open/api/versioned-file/accumulated-size`, req)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: ApiQryVerFileAccuSizeRes = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```

- POST /open/api/versioned-file/create
  - Description: Create versioned file
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "filename": (string) 
    - "fstoreFileId": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (ApiCreateVerFileRes) response data
      - "verFileId": (string) Versioned File Id
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/versioned-file/create' \
      -H 'Content-Type: application/json' \
      -d '{"filename":"","fstoreFileId":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ApiCreateVerFileReq {
      filename?: string
      fstoreFileId?: string
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string             // error code
      msg?: string                   // message
      error?: boolean                // whether the request was successful
      data?: ApiCreateVerFileRes
    }

    export interface ApiCreateVerFileRes {
      verFileId?: string             // Versioned File Id
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

    let req: ApiCreateVerFileReq | null = null;
    this.http.post<any>(`/vfm/open/api/versioned-file/create`, req)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: ApiCreateVerFileRes = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```

- POST /open/api/versioned-file/update
  - Description: Update versioned file
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "verFileId": (string) versioned file id
    - "filename": (string) 
    - "fstoreFileId": (string) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/versioned-file/update' \
      -H 'Content-Type: application/json' \
      -d '{"filename":"","fstoreFileId":"","verFileId":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ApiUpdateVerFileReq {
      verFileId?: string             // versioned file id
      filename?: string
      fstoreFileId?: string
    }
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

    let req: ApiUpdateVerFileReq | null = null;
    this.http.post<any>(`/vfm/open/api/versioned-file/update`, req)
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

- POST /open/api/versioned-file/delete
  - Description: Delete versioned file
  - Bound to Resource: `"manage-files"`
  - JSON Request:
    - "verFileId": (string) Versioned File Id
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/open/api/versioned-file/delete' \
      -H 'Content-Type: application/json' \
      -d '{"verFileId":""}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ApiDelVerFileReq {
      verFileId?: string             // Versioned File Id
    }
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

    let req: ApiDelVerFileReq | null = null;
    this.http.post<any>(`/vfm/open/api/versioned-file/delete`, req)
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

- POST /compensate/thumbnail
  - Description: Compensate thumbnail generation
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/compensate/thumbnail'
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

    this.http.post<any>(`/vfm/compensate/thumbnail`)
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

- POST /compensate/dir/calculate-size
  - Description: Calculate size of all directories recursively
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/compensate/dir/calculate-size'
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

    this.http.post<any>(`/vfm/compensate/dir/calculate-size`)
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

- PUT /bookmark/file/upload
  - Description: Upload bookmark file
  - Bound to Resource: `"manage-bookmarks"`
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X PUT 'http://localhost:8086/bookmark/file/upload'
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

    this.http.put<any>(`/vfm/bookmark/file/upload`)
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

- POST /bookmark/list
  - Description: List bookmarks
  - Bound to Resource: `"manage-bookmarks"`
  - JSON Request:
    - "name": (*string) 
    - "paging": (Paging) 
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/bookmark/list' \
      -H 'Content-Type: application/json' \
      -d '{"name":"","paging":{"limit":0,"page":0,"total":0}}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ListBookmarksReq {
      name?: string
      paging?: Paging
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

    let req: ListBookmarksReq | null = null;
    this.http.post<any>(`/vfm/bookmark/list`, req)
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

- POST /bookmark/remove
  - Description: Remove bookmark
  - Bound to Resource: `"manage-bookmarks"`
  - JSON Request:
    - "id": (int64) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/bookmark/remove' \
      -H 'Content-Type: application/json' \
      -d '{"id":0}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface RemoveBookmarkReq {
      id?: number
    }
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

    let req: RemoveBookmarkReq | null = null;
    this.http.post<any>(`/vfm/bookmark/remove`, req)
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

- POST /bookmark/blacklist/list
  - Description: List bookmark blacklist
  - Bound to Resource: `"manage-bookmarks"`
  - JSON Request:
    - "name": (*string) 
    - "paging": (Paging) 
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/bookmark/blacklist/list' \
      -H 'Content-Type: application/json' \
      -d '{"name":"","paging":{"limit":0,"page":0,"total":0}}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ListBookmarksReq {
      name?: string
      paging?: Paging
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

    let req: ListBookmarksReq | null = null;
    this.http.post<any>(`/vfm/bookmark/blacklist/list`, req)
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

- POST /bookmark/blacklist/remove
  - Description: Remove bookmark blacklist
  - Bound to Resource: `"manage-bookmarks"`
  - JSON Request:
    - "id": (int64) 
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8086/bookmark/blacklist/remove' \
      -H 'Content-Type: application/json' \
      -d '{"id":0}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface RemoveBookmarkReq {
      id?: number
    }
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

    let req: RemoveBookmarkReq | null = null;
    this.http.post<any>(`/vfm/bookmark/blacklist/remove`, req)
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
    curl -X GET 'http://localhost:8086/auth/resource'
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

    this.http.get<any>(`/vfm/auth/resource`)
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
    curl -X GET 'http://localhost:8086/metrics' \
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
    this.http.get<any>(`/vfm/metrics`,
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
    curl -X GET 'http://localhost:8086/debug/pprof'
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    this.http.get<any>(`/vfm/debug/pprof`)
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
    curl -X GET 'http://localhost:8086/debug/pprof/:name'
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    this.http.get<any>(`/vfm/debug/pprof/:name`)
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
    curl -X GET 'http://localhost:8086/debug/pprof/cmdline'
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    this.http.get<any>(`/vfm/debug/pprof/cmdline`)
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
    curl -X GET 'http://localhost:8086/debug/pprof/profile'
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    this.http.get<any>(`/vfm/debug/pprof/profile`)
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
    curl -X GET 'http://localhost:8086/debug/pprof/symbol'
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    this.http.get<any>(`/vfm/debug/pprof/symbol`)
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
    curl -X GET 'http://localhost:8086/debug/pprof/trace'
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    this.http.get<any>(`/vfm/debug/pprof/trace`)
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
    curl -X GET 'http://localhost:8086/doc/api'
    ```

  - Angular HttpClient Demo:
    ```ts
    import { MatSnackBar } from "@angular/material/snack-bar";
    import { HttpClient } from "@angular/common/http";

    constructor(
      private snackBar: MatSnackBar,
      private http: HttpClient
    ) {}

    this.http.get<any>(`/vfm/doc/api`)
      .subscribe({
        next: () => {
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
    ```
