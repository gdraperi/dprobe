// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.8

package http2

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"
)

// Tests that http2.Server.IdleTimeout is initialized from
// http.Server.***REMOVED***Idle,Read***REMOVED***Timeout. http.Server.IdleTimeout was
// added in Go 1.8.
func TestConfigureServerIdleTimeout_Go18(t *testing.T) ***REMOVED***
	const timeout = 5 * time.Second
	const notThisOne = 1 * time.Second

	// With a zero http2.Server, verify that it copies IdleTimeout:
	***REMOVED***
		s1 := &http.Server***REMOVED***
			IdleTimeout: timeout,
			ReadTimeout: notThisOne,
		***REMOVED***
		s2 := &Server***REMOVED******REMOVED***
		if err := ConfigureServer(s1, s2); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if s2.IdleTimeout != timeout ***REMOVED***
			t.Errorf("s2.IdleTimeout = %v; want %v", s2.IdleTimeout, timeout)
		***REMOVED***
	***REMOVED***

	// And that it falls back to ReadTimeout:
	***REMOVED***
		s1 := &http.Server***REMOVED***
			ReadTimeout: timeout,
		***REMOVED***
		s2 := &Server***REMOVED******REMOVED***
		if err := ConfigureServer(s1, s2); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if s2.IdleTimeout != timeout ***REMOVED***
			t.Errorf("s2.IdleTimeout = %v; want %v", s2.IdleTimeout, timeout)
		***REMOVED***
	***REMOVED***

	// Verify that s1's IdleTimeout doesn't overwrite an existing setting:
	***REMOVED***
		s1 := &http.Server***REMOVED***
			IdleTimeout: notThisOne,
		***REMOVED***
		s2 := &Server***REMOVED***
			IdleTimeout: timeout,
		***REMOVED***
		if err := ConfigureServer(s1, s2); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if s2.IdleTimeout != timeout ***REMOVED***
			t.Errorf("s2.IdleTimeout = %v; want %v", s2.IdleTimeout, timeout)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestCertClone(t *testing.T) ***REMOVED***
	c := &tls.Config***REMOVED***
		GetClientCertificate: func(*tls.CertificateRequestInfo) (*tls.Certificate, error) ***REMOVED***
			panic("shouldn't be called")
		***REMOVED***,
	***REMOVED***
	c2 := cloneTLSConfig(c)
	if c2.GetClientCertificate == nil ***REMOVED***
		t.Error("GetClientCertificate is nil")
	***REMOVED***
***REMOVED***
