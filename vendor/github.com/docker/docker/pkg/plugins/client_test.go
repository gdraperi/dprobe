package plugins

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/pkg/plugins/transport"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/stretchr/testify/assert"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
)

func setupRemotePluginServer() string ***REMOVED***
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	return server.URL
***REMOVED***

func teardownRemotePluginServer() ***REMOVED***
	if server != nil ***REMOVED***
		server.Close()
	***REMOVED***
***REMOVED***

func TestFailedConnection(t *testing.T) ***REMOVED***
	c, _ := NewClient("tcp://127.0.0.1:1", &tlsconfig.Options***REMOVED***InsecureSkipVerify: true***REMOVED***)
	_, err := c.callWithRetry("Service.Method", nil, false)
	if err == nil ***REMOVED***
		t.Fatal("Unexpected successful connection")
	***REMOVED***
***REMOVED***

func TestFailOnce(t *testing.T) ***REMOVED***
	addr := setupRemotePluginServer()
	defer teardownRemotePluginServer()

	failed := false
	mux.HandleFunc("/Test.FailOnce", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if !failed ***REMOVED***
			failed = true
			panic("Plugin not ready")
		***REMOVED***
	***REMOVED***)

	c, _ := NewClient(addr, &tlsconfig.Options***REMOVED***InsecureSkipVerify: true***REMOVED***)
	b := strings.NewReader("body")
	_, err := c.callWithRetry("Test.FailOnce", b, true)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestEchoInputOutput(t *testing.T) ***REMOVED***
	addr := setupRemotePluginServer()
	defer teardownRemotePluginServer()

	m := Manifest***REMOVED***[]string***REMOVED***"VolumeDriver", "NetworkDriver"***REMOVED******REMOVED***

	mux.HandleFunc("/Test.Echo", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method != "POST" ***REMOVED***
			t.Fatalf("Expected POST, got %s\n", r.Method)
		***REMOVED***

		header := w.Header()
		header.Set("Content-Type", transport.VersionMimetype)

		io.Copy(w, r.Body)
	***REMOVED***)

	c, _ := NewClient(addr, &tlsconfig.Options***REMOVED***InsecureSkipVerify: true***REMOVED***)
	var output Manifest
	err := c.Call("Test.Echo", m, &output)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	assert.Equal(t, m, output)
	err = c.Call("Test.Echo", nil, nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestBackoff(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		retries    int
		expTimeOff time.Duration
	***REMOVED******REMOVED***
		***REMOVED***0, time.Duration(1)***REMOVED***,
		***REMOVED***1, time.Duration(2)***REMOVED***,
		***REMOVED***2, time.Duration(4)***REMOVED***,
		***REMOVED***4, time.Duration(16)***REMOVED***,
		***REMOVED***6, time.Duration(30)***REMOVED***,
		***REMOVED***10, time.Duration(30)***REMOVED***,
	***REMOVED***

	for _, c := range cases ***REMOVED***
		s := c.expTimeOff * time.Second
		if d := backoff(c.retries); d != s ***REMOVED***
			t.Fatalf("Retry %v, expected %v, was %v\n", c.retries, s, d)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestAbortRetry(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		timeOff  time.Duration
		expAbort bool
	***REMOVED******REMOVED***
		***REMOVED***time.Duration(1), false***REMOVED***,
		***REMOVED***time.Duration(2), false***REMOVED***,
		***REMOVED***time.Duration(10), false***REMOVED***,
		***REMOVED***time.Duration(30), true***REMOVED***,
		***REMOVED***time.Duration(40), true***REMOVED***,
	***REMOVED***

	for _, c := range cases ***REMOVED***
		s := c.timeOff * time.Second
		if a := abort(time.Now(), s); a != c.expAbort ***REMOVED***
			t.Fatalf("Duration %v, expected %v, was %v\n", c.timeOff, s, a)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestClientScheme(t *testing.T) ***REMOVED***
	cases := map[string]string***REMOVED***
		"tcp://127.0.0.1:8080":          "http",
		"unix:///usr/local/plugins/foo": "http",
		"http://127.0.0.1:8080":         "http",
		"https://127.0.0.1:8080":        "https",
	***REMOVED***

	for addr, scheme := range cases ***REMOVED***
		u, err := url.Parse(addr)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		s := httpScheme(u)

		if s != scheme ***REMOVED***
			t.Fatalf("URL scheme mismatch, expected %s, got %s", scheme, s)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestNewClientWithTimeout(t *testing.T) ***REMOVED***
	addr := setupRemotePluginServer()
	defer teardownRemotePluginServer()

	m := Manifest***REMOVED***[]string***REMOVED***"VolumeDriver", "NetworkDriver"***REMOVED******REMOVED***

	mux.HandleFunc("/Test.Echo", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		time.Sleep(time.Duration(600) * time.Millisecond)
		io.Copy(w, r.Body)
	***REMOVED***)

	// setting timeout of 500ms
	timeout := time.Duration(500) * time.Millisecond
	c, _ := NewClientWithTimeout(addr, &tlsconfig.Options***REMOVED***InsecureSkipVerify: true***REMOVED***, timeout)
	var output Manifest
	err := c.Call("Test.Echo", m, &output)
	if err == nil ***REMOVED***
		t.Fatal("Expected timeout error")
	***REMOVED***
***REMOVED***

func TestClientStream(t *testing.T) ***REMOVED***
	addr := setupRemotePluginServer()
	defer teardownRemotePluginServer()

	m := Manifest***REMOVED***[]string***REMOVED***"VolumeDriver", "NetworkDriver"***REMOVED******REMOVED***
	var output Manifest

	mux.HandleFunc("/Test.Echo", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method != "POST" ***REMOVED***
			t.Fatalf("Expected POST, got %s", r.Method)
		***REMOVED***

		header := w.Header()
		header.Set("Content-Type", transport.VersionMimetype)

		io.Copy(w, r.Body)
	***REMOVED***)

	c, _ := NewClient(addr, &tlsconfig.Options***REMOVED***InsecureSkipVerify: true***REMOVED***)
	body, err := c.Stream("Test.Echo", m)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer body.Close()
	if err := json.NewDecoder(body).Decode(&output); err != nil ***REMOVED***
		t.Fatalf("Test.Echo: error reading plugin resp: %v", err)
	***REMOVED***
	assert.Equal(t, m, output)
***REMOVED***

func TestClientSendFile(t *testing.T) ***REMOVED***
	addr := setupRemotePluginServer()
	defer teardownRemotePluginServer()

	m := Manifest***REMOVED***[]string***REMOVED***"VolumeDriver", "NetworkDriver"***REMOVED******REMOVED***
	var output Manifest
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(m); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	mux.HandleFunc("/Test.Echo", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method != "POST" ***REMOVED***
			t.Fatalf("Expected POST, got %s\n", r.Method)
		***REMOVED***

		header := w.Header()
		header.Set("Content-Type", transport.VersionMimetype)

		io.Copy(w, r.Body)
	***REMOVED***)

	c, _ := NewClient(addr, &tlsconfig.Options***REMOVED***InsecureSkipVerify: true***REMOVED***)
	if err := c.SendFile("Test.Echo", &buf, &output); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	assert.Equal(t, m, output)
***REMOVED***
