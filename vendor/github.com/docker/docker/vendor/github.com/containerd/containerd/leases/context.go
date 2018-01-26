package leases

import "context"

type leaseKey struct***REMOVED******REMOVED***

// WithLease sets a given lease on the context
func WithLease(ctx context.Context, lid string) context.Context ***REMOVED***
	ctx = context.WithValue(ctx, leaseKey***REMOVED******REMOVED***, lid)

	// also store on the grpc headers so it gets picked up by any clients that
	// are using this.
	return withGRPCLeaseHeader(ctx, lid)
***REMOVED***

// Lease returns the lease from the context.
func Lease(ctx context.Context) (string, bool) ***REMOVED***
	lid, ok := ctx.Value(leaseKey***REMOVED******REMOVED***).(string)
	if !ok ***REMOVED***
		return fromGRPCHeader(ctx)
	***REMOVED***

	return lid, ok
***REMOVED***
