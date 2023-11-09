package backoff

import (
	"context"
	"sync/atomic"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
)

type Backoff struct {
	backoffAlgo   *backoff.ExponentialBackOff
	Attempts      int32
	MaxAttempts   int32
	linearEnabled bool
}

const (
	MAX_DEFAULT_ATTEMPTS = 3
)

type BackoffCaller func() error

func setBackoffLinearVars(backoffAlgo *backoff.ExponentialBackOff) {
	backoffAlgo.InitialInterval = 600 * time.Millisecond
	backoffAlgo.RandomizationFactor = 0.55
	backoffAlgo.Multiplier = 1.6
}

// NewBackoff - Create a new backoff executor
//
// linearEnabled (bool) - true enable linear ticker runner
// maxAttempts - set the max attempts - if 0 sets to default
func NewBackoff(linearEnabled bool, maxAttempts int32) *Backoff {
	maxAllowedAttempts := maxAttempts
	backoffAlgo := backoff.NewExponentialBackOff()

	if linearEnabled {
		setBackoffLinearVars(backoffAlgo)
	}

	if maxAttempts <= 0 {
		maxAllowedAttempts = MAX_DEFAULT_ATTEMPTS
	}

	return &Backoff{
		backoffAlgo:   backoffAlgo,
		Attempts:      0,
		MaxAttempts:   maxAllowedAttempts,
		linearEnabled: linearEnabled,
	}
}

// ExecuteBackoff - Execute a backoff runner, passing the function to execute
//
// ctx context.Context - current Context
// fn BackoffCaller - func() error - Function with current runtime to execute
func (b *Backoff) ExecuteBackoff(ctx context.Context, fn BackoffCaller) error {
	err := fn()
	if err == nil {
		b.Reset()
		return nil
	}

	time.Sleep(b.backoffAlgo.NextBackOff())

	if b.linearEnabled {
		ticker := NewTicker(b.Attempts)
		ticker.run(ctx)
	}

	if b.Attempts < b.MaxAttempts {
		atomic.AddInt32(&b.Attempts, 1)
		return b.ExecuteBackoff(ctx, fn)
	}

	b.Reset()
	return err
}

func (b *Backoff) Reset() {
	b.Attempts = 0
}
