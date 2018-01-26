package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"golang.org/x/net/context"
)

func TestSwarmLeaveError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	err := client.SwarmLeave(context.Background(), false)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestSwarmLeave(t *testing.T) ***REMOVED***
	expectedURL := "/swarm/leave"

	leaveCases := []struct ***REMOVED***
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

	for _, leaveCase := range leaveCases ***REMOVED***
		client := &Client***REMOVED***
			client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
				if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
					return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
				***REMOVED***
				if req.Method != "POST" ***REMOVED***
					return nil, fmt.Errorf("expected POST method, got %s", req.Method)
				***REMOVED***
				force := req.URL.Query().Get("force")
				if force != leaveCase.expectedForce ***REMOVED***
					return nil, fmt.Errorf("force not set in URL query properly. expected '%s', got %s", leaveCase.expectedForce, force)
				***REMOVED***
				return &http.Response***REMOVED***
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
				***REMOVED***, nil
			***REMOVED***),
		***REMOVED***

		err := client.SwarmLeave(context.Background(), leaveCase.force)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***
