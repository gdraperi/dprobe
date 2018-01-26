package containerd

import (
	"context"
	"time"

	leasesapi "github.com/containerd/containerd/api/services/leases/v1"
	"github.com/containerd/containerd/leases"
)

// Lease is used to hold a reference to active resources which have not been
// referenced by a root resource. This is useful for preventing garbage
// collection of resources while they are actively being updated.
type Lease struct ***REMOVED***
	id        string
	createdAt time.Time

	client *Client
***REMOVED***

// CreateLease creates a new lease
func (c *Client) CreateLease(ctx context.Context) (Lease, error) ***REMOVED***
	lapi := leasesapi.NewLeasesClient(c.conn)
	resp, err := lapi.Create(ctx, &leasesapi.CreateRequest***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		return Lease***REMOVED******REMOVED***, err
	***REMOVED***

	return Lease***REMOVED***
		id:     resp.Lease.ID,
		client: c,
	***REMOVED***, nil
***REMOVED***

// ListLeases lists active leases
func (c *Client) ListLeases(ctx context.Context) ([]Lease, error) ***REMOVED***
	lapi := leasesapi.NewLeasesClient(c.conn)
	resp, err := lapi.List(ctx, &leasesapi.ListRequest***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	leases := make([]Lease, len(resp.Leases))
	for i := range resp.Leases ***REMOVED***
		leases[i] = Lease***REMOVED***
			id:        resp.Leases[i].ID,
			createdAt: resp.Leases[i].CreatedAt,
			client:    c,
		***REMOVED***
	***REMOVED***

	return leases, nil
***REMOVED***

// WithLease attaches a lease on the context
func (c *Client) WithLease(ctx context.Context) (context.Context, func() error, error) ***REMOVED***
	_, ok := leases.Lease(ctx)
	if ok ***REMOVED***
		return ctx, func() error ***REMOVED***
			return nil
		***REMOVED***, nil
	***REMOVED***

	l, err := c.CreateLease(ctx)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	ctx = leases.WithLease(ctx, l.ID())
	return ctx, func() error ***REMOVED***
		return l.Delete(ctx)
	***REMOVED***, nil
***REMOVED***

// ID returns the lease ID
func (l Lease) ID() string ***REMOVED***
	return l.id
***REMOVED***

// CreatedAt returns the time at which the lease was created
func (l Lease) CreatedAt() time.Time ***REMOVED***
	return l.createdAt
***REMOVED***

// Delete deletes the lease, removing the reference to all resources created
// during the lease.
func (l Lease) Delete(ctx context.Context) error ***REMOVED***
	lapi := leasesapi.NewLeasesClient(l.client.conn)
	_, err := lapi.Delete(ctx, &leasesapi.DeleteRequest***REMOVED***
		ID: l.id,
	***REMOVED***)
	return err
***REMOVED***
