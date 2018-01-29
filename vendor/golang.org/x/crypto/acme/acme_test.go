// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acme

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

// Decodes a JWS-encoded request and unmarshals the decoded JSON into a provided
// interface.
func decodeJWSRequest(t *testing.T, v interface***REMOVED******REMOVED***, r *http.Request) ***REMOVED***
	// Decode request
	var req struct***REMOVED*** Payload string ***REMOVED***
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	payload, err := base64.RawURLEncoding.DecodeString(req.Payload)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	err = json.Unmarshal(payload, v)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

type jwsHead struct ***REMOVED***
	Alg   string
	Nonce string
	JWK   map[string]string `json:"jwk"`
***REMOVED***

func decodeJWSHead(r *http.Request) (*jwsHead, error) ***REMOVED***
	var req struct***REMOVED*** Protected string ***REMOVED***
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	b, err := base64.RawURLEncoding.DecodeString(req.Protected)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var head jwsHead
	if err := json.Unmarshal(b, &head); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &head, nil
***REMOVED***

func TestDiscover(t *testing.T) ***REMOVED***
	const (
		reg    = "https://example.com/acme/new-reg"
		authz  = "https://example.com/acme/new-authz"
		cert   = "https://example.com/acme/new-cert"
		revoke = "https://example.com/acme/revoke-cert"
	)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `***REMOVED***
			"new-reg": %q,
			"new-authz": %q,
			"new-cert": %q,
			"revoke-cert": %q
		***REMOVED***`, reg, authz, cert, revoke)
	***REMOVED***))
	defer ts.Close()
	c := Client***REMOVED***DirectoryURL: ts.URL***REMOVED***
	dir, err := c.Discover(context.Background())
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if dir.RegURL != reg ***REMOVED***
		t.Errorf("dir.RegURL = %q; want %q", dir.RegURL, reg)
	***REMOVED***
	if dir.AuthzURL != authz ***REMOVED***
		t.Errorf("dir.AuthzURL = %q; want %q", dir.AuthzURL, authz)
	***REMOVED***
	if dir.CertURL != cert ***REMOVED***
		t.Errorf("dir.CertURL = %q; want %q", dir.CertURL, cert)
	***REMOVED***
	if dir.RevokeURL != revoke ***REMOVED***
		t.Errorf("dir.RevokeURL = %q; want %q", dir.RevokeURL, revoke)
	***REMOVED***
***REMOVED***

func TestRegister(t *testing.T) ***REMOVED***
	contacts := []string***REMOVED***"mailto:admin@example.com"***REMOVED***

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method == "HEAD" ***REMOVED***
			w.Header().Set("Replay-Nonce", "test-nonce")
			return
		***REMOVED***
		if r.Method != "POST" ***REMOVED***
			t.Errorf("r.Method = %q; want POST", r.Method)
		***REMOVED***

		var j struct ***REMOVED***
			Resource  string
			Contact   []string
			Agreement string
		***REMOVED***
		decodeJWSRequest(t, &j, r)

		// Test request
		if j.Resource != "new-reg" ***REMOVED***
			t.Errorf("j.Resource = %q; want new-reg", j.Resource)
		***REMOVED***
		if !reflect.DeepEqual(j.Contact, contacts) ***REMOVED***
			t.Errorf("j.Contact = %v; want %v", j.Contact, contacts)
		***REMOVED***

		w.Header().Set("Location", "https://ca.tld/acme/reg/1")
		w.Header().Set("Link", `<https://ca.tld/acme/new-authz>;rel="next"`)
		w.Header().Add("Link", `<https://ca.tld/acme/recover-reg>;rel="recover"`)
		w.Header().Add("Link", `<https://ca.tld/acme/terms>;rel="terms-of-service"`)
		w.WriteHeader(http.StatusCreated)
		b, _ := json.Marshal(contacts)
		fmt.Fprintf(w, `***REMOVED***"contact": %s***REMOVED***`, b)
	***REMOVED***))
	defer ts.Close()

	prompt := func(url string) bool ***REMOVED***
		const terms = "https://ca.tld/acme/terms"
		if url != terms ***REMOVED***
			t.Errorf("prompt url = %q; want %q", url, terms)
		***REMOVED***
		return false
	***REMOVED***

	c := Client***REMOVED***Key: testKeyEC, dir: &Directory***REMOVED***RegURL: ts.URL***REMOVED******REMOVED***
	a := &Account***REMOVED***Contact: contacts***REMOVED***
	var err error
	if a, err = c.Register(context.Background(), a, prompt); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if a.URI != "https://ca.tld/acme/reg/1" ***REMOVED***
		t.Errorf("a.URI = %q; want https://ca.tld/acme/reg/1", a.URI)
	***REMOVED***
	if a.Authz != "https://ca.tld/acme/new-authz" ***REMOVED***
		t.Errorf("a.Authz = %q; want https://ca.tld/acme/new-authz", a.Authz)
	***REMOVED***
	if a.CurrentTerms != "https://ca.tld/acme/terms" ***REMOVED***
		t.Errorf("a.CurrentTerms = %q; want https://ca.tld/acme/terms", a.CurrentTerms)
	***REMOVED***
	if !reflect.DeepEqual(a.Contact, contacts) ***REMOVED***
		t.Errorf("a.Contact = %v; want %v", a.Contact, contacts)
	***REMOVED***
***REMOVED***

func TestUpdateReg(t *testing.T) ***REMOVED***
	const terms = "https://ca.tld/acme/terms"
	contacts := []string***REMOVED***"mailto:admin@example.com"***REMOVED***

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method == "HEAD" ***REMOVED***
			w.Header().Set("Replay-Nonce", "test-nonce")
			return
		***REMOVED***
		if r.Method != "POST" ***REMOVED***
			t.Errorf("r.Method = %q; want POST", r.Method)
		***REMOVED***

		var j struct ***REMOVED***
			Resource  string
			Contact   []string
			Agreement string
		***REMOVED***
		decodeJWSRequest(t, &j, r)

		// Test request
		if j.Resource != "reg" ***REMOVED***
			t.Errorf("j.Resource = %q; want reg", j.Resource)
		***REMOVED***
		if j.Agreement != terms ***REMOVED***
			t.Errorf("j.Agreement = %q; want %q", j.Agreement, terms)
		***REMOVED***
		if !reflect.DeepEqual(j.Contact, contacts) ***REMOVED***
			t.Errorf("j.Contact = %v; want %v", j.Contact, contacts)
		***REMOVED***

		w.Header().Set("Link", `<https://ca.tld/acme/new-authz>;rel="next"`)
		w.Header().Add("Link", `<https://ca.tld/acme/recover-reg>;rel="recover"`)
		w.Header().Add("Link", fmt.Sprintf(`<%s>;rel="terms-of-service"`, terms))
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(contacts)
		fmt.Fprintf(w, `***REMOVED***"contact":%s, "agreement":%q***REMOVED***`, b, terms)
	***REMOVED***))
	defer ts.Close()

	c := Client***REMOVED***Key: testKeyEC***REMOVED***
	a := &Account***REMOVED***URI: ts.URL, Contact: contacts, AgreedTerms: terms***REMOVED***
	var err error
	if a, err = c.UpdateReg(context.Background(), a); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if a.Authz != "https://ca.tld/acme/new-authz" ***REMOVED***
		t.Errorf("a.Authz = %q; want https://ca.tld/acme/new-authz", a.Authz)
	***REMOVED***
	if a.AgreedTerms != terms ***REMOVED***
		t.Errorf("a.AgreedTerms = %q; want %q", a.AgreedTerms, terms)
	***REMOVED***
	if a.CurrentTerms != terms ***REMOVED***
		t.Errorf("a.CurrentTerms = %q; want %q", a.CurrentTerms, terms)
	***REMOVED***
	if a.URI != ts.URL ***REMOVED***
		t.Errorf("a.URI = %q; want %q", a.URI, ts.URL)
	***REMOVED***
***REMOVED***

func TestGetReg(t *testing.T) ***REMOVED***
	const terms = "https://ca.tld/acme/terms"
	const newTerms = "https://ca.tld/acme/new-terms"
	contacts := []string***REMOVED***"mailto:admin@example.com"***REMOVED***

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method == "HEAD" ***REMOVED***
			w.Header().Set("Replay-Nonce", "test-nonce")
			return
		***REMOVED***
		if r.Method != "POST" ***REMOVED***
			t.Errorf("r.Method = %q; want POST", r.Method)
		***REMOVED***

		var j struct ***REMOVED***
			Resource  string
			Contact   []string
			Agreement string
		***REMOVED***
		decodeJWSRequest(t, &j, r)

		// Test request
		if j.Resource != "reg" ***REMOVED***
			t.Errorf("j.Resource = %q; want reg", j.Resource)
		***REMOVED***
		if len(j.Contact) != 0 ***REMOVED***
			t.Errorf("j.Contact = %v", j.Contact)
		***REMOVED***
		if j.Agreement != "" ***REMOVED***
			t.Errorf("j.Agreement = %q", j.Agreement)
		***REMOVED***

		w.Header().Set("Link", `<https://ca.tld/acme/new-authz>;rel="next"`)
		w.Header().Add("Link", `<https://ca.tld/acme/recover-reg>;rel="recover"`)
		w.Header().Add("Link", fmt.Sprintf(`<%s>;rel="terms-of-service"`, newTerms))
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(contacts)
		fmt.Fprintf(w, `***REMOVED***"contact":%s, "agreement":%q***REMOVED***`, b, terms)
	***REMOVED***))
	defer ts.Close()

	c := Client***REMOVED***Key: testKeyEC***REMOVED***
	a, err := c.GetReg(context.Background(), ts.URL)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if a.Authz != "https://ca.tld/acme/new-authz" ***REMOVED***
		t.Errorf("a.AuthzURL = %q; want https://ca.tld/acme/new-authz", a.Authz)
	***REMOVED***
	if a.AgreedTerms != terms ***REMOVED***
		t.Errorf("a.AgreedTerms = %q; want %q", a.AgreedTerms, terms)
	***REMOVED***
	if a.CurrentTerms != newTerms ***REMOVED***
		t.Errorf("a.CurrentTerms = %q; want %q", a.CurrentTerms, newTerms)
	***REMOVED***
	if a.URI != ts.URL ***REMOVED***
		t.Errorf("a.URI = %q; want %q", a.URI, ts.URL)
	***REMOVED***
***REMOVED***

func TestAuthorize(t *testing.T) ***REMOVED***
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method == "HEAD" ***REMOVED***
			w.Header().Set("Replay-Nonce", "test-nonce")
			return
		***REMOVED***
		if r.Method != "POST" ***REMOVED***
			t.Errorf("r.Method = %q; want POST", r.Method)
		***REMOVED***

		var j struct ***REMOVED***
			Resource   string
			Identifier struct ***REMOVED***
				Type  string
				Value string
			***REMOVED***
		***REMOVED***
		decodeJWSRequest(t, &j, r)

		// Test request
		if j.Resource != "new-authz" ***REMOVED***
			t.Errorf("j.Resource = %q; want new-authz", j.Resource)
		***REMOVED***
		if j.Identifier.Type != "dns" ***REMOVED***
			t.Errorf("j.Identifier.Type = %q; want dns", j.Identifier.Type)
		***REMOVED***
		if j.Identifier.Value != "example.com" ***REMOVED***
			t.Errorf("j.Identifier.Value = %q; want example.com", j.Identifier.Value)
		***REMOVED***

		w.Header().Set("Location", "https://ca.tld/acme/auth/1")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `***REMOVED***
			"identifier": ***REMOVED***"type":"dns","value":"example.com"***REMOVED***,
			"status":"pending",
			"challenges":[
				***REMOVED***
					"type":"http-01",
					"status":"pending",
					"uri":"https://ca.tld/acme/challenge/publickey/id1",
					"token":"token1"
				***REMOVED***,
				***REMOVED***
					"type":"tls-sni-01",
					"status":"pending",
					"uri":"https://ca.tld/acme/challenge/publickey/id2",
					"token":"token2"
				***REMOVED***
			],
			"combinations":[[0],[1]]***REMOVED***`)
	***REMOVED***))
	defer ts.Close()

	cl := Client***REMOVED***Key: testKeyEC, dir: &Directory***REMOVED***AuthzURL: ts.URL***REMOVED******REMOVED***
	auth, err := cl.Authorize(context.Background(), "example.com")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if auth.URI != "https://ca.tld/acme/auth/1" ***REMOVED***
		t.Errorf("URI = %q; want https://ca.tld/acme/auth/1", auth.URI)
	***REMOVED***
	if auth.Status != "pending" ***REMOVED***
		t.Errorf("Status = %q; want pending", auth.Status)
	***REMOVED***
	if auth.Identifier.Type != "dns" ***REMOVED***
		t.Errorf("Identifier.Type = %q; want dns", auth.Identifier.Type)
	***REMOVED***
	if auth.Identifier.Value != "example.com" ***REMOVED***
		t.Errorf("Identifier.Value = %q; want example.com", auth.Identifier.Value)
	***REMOVED***

	if n := len(auth.Challenges); n != 2 ***REMOVED***
		t.Fatalf("len(auth.Challenges) = %d; want 2", n)
	***REMOVED***

	c := auth.Challenges[0]
	if c.Type != "http-01" ***REMOVED***
		t.Errorf("c.Type = %q; want http-01", c.Type)
	***REMOVED***
	if c.URI != "https://ca.tld/acme/challenge/publickey/id1" ***REMOVED***
		t.Errorf("c.URI = %q; want https://ca.tld/acme/challenge/publickey/id1", c.URI)
	***REMOVED***
	if c.Token != "token1" ***REMOVED***
		t.Errorf("c.Token = %q; want token1", c.Token)
	***REMOVED***

	c = auth.Challenges[1]
	if c.Type != "tls-sni-01" ***REMOVED***
		t.Errorf("c.Type = %q; want tls-sni-01", c.Type)
	***REMOVED***
	if c.URI != "https://ca.tld/acme/challenge/publickey/id2" ***REMOVED***
		t.Errorf("c.URI = %q; want https://ca.tld/acme/challenge/publickey/id2", c.URI)
	***REMOVED***
	if c.Token != "token2" ***REMOVED***
		t.Errorf("c.Token = %q; want token2", c.Token)
	***REMOVED***

	combs := [][]int***REMOVED******REMOVED***0***REMOVED***, ***REMOVED***1***REMOVED******REMOVED***
	if !reflect.DeepEqual(auth.Combinations, combs) ***REMOVED***
		t.Errorf("auth.Combinations: %+v\nwant: %+v\n", auth.Combinations, combs)
	***REMOVED***
***REMOVED***

func TestAuthorizeValid(t *testing.T) ***REMOVED***
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method == "HEAD" ***REMOVED***
			w.Header().Set("Replay-Nonce", "nonce")
			return
		***REMOVED***
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`***REMOVED***"status":"valid"***REMOVED***`))
	***REMOVED***))
	defer ts.Close()
	client := Client***REMOVED***Key: testKey, dir: &Directory***REMOVED***AuthzURL: ts.URL***REMOVED******REMOVED***
	_, err := client.Authorize(context.Background(), "example.com")
	if err != nil ***REMOVED***
		t.Errorf("err = %v", err)
	***REMOVED***
***REMOVED***

func TestGetAuthorization(t *testing.T) ***REMOVED***
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method != "GET" ***REMOVED***
			t.Errorf("r.Method = %q; want GET", r.Method)
		***REMOVED***

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `***REMOVED***
			"identifier": ***REMOVED***"type":"dns","value":"example.com"***REMOVED***,
			"status":"pending",
			"challenges":[
				***REMOVED***
					"type":"http-01",
					"status":"pending",
					"uri":"https://ca.tld/acme/challenge/publickey/id1",
					"token":"token1"
				***REMOVED***,
				***REMOVED***
					"type":"tls-sni-01",
					"status":"pending",
					"uri":"https://ca.tld/acme/challenge/publickey/id2",
					"token":"token2"
				***REMOVED***
			],
			"combinations":[[0],[1]]***REMOVED***`)
	***REMOVED***))
	defer ts.Close()

	cl := Client***REMOVED***Key: testKeyEC***REMOVED***
	auth, err := cl.GetAuthorization(context.Background(), ts.URL)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if auth.Status != "pending" ***REMOVED***
		t.Errorf("Status = %q; want pending", auth.Status)
	***REMOVED***
	if auth.Identifier.Type != "dns" ***REMOVED***
		t.Errorf("Identifier.Type = %q; want dns", auth.Identifier.Type)
	***REMOVED***
	if auth.Identifier.Value != "example.com" ***REMOVED***
		t.Errorf("Identifier.Value = %q; want example.com", auth.Identifier.Value)
	***REMOVED***

	if n := len(auth.Challenges); n != 2 ***REMOVED***
		t.Fatalf("len(set.Challenges) = %d; want 2", n)
	***REMOVED***

	c := auth.Challenges[0]
	if c.Type != "http-01" ***REMOVED***
		t.Errorf("c.Type = %q; want http-01", c.Type)
	***REMOVED***
	if c.URI != "https://ca.tld/acme/challenge/publickey/id1" ***REMOVED***
		t.Errorf("c.URI = %q; want https://ca.tld/acme/challenge/publickey/id1", c.URI)
	***REMOVED***
	if c.Token != "token1" ***REMOVED***
		t.Errorf("c.Token = %q; want token1", c.Token)
	***REMOVED***

	c = auth.Challenges[1]
	if c.Type != "tls-sni-01" ***REMOVED***
		t.Errorf("c.Type = %q; want tls-sni-01", c.Type)
	***REMOVED***
	if c.URI != "https://ca.tld/acme/challenge/publickey/id2" ***REMOVED***
		t.Errorf("c.URI = %q; want https://ca.tld/acme/challenge/publickey/id2", c.URI)
	***REMOVED***
	if c.Token != "token2" ***REMOVED***
		t.Errorf("c.Token = %q; want token2", c.Token)
	***REMOVED***

	combs := [][]int***REMOVED******REMOVED***0***REMOVED***, ***REMOVED***1***REMOVED******REMOVED***
	if !reflect.DeepEqual(auth.Combinations, combs) ***REMOVED***
		t.Errorf("auth.Combinations: %+v\nwant: %+v\n", auth.Combinations, combs)
	***REMOVED***
***REMOVED***

func TestWaitAuthorization(t *testing.T) ***REMOVED***
	var count int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		count++
		w.Header().Set("Retry-After", "0")
		if count > 1 ***REMOVED***
			fmt.Fprintf(w, `***REMOVED***"status":"valid"***REMOVED***`)
			return
		***REMOVED***
		fmt.Fprintf(w, `***REMOVED***"status":"pending"***REMOVED***`)
	***REMOVED***))
	defer ts.Close()

	type res struct ***REMOVED***
		authz *Authorization
		err   error
	***REMOVED***
	done := make(chan res)
	defer close(done)
	go func() ***REMOVED***
		var client Client
		a, err := client.WaitAuthorization(context.Background(), ts.URL)
		done <- res***REMOVED***a, err***REMOVED***
	***REMOVED***()

	select ***REMOVED***
	case <-time.After(5 * time.Second):
		t.Fatal("WaitAuthz took too long to return")
	case res := <-done:
		if res.err != nil ***REMOVED***
			t.Fatalf("res.err =  %v", res.err)
		***REMOVED***
		if res.authz == nil ***REMOVED***
			t.Fatal("res.authz is nil")
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestWaitAuthorizationInvalid(t *testing.T) ***REMOVED***
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		fmt.Fprintf(w, `***REMOVED***"status":"invalid"***REMOVED***`)
	***REMOVED***))
	defer ts.Close()

	res := make(chan error)
	defer close(res)
	go func() ***REMOVED***
		var client Client
		_, err := client.WaitAuthorization(context.Background(), ts.URL)
		res <- err
	***REMOVED***()

	select ***REMOVED***
	case <-time.After(3 * time.Second):
		t.Fatal("WaitAuthz took too long to return")
	case err := <-res:
		if err == nil ***REMOVED***
			t.Error("err is nil")
		***REMOVED***
		if _, ok := err.(*AuthorizationError); !ok ***REMOVED***
			t.Errorf("err is %T; want *AuthorizationError", err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestWaitAuthorizationCancel(t *testing.T) ***REMOVED***
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Retry-After", "60")
		fmt.Fprintf(w, `***REMOVED***"status":"pending"***REMOVED***`)
	***REMOVED***))
	defer ts.Close()

	res := make(chan error)
	defer close(res)
	go func() ***REMOVED***
		var client Client
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		_, err := client.WaitAuthorization(ctx, ts.URL)
		res <- err
	***REMOVED***()

	select ***REMOVED***
	case <-time.After(time.Second):
		t.Fatal("WaitAuthz took too long to return")
	case err := <-res:
		if err == nil ***REMOVED***
			t.Error("err is nil")
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestRevokeAuthorization(t *testing.T) ***REMOVED***
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method == "HEAD" ***REMOVED***
			w.Header().Set("Replay-Nonce", "nonce")
			return
		***REMOVED***
		switch r.URL.Path ***REMOVED***
		case "/1":
			var req struct ***REMOVED***
				Resource string
				Status   string
				Delete   bool
			***REMOVED***
			decodeJWSRequest(t, &req, r)
			if req.Resource != "authz" ***REMOVED***
				t.Errorf("req.Resource = %q; want authz", req.Resource)
			***REMOVED***
			if req.Status != "deactivated" ***REMOVED***
				t.Errorf("req.Status = %q; want deactivated", req.Status)
			***REMOVED***
			if !req.Delete ***REMOVED***
				t.Errorf("req.Delete is false")
			***REMOVED***
		case "/2":
			w.WriteHeader(http.StatusInternalServerError)
		***REMOVED***
	***REMOVED***))
	defer ts.Close()
	client := &Client***REMOVED***Key: testKey***REMOVED***
	ctx := context.Background()
	if err := client.RevokeAuthorization(ctx, ts.URL+"/1"); err != nil ***REMOVED***
		t.Errorf("err = %v", err)
	***REMOVED***
	if client.RevokeAuthorization(ctx, ts.URL+"/2") == nil ***REMOVED***
		t.Error("nil error")
	***REMOVED***
***REMOVED***

func TestPollChallenge(t *testing.T) ***REMOVED***
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method != "GET" ***REMOVED***
			t.Errorf("r.Method = %q; want GET", r.Method)
		***REMOVED***

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `***REMOVED***
			"type":"http-01",
			"status":"pending",
			"uri":"https://ca.tld/acme/challenge/publickey/id1",
			"token":"token1"***REMOVED***`)
	***REMOVED***))
	defer ts.Close()

	cl := Client***REMOVED***Key: testKeyEC***REMOVED***
	chall, err := cl.GetChallenge(context.Background(), ts.URL)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if chall.Status != "pending" ***REMOVED***
		t.Errorf("Status = %q; want pending", chall.Status)
	***REMOVED***
	if chall.Type != "http-01" ***REMOVED***
		t.Errorf("c.Type = %q; want http-01", chall.Type)
	***REMOVED***
	if chall.URI != "https://ca.tld/acme/challenge/publickey/id1" ***REMOVED***
		t.Errorf("c.URI = %q; want https://ca.tld/acme/challenge/publickey/id1", chall.URI)
	***REMOVED***
	if chall.Token != "token1" ***REMOVED***
		t.Errorf("c.Token = %q; want token1", chall.Token)
	***REMOVED***
***REMOVED***

func TestAcceptChallenge(t *testing.T) ***REMOVED***
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method == "HEAD" ***REMOVED***
			w.Header().Set("Replay-Nonce", "test-nonce")
			return
		***REMOVED***
		if r.Method != "POST" ***REMOVED***
			t.Errorf("r.Method = %q; want POST", r.Method)
		***REMOVED***

		var j struct ***REMOVED***
			Resource string
			Type     string
			Auth     string `json:"keyAuthorization"`
		***REMOVED***
		decodeJWSRequest(t, &j, r)

		// Test request
		if j.Resource != "challenge" ***REMOVED***
			t.Errorf(`resource = %q; want "challenge"`, j.Resource)
		***REMOVED***
		if j.Type != "http-01" ***REMOVED***
			t.Errorf(`type = %q; want "http-01"`, j.Type)
		***REMOVED***
		keyAuth := "token1." + testKeyECThumbprint
		if j.Auth != keyAuth ***REMOVED***
			t.Errorf(`keyAuthorization = %q; want %q`, j.Auth, keyAuth)
		***REMOVED***

		// Respond to request
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, `***REMOVED***
			"type":"http-01",
			"status":"pending",
			"uri":"https://ca.tld/acme/challenge/publickey/id1",
			"token":"token1",
			"keyAuthorization":%q
		***REMOVED***`, keyAuth)
	***REMOVED***))
	defer ts.Close()

	cl := Client***REMOVED***Key: testKeyEC***REMOVED***
	c, err := cl.Accept(context.Background(), &Challenge***REMOVED***
		URI:   ts.URL,
		Token: "token1",
		Type:  "http-01",
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if c.Type != "http-01" ***REMOVED***
		t.Errorf("c.Type = %q; want http-01", c.Type)
	***REMOVED***
	if c.URI != "https://ca.tld/acme/challenge/publickey/id1" ***REMOVED***
		t.Errorf("c.URI = %q; want https://ca.tld/acme/challenge/publickey/id1", c.URI)
	***REMOVED***
	if c.Token != "token1" ***REMOVED***
		t.Errorf("c.Token = %q; want token1", c.Token)
	***REMOVED***
***REMOVED***

func TestNewCert(t *testing.T) ***REMOVED***
	notBefore := time.Now()
	notAfter := notBefore.AddDate(0, 2, 0)
	timeNow = func() time.Time ***REMOVED*** return notBefore ***REMOVED***

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method == "HEAD" ***REMOVED***
			w.Header().Set("Replay-Nonce", "test-nonce")
			return
		***REMOVED***
		if r.Method != "POST" ***REMOVED***
			t.Errorf("r.Method = %q; want POST", r.Method)
		***REMOVED***

		var j struct ***REMOVED***
			Resource  string `json:"resource"`
			CSR       string `json:"csr"`
			NotBefore string `json:"notBefore,omitempty"`
			NotAfter  string `json:"notAfter,omitempty"`
		***REMOVED***
		decodeJWSRequest(t, &j, r)

		// Test request
		if j.Resource != "new-cert" ***REMOVED***
			t.Errorf(`resource = %q; want "new-cert"`, j.Resource)
		***REMOVED***
		if j.NotBefore != notBefore.Format(time.RFC3339) ***REMOVED***
			t.Errorf(`notBefore = %q; wanted %q`, j.NotBefore, notBefore.Format(time.RFC3339))
		***REMOVED***
		if j.NotAfter != notAfter.Format(time.RFC3339) ***REMOVED***
			t.Errorf(`notAfter = %q; wanted %q`, j.NotAfter, notAfter.Format(time.RFC3339))
		***REMOVED***

		// Respond to request
		template := x509.Certificate***REMOVED***
			SerialNumber: big.NewInt(int64(1)),
			Subject: pkix.Name***REMOVED***
				Organization: []string***REMOVED***"goacme"***REMOVED***,
			***REMOVED***,
			NotBefore: notBefore,
			NotAfter:  notAfter,

			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:           []x509.ExtKeyUsage***REMOVED***x509.ExtKeyUsageServerAuth***REMOVED***,
			BasicConstraintsValid: true,
		***REMOVED***

		sampleCert, err := x509.CreateCertificate(rand.Reader, &template, &template, &testKeyEC.PublicKey, testKeyEC)
		if err != nil ***REMOVED***
			t.Fatalf("Error creating certificate: %v", err)
		***REMOVED***

		w.Header().Set("Location", "https://ca.tld/acme/cert/1")
		w.WriteHeader(http.StatusCreated)
		w.Write(sampleCert)
	***REMOVED***))
	defer ts.Close()

	csr := x509.CertificateRequest***REMOVED***
		Version: 0,
		Subject: pkix.Name***REMOVED***
			CommonName:   "example.com",
			Organization: []string***REMOVED***"goacme"***REMOVED***,
		***REMOVED***,
	***REMOVED***
	csrb, err := x509.CreateCertificateRequest(rand.Reader, &csr, testKeyEC)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	c := Client***REMOVED***Key: testKeyEC, dir: &Directory***REMOVED***CertURL: ts.URL***REMOVED******REMOVED***
	cert, certURL, err := c.CreateCert(context.Background(), csrb, notAfter.Sub(notBefore), false)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if cert == nil ***REMOVED***
		t.Errorf("cert is nil")
	***REMOVED***
	if certURL != "https://ca.tld/acme/cert/1" ***REMOVED***
		t.Errorf("certURL = %q; want https://ca.tld/acme/cert/1", certURL)
	***REMOVED***
***REMOVED***

func TestFetchCert(t *testing.T) ***REMOVED***
	var count byte
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		count++
		if count < 3 ***REMOVED***
			up := fmt.Sprintf("<%s>;rel=up", ts.URL)
			w.Header().Set("Link", up)
		***REMOVED***
		w.Write([]byte***REMOVED***count***REMOVED***)
	***REMOVED***))
	defer ts.Close()
	res, err := (&Client***REMOVED******REMOVED***).FetchCert(context.Background(), ts.URL, true)
	if err != nil ***REMOVED***
		t.Fatalf("FetchCert: %v", err)
	***REMOVED***
	cert := [][]byte***REMOVED******REMOVED***1***REMOVED***, ***REMOVED***2***REMOVED***, ***REMOVED***3***REMOVED******REMOVED***
	if !reflect.DeepEqual(res, cert) ***REMOVED***
		t.Errorf("res = %v; want %v", res, cert)
	***REMOVED***
***REMOVED***

func TestFetchCertRetry(t *testing.T) ***REMOVED***
	var count int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if count < 1 ***REMOVED***
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusAccepted)
			count++
			return
		***REMOVED***
		w.Write([]byte***REMOVED***1***REMOVED***)
	***REMOVED***))
	defer ts.Close()
	res, err := (&Client***REMOVED******REMOVED***).FetchCert(context.Background(), ts.URL, false)
	if err != nil ***REMOVED***
		t.Fatalf("FetchCert: %v", err)
	***REMOVED***
	cert := [][]byte***REMOVED******REMOVED***1***REMOVED******REMOVED***
	if !reflect.DeepEqual(res, cert) ***REMOVED***
		t.Errorf("res = %v; want %v", res, cert)
	***REMOVED***
***REMOVED***

func TestFetchCertCancel(t *testing.T) ***REMOVED***
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusAccepted)
	***REMOVED***))
	defer ts.Close()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct***REMOVED******REMOVED***)
	var err error
	go func() ***REMOVED***
		_, err = (&Client***REMOVED******REMOVED***).FetchCert(ctx, ts.URL, false)
		close(done)
	***REMOVED***()
	cancel()
	<-done
	if err != context.Canceled ***REMOVED***
		t.Errorf("err = %v; want %v", err, context.Canceled)
	***REMOVED***
***REMOVED***

func TestFetchCertDepth(t *testing.T) ***REMOVED***
	var count byte
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		count++
		if count > maxChainLen+1 ***REMOVED***
			t.Errorf("count = %d; want at most %d", count, maxChainLen+1)
			w.WriteHeader(http.StatusInternalServerError)
		***REMOVED***
		w.Header().Set("Link", fmt.Sprintf("<%s>;rel=up", ts.URL))
		w.Write([]byte***REMOVED***count***REMOVED***)
	***REMOVED***))
	defer ts.Close()
	_, err := (&Client***REMOVED******REMOVED***).FetchCert(context.Background(), ts.URL, true)
	if err == nil ***REMOVED***
		t.Errorf("err is nil")
	***REMOVED***
***REMOVED***

func TestFetchCertBreadth(t *testing.T) ***REMOVED***
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		for i := 0; i < maxChainLen+1; i++ ***REMOVED***
			w.Header().Add("Link", fmt.Sprintf("<%s>;rel=up", ts.URL))
		***REMOVED***
		w.Write([]byte***REMOVED***1***REMOVED***)
	***REMOVED***))
	defer ts.Close()
	_, err := (&Client***REMOVED******REMOVED***).FetchCert(context.Background(), ts.URL, true)
	if err == nil ***REMOVED***
		t.Errorf("err is nil")
	***REMOVED***
***REMOVED***

func TestFetchCertSize(t *testing.T) ***REMOVED***
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		b := bytes.Repeat([]byte***REMOVED***1***REMOVED***, maxCertSize+1)
		w.Write(b)
	***REMOVED***))
	defer ts.Close()
	_, err := (&Client***REMOVED******REMOVED***).FetchCert(context.Background(), ts.URL, false)
	if err == nil ***REMOVED***
		t.Errorf("err is nil")
	***REMOVED***
***REMOVED***

func TestRevokeCert(t *testing.T) ***REMOVED***
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method == "HEAD" ***REMOVED***
			w.Header().Set("Replay-Nonce", "nonce")
			return
		***REMOVED***

		var req struct ***REMOVED***
			Resource    string
			Certificate string
			Reason      int
		***REMOVED***
		decodeJWSRequest(t, &req, r)
		if req.Resource != "revoke-cert" ***REMOVED***
			t.Errorf("req.Resource = %q; want revoke-cert", req.Resource)
		***REMOVED***
		if req.Reason != 1 ***REMOVED***
			t.Errorf("req.Reason = %d; want 1", req.Reason)
		***REMOVED***
		// echo -n cert | base64 | tr -d '=' | tr '/+' '_-'
		cert := "Y2VydA"
		if req.Certificate != cert ***REMOVED***
			t.Errorf("req.Certificate = %q; want %q", req.Certificate, cert)
		***REMOVED***
	***REMOVED***))
	defer ts.Close()
	client := &Client***REMOVED***
		Key: testKeyEC,
		dir: &Directory***REMOVED***RevokeURL: ts.URL***REMOVED***,
	***REMOVED***
	ctx := context.Background()
	if err := client.RevokeCert(ctx, nil, []byte("cert"), CRLReasonKeyCompromise); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestNonce_add(t *testing.T) ***REMOVED***
	var c Client
	c.addNonce(http.Header***REMOVED***"Replay-Nonce": ***REMOVED***"nonce"***REMOVED******REMOVED***)
	c.addNonce(http.Header***REMOVED***"Replay-Nonce": ***REMOVED******REMOVED******REMOVED***)
	c.addNonce(http.Header***REMOVED***"Replay-Nonce": ***REMOVED***"nonce"***REMOVED******REMOVED***)

	nonces := map[string]struct***REMOVED******REMOVED******REMOVED***"nonce": ***REMOVED******REMOVED******REMOVED***
	if !reflect.DeepEqual(c.nonces, nonces) ***REMOVED***
		t.Errorf("c.nonces = %q; want %q", c.nonces, nonces)
	***REMOVED***
***REMOVED***

func TestNonce_addMax(t *testing.T) ***REMOVED***
	c := &Client***REMOVED***nonces: make(map[string]struct***REMOVED******REMOVED***)***REMOVED***
	for i := 0; i < maxNonces; i++ ***REMOVED***
		c.nonces[fmt.Sprintf("%d", i)] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	c.addNonce(http.Header***REMOVED***"Replay-Nonce": ***REMOVED***"nonce"***REMOVED******REMOVED***)
	if n := len(c.nonces); n != maxNonces ***REMOVED***
		t.Errorf("len(c.nonces) = %d; want %d", n, maxNonces)
	***REMOVED***
***REMOVED***

func TestNonce_fetch(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		code  int
		nonce string
	***REMOVED******REMOVED***
		***REMOVED***http.StatusOK, "nonce1"***REMOVED***,
		***REMOVED***http.StatusBadRequest, "nonce2"***REMOVED***,
		***REMOVED***http.StatusOK, ""***REMOVED***,
	***REMOVED***
	var i int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method != "HEAD" ***REMOVED***
			t.Errorf("%d: r.Method = %q; want HEAD", i, r.Method)
		***REMOVED***
		w.Header().Set("Replay-Nonce", tests[i].nonce)
		w.WriteHeader(tests[i].code)
	***REMOVED***))
	defer ts.Close()
	for ; i < len(tests); i++ ***REMOVED***
		test := tests[i]
		c := &Client***REMOVED******REMOVED***
		n, err := c.fetchNonce(context.Background(), ts.URL)
		if n != test.nonce ***REMOVED***
			t.Errorf("%d: n=%q; want %q", i, n, test.nonce)
		***REMOVED***
		switch ***REMOVED***
		case err == nil && test.nonce == "":
			t.Errorf("%d: n=%q, err=%v; want non-nil error", i, n, err)
		case err != nil && test.nonce != "":
			t.Errorf("%d: n=%q, err=%v; want %q", i, n, err, test.nonce)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestNonce_fetchError(t *testing.T) ***REMOVED***
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.WriteHeader(http.StatusTooManyRequests)
	***REMOVED***))
	defer ts.Close()
	c := &Client***REMOVED******REMOVED***
	_, err := c.fetchNonce(context.Background(), ts.URL)
	e, ok := err.(*Error)
	if !ok ***REMOVED***
		t.Fatalf("err is %T; want *Error", err)
	***REMOVED***
	if e.StatusCode != http.StatusTooManyRequests ***REMOVED***
		t.Errorf("e.StatusCode = %d; want %d", e.StatusCode, http.StatusTooManyRequests)
	***REMOVED***
***REMOVED***

func TestNonce_postJWS(t *testing.T) ***REMOVED***
	var count int
	seen := make(map[string]bool)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		count++
		w.Header().Set("Replay-Nonce", fmt.Sprintf("nonce%d", count))
		if r.Method == "HEAD" ***REMOVED***
			// We expect the client do a HEAD request
			// but only to fetch the first nonce.
			return
		***REMOVED***
		// Make client.Authorize happy; we're not testing its result.
		defer func() ***REMOVED***
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`***REMOVED***"status":"valid"***REMOVED***`))
		***REMOVED***()

		head, err := decodeJWSHead(r)
		if err != nil ***REMOVED***
			t.Errorf("decodeJWSHead: %v", err)
			return
		***REMOVED***
		if head.Nonce == "" ***REMOVED***
			t.Error("head.Nonce is empty")
			return
		***REMOVED***
		if seen[head.Nonce] ***REMOVED***
			t.Errorf("nonce is already used: %q", head.Nonce)
		***REMOVED***
		seen[head.Nonce] = true
	***REMOVED***))
	defer ts.Close()

	client := Client***REMOVED***Key: testKey, dir: &Directory***REMOVED***AuthzURL: ts.URL***REMOVED******REMOVED***
	if _, err := client.Authorize(context.Background(), "example.com"); err != nil ***REMOVED***
		t.Errorf("client.Authorize 1: %v", err)
	***REMOVED***
	// The second call should not generate another extra HEAD request.
	if _, err := client.Authorize(context.Background(), "example.com"); err != nil ***REMOVED***
		t.Errorf("client.Authorize 2: %v", err)
	***REMOVED***

	if count != 3 ***REMOVED***
		t.Errorf("total requests count: %d; want 3", count)
	***REMOVED***
	if n := len(client.nonces); n != 1 ***REMOVED***
		t.Errorf("len(client.nonces) = %d; want 1", n)
	***REMOVED***
	for k := range seen ***REMOVED***
		if _, exist := client.nonces[k]; exist ***REMOVED***
			t.Errorf("used nonce %q in client.nonces", k)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestRetryPostJWS(t *testing.T) ***REMOVED***
	var count int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		count++
		w.Header().Set("Replay-Nonce", fmt.Sprintf("nonce%d", count))
		if r.Method == "HEAD" ***REMOVED***
			// We expect the client to do 2 head requests to fetch
			// nonces, one to start and another after getting badNonce
			return
		***REMOVED***

		head, err := decodeJWSHead(r)
		if err != nil ***REMOVED***
			t.Errorf("decodeJWSHead: %v", err)
		***REMOVED*** else if head.Nonce == "" ***REMOVED***
			t.Error("head.Nonce is empty")
		***REMOVED*** else if head.Nonce == "nonce1" ***REMOVED***
			// return a badNonce error to force the call to retry
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`***REMOVED***"type":"urn:ietf:params:acme:error:badNonce"***REMOVED***`))
			return
		***REMOVED***
		// Make client.Authorize happy; we're not testing its result.
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`***REMOVED***"status":"valid"***REMOVED***`))
	***REMOVED***))
	defer ts.Close()

	client := Client***REMOVED***Key: testKey, dir: &Directory***REMOVED***AuthzURL: ts.URL***REMOVED******REMOVED***
	// This call will fail with badNonce, causing a retry
	if _, err := client.Authorize(context.Background(), "example.com"); err != nil ***REMOVED***
		t.Errorf("client.Authorize 1: %v", err)
	***REMOVED***
	if count != 4 ***REMOVED***
		t.Errorf("total requests count: %d; want 4", count)
	***REMOVED***
***REMOVED***

func TestLinkHeader(t *testing.T) ***REMOVED***
	h := http.Header***REMOVED***"Link": ***REMOVED***
		`<https://example.com/acme/new-authz>;rel="next"`,
		`<https://example.com/acme/recover-reg>; rel=recover`,
		`<https://example.com/acme/terms>; foo=bar; rel="terms-of-service"`,
		`<dup>;rel="next"`,
	***REMOVED******REMOVED***
	tests := []struct ***REMOVED***
		rel string
		out []string
	***REMOVED******REMOVED***
		***REMOVED***"next", []string***REMOVED***"https://example.com/acme/new-authz", "dup"***REMOVED******REMOVED***,
		***REMOVED***"recover", []string***REMOVED***"https://example.com/acme/recover-reg"***REMOVED******REMOVED***,
		***REMOVED***"terms-of-service", []string***REMOVED***"https://example.com/acme/terms"***REMOVED******REMOVED***,
		***REMOVED***"empty", nil***REMOVED***,
	***REMOVED***
	for i, test := range tests ***REMOVED***
		if v := linkHeader(h, test.rel); !reflect.DeepEqual(v, test.out) ***REMOVED***
			t.Errorf("%d: linkHeader(%q): %v; want %v", i, test.rel, v, test.out)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestErrorResponse(t *testing.T) ***REMOVED***
	s := `***REMOVED***
		"status": 400,
		"type": "urn:acme:error:xxx",
		"detail": "text"
	***REMOVED***`
	res := &http.Response***REMOVED***
		StatusCode: 400,
		Status:     "400 Bad Request",
		Body:       ioutil.NopCloser(strings.NewReader(s)),
		Header:     http.Header***REMOVED***"X-Foo": ***REMOVED***"bar"***REMOVED******REMOVED***,
	***REMOVED***
	err := responseError(res)
	v, ok := err.(*Error)
	if !ok ***REMOVED***
		t.Fatalf("err = %+v (%T); want *Error type", err, err)
	***REMOVED***
	if v.StatusCode != 400 ***REMOVED***
		t.Errorf("v.StatusCode = %v; want 400", v.StatusCode)
	***REMOVED***
	if v.ProblemType != "urn:acme:error:xxx" ***REMOVED***
		t.Errorf("v.ProblemType = %q; want urn:acme:error:xxx", v.ProblemType)
	***REMOVED***
	if v.Detail != "text" ***REMOVED***
		t.Errorf("v.Detail = %q; want text", v.Detail)
	***REMOVED***
	if !reflect.DeepEqual(v.Header, res.Header) ***REMOVED***
		t.Errorf("v.Header = %+v; want %+v", v.Header, res.Header)
	***REMOVED***
***REMOVED***

func TestTLSSNI01ChallengeCert(t *testing.T) ***REMOVED***
	const (
		token = "evaGxfADs6pSRb2LAv9IZf17Dt3juxGJ-PCt92wr-oA"
		// echo -n <token.testKeyECThumbprint> | shasum -a 256
		san = "dbbd5eefe7b4d06eb9d1d9f5acb4c7cd.a27d320e4b30332f0b6cb441734ad7b0.acme.invalid"
	)

	client := &Client***REMOVED***Key: testKeyEC***REMOVED***
	tlscert, name, err := client.TLSSNI01ChallengeCert(token)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if n := len(tlscert.Certificate); n != 1 ***REMOVED***
		t.Fatalf("len(tlscert.Certificate) = %d; want 1", n)
	***REMOVED***
	cert, err := x509.ParseCertificate(tlscert.Certificate[0])
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if len(cert.DNSNames) != 1 || cert.DNSNames[0] != san ***REMOVED***
		t.Fatalf("cert.DNSNames = %v; want %q", cert.DNSNames, san)
	***REMOVED***
	if cert.DNSNames[0] != name ***REMOVED***
		t.Errorf("cert.DNSNames[0] != name: %q vs %q", cert.DNSNames[0], name)
	***REMOVED***
	if cn := cert.Subject.CommonName; cn != san ***REMOVED***
		t.Errorf("cert.Subject.CommonName = %q; want %q", cn, san)
	***REMOVED***
***REMOVED***

func TestTLSSNI02ChallengeCert(t *testing.T) ***REMOVED***
	const (
		token = "evaGxfADs6pSRb2LAv9IZf17Dt3juxGJ-PCt92wr-oA"
		// echo -n evaGxfADs6pSRb2LAv9IZf17Dt3juxGJ-PCt92wr-oA | shasum -a 256
		sanA = "7ea0aaa69214e71e02cebb18bb867736.09b730209baabf60e43d4999979ff139.token.acme.invalid"
		// echo -n <token.testKeyECThumbprint> | shasum -a 256
		sanB = "dbbd5eefe7b4d06eb9d1d9f5acb4c7cd.a27d320e4b30332f0b6cb441734ad7b0.ka.acme.invalid"
	)

	client := &Client***REMOVED***Key: testKeyEC***REMOVED***
	tlscert, name, err := client.TLSSNI02ChallengeCert(token)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if n := len(tlscert.Certificate); n != 1 ***REMOVED***
		t.Fatalf("len(tlscert.Certificate) = %d; want 1", n)
	***REMOVED***
	cert, err := x509.ParseCertificate(tlscert.Certificate[0])
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	names := []string***REMOVED***sanA, sanB***REMOVED***
	if !reflect.DeepEqual(cert.DNSNames, names) ***REMOVED***
		t.Fatalf("cert.DNSNames = %v;\nwant %v", cert.DNSNames, names)
	***REMOVED***
	sort.Strings(cert.DNSNames)
	i := sort.SearchStrings(cert.DNSNames, name)
	if i >= len(cert.DNSNames) || cert.DNSNames[i] != name ***REMOVED***
		t.Errorf("%v doesn't have %q", cert.DNSNames, name)
	***REMOVED***
	if cn := cert.Subject.CommonName; cn != sanA ***REMOVED***
		t.Errorf("CommonName = %q; want %q", cn, sanA)
	***REMOVED***
***REMOVED***

func TestTLSChallengeCertOpt(t *testing.T) ***REMOVED***
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	tmpl := &x509.Certificate***REMOVED***
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name***REMOVED***Organization: []string***REMOVED***"Test"***REMOVED******REMOVED***,
		DNSNames:     []string***REMOVED***"should-be-overwritten"***REMOVED***,
	***REMOVED***
	opts := []CertOption***REMOVED***WithKey(key), WithTemplate(tmpl)***REMOVED***

	client := &Client***REMOVED***Key: testKeyEC***REMOVED***
	cert1, _, err := client.TLSSNI01ChallengeCert("token", opts...)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	cert2, _, err := client.TLSSNI02ChallengeCert("token", opts...)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	for i, tlscert := range []tls.Certificate***REMOVED***cert1, cert2***REMOVED*** ***REMOVED***
		// verify generated cert private key
		tlskey, ok := tlscert.PrivateKey.(*rsa.PrivateKey)
		if !ok ***REMOVED***
			t.Errorf("%d: tlscert.PrivateKey is %T; want *rsa.PrivateKey", i, tlscert.PrivateKey)
			continue
		***REMOVED***
		if tlskey.D.Cmp(key.D) != 0 ***REMOVED***
			t.Errorf("%d: tlskey.D = %v; want %v", i, tlskey.D, key.D)
		***REMOVED***
		// verify generated cert public key
		x509Cert, err := x509.ParseCertificate(tlscert.Certificate[0])
		if err != nil ***REMOVED***
			t.Errorf("%d: %v", i, err)
			continue
		***REMOVED***
		tlspub, ok := x509Cert.PublicKey.(*rsa.PublicKey)
		if !ok ***REMOVED***
			t.Errorf("%d: x509Cert.PublicKey is %T; want *rsa.PublicKey", i, x509Cert.PublicKey)
			continue
		***REMOVED***
		if tlspub.N.Cmp(key.N) != 0 ***REMOVED***
			t.Errorf("%d: tlspub.N = %v; want %v", i, tlspub.N, key.N)
		***REMOVED***
		// verify template option
		sn := big.NewInt(2)
		if x509Cert.SerialNumber.Cmp(sn) != 0 ***REMOVED***
			t.Errorf("%d: SerialNumber = %v; want %v", i, x509Cert.SerialNumber, sn)
		***REMOVED***
		org := []string***REMOVED***"Test"***REMOVED***
		if !reflect.DeepEqual(x509Cert.Subject.Organization, org) ***REMOVED***
			t.Errorf("%d: Subject.Organization = %+v; want %+v", i, x509Cert.Subject.Organization, org)
		***REMOVED***
		for _, v := range x509Cert.DNSNames ***REMOVED***
			if !strings.HasSuffix(v, ".acme.invalid") ***REMOVED***
				t.Errorf("%d: invalid DNSNames element: %q", i, v)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestHTTP01Challenge(t *testing.T) ***REMOVED***
	const (
		token = "xxx"
		// thumbprint is precomputed for testKeyEC in jws_test.go
		value   = token + "." + testKeyECThumbprint
		urlpath = "/.well-known/acme-challenge/" + token
	)
	client := &Client***REMOVED***Key: testKeyEC***REMOVED***
	val, err := client.HTTP01ChallengeResponse(token)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if val != value ***REMOVED***
		t.Errorf("val = %q; want %q", val, value)
	***REMOVED***
	if path := client.HTTP01ChallengePath(token); path != urlpath ***REMOVED***
		t.Errorf("path = %q; want %q", path, urlpath)
	***REMOVED***
***REMOVED***

func TestDNS01ChallengeRecord(t *testing.T) ***REMOVED***
	// echo -n xxx.<testKeyECThumbprint> | \
	//      openssl dgst -binary -sha256 | \
	//      base64 | tr -d '=' | tr '/+' '_-'
	const value = "8DERMexQ5VcdJ_prpPiA0mVdp7imgbCgjsG4SqqNMIo"

	client := &Client***REMOVED***Key: testKeyEC***REMOVED***
	val, err := client.DNS01ChallengeRecord("xxx")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if val != value ***REMOVED***
		t.Errorf("val = %q; want %q", val, value)
	***REMOVED***
***REMOVED***

func TestBackoff(t *testing.T) ***REMOVED***
	tt := []struct***REMOVED*** min, max time.Duration ***REMOVED******REMOVED***
		***REMOVED***time.Second, 2 * time.Second***REMOVED***,
		***REMOVED***2 * time.Second, 3 * time.Second***REMOVED***,
		***REMOVED***4 * time.Second, 5 * time.Second***REMOVED***,
		***REMOVED***8 * time.Second, 9 * time.Second***REMOVED***,
	***REMOVED***
	for i, test := range tt ***REMOVED***
		d := backoff(i, time.Minute)
		if d < test.min || test.max < d ***REMOVED***
			t.Errorf("%d: d = %v; want between %v and %v", i, d, test.min, test.max)
		***REMOVED***
	***REMOVED***

	min, max := time.Second, 2*time.Second
	if d := backoff(-1, time.Minute); d < min || max < d ***REMOVED***
		t.Errorf("d = %v; want between %v and %v", d, min, max)
	***REMOVED***

	bound := 10 * time.Second
	if d := backoff(100, bound); d != bound ***REMOVED***
		t.Errorf("d = %v; want %v", d, bound)
	***REMOVED***
***REMOVED***
