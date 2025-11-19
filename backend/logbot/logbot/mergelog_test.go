package logbot

import (
	"testing"
	"time"

	"github.com/curtisnewbie/miso/util/atom"
	"github.com/curtisnewbie/miso/util/heap"
)

func TestETimeHeap(t *testing.T) {
	start := atom.Now()
	h := heap.New[atom.Time](50, func(iv, jv atom.Time) bool {
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
