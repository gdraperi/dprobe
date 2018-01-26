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

func TestContainerRenameError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	err := client.ContainerRename(context.Background(), "nothing", "newNothing")
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestContainerRename(t *testing.T) ***REMOVED***
	expectedURL := "/containers/container_id/rename"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			name := req.URL.Query().Get("name")
			if name != "newName" ***REMOVED***
				return nil, fmt.Errorf("name not set in URL query properly. Expected 'newName', got %s", name)
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	err := client.ContainerRename(context.Background(), "container_id", "newName")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
