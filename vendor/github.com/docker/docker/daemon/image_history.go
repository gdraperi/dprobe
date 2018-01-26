package daemon

import (
	"fmt"
	"time"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/layer"
)

// ImageHistory returns a slice of ImageHistory structures for the specified image
// name by walking the image lineage.
func (daemon *Daemon) ImageHistory(name string) ([]*image.HistoryResponseItem, error) ***REMOVED***
	start := time.Now()
	img, err := daemon.GetImage(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	history := []*image.HistoryResponseItem***REMOVED******REMOVED***

	layerCounter := 0
	rootFS := *img.RootFS
	rootFS.DiffIDs = nil

	for _, h := range img.History ***REMOVED***
		var layerSize int64

		if !h.EmptyLayer ***REMOVED***
			if len(img.RootFS.DiffIDs) <= layerCounter ***REMOVED***
				return nil, fmt.Errorf("too many non-empty layers in History section")
			***REMOVED***

			rootFS.Append(img.RootFS.DiffIDs[layerCounter])
			l, err := daemon.layerStores[img.OperatingSystem()].Get(rootFS.ChainID())
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			layerSize, err = l.DiffSize()
			layer.ReleaseAndLog(daemon.layerStores[img.OperatingSystem()], l)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			layerCounter++
		***REMOVED***

		history = append([]*image.HistoryResponseItem***REMOVED******REMOVED***
			ID:        "<missing>",
			Created:   h.Created.Unix(),
			CreatedBy: h.CreatedBy,
			Comment:   h.Comment,
			Size:      layerSize,
		***REMOVED******REMOVED***, history...)
	***REMOVED***

	// Fill in image IDs and tags
	histImg := img
	id := img.ID()
	for _, h := range history ***REMOVED***
		h.ID = id.String()

		var tags []string
		for _, r := range daemon.referenceStore.References(id.Digest()) ***REMOVED***
			if _, ok := r.(reference.NamedTagged); ok ***REMOVED***
				tags = append(tags, reference.FamiliarString(r))
			***REMOVED***
		***REMOVED***

		h.Tags = tags

		id = histImg.Parent
		if id == "" ***REMOVED***
			break
		***REMOVED***
		histImg, err = daemon.GetImage(id.String())
		if err != nil ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	imageActions.WithValues("history").UpdateSince(start)
	return history, nil
***REMOVED***
