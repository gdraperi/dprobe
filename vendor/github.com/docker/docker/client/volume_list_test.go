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
	volumetypes "github.com/docker/docker/api/types/volume"
	"golang.org/x/net/context"
)

func TestVolumeListError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, err := client.VolumeList(context.Background(), filters.NewArgs())
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestVolumeList(t *testing.T) ***REMOVED***
	expectedURL := "/volumes"

	noDanglingFilters := filters.NewArgs()
	noDanglingFilters.Add("dangling", "false")

	danglingFilters := filters.NewArgs()
	danglingFilters.Add("dangling", "true")

	labelFilters := filters.NewArgs()
	labelFilters.Add("label", "label1")
	labelFilters.Add("label", "label2")

	listCases := []struct ***REMOVED***
		filters         filters.Args
		expectedFilters string
	***REMOVED******REMOVED***
		***REMOVED***
			filters:         filters.NewArgs(),
			expectedFilters: "",
		***REMOVED***, ***REMOVED***
			filters:         noDanglingFilters,
			expectedFilters: `***REMOVED***"dangling":***REMOVED***"false":true***REMOVED******REMOVED***`,
		***REMOVED***, ***REMOVED***
			filters:         danglingFilters,
			expectedFilters: `***REMOVED***"dangling":***REMOVED***"true":true***REMOVED******REMOVED***`,
		***REMOVED***, ***REMOVED***
			filters:         labelFilters,
			expectedFilters: `***REMOVED***"label":***REMOVED***"label1":true,"label2":true***REMOVED******REMOVED***`,
		***REMOVED***,
	***REMOVED***

	for _, listCase := range listCases ***REMOVED***
		client := &Client***REMOVED***
			client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
				if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
					return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
				***REMOVED***
				query := req.URL.Query()
				actualFilters := query.Get("filters")
				if actualFilters != listCase.expectedFilters ***REMOVED***
					return nil, fmt.Errorf("filters not set in URL query properly. Expected '%s', got %s", listCase.expectedFilters, actualFilters)
				***REMOVED***
				content, err := json.Marshal(volumetypes.VolumesListOKBody***REMOVED***
					Volumes: []*types.Volume***REMOVED***
						***REMOVED***
							Name:   "volume",
							Driver: "local",
						***REMOVED***,
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

		volumeResponse, err := client.VolumeList(context.Background(), listCase.filters)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if len(volumeResponse.Volumes) != 1 ***REMOVED***
			t.Fatalf("expected 1 volume, got %v", volumeResponse.Volumes)
		***REMOVED***
	***REMOVED***
***REMOVED***
