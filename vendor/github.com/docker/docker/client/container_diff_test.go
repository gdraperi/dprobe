package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/container"
	"golang.org/x/net/context"
)

func TestContainerDiffError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ContainerDiff(context.Background(), "nothing")
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***

***REMOVED***

func TestContainerDiff(t *testing.T) ***REMOVED***
	expectedURL := "/containers/container_id/changes"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			b, err := json.Marshal([]container.ContainerChangeResponseItem***REMOVED***
				***REMOVED***
					Kind: 0,
					Path: "/path/1",
				***REMOVED***,
				***REMOVED***
					Kind: 1,
					Path: "/path/2",
				***REMOVED***,
			***REMOVED***)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(b)),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	changes, err := client.ContainerDiff(context.Background(), "container_id")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if len(changes) != 2 ***REMOVED***
		t.Fatalf("expected an array of 2 changes, got %v", changes)
	***REMOVED***
***REMOVED***
