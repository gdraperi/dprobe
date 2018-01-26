package containerd

import (
	"context"
	"io"

	snapshotsapi "github.com/containerd/containerd/api/services/snapshots/v1"
	"github.com/containerd/containerd/api/types"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/snapshots"
	protobuftypes "github.com/gogo/protobuf/types"
)

// NewSnapshotterFromClient returns a new Snapshotter which communicates
// over a GRPC connection.
func NewSnapshotterFromClient(client snapshotsapi.SnapshotsClient, snapshotterName string) snapshots.Snapshotter ***REMOVED***
	return &remoteSnapshotter***REMOVED***
		client:          client,
		snapshotterName: snapshotterName,
	***REMOVED***
***REMOVED***

type remoteSnapshotter struct ***REMOVED***
	client          snapshotsapi.SnapshotsClient
	snapshotterName string
***REMOVED***

func (r *remoteSnapshotter) Stat(ctx context.Context, key string) (snapshots.Info, error) ***REMOVED***
	resp, err := r.client.Stat(ctx,
		&snapshotsapi.StatSnapshotRequest***REMOVED***
			Snapshotter: r.snapshotterName,
			Key:         key,
		***REMOVED***)
	if err != nil ***REMOVED***
		return snapshots.Info***REMOVED******REMOVED***, errdefs.FromGRPC(err)
	***REMOVED***
	return toInfo(resp.Info), nil
***REMOVED***

func (r *remoteSnapshotter) Update(ctx context.Context, info snapshots.Info, fieldpaths ...string) (snapshots.Info, error) ***REMOVED***
	resp, err := r.client.Update(ctx,
		&snapshotsapi.UpdateSnapshotRequest***REMOVED***
			Snapshotter: r.snapshotterName,
			Info:        fromInfo(info),
			UpdateMask: &protobuftypes.FieldMask***REMOVED***
				Paths: fieldpaths,
			***REMOVED***,
		***REMOVED***)
	if err != nil ***REMOVED***
		return snapshots.Info***REMOVED******REMOVED***, errdefs.FromGRPC(err)
	***REMOVED***
	return toInfo(resp.Info), nil
***REMOVED***

func (r *remoteSnapshotter) Usage(ctx context.Context, key string) (snapshots.Usage, error) ***REMOVED***
	resp, err := r.client.Usage(ctx, &snapshotsapi.UsageRequest***REMOVED***
		Snapshotter: r.snapshotterName,
		Key:         key,
	***REMOVED***)
	if err != nil ***REMOVED***
		return snapshots.Usage***REMOVED******REMOVED***, errdefs.FromGRPC(err)
	***REMOVED***
	return toUsage(resp), nil
***REMOVED***

func (r *remoteSnapshotter) Mounts(ctx context.Context, key string) ([]mount.Mount, error) ***REMOVED***
	resp, err := r.client.Mounts(ctx, &snapshotsapi.MountsRequest***REMOVED***
		Snapshotter: r.snapshotterName,
		Key:         key,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, errdefs.FromGRPC(err)
	***REMOVED***
	return toMounts(resp.Mounts), nil
***REMOVED***

func (r *remoteSnapshotter) Prepare(ctx context.Context, key, parent string, opts ...snapshots.Opt) ([]mount.Mount, error) ***REMOVED***
	var local snapshots.Info
	for _, opt := range opts ***REMOVED***
		if err := opt(&local); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	resp, err := r.client.Prepare(ctx, &snapshotsapi.PrepareSnapshotRequest***REMOVED***
		Snapshotter: r.snapshotterName,
		Key:         key,
		Parent:      parent,
		Labels:      local.Labels,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, errdefs.FromGRPC(err)
	***REMOVED***
	return toMounts(resp.Mounts), nil
***REMOVED***

func (r *remoteSnapshotter) View(ctx context.Context, key, parent string, opts ...snapshots.Opt) ([]mount.Mount, error) ***REMOVED***
	var local snapshots.Info
	for _, opt := range opts ***REMOVED***
		if err := opt(&local); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	resp, err := r.client.View(ctx, &snapshotsapi.ViewSnapshotRequest***REMOVED***
		Snapshotter: r.snapshotterName,
		Key:         key,
		Parent:      parent,
		Labels:      local.Labels,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, errdefs.FromGRPC(err)
	***REMOVED***
	return toMounts(resp.Mounts), nil
***REMOVED***

func (r *remoteSnapshotter) Commit(ctx context.Context, name, key string, opts ...snapshots.Opt) error ***REMOVED***
	var local snapshots.Info
	for _, opt := range opts ***REMOVED***
		if err := opt(&local); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	_, err := r.client.Commit(ctx, &snapshotsapi.CommitSnapshotRequest***REMOVED***
		Snapshotter: r.snapshotterName,
		Name:        name,
		Key:         key,
		Labels:      local.Labels,
	***REMOVED***)
	return errdefs.FromGRPC(err)
***REMOVED***

func (r *remoteSnapshotter) Remove(ctx context.Context, key string) error ***REMOVED***
	_, err := r.client.Remove(ctx, &snapshotsapi.RemoveSnapshotRequest***REMOVED***
		Snapshotter: r.snapshotterName,
		Key:         key,
	***REMOVED***)
	return errdefs.FromGRPC(err)
***REMOVED***

func (r *remoteSnapshotter) Walk(ctx context.Context, fn func(context.Context, snapshots.Info) error) error ***REMOVED***
	sc, err := r.client.List(ctx, &snapshotsapi.ListSnapshotsRequest***REMOVED***
		Snapshotter: r.snapshotterName,
	***REMOVED***)
	if err != nil ***REMOVED***
		return errdefs.FromGRPC(err)
	***REMOVED***
	for ***REMOVED***
		resp, err := sc.Recv()
		if err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				return nil
			***REMOVED***
			return errdefs.FromGRPC(err)
		***REMOVED***
		if resp == nil ***REMOVED***
			return nil
		***REMOVED***
		for _, info := range resp.Info ***REMOVED***
			if err := fn(ctx, toInfo(info)); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *remoteSnapshotter) Close() error ***REMOVED***
	return nil
***REMOVED***

func toKind(kind snapshotsapi.Kind) snapshots.Kind ***REMOVED***
	if kind == snapshotsapi.KindActive ***REMOVED***
		return snapshots.KindActive
	***REMOVED***
	if kind == snapshotsapi.KindView ***REMOVED***
		return snapshots.KindView
	***REMOVED***
	return snapshots.KindCommitted
***REMOVED***

func toInfo(info snapshotsapi.Info) snapshots.Info ***REMOVED***
	return snapshots.Info***REMOVED***
		Name:    info.Name,
		Parent:  info.Parent,
		Kind:    toKind(info.Kind),
		Created: info.CreatedAt,
		Updated: info.UpdatedAt,
		Labels:  info.Labels,
	***REMOVED***
***REMOVED***

func toUsage(resp *snapshotsapi.UsageResponse) snapshots.Usage ***REMOVED***
	return snapshots.Usage***REMOVED***
		Inodes: resp.Inodes,
		Size:   resp.Size_,
	***REMOVED***
***REMOVED***

func toMounts(mm []*types.Mount) []mount.Mount ***REMOVED***
	mounts := make([]mount.Mount, len(mm))
	for i, m := range mm ***REMOVED***
		mounts[i] = mount.Mount***REMOVED***
			Type:    m.Type,
			Source:  m.Source,
			Options: m.Options,
		***REMOVED***
	***REMOVED***
	return mounts
***REMOVED***

func fromKind(kind snapshots.Kind) snapshotsapi.Kind ***REMOVED***
	if kind == snapshots.KindActive ***REMOVED***
		return snapshotsapi.KindActive
	***REMOVED***
	if kind == snapshots.KindView ***REMOVED***
		return snapshotsapi.KindView
	***REMOVED***
	return snapshotsapi.KindCommitted
***REMOVED***

func fromInfo(info snapshots.Info) snapshotsapi.Info ***REMOVED***
	return snapshotsapi.Info***REMOVED***
		Name:      info.Name,
		Parent:    info.Parent,
		Kind:      fromKind(info.Kind),
		CreatedAt: info.Created,
		UpdatedAt: info.Updated,
		Labels:    info.Labels,
	***REMOVED***
***REMOVED***
