package postbox

import (
	"sync"
	"testing"
	"time"

	"github.com/curtisnewbie/miso/errs"
	"github.com/curtisnewbie/miso/miso"
)

var errTestEvict = errs.NewErrf("evict failed")

// fakeClock is a deterministic scheduler for the coalescer: it records scheduled
// callbacks and fires those whose deadline has passed when time is advanced.
type fakeClock struct {
	mu     sync.Mutex
	now    time.Duration
	timers []*fakeTimer
}

type fakeTimer struct {
	at    time.Duration
	f     func()
	fired bool
}

func newFakeClock() *fakeClock { return &fakeClock{} }

func (c *fakeClock) schedule(d time.Duration, f func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.timers = append(c.timers, &fakeTimer{at: c.now + d, f: f})
}

func (c *fakeClock) advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now += d
	for _, t := range c.timers {
		if !t.fired && t.at <= c.now {
			t.fired = true
			t.f()
		}
	}
}

// fakeWindowGate stands in for the Redis SET NX EX gate: it opens a window for a user on
// the first claim and rejects further claims until the window elapses.
type fakeWindowGate struct {
	clock    *fakeClock
	window   time.Duration
	open     map[string]time.Duration
	consults int
	claims   int
}

func newFakeWindowGate(clock *fakeClock, window time.Duration) *fakeWindowGate {
	return &fakeWindowGate{clock: clock, window: window, open: map[string]time.Duration{}}
}

func (g *fakeWindowGate) claim(rail miso.Rail, userNo string) (bool, error) {
	g.consults++
	if exp, ok := g.open[userNo]; ok && g.clock.now < exp {
		return false, nil
	}
	g.open[userNo] = g.clock.now + g.window
	g.claims++
	return true, nil
}

// newTestCoalescer constructs a coalescer and assigns the gate, mirroring how production
// builds the gate after the window is resolved.
func newTestCoalescer(window time.Duration, gate *fakeWindowGate, flush func(rail miso.Rail, userNo string), schedule func(time.Duration, func())) *countCoalescer {
	c := newCountCoalescer(window, flush, schedule)
	c.gate = gate.claim
	return c
}

func TestCountCoalescerBurstInOneWindowFlushesOnce(t *testing.T) {
	clock := newFakeClock()
	gate := newFakeWindowGate(clock, time.Second)
	var flushes []string
	flush := func(rail miso.Rail, userNo string) { flushes = append(flushes, userNo) }

	c := newTestCoalescer(time.Second, gate, flush, clock.schedule)

	for i := 0; i < 1000; i++ {
		c.Notify(miso.EmptyRail(), "user-1")
	}
	if len(flushes) != 0 {
		t.Fatalf("expected no flush before the window elapses, got %d", len(flushes))
	}

	clock.advance(time.Second)
	if len(flushes) != 1 {
		t.Fatalf("expected exactly one flush for a burst in one window, got %d", len(flushes))
	}
	if flushes[0] != "user-1" {
		t.Fatalf("expected the flush to carry user-1, got %q", flushes[0])
	}
}

// A single isolated change must still reach the client: the window is a delay, not a
// suppressor, so the flush fires once when the window elapses.
func TestCountCoalescerSingleChangeFlushesOnce(t *testing.T) {
	clock := newFakeClock()
	gate := newFakeWindowGate(clock, time.Second)
	var flushes int
	c := newTestCoalescer(time.Second, gate, func(rail miso.Rail, userNo string) { flushes++ }, clock.schedule)

	c.Notify(miso.EmptyRail(), "user-1")
	if flushes != 0 {
		t.Fatalf("expected no flush before the window elapses, got %d", flushes)
	}

	clock.advance(time.Second)
	if flushes != 1 {
		t.Fatalf("expected a single change to flush once, got %d", flushes)
	}
}

func TestCountCoalescerBurstSpanningWindowsFlushesPerWindow(t *testing.T) {
	clock := newFakeClock()
	gate := newFakeWindowGate(clock, time.Second)
	var flushes []string
	flush := func(rail miso.Rail, userNo string) { flushes = append(flushes, userNo) }

	c := newTestCoalescer(time.Second, gate, flush, clock.schedule)

	const windows = 3
	for w := 0; w < windows; w++ {
		for i := 0; i < 5; i++ {
			c.Notify(miso.EmptyRail(), "user-1")
		}
		clock.advance(time.Second)
	}

	if len(flushes) != windows {
		t.Fatalf("expected one flush per window (%d), got %d", windows, len(flushes))
	}
}

// The Redis gate is authoritative and shared across replicas, so it must see every change
// — the coalescer keeps no in-memory per-user state. A burst therefore consults the gate
// once per change while claiming (and flushing) only once per window.
func TestCountCoalescerGateSeesEveryChangeClaimsOncePerWindow(t *testing.T) {
	clock := newFakeClock()
	gate := newFakeWindowGate(clock, time.Second)
	var flushes int
	flush := func(rail miso.Rail, userNo string) { flushes++ }

	c := newTestCoalescer(time.Second, gate, flush, clock.schedule)

	const changes = 1000
	for i := 0; i < changes; i++ {
		c.Notify(miso.EmptyRail(), "user-1")
	}
	clock.advance(time.Second)

	if gate.consults != changes {
		t.Fatalf("expected the gate to be consulted once per change (%d), got %d", changes, gate.consults)
	}
	if gate.claims != 1 {
		t.Fatalf("expected one claim per window, got %d", gate.claims)
	}
	if flushes != 1 {
		t.Fatalf("expected one flush, got %d", flushes)
	}
}

// A non-positive window would make the gate's SET NX TTL zero, which means "no expiry" in
// go-redis and would suppress that user forever. It must be clamped to the default.
func TestCountCoalescerClampsNonPositiveWindow(t *testing.T) {
	for _, window := range []time.Duration{0, -time.Second} {
		var scheduled time.Duration
		schedule := func(d time.Duration, f func()) { scheduled = d }

		c := newCountCoalescer(window, func(rail miso.Rail, userNo string) {}, schedule)
		c.gate = func(rail miso.Rail, userNo string) (bool, error) { return true, nil }
		c.Notify(miso.EmptyRail(), "user-1")

		if scheduled != defaultCountCoalesceWindow {
			t.Fatalf("window %v: expected scheduler delay clamped to %v, got %v", window, defaultCountCoalesceWindow, scheduled)
		}
	}
}

func TestCountFlushEvictsBeforePublish(t *testing.T) {
	var order []string
	evict := func(rail miso.Rail, userNo string) error { order = append(order, "evict"); return nil }
	publish := func(rail miso.Rail, userNo string) error { order = append(order, "publish"); return nil }

	newCountFlush(evict, publish)(miso.EmptyRail(), "user-1")

	if len(order) != 2 || order[0] != "evict" || order[1] != "publish" {
		t.Fatalf("expected evict before publish, got %v", order)
	}
}

// Publishing after a failed eviction would make subscribers recompute from the still-stale
// cache and push a wrong count until the cache expires. The flush must skip the publish.
func TestCountFlushSkipsPublishWhenEvictFails(t *testing.T) {
	published := false
	evict := func(rail miso.Rail, userNo string) error { return errTestEvict }
	publish := func(rail miso.Rail, userNo string) error { published = true; return nil }

	newCountFlush(evict, publish)(miso.EmptyRail(), "user-1")

	if published {
		t.Fatal("expected the publish to be skipped when the cache eviction fails")
	}
}
