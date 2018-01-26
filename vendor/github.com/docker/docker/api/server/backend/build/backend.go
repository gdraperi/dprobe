package build

import (
	"fmt"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/builder"
	"github.com/docker/docker/builder/fscache"
	"github.com/docker/docker/image"
	"github.com/docker/docker/pkg/stringid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// ImageComponent provides an interface for working with images
type ImageComponent interface ***REMOVED***
	SquashImage(from string, to string) (string, error)
	TagImageWithReference(image.ID, reference.Named) error
***REMOVED***

// Builder defines interface for running a build
type Builder interface ***REMOVED***
	Build(context.Context, backend.BuildConfig) (*builder.Result, error)
***REMOVED***

// Backend provides build functionality to the API router
type Backend struct ***REMOVED***
	builder        Builder
	fsCache        *fscache.FSCache
	imageComponent ImageComponent
***REMOVED***

// NewBackend creates a new build backend from components
func NewBackend(components ImageComponent, builder Builder, fsCache *fscache.FSCache) (*Backend, error) ***REMOVED***
	return &Backend***REMOVED***imageComponent: components, builder: builder, fsCache: fsCache***REMOVED***, nil
***REMOVED***

// Build builds an image from a Source
func (b *Backend) Build(ctx context.Context, config backend.BuildConfig) (string, error) ***REMOVED***
	options := config.Options
	tagger, err := NewTagger(b.imageComponent, config.ProgressWriter.StdoutFormatter, options.Tags)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	build, err := b.builder.Build(ctx, config)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	var imageID = build.ImageID
	if options.Squash ***REMOVED***
		if imageID, err = squashBuild(build, b.imageComponent); err != nil ***REMOVED***
			return "", err
		***REMOVED***
		if config.ProgressWriter.AuxFormatter != nil ***REMOVED***
			if err = config.ProgressWriter.AuxFormatter.Emit(types.BuildResult***REMOVED***ID: imageID***REMOVED***); err != nil ***REMOVED***
				return "", err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	stdout := config.ProgressWriter.StdoutFormatter
	fmt.Fprintf(stdout, "Successfully built %s\n", stringid.TruncateID(imageID))
	err = tagger.TagImages(image.ID(imageID))
	return imageID, err
***REMOVED***

// PruneCache removes all cached build sources
func (b *Backend) PruneCache(ctx context.Context) (*types.BuildCachePruneReport, error) ***REMOVED***
	size, err := b.fsCache.Prune(ctx)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "failed to prune build cache")
	***REMOVED***
	return &types.BuildCachePruneReport***REMOVED***SpaceReclaimed: size***REMOVED***, nil
***REMOVED***

func squashBuild(build *builder.Result, imageComponent ImageComponent) (string, error) ***REMOVED***
	var fromID string
	if build.FromImage != nil ***REMOVED***
		fromID = build.FromImage.ImageID()
	***REMOVED***
	imageID, err := imageComponent.SquashImage(build.ImageID, fromID)
	if err != nil ***REMOVED***
		return "", errors.Wrap(err, "error squashing image")
	***REMOVED***
	return imageID, nil
***REMOVED***
