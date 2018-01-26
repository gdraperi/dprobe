// Package connectionbroker is a layer on top of remotes that returns
// a gRPC connection to a manager. The connection may be a local connection
// using a local socket such as a UNIX socket.
package connectionbroker

import (
	"net"
	"sync"
	"time"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/remotes"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

// Broker is a simple connection broker. It can either return a fresh
// connection to a remote manager selected with weighted randomization, or a
// local gRPC connection to the local manager.
type Broker struct ***REMOVED***
	mu        sync.Mutex
	remotes   remotes.Remotes
	localConn *grpc.ClientConn
***REMOVED***

// New creates a new connection broker.
func New(remotes remotes.Remotes) *Broker ***REMOVED***
	return &Broker***REMOVED***
		remotes: remotes,
	***REMOVED***
***REMOVED***

// SetLocalConn changes the local gRPC connection used by the connection broker.
func (b *Broker) SetLocalConn(localConn *grpc.ClientConn) ***REMOVED***
	b.mu.Lock()
	defer b.mu.Unlock()

	b.localConn = localConn
***REMOVED***

// Select a manager from the set of available managers, and return a connection.
func (b *Broker) Select(dialOpts ...grpc.DialOption) (*Conn, error) ***REMOVED***
	b.mu.Lock()
	localConn := b.localConn
	b.mu.Unlock()

	if localConn != nil ***REMOVED***
		return &Conn***REMOVED***
			ClientConn: localConn,
			isLocal:    true,
		***REMOVED***, nil
	***REMOVED***

	return b.SelectRemote(dialOpts...)
***REMOVED***

// SelectRemote chooses a manager from the remotes, and returns a TCP
// connection.
func (b *Broker) SelectRemote(dialOpts ...grpc.DialOption) (*Conn, error) ***REMOVED***
	peer, err := b.remotes.Select()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// gRPC dialer connects to proxy first. Provide a custom dialer here avoid that.
	// TODO(anshul) Add an option to configure this.
	dialOpts = append(dialOpts,
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) ***REMOVED***
			return net.DialTimeout("tcp", addr, timeout)
		***REMOVED***))

	cc, err := grpc.Dial(peer.Addr, dialOpts...)
	if err != nil ***REMOVED***
		b.remotes.ObserveIfExists(peer, -remotes.DefaultObservationWeight)
		return nil, err
	***REMOVED***

	return &Conn***REMOVED***
		ClientConn: cc,
		remotes:    b.remotes,
		peer:       peer,
	***REMOVED***, nil
***REMOVED***

// Remotes returns the remotes interface used by the broker, so the caller
// can make observations or see weights directly.
func (b *Broker) Remotes() remotes.Remotes ***REMOVED***
	return b.remotes
***REMOVED***

// Conn is a wrapper around a gRPC client connection.
type Conn struct ***REMOVED***
	*grpc.ClientConn
	isLocal bool
	remotes remotes.Remotes
	peer    api.Peer
***REMOVED***

// Close closes the client connection if it is a remote connection. It also
// records a positive experience with the remote peer if success is true,
// otherwise it records a negative experience. If a local connection is in use,
// Close is a noop.
func (c *Conn) Close(success bool) error ***REMOVED***
	if c.isLocal ***REMOVED***
		return nil
	***REMOVED***

	if success ***REMOVED***
		c.remotes.ObserveIfExists(c.peer, remotes.DefaultObservationWeight)
	***REMOVED*** else ***REMOVED***
		c.remotes.ObserveIfExists(c.peer, -remotes.DefaultObservationWeight)
	***REMOVED***

	return c.ClientConn.Close()
***REMOVED***
