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
	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

func TestNodeListError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, err := client.NodeList(context.Background(), types.NodeListOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestNodeList(t *testing.T) ***REMOVED***
	expectedURL := "/nodes"

	filters := filters.NewArgs()
	filters.Add("label", "label1")
	filters.Add("label", "label2")

	listCases := []struct ***REMOVED***
		options             types.NodeListOptions
		expectedQueryParams map[string]string
	***REMOVED******REMOVED***
		***REMOVED***
			options: types.NodeListOptions***REMOVED******REMOVED***,
			expectedQueryParams: map[string]string***REMOVED***
				"filters": "",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			options: types.NodeListOptions***REMOVED***
				Filters: filters,
			***REMOVED***,
			expectedQueryParams: map[string]string***REMOVED***
				"filters": `***REMOVED***"label":***REMOVED***"label1":true,"label2":true***REMOVED******REMOVED***`,
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
				content, err := json.Marshal([]swarm.Node***REMOVED***
					***REMOVED***
						ID: "node_id1",
					***REMOVED***,
					***REMOVED***
						ID: "node_id2",
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

		nodes, err := client.NodeList(context.Background(), listCase.options)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if len(nodes) != 2 ***REMOVED***
			t.Fatalf("expected 2 nodes, got %v", nodes)
		***REMOVED***
	***REMOVED***
***REMOVED***
