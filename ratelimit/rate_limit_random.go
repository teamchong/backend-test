package ratelimit

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/benthosdev/benthos/v4/public/service"
)

func randomRatelimitConfig() *service.ConfigSpec {
	spec := service.NewConfigSpec().
		Summary(`The random rate limit is X every (Y1 - Y2)ms rate limit that can be shared across any number of components within the pipeline but does not support distributed rate limits across multiple running instances of Benthos.`).
		Field(service.NewIntField("count").
			Description("The maximum number of requests to allow for a given period of time.").
			Default(1)).
		Field(service.NewDurationField("min_interval").
			Description("The min time window to limit requests by.").
			Default("250ms")).
		Field(service.NewDurationField("max_interval").
			Description("The max time window to limit requests by.").
			Default("750ms"))
	return spec
}

func init() {
	err := service.RegisterRateLimit(
		"random", randomRatelimitConfig(),
		func(conf *service.ParsedConfig, mgr *service.Resources) (service.RateLimit, error) {
			return NewRandomRateLimitFromConfig(conf)
		})
	if err != nil {
		panic(err)
	}
}

// NewRandomRateLimitFromConfig create RandomRateLimit from the provided ParsedConfig
func NewRandomRateLimitFromConfig(conf *service.ParsedConfig) (service.RateLimit, error) {
	count, err := conf.FieldInt("count")
	if err != nil {
		return nil, err
	} else if count <= 0 {
		return nil, fmt.Errorf("value for count is invalid %v", count)
	}
	minDuration, err := conf.FieldDuration("min_interval")
	if err != nil {
		return nil, err
	} else if minDuration < 0 {
		return nil, fmt.Errorf("value for min_interval is invalid %v", minDuration)
	}
	maxDuration, err := conf.FieldDuration("max_interval")
	if err != nil {
		return nil, err
	} else if maxDuration < 0 {
		return nil, fmt.Errorf("value for max_interval is invalid %v", maxDuration)
	}
	if minDuration > maxDuration {
		return nil, fmt.Errorf("value for max_interval %v should be larger than min_interval %v", maxDuration, minDuration)
	}
	return NewRandomRateLimit(count, minDuration, maxDuration)
}

//------------------------------------------------------------------------------

// RandomRateLimit store the state for this rate limiter
type RandomRateLimit struct {
	mut         sync.Mutex
	bucket      int
	lastRefresh time.Time
	size        int
	period      time.Duration
}

// NewRandomRateLimit return a new RandomRateLimit using the provided parameters
func NewRandomRateLimit(count int, minInterval time.Duration, maxInterval time.Duration) (*RandomRateLimit, error) {
	if count <= 0 {
		count = 1
	}
	period := time.Duration(0)
	if minInterval > period {
		period = minInterval
	}
	if maxInterval > minInterval {
		period += time.Duration(rand.Int() % (int(maxInterval) - int(minInterval)))
	}
	return &RandomRateLimit{
		bucket:      count,
		lastRefresh: time.Now(),
		size:        count,
		period:      period,
	}, nil
}

// Access will return 0 when the limit is reached
func (r *RandomRateLimit) Access(context.Context) (time.Duration, error) {
	r.mut.Lock()
	r.bucket--

	if r.bucket < 0 {
		r.bucket = 0
		remaining := r.period - time.Since(r.lastRefresh)

		if remaining > 0 {
			r.mut.Unlock()
			return remaining, nil
		}
		r.bucket = r.size - 1
		r.lastRefresh = time.Now()
	}
	r.mut.Unlock()
	return 0, nil
}

// Close will be called to cleanup
func (r *RandomRateLimit) Close(ctx context.Context) error {
	return nil
}
