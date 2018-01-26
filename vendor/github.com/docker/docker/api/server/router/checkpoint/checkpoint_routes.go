package checkpoint

import (
	"encoding/json"
	"net/http"

	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

func (s *checkpointRouter) postContainerCheckpoint(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	var options types.CheckpointCreateOptions

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&options); err != nil ***REMOVED***
		return err
	***REMOVED***

	err := s.backend.CheckpointCreate(vars["name"], options)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	w.WriteHeader(http.StatusCreated)
	return nil
***REMOVED***

func (s *checkpointRouter) getContainerCheckpoints(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	checkpoints, err := s.backend.CheckpointList(vars["name"], types.CheckpointListOptions***REMOVED***
		CheckpointDir: r.Form.Get("dir"),
	***REMOVED***)

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, checkpoints)
***REMOVED***

func (s *checkpointRouter) deleteContainerCheckpoint(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	err := s.backend.CheckpointDelete(vars["name"], types.CheckpointDeleteOptions***REMOVED***
		CheckpointDir: r.Form.Get("dir"),
		CheckpointID:  vars["checkpoint"],
	***REMOVED***)

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	w.WriteHeader(http.StatusNoContent)
	return nil
***REMOVED***
