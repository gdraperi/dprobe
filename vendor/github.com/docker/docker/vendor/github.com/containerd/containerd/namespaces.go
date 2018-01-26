package containerd

import (
	"context"
	"strings"

	api "github.com/containerd/containerd/api/services/namespaces/v1"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/namespaces"
	"github.com/gogo/protobuf/types"
)

// NewNamespaceStoreFromClient returns a new namespace store
func NewNamespaceStoreFromClient(client api.NamespacesClient) namespaces.Store ***REMOVED***
	return &remoteNamespaces***REMOVED***client: client***REMOVED***
***REMOVED***

type remoteNamespaces struct ***REMOVED***
	client api.NamespacesClient
***REMOVED***

func (r *remoteNamespaces) Create(ctx context.Context, namespace string, labels map[string]string) error ***REMOVED***
	var req api.CreateNamespaceRequest

	req.Namespace = api.Namespace***REMOVED***
		Name:   namespace,
		Labels: labels,
	***REMOVED***

	_, err := r.client.Create(ctx, &req)
	if err != nil ***REMOVED***
		return errdefs.FromGRPC(err)
	***REMOVED***

	return nil
***REMOVED***

func (r *remoteNamespaces) Labels(ctx context.Context, namespace string) (map[string]string, error) ***REMOVED***
	var req api.GetNamespaceRequest
	req.Name = namespace

	resp, err := r.client.Get(ctx, &req)
	if err != nil ***REMOVED***
		return nil, errdefs.FromGRPC(err)
	***REMOVED***

	return resp.Namespace.Labels, nil
***REMOVED***

func (r *remoteNamespaces) SetLabel(ctx context.Context, namespace, key, value string) error ***REMOVED***
	var req api.UpdateNamespaceRequest

	req.Namespace = api.Namespace***REMOVED***
		Name:   namespace,
		Labels: map[string]string***REMOVED***key: value***REMOVED***,
	***REMOVED***

	req.UpdateMask = &types.FieldMask***REMOVED***
		Paths: []string***REMOVED***strings.Join([]string***REMOVED***"labels", key***REMOVED***, ".")***REMOVED***,
	***REMOVED***

	_, err := r.client.Update(ctx, &req)
	if err != nil ***REMOVED***
		return errdefs.FromGRPC(err)
	***REMOVED***

	return nil
***REMOVED***

func (r *remoteNamespaces) List(ctx context.Context) ([]string, error) ***REMOVED***
	var req api.ListNamespacesRequest

	resp, err := r.client.List(ctx, &req)
	if err != nil ***REMOVED***
		return nil, errdefs.FromGRPC(err)
	***REMOVED***

	var namespaces []string

	for _, ns := range resp.Namespaces ***REMOVED***
		namespaces = append(namespaces, ns.Name)
	***REMOVED***

	return namespaces, nil
***REMOVED***

func (r *remoteNamespaces) Delete(ctx context.Context, namespace string) error ***REMOVED***
	var req api.DeleteNamespaceRequest

	req.Name = namespace
	_, err := r.client.Delete(ctx, &req)
	if err != nil ***REMOVED***
		return errdefs.FromGRPC(err)
	***REMOVED***

	return nil
***REMOVED***
