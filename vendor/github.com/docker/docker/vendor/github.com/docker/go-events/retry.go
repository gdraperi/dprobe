package events

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

// RetryingSink retries the write until success or an ErrSinkClosed is
// returned. Underlying sink must have p > 0 of succeeding or the sink will
// block. Retry is configured with a RetryStrategy.  Concurrent calls to a
// retrying sink are serialized through the sink, meaning that if one is
// in-flight, another will not proceed.
type RetryingSink struct ***REMOVED***
	sink     Sink
	strategy RetryStrategy
	closed   chan struct***REMOVED******REMOVED***
	once     sync.Once
***REMOVED***

// NewRetryingSink returns a sink that will retry writes to a sink, backing
// off on failure. Parameters threshold and backoff adjust the behavior of the
// circuit breaker.
func NewRetryingSink(sink Sink, strategy RetryStrategy) *RetryingSink ***REMOVED***
	rs := &RetryingSink***REMOVED***
		sink:     sink,
		strategy: strategy,
		closed:   make(chan struct***REMOVED******REMOVED***),
	***REMOVED***

	return rs
***REMOVED***

// Write attempts to flush the events to the downstream sink until it succeeds
// or the sink is closed.
func (rs *RetryingSink) Write(event Event) error ***REMOVED***
	logger := logrus.WithField("event", event)

retry:
	select ***REMOVED***
	case <-rs.closed:
		return ErrSinkClosed
	default:
	***REMOVED***

	if backoff := rs.strategy.Proceed(event); backoff > 0 ***REMOVED***
		select ***REMOVED***
		case <-time.After(backoff):
			// TODO(stevvooe): This branch holds up the next try. Before, we
			// would simply break to the "retry" label and then possibly wait
			// again. However, this requires all retry strategies to have a
			// large probability of probing the sync for success, rather than
			// just backing off and sending the request.
		case <-rs.closed:
			return ErrSinkClosed
		***REMOVED***
	***REMOVED***

	if err := rs.sink.Write(event); err != nil ***REMOVED***
		if err == ErrSinkClosed ***REMOVED***
			// terminal!
			return err
		***REMOVED***

		logger := logger.WithError(err) // shadow!!

		if rs.strategy.Failure(event, err) ***REMOVED***
			logger.Errorf("retryingsink: dropped event")
			return nil
		***REMOVED***

		logger.Errorf("retryingsink: error writing event, retrying")
		goto retry
	***REMOVED***

	rs.strategy.Success(event)
	return nil
***REMOVED***

// Close closes the sink and the underlying sink.
func (rs *RetryingSink) Close() error ***REMOVED***
	rs.once.Do(func() ***REMOVED***
		close(rs.closed)
	***REMOVED***)

	return nil
***REMOVED***

func (rs *RetryingSink) String() string ***REMOVED***
	// Serialize a copy of the RetryingSink without the sync.Once, to avoid
	// a data race.
	rs2 := map[string]interface***REMOVED******REMOVED******REMOVED***
		"sink":     rs.sink,
		"strategy": rs.strategy,
		"closed":   rs.closed,
	***REMOVED***
	return fmt.Sprint(rs2)
***REMOVED***

// RetryStrategy defines a strategy for retrying event sink writes.
//
// All methods should be goroutine safe.
type RetryStrategy interface ***REMOVED***
	// Proceed is called before every event send. If proceed returns a
	// positive, non-zero integer, the retryer will back off by the provided
	// duration.
	//
	// An event is provided, by may be ignored.
	Proceed(event Event) time.Duration

	// Failure reports a failure to the strategy. If this method returns true,
	// the event should be dropped.
	Failure(event Event, err error) bool

	// Success should be called when an event is sent successfully.
	Success(event Event)
***REMOVED***

// Breaker implements a circuit breaker retry strategy.
//
// The current implementation never drops events.
type Breaker struct ***REMOVED***
	threshold int
	recent    int
	last      time.Time
	backoff   time.Duration // time after which we retry after failure.
	mu        sync.Mutex
***REMOVED***

var _ RetryStrategy = &Breaker***REMOVED******REMOVED***

// NewBreaker returns a breaker that will backoff after the threshold has been
// tripped. A Breaker is thread safe and may be shared by many goroutines.
func NewBreaker(threshold int, backoff time.Duration) *Breaker ***REMOVED***
	return &Breaker***REMOVED***
		threshold: threshold,
		backoff:   backoff,
	***REMOVED***
***REMOVED***

// Proceed checks the failures against the threshold.
func (b *Breaker) Proceed(event Event) time.Duration ***REMOVED***
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.recent < b.threshold ***REMOVED***
		return 0
	***REMOVED***

	return b.last.Add(b.backoff).Sub(time.Now())
***REMOVED***

// Success resets the breaker.
func (b *Breaker) Success(event Event) ***REMOVED***
	b.mu.Lock()
	defer b.mu.Unlock()

	b.recent = 0
	b.last = time.Time***REMOVED******REMOVED***
***REMOVED***

// Failure records the failure and latest failure time.
func (b *Breaker) Failure(event Event, err error) bool ***REMOVED***
	b.mu.Lock()
	defer b.mu.Unlock()

	b.recent++
	b.last = time.Now().UTC()
	return false // never drop events.
***REMOVED***

var (
	// DefaultExponentialBackoffConfig provides a default configuration for
	// exponential backoff.
	DefaultExponentialBackoffConfig = ExponentialBackoffConfig***REMOVED***
		Base:   time.Second,
		Factor: time.Second,
		Max:    20 * time.Second,
	***REMOVED***
)

// ExponentialBackoffConfig configures backoff parameters.
//
// Note that these parameters operate on the upper bound for choosing a random
// value. For example, at Base=1s, a random value in [0,1s) will be chosen for
// the backoff value.
type ExponentialBackoffConfig struct ***REMOVED***
	// Base is the minimum bound for backing off after failure.
	Base time.Duration

	// Factor sets the amount of time by which the backoff grows with each
	// failure.
	Factor time.Duration

	// Max is the absolute maxiumum bound for a single backoff.
	Max time.Duration
***REMOVED***

// ExponentialBackoff implements random backoff with exponentially increasing
// bounds as the number consecutive failures increase.
type ExponentialBackoff struct ***REMOVED***
	config   ExponentialBackoffConfig
	failures uint64 // consecutive failure counter.
***REMOVED***

// NewExponentialBackoff returns an exponential backoff strategy with the
// desired config. If config is nil, the default is returned.
func NewExponentialBackoff(config ExponentialBackoffConfig) *ExponentialBackoff ***REMOVED***
	return &ExponentialBackoff***REMOVED***
		config: config,
	***REMOVED***
***REMOVED***

// Proceed returns the next randomly bound exponential backoff time.
func (b *ExponentialBackoff) Proceed(event Event) time.Duration ***REMOVED***
	return b.backoff(atomic.LoadUint64(&b.failures))
***REMOVED***

// Success resets the failures counter.
func (b *ExponentialBackoff) Success(event Event) ***REMOVED***
	atomic.StoreUint64(&b.failures, 0)
***REMOVED***

// Failure increments the failure counter.
func (b *ExponentialBackoff) Failure(event Event, err error) bool ***REMOVED***
	atomic.AddUint64(&b.failures, 1)
	return false
***REMOVED***

// backoff calculates the amount of time to wait based on the number of
// consecutive failures.
func (b *ExponentialBackoff) backoff(failures uint64) time.Duration ***REMOVED***
	if failures <= 0 ***REMOVED***
		// proceed normally when there are no failures.
		return 0
	***REMOVED***

	factor := b.config.Factor
	if factor <= 0 ***REMOVED***
		factor = DefaultExponentialBackoffConfig.Factor
	***REMOVED***

	backoff := b.config.Base + factor*time.Duration(1<<(failures-1))

	max := b.config.Max
	if max <= 0 ***REMOVED***
		max = DefaultExponentialBackoffConfig.Max
	***REMOVED***

	if backoff > max || backoff < 0 ***REMOVED***
		backoff = max
	***REMOVED***

	// Choose a uniformly distributed value from [0, backoff).
	return time.Duration(rand.Int63n(int64(backoff)))
***REMOVED***
