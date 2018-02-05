// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package colltab

// testWeighter is a simple Weighter that returns weights from a user-defined map.
type testWeighter map[string][]Elem

func (t testWeighter) Start(int, []byte) int       ***REMOVED*** return 0 ***REMOVED***
func (t testWeighter) StartString(int, string) int ***REMOVED*** return 0 ***REMOVED***
func (t testWeighter) Domain() []string            ***REMOVED*** return nil ***REMOVED***
func (t testWeighter) Top() uint32                 ***REMOVED*** return 0 ***REMOVED***

// maxContractBytes is the maximum length of any key in the map.
const maxContractBytes = 10

func (t testWeighter) AppendNext(buf []Elem, s []byte) ([]Elem, int) ***REMOVED***
	n := len(s)
	if n > maxContractBytes ***REMOVED***
		n = maxContractBytes
	***REMOVED***
	for i := n; i > 0; i-- ***REMOVED***
		if e, ok := t[string(s[:i])]; ok ***REMOVED***
			return append(buf, e...), i
		***REMOVED***
	***REMOVED***
	panic("incomplete testWeighter: could not find " + string(s))
***REMOVED***

func (t testWeighter) AppendNextString(buf []Elem, s string) ([]Elem, int) ***REMOVED***
	n := len(s)
	if n > maxContractBytes ***REMOVED***
		n = maxContractBytes
	***REMOVED***
	for i := n; i > 0; i-- ***REMOVED***
		if e, ok := t[s[:i]]; ok ***REMOVED***
			return append(buf, e...), i
		***REMOVED***
	***REMOVED***
	panic("incomplete testWeighter: could not find " + s)
***REMOVED***
