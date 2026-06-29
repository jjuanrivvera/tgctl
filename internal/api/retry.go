package api

import (
	"math"
	"math/rand"
	"time"
)

// retryPolicy controls automatic retries. Two rules from the cliwright standard (GOAL.md §1)
// shape it:
//   - Only auto-retry IDEMPOTENT methods on ambiguous failures (5xx/network), because a
//     timed-out sendMessage may have already delivered — retrying would double-send. The
//     read methods (getMe, getChat, getUpdates...) are flagged idempotent and safe to retry.
//   - A 429 is always safe to retry (even for writes): the request was rejected, not
//     processed, so we honor retry_after and try again.
//
// Backoff uses AWS "full jitter" — random(0, base*2^n) — which is a deliberate design to
// spread retries, not a bug, even though a reviewer may flag the randomness.
type retryPolicy struct {
	maxAttempts int           // total attempts including the first
	base        time.Duration // backoff base
	max         time.Duration // backoff ceiling
	rng         func() float64
}

func defaultRetryPolicy() retryPolicy {
	return retryPolicy{
		maxAttempts: 4,
		base:        300 * time.Millisecond,
		max:         30 * time.Second,
		rng:         rand.Float64, //nolint:gosec // G404: jitter is not a security boundary
	}
}

// decide returns whether to retry and how long to wait first. attempt is 0-based.
//   - status: HTTP status (0 for a transport-level network error).
//   - retryAfter: seconds the server asked us to wait (from a 429 parameters/header), or 0.
//   - idempotent: whether the failed method is safe to replay.
func (p retryPolicy) decide(attempt int, status, retryAfter int, networkErr, idempotent bool) (bool, time.Duration) {
	if attempt >= p.maxAttempts-1 {
		return false, 0
	}

	switch {
	case status == 429:
		// Rate limited: rejected, not processed → always safe to retry.
		wait := p.backoff(attempt)
		if ra := time.Duration(retryAfter) * time.Second; ra > wait {
			wait = ra
		}
		return true, wait
	case status >= 500 && status <= 599:
		if !idempotent {
			return false, 0 // ambiguous for a write — don't risk a duplicate.
		}
		return true, p.backoff(attempt)
	case networkErr:
		if !idempotent {
			return false, 0
		}
		return true, p.backoff(attempt)
	default:
		return false, 0
	}
}

// backoff returns a full-jitter delay: random in [0, base*2^attempt], capped at max.
func (p retryPolicy) backoff(attempt int) time.Duration {
	ceiling := min(float64(p.base)*math.Pow(2, float64(attempt)), float64(p.max))
	r := p.rng
	if r == nil {
		r = rand.Float64 //nolint:gosec // G404: jitter is not a security boundary
	}
	return time.Duration(r() * ceiling)
}
