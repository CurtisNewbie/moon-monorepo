package logbot

import (
	"github.com/curtisnewbie/miso/miso"
	"github.com/sirupsen/logrus"
)

const (
	PropMergedLogFilename = "log.merged-file-name"
)

var (
	mergedLogger = logrus.New()
)

func InitMergedLogger() {
	fn := miso.GetPropStr(PropMergedLogFilename)
	if fn == "" {
		return
	}
	out := miso.BuildRollingLogFileWriter(miso.NewRollingLogFileParam{
		Filename:   fn,
		MaxSize:    500,
		MaxAge:     2,
		MaxBackups: 3,
	})
	mergedLogger.SetFormatter(PlainStrFormatter{})
	mergedLogger.SetOutput(out)
}

type PlainStrFormatter struct {
}

func (p PlainStrFormatter) Format(e *logrus.Entry) ([]byte, error) {
	return []byte(e.Message), nil
}
