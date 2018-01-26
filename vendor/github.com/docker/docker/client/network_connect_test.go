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
	"github.com/docker/docker/api/types/network"
)

func TestNetworkConnectError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	err := client.NetworkConnect(context.Background(), "network_id", "container_id", nil)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestNetworkConnectEmptyNilEndpointSettings(t *testing.T) ***REMOVED***
	expectedURL := "/networks/network_id/connect"

	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***

			if req.Method != "POST" ***REMOVED***
				return nil, fmt.Errorf("expected POST method, got %s", req.Method)
			***REMOVED***

			var connect types.NetworkConnect
			if err := json.NewDecoder(req.Body).Decode(&connect); err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			if connect.Container != "container_id" ***REMOVED***
				return nil, fmt.Errorf("expected 'container_id', got %s", connect.Container)
			***REMOVED***

			if connect.EndpointConfig != nil ***REMOVED***
				return nil, fmt.Errorf("expected connect.EndpointConfig to be nil, got %v", connect.EndpointConfig)
			***REMOVED***

			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	err := client.NetworkConnect(context.Background(), "network_id", "container_id", nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestNetworkConnect(t *testing.T) ***REMOVED***
	expectedURL := "/networks/network_id/connect"

	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***

			if req.Method != "POST" ***REMOVED***
				return nil, fmt.Errorf("expected POST method, got %s", req.Method)
			***REMOVED***

			var connect types.NetworkConnect
			if err := json.NewDecoder(req.Body).Decode(&connect); err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			if connect.Container != "container_id" ***REMOVED***
				return nil, fmt.Errorf("expected 'container_id', got %s", connect.Container)
			***REMOVED***

			if connect.EndpointConfig == nil ***REMOVED***
				return nil, fmt.Errorf("expected connect.EndpointConfig to be not nil, got %v", connect.EndpointConfig)
			***REMOVED***

			if connect.EndpointConfig.NetworkID != "NetworkID" ***REMOVED***
				return nil, fmt.Errorf("expected 'NetworkID', got %s", connect.EndpointConfig.NetworkID)
			***REMOVED***

			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	err := client.NetworkConnect(context.Background(), "network_id", "container_id", &network.EndpointSettings***REMOVED***
		NetworkID: "NetworkID",
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
