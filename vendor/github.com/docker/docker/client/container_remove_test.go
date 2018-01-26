package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestContainerRemoveError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	err := client.ContainerRemove(context.Background(), "container_id", types.ContainerRemoveOptions***REMOVED******REMOVED***)
	assert.EqualError(t, err, "Error response from daemon: Server error")
***REMOVED***

func TestContainerRemoveNotFoundError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusNotFound, "missing")),
	***REMOVED***
	err := client.ContainerRemove(context.Background(), "container_id", types.ContainerRemoveOptions***REMOVED******REMOVED***)
	assert.EqualError(t, err, "Error: No such container: container_id")
	assert.True(t, IsErrNotFound(err))
***REMOVED***

func TestContainerRemove(t *testing.T) ***REMOVED***
	expectedURL := "/containers/container_id"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			query := req.URL.Query()
			volume := query.Get("v")
			if volume != "1" ***REMOVED***
				return nil, fmt.Errorf("v (volume) not set in URL query properly. Expected '1', got %s", volume)
			***REMOVED***
			force := query.Get("force")
			if force != "1" ***REMOVED***
				return nil, fmt.Errorf("force not set in URL query properly. Expected '1', got %s", force)
			***REMOVED***
			link := query.Get("link")
			if link != "" ***REMOVED***
				return nil, fmt.Errorf("link should have not be present in query, go %s", link)
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	err := client.ContainerRemove(context.Background(), "container_id", types.ContainerRemoveOptions***REMOVED***
		RemoveVolumes: true,
		Force:         true,
	***REMOVED***)
	assert.NoError(t, err)
***REMOVED***
