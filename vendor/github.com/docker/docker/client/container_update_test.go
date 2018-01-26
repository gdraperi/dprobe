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

func TestContainerUpdateError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ContainerUpdate(context.Background(), "nothing", container.UpdateConfig***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestContainerUpdate(t *testing.T) ***REMOVED***
	expectedURL := "/containers/container_id/update"

	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***

			b, err := json.Marshal(container.ContainerUpdateOKBody***REMOVED******REMOVED***)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(b)),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	_, err := client.ContainerUpdate(context.Background(), "container_id", container.UpdateConfig***REMOVED***
		Resources: container.Resources***REMOVED***
			CPUPeriod: 1,
		***REMOVED***,
		RestartPolicy: container.RestartPolicy***REMOVED***
			Name: "always",
		***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
