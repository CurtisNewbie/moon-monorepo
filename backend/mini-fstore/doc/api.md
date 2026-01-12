# API Endpoints

## Contents

- [GET /file/stream](#get-filestream)
- [GET /file/raw](#get-fileraw)
- [PUT /file](#put-file)
- [GET /file/info](#get-fileinfo)
- [GET /file/key](#get-filekey)
- [GET /file/direct](#get-filedirect)
- [DELETE /file](#delete-file)
- [POST /file/unzip](#post-fileunzip)
- [POST /backup/file/list](#post-backupfilelist)
- [GET /backup/file/raw](#get-backupfileraw)
- [POST /maintenance/remove-deleted](#post-maintenanceremove-deleted)
- [POST /maintenance/sanitize-storage](#post-maintenancesanitize-storage)
- [POST /maintenance/compute-checksum](#post-maintenancecompute-checksum)
- [GET /storage/info](#get-storageinfo)
- [GET /storage/usage-info](#get-storageusage-info)
- [GET /maintenance/status](#get-maintenancestatus)
- [GET /auth/resource](#get-authresource)
- [GET /debug/trace/recorder/run](#get-debugtracerecorderrun)
- [GET /debug/trace/recorder/snapshot](#get-debugtracerecordersnapshot)
- [GET /debug/trace/recorder/stop](#get-debugtracerecorderstop)

## GET /file/stream

- Description: Media streaming using temporary file key, the file_key's ttl is extended with each subsequent request. This endpoint is expected to be accessible publicly without authorization, since a temporary file_key is generated and used.
- Expected Access Scope: PUBLIC
- Query Parameter:
  - "key": temporary file key
- cURL:
  ```sh
  curl -X GET 'http://localhost:8084/file/stream?key='
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Media streaming using temporary file key, the file_key's ttl is extended with each subsequent request. This endpoint is expected to be accessible publicly without authorization, since a temporary file_key is generated and used.
  func ApiTempKeyStreamFile(rail miso.Rail, key string) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/file/stream", "fstore").
  		AddQuery("key", key).
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

  tempKeyStreamFile() {
    let key: any | null = null;
    this.http.get<any>(`/fstore/file/stream?key=${key}`)
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

## GET /file/raw

- Description: Download file using temporary file key. This endpoint is expected to be accessible publicly without authorization, since a temporary file_key is generated and used.
- Expected Access Scope: PUBLIC
- Query Parameter:
  - "key": temporary file key
- cURL:
  ```sh
  curl -X GET 'http://localhost:8084/file/raw?key='
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Download file using temporary file key. This endpoint is expected to be accessible publicly without authorization, since a temporary file_key is generated and used.
  func ApiTempKeyDownloadFile(rail miso.Rail, key string) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/file/raw", "fstore").
  		AddQuery("key", key).
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

  tempKeyDownloadFile() {
    let key: any | null = null;
    this.http.get<any>(`/fstore/file/raw?key=${key}`)
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

## PUT /file

- Description: Upload file. A temporary file_id is returned, which should be used to exchange the real file_id
- Bound to Resource: `"fstore-upload"`
- Header Parameter:
  - "filename": name of the uploaded file
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (string) response data
- cURL:
  ```sh
  curl -X PUT 'http://localhost:8084/file' \
    -H 'filename: '
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Upload file. A temporary file_id is returned, which should be used to exchange the real file_id
  func ApiUploadFile(rail miso.Rail, filename string) (string, error) {
  	var res miso.GnResp[string]
  	err := miso.NewDynClient(rail, "/file", "fstore").
  		AddHeader("filename", filename).
  		Put(nil).
  		Json(&res)
  	if err != nil {
  		return "", err
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

  uploadFile() {
    let filename: any | null = null;
    this.http.put<any>(`/fstore/file`, null,
      {
        headers: {
          "filename": filename
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

## GET /file/info

- Description: Fetch file info
- Query Parameter:
  - "fileId": actual file_id of the file record
  - "uploadFileId": temporary file_id returned when uploading files
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (FstoreFile) response data
      - "fileId": (string) file unique identifier
      - "name": (string) file name
      - "status": (string) status, 'NORMAL', 'LOG_DEL' (logically deleted), 'PHY_DEL' (physically deleted)
      - "size": (int64) file size in bytes
      - "md5": (string) MD5 checksum
      - "uplTime": (int64) upload time
      - "logDelTime": (int64) logically deleted at
      - "phyDelTime": (int64) physically deleted at
- cURL:
  ```sh
  curl -X GET 'http://localhost:8084/file/info?fileId=&uploadFileId='
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type FstoreFile struct {
  	FileId string `json:"fileId"`  // file unique identifier
  	Name string `json:"name"`      // file name
  	Status string `json:"status"`  // status, 'NORMAL', 'LOG_DEL' (logically deleted), 'PHY_DEL' (physically deleted)
  	Size int64 `json:"size"`       // file size in bytes
  	Md5 string `json:"md5"`        // MD5 checksum
  	UplTime atom.Time `json:"uplTime"` // upload time
  	LogDelTime *atom.Time `json:"logDelTime"` // logically deleted at
  	PhyDelTime *atom.Time `json:"phyDelTime"` // physically deleted at
  }

  // Fetch file info
  func ApiGetFileInfo(rail miso.Rail, fileId string, uploadFileId string) (FstoreFile, error) {
  	var res miso.GnResp[FstoreFile]
  	err := miso.NewDynClient(rail, "/file/info", "fstore").
  		AddQuery("fileId", fileId).
  		AddQuery("uploadFileId", uploadFileId).
  		Get().
  		Json(&res)
  	if err != nil {
  		var dat FstoreFile
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
    data?: FstoreFile;
  }

  export interface FstoreFile {
    fileId?: string;               // file unique identifier
    name?: string;                 // file name
    status?: string;               // status, 'NORMAL', 'LOG_DEL' (logically deleted), 'PHY_DEL' (physically deleted)
    size?: number;                 // file size in bytes
    md5?: string;                  // MD5 checksum
    uplTime?: number;              // upload time
    logDelTime?: number;           // logically deleted at
    phyDelTime?: number;           // physically deleted at
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

  getFileInfo() {
    let fileId: any | null = null;
    let uploadFileId: any | null = null;
    this.http.get<any>(`/fstore/file/info?fileId=${fileId}&uploadFileId=${uploadFileId}`)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: FstoreFile = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
  }
  ```

## GET /file/key

- Description: Generate temporary file key for downloading and streaming. This endpoint is expected to be called internally by another backend service that validates the ownership of the file properly.
- Query Parameter:
  - "fileId": actual file_id of the file record
  - "filename": the name that will be used when downloading the file
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (string) response data
- cURL:
  ```sh
  curl -X GET 'http://localhost:8084/file/key?fileId=&filename='
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Generate temporary file key for downloading and streaming. This endpoint is expected to be called internally by another backend service that validates the ownership of the file properly.
  func ApiGenFileKey(rail miso.Rail, fileId string, filename string) (string, error) {
  	var res miso.GnResp[string]
  	err := miso.NewDynClient(rail, "/file/key", "fstore").
  		AddQuery("fileId", fileId).
  		AddQuery("filename", filename).
  		Get().
  		Json(&res)
  	if err != nil {
  		return "", err
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

  genFileKey() {
    let fileId: any | null = null;
    let filename: any | null = null;
    this.http.get<any>(`/fstore/file/key?fileId=${fileId}&filename=${filename}`)
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

## GET /file/direct

- Description: Download files directly using file_id. This endpoint is expected to be protected and only used internally by another backend service. Users can eaily steal others file_id and attempt to download the file, so it's better not be exposed to the end users.
- Query Parameter:
  - "fileId": actual file_id of the file record
- cURL:
  ```sh
  curl -X GET 'http://localhost:8084/file/direct?fileId='
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Download files directly using file_id. This endpoint is expected to be protected and only used internally by another backend service. Users can eaily steal others file_id and attempt to download the file, so it's better not be exposed to the end users.
  func ApiDirectDownloadFile(rail miso.Rail, fileId string) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/file/direct", "fstore").
  		AddQuery("fileId", fileId).
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

  directDownloadFile() {
    let fileId: any | null = null;
    this.http.get<any>(`/fstore/file/direct?fileId=${fileId}`)
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

## DELETE /file

- Description: Mark file as deleted.
- Query Parameter:
  - "fileId": actual file_id of the file record
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
- cURL:
  ```sh
  curl -X DELETE 'http://localhost:8084/file?fileId='
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Mark file as deleted.
  func ApiDeleteFile(rail miso.Rail, fileId string) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/file", "fstore").
  		AddQuery("fileId", fileId).
  		Delete().
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
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

  deleteFile() {
    let fileId: any | null = null;
    this.http.delete<any>(`/fstore/file?fileId=${fileId}`)
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

## POST /file/unzip

- Description: Unzip archive, upload all the zip entries, and reply the final results back to the caller asynchronously
- JSON Request:
    - "fileId": (string) file_id of zip file. Required.
    - "replyToEventBus": (string) name of the rabbitmq exchange to reply to, routing_key is '#'. Required.
    - "extra": (string) extra information that will be passed around for the caller
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
- cURL:
  ```sh
  curl -X POST 'http://localhost:8084/file/unzip' \
    -H 'Content-Type: application/json' \
    -d '{"extra":"","fileId":"","replyToEventBus":""}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type UnzipFileReq struct {
  	FileId string `json:"fileId"`  // file_id of zip file. Required.
  	ReplyToEventBus string `json:"replyToEventBus"` // name of the rabbitmq exchange to reply to, routing_key is '#'. Required.
  	Extra string `json:"extra"`    // extra information that will be passed around for the caller
  }

  // Unzip archive, upload all the zip entries, and reply the final results back to the caller asynchronously
  func ApiUnzipFile(rail miso.Rail, req UnzipFileReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/file/unzip", "fstore").
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
  export interface UnzipFileReq {
    fileId?: string;               // file_id of zip file. Required.
    replyToEventBus?: string;      // name of the rabbitmq exchange to reply to, routing_key is '#'. Required.
    extra?: string;                // extra information that will be passed around for the caller
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

  unzipFile() {
    let req: UnzipFileReq | null = null;
    this.http.post<any>(`/fstore/file/unzip`, req)
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

## POST /backup/file/list

- Description: Backup tool list files
- Expected Access Scope: PUBLIC
- Header Parameter:
  - "Authorization": Basic Authorization
- JSON Request:
    - "limit": (int64) 
    - "idOffset": (int) 
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (ListBackupFileResp) response data
      - "files": ([]fstore.BackupFileInf) 
        - "id": (int64) 
        - "fileId": (string) 
        - "name": (string) 
        - "status": (string) 
        - "size": (int64) 
        - "md5": (string) 
- cURL:
  ```sh
  curl -X POST 'http://localhost:8084/backup/file/list' \
    -H 'Authorization: ' \
    -H 'Content-Type: application/json' \
    -d '{"idOffset":0,"limit":0}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListBackupFileReq struct {
  	Limit int64 `json:"limit"`
  	IdOffset int `json:"idOffset"`
  }

  type ListBackupFileResp struct {
  	Files []BackupFileInf `json:"files"`
  }

  type BackupFileInf struct {
  	Id int64 `json:"id"`
  	FileId string `json:"fileId"`
  	Name string `json:"name"`
  	Status string `json:"status"`
  	Size int64 `json:"size"`
  	Md5 string `json:"md5"`
  }

  // Backup tool list files
  func ApiBackupListFiles(rail miso.Rail, req ListBackupFileReq, authorization string) (ListBackupFileResp, error) {
  	var res miso.GnResp[ListBackupFileResp]
  	err := miso.NewDynClient(rail, "/backup/file/list", "fstore").
  		AddHeader("authorization", authorization).
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		var dat ListBackupFileResp
  		return dat, err
  	}
  	return res.Data, nil
  }
  ```

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface ListBackupFileReq {
    limit?: number;
    idOffset?: number;
  }

  export interface Resp {
    errorCode?: string;            // error code
    msg?: string;                  // message
    error?: boolean;               // whether the request was successful
    data?: ListBackupFileResp;
  }

  export interface ListBackupFileResp {
    files?: BackupFileInf[];
  }

  export interface BackupFileInf {
    id?: number;
    fileId?: string;
    name?: string;
    status?: string;
    size?: number;
    md5?: string;
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

  backupListFiles() {
    let authorization: any | null = null;
    let req: ListBackupFileReq | null = null;
    this.http.post<any>(`/fstore/backup/file/list`, req,
      {
        headers: {
          "Authorization": authorization
        }
      })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: ListBackupFileResp = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
  }
  ```

## GET /backup/file/raw

- Description: Backup tool download file
- Expected Access Scope: PUBLIC
- Header Parameter:
  - "Authorization": Basic Authorization
- Query Parameter:
  - "fileId": actual file_id of the file record
- cURL:
  ```sh
  curl -X GET 'http://localhost:8084/backup/file/raw?fileId=' \
    -H 'Authorization: '
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Backup tool download file
  func ApiBackupDownFile(rail miso.Rail, fileId string, authorization string) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/backup/file/raw", "fstore").
  		AddQuery("fileId", fileId).
  		AddHeader("authorization", authorization).
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

  backupDownFile() {
    let fileId: any | null = null;
    let authorization: any | null = null;
    this.http.get<any>(`/fstore/backup/file/raw?fileId=${fileId}`,
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

## POST /maintenance/remove-deleted

- Description: Remove files that are logically deleted and not linked (symbolically)
- Bound to Resource: `"fstore:server:maintenance"`
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
- cURL:
  ```sh
  curl -X POST 'http://localhost:8084/maintenance/remove-deleted'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Remove files that are logically deleted and not linked (symbolically)
  func ApiRemoveDeletedFiles(rail miso.Rail) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/maintenance/remove-deleted", "fstore").
  		Post(nil).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
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

  removeDeletedFiles() {
    this.http.post<any>(`/fstore/maintenance/remove-deleted`, null)
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

## POST /maintenance/sanitize-storage

- Description: Sanitize storage, remove files in storage directory that don't exist in database
- Bound to Resource: `"fstore:server:maintenance"`
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
- cURL:
  ```sh
  curl -X POST 'http://localhost:8084/maintenance/sanitize-storage'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Sanitize storage, remove files in storage directory that don't exist in database
  func ApiSanitizeStorage(rail miso.Rail) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/maintenance/sanitize-storage", "fstore").
  		Post(nil).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
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

  sanitizeStorage() {
    this.http.post<any>(`/fstore/maintenance/sanitize-storage`, null)
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

## POST /maintenance/compute-checksum

- Description: Compute files' checksum if absent
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
- cURL:
  ```sh
  curl -X POST 'http://localhost:8084/maintenance/compute-checksum'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Compute files' checksum if absent
  func ApiComputeChecksum(rail miso.Rail) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/maintenance/compute-checksum", "fstore").
  		Post(nil).
  		Json(&res)
  	if err != nil {
  		return err
  	}
  	return nil
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

  computeChecksum() {
    this.http.post<any>(`/fstore/maintenance/compute-checksum`, null)
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

## GET /storage/info

- Description: Fetch storage info
- Bound to Resource: `"fstore:server:maintenance"`
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (StorageInfo) response data
      - "volumns": ([]fstore.VolumnInfo) 
        - "mounted": (string) 
        - "total": (uint64) 
        - "used": (uint64) 
        - "available": (uint64) 
        - "usedPercent": (float64) 
        - "totalText": (string) 
        - "usedText": (string) 
        - "availableText": (string) 
        - "usedPercentText": (string) 
- cURL:
  ```sh
  curl -X GET 'http://localhost:8084/storage/info'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type StorageInfo struct {
  	Volumns []VolumnInfo `json:"volumns"`
  }

  type VolumnInfo struct {
  	Mounted string `json:"mounted"`
  	Total uint64 `json:"total"`
  	Used uint64 `json:"used"`
  	Available uint64 `json:"available"`
  	UsedPercent float64 `json:"usedPercent"`
  	TotalText string `json:"totalText"`
  	UsedText string `json:"usedText"`
  	AvailableText string `json:"availableText"`
  	UsedPercentText string `json:"usedPercentText"`
  }

  // Fetch storage info
  func ApiFetchStorageInfo(rail miso.Rail) (StorageInfo, error) {
  	var res miso.GnResp[StorageInfo]
  	err := miso.NewDynClient(rail, "/storage/info", "fstore").
  		Get().
  		Json(&res)
  	if err != nil {
  		var dat StorageInfo
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
    data?: StorageInfo;
  }

  export interface StorageInfo {
    volumns?: VolumnInfo[];
  }

  export interface VolumnInfo {
    mounted?: string;
    total?: number;
    used?: number;
    available?: number;
    usedPercent?: number;
    totalText?: string;
    usedText?: string;
    availableText?: string;
    usedPercentText?: string;
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

  fetchStorageInfo() {
    this.http.get<any>(`/fstore/storage/info`)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: StorageInfo = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
  }
  ```

## GET /storage/usage-info

- Description: Fetch storage usage info
- Bound to Resource: `"fstore:server:maintenance"`
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": ([]fstore.StorageUsageInfo) response data
      - "type": (string) 
      - "path": (string) 
      - "used": (uint64) 
      - "usedText": (string) 
- cURL:
  ```sh
  curl -X GET 'http://localhost:8084/storage/usage-info'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type StorageUsageInfo struct {
  	Type string `json:"type"`
  	Path string `json:"path"`
  	Used uint64 `json:"used"`
  	UsedText string `json:"usedText"`
  }

  // Fetch storage usage info
  func ApiFetchStorageUsageInfo(rail miso.Rail) ([]StorageUsageInfo, error) {
  	var res miso.GnResp[[]StorageUsageInfo]
  	err := miso.NewDynClient(rail, "/storage/usage-info", "fstore").
  		Get().
  		Json(&res)
  	if err != nil {
  		var dat []StorageUsageInfo
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
    data?: StorageUsageInfo[];
  }

  export interface StorageUsageInfo {
    type?: string;
    path?: string;
    used?: number;
    usedText?: string;
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

  fetchStorageUsageInfo() {
    this.http.get<any>(`/fstore/storage/usage-info`)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: StorageUsageInfo[] = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
  }
  ```

## GET /maintenance/status

- Description: Check server maintenance status
- Bound to Resource: `"fstore:server:maintenance"`
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (MaintenanceStatus) response data
      - "underMaintenance": (bool) 
- cURL:
  ```sh
  curl -X GET 'http://localhost:8084/maintenance/status'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type MaintenanceStatus struct {
  	UnderMaintenance bool `json:"underMaintenance"`
  }

  // Check server maintenance status
  func ApiFetchMaintenanceStatus(rail miso.Rail) (MaintenanceStatus, error) {
  	var res miso.GnResp[MaintenanceStatus]
  	err := miso.NewDynClient(rail, "/maintenance/status", "fstore").
  		Get().
  		Json(&res)
  	if err != nil {
  		var dat MaintenanceStatus
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
    data?: MaintenanceStatus;
  }

  export interface MaintenanceStatus {
    underMaintenance?: boolean;
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

  fetchMaintenanceStatus() {
    this.http.get<any>(`/fstore/maintenance/status`)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: MaintenanceStatus = resp.data;
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
  curl -X GET 'http://localhost:8084/auth/resource'
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
  	err := miso.NewDynClient(rail, "/auth/resource", "fstore").
  		Get().
  		Json(&res)
  	if err != nil {
  		var dat ResourceInfoRes
  		return dat, err
  	}
  	return res.Data, nil
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
    this.http.get<ResourceInfoRes>(`/fstore/auth/resource`)
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
  curl -X GET 'http://localhost:8084/debug/trace/recorder/run?duration='
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Start FlightRecorder. Recorded result is written to trace.out when it's finished or stopped.
  func SendRequest(rail miso.Rail, duration string) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/debug/trace/recorder/run", "fstore").
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
    this.http.get<any>(`/fstore/debug/trace/recorder/run?duration=${duration}`)
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
  curl -X GET 'http://localhost:8084/debug/trace/recorder/snapshot'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // FlightRecorder take snapshot. Recorded result is written to trace.out.
  func SendRequest(rail miso.Rail) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/debug/trace/recorder/snapshot", "fstore").
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
    this.http.get<any>(`/fstore/debug/trace/recorder/snapshot`)
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
  curl -X GET 'http://localhost:8084/debug/trace/recorder/stop'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  // Stop existing FlightRecorder session.
  func SendRequest(rail miso.Rail) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynClient(rail, "/debug/trace/recorder/stop", "fstore").
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
    this.http.get<any>(`/fstore/debug/trace/recorder/stop`)
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

- GenImgThumbnailPipeline
  - Description: Pipeline to trigger async image thumbnail generation, will reply api.ImageCompressReplyEvent when the processing succeeds.
  - RabbitMQ Queue: `event.bus.fstore.image.compress.processing`
  - RabbitMQ Exchange: `event.bus.fstore.image.compress.processing`
  - RabbitMQ RoutingKey: `#`
  - Event Payload:
    - "identifier": (string) identifier
    - "fileId": (string) file id from mini-fstore
    - "replyTo": (string) event bus that will receive event about the generated image thumbnail.

- GenVidThumbnailPipeline
  - Description: Pipeline to trigger async video thumbnail generation, will reply api.GenVideoThumbnailReplyEvent when the processing succeeds.
  - RabbitMQ Queue: `event.bus.fstore.video.thumbnail.processing`
  - RabbitMQ Exchange: `event.bus.fstore.video.thumbnail.processing`
  - RabbitMQ RoutingKey: `#`
  - Event Payload:
    - "identifier": (string) dentifier
    - "fileId": (string) file id from mini-fstore
    - "replyTo": (string) event bus that will receive event about the generated video thumbnail.
