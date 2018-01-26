package controlapi

import (
	"net"

	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/ipamapi"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/identity"
	"github.com/docker/swarmkit/manager/allocator"
	"github.com/docker/swarmkit/manager/allocator/networkallocator"
	"github.com/docker/swarmkit/manager/state/store"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateIPAMConfiguration(ipamConf *api.IPAMConfig) error ***REMOVED***
	if ipamConf == nil ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "ipam configuration: cannot be empty")
	***REMOVED***

	_, subnet, err := net.ParseCIDR(ipamConf.Subnet)
	if err != nil ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "ipam configuration: invalid subnet %s", ipamConf.Subnet)
	***REMOVED***

	if ipamConf.Range != "" ***REMOVED***
		ip, _, err := net.ParseCIDR(ipamConf.Range)
		if err != nil ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "ipam configuration: invalid range %s", ipamConf.Range)
		***REMOVED***

		if !subnet.Contains(ip) ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "ipam configuration: subnet %s does not contain range %s", ipamConf.Subnet, ipamConf.Range)
		***REMOVED***
	***REMOVED***

	if ipamConf.Gateway != "" ***REMOVED***
		ip := net.ParseIP(ipamConf.Gateway)
		if ip == nil ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "ipam configuration: invalid gateway %s", ipamConf.Gateway)
		***REMOVED***

		if !subnet.Contains(ip) ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "ipam configuration: subnet %s does not contain gateway %s", ipamConf.Subnet, ipamConf.Gateway)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func validateIPAM(ipam *api.IPAMOptions, pg plugingetter.PluginGetter) error ***REMOVED***
	if ipam == nil ***REMOVED***
		// It is ok to not specify any IPAM configurations. We
		// will choose good defaults.
		return nil
	***REMOVED***

	if err := validateDriver(ipam.Driver, pg, ipamapi.PluginEndpointType); err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, ipamConf := range ipam.Configs ***REMOVED***
		if err := validateIPAMConfiguration(ipamConf); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func validateNetworkSpec(spec *api.NetworkSpec, pg plugingetter.PluginGetter) error ***REMOVED***
	if spec == nil ***REMOVED***
		return status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***

	if spec.Ingress && spec.DriverConfig != nil && spec.DriverConfig.Name != "overlay" ***REMOVED***
		return status.Errorf(codes.Unimplemented, "only overlay driver is currently supported for ingress network")
	***REMOVED***

	if spec.Attachable && spec.Ingress ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "ingress network cannot be attachable")
	***REMOVED***

	if err := validateAnnotations(spec.Annotations); err != nil ***REMOVED***
		return err
	***REMOVED***

	if _, ok := spec.Annotations.Labels[networkallocator.PredefinedLabel]; ok ***REMOVED***
		return status.Errorf(codes.PermissionDenied, "label %s is for internally created predefined networks and cannot be applied by users",
			networkallocator.PredefinedLabel)
	***REMOVED***
	if err := validateDriver(spec.DriverConfig, pg, driverapi.NetworkPluginEndpointType); err != nil ***REMOVED***
		return err
	***REMOVED***

	return validateIPAM(spec.IPAM, pg)
***REMOVED***

// CreateNetwork creates and returns a Network based on the provided NetworkSpec.
// - Returns `InvalidArgument` if the NetworkSpec is malformed.
// - Returns an error if the creation fails.
func (s *Server) CreateNetwork(ctx context.Context, request *api.CreateNetworkRequest) (*api.CreateNetworkResponse, error) ***REMOVED***
	if err := validateNetworkSpec(request.Spec, s.pg); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// TODO(mrjana): Consider using `Name` as a primary key to handle
	// duplicate creations. See #65
	n := &api.Network***REMOVED***
		ID:   identity.NewID(),
		Spec: *request.Spec,
	***REMOVED***

	err := s.store.Update(func(tx store.Tx) error ***REMOVED***
		if request.Spec.Ingress ***REMOVED***
			if n, err := allocator.GetIngressNetwork(s.store); err == nil ***REMOVED***
				return status.Errorf(codes.AlreadyExists, "ingress network (%s) is already present", n.ID)
			***REMOVED*** else if err != allocator.ErrNoIngress ***REMOVED***
				return status.Errorf(codes.Internal, "failed ingress network presence check: %v", err)
			***REMOVED***
		***REMOVED***
		return store.CreateNetwork(tx, n)
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &api.CreateNetworkResponse***REMOVED***
		Network: n,
	***REMOVED***, nil
***REMOVED***

// GetNetwork returns a Network given a NetworkID.
// - Returns `InvalidArgument` if NetworkID is not provided.
// - Returns `NotFound` if the Network is not found.
func (s *Server) GetNetwork(ctx context.Context, request *api.GetNetworkRequest) (*api.GetNetworkResponse, error) ***REMOVED***
	if request.NetworkID == "" ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***

	var n *api.Network
	s.store.View(func(tx store.ReadTx) ***REMOVED***
		n = store.GetNetwork(tx, request.NetworkID)
	***REMOVED***)
	if n == nil ***REMOVED***
		return nil, status.Errorf(codes.NotFound, "network %s not found", request.NetworkID)
	***REMOVED***
	return &api.GetNetworkResponse***REMOVED***
		Network: n,
	***REMOVED***, nil
***REMOVED***

// RemoveNetwork removes a Network referenced by NetworkID.
// - Returns `InvalidArgument` if NetworkID is not provided.
// - Returns `NotFound` if the Network is not found.
// - Returns an error if the deletion fails.
func (s *Server) RemoveNetwork(ctx context.Context, request *api.RemoveNetworkRequest) (*api.RemoveNetworkResponse, error) ***REMOVED***
	if request.NetworkID == "" ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***

	var (
		n  *api.Network
		rm = s.removeNetwork
	)

	s.store.View(func(tx store.ReadTx) ***REMOVED***
		n = store.GetNetwork(tx, request.NetworkID)
	***REMOVED***)
	if n == nil ***REMOVED***
		return nil, status.Errorf(codes.NotFound, "network %s not found", request.NetworkID)
	***REMOVED***

	if allocator.IsIngressNetwork(n) ***REMOVED***
		rm = s.removeIngressNetwork
	***REMOVED***

	if v, ok := n.Spec.Annotations.Labels[networkallocator.PredefinedLabel]; ok && v == "true" ***REMOVED***
		return nil, status.Errorf(codes.FailedPrecondition, "network %s (%s) is a swarm predefined network and cannot be removed",
			request.NetworkID, n.Spec.Annotations.Name)
	***REMOVED***

	if err := rm(n.ID); err != nil ***REMOVED***
		if err == store.ErrNotExist ***REMOVED***
			return nil, status.Errorf(codes.NotFound, "network %s not found", request.NetworkID)
		***REMOVED***
		return nil, err
	***REMOVED***
	return &api.RemoveNetworkResponse***REMOVED******REMOVED***, nil
***REMOVED***

func (s *Server) removeNetwork(id string) error ***REMOVED***
	return s.store.Update(func(tx store.Tx) error ***REMOVED***
		services, err := store.FindServices(tx, store.ByReferencedNetworkID(id))
		if err != nil ***REMOVED***
			return status.Errorf(codes.Internal, "could not find services using network %s: %v", id, err)
		***REMOVED***

		if len(services) != 0 ***REMOVED***
			return status.Errorf(codes.FailedPrecondition, "network %s is in use by service %s", id, services[0].ID)
		***REMOVED***

		tasks, err := store.FindTasks(tx, store.ByReferencedNetworkID(id))
		if err != nil ***REMOVED***
			return status.Errorf(codes.Internal, "could not find tasks using network %s: %v", id, err)
		***REMOVED***

		for _, t := range tasks ***REMOVED***
			if t.DesiredState <= api.TaskStateRunning && t.Status.State <= api.TaskStateRunning ***REMOVED***
				return status.Errorf(codes.FailedPrecondition, "network %s is in use by task %s", id, t.ID)
			***REMOVED***
		***REMOVED***

		return store.DeleteNetwork(tx, id)
	***REMOVED***)
***REMOVED***

func (s *Server) removeIngressNetwork(id string) error ***REMOVED***
	return s.store.Update(func(tx store.Tx) error ***REMOVED***
		services, err := store.FindServices(tx, store.All)
		if err != nil ***REMOVED***
			return status.Errorf(codes.Internal, "could not find services using network %s: %v", id, err)
		***REMOVED***
		for _, srv := range services ***REMOVED***
			if allocator.IsIngressNetworkNeeded(srv) ***REMOVED***
				return status.Errorf(codes.FailedPrecondition, "ingress network cannot be removed because service %s depends on it", srv.ID)
			***REMOVED***
		***REMOVED***
		return store.DeleteNetwork(tx, id)
	***REMOVED***)
***REMOVED***

func filterNetworks(candidates []*api.Network, filters ...func(*api.Network) bool) []*api.Network ***REMOVED***
	result := []*api.Network***REMOVED******REMOVED***

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

// ListNetworks returns a list of all networks.
func (s *Server) ListNetworks(ctx context.Context, request *api.ListNetworksRequest) (*api.ListNetworksResponse, error) ***REMOVED***
	var (
		networks []*api.Network
		err      error
	)

	s.store.View(func(tx store.ReadTx) ***REMOVED***
		switch ***REMOVED***
		case request.Filters != nil && len(request.Filters.Names) > 0:
			networks, err = store.FindNetworks(tx, buildFilters(store.ByName, request.Filters.Names))
		case request.Filters != nil && len(request.Filters.NamePrefixes) > 0:
			networks, err = store.FindNetworks(tx, buildFilters(store.ByNamePrefix, request.Filters.NamePrefixes))
		case request.Filters != nil && len(request.Filters.IDPrefixes) > 0:
			networks, err = store.FindNetworks(tx, buildFilters(store.ByIDPrefix, request.Filters.IDPrefixes))
		default:
			networks, err = store.FindNetworks(tx, store.All)
		***REMOVED***
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if request.Filters != nil ***REMOVED***
		networks = filterNetworks(networks,
			func(e *api.Network) bool ***REMOVED***
				return filterContains(e.Spec.Annotations.Name, request.Filters.Names)
			***REMOVED***,
			func(e *api.Network) bool ***REMOVED***
				return filterContainsPrefix(e.Spec.Annotations.Name, request.Filters.NamePrefixes)
			***REMOVED***,
			func(e *api.Network) bool ***REMOVED***
				return filterContainsPrefix(e.ID, request.Filters.IDPrefixes)
			***REMOVED***,
			func(e *api.Network) bool ***REMOVED***
				return filterMatchLabels(e.Spec.Annotations.Labels, request.Filters.Labels)
			***REMOVED***,
		)
	***REMOVED***

	return &api.ListNetworksResponse***REMOVED***
		Networks: networks,
	***REMOVED***, nil
***REMOVED***
