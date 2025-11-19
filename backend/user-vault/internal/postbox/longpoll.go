package postbox

import (
	"context"
	"net/http"
	"sync"
	"time"

	red "github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/async"
	"github.com/curtisnewbie/miso/util/json"
	"github.com/curtisnewbie/miso/util/randutil"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	longPollingHandler           *LongPolling
	userNotifCountChangedChannel = "user-vault:channel:notification:count:changed"
)

func PrepareLongPollHandler(rail miso.Rail) error {
	longPollingHandler = newLongPolling()
	miso.AddShutdownHook(func() { longPollingHandler.Shutdown() })
	miso.PostServerBootstrap(func(rail miso.Rail) error {
		return listenUserNotificationCountChanges(red.GetRedis())
	})
	return nil
}

func listenUserNotificationCountChanges(r *redis.Client) error {
	pubsub := r.Subscribe(context.Background(), userNotifCountChangedChannel)
	c, cancel := context.WithCancel(context.Background())
	miso.AddShutdownHook(func() {
		cancel()
		pubsub.Close()
	})
	go func(ctx context.Context, recv <-chan *redis.Message) {
		for {
			select {
			case <-ctx.Done():
				return
			case m := <-recv:
				userNo := m.Payload
				if userNo != "" {
					longPollingHandler.Notify(userNo)
				}
			}
		}
	}(c, pubsub.Channel())
	return nil
}

func newLongPolling() *LongPolling {
	return &LongPolling{
		mu:   sync.RWMutex{},
		pool: async.NewIOAsyncPool(),
		sub:  map[string]map[string]*LPSub{},
	}
}

type LongPolling struct {
	mu   sync.RWMutex
	pool async.AsyncPool
	sub  map[string]map[string]*LPSub
}

func (l *LongPolling) Shutdown() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for un, usub := range l.sub {
		for k, lps := range usub {
			lps.Write(0)
			delete(usub, k)
		}
		delete(l.sub, un)
	}
}

func (l *LongPolling) Notify(userNo string) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if v, ok := l.sub[userNo]; ok {
		for _, lps := range v {
			select {
			case lps.notified <- struct{}{}:
			default:
			}
		}
	}
}

func (l *LongPolling) Poll(rail miso.Rail, user common.User, db *gorm.DB, w http.ResponseWriter, curr int) {
	userNo := user.UserNo
	lps := NewLongPollSub(user.UserNo, w)
	rail.Infof("User %v subscribes notifications using LongPolling, %v", userNo, lps.id)
	defer rail.Infof("User %v unsubscribed notifications using LongPolling, %v", userNo, lps.id)

	l.mu.Lock()
	if submap, ok := l.sub[user.UserNo]; ok {
		submap[lps.id] = lps
	} else {
		l.sub[user.UserNo] = map[string]*LPSub{lps.id: lps}
	}
	l.mu.Unlock()

	l.pool.Go(func() {
		loadCount := func(alwaysClose bool) (exit bool) {
			next, err := CachedCountNotification(rail, db, user)
			if err == nil && next != curr {
				if err := lps.Write(next); err != nil {
					rail.Errorf("Failed to write to LPSub, lps.id: %v, userId: %v, %v", lps.id, userNo, err)
				}
				exit = true
			} else if alwaysClose {
				v := curr
				if err == nil {
					v = next
				}
				if err := lps.Write(v); err != nil {
					rail.Errorf("Failed to write to LPSub, lps.id: %v, userId: %v, %v", lps.id, userNo, err)
				}
				exit = true
			}
			return
		}

		defer func() {
			l.mu.Lock()
			defer l.mu.Unlock()

			lps.Close()
			delete(l.sub[user.UserNo], lps.id)
			rail.Infof("Remove LongPollSub, response has been written, lps.id: %v, userNo: %v", lps.id, userNo)
		}()

		exit := loadCount(false)
		if exit {
			return
		}

		for {
			select {
			case <-time.After(15 * time.Second):
				loadCount(true) // close no matter what
				return
			case <-lps.notified:
				rail.Infof("LongPolling notified, query latest unread notification count for %v, %v", lps.id, user.UserNo)
				if loadCount(false) {
					return
				}
			case <-rail.Context().Done():
				rail.Infof("Client disconnected, %v, %v", lps.id, userNo)
				return
			}
		}
	})

	lps.Wait() // block until we write response to the client
}

type LPSub struct {
	id          string
	mu          sync.Mutex
	w           http.ResponseWriter
	notified    chan struct{}
	untilClosed chan struct{}
	closed      bool
}

func (l *LPSub) Write(m any) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed {
		return nil
	}

	defer func() {
		l.closed = true
		l.untilClosed <- struct{}{}
	}()

	r := miso.Resp{
		Data:  m,
		Error: false,
	}
	l.w.WriteHeader(http.StatusOK)
	l.w.Header().Add("Content-Type", "application/json")
	return json.EncodeJson(l.w, r)
}

func (l *LPSub) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed {
		return
	}
	l.closed = true
	l.untilClosed <- struct{}{}
}

func (l *LPSub) Wait() {
	<-l.untilClosed
}

func NewLongPollSub(userId string, w http.ResponseWriter) *LPSub {
	return &LPSub{
		id:          randutil.ERand(30),
		mu:          sync.Mutex{},
		w:           w,
		untilClosed: make(chan struct{}),
		notified:    make(chan struct{}, 10),
	}
}

func Poll(rail miso.Rail, user common.User, db *gorm.DB, w http.ResponseWriter, curr int) {
	longPollingHandler.Poll(rail, user, db, w, curr)
}
