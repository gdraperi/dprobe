package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/context"
)

func TestContainerRestartError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	timeout := 0 * time.Second
	err := client.ContainerRestart(context.Background(), "nothing", &timeout)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestContainerRestart(t *testing.T) ***REMOVED***
	expectedURL := "/containers/container_id/restart"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			t := req.URL.Query().Get("t")
			if t != "100" ***REMOVED***
				return nil, fmt.Errorf("t (timeout) not set in URL query properly. Expected '100', got %s", t)
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***
	timeout := 100 * time.Second
	err := client.ContainerRestart(context.Background(), "container_id", &timeout)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
