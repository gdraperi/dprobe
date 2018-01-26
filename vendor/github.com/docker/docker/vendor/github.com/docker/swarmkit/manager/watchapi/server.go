package watchapi

import (
	"errors"
	"sync"

	"github.com/docker/swarmkit/manager/state/store"
	"golang.org/x/net/context"
)

var (
	errAlreadyRunning = errors.New("broker is already running")
	errNotRunning     = errors.New("broker is not running")
)

// Server is the store API gRPC server.
type Server struct ***REMOVED***
	store     *store.MemoryStore
	mu        sync.Mutex
	pctx      context.Context
	cancelAll func()
***REMOVED***

// NewServer creates a store API server.
func NewServer(store *store.MemoryStore) *Server ***REMOVED***
	return &Server***REMOVED***
		store: store,
	***REMOVED***
***REMOVED***

// Start starts the watch server.
func (s *Server) Start(ctx context.Context) error ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancelAll != nil ***REMOVED***
		return errAlreadyRunning
	***REMOVED***

	s.pctx, s.cancelAll = context.WithCancel(ctx)
	return nil
***REMOVED***

// Stop stops the watch server.
func (s *Server) Stop() error ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancelAll == nil ***REMOVED***
		return errNotRunning
	***REMOVED***
	s.cancelAll()
	s.cancelAll = nil

	return nil
***REMOVED***
