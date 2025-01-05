package logbot

import (
	"testing"
	"time"

	"github.com/curtisnewbie/miso/util"
)

func TestETimeHeap(t *testing.T) {
	start := util.Now()
	h := util.NewHeap[util.ETime](50, func(iv, jv util.ETime) bool {
		return iv.Before(jv)
	})

	for i := 0; i < 15; i++ {
		tt := start.Add(time.Duration(15-i) * time.Millisecond)
		t.Logf("tt.%v -> %v", i, tt.ToTime().Format("2006-01-02 15:04:05.000"))
		h.Push(tt)
	}

	_ = h.Pop()

	for i := 0; i < 15; i++ {
		tt := start.Add(time.Duration(15-i) * time.Millisecond)
		t.Logf("tt.%v -> %v", i, tt.ToTime().Format("2006-01-02 15:04:05.000"))
		h.Push(tt)
	}

	last := start
	for h.Len() > 0 {
		p := h.Pop()
		t.Log(p.ToTime().Format("2006-01-02 15:04:05.000"))
		if p.Before(last) {
			t.Fatal("wrong order")
		}
		last = p
	}
}
