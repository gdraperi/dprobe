package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"golang.org/x/net/context"
)

func TestImageListError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, err := client.ImageList(context.Background(), types.ImageListOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImageList(t *testing.T) ***REMOVED***
	expectedURL := "/images/json"

	noDanglingfilters := filters.NewArgs()
	noDanglingfilters.Add("dangling", "false")

	filters := filters.NewArgs()
	filters.Add("label", "label1")
	filters.Add("label", "label2")
	filters.Add("dangling", "true")

	listCases := []struct ***REMOVED***
		options             types.ImageListOptions
		expectedQueryParams map[string]string
	***REMOVED******REMOVED***
		***REMOVED***
			options: types.ImageListOptions***REMOVED******REMOVED***,
			expectedQueryParams: map[string]string***REMOVED***
				"all":     "",
				"filter":  "",
				"filters": "",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			options: types.ImageListOptions***REMOVED***
				Filters: filters,
			***REMOVED***,
			expectedQueryParams: map[string]string***REMOVED***
				"all":     "",
				"filter":  "",
				"filters": `***REMOVED***"dangling":***REMOVED***"true":true***REMOVED***,"label":***REMOVED***"label1":true,"label2":true***REMOVED******REMOVED***`,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			options: types.ImageListOptions***REMOVED***
				Filters: noDanglingfilters,
			***REMOVED***,
			expectedQueryParams: map[string]string***REMOVED***
				"all":     "",
				"filter":  "",
				"filters": `***REMOVED***"dangling":***REMOVED***"false":true***REMOVED******REMOVED***`,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for _, listCase := range listCases ***REMOVED***
		client := &Client***REMOVED***
			client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
				if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
					return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
				***REMOVED***
				query := req.URL.Query()
				for key, expected := range listCase.expectedQueryParams ***REMOVED***
					actual := query.Get(key)
					if actual != expected ***REMOVED***
						return nil, fmt.Errorf("%s not set in URL query properly. Expected '%s', got %s", key, expected, actual)
					***REMOVED***
				***REMOVED***
				content, err := json.Marshal([]types.ImageSummary***REMOVED***
					***REMOVED***
						ID: "image_id2",
					***REMOVED***,
					***REMOVED***
						ID: "image_id2",
					***REMOVED***,
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

		images, err := client.ImageList(context.Background(), listCase.options)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if len(images) != 2 ***REMOVED***
			t.Fatalf("expected 2 images, got %v", images)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestImageListApiBefore125(t *testing.T) ***REMOVED***
	expectedFilter := "image:tag"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			query := req.URL.Query()
			actualFilter := query.Get("filter")
			if actualFilter != expectedFilter ***REMOVED***
				return nil, fmt.Errorf("filter not set in URL query properly. Expected '%s', got %s", expectedFilter, actualFilter)
			***REMOVED***
			actualFilters := query.Get("filters")
			if actualFilters != "" ***REMOVED***
				return nil, fmt.Errorf("filters should have not been present, were with value: %s", actualFilters)
			***REMOVED***
			content, err := json.Marshal([]types.ImageSummary***REMOVED***
				***REMOVED***
					ID: "image_id2",
				***REMOVED***,
				***REMOVED***
					ID: "image_id2",
				***REMOVED***,
			***REMOVED***)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(content)),
			***REMOVED***, nil
		***REMOVED***),
		version: "1.24",
	***REMOVED***

	filters := filters.NewArgs()
	filters.Add("reference", "image:tag")

	options := types.ImageListOptions***REMOVED***
		Filters: filters,
	***REMOVED***

	images, err := client.ImageList(context.Background(), options)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if len(images) != 2 ***REMOVED***
		t.Fatalf("expected 2 images, got %v", images)
	***REMOVED***
***REMOVED***
