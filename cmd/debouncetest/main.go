package main

import (
	"fmt"
	"time"

	"github.com/jncornett/flow"
)

const unit = 2 * time.Second

type message struct {
	count int
}

func generator(cancel chan bool) <-chan interface{} {
	out := make(chan interface{})
	t := time.Tick(unit)
	ctr := 0
	go func() {
		for {
			select {
			case <-t:
				out <- &message{ctr}
				ctr++
			case <-cancel:
				break
			}
		}
	}()
	return out
}

func main() {
	fmt.Println("Time unit:", unit)
	g := generator(nil)
	dlq := make(chan interface{}, 1024)
	d := flow.DebounceDLQ(g, 3*unit, dlq)
	for {
		select {
		case x := <-d:
			fmt.Println("Dequeue", x)
		case dl := <-dlq:
			fmt.Println("Discard", dl)
		}
	}
}
