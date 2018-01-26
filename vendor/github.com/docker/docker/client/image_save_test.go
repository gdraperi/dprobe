package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"golang.org/x/net/context"

	"strings"
)

func TestImageSaveError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ImageSave(context.Background(), []string***REMOVED***"nothing"***REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImageSave(t *testing.T) ***REMOVED***
	expectedURL := "/images/get"
	client := &Client***REMOVED***
		client: newMockClient(func(r *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(r.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, r.URL)
			***REMOVED***
			query := r.URL.Query()
			names := query["names"]
			expectedNames := []string***REMOVED***"image_id1", "image_id2"***REMOVED***
			if !reflect.DeepEqual(names, expectedNames) ***REMOVED***
				return nil, fmt.Errorf("names not set in URL query properly. Expected %v, got %v", names, expectedNames)
			***REMOVED***

			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("response"))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***
	saveResponse, err := client.ImageSave(context.Background(), []string***REMOVED***"image_id1", "image_id2"***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	response, err := ioutil.ReadAll(saveResponse)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	saveResponse.Close()
	if string(response) != "response" ***REMOVED***
		t.Fatalf("expected response to contain 'response', got %s", string(response))
	***REMOVED***
***REMOVED***
