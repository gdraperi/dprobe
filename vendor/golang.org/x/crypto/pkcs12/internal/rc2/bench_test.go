// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rc2

import (
	"testing"
)

func BenchmarkEncrypt(b *testing.B) ***REMOVED***
	r, _ := New([]byte***REMOVED***0, 0, 0, 0, 0, 0, 0, 0***REMOVED***, 64)
	b.ResetTimer()
	var src [8]byte
	for i := 0; i < b.N; i++ ***REMOVED***
		r.Encrypt(src[:], src[:])
	***REMOVED***
***REMOVED***

func BenchmarkDecrypt(b *testing.B) ***REMOVED***
	r, _ := New([]byte***REMOVED***0, 0, 0, 0, 0, 0, 0, 0***REMOVED***, 64)
	b.ResetTimer()
	var src [8]byte
	for i := 0; i < b.N; i++ ***REMOVED***
		r.Decrypt(src[:], src[:])
	***REMOVED***
***REMOVED***
