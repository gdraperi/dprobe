package dockerfile

import (
	"runtime"

	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/builder"
	"github.com/docker/docker/builder/remotecontext"
	dockerimage "github.com/docker/docker/image"
	"github.com/docker/docker/pkg/system"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type getAndMountFunc func(string, bool) (builder.Image, builder.ReleaseableLayer, error)

// imageSources mounts images and provides a cache for mounted images. It tracks
// all images so they can be unmounted at the end of the build.
type imageSources struct ***REMOVED***
	byImageID map[string]*imageMount
	mounts    []*imageMount
	getImage  getAndMountFunc
***REMOVED***

func newImageSources(ctx context.Context, options builderOptions) *imageSources ***REMOVED***
	getAndMount := func(idOrRef string, localOnly bool) (builder.Image, builder.ReleaseableLayer, error) ***REMOVED***
		pullOption := backend.PullOptionNoPull
		if !localOnly ***REMOVED***
			if options.Options.PullParent ***REMOVED***
				pullOption = backend.PullOptionForcePull
			***REMOVED*** else ***REMOVED***
				pullOption = backend.PullOptionPreferLocal
			***REMOVED***
		***REMOVED***
		optionsPlatform := system.ParsePlatform(options.Options.Platform)
		return options.Backend.GetImageAndReleasableLayer(ctx, idOrRef, backend.GetImageAndLayerOptions***REMOVED***
			PullOption: pullOption,
			AuthConfig: options.Options.AuthConfigs,
			Output:     options.ProgressWriter.Output,
			OS:         optionsPlatform.OS,
		***REMOVED***)
	***REMOVED***

	return &imageSources***REMOVED***
		byImageID: make(map[string]*imageMount),
		getImage:  getAndMount,
	***REMOVED***
***REMOVED***

func (m *imageSources) Get(idOrRef string, localOnly bool) (*imageMount, error) ***REMOVED***
	if im, ok := m.byImageID[idOrRef]; ok ***REMOVED***
		return im, nil
	***REMOVED***

	image, layer, err := m.getImage(idOrRef, localOnly)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	im := newImageMount(image, layer)
	m.Add(im)
	return im, nil
***REMOVED***

func (m *imageSources) Unmount() (retErr error) ***REMOVED***
	for _, im := range m.mounts ***REMOVED***
		if err := im.unmount(); err != nil ***REMOVED***
			logrus.Error(err)
			retErr = err
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (m *imageSources) Add(im *imageMount) ***REMOVED***
	switch im.image ***REMOVED***
	case nil:
		// set the OS for scratch images
		os := runtime.GOOS
		// Windows does not support scratch except for LCOW
		if runtime.GOOS == "windows" ***REMOVED***
			os = "linux"
		***REMOVED***
		im.image = &dockerimage.Image***REMOVED***V1Image: dockerimage.V1Image***REMOVED***OS: os***REMOVED******REMOVED***
	default:
		m.byImageID[im.image.ImageID()] = im
	***REMOVED***
	m.mounts = append(m.mounts, im)
***REMOVED***

// imageMount is a reference to an image that can be used as a builder.Source
type imageMount struct ***REMOVED***
	image  builder.Image
	source builder.Source
	layer  builder.ReleaseableLayer
***REMOVED***

func newImageMount(image builder.Image, layer builder.ReleaseableLayer) *imageMount ***REMOVED***
	im := &imageMount***REMOVED***image: image, layer: layer***REMOVED***
	return im
***REMOVED***

func (im *imageMount) Source() (builder.Source, error) ***REMOVED***
	if im.source == nil ***REMOVED***
		if im.layer == nil ***REMOVED***
			return nil, errors.Errorf("empty context")
		***REMOVED***
		mountPath, err := im.layer.Mount()
		if err != nil ***REMOVED***
			return nil, errors.Wrapf(err, "failed to mount %s", im.image.ImageID())
		***REMOVED***
		source, err := remotecontext.NewLazySource(mountPath)
		if err != nil ***REMOVED***
			return nil, errors.Wrapf(err, "failed to create lazycontext for %s", mountPath)
		***REMOVED***
		im.source = source
	***REMOVED***
	return im.source, nil
***REMOVED***

func (im *imageMount) unmount() error ***REMOVED***
	if im.layer == nil ***REMOVED***
		return nil
	***REMOVED***
	if err := im.layer.Release(); err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to unmount previous build image %s", im.image.ImageID())
	***REMOVED***
	im.layer = nil
	return nil
***REMOVED***

func (im *imageMount) Image() builder.Image ***REMOVED***
	return im.image
***REMOVED***

func (im *imageMount) Layer() builder.ReleaseableLayer ***REMOVED***
	return im.layer
***REMOVED***

func (im *imageMount) ImageID() string ***REMOVED***
	return im.image.ImageID()
***REMOVED***
