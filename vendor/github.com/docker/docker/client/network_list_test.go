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

func TestNetworkListError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, err := client.NetworkList(context.Background(), types.NetworkListOptions***REMOVED***
		Filters: filters.NewArgs(),
	***REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestNetworkList(t *testing.T) ***REMOVED***
	expectedURL := "/networks"

	noDanglingFilters := filters.NewArgs()
	noDanglingFilters.Add("dangling", "false")

	danglingFilters := filters.NewArgs()
	danglingFilters.Add("dangling", "true")

	labelFilters := filters.NewArgs()
	labelFilters.Add("label", "label1")
	labelFilters.Add("label", "label2")

	listCases := []struct ***REMOVED***
		options         types.NetworkListOptions
		expectedFilters string
	***REMOVED******REMOVED***
		***REMOVED***
			options: types.NetworkListOptions***REMOVED***
				Filters: filters.NewArgs(),
			***REMOVED***,
			expectedFilters: "",
		***REMOVED***, ***REMOVED***
			options: types.NetworkListOptions***REMOVED***
				Filters: noDanglingFilters,
			***REMOVED***,
			expectedFilters: `***REMOVED***"dangling":***REMOVED***"false":true***REMOVED******REMOVED***`,
		***REMOVED***, ***REMOVED***
			options: types.NetworkListOptions***REMOVED***
				Filters: danglingFilters,
			***REMOVED***,
			expectedFilters: `***REMOVED***"dangling":***REMOVED***"true":true***REMOVED******REMOVED***`,
		***REMOVED***, ***REMOVED***
			options: types.NetworkListOptions***REMOVED***
				Filters: labelFilters,
			***REMOVED***,
			expectedFilters: `***REMOVED***"label":***REMOVED***"label1":true,"label2":true***REMOVED******REMOVED***`,
		***REMOVED***,
	***REMOVED***

	for _, listCase := range listCases ***REMOVED***
		client := &Client***REMOVED***
			client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
				if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
					return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
				***REMOVED***
				if req.Method != "GET" ***REMOVED***
					return nil, fmt.Errorf("expected GET method, got %s", req.Method)
				***REMOVED***
				query := req.URL.Query()
				actualFilters := query.Get("filters")
				if actualFilters != listCase.expectedFilters ***REMOVED***
					return nil, fmt.Errorf("filters not set in URL query properly. Expected '%s', got %s", listCase.expectedFilters, actualFilters)
				***REMOVED***
				content, err := json.Marshal([]types.NetworkResource***REMOVED***
					***REMOVED***
						Name:   "network",
						Driver: "bridge",
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

		networkResources, err := client.NetworkList(context.Background(), listCase.options)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if len(networkResources) != 1 ***REMOVED***
			t.Fatalf("expected 1 network resource, got %v", networkResources)
		***REMOVED***
	***REMOVED***
***REMOVED***
