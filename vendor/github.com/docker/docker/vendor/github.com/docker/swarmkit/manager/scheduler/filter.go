package scheduler

import (
	"fmt"
	"strings"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/genericresource"
	"github.com/docker/swarmkit/manager/constraint"
)

// Filter checks whether the given task can run on the given node.
// A filter may only operate
type Filter interface ***REMOVED***
	// SetTask returns true when the filter is enabled for a given task
	// and assigns the task to the filter. It returns false if the filter
	// isn't applicable to this task.  For instance, a constraints filter
	// would return `false` if the task doesn't contain any constraints.
	SetTask(*api.Task) bool

	// Check returns true if the task assigned by SetTask can be scheduled
	// into the given node. This function should not be called if SetTask
	// returned false.
	Check(*NodeInfo) bool

	// Explain what a failure of this filter means
	Explain(nodes int) string
***REMOVED***

// ReadyFilter checks that the node is ready to schedule tasks.
type ReadyFilter struct ***REMOVED***
***REMOVED***

// SetTask returns true when the filter is enabled for a given task.
func (f *ReadyFilter) SetTask(_ *api.Task) bool ***REMOVED***
	return true
***REMOVED***

// Check returns true if the task can be scheduled into the given node.
func (f *ReadyFilter) Check(n *NodeInfo) bool ***REMOVED***
	return n.Status.State == api.NodeStatus_READY &&
		n.Spec.Availability == api.NodeAvailabilityActive
***REMOVED***

// Explain returns an explanation of a failure.
func (f *ReadyFilter) Explain(nodes int) string ***REMOVED***
	if nodes == 1 ***REMOVED***
		return "1 node not available for new tasks"
	***REMOVED***
	return fmt.Sprintf("%d nodes not available for new tasks", nodes)
***REMOVED***

// ResourceFilter checks that the node has enough resources available to run
// the task.
type ResourceFilter struct ***REMOVED***
	reservations *api.Resources
***REMOVED***

// SetTask returns true when the filter is enabled for a given task.
func (f *ResourceFilter) SetTask(t *api.Task) bool ***REMOVED***
	r := t.Spec.Resources
	if r == nil || r.Reservations == nil ***REMOVED***
		return false
	***REMOVED***

	res := r.Reservations
	if res.NanoCPUs == 0 && res.MemoryBytes == 0 && len(res.Generic) == 0 ***REMOVED***
		return false
	***REMOVED***

	f.reservations = r.Reservations
	return true
***REMOVED***

// Check returns true if the task can be scheduled into the given node.
func (f *ResourceFilter) Check(n *NodeInfo) bool ***REMOVED***
	if f.reservations.NanoCPUs > n.AvailableResources.NanoCPUs ***REMOVED***
		return false
	***REMOVED***

	if f.reservations.MemoryBytes > n.AvailableResources.MemoryBytes ***REMOVED***
		return false
	***REMOVED***

	for _, v := range f.reservations.Generic ***REMOVED***
		enough, err := genericresource.HasEnough(n.AvailableResources.Generic, v)
		if err != nil || !enough ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

// Explain returns an explanation of a failure.
func (f *ResourceFilter) Explain(nodes int) string ***REMOVED***
	if nodes == 1 ***REMOVED***
		return "insufficient resources on 1 node"
	***REMOVED***
	return fmt.Sprintf("insufficient resources on %d nodes", nodes)
***REMOVED***

// PluginFilter checks that the node has a specific volume plugin installed
type PluginFilter struct ***REMOVED***
	t *api.Task
***REMOVED***

func referencesVolumePlugin(mount api.Mount) bool ***REMOVED***
	return mount.Type == api.MountTypeVolume &&
		mount.VolumeOptions != nil &&
		mount.VolumeOptions.DriverConfig != nil &&
		mount.VolumeOptions.DriverConfig.Name != "" &&
		mount.VolumeOptions.DriverConfig.Name != "local"

***REMOVED***

// SetTask returns true when the filter is enabled for a given task.
func (f *PluginFilter) SetTask(t *api.Task) bool ***REMOVED***
	c := t.Spec.GetContainer()

	var volumeTemplates bool
	if c != nil ***REMOVED***
		for _, mount := range c.Mounts ***REMOVED***
			if referencesVolumePlugin(mount) ***REMOVED***
				volumeTemplates = true
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if (c != nil && volumeTemplates) || len(t.Networks) > 0 || t.Spec.LogDriver != nil ***REMOVED***
		f.t = t
		return true
	***REMOVED***

	return false
***REMOVED***

// Check returns true if the task can be scheduled into the given node.
// TODO(amitshukla): investigate storing Plugins as a map so it can be easily probed
func (f *PluginFilter) Check(n *NodeInfo) bool ***REMOVED***
	if n.Description == nil || n.Description.Engine == nil ***REMOVED***
		// If the node is not running Engine, plugins are not
		// supported.
		return true
	***REMOVED***

	// Get list of plugins on the node
	nodePlugins := n.Description.Engine.Plugins

	// Check if all volume plugins required by task are installed on node
	container := f.t.Spec.GetContainer()
	if container != nil ***REMOVED***
		for _, mount := range container.Mounts ***REMOVED***
			if referencesVolumePlugin(mount) ***REMOVED***
				if _, exists := f.pluginExistsOnNode("Volume", mount.VolumeOptions.DriverConfig.Name, nodePlugins); !exists ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Check if all network plugins required by task are installed on node
	for _, tn := range f.t.Networks ***REMOVED***
		if tn.Network != nil && tn.Network.DriverState != nil && tn.Network.DriverState.Name != "" ***REMOVED***
			if _, exists := f.pluginExistsOnNode("Network", tn.Network.DriverState.Name, nodePlugins); !exists ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// It's possible that the LogDriver object does not carry a name, just some
	// configuration options. In that case, the plugin filter shouldn't fail to
	// schedule the task
	if f.t.Spec.LogDriver != nil && f.t.Spec.LogDriver.Name != "none" && f.t.Spec.LogDriver.Name != "" ***REMOVED***
		// If there are no log driver types in the list at all, most likely this is
		// an older daemon that did not report this information. In this case don't filter
		if typeFound, exists := f.pluginExistsOnNode("Log", f.t.Spec.LogDriver.Name, nodePlugins); !exists && typeFound ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// pluginExistsOnNode returns true if the (pluginName, pluginType) pair is present in nodePlugins
func (f *PluginFilter) pluginExistsOnNode(pluginType string, pluginName string, nodePlugins []api.PluginDescription) (bool, bool) ***REMOVED***
	var typeFound bool

	for _, np := range nodePlugins ***REMOVED***
		if pluginType != np.Type ***REMOVED***
			continue
		***REMOVED***
		typeFound = true

		if pluginName == np.Name ***REMOVED***
			return true, true
		***REMOVED***
		// This does not use the reference package to avoid the
		// overhead of parsing references as part of the scheduling
		// loop. This is okay only because plugin names are a very
		// strict subset of the reference grammar that is always
		// name:tag.
		if strings.HasPrefix(np.Name, pluginName) && np.Name[len(pluginName):] == ":latest" ***REMOVED***
			return true, true
		***REMOVED***
	***REMOVED***
	return typeFound, false
***REMOVED***

// Explain returns an explanation of a failure.
func (f *PluginFilter) Explain(nodes int) string ***REMOVED***
	if nodes == 1 ***REMOVED***
		return "missing plugin on 1 node"
	***REMOVED***
	return fmt.Sprintf("missing plugin on %d nodes", nodes)
***REMOVED***

// ConstraintFilter selects only nodes that match certain labels.
type ConstraintFilter struct ***REMOVED***
	constraints []constraint.Constraint
***REMOVED***

// SetTask returns true when the filter is enable for a given task.
func (f *ConstraintFilter) SetTask(t *api.Task) bool ***REMOVED***
	if t.Spec.Placement == nil || len(t.Spec.Placement.Constraints) == 0 ***REMOVED***
		return false
	***REMOVED***

	constraints, err := constraint.Parse(t.Spec.Placement.Constraints)
	if err != nil ***REMOVED***
		// constraints have been validated at controlapi
		// if in any case it finds an error here, treat this task
		// as constraint filter disabled.
		return false
	***REMOVED***
	f.constraints = constraints
	return true
***REMOVED***

// Check returns true if the task's constraint is supported by the given node.
func (f *ConstraintFilter) Check(n *NodeInfo) bool ***REMOVED***
	return constraint.NodeMatches(f.constraints, n.Node)
***REMOVED***

// Explain returns an explanation of a failure.
func (f *ConstraintFilter) Explain(nodes int) string ***REMOVED***
	if nodes == 1 ***REMOVED***
		return "scheduling constraints not satisfied on 1 node"
	***REMOVED***
	return fmt.Sprintf("scheduling constraints not satisfied on %d nodes", nodes)
***REMOVED***

// PlatformFilter selects only nodes that run the required platform.
type PlatformFilter struct ***REMOVED***
	supportedPlatforms []*api.Platform
***REMOVED***

// SetTask returns true when the filter is enabled for a given task.
func (f *PlatformFilter) SetTask(t *api.Task) bool ***REMOVED***
	placement := t.Spec.Placement
	if placement != nil ***REMOVED***
		// copy the platform information
		f.supportedPlatforms = placement.Platforms
		if len(placement.Platforms) > 0 ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// Check returns true if the task can be scheduled into the given node.
func (f *PlatformFilter) Check(n *NodeInfo) bool ***REMOVED***
	// if the supportedPlatforms field is empty, then either it wasn't
	// provided or there are no constraints
	if len(f.supportedPlatforms) == 0 ***REMOVED***
		return true
	***REMOVED***
	// check if the platform for the node is supported
	if n.Description != nil ***REMOVED***
		if nodePlatform := n.Description.Platform; nodePlatform != nil ***REMOVED***
			for _, p := range f.supportedPlatforms ***REMOVED***
				if f.platformEqual(*p, *nodePlatform) ***REMOVED***
					return true
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (f *PlatformFilter) platformEqual(imgPlatform, nodePlatform api.Platform) bool ***REMOVED***
	// normalize "x86_64" architectures to "amd64"
	if imgPlatform.Architecture == "x86_64" ***REMOVED***
		imgPlatform.Architecture = "amd64"
	***REMOVED***
	if nodePlatform.Architecture == "x86_64" ***REMOVED***
		nodePlatform.Architecture = "amd64"
	***REMOVED***

	// normalize "aarch64" architectures to "arm64"
	if imgPlatform.Architecture == "aarch64" ***REMOVED***
		imgPlatform.Architecture = "arm64"
	***REMOVED***
	if nodePlatform.Architecture == "aarch64" ***REMOVED***
		nodePlatform.Architecture = "arm64"
	***REMOVED***

	if (imgPlatform.Architecture == "" || imgPlatform.Architecture == nodePlatform.Architecture) && (imgPlatform.OS == "" || imgPlatform.OS == nodePlatform.OS) ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

// Explain returns an explanation of a failure.
func (f *PlatformFilter) Explain(nodes int) string ***REMOVED***
	if nodes == 1 ***REMOVED***
		return "unsupported platform on 1 node"
	***REMOVED***
	return fmt.Sprintf("unsupported platform on %d nodes", nodes)
***REMOVED***

// HostPortFilter checks that the node has a specific port available.
type HostPortFilter struct ***REMOVED***
	t *api.Task
***REMOVED***

// SetTask returns true when the filter is enabled for a given task.
func (f *HostPortFilter) SetTask(t *api.Task) bool ***REMOVED***
	if t.Endpoint != nil ***REMOVED***
		for _, port := range t.Endpoint.Ports ***REMOVED***
			if port.PublishMode == api.PublishModeHost && port.PublishedPort != 0 ***REMOVED***
				f.t = t
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

// Check returns true if the task can be scheduled into the given node.
func (f *HostPortFilter) Check(n *NodeInfo) bool ***REMOVED***
	for _, port := range f.t.Endpoint.Ports ***REMOVED***
		if port.PublishMode == api.PublishModeHost && port.PublishedPort != 0 ***REMOVED***
			portSpec := hostPortSpec***REMOVED***protocol: port.Protocol, publishedPort: port.PublishedPort***REMOVED***
			if _, ok := n.usedHostPorts[portSpec]; ok ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

// Explain returns an explanation of a failure.
func (f *HostPortFilter) Explain(nodes int) string ***REMOVED***
	if nodes == 1 ***REMOVED***
		return "host-mode port already in use on 1 node"
	***REMOVED***
	return fmt.Sprintf("host-mode port already in use on %d nodes", nodes)
***REMOVED***
