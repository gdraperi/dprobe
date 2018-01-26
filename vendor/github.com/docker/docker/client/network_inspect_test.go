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
	"github.com/docker/docker/api/types/network"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestNetworkInspectError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, err := client.NetworkInspect(context.Background(), "nothing", types.NetworkInspectOptions***REMOVED******REMOVED***)
	assert.EqualError(t, err, "Error response from daemon: Server error")
***REMOVED***

func TestNetworkInspectNotFoundError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusNotFound, "missing")),
	***REMOVED***

	_, err := client.NetworkInspect(context.Background(), "unknown", types.NetworkInspectOptions***REMOVED******REMOVED***)
	assert.EqualError(t, err, "Error: No such network: unknown")
	assert.True(t, IsErrNotFound(err))
***REMOVED***

func TestNetworkInspect(t *testing.T) ***REMOVED***
	expectedURL := "/networks/network_id"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			if req.Method != "GET" ***REMOVED***
				return nil, fmt.Errorf("expected GET method, got %s", req.Method)
			***REMOVED***

			var (
				content []byte
				err     error
			)
			if strings.Contains(req.URL.RawQuery, "scope=global") ***REMOVED***
				return &http.Response***REMOVED***
					StatusCode: http.StatusNotFound,
					Body:       ioutil.NopCloser(bytes.NewReader(content)),
				***REMOVED***, nil
			***REMOVED***

			if strings.Contains(req.URL.RawQuery, "verbose=true") ***REMOVED***
				s := map[string]network.ServiceInfo***REMOVED***
					"web": ***REMOVED******REMOVED***,
				***REMOVED***
				content, err = json.Marshal(types.NetworkResource***REMOVED***
					Name:     "mynetwork",
					Services: s,
				***REMOVED***)
			***REMOVED*** else ***REMOVED***
				content, err = json.Marshal(types.NetworkResource***REMOVED***
					Name: "mynetwork",
				***REMOVED***)
			***REMOVED***
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(content)),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	r, err := client.NetworkInspect(context.Background(), "network_id", types.NetworkInspectOptions***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if r.Name != "mynetwork" ***REMOVED***
		t.Fatalf("expected `mynetwork`, got %s", r.Name)
	***REMOVED***

	r, err = client.NetworkInspect(context.Background(), "network_id", types.NetworkInspectOptions***REMOVED***Verbose: true***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if r.Name != "mynetwork" ***REMOVED***
		t.Fatalf("expected `mynetwork`, got %s", r.Name)
	***REMOVED***
	_, ok := r.Services["web"]
	if !ok ***REMOVED***
		t.Fatalf("expected service `web` missing in the verbose output")
	***REMOVED***

	_, err = client.NetworkInspect(context.Background(), "network_id", types.NetworkInspectOptions***REMOVED***Scope: "global"***REMOVED***)
	assert.EqualError(t, err, "Error: No such network: network_id")
***REMOVED***
