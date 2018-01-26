package session

import (
	"net"

	"github.com/docker/docker/pkg/stringid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	headerSessionID        = "X-Docker-Expose-Session-Uuid"
	headerSessionName      = "X-Docker-Expose-Session-Name"
	headerSessionSharedKey = "X-Docker-Expose-Session-Sharedkey"
	headerSessionMethod    = "X-Docker-Expose-Session-Grpc-Method"
)

// Dialer returns a connection that can be used by the session
type Dialer func(ctx context.Context, proto string, meta map[string][]string) (net.Conn, error)

// Attachable defines a feature that can be expsed on a session
type Attachable interface ***REMOVED***
	Register(*grpc.Server)
***REMOVED***

// Session is a long running connection between client and a daemon
type Session struct ***REMOVED***
	id         string
	name       string
	sharedKey  string
	ctx        context.Context
	cancelCtx  func()
	done       chan struct***REMOVED******REMOVED***
	grpcServer *grpc.Server
***REMOVED***

// NewSession returns a new long running session
func NewSession(name, sharedKey string) (*Session, error) ***REMOVED***
	id := stringid.GenerateRandomID()
	s := &Session***REMOVED***
		id:         id,
		name:       name,
		sharedKey:  sharedKey,
		grpcServer: grpc.NewServer(),
	***REMOVED***

	grpc_health_v1.RegisterHealthServer(s.grpcServer, health.NewServer())

	return s, nil
***REMOVED***

// Allow enable a given service to be reachable through the grpc session
func (s *Session) Allow(a Attachable) ***REMOVED***
	a.Register(s.grpcServer)
***REMOVED***

// ID returns unique identifier for the session
func (s *Session) ID() string ***REMOVED***
	return s.id
***REMOVED***

// Run activates the session
func (s *Session) Run(ctx context.Context, dialer Dialer) error ***REMOVED***
	ctx, cancel := context.WithCancel(ctx)
	s.cancelCtx = cancel
	s.done = make(chan struct***REMOVED******REMOVED***)

	defer cancel()
	defer close(s.done)

	meta := make(map[string][]string)
	meta[headerSessionID] = []string***REMOVED***s.id***REMOVED***
	meta[headerSessionName] = []string***REMOVED***s.name***REMOVED***
	meta[headerSessionSharedKey] = []string***REMOVED***s.sharedKey***REMOVED***

	for name, svc := range s.grpcServer.GetServiceInfo() ***REMOVED***
		for _, method := range svc.Methods ***REMOVED***
			meta[headerSessionMethod] = append(meta[headerSessionMethod], MethodURL(name, method.Name))
		***REMOVED***
	***REMOVED***
	conn, err := dialer(ctx, "h2c", meta)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to dial gRPC")
	***REMOVED***
	serve(ctx, s.grpcServer, conn)
	return nil
***REMOVED***

// Close closes the session
func (s *Session) Close() error ***REMOVED***
	if s.cancelCtx != nil && s.done != nil ***REMOVED***
		s.grpcServer.Stop()
		s.cancelCtx()
		<-s.done
	***REMOVED***
	return nil
***REMOVED***

func (s *Session) context() context.Context ***REMOVED***
	return s.ctx
***REMOVED***

func (s *Session) closed() bool ***REMOVED***
	select ***REMOVED***
	case <-s.context().Done():
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

// MethodURL returns a gRPC method URL for service and method name
func MethodURL(s, m string) string ***REMOVED***
	return "/" + s + "/" + m
***REMOVED***
