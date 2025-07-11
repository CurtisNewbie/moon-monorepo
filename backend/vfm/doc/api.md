# API Endpoints

## Contents

- [GET /open/api/file/upload/duplication/preflight](#get-openapifileuploadduplicationpreflight)
- [GET /open/api/file/parent](#get-openapifileparent)
- [POST /open/api/file/move-to-dir](#post-openapifilemove-to-dir)
- [POST /open/api/file/batch-move-to-dir](#post-openapifilebatch-move-to-dir)
- [POST /open/api/file/make-dir](#post-openapifilemake-dir)
- [GET /open/api/file/dir/list](#get-openapifiledirlist)
- [POST /open/api/file/list](#post-openapifilelist)
- [POST /open/api/file/delete](#post-openapifiledelete)
- [POST /open/api/file/dir/truncate](#post-openapifiledirtruncate)
- [POST /open/api/file/dir/bottom-up-tree](#post-openapifiledirbottom-up-tree)
- [GET /open/api/file/dir/top-down-tree](#get-openapifiledirtop-down-tree)
- [POST /open/api/file/delete/batch](#post-openapifiledeletebatch)
- [POST /open/api/file/create](#post-openapifilecreate)
- [POST /open/api/file/info/update](#post-openapifileinfoupdate)
- [POST /open/api/file/token/generate](#post-openapifiletokengenerate)
- [POST /open/api/file/unpack](#post-openapifileunpack)
- [GET /open/api/file/token/qrcode](#get-openapifiletokenqrcode)
- [GET /open/api/vfolder/brief/owned](#get-openapivfolderbriefowned)
- [POST /open/api/vfolder/list](#post-openapivfolderlist)
- [POST /open/api/vfolder/create](#post-openapivfoldercreate)
- [POST /open/api/vfolder/file/add](#post-openapivfolderfileadd)
- [POST /open/api/vfolder/file/remove](#post-openapivfolderfileremove)
- [POST /open/api/vfolder/share](#post-openapivfoldershare)
- [POST /open/api/vfolder/access/remove](#post-openapivfolderaccessremove)
- [POST /open/api/vfolder/granted/list](#post-openapivfoldergrantedlist)
- [POST /open/api/vfolder/remove](#post-openapivfolderremove)
- [GET /open/api/gallery/brief/owned](#get-openapigallerybriefowned)
- [POST /open/api/gallery/new](#post-openapigallerynew)
- [POST /open/api/gallery/update](#post-openapigalleryupdate)
- [POST /open/api/gallery/delete](#post-openapigallerydelete)
- [POST /open/api/gallery/list](#post-openapigallerylist)
- [POST /open/api/gallery/access/grant](#post-openapigalleryaccessgrant)
- [POST /open/api/gallery/access/remove](#post-openapigalleryaccessremove)
- [POST /open/api/gallery/access/list](#post-openapigalleryaccesslist)
- [POST /open/api/gallery/images](#post-openapigalleryimages)
- [POST /open/api/gallery/image/transfer](#post-openapigalleryimagetransfer)
- [POST /open/api/versioned-file/list](#post-openapiversioned-filelist)
- [POST /open/api/versioned-file/history](#post-openapiversioned-filehistory)
- [POST /open/api/versioned-file/accumulated-size](#post-openapiversioned-fileaccumulated-size)
- [POST /open/api/versioned-file/create](#post-openapiversioned-filecreate)
- [POST /open/api/versioned-file/update](#post-openapiversioned-fileupdate)
- [POST /open/api/versioned-file/delete](#post-openapiversioned-filedelete)
- [POST /compensate/thumbnail](#post-compensatethumbnail)
- [POST /compensate/regenerate-video-thumbnails](#post-compensateregenerate-video-thumbnails)
- [PUT /bookmark/file/upload](#put-bookmarkfileupload)
- [POST /bookmark/list](#post-bookmarklist)
- [POST /bookmark/remove](#post-bookmarkremove)
- [POST /bookmark/blacklist/list](#post-bookmarkblacklistlist)
- [POST /bookmark/blacklist/remove](#post-bookmarkblacklistremove)
- [GET /history/list-browse-history](#get-historylist-browse-history)
- [POST /history/record-browse-history](#post-historyrecord-browse-history)
- [GET /maintenance/status](#get-maintenancestatus)
- [POST /internal/v1/file/create](#post-internalv1filecreate)
- [GET /internal/file/upload/duplication/preflight](#get-internalfileuploadduplicationpreflight)
- [POST /internal/file/check-access](#post-internalfilecheck-access)
- [POST /internal/file/fetch-info](#post-internalfilefetch-info)
- [POST /internal/v1/file/make-dir](#post-internalv1filemake-dir)
- [GET /auth/resource](#get-authresource)

## GET /open/api/file/upload/duplication/preflight

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  func ApiPreflightCheckDuplicate(rail miso.Rail, fileName string, parentFileKey string) (bool, error) {
  	var res miso.GnResp[bool]
  	err := miso.NewDynTClient(rail, "/open/api/file/upload/duplication/preflight", "vfm").
  		AddQueryParams("fileName", fileName).
  		AddQueryParams("parentFileKey", parentFileKey).
  		Get().
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		return false, err
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
    data?: boolean;                // response data
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

  preflightCheckDuplicate() {
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
  }
  ```

## GET /open/api/file/parent

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ParentFileInfo struct {
  	FileKey string `json:"fileKey"`
  	Filename string `json:"fileName"`
  }

  func ApiGetParentFile(rail miso.Rail, fileKey string) (*ParentFileInfo, error) {
  	var res miso.GnResp[*ParentFileInfo]
  	err := miso.NewDynTClient(rail, "/open/api/file/parent", "vfm").
  		AddQueryParams("fileKey", fileKey).
  		Get().
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		return nil, err
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
    data?: ParentFileInfo;
  }

  export interface ParentFileInfo {
    fileKey?: string;
    fileName?: string;
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

  getParentFile() {
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
  }
  ```

## POST /open/api/file/move-to-dir

- Description: User move file into directory
- Bound to Resource: `"manage-files"`
- JSON Request:
    - "uuid": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type MoveIntoDirReq struct {
  	Uuid string `json:"uuid"`      // Required.
  	ParentFileUuid string `json:"parentFileUuid"`
  }

  func ApiMoveFileToDir(rail miso.Rail, req MoveIntoDirReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/file/move-to-dir", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface MoveIntoDirReq {
    uuid?: string;                 // Required.
    parentFileUuid?: string;
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

  moveFileToDir() {
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
  }
  ```

## POST /open/api/file/batch-move-to-dir

- Description: User move files into directory
- Bound to Resource: `"manage-files"`
- JSON Request:
    - "instructions": ([]vfm.MoveIntoDirReq) 
      - "uuid": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type BatchMoveIntoDirReq struct {
  	Instructions []MoveIntoDirReq `json:"instructions"`
  }

  type MoveIntoDirReq struct {
  	Uuid string `json:"uuid"`      // Required.
  	ParentFileUuid string `json:"parentFileUuid"`
  }

  func ApiBatchMoveFileToDir(rail miso.Rail, req BatchMoveIntoDirReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/file/batch-move-to-dir", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface BatchMoveIntoDirReq {
    instructions?: MoveIntoDirReq[];
  }

  export interface MoveIntoDirReq {
    uuid?: string;                 // Required.
    parentFileUuid?: string;
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

  batchMoveFileToDir() {
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
  }
  ```

## POST /open/api/file/make-dir

- Description: User make directory
- Bound to Resource: `"manage-files"`
- JSON Request:
    - "parentFile": (string) 
    - "name": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type MakeDirReq struct {
  	ParentFile string `json:"parentFile"`
  	Name string `json:"name"`      // Required.
  }

  func ApiMakeDir(rail miso.Rail, req MakeDirReq) (string, error) {
  	var res miso.GnResp[string]
  	err := miso.NewDynTClient(rail, "/open/api/file/make-dir", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface MakeDirReq {
    parentFile?: string;
    name?: string;                 // Required.
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

  makeDir() {
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
  }
  ```

## GET /open/api/file/dir/list

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListedDir struct {
  	Id int `json:"id"`
  	Uuid string `json:"uuid"`
  	Name string `json:"name"`
  }

  func ApiListDir(rail miso.Rail) ([]ListedDir, error) {
  	var res miso.GnResp[[]ListedDir]
  	err := miso.NewDynTClient(rail, "/open/api/file/dir/list", "vfm").
  		Get().
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat []ListedDir
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
    data?: ListedDir[];
  }

  export interface ListedDir {
    id?: number;
    uuid?: string;
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

  listDir() {
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
  }
  ```

## POST /open/api/file/list

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
    - "fileKey": (*string) 
    - "orderByName": (bool) 
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
    -d '{"fileKey":"","fileType":"","filename":"","folderNo":"","orderByName":false,"paging":{"limit":0,"page":0,"total":0},"parentFile":"","sensitive":false}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListFileReq struct {
  	Page miso.Paging `json:"paging"`
  	Filename *string `json:"filename"`
  	FolderNo *string `json:"folderNo"`
  	FileType *string `json:"fileType"`
  	ParentFile *string `json:"parentFile"`
  	Sensitive *bool `json:"sensitive"`
  	FileKey *string `json:"fileKey"`
  	OrderByName bool `json:"orderByName"`
  }


  type ListedFile struct {
  	Id int `json:"id"`
  	Uuid string `json:"uuid"`
  	Name string `json:"name"`
  	UploadTime util.ETime `json:"uploadTime"`
  	UploaderName string `json:"uploaderName"`
  	SizeInBytes int64 `json:"sizeInBytes"`
  	FileType string `json:"fileType"`
  	UpdateTime util.ETime `json:"updateTime"`
  	ParentFileName string `json:"parentFileName"`
  	SensitiveMode string `json:"sensitiveMode"`
  	ThumbnailToken string `json:"thumbnailToken"`
  }

  func ApiListFiles(rail miso.Rail, req ListFileReq) (miso.PageRes[ListedFile], error) {
  	var res miso.GnResp[miso.PageRes[ListedFile]]
  	err := miso.NewDynTClient(rail, "/open/api/file/list", "vfm").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat miso.PageRes[ListedFile]
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
  export interface ListFileReq {
    paging?: Paging;
    filename?: string;
    folderNo?: string;
    fileType?: string;
    parentFile?: string;
    sensitive?: boolean;
    fileKey?: string;
    orderByName?: boolean;
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
    payload?: ListedFile[];
  }

  export interface Paging {
    limit?: number;                // page limit
    page?: number;                 // page number, 1-based
    total?: number;                // total count
  }

  export interface ListedFile {
    id?: number;
    uuid?: string;
    name?: string;
    uploadTime?: number;
    uploaderName?: string;
    sizeInBytes?: number;
    fileType?: string;
    updateTime?: number;
    parentFileName?: string;
    sensitiveMode?: string;
    thumbnailToken?: string;
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

  listFiles() {
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
  }
  ```

## POST /open/api/file/delete

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type DeleteFileReq struct {
  	Uuid string `json:"uuid"`
  }

  func ApiDeleteFiles(rail miso.Rail, req DeleteFileReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/file/delete", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface DeleteFileReq {
    uuid?: string;
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

  deleteFiles() {
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
  }
  ```

## POST /open/api/file/dir/truncate

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type DeleteFileReq struct {
  	Uuid string `json:"uuid"`
  }

  func ApiTruncateDir(rail miso.Rail, req DeleteFileReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/file/dir/truncate", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface DeleteFileReq {
    uuid?: string;
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

  truncateDir() {
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
  }
  ```

## POST /open/api/file/dir/bottom-up-tree

- Description: Fetch directory tree bottom up.
- Bound to Resource: `"manage-files"`
- JSON Request:
    - "fileKey": (string) 
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (*vfm.DirBottomUpTreeNode) response data
      - "fileKey": (string) 
      - "name": (string) 
      - "child": (*vfm.DirBottomUpTreeNode) 
- cURL:
  ```sh
  curl -X POST 'http://localhost:8086/open/api/file/dir/bottom-up-tree' \
    -H 'Content-Type: application/json' \
    -d '{"fileKey":""}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type FetchDirTreeReq struct {
  	FileKey string `json:"fileKey"`
  }

  type DirBottomUpTreeNode struct {
  	FileKey string `json:"fileKey"`
  	Name string `json:"name"`
  	Child *DirBottomUpTreeNode `json:"child"`
  }

  func ApiFetchDirBottomUpTree(rail miso.Rail, req FetchDirTreeReq) (*DirBottomUpTreeNode, error) {
  	var res miso.GnResp[*DirBottomUpTreeNode]
  	err := miso.NewDynTClient(rail, "/open/api/file/dir/bottom-up-tree", "vfm").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		return nil, err
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
  export interface FetchDirTreeReq {
    fileKey?: string;
  }

  export interface Resp {
    errorCode?: string;            // error code
    msg?: string;                  // message
    error?: boolean;               // whether the request was successful
    data?: DirBottomUpTreeNode;
  }

  export interface DirBottomUpTreeNode {
    fileKey?: string;
    name?: string;
    child?: DirBottomUpTreeNode;
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

  fetchDirBottomUpTree() {
    let req: FetchDirTreeReq | null = null;
    this.http.post<any>(`/vfm/open/api/file/dir/bottom-up-tree`, req)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: DirBottomUpTreeNode = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
  }
  ```

## GET /open/api/file/dir/top-down-tree

- Description: Fetch directory tree top down.
- Bound to Resource: `"manage-files"`
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (*vfm.DirTopDownTreeNode) response data
      - "fileKey": (string) 
      - "name": (string) 
      - "child": ([]*vfm.DirTopDownTreeNode) 
- cURL:
  ```sh
  curl -X GET 'http://localhost:8086/open/api/file/dir/top-down-tree'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type DirTopDownTreeNode struct {
  	FileKey string `json:"fileKey"`
  	Name string `json:"name"`
  	Child []*DirTopDownTreeNode `json:"child"`
  }

  func ApiFetchDirTopDownTree(rail miso.Rail) (*DirTopDownTreeNode, error) {
  	var res miso.GnResp[*DirTopDownTreeNode]
  	err := miso.NewDynTClient(rail, "/open/api/file/dir/top-down-tree", "vfm").
  		Get().
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		return nil, err
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
    data?: DirTopDownTreeNode;
  }

  export interface DirTopDownTreeNode {
    fileKey?: string;
    name?: string;
    child?: DirTopDownTreeNode[];
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

  fetchDirTopDownTree() {
    this.http.get<any>(`/vfm/open/api/file/dir/top-down-tree`)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: DirTopDownTreeNode = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
  }
  ```

## POST /open/api/file/delete/batch

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type BatchDeleteFileReq struct {
  	FileKeys []string `json:"fileKeys"`
  }

  func ApiBatchDeleteFile(rail miso.Rail, req BatchDeleteFileReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/file/delete/batch", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface BatchDeleteFileReq {
    fileKeys?: string[];
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

  batchDeleteFile() {
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
  }
  ```

## POST /open/api/file/create

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type CreateFileReq struct {
  	Filename string `json:"filename"`
  	FakeFstoreFileId string `json:"fstoreFileId"`
  	ParentFile string `json:"parentFile"`
  }

  func ApiCreateFile(rail miso.Rail, req CreateFileReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/file/create", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface CreateFileReq {
    filename?: string;
    fstoreFileId?: string;
    parentFile?: string;
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

  createFile() {
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
  }
  ```

## POST /open/api/file/info/update

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type UpdateFileReq struct {
  	Id int `json:"id"`
  	Name string `json:"name"`
  	SensitiveMode string `json:"sensitiveMode"`
  }

  func ApiUpdateFile(rail miso.Rail, req UpdateFileReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/file/info/update", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface UpdateFileReq {
    id?: number;
    name?: string;
    sensitiveMode?: string;
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

  updateFile() {
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
  }
  ```

## POST /open/api/file/token/generate

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type GenerateTempTokenReq struct {
  	FileKey string `json:"fileKey"`
  }

  func ApiGenFileTkn(rail miso.Rail, req GenerateTempTokenReq) (string, error) {
  	var res miso.GnResp[string]
  	err := miso.NewDynTClient(rail, "/open/api/file/token/generate", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface GenerateTempTokenReq {
    fileKey?: string;
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

  genFileTkn() {
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
  }
  ```

## POST /open/api/file/unpack

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type UnpackZipReq struct {
  	FileKey string `json:"fileKey"`
  	ParentFileKey string `json:"parentFileKey"`
  }

  func ApiUnpackZip(rail miso.Rail, req UnpackZipReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/file/unpack", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface UnpackZipReq {
    fileKey?: string;
    parentFileKey?: string;
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

  unpackZip() {
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
  }
  ```

## GET /open/api/file/token/qrcode

- Description: User generate qrcode image for temporary token
- Expected Access Scope: PUBLIC
  - Query Parameter:
  - "token": Generated temporary file key
- cURL:
  ```sh
  curl -X GET 'http://localhost:8086/open/api/file/token/qrcode?token='
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  func ApiGenFileTknQRCode(rail miso.Rail, token string) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/file/token/qrcode", "vfm").
  		AddQueryParams("token", token).
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

  genFileTknQRCode() {
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
  }
  ```

## GET /open/api/vfolder/brief/owned

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type VFolderBrief struct {
  	FolderNo string `json:"folderNo"`
  	Name string `json:"name"`
  }

  func ApiListVFolderBrief(rail miso.Rail) ([]VFolderBrief, error) {
  	var res miso.GnResp[[]VFolderBrief]
  	err := miso.NewDynTClient(rail, "/open/api/vfolder/brief/owned", "vfm").
  		Get().
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat []VFolderBrief
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
    data?: VFolderBrief[];
  }

  export interface VFolderBrief {
    folderNo?: string;
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

  listVFolderBrief() {
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
  }
  ```

## POST /open/api/vfolder/list

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListVFolderReq struct {
  	Page miso.Paging `json:"paging"`
  	Name string `json:"name"`
  }

  type ListVFolderRes struct {
  	Page miso.Paging `json:"paging"`
  	Payload []ListedVFolder `json:"payload"`
  }

  type ListedVFolder struct {
  	Id int `json:"id"`
  	FolderNo string `json:"folderNo"`
  	Name string `json:"name"`
  	CreateTime util.ETime `json:"createTime"`
  	CreateBy string `json:"createBy"`
  	UpdateTime util.ETime `json:"updateTime"`
  	UpdateBy string `json:"updateBy"`
  	Ownership string `json:"ownership"`
  }

  func ApiListVFolders(rail miso.Rail, req ListVFolderReq) (ListVFolderRes, error) {
  	var res miso.GnResp[ListVFolderRes]
  	err := miso.NewDynTClient(rail, "/open/api/vfolder/list", "vfm").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat ListVFolderRes
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
  export interface ListVFolderReq {
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
    data?: ListVFolderRes;
  }

  export interface ListVFolderRes {
    paging?: Paging;
    payload?: ListedVFolder[];
  }

  export interface Paging {
    limit?: number;                // page limit
    page?: number;                 // page number, 1-based
    total?: number;                // total count
  }

  export interface ListedVFolder {
    id?: number;
    folderNo?: string;
    name?: string;
    createTime?: number;
    createBy?: string;
    updateTime?: number;
    updateBy?: string;
    ownership?: string;
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

  listVFolders() {
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
  }
  ```

## POST /open/api/vfolder/create

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type CreateVFolderReq struct {
  	Name string `json:"name"`
  }

  func ApiCreateVFolder(rail miso.Rail, req CreateVFolderReq) (string, error) {
  	var res miso.GnResp[string]
  	err := miso.NewDynTClient(rail, "/open/api/vfolder/create", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface CreateVFolderReq {
    name?: string;
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

  createVFolder() {
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
  }
  ```

## POST /open/api/vfolder/file/add

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type AddFileToVfolderReq struct {
  	FolderNo string `json:"folderNo"`
  	FileKeys []string `json:"fileKeys"`
  }

  func ApiVFolderAddFile(rail miso.Rail, req AddFileToVfolderReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/vfolder/file/add", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface AddFileToVfolderReq {
    folderNo?: string;
    fileKeys?: string[];
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

  vFolderAddFile() {
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
  }
  ```

## POST /open/api/vfolder/file/remove

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type RemoveFileFromVfolderReq struct {
  	FolderNo string `json:"folderNo"`
  	FileKeys []string `json:"fileKeys"`
  }

  func ApiVFolderRemoveFile(rail miso.Rail, req RemoveFileFromVfolderReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/vfolder/file/remove", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface RemoveFileFromVfolderReq {
    folderNo?: string;
    fileKeys?: string[];
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

  vFolderRemoveFile() {
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
  }
  ```

## POST /open/api/vfolder/share

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ShareVfolderReq struct {
  	FolderNo string `json:"folderNo"`
  	Username string `json:"username"`
  }

  func ApiShareVFolder(rail miso.Rail, req ShareVfolderReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/vfolder/share", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface ShareVfolderReq {
    folderNo?: string;
    username?: string;
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

  shareVFolder() {
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
  }
  ```

## POST /open/api/vfolder/access/remove

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type RemoveGrantedFolderAccessReq struct {
  	FolderNo string `json:"folderNo"`
  	UserNo string `json:"userNo"`
  }

  func ApiRemoveVFolderAccess(rail miso.Rail, req RemoveGrantedFolderAccessReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/vfolder/access/remove", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface RemoveGrantedFolderAccessReq {
    folderNo?: string;
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

  removeVFolderAccess() {
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
  }
  ```

## POST /open/api/vfolder/granted/list

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListGrantedFolderAccessReq struct {
  	Page miso.Paging `json:"paging"`
  	FolderNo string `json:"folderNo"`
  }

  type ListGrantedFolderAccessRes struct {
  	Page miso.Paging `json:"paging"`
  	Payload []ListedFolderAccess `json:"payload"`
  }

  type ListedFolderAccess struct {
  	UserNo string `json:"userNo"`
  	Username string `json:"username"`
  	CreateTime util.ETime `json:"createTime"`
  }

  func ApiListVFolderAccess(rail miso.Rail, req ListGrantedFolderAccessReq) (ListGrantedFolderAccessRes, error) {
  	var res miso.GnResp[ListGrantedFolderAccessRes]
  	err := miso.NewDynTClient(rail, "/open/api/vfolder/granted/list", "vfm").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat ListGrantedFolderAccessRes
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
  export interface ListGrantedFolderAccessReq {
    paging?: Paging;
    folderNo?: string;
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
    data?: ListGrantedFolderAccessRes;
  }

  export interface ListGrantedFolderAccessRes {
    paging?: Paging;
    payload?: ListedFolderAccess[];
  }

  export interface Paging {
    limit?: number;                // page limit
    page?: number;                 // page number, 1-based
    total?: number;                // total count
  }

  export interface ListedFolderAccess {
    userNo?: string;
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

  listVFolderAccess() {
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
  }
  ```

## POST /open/api/vfolder/remove

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type RemoveVFolderReq struct {
  	FolderNo string `json:"folderNo"`
  }

  func ApiRemoveVFolder(rail miso.Rail, req RemoveVFolderReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/vfolder/remove", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface RemoveVFolderReq {
    folderNo?: string;
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

  removeVFolder() {
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
  }
  ```

## GET /open/api/gallery/brief/owned

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type VGalleryBrief struct {
  	GalleryNo string `json:"galleryNo"`
  	Name string `json:"name"`
  }

  func ApiListGalleryBriefs(rail miso.Rail) ([]VGalleryBrief, error) {
  	var res miso.GnResp[[]VGalleryBrief]
  	err := miso.NewDynTClient(rail, "/open/api/gallery/brief/owned", "vfm").
  		Get().
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat []VGalleryBrief
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
    data?: VGalleryBrief[];
  }

  export interface VGalleryBrief {
    galleryNo?: string;
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

  listGalleryBriefs() {
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
  }
  ```

## POST /open/api/gallery/new

- Description: Create new gallery
- Bound to Resource: `"manage-files"`
- JSON Request:
    - "name": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type CreateGalleryCmd struct {
  	Name string `json:"name"`      // Required.
  }

  type Gallery struct {
  	Id int64 `json:"id"`
  	GalleryNo string `json:"galleryNo"`
  	UserNo string `json:"userNo"`
  	Name string `json:"name"`
  	DirFileKey string `json:"dirFileKey"`
  	CreateTime util.ETime `json:"createTime"`
  	CreateBy string `json:"createBy"`
  	UpdateTime util.ETime `json:"updateTime"`
  	UpdateBy string `json:"updateBy"`
  	IsDel bool `json:"isDel"`
  }

  func ApiCreateGallery(rail miso.Rail, req CreateGalleryCmd) (*Gallery, error) {
  	var res miso.GnResp[*Gallery]
  	err := miso.NewDynTClient(rail, "/open/api/gallery/new", "vfm").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		return nil, err
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
  export interface CreateGalleryCmd {
    name?: string;                 // Required.
  }

  export interface Resp {
    errorCode?: string;            // error code
    msg?: string;                  // message
    error?: boolean;               // whether the request was successful
    data?: Gallery;
  }

  export interface Gallery {
    id?: number;
    galleryNo?: string;
    userNo?: string;
    name?: string;
    dirFileKey?: string;
    createTime?: number;
    createBy?: string;
    updateTime?: number;
    updateBy?: string;
    isDel?: boolean;
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

  createGallery() {
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
  }
  ```

## POST /open/api/gallery/update

- Description: Update gallery
- Bound to Resource: `"manage-files"`
- JSON Request:
    - "galleryNo": (string) Required.
    - "name": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type UpdateGalleryCmd struct {
  	GalleryNo string `json:"galleryNo"` // Required.
  	Name string `json:"name"`      // Required.
  }

  func ApiUpdateGallery(rail miso.Rail, req UpdateGalleryCmd) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/gallery/update", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface UpdateGalleryCmd {
    galleryNo?: string;            // Required.
    name?: string;                 // Required.
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

  updateGallery() {
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
  }
  ```

## POST /open/api/gallery/delete

- Description: Delete gallery
- Bound to Resource: `"manage-files"`
- JSON Request:
    - "galleryNo": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type DeleteGalleryCmd struct {
  	GalleryNo string `json:"galleryNo"` // Required.
  }

  func ApiDeleteGallery(rail miso.Rail, req DeleteGalleryCmd) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/gallery/delete", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface DeleteGalleryCmd {
    galleryNo?: string;            // Required.
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

  deleteGallery() {
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
  }
  ```

## POST /open/api/gallery/list

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
        - "dirFileKey": (string) 
- cURL:
  ```sh
  curl -X POST 'http://localhost:8086/open/api/gallery/list' \
    -H 'Content-Type: application/json' \
    -d '{"paging":{"limit":0,"page":0,"total":0}}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListGalleriesCmd struct {
  	Paging miso.Paging `json:"paging"`
  }


  type VGallery struct {
  	ID int64 `json:"id"`
  	GalleryNo string `json:"galleryNo"`
  	UserNo string `json:"userNo"`
  	Name string `json:"name"`
  	CreateBy string `json:"createBy"`
  	UpdateBy string `json:"updateBy"`
  	IsOwner bool `json:"isOwner"`
  	CreateTimeStr string `json:"createTime"`
  	UpdateTimeStr string `json:"updateTime"`
  	DirFileKey string `json:"dirFileKey"`
  }

  func ApiListGalleries(rail miso.Rail, req ListGalleriesCmd) (miso.PageRes[VGallery], error) {
  	var res miso.GnResp[miso.PageRes[VGallery]]
  	err := miso.NewDynTClient(rail, "/open/api/gallery/list", "vfm").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat miso.PageRes[VGallery]
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
  export interface ListGalleriesCmd {
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
    payload?: VGallery[];
  }

  export interface Paging {
    limit?: number;                // page limit
    page?: number;                 // page number, 1-based
    total?: number;                // total count
  }

  export interface VGallery {
    id?: number;
    galleryNo?: string;
    userNo?: string;
    name?: string;
    createBy?: string;
    updateBy?: string;
    isOwner?: boolean;
    createTime?: string;
    updateTime?: string;
    dirFileKey?: string;
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

  listGalleries() {
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
  }
  ```

## POST /open/api/gallery/access/grant

- Description: Grant access to the galleries
- Bound to Resource: `"manage-files"`
- JSON Request:
    - "galleryNo": (string) Required.
    - "username": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type PermitGalleryAccessCmd struct {
  	GalleryNo string `json:"galleryNo"` // Required.
  	Username string `json:"username"` // Required.
  }

  func ApiGranteGalleryAccess(rail miso.Rail, req PermitGalleryAccessCmd) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/gallery/access/grant", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface PermitGalleryAccessCmd {
    galleryNo?: string;            // Required.
    username?: string;             // Required.
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

  granteGalleryAccess() {
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
  }
  ```

## POST /open/api/gallery/access/remove

- Description: Remove access to the galleries
- Bound to Resource: `"manage-files"`
- JSON Request:
    - "galleryNo": (string) Required.
    - "userNo": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type RemoveGalleryAccessCmd struct {
  	GalleryNo string `json:"galleryNo"` // Required.
  	UserNo string `json:"userNo"`  // Required.
  }

  func ApiRemoveGalleryAccess(rail miso.Rail, req RemoveGalleryAccessCmd) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/gallery/access/remove", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface RemoveGalleryAccessCmd {
    galleryNo?: string;            // Required.
    userNo?: string;               // Required.
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

  removeGalleryAccess() {
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
  }
  ```

## POST /open/api/gallery/access/list

- Description: List granted access to the galleries
- Bound to Resource: `"manage-files"`
- JSON Request:
    - "galleryNo": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListGrantedGalleryAccessCmd struct {
  	GalleryNo string `json:"galleryNo"` // Required.
  	Paging miso.Paging `json:"paging"`
  }


  type ListedGalleryAccessRes struct {
  	Id int `json:"id"`
  	GalleryNo string `json:"galleryNo"`
  	UserNo string `json:"userNo"`
  	Username string `json:"username"`
  	CreateTime util.ETime `json:"createTime"`
  }

  func ApiListGalleryAccess(rail miso.Rail, req ListGrantedGalleryAccessCmd) (miso.PageRes[ListedGalleryAccessRes], error) {
  	var res miso.GnResp[miso.PageRes[ListedGalleryAccessRes]]
  	err := miso.NewDynTClient(rail, "/open/api/gallery/access/list", "vfm").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat miso.PageRes[ListedGalleryAccessRes]
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
  export interface ListGrantedGalleryAccessCmd {
    galleryNo?: string;            // Required.
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
    payload?: ListedGalleryAccessRes[];
  }

  export interface Paging {
    limit?: number;                // page limit
    page?: number;                 // page number, 1-based
    total?: number;                // total count
  }

  export interface ListedGalleryAccessRes {
    id?: number;
    galleryNo?: string;
    userNo?: string;
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

  listGalleryAccess() {
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
  }
  ```

## POST /open/api/gallery/images

- Description: List images of gallery
- Bound to Resource: `"manage-files"`
- JSON Request:
    - "galleryNo": (string) Required.
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
        - "fileKey": (string) 
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListGalleryImagesCmd struct {
  	GalleryNo string `json:"galleryNo"` // Required.
  	Paging miso.Paging `json:"paging"`
  }

  type ListGalleryImagesResp struct {
  	Images []ImageInfo `json:"images"`
  	Paging miso.Paging `json:"paging"`
  }

  type ImageInfo struct {
  	FileKey string `json:"fileKey"`
  	ThumbnailToken string `json:"thumbnailToken"`
  	FileTempToken string `json:"fileTempToken"`
  }

  func ApiListGalleryImages(rail miso.Rail, req ListGalleryImagesCmd) (*ListGalleryImagesResp, error) {
  	var res miso.GnResp[*ListGalleryImagesResp]
  	err := miso.NewDynTClient(rail, "/open/api/gallery/images", "vfm").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		return nil, err
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
  export interface ListGalleryImagesCmd {
    galleryNo?: string;            // Required.
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
    data?: ListGalleryImagesResp;
  }

  export interface ListGalleryImagesResp {
    images?: ImageInfo[];
    paging?: Paging;
  }

  export interface ImageInfo {
    fileKey?: string;
    thumbnailToken?: string;
    fileTempToken?: string;
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

  listGalleryImages() {
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
  }
  ```

## POST /open/api/gallery/image/transfer

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type TransferGalleryImageReq struct {
  	Images []CreateGalleryImageCmd `json:"images"`
  }

  type CreateGalleryImageCmd struct {
  	GalleryNo string `json:"galleryNo"`
  	Name string `json:"name"`
  	FileKey string `json:"fileKey"`
  }

  func ApiTransferGalleryImage(rail miso.Rail, req TransferGalleryImageReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/gallery/image/transfer", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface TransferGalleryImageReq {
    images?: CreateGalleryImageCmd[];
  }

  export interface CreateGalleryImageCmd {
    galleryNo?: string;
    name?: string;
    fileKey?: string;
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

  transferGalleryImage() {
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
  }
  ```

## POST /open/api/versioned-file/list

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ApiListVerFileReq struct {
  	Paging miso.Paging `json:"paging"`
  	Name *string `json:"name"`     // file name
  }


  type ApiListVerFileRes struct {
  	VerFileId string `json:"verFileId"` // versioned file id
  	Name string `json:"name"`      // file name
  	FileKey string `json:"fileKey"` // file key
  	SizeInBytes int64 `json:"sizeInBytes"` // size in bytes
  	UploadTime util.ETime `json:"uploadTime"` // last upload time
  	CreateTime util.ETime `json:"createTime"` // create time of the versioned file record
  	UpdateTime util.ETime `json:"updateTime"` // Update time of the versioned file record
  	Thumbnail string `json:"thumbnail"` // thumbnail token
  }

  func ApiListVersionedFile(rail miso.Rail, req ApiListVerFileReq) (miso.PageRes[ApiListVerFileRes], error) {
  	var res miso.GnResp[miso.PageRes[ApiListVerFileRes]]
  	err := miso.NewDynTClient(rail, "/open/api/versioned-file/list", "vfm").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat miso.PageRes[ApiListVerFileRes]
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
  export interface ApiListVerFileReq {
    paging?: Paging;
    name?: string;                 // file name
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
    payload?: ApiListVerFileRes[];
  }

  export interface Paging {
    limit?: number;                // page limit
    page?: number;                 // page number, 1-based
    total?: number;                // total count
  }

  export interface ApiListVerFileRes {
    verFileId?: string;            // versioned file id
    name?: string;                 // file name
    fileKey?: string;              // file key
    sizeInBytes?: number;          // size in bytes
    uploadTime?: number;           // last upload time
    createTime?: number;           // create time of the versioned file record
    updateTime?: number;           // Update time of the versioned file record
    thumbnail?: string;            // thumbnail token
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

  listVersionedFile() {
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
  }
  ```

## POST /open/api/versioned-file/history

- Description: List versioned file history
- Bound to Resource: `"manage-files"`
- JSON Request:
    - "paging": (Paging) paging params
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
    - "verFileId": (string) versioned file id. Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ApiListVerFileHistoryReq struct {
  	Paging miso.Paging `json:"paging"`
  	VerFileId string `json:"verFileId"` // versioned file id. Required.
  }


  type ApiListVerFileHistoryRes struct {
  	Name string `json:"name"`      // file name
  	FileKey string `json:"fileKey"` // file key
  	SizeInBytes int64 `json:"sizeInBytes"` // size in bytes
  	UploadTime util.ETime `json:"uploadTime"` // last upload time
  	Thumbnail string `json:"thumbnail"` // thumbnail token
  }

  func ApiListVersionedFileHistory(rail miso.Rail, req ApiListVerFileHistoryReq) (miso.PageRes[ApiListVerFileHistoryRes], error) {
  	var res miso.GnResp[miso.PageRes[ApiListVerFileHistoryRes]]
  	err := miso.NewDynTClient(rail, "/open/api/versioned-file/history", "vfm").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat miso.PageRes[ApiListVerFileHistoryRes]
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
  export interface ApiListVerFileHistoryReq {
    paging?: Paging;
    verFileId?: string;            // versioned file id. Required.
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
    payload?: ApiListVerFileHistoryRes[];
  }

  export interface Paging {
    limit?: number;                // page limit
    page?: number;                 // page number, 1-based
    total?: number;                // total count
  }

  export interface ApiListVerFileHistoryRes {
    name?: string;                 // file name
    fileKey?: string;              // file key
    sizeInBytes?: number;          // size in bytes
    uploadTime?: number;           // last upload time
    thumbnail?: string;            // thumbnail token
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

  listVersionedFileHistory() {
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
  }
  ```

## POST /open/api/versioned-file/accumulated-size

- Description: Query versioned file log accumulated size
- Bound to Resource: `"manage-files"`
- JSON Request:
    - "verFileId": (string) versioned file id. Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ApiQryVerFileAccuSizeReq struct {
  	VerFileId string `json:"verFileId"` // versioned file id. Required.
  }

  type ApiQryVerFileAccuSizeRes struct {
  	SizeInBytes int64 `json:"sizeInBytes"` // total size in bytes
  }

  func ApiQryVersionedFileAccuSize(rail miso.Rail, req ApiQryVerFileAccuSizeReq) (ApiQryVerFileAccuSizeRes, error) {
  	var res miso.GnResp[ApiQryVerFileAccuSizeRes]
  	err := miso.NewDynTClient(rail, "/open/api/versioned-file/accumulated-size", "vfm").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat ApiQryVerFileAccuSizeRes
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
  export interface ApiQryVerFileAccuSizeReq {
    verFileId?: string;            // versioned file id. Required.
  }

  export interface Resp {
    errorCode?: string;            // error code
    msg?: string;                  // message
    error?: boolean;               // whether the request was successful
    data?: ApiQryVerFileAccuSizeRes;
  }

  export interface ApiQryVerFileAccuSizeRes {
    sizeInBytes?: number;          // total size in bytes
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

  qryVersionedFileAccuSize() {
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
  }
  ```

## POST /open/api/versioned-file/create

- Description: Create versioned file
- Bound to Resource: `"manage-files"`
- JSON Request:
    - "filename": (string) Required.
    - "fstoreFileId": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ApiCreateVerFileReq struct {
  	Filename string `json:"filename"` // Required.
  	FakeFstoreFileId string `json:"fstoreFileId"` // Required.
  }

  type ApiCreateVerFileRes struct {
  	VerFileId string `json:"verFileId"` // Versioned File Id
  }

  func ApiCreateVersionedFile(rail miso.Rail, req ApiCreateVerFileReq) (ApiCreateVerFileRes, error) {
  	var res miso.GnResp[ApiCreateVerFileRes]
  	err := miso.NewDynTClient(rail, "/open/api/versioned-file/create", "vfm").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat ApiCreateVerFileRes
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
  export interface ApiCreateVerFileReq {
    filename?: string;             // Required.
    fstoreFileId?: string;         // Required.
  }

  export interface Resp {
    errorCode?: string;            // error code
    msg?: string;                  // message
    error?: boolean;               // whether the request was successful
    data?: ApiCreateVerFileRes;
  }

  export interface ApiCreateVerFileRes {
    verFileId?: string;            // Versioned File Id
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

  createVersionedFile() {
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
  }
  ```

## POST /open/api/versioned-file/update

- Description: Update versioned file
- Bound to Resource: `"manage-files"`
- JSON Request:
    - "verFileId": (string) versioned file id. Required.
    - "filename": (string) Required.
    - "fstoreFileId": (string) Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ApiUpdateVerFileReq struct {
  	VerFileId string `json:"verFileId"` // versioned file id. Required.
  	Filename string `json:"filename"` // Required.
  	FakeFstoreFileId string `json:"fstoreFileId"` // Required.
  }

  func ApiUpdateVersionedFile(rail miso.Rail, req ApiUpdateVerFileReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/versioned-file/update", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface ApiUpdateVerFileReq {
    verFileId?: string;            // versioned file id. Required.
    filename?: string;             // Required.
    fstoreFileId?: string;         // Required.
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

  updateVersionedFile() {
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
  }
  ```

## POST /open/api/versioned-file/delete

- Description: Delete versioned file
- Bound to Resource: `"manage-files"`
- JSON Request:
    - "verFileId": (string) Versioned File Id. Required.
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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ApiDelVerFileReq struct {
  	VerFileId string `json:"verFileId"` // Versioned File Id. Required.
  }

  func ApiDelVersionedFile(rail miso.Rail, req ApiDelVerFileReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/open/api/versioned-file/delete", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface ApiDelVerFileReq {
    verFileId?: string;            // Versioned File Id. Required.
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

  delVersionedFile() {
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
  }
  ```

## POST /compensate/thumbnail

- Description: Compensate thumbnail generation
- Bound to Resource: `"vfm:server:maintenance"`
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
- cURL:
  ```sh
  curl -X POST 'http://localhost:8086/compensate/thumbnail'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  func ApiCompensateThumbnail(rail miso.Rail) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/compensate/thumbnail", "vfm").
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

  compensateThumbnail() {
    this.http.post<any>(`/vfm/compensate/thumbnail`, null)
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

## POST /compensate/regenerate-video-thumbnails

- Description: Regenerate video thumbnails
- Bound to Resource: `"vfm:server:maintenance"`
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
- cURL:
  ```sh
  curl -X POST 'http://localhost:8086/compensate/regenerate-video-thumbnails'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  func ApiRegenerateVideoThumbnail(rail miso.Rail) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/compensate/regenerate-video-thumbnails", "vfm").
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

  regenerateVideoThumbnail() {
    this.http.post<any>(`/vfm/compensate/regenerate-video-thumbnails`, null)
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

## PUT /bookmark/file/upload

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  func ApiUploadBookmarkFile(rail miso.Rail) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/bookmark/file/upload", "vfm").
  		Put(nil).
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

  uploadBookmarkFile() {
    this.http.put<any>(`/vfm/bookmark/file/upload`, null)
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

## POST /bookmark/list

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListBookmarksReq struct {
  	Name *string `json:"name"`
  	Paging miso.Paging `json:"paging"`
  }

  func ApiListBookmarks(rail miso.Rail, req ListBookmarksReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/bookmark/list", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface ListBookmarksReq {
    name?: string;
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

  listBookmarks() {
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
  }
  ```

## POST /bookmark/remove

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type RemoveBookmarkReq struct {
  	Id int64 `json:"id"`
  }

  func ApiRemoveBookmark(rail miso.Rail, req RemoveBookmarkReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/bookmark/remove", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface RemoveBookmarkReq {
    id?: number;
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

  removeBookmark() {
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
  }
  ```

## POST /bookmark/blacklist/list

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListBookmarksReq struct {
  	Name *string `json:"name"`
  	Paging miso.Paging `json:"paging"`
  }

  func ApiListBlacklistedBookmarks(rail miso.Rail, req ListBookmarksReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/bookmark/blacklist/list", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface ListBookmarksReq {
    name?: string;
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

  listBlacklistedBookmarks() {
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
  }
  ```

## POST /bookmark/blacklist/remove

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

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type RemoveBookmarkReq struct {
  	Id int64 `json:"id"`
  }

  func ApiRemoveBookmarkBlacklist(rail miso.Rail, req RemoveBookmarkReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/bookmark/blacklist/remove", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface RemoveBookmarkReq {
    id?: number;
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

  removeBookmarkBlacklist() {
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
  }
  ```

## GET /history/list-browse-history

- Description: List user browse history
- Bound to Resource: `"manage-files"`
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": ([]vfm.ListBrowseRecordRes) response data
      - "time": (int64) 
      - "fileKey": (string) 
      - "name": (string) 
      - "thumbnailToken": (string) 
      - "deleted": (bool) 
- cURL:
  ```sh
  curl -X GET 'http://localhost:8086/history/list-browse-history'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type ListBrowseRecordRes struct {
  	Time util.ETime `json:"time"`
  	FileKey string `json:"fileKey"`
  	Name string `json:"name"`
  	ThumbnailToken string `json:"thumbnailToken"`
  	Deleted bool `json:"deleted"`
  }

  func ApiListBrowseHistory(rail miso.Rail) ([]ListBrowseRecordRes, error) {
  	var res miso.GnResp[[]ListBrowseRecordRes]
  	err := miso.NewDynTClient(rail, "/history/list-browse-history", "vfm").
  		Get().
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat []ListBrowseRecordRes
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
    data?: ListBrowseRecordRes[];
  }

  export interface ListBrowseRecordRes {
    time?: number;
    fileKey?: string;
    name?: string;
    thumbnailToken?: string;
    deleted?: boolean;
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

  listBrowseHistory() {
    this.http.get<any>(`/vfm/history/list-browse-history`)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: ListBrowseRecordRes[] = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
  }
  ```

## POST /history/record-browse-history

- Description: Record user browse history, only files that are directly owned by the user is recorded
- Bound to Resource: `"manage-files"`
- JSON Request:
    - "fileKey": (string) Required.
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
- cURL:
  ```sh
  curl -X POST 'http://localhost:8086/history/record-browse-history' \
    -H 'Content-Type: application/json' \
    -d '{"fileKey":""}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type RecordBrowseHistoryReq struct {
  	FileKey string `json:"fileKey"` // Required.
  }

  func ApiRecordBrowseHistory(rail miso.Rail, req RecordBrowseHistoryReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/history/record-browse-history", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface RecordBrowseHistoryReq {
    fileKey?: string;              // Required.
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

  recordBrowseHistory() {
    let req: RecordBrowseHistoryReq | null = null;
    this.http.post<any>(`/vfm/history/record-browse-history`, req)
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

## GET /maintenance/status

- Description: Check server maintenance status
- Bound to Resource: `"vfm:server:maintenance"`
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (MaintenanceStatus) response data
      - "underMaintenance": (bool) 
- cURL:
  ```sh
  curl -X GET 'http://localhost:8086/maintenance/status'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type MaintenanceStatus struct {
  	UnderMaintenance bool `json:"underMaintenance"`
  }

  func ApiFetchMaintenanceStatus(rail miso.Rail) (MaintenanceStatus, error) {
  	var res miso.GnResp[MaintenanceStatus]
  	err := miso.NewDynTClient(rail, "/maintenance/status", "vfm").
  		Get().
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat MaintenanceStatus
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
    this.http.get<any>(`/vfm/maintenance/status`)
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

## POST /internal/v1/file/create

- Description: Internal endpoint, System create file
- JSON Request:
    - "filename": (string) 
    - "fstoreFileId": (string) 
    - "parentFile": (string) 
    - "userNo": (string) 
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (string) response data
- cURL:
  ```sh
  curl -X POST 'http://localhost:8086/internal/v1/file/create' \
    -H 'Content-Type: application/json' \
    -d '{"filename":"","fstoreFileId":"","parentFile":"","userNo":""}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type SysCreateFileReq struct {
  	Filename string `json:"filename"`
  	FakeFstoreFileId string `json:"fstoreFileId"`
  	ParentFile string `json:"parentFile"`
  	UserNo string `json:"userNo"`
  }

  func ApiSysCreateFile(rail miso.Rail, req SysCreateFileReq) (string, error) {
  	var res miso.GnResp[string]
  	err := miso.NewDynTClient(rail, "/internal/v1/file/create", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface SysCreateFileReq {
    filename?: string;
    fstoreFileId?: string;
    parentFile?: string;
    userNo?: string;
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

  sysCreateFile() {
    let req: SysCreateFileReq | null = null;
    this.http.post<any>(`/vfm/internal/v1/file/create`, req)
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

## GET /internal/file/upload/duplication/preflight

- Description: Internal endpoint, Preflight check for duplicate file uploads
  - Query Parameter:
  - "fileName": 
  - "parentFileKey": 
  - "userNo": 
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (bool) response data
- cURL:
  ```sh
  curl -X GET 'http://localhost:8086/internal/file/upload/duplication/preflight?fileName=&parentFileKey=&userNo='
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  func ApiInternalCheckDuplicate(rail miso.Rail, fileName string, parentFileKey string, userNo string) (bool, error) {
  	var res miso.GnResp[bool]
  	err := miso.NewDynTClient(rail, "/internal/file/upload/duplication/preflight", "vfm").
  		AddQueryParams("fileName", fileName).
  		AddQueryParams("parentFileKey", parentFileKey).
  		AddQueryParams("userNo", userNo).
  		Get().
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		return false, err
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
    data?: boolean;                // response data
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

  internalCheckDuplicate() {
    let fileName: any | null = null;
    let parentFileKey: any | null = null;
    let userNo: any | null = null;
    this.http.get<any>(`/vfm/internal/file/upload/duplication/preflight?fileName=${fileName}&parentFileKey=${parentFileKey}&userNo=${userNo}`)
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
  }
  ```

## POST /internal/file/check-access

- Description: Internal endpoint, Check if user has access to the file
- JSON Request:
    - "fileKey": (string) 
    - "userNo": (string) 
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
- cURL:
  ```sh
  curl -X POST 'http://localhost:8086/internal/file/check-access' \
    -H 'Content-Type: application/json' \
    -d '{"fileKey":"","userNo":""}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type InternalCheckFileAccessReq struct {
  	FileKey string `json:"fileKey"`
  	UserNo string `json:"userNo"`
  }

  func ApiInternalCheckFileAccess(rail miso.Rail, req InternalCheckFileAccessReq) error {
  	var res miso.GnResp[any]
  	err := miso.NewDynTClient(rail, "/internal/file/check-access", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface InternalCheckFileAccessReq {
    fileKey?: string;
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

  internalCheckFileAccess() {
    let req: InternalCheckFileAccessReq | null = null;
    this.http.post<any>(`/vfm/internal/file/check-access`, req)
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

## POST /internal/file/fetch-info

- Description: Internal endpoint. Fetch file info.
- JSON Request:
    - "fileKey": (string) 
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (InternalFetchFileInfoRes) response data
      - "name": (string) 
      - "uploadTime": (int64) 
      - "sizeInBytes": (int64) 
      - "fileType": (string) 
- cURL:
  ```sh
  curl -X POST 'http://localhost:8086/internal/file/fetch-info' \
    -H 'Content-Type: application/json' \
    -d '{"fileKey":""}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type InternalFetchFileInfoReq struct {
  	FileKey string `json:"fileKey"`
  }

  type InternalFetchFileInfoRes struct {
  	Name string `json:"name"`
  	UploadTime util.ETime `json:"uploadTime"`
  	SizeInBytes int64 `json:"sizeInBytes"`
  	FileType string `json:"fileType"`
  }

  func ApiInternalFetchFileInfo(rail miso.Rail, req InternalFetchFileInfoReq) (InternalFetchFileInfoRes, error) {
  	var res miso.GnResp[InternalFetchFileInfoRes]
  	err := miso.NewDynTClient(rail, "/internal/file/fetch-info", "vfm").
  		PostJson(req).
  		Json(&res)
  	if err != nil {
  		rail.Errorf("Request failed, %v", err)
  		var dat InternalFetchFileInfoRes
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
  export interface InternalFetchFileInfoReq {
    fileKey?: string;
  }

  export interface Resp {
    errorCode?: string;            // error code
    msg?: string;                  // message
    error?: boolean;               // whether the request was successful
    data?: InternalFetchFileInfoRes;
  }

  export interface InternalFetchFileInfoRes {
    name?: string;
    uploadTime?: number;
    sizeInBytes?: number;
    fileType?: string;
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

  internalFetchFileInfo() {
    let req: InternalFetchFileInfoReq | null = null;
    this.http.post<any>(`/vfm/internal/file/fetch-info`, req)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 })
            return;
          }
          let dat: InternalFetchFileInfoRes = resp.data;
        },
        error: (err) => {
          console.log(err)
          this.snackBar.open("Request failed, unknown error", "ok", { duration: 3000 })
        }
      });
  }
  ```

## POST /internal/v1/file/make-dir

- Description: Internal endpoint, System make directory.
- JSON Request:
    - "parentFile": (string) 
    - "userNo": (string) Required.
    - "name": (string) Required.
- JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (string) response data
- cURL:
  ```sh
  curl -X POST 'http://localhost:8086/internal/v1/file/make-dir' \
    -H 'Content-Type: application/json' \
    -d '{"name":"","parentFile":"","userNo":""}'
  ```

- Miso HTTP Client (experimental, demo may not work):
  ```go
  type SysMakeDirReq struct {
  	ParentFile string `json:"parentFile"`
  	UserNo string `json:"userNo"`  // Required.
  	Name string `json:"name"`      // Required.
  }

  func ApiSysMakeDir(rail miso.Rail, req SysMakeDirReq) (string, error) {
  	var res miso.GnResp[string]
  	err := miso.NewDynTClient(rail, "/internal/v1/file/make-dir", "vfm").
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

- JSON Request / Response Object In TypeScript:
  ```ts
  export interface SysMakeDirReq {
    parentFile?: string;
    userNo?: string;               // Required.
    name?: string;                 // Required.
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

  sysMakeDir() {
    let req: SysMakeDirReq | null = null;
    this.http.post<any>(`/vfm/internal/v1/file/make-dir`, req)
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
  curl -X GET 'http://localhost:8086/auth/resource'
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

  func SendRequest(rail miso.Rail) (ResourceInfoRes, error) {
  	var res miso.GnResp[ResourceInfoRes]
  	err := miso.NewDynTClient(rail, "/auth/resource", "vfm").
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
    this.http.get<ResourceInfoRes>(`/vfm/auth/resource`)
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
