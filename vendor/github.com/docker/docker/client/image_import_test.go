package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

func TestImageImportError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ImageImport(context.Background(), types.ImageImportSource***REMOVED******REMOVED***, "image:tag", types.ImageImportOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImageImport(t *testing.T) ***REMOVED***
	expectedURL := "/images/create"
	client := &Client***REMOVED***
		client: newMockClient(func(r *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(r.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, r.URL)
			***REMOVED***
			query := r.URL.Query()
			fromSrc := query.Get("fromSrc")
			if fromSrc != "image_source" ***REMOVED***
				return nil, fmt.Errorf("fromSrc not set in URL query properly. Expected 'image_source', got %s", fromSrc)
			***REMOVED***
			repo := query.Get("repo")
			if repo != "repository_name:imported" ***REMOVED***
				return nil, fmt.Errorf("repo not set in URL query properly. Expected 'repository_name:imported', got %s", repo)
			***REMOVED***
			tag := query.Get("tag")
			if tag != "imported" ***REMOVED***
				return nil, fmt.Errorf("tag not set in URL query properly. Expected 'imported', got %s", tag)
			***REMOVED***
			message := query.Get("message")
			if message != "A message" ***REMOVED***
				return nil, fmt.Errorf("message not set in URL query properly. Expected 'A message', got %s", message)
			***REMOVED***
			changes := query["changes"]
			expectedChanges := []string***REMOVED***"change1", "change2"***REMOVED***
			if !reflect.DeepEqual(expectedChanges, changes) ***REMOVED***
				return nil, fmt.Errorf("changes not set in URL query properly. Expected %v, got %v", expectedChanges, changes)
			***REMOVED***

			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("response"))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***
	importResponse, err := client.ImageImport(context.Background(), types.ImageImportSource***REMOVED***
		Source:     strings.NewReader("source"),
		SourceName: "image_source",
	***REMOVED***, "repository_name:imported", types.ImageImportOptions***REMOVED***
		Tag:     "imported",
		Message: "A message",
		Changes: []string***REMOVED***"change1", "change2"***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	response, err := ioutil.ReadAll(importResponse)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	importResponse.Close()
	if string(response) != "response" ***REMOVED***
		t.Fatalf("expected response to contain 'response', got %s", string(response))
	***REMOVED***
***REMOVED***
