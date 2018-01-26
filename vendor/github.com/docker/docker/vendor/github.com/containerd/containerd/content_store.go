package containerd

import (
	"context"
	"io"

	contentapi "github.com/containerd/containerd/api/services/content/v1"
	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/errdefs"
	protobuftypes "github.com/gogo/protobuf/types"
	digest "github.com/opencontainers/go-digest"
)

type remoteContent struct ***REMOVED***
	client contentapi.ContentClient
***REMOVED***

// NewContentStoreFromClient returns a new content store
func NewContentStoreFromClient(client contentapi.ContentClient) content.Store ***REMOVED***
	return &remoteContent***REMOVED***
		client: client,
	***REMOVED***
***REMOVED***

func (rs *remoteContent) Info(ctx context.Context, dgst digest.Digest) (content.Info, error) ***REMOVED***
	resp, err := rs.client.Info(ctx, &contentapi.InfoRequest***REMOVED***
		Digest: dgst,
	***REMOVED***)
	if err != nil ***REMOVED***
		return content.Info***REMOVED******REMOVED***, errdefs.FromGRPC(err)
	***REMOVED***

	return infoFromGRPC(resp.Info), nil
***REMOVED***

func (rs *remoteContent) Walk(ctx context.Context, fn content.WalkFunc, filters ...string) error ***REMOVED***
	session, err := rs.client.List(ctx, &contentapi.ListContentRequest***REMOVED***
		Filters: filters,
	***REMOVED***)
	if err != nil ***REMOVED***
		return errdefs.FromGRPC(err)
	***REMOVED***

	for ***REMOVED***
		msg, err := session.Recv()
		if err != nil ***REMOVED***
			if err != io.EOF ***REMOVED***
				return errdefs.FromGRPC(err)
			***REMOVED***

			break
		***REMOVED***

		for _, info := range msg.Info ***REMOVED***
			if err := fn(infoFromGRPC(info)); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (rs *remoteContent) Delete(ctx context.Context, dgst digest.Digest) error ***REMOVED***
	if _, err := rs.client.Delete(ctx, &contentapi.DeleteContentRequest***REMOVED***
		Digest: dgst,
	***REMOVED***); err != nil ***REMOVED***
		return errdefs.FromGRPC(err)
	***REMOVED***

	return nil
***REMOVED***

func (rs *remoteContent) ReaderAt(ctx context.Context, dgst digest.Digest) (content.ReaderAt, error) ***REMOVED***
	i, err := rs.Info(ctx, dgst)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &remoteReaderAt***REMOVED***
		ctx:    ctx,
		digest: dgst,
		size:   i.Size,
		client: rs.client,
	***REMOVED***, nil
***REMOVED***

func (rs *remoteContent) Status(ctx context.Context, ref string) (content.Status, error) ***REMOVED***
	resp, err := rs.client.Status(ctx, &contentapi.StatusRequest***REMOVED***
		Ref: ref,
	***REMOVED***)
	if err != nil ***REMOVED***
		return content.Status***REMOVED******REMOVED***, errdefs.FromGRPC(err)
	***REMOVED***

	status := resp.Status
	return content.Status***REMOVED***
		Ref:       status.Ref,
		StartedAt: status.StartedAt,
		UpdatedAt: status.UpdatedAt,
		Offset:    status.Offset,
		Total:     status.Total,
		Expected:  status.Expected,
	***REMOVED***, nil
***REMOVED***

func (rs *remoteContent) Update(ctx context.Context, info content.Info, fieldpaths ...string) (content.Info, error) ***REMOVED***
	resp, err := rs.client.Update(ctx, &contentapi.UpdateRequest***REMOVED***
		Info: infoToGRPC(info),
		UpdateMask: &protobuftypes.FieldMask***REMOVED***
			Paths: fieldpaths,
		***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		return content.Info***REMOVED******REMOVED***, errdefs.FromGRPC(err)
	***REMOVED***
	return infoFromGRPC(resp.Info), nil
***REMOVED***

func (rs *remoteContent) ListStatuses(ctx context.Context, filters ...string) ([]content.Status, error) ***REMOVED***
	resp, err := rs.client.ListStatuses(ctx, &contentapi.ListStatusesRequest***REMOVED***
		Filters: filters,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, errdefs.FromGRPC(err)
	***REMOVED***

	var statuses []content.Status
	for _, status := range resp.Statuses ***REMOVED***
		statuses = append(statuses, content.Status***REMOVED***
			Ref:       status.Ref,
			StartedAt: status.StartedAt,
			UpdatedAt: status.UpdatedAt,
			Offset:    status.Offset,
			Total:     status.Total,
			Expected:  status.Expected,
		***REMOVED***)
	***REMOVED***

	return statuses, nil
***REMOVED***

func (rs *remoteContent) Writer(ctx context.Context, ref string, size int64, expected digest.Digest) (content.Writer, error) ***REMOVED***
	wrclient, offset, err := rs.negotiate(ctx, ref, size, expected)
	if err != nil ***REMOVED***
		return nil, errdefs.FromGRPC(err)
	***REMOVED***

	return &remoteWriter***REMOVED***
		ref:    ref,
		client: wrclient,
		offset: offset,
	***REMOVED***, nil
***REMOVED***

// Abort implements asynchronous abort. It starts a new write session on the ref l
func (rs *remoteContent) Abort(ctx context.Context, ref string) error ***REMOVED***
	if _, err := rs.client.Abort(ctx, &contentapi.AbortRequest***REMOVED***
		Ref: ref,
	***REMOVED***); err != nil ***REMOVED***
		return errdefs.FromGRPC(err)
	***REMOVED***

	return nil
***REMOVED***

func (rs *remoteContent) negotiate(ctx context.Context, ref string, size int64, expected digest.Digest) (contentapi.Content_WriteClient, int64, error) ***REMOVED***
	wrclient, err := rs.client.Write(ctx)
	if err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***

	if err := wrclient.Send(&contentapi.WriteContentRequest***REMOVED***
		Action:   contentapi.WriteActionStat,
		Ref:      ref,
		Total:    size,
		Expected: expected,
	***REMOVED***); err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***

	resp, err := wrclient.Recv()
	if err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***

	return wrclient, resp.Offset, nil
***REMOVED***

func infoToGRPC(info content.Info) contentapi.Info ***REMOVED***
	return contentapi.Info***REMOVED***
		Digest:    info.Digest,
		Size_:     info.Size,
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,
		Labels:    info.Labels,
	***REMOVED***
***REMOVED***

func infoFromGRPC(info contentapi.Info) content.Info ***REMOVED***
	return content.Info***REMOVED***
		Digest:    info.Digest,
		Size:      info.Size_,
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,
		Labels:    info.Labels,
	***REMOVED***
***REMOVED***
