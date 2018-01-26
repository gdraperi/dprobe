package cache

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/dockerversion"
	"github.com/docker/docker/image"
	"github.com/docker/docker/layer"
	"github.com/pkg/errors"
)

// NewLocal returns a local image cache, based on parent chain
func NewLocal(store image.Store) *LocalImageCache ***REMOVED***
	return &LocalImageCache***REMOVED***
		store: store,
	***REMOVED***
***REMOVED***

// LocalImageCache is cache based on parent chain.
type LocalImageCache struct ***REMOVED***
	store image.Store
***REMOVED***

// GetCache returns the image id found in the cache
func (lic *LocalImageCache) GetCache(imgID string, config *containertypes.Config) (string, error) ***REMOVED***
	return getImageIDAndError(getLocalCachedImage(lic.store, image.ID(imgID), config))
***REMOVED***

// New returns an image cache, based on history objects
func New(store image.Store) *ImageCache ***REMOVED***
	return &ImageCache***REMOVED***
		store:           store,
		localImageCache: NewLocal(store),
	***REMOVED***
***REMOVED***

// ImageCache is cache based on history objects. Requires initial set of images.
type ImageCache struct ***REMOVED***
	sources         []*image.Image
	store           image.Store
	localImageCache *LocalImageCache
***REMOVED***

// Populate adds an image to the cache (to be queried later)
func (ic *ImageCache) Populate(image *image.Image) ***REMOVED***
	ic.sources = append(ic.sources, image)
***REMOVED***

// GetCache returns the image id found in the cache
func (ic *ImageCache) GetCache(parentID string, cfg *containertypes.Config) (string, error) ***REMOVED***
	imgID, err := ic.localImageCache.GetCache(parentID, cfg)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if imgID != "" ***REMOVED***
		for _, s := range ic.sources ***REMOVED***
			if ic.isParent(s.ID(), image.ID(imgID)) ***REMOVED***
				return imgID, nil
			***REMOVED***
		***REMOVED***
	***REMOVED***

	var parent *image.Image
	lenHistory := 0
	if parentID != "" ***REMOVED***
		parent, err = ic.store.Get(image.ID(parentID))
		if err != nil ***REMOVED***
			return "", errors.Wrapf(err, "unable to find image %v", parentID)
		***REMOVED***
		lenHistory = len(parent.History)
	***REMOVED***

	for _, target := range ic.sources ***REMOVED***
		if !isValidParent(target, parent) || !isValidConfig(cfg, target.History[lenHistory]) ***REMOVED***
			continue
		***REMOVED***

		if len(target.History)-1 == lenHistory ***REMOVED*** // last
			if parent != nil ***REMOVED***
				if err := ic.store.SetParent(target.ID(), parent.ID()); err != nil ***REMOVED***
					return "", errors.Wrapf(err, "failed to set parent for %v to %v", target.ID(), parent.ID())
				***REMOVED***
			***REMOVED***
			return target.ID().String(), nil
		***REMOVED***

		imgID, err := ic.restoreCachedImage(parent, target, cfg)
		if err != nil ***REMOVED***
			return "", errors.Wrapf(err, "failed to restore cached image from %q to %v", parentID, target.ID())
		***REMOVED***

		ic.sources = []*image.Image***REMOVED***target***REMOVED*** // avoid jumping to different target, tuned for safety atm
		return imgID.String(), nil
	***REMOVED***

	return "", nil
***REMOVED***

func (ic *ImageCache) restoreCachedImage(parent, target *image.Image, cfg *containertypes.Config) (image.ID, error) ***REMOVED***
	var history []image.History
	rootFS := image.NewRootFS()
	lenHistory := 0
	if parent != nil ***REMOVED***
		history = parent.History
		rootFS = parent.RootFS
		lenHistory = len(parent.History)
	***REMOVED***
	history = append(history, target.History[lenHistory])
	if layer := getLayerForHistoryIndex(target, lenHistory); layer != "" ***REMOVED***
		rootFS.Append(layer)
	***REMOVED***

	config, err := json.Marshal(&image.Image***REMOVED***
		V1Image: image.V1Image***REMOVED***
			DockerVersion: dockerversion.Version,
			Config:        cfg,
			Architecture:  target.Architecture,
			OS:            target.OS,
			Author:        target.Author,
			Created:       history[len(history)-1].Created,
		***REMOVED***,
		RootFS:     rootFS,
		History:    history,
		OSFeatures: target.OSFeatures,
		OSVersion:  target.OSVersion,
	***REMOVED***)
	if err != nil ***REMOVED***
		return "", errors.Wrap(err, "failed to marshal image config")
	***REMOVED***

	imgID, err := ic.store.Create(config)
	if err != nil ***REMOVED***
		return "", errors.Wrap(err, "failed to create cache image")
	***REMOVED***

	if parent != nil ***REMOVED***
		if err := ic.store.SetParent(imgID, parent.ID()); err != nil ***REMOVED***
			return "", errors.Wrapf(err, "failed to set parent for %v to %v", target.ID(), parent.ID())
		***REMOVED***
	***REMOVED***
	return imgID, nil
***REMOVED***

func (ic *ImageCache) isParent(imgID, parentID image.ID) bool ***REMOVED***
	nextParent, err := ic.store.GetParent(imgID)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	if nextParent == parentID ***REMOVED***
		return true
	***REMOVED***
	return ic.isParent(nextParent, parentID)
***REMOVED***

func getLayerForHistoryIndex(image *image.Image, index int) layer.DiffID ***REMOVED***
	layerIndex := 0
	for i, h := range image.History ***REMOVED***
		if i == index ***REMOVED***
			if h.EmptyLayer ***REMOVED***
				return ""
			***REMOVED***
			break
		***REMOVED***
		if !h.EmptyLayer ***REMOVED***
			layerIndex++
		***REMOVED***
	***REMOVED***
	return image.RootFS.DiffIDs[layerIndex] // validate?
***REMOVED***

func isValidConfig(cfg *containertypes.Config, h image.History) bool ***REMOVED***
	// todo: make this format better than join that loses data
	return strings.Join(cfg.Cmd, " ") == h.CreatedBy
***REMOVED***

func isValidParent(img, parent *image.Image) bool ***REMOVED***
	if len(img.History) == 0 ***REMOVED***
		return false
	***REMOVED***
	if parent == nil || len(parent.History) == 0 && len(parent.RootFS.DiffIDs) == 0 ***REMOVED***
		return true
	***REMOVED***
	if len(parent.History) >= len(img.History) ***REMOVED***
		return false
	***REMOVED***
	if len(parent.RootFS.DiffIDs) > len(img.RootFS.DiffIDs) ***REMOVED***
		return false
	***REMOVED***

	for i, h := range parent.History ***REMOVED***
		if !reflect.DeepEqual(h, img.History[i]) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	for i, d := range parent.RootFS.DiffIDs ***REMOVED***
		if d != img.RootFS.DiffIDs[i] ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func getImageIDAndError(img *image.Image, err error) (string, error) ***REMOVED***
	if img == nil || err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return img.ID().String(), nil
***REMOVED***

// getLocalCachedImage returns the most recent created image that is a child
// of the image with imgID, that had the same config when it was
// created. nil is returned if a child cannot be found. An error is
// returned if the parent image cannot be found.
func getLocalCachedImage(imageStore image.Store, imgID image.ID, config *containertypes.Config) (*image.Image, error) ***REMOVED***
	// Loop on the children of the given image and check the config
	getMatch := func(siblings []image.ID) (*image.Image, error) ***REMOVED***
		var match *image.Image
		for _, id := range siblings ***REMOVED***
			img, err := imageStore.Get(id)
			if err != nil ***REMOVED***
				return nil, fmt.Errorf("unable to find image %q", id)
			***REMOVED***

			if compare(&img.ContainerConfig, config) ***REMOVED***
				// check for the most up to date match
				if match == nil || match.Created.Before(img.Created) ***REMOVED***
					match = img
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return match, nil
	***REMOVED***

	// In this case, this is `FROM scratch`, which isn't an actual image.
	if imgID == "" ***REMOVED***
		images := imageStore.Map()
		var siblings []image.ID
		for id, img := range images ***REMOVED***
			if img.Parent == imgID ***REMOVED***
				siblings = append(siblings, id)
			***REMOVED***
		***REMOVED***
		return getMatch(siblings)
	***REMOVED***

	// find match from child images
	siblings := imageStore.Children(imgID)
	return getMatch(siblings)
***REMOVED***
