package logbot

import (
	"sync"
	"time"

	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/atom"
	"github.com/curtisnewbie/miso/util/heap"
	"github.com/sirupsen/logrus"
)

var (
	mergedLogger                     = logrus.New()
	mergedLogMu                      = sync.Mutex{}
	mergedLogs   *heap.Heap[LogLine] = heap.New(1024, func(iv, jv LogLine) bool {
		return iv.Time.Before(jv.Time)
	})
	mergedLogFlushTicker = miso.NewTickRuner(1*time.Second, flushMergedLogs)
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
	mergedLogFlushTicker.Start()
	miso.AddShutdownHook(func() {
		mergedLogFlushTicker.Stop()
		flushMergedLogs()
	})
}

type PlainStrFormatter struct {
}

func (p PlainStrFormatter) Format(e *logrus.Entry) ([]byte, error) {
	return []byte(e.Message), nil
}

func AppendMergedLog(ll LogLine) {
	mergedLogMu.Lock()
	defer mergedLogMu.Unlock()
	mergedLogs.Push(ll)
}

func flushMergedLogs() {
	mergedLogMu.Lock()
	defer mergedLogMu.Unlock()

	now := atom.Now()
	offset := now.Add(-3 * time.Second)

	for mergedLogs.Len() > 0 {
		if mergedLogs.Peek().Time.After(offset) {
			return
		}
		ll := mergedLogs.Pop()
		mergedLogger.Printf("[%11s] - %v", ll.App, ll.OriginLine)
	}
}
