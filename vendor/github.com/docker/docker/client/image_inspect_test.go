package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

func TestImageInspectError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, _, err := client.ImageInspectWithRaw(context.Background(), "nothing")
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImageInspectImageNotFound(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusNotFound, "Server error")),
	***REMOVED***

	_, _, err := client.ImageInspectWithRaw(context.Background(), "unknown")
	if err == nil || !IsErrNotFound(err) ***REMOVED***
		t.Fatalf("expected an imageNotFound error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImageInspect(t *testing.T) ***REMOVED***
	expectedURL := "/images/image_id/json"
	expectedTags := []string***REMOVED***"tag1", "tag2"***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			content, err := json.Marshal(types.ImageInspect***REMOVED***
				ID:       "image_id",
				RepoTags: expectedTags,
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

	imageInspect, _, err := client.ImageInspectWithRaw(context.Background(), "image_id")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if imageInspect.ID != "image_id" ***REMOVED***
		t.Fatalf("expected `image_id`, got %s", imageInspect.ID)
	***REMOVED***
	if !reflect.DeepEqual(imageInspect.RepoTags, expectedTags) ***REMOVED***
		t.Fatalf("expected `%v`, got %v", expectedTags, imageInspect.RepoTags)
	***REMOVED***
***REMOVED***
