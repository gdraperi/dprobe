package client

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

// TestPingFail tests that when a server sends a non-successful response that we
// can still grab API details, when set.
// Some of this is just excercising the code paths to make sure there are no
// panics.
func TestPingFail(t *testing.T) ***REMOVED***
	var withHeader bool
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			resp := &http.Response***REMOVED***StatusCode: http.StatusInternalServerError***REMOVED***
			if withHeader ***REMOVED***
				resp.Header = http.Header***REMOVED******REMOVED***
				resp.Header.Set("API-Version", "awesome")
				resp.Header.Set("Docker-Experimental", "true")
			***REMOVED***
			resp.Body = ioutil.NopCloser(strings.NewReader("some error with the server"))
			return resp, nil
		***REMOVED***),
	***REMOVED***

	ping, err := client.Ping(context.Background())
	assert.Error(t, err)
	assert.Equal(t, false, ping.Experimental)
	assert.Equal(t, "", ping.APIVersion)

	withHeader = true
	ping2, err := client.Ping(context.Background())
	assert.Error(t, err)
	assert.Equal(t, true, ping2.Experimental)
	assert.Equal(t, "awesome", ping2.APIVersion)
***REMOVED***

// TestPingWithError tests the case where there is a protocol error in the ping.
// This test is mostly just testing that there are no panics in this code path.
func TestPingWithError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			resp := &http.Response***REMOVED***StatusCode: http.StatusInternalServerError***REMOVED***
			resp.Header = http.Header***REMOVED******REMOVED***
			resp.Header.Set("API-Version", "awesome")
			resp.Header.Set("Docker-Experimental", "true")
			resp.Body = ioutil.NopCloser(strings.NewReader("some error with the server"))
			return resp, errors.New("some error")
		***REMOVED***),
	***REMOVED***

	ping, err := client.Ping(context.Background())
	assert.Error(t, err)
	assert.Equal(t, false, ping.Experimental)
	assert.Equal(t, "", ping.APIVersion)
***REMOVED***

// TestPingSuccess tests that we are able to get the expected API headers/ping
// details on success.
func TestPingSuccess(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			resp := &http.Response***REMOVED***StatusCode: http.StatusInternalServerError***REMOVED***
			resp.Header = http.Header***REMOVED******REMOVED***
			resp.Header.Set("API-Version", "awesome")
			resp.Header.Set("Docker-Experimental", "true")
			resp.Body = ioutil.NopCloser(strings.NewReader("some error with the server"))
			return resp, nil
		***REMOVED***),
	***REMOVED***
	ping, err := client.Ping(context.Background())
	assert.Error(t, err)
	assert.Equal(t, true, ping.Experimental)
	assert.Equal(t, "awesome", ping.APIVersion)
***REMOVED***
