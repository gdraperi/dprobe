package metrics

import "github.com/prometheus/client_golang/prometheus"

// Counter is a metrics that can only increment its current count
type Counter interface ***REMOVED***
	// Inc adds Sum(vs) to the counter. Sum(vs) must be positive.
	//
	// If len(vs) == 0, increments the counter by 1.
	Inc(vs ...float64)
***REMOVED***

// LabeledCounter is counter that must have labels populated before use.
type LabeledCounter interface ***REMOVED***
	WithValues(vs ...string) Counter
***REMOVED***

type labeledCounter struct ***REMOVED***
	pc *prometheus.CounterVec
***REMOVED***

func (lc *labeledCounter) WithValues(vs ...string) Counter ***REMOVED***
	return &counter***REMOVED***pc: lc.pc.WithLabelValues(vs...)***REMOVED***
***REMOVED***

func (lc *labeledCounter) Describe(ch chan<- *prometheus.Desc) ***REMOVED***
	lc.pc.Describe(ch)
***REMOVED***

func (lc *labeledCounter) Collect(ch chan<- prometheus.Metric) ***REMOVED***
	lc.pc.Collect(ch)
***REMOVED***

type counter struct ***REMOVED***
	pc prometheus.Counter
***REMOVED***

func (c *counter) Inc(vs ...float64) ***REMOVED***
	if len(vs) == 0 ***REMOVED***
		c.pc.Inc()
	***REMOVED***

	c.pc.Add(sumFloat64(vs...))
***REMOVED***

func (c *counter) Describe(ch chan<- *prometheus.Desc) ***REMOVED***
	c.pc.Describe(ch)
***REMOVED***

func (c *counter) Collect(ch chan<- prometheus.Metric) ***REMOVED***
	c.pc.Collect(ch)
***REMOVED***
