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

func TestCheckpointCreateError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	err := client.CheckpointCreate(context.Background(), "nothing", types.CheckpointCreateOptions***REMOVED***
		CheckpointID: "noting",
		Exit:         true,
	***REMOVED***)

	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestCheckpointCreate(t *testing.T) ***REMOVED***
	expectedContainerID := "container_id"
	expectedCheckpointID := "checkpoint_id"
	expectedURL := "/containers/container_id/checkpoints"

	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***

			if req.Method != "POST" ***REMOVED***
				return nil, fmt.Errorf("expected POST method, got %s", req.Method)
			***REMOVED***

			createOptions := &types.CheckpointCreateOptions***REMOVED******REMOVED***
			if err := json.NewDecoder(req.Body).Decode(createOptions); err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			if createOptions.CheckpointID != expectedCheckpointID ***REMOVED***
				return nil, fmt.Errorf("expected CheckpointID to be 'checkpoint_id', got %v", createOptions.CheckpointID)
			***REMOVED***

			if !createOptions.Exit ***REMOVED***
				return nil, fmt.Errorf("expected Exit to be true")
			***REMOVED***

			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	err := client.CheckpointCreate(context.Background(), expectedContainerID, types.CheckpointCreateOptions***REMOVED***
		CheckpointID: expectedCheckpointID,
		Exit:         true,
	***REMOVED***)

	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
