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
	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

func TestServiceInspectError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, _, err := client.ServiceInspectWithRaw(context.Background(), "nothing", types.ServiceInspectOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestServiceInspectServiceNotFound(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusNotFound, "Server error")),
	***REMOVED***

	_, _, err := client.ServiceInspectWithRaw(context.Background(), "unknown", types.ServiceInspectOptions***REMOVED******REMOVED***)
	if err == nil || !IsErrNotFound(err) ***REMOVED***
		t.Fatalf("expected a serviceNotFoundError error, got %v", err)
	***REMOVED***
***REMOVED***

func TestServiceInspect(t *testing.T) ***REMOVED***
	expectedURL := "/services/service_id"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			content, err := json.Marshal(swarm.Service***REMOVED***
				ID: "service_id",
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

	serviceInspect, _, err := client.ServiceInspectWithRaw(context.Background(), "service_id", types.ServiceInspectOptions***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if serviceInspect.ID != "service_id" ***REMOVED***
		t.Fatalf("expected `service_id`, got %s", serviceInspect.ID)
	***REMOVED***
***REMOVED***
