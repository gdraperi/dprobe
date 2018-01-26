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

package model

import (
	"sort"
)

// SeparatorByte is a byte that cannot occur in valid UTF-8 sequences and is
// used to separate label names, label values, and other strings from each other
// when calculating their combined hash value (aka signature aka fingerprint).
const SeparatorByte byte = 255

var (
	// cache the signature of an empty label set.
	emptyLabelSignature = hashNew()
)

// LabelsToSignature returns a quasi-unique signature (i.e., fingerprint) for a
// given label set. (Collisions are possible but unlikely if the number of label
// sets the function is applied to is small.)
func LabelsToSignature(labels map[string]string) uint64 ***REMOVED***
	if len(labels) == 0 ***REMOVED***
		return emptyLabelSignature
	***REMOVED***

	labelNames := make([]string, 0, len(labels))
	for labelName := range labels ***REMOVED***
		labelNames = append(labelNames, labelName)
	***REMOVED***
	sort.Strings(labelNames)

	sum := hashNew()
	for _, labelName := range labelNames ***REMOVED***
		sum = hashAdd(sum, labelName)
		sum = hashAddByte(sum, SeparatorByte)
		sum = hashAdd(sum, labels[labelName])
		sum = hashAddByte(sum, SeparatorByte)
	***REMOVED***
	return sum
***REMOVED***

// labelSetToFingerprint works exactly as LabelsToSignature but takes a LabelSet as
// parameter (rather than a label map) and returns a Fingerprint.
func labelSetToFingerprint(ls LabelSet) Fingerprint ***REMOVED***
	if len(ls) == 0 ***REMOVED***
		return Fingerprint(emptyLabelSignature)
	***REMOVED***

	labelNames := make(LabelNames, 0, len(ls))
	for labelName := range ls ***REMOVED***
		labelNames = append(labelNames, labelName)
	***REMOVED***
	sort.Sort(labelNames)

	sum := hashNew()
	for _, labelName := range labelNames ***REMOVED***
		sum = hashAdd(sum, string(labelName))
		sum = hashAddByte(sum, SeparatorByte)
		sum = hashAdd(sum, string(ls[labelName]))
		sum = hashAddByte(sum, SeparatorByte)
	***REMOVED***
	return Fingerprint(sum)
***REMOVED***

// labelSetToFastFingerprint works similar to labelSetToFingerprint but uses a
// faster and less allocation-heavy hash function, which is more susceptible to
// create hash collisions. Therefore, collision detection should be applied.
func labelSetToFastFingerprint(ls LabelSet) Fingerprint ***REMOVED***
	if len(ls) == 0 ***REMOVED***
		return Fingerprint(emptyLabelSignature)
	***REMOVED***

	var result uint64
	for labelName, labelValue := range ls ***REMOVED***
		sum := hashNew()
		sum = hashAdd(sum, string(labelName))
		sum = hashAddByte(sum, SeparatorByte)
		sum = hashAdd(sum, string(labelValue))
		result ^= sum
	***REMOVED***
	return Fingerprint(result)
***REMOVED***

// SignatureForLabels works like LabelsToSignature but takes a Metric as
// parameter (rather than a label map) and only includes the labels with the
// specified LabelNames into the signature calculation. The labels passed in
// will be sorted by this function.
func SignatureForLabels(m Metric, labels ...LabelName) uint64 ***REMOVED***
	if len(labels) == 0 ***REMOVED***
		return emptyLabelSignature
	***REMOVED***

	sort.Sort(LabelNames(labels))

	sum := hashNew()
	for _, label := range labels ***REMOVED***
		sum = hashAdd(sum, string(label))
		sum = hashAddByte(sum, SeparatorByte)
		sum = hashAdd(sum, string(m[label]))
		sum = hashAddByte(sum, SeparatorByte)
	***REMOVED***
	return sum
***REMOVED***

// SignatureWithoutLabels works like LabelsToSignature but takes a Metric as
// parameter (rather than a label map) and excludes the labels with any of the
// specified LabelNames from the signature calculation.
func SignatureWithoutLabels(m Metric, labels map[LabelName]struct***REMOVED******REMOVED***) uint64 ***REMOVED***
	if len(m) == 0 ***REMOVED***
		return emptyLabelSignature
	***REMOVED***

	labelNames := make(LabelNames, 0, len(m))
	for labelName := range m ***REMOVED***
		if _, exclude := labels[labelName]; !exclude ***REMOVED***
			labelNames = append(labelNames, labelName)
		***REMOVED***
	***REMOVED***
	if len(labelNames) == 0 ***REMOVED***
		return emptyLabelSignature
	***REMOVED***
	sort.Sort(labelNames)

	sum := hashNew()
	for _, labelName := range labelNames ***REMOVED***
		sum = hashAdd(sum, string(labelName))
		sum = hashAddByte(sum, SeparatorByte)
		sum = hashAdd(sum, string(m[labelName]))
		sum = hashAddByte(sum, SeparatorByte)
	***REMOVED***
	return sum
***REMOVED***
