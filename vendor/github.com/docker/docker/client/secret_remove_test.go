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

func TestSecretRemoveUnsupported(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.24",
		client:  &http.Client***REMOVED******REMOVED***,
	***REMOVED***
	err := client.SecretRemove(context.Background(), "secret_id")
	assert.EqualError(t, err, `"secret remove" requires API version 1.25, but the Docker daemon API version is 1.24`)
***REMOVED***

func TestSecretRemoveError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.25",
		client:  newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	err := client.SecretRemove(context.Background(), "secret_id")
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestSecretRemove(t *testing.T) ***REMOVED***
	expectedURL := "/v1.25/secrets/secret_id"

	client := &Client***REMOVED***
		version: "1.25",
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

	err := client.SecretRemove(context.Background(), "secret_id")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
