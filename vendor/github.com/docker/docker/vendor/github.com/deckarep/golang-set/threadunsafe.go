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

import (
	"fmt"
	"reflect"
	"strings"
)

type threadUnsafeSet map[interface***REMOVED******REMOVED***]struct***REMOVED******REMOVED***

type orderedPair struct ***REMOVED***
	first  interface***REMOVED******REMOVED***
	second interface***REMOVED******REMOVED***
***REMOVED***

func newThreadUnsafeSet() threadUnsafeSet ***REMOVED***
	return make(threadUnsafeSet)
***REMOVED***

func (pair *orderedPair) Equal(other orderedPair) bool ***REMOVED***
	if pair.first == other.first &&
		pair.second == other.second ***REMOVED***
		return true
	***REMOVED***

	return false
***REMOVED***

func (set *threadUnsafeSet) Add(i interface***REMOVED******REMOVED***) bool ***REMOVED***
	_, found := (*set)[i]
	(*set)[i] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	return !found //False if it existed already
***REMOVED***

func (set *threadUnsafeSet) Contains(i ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	for _, val := range i ***REMOVED***
		if _, ok := (*set)[val]; !ok ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (set *threadUnsafeSet) IsSubset(other Set) bool ***REMOVED***
	_ = other.(*threadUnsafeSet)
	for elem := range *set ***REMOVED***
		if !other.Contains(elem) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (set *threadUnsafeSet) IsSuperset(other Set) bool ***REMOVED***
	return other.IsSubset(set)
***REMOVED***

func (set *threadUnsafeSet) Union(other Set) Set ***REMOVED***
	o := other.(*threadUnsafeSet)

	unionedSet := newThreadUnsafeSet()

	for elem := range *set ***REMOVED***
		unionedSet.Add(elem)
	***REMOVED***
	for elem := range *o ***REMOVED***
		unionedSet.Add(elem)
	***REMOVED***
	return &unionedSet
***REMOVED***

func (set *threadUnsafeSet) Intersect(other Set) Set ***REMOVED***
	o := other.(*threadUnsafeSet)

	intersection := newThreadUnsafeSet()
	// loop over smaller set
	if set.Cardinality() < other.Cardinality() ***REMOVED***
		for elem := range *set ***REMOVED***
			if other.Contains(elem) ***REMOVED***
				intersection.Add(elem)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for elem := range *o ***REMOVED***
			if set.Contains(elem) ***REMOVED***
				intersection.Add(elem)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return &intersection
***REMOVED***

func (set *threadUnsafeSet) Difference(other Set) Set ***REMOVED***
	_ = other.(*threadUnsafeSet)

	difference := newThreadUnsafeSet()
	for elem := range *set ***REMOVED***
		if !other.Contains(elem) ***REMOVED***
			difference.Add(elem)
		***REMOVED***
	***REMOVED***
	return &difference
***REMOVED***

func (set *threadUnsafeSet) SymmetricDifference(other Set) Set ***REMOVED***
	_ = other.(*threadUnsafeSet)

	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
***REMOVED***

func (set *threadUnsafeSet) Clear() ***REMOVED***
	*set = newThreadUnsafeSet()
***REMOVED***

func (set *threadUnsafeSet) Remove(i interface***REMOVED******REMOVED***) ***REMOVED***
	delete(*set, i)
***REMOVED***

func (set *threadUnsafeSet) Cardinality() int ***REMOVED***
	return len(*set)
***REMOVED***

func (set *threadUnsafeSet) Iter() <-chan interface***REMOVED******REMOVED*** ***REMOVED***
	ch := make(chan interface***REMOVED******REMOVED***)
	go func() ***REMOVED***
		for elem := range *set ***REMOVED***
			ch <- elem
		***REMOVED***
		close(ch)
	***REMOVED***()

	return ch
***REMOVED***

func (set *threadUnsafeSet) Equal(other Set) bool ***REMOVED***
	_ = other.(*threadUnsafeSet)

	if set.Cardinality() != other.Cardinality() ***REMOVED***
		return false
	***REMOVED***
	for elem := range *set ***REMOVED***
		if !other.Contains(elem) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (set *threadUnsafeSet) Clone() Set ***REMOVED***
	clonedSet := newThreadUnsafeSet()
	for elem := range *set ***REMOVED***
		clonedSet.Add(elem)
	***REMOVED***
	return &clonedSet
***REMOVED***

func (set *threadUnsafeSet) String() string ***REMOVED***
	items := make([]string, 0, len(*set))

	for elem := range *set ***REMOVED***
		items = append(items, fmt.Sprintf("%v", elem))
	***REMOVED***
	return fmt.Sprintf("Set***REMOVED***%s***REMOVED***", strings.Join(items, ", "))
***REMOVED***

func (pair orderedPair) String() string ***REMOVED***
	return fmt.Sprintf("(%v, %v)", pair.first, pair.second)
***REMOVED***

func (set *threadUnsafeSet) PowerSet() Set ***REMOVED***
	powSet := NewThreadUnsafeSet()
	nullset := newThreadUnsafeSet()
	powSet.Add(&nullset)

	for es := range *set ***REMOVED***
		u := newThreadUnsafeSet()
		j := powSet.Iter()
		for er := range j ***REMOVED***
			p := newThreadUnsafeSet()
			if reflect.TypeOf(er).Name() == "" ***REMOVED***
				k := er.(*threadUnsafeSet)
				for ek := range *(k) ***REMOVED***
					p.Add(ek)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				p.Add(er)
			***REMOVED***
			p.Add(es)
			u.Add(&p)
		***REMOVED***

		powSet = powSet.Union(&u)
	***REMOVED***

	return powSet
***REMOVED***

func (set *threadUnsafeSet) CartesianProduct(other Set) Set ***REMOVED***
	o := other.(*threadUnsafeSet)
	cartProduct := NewThreadUnsafeSet()

	for i := range *set ***REMOVED***
		for j := range *o ***REMOVED***
			elem := orderedPair***REMOVED***first: i, second: j***REMOVED***
			cartProduct.Add(elem)
		***REMOVED***
	***REMOVED***

	return cartProduct
***REMOVED***

func (set *threadUnsafeSet) ToSlice() []interface***REMOVED******REMOVED*** ***REMOVED***
	keys := make([]interface***REMOVED******REMOVED***, 0, set.Cardinality())
	for elem := range *set ***REMOVED***
		keys = append(keys, elem)
	***REMOVED***

	return keys
***REMOVED***
