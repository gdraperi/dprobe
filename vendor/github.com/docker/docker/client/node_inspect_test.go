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

func TestNodeInspectError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, _, err := client.NodeInspectWithRaw(context.Background(), "nothing")
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestNodeInspectNodeNotFound(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusNotFound, "Server error")),
	***REMOVED***

	_, _, err := client.NodeInspectWithRaw(context.Background(), "unknown")
	if err == nil || !IsErrNotFound(err) ***REMOVED***
		t.Fatalf("expected a nodeNotFoundError error, got %v", err)
	***REMOVED***
***REMOVED***

func TestNodeInspect(t *testing.T) ***REMOVED***
	expectedURL := "/nodes/node_id"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			content, err := json.Marshal(swarm.Node***REMOVED***
				ID: "node_id",
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

	nodeInspect, _, err := client.NodeInspectWithRaw(context.Background(), "node_id")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if nodeInspect.ID != "node_id" ***REMOVED***
		t.Fatalf("expected `node_id`, got %s", nodeInspect.ID)
	***REMOVED***
***REMOVED***
