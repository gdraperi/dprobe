// Copyright 2013 The Prometheus Authors
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

package model

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

// A SampleValue is a representation of a value for a given sample at a given
// time.
type SampleValue float64

// MarshalJSON implements json.Marshaler.
func (v SampleValue) MarshalJSON() ([]byte, error) ***REMOVED***
	return json.Marshal(v.String())
***REMOVED***

// UnmarshalJSON implements json.Unmarshaler.
func (v *SampleValue) UnmarshalJSON(b []byte) error ***REMOVED***
	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' ***REMOVED***
		return fmt.Errorf("sample value must be a quoted string")
	***REMOVED***
	f, err := strconv.ParseFloat(string(b[1:len(b)-1]), 64)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*v = SampleValue(f)
	return nil
***REMOVED***

// Equal returns true if the value of v and o is equal or if both are NaN. Note
// that v==o is false if both are NaN. If you want the conventional float
// behavior, use == to compare two SampleValues.
func (v SampleValue) Equal(o SampleValue) bool ***REMOVED***
	if v == o ***REMOVED***
		return true
	***REMOVED***
	return math.IsNaN(float64(v)) && math.IsNaN(float64(o))
***REMOVED***

func (v SampleValue) String() string ***REMOVED***
	return strconv.FormatFloat(float64(v), 'f', -1, 64)
***REMOVED***

// SamplePair pairs a SampleValue with a Timestamp.
type SamplePair struct ***REMOVED***
	Timestamp Time
	Value     SampleValue
***REMOVED***

// MarshalJSON implements json.Marshaler.
func (s SamplePair) MarshalJSON() ([]byte, error) ***REMOVED***
	t, err := json.Marshal(s.Timestamp)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	v, err := json.Marshal(s.Value)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return []byte(fmt.Sprintf("[%s,%s]", t, v)), nil
***REMOVED***

// UnmarshalJSON implements json.Unmarshaler.
func (s *SamplePair) UnmarshalJSON(b []byte) error ***REMOVED***
	v := [...]json.Unmarshaler***REMOVED***&s.Timestamp, &s.Value***REMOVED***
	return json.Unmarshal(b, &v)
***REMOVED***

// Equal returns true if this SamplePair and o have equal Values and equal
// Timestamps. The sematics of Value equality is defined by SampleValue.Equal.
func (s *SamplePair) Equal(o *SamplePair) bool ***REMOVED***
	return s == o || (s.Value.Equal(o.Value) && s.Timestamp.Equal(o.Timestamp))
***REMOVED***

func (s SamplePair) String() string ***REMOVED***
	return fmt.Sprintf("%s @[%s]", s.Value, s.Timestamp)
***REMOVED***

// Sample is a sample pair associated with a metric.
type Sample struct ***REMOVED***
	Metric    Metric      `json:"metric"`
	Value     SampleValue `json:"value"`
	Timestamp Time        `json:"timestamp"`
***REMOVED***

// Equal compares first the metrics, then the timestamp, then the value. The
// sematics of value equality is defined by SampleValue.Equal.
func (s *Sample) Equal(o *Sample) bool ***REMOVED***
	if s == o ***REMOVED***
		return true
	***REMOVED***

	if !s.Metric.Equal(o.Metric) ***REMOVED***
		return false
	***REMOVED***
	if !s.Timestamp.Equal(o.Timestamp) ***REMOVED***
		return false
	***REMOVED***
	if s.Value.Equal(o.Value) ***REMOVED***
		return false
	***REMOVED***

	return true
***REMOVED***

func (s Sample) String() string ***REMOVED***
	return fmt.Sprintf("%s => %s", s.Metric, SamplePair***REMOVED***
		Timestamp: s.Timestamp,
		Value:     s.Value,
	***REMOVED***)
***REMOVED***

// MarshalJSON implements json.Marshaler.
func (s Sample) MarshalJSON() ([]byte, error) ***REMOVED***
	v := struct ***REMOVED***
		Metric Metric     `json:"metric"`
		Value  SamplePair `json:"value"`
	***REMOVED******REMOVED***
		Metric: s.Metric,
		Value: SamplePair***REMOVED***
			Timestamp: s.Timestamp,
			Value:     s.Value,
		***REMOVED***,
	***REMOVED***

	return json.Marshal(&v)
***REMOVED***

// UnmarshalJSON implements json.Unmarshaler.
func (s *Sample) UnmarshalJSON(b []byte) error ***REMOVED***
	v := struct ***REMOVED***
		Metric Metric     `json:"metric"`
		Value  SamplePair `json:"value"`
	***REMOVED******REMOVED***
		Metric: s.Metric,
		Value: SamplePair***REMOVED***
			Timestamp: s.Timestamp,
			Value:     s.Value,
		***REMOVED***,
	***REMOVED***

	if err := json.Unmarshal(b, &v); err != nil ***REMOVED***
		return err
	***REMOVED***

	s.Metric = v.Metric
	s.Timestamp = v.Value.Timestamp
	s.Value = v.Value.Value

	return nil
***REMOVED***

// Samples is a sortable Sample slice. It implements sort.Interface.
type Samples []*Sample

func (s Samples) Len() int ***REMOVED***
	return len(s)
***REMOVED***

// Less compares first the metrics, then the timestamp.
func (s Samples) Less(i, j int) bool ***REMOVED***
	switch ***REMOVED***
	case s[i].Metric.Before(s[j].Metric):
		return true
	case s[j].Metric.Before(s[i].Metric):
		return false
	case s[i].Timestamp.Before(s[j].Timestamp):
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

func (s Samples) Swap(i, j int) ***REMOVED***
	s[i], s[j] = s[j], s[i]
***REMOVED***

// Equal compares two sets of samples and returns true if they are equal.
func (s Samples) Equal(o Samples) bool ***REMOVED***
	if len(s) != len(o) ***REMOVED***
		return false
	***REMOVED***

	for i, sample := range s ***REMOVED***
		if !sample.Equal(o[i]) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// SampleStream is a stream of Values belonging to an attached COWMetric.
type SampleStream struct ***REMOVED***
	Metric Metric       `json:"metric"`
	Values []SamplePair `json:"values"`
***REMOVED***

func (ss SampleStream) String() string ***REMOVED***
	vals := make([]string, len(ss.Values))
	for i, v := range ss.Values ***REMOVED***
		vals[i] = v.String()
	***REMOVED***
	return fmt.Sprintf("%s =>\n%s", ss.Metric, strings.Join(vals, "\n"))
***REMOVED***

// Value is a generic interface for values resulting from a query evaluation.
type Value interface ***REMOVED***
	Type() ValueType
	String() string
***REMOVED***

func (Matrix) Type() ValueType  ***REMOVED*** return ValMatrix ***REMOVED***
func (Vector) Type() ValueType  ***REMOVED*** return ValVector ***REMOVED***
func (*Scalar) Type() ValueType ***REMOVED*** return ValScalar ***REMOVED***
func (*String) Type() ValueType ***REMOVED*** return ValString ***REMOVED***

type ValueType int

const (
	ValNone ValueType = iota
	ValScalar
	ValVector
	ValMatrix
	ValString
)

// MarshalJSON implements json.Marshaler.
func (et ValueType) MarshalJSON() ([]byte, error) ***REMOVED***
	return json.Marshal(et.String())
***REMOVED***

func (et *ValueType) UnmarshalJSON(b []byte) error ***REMOVED***
	var s string
	if err := json.Unmarshal(b, &s); err != nil ***REMOVED***
		return err
	***REMOVED***
	switch s ***REMOVED***
	case "<ValNone>":
		*et = ValNone
	case "scalar":
		*et = ValScalar
	case "vector":
		*et = ValVector
	case "matrix":
		*et = ValMatrix
	case "string":
		*et = ValString
	default:
		return fmt.Errorf("unknown value type %q", s)
	***REMOVED***
	return nil
***REMOVED***

func (e ValueType) String() string ***REMOVED***
	switch e ***REMOVED***
	case ValNone:
		return "<ValNone>"
	case ValScalar:
		return "scalar"
	case ValVector:
		return "vector"
	case ValMatrix:
		return "matrix"
	case ValString:
		return "string"
	***REMOVED***
	panic("ValueType.String: unhandled value type")
***REMOVED***

// Scalar is a scalar value evaluated at the set timestamp.
type Scalar struct ***REMOVED***
	Value     SampleValue `json:"value"`
	Timestamp Time        `json:"timestamp"`
***REMOVED***

func (s Scalar) String() string ***REMOVED***
	return fmt.Sprintf("scalar: %v @[%v]", s.Value, s.Timestamp)
***REMOVED***

// MarshalJSON implements json.Marshaler.
func (s Scalar) MarshalJSON() ([]byte, error) ***REMOVED***
	v := strconv.FormatFloat(float64(s.Value), 'f', -1, 64)
	return json.Marshal([...]interface***REMOVED******REMOVED******REMOVED***s.Timestamp, string(v)***REMOVED***)
***REMOVED***

// UnmarshalJSON implements json.Unmarshaler.
func (s *Scalar) UnmarshalJSON(b []byte) error ***REMOVED***
	var f string
	v := [...]interface***REMOVED******REMOVED******REMOVED***&s.Timestamp, &f***REMOVED***

	if err := json.Unmarshal(b, &v); err != nil ***REMOVED***
		return err
	***REMOVED***

	value, err := strconv.ParseFloat(f, 64)
	if err != nil ***REMOVED***
		return fmt.Errorf("error parsing sample value: %s", err)
	***REMOVED***
	s.Value = SampleValue(value)
	return nil
***REMOVED***

// String is a string value evaluated at the set timestamp.
type String struct ***REMOVED***
	Value     string `json:"value"`
	Timestamp Time   `json:"timestamp"`
***REMOVED***

func (s *String) String() string ***REMOVED***
	return s.Value
***REMOVED***

// MarshalJSON implements json.Marshaler.
func (s String) MarshalJSON() ([]byte, error) ***REMOVED***
	return json.Marshal([]interface***REMOVED******REMOVED******REMOVED***s.Timestamp, s.Value***REMOVED***)
***REMOVED***

// UnmarshalJSON implements json.Unmarshaler.
func (s *String) UnmarshalJSON(b []byte) error ***REMOVED***
	v := [...]interface***REMOVED******REMOVED******REMOVED***&s.Timestamp, &s.Value***REMOVED***
	return json.Unmarshal(b, &v)
***REMOVED***

// Vector is basically only an alias for Samples, but the
// contract is that in a Vector, all Samples have the same timestamp.
type Vector []*Sample

func (vec Vector) String() string ***REMOVED***
	entries := make([]string, len(vec))
	for i, s := range vec ***REMOVED***
		entries[i] = s.String()
	***REMOVED***
	return strings.Join(entries, "\n")
***REMOVED***

func (vec Vector) Len() int      ***REMOVED*** return len(vec) ***REMOVED***
func (vec Vector) Swap(i, j int) ***REMOVED*** vec[i], vec[j] = vec[j], vec[i] ***REMOVED***

// Less compares first the metrics, then the timestamp.
func (vec Vector) Less(i, j int) bool ***REMOVED***
	switch ***REMOVED***
	case vec[i].Metric.Before(vec[j].Metric):
		return true
	case vec[j].Metric.Before(vec[i].Metric):
		return false
	case vec[i].Timestamp.Before(vec[j].Timestamp):
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

// Equal compares two sets of samples and returns true if they are equal.
func (vec Vector) Equal(o Vector) bool ***REMOVED***
	if len(vec) != len(o) ***REMOVED***
		return false
	***REMOVED***

	for i, sample := range vec ***REMOVED***
		if !sample.Equal(o[i]) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// Matrix is a list of time series.
type Matrix []*SampleStream

func (m Matrix) Len() int           ***REMOVED*** return len(m) ***REMOVED***
func (m Matrix) Less(i, j int) bool ***REMOVED*** return m[i].Metric.Before(m[j].Metric) ***REMOVED***
func (m Matrix) Swap(i, j int)      ***REMOVED*** m[i], m[j] = m[j], m[i] ***REMOVED***

func (mat Matrix) String() string ***REMOVED***
	matCp := make(Matrix, len(mat))
	copy(matCp, mat)
	sort.Sort(matCp)

	strs := make([]string, len(matCp))

	for i, ss := range matCp ***REMOVED***
		strs[i] = ss.String()
	***REMOVED***

	return strings.Join(strs, "\n")
***REMOVED***
