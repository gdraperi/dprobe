// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package norm

import (
	"bytes"
	"fmt"
	"testing"
)

var bufSizes = []int***REMOVED***1, 2, 3, 4, 5, 6, 7, 8, 100, 101, 102, 103, 4000, 4001, 4002, 4003***REMOVED***

func readFunc(size int) appendFunc ***REMOVED***
	return func(f Form, out []byte, s string) []byte ***REMOVED***
		out = append(out, s...)
		r := f.Reader(bytes.NewBuffer(out))
		buf := make([]byte, size)
		result := []byte***REMOVED******REMOVED***
		for n, err := 0, error(nil); err == nil; ***REMOVED***
			n, err = r.Read(buf)
			result = append(result, buf[:n]...)
		***REMOVED***
		return result
	***REMOVED***
***REMOVED***

func TestReader(t *testing.T) ***REMOVED***
	for _, s := range bufSizes ***REMOVED***
		name := fmt.Sprintf("TestReader%d", s)
		runNormTests(t, name, readFunc(s))
	***REMOVED***
***REMOVED***

func writeFunc(size int) appendFunc ***REMOVED***
	return func(f Form, out []byte, s string) []byte ***REMOVED***
		in := append(out, s...)
		result := new(bytes.Buffer)
		w := f.Writer(result)
		buf := make([]byte, size)
		for n := 0; len(in) > 0; in = in[n:] ***REMOVED***
			n = copy(buf, in)
			_, _ = w.Write(buf[:n])
		***REMOVED***
		w.Close()
		return result.Bytes()
	***REMOVED***
***REMOVED***

func TestWriter(t *testing.T) ***REMOVED***
	for _, s := range bufSizes ***REMOVED***
		name := fmt.Sprintf("TestWriter%d", s)
		runNormTests(t, name, writeFunc(s))
	***REMOVED***
***REMOVED***
