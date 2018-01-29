// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package autocert

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/crypto/acme"
)

var discoTmpl = template.Must(template.New("disco").Parse(`***REMOVED***
	"new-reg": "***REMOVED******REMOVED***.***REMOVED******REMOVED***/new-reg",
	"new-authz": "***REMOVED******REMOVED***.***REMOVED******REMOVED***/new-authz",
	"new-cert": "***REMOVED******REMOVED***.***REMOVED******REMOVED***/new-cert"
***REMOVED***`))

var authzTmpl = template.Must(template.New("authz").Parse(`***REMOVED***
	"status": "pending",
	"challenges": [
		***REMOVED***
			"uri": "***REMOVED******REMOVED***.***REMOVED******REMOVED***/challenge/1",
			"type": "tls-sni-01",
			"token": "token-01"
		***REMOVED***,
		***REMOVED***
			"uri": "***REMOVED******REMOVED***.***REMOVED******REMOVED***/challenge/2",
			"type": "tls-sni-02",
			"token": "token-02"
		***REMOVED***,
		***REMOVED***
			"uri": "***REMOVED******REMOVED***.***REMOVED******REMOVED***/challenge/dns-01",
			"type": "dns-01",
			"token": "token-dns-01"
		***REMOVED***,
		***REMOVED***
			"uri": "***REMOVED******REMOVED***.***REMOVED******REMOVED***/challenge/http-01",
			"type": "http-01",
			"token": "token-http-01"
		***REMOVED***
	]
***REMOVED***`))

type memCache struct ***REMOVED***
	mu      sync.Mutex
	keyData map[string][]byte
***REMOVED***

func (m *memCache) Get(ctx context.Context, key string) ([]byte, error) ***REMOVED***
	m.mu.Lock()
	defer m.mu.Unlock()

	v, ok := m.keyData[key]
	if !ok ***REMOVED***
		return nil, ErrCacheMiss
	***REMOVED***
	return v, nil
***REMOVED***

func (m *memCache) Put(ctx context.Context, key string, data []byte) error ***REMOVED***
	m.mu.Lock()
	defer m.mu.Unlock()

	m.keyData[key] = data
	return nil
***REMOVED***

func (m *memCache) Delete(ctx context.Context, key string) error ***REMOVED***
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.keyData, key)
	return nil
***REMOVED***

func newMemCache() *memCache ***REMOVED***
	return &memCache***REMOVED***
		keyData: make(map[string][]byte),
	***REMOVED***
***REMOVED***

func dummyCert(pub interface***REMOVED******REMOVED***, san ...string) ([]byte, error) ***REMOVED***
	return dateDummyCert(pub, time.Now(), time.Now().Add(90*24*time.Hour), san...)
***REMOVED***

func dateDummyCert(pub interface***REMOVED******REMOVED***, start, end time.Time, san ...string) ([]byte, error) ***REMOVED***
	// use EC key to run faster on 386
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	t := &x509.Certificate***REMOVED***
		SerialNumber:          big.NewInt(1),
		NotBefore:             start,
		NotAfter:              end,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageKeyEncipherment,
		DNSNames:              san,
	***REMOVED***
	if pub == nil ***REMOVED***
		pub = &key.PublicKey
	***REMOVED***
	return x509.CreateCertificate(rand.Reader, t, t, pub, key)
***REMOVED***

func decodePayload(v interface***REMOVED******REMOVED***, r io.Reader) error ***REMOVED***
	var req struct***REMOVED*** Payload string ***REMOVED***
	if err := json.NewDecoder(r).Decode(&req); err != nil ***REMOVED***
		return err
	***REMOVED***
	payload, err := base64.RawURLEncoding.DecodeString(req.Payload)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return json.Unmarshal(payload, v)
***REMOVED***

func TestGetCertificate(t *testing.T) ***REMOVED***
	man := &Manager***REMOVED***Prompt: AcceptTOS***REMOVED***
	defer man.stopRenew()
	hello := &tls.ClientHelloInfo***REMOVED***ServerName: "example.org"***REMOVED***
	testGetCertificate(t, man, "example.org", hello)
***REMOVED***

func TestGetCertificate_trailingDot(t *testing.T) ***REMOVED***
	man := &Manager***REMOVED***Prompt: AcceptTOS***REMOVED***
	defer man.stopRenew()
	hello := &tls.ClientHelloInfo***REMOVED***ServerName: "example.org."***REMOVED***
	testGetCertificate(t, man, "example.org", hello)
***REMOVED***

func TestGetCertificate_ForceRSA(t *testing.T) ***REMOVED***
	man := &Manager***REMOVED***
		Prompt:   AcceptTOS,
		Cache:    newMemCache(),
		ForceRSA: true,
	***REMOVED***
	defer man.stopRenew()
	hello := &tls.ClientHelloInfo***REMOVED***ServerName: "example.org"***REMOVED***
	testGetCertificate(t, man, "example.org", hello)

	cert, err := man.cacheGet(context.Background(), "example.org")
	if err != nil ***REMOVED***
		t.Fatalf("man.cacheGet: %v", err)
	***REMOVED***
	if _, ok := cert.PrivateKey.(*rsa.PrivateKey); !ok ***REMOVED***
		t.Errorf("cert.PrivateKey is %T; want *rsa.PrivateKey", cert.PrivateKey)
	***REMOVED***
***REMOVED***

func TestGetCertificate_nilPrompt(t *testing.T) ***REMOVED***
	man := &Manager***REMOVED******REMOVED***
	defer man.stopRenew()
	url, finish := startACMEServerStub(t, man, "example.org")
	defer finish()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	man.Client = &acme.Client***REMOVED***
		Key:          key,
		DirectoryURL: url,
	***REMOVED***
	hello := &tls.ClientHelloInfo***REMOVED***ServerName: "example.org"***REMOVED***
	if _, err := man.GetCertificate(hello); err == nil ***REMOVED***
		t.Error("got certificate for example.org; wanted error")
	***REMOVED***
***REMOVED***

func TestGetCertificate_expiredCache(t *testing.T) ***REMOVED***
	// Make an expired cert and cache it.
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	tmpl := &x509.Certificate***REMOVED***
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name***REMOVED***CommonName: "example.org"***REMOVED***,
		NotAfter:     time.Now(),
	***REMOVED***
	pub, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &pk.PublicKey, pk)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	tlscert := &tls.Certificate***REMOVED***
		Certificate: [][]byte***REMOVED***pub***REMOVED***,
		PrivateKey:  pk,
	***REMOVED***

	man := &Manager***REMOVED***Prompt: AcceptTOS, Cache: newMemCache()***REMOVED***
	defer man.stopRenew()
	if err := man.cachePut(context.Background(), "example.org", tlscert); err != nil ***REMOVED***
		t.Fatalf("man.cachePut: %v", err)
	***REMOVED***

	// The expired cached cert should trigger a new cert issuance
	// and return without an error.
	hello := &tls.ClientHelloInfo***REMOVED***ServerName: "example.org"***REMOVED***
	testGetCertificate(t, man, "example.org", hello)
***REMOVED***

func TestGetCertificate_failedAttempt(t *testing.T) ***REMOVED***
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.WriteHeader(http.StatusBadRequest)
	***REMOVED***))
	defer ts.Close()

	const example = "example.org"
	d := createCertRetryAfter
	f := testDidRemoveState
	defer func() ***REMOVED***
		createCertRetryAfter = d
		testDidRemoveState = f
	***REMOVED***()
	createCertRetryAfter = 0
	done := make(chan struct***REMOVED******REMOVED***)
	testDidRemoveState = func(domain string) ***REMOVED***
		if domain != example ***REMOVED***
			t.Errorf("testDidRemoveState: domain = %q; want %q", domain, example)
		***REMOVED***
		close(done)
	***REMOVED***

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	man := &Manager***REMOVED***
		Prompt: AcceptTOS,
		Client: &acme.Client***REMOVED***
			Key:          key,
			DirectoryURL: ts.URL,
		***REMOVED***,
	***REMOVED***
	defer man.stopRenew()
	hello := &tls.ClientHelloInfo***REMOVED***ServerName: example***REMOVED***
	if _, err := man.GetCertificate(hello); err == nil ***REMOVED***
		t.Error("GetCertificate: err is nil")
	***REMOVED***
	select ***REMOVED***
	case <-time.After(5 * time.Second):
		t.Errorf("took too long to remove the %q state", example)
	case <-done:
		man.stateMu.Lock()
		defer man.stateMu.Unlock()
		if v, exist := man.state[example]; exist ***REMOVED***
			t.Errorf("state exists for %q: %+v", example, v)
		***REMOVED***
	***REMOVED***
***REMOVED***

// startACMEServerStub runs an ACME server
// The domain argument is the expected domain name of a certificate request.
func startACMEServerStub(t *testing.T, man *Manager, domain string) (url string, finish func()) ***REMOVED***
	// echo token-02 | shasum -a 256
	// then divide result in 2 parts separated by dot
	tokenCertName := "4e8eb87631187e9ff2153b56b13a4dec.13a35d002e485d60ff37354b32f665d9.token.acme.invalid"
	verifyTokenCert := func() ***REMOVED***
		hello := &tls.ClientHelloInfo***REMOVED***ServerName: tokenCertName***REMOVED***
		_, err := man.GetCertificate(hello)
		if err != nil ***REMOVED***
			t.Errorf("verifyTokenCert: GetCertificate(%q): %v", tokenCertName, err)
			return
		***REMOVED***
	***REMOVED***

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
				t.Errorf("discoTmpl: %v", err)
			***REMOVED***
		// client key registration
		case "/new-reg":
			w.Write([]byte("***REMOVED******REMOVED***"))
		// domain authorization
		case "/new-authz":
			w.Header().Set("Location", ca.URL+"/authz/1")
			w.WriteHeader(http.StatusCreated)
			if err := authzTmpl.Execute(w, ca.URL); err != nil ***REMOVED***
				t.Errorf("authzTmpl: %v", err)
			***REMOVED***
		// accept tls-sni-02 challenge
		case "/challenge/2":
			verifyTokenCert()
			w.Write([]byte("***REMOVED******REMOVED***"))
		// authorization status
		case "/authz/1":
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
				t.Errorf("new-cert: CSR: %v", err)
			***REMOVED***
			if csr.Subject.CommonName != domain ***REMOVED***
				t.Errorf("CommonName in CSR = %q; want %q", csr.Subject.CommonName, domain)
			***REMOVED***
			der, err := dummyCert(csr.PublicKey, domain)
			if err != nil ***REMOVED***
				t.Errorf("new-cert: dummyCert: %v", err)
			***REMOVED***
			chainUp := fmt.Sprintf("<%s/ca-cert>; rel=up", ca.URL)
			w.Header().Set("Link", chainUp)
			w.WriteHeader(http.StatusCreated)
			w.Write(der)
		// CA chain cert
		case "/ca-cert":
			der, err := dummyCert(nil, "ca")
			if err != nil ***REMOVED***
				t.Errorf("ca-cert: dummyCert: %v", err)
			***REMOVED***
			w.Write(der)
		default:
			t.Errorf("unrecognized r.URL.Path: %s", r.URL.Path)
		***REMOVED***
	***REMOVED***))
	finish = func() ***REMOVED***
		ca.Close()

		// make sure token cert was removed
		cancel := make(chan struct***REMOVED******REMOVED***)
		done := make(chan struct***REMOVED******REMOVED***)
		go func() ***REMOVED***
			defer close(done)
			tick := time.NewTicker(100 * time.Millisecond)
			defer tick.Stop()
			for ***REMOVED***
				hello := &tls.ClientHelloInfo***REMOVED***ServerName: tokenCertName***REMOVED***
				if _, err := man.GetCertificate(hello); err != nil ***REMOVED***
					return
				***REMOVED***
				select ***REMOVED***
				case <-tick.C:
				case <-cancel:
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***()
		select ***REMOVED***
		case <-done:
		case <-time.After(5 * time.Second):
			close(cancel)
			t.Error("token cert was not removed")
			<-done
		***REMOVED***
	***REMOVED***
	return ca.URL, finish
***REMOVED***

// tests man.GetCertificate flow using the provided hello argument.
// The domain argument is the expected domain name of a certificate request.
func testGetCertificate(t *testing.T, man *Manager, domain string, hello *tls.ClientHelloInfo) ***REMOVED***
	url, finish := startACMEServerStub(t, man, domain)
	defer finish()

	// use EC key to run faster on 386
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	man.Client = &acme.Client***REMOVED***
		Key:          key,
		DirectoryURL: url,
	***REMOVED***

	// simulate tls.Config.GetCertificate
	var tlscert *tls.Certificate
	done := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		tlscert, err = man.GetCertificate(hello)
		close(done)
	***REMOVED***()
	select ***REMOVED***
	case <-time.After(time.Minute):
		t.Fatal("man.GetCertificate took too long to return")
	case <-done:
	***REMOVED***
	if err != nil ***REMOVED***
		t.Fatalf("man.GetCertificate: %v", err)
	***REMOVED***

	// verify the tlscert is the same we responded with from the CA stub
	if len(tlscert.Certificate) == 0 ***REMOVED***
		t.Fatal("len(tlscert.Certificate) is 0")
	***REMOVED***
	cert, err := x509.ParseCertificate(tlscert.Certificate[0])
	if err != nil ***REMOVED***
		t.Fatalf("x509.ParseCertificate: %v", err)
	***REMOVED***
	if len(cert.DNSNames) == 0 || cert.DNSNames[0] != domain ***REMOVED***
		t.Errorf("cert.DNSNames = %v; want %q", cert.DNSNames, domain)
	***REMOVED***

***REMOVED***

func TestVerifyHTTP01(t *testing.T) ***REMOVED***
	var (
		http01 http.Handler

		authzCount      int // num. of created authorizations
		didAcceptHTTP01 bool
	)

	verifyHTTPToken := func() ***REMOVED***
		r := httptest.NewRequest("GET", "/.well-known/acme-challenge/token-http-01", nil)
		w := httptest.NewRecorder()
		http01.ServeHTTP(w, r)
		if w.Code != http.StatusOK ***REMOVED***
			t.Errorf("http token: w.Code = %d; want %d", w.Code, http.StatusOK)
		***REMOVED***
		if v := string(w.Body.Bytes()); !strings.HasPrefix(v, "token-http-01.") ***REMOVED***
			t.Errorf("http token value = %q; want 'token-http-01.' prefix", v)
		***REMOVED***
	***REMOVED***

	// ACME CA server stub, only the needed bits.
	// TODO: Merge this with startACMEServerStub, making it a configurable CA for testing.
	var ca *httptest.Server
	ca = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Replay-Nonce", "nonce")
		if r.Method == "HEAD" ***REMOVED***
			// a nonce request
			return
		***REMOVED***

		switch r.URL.Path ***REMOVED***
		// Discovery.
		case "/":
			if err := discoTmpl.Execute(w, ca.URL); err != nil ***REMOVED***
				t.Errorf("discoTmpl: %v", err)
			***REMOVED***
		// Client key registration.
		case "/new-reg":
			w.Write([]byte("***REMOVED******REMOVED***"))
		// New domain authorization.
		case "/new-authz":
			authzCount++
			w.Header().Set("Location", fmt.Sprintf("%s/authz/%d", ca.URL, authzCount))
			w.WriteHeader(http.StatusCreated)
			if err := authzTmpl.Execute(w, ca.URL); err != nil ***REMOVED***
				t.Errorf("authzTmpl: %v", err)
			***REMOVED***
		// Accept tls-sni-02.
		case "/challenge/2":
			w.Write([]byte("***REMOVED******REMOVED***"))
		// Reject tls-sni-01.
		case "/challenge/1":
			http.Error(w, "won't accept tls-sni-01", http.StatusBadRequest)
		// Should not accept dns-01.
		case "/challenge/dns-01":
			t.Errorf("dns-01 challenge was accepted")
			http.Error(w, "won't accept dns-01", http.StatusBadRequest)
		// Accept http-01.
		case "/challenge/http-01":
			didAcceptHTTP01 = true
			verifyHTTPToken()
			w.Write([]byte("***REMOVED******REMOVED***"))
		// Authorization statuses.
		// Make tls-sni-xxx invalid.
		case "/authz/1", "/authz/2":
			w.Write([]byte(`***REMOVED***"status": "invalid"***REMOVED***`))
		case "/authz/3", "/authz/4":
			w.Write([]byte(`***REMOVED***"status": "valid"***REMOVED***`))
		default:
			http.NotFound(w, r)
			t.Errorf("unrecognized r.URL.Path: %s", r.URL.Path)
		***REMOVED***
	***REMOVED***))
	defer ca.Close()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	m := &Manager***REMOVED***
		Client: &acme.Client***REMOVED***
			Key:          key,
			DirectoryURL: ca.URL,
		***REMOVED***,
	***REMOVED***
	http01 = m.HTTPHandler(nil)
	if err := m.verify(context.Background(), m.Client, "example.org"); err != nil ***REMOVED***
		t.Errorf("m.verify: %v", err)
	***REMOVED***
	// Only tls-sni-01, tls-sni-02 and http-01 must be accepted
	// The dns-01 challenge is unsupported.
	if authzCount != 3 ***REMOVED***
		t.Errorf("authzCount = %d; want 3", authzCount)
	***REMOVED***
	if !didAcceptHTTP01 ***REMOVED***
		t.Error("did not accept http-01 challenge")
	***REMOVED***
***REMOVED***

func TestHTTPHandlerDefaultFallback(t *testing.T) ***REMOVED***
	tt := []struct ***REMOVED***
		method, url  string
		wantCode     int
		wantLocation string
	***REMOVED******REMOVED***
		***REMOVED***"GET", "http://example.org", 302, "https://example.org/"***REMOVED***,
		***REMOVED***"GET", "http://example.org/foo", 302, "https://example.org/foo"***REMOVED***,
		***REMOVED***"GET", "http://example.org/foo/bar/", 302, "https://example.org/foo/bar/"***REMOVED***,
		***REMOVED***"GET", "http://example.org/?a=b", 302, "https://example.org/?a=b"***REMOVED***,
		***REMOVED***"GET", "http://example.org/foo?a=b", 302, "https://example.org/foo?a=b"***REMOVED***,
		***REMOVED***"GET", "http://example.org:80/foo?a=b", 302, "https://example.org:443/foo?a=b"***REMOVED***,
		***REMOVED***"GET", "http://example.org:80/foo%20bar", 302, "https://example.org:443/foo%20bar"***REMOVED***,
		***REMOVED***"GET", "http://[2602:d1:xxxx::c60a]:1234", 302, "https://[2602:d1:xxxx::c60a]:443/"***REMOVED***,
		***REMOVED***"GET", "http://[2602:d1:xxxx::c60a]", 302, "https://[2602:d1:xxxx::c60a]/"***REMOVED***,
		***REMOVED***"GET", "http://[2602:d1:xxxx::c60a]/foo?a=b", 302, "https://[2602:d1:xxxx::c60a]/foo?a=b"***REMOVED***,
		***REMOVED***"HEAD", "http://example.org", 302, "https://example.org/"***REMOVED***,
		***REMOVED***"HEAD", "http://example.org/foo", 302, "https://example.org/foo"***REMOVED***,
		***REMOVED***"HEAD", "http://example.org/foo/bar/", 302, "https://example.org/foo/bar/"***REMOVED***,
		***REMOVED***"HEAD", "http://example.org/?a=b", 302, "https://example.org/?a=b"***REMOVED***,
		***REMOVED***"HEAD", "http://example.org/foo?a=b", 302, "https://example.org/foo?a=b"***REMOVED***,
		***REMOVED***"POST", "http://example.org", 400, ""***REMOVED***,
		***REMOVED***"PUT", "http://example.org", 400, ""***REMOVED***,
		***REMOVED***"GET", "http://example.org/.well-known/acme-challenge/x", 404, ""***REMOVED***,
	***REMOVED***
	var m Manager
	h := m.HTTPHandler(nil)
	for i, test := range tt ***REMOVED***
		r := httptest.NewRequest(test.method, test.url, nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != test.wantCode ***REMOVED***
			t.Errorf("%d: w.Code = %d; want %d", i, w.Code, test.wantCode)
			t.Errorf("%d: body: %s", i, w.Body.Bytes())
		***REMOVED***
		if v := w.Header().Get("Location"); v != test.wantLocation ***REMOVED***
			t.Errorf("%d: Location = %q; want %q", i, v, test.wantLocation)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestAccountKeyCache(t *testing.T) ***REMOVED***
	m := Manager***REMOVED***Cache: newMemCache()***REMOVED***
	ctx := context.Background()
	k1, err := m.accountKey(ctx)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	k2, err := m.accountKey(ctx)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !reflect.DeepEqual(k1, k2) ***REMOVED***
		t.Errorf("account keys don't match: k1 = %#v; k2 = %#v", k1, k2)
	***REMOVED***
***REMOVED***

func TestCache(t *testing.T) ***REMOVED***
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	tmpl := &x509.Certificate***REMOVED***
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name***REMOVED***CommonName: "example.org"***REMOVED***,
		NotAfter:     time.Now().Add(time.Hour),
	***REMOVED***
	pub, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &privKey.PublicKey, privKey)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	tlscert := &tls.Certificate***REMOVED***
		Certificate: [][]byte***REMOVED***pub***REMOVED***,
		PrivateKey:  privKey,
	***REMOVED***

	man := &Manager***REMOVED***Cache: newMemCache()***REMOVED***
	defer man.stopRenew()
	ctx := context.Background()
	if err := man.cachePut(ctx, "example.org", tlscert); err != nil ***REMOVED***
		t.Fatalf("man.cachePut: %v", err)
	***REMOVED***
	res, err := man.cacheGet(ctx, "example.org")
	if err != nil ***REMOVED***
		t.Fatalf("man.cacheGet: %v", err)
	***REMOVED***
	if res == nil ***REMOVED***
		t.Fatal("res is nil")
	***REMOVED***
***REMOVED***

func TestHostWhitelist(t *testing.T) ***REMOVED***
	policy := HostWhitelist("example.com", "example.org", "*.example.net")
	tt := []struct ***REMOVED***
		host  string
		allow bool
	***REMOVED******REMOVED***
		***REMOVED***"example.com", true***REMOVED***,
		***REMOVED***"example.org", true***REMOVED***,
		***REMOVED***"one.example.com", false***REMOVED***,
		***REMOVED***"two.example.org", false***REMOVED***,
		***REMOVED***"three.example.net", false***REMOVED***,
		***REMOVED***"dummy", false***REMOVED***,
	***REMOVED***
	for i, test := range tt ***REMOVED***
		err := policy(nil, test.host)
		if err != nil && test.allow ***REMOVED***
			t.Errorf("%d: policy(%q): %v; want nil", i, test.host, err)
		***REMOVED***
		if err == nil && !test.allow ***REMOVED***
			t.Errorf("%d: policy(%q): nil; want an error", i, test.host)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestValidCert(t *testing.T) ***REMOVED***
	key1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	key2, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	key3, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	cert1, err := dummyCert(key1.Public(), "example.org")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	cert2, err := dummyCert(key2.Public(), "example.org")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	cert3, err := dummyCert(key3.Public(), "example.org")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	now := time.Now()
	early, err := dateDummyCert(key1.Public(), now.Add(time.Hour), now.Add(2*time.Hour), "example.org")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expired, err := dateDummyCert(key1.Public(), now.Add(-2*time.Hour), now.Add(-time.Hour), "example.org")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	tt := []struct ***REMOVED***
		domain string
		key    crypto.Signer
		cert   [][]byte
		ok     bool
	***REMOVED******REMOVED***
		***REMOVED***"example.org", key1, [][]byte***REMOVED***cert1***REMOVED***, true***REMOVED***,
		***REMOVED***"example.org", key3, [][]byte***REMOVED***cert3***REMOVED***, true***REMOVED***,
		***REMOVED***"example.org", key1, [][]byte***REMOVED***cert1, cert2, cert3***REMOVED***, true***REMOVED***,
		***REMOVED***"example.org", key1, [][]byte***REMOVED***cert1, ***REMOVED***1***REMOVED******REMOVED***, false***REMOVED***,
		***REMOVED***"example.org", key1, [][]byte***REMOVED******REMOVED***1***REMOVED******REMOVED***, false***REMOVED***,
		***REMOVED***"example.org", key1, [][]byte***REMOVED***cert2***REMOVED***, false***REMOVED***,
		***REMOVED***"example.org", key2, [][]byte***REMOVED***cert1***REMOVED***, false***REMOVED***,
		***REMOVED***"example.org", key1, [][]byte***REMOVED***cert3***REMOVED***, false***REMOVED***,
		***REMOVED***"example.org", key3, [][]byte***REMOVED***cert1***REMOVED***, false***REMOVED***,
		***REMOVED***"example.net", key1, [][]byte***REMOVED***cert1***REMOVED***, false***REMOVED***,
		***REMOVED***"example.org", key1, [][]byte***REMOVED***early***REMOVED***, false***REMOVED***,
		***REMOVED***"example.org", key1, [][]byte***REMOVED***expired***REMOVED***, false***REMOVED***,
	***REMOVED***
	for i, test := range tt ***REMOVED***
		leaf, err := validCert(test.domain, test.cert, test.key)
		if err != nil && test.ok ***REMOVED***
			t.Errorf("%d: err = %v", i, err)
		***REMOVED***
		if err == nil && !test.ok ***REMOVED***
			t.Errorf("%d: err is nil", i)
		***REMOVED***
		if err == nil && test.ok && leaf == nil ***REMOVED***
			t.Errorf("%d: leaf is nil", i)
		***REMOVED***
	***REMOVED***
***REMOVED***

type cacheGetFunc func(ctx context.Context, key string) ([]byte, error)

func (f cacheGetFunc) Get(ctx context.Context, key string) ([]byte, error) ***REMOVED***
	return f(ctx, key)
***REMOVED***

func (f cacheGetFunc) Put(ctx context.Context, key string, data []byte) error ***REMOVED***
	return fmt.Errorf("unsupported Put of %q = %q", key, data)
***REMOVED***

func (f cacheGetFunc) Delete(ctx context.Context, key string) error ***REMOVED***
	return fmt.Errorf("unsupported Delete of %q", key)
***REMOVED***

func TestManagerGetCertificateBogusSNI(t *testing.T) ***REMOVED***
	m := Manager***REMOVED***
		Prompt: AcceptTOS,
		Cache: cacheGetFunc(func(ctx context.Context, key string) ([]byte, error) ***REMOVED***
			return nil, fmt.Errorf("cache.Get of %s", key)
		***REMOVED***),
	***REMOVED***
	tests := []struct ***REMOVED***
		name    string
		wantErr string
	***REMOVED******REMOVED***
		***REMOVED***"foo.com", "cache.Get of foo.com"***REMOVED***,
		***REMOVED***"foo.com.", "cache.Get of foo.com"***REMOVED***,
		***REMOVED***`a\b.com`, "acme/autocert: server name contains invalid character"***REMOVED***,
		***REMOVED***`a/b.com`, "acme/autocert: server name contains invalid character"***REMOVED***,
		***REMOVED***"", "acme/autocert: missing server name"***REMOVED***,
		***REMOVED***"foo", "acme/autocert: server name component count invalid"***REMOVED***,
		***REMOVED***".foo", "acme/autocert: server name component count invalid"***REMOVED***,
		***REMOVED***"foo.", "acme/autocert: server name component count invalid"***REMOVED***,
		***REMOVED***"fo.o", "cache.Get of fo.o"***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		_, err := m.GetCertificate(&tls.ClientHelloInfo***REMOVED***ServerName: tt.name***REMOVED***)
		got := fmt.Sprint(err)
		if got != tt.wantErr ***REMOVED***
			t.Errorf("GetCertificate(SNI = %q) = %q; want %q", tt.name, got, tt.wantErr)
		***REMOVED***
	***REMOVED***
***REMOVED***
