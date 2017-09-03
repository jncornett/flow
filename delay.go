package flow

import (
	"time"
)

// Delay works like DelayDLQ, but does not maintain a dead-letter queue.
func Delay(in <-chan interface{}, t time.Duration) <-chan interface{} {
	return DelayDLQ(in, t, nil)
}

// DelayDLQ delays input.
//  x-y-z-------a-b-c
//  x           a
func DelayDLQ(in <-chan interface{}, t time.Duration, dlq chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	rl := make(chan struct{}, 1)
	go func() {
		for x := range in {
			select {
			case rl <- struct{}{}: // 'lock'
				time.AfterFunc(t, func() { <-rl }) // 'unlock'
				out <- x
			default: // drop it on the flo-o-or
				if dlq != nil {
					dlq <- x
				}
			}
		}
	}()
	return out
}
