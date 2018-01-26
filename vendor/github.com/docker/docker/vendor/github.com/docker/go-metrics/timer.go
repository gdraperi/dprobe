package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// StartTimer begins a timer observation at the callsite. When the target
// operation is completed, the caller should call the return done func().
func StartTimer(timer Timer) (done func()) ***REMOVED***
	start := time.Now()
	return func() ***REMOVED***
		timer.Update(time.Since(start))
	***REMOVED***
***REMOVED***

// Timer is a metric that allows collecting the duration of an action in seconds
type Timer interface ***REMOVED***
	// Update records an observation, duration, and converts to the target
	// units.
	Update(duration time.Duration)

	// UpdateSince will add the duration from the provided starting time to the
	// timer's summary with the precisions that was used in creation of the timer
	UpdateSince(time.Time)
***REMOVED***

// LabeledTimer is a timer that must have label values populated before use.
type LabeledTimer interface ***REMOVED***
	WithValues(labels ...string) Timer
***REMOVED***

type labeledTimer struct ***REMOVED***
	m *prometheus.HistogramVec
***REMOVED***

func (lt *labeledTimer) WithValues(labels ...string) Timer ***REMOVED***
	return &timer***REMOVED***m: lt.m.WithLabelValues(labels...)***REMOVED***
***REMOVED***

func (lt *labeledTimer) Describe(c chan<- *prometheus.Desc) ***REMOVED***
	lt.m.Describe(c)
***REMOVED***

func (lt *labeledTimer) Collect(c chan<- prometheus.Metric) ***REMOVED***
	lt.m.Collect(c)
***REMOVED***

type timer struct ***REMOVED***
	m prometheus.Histogram
***REMOVED***

func (t *timer) Update(duration time.Duration) ***REMOVED***
	t.m.Observe(duration.Seconds())
***REMOVED***

func (t *timer) UpdateSince(since time.Time) ***REMOVED***
	t.m.Observe(time.Since(since).Seconds())
***REMOVED***

func (t *timer) Describe(c chan<- *prometheus.Desc) ***REMOVED***
	t.m.Describe(c)
***REMOVED***

func (t *timer) Collect(c chan<- prometheus.Metric) ***REMOVED***
	t.m.Collect(c)
***REMOVED***
