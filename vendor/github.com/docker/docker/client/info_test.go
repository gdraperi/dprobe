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

func TestInfoServerError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.Info(context.Background())
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestInfoInvalidResponseJSONError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("invalid json"))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***
	_, err := client.Info(context.Background())
	if err == nil || !strings.Contains(err.Error(), "invalid character") ***REMOVED***
		t.Fatalf("expected a 'invalid character' error, got %v", err)
	***REMOVED***
***REMOVED***

func TestInfo(t *testing.T) ***REMOVED***
	expectedURL := "/info"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			info := &types.Info***REMOVED***
				ID:         "daemonID",
				Containers: 3,
			***REMOVED***
			b, err := json.Marshal(info)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(b)),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	info, err := client.Info(context.Background())
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if info.ID != "daemonID" ***REMOVED***
		t.Fatalf("expected daemonID, got %s", info.ID)
	***REMOVED***

	if info.Containers != 3 ***REMOVED***
		t.Fatalf("expected 3 containers, got %d", info.Containers)
	***REMOVED***
***REMOVED***
