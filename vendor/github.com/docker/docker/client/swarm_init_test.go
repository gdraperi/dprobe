package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types/swarm"
)

func TestSwarmInitError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, err := client.SwarmInit(context.Background(), swarm.InitRequest***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestSwarmInit(t *testing.T) ***REMOVED***
	expectedURL := "/swarm/init"

	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			if req.Method != "POST" ***REMOVED***
				return nil, fmt.Errorf("expected POST method, got %s", req.Method)
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(`"body"`))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	resp, err := client.SwarmInit(context.Background(), swarm.InitRequest***REMOVED***
		ListenAddr: "0.0.0.0:2377",
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if resp != "body" ***REMOVED***
		t.Fatalf("Expected 'body', got %s", resp)
	***REMOVED***
***REMOVED***
