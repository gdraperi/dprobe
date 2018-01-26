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

func TestNetworkDisconnectError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	err := client.NetworkDisconnect(context.Background(), "network_id", "container_id", false)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestNetworkDisconnect(t *testing.T) ***REMOVED***
	expectedURL := "/networks/network_id/disconnect"

	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***

			if req.Method != "POST" ***REMOVED***
				return nil, fmt.Errorf("expected POST method, got %s", req.Method)
			***REMOVED***

			var disconnect types.NetworkDisconnect
			if err := json.NewDecoder(req.Body).Decode(&disconnect); err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			if disconnect.Container != "container_id" ***REMOVED***
				return nil, fmt.Errorf("expected 'container_id', got %s", disconnect.Container)
			***REMOVED***

			if !disconnect.Force ***REMOVED***
				return nil, fmt.Errorf("expected Force to be true, got %v", disconnect.Force)
			***REMOVED***

			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	err := client.NetworkDisconnect(context.Background(), "network_id", "container_id", true)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
