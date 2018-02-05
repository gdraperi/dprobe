// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package precis

import "golang.org/x/text/transform"

// Transformer implements the transform.Transformer interface.
type Transformer struct ***REMOVED***
	t transform.Transformer
***REMOVED***

// Reset implements the transform.Transformer interface.
func (t Transformer) Reset() ***REMOVED*** t.t.Reset() ***REMOVED***

// Transform implements the transform.Transformer interface.
func (t Transformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	return t.t.Transform(dst, src, atEOF)
***REMOVED***

// Bytes returns a new byte slice with the result of applying t to b.
func (t Transformer) Bytes(b []byte) []byte ***REMOVED***
	b, _, _ = transform.Bytes(t, b)
	return b
***REMOVED***

// String returns a string with the result of applying t to s.
func (t Transformer) String(s string) string ***REMOVED***
	s, _, _ = transform.String(t, s)
	return s
***REMOVED***
