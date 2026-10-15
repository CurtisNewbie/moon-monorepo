package postbox

import (
	"context"
	"testing"
	"time"

	"github.com/curtisnewbie/event-pump/client"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/user-vault/internal/config"
	goredis "github.com/redis/go-redis/v9"
)

func prePostboxTest(t *testing.T) miso.Rail {
	rail := miso.EmptyRail()
	miso.DefaultReadConfig([]string{"configFile=../../conf.yml"}, rail)
	if _, err := redis.InitRedisFromProp(rail); err != nil {
		t.Log(err)
		t.FailNow()
	}
	return rail
}

// Exercises the Redis-backed gate against a live Redis: the key name, the SET NX claim
// (the second claim in the same window must fail) and the TTL expiry (after the window a
// claim succeeds again). This is the behaviour the unit seam fakes.
func TestNotificationCountCoalesceGateRedis(t *testing.T) {
	rail := prePostboxTest(t)

	const window = time.Second
	userNo := "test-notification-count-coalesce"
	key := countCoalesceKey(userNo)
	gate := newRedisCountGate(window)

	if _, err := redis.GetRedis().Del(rail.Context(), key).Result(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { redis.GetRedis().Del(rail.Context(), key) })

	claimed, err := gate(rail, userNo)
	if err != nil {
		t.Fatal(err)
	}
	if !claimed {
		t.Fatal("expected the first claim to win the window")
	}

	if ok, err := redis.Exists(rail, key); err != nil || !ok {
		t.Fatalf("expected the gate key %q to exist, ok=%v, err=%v", key, ok, err)
	}
	ttl, err := redis.GetRedis().TTL(rail.Context(), key).Result()
	if err != nil {
		t.Fatal(err)
	}
	if ttl <= 0 || ttl > window {
		t.Fatalf("expected a positive TTL within the window, got %v", ttl)
	}

	claimed, err = gate(rail, userNo)
	if err != nil {
		t.Fatal(err)
	}
	if claimed {
		t.Fatal("expected a second claim in the same window to fail")
	}

	time.Sleep(window + 200*time.Millisecond)

	claimed, err = gate(rail, userNo)
	if err != nil {
		t.Fatal(err)
	}
	if !claimed {
		t.Fatal("expected a claim after the window expired to succeed")
	}
}

// Drives the real listener entry point with a synthetic burst of stream events carrying the
// same user_no, against live Redis, and asserts exactly one count-changed announcement
// arrives once the window has elapsed. This covers the real gate + real flush + real
// publish; the binlog -> event-pump path is not exercised (event-pump is not running).
func TestNotificationCountCoalesceWiredBurstAnnouncesOnce(t *testing.T) {
	rail := prePostboxTest(t)
	if err := InitNotifCountCoalescer(rail); err != nil {
		t.Fatal(err)
	}

	got := runCoalesceAnnouncementScenario(t, rail, "test-wired-coalesce-burst", 25)

	if len(got) != 1 {
		t.Fatalf("expected exactly one announcement for a burst, got %d", len(got))
	}
}

// A single event must still reach the client within one window: coalescing delays, it does
// not suppress.
func TestNotificationCountCoalesceWiredSingleEventAnnouncesOnce(t *testing.T) {
	rail := prePostboxTest(t)
	if err := InitNotifCountCoalescer(rail); err != nil {
		t.Fatal(err)
	}

	got := runCoalesceAnnouncementScenario(t, rail, "test-wired-coalesce-single", 1)

	if len(got) != 1 {
		t.Fatalf("expected exactly one announcement for a single event, got %d", len(got))
	}
}

func runCoalesceAnnouncementScenario(t *testing.T, rail miso.Rail, userNo string, events int) []string {
	key := countCoalesceKey(userNo)
	if _, err := redis.GetRedis().Del(rail.Context(), key).Result(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { redis.GetRedis().Del(rail.Context(), key) })

	sub := redis.GetRedis().Subscribe(context.Background(), userNotifCountChangedChannel)
	defer sub.Close()
	if _, err := sub.Receive(context.Background()); err != nil {
		t.Fatal(err)
	}
	ch := sub.Channel()

	for i := 0; i < events; i++ {
		if err := NotifyNotifCountChanged(rail, coalesceStreamEvent(userNo)); err != nil {
			t.Fatal(err)
		}
	}

	window := miso.GetPropDuration(config.PropNotificationCountCoalesceWindow)
	return collectAnnouncements(ch, userNo, window+time.Second)
}

func coalesceStreamEvent(userNo string) client.StreamEvent {
	return client.StreamEvent{
		Table: "notification",
		Type:  "UPD",
		Columns: map[string]client.EventColumn{
			"user_no": {DataType: "varchar", After: userNo},
		},
	}
}

func collectAnnouncements(ch <-chan *goredis.Message, userNo string, wait time.Duration) []string {
	var got []string
	deadline := time.After(wait)
	for {
		select {
		case m := <-ch:
			if m.Payload == userNo {
				got = append(got, m.Payload)
			}
		case <-deadline:
			return got
		}
	}
}
