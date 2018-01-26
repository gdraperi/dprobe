package namespaces

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

const (
	// GRPCHeader defines the header name for specifying a containerd namespace.
	GRPCHeader = "containerd-namespace"
)

// NOTE(stevvooe): We can stub this file out if we don't want a grpc dependency here.

func withGRPCNamespaceHeader(ctx context.Context, namespace string) context.Context ***REMOVED***
	// also store on the grpc headers so it gets picked up by any clients that
	// are using this.
	nsheader := metadata.Pairs(GRPCHeader, namespace)
	md, ok := metadata.FromOutgoingContext(ctx) // merge with outgoing context.
	if !ok ***REMOVED***
		md = nsheader
	***REMOVED*** else ***REMOVED***
		// order ensures the latest is first in this list.
		md = metadata.Join(nsheader, md)
	***REMOVED***

	return metadata.NewOutgoingContext(ctx, md)
***REMOVED***

func fromGRPCHeader(ctx context.Context) (string, bool) ***REMOVED***
	// try to extract for use in grpc servers.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok ***REMOVED***
		// TODO(stevvooe): Check outgoing context?
		return "", false
	***REMOVED***

	values := md[GRPCHeader]
	if len(values) == 0 ***REMOVED***
		return "", false
	***REMOVED***

	return values[0], true
***REMOVED***
