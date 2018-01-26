package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestConfigRemoveUnsupported(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.29",
		client:  &http.Client***REMOVED******REMOVED***,
	***REMOVED***
	err := client.ConfigRemove(context.Background(), "config_id")
	assert.EqualError(t, err, `"config remove" requires API version 1.30, but the Docker daemon API version is 1.29`)
***REMOVED***

func TestConfigRemoveError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.30",
		client:  newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	err := client.ConfigRemove(context.Background(), "config_id")
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestConfigRemove(t *testing.T) ***REMOVED***
	expectedURL := "/v1.30/configs/config_id"

	client := &Client***REMOVED***
		version: "1.30",
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

	err := client.ConfigRemove(context.Background(), "config_id")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
