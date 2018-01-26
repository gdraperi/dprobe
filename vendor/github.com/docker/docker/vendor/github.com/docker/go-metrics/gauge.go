package metrics

import "github.com/prometheus/client_golang/prometheus"

// Gauge is a metric that allows incrementing and decrementing a value
type Gauge interface ***REMOVED***
	Inc(...float64)
	Dec(...float64)

	// Add adds the provided value to the gauge's current value
	Add(float64)

	// Set replaces the gauge's current value with the provided value
	Set(float64)
***REMOVED***

// LabeledGauge describes a gauge the must have values populated before use.
type LabeledGauge interface ***REMOVED***
	WithValues(labels ...string) Gauge
***REMOVED***

type labeledGauge struct ***REMOVED***
	pg *prometheus.GaugeVec
***REMOVED***

func (lg *labeledGauge) WithValues(labels ...string) Gauge ***REMOVED***
	return &gauge***REMOVED***pg: lg.pg.WithLabelValues(labels...)***REMOVED***
***REMOVED***

func (lg *labeledGauge) Describe(c chan<- *prometheus.Desc) ***REMOVED***
	lg.pg.Describe(c)
***REMOVED***

func (lg *labeledGauge) Collect(c chan<- prometheus.Metric) ***REMOVED***
	lg.pg.Collect(c)
***REMOVED***

type gauge struct ***REMOVED***
	pg prometheus.Gauge
***REMOVED***

func (g *gauge) Inc(vs ...float64) ***REMOVED***
	if len(vs) == 0 ***REMOVED***
		g.pg.Inc()
	***REMOVED***

	g.Add(sumFloat64(vs...))
***REMOVED***

func (g *gauge) Dec(vs ...float64) ***REMOVED***
	if len(vs) == 0 ***REMOVED***
		g.pg.Dec()
	***REMOVED***

	g.Add(-sumFloat64(vs...))
***REMOVED***

func (g *gauge) Add(v float64) ***REMOVED***
	g.pg.Add(v)
***REMOVED***

func (g *gauge) Set(v float64) ***REMOVED***
	g.pg.Set(v)
***REMOVED***

func (g *gauge) Describe(c chan<- *prometheus.Desc) ***REMOVED***
	g.pg.Describe(c)
***REMOVED***

func (g *gauge) Collect(c chan<- prometheus.Metric) ***REMOVED***
	g.pg.Collect(c)
***REMOVED***
