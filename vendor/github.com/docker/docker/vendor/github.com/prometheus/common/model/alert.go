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
	"time"
)

type AlertStatus string

const (
	AlertFiring   AlertStatus = "firing"
	AlertResolved AlertStatus = "resolved"
)

// Alert is a generic representation of an alert in the Prometheus eco-system.
type Alert struct ***REMOVED***
	// Label value pairs for purpose of aggregation, matching, and disposition
	// dispatching. This must minimally include an "alertname" label.
	Labels LabelSet `json:"labels"`

	// Extra key/value information which does not define alert identity.
	Annotations LabelSet `json:"annotations"`

	// The known time range for this alert. Both ends are optional.
	StartsAt     time.Time `json:"startsAt,omitempty"`
	EndsAt       time.Time `json:"endsAt,omitempty"`
	GeneratorURL string    `json:"generatorURL"`
***REMOVED***

// Name returns the name of the alert. It is equivalent to the "alertname" label.
func (a *Alert) Name() string ***REMOVED***
	return string(a.Labels[AlertNameLabel])
***REMOVED***

// Fingerprint returns a unique hash for the alert. It is equivalent to
// the fingerprint of the alert's label set.
func (a *Alert) Fingerprint() Fingerprint ***REMOVED***
	return a.Labels.Fingerprint()
***REMOVED***

func (a *Alert) String() string ***REMOVED***
	s := fmt.Sprintf("%s[%s]", a.Name(), a.Fingerprint().String()[:7])
	if a.Resolved() ***REMOVED***
		return s + "[resolved]"
	***REMOVED***
	return s + "[active]"
***REMOVED***

// Resolved returns true iff the activity interval ended in the past.
func (a *Alert) Resolved() bool ***REMOVED***
	return a.ResolvedAt(time.Now())
***REMOVED***

// ResolvedAt returns true off the activity interval ended before
// the given timestamp.
func (a *Alert) ResolvedAt(ts time.Time) bool ***REMOVED***
	if a.EndsAt.IsZero() ***REMOVED***
		return false
	***REMOVED***
	return !a.EndsAt.After(ts)
***REMOVED***

// Status returns the status of the alert.
func (a *Alert) Status() AlertStatus ***REMOVED***
	if a.Resolved() ***REMOVED***
		return AlertResolved
	***REMOVED***
	return AlertFiring
***REMOVED***

// Validate checks whether the alert data is inconsistent.
func (a *Alert) Validate() error ***REMOVED***
	if a.StartsAt.IsZero() ***REMOVED***
		return fmt.Errorf("start time missing")
	***REMOVED***
	if !a.EndsAt.IsZero() && a.EndsAt.Before(a.StartsAt) ***REMOVED***
		return fmt.Errorf("start time must be before end time")
	***REMOVED***
	if err := a.Labels.Validate(); err != nil ***REMOVED***
		return fmt.Errorf("invalid label set: %s", err)
	***REMOVED***
	if len(a.Labels) == 0 ***REMOVED***
		return fmt.Errorf("at least one label pair required")
	***REMOVED***
	if err := a.Annotations.Validate(); err != nil ***REMOVED***
		return fmt.Errorf("invalid annotations: %s", err)
	***REMOVED***
	return nil
***REMOVED***

// Alert is a list of alerts that can be sorted in chronological order.
type Alerts []*Alert

func (as Alerts) Len() int      ***REMOVED*** return len(as) ***REMOVED***
func (as Alerts) Swap(i, j int) ***REMOVED*** as[i], as[j] = as[j], as[i] ***REMOVED***

func (as Alerts) Less(i, j int) bool ***REMOVED***
	if as[i].StartsAt.Before(as[j].StartsAt) ***REMOVED***
		return true
	***REMOVED***
	if as[i].EndsAt.Before(as[j].EndsAt) ***REMOVED***
		return true
	***REMOVED***
	return as[i].Fingerprint() < as[j].Fingerprint()
***REMOVED***

// HasFiring returns true iff one of the alerts is not resolved.
func (as Alerts) HasFiring() bool ***REMOVED***
	for _, a := range as ***REMOVED***
		if !a.Resolved() ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// Status returns StatusFiring iff at least one of the alerts is firing.
func (as Alerts) Status() AlertStatus ***REMOVED***
	if as.HasFiring() ***REMOVED***
		return AlertFiring
	***REMOVED***
	return AlertResolved
***REMOVED***
