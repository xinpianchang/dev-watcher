package watcher

import (
	"time"
)

// NewDebounce returns a debounced function
func NewDebounce(d time.Duration, fn func()) func() {
	last := time.Now().Add(-d)
	return func() {
		now := time.Now()
		check := last.Add(d)
		if check.Before(now) {
			last = now
			go fn()
		}
	}
}
