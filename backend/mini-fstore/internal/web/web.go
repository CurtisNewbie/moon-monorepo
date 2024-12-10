package web

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/mini-fstore/internal/config"
	"github.com/curtisnewbie/mini-fstore/internal/fstore"
	"github.com/curtisnewbie/mini-fstore/internal/metrics"
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/auth"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"github.com/gin-gonic/gin"
)

const (
	headerAuthorization = "Authorization"
)

const (
	ResCodeFstoreUpload      = "fstore-upload"
	ResCodeFstoreMaintenance = "fstore:server:maintenance"
)

func PrepareWebServer(rail miso.Rail) error {

	miso.AddInterceptor(func(c *gin.Context, next func()) {
		url := c.Request.RequestURI

		// endpoints for file backup
		if strings.HasPrefix(url, "/backup") {

			// disabled
			if !miso.GetPropBool(config.PropEnableFstoreBackup) || miso.GetPropStr(config.PropBackupAuthSecret) == "" {
				miso.Infof("Reject request to %v, backup endpoint disabled", url)
				c.AbortWithStatus(404)
				return
			}
			// not authorized
			if err := fstore.CheckBackupAuth(rail, c.Request.Header.Get(headerAuthorization)); err != nil {
				miso.Infof("Reject request to %v, request not authorized", url)
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		next()
	})

	auth.ExposeResourceInfo([]auth.Resource{
		{Name: "Fstore File Upload", Code: ResCodeFstoreUpload},
		{Name: "Fstore Server Maintenance", Code: ResCodeFstoreMaintenance},
	})

	return nil
}

// Streaming file.
//
//   - misoapi-http: GET /file/stream
//   - misoapi-scope: PUBLIC
//   - misoapi-query-doc: key: temporary file key
//   - misoapi-desc: Media streaming using temporary file key, the file_key's ttl is extended with each subsequent
//     request. This endpoint is expected to be accessible publicly without authorization, since a temporary
//     file_key is generated and used.
func TempKeyStreamFileEp(inb *miso.Inbound) {
	rail := inb.Rail()
	w, r := inb.Unwrap()
	query := r.URL.Query()
	key := strings.TrimSpace(query.Get("key"))
	if key == "" {
		w.WriteHeader(404)
		return
	}

	if e := fstore.StreamFileKey(rail, w, key, parseByteRangeRequest(r)); e != nil {
		rail.Warnf("Failed to stream by fileKey, %v", e)
		w.WriteHeader(404)
		return
	}
}

// Download file.
//
//   - misoapi-http: GET /file/raw
//   - misoapi-scope: PUBLIC
//   - misoapi-query-doc: key: temporary file key
//   - misoapi-desc: Download file using temporary file key. This endpoint is expected to be accessible
//     publicly without authorization, since a temporary file_key is generated and used.
func TempKeyDownloadFileEp(inb *miso.Inbound) {
	rail := inb.Rail()
	w, r := inb.Unwrap()
	key := strings.TrimSpace(r.URL.Query().Get("key"))
	if key == "" {
		w.WriteHeader(404)
		return
	}

	if e := fstore.DownloadFileKey(rail, w, key); e != nil {
		rail.Warnf("Failed to download by fileKey, %v", e)
		w.WriteHeader(404)
		return
	}
}

// Upload file.
//
//   - misoapi-http: PUT /file
//   - misoapi-resource: ref(ResCodeFstoreUpload)
//   - misoapi-header: filename: name of the uploaded file
//   - misoapi-desc: Upload file. A temporary file_id is returned, which should be used to exchange
//     the real file_id
func UploadFileEp(inb *miso.Inbound) (string, error) {
	rail := inb.Rail()
	_, r := inb.Unwrap()
	fname := strings.TrimSpace(r.Header.Get("filename"))
	if fname == "" {
		return "", fstore.ErrFilenameRequired
	}
	if n, err := url.QueryUnescape(fname); err == nil {
		fname = n
	}

	fileId, e := fstore.UploadFile(rail, r.Body, fname)
	if e != nil {
		return "", e
	}

	// generate a random file key for the backend server to retrieve the
	// actual fileId later (this is to prevent user guessing others files' fileId,
	// the fileId should be used internally within the system)
	tempFileId := util.ERand(40)

	cmd := redis.GetRedis().Set("mini-fstore:upload:fileId:"+tempFileId, fileId, 6*time.Hour)
	if cmd.Err() != nil {
		return "", fmt.Errorf("failed to cache the generated fake fileId, %v", e)
	}
	rail.Infof("Generated fake fileId '%v' for '%v'", tempFileId, fileId)

	return tempFileId, nil
}

type FileInfoReq struct {
	FileId       string `form:"fileId" desc:"actual file_id of the file record"`
	UploadFileId string `form:"uploadFileId" desc:"temporary file_id returned when uploading files"`
}

// Fetch file info.
//
//   - misoapi-http: GET /file/info
//   - misoapi-desc: Fetch file info
func GetFileInfoEp(inb *miso.Inbound, req FileInfoReq) (api.FstoreFile, error) {
	// fake fileId for uploaded file
	if req.UploadFileId != "" {
		rcmd := redis.GetRedis().Get("mini-fstore:upload:fileId:" + req.UploadFileId)
		if rcmd.Err() != nil {
			if redis.IsNil(rcmd.Err()) {
				// invalid fileId, or the uploadFileId has expired
				return api.FstoreFile{}, fstore.ErrFileNotFound
			}
			return api.FstoreFile{}, rcmd.Err()
		}
		req.FileId = rcmd.Val() // the cached fileId, the real one
	}

	// using real fileId
	if req.FileId == "" {
		return api.FstoreFile{}, fstore.ErrFileNotFound
	}

	f, ef := fstore.FindFile(mysql.GetMySQL(), req.FileId)
	if ef != nil {
		return api.FstoreFile{}, ef
	}
	if f.IsZero() {
		return api.FstoreFile{}, fstore.ErrFileNotFound
	}
	return api.FstoreFile{
		FileId:     f.FileId,
		Name:       f.Name,
		Status:     f.Status,
		Size:       f.Size,
		Md5:        f.Md5,
		UplTime:    f.UplTime,
		LogDelTime: f.LogDelTime,
		PhyDelTime: f.PhyDelTime,
	}, nil
}

type DownloadFileReq struct {
	FileId   string `form:"fileId" desc:"actual file_id of the file record"`
	Filename string `form:"filename" desc:"the name that will be used when downloading the file"`
}

// Generate random file key for downloading or streaming the file.
//
//   - misoapi-http: GET /file/key
//   - misoapi-desc: Generate temporary file key for downloading and streaming. This endpoint is expected
//     to be called internally by another backend service that validates the ownership of the file properly.
func GenFileKeyEp(inb *miso.Inbound, req DownloadFileReq) (string, error) {
	rail := inb.Rail()
	timer := metrics.GenFileKeyTimer()
	defer timer.ObserveDuration()

	fileId := strings.TrimSpace(req.FileId)
	if fileId == "" {
		return "", fstore.ErrFileNotFound
	}

	filename := req.Filename
	unescaped, err := url.QueryUnescape(req.Filename)
	if err == nil {
		filename = unescaped
	}
	filename = strings.TrimSpace(filename)

	k, re := fstore.RandFileKey(rail, filename, fileId)
	rail.Infof("Generated random key %s for fileId %s (using filename: '%s')", k, fileId, filename)
	return k, re
}

// Parse ByteRange Request.
//
// e.g., bytes = 123-124
func parseByteRangeRequest(r *http.Request) fstore.ByteRange {
	headers := r.Header
	rg := headers.Get("Range") // e.g., Range: bytes = 1-2
	if rg == "" {
		return fstore.ByteRange{Start: 0, End: math.MaxInt64}
	}
	return parseByteRangeHeader(rg)
}

// Parse ByteRange Header.
//
// e.g., bytes=123-124
func parseByteRangeHeader(rangeHeader string) fstore.ByteRange {
	var start int64 = 0
	var end int64 = math.MaxInt64

	eqSplit := strings.Split(rangeHeader, "=") // split by '='
	if len(eqSplit) <= 1 {                     // 'bytes=' or '=1-2', both are illegal
		return fstore.ByteRange{Start: start, End: end}
	}

	brr := []rune(strings.TrimSpace(eqSplit[1]))
	if len(brr) < 1 { // empty byte ranges, illegal
		return fstore.ByteRange{Start: start, End: end}
	}

	dash := -1
	for i := 0; i < len(brr); i++ { // try to find the first '-'
		if brr[i] == '-' {
			dash = i
			break
		}
	}

	if dash == 0 { // the '-2' case, only the end is specified, start will still be 0
		if v, e := strconv.ParseInt(string(brr[dash+1:]), 10, 0); e == nil {
			end = v
		}
	} else if dash == len(brr)-1 { // the '1-' case, only the start is specified, end will be MaxInt64
		if v, e := strconv.ParseInt(string(brr[:dash]), 10, 0); e == nil {
			start = v
		}

	} else if dash < 0 { // the '-' case, both start and end are not specified
		// do nothing

	} else { // '1-2' normal case
		if v, e := strconv.ParseInt(string(brr[:dash]), 10, 0); e == nil {
			start = v
		}

		if v, e := strconv.ParseInt(string(brr[dash+1:]), 10, 0); e == nil {
			end = v
		}
	}
	return fstore.ByteRange{Start: start, End: end}
}

// Download file directly without temporary key.
//
//   - misoapi-http: GET /file/direct
//   - misoapi-query: fileId: actual file_id of the file record
//   - misoapi-desc: Download files directly using file_id. This endpoint is expected to be protected and only used
//     internally by another backend service. Users can eaily steal others file_id and attempt to download the file,
//     so it's better not be exposed to the end users.
func DirectDownloadFileEp(inb *miso.Inbound) {
	rail := inb.Rail()
	w, r := inb.Unwrap()
	fileId := r.URL.Query().Get("fileId")
	if fileId == "" {
		rail.Warnf("missing fileId")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if e := fstore.DownloadFile(rail, w, fileId); e != nil {
		rail.Warnf("Failed to DownloadFile using fileId: %v, %v", fileId, e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type DeleteFileReq struct {
	FileId string `form:"fileId" valid:"notEmpty" desc:"actual file_id of the file record"`
}

// Mark file deleted.
//
//   - misoapi-http: DELETE /file
//   - misoapi-desc: Mark file as deleted.
func DeleteFileEp(inb *miso.Inbound, req DeleteFileReq) (any, error) {
	rail := inb.Rail()
	fileId := strings.TrimSpace(req.FileId)
	if fileId == "" {
		return nil, fstore.ErrFileNotFound
	}
	return nil, fstore.LDelFile(rail, mysql.GetMySQL(), fileId)
}

// Unzip files.
//
//   - misoapi-http: POST /file/unzip
//   - misoapi-desc: Unzip archive, upload all the zip entries, and reply the final results back to the caller
//     asynchronously
func UnzipFileEp(inb *miso.Inbound, req api.UnzipFileReq) (any, error) {
	rail := inb.Rail()
	return nil, fstore.TriggerUnzipFilePipeline(rail, mysql.GetMySQL(), req)
}

// List files (for backup).
//
//   - misoapi-http: POST /backup/file/list
//   - misoapi-desc: Backup tool list files
//   - misoapi-scope: PUBLIC
//   - misoapi-header: Authorization: Basic Authorization
func BackupListFilesEp(inb *miso.Inbound, req fstore.ListBackupFileReq) (fstore.ListBackupFileResp, error) {
	rail := inb.Rail()
	rail.Infof("Backup tool listing files %+v", req)
	return fstore.ListBackupFiles(rail, mysql.GetMySQL(), req)
}

// Download file (for backup).
//
//   - misoapi-http: GET /backup/file/raw
//   - misoapi-desc: Backup tool download file
//   - misoapi-scope: PUBLIC
//   - misoapi-header: Authorization: Basic Authorization
//   - misoapi-query: fileId: actual file_id of the file record
func BackupDownFileEp(inb *miso.Inbound) {
	rail := inb.Rail()
	w, r := inb.Unwrap()
	fileId := strings.TrimSpace(r.URL.Query().Get("fileId"))
	rail.Infof("Backup tool download file, fileId: %v", fileId)

	if e := fstore.DownloadFile(rail, w, fileId); e != nil {
		rail.Errorf("Download file failed, %v", e)
		if errors.Is(e, fstore.ErrFileNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusForbidden)
		return
	}
}

// Delete file that were marked deleted (for maintenance).
//
//   - misoapi-http: POST /maintenance/remove-deleted
//   - misoapi-desc: Remove files that are logically deleted and not linked (symbolically)
//   - misoapi-resource: ref(ResCodeFstoreMaintenance)
func RemoveDeletedFilesEp(inb *miso.Inbound) (any, error) {
	rail := inb.Rail()
	return nil, fstore.RemoveDeletedFiles(rail, mysql.GetMySQL())
}

// Sanitize storage, cleanup dangling files (for maintenance).
//
//   - misoapi-http: POST /maintenance/sanitize-storage
//   - misoapi-desc: Sanitize storage, remove files in storage directory that don't exist in database
//   - misoapi-resource: ref(ResCodeFstoreMaintenance)
func SanitizeStorageEp(inb *miso.Inbound) (any, error) {
	rail := inb.Rail()
	return nil, fstore.SanitizeStorage(rail)
}

// Compute file checksum (for maintenance).
//
//   - misoapi-http: POST /maintenance/compute-checksum
//   - misoapi-desc: Compute files' checksum if absent
func ComputeChecksumEp(inb *miso.Inbound) (any, error) {
	rail := inb.Rail()
	return nil, fstore.ComputeFilesChecksum(rail, mysql.GetMySQL())
}

// Fetch storage info.
//
//   - misoapi-http: Get /storage/info
//   - misoapi-desc: Fetch storage info
//   - misoapi-resource: ref(ResCodeFstoreMaintenance)
func FetchStorageInfoEp(inb *miso.Inbound) (fstore.StorageInfo, error) {
	return fstore.LoadStorageInfo(), nil
}

// Fetch storage usage info.
//
//   - misoapi-http: Get /storage/usage-info
//   - misoapi-desc: Fetch storage usage info
//   - misoapi-resource: ref(ResCodeFstoreMaintenance)
func FetchStorageUsageInfoEp(inb *miso.Inbound) ([]fstore.StorageUsageInfo, error) {
	return fstore.LoadStorageUsageInfo(inb.Rail())
}
