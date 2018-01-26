package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"golang.org/x/net/context"
)

func TestNetworkRemoveError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	err := client.NetworkRemove(context.Background(), "network_id")
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestNetworkRemove(t *testing.T) ***REMOVED***
	expectedURL := "/networks/network_id"

	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			if req.Method != "DELETE" ***REMOVED***
				return nil, fmt.Errorf("expected DELETE method, got %s", req.Method)
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("body"))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	err := client.NetworkRemove(context.Background(), "network_id")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
