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
	"github.com/docker/docker/api/types/swarm"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestSecretCreateUnsupported(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.24",
		client:  &http.Client***REMOVED******REMOVED***,
	***REMOVED***
	_, err := client.SecretCreate(context.Background(), swarm.SecretSpec***REMOVED******REMOVED***)
	assert.EqualError(t, err, `"secret create" requires API version 1.25, but the Docker daemon API version is 1.24`)
***REMOVED***

func TestSecretCreateError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.25",
		client:  newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.SecretCreate(context.Background(), swarm.SecretSpec***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestSecretCreate(t *testing.T) ***REMOVED***
	expectedURL := "/v1.25/secrets/create"
	client := &Client***REMOVED***
		version: "1.25",
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			if req.Method != "POST" ***REMOVED***
				return nil, fmt.Errorf("expected POST method, got %s", req.Method)
			***REMOVED***
			b, err := json.Marshal(types.SecretCreateResponse***REMOVED***
				ID: "test_secret",
			***REMOVED***)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusCreated,
				Body:       ioutil.NopCloser(bytes.NewReader(b)),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	r, err := client.SecretCreate(context.Background(), swarm.SecretSpec***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if r.ID != "test_secret" ***REMOVED***
		t.Fatalf("expected `test_secret`, got %s", r.ID)
	***REMOVED***
***REMOVED***
