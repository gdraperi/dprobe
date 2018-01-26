package convert

import (
	"fmt"
	"strings"

	types "github.com/docker/docker/api/types/swarm"
	swarmapi "github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/ca"
	gogotypes "github.com/gogo/protobuf/types"
)

// SwarmFromGRPC converts a grpc Cluster to a Swarm.
func SwarmFromGRPC(c swarmapi.Cluster) types.Swarm ***REMOVED***
	swarm := types.Swarm***REMOVED***
		ClusterInfo: types.ClusterInfo***REMOVED***
			ID: c.ID,
			Spec: types.Spec***REMOVED***
				Orchestration: types.OrchestrationConfig***REMOVED***
					TaskHistoryRetentionLimit: &c.Spec.Orchestration.TaskHistoryRetentionLimit,
				***REMOVED***,
				Raft: types.RaftConfig***REMOVED***
					SnapshotInterval:           c.Spec.Raft.SnapshotInterval,
					KeepOldSnapshots:           &c.Spec.Raft.KeepOldSnapshots,
					LogEntriesForSlowFollowers: c.Spec.Raft.LogEntriesForSlowFollowers,
					HeartbeatTick:              int(c.Spec.Raft.HeartbeatTick),
					ElectionTick:               int(c.Spec.Raft.ElectionTick),
				***REMOVED***,
				EncryptionConfig: types.EncryptionConfig***REMOVED***
					AutoLockManagers: c.Spec.EncryptionConfig.AutoLockManagers,
				***REMOVED***,
				CAConfig: types.CAConfig***REMOVED***
					// do not include the signing CA cert or key (it should already be redacted via the swarm APIs) -
					// the key because it's secret, and the cert because otherwise doing a get + update on the spec
					// can cause issues because the key would be missing and the cert wouldn't
					ForceRotate: c.Spec.CAConfig.ForceRotate,
				***REMOVED***,
			***REMOVED***,
			TLSInfo: types.TLSInfo***REMOVED***
				TrustRoot: string(c.RootCA.CACert),
			***REMOVED***,
			RootRotationInProgress: c.RootCA.RootRotation != nil,
		***REMOVED***,
		JoinTokens: types.JoinTokens***REMOVED***
			Worker:  c.RootCA.JoinTokens.Worker,
			Manager: c.RootCA.JoinTokens.Manager,
		***REMOVED***,
	***REMOVED***

	issuerInfo, err := ca.IssuerFromAPIRootCA(&c.RootCA)
	if err == nil && issuerInfo != nil ***REMOVED***
		swarm.TLSInfo.CertIssuerSubject = issuerInfo.Subject
		swarm.TLSInfo.CertIssuerPublicKey = issuerInfo.PublicKey
	***REMOVED***

	heartbeatPeriod, _ := gogotypes.DurationFromProto(c.Spec.Dispatcher.HeartbeatPeriod)
	swarm.Spec.Dispatcher.HeartbeatPeriod = heartbeatPeriod

	swarm.Spec.CAConfig.NodeCertExpiry, _ = gogotypes.DurationFromProto(c.Spec.CAConfig.NodeCertExpiry)

	for _, ca := range c.Spec.CAConfig.ExternalCAs ***REMOVED***
		swarm.Spec.CAConfig.ExternalCAs = append(swarm.Spec.CAConfig.ExternalCAs, &types.ExternalCA***REMOVED***
			Protocol: types.ExternalCAProtocol(strings.ToLower(ca.Protocol.String())),
			URL:      ca.URL,
			Options:  ca.Options,
			CACert:   string(ca.CACert),
		***REMOVED***)
	***REMOVED***

	// Meta
	swarm.Version.Index = c.Meta.Version.Index
	swarm.CreatedAt, _ = gogotypes.TimestampFromProto(c.Meta.CreatedAt)
	swarm.UpdatedAt, _ = gogotypes.TimestampFromProto(c.Meta.UpdatedAt)

	// Annotations
	swarm.Spec.Annotations = annotationsFromGRPC(c.Spec.Annotations)

	return swarm
***REMOVED***

// SwarmSpecToGRPC converts a Spec to a grpc ClusterSpec.
func SwarmSpecToGRPC(s types.Spec) (swarmapi.ClusterSpec, error) ***REMOVED***
	return MergeSwarmSpecToGRPC(s, swarmapi.ClusterSpec***REMOVED******REMOVED***)
***REMOVED***

// MergeSwarmSpecToGRPC merges a Spec with an initial grpc ClusterSpec
func MergeSwarmSpecToGRPC(s types.Spec, spec swarmapi.ClusterSpec) (swarmapi.ClusterSpec, error) ***REMOVED***
	// We take the initSpec (either created from scratch, or returned by swarmkit),
	// and will only change the value if the one taken from types.Spec is not nil or 0.
	// In other words, if the value taken from types.Spec is nil or 0, we will maintain the status quo.
	if s.Annotations.Name != "" ***REMOVED***
		spec.Annotations.Name = s.Annotations.Name
	***REMOVED***
	if len(s.Annotations.Labels) != 0 ***REMOVED***
		spec.Annotations.Labels = s.Annotations.Labels
	***REMOVED***

	if s.Orchestration.TaskHistoryRetentionLimit != nil ***REMOVED***
		spec.Orchestration.TaskHistoryRetentionLimit = *s.Orchestration.TaskHistoryRetentionLimit
	***REMOVED***
	if s.Raft.SnapshotInterval != 0 ***REMOVED***
		spec.Raft.SnapshotInterval = s.Raft.SnapshotInterval
	***REMOVED***
	if s.Raft.KeepOldSnapshots != nil ***REMOVED***
		spec.Raft.KeepOldSnapshots = *s.Raft.KeepOldSnapshots
	***REMOVED***
	if s.Raft.LogEntriesForSlowFollowers != 0 ***REMOVED***
		spec.Raft.LogEntriesForSlowFollowers = s.Raft.LogEntriesForSlowFollowers
	***REMOVED***
	if s.Raft.HeartbeatTick != 0 ***REMOVED***
		spec.Raft.HeartbeatTick = uint32(s.Raft.HeartbeatTick)
	***REMOVED***
	if s.Raft.ElectionTick != 0 ***REMOVED***
		spec.Raft.ElectionTick = uint32(s.Raft.ElectionTick)
	***REMOVED***
	if s.Dispatcher.HeartbeatPeriod != 0 ***REMOVED***
		spec.Dispatcher.HeartbeatPeriod = gogotypes.DurationProto(s.Dispatcher.HeartbeatPeriod)
	***REMOVED***
	if s.CAConfig.NodeCertExpiry != 0 ***REMOVED***
		spec.CAConfig.NodeCertExpiry = gogotypes.DurationProto(s.CAConfig.NodeCertExpiry)
	***REMOVED***
	if s.CAConfig.SigningCACert != "" ***REMOVED***
		spec.CAConfig.SigningCACert = []byte(s.CAConfig.SigningCACert)
	***REMOVED***
	if s.CAConfig.SigningCAKey != "" ***REMOVED***
		// do propagate the signing CA key here because we want to provide it TO the swarm APIs
		spec.CAConfig.SigningCAKey = []byte(s.CAConfig.SigningCAKey)
	***REMOVED***
	spec.CAConfig.ForceRotate = s.CAConfig.ForceRotate

	for _, ca := range s.CAConfig.ExternalCAs ***REMOVED***
		protocol, ok := swarmapi.ExternalCA_CAProtocol_value[strings.ToUpper(string(ca.Protocol))]
		if !ok ***REMOVED***
			return swarmapi.ClusterSpec***REMOVED******REMOVED***, fmt.Errorf("invalid protocol: %q", ca.Protocol)
		***REMOVED***
		spec.CAConfig.ExternalCAs = append(spec.CAConfig.ExternalCAs, &swarmapi.ExternalCA***REMOVED***
			Protocol: swarmapi.ExternalCA_CAProtocol(protocol),
			URL:      ca.URL,
			Options:  ca.Options,
			CACert:   []byte(ca.CACert),
		***REMOVED***)
	***REMOVED***

	spec.EncryptionConfig.AutoLockManagers = s.EncryptionConfig.AutoLockManagers

	return spec, nil
***REMOVED***
