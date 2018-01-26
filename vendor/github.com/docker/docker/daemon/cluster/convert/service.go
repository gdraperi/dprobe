package convert

import (
	"fmt"
	"strings"

	types "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/api/types/swarm/runtime"
	"github.com/docker/docker/pkg/namesgenerator"
	swarmapi "github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/genericresource"
	"github.com/gogo/protobuf/proto"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
)

var (
	// ErrUnsupportedRuntime returns an error if the runtime is not supported by the daemon
	ErrUnsupportedRuntime = errors.New("unsupported runtime")
)

// ServiceFromGRPC converts a grpc Service to a Service.
func ServiceFromGRPC(s swarmapi.Service) (types.Service, error) ***REMOVED***
	curSpec, err := serviceSpecFromGRPC(&s.Spec)
	if err != nil ***REMOVED***
		return types.Service***REMOVED******REMOVED***, err
	***REMOVED***
	prevSpec, err := serviceSpecFromGRPC(s.PreviousSpec)
	if err != nil ***REMOVED***
		return types.Service***REMOVED******REMOVED***, err
	***REMOVED***
	service := types.Service***REMOVED***
		ID:           s.ID,
		Spec:         *curSpec,
		PreviousSpec: prevSpec,

		Endpoint: endpointFromGRPC(s.Endpoint),
	***REMOVED***

	// Meta
	service.Version.Index = s.Meta.Version.Index
	service.CreatedAt, _ = gogotypes.TimestampFromProto(s.Meta.CreatedAt)
	service.UpdatedAt, _ = gogotypes.TimestampFromProto(s.Meta.UpdatedAt)

	// UpdateStatus
	if s.UpdateStatus != nil ***REMOVED***
		service.UpdateStatus = &types.UpdateStatus***REMOVED******REMOVED***
		switch s.UpdateStatus.State ***REMOVED***
		case swarmapi.UpdateStatus_UPDATING:
			service.UpdateStatus.State = types.UpdateStateUpdating
		case swarmapi.UpdateStatus_PAUSED:
			service.UpdateStatus.State = types.UpdateStatePaused
		case swarmapi.UpdateStatus_COMPLETED:
			service.UpdateStatus.State = types.UpdateStateCompleted
		case swarmapi.UpdateStatus_ROLLBACK_STARTED:
			service.UpdateStatus.State = types.UpdateStateRollbackStarted
		case swarmapi.UpdateStatus_ROLLBACK_PAUSED:
			service.UpdateStatus.State = types.UpdateStateRollbackPaused
		case swarmapi.UpdateStatus_ROLLBACK_COMPLETED:
			service.UpdateStatus.State = types.UpdateStateRollbackCompleted
		***REMOVED***

		startedAt, _ := gogotypes.TimestampFromProto(s.UpdateStatus.StartedAt)
		if !startedAt.IsZero() && startedAt.Unix() != 0 ***REMOVED***
			service.UpdateStatus.StartedAt = &startedAt
		***REMOVED***

		completedAt, _ := gogotypes.TimestampFromProto(s.UpdateStatus.CompletedAt)
		if !completedAt.IsZero() && completedAt.Unix() != 0 ***REMOVED***
			service.UpdateStatus.CompletedAt = &completedAt
		***REMOVED***

		service.UpdateStatus.Message = s.UpdateStatus.Message
	***REMOVED***

	return service, nil
***REMOVED***

func serviceSpecFromGRPC(spec *swarmapi.ServiceSpec) (*types.ServiceSpec, error) ***REMOVED***
	if spec == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	serviceNetworks := make([]types.NetworkAttachmentConfig, 0, len(spec.Networks))
	for _, n := range spec.Networks ***REMOVED***
		netConfig := types.NetworkAttachmentConfig***REMOVED***Target: n.Target, Aliases: n.Aliases, DriverOpts: n.DriverAttachmentOpts***REMOVED***
		serviceNetworks = append(serviceNetworks, netConfig)

	***REMOVED***

	taskTemplate, err := taskSpecFromGRPC(spec.Task)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	switch t := spec.Task.GetRuntime().(type) ***REMOVED***
	case *swarmapi.TaskSpec_Container:
		containerConfig := t.Container
		taskTemplate.ContainerSpec = containerSpecFromGRPC(containerConfig)
		taskTemplate.Runtime = types.RuntimeContainer
	case *swarmapi.TaskSpec_Generic:
		switch t.Generic.Kind ***REMOVED***
		case string(types.RuntimePlugin):
			taskTemplate.Runtime = types.RuntimePlugin
		default:
			return nil, fmt.Errorf("unknown task runtime type: %s", t.Generic.Payload.TypeUrl)
		***REMOVED***

	default:
		return nil, fmt.Errorf("error creating service; unsupported runtime %T", t)
	***REMOVED***

	convertedSpec := &types.ServiceSpec***REMOVED***
		Annotations:  annotationsFromGRPC(spec.Annotations),
		TaskTemplate: taskTemplate,
		Networks:     serviceNetworks,
		EndpointSpec: endpointSpecFromGRPC(spec.Endpoint),
	***REMOVED***

	// UpdateConfig
	convertedSpec.UpdateConfig = updateConfigFromGRPC(spec.Update)
	convertedSpec.RollbackConfig = updateConfigFromGRPC(spec.Rollback)

	// Mode
	switch t := spec.GetMode().(type) ***REMOVED***
	case *swarmapi.ServiceSpec_Global:
		convertedSpec.Mode.Global = &types.GlobalService***REMOVED******REMOVED***
	case *swarmapi.ServiceSpec_Replicated:
		convertedSpec.Mode.Replicated = &types.ReplicatedService***REMOVED***
			Replicas: &t.Replicated.Replicas,
		***REMOVED***
	***REMOVED***

	return convertedSpec, nil
***REMOVED***

// ServiceSpecToGRPC converts a ServiceSpec to a grpc ServiceSpec.
func ServiceSpecToGRPC(s types.ServiceSpec) (swarmapi.ServiceSpec, error) ***REMOVED***
	name := s.Name
	if name == "" ***REMOVED***
		name = namesgenerator.GetRandomName(0)
	***REMOVED***

	serviceNetworks := make([]*swarmapi.NetworkAttachmentConfig, 0, len(s.Networks))
	for _, n := range s.Networks ***REMOVED***
		netConfig := &swarmapi.NetworkAttachmentConfig***REMOVED***Target: n.Target, Aliases: n.Aliases, DriverAttachmentOpts: n.DriverOpts***REMOVED***
		serviceNetworks = append(serviceNetworks, netConfig)
	***REMOVED***

	taskNetworks := make([]*swarmapi.NetworkAttachmentConfig, 0, len(s.TaskTemplate.Networks))
	for _, n := range s.TaskTemplate.Networks ***REMOVED***
		netConfig := &swarmapi.NetworkAttachmentConfig***REMOVED***Target: n.Target, Aliases: n.Aliases, DriverAttachmentOpts: n.DriverOpts***REMOVED***
		taskNetworks = append(taskNetworks, netConfig)

	***REMOVED***

	spec := swarmapi.ServiceSpec***REMOVED***
		Annotations: swarmapi.Annotations***REMOVED***
			Name:   name,
			Labels: s.Labels,
		***REMOVED***,
		Task: swarmapi.TaskSpec***REMOVED***
			Resources:   resourcesToGRPC(s.TaskTemplate.Resources),
			LogDriver:   driverToGRPC(s.TaskTemplate.LogDriver),
			Networks:    taskNetworks,
			ForceUpdate: s.TaskTemplate.ForceUpdate,
		***REMOVED***,
		Networks: serviceNetworks,
	***REMOVED***

	switch s.TaskTemplate.Runtime ***REMOVED***
	case types.RuntimeContainer, "": // if empty runtime default to container
		if s.TaskTemplate.ContainerSpec != nil ***REMOVED***
			containerSpec, err := containerToGRPC(s.TaskTemplate.ContainerSpec)
			if err != nil ***REMOVED***
				return swarmapi.ServiceSpec***REMOVED******REMOVED***, err
			***REMOVED***
			spec.Task.Runtime = &swarmapi.TaskSpec_Container***REMOVED***Container: containerSpec***REMOVED***
		***REMOVED***
	case types.RuntimePlugin:
		if s.Mode.Replicated != nil ***REMOVED***
			return swarmapi.ServiceSpec***REMOVED******REMOVED***, errors.New("plugins must not use replicated mode")
		***REMOVED***

		s.Mode.Global = &types.GlobalService***REMOVED******REMOVED*** // must always be global

		if s.TaskTemplate.PluginSpec != nil ***REMOVED***
			pluginSpec, err := proto.Marshal(s.TaskTemplate.PluginSpec)
			if err != nil ***REMOVED***
				return swarmapi.ServiceSpec***REMOVED******REMOVED***, err
			***REMOVED***
			spec.Task.Runtime = &swarmapi.TaskSpec_Generic***REMOVED***
				Generic: &swarmapi.GenericRuntimeSpec***REMOVED***
					Kind: string(types.RuntimePlugin),
					Payload: &gogotypes.Any***REMOVED***
						TypeUrl: string(types.RuntimeURLPlugin),
						Value:   pluginSpec,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***
		***REMOVED***
	default:
		return swarmapi.ServiceSpec***REMOVED******REMOVED***, ErrUnsupportedRuntime
	***REMOVED***

	restartPolicy, err := restartPolicyToGRPC(s.TaskTemplate.RestartPolicy)
	if err != nil ***REMOVED***
		return swarmapi.ServiceSpec***REMOVED******REMOVED***, err
	***REMOVED***
	spec.Task.Restart = restartPolicy

	if s.TaskTemplate.Placement != nil ***REMOVED***
		var preferences []*swarmapi.PlacementPreference
		for _, pref := range s.TaskTemplate.Placement.Preferences ***REMOVED***
			if pref.Spread != nil ***REMOVED***
				preferences = append(preferences, &swarmapi.PlacementPreference***REMOVED***
					Preference: &swarmapi.PlacementPreference_Spread***REMOVED***
						Spread: &swarmapi.SpreadOver***REMOVED***
							SpreadDescriptor: pref.Spread.SpreadDescriptor,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***)
			***REMOVED***
		***REMOVED***
		var platforms []*swarmapi.Platform
		for _, plat := range s.TaskTemplate.Placement.Platforms ***REMOVED***
			platforms = append(platforms, &swarmapi.Platform***REMOVED***
				Architecture: plat.Architecture,
				OS:           plat.OS,
			***REMOVED***)
		***REMOVED***
		spec.Task.Placement = &swarmapi.Placement***REMOVED***
			Constraints: s.TaskTemplate.Placement.Constraints,
			Preferences: preferences,
			Platforms:   platforms,
		***REMOVED***
	***REMOVED***

	spec.Update, err = updateConfigToGRPC(s.UpdateConfig)
	if err != nil ***REMOVED***
		return swarmapi.ServiceSpec***REMOVED******REMOVED***, err
	***REMOVED***
	spec.Rollback, err = updateConfigToGRPC(s.RollbackConfig)
	if err != nil ***REMOVED***
		return swarmapi.ServiceSpec***REMOVED******REMOVED***, err
	***REMOVED***

	if s.EndpointSpec != nil ***REMOVED***
		if s.EndpointSpec.Mode != "" &&
			s.EndpointSpec.Mode != types.ResolutionModeVIP &&
			s.EndpointSpec.Mode != types.ResolutionModeDNSRR ***REMOVED***
			return swarmapi.ServiceSpec***REMOVED******REMOVED***, fmt.Errorf("invalid resolution mode: %q", s.EndpointSpec.Mode)
		***REMOVED***

		spec.Endpoint = &swarmapi.EndpointSpec***REMOVED******REMOVED***

		spec.Endpoint.Mode = swarmapi.EndpointSpec_ResolutionMode(swarmapi.EndpointSpec_ResolutionMode_value[strings.ToUpper(string(s.EndpointSpec.Mode))])

		for _, portConfig := range s.EndpointSpec.Ports ***REMOVED***
			spec.Endpoint.Ports = append(spec.Endpoint.Ports, &swarmapi.PortConfig***REMOVED***
				Name:          portConfig.Name,
				Protocol:      swarmapi.PortConfig_Protocol(swarmapi.PortConfig_Protocol_value[strings.ToUpper(string(portConfig.Protocol))]),
				PublishMode:   swarmapi.PortConfig_PublishMode(swarmapi.PortConfig_PublishMode_value[strings.ToUpper(string(portConfig.PublishMode))]),
				TargetPort:    portConfig.TargetPort,
				PublishedPort: portConfig.PublishedPort,
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	// Mode
	if s.Mode.Global != nil && s.Mode.Replicated != nil ***REMOVED***
		return swarmapi.ServiceSpec***REMOVED******REMOVED***, fmt.Errorf("cannot specify both replicated mode and global mode")
	***REMOVED***

	if s.Mode.Global != nil ***REMOVED***
		spec.Mode = &swarmapi.ServiceSpec_Global***REMOVED***
			Global: &swarmapi.GlobalService***REMOVED******REMOVED***,
		***REMOVED***
	***REMOVED*** else if s.Mode.Replicated != nil && s.Mode.Replicated.Replicas != nil ***REMOVED***
		spec.Mode = &swarmapi.ServiceSpec_Replicated***REMOVED***
			Replicated: &swarmapi.ReplicatedService***REMOVED***Replicas: *s.Mode.Replicated.Replicas***REMOVED***,
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		spec.Mode = &swarmapi.ServiceSpec_Replicated***REMOVED***
			Replicated: &swarmapi.ReplicatedService***REMOVED***Replicas: 1***REMOVED***,
		***REMOVED***
	***REMOVED***

	return spec, nil
***REMOVED***

func annotationsFromGRPC(ann swarmapi.Annotations) types.Annotations ***REMOVED***
	a := types.Annotations***REMOVED***
		Name:   ann.Name,
		Labels: ann.Labels,
	***REMOVED***

	if a.Labels == nil ***REMOVED***
		a.Labels = make(map[string]string)
	***REMOVED***

	return a
***REMOVED***

// GenericResourcesFromGRPC converts a GRPC GenericResource to a GenericResource
func GenericResourcesFromGRPC(genericRes []*swarmapi.GenericResource) []types.GenericResource ***REMOVED***
	var generic []types.GenericResource
	for _, res := range genericRes ***REMOVED***
		var current types.GenericResource

		switch r := res.Resource.(type) ***REMOVED***
		case *swarmapi.GenericResource_DiscreteResourceSpec:
			current.DiscreteResourceSpec = &types.DiscreteGenericResource***REMOVED***
				Kind:  r.DiscreteResourceSpec.Kind,
				Value: r.DiscreteResourceSpec.Value,
			***REMOVED***
		case *swarmapi.GenericResource_NamedResourceSpec:
			current.NamedResourceSpec = &types.NamedGenericResource***REMOVED***
				Kind:  r.NamedResourceSpec.Kind,
				Value: r.NamedResourceSpec.Value,
			***REMOVED***
		***REMOVED***

		generic = append(generic, current)
	***REMOVED***

	return generic
***REMOVED***

func resourcesFromGRPC(res *swarmapi.ResourceRequirements) *types.ResourceRequirements ***REMOVED***
	var resources *types.ResourceRequirements
	if res != nil ***REMOVED***
		resources = &types.ResourceRequirements***REMOVED******REMOVED***
		if res.Limits != nil ***REMOVED***
			resources.Limits = &types.Resources***REMOVED***
				NanoCPUs:    res.Limits.NanoCPUs,
				MemoryBytes: res.Limits.MemoryBytes,
			***REMOVED***
		***REMOVED***
		if res.Reservations != nil ***REMOVED***
			resources.Reservations = &types.Resources***REMOVED***
				NanoCPUs:         res.Reservations.NanoCPUs,
				MemoryBytes:      res.Reservations.MemoryBytes,
				GenericResources: GenericResourcesFromGRPC(res.Reservations.Generic),
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return resources
***REMOVED***

// GenericResourcesToGRPC converts a GenericResource to a GRPC GenericResource
func GenericResourcesToGRPC(genericRes []types.GenericResource) []*swarmapi.GenericResource ***REMOVED***
	var generic []*swarmapi.GenericResource
	for _, res := range genericRes ***REMOVED***
		var r *swarmapi.GenericResource

		if res.DiscreteResourceSpec != nil ***REMOVED***
			r = genericresource.NewDiscrete(res.DiscreteResourceSpec.Kind, res.DiscreteResourceSpec.Value)
		***REMOVED*** else if res.NamedResourceSpec != nil ***REMOVED***
			r = genericresource.NewString(res.NamedResourceSpec.Kind, res.NamedResourceSpec.Value)
		***REMOVED***

		generic = append(generic, r)
	***REMOVED***

	return generic
***REMOVED***

func resourcesToGRPC(res *types.ResourceRequirements) *swarmapi.ResourceRequirements ***REMOVED***
	var reqs *swarmapi.ResourceRequirements
	if res != nil ***REMOVED***
		reqs = &swarmapi.ResourceRequirements***REMOVED******REMOVED***
		if res.Limits != nil ***REMOVED***
			reqs.Limits = &swarmapi.Resources***REMOVED***
				NanoCPUs:    res.Limits.NanoCPUs,
				MemoryBytes: res.Limits.MemoryBytes,
			***REMOVED***
		***REMOVED***
		if res.Reservations != nil ***REMOVED***
			reqs.Reservations = &swarmapi.Resources***REMOVED***
				NanoCPUs:    res.Reservations.NanoCPUs,
				MemoryBytes: res.Reservations.MemoryBytes,
				Generic:     GenericResourcesToGRPC(res.Reservations.GenericResources),
			***REMOVED***

		***REMOVED***
	***REMOVED***
	return reqs
***REMOVED***

func restartPolicyFromGRPC(p *swarmapi.RestartPolicy) *types.RestartPolicy ***REMOVED***
	var rp *types.RestartPolicy
	if p != nil ***REMOVED***
		rp = &types.RestartPolicy***REMOVED******REMOVED***

		switch p.Condition ***REMOVED***
		case swarmapi.RestartOnNone:
			rp.Condition = types.RestartPolicyConditionNone
		case swarmapi.RestartOnFailure:
			rp.Condition = types.RestartPolicyConditionOnFailure
		case swarmapi.RestartOnAny:
			rp.Condition = types.RestartPolicyConditionAny
		default:
			rp.Condition = types.RestartPolicyConditionAny
		***REMOVED***

		if p.Delay != nil ***REMOVED***
			delay, _ := gogotypes.DurationFromProto(p.Delay)
			rp.Delay = &delay
		***REMOVED***
		if p.Window != nil ***REMOVED***
			window, _ := gogotypes.DurationFromProto(p.Window)
			rp.Window = &window
		***REMOVED***

		rp.MaxAttempts = &p.MaxAttempts
	***REMOVED***
	return rp
***REMOVED***

func restartPolicyToGRPC(p *types.RestartPolicy) (*swarmapi.RestartPolicy, error) ***REMOVED***
	var rp *swarmapi.RestartPolicy
	if p != nil ***REMOVED***
		rp = &swarmapi.RestartPolicy***REMOVED******REMOVED***

		switch p.Condition ***REMOVED***
		case types.RestartPolicyConditionNone:
			rp.Condition = swarmapi.RestartOnNone
		case types.RestartPolicyConditionOnFailure:
			rp.Condition = swarmapi.RestartOnFailure
		case types.RestartPolicyConditionAny:
			rp.Condition = swarmapi.RestartOnAny
		default:
			if string(p.Condition) != "" ***REMOVED***
				return nil, fmt.Errorf("invalid RestartCondition: %q", p.Condition)
			***REMOVED***
			rp.Condition = swarmapi.RestartOnAny
		***REMOVED***

		if p.Delay != nil ***REMOVED***
			rp.Delay = gogotypes.DurationProto(*p.Delay)
		***REMOVED***
		if p.Window != nil ***REMOVED***
			rp.Window = gogotypes.DurationProto(*p.Window)
		***REMOVED***
		if p.MaxAttempts != nil ***REMOVED***
			rp.MaxAttempts = *p.MaxAttempts

		***REMOVED***
	***REMOVED***
	return rp, nil
***REMOVED***

func placementFromGRPC(p *swarmapi.Placement) *types.Placement ***REMOVED***
	if p == nil ***REMOVED***
		return nil
	***REMOVED***
	r := &types.Placement***REMOVED***
		Constraints: p.Constraints,
	***REMOVED***

	for _, pref := range p.Preferences ***REMOVED***
		if spread := pref.GetSpread(); spread != nil ***REMOVED***
			r.Preferences = append(r.Preferences, types.PlacementPreference***REMOVED***
				Spread: &types.SpreadOver***REMOVED***
					SpreadDescriptor: spread.SpreadDescriptor,
				***REMOVED***,
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	for _, plat := range p.Platforms ***REMOVED***
		r.Platforms = append(r.Platforms, types.Platform***REMOVED***
			Architecture: plat.Architecture,
			OS:           plat.OS,
		***REMOVED***)
	***REMOVED***

	return r
***REMOVED***

func driverFromGRPC(p *swarmapi.Driver) *types.Driver ***REMOVED***
	if p == nil ***REMOVED***
		return nil
	***REMOVED***

	return &types.Driver***REMOVED***
		Name:    p.Name,
		Options: p.Options,
	***REMOVED***
***REMOVED***

func driverToGRPC(p *types.Driver) *swarmapi.Driver ***REMOVED***
	if p == nil ***REMOVED***
		return nil
	***REMOVED***

	return &swarmapi.Driver***REMOVED***
		Name:    p.Name,
		Options: p.Options,
	***REMOVED***
***REMOVED***

func updateConfigFromGRPC(updateConfig *swarmapi.UpdateConfig) *types.UpdateConfig ***REMOVED***
	if updateConfig == nil ***REMOVED***
		return nil
	***REMOVED***

	converted := &types.UpdateConfig***REMOVED***
		Parallelism:     updateConfig.Parallelism,
		MaxFailureRatio: updateConfig.MaxFailureRatio,
	***REMOVED***

	converted.Delay = updateConfig.Delay
	if updateConfig.Monitor != nil ***REMOVED***
		converted.Monitor, _ = gogotypes.DurationFromProto(updateConfig.Monitor)
	***REMOVED***

	switch updateConfig.FailureAction ***REMOVED***
	case swarmapi.UpdateConfig_PAUSE:
		converted.FailureAction = types.UpdateFailureActionPause
	case swarmapi.UpdateConfig_CONTINUE:
		converted.FailureAction = types.UpdateFailureActionContinue
	case swarmapi.UpdateConfig_ROLLBACK:
		converted.FailureAction = types.UpdateFailureActionRollback
	***REMOVED***

	switch updateConfig.Order ***REMOVED***
	case swarmapi.UpdateConfig_STOP_FIRST:
		converted.Order = types.UpdateOrderStopFirst
	case swarmapi.UpdateConfig_START_FIRST:
		converted.Order = types.UpdateOrderStartFirst
	***REMOVED***

	return converted
***REMOVED***

func updateConfigToGRPC(updateConfig *types.UpdateConfig) (*swarmapi.UpdateConfig, error) ***REMOVED***
	if updateConfig == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	converted := &swarmapi.UpdateConfig***REMOVED***
		Parallelism:     updateConfig.Parallelism,
		Delay:           updateConfig.Delay,
		MaxFailureRatio: updateConfig.MaxFailureRatio,
	***REMOVED***

	switch updateConfig.FailureAction ***REMOVED***
	case types.UpdateFailureActionPause, "":
		converted.FailureAction = swarmapi.UpdateConfig_PAUSE
	case types.UpdateFailureActionContinue:
		converted.FailureAction = swarmapi.UpdateConfig_CONTINUE
	case types.UpdateFailureActionRollback:
		converted.FailureAction = swarmapi.UpdateConfig_ROLLBACK
	default:
		return nil, fmt.Errorf("unrecognized update failure action %s", updateConfig.FailureAction)
	***REMOVED***
	if updateConfig.Monitor != 0 ***REMOVED***
		converted.Monitor = gogotypes.DurationProto(updateConfig.Monitor)
	***REMOVED***

	switch updateConfig.Order ***REMOVED***
	case types.UpdateOrderStopFirst, "":
		converted.Order = swarmapi.UpdateConfig_STOP_FIRST
	case types.UpdateOrderStartFirst:
		converted.Order = swarmapi.UpdateConfig_START_FIRST
	default:
		return nil, fmt.Errorf("unrecognized update order %s", updateConfig.Order)
	***REMOVED***

	return converted, nil
***REMOVED***

func taskSpecFromGRPC(taskSpec swarmapi.TaskSpec) (types.TaskSpec, error) ***REMOVED***
	taskNetworks := make([]types.NetworkAttachmentConfig, 0, len(taskSpec.Networks))
	for _, n := range taskSpec.Networks ***REMOVED***
		netConfig := types.NetworkAttachmentConfig***REMOVED***Target: n.Target, Aliases: n.Aliases, DriverOpts: n.DriverAttachmentOpts***REMOVED***
		taskNetworks = append(taskNetworks, netConfig)
	***REMOVED***

	t := types.TaskSpec***REMOVED***
		Resources:     resourcesFromGRPC(taskSpec.Resources),
		RestartPolicy: restartPolicyFromGRPC(taskSpec.Restart),
		Placement:     placementFromGRPC(taskSpec.Placement),
		LogDriver:     driverFromGRPC(taskSpec.LogDriver),
		Networks:      taskNetworks,
		ForceUpdate:   taskSpec.ForceUpdate,
	***REMOVED***

	switch taskSpec.GetRuntime().(type) ***REMOVED***
	case *swarmapi.TaskSpec_Container, nil:
		c := taskSpec.GetContainer()
		if c != nil ***REMOVED***
			t.ContainerSpec = containerSpecFromGRPC(c)
		***REMOVED***
	case *swarmapi.TaskSpec_Generic:
		g := taskSpec.GetGeneric()
		if g != nil ***REMOVED***
			switch g.Kind ***REMOVED***
			case string(types.RuntimePlugin):
				var p runtime.PluginSpec
				if err := proto.Unmarshal(g.Payload.Value, &p); err != nil ***REMOVED***
					return t, errors.Wrap(err, "error unmarshalling plugin spec")
				***REMOVED***
				t.PluginSpec = &p
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return t, nil
***REMOVED***
