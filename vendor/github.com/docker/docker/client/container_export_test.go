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

func TestContainerExportError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ContainerExport(context.Background(), "nothing")
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestContainerExport(t *testing.T) ***REMOVED***
	expectedURL := "/containers/container_id/export"
	client := &Client***REMOVED***
		client: newMockClient(func(r *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(r.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, r.URL)
			***REMOVED***

			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("response"))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***
	body, err := client.ContainerExport(context.Background(), "container_id")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer body.Close()
	content, err := ioutil.ReadAll(body)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if string(content) != "response" ***REMOVED***
		t.Fatalf("expected response to contain 'response', got %s", string(content))
	***REMOVED***
***REMOVED***
