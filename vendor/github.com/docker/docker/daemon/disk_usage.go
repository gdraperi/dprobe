package daemon

import (
	"fmt"
	"sync/atomic"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/directory"
	"github.com/docker/docker/volume"
	"github.com/opencontainers/go-digest"
	"github.com/sirupsen/logrus"
)

func (daemon *Daemon) getLayerRefs() map[layer.ChainID]int ***REMOVED***
	tmpImages := daemon.imageStore.Map()
	layerRefs := map[layer.ChainID]int***REMOVED******REMOVED***
	for id, img := range tmpImages ***REMOVED***
		dgst := digest.Digest(id)
		if len(daemon.referenceStore.References(dgst)) == 0 && len(daemon.imageStore.Children(id)) != 0 ***REMOVED***
			continue
		***REMOVED***

		rootFS := *img.RootFS
		rootFS.DiffIDs = nil
		for _, id := range img.RootFS.DiffIDs ***REMOVED***
			rootFS.Append(id)
			chid := rootFS.ChainID()
			layerRefs[chid]++
		***REMOVED***
	***REMOVED***

	return layerRefs
***REMOVED***

// SystemDiskUsage returns information about the daemon data disk usage
func (daemon *Daemon) SystemDiskUsage(ctx context.Context) (*types.DiskUsage, error) ***REMOVED***
	if !atomic.CompareAndSwapInt32(&daemon.diskUsageRunning, 0, 1) ***REMOVED***
		return nil, fmt.Errorf("a disk usage operation is already running")
	***REMOVED***
	defer atomic.StoreInt32(&daemon.diskUsageRunning, 0)

	// Retrieve container list
	allContainers, err := daemon.Containers(&types.ContainerListOptions***REMOVED***
		Size: true,
		All:  true,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to retrieve container list: %v", err)
	***REMOVED***

	// Get all top images with extra attributes
	allImages, err := daemon.Images(filters.NewArgs(), false, true)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to retrieve image list: %v", err)
	***REMOVED***

	// Get all local volumes
	allVolumes := []*types.Volume***REMOVED******REMOVED***
	getLocalVols := func(v volume.Volume) error ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return ctx.Err()
		default:
			if d, ok := v.(volume.DetailedVolume); ok ***REMOVED***
				// skip local volumes with mount options since these could have external
				// mounted filesystems that will be slow to enumerate.
				if len(d.Options()) > 0 ***REMOVED***
					return nil
				***REMOVED***
			***REMOVED***
			name := v.Name()
			refs := daemon.volumes.Refs(v)

			tv := volumeToAPIType(v)
			sz, err := directory.Size(v.Path())
			if err != nil ***REMOVED***
				logrus.Warnf("failed to determine size of volume %v", name)
				sz = -1
			***REMOVED***
			tv.UsageData = &types.VolumeUsageData***REMOVED***Size: sz, RefCount: int64(len(refs))***REMOVED***
			allVolumes = append(allVolumes, tv)
		***REMOVED***

		return nil
	***REMOVED***

	err = daemon.traverseLocalVolumes(getLocalVols)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Get total layers size on disk
	var allLayersSize int64
	layerRefs := daemon.getLayerRefs()
	for _, ls := range daemon.layerStores ***REMOVED***
		allLayers := ls.Map()
		for _, l := range allLayers ***REMOVED***
			select ***REMOVED***
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				size, err := l.DiffSize()
				if err == nil ***REMOVED***
					if _, ok := layerRefs[l.ChainID()]; ok ***REMOVED***
						allLayersSize += size
					***REMOVED*** else ***REMOVED***
						logrus.Warnf("found leaked image layer %v", l.ChainID())
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					logrus.Warnf("failed to get diff size for layer %v", l.ChainID())
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return &types.DiskUsage***REMOVED***
		LayersSize: allLayersSize,
		Containers: allContainers,
		Volumes:    allVolumes,
		Images:     allImages,
	***REMOVED***, nil
***REMOVED***
