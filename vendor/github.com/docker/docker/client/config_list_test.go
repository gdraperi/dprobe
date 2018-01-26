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
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestConfigListUnsupported(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.29",
		client:  &http.Client***REMOVED******REMOVED***,
	***REMOVED***
	_, err := client.ConfigList(context.Background(), types.ConfigListOptions***REMOVED******REMOVED***)
	assert.EqualError(t, err, `"config list" requires API version 1.30, but the Docker daemon API version is 1.29`)
***REMOVED***

func TestConfigListError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.30",
		client:  newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, err := client.ConfigList(context.Background(), types.ConfigListOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestConfigList(t *testing.T) ***REMOVED***
	expectedURL := "/v1.30/configs"

	filters := filters.NewArgs()
	filters.Add("label", "label1")
	filters.Add("label", "label2")

	listCases := []struct ***REMOVED***
		options             types.ConfigListOptions
		expectedQueryParams map[string]string
	***REMOVED******REMOVED***
		***REMOVED***
			options: types.ConfigListOptions***REMOVED******REMOVED***,
			expectedQueryParams: map[string]string***REMOVED***
				"filters": "",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			options: types.ConfigListOptions***REMOVED***
				Filters: filters,
			***REMOVED***,
			expectedQueryParams: map[string]string***REMOVED***
				"filters": `***REMOVED***"label":***REMOVED***"label1":true,"label2":true***REMOVED******REMOVED***`,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for _, listCase := range listCases ***REMOVED***
		client := &Client***REMOVED***
			version: "1.30",
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
				content, err := json.Marshal([]swarm.Config***REMOVED***
					***REMOVED***
						ID: "config_id1",
					***REMOVED***,
					***REMOVED***
						ID: "config_id2",
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

		configs, err := client.ConfigList(context.Background(), listCase.options)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if len(configs) != 2 ***REMOVED***
			t.Fatalf("expected 2 configs, got %v", configs)
		***REMOVED***
	***REMOVED***
***REMOVED***
