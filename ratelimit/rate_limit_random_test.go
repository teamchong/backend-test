package ratelimit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const Count = 1
const MinInterval = time.Duration(250 * time.Millisecond)
const MaxInterval = time.Duration(750 * time.Millisecond)

func TestRandomRatelimitConfErrors(t *testing.T) {
	conf, err := randomRatelimitConfig().ParseYAML(`count: -1`, nil)
	require.NoError(t, err)

	_, err = NewRandomRateLimitFromConfig(conf)
	require.Error(t, err)

	_, err = randomRatelimitConfig().ParseYAML(`min_interval: nope`, nil)
	require.NoError(t, err)

	_, err = NewRandomRateLimitFromConfig(conf)
	require.Error(t, err)

	_, err = randomRatelimitConfig().ParseYAML(`max_interval: nope`, nil)
	require.NoError(t, err)

	_, err = NewRandomRateLimitFromConfig(conf)
	require.Error(t, err)
}

func TestRandomRatelimitBasic(t *testing.T) {
	conf, err := randomRatelimitConfig().ParseYAML(`
count: 1
min_interval: 250ms
max_interval: 750ms
`, nil)
	require.NoError(t, err)

	rl, err := NewRandomRateLimitFromConfig(conf)
	require.NoError(t, err)

	ctx := context.Background()

	for i := 0; i < Count; i++ {
		period, _ := rl.Access(ctx)
		assert.LessOrEqual(t, period, time.Duration(0))
	}

	if period, _ := rl.Access(ctx); period == 0 {
		t.Error("Expected limit on final request")
	} else if period > MaxInterval {
		t.Errorf("Period beyond max_interval: %v", period)
	}
}

func TestRandomRatelimitRefresh(t *testing.T) {
	conf, err := randomRatelimitConfig().ParseYAML(`
count: 1
min_interval: 250ms
max_interval: 750ms
`, nil)
	require.NoError(t, err)

	rl, err := NewRandomRateLimitFromConfig(conf)
	require.NoError(t, err)

	ctx := context.Background()

	for i := 0; i < Count; i++ {
		period, _ := rl.Access(ctx)
		if period > 0 {
			t.Errorf("Period above zero: %v", period)
		}
	}

	if period, _ := rl.Access(ctx); period == 0 {
		t.Error("Expected limit on final request")
	} else if period > MaxInterval {
		t.Errorf("Period beyond max_interval: %v", period)
	}

	<-time.After(MaxInterval + 15*time.Millisecond)

	for i := 0; i < Count; i++ {
		period, _ := rl.Access(ctx)
		if period != 0 {
			t.Errorf("Rate limited on get %v", period)
		}
	}

	if period, _ := rl.Access(ctx); period == 0 {
		t.Error("Expected limit on final request")
	} else if period > MaxInterval {
		t.Errorf("Period beyond max_interval: %v", period)
	}
}
