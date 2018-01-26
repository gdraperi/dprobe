// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

// gRPC Prometheus monitoring interceptors for server-side gRPC.

package grpc_prometheus

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// PreregisterServices takes a gRPC server and pre-initializes all counters to 0.
// This allows for easier monitoring in Prometheus (no missing metrics), and should be called *after* all services have
// been registered with the server.
func Register(server *grpc.Server) ***REMOVED***
	serviceInfo := server.GetServiceInfo()
	for serviceName, info := range serviceInfo ***REMOVED***
		for _, mInfo := range info.Methods ***REMOVED***
			preRegisterMethod(serviceName, &mInfo)
		***REMOVED***
	***REMOVED***
***REMOVED***

// UnaryServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Unary RPCs.
func UnaryServerInterceptor(ctx context.Context, req interface***REMOVED******REMOVED***, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	monitor := newServerReporter(Unary, info.FullMethod)
	monitor.ReceivedMessage()
	resp, err := handler(ctx, req)
	monitor.Handled(grpc.Code(err))
	if err == nil ***REMOVED***
		monitor.SentMessage()
	***REMOVED***
	return resp, err
***REMOVED***

// StreamServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func StreamServerInterceptor(srv interface***REMOVED******REMOVED***, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error ***REMOVED***
	monitor := newServerReporter(streamRpcType(info), info.FullMethod)
	err := handler(srv, &monitoredServerStream***REMOVED***ss, monitor***REMOVED***)
	monitor.Handled(grpc.Code(err))
	return err
***REMOVED***

func streamRpcType(info *grpc.StreamServerInfo) grpcType ***REMOVED***
	if info.IsClientStream && !info.IsServerStream ***REMOVED***
		return ClientStream
	***REMOVED*** else if !info.IsClientStream && info.IsServerStream ***REMOVED***
		return ServerStream
	***REMOVED***
	return BidiStream
***REMOVED***

// monitoredStream wraps grpc.ServerStream allowing each Sent/Recv of message to increment counters.
type monitoredServerStream struct ***REMOVED***
	grpc.ServerStream
	monitor *serverReporter
***REMOVED***

func (s *monitoredServerStream) SendMsg(m interface***REMOVED******REMOVED***) error ***REMOVED***
	err := s.ServerStream.SendMsg(m)
	if err == nil ***REMOVED***
		s.monitor.SentMessage()
	***REMOVED***
	return err
***REMOVED***

func (s *monitoredServerStream) RecvMsg(m interface***REMOVED******REMOVED***) error ***REMOVED***
	err := s.ServerStream.RecvMsg(m)
	if err == nil ***REMOVED***
		s.monitor.ReceivedMessage()
	***REMOVED***
	return err
***REMOVED***
