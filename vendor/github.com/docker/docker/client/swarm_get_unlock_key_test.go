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

func TestSwarmGetUnlockKeyError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, err := client.SwarmGetUnlockKey(context.Background())
	testutil.ErrorContains(t, err, "Error response from daemon: Server error")
***REMOVED***

func TestSwarmGetUnlockKey(t *testing.T) ***REMOVED***
	expectedURL := "/swarm/unlockkey"
	unlockKey := "SWMKEY-1-y6guTZNTwpQeTL5RhUfOsdBdXoQjiB2GADHSRJvbXeE"

	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			if req.Method != "GET" ***REMOVED***
				return nil, fmt.Errorf("expected GET method, got %s", req.Method)
			***REMOVED***

			key := types.SwarmUnlockKeyResponse***REMOVED***
				UnlockKey: unlockKey,
			***REMOVED***

			b, err := json.Marshal(key)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(b)),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	resp, err := client.SwarmGetUnlockKey(context.Background())
	require.NoError(t, err)
	assert.Equal(t, unlockKey, resp.UnlockKey)
***REMOVED***
