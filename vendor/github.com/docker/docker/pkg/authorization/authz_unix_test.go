// +build !windows

// TODO Windows: This uses a Unix socket for testing. This might be possible
// to port to Windows using a named pipe instead.

package authorization

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/docker/docker/pkg/plugins"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/gorilla/mux"
)

const (
	pluginAddress = "authz-test-plugin.sock"
)

func TestAuthZRequestPluginError(t *testing.T) ***REMOVED***
	server := authZPluginTestServer***REMOVED***t: t***REMOVED***
	server.start()
	defer server.stop()

	authZPlugin := createTestPlugin(t)

	request := Request***REMOVED***
		User:           "user",
		RequestBody:    []byte("sample body"),
		RequestURI:     "www.authz.com/auth",
		RequestMethod:  "GET",
		RequestHeaders: map[string]string***REMOVED***"header": "value"***REMOVED***,
	***REMOVED***
	server.replayResponse = Response***REMOVED***
		Err: "an error",
	***REMOVED***

	actualResponse, err := authZPlugin.AuthZRequest(&request)
	if err != nil ***REMOVED***
		t.Fatalf("Failed to authorize request %v", err)
	***REMOVED***

	if !reflect.DeepEqual(server.replayResponse, *actualResponse) ***REMOVED***
		t.Fatal("Response must be equal")
	***REMOVED***
	if !reflect.DeepEqual(request, server.recordedRequest) ***REMOVED***
		t.Fatal("Requests must be equal")
	***REMOVED***
***REMOVED***

func TestAuthZRequestPlugin(t *testing.T) ***REMOVED***
	server := authZPluginTestServer***REMOVED***t: t***REMOVED***
	server.start()
	defer server.stop()

	authZPlugin := createTestPlugin(t)

	request := Request***REMOVED***
		User:           "user",
		RequestBody:    []byte("sample body"),
		RequestURI:     "www.authz.com/auth",
		RequestMethod:  "GET",
		RequestHeaders: map[string]string***REMOVED***"header": "value"***REMOVED***,
	***REMOVED***
	server.replayResponse = Response***REMOVED***
		Allow: true,
		Msg:   "Sample message",
	***REMOVED***

	actualResponse, err := authZPlugin.AuthZRequest(&request)
	if err != nil ***REMOVED***
		t.Fatalf("Failed to authorize request %v", err)
	***REMOVED***

	if !reflect.DeepEqual(server.replayResponse, *actualResponse) ***REMOVED***
		t.Fatal("Response must be equal")
	***REMOVED***
	if !reflect.DeepEqual(request, server.recordedRequest) ***REMOVED***
		t.Fatal("Requests must be equal")
	***REMOVED***
***REMOVED***

func TestAuthZResponsePlugin(t *testing.T) ***REMOVED***
	server := authZPluginTestServer***REMOVED***t: t***REMOVED***
	server.start()
	defer server.stop()

	authZPlugin := createTestPlugin(t)

	request := Request***REMOVED***
		User:        "user",
		RequestURI:  "something.com/auth",
		RequestBody: []byte("sample body"),
	***REMOVED***
	server.replayResponse = Response***REMOVED***
		Allow: true,
		Msg:   "Sample message",
	***REMOVED***

	actualResponse, err := authZPlugin.AuthZResponse(&request)
	if err != nil ***REMOVED***
		t.Fatalf("Failed to authorize request %v", err)
	***REMOVED***

	if !reflect.DeepEqual(server.replayResponse, *actualResponse) ***REMOVED***
		t.Fatal("Response must be equal")
	***REMOVED***
	if !reflect.DeepEqual(request, server.recordedRequest) ***REMOVED***
		t.Fatal("Requests must be equal")
	***REMOVED***
***REMOVED***

func TestResponseModifier(t *testing.T) ***REMOVED***
	r := httptest.NewRecorder()
	m := NewResponseModifier(r)
	m.Header().Set("h1", "v1")
	m.Write([]byte("body"))
	m.WriteHeader(http.StatusInternalServerError)

	m.FlushAll()
	if r.Header().Get("h1") != "v1" ***REMOVED***
		t.Fatalf("Header value must exists %s", r.Header().Get("h1"))
	***REMOVED***
	if !reflect.DeepEqual(r.Body.Bytes(), []byte("body")) ***REMOVED***
		t.Fatalf("Body value must exists %s", r.Body.Bytes())
	***REMOVED***
	if r.Code != http.StatusInternalServerError ***REMOVED***
		t.Fatalf("Status code must be correct %d", r.Code)
	***REMOVED***
***REMOVED***

func TestDrainBody(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		length             int // length is the message length send to drainBody
		expectedBodyLength int // expectedBodyLength is the expected body length after drainBody is called
	***REMOVED******REMOVED***
		***REMOVED***10, 10***REMOVED***, // Small message size
		***REMOVED***maxBodySize - 1, maxBodySize - 1***REMOVED***, // Max message size
		***REMOVED***maxBodySize * 2, 0***REMOVED***,               // Large message size (skip copying body)

	***REMOVED***

	for _, test := range tests ***REMOVED***
		msg := strings.Repeat("a", test.length)
		body, closer, err := drainBody(ioutil.NopCloser(bytes.NewReader([]byte(msg))))
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if len(body) != test.expectedBodyLength ***REMOVED***
			t.Fatalf("Body must be copied, actual length: '%d'", len(body))
		***REMOVED***
		if closer == nil ***REMOVED***
			t.Fatal("Closer must not be nil")
		***REMOVED***
		modified, err := ioutil.ReadAll(closer)
		if err != nil ***REMOVED***
			t.Fatalf("Error must not be nil: '%v'", err)
		***REMOVED***
		if len(modified) != len(msg) ***REMOVED***
			t.Fatalf("Result should not be truncated. Original length: '%d', new length: '%d'", len(msg), len(modified))
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestResponseModifierOverride(t *testing.T) ***REMOVED***
	r := httptest.NewRecorder()
	m := NewResponseModifier(r)
	m.Header().Set("h1", "v1")
	m.Write([]byte("body"))
	m.WriteHeader(http.StatusInternalServerError)

	overrideHeader := make(http.Header)
	overrideHeader.Add("h1", "v2")
	overrideHeaderBytes, err := json.Marshal(overrideHeader)
	if err != nil ***REMOVED***
		t.Fatalf("override header failed %v", err)
	***REMOVED***

	m.OverrideHeader(overrideHeaderBytes)
	m.OverrideBody([]byte("override body"))
	m.OverrideStatusCode(http.StatusNotFound)
	m.FlushAll()
	if r.Header().Get("h1") != "v2" ***REMOVED***
		t.Fatalf("Header value must exists %s", r.Header().Get("h1"))
	***REMOVED***
	if !reflect.DeepEqual(r.Body.Bytes(), []byte("override body")) ***REMOVED***
		t.Fatalf("Body value must exists %s", r.Body.Bytes())
	***REMOVED***
	if r.Code != http.StatusNotFound ***REMOVED***
		t.Fatalf("Status code must be correct %d", r.Code)
	***REMOVED***
***REMOVED***

// createTestPlugin creates a new sample authorization plugin
func createTestPlugin(t *testing.T) *authorizationPlugin ***REMOVED***
	pwd, err := os.Getwd()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	client, err := plugins.NewClient("unix:///"+path.Join(pwd, pluginAddress), &tlsconfig.Options***REMOVED***InsecureSkipVerify: true***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("Failed to create client %v", err)
	***REMOVED***

	return &authorizationPlugin***REMOVED***name: "plugin", plugin: client***REMOVED***
***REMOVED***

// AuthZPluginTestServer is a simple server that implements the authZ plugin interface
type authZPluginTestServer struct ***REMOVED***
	listener net.Listener
	t        *testing.T
	// request stores the request sent from the daemon to the plugin
	recordedRequest Request
	// response stores the response sent from the plugin to the daemon
	replayResponse Response
	server         *httptest.Server
***REMOVED***

// start starts the test server that implements the plugin
func (t *authZPluginTestServer) start() ***REMOVED***
	r := mux.NewRouter()
	l, err := net.Listen("unix", pluginAddress)
	if err != nil ***REMOVED***
		t.t.Fatal(err)
	***REMOVED***
	t.listener = l
	r.HandleFunc("/Plugin.Activate", t.activate)
	r.HandleFunc("/"+AuthZApiRequest, t.auth)
	r.HandleFunc("/"+AuthZApiResponse, t.auth)
	t.server = &httptest.Server***REMOVED***
		Listener: l,
		Config: &http.Server***REMOVED***
			Handler: r,
			Addr:    pluginAddress,
		***REMOVED***,
	***REMOVED***
	t.server.Start()
***REMOVED***

// stop stops the test server that implements the plugin
func (t *authZPluginTestServer) stop() ***REMOVED***
	t.server.Close()
	os.Remove(pluginAddress)
	if t.listener != nil ***REMOVED***
		t.listener.Close()
	***REMOVED***
***REMOVED***

// auth is a used to record/replay the authentication api messages
func (t *authZPluginTestServer) auth(w http.ResponseWriter, r *http.Request) ***REMOVED***
	t.recordedRequest = Request***REMOVED******REMOVED***
	body, err := ioutil.ReadAll(r.Body)
	if err != nil ***REMOVED***
		t.t.Fatal(err)
	***REMOVED***
	r.Body.Close()
	json.Unmarshal(body, &t.recordedRequest)
	b, err := json.Marshal(t.replayResponse)
	if err != nil ***REMOVED***
		t.t.Fatal(err)
	***REMOVED***
	w.Write(b)
***REMOVED***

func (t *authZPluginTestServer) activate(w http.ResponseWriter, r *http.Request) ***REMOVED***
	b, err := json.Marshal(plugins.Manifest***REMOVED***Implements: []string***REMOVED***AuthZApiImplements***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.t.Fatal(err)
	***REMOVED***
	w.Write(b)
***REMOVED***
