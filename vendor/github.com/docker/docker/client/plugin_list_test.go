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

func TestPluginListError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, err := client.PluginList(context.Background(), filters.NewArgs())
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestPluginList(t *testing.T) ***REMOVED***
	expectedURL := "/plugins"

	enabledFilters := filters.NewArgs()
	enabledFilters.Add("enabled", "true")

	capabilityFilters := filters.NewArgs()
	capabilityFilters.Add("capability", "volumedriver")
	capabilityFilters.Add("capability", "authz")

	listCases := []struct ***REMOVED***
		filters             filters.Args
		expectedQueryParams map[string]string
	***REMOVED******REMOVED***
		***REMOVED***
			filters: filters.NewArgs(),
			expectedQueryParams: map[string]string***REMOVED***
				"all":     "",
				"filter":  "",
				"filters": "",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			filters: enabledFilters,
			expectedQueryParams: map[string]string***REMOVED***
				"all":     "",
				"filter":  "",
				"filters": `***REMOVED***"enabled":***REMOVED***"true":true***REMOVED******REMOVED***`,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			filters: capabilityFilters,
			expectedQueryParams: map[string]string***REMOVED***
				"all":     "",
				"filter":  "",
				"filters": `***REMOVED***"capability":***REMOVED***"authz":true,"volumedriver":true***REMOVED******REMOVED***`,
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
				content, err := json.Marshal([]*types.Plugin***REMOVED***
					***REMOVED***
						ID: "plugin_id1",
					***REMOVED***,
					***REMOVED***
						ID: "plugin_id2",
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

		plugins, err := client.PluginList(context.Background(), listCase.filters)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if len(plugins) != 2 ***REMOVED***
			t.Fatalf("expected 2 plugins, got %v", plugins)
		***REMOVED***
	***REMOVED***
***REMOVED***
