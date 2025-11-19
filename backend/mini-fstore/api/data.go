package api

import "github.com/curtisnewbie/miso/util/atom"

const (
	FileStatusNormal    = "NORMAL"  // file.status - normal
	FileStatusLogicDel  = "LOG_DEL" // file.status - logically deleted
	FileStatusPhysicDel = "PHY_DEL" // file.status - physically deleted
)

type FetchFileInfoReq struct {
	FileId       string
	UploadFileId string
}

type FstoreFile struct {
	FileId     string     `json:"fileId" desc:"file unique identifier"`
	Name       string     `json:"name" desc:"file name"`
	Status     string     `json:"status" desc:"status, 'NORMAL', 'LOG_DEL' (logically deleted), 'PHY_DEL' (physically deleted)"`
	Size       int64      `json:"size" desc:"file size in bytes"`
	Md5        string     `json:"md5" desc:"MD5 checksum"`
	UplTime    atom.Time  `json:"uplTime" desc:"upload time"`
	LogDelTime *atom.Time `json:"logDelTime" desc:"logically deleted at"`
	PhyDelTime *atom.Time `json:"phyDelTime" desc:"physically deleted at"`
}

type UnzipFileReq struct {
	// zip file's mini-fstore file_id.
	FileId string `valid:"notEmpty" desc:"file_id of zip file" json:"fileId"`

	// rabbitmq exchange (both the exchange and queue must all use the same name, and are bound together using routing key '#').
	//
	// See UnzipFileReplyEvent (reply message body).
	ReplyToEventBus string `valid:"notEmpty" desc:"name of the rabbitmq exchange to reply to, routing_key is '#'" json:"replyToEventBus"`

	// Extra information that will be passed back to the caller in reply event.
	Extra string `desc:"extra information that will be passed around for the caller" json:"extra"`
}
