package logbot

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/rabbit"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	uvault "github.com/curtisnewbie/user-vault/api"
	red "github.com/go-redis/redis"
	"gorm.io/gorm"
)

const (
	ErrorLogEventBus = "event.bus.logbot.log.error"
)

var (
	logPatternCache = miso.NewLocalCache[*regexp.Regexp]()
)

func init() {
	miso.SetDefProp("logbot.node", "default")
}

func lastPos(rail miso.Rail, app string, nodeName string) (int64, error) {
	cmd := redis.GetRedis().Get(fmt.Sprintf("log-bot:pos:%v:%v", nodeName, app))
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), red.Nil) {
			return 0, nil
		}
		return 0, cmd.Err()
	}

	n, ea := strconv.Atoi(cmd.Val())
	if ea != nil {
		return 0, nil
	}
	if n < 0 {
		n = 0
	}
	return int64(n), nil
}

func recPos(rail miso.Rail, app string, nodeName string, pos int64) error {
	rail.Debugf("app: %v, node: %v, pos: %v", app, nodeName, pos)
	posStr := strconv.FormatInt(pos, 10)
	cmd := redis.GetRedis().Set(fmt.Sprintf("log-bot:pos:%v:%v", nodeName, app), posStr, 0)
	return cmd.Err()
}

func WatchLogFile(rail miso.Rail, wc WatchConfig, nodeName string) error {
	rail.Infof("Watching log file '%v' for app '%v'", wc.File, wc.App)
	f, err := os.Open(wc.File)

	if err != nil {
		if !os.IsNotExist(err) { // is possible that the log file doesn't exist
			return fmt.Errorf("failed to open log file, %v", err)
		}
	}

	if f != nil {
		defer f.Close() // the log file is opened
	}

	pos, el := lastPos(rail, wc.App, nodeName)
	if el != nil {
		return fmt.Errorf("failed to find last pos, %v", el)
	}

	if f != nil && pos > 0 {
		fi, es := f.Stat()
		if es != nil {
			return es
		}

		// the file was truncated
		if pos > fi.Size() {
			pos = 0
		}

		// seek pos
		if pos > 0 {
			_, e := f.Seek(pos, io.SeekStart)
			if e != nil {
				return fmt.Errorf("failed to seek pos, %v", e)
			}
			rail.Infof("Log file '%v' seek to position %v", wc.File, pos)
		}
	}

	// create reader for the file
	var rd *bufio.Reader
	if f != nil {
		rd = bufio.NewReaderSize(f, 1024*16)
	}

	lastRead := time.Now()
	var prevBytesRead int64
	var prevLine string
	var prevLogLine *LogLine // a single log can contain multiple lines

	for {
		if rd == nil {
			time.Sleep(2 * time.Second) // wait for the file to be created

			f, err = os.Open(wc.File)
			if err != nil {
				f = nil
				continue // the file is still not created
			}
			rail.Infof("Opened %v", wc.File)

			// new file, create reader and set pos = 0
			rd = bufio.NewReader(f)
			pos = 0
		}

		// check if the file is still valid
		if time.Since(lastRead) > 15*time.Second {
			rail.Debug("Checking if the file is still valid, ", wc.File)

			reopenFile := false

			fi, es := f.Stat()
			if es != nil {
				// if the file is deleted, es will still be nil
				reopenFile = true
			}

			if !reopenFile {
				// https://stackoverflow.com/questions/53184549/how-to-detect-deleted-file
				nlink := uint64(0)
				if sys := fi.Sys(); sys != nil {
					if stat, ok := sys.(*syscall.Stat_t); ok {
						nlink = uint64(stat.Nlink)
					}
				}
				if nlink < 1 { // no hard links, the underlying file is deleted already
					reopenFile = true
				}
			}

			lastRead = time.Now()

			if reopenFile {
				f.Close()
				rd = nil
				f = nil
				rail.Infof("Closed file '%v' fd", wc.File)
				continue
			}
		}

		didWaitForEOF := false
		line, err := rd.ReadString('\n')

		if err == nil {
			didWaitForEOF = false
			logLine, e := parseLogLine(rail, line, wc.Type)
			if e == nil {

				// we always report the previous log, coz a single log can contain multiple lines
				//
				// in the same log, the first line is of course valid, but the following lines are not
				// we parse each line, the invalid lines are appended to the previous log
				// once we find a valid line then we know the previous log is now complete, and we report it
				if prevLogLine != nil {

					// report the previous log
					if e := reportLine(rail, *prevLogLine, nodeName, wc); e != nil {
						rail.Errorf("Failed to reportLine, logLine: %+v, %v", *prevLogLine, e)
					}

					// append previous to merged logs
					AppendMergedLog(*prevLogLine, wc.App, line)

					// move the position only when we report the previous log
					pos += prevBytesRead
					recPos(rail, wc.App, nodeName, pos)
				}

				prevBytesRead = int64(len(util.UnsafeStr2Byt(line)))
				prevLine = line
				prevLogLine = &logLine

			} else {

				// if current line is not parseable, it's part of previous line
				// we put them together and we parse again
				//
				// 90% of the time, the log is single line
				// so it's better leave it here
				prevBytesRead += int64(len(util.UnsafeStr2Byt(line)))
				prevLine = prevLine + line
				if parsed, ep := parseLogLine(rail, prevLine, wc.Type); ep == nil {
					prevLogLine = &parsed
				}
			}

			lastRead = time.Now()
			continue

		} else if err == io.EOF {

			// report the last log line if we have already waited for EOF
			if prevLogLine != nil && didWaitForEOF {
				if e := reportLine(rail, *prevLogLine, nodeName, wc); e != nil {
					rail.Errorf("Failed to reportLine, logLine: %+v, %v", *prevLogLine, e)
				}
				prevLogLine = nil
				pos += prevBytesRead
				recPos(rail, wc.App, nodeName, pos)
			}

			// Sleep for a shorter period if we have a previous log line to avoid latency
			if prevLogLine != nil {
				time.Sleep(100 * time.Millisecond)
			} else {
				time.Sleep(500 * time.Millisecond)
			}

			didWaitForEOF = true
			continue

		} else {
			rail.Errorf("Failed to read file, %v, %v", wc.File, err)
		}

		if miso.IsShuttingDown() {
			return nil
		}
	}
}

type LogLineEvent struct {
	App     string
	Node    string
	Time    util.ETime
	Level   string
	TraceId string
	SpanId  string
	Caller  string
	Message string
}

type LogLine struct {
	ParseTime  util.ETime
	App        string
	Time       util.ETime
	TimeStr    string
	Level      string
	TraceId    string
	SpanId     string
	Caller     string
	Message    string
	OriginLine string
}

func parseLogLine(rail miso.Rail, line string, typ string) (LogLine, error) {
	patType := miso.GetPropStr("log.pattern." + typ)
	pat, _ := logPatternCache.Get(patType, func(s string) (*regexp.Regexp, error) {
		return regexp.MustCompile(s), nil
	})

	matches := pat.FindStringSubmatch(line)
	if len(matches) < 7 {
		return LogLine{}, fmt.Errorf("doesn't match pattern")
	}

	time, ep := time.ParseInLocation(`2006-01-02 15:04:05.000`, matches[1], time.Local)
	if ep != nil {
		return LogLine{}, fmt.Errorf("time format illegal, %v", ep)
	}

	// only save the first 1000 characters
	msg := matches[6]
	msgRu := []rune(msg)
	if len(msgRu) > 1000 {
		msg = string(msgRu[:1001])
	}

	ll := LogLine{
		OriginLine: line,
		ParseTime:  util.Now(),
		Time:       util.ToETime(time),
		TimeStr:    matches[1],
		Level:      matches[2],
		TraceId:    strings.TrimSpace(matches[3]),
		SpanId:     strings.TrimSpace(matches[4]),
		Caller:     matches[5],
		Message:    msg,
	}
	return ll, nil
}

func reportLine(rail miso.Rail, line LogLine, node string, wc WatchConfig) error {
	if !wc.ReportError || line.Level != "ERROR" {
		return nil
	}

	return rabbit.PubEventBus(rail,
		LogLineEvent{
			App:     wc.App,
			Node:    node,
			Time:    line.Time,
			Level:   line.Level,
			TraceId: line.TraceId,
			SpanId:  line.SpanId,
			Caller:  line.Caller,
			Message: line.Message,
		},
		ErrorLogEventBus,
	)
}

type SaveErrorLogCmd struct {
	Node    string
	App     string
	Caller  string
	TraceId string
	SpanId  string
	ErrMsg  string
	RTime   util.ETime `gorm:"column:rtime"`
}

func SaveErrorLog(rail miso.Rail, evt LogLineEvent) error {
	el := SaveErrorLogCmd{
		Node:    evt.Node,
		App:     evt.App,
		Caller:  evt.Caller,
		TraceId: evt.TraceId,
		SpanId:  evt.SpanId,
		ErrMsg:  evt.Message,
		RTime:   evt.Time,
	}
	err := mysql.GetMySQL().
		Table("error_log").
		Create(&el).
		Error

	if err == nil {
		if cerr := uvault.CreateNotifiByAccessPipeline.Send(rail, uvault.CreateNotifiByAccessEvent{
			Title:   fmt.Sprintf("Logbot - %s has error", evt.App),
			Message: fmt.Sprintf("%s [%s,%s] %s : %s", evt.Time.FormatClassic(), evt.TraceId, evt.SpanId, evt.Caller, evt.Message),
			ResCode: ResourceManageLogbot,
		}); cerr != nil {
			rail.Errorf("failed to create platform notification, %v", cerr)
		}
	}
	return err
}

type ListedErrorLog struct {
	Id      int64      `json:"id"`
	Node    string     `json:"node"`
	App     string     `json:"app"`
	Caller  string     `json:"caller"`
	TraceId string     `json:"traceId"`
	SpanId  string     `json:"spanId"`
	ErrMsg  string     `json:"errMsg"`
	RTime   util.ETime `json:"rtime" gorm:"column:rtime"`
}

type ListErrorLogReq struct {
	App  string      `json:"app"`
	Page miso.Paging `json:"page"`
}

type ListErrorLogResp struct {
	Page    miso.Paging      `json:"page"`
	Payload []ListedErrorLog `json:"payload"`
}

func newListErrorLogsQry(rail miso.Rail, r ListErrorLogReq) *gorm.DB {
	t := mysql.GetMySQL().
		Table("error_log")

	if r.App != "" {
		t = t.Where("app = ?", r.App)
	}

	return t
}

func ListErrorLogs(rail miso.Rail, r ListErrorLogReq) (ListErrorLogResp, error) {
	var listed []ListedErrorLog
	e := newListErrorLogsQry(rail, r).
		Offset(r.Page.GetOffset()).
		Limit(r.Page.GetLimit()).
		Order("rtime desc").
		Scan(&listed).Error

	if e != nil {
		return ListErrorLogResp{}, e
	}

	var total int
	e = newListErrorLogsQry(rail, r).
		Select("count(*)").
		Scan(&total).Error
	if e != nil {
		return ListErrorLogResp{}, e
	}

	return ListErrorLogResp{Page: r.Page.ToRespPage(total), Payload: listed}, nil
}

func RemoveErrorLogsBefore(rail miso.Rail, upperBound time.Time) error {
	rail.Infof("Remove error logs before %s", upperBound)
	return mysql.GetMySQL().Exec("delete from error_log where rtime < ?", upperBound).Error
}
