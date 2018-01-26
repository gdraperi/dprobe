package convert

import (
	swarmtypes "github.com/docker/docker/api/types/swarm"
	swarmapi "github.com/docker/swarmkit/api"
	gogotypes "github.com/gogo/protobuf/types"
)

// SecretFromGRPC converts a grpc Secret to a Secret.
func SecretFromGRPC(s *swarmapi.Secret) swarmtypes.Secret ***REMOVED***
	secret := swarmtypes.Secret***REMOVED***
		ID: s.ID,
		Spec: swarmtypes.SecretSpec***REMOVED***
			Annotations: annotationsFromGRPC(s.Spec.Annotations),
			Data:        s.Spec.Data,
			Driver:      driverFromGRPC(s.Spec.Driver),
		***REMOVED***,
	***REMOVED***

	secret.Version.Index = s.Meta.Version.Index
	// Meta
	secret.CreatedAt, _ = gogotypes.TimestampFromProto(s.Meta.CreatedAt)
	secret.UpdatedAt, _ = gogotypes.TimestampFromProto(s.Meta.UpdatedAt)

	return secret
***REMOVED***

// SecretSpecToGRPC converts Secret to a grpc Secret.
func SecretSpecToGRPC(s swarmtypes.SecretSpec) swarmapi.SecretSpec ***REMOVED***
	return swarmapi.SecretSpec***REMOVED***
		Annotations: swarmapi.Annotations***REMOVED***
			Name:   s.Name,
			Labels: s.Labels,
		***REMOVED***,
		Data:   s.Data,
		Driver: driverToGRPC(s.Driver),
	***REMOVED***
***REMOVED***

// SecretReferencesFromGRPC converts a slice of grpc SecretReference to SecretReference
func SecretReferencesFromGRPC(s []*swarmapi.SecretReference) []*swarmtypes.SecretReference ***REMOVED***
	refs := []*swarmtypes.SecretReference***REMOVED******REMOVED***

	for _, r := range s ***REMOVED***
		ref := &swarmtypes.SecretReference***REMOVED***
			SecretID:   r.SecretID,
			SecretName: r.SecretName,
		***REMOVED***

		if t, ok := r.Target.(*swarmapi.SecretReference_File); ok ***REMOVED***
			ref.File = &swarmtypes.SecretReferenceFileTarget***REMOVED***
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
