package daemon

import (
	"fmt"
	"regexp"
	"sync/atomic"
	"time"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	timetypes "github.com/docker/docker/api/types/time"
	"github.com/docker/docker/image"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/directory"
	"github.com/docker/docker/runconfig"
	"github.com/docker/docker/volume"
	"github.com/docker/libnetwork"
	digest "github.com/opencontainers/go-digest"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

var (
	// errPruneRunning is returned when a prune request is received while
	// one is in progress
	errPruneRunning = fmt.Errorf("a prune operation is already running")

	containersAcceptedFilters = map[string]bool***REMOVED***
		"label":  true,
		"label!": true,
		"until":  true,
	***REMOVED***
	volumesAcceptedFilters = map[string]bool***REMOVED***
		"label":  true,
		"label!": true,
	***REMOVED***
	imagesAcceptedFilters = map[string]bool***REMOVED***
		"dangling": true,
		"label":    true,
		"label!":   true,
		"until":    true,
	***REMOVED***
	networksAcceptedFilters = map[string]bool***REMOVED***
		"label":  true,
		"label!": true,
		"until":  true,
	***REMOVED***
)

// ContainersPrune removes unused containers
func (daemon *Daemon) ContainersPrune(ctx context.Context, pruneFilters filters.Args) (*types.ContainersPruneReport, error) ***REMOVED***
	if !atomic.CompareAndSwapInt32(&daemon.pruneRunning, 0, 1) ***REMOVED***
		return nil, errPruneRunning
	***REMOVED***
	defer atomic.StoreInt32(&daemon.pruneRunning, 0)

	rep := &types.ContainersPruneReport***REMOVED******REMOVED***

	// make sure that only accepted filters have been received
	err := pruneFilters.Validate(containersAcceptedFilters)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	until, err := getUntilFromPruneFilters(pruneFilters)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	allContainers := daemon.List()
	for _, c := range allContainers ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			logrus.Debugf("ContainersPrune operation cancelled: %#v", *rep)
			return rep, nil
		default:
		***REMOVED***

		if !c.IsRunning() ***REMOVED***
			if !until.IsZero() && c.Created.After(until) ***REMOVED***
				continue
			***REMOVED***
			if !matchLabels(pruneFilters, c.Config.Labels) ***REMOVED***
				continue
			***REMOVED***
			cSize, _ := daemon.getSize(c.ID)
			// TODO: sets RmLink to true?
			err := daemon.ContainerRm(c.ID, &types.ContainerRmConfig***REMOVED******REMOVED***)
			if err != nil ***REMOVED***
				logrus.Warnf("failed to prune container %s: %v", c.ID, err)
				continue
			***REMOVED***
			if cSize > 0 ***REMOVED***
				rep.SpaceReclaimed += uint64(cSize)
			***REMOVED***
			rep.ContainersDeleted = append(rep.ContainersDeleted, c.ID)
		***REMOVED***
	***REMOVED***

	return rep, nil
***REMOVED***

// VolumesPrune removes unused local volumes
func (daemon *Daemon) VolumesPrune(ctx context.Context, pruneFilters filters.Args) (*types.VolumesPruneReport, error) ***REMOVED***
	if !atomic.CompareAndSwapInt32(&daemon.pruneRunning, 0, 1) ***REMOVED***
		return nil, errPruneRunning
	***REMOVED***
	defer atomic.StoreInt32(&daemon.pruneRunning, 0)

	// make sure that only accepted filters have been received
	err := pruneFilters.Validate(volumesAcceptedFilters)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	rep := &types.VolumesPruneReport***REMOVED******REMOVED***

	pruneVols := func(v volume.Volume) error ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			logrus.Debugf("VolumesPrune operation cancelled: %#v", *rep)
			return ctx.Err()
		default:
		***REMOVED***

		name := v.Name()
		refs := daemon.volumes.Refs(v)

		if len(refs) == 0 ***REMOVED***
			detailedVolume, ok := v.(volume.DetailedVolume)
			if ok ***REMOVED***
				if !matchLabels(pruneFilters, detailedVolume.Labels()) ***REMOVED***
					return nil
				***REMOVED***
			***REMOVED***
			vSize, err := directory.Size(v.Path())
			if err != nil ***REMOVED***
				logrus.Warnf("could not determine size of volume %s: %v", name, err)
			***REMOVED***
			err = daemon.volumeRm(v)
			if err != nil ***REMOVED***
				logrus.Warnf("could not remove volume %s: %v", name, err)
				return nil
			***REMOVED***
			rep.SpaceReclaimed += uint64(vSize)
			rep.VolumesDeleted = append(rep.VolumesDeleted, name)
		***REMOVED***

		return nil
	***REMOVED***

	err = daemon.traverseLocalVolumes(pruneVols)
	if err == context.Canceled ***REMOVED***
		return rep, nil
	***REMOVED***

	return rep, err
***REMOVED***

// ImagesPrune removes unused images
func (daemon *Daemon) ImagesPrune(ctx context.Context, pruneFilters filters.Args) (*types.ImagesPruneReport, error) ***REMOVED***
	if !atomic.CompareAndSwapInt32(&daemon.pruneRunning, 0, 1) ***REMOVED***
		return nil, errPruneRunning
	***REMOVED***
	defer atomic.StoreInt32(&daemon.pruneRunning, 0)

	// make sure that only accepted filters have been received
	err := pruneFilters.Validate(imagesAcceptedFilters)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	rep := &types.ImagesPruneReport***REMOVED******REMOVED***

	danglingOnly := true
	if pruneFilters.Contains("dangling") ***REMOVED***
		if pruneFilters.ExactMatch("dangling", "false") || pruneFilters.ExactMatch("dangling", "0") ***REMOVED***
			danglingOnly = false
		***REMOVED*** else if !pruneFilters.ExactMatch("dangling", "true") && !pruneFilters.ExactMatch("dangling", "1") ***REMOVED***
			return nil, invalidFilter***REMOVED***"dangling", pruneFilters.Get("dangling")***REMOVED***
		***REMOVED***
	***REMOVED***

	until, err := getUntilFromPruneFilters(pruneFilters)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var allImages map[image.ID]*image.Image
	if danglingOnly ***REMOVED***
		allImages = daemon.imageStore.Heads()
	***REMOVED*** else ***REMOVED***
		allImages = daemon.imageStore.Map()
	***REMOVED***
	allContainers := daemon.List()
	imageRefs := map[string]bool***REMOVED******REMOVED***
	for _, c := range allContainers ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			imageRefs[c.ID] = true
		***REMOVED***
	***REMOVED***

	// Filter intermediary images and get their unique size
	allLayers := make(map[layer.ChainID]layer.Layer)
	for _, ls := range daemon.layerStores ***REMOVED***
		for k, v := range ls.Map() ***REMOVED***
			allLayers[k] = v
		***REMOVED***
	***REMOVED***
	topImages := map[image.ID]*image.Image***REMOVED******REMOVED***
	for id, img := range allImages ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			dgst := digest.Digest(id)
			if len(daemon.referenceStore.References(dgst)) == 0 && len(daemon.imageStore.Children(id)) != 0 ***REMOVED***
				continue
			***REMOVED***
			if !until.IsZero() && img.Created.After(until) ***REMOVED***
				continue
			***REMOVED***
			if img.Config != nil && !matchLabels(pruneFilters, img.Config.Labels) ***REMOVED***
				continue
			***REMOVED***
			topImages[id] = img
		***REMOVED***
	***REMOVED***

	canceled := false
deleteImagesLoop:
	for id := range topImages ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			// we still want to calculate freed size and return the data
			canceled = true
			break deleteImagesLoop
		default:
		***REMOVED***

		dgst := digest.Digest(id)
		hex := dgst.Hex()
		if _, ok := imageRefs[hex]; ok ***REMOVED***
			continue
		***REMOVED***

		deletedImages := []types.ImageDeleteResponseItem***REMOVED******REMOVED***
		refs := daemon.referenceStore.References(dgst)
		if len(refs) > 0 ***REMOVED***
			shouldDelete := !danglingOnly
			if !shouldDelete ***REMOVED***
				hasTag := false
				for _, ref := range refs ***REMOVED***
					if _, ok := ref.(reference.NamedTagged); ok ***REMOVED***
						hasTag = true
						break
					***REMOVED***
				***REMOVED***

				// Only delete if it's untagged (i.e. repo:<none>)
				shouldDelete = !hasTag
			***REMOVED***

			if shouldDelete ***REMOVED***
				for _, ref := range refs ***REMOVED***
					imgDel, err := daemon.ImageDelete(ref.String(), false, true)
					if err != nil ***REMOVED***
						logrus.Warnf("could not delete reference %s: %v", ref.String(), err)
						continue
					***REMOVED***
					deletedImages = append(deletedImages, imgDel...)
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			imgDel, err := daemon.ImageDelete(hex, false, true)
			if err != nil ***REMOVED***
				logrus.Warnf("could not delete image %s: %v", hex, err)
				continue
			***REMOVED***
			deletedImages = append(deletedImages, imgDel...)
		***REMOVED***

		rep.ImagesDeleted = append(rep.ImagesDeleted, deletedImages...)
	***REMOVED***

	// Compute how much space was freed
	for _, d := range rep.ImagesDeleted ***REMOVED***
		if d.Deleted != "" ***REMOVED***
			chid := layer.ChainID(d.Deleted)
			if l, ok := allLayers[chid]; ok ***REMOVED***
				diffSize, err := l.DiffSize()
				if err != nil ***REMOVED***
					logrus.Warnf("failed to get layer %s size: %v", chid, err)
					continue
				***REMOVED***
				rep.SpaceReclaimed += uint64(diffSize)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if canceled ***REMOVED***
		logrus.Debugf("ImagesPrune operation cancelled: %#v", *rep)
	***REMOVED***

	return rep, nil
***REMOVED***

// localNetworksPrune removes unused local networks
func (daemon *Daemon) localNetworksPrune(ctx context.Context, pruneFilters filters.Args) *types.NetworksPruneReport ***REMOVED***
	rep := &types.NetworksPruneReport***REMOVED******REMOVED***

	until, _ := getUntilFromPruneFilters(pruneFilters)

	// When the function returns true, the walk will stop.
	l := func(nw libnetwork.Network) bool ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			// context cancelled
			return true
		default:
		***REMOVED***
		if nw.Info().ConfigOnly() ***REMOVED***
			return false
		***REMOVED***
		if !until.IsZero() && nw.Info().Created().After(until) ***REMOVED***
			return false
		***REMOVED***
		if !matchLabels(pruneFilters, nw.Info().Labels()) ***REMOVED***
			return false
		***REMOVED***
		nwName := nw.Name()
		if runconfig.IsPreDefinedNetwork(nwName) ***REMOVED***
			return false
		***REMOVED***
		if len(nw.Endpoints()) > 0 ***REMOVED***
			return false
		***REMOVED***
		if err := daemon.DeleteNetwork(nw.ID()); err != nil ***REMOVED***
			logrus.Warnf("could not remove local network %s: %v", nwName, err)
			return false
		***REMOVED***
		rep.NetworksDeleted = append(rep.NetworksDeleted, nwName)
		return false
	***REMOVED***
	daemon.netController.WalkNetworks(l)
	return rep
***REMOVED***

// clusterNetworksPrune removes unused cluster networks
func (daemon *Daemon) clusterNetworksPrune(ctx context.Context, pruneFilters filters.Args) (*types.NetworksPruneReport, error) ***REMOVED***
	rep := &types.NetworksPruneReport***REMOVED******REMOVED***

	until, _ := getUntilFromPruneFilters(pruneFilters)

	cluster := daemon.GetCluster()

	if !cluster.IsManager() ***REMOVED***
		return rep, nil
	***REMOVED***

	networks, err := cluster.GetNetworks()
	if err != nil ***REMOVED***
		return rep, err
	***REMOVED***
	networkIsInUse := regexp.MustCompile(`network ([[:alnum:]]+) is in use`)
	for _, nw := range networks ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return rep, nil
		default:
			if nw.Ingress ***REMOVED***
				// Routing-mesh network removal has to be explicitly invoked by user
				continue
			***REMOVED***
			if !until.IsZero() && nw.Created.After(until) ***REMOVED***
				continue
			***REMOVED***
			if !matchLabels(pruneFilters, nw.Labels) ***REMOVED***
				continue
			***REMOVED***
			// https://github.com/docker/docker/issues/24186
			// `docker network inspect` unfortunately displays ONLY those containers that are local to that node.
			// So we try to remove it anyway and check the error
			err = cluster.RemoveNetwork(nw.ID)
			if err != nil ***REMOVED***
				// we can safely ignore the "network .. is in use" error
				match := networkIsInUse.FindStringSubmatch(err.Error())
				if len(match) != 2 || match[1] != nw.ID ***REMOVED***
					logrus.Warnf("could not remove cluster network %s: %v", nw.Name, err)
				***REMOVED***
				continue
			***REMOVED***
			rep.NetworksDeleted = append(rep.NetworksDeleted, nw.Name)
		***REMOVED***
	***REMOVED***
	return rep, nil
***REMOVED***

// NetworksPrune removes unused networks
func (daemon *Daemon) NetworksPrune(ctx context.Context, pruneFilters filters.Args) (*types.NetworksPruneReport, error) ***REMOVED***
	if !atomic.CompareAndSwapInt32(&daemon.pruneRunning, 0, 1) ***REMOVED***
		return nil, errPruneRunning
	***REMOVED***
	defer atomic.StoreInt32(&daemon.pruneRunning, 0)

	// make sure that only accepted filters have been received
	err := pruneFilters.Validate(networksAcceptedFilters)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if _, err := getUntilFromPruneFilters(pruneFilters); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	rep := &types.NetworksPruneReport***REMOVED******REMOVED***
	if clusterRep, err := daemon.clusterNetworksPrune(ctx, pruneFilters); err == nil ***REMOVED***
		rep.NetworksDeleted = append(rep.NetworksDeleted, clusterRep.NetworksDeleted...)
	***REMOVED***

	localRep := daemon.localNetworksPrune(ctx, pruneFilters)
	rep.NetworksDeleted = append(rep.NetworksDeleted, localRep.NetworksDeleted...)

	select ***REMOVED***
	case <-ctx.Done():
		logrus.Debugf("NetworksPrune operation cancelled: %#v", *rep)
		return rep, nil
	default:
	***REMOVED***

	return rep, nil
***REMOVED***

func getUntilFromPruneFilters(pruneFilters filters.Args) (time.Time, error) ***REMOVED***
	until := time.Time***REMOVED******REMOVED***
	if !pruneFilters.Contains("until") ***REMOVED***
		return until, nil
	***REMOVED***
	untilFilters := pruneFilters.Get("until")
	if len(untilFilters) > 1 ***REMOVED***
		return until, fmt.Errorf("more than one until filter specified")
	***REMOVED***
	ts, err := timetypes.GetTimestamp(untilFilters[0], time.Now())
	if err != nil ***REMOVED***
		return until, err
	***REMOVED***
	seconds, nanoseconds, err := timetypes.ParseTimestamps(ts, 0)
	if err != nil ***REMOVED***
		return until, err
	***REMOVED***
	until = time.Unix(seconds, nanoseconds)
	return until, nil
***REMOVED***

func matchLabels(pruneFilters filters.Args, labels map[string]string) bool ***REMOVED***
	if !pruneFilters.MatchKVList("label", labels) ***REMOVED***
		return false
	***REMOVED***
	// By default MatchKVList will return true if field (like 'label!') does not exist
	// So we have to add additional Contains("label!") check
	if pruneFilters.Contains("label!") ***REMOVED***
		if pruneFilters.MatchKVList("label!", labels) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***
