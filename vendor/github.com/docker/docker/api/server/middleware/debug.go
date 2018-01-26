package middleware

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// DebugRequestMiddleware dumps the request to logger
func DebugRequestMiddleware(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error) func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
		logrus.Debugf("Calling %s %s", r.Method, r.RequestURI)

		if r.Method != "POST" ***REMOVED***
			return handler(ctx, w, r, vars)
		***REMOVED***
		if err := httputils.CheckForJSON(r); err != nil ***REMOVED***
			return handler(ctx, w, r, vars)
		***REMOVED***
		maxBodySize := 4096 // 4KB
		if r.ContentLength > int64(maxBodySize) ***REMOVED***
			return handler(ctx, w, r, vars)
		***REMOVED***

		body := r.Body
		bufReader := bufio.NewReaderSize(body, maxBodySize)
		r.Body = ioutils.NewReadCloserWrapper(bufReader, func() error ***REMOVED*** return body.Close() ***REMOVED***)

		b, err := bufReader.Peek(maxBodySize)
		if err != io.EOF ***REMOVED***
			// either there was an error reading, or the buffer is full (in which case the request is too large)
			return handler(ctx, w, r, vars)
		***REMOVED***

		var postForm map[string]interface***REMOVED******REMOVED***
		if err := json.Unmarshal(b, &postForm); err == nil ***REMOVED***
			maskSecretKeys(postForm, r.RequestURI)
			formStr, errMarshal := json.Marshal(postForm)
			if errMarshal == nil ***REMOVED***
				logrus.Debugf("form data: %s", string(formStr))
			***REMOVED*** else ***REMOVED***
				logrus.Debugf("form data: %q", postForm)
			***REMOVED***
		***REMOVED***

		return handler(ctx, w, r, vars)
	***REMOVED***
***REMOVED***

func maskSecretKeys(inp interface***REMOVED******REMOVED***, path string) ***REMOVED***
	// Remove any query string from the path
	idx := strings.Index(path, "?")
	if idx != -1 ***REMOVED***
		path = path[:idx]
	***REMOVED***
	// Remove trailing / characters
	path = strings.TrimRight(path, "/")

	if arr, ok := inp.([]interface***REMOVED******REMOVED***); ok ***REMOVED***
		for _, f := range arr ***REMOVED***
			maskSecretKeys(f, path)
		***REMOVED***
		return
	***REMOVED***

	if form, ok := inp.(map[string]interface***REMOVED******REMOVED***); ok ***REMOVED***
	loop0:
		for k, v := range form ***REMOVED***
			for _, m := range []string***REMOVED***"password", "secret", "jointoken", "unlockkey", "signingcakey"***REMOVED*** ***REMOVED***
				if strings.EqualFold(m, k) ***REMOVED***
					form[k] = "*****"
					continue loop0
				***REMOVED***
			***REMOVED***
			maskSecretKeys(v, path)
		***REMOVED***

		// Route-specific redactions
		if strings.HasSuffix(path, "/secrets/create") ***REMOVED***
			for k := range form ***REMOVED***
				if k == "Data" ***REMOVED***
					form[k] = "*****"
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
