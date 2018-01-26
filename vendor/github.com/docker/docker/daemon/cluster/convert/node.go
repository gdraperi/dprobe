package convert

import (
	"fmt"
	"strings"

	types "github.com/docker/docker/api/types/swarm"
	swarmapi "github.com/docker/swarmkit/api"
	gogotypes "github.com/gogo/protobuf/types"
)

// NodeFromGRPC converts a grpc Node to a Node.
func NodeFromGRPC(n swarmapi.Node) types.Node ***REMOVED***
	node := types.Node***REMOVED***
		ID: n.ID,
		Spec: types.NodeSpec***REMOVED***
			Role:         types.NodeRole(strings.ToLower(n.Spec.DesiredRole.String())),
			Availability: types.NodeAvailability(strings.ToLower(n.Spec.Availability.String())),
		***REMOVED***,
		Status: types.NodeStatus***REMOVED***
			State:   types.NodeState(strings.ToLower(n.Status.State.String())),
			Message: n.Status.Message,
			Addr:    n.Status.Addr,
		***REMOVED***,
	***REMOVED***

	// Meta
	node.Version.Index = n.Meta.Version.Index
	node.CreatedAt, _ = gogotypes.TimestampFromProto(n.Meta.CreatedAt)
	node.UpdatedAt, _ = gogotypes.TimestampFromProto(n.Meta.UpdatedAt)

	//Annotations
	node.Spec.Annotations = annotationsFromGRPC(n.Spec.Annotations)

	//Description
	if n.Description != nil ***REMOVED***
		node.Description.Hostname = n.Description.Hostname
		if n.Description.Platform != nil ***REMOVED***
			node.Description.Platform.Architecture = n.Description.Platform.Architecture
			node.Description.Platform.OS = n.Description.Platform.OS
		***REMOVED***
		if n.Description.Resources != nil ***REMOVED***
			node.Description.Resources.NanoCPUs = n.Description.Resources.NanoCPUs
			node.Description.Resources.MemoryBytes = n.Description.Resources.MemoryBytes
			node.Description.Resources.GenericResources = GenericResourcesFromGRPC(n.Description.Resources.Generic)
		***REMOVED***
		if n.Description.Engine != nil ***REMOVED***
			node.Description.Engine.EngineVersion = n.Description.Engine.EngineVersion
			node.Description.Engine.Labels = n.Description.Engine.Labels
			for _, plugin := range n.Description.Engine.Plugins ***REMOVED***
				node.Description.Engine.Plugins = append(node.Description.Engine.Plugins, types.PluginDescription***REMOVED***Type: plugin.Type, Name: plugin.Name***REMOVED***)
			***REMOVED***
		***REMOVED***
		if n.Description.TLSInfo != nil ***REMOVED***
			node.Description.TLSInfo.TrustRoot = string(n.Description.TLSInfo.TrustRoot)
			node.Description.TLSInfo.CertIssuerPublicKey = n.Description.TLSInfo.CertIssuerPublicKey
			node.Description.TLSInfo.CertIssuerSubject = n.Description.TLSInfo.CertIssuerSubject
		***REMOVED***
	***REMOVED***

	//Manager
	if n.ManagerStatus != nil ***REMOVED***
		node.ManagerStatus = &types.ManagerStatus***REMOVED***
			Leader:       n.ManagerStatus.Leader,
			Reachability: types.Reachability(strings.ToLower(n.ManagerStatus.Reachability.String())),
			Addr:         n.ManagerStatus.Addr,
		***REMOVED***
	***REMOVED***

	return node
***REMOVED***

// NodeSpecToGRPC converts a NodeSpec to a grpc NodeSpec.
func NodeSpecToGRPC(s types.NodeSpec) (swarmapi.NodeSpec, error) ***REMOVED***
	spec := swarmapi.NodeSpec***REMOVED***
		Annotations: swarmapi.Annotations***REMOVED***
			Name:   s.Name,
			Labels: s.Labels,
		***REMOVED***,
	***REMOVED***
	if role, ok := swarmapi.NodeRole_value[strings.ToUpper(string(s.Role))]; ok ***REMOVED***
		spec.DesiredRole = swarmapi.NodeRole(role)
	***REMOVED*** else ***REMOVED***
		return swarmapi.NodeSpec***REMOVED******REMOVED***, fmt.Errorf("invalid Role: %q", s.Role)
	***REMOVED***

	if availability, ok := swarmapi.NodeSpec_Availability_value[strings.ToUpper(string(s.Availability))]; ok ***REMOVED***
		spec.Availability = swarmapi.NodeSpec_Availability(availability)
	***REMOVED*** else ***REMOVED***
		return swarmapi.NodeSpec***REMOVED******REMOVED***, fmt.Errorf("invalid Availability: %q", s.Availability)
	***REMOVED***

	return spec, nil
***REMOVED***
