// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheus

import "github.com/prometheus/procfs"

type processCollector struct ***REMOVED***
	pid             int
	collectFn       func(chan<- Metric)
	pidFn           func() (int, error)
	cpuTotal        Counter
	openFDs, maxFDs Gauge
	vsize, rss      Gauge
	startTime       Gauge
***REMOVED***

// NewProcessCollector returns a collector which exports the current state of
// process metrics including cpu, memory and file descriptor usage as well as
// the process start time for the given process id under the given namespace.
func NewProcessCollector(pid int, namespace string) *processCollector ***REMOVED***
	return NewProcessCollectorPIDFn(
		func() (int, error) ***REMOVED*** return pid, nil ***REMOVED***,
		namespace,
	)
***REMOVED***

// NewProcessCollectorPIDFn returns a collector which exports the current state
// of process metrics including cpu, memory and file descriptor usage as well
// as the process start time under the given namespace. The given pidFn is
// called on each collect and is used to determine the process to export
// metrics for.
func NewProcessCollectorPIDFn(
	pidFn func() (int, error),
	namespace string,
) *processCollector ***REMOVED***
	c := processCollector***REMOVED***
		pidFn:     pidFn,
		collectFn: func(chan<- Metric) ***REMOVED******REMOVED***,

		cpuTotal: NewCounter(CounterOpts***REMOVED***
			Namespace: namespace,
			Name:      "process_cpu_seconds_total",
			Help:      "Total user and system CPU time spent in seconds.",
		***REMOVED***),
		openFDs: NewGauge(GaugeOpts***REMOVED***
			Namespace: namespace,
			Name:      "process_open_fds",
			Help:      "Number of open file descriptors.",
		***REMOVED***),
		maxFDs: NewGauge(GaugeOpts***REMOVED***
			Namespace: namespace,
			Name:      "process_max_fds",
			Help:      "Maximum number of open file descriptors.",
		***REMOVED***),
		vsize: NewGauge(GaugeOpts***REMOVED***
			Namespace: namespace,
			Name:      "process_virtual_memory_bytes",
			Help:      "Virtual memory size in bytes.",
		***REMOVED***),
		rss: NewGauge(GaugeOpts***REMOVED***
			Namespace: namespace,
			Name:      "process_resident_memory_bytes",
			Help:      "Resident memory size in bytes.",
		***REMOVED***),
		startTime: NewGauge(GaugeOpts***REMOVED***
			Namespace: namespace,
			Name:      "process_start_time_seconds",
			Help:      "Start time of the process since unix epoch in seconds.",
		***REMOVED***),
	***REMOVED***

	// Set up process metric collection if supported by the runtime.
	if _, err := procfs.NewStat(); err == nil ***REMOVED***
		c.collectFn = c.processCollect
	***REMOVED***

	return &c
***REMOVED***

// Describe returns all descriptions of the collector.
func (c *processCollector) Describe(ch chan<- *Desc) ***REMOVED***
	ch <- c.cpuTotal.Desc()
	ch <- c.openFDs.Desc()
	ch <- c.maxFDs.Desc()
	ch <- c.vsize.Desc()
	ch <- c.rss.Desc()
	ch <- c.startTime.Desc()
***REMOVED***

// Collect returns the current state of all metrics of the collector.
func (c *processCollector) Collect(ch chan<- Metric) ***REMOVED***
	c.collectFn(ch)
***REMOVED***

// TODO(ts): Bring back error reporting by reverting 7faf9e7 as soon as the
// client allows users to configure the error behavior.
func (c *processCollector) processCollect(ch chan<- Metric) ***REMOVED***
	pid, err := c.pidFn()
	if err != nil ***REMOVED***
		return
	***REMOVED***

	p, err := procfs.NewProc(pid)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	if stat, err := p.NewStat(); err == nil ***REMOVED***
		c.cpuTotal.Set(stat.CPUTime())
		ch <- c.cpuTotal
		c.vsize.Set(float64(stat.VirtualMemory()))
		ch <- c.vsize
		c.rss.Set(float64(stat.ResidentMemory()))
		ch <- c.rss

		if startTime, err := stat.StartTime(); err == nil ***REMOVED***
			c.startTime.Set(startTime)
			ch <- c.startTime
		***REMOVED***
	***REMOVED***

	if fds, err := p.FileDescriptorsLen(); err == nil ***REMOVED***
		c.openFDs.Set(float64(fds))
		ch <- c.openFDs
	***REMOVED***

	if limits, err := p.NewLimits(); err == nil ***REMOVED***
		c.maxFDs.Set(float64(limits.OpenFiles))
		ch <- c.maxFDs
	***REMOVED***
***REMOVED***
