// Copyright 2014 The Prometheus Authors
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

// Gauge is a Metric that represents a single numerical value that can
// arbitrarily go up and down.
//
// A Gauge is typically used for measured values like temperatures or current
// memory usage, but also "counts" that can go up and down, like the number of
// running goroutines.
//
// To create Gauge instances, use NewGauge.
type Gauge interface ***REMOVED***
	Metric
	Collector

	// Set sets the Gauge to an arbitrary value.
	Set(float64)
	// Inc increments the Gauge by 1.
	Inc()
	// Dec decrements the Gauge by 1.
	Dec()
	// Add adds the given value to the Gauge. (The value can be
	// negative, resulting in a decrease of the Gauge.)
	Add(float64)
	// Sub subtracts the given value from the Gauge. (The value can be
	// negative, resulting in an increase of the Gauge.)
	Sub(float64)
***REMOVED***

// GaugeOpts is an alias for Opts. See there for doc comments.
type GaugeOpts Opts

// NewGauge creates a new Gauge based on the provided GaugeOpts.
func NewGauge(opts GaugeOpts) Gauge ***REMOVED***
	return newValue(NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		nil,
		opts.ConstLabels,
	), GaugeValue, 0)
***REMOVED***

// GaugeVec is a Collector that bundles a set of Gauges that all share the same
// Desc, but have different values for their variable labels. This is used if
// you want to count the same thing partitioned by various dimensions
// (e.g. number of operations queued, partitioned by user and operation
// type). Create instances with NewGaugeVec.
type GaugeVec struct ***REMOVED***
	MetricVec
***REMOVED***

// NewGaugeVec creates a new GaugeVec based on the provided GaugeOpts and
// partitioned by the given label names. At least one label name must be
// provided.
func NewGaugeVec(opts GaugeOpts, labelNames []string) *GaugeVec ***REMOVED***
	desc := NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		labelNames,
		opts.ConstLabels,
	)
	return &GaugeVec***REMOVED***
		MetricVec: MetricVec***REMOVED***
			children: map[uint64]Metric***REMOVED******REMOVED***,
			desc:     desc,
			newMetric: func(lvs ...string) Metric ***REMOVED***
				return newValue(desc, GaugeValue, 0, lvs...)
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// GetMetricWithLabelValues replaces the method of the same name in
// MetricVec. The difference is that this method returns a Gauge and not a
// Metric so that no type conversion is required.
func (m *GaugeVec) GetMetricWithLabelValues(lvs ...string) (Gauge, error) ***REMOVED***
	metric, err := m.MetricVec.GetMetricWithLabelValues(lvs...)
	if metric != nil ***REMOVED***
		return metric.(Gauge), err
	***REMOVED***
	return nil, err
***REMOVED***

// GetMetricWith replaces the method of the same name in MetricVec. The
// difference is that this method returns a Gauge and not a Metric so that no
// type conversion is required.
func (m *GaugeVec) GetMetricWith(labels Labels) (Gauge, error) ***REMOVED***
	metric, err := m.MetricVec.GetMetricWith(labels)
	if metric != nil ***REMOVED***
		return metric.(Gauge), err
	***REMOVED***
	return nil, err
***REMOVED***

// WithLabelValues works as GetMetricWithLabelValues, but panics where
// GetMetricWithLabelValues would have returned an error. By not returning an
// error, WithLabelValues allows shortcuts like
//     myVec.WithLabelValues("404", "GET").Add(42)
func (m *GaugeVec) WithLabelValues(lvs ...string) Gauge ***REMOVED***
	return m.MetricVec.WithLabelValues(lvs...).(Gauge)
***REMOVED***

// With works as GetMetricWith, but panics where GetMetricWithLabels would have
// returned an error. By not returning an error, With allows shortcuts like
//     myVec.With(Labels***REMOVED***"code": "404", "method": "GET"***REMOVED***).Add(42)
func (m *GaugeVec) With(labels Labels) Gauge ***REMOVED***
	return m.MetricVec.With(labels).(Gauge)
***REMOVED***

// GaugeFunc is a Gauge whose value is determined at collect time by calling a
// provided function.
//
// To create GaugeFunc instances, use NewGaugeFunc.
type GaugeFunc interface ***REMOVED***
	Metric
	Collector
***REMOVED***

// NewGaugeFunc creates a new GaugeFunc based on the provided GaugeOpts. The
// value reported is determined by calling the given function from within the
// Write method. Take into account that metric collection may happen
// concurrently. If that results in concurrent calls to Write, like in the case
// where a GaugeFunc is directly registered with Prometheus, the provided
// function must be concurrency-safe.
func NewGaugeFunc(opts GaugeOpts, function func() float64) GaugeFunc ***REMOVED***
	return newValueFunc(NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		nil,
		opts.ConstLabels,
	), GaugeValue, function)
***REMOVED***
