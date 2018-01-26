package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestConfigUpdateUnsupported(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.29",
		client:  &http.Client***REMOVED******REMOVED***,
	***REMOVED***
	err := client.ConfigUpdate(context.Background(), "config_id", swarm.Version***REMOVED******REMOVED***, swarm.ConfigSpec***REMOVED******REMOVED***)
	assert.EqualError(t, err, `"config update" requires API version 1.30, but the Docker daemon API version is 1.29`)
***REMOVED***

func TestConfigUpdateError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.30",
		client:  newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	err := client.ConfigUpdate(context.Background(), "config_id", swarm.Version***REMOVED******REMOVED***, swarm.ConfigSpec***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestConfigUpdate(t *testing.T) ***REMOVED***
	expectedURL := "/v1.30/configs/config_id/update"

	client := &Client***REMOVED***
		version: "1.30",
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			if req.Method != "POST" ***REMOVED***
				return nil, fmt.Errorf("expected POST method, got %s", req.Method)
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("body"))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	err := client.ConfigUpdate(context.Background(), "config_id", swarm.Version***REMOVED******REMOVED***, swarm.ConfigSpec***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
