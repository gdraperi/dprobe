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

func TestSecretListUnsupported(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.24",
		client:  &http.Client***REMOVED******REMOVED***,
	***REMOVED***
	_, err := client.SecretList(context.Background(), types.SecretListOptions***REMOVED******REMOVED***)
	assert.EqualError(t, err, `"secret list" requires API version 1.25, but the Docker daemon API version is 1.24`)
***REMOVED***

func TestSecretListError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.25",
		client:  newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, err := client.SecretList(context.Background(), types.SecretListOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestSecretList(t *testing.T) ***REMOVED***
	expectedURL := "/v1.25/secrets"

	filters := filters.NewArgs()
	filters.Add("label", "label1")
	filters.Add("label", "label2")

	listCases := []struct ***REMOVED***
		options             types.SecretListOptions
		expectedQueryParams map[string]string
	***REMOVED******REMOVED***
		***REMOVED***
			options: types.SecretListOptions***REMOVED******REMOVED***,
			expectedQueryParams: map[string]string***REMOVED***
				"filters": "",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			options: types.SecretListOptions***REMOVED***
				Filters: filters,
			***REMOVED***,
			expectedQueryParams: map[string]string***REMOVED***
				"filters": `***REMOVED***"label":***REMOVED***"label1":true,"label2":true***REMOVED******REMOVED***`,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for _, listCase := range listCases ***REMOVED***
		client := &Client***REMOVED***
			version: "1.25",
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
				content, err := json.Marshal([]swarm.Secret***REMOVED***
					***REMOVED***
						ID: "secret_id1",
					***REMOVED***,
					***REMOVED***
						ID: "secret_id2",
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

		secrets, err := client.SecretList(context.Background(), listCase.options)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if len(secrets) != 2 ***REMOVED***
			t.Fatalf("expected 2 secrets, got %v", secrets)
		***REMOVED***
	***REMOVED***
***REMOVED***
