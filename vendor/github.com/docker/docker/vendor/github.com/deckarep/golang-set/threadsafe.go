/*
Open Source Initiative OSI - The MIT License (MIT):Licensing

The MIT License (MIT)
Copyright (c) 2013 Ralph Caraveo (deckarep@gmail.com)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package mapset

import "sync"

type threadSafeSet struct ***REMOVED***
	s threadUnsafeSet
	sync.RWMutex
***REMOVED***

func newThreadSafeSet() threadSafeSet ***REMOVED***
	return threadSafeSet***REMOVED***s: newThreadUnsafeSet()***REMOVED***
***REMOVED***

func (set *threadSafeSet) Add(i interface***REMOVED******REMOVED***) bool ***REMOVED***
	set.Lock()
	ret := set.s.Add(i)
	set.Unlock()
	return ret
***REMOVED***

func (set *threadSafeSet) Contains(i ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	set.RLock()
	ret := set.s.Contains(i...)
	set.RUnlock()
	return ret
***REMOVED***

func (set *threadSafeSet) IsSubset(other Set) bool ***REMOVED***
	o := other.(*threadSafeSet)

	set.RLock()
	o.RLock()

	ret := set.s.IsSubset(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
***REMOVED***

func (set *threadSafeSet) IsSuperset(other Set) bool ***REMOVED***
	return other.IsSubset(set)
***REMOVED***

func (set *threadSafeSet) Union(other Set) Set ***REMOVED***
	o := other.(*threadSafeSet)

	set.RLock()
	o.RLock()

	unsafeUnion := set.s.Union(&o.s).(*threadUnsafeSet)
	ret := &threadSafeSet***REMOVED***s: *unsafeUnion***REMOVED***
	set.RUnlock()
	o.RUnlock()
	return ret
***REMOVED***

func (set *threadSafeSet) Intersect(other Set) Set ***REMOVED***
	o := other.(*threadSafeSet)

	set.RLock()
	o.RLock()

	unsafeIntersection := set.s.Intersect(&o.s).(*threadUnsafeSet)
	ret := &threadSafeSet***REMOVED***s: *unsafeIntersection***REMOVED***
	set.RUnlock()
	o.RUnlock()
	return ret
***REMOVED***

func (set *threadSafeSet) Difference(other Set) Set ***REMOVED***
	o := other.(*threadSafeSet)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.Difference(&o.s).(*threadUnsafeSet)
	ret := &threadSafeSet***REMOVED***s: *unsafeDifference***REMOVED***
	set.RUnlock()
	o.RUnlock()
	return ret
***REMOVED***

func (set *threadSafeSet) SymmetricDifference(other Set) Set ***REMOVED***
	o := other.(*threadSafeSet)

	unsafeDifference := set.s.SymmetricDifference(&o.s).(*threadUnsafeSet)
	return &threadSafeSet***REMOVED***s: *unsafeDifference***REMOVED***
***REMOVED***

func (set *threadSafeSet) Clear() ***REMOVED***
	set.Lock()
	set.s = newThreadUnsafeSet()
	set.Unlock()
***REMOVED***

func (set *threadSafeSet) Remove(i interface***REMOVED******REMOVED***) ***REMOVED***
	set.Lock()
	delete(set.s, i)
	set.Unlock()
***REMOVED***

func (set *threadSafeSet) Cardinality() int ***REMOVED***
	set.RLock()
	defer set.RUnlock()
	return len(set.s)
***REMOVED***

func (set *threadSafeSet) Iter() <-chan interface***REMOVED******REMOVED*** ***REMOVED***
	ch := make(chan interface***REMOVED******REMOVED***)
	go func() ***REMOVED***
		set.RLock()

		for elem := range set.s ***REMOVED***
			ch <- elem
		***REMOVED***
		close(ch)
		set.RUnlock()
	***REMOVED***()

	return ch
***REMOVED***

func (set *threadSafeSet) Equal(other Set) bool ***REMOVED***
	o := other.(*threadSafeSet)

	set.RLock()
	o.RLock()

	ret := set.s.Equal(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
***REMOVED***

func (set *threadSafeSet) Clone() Set ***REMOVED***
	set.RLock()

	unsafeClone := set.s.Clone().(*threadUnsafeSet)
	ret := &threadSafeSet***REMOVED***s: *unsafeClone***REMOVED***
	set.RUnlock()
	return ret
***REMOVED***

func (set *threadSafeSet) String() string ***REMOVED***
	set.RLock()
	ret := set.s.String()
	set.RUnlock()
	return ret
***REMOVED***

func (set *threadSafeSet) PowerSet() Set ***REMOVED***
	set.RLock()
	ret := set.s.PowerSet()
	set.RUnlock()
	return ret
***REMOVED***

func (set *threadSafeSet) CartesianProduct(other Set) Set ***REMOVED***
	o := other.(*threadSafeSet)

	set.RLock()
	o.RLock()

	unsafeCartProduct := set.s.CartesianProduct(&o.s).(*threadUnsafeSet)
	ret := &threadSafeSet***REMOVED***s: *unsafeCartProduct***REMOVED***
	set.RUnlock()
	o.RUnlock()
	return ret
***REMOVED***

func (set *threadSafeSet) ToSlice() []interface***REMOVED******REMOVED*** ***REMOVED***
	set.RLock()
	keys := make([]interface***REMOVED******REMOVED***, 0, set.Cardinality())
	for elem := range set.s ***REMOVED***
		keys = append(keys, elem)
	***REMOVED***
	set.RUnlock()
	return keys
***REMOVED***
