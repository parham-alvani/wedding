package ratelimit_test

import (
	"testing"
	"time"

	"github.com/parham-alvani/wedding/wedback/internal/infra/ratelimit"
	"github.com/stretchr/testify/assert"
)

func TestAllowsUpToTheLimitThenRefuses(t *testing.T) {
	t.Parallel()

	limiter := ratelimit.New(3, time.Minute)

	for i := range 3 {
		assert.True(t, limiter.Allow("a"), "attempt %d should be allowed", i+1)
	}

	assert.False(t, limiter.Allow("a"), "the fourth attempt is over the limit")
}

func TestKeysAreIndependent(t *testing.T) {
	t.Parallel()

	limiter := ratelimit.New(1, time.Minute)

	assert.True(t, limiter.Allow("a"))
	assert.False(t, limiter.Allow("a"))
	assert.True(t, limiter.Allow("b"), "one guest's attempts must not lock out another")
}

func TestAttemptsExpire(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, time.June, 16, 18, 30, 0, 0, time.UTC)
	limiter := ratelimit.New(2, time.Minute).WithClock(func() time.Time { return now })

	assert.True(t, limiter.Allow("a"))
	assert.True(t, limiter.Allow("a"))
	assert.False(t, limiter.Allow("a"))

	now = now.Add(2 * time.Minute)
	assert.True(t, limiter.Allow("a"), "the window has moved on")
}

func TestResetClearsAKey(t *testing.T) {
	t.Parallel()

	limiter := ratelimit.New(1, time.Minute)

	assert.True(t, limiter.Allow("a"))
	assert.False(t, limiter.Allow("a"))

	limiter.Reset("a")
	assert.True(t, limiter.Allow("a"))
}
