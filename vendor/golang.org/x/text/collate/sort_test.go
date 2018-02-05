// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collate_test

import (
	"fmt"
	"testing"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

func ExampleCollator_Strings() ***REMOVED***
	c := collate.New(language.Und)
	strings := []string***REMOVED***
		"ad",
		"ab",
		"äb",
		"ac",
	***REMOVED***
	c.SortStrings(strings)
	fmt.Println(strings)
	// Output: [ab äb ac ad]
***REMOVED***

type sorter []string

func (s sorter) Len() int ***REMOVED***
	return len(s)
***REMOVED***

func (s sorter) Swap(i, j int) ***REMOVED***
	s[j], s[i] = s[i], s[j]
***REMOVED***

func (s sorter) Bytes(i int) []byte ***REMOVED***
	return []byte(s[i])
***REMOVED***

func TestSort(t *testing.T) ***REMOVED***
	c := collate.New(language.English)
	strings := []string***REMOVED***
		"bcd",
		"abc",
		"ddd",
	***REMOVED***
	c.Sort(sorter(strings))
	res := fmt.Sprint(strings)
	want := "[abc bcd ddd]"
	if res != want ***REMOVED***
		t.Errorf("found %s; want %s", res, want)
	***REMOVED***
***REMOVED***
