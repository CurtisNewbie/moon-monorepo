// auto generated by misoapi v0.1.12-beta.1 at 2024/11/08 17:30:39, please do not modify
package web

import (
	"github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/mini-fstore/internal/fstore"
	"github.com/curtisnewbie/miso/miso"
)

func init() {
	miso.RawGet("/file/stream", TempKeyStreamFileEp).
		Desc("Media streaming using temporary file key, the file_key's ttl is extended with each subsequent request. This endpoint is expected to be accessible publicly without authorization, since a temporary file_key is generated and used.").
		Public().
		DocQueryParam("key", "temporary file key")

	miso.RawGet("/file/raw", TempKeyDownloadFileEp).
		Desc("Download file using temporary file key. This endpoint is expected to be accessible publicly without authorization, since a temporary file_key is generated and used.").
		Public().
		DocQueryParam("key", "temporary file key")

	miso.Put("/file",
		func(inb *miso.Inbound) (string, error) {
			return UploadFileEp(inb)
		}).
		Desc("Upload file. A temporary file_id is returned, which should be used to exchange the real file_id").
		Resource(ResCodeFstoreUpload).
		DocHeader("filename", "name of the uploaded file")

	miso.IGet("/file/info",
		func(inb *miso.Inbound, req FileInfoReq) (api.FstoreFile, error) {
			return GetFileInfoEp(inb, req)
		}).
		Desc("Fetch file info")

	miso.IGet("/file/key",
		func(inb *miso.Inbound, req DownloadFileReq) (string, error) {
			return GenFileKeyEp(inb, req)
		}).
		Desc("Generate temporary file key for downloading and streaming. This endpoint is expected to be called internally by another backend service that validates the ownership of the file properly.")

	miso.RawGet("/file/direct", DirectDownloadFileEp).
		Desc("Download files directly using file_id. This endpoint is expected to be protected and only used internally by another backend service. Users can eaily steal others file_id and attempt to download the file, so it's better not be exposed to the end users.").
		DocQueryParam("fileId", "actual file_id of the file record")

	miso.IDelete("/file",
		func(inb *miso.Inbound, req DeleteFileReq) (any, error) {
			return DeleteFileEp(inb, req)
		}).
		Desc("Mark file as deleted.")

	miso.IPost("/file/unzip",
		func(inb *miso.Inbound, req api.UnzipFileReq) (any, error) {
			return UnzipFileEp(inb, req)
		}).
		Desc("Unzip archive, upload all the zip entries, and reply the final results back to the caller asynchronously")

	miso.IPost("/backup/file/list",
		func(inb *miso.Inbound, req fstore.ListBackupFileReq) (fstore.ListBackupFileResp, error) {
			return BackupListFilesEp(inb, req)
		}).
		Desc("Backup tool list files").
		Public().
		DocHeader("Authorization", "Basic Authorization")

	miso.RawGet("/backup/file/raw", BackupDownFileEp).
		Desc("Backup tool download file").
		Public().
		DocHeader("Authorization", "Basic Authorization").
		DocQueryParam("fileId", "actual file_id of the file record")

	miso.Post("/maintenance/remove-deleted",
		func(inb *miso.Inbound) (any, error) {
			return RemoveDeletedFilesEp(inb)
		}).
		Desc("Remove files that are logically deleted and not linked (symbolically)")

	miso.Post("/maintenance/sanitize-storage",
		func(inb *miso.Inbound) (any, error) {
			return SanitizeStorageEp(inb)
		}).
		Desc("Sanitize storage, remove files in storage directory that don't exist in database")

	miso.Post("/maintenance/compute-checksum",
		func(inb *miso.Inbound) (any, error) {
			return ComputeChecksumEp(inb)
		}).
		Desc("Compute files' checksum if absent")

	miso.Get("/storage/info",
		func(inb *miso.Inbound) (fstore.StorageInfo, error) {
			return FetchStorageInfoEp(inb)
		}).
		Desc("Fetch storage info").
		Resource(ResCodeFstoreFetchStorageInfo)

}