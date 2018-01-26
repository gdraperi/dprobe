package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
)

func TestImageCreateError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ImageCreate(context.Background(), "reference", types.ImageCreateOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImageCreate(t *testing.T) ***REMOVED***
	expectedURL := "/images/create"
	expectedImage := "test:5000/my_image"
	expectedTag := "sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	expectedReference := fmt.Sprintf("%s@%s", expectedImage, expectedTag)
	expectedRegistryAuth := "eyJodHRwczovL2luZGV4LmRvY2tlci5pby92MS8iOnsiYXV0aCI6ImRHOTBid289IiwiZW1haWwiOiJqb2huQGRvZS5jb20ifX0="
	client := &Client***REMOVED***
		client: newMockClient(func(r *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(r.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, r.URL)
			***REMOVED***
			registryAuth := r.Header.Get("X-Registry-Auth")
			if registryAuth != expectedRegistryAuth ***REMOVED***
				return nil, fmt.Errorf("X-Registry-Auth header not properly set in the request. Expected '%s', got %s", expectedRegistryAuth, registryAuth)
			***REMOVED***

			query := r.URL.Query()
			fromImage := query.Get("fromImage")
			if fromImage != expectedImage ***REMOVED***
				return nil, fmt.Errorf("fromImage not set in URL query properly. Expected '%s', got %s", expectedImage, fromImage)
			***REMOVED***

			tag := query.Get("tag")
			if tag != expectedTag ***REMOVED***
				return nil, fmt.Errorf("tag not set in URL query properly. Expected '%s', got %s", expectedTag, tag)
			***REMOVED***

			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("body"))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	createResponse, err := client.ImageCreate(context.Background(), expectedReference, types.ImageCreateOptions***REMOVED***
		RegistryAuth: expectedRegistryAuth,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	response, err := ioutil.ReadAll(createResponse)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err = createResponse.Close(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if string(response) != "body" ***REMOVED***
		t.Fatalf("expected Body to contain 'body' string, got %s", response)
	***REMOVED***
***REMOVED***
