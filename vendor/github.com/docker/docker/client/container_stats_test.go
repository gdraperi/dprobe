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

func TestContainerStatsError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ContainerStats(context.Background(), "nothing", false)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestContainerStats(t *testing.T) ***REMOVED***
	expectedURL := "/containers/container_id/stats"
	cases := []struct ***REMOVED***
		stream         bool
		expectedStream string
	***REMOVED******REMOVED***
		***REMOVED***
			expectedStream: "0",
		***REMOVED***,
		***REMOVED***
			stream:         true,
			expectedStream: "1",
		***REMOVED***,
	***REMOVED***
	for _, c := range cases ***REMOVED***
		client := &Client***REMOVED***
			client: newMockClient(func(r *http.Request) (*http.Response, error) ***REMOVED***
				if !strings.HasPrefix(r.URL.Path, expectedURL) ***REMOVED***
					return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, r.URL)
				***REMOVED***

				query := r.URL.Query()
				stream := query.Get("stream")
				if stream != c.expectedStream ***REMOVED***
					return nil, fmt.Errorf("stream not set in URL query properly. Expected '%s', got %s", c.expectedStream, stream)
				***REMOVED***

				return &http.Response***REMOVED***
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("response"))),
				***REMOVED***, nil
			***REMOVED***),
		***REMOVED***
		resp, err := client.ContainerStats(context.Background(), "container_id", c.stream)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer resp.Body.Close()
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if string(content) != "response" ***REMOVED***
			t.Fatalf("expected response to contain 'response', got %s", string(content))
		***REMOVED***
	***REMOVED***
***REMOVED***
