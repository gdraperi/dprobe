package plugin

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func parseHeaders(headers http.Header) (map[string][]string, *types.AuthConfig) ***REMOVED***

	metaHeaders := map[string][]string***REMOVED******REMOVED***
	for k, v := range headers ***REMOVED***
		if strings.HasPrefix(k, "X-Meta-") ***REMOVED***
			metaHeaders[k] = v
		***REMOVED***
	***REMOVED***

	// Get X-Registry-Auth
	authEncoded := headers.Get("X-Registry-Auth")
	authConfig := &types.AuthConfig***REMOVED******REMOVED***
	if authEncoded != "" ***REMOVED***
		authJSON := base64.NewDecoder(base64.URLEncoding, strings.NewReader(authEncoded))
		if err := json.NewDecoder(authJSON).Decode(authConfig); err != nil ***REMOVED***
			authConfig = &types.AuthConfig***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	return metaHeaders, authConfig
***REMOVED***

// parseRemoteRef parses the remote reference into a reference.Named
// returning the tag associated with the reference. In the case the
// given reference string includes both digest and tag, the returned
// reference will have the digest without the tag, but the tag will
// be returned.
func parseRemoteRef(remote string) (reference.Named, string, error) ***REMOVED***
	// Parse remote reference, supporting remotes with name and tag
	remoteRef, err := reference.ParseNormalizedNamed(remote)
	if err != nil ***REMOVED***
		return nil, "", err
	***REMOVED***

	type canonicalWithTag interface ***REMOVED***
		reference.Canonical
		Tag() string
	***REMOVED***

	if canonical, ok := remoteRef.(canonicalWithTag); ok ***REMOVED***
		remoteRef, err = reference.WithDigest(reference.TrimNamed(remoteRef), canonical.Digest())
		if err != nil ***REMOVED***
			return nil, "", err
		***REMOVED***
		return remoteRef, canonical.Tag(), nil
	***REMOVED***

	remoteRef = reference.TagNameOnly(remoteRef)

	return remoteRef, "", nil
***REMOVED***

func (pr *pluginRouter) getPrivileges(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	metaHeaders, authConfig := parseHeaders(r.Header)

	ref, _, err := parseRemoteRef(r.FormValue("remote"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	privileges, err := pr.backend.Privileges(ctx, ref, metaHeaders, authConfig)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return httputils.WriteJSON(w, http.StatusOK, privileges)
***REMOVED***

func (pr *pluginRouter) upgradePlugin(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return errors.Wrap(err, "failed to parse form")
	***REMOVED***

	var privileges types.PluginPrivileges
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&privileges); err != nil ***REMOVED***
		return errors.Wrap(err, "failed to parse privileges")
	***REMOVED***
	if dec.More() ***REMOVED***
		return errors.New("invalid privileges")
	***REMOVED***

	metaHeaders, authConfig := parseHeaders(r.Header)
	ref, tag, err := parseRemoteRef(r.FormValue("remote"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	name, err := getName(ref, tag, vars["name"])
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	w.Header().Set("Docker-Plugin-Name", name)

	w.Header().Set("Content-Type", "application/json")
	output := ioutils.NewWriteFlusher(w)

	if err := pr.backend.Upgrade(ctx, ref, name, metaHeaders, authConfig, privileges, output); err != nil ***REMOVED***
		if !output.Flushed() ***REMOVED***
			return err
		***REMOVED***
		output.Write(streamformatter.FormatError(err))
	***REMOVED***

	return nil
***REMOVED***

func (pr *pluginRouter) pullPlugin(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return errors.Wrap(err, "failed to parse form")
	***REMOVED***

	var privileges types.PluginPrivileges
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&privileges); err != nil ***REMOVED***
		return errors.Wrap(err, "failed to parse privileges")
	***REMOVED***
	if dec.More() ***REMOVED***
		return errors.New("invalid privileges")
	***REMOVED***

	metaHeaders, authConfig := parseHeaders(r.Header)
	ref, tag, err := parseRemoteRef(r.FormValue("remote"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	name, err := getName(ref, tag, r.FormValue("name"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	w.Header().Set("Docker-Plugin-Name", name)

	w.Header().Set("Content-Type", "application/json")
	output := ioutils.NewWriteFlusher(w)

	if err := pr.backend.Pull(ctx, ref, name, metaHeaders, authConfig, privileges, output); err != nil ***REMOVED***
		if !output.Flushed() ***REMOVED***
			return err
		***REMOVED***
		output.Write(streamformatter.FormatError(err))
	***REMOVED***

	return nil
***REMOVED***

func getName(ref reference.Named, tag, name string) (string, error) ***REMOVED***
	if name == "" ***REMOVED***
		if _, ok := ref.(reference.Canonical); ok ***REMOVED***
			trimmed := reference.TrimNamed(ref)
			if tag != "" ***REMOVED***
				nt, err := reference.WithTag(trimmed, tag)
				if err != nil ***REMOVED***
					return "", err
				***REMOVED***
				name = reference.FamiliarString(nt)
			***REMOVED*** else ***REMOVED***
				name = reference.FamiliarString(reference.TagNameOnly(trimmed))
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			name = reference.FamiliarString(ref)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		localRef, err := reference.ParseNormalizedNamed(name)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		if _, ok := localRef.(reference.Canonical); ok ***REMOVED***
			return "", errors.New("cannot use digest in plugin tag")
		***REMOVED***
		if reference.IsNameOnly(localRef) ***REMOVED***
			// TODO: log change in name to out stream
			name = reference.FamiliarString(reference.TagNameOnly(localRef))
		***REMOVED***
	***REMOVED***
	return name, nil
***REMOVED***

func (pr *pluginRouter) createPlugin(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	options := &types.PluginCreateOptions***REMOVED***
		RepoName: r.FormValue("name")***REMOVED***

	if err := pr.backend.CreateFromContext(ctx, r.Body, options); err != nil ***REMOVED***
		return err
	***REMOVED***
	//TODO: send progress bar
	w.WriteHeader(http.StatusNoContent)
	return nil
***REMOVED***

func (pr *pluginRouter) enablePlugin(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	name := vars["name"]
	timeout, err := strconv.Atoi(r.Form.Get("timeout"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	config := &types.PluginEnableConfig***REMOVED***Timeout: timeout***REMOVED***

	return pr.backend.Enable(name, config)
***REMOVED***

func (pr *pluginRouter) disablePlugin(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	name := vars["name"]
	config := &types.PluginDisableConfig***REMOVED***
		ForceDisable: httputils.BoolValue(r, "force"),
	***REMOVED***

	return pr.backend.Disable(name, config)
***REMOVED***

func (pr *pluginRouter) removePlugin(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	name := vars["name"]
	config := &types.PluginRmConfig***REMOVED***
		ForceRemove: httputils.BoolValue(r, "force"),
	***REMOVED***
	return pr.backend.Remove(name, config)
***REMOVED***

func (pr *pluginRouter) pushPlugin(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return errors.Wrap(err, "failed to parse form")
	***REMOVED***

	metaHeaders, authConfig := parseHeaders(r.Header)

	w.Header().Set("Content-Type", "application/json")
	output := ioutils.NewWriteFlusher(w)

	if err := pr.backend.Push(ctx, vars["name"], metaHeaders, authConfig, output); err != nil ***REMOVED***
		if !output.Flushed() ***REMOVED***
			return err
		***REMOVED***
		output.Write(streamformatter.FormatError(err))
	***REMOVED***
	return nil
***REMOVED***

func (pr *pluginRouter) setPlugin(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	var args []string
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := pr.backend.Set(vars["name"], args); err != nil ***REMOVED***
		return err
	***REMOVED***
	w.WriteHeader(http.StatusNoContent)
	return nil
***REMOVED***

func (pr *pluginRouter) listPlugins(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	pluginFilters, err := filters.FromJSON(r.Form.Get("filters"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	l, err := pr.backend.List(pluginFilters)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return httputils.WriteJSON(w, http.StatusOK, l)
***REMOVED***

func (pr *pluginRouter) inspectPlugin(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	result, err := pr.backend.Inspect(vars["name"])
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return httputils.WriteJSON(w, http.StatusOK, result)
***REMOVED***
