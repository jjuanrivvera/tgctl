package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetryPolicy_Decide(t *testing.T) {
	p := retryPolicy{maxAttempts: 4, base: time.Second, max: 10 * time.Second, rng: func() float64 { return 1 }}

	t.Run("429 retries regardless of idempotency", func(t *testing.T) {
		retry, _ := p.decide(0, 429, 0, false, false)
		assert.True(t, retry)
	})
	t.Run("429 honors retry_after floor", func(t *testing.T) {
		_, wait := p.decide(0, 429, 7, false, false)
		assert.GreaterOrEqual(t, wait, 7*time.Second)
	})
	t.Run("5xx retries only idempotent", func(t *testing.T) {
		r1, _ := p.decide(0, 500, 0, false, true)
		r2, _ := p.decide(0, 500, 0, false, false)
		assert.True(t, r1)
		assert.False(t, r2)
	})
	t.Run("network error retries only idempotent", func(t *testing.T) {
		r1, _ := p.decide(0, 0, 0, true, true)
		r2, _ := p.decide(0, 0, 0, true, false)
		assert.True(t, r1)
		assert.False(t, r2)
	})
	t.Run("4xx (non-429) never retries", func(t *testing.T) {
		retry, _ := p.decide(0, 400, 0, false, true)
		assert.False(t, retry)
	})
	t.Run("stops at maxAttempts", func(t *testing.T) {
		retry, _ := p.decide(3, 500, 0, false, true)
		assert.False(t, retry)
	})
}

func TestRetryPolicy_BackoffBounded(t *testing.T) {
	p := retryPolicy{maxAttempts: 10, base: time.Second, max: 4 * time.Second, rng: func() float64 { return 1 }}
	// rng=1 yields the full ceiling; it must never exceed max.
	for attempt := range 10 {
		assert.LessOrEqual(t, p.backoff(attempt), 4*time.Second)
	}
	// rng=0 yields zero jitter.
	p.rng = func() float64 { return 0 }
	assert.Equal(t, time.Duration(0), p.backoff(5))
}

func TestDefaultRetryPolicy(t *testing.T) {
	p := defaultRetryPolicy()
	assert.Equal(t, 4, p.maxAttempts)
	assert.NotNil(t, p.rng)
}
