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

import (
	"strings"

	dto "github.com/prometheus/client_model/go"
)

const separatorByte byte = 255

// A Metric models a single sample value with its meta data being exported to
// Prometheus. Implementers of Metric in this package inclued Gauge, Counter,
// Untyped, and Summary. Users can implement their own Metric types, but that
// should be rarely needed. See the example for SelfCollector, which is also an
// example for a user-implemented Metric.
type Metric interface ***REMOVED***
	// Desc returns the descriptor for the Metric. This method idempotently
	// returns the same descriptor throughout the lifetime of the
	// Metric. The returned descriptor is immutable by contract. A Metric
	// unable to describe itself must return an invalid descriptor (created
	// with NewInvalidDesc).
	Desc() *Desc
	// Write encodes the Metric into a "Metric" Protocol Buffer data
	// transmission object.
	//
	// Implementers of custom Metric types must observe concurrency safety
	// as reads of this metric may occur at any time, and any blocking
	// occurs at the expense of total performance of rendering all
	// registered metrics. Ideally Metric implementations should support
	// concurrent readers.
	//
	// The Prometheus client library attempts to minimize memory allocations
	// and will provide a pre-existing reset dto.Metric pointer. Prometheus
	// may recycle the dto.Metric proto message, so Metric implementations
	// should just populate the provided dto.Metric and then should not keep
	// any reference to it.
	//
	// While populating dto.Metric, labels must be sorted lexicographically.
	// (Implementers may find LabelPairSorter useful for that.)
	Write(*dto.Metric) error
***REMOVED***

// Opts bundles the options for creating most Metric types. Each metric
// implementation XXX has its own XXXOpts type, but in most cases, it is just be
// an alias of this type (which might change when the requirement arises.)
//
// It is mandatory to set Name and Help to a non-empty string. All other fields
// are optional and can safely be left at their zero value.
type Opts struct ***REMOVED***
	// Namespace, Subsystem, and Name are components of the fully-qualified
	// name of the Metric (created by joining these components with
	// "_"). Only Name is mandatory, the others merely help structuring the
	// name. Note that the fully-qualified name of the metric must be a
	// valid Prometheus metric name.
	Namespace string
	Subsystem string
	Name      string

	// Help provides information about this metric. Mandatory!
	//
	// Metrics with the same fully-qualified name must have the same Help
	// string.
	Help string

	// ConstLabels are used to attach fixed labels to this metric. Metrics
	// with the same fully-qualified name must have the same label names in
	// their ConstLabels.
	//
	// Note that in most cases, labels have a value that varies during the
	// lifetime of a process. Those labels are usually managed with a metric
	// vector collector (like CounterVec, GaugeVec, UntypedVec). ConstLabels
	// serve only special purposes. One is for the special case where the
	// value of a label does not change during the lifetime of a process,
	// e.g. if the revision of the running binary is put into a
	// label. Another, more advanced purpose is if more than one Collector
	// needs to collect Metrics with the same fully-qualified name. In that
	// case, those Metrics must differ in the values of their
	// ConstLabels. See the Collector examples.
	//
	// If the value of a label never changes (not even between binaries),
	// that label most likely should not be a label at all (but part of the
	// metric name).
	ConstLabels Labels
***REMOVED***

// BuildFQName joins the given three name components by "_". Empty name
// components are ignored. If the name parameter itself is empty, an empty
// string is returned, no matter what. Metric implementations included in this
// library use this function internally to generate the fully-qualified metric
// name from the name component in their Opts. Users of the library will only
// need this function if they implement their own Metric or instantiate a Desc
// (with NewDesc) directly.
func BuildFQName(namespace, subsystem, name string) string ***REMOVED***
	if name == "" ***REMOVED***
		return ""
	***REMOVED***
	switch ***REMOVED***
	case namespace != "" && subsystem != "":
		return strings.Join([]string***REMOVED***namespace, subsystem, name***REMOVED***, "_")
	case namespace != "":
		return strings.Join([]string***REMOVED***namespace, name***REMOVED***, "_")
	case subsystem != "":
		return strings.Join([]string***REMOVED***subsystem, name***REMOVED***, "_")
	***REMOVED***
	return name
***REMOVED***

// LabelPairSorter implements sort.Interface. It is used to sort a slice of
// dto.LabelPair pointers. This is useful for implementing the Write method of
// custom metrics.
type LabelPairSorter []*dto.LabelPair

func (s LabelPairSorter) Len() int ***REMOVED***
	return len(s)
***REMOVED***

func (s LabelPairSorter) Swap(i, j int) ***REMOVED***
	s[i], s[j] = s[j], s[i]
***REMOVED***

func (s LabelPairSorter) Less(i, j int) bool ***REMOVED***
	return s[i].GetName() < s[j].GetName()
***REMOVED***

type hashSorter []uint64

func (s hashSorter) Len() int ***REMOVED***
	return len(s)
***REMOVED***

func (s hashSorter) Swap(i, j int) ***REMOVED***
	s[i], s[j] = s[j], s[i]
***REMOVED***

func (s hashSorter) Less(i, j int) bool ***REMOVED***
	return s[i] < s[j]
***REMOVED***

type invalidMetric struct ***REMOVED***
	desc *Desc
	err  error
***REMOVED***

// NewInvalidMetric returns a metric whose Write method always returns the
// provided error. It is useful if a Collector finds itself unable to collect
// a metric and wishes to report an error to the registry.
func NewInvalidMetric(desc *Desc, err error) Metric ***REMOVED***
	return &invalidMetric***REMOVED***desc, err***REMOVED***
***REMOVED***

func (m *invalidMetric) Desc() *Desc ***REMOVED*** return m.desc ***REMOVED***

func (m *invalidMetric) Write(*dto.Metric) error ***REMOVED*** return m.err ***REMOVED***
