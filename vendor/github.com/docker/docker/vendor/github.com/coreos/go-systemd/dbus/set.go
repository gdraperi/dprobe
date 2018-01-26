// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dbus

type set struct ***REMOVED***
	data map[string]bool
***REMOVED***

func (s *set) Add(value string) ***REMOVED***
	s.data[value] = true
***REMOVED***

func (s *set) Remove(value string) ***REMOVED***
	delete(s.data, value)
***REMOVED***

func (s *set) Contains(value string) (exists bool) ***REMOVED***
	_, exists = s.data[value]
	return
***REMOVED***

func (s *set) Length() int ***REMOVED***
	return len(s.data)
***REMOVED***

func (s *set) Values() (values []string) ***REMOVED***
	for val, _ := range s.data ***REMOVED***
		values = append(values, val)
	***REMOVED***
	return
***REMOVED***

func newSet() *set ***REMOVED***
	return &set***REMOVED***make(map[string]bool)***REMOVED***
***REMOVED***
