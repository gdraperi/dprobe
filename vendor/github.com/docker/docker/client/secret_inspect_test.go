package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestSecretInspectUnsupported(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.24",
		client:  &http.Client***REMOVED******REMOVED***,
	***REMOVED***
	_, _, err := client.SecretInspectWithRaw(context.Background(), "nothing")
	assert.EqualError(t, err, `"secret inspect" requires API version 1.25, but the Docker daemon API version is 1.24`)
***REMOVED***

func TestSecretInspectError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.25",
		client:  newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, _, err := client.SecretInspectWithRaw(context.Background(), "nothing")
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestSecretInspectSecretNotFound(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.25",
		client:  newMockClient(errorMock(http.StatusNotFound, "Server error")),
	***REMOVED***

	_, _, err := client.SecretInspectWithRaw(context.Background(), "unknown")
	if err == nil || !IsErrNotFound(err) ***REMOVED***
		t.Fatalf("expected a secretNotFoundError error, got %v", err)
	***REMOVED***
***REMOVED***

func TestSecretInspect(t *testing.T) ***REMOVED***
	expectedURL := "/v1.25/secrets/secret_id"
	client := &Client***REMOVED***
		version: "1.25",
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			content, err := json.Marshal(swarm.Secret***REMOVED***
				ID: "secret_id",
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

	secretInspect, _, err := client.SecretInspectWithRaw(context.Background(), "secret_id")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if secretInspect.ID != "secret_id" ***REMOVED***
		t.Fatalf("expected `secret_id`, got %s", secretInspect.ID)
	***REMOVED***
***REMOVED***
