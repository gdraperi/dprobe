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

package model

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)

// Matcher describes a matches the value of a given label.
type Matcher struct ***REMOVED***
	Name    LabelName `json:"name"`
	Value   string    `json:"value"`
	IsRegex bool      `json:"isRegex"`
***REMOVED***

func (m *Matcher) UnmarshalJSON(b []byte) error ***REMOVED***
	type plain Matcher
	if err := json.Unmarshal(b, (*plain)(m)); err != nil ***REMOVED***
		return err
	***REMOVED***

	if len(m.Name) == 0 ***REMOVED***
		return fmt.Errorf("label name in matcher must not be empty")
	***REMOVED***
	if m.IsRegex ***REMOVED***
		if _, err := regexp.Compile(m.Value); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Validate returns true iff all fields of the matcher have valid values.
func (m *Matcher) Validate() error ***REMOVED***
	if !m.Name.IsValid() ***REMOVED***
		return fmt.Errorf("invalid name %q", m.Name)
	***REMOVED***
	if m.IsRegex ***REMOVED***
		if _, err := regexp.Compile(m.Value); err != nil ***REMOVED***
			return fmt.Errorf("invalid regular expression %q", m.Value)
		***REMOVED***
	***REMOVED*** else if !LabelValue(m.Value).IsValid() || len(m.Value) == 0 ***REMOVED***
		return fmt.Errorf("invalid value %q", m.Value)
	***REMOVED***
	return nil
***REMOVED***

// Silence defines the representation of a silence definiton
// in the Prometheus eco-system.
type Silence struct ***REMOVED***
	ID uint64 `json:"id,omitempty"`

	Matchers []*Matcher `json:"matchers"`

	StartsAt time.Time `json:"startsAt"`
	EndsAt   time.Time `json:"endsAt"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	CreatedBy string    `json:"createdBy"`
	Comment   string    `json:"comment,omitempty"`
***REMOVED***

// Validate returns true iff all fields of the silence have valid values.
func (s *Silence) Validate() error ***REMOVED***
	if len(s.Matchers) == 0 ***REMOVED***
		return fmt.Errorf("at least one matcher required")
	***REMOVED***
	for _, m := range s.Matchers ***REMOVED***
		if err := m.Validate(); err != nil ***REMOVED***
			return fmt.Errorf("invalid matcher: %s", err)
		***REMOVED***
	***REMOVED***
	if s.StartsAt.IsZero() ***REMOVED***
		return fmt.Errorf("start time missing")
	***REMOVED***
	if s.EndsAt.IsZero() ***REMOVED***
		return fmt.Errorf("end time missing")
	***REMOVED***
	if s.EndsAt.Before(s.StartsAt) ***REMOVED***
		return fmt.Errorf("start time must be before end time")
	***REMOVED***
	if s.CreatedBy == "" ***REMOVED***
		return fmt.Errorf("creator information missing")
	***REMOVED***
	if s.Comment == "" ***REMOVED***
		return fmt.Errorf("comment missing")
	***REMOVED***
	if s.CreatedAt.IsZero() ***REMOVED***
		return fmt.Errorf("creation timestamp missing")
	***REMOVED***
	return nil
***REMOVED***
