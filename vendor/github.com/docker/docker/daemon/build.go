package daemon

import (
	"io"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/builder"
	"github.com/docker/docker/image"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/docker/registry"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type releaseableLayer struct ***REMOVED***
	released   bool
	layerStore layer.Store
	roLayer    layer.Layer
	rwLayer    layer.RWLayer
***REMOVED***

func (rl *releaseableLayer) Mount() (containerfs.ContainerFS, error) ***REMOVED***
	var err error
	var mountPath containerfs.ContainerFS
	var chainID layer.ChainID
	if rl.roLayer != nil ***REMOVED***
		chainID = rl.roLayer.ChainID()
	***REMOVED***

	mountID := stringid.GenerateRandomID()
	rl.rwLayer, err = rl.layerStore.CreateRWLayer(mountID, chainID, nil)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "failed to create rwlayer")
	***REMOVED***

	mountPath, err = rl.rwLayer.Mount("")
	if err != nil ***REMOVED***
		// Clean up the layer if we fail to mount it here.
		metadata, err := rl.layerStore.ReleaseRWLayer(rl.rwLayer)
		layer.LogReleaseMetadata(metadata)
		if err != nil ***REMOVED***
			logrus.Errorf("Failed to release RWLayer: %s", err)
		***REMOVED***
		rl.rwLayer = nil
		return nil, err
	***REMOVED***

	return mountPath, nil
***REMOVED***

func (rl *releaseableLayer) Commit() (builder.ReleaseableLayer, error) ***REMOVED***
	var chainID layer.ChainID
	if rl.roLayer != nil ***REMOVED***
		chainID = rl.roLayer.ChainID()
	***REMOVED***

	stream, err := rl.rwLayer.TarStream()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer stream.Close()

	newLayer, err := rl.layerStore.Register(stream, chainID)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// TODO: An optimization would be to handle empty layers before returning
	return &releaseableLayer***REMOVED***layerStore: rl.layerStore, roLayer: newLayer***REMOVED***, nil
***REMOVED***

func (rl *releaseableLayer) DiffID() layer.DiffID ***REMOVED***
	if rl.roLayer == nil ***REMOVED***
		return layer.DigestSHA256EmptyTar
	***REMOVED***
	return rl.roLayer.DiffID()
***REMOVED***

func (rl *releaseableLayer) Release() error ***REMOVED***
	if rl.released ***REMOVED***
		return nil
	***REMOVED***
	if err := rl.releaseRWLayer(); err != nil ***REMOVED***
		// Best effort attempt at releasing read-only layer before returning original error.
		rl.releaseROLayer()
		return err
	***REMOVED***
	if err := rl.releaseROLayer(); err != nil ***REMOVED***
		return err
	***REMOVED***
	rl.released = true
	return nil
***REMOVED***

func (rl *releaseableLayer) releaseRWLayer() error ***REMOVED***
	if rl.rwLayer == nil ***REMOVED***
		return nil
	***REMOVED***
	if err := rl.rwLayer.Unmount(); err != nil ***REMOVED***
		logrus.Errorf("Failed to unmount RWLayer: %s", err)
		return err
	***REMOVED***
	metadata, err := rl.layerStore.ReleaseRWLayer(rl.rwLayer)
	layer.LogReleaseMetadata(metadata)
	if err != nil ***REMOVED***
		logrus.Errorf("Failed to release RWLayer: %s", err)
	***REMOVED***
	rl.rwLayer = nil
	return err
***REMOVED***

func (rl *releaseableLayer) releaseROLayer() error ***REMOVED***
	if rl.roLayer == nil ***REMOVED***
		return nil
	***REMOVED***
	metadata, err := rl.layerStore.Release(rl.roLayer)
	layer.LogReleaseMetadata(metadata)
	if err != nil ***REMOVED***
		logrus.Errorf("Failed to release ROLayer: %s", err)
	***REMOVED***
	rl.roLayer = nil
	return err
***REMOVED***

func newReleasableLayerForImage(img *image.Image, layerStore layer.Store) (builder.ReleaseableLayer, error) ***REMOVED***
	if img == nil || img.RootFS.ChainID() == "" ***REMOVED***
		return &releaseableLayer***REMOVED***layerStore: layerStore***REMOVED***, nil
	***REMOVED***
	// Hold a reference to the image layer so that it can't be removed before
	// it is released
	roLayer, err := layerStore.Get(img.RootFS.ChainID())
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "failed to get layer for image %s", img.ImageID())
	***REMOVED***
	return &releaseableLayer***REMOVED***layerStore: layerStore, roLayer: roLayer***REMOVED***, nil
***REMOVED***

// TODO: could this use the regular daemon PullImage ?
func (daemon *Daemon) pullForBuilder(ctx context.Context, name string, authConfigs map[string]types.AuthConfig, output io.Writer, os string) (*image.Image, error) ***REMOVED***
	ref, err := reference.ParseNormalizedNamed(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ref = reference.TagNameOnly(ref)

	pullRegistryAuth := &types.AuthConfig***REMOVED******REMOVED***
	if len(authConfigs) > 0 ***REMOVED***
		// The request came with a full auth config, use it
		repoInfo, err := daemon.RegistryService.ResolveRepository(ref)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		resolvedConfig := registry.ResolveAuthConfig(authConfigs, repoInfo.Index)
		pullRegistryAuth = &resolvedConfig
	***REMOVED***

	if err := daemon.pullImageWithReference(ctx, ref, os, nil, pullRegistryAuth, output); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return daemon.GetImage(name)
***REMOVED***

// GetImageAndReleasableLayer returns an image and releaseable layer for a reference or ID.
// Every call to GetImageAndReleasableLayer MUST call releasableLayer.Release() to prevent
// leaking of layers.
func (daemon *Daemon) GetImageAndReleasableLayer(ctx context.Context, refOrID string, opts backend.GetImageAndLayerOptions) (builder.Image, builder.ReleaseableLayer, error) ***REMOVED***
	if refOrID == "" ***REMOVED***
		if !system.IsOSSupported(opts.OS) ***REMOVED***
			return nil, nil, system.ErrNotSupportedOperatingSystem
		***REMOVED***
		layer, err := newReleasableLayerForImage(nil, daemon.layerStores[opts.OS])
		return nil, layer, err
	***REMOVED***

	if opts.PullOption != backend.PullOptionForcePull ***REMOVED***
		image, err := daemon.GetImage(refOrID)
		if err != nil && opts.PullOption == backend.PullOptionNoPull ***REMOVED***
			return nil, nil, err
		***REMOVED***
		// TODO: shouldn't we error out if error is different from "not found" ?
		if image != nil ***REMOVED***
			if !system.IsOSSupported(image.OperatingSystem()) ***REMOVED***
				return nil, nil, system.ErrNotSupportedOperatingSystem
			***REMOVED***
			layer, err := newReleasableLayerForImage(image, daemon.layerStores[image.OperatingSystem()])
			return image, layer, err
		***REMOVED***
	***REMOVED***

	image, err := daemon.pullForBuilder(ctx, refOrID, opts.AuthConfig, opts.Output, opts.OS)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	if !system.IsOSSupported(image.OperatingSystem()) ***REMOVED***
		return nil, nil, system.ErrNotSupportedOperatingSystem
	***REMOVED***
	layer, err := newReleasableLayerForImage(image, daemon.layerStores[image.OperatingSystem()])
	return image, layer, err
***REMOVED***

// CreateImage creates a new image by adding a config and ID to the image store.
// This is similar to LoadImage() except that it receives JSON encoded bytes of
// an image instead of a tar archive.
func (daemon *Daemon) CreateImage(config []byte, parent string) (builder.Image, error) ***REMOVED***
	id, err := daemon.imageStore.Create(config)
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "failed to create image")
	***REMOVED***

	if parent != "" ***REMOVED***
		if err := daemon.imageStore.SetParent(id, image.ID(parent)); err != nil ***REMOVED***
			return nil, errors.Wrapf(err, "failed to set parent %s", parent)
		***REMOVED***
	***REMOVED***

	return daemon.imageStore.Get(id)
***REMOVED***

// IDMappings returns uid/gid mappings for the builder
func (daemon *Daemon) IDMappings() *idtools.IDMappings ***REMOVED***
	return daemon.idMappings
***REMOVED***
