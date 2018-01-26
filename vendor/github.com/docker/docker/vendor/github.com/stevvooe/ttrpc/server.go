package ttrpc

import (
	"context"
	"io"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/containerd/containerd/log"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrServerClosed = errors.New("ttrpc: server close")
)

type Server struct ***REMOVED***
	config   *serverConfig
	services *serviceSet
	codec    codec

	mu          sync.Mutex
	listeners   map[net.Listener]struct***REMOVED******REMOVED***
	connections map[*serverConn]struct***REMOVED******REMOVED*** // all connections to current state
	done        chan struct***REMOVED******REMOVED***            // marks point at which we stop serving requests
***REMOVED***

func NewServer(opts ...ServerOpt) (*Server, error) ***REMOVED***
	config := &serverConfig***REMOVED******REMOVED***
	for _, opt := range opts ***REMOVED***
		if err := opt(config); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return &Server***REMOVED***
		config:      config,
		services:    newServiceSet(),
		done:        make(chan struct***REMOVED******REMOVED***),
		listeners:   make(map[net.Listener]struct***REMOVED******REMOVED***),
		connections: make(map[*serverConn]struct***REMOVED******REMOVED***),
	***REMOVED***, nil
***REMOVED***

func (s *Server) Register(name string, methods map[string]Method) ***REMOVED***
	s.services.register(name, methods)
***REMOVED***

func (s *Server) Serve(l net.Listener) error ***REMOVED***
	s.addListener(l)
	defer s.closeListener(l)

	var (
		ctx        = context.Background()
		backoff    time.Duration
		handshaker = s.config.handshaker
	)

	if handshaker == nil ***REMOVED***
		handshaker = handshakerFunc(noopHandshake)
	***REMOVED***

	for ***REMOVED***
		conn, err := l.Accept()
		if err != nil ***REMOVED***
			select ***REMOVED***
			case <-s.done:
				return ErrServerClosed
			default:
			***REMOVED***

			if terr, ok := err.(interface ***REMOVED***
				Temporary() bool
			***REMOVED***); ok && terr.Temporary() ***REMOVED***
				if backoff == 0 ***REMOVED***
					backoff = time.Millisecond
				***REMOVED*** else ***REMOVED***
					backoff *= 2
				***REMOVED***

				if max := time.Second; backoff > max ***REMOVED***
					backoff = max
				***REMOVED***

				sleep := time.Duration(rand.Int63n(int64(backoff)))
				log.L.WithError(err).Errorf("ttrpc: failed accept; backoff %v", sleep)
				time.Sleep(sleep)
				continue
			***REMOVED***

			return err
		***REMOVED***

		backoff = 0

		approved, handshake, err := handshaker.Handshake(ctx, conn)
		if err != nil ***REMOVED***
			log.L.WithError(err).Errorf("ttrpc: refusing connection after handshake")
			conn.Close()
			continue
		***REMOVED***

		sc := s.newConn(approved, handshake)
		go sc.run(ctx)
	***REMOVED***
***REMOVED***

func (s *Server) Shutdown(ctx context.Context) error ***REMOVED***
	s.mu.Lock()
	lnerr := s.closeListeners()
	select ***REMOVED***
	case <-s.done:
	default:
		// protected by mutex
		close(s.done)
	***REMOVED***
	s.mu.Unlock()

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	for ***REMOVED***
		if s.closeIdleConns() ***REMOVED***
			return lnerr
		***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		***REMOVED***
	***REMOVED***
***REMOVED***

// Close the server without waiting for active connections.
func (s *Server) Close() error ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	select ***REMOVED***
	case <-s.done:
	default:
		// protected by mutex
		close(s.done)
	***REMOVED***

	err := s.closeListeners()
	for c := range s.connections ***REMOVED***
		c.close()
		delete(s.connections, c)
	***REMOVED***

	return err
***REMOVED***

func (s *Server) addListener(l net.Listener) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	s.listeners[l] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
***REMOVED***

func (s *Server) closeListener(l net.Listener) error ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.closeListenerLocked(l)
***REMOVED***

func (s *Server) closeListenerLocked(l net.Listener) error ***REMOVED***
	defer delete(s.listeners, l)
	return l.Close()
***REMOVED***

func (s *Server) closeListeners() error ***REMOVED***
	var err error
	for l := range s.listeners ***REMOVED***
		if cerr := s.closeListenerLocked(l); cerr != nil && err == nil ***REMOVED***
			err = cerr
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***

func (s *Server) addConnection(c *serverConn) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	s.connections[c] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
***REMOVED***

func (s *Server) closeIdleConns() bool ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	quiescent := true
	for c := range s.connections ***REMOVED***
		st, ok := c.getState()
		if !ok || st != connStateIdle ***REMOVED***
			quiescent = false
			continue
		***REMOVED***
		c.close()
		delete(s.connections, c)
	***REMOVED***
	return quiescent
***REMOVED***

type connState int

const (
	connStateActive = iota + 1 // outstanding requests
	connStateIdle              // no requests
	connStateClosed            // closed connection
)

func (cs connState) String() string ***REMOVED***
	switch cs ***REMOVED***
	case connStateActive:
		return "active"
	case connStateIdle:
		return "idle"
	case connStateClosed:
		return "closed"
	default:
		return "unknown"
	***REMOVED***
***REMOVED***

func (s *Server) newConn(conn net.Conn, handshake interface***REMOVED******REMOVED***) *serverConn ***REMOVED***
	c := &serverConn***REMOVED***
		server:    s,
		conn:      conn,
		handshake: handshake,
		shutdown:  make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
	c.setState(connStateIdle)
	s.addConnection(c)
	return c
***REMOVED***

type serverConn struct ***REMOVED***
	server    *Server
	conn      net.Conn
	handshake interface***REMOVED******REMOVED*** // data from handshake, not used for now
	state     atomic.Value

	shutdownOnce sync.Once
	shutdown     chan struct***REMOVED******REMOVED*** // forced shutdown, used by close
***REMOVED***

func (c *serverConn) getState() (connState, bool) ***REMOVED***
	cs, ok := c.state.Load().(connState)
	return cs, ok
***REMOVED***

func (c *serverConn) setState(newstate connState) ***REMOVED***
	c.state.Store(newstate)
***REMOVED***

func (c *serverConn) close() error ***REMOVED***
	c.shutdownOnce.Do(func() ***REMOVED***
		close(c.shutdown)
	***REMOVED***)

	return nil
***REMOVED***

func (c *serverConn) run(sctx context.Context) ***REMOVED***
	type (
		request struct ***REMOVED***
			id  uint32
			req *Request
		***REMOVED***

		response struct ***REMOVED***
			id   uint32
			resp *Response
		***REMOVED***
	)

	var (
		ch          = newChannel(c.conn, c.conn)
		ctx, cancel = context.WithCancel(sctx)
		active      int
		state       connState = connStateIdle
		responses             = make(chan response)
		requests              = make(chan request)
		recvErr               = make(chan error, 1)
		shutdown              = c.shutdown
		done                  = make(chan struct***REMOVED******REMOVED***)
	)

	defer c.conn.Close()
	defer cancel()
	defer close(done)

	go func(recvErr chan error) ***REMOVED***
		defer close(recvErr)
		sendImmediate := func(id uint32, st *status.Status) bool ***REMOVED***
			select ***REMOVED***
			case responses <- response***REMOVED***
				// even though we've had an invalid stream id, we send it
				// back on the same stream id so the client knows which
				// stream id was bad.
				id: id,
				resp: &Response***REMOVED***
					Status: st.Proto(),
				***REMOVED***,
			***REMOVED***:
				return true
			case <-c.shutdown:
				return false
			case <-done:
				return false
			***REMOVED***
		***REMOVED***

		for ***REMOVED***
			select ***REMOVED***
			case <-c.shutdown:
				return
			case <-done:
				return
			default: // proceed
			***REMOVED***

			mh, p, err := ch.recv(ctx)
			if err != nil ***REMOVED***
				status, ok := status.FromError(err)
				if !ok ***REMOVED***
					recvErr <- err
					return
				***REMOVED***

				// in this case, we send an error for that particular message
				// when the status is defined.
				if !sendImmediate(mh.StreamID, status) ***REMOVED***
					return
				***REMOVED***

				continue
			***REMOVED***

			if mh.Type != messageTypeRequest ***REMOVED***
				// we must ignore this for future compat.
				continue
			***REMOVED***

			var req Request
			if err := c.server.codec.Unmarshal(p, &req); err != nil ***REMOVED***
				ch.putmbuf(p)
				if !sendImmediate(mh.StreamID, status.Newf(codes.InvalidArgument, "unmarshal request error: %v", err)) ***REMOVED***
					return
				***REMOVED***
				continue
			***REMOVED***
			ch.putmbuf(p)

			if mh.StreamID%2 != 1 ***REMOVED***
				// enforce odd client initiated identifiers.
				if !sendImmediate(mh.StreamID, status.Newf(codes.InvalidArgument, "StreamID must be odd for client initiated streams")) ***REMOVED***
					return
				***REMOVED***
				continue
			***REMOVED***

			// Forward the request to the main loop. We don't wait on s.done
			// because we have already accepted the client request.
			select ***REMOVED***
			case requests <- request***REMOVED***
				id:  mh.StreamID,
				req: &req,
			***REMOVED***:
			case <-done:
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***(recvErr)

	for ***REMOVED***
		newstate := state
		switch ***REMOVED***
		case active > 0:
			newstate = connStateActive
			shutdown = nil
		case active == 0:
			newstate = connStateIdle
			shutdown = c.shutdown // only enable this branch in idle mode
		***REMOVED***

		if newstate != state ***REMOVED***
			c.setState(newstate)
			state = newstate
		***REMOVED***

		select ***REMOVED***
		case request := <-requests:
			active++
			go func(id uint32) ***REMOVED***
				p, status := c.server.services.call(ctx, request.req.Service, request.req.Method, request.req.Payload)
				resp := &Response***REMOVED***
					Status:  status.Proto(),
					Payload: p,
				***REMOVED***

				select ***REMOVED***
				case responses <- response***REMOVED***
					id:   id,
					resp: resp,
				***REMOVED***:
				case <-done:
				***REMOVED***
			***REMOVED***(request.id)
		case response := <-responses:
			p, err := c.server.codec.Marshal(response.resp)
			if err != nil ***REMOVED***
				log.L.WithError(err).Error("failed marshaling response")
				return
			***REMOVED***

			if err := ch.send(ctx, response.id, messageTypeResponse, p); err != nil ***REMOVED***
				log.L.WithError(err).Error("failed sending message on channel")
				return
			***REMOVED***

			active--
		case err := <-recvErr:
			// TODO(stevvooe): Not wildly clear what we should do in this
			// branch. Basically, it means that we are no longer receiving
			// requests due to a terminal error.
			recvErr = nil // connection is now "closing"
			if err != nil && err != io.EOF ***REMOVED***
				log.L.WithError(err).Error("error receiving message")
			***REMOVED***
		case <-shutdown:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***
