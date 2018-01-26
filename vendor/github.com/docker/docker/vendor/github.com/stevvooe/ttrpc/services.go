package ttrpc

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Method func(ctx context.Context, unmarshal func(interface***REMOVED******REMOVED***) error) (interface***REMOVED******REMOVED***, error)

type ServiceDesc struct ***REMOVED***
	Methods map[string]Method

	// TODO(stevvooe): Add stream support.
***REMOVED***

type serviceSet struct ***REMOVED***
	services map[string]ServiceDesc
***REMOVED***

func newServiceSet() *serviceSet ***REMOVED***
	return &serviceSet***REMOVED***
		services: make(map[string]ServiceDesc),
	***REMOVED***
***REMOVED***

func (s *serviceSet) register(name string, methods map[string]Method) ***REMOVED***
	if _, ok := s.services[name]; ok ***REMOVED***
		panic(errors.Errorf("duplicate service %v registered", name))
	***REMOVED***

	s.services[name] = ServiceDesc***REMOVED***
		Methods: methods,
	***REMOVED***
***REMOVED***

func (s *serviceSet) call(ctx context.Context, serviceName, methodName string, p []byte) ([]byte, *status.Status) ***REMOVED***
	p, err := s.dispatch(ctx, serviceName, methodName, p)
	st, ok := status.FromError(err)
	if !ok ***REMOVED***
		st = status.New(convertCode(err), err.Error())
	***REMOVED***

	return p, st
***REMOVED***

func (s *serviceSet) dispatch(ctx context.Context, serviceName, methodName string, p []byte) ([]byte, error) ***REMOVED***
	method, err := s.resolve(serviceName, methodName)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	unmarshal := func(obj interface***REMOVED******REMOVED***) error ***REMOVED***
		switch v := obj.(type) ***REMOVED***
		case proto.Message:
			if err := proto.Unmarshal(p, v); err != nil ***REMOVED***
				return status.Errorf(codes.Internal, "ttrpc: error unmarshaling payload: %v", err.Error())
			***REMOVED***
		default:
			return status.Errorf(codes.Internal, "ttrpc: error unsupported request type: %T", v)
		***REMOVED***
		return nil
	***REMOVED***

	resp, err := method(ctx, unmarshal)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	switch v := resp.(type) ***REMOVED***
	case proto.Message:
		r, err := proto.Marshal(v)
		if err != nil ***REMOVED***
			return nil, status.Errorf(codes.Internal, "ttrpc: error marshaling payload: %v", err.Error())
		***REMOVED***

		return r, nil
	default:
		return nil, status.Errorf(codes.Internal, "ttrpc: error unsupported response type: %T", v)
	***REMOVED***
***REMOVED***

func (s *serviceSet) resolve(service, method string) (Method, error) ***REMOVED***
	srv, ok := s.services[service]
	if !ok ***REMOVED***
		return nil, status.Errorf(codes.NotFound, "service %v", service)
	***REMOVED***

	mthd, ok := srv.Methods[method]
	if !ok ***REMOVED***
		return nil, status.Errorf(codes.NotFound, "method %v", method)
	***REMOVED***

	return mthd, nil
***REMOVED***

// convertCode maps stdlib go errors into grpc space.
//
// This is ripped from the grpc-go code base.
func convertCode(err error) codes.Code ***REMOVED***
	switch err ***REMOVED***
	case nil:
		return codes.OK
	case io.EOF:
		return codes.OutOfRange
	case io.ErrClosedPipe, io.ErrNoProgress, io.ErrShortBuffer, io.ErrShortWrite, io.ErrUnexpectedEOF:
		return codes.FailedPrecondition
	case os.ErrInvalid:
		return codes.InvalidArgument
	case context.Canceled:
		return codes.Canceled
	case context.DeadlineExceeded:
		return codes.DeadlineExceeded
	***REMOVED***
	switch ***REMOVED***
	case os.IsExist(err):
		return codes.AlreadyExists
	case os.IsNotExist(err):
		return codes.NotFound
	case os.IsPermission(err):
		return codes.PermissionDenied
	***REMOVED***
	return codes.Unknown
***REMOVED***

func fullPath(service, method string) string ***REMOVED***
	return "/" + path.Join("/", service, method)
***REMOVED***
