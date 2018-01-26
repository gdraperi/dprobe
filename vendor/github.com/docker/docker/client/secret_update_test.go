package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestSecretUpdateUnsupported(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.24",
		client:  &http.Client***REMOVED******REMOVED***,
	***REMOVED***
	err := client.SecretUpdate(context.Background(), "secret_id", swarm.Version***REMOVED******REMOVED***, swarm.SecretSpec***REMOVED******REMOVED***)
	assert.EqualError(t, err, `"secret update" requires API version 1.25, but the Docker daemon API version is 1.24`)
***REMOVED***

func TestSecretUpdateError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.25",
		client:  newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	err := client.SecretUpdate(context.Background(), "secret_id", swarm.Version***REMOVED******REMOVED***, swarm.SecretSpec***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestSecretUpdate(t *testing.T) ***REMOVED***
	expectedURL := "/v1.25/secrets/secret_id/update"

	client := &Client***REMOVED***
		version: "1.25",
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			if req.Method != "POST" ***REMOVED***
				return nil, fmt.Errorf("expected POST method, got %s", req.Method)
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("body"))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	err := client.SecretUpdate(context.Background(), "secret_id", swarm.Version***REMOVED******REMOVED***, swarm.SecretSpec***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
