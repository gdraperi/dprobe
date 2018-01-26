package containerd

import (
	diffapi "github.com/containerd/containerd/api/services/diff/v1"
	"github.com/containerd/containerd/api/types"
	"github.com/containerd/containerd/diff"
	"github.com/containerd/containerd/mount"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"golang.org/x/net/context"
)

// NewDiffServiceFromClient returns a new diff service which communicates
// over a GRPC connection.
func NewDiffServiceFromClient(client diffapi.DiffClient) diff.Differ ***REMOVED***
	return &diffRemote***REMOVED***
		client: client,
	***REMOVED***
***REMOVED***

type diffRemote struct ***REMOVED***
	client diffapi.DiffClient
***REMOVED***

func (r *diffRemote) Apply(ctx context.Context, diff ocispec.Descriptor, mounts []mount.Mount) (ocispec.Descriptor, error) ***REMOVED***
	req := &diffapi.ApplyRequest***REMOVED***
		Diff:   fromDescriptor(diff),
		Mounts: fromMounts(mounts),
	***REMOVED***
	resp, err := r.client.Apply(ctx, req)
	if err != nil ***REMOVED***
		return ocispec.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***
	return toDescriptor(resp.Applied), nil
***REMOVED***

func (r *diffRemote) DiffMounts(ctx context.Context, a, b []mount.Mount, opts ...diff.Opt) (ocispec.Descriptor, error) ***REMOVED***
	var config diff.Config
	for _, opt := range opts ***REMOVED***
		if err := opt(&config); err != nil ***REMOVED***
			return ocispec.Descriptor***REMOVED******REMOVED***, err
		***REMOVED***
	***REMOVED***
	req := &diffapi.DiffRequest***REMOVED***
		Left:      fromMounts(a),
		Right:     fromMounts(b),
		MediaType: config.MediaType,
		Ref:       config.Reference,
		Labels:    config.Labels,
	***REMOVED***
	resp, err := r.client.Diff(ctx, req)
	if err != nil ***REMOVED***
		return ocispec.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***
	return toDescriptor(resp.Diff), nil
***REMOVED***

func toDescriptor(d *types.Descriptor) ocispec.Descriptor ***REMOVED***
	return ocispec.Descriptor***REMOVED***
		MediaType: d.MediaType,
		Digest:    d.Digest,
		Size:      d.Size_,
	***REMOVED***
***REMOVED***

func fromDescriptor(d ocispec.Descriptor) *types.Descriptor ***REMOVED***
	return &types.Descriptor***REMOVED***
		MediaType: d.MediaType,
		Digest:    d.Digest,
		Size_:     d.Size,
	***REMOVED***
***REMOVED***

func fromMounts(mounts []mount.Mount) []*types.Mount ***REMOVED***
	apiMounts := make([]*types.Mount, len(mounts))
	for i, m := range mounts ***REMOVED***
		apiMounts[i] = &types.Mount***REMOVED***
			Type:    m.Type,
			Source:  m.Source,
			Options: m.Options,
		***REMOVED***
	***REMOVED***
	return apiMounts
***REMOVED***
