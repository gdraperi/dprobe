package leases

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

const (
	// GRPCHeader defines the header name for specifying a containerd lease.
	GRPCHeader = "containerd-lease"
)

func withGRPCLeaseHeader(ctx context.Context, lid string) context.Context ***REMOVED***
	// also store on the grpc headers so it gets picked up by any clients
	// that are using this.
	txheader := metadata.Pairs(GRPCHeader, lid)
	md, ok := metadata.FromOutgoingContext(ctx) // merge with outgoing context.
	if !ok ***REMOVED***
		md = txheader
	***REMOVED*** else ***REMOVED***
		// order ensures the latest is first in this list.
		md = metadata.Join(txheader, md)
	***REMOVED***

	return metadata.NewOutgoingContext(ctx, md)
***REMOVED***

func fromGRPCHeader(ctx context.Context) (string, bool) ***REMOVED***
	// try to extract for use in grpc servers.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok ***REMOVED***
		return "", false
	***REMOVED***

	values := md[GRPCHeader]
	if len(values) == 0 ***REMOVED***
		return "", false
	***REMOVED***

	return values[0], true
***REMOVED***
