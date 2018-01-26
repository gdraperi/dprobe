// Package health provides some utility functions to health-check a server. The implementation
// is based on protobuf. Users need to write their own implementations if other IDLs are used.
package health

import (
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Server implements `service Health`.
type Server struct ***REMOVED***
	mu sync.Mutex
	// statusMap stores the serving status of the services this Server monitors.
	statusMap map[string]healthpb.HealthCheckResponse_ServingStatus
***REMOVED***

// NewServer returns a new Server.
func NewServer() *Server ***REMOVED***
	return &Server***REMOVED***
		statusMap: make(map[string]healthpb.HealthCheckResponse_ServingStatus),
	***REMOVED***
***REMOVED***

// Check implements `service Health`.
func (s *Server) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	if in.Service == "" ***REMOVED***
		// check the server overall health status.
		return &healthpb.HealthCheckResponse***REMOVED***
			Status: healthpb.HealthCheckResponse_SERVING,
		***REMOVED***, nil
	***REMOVED***
	if status, ok := s.statusMap[in.Service]; ok ***REMOVED***
		return &healthpb.HealthCheckResponse***REMOVED***
			Status: status,
		***REMOVED***, nil
	***REMOVED***
	return nil, grpc.Errorf(codes.NotFound, "unknown service")
***REMOVED***

// SetServingStatus is called when need to reset the serving status of a service
// or insert a new service entry into the statusMap.
func (s *Server) SetServingStatus(service string, status healthpb.HealthCheckResponse_ServingStatus) ***REMOVED***
	s.mu.Lock()
	s.statusMap[service] = status
	s.mu.Unlock()
***REMOVED***
