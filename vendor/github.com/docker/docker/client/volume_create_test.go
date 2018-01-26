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
	volumetypes "github.com/docker/docker/api/types/volume"
	"golang.org/x/net/context"
)

func TestVolumeCreateError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, err := client.VolumeCreate(context.Background(), volumetypes.VolumesCreateBody***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestVolumeCreate(t *testing.T) ***REMOVED***
	expectedURL := "/volumes/create"

	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***

			if req.Method != "POST" ***REMOVED***
				return nil, fmt.Errorf("expected POST method, got %s", req.Method)
			***REMOVED***

			content, err := json.Marshal(types.Volume***REMOVED***
				Name:       "volume",
				Driver:     "local",
				Mountpoint: "mountpoint",
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

	volume, err := client.VolumeCreate(context.Background(), volumetypes.VolumesCreateBody***REMOVED***
		Name:   "myvolume",
		Driver: "mydriver",
		DriverOpts: map[string]string***REMOVED***
			"opt-key": "opt-value",
		***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if volume.Name != "volume" ***REMOVED***
		t.Fatalf("expected volume.Name to be 'volume', got %s", volume.Name)
	***REMOVED***
	if volume.Driver != "local" ***REMOVED***
		t.Fatalf("expected volume.Driver to be 'local', got %s", volume.Driver)
	***REMOVED***
	if volume.Mountpoint != "mountpoint" ***REMOVED***
		t.Fatalf("expected volume.Mountpoint to be 'mountpoint', got %s", volume.Mountpoint)
	***REMOVED***
***REMOVED***
