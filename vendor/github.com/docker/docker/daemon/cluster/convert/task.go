package convert

import (
	"strings"

	types "github.com/docker/docker/api/types/swarm"
	swarmapi "github.com/docker/swarmkit/api"
	gogotypes "github.com/gogo/protobuf/types"
)

// TaskFromGRPC converts a grpc Task to a Task.
func TaskFromGRPC(t swarmapi.Task) (types.Task, error) ***REMOVED***
	if t.Spec.GetAttachment() != nil ***REMOVED***
		return types.Task***REMOVED******REMOVED***, nil
	***REMOVED***
	containerStatus := t.Status.GetContainer()
	taskSpec, err := taskSpecFromGRPC(t.Spec)
	if err != nil ***REMOVED***
		return types.Task***REMOVED******REMOVED***, err
	***REMOVED***
	task := types.Task***REMOVED***
		ID:          t.ID,
		Annotations: annotationsFromGRPC(t.Annotations),
		ServiceID:   t.ServiceID,
		Slot:        int(t.Slot),
		NodeID:      t.NodeID,
		Spec:        taskSpec,
		Status: types.TaskStatus***REMOVED***
			State:   types.TaskState(strings.ToLower(t.Status.State.String())),
			Message: t.Status.Message,
			Err:     t.Status.Err,
		***REMOVED***,
		DesiredState:     types.TaskState(strings.ToLower(t.DesiredState.String())),
		GenericResources: GenericResourcesFromGRPC(t.AssignedGenericResources),
	***REMOVED***

	// Meta
	task.Version.Index = t.Meta.Version.Index
	task.CreatedAt, _ = gogotypes.TimestampFromProto(t.Meta.CreatedAt)
	task.UpdatedAt, _ = gogotypes.TimestampFromProto(t.Meta.UpdatedAt)

	task.Status.Timestamp, _ = gogotypes.TimestampFromProto(t.Status.Timestamp)

	if containerStatus != nil ***REMOVED***
		task.Status.ContainerStatus.ContainerID = containerStatus.ContainerID
		task.Status.ContainerStatus.PID = int(containerStatus.PID)
		task.Status.ContainerStatus.ExitCode = int(containerStatus.ExitCode)
	***REMOVED***

	// NetworksAttachments
	for _, na := range t.Networks ***REMOVED***
		task.NetworksAttachments = append(task.NetworksAttachments, networkAttachmentFromGRPC(na))
	***REMOVED***

	if t.Status.PortStatus == nil ***REMOVED***
		return task, nil
	***REMOVED***

	for _, p := range t.Status.PortStatus.Ports ***REMOVED***
		task.Status.PortStatus.Ports = append(task.Status.PortStatus.Ports, types.PortConfig***REMOVED***
			Name:          p.Name,
			Protocol:      types.PortConfigProtocol(strings.ToLower(swarmapi.PortConfig_Protocol_name[int32(p.Protocol)])),
			PublishMode:   types.PortConfigPublishMode(strings.ToLower(swarmapi.PortConfig_PublishMode_name[int32(p.PublishMode)])),
			TargetPort:    p.TargetPort,
			PublishedPort: p.PublishedPort,
		***REMOVED***)
	***REMOVED***

	return task, nil
***REMOVED***
