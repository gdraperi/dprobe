package ttrpc

import (
	"context"
	"net"
)

// Handshaker defines the interface for connection handshakes performed on the
// server or client when first connecting.
type Handshaker interface ***REMOVED***
	// Handshake should confirm or decorate a connection that may be incoming
	// to a server or outgoing from a client.
	//
	// If this returns without an error, the caller should use the connection
	// in place of the original connection.
	//
	// The second return value can contain credential specific data, such as
	// unix socket credentials or TLS information.
	//
	// While we currently only have implementations on the server-side, this
	// interface should be sufficient to implement similar handshakes on the
	// client-side.
	Handshake(ctx context.Context, conn net.Conn) (net.Conn, interface***REMOVED******REMOVED***, error)
***REMOVED***

type handshakerFunc func(ctx context.Context, conn net.Conn) (net.Conn, interface***REMOVED******REMOVED***, error)

func (fn handshakerFunc) Handshake(ctx context.Context, conn net.Conn) (net.Conn, interface***REMOVED******REMOVED***, error) ***REMOVED***
	return fn(ctx, conn)
***REMOVED***

func noopHandshake(ctx context.Context, conn net.Conn) (net.Conn, interface***REMOVED******REMOVED***, error) ***REMOVED***
	return conn, nil, nil
***REMOVED***
