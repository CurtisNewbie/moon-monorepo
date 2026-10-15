package postbox

import (
	"time"

	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/user-vault/internal/config"
)

// countCoalescer folds a burst of notification count changes for one user into a single
// trailing flush per coalescing window.
//
// The window is owned by Redis, not by this process: the notification change stream is a
// shared durable queue with competing consumers, so with more than one replica a user's
// events can land on different replicas, and only a shared gate sees the whole burst.
// There is deliberately no in-memory per-user map — the injected gate is authoritative
// and is consulted for every change.
type countCoalescer struct {
	window   time.Duration
	gate     func(rail miso.Rail, userNo string) (bool, error)
	flush    func(rail miso.Rail, userNo string)
	schedule func(d time.Duration, f func())
}

// defaultCountCoalesceWindow is the floor a resolved window is clamped to: a Redis TTL of
// zero means "no expiry", which would make the coalescing key never expire and suppress that
// user's announcements permanently, so a non-positive window must never reach the gate.
const defaultCountCoalesceWindow = time.Second

// newCountCoalescer resolves the coalescing window, clamping a missing or non-positive value
// to defaultCountCoalesceWindow. This is the single place the window is resolved: the caller
// builds the gate from the coalescer's resolved window, so the Redis TTL can never be zero.
// The gate is assigned by the caller once the window is known.
func newCountCoalescer(window time.Duration, flush func(rail miso.Rail, userNo string), schedule func(d time.Duration, f func())) *countCoalescer {
	if window <= 0 {
		window = defaultCountCoalesceWindow
	}
	return &countCoalescer{
		window:   window,
		flush:    flush,
		schedule: schedule,
	}
}

// Notify announces a count change for a user. It claims the user's coalescing window
// through the gate and, only if this caller owns the window, schedules one flush a window
// later. A caller that does not own the window drops the change.
func (c *countCoalescer) Notify(rail miso.Rail, userNo string) {
	claimed, err := c.gate(rail, userNo)
	if err != nil {
		rail.Errorf("Failed to claim notification count coalesce window, %v, %v", userNo, err)
		return
	}
	if !claimed {
		return
	}
	// The flush runs after the change-stream listener returns, so it must not inherit that
	// listener's context: detach and start a new span (miso rail async idiom).
	deferredRail := rail.NewCtx().NextSpanId()
	c.schedule(c.window, func() { c.flush(deferredRail, userNo) })
}

// notifCountCoalescer is the process-wide coalescer wired to the notification change
// stream. It is initialized at bootstrap by InitNotifCountCoalescer.
var notifCountCoalescer *countCoalescer

// InitNotifCountCoalescer wires the coalescer to the real Redis gate, the cache-evict +
// publish flush, and a time.AfterFunc scheduler. Coalescing is always on; there is no off
// switch. The window comes from user-vault.notification-count-coalesce-window with a
// one-second default, clamped so a missing or non-positive value can never become a zero
// (no-expiry) Redis TTL.
func InitNotifCountCoalescer(rail miso.Rail) error {
	c := newCountCoalescer(
		miso.GetPropDuration(config.PropNotificationCountCoalesceWindow),
		newCountFlush(evictUserNotifCount, publishUserNotifCountChanged),
		func(d time.Duration, f func()) { time.AfterFunc(d, f) },
	)
	// Build the gate from the already-resolved window so the SetNX TTL shares the clamp.
	c.gate = newRedisCountGate(c.window)
	notifCountCoalescer = c
	return nil
}

func countCoalesceKey(userNo string) string {
	return "user-vault:notification:count-coalesce:" + userNo
}

// newRedisCountGate is the authoritative, cross-replica window gate: the first caller to
// claim a user's key owns the window, and the key expires with the window so expiry needs
// no separate cleanup.
func newRedisCountGate(window time.Duration) func(rail miso.Rail, userNo string) (bool, error) {
	return func(rail miso.Rail, userNo string) (bool, error) {
		return redis.SetNX(rail, countCoalesceKey(userNo), "1", window)
	}
}

// newCountFlush evicts the user's count cache and, only if the eviction succeeds, publishes
// the user id to the count-changed channel. Publishing after a failed eviction would make
// subscribers recompute from the still-stale cache and push a wrong count that persists
// until the cache expires — the opposite of why eviction-before-publish exists.
func newCountFlush(evict func(rail miso.Rail, userNo string) error, publish func(rail miso.Rail, userNo string) error) func(rail miso.Rail, userNo string) {
	return func(rail miso.Rail, userNo string) {
		if err := evict(rail, userNo); err != nil {
			rail.Errorf("Failed to evict user notification count cache, skipping count-change publish, %v, %v", userNo, err)
			return
		}
		if err := publish(rail, userNo); err != nil {
			rail.Errorf("Failed to publish user notification count change to %v, %v, %v", userNotifCountChangedChannel, userNo, err)
		}
	}
}

func evictUserNotifCount(rail miso.Rail, userNo string) error {
	return userNotifCountCache.Del(rail, userNo)
}

func publishUserNotifCountChanged(rail miso.Rail, userNo string) error {
	return redis.GetRedis().Publish(rail.Context(), userNotifCountChangedChannel, userNo).Err()
}
