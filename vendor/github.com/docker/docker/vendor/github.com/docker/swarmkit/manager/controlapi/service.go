package controlapi

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/docker/distribution/reference"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/defaults"
	"github.com/docker/swarmkit/api/genericresource"
	"github.com/docker/swarmkit/api/naming"
	"github.com/docker/swarmkit/identity"
	"github.com/docker/swarmkit/manager/allocator"
	"github.com/docker/swarmkit/manager/constraint"
	"github.com/docker/swarmkit/manager/state/store"
	"github.com/docker/swarmkit/protobuf/ptypes"
	"github.com/docker/swarmkit/template"
	gogotypes "github.com/gogo/protobuf/types"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errNetworkUpdateNotSupported = errors.New("networks must be migrated to TaskSpec before being changed")
	errRenameNotSupported        = errors.New("renaming services is not supported")
	errModeChangeNotAllowed      = errors.New("service mode change is not allowed")
)

const minimumDuration = 1 * time.Millisecond

func validateResources(r *api.Resources) error ***REMOVED***
	if r == nil ***REMOVED***
		return nil
	***REMOVED***

	if r.NanoCPUs != 0 && r.NanoCPUs < 1e6 ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "invalid cpu value %g: Must be at least %g", float64(r.NanoCPUs)/1e9, 1e6/1e9)
	***REMOVED***

	if r.MemoryBytes != 0 && r.MemoryBytes < 4*1024*1024 ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "invalid memory value %d: Must be at least 4MiB", r.MemoryBytes)
	***REMOVED***
	if err := genericresource.ValidateTask(r); err != nil ***REMOVED***
		return nil
	***REMOVED***
	return nil
***REMOVED***

func validateResourceRequirements(r *api.ResourceRequirements) error ***REMOVED***
	if r == nil ***REMOVED***
		return nil
	***REMOVED***
	if err := validateResources(r.Limits); err != nil ***REMOVED***
		return err
	***REMOVED***
	return validateResources(r.Reservations)
***REMOVED***

func validateRestartPolicy(rp *api.RestartPolicy) error ***REMOVED***
	if rp == nil ***REMOVED***
		return nil
	***REMOVED***

	if rp.Delay != nil ***REMOVED***
		delay, err := gogotypes.DurationFromProto(rp.Delay)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if delay < 0 ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "TaskSpec: restart-delay cannot be negative")
		***REMOVED***
	***REMOVED***

	if rp.Window != nil ***REMOVED***
		win, err := gogotypes.DurationFromProto(rp.Window)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if win < 0 ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "TaskSpec: restart-window cannot be negative")
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func validatePlacement(placement *api.Placement) error ***REMOVED***
	if placement == nil ***REMOVED***
		return nil
	***REMOVED***
	_, err := constraint.Parse(placement.Constraints)
	return err
***REMOVED***

func validateUpdate(uc *api.UpdateConfig) error ***REMOVED***
	if uc == nil ***REMOVED***
		return nil
	***REMOVED***

	if uc.Delay < 0 ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "TaskSpec: update-delay cannot be negative")
	***REMOVED***

	if uc.Monitor != nil ***REMOVED***
		monitor, err := gogotypes.DurationFromProto(uc.Monitor)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if monitor < 0 ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "TaskSpec: update-monitor cannot be negative")
		***REMOVED***
	***REMOVED***

	if uc.MaxFailureRatio < 0 || uc.MaxFailureRatio > 1 ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "TaskSpec: update-maxfailureratio cannot be less than 0 or bigger than 1")
	***REMOVED***

	return nil
***REMOVED***

func validateContainerSpec(taskSpec api.TaskSpec) error ***REMOVED***
	// Building a empty/dummy Task to validate the templating and
	// the resulting container spec as well. This is a *best effort*
	// validation.
	container, err := template.ExpandContainerSpec(&api.NodeDescription***REMOVED***
		Hostname: "nodeHostname",
		Platform: &api.Platform***REMOVED***
			OS:           "os",
			Architecture: "architecture",
		***REMOVED***,
	***REMOVED***, &api.Task***REMOVED***
		Spec:      taskSpec,
		ServiceID: "serviceid",
		Slot:      1,
		NodeID:    "nodeid",
		Networks:  []*api.NetworkAttachment***REMOVED******REMOVED***,
		Annotations: api.Annotations***REMOVED***
			Name: "taskname",
		***REMOVED***,
		ServiceAnnotations: api.Annotations***REMOVED***
			Name: "servicename",
		***REMOVED***,
		Endpoint:  &api.Endpoint***REMOVED******REMOVED***,
		LogDriver: taskSpec.LogDriver,
	***REMOVED***)
	if err != nil ***REMOVED***
		return status.Errorf(codes.InvalidArgument, err.Error())
	***REMOVED***

	if err := validateImage(container.Image); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := validateMounts(container.Mounts); err != nil ***REMOVED***
		return err
	***REMOVED***

	return validateHealthCheck(container.Healthcheck)
***REMOVED***

// validateImage validates image name in containerSpec
func validateImage(image string) error ***REMOVED***
	if image == "" ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "ContainerSpec: image reference must be provided")
	***REMOVED***

	if _, err := reference.ParseNormalizedNamed(image); err != nil ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "ContainerSpec: %q is not a valid repository/tag", image)
	***REMOVED***
	return nil
***REMOVED***

// validateMounts validates if there are duplicate mounts in containerSpec
func validateMounts(mounts []api.Mount) error ***REMOVED***
	mountMap := make(map[string]bool)
	for _, mount := range mounts ***REMOVED***
		if _, exists := mountMap[mount.Target]; exists ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "ContainerSpec: duplicate mount point: %s", mount.Target)
		***REMOVED***
		mountMap[mount.Target] = true
	***REMOVED***

	return nil
***REMOVED***

// validateHealthCheck validates configs about container's health check
func validateHealthCheck(hc *api.HealthConfig) error ***REMOVED***
	if hc == nil ***REMOVED***
		return nil
	***REMOVED***

	if hc.Interval != nil ***REMOVED***
		interval, err := gogotypes.DurationFromProto(hc.Interval)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if interval != 0 && interval < time.Duration(minimumDuration) ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "ContainerSpec: Interval in HealthConfig cannot be less than %s", minimumDuration)
		***REMOVED***
	***REMOVED***

	if hc.Timeout != nil ***REMOVED***
		timeout, err := gogotypes.DurationFromProto(hc.Timeout)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if timeout != 0 && timeout < time.Duration(minimumDuration) ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "ContainerSpec: Timeout in HealthConfig cannot be less than %s", minimumDuration)
		***REMOVED***
	***REMOVED***

	if hc.StartPeriod != nil ***REMOVED***
		sp, err := gogotypes.DurationFromProto(hc.StartPeriod)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if sp != 0 && sp < time.Duration(minimumDuration) ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "ContainerSpec: StartPeriod in HealthConfig cannot be less than %s", minimumDuration)
		***REMOVED***
	***REMOVED***

	if hc.Retries < 0 ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "ContainerSpec: Retries in HealthConfig cannot be negative")
	***REMOVED***

	return nil
***REMOVED***

func validateGenericRuntimeSpec(taskSpec api.TaskSpec) error ***REMOVED***
	generic := taskSpec.GetGeneric()

	if len(generic.Kind) < 3 ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "Generic runtime: Invalid name %q", generic.Kind)
	***REMOVED***

	reservedNames := []string***REMOVED***"container", "attachment"***REMOVED***
	for _, n := range reservedNames ***REMOVED***
		if strings.ToLower(generic.Kind) == n ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "Generic runtime: %q is a reserved name", generic.Kind)
		***REMOVED***
	***REMOVED***

	payload := generic.Payload

	if payload == nil ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "Generic runtime is missing payload")
	***REMOVED***

	if payload.TypeUrl == "" ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "Generic runtime is missing payload type")
	***REMOVED***

	if len(payload.Value) == 0 ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "Generic runtime has an empty payload")
	***REMOVED***

	return nil
***REMOVED***

func validateTaskSpec(taskSpec api.TaskSpec) error ***REMOVED***
	if err := validateResourceRequirements(taskSpec.Resources); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := validateRestartPolicy(taskSpec.Restart); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := validatePlacement(taskSpec.Placement); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Check to see if the secret reference portion of the spec is valid
	if err := validateSecretRefsSpec(taskSpec); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Check to see if the config reference portion of the spec is valid
	if err := validateConfigRefsSpec(taskSpec); err != nil ***REMOVED***
		return err
	***REMOVED***

	if taskSpec.GetRuntime() == nil ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "TaskSpec: missing runtime")
	***REMOVED***

	switch taskSpec.GetRuntime().(type) ***REMOVED***
	case *api.TaskSpec_Container:
		if err := validateContainerSpec(taskSpec); err != nil ***REMOVED***
			return err
		***REMOVED***
	case *api.TaskSpec_Generic:
		if err := validateGenericRuntimeSpec(taskSpec); err != nil ***REMOVED***
			return err
		***REMOVED***
	default:
		return status.Errorf(codes.Unimplemented, "RuntimeSpec: unimplemented runtime in service spec")
	***REMOVED***

	return nil
***REMOVED***

func validateEndpointSpec(epSpec *api.EndpointSpec) error ***REMOVED***
	// Endpoint spec is optional
	if epSpec == nil ***REMOVED***
		return nil
	***REMOVED***

	type portSpec struct ***REMOVED***
		publishedPort uint32
		protocol      api.PortConfig_Protocol
	***REMOVED***

	portSet := make(map[portSpec]struct***REMOVED******REMOVED***)
	for _, port := range epSpec.Ports ***REMOVED***
		// Publish mode = "ingress" represents Routing-Mesh and current implementation
		// of routing-mesh relies on IPVS based load-balancing with input=published-port.
		// But Endpoint-Spec mode of DNSRR relies on multiple A records and cannot be used
		// with routing-mesh (PublishMode="ingress") which cannot rely on DNSRR.
		// But PublishMode="host" doesn't provide Routing-Mesh and the DNSRR is applicable
		// for the backend network and hence we accept that configuration.

		if epSpec.Mode == api.ResolutionModeDNSRoundRobin && port.PublishMode == api.PublishModeIngress ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "EndpointSpec: port published with ingress mode can't be used with dnsrr mode")
		***REMOVED***

		// If published port is not specified, it does not conflict
		// with any others.
		if port.PublishedPort == 0 ***REMOVED***
			continue
		***REMOVED***

		portSpec := portSpec***REMOVED***publishedPort: port.PublishedPort, protocol: port.Protocol***REMOVED***
		if _, ok := portSet[portSpec]; ok ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "EndpointSpec: duplicate published ports provided")
		***REMOVED***

		portSet[portSpec] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// validateSecretRefsSpec finds if the secrets passed in spec are valid and have no
// conflicting targets.
func validateSecretRefsSpec(spec api.TaskSpec) error ***REMOVED***
	container := spec.GetContainer()
	if container == nil ***REMOVED***
		return nil
	***REMOVED***

	// Keep a map to track all the targets that will be exposed
	// The string returned is only used for logging. It could as well be struct***REMOVED******REMOVED******REMOVED******REMOVED***
	existingTargets := make(map[string]string)
	for _, secretRef := range container.Secrets ***REMOVED***
		// SecretID and SecretName are mandatory, we have invalid references without them
		if secretRef.SecretID == "" || secretRef.SecretName == "" ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "malformed secret reference")
		***REMOVED***

		// Every secret reference requires a Target
		if secretRef.GetTarget() == nil ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "malformed secret reference, no target provided")
		***REMOVED***

		// If this is a file target, we will ensure filename uniqueness
		if secretRef.GetFile() != nil ***REMOVED***
			fileName := secretRef.GetFile().Name
			if fileName == "" ***REMOVED***
				return status.Errorf(codes.InvalidArgument, "malformed file secret reference, invalid target file name provided")
			***REMOVED***
			// If this target is already in use, we have conflicting targets
			if prevSecretName, ok := existingTargets[fileName]; ok ***REMOVED***
				return status.Errorf(codes.InvalidArgument, "secret references '%s' and '%s' have a conflicting target: '%s'", prevSecretName, secretRef.SecretName, fileName)
			***REMOVED***

			existingTargets[fileName] = secretRef.SecretName
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// validateConfigRefsSpec finds if the configs passed in spec are valid and have no
// conflicting targets.
func validateConfigRefsSpec(spec api.TaskSpec) error ***REMOVED***
	container := spec.GetContainer()
	if container == nil ***REMOVED***
		return nil
	***REMOVED***

	// Keep a map to track all the targets that will be exposed
	// The string returned is only used for logging. It could as well be struct***REMOVED******REMOVED******REMOVED******REMOVED***
	existingTargets := make(map[string]string)
	for _, configRef := range container.Configs ***REMOVED***
		// ConfigID and ConfigName are mandatory, we have invalid references without them
		if configRef.ConfigID == "" || configRef.ConfigName == "" ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "malformed config reference")
		***REMOVED***

		// Every config reference requires a Target
		if configRef.GetTarget() == nil ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "malformed config reference, no target provided")
		***REMOVED***

		// If this is a file target, we will ensure filename uniqueness
		if configRef.GetFile() != nil ***REMOVED***
			fileName := configRef.GetFile().Name
			// Validate the file name
			if fileName == "" ***REMOVED***
				return status.Errorf(codes.InvalidArgument, "malformed file config reference, invalid target file name provided")
			***REMOVED***

			// If this target is already in use, we have conflicting targets
			if prevConfigName, ok := existingTargets[fileName]; ok ***REMOVED***
				return status.Errorf(codes.InvalidArgument, "config references '%s' and '%s' have a conflicting target: '%s'", prevConfigName, configRef.ConfigName, fileName)
			***REMOVED***

			existingTargets[fileName] = configRef.ConfigName
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (s *Server) validateNetworks(networks []*api.NetworkAttachmentConfig) error ***REMOVED***
	for _, na := range networks ***REMOVED***
		var network *api.Network
		s.store.View(func(tx store.ReadTx) ***REMOVED***
			network = store.GetNetwork(tx, na.Target)
		***REMOVED***)
		if network == nil ***REMOVED***
			continue
		***REMOVED***
		if allocator.IsIngressNetwork(network) ***REMOVED***
			return status.Errorf(codes.InvalidArgument,
				"Service cannot be explicitly attached to the ingress network %q", network.Spec.Annotations.Name)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func validateMode(s *api.ServiceSpec) error ***REMOVED***
	m := s.GetMode()
	switch m.(type) ***REMOVED***
	case *api.ServiceSpec_Replicated:
		if int64(m.(*api.ServiceSpec_Replicated).Replicated.Replicas) < 0 ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "Number of replicas must be non-negative")
		***REMOVED***
	case *api.ServiceSpec_Global:
	default:
		return status.Errorf(codes.InvalidArgument, "Unrecognized service mode")
	***REMOVED***

	return nil
***REMOVED***

func validateServiceSpec(spec *api.ServiceSpec) error ***REMOVED***
	if spec == nil ***REMOVED***
		return status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***
	if err := validateAnnotations(spec.Annotations); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := validateTaskSpec(spec.Task); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := validateUpdate(spec.Update); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := validateEndpointSpec(spec.Endpoint); err != nil ***REMOVED***
		return err
	***REMOVED***
	return validateMode(spec)
***REMOVED***

// checkPortConflicts does a best effort to find if the passed in spec has port
// conflicts with existing services.
// `serviceID string` is the service ID of the spec in service update. If
// `serviceID` is not "", then conflicts check will be skipped against this
// service (the service being updated).
func (s *Server) checkPortConflicts(spec *api.ServiceSpec, serviceID string) error ***REMOVED***
	if spec.Endpoint == nil ***REMOVED***
		return nil
	***REMOVED***

	type portSpec struct ***REMOVED***
		protocol      api.PortConfig_Protocol
		publishedPort uint32
	***REMOVED***

	pcToStruct := func(pc *api.PortConfig) portSpec ***REMOVED***
		return portSpec***REMOVED***
			protocol:      pc.Protocol,
			publishedPort: pc.PublishedPort,
		***REMOVED***
	***REMOVED***

	ingressPorts := make(map[portSpec]struct***REMOVED******REMOVED***)
	hostModePorts := make(map[portSpec]struct***REMOVED******REMOVED***)
	for _, pc := range spec.Endpoint.Ports ***REMOVED***
		if pc.PublishedPort == 0 ***REMOVED***
			continue
		***REMOVED***
		switch pc.PublishMode ***REMOVED***
		case api.PublishModeIngress:
			ingressPorts[pcToStruct(pc)] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		case api.PublishModeHost:
			hostModePorts[pcToStruct(pc)] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
	if len(ingressPorts) == 0 && len(hostModePorts) == 0 ***REMOVED***
		return nil
	***REMOVED***

	var (
		services []*api.Service
		err      error
	)

	s.store.View(func(tx store.ReadTx) ***REMOVED***
		services, err = store.FindServices(tx, store.All)
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	isPortInUse := func(pc *api.PortConfig, service *api.Service) error ***REMOVED***
		if pc.PublishedPort == 0 ***REMOVED***
			return nil
		***REMOVED***

		switch pc.PublishMode ***REMOVED***
		case api.PublishModeHost:
			if _, ok := ingressPorts[pcToStruct(pc)]; ok ***REMOVED***
				return status.Errorf(codes.InvalidArgument, "port '%d' is already in use by service '%s' (%s) as a host-published port", pc.PublishedPort, service.Spec.Annotations.Name, service.ID)
			***REMOVED***

			// Multiple services with same port in host publish mode can
			// coexist - this is handled by the scheduler.
			return nil
		case api.PublishModeIngress:
			_, ingressConflict := ingressPorts[pcToStruct(pc)]
			_, hostModeConflict := hostModePorts[pcToStruct(pc)]
			if ingressConflict || hostModeConflict ***REMOVED***
				return status.Errorf(codes.InvalidArgument, "port '%d' is already in use by service '%s' (%s) as an ingress port", pc.PublishedPort, service.Spec.Annotations.Name, service.ID)
			***REMOVED***
		***REMOVED***

		return nil
	***REMOVED***

	for _, service := range services ***REMOVED***
		// If service ID is the same (and not "") then this is an update
		if serviceID != "" && serviceID == service.ID ***REMOVED***
			continue
		***REMOVED***
		if service.Spec.Endpoint != nil ***REMOVED***
			for _, pc := range service.Spec.Endpoint.Ports ***REMOVED***
				if err := isPortInUse(pc, service); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if service.Endpoint != nil ***REMOVED***
			for _, pc := range service.Endpoint.Ports ***REMOVED***
				if err := isPortInUse(pc, service); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// checkSecretExistence finds if the secret exists
func (s *Server) checkSecretExistence(tx store.Tx, spec *api.ServiceSpec) error ***REMOVED***
	container := spec.Task.GetContainer()
	if container == nil ***REMOVED***
		return nil
	***REMOVED***

	var failedSecrets []string
	for _, secretRef := range container.Secrets ***REMOVED***
		secret := store.GetSecret(tx, secretRef.SecretID)
		// Check to see if the secret exists and secretRef.SecretName matches the actual secretName
		if secret == nil || secret.Spec.Annotations.Name != secretRef.SecretName ***REMOVED***
			failedSecrets = append(failedSecrets, secretRef.SecretName)
		***REMOVED***
	***REMOVED***

	if len(failedSecrets) > 0 ***REMOVED***
		secretStr := "secrets"
		if len(failedSecrets) == 1 ***REMOVED***
			secretStr = "secret"
		***REMOVED***

		return status.Errorf(codes.InvalidArgument, "%s not found: %v", secretStr, strings.Join(failedSecrets, ", "))

	***REMOVED***

	return nil
***REMOVED***

// checkConfigExistence finds if the config exists
func (s *Server) checkConfigExistence(tx store.Tx, spec *api.ServiceSpec) error ***REMOVED***
	container := spec.Task.GetContainer()
	if container == nil ***REMOVED***
		return nil
	***REMOVED***

	var failedConfigs []string
	for _, configRef := range container.Configs ***REMOVED***
		config := store.GetConfig(tx, configRef.ConfigID)
		// Check to see if the config exists and configRef.ConfigName matches the actual configName
		if config == nil || config.Spec.Annotations.Name != configRef.ConfigName ***REMOVED***
			failedConfigs = append(failedConfigs, configRef.ConfigName)
		***REMOVED***
	***REMOVED***

	if len(failedConfigs) > 0 ***REMOVED***
		configStr := "configs"
		if len(failedConfigs) == 1 ***REMOVED***
			configStr = "config"
		***REMOVED***

		return status.Errorf(codes.InvalidArgument, "%s not found: %v", configStr, strings.Join(failedConfigs, ", "))

	***REMOVED***

	return nil
***REMOVED***

// CreateService creates and returns a Service based on the provided ServiceSpec.
// - Returns `InvalidArgument` if the ServiceSpec is malformed.
// - Returns `Unimplemented` if the ServiceSpec references unimplemented features.
// - Returns `AlreadyExists` if the ServiceID conflicts.
// - Returns an error if the creation fails.
func (s *Server) CreateService(ctx context.Context, request *api.CreateServiceRequest) (*api.CreateServiceResponse, error) ***REMOVED***
	if err := validateServiceSpec(request.Spec); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := s.validateNetworks(request.Spec.Networks); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := s.checkPortConflicts(request.Spec, ""); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// TODO(aluzzardi): Consider using `Name` as a primary key to handle
	// duplicate creations. See #65
	service := &api.Service***REMOVED***
		ID:          identity.NewID(),
		Spec:        *request.Spec,
		SpecVersion: &api.Version***REMOVED******REMOVED***,
	***REMOVED***

	if allocator.IsIngressNetworkNeeded(service) ***REMOVED***
		if _, err := allocator.GetIngressNetwork(s.store); err == allocator.ErrNoIngress ***REMOVED***
			return nil, status.Errorf(codes.FailedPrecondition, "service needs ingress network, but no ingress network is present")
		***REMOVED***
	***REMOVED***

	err := s.store.Update(func(tx store.Tx) error ***REMOVED***
		// Check to see if all the secrets being added exist as objects
		// in our datastore
		err := s.checkSecretExistence(tx, request.Spec)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		err = s.checkConfigExistence(tx, request.Spec)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		return store.CreateService(tx, service)
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &api.CreateServiceResponse***REMOVED***
		Service: service,
	***REMOVED***, nil
***REMOVED***

// GetService returns a Service given a ServiceID.
// - Returns `InvalidArgument` if ServiceID is not provided.
// - Returns `NotFound` if the Service is not found.
func (s *Server) GetService(ctx context.Context, request *api.GetServiceRequest) (*api.GetServiceResponse, error) ***REMOVED***
	if request.ServiceID == "" ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***

	var service *api.Service
	s.store.View(func(tx store.ReadTx) ***REMOVED***
		service = store.GetService(tx, request.ServiceID)
	***REMOVED***)
	if service == nil ***REMOVED***
		return nil, status.Errorf(codes.NotFound, "service %s not found", request.ServiceID)
	***REMOVED***

	if request.InsertDefaults ***REMOVED***
		service.Spec = *defaults.InterpolateService(&service.Spec)
	***REMOVED***

	return &api.GetServiceResponse***REMOVED***
		Service: service,
	***REMOVED***, nil
***REMOVED***

// UpdateService updates a Service referenced by ServiceID with the given ServiceSpec.
// - Returns `NotFound` if the Service is not found.
// - Returns `InvalidArgument` if the ServiceSpec is malformed.
// - Returns `Unimplemented` if the ServiceSpec references unimplemented features.
// - Returns an error if the update fails.
func (s *Server) UpdateService(ctx context.Context, request *api.UpdateServiceRequest) (*api.UpdateServiceResponse, error) ***REMOVED***
	if request.ServiceID == "" || request.ServiceVersion == nil ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***
	if err := validateServiceSpec(request.Spec); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var service *api.Service
	s.store.View(func(tx store.ReadTx) ***REMOVED***
		service = store.GetService(tx, request.ServiceID)
	***REMOVED***)
	if service == nil ***REMOVED***
		return nil, status.Errorf(codes.NotFound, "service %s not found", request.ServiceID)
	***REMOVED***

	if request.Spec.Endpoint != nil && !reflect.DeepEqual(request.Spec.Endpoint, service.Spec.Endpoint) ***REMOVED***
		if err := s.checkPortConflicts(request.Spec, request.ServiceID); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	err := s.store.Update(func(tx store.Tx) error ***REMOVED***
		service = store.GetService(tx, request.ServiceID)
		if service == nil ***REMOVED***
			return status.Errorf(codes.NotFound, "service %s not found", request.ServiceID)
		***REMOVED***

		// It's not okay to update Service.Spec.Networks on its own.
		// However, if Service.Spec.Task.Networks is also being
		// updated, that's okay (for example when migrating from the
		// deprecated Spec.Networks field to Spec.Task.Networks).
		if (len(request.Spec.Networks) != 0 || len(service.Spec.Networks) != 0) &&
			!reflect.DeepEqual(request.Spec.Networks, service.Spec.Networks) &&
			reflect.DeepEqual(request.Spec.Task.Networks, service.Spec.Task.Networks) ***REMOVED***
			return status.Errorf(codes.Unimplemented, errNetworkUpdateNotSupported.Error())
		***REMOVED***

		// Check to see if all the secrets being added exist as objects
		// in our datastore
		err := s.checkSecretExistence(tx, request.Spec)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		err = s.checkConfigExistence(tx, request.Spec)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// orchestrator is designed to be stateless, so it should not deal
		// with service mode change (comparing current config with previous config).
		// proper way to change service mode is to delete and re-add.
		if reflect.TypeOf(service.Spec.Mode) != reflect.TypeOf(request.Spec.Mode) ***REMOVED***
			return status.Errorf(codes.Unimplemented, errModeChangeNotAllowed.Error())
		***REMOVED***

		if service.Spec.Annotations.Name != request.Spec.Annotations.Name ***REMOVED***
			return status.Errorf(codes.Unimplemented, errRenameNotSupported.Error())
		***REMOVED***

		service.Meta.Version = *request.ServiceVersion

		if request.Rollback == api.UpdateServiceRequest_PREVIOUS ***REMOVED***
			if service.PreviousSpec == nil ***REMOVED***
				return status.Errorf(codes.FailedPrecondition, "service %s does not have a previous spec", request.ServiceID)
			***REMOVED***

			curSpec := service.Spec.Copy()
			curSpecVersion := service.SpecVersion
			service.Spec = *service.PreviousSpec.Copy()
			service.SpecVersion = service.PreviousSpecVersion.Copy()
			service.PreviousSpec = curSpec
			service.PreviousSpecVersion = curSpecVersion

			service.UpdateStatus = &api.UpdateStatus***REMOVED***
				State:     api.UpdateStatus_ROLLBACK_STARTED,
				Message:   "manually requested rollback",
				StartedAt: ptypes.MustTimestampProto(time.Now()),
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			service.PreviousSpec = service.Spec.Copy()
			service.PreviousSpecVersion = service.SpecVersion
			service.Spec = *request.Spec.Copy()
			// Set spec version. Note that this will not match the
			// service's Meta.Version after the store update. The
			// versions for the spec and the service itself are not
			// meant to be directly comparable.
			service.SpecVersion = service.Meta.Version.Copy()

			// Reset update status
			service.UpdateStatus = nil
		***REMOVED***

		if allocator.IsIngressNetworkNeeded(service) ***REMOVED***
			if _, err := allocator.GetIngressNetwork(s.store); err == allocator.ErrNoIngress ***REMOVED***
				return status.Errorf(codes.FailedPrecondition, "service needs ingress network, but no ingress network is present")
			***REMOVED***
		***REMOVED***

		return store.UpdateService(tx, service)
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &api.UpdateServiceResponse***REMOVED***
		Service: service,
	***REMOVED***, nil
***REMOVED***

// RemoveService removes a Service referenced by ServiceID.
// - Returns `InvalidArgument` if ServiceID is not provided.
// - Returns `NotFound` if the Service is not found.
// - Returns an error if the deletion fails.
func (s *Server) RemoveService(ctx context.Context, request *api.RemoveServiceRequest) (*api.RemoveServiceResponse, error) ***REMOVED***
	if request.ServiceID == "" ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***

	err := s.store.Update(func(tx store.Tx) error ***REMOVED***
		return store.DeleteService(tx, request.ServiceID)
	***REMOVED***)
	if err != nil ***REMOVED***
		if err == store.ErrNotExist ***REMOVED***
			return nil, status.Errorf(codes.NotFound, "service %s not found", request.ServiceID)
		***REMOVED***
		return nil, err
	***REMOVED***
	return &api.RemoveServiceResponse***REMOVED******REMOVED***, nil
***REMOVED***

func filterServices(candidates []*api.Service, filters ...func(*api.Service) bool) []*api.Service ***REMOVED***
	result := []*api.Service***REMOVED******REMOVED***

	for _, c := range candidates ***REMOVED***
		match := true
		for _, f := range filters ***REMOVED***
			if !f(c) ***REMOVED***
				match = false
				break
			***REMOVED***
		***REMOVED***
		if match ***REMOVED***
			result = append(result, c)
		***REMOVED***
	***REMOVED***

	return result
***REMOVED***

// ListServices returns a list of all services.
func (s *Server) ListServices(ctx context.Context, request *api.ListServicesRequest) (*api.ListServicesResponse, error) ***REMOVED***
	var (
		services []*api.Service
		err      error
	)

	s.store.View(func(tx store.ReadTx) ***REMOVED***
		switch ***REMOVED***
		case request.Filters != nil && len(request.Filters.Names) > 0:
			services, err = store.FindServices(tx, buildFilters(store.ByName, request.Filters.Names))
		case request.Filters != nil && len(request.Filters.NamePrefixes) > 0:
			services, err = store.FindServices(tx, buildFilters(store.ByNamePrefix, request.Filters.NamePrefixes))
		case request.Filters != nil && len(request.Filters.IDPrefixes) > 0:
			services, err = store.FindServices(tx, buildFilters(store.ByIDPrefix, request.Filters.IDPrefixes))
		case request.Filters != nil && len(request.Filters.Runtimes) > 0:
			services, err = store.FindServices(tx, buildFilters(store.ByRuntime, request.Filters.Runtimes))
		default:
			services, err = store.FindServices(tx, store.All)
		***REMOVED***
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if request.Filters != nil ***REMOVED***
		services = filterServices(services,
			func(e *api.Service) bool ***REMOVED***
				return filterContains(e.Spec.Annotations.Name, request.Filters.Names)
			***REMOVED***,
			func(e *api.Service) bool ***REMOVED***
				return filterContainsPrefix(e.Spec.Annotations.Name, request.Filters.NamePrefixes)
			***REMOVED***,
			func(e *api.Service) bool ***REMOVED***
				return filterContainsPrefix(e.ID, request.Filters.IDPrefixes)
			***REMOVED***,
			func(e *api.Service) bool ***REMOVED***
				return filterMatchLabels(e.Spec.Annotations.Labels, request.Filters.Labels)
			***REMOVED***,
			func(e *api.Service) bool ***REMOVED***
				if len(request.Filters.Runtimes) == 0 ***REMOVED***
					return true
				***REMOVED***
				r, err := naming.Runtime(e.Spec.Task)
				if err != nil ***REMOVED***
					return false
				***REMOVED***
				return filterContains(r, request.Filters.Runtimes)
			***REMOVED***,
		)
	***REMOVED***

	return &api.ListServicesResponse***REMOVED***
		Services: services,
	***REMOVED***, nil
***REMOVED***
