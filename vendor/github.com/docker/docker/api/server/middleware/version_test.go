package middleware

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/docker/docker/api/server/httputils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestVersionMiddlewareVersion(t *testing.T) ***REMOVED***
	defaultVersion := "1.10.0"
	minVersion := "1.2.0"
	expectedVersion := defaultVersion
	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
		v := httputils.VersionFromContext(ctx)
		assert.Equal(t, expectedVersion, v)
		return nil
	***REMOVED***

	m := NewVersionMiddleware(defaultVersion, defaultVersion, minVersion)
	h := m.WrapHandler(handler)

	req, _ := http.NewRequest("GET", "/containers/json", nil)
	resp := httptest.NewRecorder()
	ctx := context.Background()

	tests := []struct ***REMOVED***
		reqVersion      string
		expectedVersion string
		errString       string
	***REMOVED******REMOVED***
		***REMOVED***
			expectedVersion: "1.10.0",
		***REMOVED***,
		***REMOVED***
			reqVersion:      "1.9.0",
			expectedVersion: "1.9.0",
		***REMOVED***,
		***REMOVED***
			reqVersion: "0.1",
			errString:  "client version 0.1 is too old. Minimum supported API version is 1.2.0, please upgrade your client to a newer version",
		***REMOVED***,
		***REMOVED***
			reqVersion: "9999.9999",
			errString:  "client version 9999.9999 is too new. Maximum supported API version is 1.10.0",
		***REMOVED***,
	***REMOVED***

	for _, test := range tests ***REMOVED***
		expectedVersion = test.expectedVersion

		err := h(ctx, resp, req, map[string]string***REMOVED***"version": test.reqVersion***REMOVED***)

		if test.errString != "" ***REMOVED***
			assert.EqualError(t, err, test.errString)
		***REMOVED*** else ***REMOVED***
			assert.NoError(t, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestVersionMiddlewareWithErrorsReturnsHeaders(t *testing.T) ***REMOVED***
	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
		v := httputils.VersionFromContext(ctx)
		assert.NotEmpty(t, v)
		return nil
	***REMOVED***

	defaultVersion := "1.10.0"
	minVersion := "1.2.0"
	m := NewVersionMiddleware(defaultVersion, defaultVersion, minVersion)
	h := m.WrapHandler(handler)

	req, _ := http.NewRequest("GET", "/containers/json", nil)
	resp := httptest.NewRecorder()
	ctx := context.Background()

	vars := map[string]string***REMOVED***"version": "0.1"***REMOVED***
	err := h(ctx, resp, req, vars)
	assert.Error(t, err)

	hdr := resp.Result().Header
	assert.Contains(t, hdr.Get("Server"), "Docker/"+defaultVersion)
	assert.Contains(t, hdr.Get("Server"), runtime.GOOS)
	assert.Equal(t, hdr.Get("API-Version"), defaultVersion)
	assert.Equal(t, hdr.Get("OSType"), runtime.GOOS)
***REMOVED***
