package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

func TestSwarmInspectError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, err := client.SwarmInspect(context.Background())
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestSwarmInspect(t *testing.T) ***REMOVED***
	expectedURL := "/swarm"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			content, err := json.Marshal(swarm.Swarm***REMOVED***
				ClusterInfo: swarm.ClusterInfo***REMOVED***
					ID: "swarm_id",
				***REMOVED***,
			***REMOVED***)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(content)),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	swarmInspect, err := client.SwarmInspect(context.Background())
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if swarmInspect.ID != "swarm_id" ***REMOVED***
		t.Fatalf("expected `swarm_id`, got %s", swarmInspect.ID)
	***REMOVED***
***REMOVED***
