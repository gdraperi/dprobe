package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

func TestCheckpointListError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, err := client.CheckpointList(context.Background(), "container_id", types.CheckpointListOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestCheckpointList(t *testing.T) ***REMOVED***
	expectedURL := "/containers/container_id/checkpoints"

	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			content, err := json.Marshal([]types.Checkpoint***REMOVED***
				***REMOVED***
					Name: "checkpoint",
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

	checkpoints, err := client.CheckpointList(context.Background(), "container_id", types.CheckpointListOptions***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if len(checkpoints) != 1 ***REMOVED***
		t.Fatalf("expected 1 checkpoint, got %v", checkpoints)
	***REMOVED***
***REMOVED***

func TestCheckpointListContainerNotFound(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusNotFound, "Server error")),
	***REMOVED***

	_, err := client.CheckpointList(context.Background(), "unknown", types.CheckpointListOptions***REMOVED******REMOVED***)
	if err == nil || !IsErrNotFound(err) ***REMOVED***
		t.Fatalf("expected a containerNotFound error, got %v", err)
	***REMOVED***
***REMOVED***
