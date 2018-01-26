// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

// gRPC Prometheus monitoring interceptors for client-side gRPC.

package grpc_prometheus

import (
	"io"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// UnaryClientInterceptor is a gRPC client-side interceptor that provides Prometheus monitoring for Unary RPCs.
func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface***REMOVED******REMOVED***, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error ***REMOVED***
	monitor := newClientReporter(Unary, method)
	monitor.SentMessage()
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil ***REMOVED***
		monitor.ReceivedMessage()
	***REMOVED***
	monitor.Handled(grpc.Code(err))
	return err
***REMOVED***

// StreamServerInterceptor is a gRPC client-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func StreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) ***REMOVED***
	monitor := newClientReporter(clientStreamType(desc), method)
	clientStream, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil ***REMOVED***
		monitor.Handled(grpc.Code(err))
		return nil, err
	***REMOVED***
	return &monitoredClientStream***REMOVED***clientStream, monitor***REMOVED***, nil
***REMOVED***

func clientStreamType(desc *grpc.StreamDesc) grpcType ***REMOVED***
	if desc.ClientStreams && !desc.ServerStreams ***REMOVED***
		return ClientStream
	***REMOVED*** else if !desc.ClientStreams && desc.ServerStreams ***REMOVED***
		return ServerStream
	***REMOVED***
	return BidiStream
***REMOVED***

// monitoredClientStream wraps grpc.ClientStream allowing each Sent/Recv of message to increment counters.
type monitoredClientStream struct ***REMOVED***
	grpc.ClientStream
	monitor *clientReporter
***REMOVED***

func (s *monitoredClientStream) SendMsg(m interface***REMOVED******REMOVED***) error ***REMOVED***
	err := s.ClientStream.SendMsg(m)
	if err == nil ***REMOVED***
		s.monitor.SentMessage()
	***REMOVED***
	return err
***REMOVED***

func (s *monitoredClientStream) RecvMsg(m interface***REMOVED******REMOVED***) error ***REMOVED***
	err := s.ClientStream.RecvMsg(m)
	if err == nil ***REMOVED***
		s.monitor.ReceivedMessage()
	***REMOVED*** else if err == io.EOF ***REMOVED***
		s.monitor.Handled(codes.OK)
	***REMOVED*** else ***REMOVED***
		s.monitor.Handled(grpc.Code(err))
	***REMOVED***
	return err
***REMOVED***
