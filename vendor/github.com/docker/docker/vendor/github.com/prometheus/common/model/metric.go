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
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var (
	separator    = []byte***REMOVED***0***REMOVED***
	MetricNameRE = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_:]*$`)
)

// A Metric is similar to a LabelSet, but the key difference is that a Metric is
// a singleton and refers to one and only one stream of samples.
type Metric LabelSet

// Equal compares the metrics.
func (m Metric) Equal(o Metric) bool ***REMOVED***
	return LabelSet(m).Equal(LabelSet(o))
***REMOVED***

// Before compares the metrics' underlying label sets.
func (m Metric) Before(o Metric) bool ***REMOVED***
	return LabelSet(m).Before(LabelSet(o))
***REMOVED***

// Clone returns a copy of the Metric.
func (m Metric) Clone() Metric ***REMOVED***
	clone := Metric***REMOVED******REMOVED***
	for k, v := range m ***REMOVED***
		clone[k] = v
	***REMOVED***
	return clone
***REMOVED***

func (m Metric) String() string ***REMOVED***
	metricName, hasName := m[MetricNameLabel]
	numLabels := len(m) - 1
	if !hasName ***REMOVED***
		numLabels = len(m)
	***REMOVED***
	labelStrings := make([]string, 0, numLabels)
	for label, value := range m ***REMOVED***
		if label != MetricNameLabel ***REMOVED***
			labelStrings = append(labelStrings, fmt.Sprintf("%s=%q", label, value))
		***REMOVED***
	***REMOVED***

	switch numLabels ***REMOVED***
	case 0:
		if hasName ***REMOVED***
			return string(metricName)
		***REMOVED***
		return "***REMOVED******REMOVED***"
	default:
		sort.Strings(labelStrings)
		return fmt.Sprintf("%s***REMOVED***%s***REMOVED***", metricName, strings.Join(labelStrings, ", "))
	***REMOVED***
***REMOVED***

// Fingerprint returns a Metric's Fingerprint.
func (m Metric) Fingerprint() Fingerprint ***REMOVED***
	return LabelSet(m).Fingerprint()
***REMOVED***

// FastFingerprint returns a Metric's Fingerprint calculated by a faster hashing
// algorithm, which is, however, more susceptible to hash collisions.
func (m Metric) FastFingerprint() Fingerprint ***REMOVED***
	return LabelSet(m).FastFingerprint()
***REMOVED***

// IsValidMetricName returns true iff name matches the pattern of MetricNameRE.
func IsValidMetricName(n LabelValue) bool ***REMOVED***
	if len(n) == 0 ***REMOVED***
		return false
	***REMOVED***
	for i, b := range n ***REMOVED***
		if !((b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_' || b == ':' || (b >= '0' && b <= '9' && i > 0)) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***
