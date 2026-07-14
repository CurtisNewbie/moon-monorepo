package postbox

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/miso"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

const (
	wsWriteWait  = 10 * time.Second
	wsPongWait   = 60 * time.Second
	wsPingPeriod = (wsPongWait * 9) / 10
	wsMaxMsgSize = 512
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type wsMessage struct {
	Type string `json:"type"`
	Data int    `json:"data"`
}

type wsConn struct {
	conn   *websocket.Conn
	userNo string
	send   chan wsMessage
}

type WSManager struct {
	mu    sync.RWMutex
	conns map[string]map[*wsConn]struct{}
}

var wsMgr *WSManager

func InitWSManager() {
	wsMgr = &WSManager{
		conns: make(map[string]map[*wsConn]struct{}),
	}
}

func (m *WSManager) subscribe(userNo string, c *wsConn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.conns[userNo] == nil {
		m.conns[userNo] = make(map[*wsConn]struct{})
	}
	m.conns[userNo][c] = struct{}{}
}

func (m *WSManager) unsubscribe(userNo string, c *wsConn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if cs, ok := m.conns[userNo]; ok {
		delete(cs, c)
		close(c.send)
		if len(cs) == 0 {
			delete(m.conns, userNo)
		}
	}
}

func (m *WSManager) connCount(userNo string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.conns[userNo])
}

// PushCount pushes a count update to all WS connections for a user.
func (m *WSManager) PushCount(userNo string, count int) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if cs, ok := m.conns[userNo]; ok {
		msg := wsMessage{Type: "count", Data: count}
		miso.Infof("WS notification: pushing count=%d, user=%v, conns=%d", count, userNo, len(cs))
		for c := range cs {
			select {
			case c.send <- msg:
			default:
			}
		}
	}
}

// HandleWS handles a WebSocket connection for notification count updates.
func HandleWS(rail miso.Rail, user flow.User, db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	rail.Infof("Upgrading ws, %v", user.Username)
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		rail.Errorf("WS notification: upgrade failed, user=%v, err=%v", user.UserNo, err)
		return
	}
	rail.Infof("Upgraded ws")

	userNo := user.UserNo
	rail.Infof("WS notification: connected, user=%v", userNo)

	wsc := &wsConn{
		conn:   conn,
		userNo: userNo,
		send:   make(chan wsMessage, 256),
	}

	wsMgr.subscribe(userNo, wsc)
	rail.Infof("WS notification: subscribed, user=%v, conns=%d", userNo, wsMgr.connCount(userNo))

	defer func() {
		wsMgr.unsubscribe(userNo, wsc)
		conn.Close()
		rail.Infof("WS notification: unsubscribed, user=%v, conns=%d", userNo, wsMgr.connCount(userNo))
	}()

	// Push initial count
	if count, err := CachedCountNotification(rail, db, user); err == nil {
		rail.Infof("WS notification: pushing initial count=%d, user=%v", count, userNo)
		wsc.send <- wsMessage{Type: "count", Data: count}
	} else {
		rail.Errorf("WS notification: failed to load initial count, user=%v, err=%v", userNo, err)
	}

	// Start writer goroutine
	go wsc.writePump(rail)

	// Block on reader (handles pong, close messages)
	wsc.readPump(rail)
}

func (c *wsConn) readPump(rail miso.Rail) {
	defer c.conn.Close()
	rail.Debugf("WS notification: read pump started, user=%v", c.userNo)

	c.conn.SetReadLimit(wsMaxMsgSize)
	c.conn.SetReadDeadline(time.Now().Add(wsPongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(wsPongWait))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				rail.Infof("WS notification: read error, user=%v, err=%v", c.userNo, err)
			}
			break
		}
	}
	rail.Debugf("WS notification: read pump stopped, user=%v", c.userNo)
}

func (c *wsConn) writePump(rail miso.Rail) {
	ticker := time.NewTicker(wsPingPeriod)
	rail.Debugf("WS notification: write pump started, user=%v", c.userNo)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		rail.Debugf("WS notification: write pump stopped, user=%v", c.userNo)
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
			data, err := json.Marshal(msg)
			if err != nil {
				rail.Errorf("WS notification: marshal error, user=%v, err=%v", c.userNo, err)
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				rail.Infof("WS notification: write error, user=%v, err=%v", c.userNo, err)
				return
			}
			rail.Debugf("WS notification: sent %v=%d, user=%v", msg.Type, msg.Data, c.userNo)
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
