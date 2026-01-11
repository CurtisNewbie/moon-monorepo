package api

import (
	"io"
	"net/url"

	"github.com/curtisnewbie/miso/errs"
	"github.com/curtisnewbie/miso/miso"
)

var (
	ErrFileNotFound  = errs.NewErrfCode(FileNotFound, "File is not found")
	ErrFileDeleted   = errs.NewErrfCode(FileDeleted, "File has been deleted")
	ErrIllegalFormat = errs.NewErrfCode(IllegalFormat, "Illegal format")
)

func FetchFileInfo(rail miso.Rail, req FetchFileInfoReq) (FstoreFile, error) {
	var r miso.GnResp[FstoreFile]
	err := miso.NewDynClient(rail, "/file/info", "fstore").
		AddQuery("fileId", req.FileId).
		AddQuery("uploadFileId", req.UploadFileId).
		Get().
		Json(&r)
	return r.Data, err
}

func DeleteFile(rail miso.Rail, fileId string) error {
	var r miso.GnResp[any]
	err := miso.NewDynClient(rail, "/file", "fstore").
		AddQuery("fileId", fileId).
		Delete().
		Json(&r)
	return err
}

func GenTempFileKey(rail miso.Rail, fileId string, filename string) (string, error) {
	var r miso.GnResp[string]
	err := miso.NewDynClient(rail, "/file/key", "fstore").
		AddQuery("fileId", fileId).
		AddQuery("filename", url.QueryEscape(filename)).
		Get().
		Json(&r)
	if err != nil {
		return "", err
	}
	return r.Data, nil
}

func DownloadFile(rail miso.Rail, tmpToken string, writer io.Writer) error {
	_, err := miso.NewDynClient(rail, "/file/raw", "fstore").
		AddQuery("key", tmpToken).
		Get().
		WriteTo(writer)
	return err
}

func UploadFile(rail miso.Rail, filename string, dat io.Reader) (string /* uploadFileId */, error) {
	var res miso.GnResp[string]
	err := miso.NewDynClient(rail, "/file", "fstore").
		AddHeaders(map[string]string{"filename": filename}).
		Put(dat).
		Json(&res)
	if err != nil {
		return "", err
	}
	return res.Res()
}

func TriggerFileUnzip(rail miso.Rail, req UnzipFileReq) error {
	var r miso.GnResp[any]
	err := miso.NewDynClient(rail, "/file/unzip", "fstore").
		PostJson(req).
		Json(&r)
	return err
}

type DirectDownloadFileReq struct {
	FileId string
}

func DownloadFileDirect(rail miso.Rail, fileId string, writer io.Writer) error {
	_, err := miso.NewDynClient(rail, "/file/direct", "fstore").
		AddQuery("fileId", fileId).
		Get().
		WriteTo(writer)
	return err
}
