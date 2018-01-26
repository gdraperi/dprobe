package daemon

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/pkg/errors"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/container"
	"github.com/docker/docker/image"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/system"
)

var acceptedImageFilterTags = map[string]bool***REMOVED***
	"dangling":  true,
	"label":     true,
	"before":    true,
	"since":     true,
	"reference": true,
***REMOVED***

// byCreated is a temporary type used to sort a list of images by creation
// time.
type byCreated []*types.ImageSummary

func (r byCreated) Len() int           ***REMOVED*** return len(r) ***REMOVED***
func (r byCreated) Swap(i, j int)      ***REMOVED*** r[i], r[j] = r[j], r[i] ***REMOVED***
func (r byCreated) Less(i, j int) bool ***REMOVED*** return r[i].Created < r[j].Created ***REMOVED***

// Map returns a map of all images in the ImageStore
func (daemon *Daemon) Map() map[image.ID]*image.Image ***REMOVED***
	return daemon.imageStore.Map()
***REMOVED***

// Images returns a filtered list of images. filterArgs is a JSON-encoded set
// of filter arguments which will be interpreted by api/types/filters.
// filter is a shell glob string applied to repository names. The argument
// named all controls whether all images in the graph are filtered, or just
// the heads.
func (daemon *Daemon) Images(imageFilters filters.Args, all bool, withExtraAttrs bool) ([]*types.ImageSummary, error) ***REMOVED***
	var (
		allImages    map[image.ID]*image.Image
		err          error
		danglingOnly = false
	)

	if err := imageFilters.Validate(acceptedImageFilterTags); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if imageFilters.Contains("dangling") ***REMOVED***
		if imageFilters.ExactMatch("dangling", "true") ***REMOVED***
			danglingOnly = true
		***REMOVED*** else if !imageFilters.ExactMatch("dangling", "false") ***REMOVED***
			return nil, invalidFilter***REMOVED***"dangling", imageFilters.Get("dangling")***REMOVED***
		***REMOVED***
	***REMOVED***
	if danglingOnly ***REMOVED***
		allImages = daemon.imageStore.Heads()
	***REMOVED*** else ***REMOVED***
		allImages = daemon.imageStore.Map()
	***REMOVED***

	var beforeFilter, sinceFilter *image.Image
	err = imageFilters.WalkValues("before", func(value string) error ***REMOVED***
		beforeFilter, err = daemon.GetImage(value)
		return err
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err = imageFilters.WalkValues("since", func(value string) error ***REMOVED***
		sinceFilter, err = daemon.GetImage(value)
		return err
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	images := []*types.ImageSummary***REMOVED******REMOVED***
	var imagesMap map[*image.Image]*types.ImageSummary
	var layerRefs map[layer.ChainID]int
	var allLayers map[layer.ChainID]layer.Layer
	var allContainers []*container.Container

	for id, img := range allImages ***REMOVED***
		if beforeFilter != nil ***REMOVED***
			if img.Created.Equal(beforeFilter.Created) || img.Created.After(beforeFilter.Created) ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***

		if sinceFilter != nil ***REMOVED***
			if img.Created.Equal(sinceFilter.Created) || img.Created.Before(sinceFilter.Created) ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***

		if imageFilters.Contains("label") ***REMOVED***
			// Very old image that do not have image.Config (or even labels)
			if img.Config == nil ***REMOVED***
				continue
			***REMOVED***
			// We are now sure image.Config is not nil
			if !imageFilters.MatchKVList("label", img.Config.Labels) ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***

		// Skip any images with an unsupported operating system to avoid a potential
		// panic when indexing through the layerstore. Don't error as we want to list
		// the other images. This should never happen, but here as a safety precaution.
		if !system.IsOSSupported(img.OperatingSystem()) ***REMOVED***
			continue
		***REMOVED***

		layerID := img.RootFS.ChainID()
		var size int64
		if layerID != "" ***REMOVED***
			l, err := daemon.layerStores[img.OperatingSystem()].Get(layerID)
			if err != nil ***REMOVED***
				// The layer may have been deleted between the call to `Map()` or
				// `Heads()` and the call to `Get()`, so we just ignore this error
				if err == layer.ErrLayerDoesNotExist ***REMOVED***
					continue
				***REMOVED***
				return nil, err
			***REMOVED***

			size, err = l.Size()
			layer.ReleaseAndLog(daemon.layerStores[img.OperatingSystem()], l)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***

		newImage := newImage(img, size)

		for _, ref := range daemon.referenceStore.References(id.Digest()) ***REMOVED***
			if imageFilters.Contains("reference") ***REMOVED***
				var found bool
				var matchErr error
				for _, pattern := range imageFilters.Get("reference") ***REMOVED***
					found, matchErr = reference.FamiliarMatch(pattern, ref)
					if matchErr != nil ***REMOVED***
						return nil, matchErr
					***REMOVED***
				***REMOVED***
				if !found ***REMOVED***
					continue
				***REMOVED***
			***REMOVED***
			if _, ok := ref.(reference.Canonical); ok ***REMOVED***
				newImage.RepoDigests = append(newImage.RepoDigests, reference.FamiliarString(ref))
			***REMOVED***
			if _, ok := ref.(reference.NamedTagged); ok ***REMOVED***
				newImage.RepoTags = append(newImage.RepoTags, reference.FamiliarString(ref))
			***REMOVED***
		***REMOVED***
		if newImage.RepoDigests == nil && newImage.RepoTags == nil ***REMOVED***
			if all || len(daemon.imageStore.Children(id)) == 0 ***REMOVED***

				if imageFilters.Contains("dangling") && !danglingOnly ***REMOVED***
					//dangling=false case, so dangling image is not needed
					continue
				***REMOVED***
				if imageFilters.Contains("reference") ***REMOVED*** // skip images with no references if filtering by reference
					continue
				***REMOVED***
				newImage.RepoDigests = []string***REMOVED***"<none>@<none>"***REMOVED***
				newImage.RepoTags = []string***REMOVED***"<none>:<none>"***REMOVED***
			***REMOVED*** else ***REMOVED***
				continue
			***REMOVED***
		***REMOVED*** else if danglingOnly && len(newImage.RepoTags) > 0 ***REMOVED***
			continue
		***REMOVED***

		if withExtraAttrs ***REMOVED***
			// lazily init variables
			if imagesMap == nil ***REMOVED***
				allContainers = daemon.List()
				allLayers = daemon.layerStores[img.OperatingSystem()].Map()
				imagesMap = make(map[*image.Image]*types.ImageSummary)
				layerRefs = make(map[layer.ChainID]int)
			***REMOVED***

			// Get container count
			newImage.Containers = 0
			for _, c := range allContainers ***REMOVED***
				if c.ImageID == id ***REMOVED***
					newImage.Containers++
				***REMOVED***
			***REMOVED***

			// count layer references
			rootFS := *img.RootFS
			rootFS.DiffIDs = nil
			for _, id := range img.RootFS.DiffIDs ***REMOVED***
				rootFS.Append(id)
				chid := rootFS.ChainID()
				layerRefs[chid]++
				if _, ok := allLayers[chid]; !ok ***REMOVED***
					return nil, fmt.Errorf("layer %v was not found (corruption?)", chid)
				***REMOVED***
			***REMOVED***
			imagesMap[img] = newImage
		***REMOVED***

		images = append(images, newImage)
	***REMOVED***

	if withExtraAttrs ***REMOVED***
		// Get Shared sizes
		for img, newImage := range imagesMap ***REMOVED***
			rootFS := *img.RootFS
			rootFS.DiffIDs = nil

			newImage.SharedSize = 0
			for _, id := range img.RootFS.DiffIDs ***REMOVED***
				rootFS.Append(id)
				chid := rootFS.ChainID()

				diffSize, err := allLayers[chid].DiffSize()
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***

				if layerRefs[chid] > 1 ***REMOVED***
					newImage.SharedSize += diffSize
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	sort.Sort(sort.Reverse(byCreated(images)))

	return images, nil
***REMOVED***

// SquashImage creates a new image with the diff of the specified image and the specified parent.
// This new image contains only the layers from it's parent + 1 extra layer which contains the diff of all the layers in between.
// The existing image(s) is not destroyed.
// If no parent is specified, a new image with the diff of all the specified image's layers merged into a new layer that has no parents.
func (daemon *Daemon) SquashImage(id, parent string) (string, error) ***REMOVED***

	var (
		img *image.Image
		err error
	)
	if img, err = daemon.imageStore.Get(image.ID(id)); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	var parentImg *image.Image
	var parentChainID layer.ChainID
	if len(parent) != 0 ***REMOVED***
		parentImg, err = daemon.imageStore.Get(image.ID(parent))
		if err != nil ***REMOVED***
			return "", errors.Wrap(err, "error getting specified parent layer")
		***REMOVED***
		parentChainID = parentImg.RootFS.ChainID()
	***REMOVED*** else ***REMOVED***
		rootFS := image.NewRootFS()
		parentImg = &image.Image***REMOVED***RootFS: rootFS***REMOVED***
	***REMOVED***

	l, err := daemon.layerStores[img.OperatingSystem()].Get(img.RootFS.ChainID())
	if err != nil ***REMOVED***
		return "", errors.Wrap(err, "error getting image layer")
	***REMOVED***
	defer daemon.layerStores[img.OperatingSystem()].Release(l)

	ts, err := l.TarStreamFrom(parentChainID)
	if err != nil ***REMOVED***
		return "", errors.Wrapf(err, "error getting tar stream to parent")
	***REMOVED***
	defer ts.Close()

	newL, err := daemon.layerStores[img.OperatingSystem()].Register(ts, parentChainID)
	if err != nil ***REMOVED***
		return "", errors.Wrap(err, "error registering layer")
	***REMOVED***
	defer daemon.layerStores[img.OperatingSystem()].Release(newL)

	newImage := *img
	newImage.RootFS = nil

	rootFS := *parentImg.RootFS
	rootFS.DiffIDs = append(rootFS.DiffIDs, newL.DiffID())
	newImage.RootFS = &rootFS

	for i, hi := range newImage.History ***REMOVED***
		if i >= len(parentImg.History) ***REMOVED***
			hi.EmptyLayer = true
		***REMOVED***
		newImage.History[i] = hi
	***REMOVED***

	now := time.Now()
	var historyComment string
	if len(parent) > 0 ***REMOVED***
		historyComment = fmt.Sprintf("merge %s to %s", id, parent)
	***REMOVED*** else ***REMOVED***
		historyComment = fmt.Sprintf("create new from %s", id)
	***REMOVED***

	newImage.History = append(newImage.History, image.History***REMOVED***
		Created: now,
		Comment: historyComment,
	***REMOVED***)
	newImage.Created = now

	b, err := json.Marshal(&newImage)
	if err != nil ***REMOVED***
		return "", errors.Wrap(err, "error marshalling image config")
	***REMOVED***

	newImgID, err := daemon.imageStore.Create(b)
	if err != nil ***REMOVED***
		return "", errors.Wrap(err, "error creating new image after squash")
	***REMOVED***
	return string(newImgID), nil
***REMOVED***

func newImage(image *image.Image, size int64) *types.ImageSummary ***REMOVED***
	newImage := new(types.ImageSummary)
	newImage.ParentID = image.Parent.String()
	newImage.ID = image.ID().String()
	newImage.Created = image.Created.Unix()
	newImage.Size = size
	newImage.VirtualSize = size
	newImage.SharedSize = -1
	newImage.Containers = -1
	if image.Config != nil ***REMOVED***
		newImage.Labels = image.Config.Labels
	***REMOVED***
	return newImage
***REMOVED***
