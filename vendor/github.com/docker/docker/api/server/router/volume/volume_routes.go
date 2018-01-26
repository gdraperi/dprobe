package volume

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/types/filters"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/errdefs"
	"golang.org/x/net/context"
)

func (v *volumeRouter) getVolumesList(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	volumes, warnings, err := v.backend.Volumes(r.Form.Get("filters"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return httputils.WriteJSON(w, http.StatusOK, &volumetypes.VolumesListOKBody***REMOVED***Volumes: volumes, Warnings: warnings***REMOVED***)
***REMOVED***

func (v *volumeRouter) getVolumeByName(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	volume, err := v.backend.VolumeInspect(vars["name"])
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return httputils.WriteJSON(w, http.StatusOK, volume)
***REMOVED***

func (v *volumeRouter) postVolumesCreate(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := httputils.CheckForJSON(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	var req volumetypes.VolumesCreateBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil ***REMOVED***
		if err == io.EOF ***REMOVED***
			return errdefs.InvalidParameter(errors.New("got EOF while reading request body"))
		***REMOVED***
		return err
	***REMOVED***

	volume, err := v.backend.VolumeCreate(req.Name, req.Driver, req.DriverOpts, req.Labels)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return httputils.WriteJSON(w, http.StatusCreated, volume)
***REMOVED***

func (v *volumeRouter) deleteVolumes(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***
	force := httputils.BoolValue(r, "force")
	if err := v.backend.VolumeRm(vars["name"], force); err != nil ***REMOVED***
		return err
	***REMOVED***
	w.WriteHeader(http.StatusNoContent)
	return nil
***REMOVED***

func (v *volumeRouter) postVolumesPrune(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	pruneFilters, err := filters.FromJSON(r.Form.Get("filters"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	pruneReport, err := v.backend.VolumesPrune(ctx, pruneFilters)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return httputils.WriteJSON(w, http.StatusOK, pruneReport)
***REMOVED***
