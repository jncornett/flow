package flow

import (
	"time"
)

func Debounce(in <-chan interface{}, t time.Duration) <-chan interface{} {
	return DebounceDLQ(in, t, nil)
}

func DebounceDLQ(in <-chan interface{}, t time.Duration, dlq chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		var saved interface{}
		var after <-chan time.Time
		for {
			select {
			case tmp := <-in: // save for later
				if saved != nil && dlq != nil {
					dlq <- saved
				}
				saved = tmp
				if after == nil {
					after = time.After(t)
				}
			case <-after:
				out <- saved
				saved = nil
				after = nil
			}
		}
	}()
	return out
}
