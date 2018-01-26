package controlapi

import (
	"strings"
	"time"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/ca"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/encryption"
	"github.com/docker/swarmkit/manager/state/store"
	gogotypes "github.com/gogo/protobuf/types"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// expiredCertGrace is the amount of time to keep a node in the
	// blacklist beyond its certificate expiration timestamp.
	expiredCertGrace = 24 * time.Hour * 7
)

func validateClusterSpec(spec *api.ClusterSpec) error ***REMOVED***
	if spec == nil ***REMOVED***
		return status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***

	// Validate that expiry time being provided is valid, and over our minimum
	if spec.CAConfig.NodeCertExpiry != nil ***REMOVED***
		expiry, err := gogotypes.DurationFromProto(spec.CAConfig.NodeCertExpiry)
		if err != nil ***REMOVED***
			return status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
		***REMOVED***
		if expiry < ca.MinNodeCertExpiration ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "minimum certificate expiry time is: %s", ca.MinNodeCertExpiration)
		***REMOVED***
	***REMOVED***

	// Validate that AcceptancePolicies only include Secrets that are bcrypted
	// TODO(diogo): Add a global list of acceptance algorithms. We only support bcrypt for now.
	if len(spec.AcceptancePolicy.Policies) > 0 ***REMOVED***
		for _, policy := range spec.AcceptancePolicy.Policies ***REMOVED***
			if policy.Secret != nil && strings.ToLower(policy.Secret.Alg) != "bcrypt" ***REMOVED***
				return status.Errorf(codes.InvalidArgument, "hashing algorithm is not supported: %s", policy.Secret.Alg)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Validate that heartbeatPeriod time being provided is valid
	if spec.Dispatcher.HeartbeatPeriod != nil ***REMOVED***
		heartbeatPeriod, err := gogotypes.DurationFromProto(spec.Dispatcher.HeartbeatPeriod)
		if err != nil ***REMOVED***
			return status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
		***REMOVED***
		if heartbeatPeriod < 0 ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "heartbeat time period cannot be a negative duration")
		***REMOVED***
	***REMOVED***

	if spec.Annotations.Name != store.DefaultClusterName ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "modification of cluster name is not allowed")
	***REMOVED***

	return nil
***REMOVED***

// GetCluster returns a Cluster given a ClusterID.
// - Returns `InvalidArgument` if ClusterID is not provided.
// - Returns `NotFound` if the Cluster is not found.
func (s *Server) GetCluster(ctx context.Context, request *api.GetClusterRequest) (*api.GetClusterResponse, error) ***REMOVED***
	if request.ClusterID == "" ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***

	var cluster *api.Cluster
	s.store.View(func(tx store.ReadTx) ***REMOVED***
		cluster = store.GetCluster(tx, request.ClusterID)
	***REMOVED***)
	if cluster == nil ***REMOVED***
		return nil, status.Errorf(codes.NotFound, "cluster %s not found", request.ClusterID)
	***REMOVED***

	redactedClusters := redactClusters([]*api.Cluster***REMOVED***cluster***REMOVED***)

	// WARN: we should never return cluster here. We need to redact the private fields first.
	return &api.GetClusterResponse***REMOVED***
		Cluster: redactedClusters[0],
	***REMOVED***, nil
***REMOVED***

// UpdateCluster updates a Cluster referenced by ClusterID with the given ClusterSpec.
// - Returns `NotFound` if the Cluster is not found.
// - Returns `InvalidArgument` if the ClusterSpec is malformed.
// - Returns `Unimplemented` if the ClusterSpec references unimplemented features.
// - Returns an error if the update fails.
func (s *Server) UpdateCluster(ctx context.Context, request *api.UpdateClusterRequest) (*api.UpdateClusterResponse, error) ***REMOVED***
	if request.ClusterID == "" || request.ClusterVersion == nil ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***
	if err := validateClusterSpec(request.Spec); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var cluster *api.Cluster
	err := s.store.Update(func(tx store.Tx) error ***REMOVED***
		cluster = store.GetCluster(tx, request.ClusterID)
		if cluster == nil ***REMOVED***
			return status.Errorf(codes.NotFound, "cluster %s not found", request.ClusterID)
		***REMOVED***
		// This ensures that we have the current rootCA with which to generate tokens (expiration doesn't matter
		// for generating the tokens)
		rootCA, err := ca.RootCAFromAPI(ctx, &cluster.RootCA, ca.DefaultNodeCertExpiration)
		if err != nil ***REMOVED***
			log.G(ctx).WithField(
				"method", "(*controlapi.Server).UpdateCluster").WithError(err).Error("invalid cluster root CA")
			return status.Errorf(codes.Internal, "error loading cluster rootCA for update")
		***REMOVED***

		cluster.Meta.Version = *request.ClusterVersion
		cluster.Spec = *request.Spec.Copy()

		expireBlacklistedCerts(cluster)

		if request.Rotation.WorkerJoinToken ***REMOVED***
			cluster.RootCA.JoinTokens.Worker = ca.GenerateJoinToken(&rootCA)
		***REMOVED***
		if request.Rotation.ManagerJoinToken ***REMOVED***
			cluster.RootCA.JoinTokens.Manager = ca.GenerateJoinToken(&rootCA)
		***REMOVED***

		updatedRootCA, err := validateCAConfig(ctx, s.securityConfig, cluster)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		cluster.RootCA = *updatedRootCA

		var unlockKeys []*api.EncryptionKey
		var managerKey *api.EncryptionKey
		for _, eKey := range cluster.UnlockKeys ***REMOVED***
			if eKey.Subsystem == ca.ManagerRole ***REMOVED***
				if !cluster.Spec.EncryptionConfig.AutoLockManagers ***REMOVED***
					continue
				***REMOVED***
				managerKey = eKey
			***REMOVED***
			unlockKeys = append(unlockKeys, eKey)
		***REMOVED***

		switch ***REMOVED***
		case !cluster.Spec.EncryptionConfig.AutoLockManagers:
			break
		case managerKey == nil:
			unlockKeys = append(unlockKeys, &api.EncryptionKey***REMOVED***
				Subsystem: ca.ManagerRole,
				Key:       encryption.GenerateSecretKey(),
			***REMOVED***)
		case request.Rotation.ManagerUnlockKey:
			managerKey.Key = encryption.GenerateSecretKey()
		***REMOVED***
		cluster.UnlockKeys = unlockKeys

		return store.UpdateCluster(tx, cluster)
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	redactedClusters := redactClusters([]*api.Cluster***REMOVED***cluster***REMOVED***)

	// WARN: we should never return cluster here. We need to redact the private fields first.
	return &api.UpdateClusterResponse***REMOVED***
		Cluster: redactedClusters[0],
	***REMOVED***, nil
***REMOVED***

func filterClusters(candidates []*api.Cluster, filters ...func(*api.Cluster) bool) []*api.Cluster ***REMOVED***
	result := []*api.Cluster***REMOVED******REMOVED***

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

// ListClusters returns a list of all clusters.
func (s *Server) ListClusters(ctx context.Context, request *api.ListClustersRequest) (*api.ListClustersResponse, error) ***REMOVED***
	var (
		clusters []*api.Cluster
		err      error
	)
	s.store.View(func(tx store.ReadTx) ***REMOVED***
		switch ***REMOVED***
		case request.Filters != nil && len(request.Filters.Names) > 0:
			clusters, err = store.FindClusters(tx, buildFilters(store.ByName, request.Filters.Names))
		case request.Filters != nil && len(request.Filters.NamePrefixes) > 0:
			clusters, err = store.FindClusters(tx, buildFilters(store.ByNamePrefix, request.Filters.NamePrefixes))
		case request.Filters != nil && len(request.Filters.IDPrefixes) > 0:
			clusters, err = store.FindClusters(tx, buildFilters(store.ByIDPrefix, request.Filters.IDPrefixes))
		default:
			clusters, err = store.FindClusters(tx, store.All)
		***REMOVED***
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if request.Filters != nil ***REMOVED***
		clusters = filterClusters(clusters,
			func(e *api.Cluster) bool ***REMOVED***
				return filterContains(e.Spec.Annotations.Name, request.Filters.Names)
			***REMOVED***,
			func(e *api.Cluster) bool ***REMOVED***
				return filterContainsPrefix(e.Spec.Annotations.Name, request.Filters.NamePrefixes)
			***REMOVED***,
			func(e *api.Cluster) bool ***REMOVED***
				return filterContainsPrefix(e.ID, request.Filters.IDPrefixes)
			***REMOVED***,
			func(e *api.Cluster) bool ***REMOVED***
				return filterMatchLabels(e.Spec.Annotations.Labels, request.Filters.Labels)
			***REMOVED***,
		)
	***REMOVED***

	// WARN: we should never return cluster here. We need to redact the private fields first.
	return &api.ListClustersResponse***REMOVED***
		Clusters: redactClusters(clusters),
	***REMOVED***, nil
***REMOVED***

// redactClusters is a method that enforces a whitelist of fields that are ok to be
// returned in the Cluster object. It should filter out all sensitive information.
func redactClusters(clusters []*api.Cluster) []*api.Cluster ***REMOVED***
	var redactedClusters []*api.Cluster
	// Only add public fields to the new clusters
	for _, cluster := range clusters ***REMOVED***
		// Copy all the mandatory fields
		// Do not copy secret keys
		redactedSpec := cluster.Spec.Copy()
		redactedSpec.CAConfig.SigningCAKey = nil
		// the cert is not a secret, but if API users get the cluster spec and then update,
		// then because the cert is included but not the key, the user can get update errors
		// or unintended consequences (such as telling swarm to forget about the key so long
		// as there is a corresponding external CA)
		redactedSpec.CAConfig.SigningCACert = nil

		redactedRootCA := cluster.RootCA.Copy()
		redactedRootCA.CAKey = nil
		if r := redactedRootCA.RootRotation; r != nil ***REMOVED***
			r.CAKey = nil
		***REMOVED***
		newCluster := &api.Cluster***REMOVED***
			ID:                      cluster.ID,
			Meta:                    cluster.Meta,
			Spec:                    *redactedSpec,
			RootCA:                  *redactedRootCA,
			BlacklistedCertificates: cluster.BlacklistedCertificates,
		***REMOVED***
		redactedClusters = append(redactedClusters, newCluster)
	***REMOVED***

	return redactedClusters
***REMOVED***

func expireBlacklistedCerts(cluster *api.Cluster) ***REMOVED***
	nowMinusGrace := time.Now().Add(-expiredCertGrace)

	for cn, blacklistedCert := range cluster.BlacklistedCertificates ***REMOVED***
		if blacklistedCert.Expiry == nil ***REMOVED***
			continue
		***REMOVED***

		expiry, err := gogotypes.TimestampFromProto(blacklistedCert.Expiry)
		if err == nil && nowMinusGrace.After(expiry) ***REMOVED***
			delete(cluster.BlacklistedCertificates, cn)
		***REMOVED***
	***REMOVED***
***REMOVED***
