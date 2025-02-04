package metrics

import (
	"github.com/curtisnewbie/miso/miso"
)

var (
	genFileKeyPromHisto        = miso.NewPromHisto("mini_fstore_generate_file_key_duration")
	genImgThumbnailPromHisto   = miso.NewPromHisto("mini_fstore_generate_img_thumbnail_duration")
	genVideoThumbnailPromHisto = miso.NewPromHisto("mini_fstore_generate_video_thumbnail_duration")

	genImgThumbnailPromCounterName   = "mini_fstore_generate_img_thumbnail_count"
	genVideoThumbnailPromCounterName = "mini_fstore_generate_video_thumbnail_count"
)

func GenFileKeyTimer() *miso.HistTimer {
	return miso.NewHistTimer(genFileKeyPromHisto)
}

func IncGenImgThumbnailCounter() {
	miso.NewPromCounter(genImgThumbnailPromCounterName).Inc()
}

func GenImgThumbnailTimer() *miso.HistTimer {
	return miso.NewHistTimer(genImgThumbnailPromHisto)
}

func GenVideoThumbnailTimer() *miso.HistTimer {
	return miso.NewHistTimer(genVideoThumbnailPromHisto)
}

func IncGenVideoThumbnailCounter() {
	miso.NewPromCounter(genVideoThumbnailPromCounterName).Inc()
}
