package containerd

import (
	"context"

	containersapi "github.com/containerd/containerd/api/services/containers/v1"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/errdefs"
	ptypes "github.com/gogo/protobuf/types"
)

type remoteContainers struct ***REMOVED***
	client containersapi.ContainersClient
***REMOVED***

var _ containers.Store = &remoteContainers***REMOVED******REMOVED***

// NewRemoteContainerStore returns the container Store connected with the provided client
func NewRemoteContainerStore(client containersapi.ContainersClient) containers.Store ***REMOVED***
	return &remoteContainers***REMOVED***
		client: client,
	***REMOVED***
***REMOVED***

func (r *remoteContainers) Get(ctx context.Context, id string) (containers.Container, error) ***REMOVED***
	resp, err := r.client.Get(ctx, &containersapi.GetContainerRequest***REMOVED***
		ID: id,
	***REMOVED***)
	if err != nil ***REMOVED***
		return containers.Container***REMOVED******REMOVED***, errdefs.FromGRPC(err)
	***REMOVED***

	return containerFromProto(&resp.Container), nil
***REMOVED***

func (r *remoteContainers) List(ctx context.Context, filters ...string) ([]containers.Container, error) ***REMOVED***
	resp, err := r.client.List(ctx, &containersapi.ListContainersRequest***REMOVED***
		Filters: filters,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, errdefs.FromGRPC(err)
	***REMOVED***

	return containersFromProto(resp.Containers), nil

***REMOVED***

func (r *remoteContainers) Create(ctx context.Context, container containers.Container) (containers.Container, error) ***REMOVED***
	created, err := r.client.Create(ctx, &containersapi.CreateContainerRequest***REMOVED***
		Container: containerToProto(&container),
	***REMOVED***)
	if err != nil ***REMOVED***
		return containers.Container***REMOVED******REMOVED***, errdefs.FromGRPC(err)
	***REMOVED***

	return containerFromProto(&created.Container), nil

***REMOVED***

func (r *remoteContainers) Update(ctx context.Context, container containers.Container, fieldpaths ...string) (containers.Container, error) ***REMOVED***
	var updateMask *ptypes.FieldMask
	if len(fieldpaths) > 0 ***REMOVED***
		updateMask = &ptypes.FieldMask***REMOVED***
			Paths: fieldpaths,
		***REMOVED***
	***REMOVED***

	updated, err := r.client.Update(ctx, &containersapi.UpdateContainerRequest***REMOVED***
		Container:  containerToProto(&container),
		UpdateMask: updateMask,
	***REMOVED***)
	if err != nil ***REMOVED***
		return containers.Container***REMOVED******REMOVED***, errdefs.FromGRPC(err)
	***REMOVED***

	return containerFromProto(&updated.Container), nil

***REMOVED***

func (r *remoteContainers) Delete(ctx context.Context, id string) error ***REMOVED***
	_, err := r.client.Delete(ctx, &containersapi.DeleteContainerRequest***REMOVED***
		ID: id,
	***REMOVED***)

	return errdefs.FromGRPC(err)

***REMOVED***

func containerToProto(container *containers.Container) containersapi.Container ***REMOVED***
	return containersapi.Container***REMOVED***
		ID:     container.ID,
		Labels: container.Labels,
		Image:  container.Image,
		Runtime: &containersapi.Container_Runtime***REMOVED***
			Name:    container.Runtime.Name,
			Options: container.Runtime.Options,
		***REMOVED***,
		Spec:        container.Spec,
		Snapshotter: container.Snapshotter,
		SnapshotKey: container.SnapshotKey,
		Extensions:  container.Extensions,
	***REMOVED***
***REMOVED***

func containerFromProto(containerpb *containersapi.Container) containers.Container ***REMOVED***
	var runtime containers.RuntimeInfo
	if containerpb.Runtime != nil ***REMOVED***
		runtime = containers.RuntimeInfo***REMOVED***
			Name:    containerpb.Runtime.Name,
			Options: containerpb.Runtime.Options,
		***REMOVED***
	***REMOVED***
	return containers.Container***REMOVED***
		ID:          containerpb.ID,
		Labels:      containerpb.Labels,
		Image:       containerpb.Image,
		Runtime:     runtime,
		Spec:        containerpb.Spec,
		Snapshotter: containerpb.Snapshotter,
		SnapshotKey: containerpb.SnapshotKey,
		CreatedAt:   containerpb.CreatedAt,
		UpdatedAt:   containerpb.UpdatedAt,
		Extensions:  containerpb.Extensions,
	***REMOVED***
***REMOVED***

func containersFromProto(containerspb []containersapi.Container) []containers.Container ***REMOVED***
	var containers []containers.Container

	for _, container := range containerspb ***REMOVED***
		containers = append(containers, containerFromProto(&container))
	***REMOVED***

	return containers
***REMOVED***
