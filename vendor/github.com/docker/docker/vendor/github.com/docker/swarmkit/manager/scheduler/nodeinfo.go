package scheduler

import (
	"time"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/genericresource"
	"github.com/docker/swarmkit/log"
	"golang.org/x/net/context"
)

// hostPortSpec specifies a used host port.
type hostPortSpec struct ***REMOVED***
	protocol      api.PortConfig_Protocol
	publishedPort uint32
***REMOVED***

// versionedService defines a tuple that contains a service ID and a spec
// version, so that failures can be tracked per spec version. Note that if the
// task predates spec versioning, specVersion will contain the zero value, and
// this will still work correctly.
type versionedService struct ***REMOVED***
	serviceID   string
	specVersion api.Version
***REMOVED***

// NodeInfo contains a node and some additional metadata.
type NodeInfo struct ***REMOVED***
	*api.Node
	Tasks                     map[string]*api.Task
	ActiveTasksCount          int
	ActiveTasksCountByService map[string]int
	AvailableResources        *api.Resources
	usedHostPorts             map[hostPortSpec]struct***REMOVED******REMOVED***

	// recentFailures is a map from service ID/version to the timestamps of
	// the most recent failures the node has experienced from replicas of
	// that service.
	recentFailures map[versionedService][]time.Time

	// lastCleanup is the last time recentFailures was cleaned up. This is
	// done periodically to avoid recentFailures growing without any limit.
	lastCleanup time.Time
***REMOVED***

func newNodeInfo(n *api.Node, tasks map[string]*api.Task, availableResources api.Resources) NodeInfo ***REMOVED***
	nodeInfo := NodeInfo***REMOVED***
		Node:  n,
		Tasks: make(map[string]*api.Task),
		ActiveTasksCountByService: make(map[string]int),
		AvailableResources:        availableResources.Copy(),
		usedHostPorts:             make(map[hostPortSpec]struct***REMOVED******REMOVED***),
		recentFailures:            make(map[versionedService][]time.Time),
		lastCleanup:               time.Now(),
	***REMOVED***

	for _, t := range tasks ***REMOVED***
		nodeInfo.addTask(t)
	***REMOVED***

	return nodeInfo
***REMOVED***

// removeTask removes a task from nodeInfo if it's tracked there, and returns true
// if nodeInfo was modified.
func (nodeInfo *NodeInfo) removeTask(t *api.Task) bool ***REMOVED***
	oldTask, ok := nodeInfo.Tasks[t.ID]
	if !ok ***REMOVED***
		return false
	***REMOVED***

	delete(nodeInfo.Tasks, t.ID)
	if oldTask.DesiredState <= api.TaskStateRunning ***REMOVED***
		nodeInfo.ActiveTasksCount--
		nodeInfo.ActiveTasksCountByService[t.ServiceID]--
	***REMOVED***

	if t.Endpoint != nil ***REMOVED***
		for _, port := range t.Endpoint.Ports ***REMOVED***
			if port.PublishMode == api.PublishModeHost && port.PublishedPort != 0 ***REMOVED***
				portSpec := hostPortSpec***REMOVED***protocol: port.Protocol, publishedPort: port.PublishedPort***REMOVED***
				delete(nodeInfo.usedHostPorts, portSpec)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	reservations := taskReservations(t.Spec)
	resources := nodeInfo.AvailableResources

	resources.MemoryBytes += reservations.MemoryBytes
	resources.NanoCPUs += reservations.NanoCPUs

	if nodeInfo.Description == nil || nodeInfo.Description.Resources == nil ||
		nodeInfo.Description.Resources.Generic == nil ***REMOVED***
		return true
	***REMOVED***

	taskAssigned := t.AssignedGenericResources
	nodeAvailableResources := &resources.Generic
	nodeRes := nodeInfo.Description.Resources.Generic
	genericresource.Reclaim(nodeAvailableResources, taskAssigned, nodeRes)

	return true
***REMOVED***

// addTask adds or updates a task on nodeInfo, and returns true if nodeInfo was
// modified.
func (nodeInfo *NodeInfo) addTask(t *api.Task) bool ***REMOVED***
	oldTask, ok := nodeInfo.Tasks[t.ID]
	if ok ***REMOVED***
		if t.DesiredState <= api.TaskStateRunning && oldTask.DesiredState > api.TaskStateRunning ***REMOVED***
			nodeInfo.Tasks[t.ID] = t
			nodeInfo.ActiveTasksCount++
			nodeInfo.ActiveTasksCountByService[t.ServiceID]++
			return true
		***REMOVED*** else if t.DesiredState > api.TaskStateRunning && oldTask.DesiredState <= api.TaskStateRunning ***REMOVED***
			nodeInfo.Tasks[t.ID] = t
			nodeInfo.ActiveTasksCount--
			nodeInfo.ActiveTasksCountByService[t.ServiceID]--
			return true
		***REMOVED***
		return false
	***REMOVED***

	nodeInfo.Tasks[t.ID] = t

	reservations := taskReservations(t.Spec)
	resources := nodeInfo.AvailableResources

	resources.MemoryBytes -= reservations.MemoryBytes
	resources.NanoCPUs -= reservations.NanoCPUs

	// minimum size required
	t.AssignedGenericResources = make([]*api.GenericResource, 0, len(resources.Generic))
	taskAssigned := &t.AssignedGenericResources

	genericresource.Claim(&resources.Generic, taskAssigned, reservations.Generic)

	if t.Endpoint != nil ***REMOVED***
		for _, port := range t.Endpoint.Ports ***REMOVED***
			if port.PublishMode == api.PublishModeHost && port.PublishedPort != 0 ***REMOVED***
				portSpec := hostPortSpec***REMOVED***protocol: port.Protocol, publishedPort: port.PublishedPort***REMOVED***
				nodeInfo.usedHostPorts[portSpec] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if t.DesiredState <= api.TaskStateRunning ***REMOVED***
		nodeInfo.ActiveTasksCount++
		nodeInfo.ActiveTasksCountByService[t.ServiceID]++
	***REMOVED***

	return true
***REMOVED***

func taskReservations(spec api.TaskSpec) (reservations api.Resources) ***REMOVED***
	if spec.Resources != nil && spec.Resources.Reservations != nil ***REMOVED***
		reservations = *spec.Resources.Reservations
	***REMOVED***
	return
***REMOVED***

func (nodeInfo *NodeInfo) cleanupFailures(now time.Time) ***REMOVED***
entriesLoop:
	for key, failuresEntry := range nodeInfo.recentFailures ***REMOVED***
		for _, timestamp := range failuresEntry ***REMOVED***
			if now.Sub(timestamp) < monitorFailures ***REMOVED***
				continue entriesLoop
			***REMOVED***
		***REMOVED***
		delete(nodeInfo.recentFailures, key)
	***REMOVED***
	nodeInfo.lastCleanup = now
***REMOVED***

// taskFailed records a task failure from a given service.
func (nodeInfo *NodeInfo) taskFailed(ctx context.Context, t *api.Task) ***REMOVED***
	expired := 0
	now := time.Now()

	if now.Sub(nodeInfo.lastCleanup) >= monitorFailures ***REMOVED***
		nodeInfo.cleanupFailures(now)
	***REMOVED***

	versionedService := versionedService***REMOVED***serviceID: t.ServiceID***REMOVED***
	if t.SpecVersion != nil ***REMOVED***
		versionedService.specVersion = *t.SpecVersion
	***REMOVED***

	for _, timestamp := range nodeInfo.recentFailures[versionedService] ***REMOVED***
		if now.Sub(timestamp) < monitorFailures ***REMOVED***
			break
		***REMOVED***
		expired++
	***REMOVED***

	if len(nodeInfo.recentFailures[versionedService])-expired == maxFailures-1 ***REMOVED***
		log.G(ctx).Warnf("underweighting node %s for service %s because it experienced %d failures or rejections within %s", nodeInfo.ID, t.ServiceID, maxFailures, monitorFailures.String())
	***REMOVED***

	nodeInfo.recentFailures[versionedService] = append(nodeInfo.recentFailures[versionedService][expired:], now)
***REMOVED***

// countRecentFailures returns the number of times the service has failed on
// this node within the lookback window monitorFailures.
func (nodeInfo *NodeInfo) countRecentFailures(now time.Time, t *api.Task) int ***REMOVED***
	versionedService := versionedService***REMOVED***serviceID: t.ServiceID***REMOVED***
	if t.SpecVersion != nil ***REMOVED***
		versionedService.specVersion = *t.SpecVersion
	***REMOVED***

	recentFailureCount := len(nodeInfo.recentFailures[versionedService])
	for i := recentFailureCount - 1; i >= 0; i-- ***REMOVED***
		if now.Sub(nodeInfo.recentFailures[versionedService][i]) > monitorFailures ***REMOVED***
			recentFailureCount -= i + 1
			break
		***REMOVED***
	***REMOVED***

	return recentFailureCount
***REMOVED***
