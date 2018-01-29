// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package autocert

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/crypto/acme"
)

func TestRenewalNext(t *testing.T) ***REMOVED***
	now := time.Now()
	timeNow = func() time.Time ***REMOVED*** return now ***REMOVED***
	defer func() ***REMOVED*** timeNow = time.Now ***REMOVED***()

	man := &Manager***REMOVED***RenewBefore: 7 * 24 * time.Hour***REMOVED***
	defer man.stopRenew()
	tt := []struct ***REMOVED***
		expiry   time.Time
		min, max time.Duration
	***REMOVED******REMOVED***
		***REMOVED***now.Add(90 * 24 * time.Hour), 83*24*time.Hour - renewJitter, 83 * 24 * time.Hour***REMOVED***,
		***REMOVED***now.Add(time.Hour), 0, 1***REMOVED***,
		***REMOVED***now, 0, 1***REMOVED***,
		***REMOVED***now.Add(-time.Hour), 0, 1***REMOVED***,
	***REMOVED***

	dr := &domainRenewal***REMOVED***m: man***REMOVED***
	for i, test := range tt ***REMOVED***
		next := dr.next(test.expiry)
		if next < test.min || test.max < next ***REMOVED***
			t.Errorf("%d: next = %v; want between %v and %v", i, next, test.min, test.max)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestRenewFromCache(t *testing.T) ***REMOVED***
	const domain = "example.org"

	// ACME CA server stub
	var ca *httptest.Server
	ca = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Replay-Nonce", "nonce")
		if r.Method == "HEAD" ***REMOVED***
			// a nonce request
			return
		***REMOVED***

		switch r.URL.Path ***REMOVED***
		// discovery
		case "/":
			if err := discoTmpl.Execute(w, ca.URL); err != nil ***REMOVED***
				t.Fatalf("discoTmpl: %v", err)
			***REMOVED***
		// client key registration
		case "/new-reg":
			w.Write([]byte("***REMOVED******REMOVED***"))
		// domain authorization
		case "/new-authz":
			w.Header().Set("Location", ca.URL+"/authz/1")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`***REMOVED***"status": "valid"***REMOVED***`))
		// cert request
		case "/new-cert":
			var req struct ***REMOVED***
				CSR string `json:"csr"`
			***REMOVED***
			decodePayload(&req, r.Body)
			b, _ := base64.RawURLEncoding.DecodeString(req.CSR)
			csr, err := x509.ParseCertificateRequest(b)
			if err != nil ***REMOVED***
				t.Fatalf("new-cert: CSR: %v", err)
			***REMOVED***
			der, err := dummyCert(csr.PublicKey, domain)
			if err != nil ***REMOVED***
				t.Fatalf("new-cert: dummyCert: %v", err)
			***REMOVED***
			chainUp := fmt.Sprintf("<%s/ca-cert>; rel=up", ca.URL)
			w.Header().Set("Link", chainUp)
			w.WriteHeader(http.StatusCreated)
			w.Write(der)
		// CA chain cert
		case "/ca-cert":
			der, err := dummyCert(nil, "ca")
			if err != nil ***REMOVED***
				t.Fatalf("ca-cert: dummyCert: %v", err)
			***REMOVED***
			w.Write(der)
		default:
			t.Errorf("unrecognized r.URL.Path: %s", r.URL.Path)
		***REMOVED***
	***REMOVED***))
	defer ca.Close()

	// use EC key to run faster on 386
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	man := &Manager***REMOVED***
		Prompt:      AcceptTOS,
		Cache:       newMemCache(),
		RenewBefore: 24 * time.Hour,
		Client: &acme.Client***REMOVED***
			Key:          key,
			DirectoryURL: ca.URL,
		***REMOVED***,
	***REMOVED***
	defer man.stopRenew()

	// cache an almost expired cert
	now := time.Now()
	cert, err := dateDummyCert(key.Public(), now.Add(-2*time.Hour), now.Add(time.Minute), domain)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	tlscert := &tls.Certificate***REMOVED***PrivateKey: key, Certificate: [][]byte***REMOVED***cert***REMOVED******REMOVED***
	if err := man.cachePut(context.Background(), domain, tlscert); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// veriy the renewal happened
	defer func() ***REMOVED***
		testDidRenewLoop = func(next time.Duration, err error) ***REMOVED******REMOVED***
	***REMOVED***()
	done := make(chan struct***REMOVED******REMOVED***)
	testDidRenewLoop = func(next time.Duration, err error) ***REMOVED***
		defer close(done)
		if err != nil ***REMOVED***
			t.Errorf("testDidRenewLoop: %v", err)
		***REMOVED***
		// Next should be about 90 days:
		// dummyCert creates 90days expiry + account for man.RenewBefore.
		// Previous expiration was within 1 min.
		future := 88 * 24 * time.Hour
		if next < future ***REMOVED***
			t.Errorf("testDidRenewLoop: next = %v; want >= %v", next, future)
		***REMOVED***

		// ensure the new cert is cached
		after := time.Now().Add(future)
		tlscert, err := man.cacheGet(context.Background(), domain)
		if err != nil ***REMOVED***
			t.Fatalf("man.cacheGet: %v", err)
		***REMOVED***
		if !tlscert.Leaf.NotAfter.After(after) ***REMOVED***
			t.Errorf("cache leaf.NotAfter = %v; want > %v", tlscert.Leaf.NotAfter, after)
		***REMOVED***

		// verify the old cert is also replaced in memory
		man.stateMu.Lock()
		defer man.stateMu.Unlock()
		s := man.state[domain]
		if s == nil ***REMOVED***
			t.Fatalf("m.state[%q] is nil", domain)
		***REMOVED***
		tlscert, err = s.tlscert()
		if err != nil ***REMOVED***
			t.Fatalf("s.tlscert: %v", err)
		***REMOVED***
		if !tlscert.Leaf.NotAfter.After(after) ***REMOVED***
			t.Errorf("state leaf.NotAfter = %v; want > %v", tlscert.Leaf.NotAfter, after)
		***REMOVED***
	***REMOVED***

	// trigger renew
	hello := &tls.ClientHelloInfo***REMOVED***ServerName: domain***REMOVED***
	if _, err := man.GetCertificate(hello); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// wait for renew loop
	select ***REMOVED***
	case <-time.After(10 * time.Second):
		t.Fatal("renew took too long to occur")
	case <-done:
	***REMOVED***
***REMOVED***
