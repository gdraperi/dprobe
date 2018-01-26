package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
)

func TestContainerStartError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	err := client.ContainerStart(context.Background(), "nothing", types.ContainerStartOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestContainerStart(t *testing.T) ***REMOVED***
	expectedURL := "/containers/container_id/start"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			// we're not expecting any payload, but if one is supplied, check it is valid.
			if req.Header.Get("Content-Type") == "application/json" ***REMOVED***
				var startConfig interface***REMOVED******REMOVED***
				if err := json.NewDecoder(req.Body).Decode(&startConfig); err != nil ***REMOVED***
					return nil, fmt.Errorf("Unable to parse json: %s", err)
				***REMOVED***
			***REMOVED***

			checkpoint := req.URL.Query().Get("checkpoint")
			if checkpoint != "checkpoint_id" ***REMOVED***
				return nil, fmt.Errorf("checkpoint not set in URL query properly. Expected 'checkpoint_id', got %s", checkpoint)
			***REMOVED***

			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	err := client.ContainerStart(context.Background(), "container_id", types.ContainerStartOptions***REMOVED***CheckpointID: "checkpoint_id"***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
