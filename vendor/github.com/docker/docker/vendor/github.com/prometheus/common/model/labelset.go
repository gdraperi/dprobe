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
	"sort"
	"strings"
)

// A LabelSet is a collection of LabelName and LabelValue pairs.  The LabelSet
// may be fully-qualified down to the point where it may resolve to a single
// Metric in the data store or not.  All operations that occur within the realm
// of a LabelSet can emit a vector of Metric entities to which the LabelSet may
// match.
type LabelSet map[LabelName]LabelValue

// Validate checks whether all names and values in the label set
// are valid.
func (ls LabelSet) Validate() error ***REMOVED***
	for ln, lv := range ls ***REMOVED***
		if !ln.IsValid() ***REMOVED***
			return fmt.Errorf("invalid name %q", ln)
		***REMOVED***
		if !lv.IsValid() ***REMOVED***
			return fmt.Errorf("invalid value %q", lv)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Equal returns true iff both label sets have exactly the same key/value pairs.
func (ls LabelSet) Equal(o LabelSet) bool ***REMOVED***
	if len(ls) != len(o) ***REMOVED***
		return false
	***REMOVED***
	for ln, lv := range ls ***REMOVED***
		olv, ok := o[ln]
		if !ok ***REMOVED***
			return false
		***REMOVED***
		if olv != lv ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// Before compares the metrics, using the following criteria:
//
// If m has fewer labels than o, it is before o. If it has more, it is not.
//
// If the number of labels is the same, the superset of all label names is
// sorted alphanumerically. The first differing label pair found in that order
// determines the outcome: If the label does not exist at all in m, then m is
// before o, and vice versa. Otherwise the label value is compared
// alphanumerically.
//
// If m and o are equal, the method returns false.
func (ls LabelSet) Before(o LabelSet) bool ***REMOVED***
	if len(ls) < len(o) ***REMOVED***
		return true
	***REMOVED***
	if len(ls) > len(o) ***REMOVED***
		return false
	***REMOVED***

	lns := make(LabelNames, 0, len(ls)+len(o))
	for ln := range ls ***REMOVED***
		lns = append(lns, ln)
	***REMOVED***
	for ln := range o ***REMOVED***
		lns = append(lns, ln)
	***REMOVED***
	// It's probably not worth it to de-dup lns.
	sort.Sort(lns)
	for _, ln := range lns ***REMOVED***
		mlv, ok := ls[ln]
		if !ok ***REMOVED***
			return true
		***REMOVED***
		olv, ok := o[ln]
		if !ok ***REMOVED***
			return false
		***REMOVED***
		if mlv < olv ***REMOVED***
			return true
		***REMOVED***
		if mlv > olv ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// Clone returns a copy of the label set.
func (ls LabelSet) Clone() LabelSet ***REMOVED***
	lsn := make(LabelSet, len(ls))
	for ln, lv := range ls ***REMOVED***
		lsn[ln] = lv
	***REMOVED***
	return lsn
***REMOVED***

// Merge is a helper function to non-destructively merge two label sets.
func (l LabelSet) Merge(other LabelSet) LabelSet ***REMOVED***
	result := make(LabelSet, len(l))

	for k, v := range l ***REMOVED***
		result[k] = v
	***REMOVED***

	for k, v := range other ***REMOVED***
		result[k] = v
	***REMOVED***

	return result
***REMOVED***

func (l LabelSet) String() string ***REMOVED***
	lstrs := make([]string, 0, len(l))
	for l, v := range l ***REMOVED***
		lstrs = append(lstrs, fmt.Sprintf("%s=%q", l, v))
	***REMOVED***

	sort.Strings(lstrs)
	return fmt.Sprintf("***REMOVED***%s***REMOVED***", strings.Join(lstrs, ", "))
***REMOVED***

// Fingerprint returns the LabelSet's fingerprint.
func (ls LabelSet) Fingerprint() Fingerprint ***REMOVED***
	return labelSetToFingerprint(ls)
***REMOVED***

// FastFingerprint returns the LabelSet's Fingerprint calculated by a faster hashing
// algorithm, which is, however, more susceptible to hash collisions.
func (ls LabelSet) FastFingerprint() Fingerprint ***REMOVED***
	return labelSetToFastFingerprint(ls)
***REMOVED***

// UnmarshalJSON implements the json.Unmarshaler interface.
func (l *LabelSet) UnmarshalJSON(b []byte) error ***REMOVED***
	var m map[LabelName]LabelValue
	if err := json.Unmarshal(b, &m); err != nil ***REMOVED***
		return err
	***REMOVED***
	// encoding/json only unmarshals maps of the form map[string]T. It treats
	// LabelName as a string and does not call its UnmarshalJSON method.
	// Thus, we have to replicate the behavior here.
	for ln := range m ***REMOVED***
		if !LabelNameRE.MatchString(string(ln)) ***REMOVED***
			return fmt.Errorf("%q is not a valid label name", ln)
		***REMOVED***
	***REMOVED***
	*l = LabelSet(m)
	return nil
***REMOVED***
