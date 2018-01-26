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
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestImageRemoveError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, err := client.ImageRemove(context.Background(), "image_id", types.ImageRemoveOptions***REMOVED******REMOVED***)
	assert.EqualError(t, err, "Error response from daemon: Server error")
***REMOVED***

func TestImageRemoveImageNotFound(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusNotFound, "missing")),
	***REMOVED***

	_, err := client.ImageRemove(context.Background(), "unknown", types.ImageRemoveOptions***REMOVED******REMOVED***)
	assert.EqualError(t, err, "Error: No such image: unknown")
	assert.True(t, IsErrNotFound(err))
***REMOVED***

func TestImageRemove(t *testing.T) ***REMOVED***
	expectedURL := "/images/image_id"
	removeCases := []struct ***REMOVED***
		force               bool
		pruneChildren       bool
		expectedQueryParams map[string]string
	***REMOVED******REMOVED***
		***REMOVED***
			force:         false,
			pruneChildren: false,
			expectedQueryParams: map[string]string***REMOVED***
				"force":   "",
				"noprune": "1",
			***REMOVED***,
		***REMOVED***, ***REMOVED***
			force:         true,
			pruneChildren: true,
			expectedQueryParams: map[string]string***REMOVED***
				"force":   "1",
				"noprune": "",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for _, removeCase := range removeCases ***REMOVED***
		client := &Client***REMOVED***
			client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
				if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
					return nil, fmt.Errorf("expected URL '%s', got '%s'", expectedURL, req.URL)
				***REMOVED***
				if req.Method != "DELETE" ***REMOVED***
					return nil, fmt.Errorf("expected DELETE method, got %s", req.Method)
				***REMOVED***
				query := req.URL.Query()
				for key, expected := range removeCase.expectedQueryParams ***REMOVED***
					actual := query.Get(key)
					if actual != expected ***REMOVED***
						return nil, fmt.Errorf("%s not set in URL query properly. Expected '%s', got %s", key, expected, actual)
					***REMOVED***
				***REMOVED***
				b, err := json.Marshal([]types.ImageDeleteResponseItem***REMOVED***
					***REMOVED***
						Untagged: "image_id1",
					***REMOVED***,
					***REMOVED***
						Deleted: "image_id",
					***REMOVED***,
				***REMOVED***)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***

				return &http.Response***REMOVED***
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(b)),
				***REMOVED***, nil
			***REMOVED***),
		***REMOVED***
		imageDeletes, err := client.ImageRemove(context.Background(), "image_id", types.ImageRemoveOptions***REMOVED***
			Force:         removeCase.force,
			PruneChildren: removeCase.pruneChildren,
		***REMOVED***)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if len(imageDeletes) != 2 ***REMOVED***
			t.Fatalf("expected 2 deleted images, got %v", imageDeletes)
		***REMOVED***
	***REMOVED***
***REMOVED***
