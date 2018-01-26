package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestServiceRemoveError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	err := client.ServiceRemove(context.Background(), "service_id")
	assert.EqualError(t, err, "Error response from daemon: Server error")
***REMOVED***

func TestServiceRemoveNotFoundError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusNotFound, "missing")),
	***REMOVED***

	err := client.ServiceRemove(context.Background(), "service_id")
	assert.EqualError(t, err, "Error: No such service: service_id")
	assert.True(t, IsErrNotFound(err))
***REMOVED***

func TestServiceRemove(t *testing.T) ***REMOVED***
	expectedURL := "/services/service_id"

	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			if req.Method != "DELETE" ***REMOVED***
				return nil, fmt.Errorf("expected DELETE method, got %s", req.Method)
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("body"))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	err := client.ServiceRemove(context.Background(), "service_id")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
