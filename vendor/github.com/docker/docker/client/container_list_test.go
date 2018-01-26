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

func TestContainerListError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ContainerList(context.Background(), types.ContainerListOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestContainerList(t *testing.T) ***REMOVED***
	expectedURL := "/containers/json"
	expectedFilters := `***REMOVED***"before":***REMOVED***"container":true***REMOVED***,"label":***REMOVED***"label1":true,"label2":true***REMOVED******REMOVED***`
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			query := req.URL.Query()
			all := query.Get("all")
			if all != "1" ***REMOVED***
				return nil, fmt.Errorf("all not set in URL query properly. Expected '1', got %s", all)
			***REMOVED***
			limit := query.Get("limit")
			if limit != "0" ***REMOVED***
				return nil, fmt.Errorf("limit should have not be present in query. Expected '0', got %s", limit)
			***REMOVED***
			since := query.Get("since")
			if since != "container" ***REMOVED***
				return nil, fmt.Errorf("since not set in URL query properly. Expected 'container', got %s", since)
			***REMOVED***
			before := query.Get("before")
			if before != "" ***REMOVED***
				return nil, fmt.Errorf("before should have not be present in query, go %s", before)
			***REMOVED***
			size := query.Get("size")
			if size != "1" ***REMOVED***
				return nil, fmt.Errorf("size not set in URL query properly. Expected '1', got %s", size)
			***REMOVED***
			filters := query.Get("filters")
			if filters != expectedFilters ***REMOVED***
				return nil, fmt.Errorf("expected filters incoherent '%v' with actual filters %v", expectedFilters, filters)
			***REMOVED***

			b, err := json.Marshal([]types.Container***REMOVED***
				***REMOVED***
					ID: "container_id1",
				***REMOVED***,
				***REMOVED***
					ID: "container_id2",
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

	filters := filters.NewArgs()
	filters.Add("label", "label1")
	filters.Add("label", "label2")
	filters.Add("before", "container")
	containers, err := client.ContainerList(context.Background(), types.ContainerListOptions***REMOVED***
		Size:    true,
		All:     true,
		Since:   "container",
		Filters: filters,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if len(containers) != 2 ***REMOVED***
		t.Fatalf("expected 2 containers, got %v", containers)
	***REMOVED***
***REMOVED***
