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

func TestImageTagError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	err := client.ImageTag(context.Background(), "image_id", "repo:tag")
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

// Note: this is not testing all the InvalidReference as it's the responsibility
// of distribution/reference package.
func TestImageTagInvalidReference(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	err := client.ImageTag(context.Background(), "image_id", "aa/asdf$$^/aa")
	if err == nil || err.Error() != `Error parsing reference: "aa/asdf$$^/aa" is not a valid repository/tag: invalid reference format` ***REMOVED***
		t.Fatalf("expected ErrReferenceInvalidFormat, got %v", err)
	***REMOVED***
***REMOVED***

func TestImageTagInvalidSourceImageName(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	err := client.ImageTag(context.Background(), "invalid_source_image_name_", "repo:tag")
	if err == nil || err.Error() != "Error parsing reference: \"invalid_source_image_name_\" is not a valid repository/tag: invalid reference format" ***REMOVED***
		t.Fatalf("expected Parsing Reference Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImageTagHexSource(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusOK, "OK")),
	***REMOVED***

	err := client.ImageTag(context.Background(), "0d409d33b27e47423b049f7f863faa08655a8c901749c2b25b93ca67d01a470d", "repo:tag")
	if err != nil ***REMOVED***
		t.Fatalf("got error: %v", err)
	***REMOVED***
***REMOVED***

func TestImageTag(t *testing.T) ***REMOVED***
	expectedURL := "/images/image_id/tag"
	tagCases := []struct ***REMOVED***
		reference           string
		expectedQueryParams map[string]string
	***REMOVED******REMOVED***
		***REMOVED***
			reference: "repository:tag1",
			expectedQueryParams: map[string]string***REMOVED***
				"repo": "repository",
				"tag":  "tag1",
			***REMOVED***,
		***REMOVED***, ***REMOVED***
			reference: "another_repository:latest",
			expectedQueryParams: map[string]string***REMOVED***
				"repo": "another_repository",
				"tag":  "latest",
			***REMOVED***,
		***REMOVED***, ***REMOVED***
			reference: "another_repository",
			expectedQueryParams: map[string]string***REMOVED***
				"repo": "another_repository",
				"tag":  "latest",
			***REMOVED***,
		***REMOVED***, ***REMOVED***
			reference: "test/another_repository",
			expectedQueryParams: map[string]string***REMOVED***
				"repo": "test/another_repository",
				"tag":  "latest",
			***REMOVED***,
		***REMOVED***, ***REMOVED***
			reference: "test/another_repository:tag1",
			expectedQueryParams: map[string]string***REMOVED***
				"repo": "test/another_repository",
				"tag":  "tag1",
			***REMOVED***,
		***REMOVED***, ***REMOVED***
			reference: "test/test/another_repository:tag1",
			expectedQueryParams: map[string]string***REMOVED***
				"repo": "test/test/another_repository",
				"tag":  "tag1",
			***REMOVED***,
		***REMOVED***, ***REMOVED***
			reference: "test:5000/test/another_repository:tag1",
			expectedQueryParams: map[string]string***REMOVED***
				"repo": "test:5000/test/another_repository",
				"tag":  "tag1",
			***REMOVED***,
		***REMOVED***, ***REMOVED***
			reference: "test:5000/test/another_repository",
			expectedQueryParams: map[string]string***REMOVED***
				"repo": "test:5000/test/another_repository",
				"tag":  "latest",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for _, tagCase := range tagCases ***REMOVED***
		client := &Client***REMOVED***
			client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
				if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
					return nil, fmt.Errorf("expected URL '%s', got '%s'", expectedURL, req.URL)
				***REMOVED***
				if req.Method != "POST" ***REMOVED***
					return nil, fmt.Errorf("expected POST method, got %s", req.Method)
				***REMOVED***
				query := req.URL.Query()
				for key, expected := range tagCase.expectedQueryParams ***REMOVED***
					actual := query.Get(key)
					if actual != expected ***REMOVED***
						return nil, fmt.Errorf("%s not set in URL query properly. Expected '%s', got %s", key, expected, actual)
					***REMOVED***
				***REMOVED***
				return &http.Response***REMOVED***
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
				***REMOVED***, nil
			***REMOVED***),
		***REMOVED***
		err := client.ImageTag(context.Background(), "image_id", tagCase.reference)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***
