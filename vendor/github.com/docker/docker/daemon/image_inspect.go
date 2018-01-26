package daemon

import (
	"time"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/system"
	"github.com/pkg/errors"
)

// LookupImage looks up an image by name and returns it as an ImageInspect
// structure.
func (daemon *Daemon) LookupImage(name string) (*types.ImageInspect, error) ***REMOVED***
	img, err := daemon.GetImage(name)
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "no such image: %s", name)
	***REMOVED***
	if !system.IsOSSupported(img.OperatingSystem()) ***REMOVED***
		return nil, system.ErrNotSupportedOperatingSystem
	***REMOVED***
	refs := daemon.referenceStore.References(img.ID().Digest())
	repoTags := []string***REMOVED******REMOVED***
	repoDigests := []string***REMOVED******REMOVED***
	for _, ref := range refs ***REMOVED***
		switch ref.(type) ***REMOVED***
		case reference.NamedTagged:
			repoTags = append(repoTags, reference.FamiliarString(ref))
		case reference.Canonical:
			repoDigests = append(repoDigests, reference.FamiliarString(ref))
		***REMOVED***
	***REMOVED***

	var size int64
	var layerMetadata map[string]string
	layerID := img.RootFS.ChainID()
	if layerID != "" ***REMOVED***
		l, err := daemon.layerStores[img.OperatingSystem()].Get(layerID)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		defer layer.ReleaseAndLog(daemon.layerStores[img.OperatingSystem()], l)
		size, err = l.Size()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		layerMetadata, err = l.Metadata()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	comment := img.Comment
	if len(comment) == 0 && len(img.History) > 0 ***REMOVED***
		comment = img.History[len(img.History)-1].Comment
	***REMOVED***

	lastUpdated, err := daemon.imageStore.GetLastUpdated(img.ID())
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	imageInspect := &types.ImageInspect***REMOVED***
		ID:              img.ID().String(),
		RepoTags:        repoTags,
		RepoDigests:     repoDigests,
		Parent:          img.Parent.String(),
		Comment:         comment,
		Created:         img.Created.Format(time.RFC3339Nano),
		Container:       img.Container,
		ContainerConfig: &img.ContainerConfig,
		DockerVersion:   img.DockerVersion,
		Author:          img.Author,
		Config:          img.Config,
		Architecture:    img.Architecture,
		Os:              img.OperatingSystem(),
		OsVersion:       img.OSVersion,
		Size:            size,
		VirtualSize:     size, // TODO: field unused, deprecate
		RootFS:          rootFSToAPIType(img.RootFS),
		Metadata: types.ImageMetadata***REMOVED***
			LastTagTime: lastUpdated,
		***REMOVED***,
	***REMOVED***

	imageInspect.GraphDriver.Name = daemon.GraphDriverName(img.OperatingSystem())
	imageInspect.GraphDriver.Data = layerMetadata

	return imageInspect, nil
***REMOVED***
