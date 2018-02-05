// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"strconv"
	"testing"

	"golang.org/x/text/internal/colltab"
)

type entryTest struct ***REMOVED***
	f   func(in []int) (uint32, error)
	arg []int
	val uint32
***REMOVED***

// makeList returns a list of entries of length n+2, with n normal
// entries plus a leading and trailing anchor.
func makeList(n int) []*entry ***REMOVED***
	es := make([]*entry, n+2)
	weights := []rawCE***REMOVED******REMOVED***w: []int***REMOVED***100, 20, 5, 0***REMOVED******REMOVED******REMOVED***
	for i := range es ***REMOVED***
		runes := []rune***REMOVED***rune(i)***REMOVED***
		es[i] = &entry***REMOVED***
			runes: runes,
			elems: weights,
		***REMOVED***
		weights = nextWeight(colltab.Primary, weights)
	***REMOVED***
	for i := 1; i < len(es); i++ ***REMOVED***
		es[i-1].next = es[i]
		es[i].prev = es[i-1]
		_, es[i-1].level = compareWeights(es[i-1].elems, es[i].elems)
	***REMOVED***
	es[0].exclude = true
	es[0].logical = firstAnchor
	es[len(es)-1].exclude = true
	es[len(es)-1].logical = lastAnchor
	return es
***REMOVED***

func TestNextIndexed(t *testing.T) ***REMOVED***
	const n = 5
	es := makeList(n)
	for i := int64(0); i < 1<<n; i++ ***REMOVED***
		mask := strconv.FormatInt(i+(1<<n), 2)
		for i, c := range mask ***REMOVED***
			es[i].exclude = c == '1'
		***REMOVED***
		e := es[0]
		for i, c := range mask ***REMOVED***
			if c == '0' ***REMOVED***
				e, _ = e.nextIndexed()
				if e != es[i] ***REMOVED***
					t.Errorf("%d: expected entry %d; found %d", i, es[i].elems, e.elems)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if e, _ = e.nextIndexed(); e != nil ***REMOVED***
			t.Errorf("%d: expected nil entry; found %d", i, e.elems)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestRemove(t *testing.T) ***REMOVED***
	const n = 5
	for i := int64(0); i < 1<<n; i++ ***REMOVED***
		es := makeList(n)
		mask := strconv.FormatInt(i+(1<<n), 2)
		for i, c := range mask ***REMOVED***
			if c == '0' ***REMOVED***
				es[i].remove()
			***REMOVED***
		***REMOVED***
		e := es[0]
		for i, c := range mask ***REMOVED***
			if c == '1' ***REMOVED***
				if e != es[i] ***REMOVED***
					t.Errorf("%d: expected entry %d; found %d", i, es[i].elems, e.elems)
				***REMOVED***
				e, _ = e.nextIndexed()
			***REMOVED***
		***REMOVED***
		if e != nil ***REMOVED***
			t.Errorf("%d: expected nil entry; found %d", i, e.elems)
		***REMOVED***
	***REMOVED***
***REMOVED***

// nextPerm generates the next permutation of the array.  The starting
// permutation is assumed to be a list of integers sorted in increasing order.
// It returns false if there are no more permuations left.
func nextPerm(a []int) bool ***REMOVED***
	i := len(a) - 2
	for ; i >= 0; i-- ***REMOVED***
		if a[i] < a[i+1] ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	if i < 0 ***REMOVED***
		return false
	***REMOVED***
	for j := len(a) - 1; j >= i; j-- ***REMOVED***
		if a[j] > a[i] ***REMOVED***
			a[i], a[j] = a[j], a[i]
			break
		***REMOVED***
	***REMOVED***
	for j := i + 1; j < (len(a)+i+1)/2; j++ ***REMOVED***
		a[j], a[len(a)+i-j] = a[len(a)+i-j], a[j]
	***REMOVED***
	return true
***REMOVED***

func TestInsertAfter(t *testing.T) ***REMOVED***
	const n = 5
	orig := makeList(n)
	perm := make([]int, n)
	for i := range perm ***REMOVED***
		perm[i] = i + 1
	***REMOVED***
	for ok := true; ok; ok = nextPerm(perm) ***REMOVED***
		es := makeList(n)
		last := es[0]
		for _, i := range perm ***REMOVED***
			last.insertAfter(es[i])
			last = es[i]
		***REMOVED***
		for _, e := range es ***REMOVED***
			e.elems = es[0].elems
		***REMOVED***
		e := es[0]
		for _, i := range perm ***REMOVED***
			e, _ = e.nextIndexed()
			if e.runes[0] != orig[i].runes[0] ***REMOVED***
				t.Errorf("%d:%d: expected entry %X; found %X", perm, i, orig[i].runes, e.runes)
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestInsertBefore(t *testing.T) ***REMOVED***
	const n = 5
	orig := makeList(n)
	perm := make([]int, n)
	for i := range perm ***REMOVED***
		perm[i] = i + 1
	***REMOVED***
	for ok := true; ok; ok = nextPerm(perm) ***REMOVED***
		es := makeList(n)
		last := es[len(es)-1]
		for _, i := range perm ***REMOVED***
			last.insertBefore(es[i])
			last = es[i]
		***REMOVED***
		for _, e := range es ***REMOVED***
			e.elems = es[0].elems
		***REMOVED***
		e := es[0]
		for i := n - 1; i >= 0; i-- ***REMOVED***
			e, _ = e.nextIndexed()
			if e.runes[0] != rune(perm[i]) ***REMOVED***
				t.Errorf("%d:%d: expected entry %X; found %X", perm, i, orig[i].runes, e.runes)
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

type entryLessTest struct ***REMOVED***
	a, b *entry
	res  bool
***REMOVED***

var (
	w1 = []rawCE***REMOVED******REMOVED***w: []int***REMOVED***100, 20, 5, 5***REMOVED******REMOVED******REMOVED***
	w2 = []rawCE***REMOVED******REMOVED***w: []int***REMOVED***101, 20, 5, 5***REMOVED******REMOVED******REMOVED***
)

var entryLessTests = []entryLessTest***REMOVED***
	***REMOVED***&entry***REMOVED***str: "a", elems: w1***REMOVED***,
		&entry***REMOVED***str: "a", elems: w1***REMOVED***,
		false,
	***REMOVED***,
	***REMOVED***&entry***REMOVED***str: "a", elems: w1***REMOVED***,
		&entry***REMOVED***str: "a", elems: w2***REMOVED***,
		true,
	***REMOVED***,
	***REMOVED***&entry***REMOVED***str: "a", elems: w1***REMOVED***,
		&entry***REMOVED***str: "b", elems: w1***REMOVED***,
		true,
	***REMOVED***,
	***REMOVED***&entry***REMOVED***str: "a", elems: w2***REMOVED***,
		&entry***REMOVED***str: "a", elems: w1***REMOVED***,
		false,
	***REMOVED***,
	***REMOVED***&entry***REMOVED***str: "c", elems: w1***REMOVED***,
		&entry***REMOVED***str: "b", elems: w1***REMOVED***,
		false,
	***REMOVED***,
	***REMOVED***&entry***REMOVED***str: "a", elems: w1, logical: firstAnchor***REMOVED***,
		&entry***REMOVED***str: "a", elems: w1***REMOVED***,
		true,
	***REMOVED***,
	***REMOVED***&entry***REMOVED***str: "a", elems: w1***REMOVED***,
		&entry***REMOVED***str: "b", elems: w1, logical: firstAnchor***REMOVED***,
		false,
	***REMOVED***,
	***REMOVED***&entry***REMOVED***str: "b", elems: w1***REMOVED***,
		&entry***REMOVED***str: "a", elems: w1, logical: lastAnchor***REMOVED***,
		true,
	***REMOVED***,
	***REMOVED***&entry***REMOVED***str: "a", elems: w1, logical: lastAnchor***REMOVED***,
		&entry***REMOVED***str: "c", elems: w1***REMOVED***,
		false,
	***REMOVED***,
***REMOVED***

func TestEntryLess(t *testing.T) ***REMOVED***
	for i, tt := range entryLessTests ***REMOVED***
		if res := entryLess(tt.a, tt.b); res != tt.res ***REMOVED***
			t.Errorf("%d: was %v; want %v", i, res, tt.res)
		***REMOVED***
	***REMOVED***
***REMOVED***
