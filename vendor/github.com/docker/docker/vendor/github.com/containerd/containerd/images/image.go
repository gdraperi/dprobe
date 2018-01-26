package images

import (
	"context"
	"encoding/json"
	"time"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/platforms"
	digest "github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

// Image provides the model for how containerd views container images.
type Image struct ***REMOVED***
	// Name of the image.
	//
	// To be pulled, it must be a reference compatible with resolvers.
	//
	// This field is required.
	Name string

	// Labels provide runtime decoration for the image record.
	//
	// There is no default behavior for how these labels are propagated. They
	// only decorate the static metadata object.
	//
	// This field is optional.
	Labels map[string]string

	// Target describes the root content for this image. Typically, this is
	// a manifest, index or manifest list.
	Target ocispec.Descriptor

	CreatedAt, UpdatedAt time.Time
***REMOVED***

// DeleteOptions provide options on image delete
type DeleteOptions struct ***REMOVED***
	Synchronous bool
***REMOVED***

// DeleteOpt allows configuring a delete operation
type DeleteOpt func(context.Context, *DeleteOptions) error

// SynchronousDelete is used to indicate that an image deletion and removal of
// the image resources should occur synchronously before returning a result.
func SynchronousDelete() DeleteOpt ***REMOVED***
	return func(ctx context.Context, o *DeleteOptions) error ***REMOVED***
		o.Synchronous = true
		return nil
	***REMOVED***
***REMOVED***

// Store and interact with images
type Store interface ***REMOVED***
	Get(ctx context.Context, name string) (Image, error)
	List(ctx context.Context, filters ...string) ([]Image, error)
	Create(ctx context.Context, image Image) (Image, error)

	// Update will replace the data in the store with the provided image. If
	// one or more fieldpaths are provided, only those fields will be updated.
	Update(ctx context.Context, image Image, fieldpaths ...string) (Image, error)

	Delete(ctx context.Context, name string, opts ...DeleteOpt) error
***REMOVED***

// TODO(stevvooe): Many of these functions make strong platform assumptions,
// which are untrue in a lot of cases. More refactoring must be done here to
// make this work in all cases.

// Config resolves the image configuration descriptor.
//
// The caller can then use the descriptor to resolve and process the
// configuration of the image.
func (image *Image) Config(ctx context.Context, provider content.Provider, platform string) (ocispec.Descriptor, error) ***REMOVED***
	return Config(ctx, provider, image.Target, platform)
***REMOVED***

// RootFS returns the unpacked diffids that make up and images rootfs.
//
// These are used to verify that a set of layers unpacked to the expected
// values.
func (image *Image) RootFS(ctx context.Context, provider content.Provider, platform string) ([]digest.Digest, error) ***REMOVED***
	desc, err := image.Config(ctx, provider, platform)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return RootFS(ctx, provider, desc)
***REMOVED***

// Size returns the total size of an image's packed resources.
func (image *Image) Size(ctx context.Context, provider content.Provider, platform string) (int64, error) ***REMOVED***
	var size int64
	return size, Walk(ctx, Handlers(HandlerFunc(func(ctx context.Context, desc ocispec.Descriptor) ([]ocispec.Descriptor, error) ***REMOVED***
		if desc.Size < 0 ***REMOVED***
			return nil, errors.Errorf("invalid size %v in %v (%v)", desc.Size, desc.Digest, desc.MediaType)
		***REMOVED***
		size += desc.Size
		return nil, nil
	***REMOVED***), ChildrenHandler(provider, platform)), image.Target)
***REMOVED***

// Manifest resolves a manifest from the image for the given platform.
//
// TODO(stevvooe): This violates the current platform agnostic approach to this
// package by returning a specific manifest type. We'll need to refactor this
// to return a manifest descriptor or decide that we want to bring the API in
// this direction because this abstraction is not needed.`
func Manifest(ctx context.Context, provider content.Provider, image ocispec.Descriptor, platform string) (ocispec.Manifest, error) ***REMOVED***
	var (
		matcher platforms.Matcher
		m       *ocispec.Manifest
		err     error
	)
	if platform != "" ***REMOVED***
		matcher, err = platforms.Parse(platform)
		if err != nil ***REMOVED***
			return ocispec.Manifest***REMOVED******REMOVED***, err
		***REMOVED***
	***REMOVED***

	if err := Walk(ctx, HandlerFunc(func(ctx context.Context, desc ocispec.Descriptor) ([]ocispec.Descriptor, error) ***REMOVED***
		switch desc.MediaType ***REMOVED***
		case MediaTypeDockerSchema2Manifest, ocispec.MediaTypeImageManifest:
			p, err := content.ReadBlob(ctx, provider, desc.Digest)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			var manifest ocispec.Manifest
			if err := json.Unmarshal(p, &manifest); err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			if platform != "" ***REMOVED***
				if desc.Platform != nil && !matcher.Match(*desc.Platform) ***REMOVED***
					return nil, nil
				***REMOVED***

				if desc.Platform == nil ***REMOVED***
					p, err := content.ReadBlob(ctx, provider, manifest.Config.Digest)
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***

					var image ocispec.Image
					if err := json.Unmarshal(p, &image); err != nil ***REMOVED***
						return nil, err
					***REMOVED***

					if !matcher.Match(platforms.Normalize(ocispec.Platform***REMOVED***OS: image.OS, Architecture: image.Architecture***REMOVED***)) ***REMOVED***
						return nil, nil
					***REMOVED***

				***REMOVED***
			***REMOVED***

			m = &manifest

			return nil, nil
		case MediaTypeDockerSchema2ManifestList, ocispec.MediaTypeImageIndex:
			p, err := content.ReadBlob(ctx, provider, desc.Digest)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			var idx ocispec.Index
			if err := json.Unmarshal(p, &idx); err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			if platform == "" ***REMOVED***
				return idx.Manifests, nil
			***REMOVED***

			var descs []ocispec.Descriptor
			for _, d := range idx.Manifests ***REMOVED***
				if d.Platform == nil || matcher.Match(*d.Platform) ***REMOVED***
					descs = append(descs, d)
				***REMOVED***
			***REMOVED***

			return descs, nil

		***REMOVED***
		return nil, errors.Wrapf(errdefs.ErrNotFound, "unexpected media type %v for %v", desc.MediaType, desc.Digest)
	***REMOVED***), image); err != nil ***REMOVED***
		return ocispec.Manifest***REMOVED******REMOVED***, err
	***REMOVED***

	if m == nil ***REMOVED***
		return ocispec.Manifest***REMOVED******REMOVED***, errors.Wrapf(errdefs.ErrNotFound, "manifest %v", image.Digest)
	***REMOVED***

	return *m, nil
***REMOVED***

// Config resolves the image configuration descriptor using a content provided
// to resolve child resources on the image.
//
// The caller can then use the descriptor to resolve and process the
// configuration of the image.
func Config(ctx context.Context, provider content.Provider, image ocispec.Descriptor, platform string) (ocispec.Descriptor, error) ***REMOVED***
	manifest, err := Manifest(ctx, provider, image, platform)
	if err != nil ***REMOVED***
		return ocispec.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***
	return manifest.Config, err
***REMOVED***

// Platforms returns one or more platforms supported by the image.
func Platforms(ctx context.Context, provider content.Provider, image ocispec.Descriptor) ([]ocispec.Platform, error) ***REMOVED***
	var platformSpecs []ocispec.Platform
	return platformSpecs, Walk(ctx, Handlers(HandlerFunc(func(ctx context.Context, desc ocispec.Descriptor) ([]ocispec.Descriptor, error) ***REMOVED***
		if desc.Platform != nil ***REMOVED***
			platformSpecs = append(platformSpecs, *desc.Platform)
			return nil, ErrSkipDesc
		***REMOVED***

		switch desc.MediaType ***REMOVED***
		case MediaTypeDockerSchema2Config, ocispec.MediaTypeImageConfig:
			p, err := content.ReadBlob(ctx, provider, desc.Digest)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			var image ocispec.Image
			if err := json.Unmarshal(p, &image); err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			platformSpecs = append(platformSpecs,
				platforms.Normalize(ocispec.Platform***REMOVED***OS: image.OS, Architecture: image.Architecture***REMOVED***))
		***REMOVED***
		return nil, nil
	***REMOVED***), ChildrenHandler(provider, "")), image)
***REMOVED***

// Check returns nil if the all components of an image are available in the
// provider for the specified platform.
//
// If available is true, the caller can assume that required represents the
// complete set of content required for the image.
//
// missing will have the components that are part of required but not avaiiable
// in the provider.
//
// If there is a problem resolving content, an error will be returned.
func Check(ctx context.Context, provider content.Provider, image ocispec.Descriptor, platform string) (available bool, required, present, missing []ocispec.Descriptor, err error) ***REMOVED***
	mfst, err := Manifest(ctx, provider, image, platform)
	if err != nil ***REMOVED***
		if errdefs.IsNotFound(err) ***REMOVED***
			return false, []ocispec.Descriptor***REMOVED***image***REMOVED***, nil, []ocispec.Descriptor***REMOVED***image***REMOVED***, nil
		***REMOVED***

		return false, nil, nil, nil, errors.Wrapf(err, "failed to check image %v", image.Digest)
	***REMOVED***

	// TODO(stevvooe): It is possible that referenced conponents could have
	// children, but this is rare. For now, we ignore this and only verify
	// that manfiest components are present.
	required = append([]ocispec.Descriptor***REMOVED***mfst.Config***REMOVED***, mfst.Layers...)

	for _, desc := range required ***REMOVED***
		ra, err := provider.ReaderAt(ctx, desc.Digest)
		if err != nil ***REMOVED***
			if errdefs.IsNotFound(err) ***REMOVED***
				missing = append(missing, desc)
				continue
			***REMOVED*** else ***REMOVED***
				return false, nil, nil, nil, errors.Wrapf(err, "failed to check image %v", desc.Digest)
			***REMOVED***
		***REMOVED***
		ra.Close()
		present = append(present, desc)

	***REMOVED***

	return true, required, present, missing, nil
***REMOVED***

// Children returns the immediate children of content described by the descriptor.
func Children(ctx context.Context, provider content.Provider, desc ocispec.Descriptor, platform string) ([]ocispec.Descriptor, error) ***REMOVED***
	var descs []ocispec.Descriptor
	switch desc.MediaType ***REMOVED***
	case MediaTypeDockerSchema2Manifest, ocispec.MediaTypeImageManifest:
		p, err := content.ReadBlob(ctx, provider, desc.Digest)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		// TODO(stevvooe): We just assume oci manifest, for now. There may be
		// subtle differences from the docker version.
		var manifest ocispec.Manifest
		if err := json.Unmarshal(p, &manifest); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		descs = append(descs, manifest.Config)
		descs = append(descs, manifest.Layers...)
	case MediaTypeDockerSchema2ManifestList, ocispec.MediaTypeImageIndex:
		p, err := content.ReadBlob(ctx, provider, desc.Digest)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		var index ocispec.Index
		if err := json.Unmarshal(p, &index); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		if platform != "" ***REMOVED***
			matcher, err := platforms.Parse(platform)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			for _, d := range index.Manifests ***REMOVED***
				if d.Platform == nil || matcher.Match(*d.Platform) ***REMOVED***
					descs = append(descs, d)
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			descs = append(descs, index.Manifests...)
		***REMOVED***

	case MediaTypeDockerSchema2Layer, MediaTypeDockerSchema2LayerGzip,
		MediaTypeDockerSchema2LayerForeign, MediaTypeDockerSchema2LayerForeignGzip,
		MediaTypeDockerSchema2Config, ocispec.MediaTypeImageConfig,
		ocispec.MediaTypeImageLayer, ocispec.MediaTypeImageLayerGzip,
		ocispec.MediaTypeImageLayerNonDistributable, ocispec.MediaTypeImageLayerNonDistributableGzip,
		MediaTypeContainerd1Checkpoint, MediaTypeContainerd1CheckpointConfig:
		// childless data types.
		return nil, nil
	default:
		log.G(ctx).Warnf("encountered unknown type %v; children may not be fetched", desc.MediaType)
	***REMOVED***

	return descs, nil
***REMOVED***

// RootFS returns the unpacked diffids that make up and images rootfs.
//
// These are used to verify that a set of layers unpacked to the expected
// values.
func RootFS(ctx context.Context, provider content.Provider, configDesc ocispec.Descriptor) ([]digest.Digest, error) ***REMOVED***
	p, err := content.ReadBlob(ctx, provider, configDesc.Digest)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var config ocispec.Image
	if err := json.Unmarshal(p, &config); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// TODO(stevvooe): Remove this bit when OCI structure uses correct type for
	// rootfs.DiffIDs.
	var diffIDs []digest.Digest
	for _, diffID := range config.RootFS.DiffIDs ***REMOVED***
		diffIDs = append(diffIDs, digest.Digest(diffID))
	***REMOVED***

	return diffIDs, nil
***REMOVED***
