package watcher

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestNewDebounce(t *testing.T) {
	var (
		count int32
	)

	f1 := func() {
		atomic.AddInt32(&count, 1)
	}

	fn := NewDebounce(100*time.Millisecond, f1)

	for i := 0; i < 10; i++ {
		fn()
	}

	time.Sleep(time.Millisecond * 200)

	if count > 1 {
		t.Fatal("count check fail")
	}
}
