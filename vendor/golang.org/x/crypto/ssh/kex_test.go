// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

// Key exchange tests.

import (
	"crypto/rand"
	"reflect"
	"testing"
)

func TestKexes(t *testing.T) ***REMOVED***
	type kexResultErr struct ***REMOVED***
		result *kexResult
		err    error
	***REMOVED***

	for name, kex := range kexAlgoMap ***REMOVED***
		a, b := memPipe()

		s := make(chan kexResultErr, 1)
		c := make(chan kexResultErr, 1)
		var magics handshakeMagics
		go func() ***REMOVED***
			r, e := kex.Client(a, rand.Reader, &magics)
			a.Close()
			c <- kexResultErr***REMOVED***r, e***REMOVED***
		***REMOVED***()
		go func() ***REMOVED***
			r, e := kex.Server(b, rand.Reader, &magics, testSigners["ecdsa"])
			b.Close()
			s <- kexResultErr***REMOVED***r, e***REMOVED***
		***REMOVED***()

		clientRes := <-c
		serverRes := <-s
		if clientRes.err != nil ***REMOVED***
			t.Errorf("client: %v", clientRes.err)
		***REMOVED***
		if serverRes.err != nil ***REMOVED***
			t.Errorf("server: %v", serverRes.err)
		***REMOVED***
		if !reflect.DeepEqual(clientRes.result, serverRes.result) ***REMOVED***
			t.Errorf("kex %q: mismatch %#v, %#v", name, clientRes.result, serverRes.result)
		***REMOVED***
	***REMOVED***
***REMOVED***
