package session

import (
	"net/http"

	"github.com/docker/docker/errdefs"
	"golang.org/x/net/context"
)

func (sr *sessionRouter) startSession(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	err := sr.backend.HandleHTTPRequest(ctx, w, r)
	if err != nil ***REMOVED***
		return errdefs.InvalidParameter(err)
	***REMOVED***
	return nil
***REMOVED***
