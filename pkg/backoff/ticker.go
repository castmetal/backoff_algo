package backoff

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

type Ticker struct {
	C        <-chan time.Time
	c        chan bool
	timer    *time.Ticker
	stop     chan struct{}
	stopOnce sync.Once
}

// NewTicker - create a Ticker to fire
//
// attempts: number of attempts
func NewTicker(attempts int32) *Ticker {
	defaultMultiplier := 1.5        // Linear seconds multiplier
	fixedIntervalMs := float64(500) // Fixed Interval in miliseconds

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randDelay := 50 + r.Float64()*(250-50) + fixedIntervalMs

	return &Ticker{
		C:     make(<-chan time.Time),
		c:     make(chan bool, 1),
		timer: time.NewTicker(time.Duration(defaultMultiplier*float64(attempts))*time.Second + time.Duration(randDelay)*time.Millisecond),
		stop:  make(chan struct{}),
	}
}

// run - fire a ticket to wait
//
// ctx context.Context
func (t *Ticker) run(ctx context.Context) {
	defer close(t.c)

	go t.send()
	defer t.Stop()

	for {
		select {
		case <-t.c:
			return
		case <-t.stop:
			t.c = nil
			return
		case <-ctx.Done():
			return
		}
	}
}

// send - handle the ticker channels
func (t *Ticker) send() {
	for {
		select {
		case _, ok := <-t.timer.C:
			if !ok {
				t.Stop()
				continue
			}
			t.c <- true
		case <-t.stop:
			return
		}
	}
}

// Stop - stops the current ticker runner
func (t *Ticker) Stop() {
	t.timer.Stop()
	t.stopOnce.Do(func() { close(t.stop) })
}
