package api

import (
	"github.com/curtisnewbie/miso/middleware/rabbit"
)

var (
	// Pipeline to trigger async image thumbnail generation.
	//
	// Reply api.ImageCompressReplyEvent when the processing succeeds.
	GenImgThumbnailPipeline = rabbit.NewEventPipeline[ImgThumbnailTriggerEvent]("event.bus.fstore.image.compress.processing").
				LogPayload().
				MaxRetry(10).
				Document("GenImgThumbnailPipeline", "Pipeline to trigger async image thumbnail generation, will reply api.ImageCompressReplyEvent when the processing succeeds.", "fstore")

	// Pipeline to trigger async video thumbnail generation.
	//
	// Reply api.GenVideoThumbnailReplyEvent when the processing succeeds.
	GenVidThumbnailPipeline = rabbit.NewEventPipeline[VidThumbnailTriggerEvent]("event.bus.fstore.video.thumbnail.processing").
				LogPayload().
				MaxRetry(10).
				Document("GenVidThumbnailPipeline", "Pipeline to trigger async video thumbnail generation, will reply api.GenVideoThumbnailReplyEvent when the processing succeeds.", "fstore")
)

// Event sent to hammer to trigger an vidoe thumbnail generation.
type VidThumbnailTriggerEvent struct {
	Identifier string `desc:"dentifier" json:"identifier"`
	FileId     string `desc:"file id from mini-fstore" json:"fileId"`
	ReplyTo    string `desc:"event bus that will receive event about the generated video thumbnail." json:"replyTo"`
}

// Event sent to hammer to trigger an image compression.
type ImgThumbnailTriggerEvent struct {
	Identifier string `desc:"identifier" json:"identifier"`
	FileId     string `desc:"file id from mini-fstore" json:"fileId"`
	ReplyTo    string `desc:"event bus that will receive event about the generated image thumbnail." json:"replyTo"`
}

// Event replied from hammer about the compressed image.
type ImageCompressReplyEvent struct {
	Identifier string `json:"identifier"` // identifier
	FileId     string `json:"fileId"`     // file id from mini-fstore
}

// Event replied from hammer about the generated video thumbnail.
type GenVideoThumbnailReplyEvent struct {
	Identifier string `json:"identifier"` // identifier
	FileId     string `json:"fileId"`     // file id from mini-fstore
}

type UnzipFileReplyEvent struct {
	ZipFileId  string     `json:"zipFileId"`
	ZipEntries []ZipEntry `json:"zipEntries"`
	Extra      string     `json:"extra"`
}

type ZipEntry struct {
	FileId string `json:"fileId"`
	Md5    string `json:"md5"`
	Name   string `json:"name"`
	Size   int64  `json:"size"`
}
