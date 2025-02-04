package metrics

import (
	"github.com/curtisnewbie/miso/miso"
)

var (
	genFileKeyPromHisto        = miso.NewPromHisto("mini_fstore_generate_file_key_duration")
	genImgThumbnailPromHisto   = miso.NewPromHisto("mini_fstore_generate_img_thumbnail_duration")
	genVideoThumbnailPromHisto = miso.NewPromHisto("mini_fstore_generate_video_thumbnail_duration")
)

func GenFileKeyTimer() *miso.HistTimer {
	return miso.NewHistTimer(genFileKeyPromHisto)
}

func GenImgThumbnailTimer() *miso.HistTimer {
	return miso.NewHistTimer(genImgThumbnailPromHisto)
}

func GenVideoThumbnailTimer() *miso.HistTimer {
	return miso.NewHistTimer(genVideoThumbnailPromHisto)
}
