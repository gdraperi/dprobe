package session

import (
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Caller can invoke requests on the session
type Caller interface ***REMOVED***
	Context() context.Context
	Supports(method string) bool
	Conn() *grpc.ClientConn
	Name() string
	SharedKey() string
***REMOVED***

type client struct ***REMOVED***
	Session
	cc        *grpc.ClientConn
	supported map[string]struct***REMOVED******REMOVED***
***REMOVED***

// Manager is a controller for accessing currently active sessions
type Manager struct ***REMOVED***
	sessions        map[string]*client
	mu              sync.Mutex
	updateCondition *sync.Cond
***REMOVED***

// NewManager returns a new Manager
func NewManager() (*Manager, error) ***REMOVED***
	sm := &Manager***REMOVED***
		sessions: make(map[string]*client),
	***REMOVED***
	sm.updateCondition = sync.NewCond(&sm.mu)
	return sm, nil
***REMOVED***

// HandleHTTPRequest handles an incoming HTTP request
func (sm *Manager) HandleHTTPRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) error ***REMOVED***
	hijacker, ok := w.(http.Hijacker)
	if !ok ***REMOVED***
		return errors.New("handler does not support hijack")
	***REMOVED***

	id := r.Header.Get(headerSessionID)

	proto := r.Header.Get("Upgrade")

	sm.mu.Lock()
	if _, ok := sm.sessions[id]; ok ***REMOVED***
		sm.mu.Unlock()
		return errors.Errorf("session %s already exists", id)
	***REMOVED***

	if proto == "" ***REMOVED***
		sm.mu.Unlock()
		return errors.New("no upgrade proto in request")
	***REMOVED***

	if proto != "h2c" ***REMOVED***
		sm.mu.Unlock()
		return errors.Errorf("protocol %s not supported", proto)
	***REMOVED***

	conn, _, err := hijacker.Hijack()
	if err != nil ***REMOVED***
		sm.mu.Unlock()
		return errors.Wrap(err, "failed to hijack connection")
	***REMOVED***

	resp := &http.Response***REMOVED***
		StatusCode: http.StatusSwitchingProtocols,
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     http.Header***REMOVED******REMOVED***,
	***REMOVED***
	resp.Header.Set("Connection", "Upgrade")
	resp.Header.Set("Upgrade", proto)

	// set raw mode
	conn.Write([]byte***REMOVED******REMOVED***)
	resp.Write(conn)

	return sm.handleConn(ctx, conn, r.Header)
***REMOVED***

// HandleConn handles an incoming raw connection
func (sm *Manager) HandleConn(ctx context.Context, conn net.Conn, opts map[string][]string) error ***REMOVED***
	sm.mu.Lock()
	return sm.handleConn(ctx, conn, opts)
***REMOVED***

// caller needs to take lock, this function will release it
func (sm *Manager) handleConn(ctx context.Context, conn net.Conn, opts map[string][]string) error ***REMOVED***
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opts = canonicalHeaders(opts)

	h := http.Header(opts)
	id := h.Get(headerSessionID)
	name := h.Get(headerSessionName)
	sharedKey := h.Get(headerSessionSharedKey)

	ctx, cc, err := grpcClientConn(ctx, conn)
	if err != nil ***REMOVED***
		sm.mu.Unlock()
		return err
	***REMOVED***

	c := &client***REMOVED***
		Session: Session***REMOVED***
			id:        id,
			name:      name,
			sharedKey: sharedKey,
			ctx:       ctx,
			cancelCtx: cancel,
			done:      make(chan struct***REMOVED******REMOVED***),
		***REMOVED***,
		cc:        cc,
		supported: make(map[string]struct***REMOVED******REMOVED***),
	***REMOVED***

	for _, m := range opts[headerSessionMethod] ***REMOVED***
		c.supported[strings.ToLower(m)] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	sm.sessions[id] = c
	sm.updateCondition.Broadcast()
	sm.mu.Unlock()

	defer func() ***REMOVED***
		sm.mu.Lock()
		delete(sm.sessions, id)
		sm.mu.Unlock()
	***REMOVED***()

	<-c.ctx.Done()
	conn.Close()
	close(c.done)

	return nil
***REMOVED***

// Get returns a session by ID
func (sm *Manager) Get(ctx context.Context, id string) (Caller, error) ***REMOVED***
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			sm.updateCondition.Broadcast()
		***REMOVED***
	***REMOVED***()

	var c *client

	sm.mu.Lock()
	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			sm.mu.Unlock()
			return nil, errors.Wrapf(ctx.Err(), "no active session for %s", id)
		default:
		***REMOVED***
		var ok bool
		c, ok = sm.sessions[id]
		if !ok || c.closed() ***REMOVED***
			sm.updateCondition.Wait()
			continue
		***REMOVED***
		sm.mu.Unlock()
		break
	***REMOVED***

	return c, nil
***REMOVED***

func (c *client) Context() context.Context ***REMOVED***
	return c.context()
***REMOVED***

func (c *client) Name() string ***REMOVED***
	return c.name
***REMOVED***

func (c *client) SharedKey() string ***REMOVED***
	return c.sharedKey
***REMOVED***

func (c *client) Supports(url string) bool ***REMOVED***
	_, ok := c.supported[strings.ToLower(url)]
	return ok
***REMOVED***
func (c *client) Conn() *grpc.ClientConn ***REMOVED***
	return c.cc
***REMOVED***

func canonicalHeaders(in map[string][]string) map[string][]string ***REMOVED***
	out := map[string][]string***REMOVED******REMOVED***
	for k := range in ***REMOVED***
		out[http.CanonicalHeaderKey(k)] = in[k]
	***REMOVED***
	return out
***REMOVED***
