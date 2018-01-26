package daemon

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/builder/dockerfile"
	"github.com/docker/docker/builder/remotecontext"
	"github.com/docker/docker/dockerversion"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/image"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/pkg/errors"
)

// ImportImage imports an image, getting the archived layer data either from
// inConfig (if src is "-"), or from a URI specified in src. Progress output is
// written to outStream. Repository and tag names can optionally be given in
// the repo and tag arguments, respectively.
func (daemon *Daemon) ImportImage(src string, repository, os string, tag string, msg string, inConfig io.ReadCloser, outStream io.Writer, changes []string) error ***REMOVED***
	var (
		rc     io.ReadCloser
		resp   *http.Response
		newRef reference.Named
	)

	// Default the operating system if not supplied.
	if os == "" ***REMOVED***
		os = runtime.GOOS
	***REMOVED***

	if repository != "" ***REMOVED***
		var err error
		newRef, err = reference.ParseNormalizedNamed(repository)
		if err != nil ***REMOVED***
			return errdefs.InvalidParameter(err)
		***REMOVED***
		if _, isCanonical := newRef.(reference.Canonical); isCanonical ***REMOVED***
			return errdefs.InvalidParameter(errors.New("cannot import digest reference"))
		***REMOVED***

		if tag != "" ***REMOVED***
			newRef, err = reference.WithTag(newRef, tag)
			if err != nil ***REMOVED***
				return errdefs.InvalidParameter(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	config, err := dockerfile.BuildFromConfig(&container.Config***REMOVED******REMOVED***, changes, os)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if src == "-" ***REMOVED***
		rc = inConfig
	***REMOVED*** else ***REMOVED***
		inConfig.Close()
		if len(strings.Split(src, "://")) == 1 ***REMOVED***
			src = "http://" + src
		***REMOVED***
		u, err := url.Parse(src)
		if err != nil ***REMOVED***
			return errdefs.InvalidParameter(err)
		***REMOVED***

		resp, err = remotecontext.GetWithStatusError(u.String())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		outStream.Write(streamformatter.FormatStatus("", "Downloading from %s", u))
		progressOutput := streamformatter.NewJSONProgressOutput(outStream, true)
		rc = progress.NewProgressReader(resp.Body, progressOutput, resp.ContentLength, "", "Importing")
	***REMOVED***

	defer rc.Close()
	if len(msg) == 0 ***REMOVED***
		msg = "Imported from " + src
	***REMOVED***

	inflatedLayerData, err := archive.DecompressStream(rc)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	l, err := daemon.layerStores[os].Register(inflatedLayerData, "")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer layer.ReleaseAndLog(daemon.layerStores[os], l)

	created := time.Now().UTC()
	imgConfig, err := json.Marshal(&image.Image***REMOVED***
		V1Image: image.V1Image***REMOVED***
			DockerVersion: dockerversion.Version,
			Config:        config,
			Architecture:  runtime.GOARCH,
			OS:            os,
			Created:       created,
			Comment:       msg,
		***REMOVED***,
		RootFS: &image.RootFS***REMOVED***
			Type:    "layers",
			DiffIDs: []layer.DiffID***REMOVED***l.DiffID()***REMOVED***,
		***REMOVED***,
		History: []image.History***REMOVED******REMOVED***
			Created: created,
			Comment: msg,
		***REMOVED******REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	id, err := daemon.imageStore.Create(imgConfig)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// FIXME: connect with commit code and call refstore directly
	if newRef != nil ***REMOVED***
		if err := daemon.TagImageWithReference(id, newRef); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	daemon.LogImageEvent(id.String(), id.String(), "import")
	outStream.Write(streamformatter.FormatStatus("", id.String()))
	return nil
***REMOVED***
