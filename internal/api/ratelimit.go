package api

import (
	"context"
	"sync"
	"time"
)

// rateLimiter is a minimal token-bucket-ish pacer. The Bot API exposes no quota headers
// (only 429 + retry_after), so we pace at a fixed requests-per-second and adapt reactively:
// halve the rate on a 429, then restore it gradually on sustained success (GOAL.md §1).
type rateLimiter struct {
	mu       sync.Mutex
	interval time.Duration // minimum spacing between requests at the current rate
	base     time.Duration // the configured spacing (the ceiling we restore toward)
	maxIval  time.Duration // the floor on rate (longest spacing) after repeated 429s
	last     time.Time
	okStreak int

	now   func() time.Time // injectable clock for tests
	sleep func(context.Context, time.Duration) error
}

// newRateLimiter builds a limiter for rps requests/second. rps <= 0 disables pacing.
func newRateLimiter(rps float64) *rateLimiter {
	var interval time.Duration
	if rps > 0 {
		interval = time.Duration(float64(time.Second) / rps)
	}
	return &rateLimiter{
		interval: interval,
		base:     interval,
		maxIval:  interval * 8, // never crawl slower than 1/8 the configured rate
		now:      time.Now,
		sleep:    sleepCtx,
	}
}

// wait blocks until the next request is allowed, or ctx is cancelled.
func (r *rateLimiter) wait(ctx context.Context) error {
	r.mu.Lock()
	if r.interval <= 0 {
		r.mu.Unlock()
		return ctx.Err()
	}
	now := r.now()
	var delay time.Duration
	if !r.last.IsZero() {
		if next := r.last.Add(r.interval); next.After(now) {
			delay = next.Sub(now)
		}
	}
	r.last = now.Add(delay)
	r.mu.Unlock()

	if delay <= 0 {
		return ctx.Err()
	}
	return r.sleep(ctx, delay)
}

// penalize halves the rate (doubles the spacing) after a 429, bounded by maxIval.
func (r *rateLimiter) penalize() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.interval <= 0 {
		return
	}
	r.okStreak = 0
	r.interval *= 2
	if r.interval > r.maxIval {
		r.interval = r.maxIval
	}
}

// reward restores the rate gradually after a streak of successes, never past the base rate.
func (r *rateLimiter) reward() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.interval <= 0 || r.interval <= r.base {
		return
	}
	r.okStreak++
	if r.okStreak >= 5 {
		r.okStreak = 0
		r.interval = r.interval * 3 / 4
		if r.interval < r.base {
			r.interval = r.base
		}
	}
}

// sleepCtx sleeps for d but returns early if ctx is cancelled (so Ctrl-C is honored mid-wait).
func sleepCtx(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
