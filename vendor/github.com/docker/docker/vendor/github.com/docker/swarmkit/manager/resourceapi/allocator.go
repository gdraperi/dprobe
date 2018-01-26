package resourceapi

import (
	"errors"
	"time"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/ca"
	"github.com/docker/swarmkit/identity"
	"github.com/docker/swarmkit/manager/state/store"
	"github.com/docker/swarmkit/protobuf/ptypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errInvalidArgument = errors.New("invalid argument")
)

// ResourceAllocator handles resource allocation of cluster entities.
type ResourceAllocator struct ***REMOVED***
	store *store.MemoryStore
***REMOVED***

// New returns an instance of the allocator
func New(store *store.MemoryStore) *ResourceAllocator ***REMOVED***
	return &ResourceAllocator***REMOVED***store: store***REMOVED***
***REMOVED***

// AttachNetwork allows the node to request the resources
// allocation needed for a network attachment on the specific node.
// - Returns `InvalidArgument` if the Spec is malformed.
// - Returns `NotFound` if the Network is not found.
// - Returns `PermissionDenied` if the Network is not manually attachable.
// - Returns an error if the creation fails.
func (ra *ResourceAllocator) AttachNetwork(ctx context.Context, request *api.AttachNetworkRequest) (*api.AttachNetworkResponse, error) ***REMOVED***
	nodeInfo, err := ca.RemoteNode(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var network *api.Network
	ra.store.View(func(tx store.ReadTx) ***REMOVED***
		network = store.GetNetwork(tx, request.Config.Target)
		if network == nil ***REMOVED***
			if networks, err := store.FindNetworks(tx, store.ByName(request.Config.Target)); err == nil && len(networks) == 1 ***REMOVED***
				network = networks[0]
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	if network == nil ***REMOVED***
		return nil, status.Errorf(codes.NotFound, "network %s not found", request.Config.Target)
	***REMOVED***

	if !network.Spec.Attachable ***REMOVED***
		return nil, status.Errorf(codes.PermissionDenied, "network %s not manually attachable", request.Config.Target)
	***REMOVED***

	t := &api.Task***REMOVED***
		ID:     identity.NewID(),
		NodeID: nodeInfo.NodeID,
		Spec: api.TaskSpec***REMOVED***
			Runtime: &api.TaskSpec_Attachment***REMOVED***
				Attachment: &api.NetworkAttachmentSpec***REMOVED***
					ContainerID: request.ContainerID,
				***REMOVED***,
			***REMOVED***,
			Networks: []*api.NetworkAttachmentConfig***REMOVED***
				***REMOVED***
					Target:    network.ID,
					Addresses: request.Config.Addresses,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Status: api.TaskStatus***REMOVED***
			State:     api.TaskStateNew,
			Timestamp: ptypes.MustTimestampProto(time.Now()),
			Message:   "created",
		***REMOVED***,
		DesiredState: api.TaskStateRunning,
		// TODO: Add Network attachment.
	***REMOVED***

	if err := ra.store.Update(func(tx store.Tx) error ***REMOVED***
		return store.CreateTask(tx, t)
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &api.AttachNetworkResponse***REMOVED***AttachmentID: t.ID***REMOVED***, nil
***REMOVED***

// DetachNetwork allows the node to request the release of
// the resources associated to the network attachment.
// - Returns `InvalidArgument` if attachment ID is not provided.
// - Returns `NotFound` if the attachment is not found.
// - Returns an error if the deletion fails.
func (ra *ResourceAllocator) DetachNetwork(ctx context.Context, request *api.DetachNetworkRequest) (*api.DetachNetworkResponse, error) ***REMOVED***
	if request.AttachmentID == "" ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***

	nodeInfo, err := ca.RemoteNode(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := ra.store.Update(func(tx store.Tx) error ***REMOVED***
		t := store.GetTask(tx, request.AttachmentID)
		if t == nil ***REMOVED***
			return status.Errorf(codes.NotFound, "attachment %s not found", request.AttachmentID)
		***REMOVED***
		if t.NodeID != nodeInfo.NodeID ***REMOVED***
			return status.Errorf(codes.PermissionDenied, "attachment %s doesn't belong to this node", request.AttachmentID)
		***REMOVED***

		return store.DeleteTask(tx, request.AttachmentID)
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &api.DetachNetworkResponse***REMOVED******REMOVED***, nil
***REMOVED***
