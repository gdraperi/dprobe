package containerd

import (
	"github.com/containerd/containerd/namespaces"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type namespaceInterceptor struct ***REMOVED***
	namespace string
***REMOVED***

func (ni namespaceInterceptor) unary(ctx context.Context, method string, req, reply interface***REMOVED******REMOVED***, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error ***REMOVED***
	_, ok := namespaces.Namespace(ctx)
	if !ok ***REMOVED***
		ctx = namespaces.WithNamespace(ctx, ni.namespace)
	***REMOVED***
	return invoker(ctx, method, req, reply, cc, opts...)
***REMOVED***

func (ni namespaceInterceptor) stream(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) ***REMOVED***
	_, ok := namespaces.Namespace(ctx)
	if !ok ***REMOVED***
		ctx = namespaces.WithNamespace(ctx, ni.namespace)
	***REMOVED***

	return streamer(ctx, desc, cc, method, opts...)
***REMOVED***

func newNSInterceptors(ns string) (grpc.UnaryClientInterceptor, grpc.StreamClientInterceptor) ***REMOVED***
	ni := namespaceInterceptor***REMOVED***
		namespace: ns,
	***REMOVED***
	return grpc.UnaryClientInterceptor(ni.unary), grpc.StreamClientInterceptor(ni.stream)
***REMOVED***
