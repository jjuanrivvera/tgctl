package api

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRateLimiter_Paces(t *testing.T) {
	now := time.Unix(0, 0)
	var slept []time.Duration
	r := newRateLimiter(10) // 100ms spacing
	r.now = func() time.Time { return now }
	r.sleep = func(_ context.Context, d time.Duration) error { slept = append(slept, d); return nil }

	require.NoError(t, r.wait(t.Context())) // first call: no wait
	require.NoError(t, r.wait(t.Context())) // second: must pace ~100ms
	require.Len(t, slept, 1)
	assert.InDelta(t, float64(100*time.Millisecond), float64(slept[0]), float64(5*time.Millisecond))
}

func TestRateLimiter_Disabled(t *testing.T) {
	r := newRateLimiter(0)
	require.NoError(t, r.wait(t.Context()))
	require.NoError(t, r.wait(t.Context()))
}

func TestRateLimiter_PenalizeAndReward(t *testing.T) {
	r := newRateLimiter(10)
	base := r.interval
	r.penalize()
	assert.Equal(t, base*2, r.interval)
	r.penalize()
	assert.Equal(t, base*4, r.interval)

	// penalize is bounded by maxIval (8x base).
	for range 10 {
		r.penalize()
	}
	assert.LessOrEqual(t, r.interval, base*8)

	// reward only restores after a streak, never past base.
	for range 100 {
		r.reward()
	}
	assert.Equal(t, base, r.interval)
}

func TestRateLimiter_WaitCancelled(t *testing.T) {
	r := newRateLimiter(1) // 1s spacing → second call would block ~1s
	ctx, cancel := context.WithCancel(t.Context())
	require.NoError(t, r.wait(ctx))
	cancel()
	err := r.wait(ctx)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestSleepCtx(t *testing.T) {
	assert.NoError(t, sleepCtx(t.Context(), 0))
	require.NoError(t, sleepCtx(t.Context(), time.Millisecond))

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	assert.ErrorIs(t, sleepCtx(ctx, time.Hour), context.Canceled)
}
