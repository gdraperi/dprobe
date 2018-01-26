package distribution

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"sort"
	"strings"
	"sync"

	"golang.org/x/net/context"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/client"
	apitypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/distribution/metadata"
	"github.com/docker/docker/distribution/xfer"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/registry"
	digest "github.com/opencontainers/go-digest"
	"github.com/sirupsen/logrus"
)

const (
	smallLayerMaximumSize  = 100 * (1 << 10) // 100KB
	middleLayerMaximumSize = 10 * (1 << 20)  // 10MB
)

type v2Pusher struct ***REMOVED***
	v2MetadataService metadata.V2MetadataService
	ref               reference.Named
	endpoint          registry.APIEndpoint
	repoInfo          *registry.RepositoryInfo
	config            *ImagePushConfig
	repo              distribution.Repository

	// pushState is state built by the Upload functions.
	pushState pushState
***REMOVED***

type pushState struct ***REMOVED***
	sync.Mutex
	// remoteLayers is the set of layers known to exist on the remote side.
	// This avoids redundant queries when pushing multiple tags that
	// involve the same layers. It is also used to fill in digest and size
	// information when building the manifest.
	remoteLayers map[layer.DiffID]distribution.Descriptor
	// confirmedV2 is set to true if we confirm we're talking to a v2
	// registry. This is used to limit fallbacks to the v1 protocol.
	confirmedV2 bool
***REMOVED***

func (p *v2Pusher) Push(ctx context.Context) (err error) ***REMOVED***
	p.pushState.remoteLayers = make(map[layer.DiffID]distribution.Descriptor)

	p.repo, p.pushState.confirmedV2, err = NewV2Repository(ctx, p.repoInfo, p.endpoint, p.config.MetaHeaders, p.config.AuthConfig, "push", "pull")
	if err != nil ***REMOVED***
		logrus.Debugf("Error getting v2 registry: %v", err)
		return err
	***REMOVED***

	if err = p.pushV2Repository(ctx); err != nil ***REMOVED***
		if continueOnError(err, p.endpoint.Mirror) ***REMOVED***
			return fallbackError***REMOVED***
				err:         err,
				confirmedV2: p.pushState.confirmedV2,
				transportOK: true,
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***

func (p *v2Pusher) pushV2Repository(ctx context.Context) (err error) ***REMOVED***
	if namedTagged, isNamedTagged := p.ref.(reference.NamedTagged); isNamedTagged ***REMOVED***
		imageID, err := p.config.ReferenceStore.Get(p.ref)
		if err != nil ***REMOVED***
			return fmt.Errorf("tag does not exist: %s", reference.FamiliarString(p.ref))
		***REMOVED***

		return p.pushV2Tag(ctx, namedTagged, imageID)
	***REMOVED***

	if !reference.IsNameOnly(p.ref) ***REMOVED***
		return errors.New("cannot push a digest reference")
	***REMOVED***

	// Pull all tags
	pushed := 0
	for _, association := range p.config.ReferenceStore.ReferencesByName(p.ref) ***REMOVED***
		if namedTagged, isNamedTagged := association.Ref.(reference.NamedTagged); isNamedTagged ***REMOVED***
			pushed++
			if err := p.pushV2Tag(ctx, namedTagged, association.ID); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if pushed == 0 ***REMOVED***
		return fmt.Errorf("no tags to push for %s", reference.FamiliarName(p.repoInfo.Name))
	***REMOVED***

	return nil
***REMOVED***

func (p *v2Pusher) pushV2Tag(ctx context.Context, ref reference.NamedTagged, id digest.Digest) error ***REMOVED***
	logrus.Debugf("Pushing repository: %s", reference.FamiliarString(ref))

	imgConfig, err := p.config.ImageStore.Get(id)
	if err != nil ***REMOVED***
		return fmt.Errorf("could not find image from tag %s: %v", reference.FamiliarString(ref), err)
	***REMOVED***

	rootfs, os, err := p.config.ImageStore.RootFSAndOSFromConfig(imgConfig)
	if err != nil ***REMOVED***
		return fmt.Errorf("unable to get rootfs for image %s: %s", reference.FamiliarString(ref), err)
	***REMOVED***

	l, err := p.config.LayerStores[os].Get(rootfs.ChainID())
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to get top layer from image: %v", err)
	***REMOVED***
	defer l.Release()

	hmacKey, err := metadata.ComputeV2MetadataHMACKey(p.config.AuthConfig)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to compute hmac key of auth config: %v", err)
	***REMOVED***

	var descriptors []xfer.UploadDescriptor

	descriptorTemplate := v2PushDescriptor***REMOVED***
		v2MetadataService: p.v2MetadataService,
		hmacKey:           hmacKey,
		repoInfo:          p.repoInfo.Name,
		ref:               p.ref,
		endpoint:          p.endpoint,
		repo:              p.repo,
		pushState:         &p.pushState,
	***REMOVED***

	// Loop bounds condition is to avoid pushing the base layer on Windows.
	for range rootfs.DiffIDs ***REMOVED***
		descriptor := descriptorTemplate
		descriptor.layer = l
		descriptor.checkedDigests = make(map[digest.Digest]struct***REMOVED******REMOVED***)
		descriptors = append(descriptors, &descriptor)

		l = l.Parent()
	***REMOVED***

	if err := p.config.UploadManager.Upload(ctx, descriptors, p.config.ProgressOutput); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Try schema2 first
	builder := schema2.NewManifestBuilder(p.repo.Blobs(ctx), p.config.ConfigMediaType, imgConfig)
	manifest, err := manifestFromBuilder(ctx, builder, descriptors)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	manSvc, err := p.repo.Manifests(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	putOptions := []distribution.ManifestServiceOption***REMOVED***distribution.WithTag(ref.Tag())***REMOVED***
	if _, err = manSvc.Put(ctx, manifest, putOptions...); err != nil ***REMOVED***
		if runtime.GOOS == "windows" || p.config.TrustKey == nil || p.config.RequireSchema2 ***REMOVED***
			logrus.Warnf("failed to upload schema2 manifest: %v", err)
			return err
		***REMOVED***

		logrus.Warnf("failed to upload schema2 manifest: %v - falling back to schema1", err)

		manifestRef, err := reference.WithTag(p.repo.Named(), ref.Tag())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		builder = schema1.NewConfigManifestBuilder(p.repo.Blobs(ctx), p.config.TrustKey, manifestRef, imgConfig)
		manifest, err = manifestFromBuilder(ctx, builder, descriptors)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if _, err = manSvc.Put(ctx, manifest, putOptions...); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	var canonicalManifest []byte

	switch v := manifest.(type) ***REMOVED***
	case *schema1.SignedManifest:
		canonicalManifest = v.Canonical
	case *schema2.DeserializedManifest:
		_, canonicalManifest, err = v.Payload()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	manifestDigest := digest.FromBytes(canonicalManifest)
	progress.Messagef(p.config.ProgressOutput, "", "%s: digest: %s size: %d", ref.Tag(), manifestDigest, len(canonicalManifest))

	if err := addDigestReference(p.config.ReferenceStore, ref, manifestDigest, id); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Signal digest to the trust client so it can sign the
	// push, if appropriate.
	progress.Aux(p.config.ProgressOutput, apitypes.PushResult***REMOVED***Tag: ref.Tag(), Digest: manifestDigest.String(), Size: len(canonicalManifest)***REMOVED***)

	return nil
***REMOVED***

func manifestFromBuilder(ctx context.Context, builder distribution.ManifestBuilder, descriptors []xfer.UploadDescriptor) (distribution.Manifest, error) ***REMOVED***
	// descriptors is in reverse order; iterate backwards to get references
	// appended in the right order.
	for i := len(descriptors) - 1; i >= 0; i-- ***REMOVED***
		if err := builder.AppendReference(descriptors[i].(*v2PushDescriptor)); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return builder.Build(ctx)
***REMOVED***

type v2PushDescriptor struct ***REMOVED***
	layer             PushLayer
	v2MetadataService metadata.V2MetadataService
	hmacKey           []byte
	repoInfo          reference.Named
	ref               reference.Named
	endpoint          registry.APIEndpoint
	repo              distribution.Repository
	pushState         *pushState
	remoteDescriptor  distribution.Descriptor
	// a set of digests whose presence has been checked in a target repository
	checkedDigests map[digest.Digest]struct***REMOVED******REMOVED***
***REMOVED***

func (pd *v2PushDescriptor) Key() string ***REMOVED***
	return "v2push:" + pd.ref.Name() + " " + pd.layer.DiffID().String()
***REMOVED***

func (pd *v2PushDescriptor) ID() string ***REMOVED***
	return stringid.TruncateID(pd.layer.DiffID().String())
***REMOVED***

func (pd *v2PushDescriptor) DiffID() layer.DiffID ***REMOVED***
	return pd.layer.DiffID()
***REMOVED***

func (pd *v2PushDescriptor) Upload(ctx context.Context, progressOutput progress.Output) (distribution.Descriptor, error) ***REMOVED***
	// Skip foreign layers unless this registry allows nondistributable artifacts.
	if !pd.endpoint.AllowNondistributableArtifacts ***REMOVED***
		if fs, ok := pd.layer.(distribution.Describable); ok ***REMOVED***
			if d := fs.Descriptor(); len(d.URLs) > 0 ***REMOVED***
				progress.Update(progressOutput, pd.ID(), "Skipped foreign layer")
				return d, nil
			***REMOVED***
		***REMOVED***
	***REMOVED***

	diffID := pd.DiffID()

	pd.pushState.Lock()
	if descriptor, ok := pd.pushState.remoteLayers[diffID]; ok ***REMOVED***
		// it is already known that the push is not needed and
		// therefore doing a stat is unnecessary
		pd.pushState.Unlock()
		progress.Update(progressOutput, pd.ID(), "Layer already exists")
		return descriptor, nil
	***REMOVED***
	pd.pushState.Unlock()

	maxMountAttempts, maxExistenceChecks, checkOtherRepositories := getMaxMountAndExistenceCheckAttempts(pd.layer)

	// Do we have any metadata associated with this layer's DiffID?
	v2Metadata, err := pd.v2MetadataService.GetMetadata(diffID)
	if err == nil ***REMOVED***
		// check for blob existence in the target repository
		descriptor, exists, err := pd.layerAlreadyExists(ctx, progressOutput, diffID, true, 1, v2Metadata)
		if exists || err != nil ***REMOVED***
			return descriptor, err
		***REMOVED***
	***REMOVED***

	// if digest was empty or not saved, or if blob does not exist on the remote repository,
	// then push the blob.
	bs := pd.repo.Blobs(ctx)

	var layerUpload distribution.BlobWriter

	// Attempt to find another repository in the same registry to mount the layer from to avoid an unnecessary upload
	candidates := getRepositoryMountCandidates(pd.repoInfo, pd.hmacKey, maxMountAttempts, v2Metadata)
	for _, mountCandidate := range candidates ***REMOVED***
		logrus.Debugf("attempting to mount layer %s (%s) from %s", diffID, mountCandidate.Digest, mountCandidate.SourceRepository)
		createOpts := []distribution.BlobCreateOption***REMOVED******REMOVED***

		if len(mountCandidate.SourceRepository) > 0 ***REMOVED***
			namedRef, err := reference.ParseNormalizedNamed(mountCandidate.SourceRepository)
			if err != nil ***REMOVED***
				logrus.Errorf("failed to parse source repository reference %v: %v", reference.FamiliarString(namedRef), err)
				pd.v2MetadataService.Remove(mountCandidate)
				continue
			***REMOVED***

			// Candidates are always under same domain, create remote reference
			// with only path to set mount from with
			remoteRef, err := reference.WithName(reference.Path(namedRef))
			if err != nil ***REMOVED***
				logrus.Errorf("failed to make remote reference out of %q: %v", reference.Path(namedRef), err)
				continue
			***REMOVED***

			canonicalRef, err := reference.WithDigest(reference.TrimNamed(remoteRef), mountCandidate.Digest)
			if err != nil ***REMOVED***
				logrus.Errorf("failed to make canonical reference: %v", err)
				continue
			***REMOVED***

			createOpts = append(createOpts, client.WithMountFrom(canonicalRef))
		***REMOVED***

		// send the layer
		lu, err := bs.Create(ctx, createOpts...)
		switch err := err.(type) ***REMOVED***
		case nil:
			// noop
		case distribution.ErrBlobMounted:
			progress.Updatef(progressOutput, pd.ID(), "Mounted from %s", err.From.Name())

			err.Descriptor.MediaType = schema2.MediaTypeLayer

			pd.pushState.Lock()
			pd.pushState.confirmedV2 = true
			pd.pushState.remoteLayers[diffID] = err.Descriptor
			pd.pushState.Unlock()

			// Cache mapping from this layer's DiffID to the blobsum
			if err := pd.v2MetadataService.TagAndAdd(diffID, pd.hmacKey, metadata.V2Metadata***REMOVED***
				Digest:           err.Descriptor.Digest,
				SourceRepository: pd.repoInfo.Name(),
			***REMOVED***); err != nil ***REMOVED***
				return distribution.Descriptor***REMOVED******REMOVED***, xfer.DoNotRetry***REMOVED***Err: err***REMOVED***
			***REMOVED***
			return err.Descriptor, nil
		default:
			logrus.Infof("failed to mount layer %s (%s) from %s: %v", diffID, mountCandidate.Digest, mountCandidate.SourceRepository, err)
		***REMOVED***

		if len(mountCandidate.SourceRepository) > 0 &&
			(metadata.CheckV2MetadataHMAC(&mountCandidate, pd.hmacKey) ||
				len(mountCandidate.HMAC) == 0) ***REMOVED***
			cause := "blob mount failure"
			if err != nil ***REMOVED***
				cause = fmt.Sprintf("an error: %v", err.Error())
			***REMOVED***
			logrus.Debugf("removing association between layer %s and %s due to %s", mountCandidate.Digest, mountCandidate.SourceRepository, cause)
			pd.v2MetadataService.Remove(mountCandidate)
		***REMOVED***

		if lu != nil ***REMOVED***
			// cancel previous upload
			cancelLayerUpload(ctx, mountCandidate.Digest, layerUpload)
			layerUpload = lu
		***REMOVED***
	***REMOVED***

	if maxExistenceChecks-len(pd.checkedDigests) > 0 ***REMOVED***
		// do additional layer existence checks with other known digests if any
		descriptor, exists, err := pd.layerAlreadyExists(ctx, progressOutput, diffID, checkOtherRepositories, maxExistenceChecks-len(pd.checkedDigests), v2Metadata)
		if exists || err != nil ***REMOVED***
			return descriptor, err
		***REMOVED***
	***REMOVED***

	logrus.Debugf("Pushing layer: %s", diffID)
	if layerUpload == nil ***REMOVED***
		layerUpload, err = bs.Create(ctx)
		if err != nil ***REMOVED***
			return distribution.Descriptor***REMOVED******REMOVED***, retryOnError(err)
		***REMOVED***
	***REMOVED***
	defer layerUpload.Close()

	// upload the blob
	return pd.uploadUsingSession(ctx, progressOutput, diffID, layerUpload)
***REMOVED***

func (pd *v2PushDescriptor) SetRemoteDescriptor(descriptor distribution.Descriptor) ***REMOVED***
	pd.remoteDescriptor = descriptor
***REMOVED***

func (pd *v2PushDescriptor) Descriptor() distribution.Descriptor ***REMOVED***
	return pd.remoteDescriptor
***REMOVED***

func (pd *v2PushDescriptor) uploadUsingSession(
	ctx context.Context,
	progressOutput progress.Output,
	diffID layer.DiffID,
	layerUpload distribution.BlobWriter,
) (distribution.Descriptor, error) ***REMOVED***
	var reader io.ReadCloser

	contentReader, err := pd.layer.Open()
	if err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, retryOnError(err)
	***REMOVED***

	size, _ := pd.layer.Size()

	reader = progress.NewProgressReader(ioutils.NewCancelReadCloser(ctx, contentReader), progressOutput, size, pd.ID(), "Pushing")

	switch m := pd.layer.MediaType(); m ***REMOVED***
	case schema2.MediaTypeUncompressedLayer:
		compressedReader, compressionDone := compress(reader)
		defer func(closer io.Closer) ***REMOVED***
			closer.Close()
			<-compressionDone
		***REMOVED***(reader)
		reader = compressedReader
	case schema2.MediaTypeLayer:
	default:
		reader.Close()
		return distribution.Descriptor***REMOVED******REMOVED***, fmt.Errorf("unsupported layer media type %s", m)
	***REMOVED***

	digester := digest.Canonical.Digester()
	tee := io.TeeReader(reader, digester.Hash())

	nn, err := layerUpload.ReadFrom(tee)
	reader.Close()
	if err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, retryOnError(err)
	***REMOVED***

	pushDigest := digester.Digest()
	if _, err := layerUpload.Commit(ctx, distribution.Descriptor***REMOVED***Digest: pushDigest***REMOVED***); err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, retryOnError(err)
	***REMOVED***

	logrus.Debugf("uploaded layer %s (%s), %d bytes", diffID, pushDigest, nn)
	progress.Update(progressOutput, pd.ID(), "Pushed")

	// Cache mapping from this layer's DiffID to the blobsum
	if err := pd.v2MetadataService.TagAndAdd(diffID, pd.hmacKey, metadata.V2Metadata***REMOVED***
		Digest:           pushDigest,
		SourceRepository: pd.repoInfo.Name(),
	***REMOVED***); err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, xfer.DoNotRetry***REMOVED***Err: err***REMOVED***
	***REMOVED***

	desc := distribution.Descriptor***REMOVED***
		Digest:    pushDigest,
		MediaType: schema2.MediaTypeLayer,
		Size:      nn,
	***REMOVED***

	pd.pushState.Lock()
	// If Commit succeeded, that's an indication that the remote registry speaks the v2 protocol.
	pd.pushState.confirmedV2 = true
	pd.pushState.remoteLayers[diffID] = desc
	pd.pushState.Unlock()

	return desc, nil
***REMOVED***

// layerAlreadyExists checks if the registry already knows about any of the metadata passed in the "metadata"
// slice. If it finds one that the registry knows about, it returns the known digest and "true". If
// "checkOtherRepositories" is true, stat will be performed also with digests mapped to any other repository
// (not just the target one).
func (pd *v2PushDescriptor) layerAlreadyExists(
	ctx context.Context,
	progressOutput progress.Output,
	diffID layer.DiffID,
	checkOtherRepositories bool,
	maxExistenceCheckAttempts int,
	v2Metadata []metadata.V2Metadata,
) (desc distribution.Descriptor, exists bool, err error) ***REMOVED***
	// filter the metadata
	candidates := []metadata.V2Metadata***REMOVED******REMOVED***
	for _, meta := range v2Metadata ***REMOVED***
		if len(meta.SourceRepository) > 0 && !checkOtherRepositories && meta.SourceRepository != pd.repoInfo.Name() ***REMOVED***
			continue
		***REMOVED***
		candidates = append(candidates, meta)
	***REMOVED***
	// sort the candidates by similarity
	sortV2MetadataByLikenessAndAge(pd.repoInfo, pd.hmacKey, candidates)

	digestToMetadata := make(map[digest.Digest]*metadata.V2Metadata)
	// an array of unique blob digests ordered from the best mount candidates to worst
	layerDigests := []digest.Digest***REMOVED******REMOVED***
	for i := 0; i < len(candidates); i++ ***REMOVED***
		if len(layerDigests) >= maxExistenceCheckAttempts ***REMOVED***
			break
		***REMOVED***
		meta := &candidates[i]
		if _, exists := digestToMetadata[meta.Digest]; exists ***REMOVED***
			// keep reference just to the first mapping (the best mount candidate)
			continue
		***REMOVED***
		if _, exists := pd.checkedDigests[meta.Digest]; exists ***REMOVED***
			// existence of this digest has already been tested
			continue
		***REMOVED***
		digestToMetadata[meta.Digest] = meta
		layerDigests = append(layerDigests, meta.Digest)
	***REMOVED***

attempts:
	for _, dgst := range layerDigests ***REMOVED***
		meta := digestToMetadata[dgst]
		logrus.Debugf("Checking for presence of layer %s (%s) in %s", diffID, dgst, pd.repoInfo.Name())
		desc, err = pd.repo.Blobs(ctx).Stat(ctx, dgst)
		pd.checkedDigests[meta.Digest] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		switch err ***REMOVED***
		case nil:
			if m, ok := digestToMetadata[desc.Digest]; !ok || m.SourceRepository != pd.repoInfo.Name() || !metadata.CheckV2MetadataHMAC(m, pd.hmacKey) ***REMOVED***
				// cache mapping from this layer's DiffID to the blobsum
				if err := pd.v2MetadataService.TagAndAdd(diffID, pd.hmacKey, metadata.V2Metadata***REMOVED***
					Digest:           desc.Digest,
					SourceRepository: pd.repoInfo.Name(),
				***REMOVED***); err != nil ***REMOVED***
					return distribution.Descriptor***REMOVED******REMOVED***, false, xfer.DoNotRetry***REMOVED***Err: err***REMOVED***
				***REMOVED***
			***REMOVED***
			desc.MediaType = schema2.MediaTypeLayer
			exists = true
			break attempts
		case distribution.ErrBlobUnknown:
			if meta.SourceRepository == pd.repoInfo.Name() ***REMOVED***
				// remove the mapping to the target repository
				pd.v2MetadataService.Remove(*meta)
			***REMOVED***
		default:
			logrus.WithError(err).Debugf("Failed to check for presence of layer %s (%s) in %s", diffID, dgst, pd.repoInfo.Name())
		***REMOVED***
	***REMOVED***

	if exists ***REMOVED***
		progress.Update(progressOutput, pd.ID(), "Layer already exists")
		pd.pushState.Lock()
		pd.pushState.remoteLayers[diffID] = desc
		pd.pushState.Unlock()
	***REMOVED***

	return desc, exists, nil
***REMOVED***

// getMaxMountAndExistenceCheckAttempts returns a maximum number of cross repository mount attempts from
// source repositories of target registry, maximum number of layer existence checks performed on the target
// repository and whether the check shall be done also with digests mapped to different repositories. The
// decision is based on layer size. The smaller the layer, the fewer attempts shall be made because the cost
// of upload does not outweigh a latency.
func getMaxMountAndExistenceCheckAttempts(layer PushLayer) (maxMountAttempts, maxExistenceCheckAttempts int, checkOtherRepositories bool) ***REMOVED***
	size, err := layer.Size()
	switch ***REMOVED***
	// big blob
	case size > middleLayerMaximumSize:
		// 1st attempt to mount the blob few times
		// 2nd few existence checks with digests associated to any repository
		// then fallback to upload
		return 4, 3, true

	// middle sized blobs; if we could not get the size, assume we deal with middle sized blob
	case size > smallLayerMaximumSize, err != nil:
		// 1st attempt to mount blobs of average size few times
		// 2nd try at most 1 existence check if there's an existing mapping to the target repository
		// then fallback to upload
		return 3, 1, false

	// small blobs, do a minimum number of checks
	default:
		return 1, 1, false
	***REMOVED***
***REMOVED***

// getRepositoryMountCandidates returns an array of v2 metadata items belonging to the given registry. The
// array is sorted from youngest to oldest. If requireRegistryMatch is true, the resulting array will contain
// only metadata entries having registry part of SourceRepository matching the part of repoInfo.
func getRepositoryMountCandidates(
	repoInfo reference.Named,
	hmacKey []byte,
	max int,
	v2Metadata []metadata.V2Metadata,
) []metadata.V2Metadata ***REMOVED***
	candidates := []metadata.V2Metadata***REMOVED******REMOVED***
	for _, meta := range v2Metadata ***REMOVED***
		sourceRepo, err := reference.ParseNamed(meta.SourceRepository)
		if err != nil || reference.Domain(repoInfo) != reference.Domain(sourceRepo) ***REMOVED***
			continue
		***REMOVED***
		// target repository is not a viable candidate
		if meta.SourceRepository == repoInfo.Name() ***REMOVED***
			continue
		***REMOVED***
		candidates = append(candidates, meta)
	***REMOVED***

	sortV2MetadataByLikenessAndAge(repoInfo, hmacKey, candidates)
	if max >= 0 && len(candidates) > max ***REMOVED***
		// select the youngest metadata
		candidates = candidates[:max]
	***REMOVED***

	return candidates
***REMOVED***

// byLikeness is a sorting container for v2 metadata candidates for cross repository mount. The
// candidate "a" is preferred over "b":
//
//  1. if it was hashed using the same AuthConfig as the one used to authenticate to target repository and the
//     "b" was not
//  2. if a number of its repository path components exactly matching path components of target repository is higher
type byLikeness struct ***REMOVED***
	arr            []metadata.V2Metadata
	hmacKey        []byte
	pathComponents []string
***REMOVED***

func (bla byLikeness) Less(i, j int) bool ***REMOVED***
	aMacMatch := metadata.CheckV2MetadataHMAC(&bla.arr[i], bla.hmacKey)
	bMacMatch := metadata.CheckV2MetadataHMAC(&bla.arr[j], bla.hmacKey)
	if aMacMatch != bMacMatch ***REMOVED***
		return aMacMatch
	***REMOVED***
	aMatch := numOfMatchingPathComponents(bla.arr[i].SourceRepository, bla.pathComponents)
	bMatch := numOfMatchingPathComponents(bla.arr[j].SourceRepository, bla.pathComponents)
	return aMatch > bMatch
***REMOVED***
func (bla byLikeness) Swap(i, j int) ***REMOVED***
	bla.arr[i], bla.arr[j] = bla.arr[j], bla.arr[i]
***REMOVED***
func (bla byLikeness) Len() int ***REMOVED*** return len(bla.arr) ***REMOVED***

// nolint: interfacer
func sortV2MetadataByLikenessAndAge(repoInfo reference.Named, hmacKey []byte, marr []metadata.V2Metadata) ***REMOVED***
	// reverse the metadata array to shift the newest entries to the beginning
	for i := 0; i < len(marr)/2; i++ ***REMOVED***
		marr[i], marr[len(marr)-i-1] = marr[len(marr)-i-1], marr[i]
	***REMOVED***
	// keep equal entries ordered from the youngest to the oldest
	sort.Stable(byLikeness***REMOVED***
		arr:            marr,
		hmacKey:        hmacKey,
		pathComponents: getPathComponents(repoInfo.Name()),
	***REMOVED***)
***REMOVED***

// numOfMatchingPathComponents returns a number of path components in "pth" that exactly match "matchComponents".
func numOfMatchingPathComponents(pth string, matchComponents []string) int ***REMOVED***
	pthComponents := getPathComponents(pth)
	i := 0
	for ; i < len(pthComponents) && i < len(matchComponents); i++ ***REMOVED***
		if matchComponents[i] != pthComponents[i] ***REMOVED***
			return i
		***REMOVED***
	***REMOVED***
	return i
***REMOVED***

func getPathComponents(path string) []string ***REMOVED***
	return strings.Split(path, "/")
***REMOVED***

func cancelLayerUpload(ctx context.Context, dgst digest.Digest, layerUpload distribution.BlobWriter) ***REMOVED***
	if layerUpload != nil ***REMOVED***
		logrus.Debugf("cancelling upload of blob %s", dgst)
		err := layerUpload.Cancel(ctx)
		if err != nil ***REMOVED***
			logrus.Warnf("failed to cancel upload: %v", err)
		***REMOVED***
	***REMOVED***
***REMOVED***
