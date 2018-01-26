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
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestContainersPruneError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client:  newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
		version: "1.25",
	***REMOVED***

	filters := filters.NewArgs()

	_, err := client.ContainersPrune(context.Background(), filters)
	assert.EqualError(t, err, "Error response from daemon: Server error")
***REMOVED***

func TestContainersPrune(t *testing.T) ***REMOVED***
	expectedURL := "/v1.25/containers/prune"

	danglingFilters := filters.NewArgs()
	danglingFilters.Add("dangling", "true")

	noDanglingFilters := filters.NewArgs()
	noDanglingFilters.Add("dangling", "false")

	danglingUntilFilters := filters.NewArgs()
	danglingUntilFilters.Add("dangling", "true")
	danglingUntilFilters.Add("until", "2016-12-15T14:00")

	labelFilters := filters.NewArgs()
	labelFilters.Add("dangling", "true")
	labelFilters.Add("label", "label1=foo")
	labelFilters.Add("label", "label2!=bar")

	listCases := []struct ***REMOVED***
		filters             filters.Args
		expectedQueryParams map[string]string
	***REMOVED******REMOVED***
		***REMOVED***
			filters: filters.Args***REMOVED******REMOVED***,
			expectedQueryParams: map[string]string***REMOVED***
				"until":   "",
				"filter":  "",
				"filters": "",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			filters: danglingFilters,
			expectedQueryParams: map[string]string***REMOVED***
				"until":   "",
				"filter":  "",
				"filters": `***REMOVED***"dangling":***REMOVED***"true":true***REMOVED******REMOVED***`,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			filters: danglingUntilFilters,
			expectedQueryParams: map[string]string***REMOVED***
				"until":   "",
				"filter":  "",
				"filters": `***REMOVED***"dangling":***REMOVED***"true":true***REMOVED***,"until":***REMOVED***"2016-12-15T14:00":true***REMOVED******REMOVED***`,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			filters: noDanglingFilters,
			expectedQueryParams: map[string]string***REMOVED***
				"until":   "",
				"filter":  "",
				"filters": `***REMOVED***"dangling":***REMOVED***"false":true***REMOVED******REMOVED***`,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			filters: labelFilters,
			expectedQueryParams: map[string]string***REMOVED***
				"until":   "",
				"filter":  "",
				"filters": `***REMOVED***"dangling":***REMOVED***"true":true***REMOVED***,"label":***REMOVED***"label1=foo":true,"label2!=bar":true***REMOVED******REMOVED***`,
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
					assert.Equal(t, expected, actual)
				***REMOVED***
				content, err := json.Marshal(types.ContainersPruneReport***REMOVED***
					ContainersDeleted: []string***REMOVED***"container_id1", "container_id2"***REMOVED***,
					SpaceReclaimed:    9999,
				***REMOVED***)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				return &http.Response***REMOVED***
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(content)),
				***REMOVED***, nil
			***REMOVED***),
			version: "1.25",
		***REMOVED***

		report, err := client.ContainersPrune(context.Background(), listCase.filters)
		assert.NoError(t, err)
		assert.Len(t, report.ContainersDeleted, 2)
		assert.Equal(t, uint64(9999), report.SpaceReclaimed)
	***REMOVED***
***REMOVED***
