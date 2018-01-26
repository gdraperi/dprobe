package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/image"
	"golang.org/x/net/context"
)

func TestImageHistoryError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ImageHistory(context.Background(), "nothing")
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImageHistory(t *testing.T) ***REMOVED***
	expectedURL := "/images/image_id/history"
	client := &Client***REMOVED***
		client: newMockClient(func(r *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(r.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, r.URL)
			***REMOVED***
			b, err := json.Marshal([]image.HistoryResponseItem***REMOVED***
				***REMOVED***
					ID:   "image_id1",
					Tags: []string***REMOVED***"tag1", "tag2"***REMOVED***,
				***REMOVED***,
				***REMOVED***
					ID:   "image_id2",
					Tags: []string***REMOVED***"tag1", "tag2"***REMOVED***,
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
	imageHistories, err := client.ImageHistory(context.Background(), "image_id")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if len(imageHistories) != 2 ***REMOVED***
		t.Fatalf("expected 2 containers, got %v", imageHistories)
	***REMOVED***
***REMOVED***
