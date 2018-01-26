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

func TestImageLoadError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, err := client.ImageLoad(context.Background(), nil, true)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImageLoad(t *testing.T) ***REMOVED***
	expectedURL := "/images/load"
	expectedInput := "inputBody"
	expectedOutput := "outputBody"
	loadCases := []struct ***REMOVED***
		quiet                bool
		responseContentType  string
		expectedResponseJSON bool
		expectedQueryParams  map[string]string
	***REMOVED******REMOVED***
		***REMOVED***
			quiet:                false,
			responseContentType:  "text/plain",
			expectedResponseJSON: false,
			expectedQueryParams: map[string]string***REMOVED***
				"quiet": "0",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			quiet:                true,
			responseContentType:  "application/json",
			expectedResponseJSON: true,
			expectedQueryParams: map[string]string***REMOVED***
				"quiet": "1",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for _, loadCase := range loadCases ***REMOVED***
		client := &Client***REMOVED***
			client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
				if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
					return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
				***REMOVED***
				contentType := req.Header.Get("Content-Type")
				if contentType != "application/x-tar" ***REMOVED***
					return nil, fmt.Errorf("content-type not set in URL headers properly. Expected 'application/x-tar', got %s", contentType)
				***REMOVED***
				query := req.URL.Query()
				for key, expected := range loadCase.expectedQueryParams ***REMOVED***
					actual := query.Get(key)
					if actual != expected ***REMOVED***
						return nil, fmt.Errorf("%s not set in URL query properly. Expected '%s', got %s", key, expected, actual)
					***REMOVED***
				***REMOVED***
				headers := http.Header***REMOVED******REMOVED***
				headers.Add("Content-Type", loadCase.responseContentType)
				return &http.Response***REMOVED***
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(expectedOutput))),
					Header:     headers,
				***REMOVED***, nil
			***REMOVED***),
		***REMOVED***

		input := bytes.NewReader([]byte(expectedInput))
		imageLoadResponse, err := client.ImageLoad(context.Background(), input, loadCase.quiet)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if imageLoadResponse.JSON != loadCase.expectedResponseJSON ***REMOVED***
			t.Fatalf("expected a JSON response, was not.")
		***REMOVED***
		body, err := ioutil.ReadAll(imageLoadResponse.Body)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if string(body) != expectedOutput ***REMOVED***
			t.Fatalf("expected %s, got %s", expectedOutput, string(body))
		***REMOVED***
	***REMOVED***
***REMOVED***
