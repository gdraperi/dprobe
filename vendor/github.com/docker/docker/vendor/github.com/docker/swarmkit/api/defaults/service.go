package defaults

import (
	"time"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/deepcopy"
	gogotypes "github.com/gogo/protobuf/types"
)

// Service is a ServiceSpec object with all fields filled in using default
// values.
var Service = api.ServiceSpec***REMOVED***
	Task: api.TaskSpec***REMOVED***
		Runtime: &api.TaskSpec_Container***REMOVED***
			Container: &api.ContainerSpec***REMOVED***
				StopGracePeriod: gogotypes.DurationProto(10 * time.Second),
				PullOptions:     &api.ContainerSpec_PullOptions***REMOVED******REMOVED***,
				DNSConfig:       &api.ContainerSpec_DNSConfig***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Resources: &api.ResourceRequirements***REMOVED******REMOVED***,
		Restart: &api.RestartPolicy***REMOVED***
			Condition: api.RestartOnAny,
			Delay:     gogotypes.DurationProto(5 * time.Second),
		***REMOVED***,
		Placement: &api.Placement***REMOVED******REMOVED***,
	***REMOVED***,
	Update: &api.UpdateConfig***REMOVED***
		FailureAction: api.UpdateConfig_PAUSE,
		Monitor:       gogotypes.DurationProto(5 * time.Second),
		Parallelism:   1,
		Order:         api.UpdateConfig_STOP_FIRST,
	***REMOVED***,
	Rollback: &api.UpdateConfig***REMOVED***
		FailureAction: api.UpdateConfig_PAUSE,
		Monitor:       gogotypes.DurationProto(5 * time.Second),
		Parallelism:   1,
		Order:         api.UpdateConfig_STOP_FIRST,
	***REMOVED***,
***REMOVED***

// InterpolateService returns a ServiceSpec based on the provided spec, which
// has all unspecified values filled in with default values.
func InterpolateService(origSpec *api.ServiceSpec) *api.ServiceSpec ***REMOVED***
	spec := origSpec.Copy()

	container := spec.Task.GetContainer()
	defaultContainer := Service.Task.GetContainer()
	if container != nil ***REMOVED***
		if container.StopGracePeriod == nil ***REMOVED***
			container.StopGracePeriod = &gogotypes.Duration***REMOVED******REMOVED***
			deepcopy.Copy(container.StopGracePeriod, defaultContainer.StopGracePeriod)
		***REMOVED***
		if container.PullOptions == nil ***REMOVED***
			container.PullOptions = defaultContainer.PullOptions.Copy()
		***REMOVED***
		if container.DNSConfig == nil ***REMOVED***
			container.DNSConfig = defaultContainer.DNSConfig.Copy()
		***REMOVED***
	***REMOVED***

	if spec.Task.Resources == nil ***REMOVED***
		spec.Task.Resources = Service.Task.Resources.Copy()
	***REMOVED***

	if spec.Task.Restart == nil ***REMOVED***
		spec.Task.Restart = Service.Task.Restart.Copy()
	***REMOVED*** else ***REMOVED***
		if spec.Task.Restart.Delay == nil ***REMOVED***
			spec.Task.Restart.Delay = &gogotypes.Duration***REMOVED******REMOVED***
			deepcopy.Copy(spec.Task.Restart.Delay, Service.Task.Restart.Delay)
		***REMOVED***
	***REMOVED***

	if spec.Task.Placement == nil ***REMOVED***
		spec.Task.Placement = Service.Task.Placement.Copy()
	***REMOVED***

	if spec.Update == nil ***REMOVED***
		spec.Update = Service.Update.Copy()
	***REMOVED*** else ***REMOVED***
		if spec.Update.Monitor == nil ***REMOVED***
			spec.Update.Monitor = &gogotypes.Duration***REMOVED******REMOVED***
			deepcopy.Copy(spec.Update.Monitor, Service.Update.Monitor)
		***REMOVED***
	***REMOVED***

	if spec.Rollback == nil ***REMOVED***
		spec.Rollback = Service.Rollback.Copy()
	***REMOVED*** else ***REMOVED***
		if spec.Rollback.Monitor == nil ***REMOVED***
			spec.Rollback.Monitor = &gogotypes.Duration***REMOVED******REMOVED***
			deepcopy.Copy(spec.Rollback.Monitor, Service.Rollback.Monitor)
		***REMOVED***
	***REMOVED***

	return spec
***REMOVED***
