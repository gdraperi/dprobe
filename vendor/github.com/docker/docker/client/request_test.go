package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// TestSetHostHeader should set fake host for local communications, set real host
// for normal communications.
func TestSetHostHeader(t *testing.T) ***REMOVED***
	testURL := "/test"
	testCases := []struct ***REMOVED***
		host            string
		expectedHost    string
		expectedURLHost string
	***REMOVED******REMOVED***
		***REMOVED***
			"unix:///var/run/docker.sock",
			"docker",
			"/var/run/docker.sock",
		***REMOVED***,
		***REMOVED***
			"npipe:////./pipe/docker_engine",
			"docker",
			"//./pipe/docker_engine",
		***REMOVED***,
		***REMOVED***
			"tcp://0.0.0.0:4243",
			"",
			"0.0.0.0:4243",
		***REMOVED***,
		***REMOVED***
			"tcp://localhost:4243",
			"",
			"localhost:4243",
		***REMOVED***,
	***REMOVED***

	for c, test := range testCases ***REMOVED***
		proto, addr, basePath, err := ParseHost(test.host)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		client := &Client***REMOVED***
			client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
				if !strings.HasPrefix(req.URL.Path, testURL) ***REMOVED***
					return nil, fmt.Errorf("Test Case #%d: Expected URL %q, got %q", c, testURL, req.URL)
				***REMOVED***
				if req.Host != test.expectedHost ***REMOVED***
					return nil, fmt.Errorf("Test Case #%d: Expected host %q, got %q", c, test.expectedHost, req.Host)
				***REMOVED***
				if req.URL.Host != test.expectedURLHost ***REMOVED***
					return nil, fmt.Errorf("Test Case #%d: Expected URL host %q, got %q", c, test.expectedURLHost, req.URL.Host)
				***REMOVED***
				return &http.Response***REMOVED***
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(([]byte("")))),
				***REMOVED***, nil
			***REMOVED***),

			proto:    proto,
			addr:     addr,
			basePath: basePath,
		***REMOVED***

		_, err = client.sendRequest(context.Background(), "GET", testURL, nil, nil, nil)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

// TestPlainTextError tests the server returning an error in plain text for
// backwards compatibility with API versions <1.24. All other tests use
// errors returned as JSON
func TestPlainTextError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(plainTextErrorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ContainerList(context.Background(), types.ContainerListOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***
