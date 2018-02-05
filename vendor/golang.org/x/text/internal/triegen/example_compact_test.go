// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package triegen_test

import (
	"fmt"
	"io"
	"io/ioutil"

	"golang.org/x/text/internal/triegen"
)

func ExampleCompacter() ***REMOVED***
	t := triegen.NewTrie("root")
	for r := rune(0); r < 10000; r += 64 ***REMOVED***
		t.Insert(r, 0x9015BADA55^uint64(r))
	***REMOVED***
	sz, _ := t.Gen(ioutil.Discard)

	fmt.Printf("Size normal:    %5d\n", sz)

	var c myCompacter
	sz, _ = t.Gen(ioutil.Discard, triegen.Compact(&c))

	fmt.Printf("Size compacted: %5d\n", sz)

	// Output:
	// Size normal:    81344
	// Size compacted:  3224
***REMOVED***

// A myCompacter accepts a block if only the first value is given.
type myCompacter []uint64

func (c *myCompacter) Size(values []uint64) (sz int, ok bool) ***REMOVED***
	for _, v := range values[1:] ***REMOVED***
		if v != 0 ***REMOVED***
			return 0, false
		***REMOVED***
	***REMOVED***
	return 8, true // the size of a uint64
***REMOVED***

func (c *myCompacter) Store(v []uint64) uint32 ***REMOVED***
	x := uint32(len(*c))
	*c = append(*c, v[0])
	return x
***REMOVED***

func (c *myCompacter) Print(w io.Writer) error ***REMOVED***
	fmt.Fprintln(w, "var firstValue = []uint64***REMOVED***")
	for _, v := range *c ***REMOVED***
		fmt.Fprintf(w, "\t%#x,\n", v)
	***REMOVED***
	fmt.Fprintln(w, "***REMOVED***")
	return nil
***REMOVED***

func (c *myCompacter) Handler() string ***REMOVED***
	return "getFirstValue"

	// Where getFirstValue is included along with the generated code:
	// func getFirstValue(n uint32, b byte) uint64 ***REMOVED***
	//     if b == 0x80 ***REMOVED*** // the first continuation byte
	//         return firstValue[n]
	// ***REMOVED***
	//     return 0
	//  ***REMOVED***
***REMOVED***
