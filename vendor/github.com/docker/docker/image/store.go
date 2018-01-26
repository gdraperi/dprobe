package image

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/docker/distribution/digestset"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/system"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Store is an interface for creating and accessing images
type Store interface ***REMOVED***
	Create(config []byte) (ID, error)
	Get(id ID) (*Image, error)
	Delete(id ID) ([]layer.Metadata, error)
	Search(partialID string) (ID, error)
	SetParent(id ID, parent ID) error
	GetParent(id ID) (ID, error)
	SetLastUpdated(id ID) error
	GetLastUpdated(id ID) (time.Time, error)
	Children(id ID) []ID
	Map() map[ID]*Image
	Heads() map[ID]*Image
***REMOVED***

// LayerGetReleaser is a minimal interface for getting and releasing images.
type LayerGetReleaser interface ***REMOVED***
	Get(layer.ChainID) (layer.Layer, error)
	Release(layer.Layer) ([]layer.Metadata, error)
***REMOVED***

type imageMeta struct ***REMOVED***
	layer    layer.Layer
	children map[ID]struct***REMOVED******REMOVED***
***REMOVED***

type store struct ***REMOVED***
	sync.RWMutex
	lss       map[string]LayerGetReleaser
	images    map[ID]*imageMeta
	fs        StoreBackend
	digestSet *digestset.Set
***REMOVED***

// NewImageStore returns new store object for given set of layer stores
func NewImageStore(fs StoreBackend, lss map[string]LayerGetReleaser) (Store, error) ***REMOVED***
	is := &store***REMOVED***
		lss:       lss,
		images:    make(map[ID]*imageMeta),
		fs:        fs,
		digestSet: digestset.NewSet(),
	***REMOVED***

	// load all current images and retain layers
	if err := is.restore(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return is, nil
***REMOVED***

func (is *store) restore() error ***REMOVED***
	err := is.fs.Walk(func(dgst digest.Digest) error ***REMOVED***
		img, err := is.Get(IDFromDigest(dgst))
		if err != nil ***REMOVED***
			logrus.Errorf("invalid image %v, %v", dgst, err)
			return nil
		***REMOVED***
		var l layer.Layer
		if chainID := img.RootFS.ChainID(); chainID != "" ***REMOVED***
			if !system.IsOSSupported(img.OperatingSystem()) ***REMOVED***
				return system.ErrNotSupportedOperatingSystem
			***REMOVED***
			l, err = is.lss[img.OperatingSystem()].Get(chainID)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if err := is.digestSet.Add(dgst); err != nil ***REMOVED***
			return err
		***REMOVED***

		imageMeta := &imageMeta***REMOVED***
			layer:    l,
			children: make(map[ID]struct***REMOVED******REMOVED***),
		***REMOVED***

		is.images[IDFromDigest(dgst)] = imageMeta

		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Second pass to fill in children maps
	for id := range is.images ***REMOVED***
		if parent, err := is.GetParent(id); err == nil ***REMOVED***
			if parentMeta := is.images[parent]; parentMeta != nil ***REMOVED***
				parentMeta.children[id] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (is *store) Create(config []byte) (ID, error) ***REMOVED***
	var img Image
	err := json.Unmarshal(config, &img)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	// Must reject any config that references diffIDs from the history
	// which aren't among the rootfs layers.
	rootFSLayers := make(map[layer.DiffID]struct***REMOVED******REMOVED***)
	for _, diffID := range img.RootFS.DiffIDs ***REMOVED***
		rootFSLayers[diffID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	layerCounter := 0
	for _, h := range img.History ***REMOVED***
		if !h.EmptyLayer ***REMOVED***
			layerCounter++
		***REMOVED***
	***REMOVED***
	if layerCounter > len(img.RootFS.DiffIDs) ***REMOVED***
		return "", errors.New("too many non-empty layers in History section")
	***REMOVED***

	dgst, err := is.fs.Set(config)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	imageID := IDFromDigest(dgst)

	is.Lock()
	defer is.Unlock()

	if _, exists := is.images[imageID]; exists ***REMOVED***
		return imageID, nil
	***REMOVED***

	layerID := img.RootFS.ChainID()

	var l layer.Layer
	if layerID != "" ***REMOVED***
		if !system.IsOSSupported(img.OperatingSystem()) ***REMOVED***
			return "", system.ErrNotSupportedOperatingSystem
		***REMOVED***
		l, err = is.lss[img.OperatingSystem()].Get(layerID)
		if err != nil ***REMOVED***
			return "", errors.Wrapf(err, "failed to get layer %s", layerID)
		***REMOVED***
	***REMOVED***

	imageMeta := &imageMeta***REMOVED***
		layer:    l,
		children: make(map[ID]struct***REMOVED******REMOVED***),
	***REMOVED***

	is.images[imageID] = imageMeta
	if err := is.digestSet.Add(imageID.Digest()); err != nil ***REMOVED***
		delete(is.images, imageID)
		return "", err
	***REMOVED***

	return imageID, nil
***REMOVED***

type imageNotFoundError string

func (e imageNotFoundError) Error() string ***REMOVED***
	return "No such image: " + string(e)
***REMOVED***

func (imageNotFoundError) NotFound() ***REMOVED******REMOVED***

func (is *store) Search(term string) (ID, error) ***REMOVED***
	dgst, err := is.digestSet.Lookup(term)
	if err != nil ***REMOVED***
		if err == digestset.ErrDigestNotFound ***REMOVED***
			err = imageNotFoundError(term)
		***REMOVED***
		return "", errors.WithStack(err)
	***REMOVED***
	return IDFromDigest(dgst), nil
***REMOVED***

func (is *store) Get(id ID) (*Image, error) ***REMOVED***
	// todo: Check if image is in images
	// todo: Detect manual insertions and start using them
	config, err := is.fs.Get(id.Digest())
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	img, err := NewFromJSON(config)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	img.computedID = id

	img.Parent, err = is.GetParent(id)
	if err != nil ***REMOVED***
		img.Parent = ""
	***REMOVED***

	return img, nil
***REMOVED***

func (is *store) Delete(id ID) ([]layer.Metadata, error) ***REMOVED***
	is.Lock()
	defer is.Unlock()

	imageMeta := is.images[id]
	if imageMeta == nil ***REMOVED***
		return nil, fmt.Errorf("unrecognized image ID %s", id.String())
	***REMOVED***
	img, err := is.Get(id)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("unrecognized image %s, %v", id.String(), err)
	***REMOVED***
	if !system.IsOSSupported(img.OperatingSystem()) ***REMOVED***
		return nil, fmt.Errorf("unsupported image operating system %q", img.OperatingSystem())
	***REMOVED***
	for id := range imageMeta.children ***REMOVED***
		is.fs.DeleteMetadata(id.Digest(), "parent")
	***REMOVED***
	if parent, err := is.GetParent(id); err == nil && is.images[parent] != nil ***REMOVED***
		delete(is.images[parent].children, id)
	***REMOVED***

	if err := is.digestSet.Remove(id.Digest()); err != nil ***REMOVED***
		logrus.Errorf("error removing %s from digest set: %q", id, err)
	***REMOVED***
	delete(is.images, id)
	is.fs.Delete(id.Digest())

	if imageMeta.layer != nil ***REMOVED***
		return is.lss[img.OperatingSystem()].Release(imageMeta.layer)
	***REMOVED***
	return nil, nil
***REMOVED***

func (is *store) SetParent(id, parent ID) error ***REMOVED***
	is.Lock()
	defer is.Unlock()
	parentMeta := is.images[parent]
	if parentMeta == nil ***REMOVED***
		return fmt.Errorf("unknown parent image ID %s", parent.String())
	***REMOVED***
	if parent, err := is.GetParent(id); err == nil && is.images[parent] != nil ***REMOVED***
		delete(is.images[parent].children, id)
	***REMOVED***
	parentMeta.children[id] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	return is.fs.SetMetadata(id.Digest(), "parent", []byte(parent))
***REMOVED***

func (is *store) GetParent(id ID) (ID, error) ***REMOVED***
	d, err := is.fs.GetMetadata(id.Digest(), "parent")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return ID(d), nil // todo: validate?
***REMOVED***

// SetLastUpdated time for the image ID to the current time
func (is *store) SetLastUpdated(id ID) error ***REMOVED***
	lastUpdated := []byte(time.Now().Format(time.RFC3339Nano))
	return is.fs.SetMetadata(id.Digest(), "lastUpdated", lastUpdated)
***REMOVED***

// GetLastUpdated time for the image ID
func (is *store) GetLastUpdated(id ID) (time.Time, error) ***REMOVED***
	bytes, err := is.fs.GetMetadata(id.Digest(), "lastUpdated")
	if err != nil || len(bytes) == 0 ***REMOVED***
		// No lastUpdated time
		return time.Time***REMOVED******REMOVED***, nil
	***REMOVED***
	return time.Parse(time.RFC3339Nano, string(bytes))
***REMOVED***

func (is *store) Children(id ID) []ID ***REMOVED***
	is.RLock()
	defer is.RUnlock()

	return is.children(id)
***REMOVED***

func (is *store) children(id ID) []ID ***REMOVED***
	var ids []ID
	if is.images[id] != nil ***REMOVED***
		for id := range is.images[id].children ***REMOVED***
			ids = append(ids, id)
		***REMOVED***
	***REMOVED***
	return ids
***REMOVED***

func (is *store) Heads() map[ID]*Image ***REMOVED***
	return is.imagesMap(false)
***REMOVED***

func (is *store) Map() map[ID]*Image ***REMOVED***
	return is.imagesMap(true)
***REMOVED***

func (is *store) imagesMap(all bool) map[ID]*Image ***REMOVED***
	is.RLock()
	defer is.RUnlock()

	images := make(map[ID]*Image)

	for id := range is.images ***REMOVED***
		if !all && len(is.children(id)) > 0 ***REMOVED***
			continue
		***REMOVED***
		img, err := is.Get(id)
		if err != nil ***REMOVED***
			logrus.Errorf("invalid image access: %q, error: %q", id, err)
			continue
		***REMOVED***
		images[id] = img
	***REMOVED***
	return images
***REMOVED***
