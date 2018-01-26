package distribution

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/client/auth"
	"github.com/docker/distribution/registry/client/transport"
	"github.com/docker/docker/distribution/metadata"
	"github.com/docker/docker/distribution/xfer"
	"github.com/docker/docker/dockerversion"
	"github.com/docker/docker/image"
	"github.com/docker/docker/image/v1"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/registry"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type v1Puller struct ***REMOVED***
	v1IDService *metadata.V1IDService
	endpoint    registry.APIEndpoint
	config      *ImagePullConfig
	repoInfo    *registry.RepositoryInfo
	session     *registry.Session
***REMOVED***

func (p *v1Puller) Pull(ctx context.Context, ref reference.Named, os string) error ***REMOVED***
	if _, isCanonical := ref.(reference.Canonical); isCanonical ***REMOVED***
		// Allowing fallback, because HTTPS v1 is before HTTP v2
		return fallbackError***REMOVED***err: ErrNoSupport***REMOVED***Err: errors.New("Cannot pull by digest with v1 registry")***REMOVED******REMOVED***
	***REMOVED***

	tlsConfig, err := p.config.RegistryService.TLSConfig(p.repoInfo.Index.Name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// Adds Docker-specific headers as well as user-specified headers (metaHeaders)
	tr := transport.NewTransport(
		// TODO(tiborvass): was ReceiveTimeout
		registry.NewTransport(tlsConfig),
		registry.Headers(dockerversion.DockerUserAgent(ctx), p.config.MetaHeaders)...,
	)
	client := registry.HTTPClient(tr)
	v1Endpoint := p.endpoint.ToV1Endpoint(dockerversion.DockerUserAgent(ctx), p.config.MetaHeaders)
	p.session, err = registry.NewSession(client, p.config.AuthConfig, v1Endpoint)
	if err != nil ***REMOVED***
		// TODO(dmcgowan): Check if should fallback
		logrus.Debugf("Fallback from error: %s", err)
		return fallbackError***REMOVED***err: err***REMOVED***
	***REMOVED***
	if err := p.pullRepository(ctx, ref); err != nil ***REMOVED***
		// TODO(dmcgowan): Check if should fallback
		return err
	***REMOVED***
	progress.Message(p.config.ProgressOutput, "", p.repoInfo.Name.Name()+": this image was pulled from a legacy registry.  Important: This registry version will not be supported in future versions of docker.")

	return nil
***REMOVED***

// Note use auth.Scope rather than reference.Named due to this warning causing Jenkins CI to fail:
// warning: ref can be github.com/docker/docker/vendor/github.com/docker/distribution/registry/client/auth.Scope (interfacer)
func (p *v1Puller) pullRepository(ctx context.Context, ref auth.Scope) error ***REMOVED***
	progress.Message(p.config.ProgressOutput, "", "Pulling repository "+p.repoInfo.Name.Name())

	tagged, isTagged := ref.(reference.NamedTagged)

	repoData, err := p.session.GetRepositoryData(p.repoInfo.Name)
	if err != nil ***REMOVED***
		if strings.Contains(err.Error(), "HTTP code: 404") ***REMOVED***
			if isTagged ***REMOVED***
				return fmt.Errorf("Error: image %s:%s not found", reference.Path(p.repoInfo.Name), tagged.Tag())
			***REMOVED***
			return fmt.Errorf("Error: image %s not found", reference.Path(p.repoInfo.Name))
		***REMOVED***
		// Unexpected HTTP error
		return err
	***REMOVED***

	logrus.Debug("Retrieving the tag list")
	var tagsList map[string]string
	if !isTagged ***REMOVED***
		tagsList, err = p.session.GetRemoteTags(repoData.Endpoints, p.repoInfo.Name)
	***REMOVED*** else ***REMOVED***
		var tagID string
		tagsList = make(map[string]string)
		tagID, err = p.session.GetRemoteTag(repoData.Endpoints, p.repoInfo.Name, tagged.Tag())
		if err == registry.ErrRepoNotFound ***REMOVED***
			return fmt.Errorf("Tag %s not found in repository %s", tagged.Tag(), p.repoInfo.Name.Name())
		***REMOVED***
		tagsList[tagged.Tag()] = tagID
	***REMOVED***
	if err != nil ***REMOVED***
		logrus.Errorf("unable to get remote tags: %s", err)
		return err
	***REMOVED***

	for tag, id := range tagsList ***REMOVED***
		repoData.ImgList[id] = &registry.ImgData***REMOVED***
			ID:       id,
			Tag:      tag,
			Checksum: "",
		***REMOVED***
	***REMOVED***

	layersDownloaded := false
	for _, imgData := range repoData.ImgList ***REMOVED***
		if isTagged && imgData.Tag != tagged.Tag() ***REMOVED***
			continue
		***REMOVED***

		err := p.downloadImage(ctx, repoData, imgData, &layersDownloaded)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	writeStatus(reference.FamiliarString(ref), p.config.ProgressOutput, layersDownloaded)
	return nil
***REMOVED***

func (p *v1Puller) downloadImage(ctx context.Context, repoData *registry.RepositoryData, img *registry.ImgData, layersDownloaded *bool) error ***REMOVED***
	if img.Tag == "" ***REMOVED***
		logrus.Debugf("Image (id: %s) present in this repository but untagged, skipping", img.ID)
		return nil
	***REMOVED***

	localNameRef, err := reference.WithTag(p.repoInfo.Name, img.Tag)
	if err != nil ***REMOVED***
		retErr := fmt.Errorf("Image (id: %s) has invalid tag: %s", img.ID, img.Tag)
		logrus.Debug(retErr.Error())
		return retErr
	***REMOVED***

	if err := v1.ValidateID(img.ID); err != nil ***REMOVED***
		return err
	***REMOVED***

	progress.Updatef(p.config.ProgressOutput, stringid.TruncateID(img.ID), "Pulling image (%s) from %s", img.Tag, p.repoInfo.Name.Name())
	success := false
	var lastErr error
	for _, ep := range p.repoInfo.Index.Mirrors ***REMOVED***
		ep += "v1/"
		progress.Updatef(p.config.ProgressOutput, stringid.TruncateID(img.ID), fmt.Sprintf("Pulling image (%s) from %s, mirror: %s", img.Tag, p.repoInfo.Name.Name(), ep))
		if err = p.pullImage(ctx, img.ID, ep, localNameRef, layersDownloaded); err != nil ***REMOVED***
			// Don't report errors when pulling from mirrors.
			logrus.Debugf("Error pulling image (%s) from %s, mirror: %s, %s", img.Tag, p.repoInfo.Name.Name(), ep, err)
			continue
		***REMOVED***
		success = true
		break
	***REMOVED***
	if !success ***REMOVED***
		for _, ep := range repoData.Endpoints ***REMOVED***
			progress.Updatef(p.config.ProgressOutput, stringid.TruncateID(img.ID), "Pulling image (%s) from %s, endpoint: %s", img.Tag, p.repoInfo.Name.Name(), ep)
			if err = p.pullImage(ctx, img.ID, ep, localNameRef, layersDownloaded); err != nil ***REMOVED***
				// It's not ideal that only the last error is returned, it would be better to concatenate the errors.
				// As the error is also given to the output stream the user will see the error.
				lastErr = err
				progress.Updatef(p.config.ProgressOutput, stringid.TruncateID(img.ID), "Error pulling image (%s) from %s, endpoint: %s, %s", img.Tag, p.repoInfo.Name.Name(), ep, err)
				continue
			***REMOVED***
			success = true
			break
		***REMOVED***
	***REMOVED***
	if !success ***REMOVED***
		err := fmt.Errorf("Error pulling image (%s) from %s, %v", img.Tag, p.repoInfo.Name.Name(), lastErr)
		progress.Update(p.config.ProgressOutput, stringid.TruncateID(img.ID), err.Error())
		return err
	***REMOVED***
	return nil
***REMOVED***

func (p *v1Puller) pullImage(ctx context.Context, v1ID, endpoint string, localNameRef reference.Named, layersDownloaded *bool) (err error) ***REMOVED***
	var history []string
	history, err = p.session.GetRemoteHistory(v1ID, endpoint)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if len(history) < 1 ***REMOVED***
		return fmt.Errorf("empty history for image %s", v1ID)
	***REMOVED***
	progress.Update(p.config.ProgressOutput, stringid.TruncateID(v1ID), "Pulling dependent layers")

	var (
		descriptors []xfer.DownloadDescriptor
		newHistory  []image.History
		imgJSON     []byte
		imgSize     int64
	)

	// Iterate over layers, in order from bottom-most to top-most. Download
	// config for all layers and create descriptors.
	for i := len(history) - 1; i >= 0; i-- ***REMOVED***
		v1LayerID := history[i]
		imgJSON, imgSize, err = p.downloadLayerConfig(v1LayerID, endpoint)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Create a new-style config from the legacy configs
		h, err := v1.HistoryFromConfig(imgJSON, false)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		newHistory = append(newHistory, h)

		layerDescriptor := &v1LayerDescriptor***REMOVED***
			v1LayerID:        v1LayerID,
			indexName:        p.repoInfo.Index.Name,
			endpoint:         endpoint,
			v1IDService:      p.v1IDService,
			layersDownloaded: layersDownloaded,
			layerSize:        imgSize,
			session:          p.session,
		***REMOVED***

		descriptors = append(descriptors, layerDescriptor)
	***REMOVED***

	rootFS := image.NewRootFS()
	resultRootFS, release, err := p.config.DownloadManager.Download(ctx, *rootFS, "", descriptors, p.config.ProgressOutput)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer release()

	config, err := v1.MakeConfigFromV1Config(imgJSON, &resultRootFS, newHistory)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	imageID, err := p.config.ImageStore.Put(config)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if p.config.ReferenceStore != nil ***REMOVED***
		if err := p.config.ReferenceStore.AddTag(localNameRef, imageID, true); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (p *v1Puller) downloadLayerConfig(v1LayerID, endpoint string) (imgJSON []byte, imgSize int64, err error) ***REMOVED***
	progress.Update(p.config.ProgressOutput, stringid.TruncateID(v1LayerID), "Pulling metadata")

	retries := 5
	for j := 1; j <= retries; j++ ***REMOVED***
		imgJSON, imgSize, err := p.session.GetRemoteImageJSON(v1LayerID, endpoint)
		if err != nil && j == retries ***REMOVED***
			progress.Update(p.config.ProgressOutput, stringid.TruncateID(v1LayerID), "Error pulling layer metadata")
			return nil, 0, err
		***REMOVED*** else if err != nil ***REMOVED***
			time.Sleep(time.Duration(j) * 500 * time.Millisecond)
			continue
		***REMOVED***

		return imgJSON, imgSize, nil
	***REMOVED***

	// not reached
	return nil, 0, nil
***REMOVED***

type v1LayerDescriptor struct ***REMOVED***
	v1LayerID        string
	indexName        string
	endpoint         string
	v1IDService      *metadata.V1IDService
	layersDownloaded *bool
	layerSize        int64
	session          *registry.Session
	tmpFile          *os.File
***REMOVED***

func (ld *v1LayerDescriptor) Key() string ***REMOVED***
	return "v1:" + ld.v1LayerID
***REMOVED***

func (ld *v1LayerDescriptor) ID() string ***REMOVED***
	return stringid.TruncateID(ld.v1LayerID)
***REMOVED***

func (ld *v1LayerDescriptor) DiffID() (layer.DiffID, error) ***REMOVED***
	return ld.v1IDService.Get(ld.v1LayerID, ld.indexName)
***REMOVED***

func (ld *v1LayerDescriptor) Download(ctx context.Context, progressOutput progress.Output) (io.ReadCloser, int64, error) ***REMOVED***
	progress.Update(progressOutput, ld.ID(), "Pulling fs layer")
	layerReader, err := ld.session.GetRemoteImageLayer(ld.v1LayerID, ld.endpoint, ld.layerSize)
	if err != nil ***REMOVED***
		progress.Update(progressOutput, ld.ID(), "Error pulling dependent layers")
		if uerr, ok := err.(*url.Error); ok ***REMOVED***
			err = uerr.Err
		***REMOVED***
		if terr, ok := err.(net.Error); ok && terr.Timeout() ***REMOVED***
			return nil, 0, err
		***REMOVED***
		return nil, 0, xfer.DoNotRetry***REMOVED***Err: err***REMOVED***
	***REMOVED***
	*ld.layersDownloaded = true

	ld.tmpFile, err = ioutil.TempFile("", "GetImageBlob")
	if err != nil ***REMOVED***
		layerReader.Close()
		return nil, 0, err
	***REMOVED***

	reader := progress.NewProgressReader(ioutils.NewCancelReadCloser(ctx, layerReader), progressOutput, ld.layerSize, ld.ID(), "Downloading")
	defer reader.Close()

	_, err = io.Copy(ld.tmpFile, reader)
	if err != nil ***REMOVED***
		ld.Close()
		return nil, 0, err
	***REMOVED***

	progress.Update(progressOutput, ld.ID(), "Download complete")

	logrus.Debugf("Downloaded %s to tempfile %s", ld.ID(), ld.tmpFile.Name())

	ld.tmpFile.Seek(0, 0)

	// hand off the temporary file to the download manager, so it will only
	// be closed once
	tmpFile := ld.tmpFile
	ld.tmpFile = nil

	return ioutils.NewReadCloserWrapper(tmpFile, func() error ***REMOVED***
		tmpFile.Close()
		err := os.RemoveAll(tmpFile.Name())
		if err != nil ***REMOVED***
			logrus.Errorf("Failed to remove temp file: %s", tmpFile.Name())
		***REMOVED***
		return err
	***REMOVED***), ld.layerSize, nil
***REMOVED***

func (ld *v1LayerDescriptor) Close() ***REMOVED***
	if ld.tmpFile != nil ***REMOVED***
		ld.tmpFile.Close()
		if err := os.RemoveAll(ld.tmpFile.Name()); err != nil ***REMOVED***
			logrus.Errorf("Failed to remove temp file: %s", ld.tmpFile.Name())
		***REMOVED***
		ld.tmpFile = nil
	***REMOVED***
***REMOVED***

func (ld *v1LayerDescriptor) Registered(diffID layer.DiffID) ***REMOVED***
	// Cache mapping from this layer's DiffID to the blobsum
	ld.v1IDService.Set(ld.v1LayerID, ld.indexName, diffID)
***REMOVED***
