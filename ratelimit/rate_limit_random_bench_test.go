package ratelimit

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func BenchmarkRateLimit(b *testing.B) {
	/* A rate limit is typically going to be protecting a networked resource
	 * where the request will likely be measured at least in hundreds of
	 * microseconds. It would be reasonable to assume the rate limit might be
	 * shared across tens of components.
	 *
	 * Therefore, we can probably sit comfortably with lock contention across
	 * one hundred or so parallel components adding an overhead of single digit
	 * microseconds. Since this benchmark doesn't take into account the actual
	 * request duration after receiving a rate limit I've set the number of
	 * components to ten in order to compensate.
	 */
	b.ReportAllocs()

	nParallel := 10
	startChan := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(nParallel)

	conf, err := randomRatelimitConfig().ParseYAML(`
count: 1
min_interval: 250ms
max_interval: 750ms
`, nil)
	require.NoError(b, err)

	rl, err := NewRandomRateLimitFromConfig(conf)
	require.NoError(b, err)

	ctx := context.Background()

	for i := 0; i < nParallel; i++ {
		go func() {
			<-startChan
			for j := 0; j < b.N; j++ {
				period, _ := rl.Access(ctx)
				if period > 0 {
					time.Sleep(period)
				}
			}
			wg.Done()
		}()
	}

	b.ResetTimer()
	close(startChan)
	wg.Wait()
}
