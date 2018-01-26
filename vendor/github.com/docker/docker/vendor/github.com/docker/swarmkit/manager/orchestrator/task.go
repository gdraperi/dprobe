package orchestrator

import (
	"reflect"
	"time"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/defaults"
	"github.com/docker/swarmkit/identity"
	"github.com/docker/swarmkit/protobuf/ptypes"
)

// NewTask creates a new task.
func NewTask(cluster *api.Cluster, service *api.Service, slot uint64, nodeID string) *api.Task ***REMOVED***
	var logDriver *api.Driver
	if service.Spec.Task.LogDriver != nil ***REMOVED***
		// use the log driver specific to the task, if we have it.
		logDriver = service.Spec.Task.LogDriver
	***REMOVED*** else if cluster != nil ***REMOVED***
		// pick up the cluster default, if available.
		logDriver = cluster.Spec.TaskDefaults.LogDriver // nil is okay here.
	***REMOVED***

	taskID := identity.NewID()
	task := api.Task***REMOVED***
		ID:                 taskID,
		ServiceAnnotations: service.Spec.Annotations,
		Spec:               service.Spec.Task,
		SpecVersion:        service.SpecVersion,
		ServiceID:          service.ID,
		Slot:               slot,
		Status: api.TaskStatus***REMOVED***
			State:     api.TaskStateNew,
			Timestamp: ptypes.MustTimestampProto(time.Now()),
			Message:   "created",
		***REMOVED***,
		Endpoint: &api.Endpoint***REMOVED***
			Spec: service.Spec.Endpoint.Copy(),
		***REMOVED***,
		DesiredState: api.TaskStateRunning,
		LogDriver:    logDriver,
	***REMOVED***

	// In global mode we also set the NodeID
	if nodeID != "" ***REMOVED***
		task.NodeID = nodeID
	***REMOVED***

	return &task
***REMOVED***

// RestartCondition returns the restart condition to apply to this task.
func RestartCondition(task *api.Task) api.RestartPolicy_RestartCondition ***REMOVED***
	restartCondition := defaults.Service.Task.Restart.Condition
	if task.Spec.Restart != nil ***REMOVED***
		restartCondition = task.Spec.Restart.Condition
	***REMOVED***
	return restartCondition
***REMOVED***

// IsTaskDirty determines whether a task matches the given service's spec.
func IsTaskDirty(s *api.Service, t *api.Task) bool ***REMOVED***
	// If the spec version matches, we know the task is not dirty. However,
	// if it does not match, that doesn't mean the task is dirty, since
	// only a portion of the spec is included in the comparison.
	if t.SpecVersion != nil && s.SpecVersion != nil && *s.SpecVersion == *t.SpecVersion ***REMOVED***
		return false
	***REMOVED***

	// Make a deep copy of the service and task spec for the comparison.
	serviceTaskSpec := *s.Spec.Task.Copy()

	// For non-failed tasks with a container spec runtime that have already
	// pulled the required image (i.e., current state is between READY and
	// RUNNING inclusively), ignore the value of the `PullOptions` field by
	// setting the copied service to have the same PullOptions value as the
	// task. A difference in only the `PullOptions` field should not cause
	// a running (or ready to run) task to be considered 'dirty' when we
	// handle updates.
	// See https://github.com/docker/swarmkit/issues/971
	currentState := t.Status.State
	// Ignore PullOpts if the task is desired to be in a "runnable" state
	// and its last known current state is between READY and RUNNING in
	// which case we know that the task either successfully pulled its
	// container image or didn't need to.
	ignorePullOpts := t.DesiredState <= api.TaskStateRunning && currentState >= api.TaskStateReady && currentState <= api.TaskStateRunning
	if ignorePullOpts && serviceTaskSpec.GetContainer() != nil && t.Spec.GetContainer() != nil ***REMOVED***
		// Modify the service's container spec.
		serviceTaskSpec.GetContainer().PullOptions = t.Spec.GetContainer().PullOptions
	***REMOVED***

	return !reflect.DeepEqual(serviceTaskSpec, t.Spec) ||
		(t.Endpoint != nil && !reflect.DeepEqual(s.Spec.Endpoint, t.Endpoint.Spec))
***REMOVED***

// InvalidNode is true if the node is nil, down, or drained
func InvalidNode(n *api.Node) bool ***REMOVED***
	return n == nil ||
		n.Status.State == api.NodeStatus_DOWN ||
		n.Spec.Availability == api.NodeAvailabilityDrain
***REMOVED***

// TasksByTimestamp sorts tasks by applied timestamp if available, otherwise
// status timestamp.
type TasksByTimestamp []*api.Task

// Len implements the Len method for sorting.
func (t TasksByTimestamp) Len() int ***REMOVED***
	return len(t)
***REMOVED***

// Swap implements the Swap method for sorting.
func (t TasksByTimestamp) Swap(i, j int) ***REMOVED***
	t[i], t[j] = t[j], t[i]
***REMOVED***

// Less implements the Less method for sorting.
func (t TasksByTimestamp) Less(i, j int) bool ***REMOVED***
	iTimestamp := t[i].Status.Timestamp
	if t[i].Status.AppliedAt != nil ***REMOVED***
		iTimestamp = t[i].Status.AppliedAt
	***REMOVED***

	jTimestamp := t[j].Status.Timestamp
	if t[j].Status.AppliedAt != nil ***REMOVED***
		iTimestamp = t[j].Status.AppliedAt
	***REMOVED***

	if iTimestamp == nil ***REMOVED***
		return true
	***REMOVED***
	if jTimestamp == nil ***REMOVED***
		return false
	***REMOVED***
	if iTimestamp.Seconds < jTimestamp.Seconds ***REMOVED***
		return true
	***REMOVED***
	if iTimestamp.Seconds > jTimestamp.Seconds ***REMOVED***
		return false
	***REMOVED***
	return iTimestamp.Nanos < jTimestamp.Nanos
***REMOVED***
