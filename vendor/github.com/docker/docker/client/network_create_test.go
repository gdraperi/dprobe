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

func TestNetworkCreateError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, err := client.NetworkCreate(context.Background(), "mynetwork", types.NetworkCreate***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestNetworkCreate(t *testing.T) ***REMOVED***
	expectedURL := "/networks/create"

	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***

			if req.Method != "POST" ***REMOVED***
				return nil, fmt.Errorf("expected POST method, got %s", req.Method)
			***REMOVED***

			content, err := json.Marshal(types.NetworkCreateResponse***REMOVED***
				ID:      "network_id",
				Warning: "warning",
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

	networkResponse, err := client.NetworkCreate(context.Background(), "mynetwork", types.NetworkCreate***REMOVED***
		CheckDuplicate: true,
		Driver:         "mydriver",
		EnableIPv6:     true,
		Internal:       true,
		Options: map[string]string***REMOVED***
			"opt-key": "opt-value",
		***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if networkResponse.ID != "network_id" ***REMOVED***
		t.Fatalf("expected networkResponse.ID to be 'network_id', got %s", networkResponse.ID)
	***REMOVED***
	if networkResponse.Warning != "warning" ***REMOVED***
		t.Fatalf("expected networkResponse.Warning to be 'warning', got %s", networkResponse.Warning)
	***REMOVED***
***REMOVED***
