// Copyright 2015 The etcd Authors
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

package types

import (
	"reflect"
	"sort"
	"sync"
)

type Set interface ***REMOVED***
	Add(string)
	Remove(string)
	Contains(string) bool
	Equals(Set) bool
	Length() int
	Values() []string
	Copy() Set
	Sub(Set) Set
***REMOVED***

func NewUnsafeSet(values ...string) *unsafeSet ***REMOVED***
	set := &unsafeSet***REMOVED***make(map[string]struct***REMOVED******REMOVED***)***REMOVED***
	for _, v := range values ***REMOVED***
		set.Add(v)
	***REMOVED***
	return set
***REMOVED***

func NewThreadsafeSet(values ...string) *tsafeSet ***REMOVED***
	us := NewUnsafeSet(values...)
	return &tsafeSet***REMOVED***us, sync.RWMutex***REMOVED******REMOVED******REMOVED***
***REMOVED***

type unsafeSet struct ***REMOVED***
	d map[string]struct***REMOVED******REMOVED***
***REMOVED***

// Add adds a new value to the set (no-op if the value is already present)
func (us *unsafeSet) Add(value string) ***REMOVED***
	us.d[value] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
***REMOVED***

// Remove removes the given value from the set
func (us *unsafeSet) Remove(value string) ***REMOVED***
	delete(us.d, value)
***REMOVED***

// Contains returns whether the set contains the given value
func (us *unsafeSet) Contains(value string) (exists bool) ***REMOVED***
	_, exists = us.d[value]
	return
***REMOVED***

// ContainsAll returns whether the set contains all given values
func (us *unsafeSet) ContainsAll(values []string) bool ***REMOVED***
	for _, s := range values ***REMOVED***
		if !us.Contains(s) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// Equals returns whether the contents of two sets are identical
func (us *unsafeSet) Equals(other Set) bool ***REMOVED***
	v1 := sort.StringSlice(us.Values())
	v2 := sort.StringSlice(other.Values())
	v1.Sort()
	v2.Sort()
	return reflect.DeepEqual(v1, v2)
***REMOVED***

// Length returns the number of elements in the set
func (us *unsafeSet) Length() int ***REMOVED***
	return len(us.d)
***REMOVED***

// Values returns the values of the Set in an unspecified order.
func (us *unsafeSet) Values() (values []string) ***REMOVED***
	values = make([]string, 0)
	for val := range us.d ***REMOVED***
		values = append(values, val)
	***REMOVED***
	return
***REMOVED***

// Copy creates a new Set containing the values of the first
func (us *unsafeSet) Copy() Set ***REMOVED***
	cp := NewUnsafeSet()
	for val := range us.d ***REMOVED***
		cp.Add(val)
	***REMOVED***

	return cp
***REMOVED***

// Sub removes all elements in other from the set
func (us *unsafeSet) Sub(other Set) Set ***REMOVED***
	oValues := other.Values()
	result := us.Copy().(*unsafeSet)

	for _, val := range oValues ***REMOVED***
		if _, ok := result.d[val]; !ok ***REMOVED***
			continue
		***REMOVED***
		delete(result.d, val)
	***REMOVED***

	return result
***REMOVED***

type tsafeSet struct ***REMOVED***
	us *unsafeSet
	m  sync.RWMutex
***REMOVED***

func (ts *tsafeSet) Add(value string) ***REMOVED***
	ts.m.Lock()
	defer ts.m.Unlock()
	ts.us.Add(value)
***REMOVED***

func (ts *tsafeSet) Remove(value string) ***REMOVED***
	ts.m.Lock()
	defer ts.m.Unlock()
	ts.us.Remove(value)
***REMOVED***

func (ts *tsafeSet) Contains(value string) (exists bool) ***REMOVED***
	ts.m.RLock()
	defer ts.m.RUnlock()
	return ts.us.Contains(value)
***REMOVED***

func (ts *tsafeSet) Equals(other Set) bool ***REMOVED***
	ts.m.RLock()
	defer ts.m.RUnlock()
	return ts.us.Equals(other)
***REMOVED***

func (ts *tsafeSet) Length() int ***REMOVED***
	ts.m.RLock()
	defer ts.m.RUnlock()
	return ts.us.Length()
***REMOVED***

func (ts *tsafeSet) Values() (values []string) ***REMOVED***
	ts.m.RLock()
	defer ts.m.RUnlock()
	return ts.us.Values()
***REMOVED***

func (ts *tsafeSet) Copy() Set ***REMOVED***
	ts.m.RLock()
	defer ts.m.RUnlock()
	usResult := ts.us.Copy().(*unsafeSet)
	return &tsafeSet***REMOVED***usResult, sync.RWMutex***REMOVED******REMOVED******REMOVED***
***REMOVED***

func (ts *tsafeSet) Sub(other Set) Set ***REMOVED***
	ts.m.RLock()
	defer ts.m.RUnlock()
	usResult := ts.us.Sub(other).(*unsafeSet)
	return &tsafeSet***REMOVED***usResult, sync.RWMutex***REMOVED******REMOVED******REMOVED***
***REMOVED***
