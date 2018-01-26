package convert

import (
	swarmtypes "github.com/docker/docker/api/types/swarm"
	swarmapi "github.com/docker/swarmkit/api"
	gogotypes "github.com/gogo/protobuf/types"
)

// ConfigFromGRPC converts a grpc Config to a Config.
func ConfigFromGRPC(s *swarmapi.Config) swarmtypes.Config ***REMOVED***
	config := swarmtypes.Config***REMOVED***
		ID: s.ID,
		Spec: swarmtypes.ConfigSpec***REMOVED***
			Annotations: annotationsFromGRPC(s.Spec.Annotations),
			Data:        s.Spec.Data,
		***REMOVED***,
	***REMOVED***

	config.Version.Index = s.Meta.Version.Index
	// Meta
	config.CreatedAt, _ = gogotypes.TimestampFromProto(s.Meta.CreatedAt)
	config.UpdatedAt, _ = gogotypes.TimestampFromProto(s.Meta.UpdatedAt)

	return config
***REMOVED***

// ConfigSpecToGRPC converts Config to a grpc Config.
func ConfigSpecToGRPC(s swarmtypes.ConfigSpec) swarmapi.ConfigSpec ***REMOVED***
	return swarmapi.ConfigSpec***REMOVED***
		Annotations: swarmapi.Annotations***REMOVED***
			Name:   s.Name,
			Labels: s.Labels,
		***REMOVED***,
		Data: s.Data,
	***REMOVED***
***REMOVED***

// ConfigReferencesFromGRPC converts a slice of grpc ConfigReference to ConfigReference
func ConfigReferencesFromGRPC(s []*swarmapi.ConfigReference) []*swarmtypes.ConfigReference ***REMOVED***
	refs := []*swarmtypes.ConfigReference***REMOVED******REMOVED***

	for _, r := range s ***REMOVED***
		ref := &swarmtypes.ConfigReference***REMOVED***
			ConfigID:   r.ConfigID,
			ConfigName: r.ConfigName,
		***REMOVED***

		if t, ok := r.Target.(*swarmapi.ConfigReference_File); ok ***REMOVED***
			ref.File = &swarmtypes.ConfigReferenceFileTarget***REMOVED***
				Name: t.File.Name,
				UID:  t.File.UID,
				GID:  t.File.GID,
				Mode: t.File.Mode,
			***REMOVED***
		***REMOVED***

		refs = append(refs, ref)
	***REMOVED***

	return refs
***REMOVED***
