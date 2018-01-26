package metrics

import (
	"runtime"
	"time"
)

func (m *Metrics) SetGauge(key []string, val float32) ***REMOVED***
	if m.HostName != "" && m.EnableHostname ***REMOVED***
		key = insert(0, m.HostName, key)
	***REMOVED***
	if m.EnableTypePrefix ***REMOVED***
		key = insert(0, "gauge", key)
	***REMOVED***
	if m.ServiceName != "" ***REMOVED***
		key = insert(0, m.ServiceName, key)
	***REMOVED***
	m.sink.SetGauge(key, val)
***REMOVED***

func (m *Metrics) EmitKey(key []string, val float32) ***REMOVED***
	if m.EnableTypePrefix ***REMOVED***
		key = insert(0, "kv", key)
	***REMOVED***
	if m.ServiceName != "" ***REMOVED***
		key = insert(0, m.ServiceName, key)
	***REMOVED***
	m.sink.EmitKey(key, val)
***REMOVED***

func (m *Metrics) IncrCounter(key []string, val float32) ***REMOVED***
	if m.EnableTypePrefix ***REMOVED***
		key = insert(0, "counter", key)
	***REMOVED***
	if m.ServiceName != "" ***REMOVED***
		key = insert(0, m.ServiceName, key)
	***REMOVED***
	m.sink.IncrCounter(key, val)
***REMOVED***

func (m *Metrics) AddSample(key []string, val float32) ***REMOVED***
	if m.EnableTypePrefix ***REMOVED***
		key = insert(0, "sample", key)
	***REMOVED***
	if m.ServiceName != "" ***REMOVED***
		key = insert(0, m.ServiceName, key)
	***REMOVED***
	m.sink.AddSample(key, val)
***REMOVED***

func (m *Metrics) MeasureSince(key []string, start time.Time) ***REMOVED***
	if m.EnableTypePrefix ***REMOVED***
		key = insert(0, "timer", key)
	***REMOVED***
	if m.ServiceName != "" ***REMOVED***
		key = insert(0, m.ServiceName, key)
	***REMOVED***
	now := time.Now()
	elapsed := now.Sub(start)
	msec := float32(elapsed.Nanoseconds()) / float32(m.TimerGranularity)
	m.sink.AddSample(key, msec)
***REMOVED***

// Periodically collects runtime stats to publish
func (m *Metrics) collectStats() ***REMOVED***
	for ***REMOVED***
		time.Sleep(m.ProfileInterval)
		m.emitRuntimeStats()
	***REMOVED***
***REMOVED***

// Emits various runtime statsitics
func (m *Metrics) emitRuntimeStats() ***REMOVED***
	// Export number of Goroutines
	numRoutines := runtime.NumGoroutine()
	m.SetGauge([]string***REMOVED***"runtime", "num_goroutines"***REMOVED***, float32(numRoutines))

	// Export memory stats
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	m.SetGauge([]string***REMOVED***"runtime", "alloc_bytes"***REMOVED***, float32(stats.Alloc))
	m.SetGauge([]string***REMOVED***"runtime", "sys_bytes"***REMOVED***, float32(stats.Sys))
	m.SetGauge([]string***REMOVED***"runtime", "malloc_count"***REMOVED***, float32(stats.Mallocs))
	m.SetGauge([]string***REMOVED***"runtime", "free_count"***REMOVED***, float32(stats.Frees))
	m.SetGauge([]string***REMOVED***"runtime", "heap_objects"***REMOVED***, float32(stats.HeapObjects))
	m.SetGauge([]string***REMOVED***"runtime", "total_gc_pause_ns"***REMOVED***, float32(stats.PauseTotalNs))
	m.SetGauge([]string***REMOVED***"runtime", "total_gc_runs"***REMOVED***, float32(stats.NumGC))

	// Export info about the last few GC runs
	num := stats.NumGC

	// Handle wrap around
	if num < m.lastNumGC ***REMOVED***
		m.lastNumGC = 0
	***REMOVED***

	// Ensure we don't scan more than 256
	if num-m.lastNumGC >= 256 ***REMOVED***
		m.lastNumGC = num - 255
	***REMOVED***

	for i := m.lastNumGC; i < num; i++ ***REMOVED***
		pause := stats.PauseNs[i%256]
		m.AddSample([]string***REMOVED***"runtime", "gc_pause_ns"***REMOVED***, float32(pause))
	***REMOVED***
	m.lastNumGC = num
***REMOVED***

// Inserts a string value at an index into the slice
func insert(i int, v string, s []string) []string ***REMOVED***
	s = append(s, "")
	copy(s[i+1:], s[i:])
	s[i] = v
	return s
***REMOVED***
