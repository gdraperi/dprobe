package controlapi

import (
	"crypto/x509"
	"encoding/pem"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/manager/state/raft/membership"
	"github.com/docker/swarmkit/manager/state/store"
	gogotypes "github.com/gogo/protobuf/types"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateNodeSpec(spec *api.NodeSpec) error ***REMOVED***
	if spec == nil ***REMOVED***
		return status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***
	return nil
***REMOVED***

// GetNode returns a Node given a NodeID.
// - Returns `InvalidArgument` if NodeID is not provided.
// - Returns `NotFound` if the Node is not found.
func (s *Server) GetNode(ctx context.Context, request *api.GetNodeRequest) (*api.GetNodeResponse, error) ***REMOVED***
	if request.NodeID == "" ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***

	var node *api.Node
	s.store.View(func(tx store.ReadTx) ***REMOVED***
		node = store.GetNode(tx, request.NodeID)
	***REMOVED***)
	if node == nil ***REMOVED***
		return nil, status.Errorf(codes.NotFound, "node %s not found", request.NodeID)
	***REMOVED***

	if s.raft != nil ***REMOVED***
		memberlist := s.raft.GetMemberlist()
		for _, member := range memberlist ***REMOVED***
			if member.NodeID == node.ID ***REMOVED***
				node.ManagerStatus = &api.ManagerStatus***REMOVED***
					RaftID:       member.RaftID,
					Addr:         member.Addr,
					Leader:       member.Status.Leader,
					Reachability: member.Status.Reachability,
				***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return &api.GetNodeResponse***REMOVED***
		Node: node,
	***REMOVED***, nil
***REMOVED***

func filterNodes(candidates []*api.Node, filters ...func(*api.Node) bool) []*api.Node ***REMOVED***
	result := []*api.Node***REMOVED******REMOVED***

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

// ListNodes returns a list of all nodes.
func (s *Server) ListNodes(ctx context.Context, request *api.ListNodesRequest) (*api.ListNodesResponse, error) ***REMOVED***
	var (
		nodes []*api.Node
		err   error
	)
	s.store.View(func(tx store.ReadTx) ***REMOVED***
		switch ***REMOVED***
		case request.Filters != nil && len(request.Filters.Names) > 0:
			nodes, err = store.FindNodes(tx, buildFilters(store.ByName, request.Filters.Names))
		case request.Filters != nil && len(request.Filters.NamePrefixes) > 0:
			nodes, err = store.FindNodes(tx, buildFilters(store.ByNamePrefix, request.Filters.NamePrefixes))
		case request.Filters != nil && len(request.Filters.IDPrefixes) > 0:
			nodes, err = store.FindNodes(tx, buildFilters(store.ByIDPrefix, request.Filters.IDPrefixes))
		case request.Filters != nil && len(request.Filters.Roles) > 0:
			filters := make([]store.By, 0, len(request.Filters.Roles))
			for _, v := range request.Filters.Roles ***REMOVED***
				filters = append(filters, store.ByRole(v))
			***REMOVED***
			nodes, err = store.FindNodes(tx, store.Or(filters...))
		case request.Filters != nil && len(request.Filters.Memberships) > 0:
			filters := make([]store.By, 0, len(request.Filters.Memberships))
			for _, v := range request.Filters.Memberships ***REMOVED***
				filters = append(filters, store.ByMembership(v))
			***REMOVED***
			nodes, err = store.FindNodes(tx, store.Or(filters...))
		default:
			nodes, err = store.FindNodes(tx, store.All)
		***REMOVED***
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if request.Filters != nil ***REMOVED***
		nodes = filterNodes(nodes,
			func(e *api.Node) bool ***REMOVED***
				if len(request.Filters.Names) == 0 ***REMOVED***
					return true
				***REMOVED***
				if e.Description == nil ***REMOVED***
					return false
				***REMOVED***
				return filterContains(e.Description.Hostname, request.Filters.Names)
			***REMOVED***,
			func(e *api.Node) bool ***REMOVED***
				if len(request.Filters.NamePrefixes) == 0 ***REMOVED***
					return true
				***REMOVED***
				if e.Description == nil ***REMOVED***
					return false
				***REMOVED***
				return filterContainsPrefix(e.Description.Hostname, request.Filters.NamePrefixes)
			***REMOVED***,
			func(e *api.Node) bool ***REMOVED***
				return filterContainsPrefix(e.ID, request.Filters.IDPrefixes)
			***REMOVED***,
			func(e *api.Node) bool ***REMOVED***
				if len(request.Filters.Labels) == 0 ***REMOVED***
					return true
				***REMOVED***
				if e.Description == nil ***REMOVED***
					return false
				***REMOVED***
				return filterMatchLabels(e.Description.Engine.Labels, request.Filters.Labels)
			***REMOVED***,
			func(e *api.Node) bool ***REMOVED***
				if len(request.Filters.Roles) == 0 ***REMOVED***
					return true
				***REMOVED***
				for _, c := range request.Filters.Roles ***REMOVED***
					if c == e.Role ***REMOVED***
						return true
					***REMOVED***
				***REMOVED***
				return false
			***REMOVED***,
			func(e *api.Node) bool ***REMOVED***
				if len(request.Filters.Memberships) == 0 ***REMOVED***
					return true
				***REMOVED***
				for _, c := range request.Filters.Memberships ***REMOVED***
					if c == e.Spec.Membership ***REMOVED***
						return true
					***REMOVED***
				***REMOVED***
				return false
			***REMOVED***,
		)
	***REMOVED***

	// Add in manager information on nodes that are managers
	if s.raft != nil ***REMOVED***
		memberlist := s.raft.GetMemberlist()

		for _, node := range nodes ***REMOVED***
			for _, member := range memberlist ***REMOVED***
				if member.NodeID == node.ID ***REMOVED***
					node.ManagerStatus = &api.ManagerStatus***REMOVED***
						RaftID:       member.RaftID,
						Addr:         member.Addr,
						Leader:       member.Status.Leader,
						Reachability: member.Status.Reachability,
					***REMOVED***
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return &api.ListNodesResponse***REMOVED***
		Nodes: nodes,
	***REMOVED***, nil
***REMOVED***

// UpdateNode updates a Node referenced by NodeID with the given NodeSpec.
// - Returns `NotFound` if the Node is not found.
// - Returns `InvalidArgument` if the NodeSpec is malformed.
// - Returns an error if the update fails.
func (s *Server) UpdateNode(ctx context.Context, request *api.UpdateNodeRequest) (*api.UpdateNodeResponse, error) ***REMOVED***
	if request.NodeID == "" || request.NodeVersion == nil ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***
	if err := validateNodeSpec(request.Spec); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var (
		node   *api.Node
		member *membership.Member
	)

	err := s.store.Update(func(tx store.Tx) error ***REMOVED***
		node = store.GetNode(tx, request.NodeID)
		if node == nil ***REMOVED***
			return status.Errorf(codes.NotFound, "node %s not found", request.NodeID)
		***REMOVED***

		// Demotion sanity checks.
		if node.Spec.DesiredRole == api.NodeRoleManager && request.Spec.DesiredRole == api.NodeRoleWorker ***REMOVED***
			// Check for manager entries in Store.
			managers, err := store.FindNodes(tx, store.ByRole(api.NodeRoleManager))
			if err != nil ***REMOVED***
				return status.Errorf(codes.Internal, "internal store error: %v", err)
			***REMOVED***
			if len(managers) == 1 && managers[0].ID == node.ID ***REMOVED***
				return status.Errorf(codes.FailedPrecondition, "attempting to demote the last manager of the swarm")
			***REMOVED***

			// Check for node in memberlist
			if member = s.raft.GetMemberByNodeID(request.NodeID); member == nil ***REMOVED***
				return status.Errorf(codes.NotFound, "can't find manager in raft memberlist")
			***REMOVED***

			// Quorum safeguard
			if !s.raft.CanRemoveMember(member.RaftID) ***REMOVED***
				return status.Errorf(codes.FailedPrecondition, "can't remove member from the raft: this would result in a loss of quorum")
			***REMOVED***
		***REMOVED***

		node.Meta.Version = *request.NodeVersion
		node.Spec = *request.Spec.Copy()
		return store.UpdateNode(tx, node)
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &api.UpdateNodeResponse***REMOVED***
		Node: node,
	***REMOVED***, nil
***REMOVED***

func removeNodeAttachments(tx store.Tx, nodeID string) error ***REMOVED***
	// orphan the node's attached containers. if we don't do this, the
	// network these attachments are connected to will never be removeable
	tasks, err := store.FindTasks(tx, store.ByNodeID(nodeID))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, task := range tasks ***REMOVED***
		// if the task is an attachment, then we just delete it. the allocator
		// will do the heavy lifting. basically, GetAttachment will return the
		// attachment if that's the kind of runtime, or nil if it's not.
		if task.Spec.GetAttachment() != nil ***REMOVED***
			// don't delete the task. instead, update it to `ORPHANED` so that
			// the taskreaper will clean it up.
			task.Status.State = api.TaskStateOrphaned
			if err := store.UpdateTask(tx, task); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// RemoveNode removes a Node referenced by NodeID with the given NodeSpec.
// - Returns NotFound if the Node is not found.
// - Returns FailedPrecondition if the Node has manager role (and is part of the memberlist) or is not shut down.
// - Returns InvalidArgument if NodeID or NodeVersion is not valid.
// - Returns an error if the delete fails.
func (s *Server) RemoveNode(ctx context.Context, request *api.RemoveNodeRequest) (*api.RemoveNodeResponse, error) ***REMOVED***
	if request.NodeID == "" ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***

	err := s.store.Update(func(tx store.Tx) error ***REMOVED***
		node := store.GetNode(tx, request.NodeID)
		if node == nil ***REMOVED***
			return status.Errorf(codes.NotFound, "node %s not found", request.NodeID)
		***REMOVED***
		if node.Spec.DesiredRole == api.NodeRoleManager ***REMOVED***
			if s.raft == nil ***REMOVED***
				return status.Errorf(codes.FailedPrecondition, "node %s is a manager but cannot access node information from the raft memberlist", request.NodeID)
			***REMOVED***
			if member := s.raft.GetMemberByNodeID(request.NodeID); member != nil ***REMOVED***
				return status.Errorf(codes.FailedPrecondition, "node %s is a cluster manager and is a member of the raft cluster. It must be demoted to worker before removal", request.NodeID)
			***REMOVED***
		***REMOVED***
		if !request.Force && node.Status.State == api.NodeStatus_READY ***REMOVED***
			return status.Errorf(codes.FailedPrecondition, "node %s is not down and can't be removed", request.NodeID)
		***REMOVED***

		// lookup the cluster
		clusters, err := store.FindClusters(tx, store.ByName(store.DefaultClusterName))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if len(clusters) != 1 ***REMOVED***
			return status.Errorf(codes.Internal, "could not fetch cluster object")
		***REMOVED***
		cluster := clusters[0]

		blacklistedCert := &api.BlacklistedCertificate***REMOVED******REMOVED***

		// Set an expiry time for this RemovedNode if a certificate
		// exists and can be parsed.
		if len(node.Certificate.Certificate) != 0 ***REMOVED***
			certBlock, _ := pem.Decode(node.Certificate.Certificate)
			if certBlock != nil ***REMOVED***
				X509Cert, err := x509.ParseCertificate(certBlock.Bytes)
				if err == nil && !X509Cert.NotAfter.IsZero() ***REMOVED***
					expiry, err := gogotypes.TimestampProto(X509Cert.NotAfter)
					if err == nil ***REMOVED***
						blacklistedCert.Expiry = expiry
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if cluster.BlacklistedCertificates == nil ***REMOVED***
			cluster.BlacklistedCertificates = make(map[string]*api.BlacklistedCertificate)
		***REMOVED***
		cluster.BlacklistedCertificates[node.ID] = blacklistedCert

		expireBlacklistedCerts(cluster)

		if err := store.UpdateCluster(tx, cluster); err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := removeNodeAttachments(tx, request.NodeID); err != nil ***REMOVED***
			return err
		***REMOVED***

		return store.DeleteNode(tx, request.NodeID)
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &api.RemoveNodeResponse***REMOVED******REMOVED***, nil
***REMOVED***
