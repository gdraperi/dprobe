package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/types/versions"
	"golang.org/x/net/context"
)

// VersionMiddleware is a middleware that
// validates the client and server versions.
type VersionMiddleware struct ***REMOVED***
	serverVersion  string
	defaultVersion string
	minVersion     string
***REMOVED***

// NewVersionMiddleware creates a new VersionMiddleware
// with the default versions.
func NewVersionMiddleware(s, d, m string) VersionMiddleware ***REMOVED***
	return VersionMiddleware***REMOVED***
		serverVersion:  s,
		defaultVersion: d,
		minVersion:     m,
	***REMOVED***
***REMOVED***

type versionUnsupportedError struct ***REMOVED***
	version, minVersion, maxVersion string
***REMOVED***

func (e versionUnsupportedError) Error() string ***REMOVED***
	if e.minVersion != "" ***REMOVED***
		return fmt.Sprintf("client version %s is too old. Minimum supported API version is %s, please upgrade your client to a newer version", e.version, e.minVersion)
	***REMOVED***
	return fmt.Sprintf("client version %s is too new. Maximum supported API version is %s", e.version, e.maxVersion)
***REMOVED***

func (e versionUnsupportedError) InvalidParameter() ***REMOVED******REMOVED***

// WrapHandler returns a new handler function wrapping the previous one in the request chain.
func (v VersionMiddleware) WrapHandler(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error) func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
		w.Header().Set("Server", fmt.Sprintf("Docker/%s (%s)", v.serverVersion, runtime.GOOS))
		w.Header().Set("API-Version", v.defaultVersion)
		w.Header().Set("OSType", runtime.GOOS)

		apiVersion := vars["version"]
		if apiVersion == "" ***REMOVED***
			apiVersion = v.defaultVersion
		***REMOVED***
		if versions.LessThan(apiVersion, v.minVersion) ***REMOVED***
			return versionUnsupportedError***REMOVED***version: apiVersion, minVersion: v.minVersion***REMOVED***
		***REMOVED***
		if versions.GreaterThan(apiVersion, v.defaultVersion) ***REMOVED***
			return versionUnsupportedError***REMOVED***version: apiVersion, maxVersion: v.defaultVersion***REMOVED***
		***REMOVED***
		ctx = context.WithValue(ctx, httputils.APIVersionKey, apiVersion)
		return handler(ctx, w, r, vars)
	***REMOVED***

***REMOVED***
