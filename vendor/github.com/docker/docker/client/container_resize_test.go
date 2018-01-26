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

func TestContainerResizeError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	err := client.ContainerResize(context.Background(), "container_id", types.ResizeOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestContainerExecResizeError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	err := client.ContainerExecResize(context.Background(), "exec_id", types.ResizeOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestContainerResize(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(resizeTransport("/containers/container_id/resize")),
	***REMOVED***

	err := client.ContainerResize(context.Background(), "container_id", types.ResizeOptions***REMOVED***
		Height: 500,
		Width:  600,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestContainerExecResize(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(resizeTransport("/exec/exec_id/resize")),
	***REMOVED***

	err := client.ContainerExecResize(context.Background(), "exec_id", types.ResizeOptions***REMOVED***
		Height: 500,
		Width:  600,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func resizeTransport(expectedURL string) func(req *http.Request) (*http.Response, error) ***REMOVED***
	return func(req *http.Request) (*http.Response, error) ***REMOVED***
		if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
			return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
		***REMOVED***
		query := req.URL.Query()
		h := query.Get("h")
		if h != "500" ***REMOVED***
			return nil, fmt.Errorf("h not set in URL query properly. Expected '500', got %s", h)
		***REMOVED***
		w := query.Get("w")
		if w != "600" ***REMOVED***
			return nil, fmt.Errorf("w not set in URL query properly. Expected '600', got %s", w)
		***REMOVED***
		return &http.Response***REMOVED***
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
		***REMOVED***, nil
	***REMOVED***
***REMOVED***
