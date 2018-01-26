package build

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/versions"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/docker/docker/pkg/system"
	units "github.com/docker/go-units"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type invalidIsolationError string

func (e invalidIsolationError) Error() string ***REMOVED***
	return fmt.Sprintf("Unsupported isolation: %q", string(e))
***REMOVED***

func (e invalidIsolationError) InvalidParameter() ***REMOVED******REMOVED***

func newImageBuildOptions(ctx context.Context, r *http.Request) (*types.ImageBuildOptions, error) ***REMOVED***
	version := httputils.VersionFromContext(ctx)
	options := &types.ImageBuildOptions***REMOVED******REMOVED***
	if httputils.BoolValue(r, "forcerm") && versions.GreaterThanOrEqualTo(version, "1.12") ***REMOVED***
		options.Remove = true
	***REMOVED*** else if r.FormValue("rm") == "" && versions.GreaterThanOrEqualTo(version, "1.12") ***REMOVED***
		options.Remove = true
	***REMOVED*** else ***REMOVED***
		options.Remove = httputils.BoolValue(r, "rm")
	***REMOVED***
	if httputils.BoolValue(r, "pull") && versions.GreaterThanOrEqualTo(version, "1.16") ***REMOVED***
		options.PullParent = true
	***REMOVED***

	options.Dockerfile = r.FormValue("dockerfile")
	options.SuppressOutput = httputils.BoolValue(r, "q")
	options.NoCache = httputils.BoolValue(r, "nocache")
	options.ForceRemove = httputils.BoolValue(r, "forcerm")
	options.MemorySwap = httputils.Int64ValueOrZero(r, "memswap")
	options.Memory = httputils.Int64ValueOrZero(r, "memory")
	options.CPUShares = httputils.Int64ValueOrZero(r, "cpushares")
	options.CPUPeriod = httputils.Int64ValueOrZero(r, "cpuperiod")
	options.CPUQuota = httputils.Int64ValueOrZero(r, "cpuquota")
	options.CPUSetCPUs = r.FormValue("cpusetcpus")
	options.CPUSetMems = r.FormValue("cpusetmems")
	options.CgroupParent = r.FormValue("cgroupparent")
	options.NetworkMode = r.FormValue("networkmode")
	options.Tags = r.Form["t"]
	options.ExtraHosts = r.Form["extrahosts"]
	options.SecurityOpt = r.Form["securityopt"]
	options.Squash = httputils.BoolValue(r, "squash")
	options.Target = r.FormValue("target")
	options.RemoteContext = r.FormValue("remote")
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
		p := system.ParsePlatform(apiPlatform)
		if err := system.ValidatePlatform(p); err != nil ***REMOVED***
			return nil, errdefs.InvalidParameter(errors.Errorf("invalid platform: %s", err))
		***REMOVED***
		options.Platform = p.OS
	***REMOVED***

	if r.Form.Get("shmsize") != "" ***REMOVED***
		shmSize, err := strconv.ParseInt(r.Form.Get("shmsize"), 10, 64)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		options.ShmSize = shmSize
	***REMOVED***

	if i := container.Isolation(r.FormValue("isolation")); i != "" ***REMOVED***
		if !container.Isolation.IsValid(i) ***REMOVED***
			return nil, invalidIsolationError(i)
		***REMOVED***
		options.Isolation = i
	***REMOVED***

	if runtime.GOOS != "windows" && options.SecurityOpt != nil ***REMOVED***
		return nil, errdefs.InvalidParameter(errors.New("The daemon on this platform does not support setting security options on build"))
	***REMOVED***

	var buildUlimits = []*units.Ulimit***REMOVED******REMOVED***
	ulimitsJSON := r.FormValue("ulimits")
	if ulimitsJSON != "" ***REMOVED***
		if err := json.Unmarshal([]byte(ulimitsJSON), &buildUlimits); err != nil ***REMOVED***
			return nil, errors.Wrap(errdefs.InvalidParameter(err), "error reading ulimit settings")
		***REMOVED***
		options.Ulimits = buildUlimits
	***REMOVED***

	// Note that there are two ways a --build-arg might appear in the
	// json of the query param:
	//     "foo":"bar"
	// and "foo":nil
	// The first is the normal case, ie. --build-arg foo=bar
	// or  --build-arg foo
	// where foo's value was picked up from an env var.
	// The second ("foo":nil) is where they put --build-arg foo
	// but "foo" isn't set as an env var. In that case we can't just drop
	// the fact they mentioned it, we need to pass that along to the builder
	// so that it can print a warning about "foo" being unused if there is
	// no "ARG foo" in the Dockerfile.
	buildArgsJSON := r.FormValue("buildargs")
	if buildArgsJSON != "" ***REMOVED***
		var buildArgs = map[string]*string***REMOVED******REMOVED***
		if err := json.Unmarshal([]byte(buildArgsJSON), &buildArgs); err != nil ***REMOVED***
			return nil, errors.Wrap(errdefs.InvalidParameter(err), "error reading build args")
		***REMOVED***
		options.BuildArgs = buildArgs
	***REMOVED***

	labelsJSON := r.FormValue("labels")
	if labelsJSON != "" ***REMOVED***
		var labels = map[string]string***REMOVED******REMOVED***
		if err := json.Unmarshal([]byte(labelsJSON), &labels); err != nil ***REMOVED***
			return nil, errors.Wrap(errdefs.InvalidParameter(err), "error reading labels")
		***REMOVED***
		options.Labels = labels
	***REMOVED***

	cacheFromJSON := r.FormValue("cachefrom")
	if cacheFromJSON != "" ***REMOVED***
		var cacheFrom = []string***REMOVED******REMOVED***
		if err := json.Unmarshal([]byte(cacheFromJSON), &cacheFrom); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		options.CacheFrom = cacheFrom
	***REMOVED***
	options.SessionID = r.FormValue("session")

	return options, nil
***REMOVED***

func (br *buildRouter) postPrune(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	report, err := br.backend.PruneCache(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return httputils.WriteJSON(w, http.StatusOK, report)
***REMOVED***

func (br *buildRouter) postBuild(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	var (
		notVerboseBuffer = bytes.NewBuffer(nil)
		version          = httputils.VersionFromContext(ctx)
	)

	w.Header().Set("Content-Type", "application/json")

	output := ioutils.NewWriteFlusher(w)
	defer output.Close()
	errf := func(err error) error ***REMOVED***
		if httputils.BoolValue(r, "q") && notVerboseBuffer.Len() > 0 ***REMOVED***
			output.Write(notVerboseBuffer.Bytes())
		***REMOVED***
		// Do not write the error in the http output if it's still empty.
		// This prevents from writing a 200(OK) when there is an internal error.
		if !output.Flushed() ***REMOVED***
			return err
		***REMOVED***
		_, err = w.Write(streamformatter.FormatError(err))
		if err != nil ***REMOVED***
			logrus.Warnf("could not write error response: %v", err)
		***REMOVED***
		return nil
	***REMOVED***

	buildOptions, err := newImageBuildOptions(ctx, r)
	if err != nil ***REMOVED***
		return errf(err)
	***REMOVED***
	buildOptions.AuthConfigs = getAuthConfigs(r.Header)

	if buildOptions.Squash && !br.daemon.HasExperimental() ***REMOVED***
		return errdefs.InvalidParameter(errors.New("squash is only supported with experimental mode"))
	***REMOVED***

	out := io.Writer(output)
	if buildOptions.SuppressOutput ***REMOVED***
		out = notVerboseBuffer
	***REMOVED***

	// Currently, only used if context is from a remote url.
	// Look at code in DetectContextFromRemoteURL for more information.
	createProgressReader := func(in io.ReadCloser) io.ReadCloser ***REMOVED***
		progressOutput := streamformatter.NewJSONProgressOutput(out, true)
		return progress.NewProgressReader(in, progressOutput, r.ContentLength, "Downloading context", buildOptions.RemoteContext)
	***REMOVED***

	wantAux := versions.GreaterThanOrEqualTo(version, "1.30")

	imgID, err := br.backend.Build(ctx, backend.BuildConfig***REMOVED***
		Source:         r.Body,
		Options:        buildOptions,
		ProgressWriter: buildProgressWriter(out, wantAux, createProgressReader),
	***REMOVED***)
	if err != nil ***REMOVED***
		return errf(err)
	***REMOVED***

	// Everything worked so if -q was provided the output from the daemon
	// should be just the image ID and we'll print that to stdout.
	if buildOptions.SuppressOutput ***REMOVED***
		fmt.Fprintln(streamformatter.NewStdoutWriter(output), imgID)
	***REMOVED***
	return nil
***REMOVED***

func getAuthConfigs(header http.Header) map[string]types.AuthConfig ***REMOVED***
	authConfigs := map[string]types.AuthConfig***REMOVED******REMOVED***
	authConfigsEncoded := header.Get("X-Registry-Config")

	if authConfigsEncoded == "" ***REMOVED***
		return authConfigs
	***REMOVED***

	authConfigsJSON := base64.NewDecoder(base64.URLEncoding, strings.NewReader(authConfigsEncoded))
	// Pulling an image does not error when no auth is provided so to remain
	// consistent with the existing api decode errors are ignored
	json.NewDecoder(authConfigsJSON).Decode(&authConfigs)
	return authConfigs
***REMOVED***

type syncWriter struct ***REMOVED***
	w  io.Writer
	mu sync.Mutex
***REMOVED***

func (s *syncWriter) Write(b []byte) (count int, err error) ***REMOVED***
	s.mu.Lock()
	count, err = s.w.Write(b)
	s.mu.Unlock()
	return
***REMOVED***

func buildProgressWriter(out io.Writer, wantAux bool, createProgressReader func(io.ReadCloser) io.ReadCloser) backend.ProgressWriter ***REMOVED***
	out = &syncWriter***REMOVED***w: out***REMOVED***

	var aux *streamformatter.AuxFormatter
	if wantAux ***REMOVED***
		aux = &streamformatter.AuxFormatter***REMOVED***Writer: out***REMOVED***
	***REMOVED***

	return backend.ProgressWriter***REMOVED***
		Output:             out,
		StdoutFormatter:    streamformatter.NewStdoutWriter(out),
		StderrFormatter:    streamformatter.NewStderrWriter(out),
		AuxFormatter:       aux,
		ProgressReaderFunc: createProgressReader,
	***REMOVED***
***REMOVED***
