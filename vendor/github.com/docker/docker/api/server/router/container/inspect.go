package container

import (
	"net/http"

	"github.com/docker/docker/api/server/httputils"
	"golang.org/x/net/context"
)

// getContainersByName inspects container's configuration and serializes it as json.
func (s *containerRouter) getContainersByName(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	displaySize := httputils.BoolValue(r, "size")

	version := httputils.VersionFromContext(ctx)
	json, err := s.backend.ContainerInspect(vars["name"], displaySize, version)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, json)
***REMOVED***
