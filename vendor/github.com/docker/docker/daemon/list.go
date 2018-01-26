package daemon

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/container"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/image"
	"github.com/docker/docker/volume"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var acceptedVolumeFilterTags = map[string]bool***REMOVED***
	"dangling": true,
	"name":     true,
	"driver":   true,
	"label":    true,
***REMOVED***

var acceptedPsFilterTags = map[string]bool***REMOVED***
	"ancestor":  true,
	"before":    true,
	"exited":    true,
	"id":        true,
	"isolation": true,
	"label":     true,
	"name":      true,
	"status":    true,
	"health":    true,
	"since":     true,
	"volume":    true,
	"network":   true,
	"is-task":   true,
	"publish":   true,
	"expose":    true,
***REMOVED***

// iterationAction represents possible outcomes happening during the container iteration.
type iterationAction int

// containerReducer represents a reducer for a container.
// Returns the object to serialize by the api.
type containerReducer func(*container.Snapshot, *listContext) (*types.Container, error)

const (
	// includeContainer is the action to include a container in the reducer.
	includeContainer iterationAction = iota
	// excludeContainer is the action to exclude a container in the reducer.
	excludeContainer
	// stopIteration is the action to stop iterating over the list of containers.
	stopIteration
)

// errStopIteration makes the iterator to stop without returning an error.
var errStopIteration = errors.New("container list iteration stopped")

// List returns an array of all containers registered in the daemon.
func (daemon *Daemon) List() []*container.Container ***REMOVED***
	return daemon.containers.List()
***REMOVED***

// listContext is the daemon generated filtering to iterate over containers.
// This is created based on the user specification from types.ContainerListOptions.
type listContext struct ***REMOVED***
	// idx is the container iteration index for this context
	idx int
	// ancestorFilter tells whether it should check ancestors or not
	ancestorFilter bool
	// names is a list of container names to filter with
	names map[string][]string
	// images is a list of images to filter with
	images map[image.ID]bool
	// filters is a collection of arguments to filter with, specified by the user
	filters filters.Args
	// exitAllowed is a list of exit codes allowed to filter with
	exitAllowed []int

	// beforeFilter is a filter to ignore containers that appear before the one given
	beforeFilter *container.Snapshot
	// sinceFilter is a filter to stop the filtering when the iterator arrive to the given container
	sinceFilter *container.Snapshot

	// taskFilter tells if we should filter based on wether a container is part of a task
	taskFilter bool
	// isTask tells us if the we should filter container that are a task (true) or not (false)
	isTask bool

	// publish is a list of published ports to filter with
	publish map[nat.Port]bool
	// expose is a list of exposed ports to filter with
	expose map[nat.Port]bool

	// ContainerListOptions is the filters set by the user
	*types.ContainerListOptions
***REMOVED***

// byCreatedDescending is a temporary type used to sort a list of containers by creation time.
type byCreatedDescending []container.Snapshot

func (r byCreatedDescending) Len() int      ***REMOVED*** return len(r) ***REMOVED***
func (r byCreatedDescending) Swap(i, j int) ***REMOVED*** r[i], r[j] = r[j], r[i] ***REMOVED***
func (r byCreatedDescending) Less(i, j int) bool ***REMOVED***
	return r[j].CreatedAt.UnixNano() < r[i].CreatedAt.UnixNano()
***REMOVED***

// Containers returns the list of containers to show given the user's filtering.
func (daemon *Daemon) Containers(config *types.ContainerListOptions) ([]*types.Container, error) ***REMOVED***
	return daemon.reduceContainers(config, daemon.refreshImage)
***REMOVED***

func (daemon *Daemon) filterByNameIDMatches(view container.View, ctx *listContext) ([]container.Snapshot, error) ***REMOVED***
	idSearch := false
	names := ctx.filters.Get("name")
	ids := ctx.filters.Get("id")
	if len(names)+len(ids) == 0 ***REMOVED***
		// if name or ID filters are not in use, return to
		// standard behavior of walking the entire container
		// list from the daemon's in-memory store
		all, err := view.All()
		sort.Sort(byCreatedDescending(all))
		return all, err
	***REMOVED***

	// idSearch will determine if we limit name matching to the IDs
	// matched from any IDs which were specified as filters
	if len(ids) > 0 ***REMOVED***
		idSearch = true
	***REMOVED***

	matches := make(map[string]bool)
	// find ID matches; errors represent "not found" and can be ignored
	for _, id := range ids ***REMOVED***
		if fullID, err := daemon.idIndex.Get(id); err == nil ***REMOVED***
			matches[fullID] = true
		***REMOVED***
	***REMOVED***

	// look for name matches; if ID filtering was used, then limit the
	// search space to the matches map only; errors represent "not found"
	// and can be ignored
	if len(names) > 0 ***REMOVED***
		for id, idNames := range ctx.names ***REMOVED***
			// if ID filters were used and no matches on that ID were
			// found, continue to next ID in the list
			if idSearch && !matches[id] ***REMOVED***
				continue
			***REMOVED***
			for _, eachName := range idNames ***REMOVED***
				if ctx.filters.Match("name", eachName) ***REMOVED***
					matches[id] = true
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	cntrs := make([]container.Snapshot, 0, len(matches))
	for id := range matches ***REMOVED***
		c, err := view.Get(id)
		switch err.(type) ***REMOVED***
		case nil:
			cntrs = append(cntrs, *c)
		case container.NoSuchContainerError:
			// ignore error
		default:
			return nil, err
		***REMOVED***
	***REMOVED***

	// Restore sort-order after filtering
	// Created gives us nanosec resolution for sorting
	sort.Sort(byCreatedDescending(cntrs))

	return cntrs, nil
***REMOVED***

// reduceContainers parses the user's filtering options and generates the list of containers to return based on a reducer.
func (daemon *Daemon) reduceContainers(config *types.ContainerListOptions, reducer containerReducer) ([]*types.Container, error) ***REMOVED***
	var (
		view       = daemon.containersReplica.Snapshot()
		containers = []*types.Container***REMOVED******REMOVED***
	)

	ctx, err := daemon.foldFilter(view, config)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// fastpath to only look at a subset of containers if specific name
	// or ID matches were provided by the user--otherwise we potentially
	// end up querying many more containers than intended
	containerList, err := daemon.filterByNameIDMatches(view, ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for i := range containerList ***REMOVED***
		t, err := daemon.reducePsContainer(&containerList[i], ctx, reducer)
		if err != nil ***REMOVED***
			if err != errStopIteration ***REMOVED***
				return nil, err
			***REMOVED***
			break
		***REMOVED***
		if t != nil ***REMOVED***
			containers = append(containers, t)
			ctx.idx++
		***REMOVED***
	***REMOVED***

	return containers, nil
***REMOVED***

// reducePsContainer is the basic representation for a container as expected by the ps command.
func (daemon *Daemon) reducePsContainer(container *container.Snapshot, ctx *listContext, reducer containerReducer) (*types.Container, error) ***REMOVED***
	// filter containers to return
	switch includeContainerInList(container, ctx) ***REMOVED***
	case excludeContainer:
		return nil, nil
	case stopIteration:
		return nil, errStopIteration
	***REMOVED***

	// transform internal container struct into api structs
	newC, err := reducer(container, ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// release lock because size calculation is slow
	if ctx.Size ***REMOVED***
		sizeRw, sizeRootFs := daemon.getSize(newC.ID)
		newC.SizeRw = sizeRw
		newC.SizeRootFs = sizeRootFs
	***REMOVED***
	return newC, nil
***REMOVED***

// foldFilter generates the container filter based on the user's filtering options.
func (daemon *Daemon) foldFilter(view container.View, config *types.ContainerListOptions) (*listContext, error) ***REMOVED***
	psFilters := config.Filters

	if err := psFilters.Validate(acceptedPsFilterTags); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var filtExited []int

	err := psFilters.WalkValues("exited", func(value string) error ***REMOVED***
		code, err := strconv.Atoi(value)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		filtExited = append(filtExited, code)
		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err = psFilters.WalkValues("status", func(value string) error ***REMOVED***
		if !container.IsValidStateString(value) ***REMOVED***
			return invalidFilter***REMOVED***"status", value***REMOVED***
		***REMOVED***

		config.All = true
		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var taskFilter, isTask bool
	if psFilters.Contains("is-task") ***REMOVED***
		if psFilters.ExactMatch("is-task", "true") ***REMOVED***
			taskFilter = true
			isTask = true
		***REMOVED*** else if psFilters.ExactMatch("is-task", "false") ***REMOVED***
			taskFilter = true
			isTask = false
		***REMOVED*** else ***REMOVED***
			return nil, invalidFilter***REMOVED***"is-task", psFilters.Get("is-task")***REMOVED***
		***REMOVED***
	***REMOVED***

	err = psFilters.WalkValues("health", func(value string) error ***REMOVED***
		if !container.IsValidHealthString(value) ***REMOVED***
			return errdefs.InvalidParameter(errors.Errorf("Unrecognised filter value for health: %s", value))
		***REMOVED***

		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var beforeContFilter, sinceContFilter *container.Snapshot

	err = psFilters.WalkValues("before", func(value string) error ***REMOVED***
		beforeContFilter, err = idOrNameFilter(view, value)
		return err
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err = psFilters.WalkValues("since", func(value string) error ***REMOVED***
		sinceContFilter, err = idOrNameFilter(view, value)
		return err
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	imagesFilter := map[image.ID]bool***REMOVED******REMOVED***
	var ancestorFilter bool
	if psFilters.Contains("ancestor") ***REMOVED***
		ancestorFilter = true
		psFilters.WalkValues("ancestor", func(ancestor string) error ***REMOVED***
			id, _, err := daemon.GetImageIDAndOS(ancestor)
			if err != nil ***REMOVED***
				logrus.Warnf("Error while looking up for image %v", ancestor)
				return nil
			***REMOVED***
			if imagesFilter[id] ***REMOVED***
				// Already seen this ancestor, skip it
				return nil
			***REMOVED***
			// Then walk down the graph and put the imageIds in imagesFilter
			populateImageFilterByParents(imagesFilter, id, daemon.imageStore.Children)
			return nil
		***REMOVED***)
	***REMOVED***

	publishFilter := map[nat.Port]bool***REMOVED******REMOVED***
	err = psFilters.WalkValues("publish", portOp("publish", publishFilter))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	exposeFilter := map[nat.Port]bool***REMOVED******REMOVED***
	err = psFilters.WalkValues("expose", portOp("expose", exposeFilter))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &listContext***REMOVED***
		filters:              psFilters,
		ancestorFilter:       ancestorFilter,
		images:               imagesFilter,
		exitAllowed:          filtExited,
		beforeFilter:         beforeContFilter,
		sinceFilter:          sinceContFilter,
		taskFilter:           taskFilter,
		isTask:               isTask,
		publish:              publishFilter,
		expose:               exposeFilter,
		ContainerListOptions: config,
		names:                view.GetAllNames(),
	***REMOVED***, nil
***REMOVED***

func idOrNameFilter(view container.View, value string) (*container.Snapshot, error) ***REMOVED***
	filter, err := view.Get(value)
	switch err.(type) ***REMOVED***
	case container.NoSuchContainerError:
		// Try name search instead
		found := ""
		for id, idNames := range view.GetAllNames() ***REMOVED***
			for _, eachName := range idNames ***REMOVED***
				if strings.TrimPrefix(value, "/") == strings.TrimPrefix(eachName, "/") ***REMOVED***
					if found != "" && found != id ***REMOVED***
						return nil, err
					***REMOVED***
					found = id
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if found != "" ***REMOVED***
			filter, err = view.Get(found)
		***REMOVED***
	***REMOVED***
	return filter, err
***REMOVED***

func portOp(key string, filter map[nat.Port]bool) func(value string) error ***REMOVED***
	return func(value string) error ***REMOVED***
		if strings.Contains(value, ":") ***REMOVED***
			return fmt.Errorf("filter for '%s' should not contain ':': %s", key, value)
		***REMOVED***
		//support two formats, original format <portnum>/[<proto>] or <startport-endport>/[<proto>]
		proto, port := nat.SplitProtoPort(value)
		start, end, err := nat.ParsePortRange(port)
		if err != nil ***REMOVED***
			return fmt.Errorf("error while looking up for %s %s: %s", key, value, err)
		***REMOVED***
		for i := start; i <= end; i++ ***REMOVED***
			p, err := nat.NewPort(proto, strconv.FormatUint(i, 10))
			if err != nil ***REMOVED***
				return fmt.Errorf("error while looking up for %s %s: %s", key, value, err)
			***REMOVED***
			filter[p] = true
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// includeContainerInList decides whether a container should be included in the output or not based in the filter.
// It also decides if the iteration should be stopped or not.
func includeContainerInList(container *container.Snapshot, ctx *listContext) iterationAction ***REMOVED***
	// Do not include container if it's in the list before the filter container.
	// Set the filter container to nil to include the rest of containers after this one.
	if ctx.beforeFilter != nil ***REMOVED***
		if container.ID == ctx.beforeFilter.ID ***REMOVED***
			ctx.beforeFilter = nil
		***REMOVED***
		return excludeContainer
	***REMOVED***

	// Stop iteration when the container arrives to the filter container
	if ctx.sinceFilter != nil ***REMOVED***
		if container.ID == ctx.sinceFilter.ID ***REMOVED***
			return stopIteration
		***REMOVED***
	***REMOVED***

	// Do not include container if it's stopped and we're not filters
	if !container.Running && !ctx.All && ctx.Limit <= 0 ***REMOVED***
		return excludeContainer
	***REMOVED***

	// Do not include container if the name doesn't match
	if !ctx.filters.Match("name", container.Name) ***REMOVED***
		return excludeContainer
	***REMOVED***

	// Do not include container if the id doesn't match
	if !ctx.filters.Match("id", container.ID) ***REMOVED***
		return excludeContainer
	***REMOVED***

	if ctx.taskFilter ***REMOVED***
		if ctx.isTask != container.Managed ***REMOVED***
			return excludeContainer
		***REMOVED***
	***REMOVED***

	// Do not include container if any of the labels don't match
	if !ctx.filters.MatchKVList("label", container.Labels) ***REMOVED***
		return excludeContainer
	***REMOVED***

	// Do not include container if isolation doesn't match
	if excludeContainer == excludeByIsolation(container, ctx) ***REMOVED***
		return excludeContainer
	***REMOVED***

	// Stop iteration when the index is over the limit
	if ctx.Limit > 0 && ctx.idx == ctx.Limit ***REMOVED***
		return stopIteration
	***REMOVED***

	// Do not include container if its exit code is not in the filter
	if len(ctx.exitAllowed) > 0 ***REMOVED***
		shouldSkip := true
		for _, code := range ctx.exitAllowed ***REMOVED***
			if code == container.ExitCode && !container.Running && !container.StartedAt.IsZero() ***REMOVED***
				shouldSkip = false
				break
			***REMOVED***
		***REMOVED***
		if shouldSkip ***REMOVED***
			return excludeContainer
		***REMOVED***
	***REMOVED***

	// Do not include container if its status doesn't match the filter
	if !ctx.filters.Match("status", container.State) ***REMOVED***
		return excludeContainer
	***REMOVED***

	// Do not include container if its health doesn't match the filter
	if !ctx.filters.ExactMatch("health", container.Health) ***REMOVED***
		return excludeContainer
	***REMOVED***

	if ctx.filters.Contains("volume") ***REMOVED***
		volumesByName := make(map[string]types.MountPoint)
		for _, m := range container.Mounts ***REMOVED***
			if m.Name != "" ***REMOVED***
				volumesByName[m.Name] = m
			***REMOVED*** else ***REMOVED***
				volumesByName[m.Source] = m
			***REMOVED***
		***REMOVED***
		volumesByDestination := make(map[string]types.MountPoint)
		for _, m := range container.Mounts ***REMOVED***
			if m.Destination != "" ***REMOVED***
				volumesByDestination[m.Destination] = m
			***REMOVED***
		***REMOVED***

		volumeExist := fmt.Errorf("volume mounted in container")
		err := ctx.filters.WalkValues("volume", func(value string) error ***REMOVED***
			if _, exist := volumesByDestination[value]; exist ***REMOVED***
				return volumeExist
			***REMOVED***
			if _, exist := volumesByName[value]; exist ***REMOVED***
				return volumeExist
			***REMOVED***
			return nil
		***REMOVED***)
		if err != volumeExist ***REMOVED***
			return excludeContainer
		***REMOVED***
	***REMOVED***

	if ctx.ancestorFilter ***REMOVED***
		if len(ctx.images) == 0 ***REMOVED***
			return excludeContainer
		***REMOVED***
		if !ctx.images[image.ID(container.ImageID)] ***REMOVED***
			return excludeContainer
		***REMOVED***
	***REMOVED***

	var (
		networkExist = errors.New("container part of network")
		noNetworks   = errors.New("container is not part of any networks")
	)
	if ctx.filters.Contains("network") ***REMOVED***
		err := ctx.filters.WalkValues("network", func(value string) error ***REMOVED***
			if container.NetworkSettings == nil ***REMOVED***
				return noNetworks
			***REMOVED***
			if _, ok := container.NetworkSettings.Networks[value]; ok ***REMOVED***
				return networkExist
			***REMOVED***
			for _, nw := range container.NetworkSettings.Networks ***REMOVED***
				if nw == nil ***REMOVED***
					continue
				***REMOVED***
				if strings.HasPrefix(nw.NetworkID, value) ***REMOVED***
					return networkExist
				***REMOVED***
			***REMOVED***
			return nil
		***REMOVED***)
		if err != networkExist ***REMOVED***
			return excludeContainer
		***REMOVED***
	***REMOVED***

	if len(ctx.publish) > 0 ***REMOVED***
		shouldSkip := true
		for port := range ctx.publish ***REMOVED***
			if _, ok := container.PortBindings[port]; ok ***REMOVED***
				shouldSkip = false
				break
			***REMOVED***
		***REMOVED***
		if shouldSkip ***REMOVED***
			return excludeContainer
		***REMOVED***
	***REMOVED***

	if len(ctx.expose) > 0 ***REMOVED***
		shouldSkip := true
		for port := range ctx.expose ***REMOVED***
			if _, ok := container.ExposedPorts[port]; ok ***REMOVED***
				shouldSkip = false
				break
			***REMOVED***
		***REMOVED***
		if shouldSkip ***REMOVED***
			return excludeContainer
		***REMOVED***
	***REMOVED***

	return includeContainer
***REMOVED***

// refreshImage checks if the Image ref still points to the correct ID, and updates the ref to the actual ID when it doesn't
func (daemon *Daemon) refreshImage(s *container.Snapshot, ctx *listContext) (*types.Container, error) ***REMOVED***
	c := s.Container
	image := s.Image // keep the original ref if still valid (hasn't changed)
	if image != s.ImageID ***REMOVED***
		id, _, err := daemon.GetImageIDAndOS(image)
		if _, isDNE := err.(errImageDoesNotExist); err != nil && !isDNE ***REMOVED***
			return nil, err
		***REMOVED***
		if err != nil || id.String() != s.ImageID ***REMOVED***
			// ref changed, we need to use original ID
			image = s.ImageID
		***REMOVED***
	***REMOVED***
	c.Image = image
	return &c, nil
***REMOVED***

// Volumes lists known volumes, using the filter to restrict the range
// of volumes returned.
func (daemon *Daemon) Volumes(filter string) ([]*types.Volume, []string, error) ***REMOVED***
	var (
		volumesOut []*types.Volume
	)
	volFilters, err := filters.FromJSON(filter)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	if err := volFilters.Validate(acceptedVolumeFilterTags); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	volumes, warnings, err := daemon.volumes.List()
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	filterVolumes, err := daemon.filterVolumes(volumes, volFilters)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	for _, v := range filterVolumes ***REMOVED***
		apiV := volumeToAPIType(v)
		if vv, ok := v.(interface ***REMOVED***
			CachedPath() string
		***REMOVED***); ok ***REMOVED***
			apiV.Mountpoint = vv.CachedPath()
		***REMOVED*** else ***REMOVED***
			apiV.Mountpoint = v.Path()
		***REMOVED***
		volumesOut = append(volumesOut, apiV)
	***REMOVED***
	return volumesOut, warnings, nil
***REMOVED***

// filterVolumes filters volume list according to user specified filter
// and returns user chosen volumes
func (daemon *Daemon) filterVolumes(vols []volume.Volume, filter filters.Args) ([]volume.Volume, error) ***REMOVED***
	// if filter is empty, return original volume list
	if filter.Len() == 0 ***REMOVED***
		return vols, nil
	***REMOVED***

	var retVols []volume.Volume
	for _, vol := range vols ***REMOVED***
		if filter.Contains("name") ***REMOVED***
			if !filter.Match("name", vol.Name()) ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
		if filter.Contains("driver") ***REMOVED***
			if !filter.ExactMatch("driver", vol.DriverName()) ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
		if filter.Contains("label") ***REMOVED***
			v, ok := vol.(volume.DetailedVolume)
			if !ok ***REMOVED***
				continue
			***REMOVED***
			if !filter.MatchKVList("label", v.Labels()) ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
		retVols = append(retVols, vol)
	***REMOVED***
	danglingOnly := false
	if filter.Contains("dangling") ***REMOVED***
		if filter.ExactMatch("dangling", "true") || filter.ExactMatch("dangling", "1") ***REMOVED***
			danglingOnly = true
		***REMOVED*** else if !filter.ExactMatch("dangling", "false") && !filter.ExactMatch("dangling", "0") ***REMOVED***
			return nil, invalidFilter***REMOVED***"dangling", filter.Get("dangling")***REMOVED***
		***REMOVED***
		retVols = daemon.volumes.FilterByUsed(retVols, !danglingOnly)
	***REMOVED***
	return retVols, nil
***REMOVED***

func populateImageFilterByParents(ancestorMap map[image.ID]bool, imageID image.ID, getChildren func(image.ID) []image.ID) ***REMOVED***
	if !ancestorMap[imageID] ***REMOVED***
		for _, id := range getChildren(imageID) ***REMOVED***
			populateImageFilterByParents(ancestorMap, id, getChildren)
		***REMOVED***
		ancestorMap[imageID] = true
	***REMOVED***
***REMOVED***
