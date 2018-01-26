package container

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/versions"
	"golang.org/x/net/context"
)

type pathError struct***REMOVED******REMOVED***

func (pathError) Error() string ***REMOVED***
	return "Path cannot be empty"
***REMOVED***

func (pathError) InvalidParameter() ***REMOVED******REMOVED***

// postContainersCopy is deprecated in favor of getContainersArchive.
func (s *containerRouter) postContainersCopy(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	// Deprecated since 1.8, Errors out since 1.12
	version := httputils.VersionFromContext(ctx)
	if versions.GreaterThanOrEqualTo(version, "1.24") ***REMOVED***
		w.WriteHeader(http.StatusNotFound)
		return nil
	***REMOVED***
	if err := httputils.CheckForJSON(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	cfg := types.CopyConfig***REMOVED******REMOVED***
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil ***REMOVED***
		return err
	***REMOVED***

	if cfg.Resource == "" ***REMOVED***
		return pathError***REMOVED******REMOVED***
	***REMOVED***

	data, err := s.backend.ContainerCopy(vars["name"], cfg.Resource)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer data.Close()

	w.Header().Set("Content-Type", "application/x-tar")
	_, err = io.Copy(w, data)
	return err
***REMOVED***

// // Encode the stat to JSON, base64 encode, and place in a header.
func setContainerPathStatHeader(stat *types.ContainerPathStat, header http.Header) error ***REMOVED***
	statJSON, err := json.Marshal(stat)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	header.Set(
		"X-Docker-Container-Path-Stat",
		base64.StdEncoding.EncodeToString(statJSON),
	)

	return nil
***REMOVED***

func (s *containerRouter) headContainersArchive(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	v, err := httputils.ArchiveFormValues(r, vars)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	stat, err := s.backend.ContainerStatPath(v.Name, v.Path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return setContainerPathStatHeader(stat, w.Header())
***REMOVED***

func (s *containerRouter) getContainersArchive(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	v, err := httputils.ArchiveFormValues(r, vars)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	tarArchive, stat, err := s.backend.ContainerArchivePath(v.Name, v.Path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer tarArchive.Close()

	if err := setContainerPathStatHeader(stat, w.Header()); err != nil ***REMOVED***
		return err
	***REMOVED***

	w.Header().Set("Content-Type", "application/x-tar")
	_, err = io.Copy(w, tarArchive)

	return err
***REMOVED***

func (s *containerRouter) putContainersArchive(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	v, err := httputils.ArchiveFormValues(r, vars)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	noOverwriteDirNonDir := httputils.BoolValue(r, "noOverwriteDirNonDir")
	copyUIDGID := httputils.BoolValue(r, "copyUIDGID")

	return s.backend.ContainerExtractToDir(v.Name, v.Path, copyUIDGID, noOverwriteDirNonDir, r.Body)
***REMOVED***
