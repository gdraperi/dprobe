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
	"github.com/docker/docker/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestVolumeInspectError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, err := client.VolumeInspect(context.Background(), "nothing")
	testutil.ErrorContains(t, err, "Error response from daemon: Server error")
***REMOVED***

func TestVolumeInspectNotFound(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusNotFound, "Server error")),
	***REMOVED***

	_, err := client.VolumeInspect(context.Background(), "unknown")
	assert.True(t, IsErrNotFound(err))
***REMOVED***

func TestVolumeInspectWithEmptyID(t *testing.T) ***REMOVED***
	expectedURL := "/volumes/"

	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			assert.Equal(t, req.URL.Path, expectedURL)
			return &http.Response***REMOVED***
				StatusCode: http.StatusNotFound,
				Body:       ioutil.NopCloser(bytes.NewReader(nil)),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***
	_, err := client.VolumeInspect(context.Background(), "")
	testutil.ErrorContains(t, err, "No such volume: ")

***REMOVED***

func TestVolumeInspect(t *testing.T) ***REMOVED***
	expectedURL := "/volumes/volume_id"
	expected := types.Volume***REMOVED***
		Name:       "name",
		Driver:     "driver",
		Mountpoint: "mountpoint",
	***REMOVED***

	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			if req.Method != "GET" ***REMOVED***
				return nil, fmt.Errorf("expected GET method, got %s", req.Method)
			***REMOVED***
			content, err := json.Marshal(expected)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(content)),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	volume, err := client.VolumeInspect(context.Background(), "volume_id")
	require.NoError(t, err)
	assert.Equal(t, expected, volume)
***REMOVED***
