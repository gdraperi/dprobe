package registry

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestEndpointParse(t *testing.T) ***REMOVED***
	testData := []struct ***REMOVED***
		str      string
		expected string
	***REMOVED******REMOVED***
		***REMOVED***IndexServer, IndexServer***REMOVED***,
		***REMOVED***"http://0.0.0.0:5000/v1/", "http://0.0.0.0:5000/v1/"***REMOVED***,
		***REMOVED***"http://0.0.0.0:5000", "http://0.0.0.0:5000/v1/"***REMOVED***,
		***REMOVED***"0.0.0.0:5000", "https://0.0.0.0:5000/v1/"***REMOVED***,
		***REMOVED***"http://0.0.0.0:5000/nonversion/", "http://0.0.0.0:5000/nonversion/v1/"***REMOVED***,
		***REMOVED***"http://0.0.0.0:5000/v0/", "http://0.0.0.0:5000/v0/v1/"***REMOVED***,
	***REMOVED***
	for _, td := range testData ***REMOVED***
		e, err := newV1EndpointFromStr(td.str, nil, "", nil)
		if err != nil ***REMOVED***
			t.Errorf("%q: %s", td.str, err)
		***REMOVED***
		if e == nil ***REMOVED***
			t.Logf("something's fishy, endpoint for %q is nil", td.str)
			continue
		***REMOVED***
		if e.String() != td.expected ***REMOVED***
			t.Errorf("expected %q, got %q", td.expected, e.String())
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestEndpointParseInvalid(t *testing.T) ***REMOVED***
	testData := []string***REMOVED***
		"http://0.0.0.0:5000/v2/",
	***REMOVED***
	for _, td := range testData ***REMOVED***
		e, err := newV1EndpointFromStr(td, nil, "", nil)
		if err == nil ***REMOVED***
			t.Errorf("expected error parsing %q: parsed as %q", td, e)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Ensure that a registry endpoint that responds with a 401 only is determined
// to be a valid v1 registry endpoint
func TestValidateEndpoint(t *testing.T) ***REMOVED***
	requireBasicAuthHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Add("WWW-Authenticate", `Basic realm="localhost"`)
		w.WriteHeader(http.StatusUnauthorized)
	***REMOVED***)

	// Make a test server which should validate as a v1 server.
	testServer := httptest.NewServer(requireBasicAuthHandler)
	defer testServer.Close()

	testServerURL, err := url.Parse(testServer.URL)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	testEndpoint := V1Endpoint***REMOVED***
		URL:    testServerURL,
		client: HTTPClient(NewTransport(nil)),
	***REMOVED***

	if err = validateEndpoint(&testEndpoint); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if testEndpoint.URL.Scheme != "http" ***REMOVED***
		t.Fatalf("expecting to validate endpoint as http, got url %s", testEndpoint.String())
	***REMOVED***
***REMOVED***
