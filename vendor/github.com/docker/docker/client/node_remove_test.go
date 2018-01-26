package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"

	"golang.org/x/net/context"
)

func TestNodeRemoveError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	err := client.NodeRemove(context.Background(), "node_id", types.NodeRemoveOptions***REMOVED***Force: false***REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestNodeRemove(t *testing.T) ***REMOVED***
	expectedURL := "/nodes/node_id"

	removeCases := []struct ***REMOVED***
		force         bool
		expectedForce string
	***REMOVED******REMOVED***
		***REMOVED***
			expectedForce: "",
		***REMOVED***,
		***REMOVED***
			force:         true,
			expectedForce: "1",
		***REMOVED***,
	***REMOVED***

	for _, removeCase := range removeCases ***REMOVED***
		client := &Client***REMOVED***
			client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
				if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
					return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
				***REMOVED***
				if req.Method != "DELETE" ***REMOVED***
					return nil, fmt.Errorf("expected DELETE method, got %s", req.Method)
				***REMOVED***
				force := req.URL.Query().Get("force")
				if force != removeCase.expectedForce ***REMOVED***
					return nil, fmt.Errorf("force not set in URL query properly. expected '%s', got %s", removeCase.expectedForce, force)
				***REMOVED***

				return &http.Response***REMOVED***
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("body"))),
				***REMOVED***, nil
			***REMOVED***),
		***REMOVED***

		err := client.NodeRemove(context.Background(), "node_id", types.NodeRemoveOptions***REMOVED***Force: removeCase.force***REMOVED***)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***
