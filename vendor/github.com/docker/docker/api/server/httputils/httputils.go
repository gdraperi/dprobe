package httputils

import (
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/docker/docker/errdefs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type contextKey string

// APIVersionKey is the client's requested API version.
const APIVersionKey contextKey = "api-version"

// APIFunc is an adapter to allow the use of ordinary functions as Docker API endpoints.
// Any function that has the appropriate signature can be registered as an API endpoint (e.g. getVersion).
type APIFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error

// HijackConnection interrupts the http response writer to get the
// underlying connection and operate with it.
func HijackConnection(w http.ResponseWriter) (io.ReadCloser, io.Writer, error) ***REMOVED***
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	// Flush the options to make sure the client sets the raw mode
	conn.Write([]byte***REMOVED******REMOVED***)
	return conn, conn, nil
***REMOVED***

// CloseStreams ensures that a list for http streams are properly closed.
func CloseStreams(streams ...interface***REMOVED******REMOVED***) ***REMOVED***
	for _, stream := range streams ***REMOVED***
		if tcpc, ok := stream.(interface ***REMOVED***
			CloseWrite() error
		***REMOVED***); ok ***REMOVED***
			tcpc.CloseWrite()
		***REMOVED*** else if closer, ok := stream.(io.Closer); ok ***REMOVED***
			closer.Close()
		***REMOVED***
	***REMOVED***
***REMOVED***

// CheckForJSON makes sure that the request's Content-Type is application/json.
func CheckForJSON(r *http.Request) error ***REMOVED***
	ct := r.Header.Get("Content-Type")

	// No Content-Type header is ok as long as there's no Body
	if ct == "" ***REMOVED***
		if r.Body == nil || r.ContentLength == 0 ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	// Otherwise it better be json
	if matchesContentType(ct, "application/json") ***REMOVED***
		return nil
	***REMOVED***
	return errdefs.InvalidParameter(errors.Errorf("Content-Type specified (%s) must be 'application/json'", ct))
***REMOVED***

// ParseForm ensures the request form is parsed even with invalid content types.
// If we don't do this, POST method without Content-type (even with empty body) will fail.
func ParseForm(r *http.Request) error ***REMOVED***
	if r == nil ***REMOVED***
		return nil
	***REMOVED***
	if err := r.ParseForm(); err != nil && !strings.HasPrefix(err.Error(), "mime:") ***REMOVED***
		return errdefs.InvalidParameter(err)
	***REMOVED***
	return nil
***REMOVED***

// VersionFromContext returns an API version from the context using APIVersionKey.
// It panics if the context value does not have version.Version type.
func VersionFromContext(ctx context.Context) string ***REMOVED***
	if ctx == nil ***REMOVED***
		return ""
	***REMOVED***

	if val := ctx.Value(APIVersionKey); val != nil ***REMOVED***
		return val.(string)
	***REMOVED***

	return ""
***REMOVED***

// matchesContentType validates the content type against the expected one
func matchesContentType(contentType, expectedType string) bool ***REMOVED***
	mimetype, _, err := mime.ParseMediaType(contentType)
	if err != nil ***REMOVED***
		logrus.Errorf("Error parsing media type: %s error: %v", contentType, err)
	***REMOVED***
	return err == nil && mimetype == expectedType
***REMOVED***
