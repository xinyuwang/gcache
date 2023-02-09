package gcache

import (
	"fmt"
	"runtime"
	"time"
)

type monitor[T any] struct {
	interval time.Duration
	stop     chan bool
	ticker   *time.Ticker
}

func (m *monitor[T]) Start(g *gcache[T]) {

	m.ticker = time.NewTicker(m.interval)
	for {
		select {
		case <-m.ticker.C:
			g.clearExpiredObject()

		case <-m.stop:
			m.ticker.Stop()
			return
		}
	}
}

func (m *monitor[T]) Stop() {
	m.stop <- true
}

func newMonitor[T any](g *gcache[T], interval time.Duration) *monitor[T] {
	m := &monitor[T]{
		interval: interval,
		stop:     make(chan bool),
		ticker:   nil,
	}

	runtime.SetFinalizer(m, monitorStop[T])
	return m
}

func monitorStop[T any](m *monitor[T]) {
	fmt.Printf("%v\n", m)
	m.Stop()
}
