package distribution

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/manifestlist"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/api/errcode"
	"github.com/docker/distribution/registry/client/auth"
	"github.com/docker/distribution/registry/client/transport"
	"github.com/docker/docker/distribution/metadata"
	"github.com/docker/docker/distribution/xfer"
	"github.com/docker/docker/image"
	"github.com/docker/docker/image/v1"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/system"
	refstore "github.com/docker/docker/reference"
	"github.com/docker/docker/registry"
	digest "github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

var (
	errRootFSMismatch = errors.New("layers from manifest don't match image configuration")
	errRootFSInvalid  = errors.New("invalid rootfs in image configuration")
)

// ImageConfigPullError is an error pulling the image config blob
// (only applies to schema2).
type ImageConfigPullError struct ***REMOVED***
	Err error
***REMOVED***

// Error returns the error string for ImageConfigPullError.
func (e ImageConfigPullError) Error() string ***REMOVED***
	return "error pulling image configuration: " + e.Err.Error()
***REMOVED***

type v2Puller struct ***REMOVED***
	V2MetadataService metadata.V2MetadataService
	endpoint          registry.APIEndpoint
	config            *ImagePullConfig
	repoInfo          *registry.RepositoryInfo
	repo              distribution.Repository
	// confirmedV2 is set to true if we confirm we're talking to a v2
	// registry. This is used to limit fallbacks to the v1 protocol.
	confirmedV2 bool
***REMOVED***

func (p *v2Puller) Pull(ctx context.Context, ref reference.Named, os string) (err error) ***REMOVED***
	// TODO(tiborvass): was ReceiveTimeout
	p.repo, p.confirmedV2, err = NewV2Repository(ctx, p.repoInfo, p.endpoint, p.config.MetaHeaders, p.config.AuthConfig, "pull")
	if err != nil ***REMOVED***
		logrus.Warnf("Error getting v2 registry: %v", err)
		return err
	***REMOVED***

	if err = p.pullV2Repository(ctx, ref, os); err != nil ***REMOVED***
		if _, ok := err.(fallbackError); ok ***REMOVED***
			return err
		***REMOVED***
		if continueOnError(err, p.endpoint.Mirror) ***REMOVED***
			return fallbackError***REMOVED***
				err:         err,
				confirmedV2: p.confirmedV2,
				transportOK: true,
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***

func (p *v2Puller) pullV2Repository(ctx context.Context, ref reference.Named, os string) (err error) ***REMOVED***
	var layersDownloaded bool
	if !reference.IsNameOnly(ref) ***REMOVED***
		layersDownloaded, err = p.pullV2Tag(ctx, ref, os)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		tags, err := p.repo.Tags(ctx).All(ctx)
		if err != nil ***REMOVED***
			// If this repository doesn't exist on V2, we should
			// permit a fallback to V1.
			return allowV1Fallback(err)
		***REMOVED***

		// The v2 registry knows about this repository, so we will not
		// allow fallback to the v1 protocol even if we encounter an
		// error later on.
		p.confirmedV2 = true

		for _, tag := range tags ***REMOVED***
			tagRef, err := reference.WithTag(ref, tag)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			pulledNew, err := p.pullV2Tag(ctx, tagRef, os)
			if err != nil ***REMOVED***
				// Since this is the pull-all-tags case, don't
				// allow an error pulling a particular tag to
				// make the whole pull fall back to v1.
				if fallbackErr, ok := err.(fallbackError); ok ***REMOVED***
					return fallbackErr.err
				***REMOVED***
				return err
			***REMOVED***
			// pulledNew is true if either new layers were downloaded OR if existing images were newly tagged
			// TODO(tiborvass): should we change the name of `layersDownload`? What about message in WriteStatus?
			layersDownloaded = layersDownloaded || pulledNew
		***REMOVED***
	***REMOVED***

	writeStatus(reference.FamiliarString(ref), p.config.ProgressOutput, layersDownloaded)

	return nil
***REMOVED***

type v2LayerDescriptor struct ***REMOVED***
	digest            digest.Digest
	diffID            layer.DiffID
	repoInfo          *registry.RepositoryInfo
	repo              distribution.Repository
	V2MetadataService metadata.V2MetadataService
	tmpFile           *os.File
	verifier          digest.Verifier
	src               distribution.Descriptor
***REMOVED***

func (ld *v2LayerDescriptor) Key() string ***REMOVED***
	return "v2:" + ld.digest.String()
***REMOVED***

func (ld *v2LayerDescriptor) ID() string ***REMOVED***
	return stringid.TruncateID(ld.digest.String())
***REMOVED***

func (ld *v2LayerDescriptor) DiffID() (layer.DiffID, error) ***REMOVED***
	if ld.diffID != "" ***REMOVED***
		return ld.diffID, nil
	***REMOVED***
	return ld.V2MetadataService.GetDiffID(ld.digest)
***REMOVED***

func (ld *v2LayerDescriptor) Download(ctx context.Context, progressOutput progress.Output) (io.ReadCloser, int64, error) ***REMOVED***
	logrus.Debugf("pulling blob %q", ld.digest)

	var (
		err    error
		offset int64
	)

	if ld.tmpFile == nil ***REMOVED***
		ld.tmpFile, err = createDownloadFile()
		if err != nil ***REMOVED***
			return nil, 0, xfer.DoNotRetry***REMOVED***Err: err***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		offset, err = ld.tmpFile.Seek(0, os.SEEK_END)
		if err != nil ***REMOVED***
			logrus.Debugf("error seeking to end of download file: %v", err)
			offset = 0

			ld.tmpFile.Close()
			if err := os.Remove(ld.tmpFile.Name()); err != nil ***REMOVED***
				logrus.Errorf("Failed to remove temp file: %s", ld.tmpFile.Name())
			***REMOVED***
			ld.tmpFile, err = createDownloadFile()
			if err != nil ***REMOVED***
				return nil, 0, xfer.DoNotRetry***REMOVED***Err: err***REMOVED***
			***REMOVED***
		***REMOVED*** else if offset != 0 ***REMOVED***
			logrus.Debugf("attempting to resume download of %q from %d bytes", ld.digest, offset)
		***REMOVED***
	***REMOVED***

	tmpFile := ld.tmpFile

	layerDownload, err := ld.open(ctx)
	if err != nil ***REMOVED***
		logrus.Errorf("Error initiating layer download: %v", err)
		return nil, 0, retryOnError(err)
	***REMOVED***

	if offset != 0 ***REMOVED***
		_, err := layerDownload.Seek(offset, os.SEEK_SET)
		if err != nil ***REMOVED***
			if err := ld.truncateDownloadFile(); err != nil ***REMOVED***
				return nil, 0, xfer.DoNotRetry***REMOVED***Err: err***REMOVED***
			***REMOVED***
			return nil, 0, err
		***REMOVED***
	***REMOVED***
	size, err := layerDownload.Seek(0, os.SEEK_END)
	if err != nil ***REMOVED***
		// Seek failed, perhaps because there was no Content-Length
		// header. This shouldn't fail the download, because we can
		// still continue without a progress bar.
		size = 0
	***REMOVED*** else ***REMOVED***
		if size != 0 && offset > size ***REMOVED***
			logrus.Debug("Partial download is larger than full blob. Starting over")
			offset = 0
			if err := ld.truncateDownloadFile(); err != nil ***REMOVED***
				return nil, 0, xfer.DoNotRetry***REMOVED***Err: err***REMOVED***
			***REMOVED***
		***REMOVED***

		// Restore the seek offset either at the beginning of the
		// stream, or just after the last byte we have from previous
		// attempts.
		_, err = layerDownload.Seek(offset, os.SEEK_SET)
		if err != nil ***REMOVED***
			return nil, 0, err
		***REMOVED***
	***REMOVED***

	reader := progress.NewProgressReader(ioutils.NewCancelReadCloser(ctx, layerDownload), progressOutput, size-offset, ld.ID(), "Downloading")
	defer reader.Close()

	if ld.verifier == nil ***REMOVED***
		ld.verifier = ld.digest.Verifier()
	***REMOVED***

	_, err = io.Copy(tmpFile, io.TeeReader(reader, ld.verifier))
	if err != nil ***REMOVED***
		if err == transport.ErrWrongCodeForByteRange ***REMOVED***
			if err := ld.truncateDownloadFile(); err != nil ***REMOVED***
				return nil, 0, xfer.DoNotRetry***REMOVED***Err: err***REMOVED***
			***REMOVED***
			return nil, 0, err
		***REMOVED***
		return nil, 0, retryOnError(err)
	***REMOVED***

	progress.Update(progressOutput, ld.ID(), "Verifying Checksum")

	if !ld.verifier.Verified() ***REMOVED***
		err = fmt.Errorf("filesystem layer verification failed for digest %s", ld.digest)
		logrus.Error(err)

		// Allow a retry if this digest verification error happened
		// after a resumed download.
		if offset != 0 ***REMOVED***
			if err := ld.truncateDownloadFile(); err != nil ***REMOVED***
				return nil, 0, xfer.DoNotRetry***REMOVED***Err: err***REMOVED***
			***REMOVED***

			return nil, 0, err
		***REMOVED***
		return nil, 0, xfer.DoNotRetry***REMOVED***Err: err***REMOVED***
	***REMOVED***

	progress.Update(progressOutput, ld.ID(), "Download complete")

	logrus.Debugf("Downloaded %s to tempfile %s", ld.ID(), tmpFile.Name())

	_, err = tmpFile.Seek(0, os.SEEK_SET)
	if err != nil ***REMOVED***
		tmpFile.Close()
		if err := os.Remove(tmpFile.Name()); err != nil ***REMOVED***
			logrus.Errorf("Failed to remove temp file: %s", tmpFile.Name())
		***REMOVED***
		ld.tmpFile = nil
		ld.verifier = nil
		return nil, 0, xfer.DoNotRetry***REMOVED***Err: err***REMOVED***
	***REMOVED***

	// hand off the temporary file to the download manager, so it will only
	// be closed once
	ld.tmpFile = nil

	return ioutils.NewReadCloserWrapper(tmpFile, func() error ***REMOVED***
		tmpFile.Close()
		err := os.RemoveAll(tmpFile.Name())
		if err != nil ***REMOVED***
			logrus.Errorf("Failed to remove temp file: %s", tmpFile.Name())
		***REMOVED***
		return err
	***REMOVED***), size, nil
***REMOVED***

func (ld *v2LayerDescriptor) Close() ***REMOVED***
	if ld.tmpFile != nil ***REMOVED***
		ld.tmpFile.Close()
		if err := os.RemoveAll(ld.tmpFile.Name()); err != nil ***REMOVED***
			logrus.Errorf("Failed to remove temp file: %s", ld.tmpFile.Name())
		***REMOVED***
	***REMOVED***
***REMOVED***

func (ld *v2LayerDescriptor) truncateDownloadFile() error ***REMOVED***
	// Need a new hash context since we will be redoing the download
	ld.verifier = nil

	if _, err := ld.tmpFile.Seek(0, os.SEEK_SET); err != nil ***REMOVED***
		logrus.Errorf("error seeking to beginning of download file: %v", err)
		return err
	***REMOVED***

	if err := ld.tmpFile.Truncate(0); err != nil ***REMOVED***
		logrus.Errorf("error truncating download file: %v", err)
		return err
	***REMOVED***

	return nil
***REMOVED***

func (ld *v2LayerDescriptor) Registered(diffID layer.DiffID) ***REMOVED***
	// Cache mapping from this layer's DiffID to the blobsum
	ld.V2MetadataService.Add(diffID, metadata.V2Metadata***REMOVED***Digest: ld.digest, SourceRepository: ld.repoInfo.Name.Name()***REMOVED***)
***REMOVED***

func (p *v2Puller) pullV2Tag(ctx context.Context, ref reference.Named, os string) (tagUpdated bool, err error) ***REMOVED***
	manSvc, err := p.repo.Manifests(ctx)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	var (
		manifest    distribution.Manifest
		tagOrDigest string // Used for logging/progress only
	)
	if digested, isDigested := ref.(reference.Canonical); isDigested ***REMOVED***
		manifest, err = manSvc.Get(ctx, digested.Digest())
		if err != nil ***REMOVED***
			return false, err
		***REMOVED***
		tagOrDigest = digested.Digest().String()
	***REMOVED*** else if tagged, isTagged := ref.(reference.NamedTagged); isTagged ***REMOVED***
		manifest, err = manSvc.Get(ctx, "", distribution.WithTag(tagged.Tag()))
		if err != nil ***REMOVED***
			return false, allowV1Fallback(err)
		***REMOVED***
		tagOrDigest = tagged.Tag()
	***REMOVED*** else ***REMOVED***
		return false, fmt.Errorf("internal error: reference has neither a tag nor a digest: %s", reference.FamiliarString(ref))
	***REMOVED***

	if manifest == nil ***REMOVED***
		return false, fmt.Errorf("image manifest does not exist for tag or digest %q", tagOrDigest)
	***REMOVED***

	if m, ok := manifest.(*schema2.DeserializedManifest); ok ***REMOVED***
		var allowedMediatype bool
		for _, t := range p.config.Schema2Types ***REMOVED***
			if m.Manifest.Config.MediaType == t ***REMOVED***
				allowedMediatype = true
				break
			***REMOVED***
		***REMOVED***
		if !allowedMediatype ***REMOVED***
			configClass := mediaTypeClasses[m.Manifest.Config.MediaType]
			if configClass == "" ***REMOVED***
				configClass = "unknown"
			***REMOVED***
			return false, invalidManifestClassError***REMOVED***m.Manifest.Config.MediaType, configClass***REMOVED***
		***REMOVED***
	***REMOVED***

	// If manSvc.Get succeeded, we can be confident that the registry on
	// the other side speaks the v2 protocol.
	p.confirmedV2 = true

	logrus.Debugf("Pulling ref from V2 registry: %s", reference.FamiliarString(ref))
	progress.Message(p.config.ProgressOutput, tagOrDigest, "Pulling from "+reference.FamiliarName(p.repo.Named()))

	var (
		id             digest.Digest
		manifestDigest digest.Digest
	)

	switch v := manifest.(type) ***REMOVED***
	case *schema1.SignedManifest:
		if p.config.RequireSchema2 ***REMOVED***
			return false, fmt.Errorf("invalid manifest: not schema2")
		***REMOVED***
		id, manifestDigest, err = p.pullSchema1(ctx, ref, v, os)
		if err != nil ***REMOVED***
			return false, err
		***REMOVED***
	case *schema2.DeserializedManifest:
		id, manifestDigest, err = p.pullSchema2(ctx, ref, v, os)
		if err != nil ***REMOVED***
			return false, err
		***REMOVED***
	case *manifestlist.DeserializedManifestList:
		id, manifestDigest, err = p.pullManifestList(ctx, ref, v, os)
		if err != nil ***REMOVED***
			return false, err
		***REMOVED***
	default:
		return false, invalidManifestFormatError***REMOVED******REMOVED***
	***REMOVED***

	progress.Message(p.config.ProgressOutput, "", "Digest: "+manifestDigest.String())

	if p.config.ReferenceStore != nil ***REMOVED***
		oldTagID, err := p.config.ReferenceStore.Get(ref)
		if err == nil ***REMOVED***
			if oldTagID == id ***REMOVED***
				return false, addDigestReference(p.config.ReferenceStore, ref, manifestDigest, id)
			***REMOVED***
		***REMOVED*** else if err != refstore.ErrDoesNotExist ***REMOVED***
			return false, err
		***REMOVED***

		if canonical, ok := ref.(reference.Canonical); ok ***REMOVED***
			if err = p.config.ReferenceStore.AddDigest(canonical, id, true); err != nil ***REMOVED***
				return false, err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if err = addDigestReference(p.config.ReferenceStore, ref, manifestDigest, id); err != nil ***REMOVED***
				return false, err
			***REMOVED***
			if err = p.config.ReferenceStore.AddTag(ref, id, true); err != nil ***REMOVED***
				return false, err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return true, nil
***REMOVED***

func (p *v2Puller) pullSchema1(ctx context.Context, ref reference.Reference, unverifiedManifest *schema1.SignedManifest, requestedOS string) (id digest.Digest, manifestDigest digest.Digest, err error) ***REMOVED***
	var verifiedManifest *schema1.Manifest
	verifiedManifest, err = verifySchema1Manifest(unverifiedManifest, ref)
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	rootFS := image.NewRootFS()

	// remove duplicate layers and check parent chain validity
	err = fixManifestLayers(verifiedManifest)
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	var descriptors []xfer.DownloadDescriptor

	// Image history converted to the new format
	var history []image.History

	// Note that the order of this loop is in the direction of bottom-most
	// to top-most, so that the downloads slice gets ordered correctly.
	for i := len(verifiedManifest.FSLayers) - 1; i >= 0; i-- ***REMOVED***
		blobSum := verifiedManifest.FSLayers[i].BlobSum

		var throwAway struct ***REMOVED***
			ThrowAway bool `json:"throwaway,omitempty"`
		***REMOVED***
		if err := json.Unmarshal([]byte(verifiedManifest.History[i].V1Compatibility), &throwAway); err != nil ***REMOVED***
			return "", "", err
		***REMOVED***

		h, err := v1.HistoryFromConfig([]byte(verifiedManifest.History[i].V1Compatibility), throwAway.ThrowAway)
		if err != nil ***REMOVED***
			return "", "", err
		***REMOVED***
		history = append(history, h)

		if throwAway.ThrowAway ***REMOVED***
			continue
		***REMOVED***

		layerDescriptor := &v2LayerDescriptor***REMOVED***
			digest:            blobSum,
			repoInfo:          p.repoInfo,
			repo:              p.repo,
			V2MetadataService: p.V2MetadataService,
		***REMOVED***

		descriptors = append(descriptors, layerDescriptor)
	***REMOVED***

	// The v1 manifest itself doesn't directly contain an OS. However,
	// the history does, but unfortunately that's a string, so search through
	// all the history until hopefully we find one which indicates the OS.
	// supertest2014/nyan is an example of a registry image with schemav1.
	configOS := runtime.GOOS
	if system.LCOWSupported() ***REMOVED***
		type config struct ***REMOVED***
			Os string `json:"os,omitempty"`
		***REMOVED***
		for _, v := range verifiedManifest.History ***REMOVED***
			var c config
			if err := json.Unmarshal([]byte(v.V1Compatibility), &c); err == nil ***REMOVED***
				if c.Os != "" ***REMOVED***
					configOS = c.Os
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Early bath if the requested OS doesn't match that of the configuration.
	// This avoids doing the download, only to potentially fail later.
	if !strings.EqualFold(configOS, requestedOS) ***REMOVED***
		return "", "", fmt.Errorf("cannot download image with operating system %q when requesting %q", configOS, requestedOS)
	***REMOVED***

	resultRootFS, release, err := p.config.DownloadManager.Download(ctx, *rootFS, configOS, descriptors, p.config.ProgressOutput)
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***
	defer release()

	config, err := v1.MakeConfigFromV1Config([]byte(verifiedManifest.History[0].V1Compatibility), &resultRootFS, history)
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	imageID, err := p.config.ImageStore.Put(config)
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	manifestDigest = digest.FromBytes(unverifiedManifest.Canonical)

	return imageID, manifestDigest, nil
***REMOVED***

func (p *v2Puller) pullSchema2(ctx context.Context, ref reference.Named, mfst *schema2.DeserializedManifest, requestedOS string) (id digest.Digest, manifestDigest digest.Digest, err error) ***REMOVED***
	manifestDigest, err = schema2ManifestDigest(ref, mfst)
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	target := mfst.Target()
	if _, err := p.config.ImageStore.Get(target.Digest); err == nil ***REMOVED***
		// If the image already exists locally, no need to pull
		// anything.
		return target.Digest, manifestDigest, nil
	***REMOVED***

	var descriptors []xfer.DownloadDescriptor

	// Note that the order of this loop is in the direction of bottom-most
	// to top-most, so that the downloads slice gets ordered correctly.
	for _, d := range mfst.Layers ***REMOVED***
		layerDescriptor := &v2LayerDescriptor***REMOVED***
			digest:            d.Digest,
			repo:              p.repo,
			repoInfo:          p.repoInfo,
			V2MetadataService: p.V2MetadataService,
			src:               d,
		***REMOVED***

		descriptors = append(descriptors, layerDescriptor)
	***REMOVED***

	configChan := make(chan []byte, 1)
	configErrChan := make(chan error, 1)
	layerErrChan := make(chan error, 1)
	downloadsDone := make(chan struct***REMOVED******REMOVED***)
	var cancel func()
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	// Pull the image config
	go func() ***REMOVED***
		configJSON, err := p.pullSchema2Config(ctx, target.Digest)
		if err != nil ***REMOVED***
			configErrChan <- ImageConfigPullError***REMOVED***Err: err***REMOVED***
			cancel()
			return
		***REMOVED***
		configChan <- configJSON
	***REMOVED***()

	var (
		configJSON       []byte        // raw serialized image config
		downloadedRootFS *image.RootFS // rootFS from registered layers
		configRootFS     *image.RootFS // rootFS from configuration
		release          func()        // release resources from rootFS download
		configOS         string        // for LCOW when registering downloaded layers
	)

	// https://github.com/docker/docker/issues/24766 - Err on the side of caution,
	// explicitly blocking images intended for linux from the Windows daemon. On
	// Windows, we do this before the attempt to download, effectively serialising
	// the download slightly slowing it down. We have to do it this way, as
	// chances are the download of layers itself would fail due to file names
	// which aren't suitable for NTFS. At some point in the future, if a similar
	// check to block Windows images being pulled on Linux is implemented, it
	// may be necessary to perform the same type of serialisation.
	if runtime.GOOS == "windows" ***REMOVED***
		configJSON, configRootFS, configOS, err = receiveConfig(p.config.ImageStore, configChan, configErrChan)
		if err != nil ***REMOVED***
			return "", "", err
		***REMOVED***

		if configRootFS == nil ***REMOVED***
			return "", "", errRootFSInvalid
		***REMOVED***

		if len(descriptors) != len(configRootFS.DiffIDs) ***REMOVED***
			return "", "", errRootFSMismatch
		***REMOVED***

		// Early bath if the requested OS doesn't match that of the configuration.
		// This avoids doing the download, only to potentially fail later.
		if !strings.EqualFold(configOS, requestedOS) ***REMOVED***
			return "", "", fmt.Errorf("cannot download image with operating system %q when requesting %q", configOS, requestedOS)
		***REMOVED***

		// Populate diff ids in descriptors to avoid downloading foreign layers
		// which have been side loaded
		for i := range descriptors ***REMOVED***
			descriptors[i].(*v2LayerDescriptor).diffID = configRootFS.DiffIDs[i]
		***REMOVED***
	***REMOVED***

	if p.config.DownloadManager != nil ***REMOVED***
		go func() ***REMOVED***
			var (
				err    error
				rootFS image.RootFS
			)
			downloadRootFS := *image.NewRootFS()
			rootFS, release, err = p.config.DownloadManager.Download(ctx, downloadRootFS, requestedOS, descriptors, p.config.ProgressOutput)
			if err != nil ***REMOVED***
				// Intentionally do not cancel the config download here
				// as the error from config download (if there is one)
				// is more interesting than the layer download error
				layerErrChan <- err
				return
			***REMOVED***

			downloadedRootFS = &rootFS
			close(downloadsDone)
		***REMOVED***()
	***REMOVED*** else ***REMOVED***
		// We have nothing to download
		close(downloadsDone)
	***REMOVED***

	if configJSON == nil ***REMOVED***
		configJSON, configRootFS, _, err = receiveConfig(p.config.ImageStore, configChan, configErrChan)
		if err == nil && configRootFS == nil ***REMOVED***
			err = errRootFSInvalid
		***REMOVED***
		if err != nil ***REMOVED***
			cancel()
			select ***REMOVED***
			case <-downloadsDone:
			case <-layerErrChan:
			***REMOVED***
			return "", "", err
		***REMOVED***
	***REMOVED***

	select ***REMOVED***
	case <-downloadsDone:
	case err = <-layerErrChan:
		return "", "", err
	***REMOVED***

	if release != nil ***REMOVED***
		defer release()
	***REMOVED***

	if downloadedRootFS != nil ***REMOVED***
		// The DiffIDs returned in rootFS MUST match those in the config.
		// Otherwise the image config could be referencing layers that aren't
		// included in the manifest.
		if len(downloadedRootFS.DiffIDs) != len(configRootFS.DiffIDs) ***REMOVED***
			return "", "", errRootFSMismatch
		***REMOVED***

		for i := range downloadedRootFS.DiffIDs ***REMOVED***
			if downloadedRootFS.DiffIDs[i] != configRootFS.DiffIDs[i] ***REMOVED***
				return "", "", errRootFSMismatch
			***REMOVED***
		***REMOVED***
	***REMOVED***

	imageID, err := p.config.ImageStore.Put(configJSON)
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	return imageID, manifestDigest, nil
***REMOVED***

func receiveConfig(s ImageConfigStore, configChan <-chan []byte, errChan <-chan error) ([]byte, *image.RootFS, string, error) ***REMOVED***
	select ***REMOVED***
	case configJSON := <-configChan:
		rootfs, os, err := s.RootFSAndOSFromConfig(configJSON)
		if err != nil ***REMOVED***
			return nil, nil, "", err
		***REMOVED***
		return configJSON, rootfs, os, nil
	case err := <-errChan:
		return nil, nil, "", err
		// Don't need a case for ctx.Done in the select because cancellation
		// will trigger an error in p.pullSchema2ImageConfig.
	***REMOVED***
***REMOVED***

// pullManifestList handles "manifest lists" which point to various
// platform-specific manifests.
func (p *v2Puller) pullManifestList(ctx context.Context, ref reference.Named, mfstList *manifestlist.DeserializedManifestList, os string) (id digest.Digest, manifestListDigest digest.Digest, err error) ***REMOVED***
	manifestListDigest, err = schema2ManifestDigest(ref, mfstList)
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	logrus.Debugf("%s resolved to a manifestList object with %d entries; looking for a %s/%s match", ref, len(mfstList.Manifests), os, runtime.GOARCH)

	manifestMatches := filterManifests(mfstList.Manifests, os)

	if len(manifestMatches) == 0 ***REMOVED***
		errMsg := fmt.Sprintf("no matching manifest for %s/%s in the manifest list entries", os, runtime.GOARCH)
		logrus.Debugf(errMsg)
		return "", "", errors.New(errMsg)
	***REMOVED***

	if len(manifestMatches) > 1 ***REMOVED***
		logrus.Debugf("found multiple matches in manifest list, choosing best match %s", manifestMatches[0].Digest.String())
	***REMOVED***
	manifestDigest := manifestMatches[0].Digest

	manSvc, err := p.repo.Manifests(ctx)
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	manifest, err := manSvc.Get(ctx, manifestDigest)
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	manifestRef, err := reference.WithDigest(reference.TrimNamed(ref), manifestDigest)
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	switch v := manifest.(type) ***REMOVED***
	case *schema1.SignedManifest:
		id, _, err = p.pullSchema1(ctx, manifestRef, v, os)
		if err != nil ***REMOVED***
			return "", "", err
		***REMOVED***
	case *schema2.DeserializedManifest:
		id, _, err = p.pullSchema2(ctx, manifestRef, v, os)
		if err != nil ***REMOVED***
			return "", "", err
		***REMOVED***
	default:
		return "", "", errors.New("unsupported manifest format")
	***REMOVED***

	return id, manifestListDigest, err
***REMOVED***

func (p *v2Puller) pullSchema2Config(ctx context.Context, dgst digest.Digest) (configJSON []byte, err error) ***REMOVED***
	blobs := p.repo.Blobs(ctx)
	configJSON, err = blobs.Get(ctx, dgst)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Verify image config digest
	verifier := dgst.Verifier()
	if _, err := verifier.Write(configJSON); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !verifier.Verified() ***REMOVED***
		err := fmt.Errorf("image config verification failed for digest %s", dgst)
		logrus.Error(err)
		return nil, err
	***REMOVED***

	return configJSON, nil
***REMOVED***

// schema2ManifestDigest computes the manifest digest, and, if pulling by
// digest, ensures that it matches the requested digest.
func schema2ManifestDigest(ref reference.Named, mfst distribution.Manifest) (digest.Digest, error) ***REMOVED***
	_, canonical, err := mfst.Payload()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	// If pull by digest, then verify the manifest digest.
	if digested, isDigested := ref.(reference.Canonical); isDigested ***REMOVED***
		verifier := digested.Digest().Verifier()
		if _, err := verifier.Write(canonical); err != nil ***REMOVED***
			return "", err
		***REMOVED***
		if !verifier.Verified() ***REMOVED***
			err := fmt.Errorf("manifest verification failed for digest %s", digested.Digest())
			logrus.Error(err)
			return "", err
		***REMOVED***
		return digested.Digest(), nil
	***REMOVED***

	return digest.FromBytes(canonical), nil
***REMOVED***

// allowV1Fallback checks if the error is a possible reason to fallback to v1
// (even if confirmedV2 has been set already), and if so, wraps the error in
// a fallbackError with confirmedV2 set to false. Otherwise, it returns the
// error unmodified.
func allowV1Fallback(err error) error ***REMOVED***
	switch v := err.(type) ***REMOVED***
	case errcode.Errors:
		if len(v) != 0 ***REMOVED***
			if v0, ok := v[0].(errcode.Error); ok && shouldV2Fallback(v0) ***REMOVED***
				return fallbackError***REMOVED***
					err:         err,
					confirmedV2: false,
					transportOK: true,
				***REMOVED***
			***REMOVED***
		***REMOVED***
	case errcode.Error:
		if shouldV2Fallback(v) ***REMOVED***
			return fallbackError***REMOVED***
				err:         err,
				confirmedV2: false,
				transportOK: true,
			***REMOVED***
		***REMOVED***
	case *url.Error:
		if v.Err == auth.ErrNoBasicAuthCredentials ***REMOVED***
			return fallbackError***REMOVED***err: err, confirmedV2: false***REMOVED***
		***REMOVED***
	***REMOVED***

	return err
***REMOVED***

func verifySchema1Manifest(signedManifest *schema1.SignedManifest, ref reference.Reference) (m *schema1.Manifest, err error) ***REMOVED***
	// If pull by digest, then verify the manifest digest. NOTE: It is
	// important to do this first, before any other content validation. If the
	// digest cannot be verified, don't even bother with those other things.
	if digested, isCanonical := ref.(reference.Canonical); isCanonical ***REMOVED***
		verifier := digested.Digest().Verifier()
		if _, err := verifier.Write(signedManifest.Canonical); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if !verifier.Verified() ***REMOVED***
			err := fmt.Errorf("image verification failed for digest %s", digested.Digest())
			logrus.Error(err)
			return nil, err
		***REMOVED***
	***REMOVED***
	m = &signedManifest.Manifest

	if m.SchemaVersion != 1 ***REMOVED***
		return nil, fmt.Errorf("unsupported schema version %d for %q", m.SchemaVersion, reference.FamiliarString(ref))
	***REMOVED***
	if len(m.FSLayers) != len(m.History) ***REMOVED***
		return nil, fmt.Errorf("length of history not equal to number of layers for %q", reference.FamiliarString(ref))
	***REMOVED***
	if len(m.FSLayers) == 0 ***REMOVED***
		return nil, fmt.Errorf("no FSLayers in manifest for %q", reference.FamiliarString(ref))
	***REMOVED***
	return m, nil
***REMOVED***

// fixManifestLayers removes repeated layers from the manifest and checks the
// correctness of the parent chain.
func fixManifestLayers(m *schema1.Manifest) error ***REMOVED***
	imgs := make([]*image.V1Image, len(m.FSLayers))
	for i := range m.FSLayers ***REMOVED***
		img := &image.V1Image***REMOVED******REMOVED***

		if err := json.Unmarshal([]byte(m.History[i].V1Compatibility), img); err != nil ***REMOVED***
			return err
		***REMOVED***

		imgs[i] = img
		if err := v1.ValidateID(img.ID); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if imgs[len(imgs)-1].Parent != "" && runtime.GOOS != "windows" ***REMOVED***
		// Windows base layer can point to a base layer parent that is not in manifest.
		return errors.New("invalid parent ID in the base layer of the image")
	***REMOVED***

	// check general duplicates to error instead of a deadlock
	idmap := make(map[string]struct***REMOVED******REMOVED***)

	var lastID string
	for _, img := range imgs ***REMOVED***
		// skip IDs that appear after each other, we handle those later
		if _, exists := idmap[img.ID]; img.ID != lastID && exists ***REMOVED***
			return fmt.Errorf("ID %+v appears multiple times in manifest", img.ID)
		***REMOVED***
		lastID = img.ID
		idmap[lastID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	// backwards loop so that we keep the remaining indexes after removing items
	for i := len(imgs) - 2; i >= 0; i-- ***REMOVED***
		if imgs[i].ID == imgs[i+1].ID ***REMOVED*** // repeated ID. remove and continue
			m.FSLayers = append(m.FSLayers[:i], m.FSLayers[i+1:]...)
			m.History = append(m.History[:i], m.History[i+1:]...)
		***REMOVED*** else if imgs[i].Parent != imgs[i+1].ID ***REMOVED***
			return fmt.Errorf("invalid parent ID. Expected %v, got %v", imgs[i+1].ID, imgs[i].Parent)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func createDownloadFile() (*os.File, error) ***REMOVED***
	return ioutil.TempFile("", "GetImageBlob")
***REMOVED***
