package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
)

func TestServiceUpdateError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, err := client.ServiceUpdate(context.Background(), "service_id", swarm.Version***REMOVED******REMOVED***, swarm.ServiceSpec***REMOVED******REMOVED***, types.ServiceUpdateOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestServiceUpdate(t *testing.T) ***REMOVED***
	expectedURL := "/services/service_id/update"

	updateCases := []struct ***REMOVED***
		swarmVersion    swarm.Version
		expectedVersion string
	***REMOVED******REMOVED***
		***REMOVED***
			expectedVersion: "0",
		***REMOVED***,
		***REMOVED***
			swarmVersion: swarm.Version***REMOVED***
				Index: 0,
			***REMOVED***,
			expectedVersion: "0",
		***REMOVED***,
		***REMOVED***
			swarmVersion: swarm.Version***REMOVED***
				Index: 10,
			***REMOVED***,
			expectedVersion: "10",
		***REMOVED***,
	***REMOVED***

	for _, updateCase := range updateCases ***REMOVED***
		client := &Client***REMOVED***
			client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
				if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
					return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
				***REMOVED***
				if req.Method != "POST" ***REMOVED***
					return nil, fmt.Errorf("expected POST method, got %s", req.Method)
				***REMOVED***
				version := req.URL.Query().Get("version")
				if version != updateCase.expectedVersion ***REMOVED***
					return nil, fmt.Errorf("version not set in URL query properly, expected '%s', got %s", updateCase.expectedVersion, version)
				***REMOVED***
				return &http.Response***REMOVED***
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("***REMOVED******REMOVED***"))),
				***REMOVED***, nil
			***REMOVED***),
		***REMOVED***

		_, err := client.ServiceUpdate(context.Background(), "service_id", updateCase.swarmVersion, swarm.ServiceSpec***REMOVED******REMOVED***, types.ServiceUpdateOptions***REMOVED******REMOVED***)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***
