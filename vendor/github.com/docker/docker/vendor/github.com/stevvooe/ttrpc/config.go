package ttrpc

import "github.com/pkg/errors"

type serverConfig struct ***REMOVED***
	handshaker Handshaker
***REMOVED***

type ServerOpt func(*serverConfig) error

// WithServerHandshaker can be passed to NewServer to ensure that the
// handshaker is called before every connection attempt.
//
// Only one handshaker is allowed per server.
func WithServerHandshaker(handshaker Handshaker) ServerOpt ***REMOVED***
	return func(c *serverConfig) error ***REMOVED***
		if c.handshaker != nil ***REMOVED***
			return errors.New("only one handshaker allowed per server")
		***REMOVED***
		c.handshaker = handshaker
		return nil
	***REMOVED***
***REMOVED***
