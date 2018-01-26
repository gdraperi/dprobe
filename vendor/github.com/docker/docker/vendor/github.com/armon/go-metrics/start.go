package metrics

import (
	"os"
	"time"
)

// Config is used to configure metrics settings
type Config struct ***REMOVED***
	ServiceName          string        // Prefixed with keys to seperate services
	HostName             string        // Hostname to use. If not provided and EnableHostname, it will be os.Hostname
	EnableHostname       bool          // Enable prefixing gauge values with hostname
	EnableRuntimeMetrics bool          // Enables profiling of runtime metrics (GC, Goroutines, Memory)
	EnableTypePrefix     bool          // Prefixes key with a type ("counter", "gauge", "timer")
	TimerGranularity     time.Duration // Granularity of timers.
	ProfileInterval      time.Duration // Interval to profile runtime metrics
***REMOVED***

// Metrics represents an instance of a metrics sink that can
// be used to emit
type Metrics struct ***REMOVED***
	Config
	lastNumGC uint32
	sink      MetricSink
***REMOVED***

// Shared global metrics instance
var globalMetrics *Metrics

func init() ***REMOVED***
	// Initialize to a blackhole sink to avoid errors
	globalMetrics = &Metrics***REMOVED***sink: &BlackholeSink***REMOVED******REMOVED******REMOVED***
***REMOVED***

// DefaultConfig provides a sane default configuration
func DefaultConfig(serviceName string) *Config ***REMOVED***
	c := &Config***REMOVED***
		ServiceName:          serviceName, // Use client provided service
		HostName:             "",
		EnableHostname:       true,             // Enable hostname prefix
		EnableRuntimeMetrics: true,             // Enable runtime profiling
		EnableTypePrefix:     false,            // Disable type prefix
		TimerGranularity:     time.Millisecond, // Timers are in milliseconds
		ProfileInterval:      time.Second,      // Poll runtime every second
	***REMOVED***

	// Try to get the hostname
	name, _ := os.Hostname()
	c.HostName = name
	return c
***REMOVED***

// New is used to create a new instance of Metrics
func New(conf *Config, sink MetricSink) (*Metrics, error) ***REMOVED***
	met := &Metrics***REMOVED******REMOVED***
	met.Config = *conf
	met.sink = sink

	// Start the runtime collector
	if conf.EnableRuntimeMetrics ***REMOVED***
		go met.collectStats()
	***REMOVED***
	return met, nil
***REMOVED***

// NewGlobal is the same as New, but it assigns the metrics object to be
// used globally as well as returning it.
func NewGlobal(conf *Config, sink MetricSink) (*Metrics, error) ***REMOVED***
	metrics, err := New(conf, sink)
	if err == nil ***REMOVED***
		globalMetrics = metrics
	***REMOVED***
	return metrics, err
***REMOVED***

// Proxy all the methods to the globalMetrics instance
func SetGauge(key []string, val float32) ***REMOVED***
	globalMetrics.SetGauge(key, val)
***REMOVED***

func EmitKey(key []string, val float32) ***REMOVED***
	globalMetrics.EmitKey(key, val)
***REMOVED***

func IncrCounter(key []string, val float32) ***REMOVED***
	globalMetrics.IncrCounter(key, val)
***REMOVED***

func AddSample(key []string, val float32) ***REMOVED***
	globalMetrics.AddSample(key, val)
***REMOVED***

func MeasureSince(key []string, start time.Time) ***REMOVED***
	globalMetrics.MeasureSince(key, start)
***REMOVED***
