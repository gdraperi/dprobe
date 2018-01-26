package metrics

// The MetricSink interface is used to transmit metrics information
// to an external system
type MetricSink interface ***REMOVED***
	// A Gauge should retain the last value it is set to
	SetGauge(key []string, val float32)

	// Should emit a Key/Value pair for each call
	EmitKey(key []string, val float32)

	// Counters should accumulate values
	IncrCounter(key []string, val float32)

	// Samples are for timing information, where quantiles are used
	AddSample(key []string, val float32)
***REMOVED***

// BlackholeSink is used to just blackhole messages
type BlackholeSink struct***REMOVED******REMOVED***

func (*BlackholeSink) SetGauge(key []string, val float32)    ***REMOVED******REMOVED***
func (*BlackholeSink) EmitKey(key []string, val float32)     ***REMOVED******REMOVED***
func (*BlackholeSink) IncrCounter(key []string, val float32) ***REMOVED******REMOVED***
func (*BlackholeSink) AddSample(key []string, val float32)   ***REMOVED******REMOVED***

// FanoutSink is used to sink to fanout values to multiple sinks
type FanoutSink []MetricSink

func (fh FanoutSink) SetGauge(key []string, val float32) ***REMOVED***
	for _, s := range fh ***REMOVED***
		s.SetGauge(key, val)
	***REMOVED***
***REMOVED***

func (fh FanoutSink) EmitKey(key []string, val float32) ***REMOVED***
	for _, s := range fh ***REMOVED***
		s.EmitKey(key, val)
	***REMOVED***
***REMOVED***

func (fh FanoutSink) IncrCounter(key []string, val float32) ***REMOVED***
	for _, s := range fh ***REMOVED***
		s.IncrCounter(key, val)
	***REMOVED***
***REMOVED***

func (fh FanoutSink) AddSample(key []string, val float32) ***REMOVED***
	for _, s := range fh ***REMOVED***
		s.AddSample(key, val)
	***REMOVED***
***REMOVED***
