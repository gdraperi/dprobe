package image

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/versions"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/docker/registry"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func (s *imageRouter) postCommit(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := httputils.CheckForJSON(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	cname := r.Form.Get("container")

	pause := httputils.BoolValue(r, "pause")
	version := httputils.VersionFromContext(ctx)
	if r.FormValue("pause") == "" && versions.GreaterThanOrEqualTo(version, "1.13") ***REMOVED***
		pause = true
	***REMOVED***

	c, _, _, err := s.decoder.DecodeConfig(r.Body)
	if err != nil && err != io.EOF ***REMOVED*** //Do not fail if body is empty.
		return err
	***REMOVED***

	commitCfg := &backend.ContainerCommitConfig***REMOVED***
		ContainerCommitConfig: types.ContainerCommitConfig***REMOVED***
			Pause:        pause,
			Repo:         r.Form.Get("repo"),
			Tag:          r.Form.Get("tag"),
			Author:       r.Form.Get("author"),
			Comment:      r.Form.Get("comment"),
			Config:       c,
			MergeConfigs: true,
		***REMOVED***,
		Changes: r.Form["changes"],
	***REMOVED***

	imgID, err := s.backend.Commit(cname, commitCfg)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusCreated, &types.IDResponse***REMOVED***ID: imgID***REMOVED***)
***REMOVED***

// Creates an image from Pull or from Import
func (s *imageRouter) postImagesCreate(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***

	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	var (
		image    = r.Form.Get("fromImage")
		repo     = r.Form.Get("repo")
		tag      = r.Form.Get("tag")
		message  = r.Form.Get("message")
		err      error
		output   = ioutils.NewWriteFlusher(w)
		platform = &specs.Platform***REMOVED******REMOVED***
	)
	defer output.Close()

	w.Header().Set("Content-Type", "application/json")

	version := httputils.VersionFromContext(ctx)
	if versions.GreaterThanOrEqualTo(version, "1.32") ***REMOVED***
		// TODO @jhowardmsft. The following environment variable is an interim
		// measure to allow the daemon to have a default platform if omitted by
		// the client. This allows LCOW and WCOW to work with a down-level CLI
		// for a short period of time, as the CLI changes can't be merged
		// until after the daemon changes have been merged. Once the CLI is
		// updated, this can be removed. PR for CLI is currently in
		// https://github.com/docker/cli/pull/474.
		apiPlatform := r.FormValue("platform")
		if system.LCOWSupported() && apiPlatform == "" ***REMOVED***
			apiPlatform = os.Getenv("LCOW_API_PLATFORM_IF_OMITTED")
		***REMOVED***
		platform = system.ParsePlatform(apiPlatform)
		if err = system.ValidatePlatform(platform); err != nil ***REMOVED***
			err = fmt.Errorf("invalid platform: %s", err)
		***REMOVED***
	***REMOVED***

	if err == nil ***REMOVED***
		if image != "" ***REMOVED*** //pull
			metaHeaders := map[string][]string***REMOVED******REMOVED***
			for k, v := range r.Header ***REMOVED***
				if strings.HasPrefix(k, "X-Meta-") ***REMOVED***
					metaHeaders[k] = v
				***REMOVED***
			***REMOVED***

			authEncoded := r.Header.Get("X-Registry-Auth")
			authConfig := &types.AuthConfig***REMOVED******REMOVED***
			if authEncoded != "" ***REMOVED***
				authJSON := base64.NewDecoder(base64.URLEncoding, strings.NewReader(authEncoded))
				if err := json.NewDecoder(authJSON).Decode(authConfig); err != nil ***REMOVED***
					// for a pull it is not an error if no auth was given
					// to increase compatibility with the existing api it is defaulting to be empty
					authConfig = &types.AuthConfig***REMOVED******REMOVED***
				***REMOVED***
			***REMOVED***
			err = s.backend.PullImage(ctx, image, tag, platform.OS, metaHeaders, authConfig, output)
		***REMOVED*** else ***REMOVED*** //import
			src := r.Form.Get("fromSrc")
			// 'err' MUST NOT be defined within this block, we need any error
			// generated from the download to be available to the output
			// stream processing below
			err = s.backend.ImportImage(src, repo, platform.OS, tag, message, r.Body, output, r.Form["changes"])
		***REMOVED***
	***REMOVED***
	if err != nil ***REMOVED***
		if !output.Flushed() ***REMOVED***
			return err
		***REMOVED***
		output.Write(streamformatter.FormatError(err))
	***REMOVED***

	return nil
***REMOVED***

func (s *imageRouter) postImagesPush(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	metaHeaders := map[string][]string***REMOVED******REMOVED***
	for k, v := range r.Header ***REMOVED***
		if strings.HasPrefix(k, "X-Meta-") ***REMOVED***
			metaHeaders[k] = v
		***REMOVED***
	***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***
	authConfig := &types.AuthConfig***REMOVED******REMOVED***

	authEncoded := r.Header.Get("X-Registry-Auth")
	if authEncoded != "" ***REMOVED***
		// the new format is to handle the authConfig as a header
		authJSON := base64.NewDecoder(base64.URLEncoding, strings.NewReader(authEncoded))
		if err := json.NewDecoder(authJSON).Decode(authConfig); err != nil ***REMOVED***
			// to increase compatibility to existing api it is defaulting to be empty
			authConfig = &types.AuthConfig***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// the old format is supported for compatibility if there was no authConfig header
		if err := json.NewDecoder(r.Body).Decode(authConfig); err != nil ***REMOVED***
			return errors.Wrap(errdefs.InvalidParameter(err), "Bad parameters and missing X-Registry-Auth")
		***REMOVED***
	***REMOVED***

	image := vars["name"]
	tag := r.Form.Get("tag")

	output := ioutils.NewWriteFlusher(w)
	defer output.Close()

	w.Header().Set("Content-Type", "application/json")

	if err := s.backend.PushImage(ctx, image, tag, metaHeaders, authConfig, output); err != nil ***REMOVED***
		if !output.Flushed() ***REMOVED***
			return err
		***REMOVED***
		output.Write(streamformatter.FormatError(err))
	***REMOVED***
	return nil
***REMOVED***

func (s *imageRouter) getImagesGet(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	w.Header().Set("Content-Type", "application/x-tar")

	output := ioutils.NewWriteFlusher(w)
	defer output.Close()
	var names []string
	if name, ok := vars["name"]; ok ***REMOVED***
		names = []string***REMOVED***name***REMOVED***
	***REMOVED*** else ***REMOVED***
		names = r.Form["names"]
	***REMOVED***

	if err := s.backend.ExportImage(names, output); err != nil ***REMOVED***
		if !output.Flushed() ***REMOVED***
			return err
		***REMOVED***
		output.Write(streamformatter.FormatError(err))
	***REMOVED***
	return nil
***REMOVED***

func (s *imageRouter) postImagesLoad(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***
	quiet := httputils.BoolValueOrDefault(r, "quiet", true)

	w.Header().Set("Content-Type", "application/json")

	output := ioutils.NewWriteFlusher(w)
	defer output.Close()
	if err := s.backend.LoadImage(r.Body, output, quiet); err != nil ***REMOVED***
		output.Write(streamformatter.FormatError(err))
	***REMOVED***
	return nil
***REMOVED***

type missingImageError struct***REMOVED******REMOVED***

func (missingImageError) Error() string ***REMOVED***
	return "image name cannot be blank"
***REMOVED***

func (missingImageError) InvalidParameter() ***REMOVED******REMOVED***

func (s *imageRouter) deleteImages(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	name := vars["name"]

	if strings.TrimSpace(name) == "" ***REMOVED***
		return missingImageError***REMOVED******REMOVED***
	***REMOVED***

	force := httputils.BoolValue(r, "force")
	prune := !httputils.BoolValue(r, "noprune")

	list, err := s.backend.ImageDelete(name, force, prune)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, list)
***REMOVED***

func (s *imageRouter) getImagesByName(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	imageInspect, err := s.backend.LookupImage(vars["name"])
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, imageInspect)
***REMOVED***

func (s *imageRouter) getImagesJSON(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	imageFilters, err := filters.FromJSON(r.Form.Get("filters"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	filterParam := r.Form.Get("filter")
	// FIXME(vdemeester) This has been deprecated in 1.13, and is target for removal for v17.12
	if filterParam != "" ***REMOVED***
		imageFilters.Add("reference", filterParam)
	***REMOVED***

	images, err := s.backend.Images(imageFilters, httputils.BoolValue(r, "all"), false)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, images)
***REMOVED***

func (s *imageRouter) getImagesHistory(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	name := vars["name"]
	history, err := s.backend.ImageHistory(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, history)
***REMOVED***

func (s *imageRouter) postImagesTag(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := s.backend.TagImage(vars["name"], r.Form.Get("repo"), r.Form.Get("tag")); err != nil ***REMOVED***
		return err
	***REMOVED***
	w.WriteHeader(http.StatusCreated)
	return nil
***REMOVED***

func (s *imageRouter) getImagesSearch(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***
	var (
		config      *types.AuthConfig
		authEncoded = r.Header.Get("X-Registry-Auth")
		headers     = map[string][]string***REMOVED******REMOVED***
	)

	if authEncoded != "" ***REMOVED***
		authJSON := base64.NewDecoder(base64.URLEncoding, strings.NewReader(authEncoded))
		if err := json.NewDecoder(authJSON).Decode(&config); err != nil ***REMOVED***
			// for a search it is not an error if no auth was given
			// to increase compatibility with the existing api it is defaulting to be empty
			config = &types.AuthConfig***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
	for k, v := range r.Header ***REMOVED***
		if strings.HasPrefix(k, "X-Meta-") ***REMOVED***
			headers[k] = v
		***REMOVED***
	***REMOVED***
	limit := registry.DefaultSearchLimit
	if r.Form.Get("limit") != "" ***REMOVED***
		limitValue, err := strconv.Atoi(r.Form.Get("limit"))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		limit = limitValue
	***REMOVED***
	query, err := s.backend.SearchRegistryForImages(ctx, r.Form.Get("filters"), r.Form.Get("term"), limit, config, headers)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return httputils.WriteJSON(w, http.StatusOK, query.Results)
***REMOVED***

func (s *imageRouter) postImagesPrune(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	pruneFilters, err := filters.FromJSON(r.Form.Get("filters"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	pruneReport, err := s.backend.ImagesPrune(ctx, pruneFilters)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return httputils.WriteJSON(w, http.StatusOK, pruneReport)
***REMOVED***
