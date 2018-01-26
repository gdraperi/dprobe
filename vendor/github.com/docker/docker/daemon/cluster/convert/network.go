package convert

import (
	"strings"

	basictypes "github.com/docker/docker/api/types"
	networktypes "github.com/docker/docker/api/types/network"
	types "github.com/docker/docker/api/types/swarm"
	netconst "github.com/docker/libnetwork/datastore"
	swarmapi "github.com/docker/swarmkit/api"
	gogotypes "github.com/gogo/protobuf/types"
)

func networkAttachmentFromGRPC(na *swarmapi.NetworkAttachment) types.NetworkAttachment ***REMOVED***
	if na != nil ***REMOVED***
		return types.NetworkAttachment***REMOVED***
			Network:   networkFromGRPC(na.Network),
			Addresses: na.Addresses,
		***REMOVED***
	***REMOVED***
	return types.NetworkAttachment***REMOVED******REMOVED***
***REMOVED***

func networkFromGRPC(n *swarmapi.Network) types.Network ***REMOVED***
	if n != nil ***REMOVED***
		network := types.Network***REMOVED***
			ID: n.ID,
			Spec: types.NetworkSpec***REMOVED***
				IPv6Enabled: n.Spec.Ipv6Enabled,
				Internal:    n.Spec.Internal,
				Attachable:  n.Spec.Attachable,
				Ingress:     IsIngressNetwork(n),
				IPAMOptions: ipamFromGRPC(n.Spec.IPAM),
				Scope:       netconst.SwarmScope,
			***REMOVED***,
			IPAMOptions: ipamFromGRPC(n.IPAM),
		***REMOVED***

		if n.Spec.GetNetwork() != "" ***REMOVED***
			network.Spec.ConfigFrom = &networktypes.ConfigReference***REMOVED***
				Network: n.Spec.GetNetwork(),
			***REMOVED***
		***REMOVED***

		// Meta
		network.Version.Index = n.Meta.Version.Index
		network.CreatedAt, _ = gogotypes.TimestampFromProto(n.Meta.CreatedAt)
		network.UpdatedAt, _ = gogotypes.TimestampFromProto(n.Meta.UpdatedAt)

		//Annotations
		network.Spec.Annotations = annotationsFromGRPC(n.Spec.Annotations)

		//DriverConfiguration
		if n.Spec.DriverConfig != nil ***REMOVED***
			network.Spec.DriverConfiguration = &types.Driver***REMOVED***
				Name:    n.Spec.DriverConfig.Name,
				Options: n.Spec.DriverConfig.Options,
			***REMOVED***
		***REMOVED***

		//DriverState
		if n.DriverState != nil ***REMOVED***
			network.DriverState = types.Driver***REMOVED***
				Name:    n.DriverState.Name,
				Options: n.DriverState.Options,
			***REMOVED***
		***REMOVED***

		return network
	***REMOVED***
	return types.Network***REMOVED******REMOVED***
***REMOVED***

func ipamFromGRPC(i *swarmapi.IPAMOptions) *types.IPAMOptions ***REMOVED***
	var ipam *types.IPAMOptions
	if i != nil ***REMOVED***
		ipam = &types.IPAMOptions***REMOVED******REMOVED***
		if i.Driver != nil ***REMOVED***
			ipam.Driver.Name = i.Driver.Name
			ipam.Driver.Options = i.Driver.Options
		***REMOVED***

		for _, config := range i.Configs ***REMOVED***
			ipam.Configs = append(ipam.Configs, types.IPAMConfig***REMOVED***
				Subnet:  config.Subnet,
				Range:   config.Range,
				Gateway: config.Gateway,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	return ipam
***REMOVED***

func endpointSpecFromGRPC(es *swarmapi.EndpointSpec) *types.EndpointSpec ***REMOVED***
	var endpointSpec *types.EndpointSpec
	if es != nil ***REMOVED***
		endpointSpec = &types.EndpointSpec***REMOVED******REMOVED***
		endpointSpec.Mode = types.ResolutionMode(strings.ToLower(es.Mode.String()))

		for _, portState := range es.Ports ***REMOVED***
			endpointSpec.Ports = append(endpointSpec.Ports, swarmPortConfigToAPIPortConfig(portState))
		***REMOVED***
	***REMOVED***
	return endpointSpec
***REMOVED***

func endpointFromGRPC(e *swarmapi.Endpoint) types.Endpoint ***REMOVED***
	endpoint := types.Endpoint***REMOVED******REMOVED***
	if e != nil ***REMOVED***
		if espec := endpointSpecFromGRPC(e.Spec); espec != nil ***REMOVED***
			endpoint.Spec = *espec
		***REMOVED***

		for _, portState := range e.Ports ***REMOVED***
			endpoint.Ports = append(endpoint.Ports, swarmPortConfigToAPIPortConfig(portState))
		***REMOVED***

		for _, v := range e.VirtualIPs ***REMOVED***
			endpoint.VirtualIPs = append(endpoint.VirtualIPs, types.EndpointVirtualIP***REMOVED***
				NetworkID: v.NetworkID,
				Addr:      v.Addr***REMOVED***)
		***REMOVED***

	***REMOVED***

	return endpoint
***REMOVED***

func swarmPortConfigToAPIPortConfig(portConfig *swarmapi.PortConfig) types.PortConfig ***REMOVED***
	return types.PortConfig***REMOVED***
		Name:          portConfig.Name,
		Protocol:      types.PortConfigProtocol(strings.ToLower(swarmapi.PortConfig_Protocol_name[int32(portConfig.Protocol)])),
		PublishMode:   types.PortConfigPublishMode(strings.ToLower(swarmapi.PortConfig_PublishMode_name[int32(portConfig.PublishMode)])),
		TargetPort:    portConfig.TargetPort,
		PublishedPort: portConfig.PublishedPort,
	***REMOVED***
***REMOVED***

// BasicNetworkFromGRPC converts a grpc Network to a NetworkResource.
func BasicNetworkFromGRPC(n swarmapi.Network) basictypes.NetworkResource ***REMOVED***
	spec := n.Spec
	var ipam networktypes.IPAM
	if spec.IPAM != nil ***REMOVED***
		if spec.IPAM.Driver != nil ***REMOVED***
			ipam.Driver = spec.IPAM.Driver.Name
			ipam.Options = spec.IPAM.Driver.Options
		***REMOVED***
		ipam.Config = make([]networktypes.IPAMConfig, 0, len(spec.IPAM.Configs))
		for _, ic := range spec.IPAM.Configs ***REMOVED***
			ipamConfig := networktypes.IPAMConfig***REMOVED***
				Subnet:     ic.Subnet,
				IPRange:    ic.Range,
				Gateway:    ic.Gateway,
				AuxAddress: ic.Reserved,
			***REMOVED***
			ipam.Config = append(ipam.Config, ipamConfig)
		***REMOVED***
	***REMOVED***

	nr := basictypes.NetworkResource***REMOVED***
		ID:         n.ID,
		Name:       n.Spec.Annotations.Name,
		Scope:      netconst.SwarmScope,
		EnableIPv6: spec.Ipv6Enabled,
		IPAM:       ipam,
		Internal:   spec.Internal,
		Attachable: spec.Attachable,
		Ingress:    IsIngressNetwork(&n),
		Labels:     n.Spec.Annotations.Labels,
	***REMOVED***

	if n.Spec.GetNetwork() != "" ***REMOVED***
		nr.ConfigFrom = networktypes.ConfigReference***REMOVED***
			Network: n.Spec.GetNetwork(),
		***REMOVED***
	***REMOVED***

	if n.DriverState != nil ***REMOVED***
		nr.Driver = n.DriverState.Name
		nr.Options = n.DriverState.Options
	***REMOVED***

	return nr
***REMOVED***

// BasicNetworkCreateToGRPC converts a NetworkCreateRequest to a grpc NetworkSpec.
func BasicNetworkCreateToGRPC(create basictypes.NetworkCreateRequest) swarmapi.NetworkSpec ***REMOVED***
	ns := swarmapi.NetworkSpec***REMOVED***
		Annotations: swarmapi.Annotations***REMOVED***
			Name:   create.Name,
			Labels: create.Labels,
		***REMOVED***,
		DriverConfig: &swarmapi.Driver***REMOVED***
			Name:    create.Driver,
			Options: create.Options,
		***REMOVED***,
		Ipv6Enabled: create.EnableIPv6,
		Internal:    create.Internal,
		Attachable:  create.Attachable,
		Ingress:     create.Ingress,
	***REMOVED***
	if create.IPAM != nil ***REMOVED***
		driver := create.IPAM.Driver
		if driver == "" ***REMOVED***
			driver = "default"
		***REMOVED***
		ns.IPAM = &swarmapi.IPAMOptions***REMOVED***
			Driver: &swarmapi.Driver***REMOVED***
				Name:    driver,
				Options: create.IPAM.Options,
			***REMOVED***,
		***REMOVED***
		ipamSpec := make([]*swarmapi.IPAMConfig, 0, len(create.IPAM.Config))
		for _, ipamConfig := range create.IPAM.Config ***REMOVED***
			ipamSpec = append(ipamSpec, &swarmapi.IPAMConfig***REMOVED***
				Subnet:  ipamConfig.Subnet,
				Range:   ipamConfig.IPRange,
				Gateway: ipamConfig.Gateway,
			***REMOVED***)
		***REMOVED***
		ns.IPAM.Configs = ipamSpec
	***REMOVED***
	if create.ConfigFrom != nil ***REMOVED***
		ns.ConfigFrom = &swarmapi.NetworkSpec_Network***REMOVED***
			Network: create.ConfigFrom.Network,
		***REMOVED***
	***REMOVED***
	return ns
***REMOVED***

// IsIngressNetwork check if the swarm network is an ingress network
func IsIngressNetwork(n *swarmapi.Network) bool ***REMOVED***
	if n.Spec.Ingress ***REMOVED***
		return true
	***REMOVED***
	// Check if legacy defined ingress network
	_, ok := n.Spec.Annotations.Labels["com.docker.swarm.internal"]
	return ok && n.Spec.Annotations.Name == "ingress"
***REMOVED***
