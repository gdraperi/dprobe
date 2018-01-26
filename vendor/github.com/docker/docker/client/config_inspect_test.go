package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestConfigInspectUnsupported(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.29",
		client:  &http.Client***REMOVED******REMOVED***,
	***REMOVED***
	_, _, err := client.ConfigInspectWithRaw(context.Background(), "nothing")
	assert.EqualError(t, err, `"config inspect" requires API version 1.30, but the Docker daemon API version is 1.29`)
***REMOVED***

func TestConfigInspectError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.30",
		client:  newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, _, err := client.ConfigInspectWithRaw(context.Background(), "nothing")
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestConfigInspectConfigNotFound(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.30",
		client:  newMockClient(errorMock(http.StatusNotFound, "Server error")),
	***REMOVED***

	_, _, err := client.ConfigInspectWithRaw(context.Background(), "unknown")
	if err == nil || !IsErrNotFound(err) ***REMOVED***
		t.Fatalf("expected a configNotFoundError error, got %v", err)
	***REMOVED***
***REMOVED***

func TestConfigInspect(t *testing.T) ***REMOVED***
	expectedURL := "/v1.30/configs/config_id"
	client := &Client***REMOVED***
		version: "1.30",
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			content, err := json.Marshal(swarm.Config***REMOVED***
				ID: "config_id",
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

	configInspect, _, err := client.ConfigInspectWithRaw(context.Background(), "config_id")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if configInspect.ID != "config_id" ***REMOVED***
		t.Fatalf("expected `config_id`, got %s", configInspect.ID)
	***REMOVED***
***REMOVED***
