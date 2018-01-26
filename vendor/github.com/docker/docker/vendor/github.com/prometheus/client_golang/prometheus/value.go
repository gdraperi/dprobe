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
	"errors"
	"fmt"
	"math"
	"sort"
	"sync/atomic"

	dto "github.com/prometheus/client_model/go"

	"github.com/golang/protobuf/proto"
)

// ValueType is an enumeration of metric types that represent a simple value.
type ValueType int

// Possible values for the ValueType enum.
const (
	_ ValueType = iota
	CounterValue
	GaugeValue
	UntypedValue
)

var errInconsistentCardinality = errors.New("inconsistent label cardinality")

// value is a generic metric for simple values. It implements Metric, Collector,
// Counter, Gauge, and Untyped. Its effective type is determined by
// ValueType. This is a low-level building block used by the library to back the
// implementations of Counter, Gauge, and Untyped.
type value struct ***REMOVED***
	// valBits containst the bits of the represented float64 value. It has
	// to go first in the struct to guarantee alignment for atomic
	// operations.  http://golang.org/pkg/sync/atomic/#pkg-note-BUG
	valBits uint64

	SelfCollector

	desc       *Desc
	valType    ValueType
	labelPairs []*dto.LabelPair
***REMOVED***

// newValue returns a newly allocated value with the given Desc, ValueType,
// sample value and label values. It panics if the number of label
// values is different from the number of variable labels in Desc.
func newValue(desc *Desc, valueType ValueType, val float64, labelValues ...string) *value ***REMOVED***
	if len(labelValues) != len(desc.variableLabels) ***REMOVED***
		panic(errInconsistentCardinality)
	***REMOVED***
	result := &value***REMOVED***
		desc:       desc,
		valType:    valueType,
		valBits:    math.Float64bits(val),
		labelPairs: makeLabelPairs(desc, labelValues),
	***REMOVED***
	result.Init(result)
	return result
***REMOVED***

func (v *value) Desc() *Desc ***REMOVED***
	return v.desc
***REMOVED***

func (v *value) Set(val float64) ***REMOVED***
	atomic.StoreUint64(&v.valBits, math.Float64bits(val))
***REMOVED***

func (v *value) Inc() ***REMOVED***
	v.Add(1)
***REMOVED***

func (v *value) Dec() ***REMOVED***
	v.Add(-1)
***REMOVED***

func (v *value) Add(val float64) ***REMOVED***
	for ***REMOVED***
		oldBits := atomic.LoadUint64(&v.valBits)
		newBits := math.Float64bits(math.Float64frombits(oldBits) + val)
		if atomic.CompareAndSwapUint64(&v.valBits, oldBits, newBits) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (v *value) Sub(val float64) ***REMOVED***
	v.Add(val * -1)
***REMOVED***

func (v *value) Write(out *dto.Metric) error ***REMOVED***
	val := math.Float64frombits(atomic.LoadUint64(&v.valBits))
	return populateMetric(v.valType, val, v.labelPairs, out)
***REMOVED***

// valueFunc is a generic metric for simple values retrieved on collect time
// from a function. It implements Metric and Collector. Its effective type is
// determined by ValueType. This is a low-level building block used by the
// library to back the implementations of CounterFunc, GaugeFunc, and
// UntypedFunc.
type valueFunc struct ***REMOVED***
	SelfCollector

	desc       *Desc
	valType    ValueType
	function   func() float64
	labelPairs []*dto.LabelPair
***REMOVED***

// newValueFunc returns a newly allocated valueFunc with the given Desc and
// ValueType. The value reported is determined by calling the given function
// from within the Write method. Take into account that metric collection may
// happen concurrently. If that results in concurrent calls to Write, like in
// the case where a valueFunc is directly registered with Prometheus, the
// provided function must be concurrency-safe.
func newValueFunc(desc *Desc, valueType ValueType, function func() float64) *valueFunc ***REMOVED***
	result := &valueFunc***REMOVED***
		desc:       desc,
		valType:    valueType,
		function:   function,
		labelPairs: makeLabelPairs(desc, nil),
	***REMOVED***
	result.Init(result)
	return result
***REMOVED***

func (v *valueFunc) Desc() *Desc ***REMOVED***
	return v.desc
***REMOVED***

func (v *valueFunc) Write(out *dto.Metric) error ***REMOVED***
	return populateMetric(v.valType, v.function(), v.labelPairs, out)
***REMOVED***

// NewConstMetric returns a metric with one fixed value that cannot be
// changed. Users of this package will not have much use for it in regular
// operations. However, when implementing custom Collectors, it is useful as a
// throw-away metric that is generated on the fly to send it to Prometheus in
// the Collect method. NewConstMetric returns an error if the length of
// labelValues is not consistent with the variable labels in Desc.
func NewConstMetric(desc *Desc, valueType ValueType, value float64, labelValues ...string) (Metric, error) ***REMOVED***
	if len(desc.variableLabels) != len(labelValues) ***REMOVED***
		return nil, errInconsistentCardinality
	***REMOVED***
	return &constMetric***REMOVED***
		desc:       desc,
		valType:    valueType,
		val:        value,
		labelPairs: makeLabelPairs(desc, labelValues),
	***REMOVED***, nil
***REMOVED***

// MustNewConstMetric is a version of NewConstMetric that panics where
// NewConstMetric would have returned an error.
func MustNewConstMetric(desc *Desc, valueType ValueType, value float64, labelValues ...string) Metric ***REMOVED***
	m, err := NewConstMetric(desc, valueType, value, labelValues...)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return m
***REMOVED***

type constMetric struct ***REMOVED***
	desc       *Desc
	valType    ValueType
	val        float64
	labelPairs []*dto.LabelPair
***REMOVED***

func (m *constMetric) Desc() *Desc ***REMOVED***
	return m.desc
***REMOVED***

func (m *constMetric) Write(out *dto.Metric) error ***REMOVED***
	return populateMetric(m.valType, m.val, m.labelPairs, out)
***REMOVED***

func populateMetric(
	t ValueType,
	v float64,
	labelPairs []*dto.LabelPair,
	m *dto.Metric,
) error ***REMOVED***
	m.Label = labelPairs
	switch t ***REMOVED***
	case CounterValue:
		m.Counter = &dto.Counter***REMOVED***Value: proto.Float64(v)***REMOVED***
	case GaugeValue:
		m.Gauge = &dto.Gauge***REMOVED***Value: proto.Float64(v)***REMOVED***
	case UntypedValue:
		m.Untyped = &dto.Untyped***REMOVED***Value: proto.Float64(v)***REMOVED***
	default:
		return fmt.Errorf("encountered unknown type %v", t)
	***REMOVED***
	return nil
***REMOVED***

func makeLabelPairs(desc *Desc, labelValues []string) []*dto.LabelPair ***REMOVED***
	totalLen := len(desc.variableLabels) + len(desc.constLabelPairs)
	if totalLen == 0 ***REMOVED***
		// Super fast path.
		return nil
	***REMOVED***
	if len(desc.variableLabels) == 0 ***REMOVED***
		// Moderately fast path.
		return desc.constLabelPairs
	***REMOVED***
	labelPairs := make([]*dto.LabelPair, 0, totalLen)
	for i, n := range desc.variableLabels ***REMOVED***
		labelPairs = append(labelPairs, &dto.LabelPair***REMOVED***
			Name:  proto.String(n),
			Value: proto.String(labelValues[i]),
		***REMOVED***)
	***REMOVED***
	for _, lp := range desc.constLabelPairs ***REMOVED***
		labelPairs = append(labelPairs, lp)
	***REMOVED***
	sort.Sort(LabelPairSorter(labelPairs))
	return labelPairs
***REMOVED***
