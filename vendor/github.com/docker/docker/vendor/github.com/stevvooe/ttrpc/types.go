package ttrpc

import (
	"fmt"

	spb "google.golang.org/genproto/googleapis/rpc/status"
)

type Request struct ***REMOVED***
	Service string `protobuf:"bytes,1,opt,name=service,proto3"`
	Method  string `protobuf:"bytes,2,opt,name=method,proto3"`
	Payload []byte `protobuf:"bytes,3,opt,name=payload,proto3"`
***REMOVED***

func (r *Request) Reset()         ***REMOVED*** *r = Request***REMOVED******REMOVED*** ***REMOVED***
func (r *Request) String() string ***REMOVED*** return fmt.Sprintf("%+#v", r) ***REMOVED***
func (r *Request) ProtoMessage()  ***REMOVED******REMOVED***

type Response struct ***REMOVED***
	Status  *spb.Status `protobuf:"bytes,1,opt,name=status,proto3"`
	Payload []byte      `protobuf:"bytes,2,opt,name=payload,proto3"`
***REMOVED***

func (r *Response) Reset()         ***REMOVED*** *r = Response***REMOVED******REMOVED*** ***REMOVED***
func (r *Response) String() string ***REMOVED*** return fmt.Sprintf("%+#v", r) ***REMOVED***
func (r *Response) ProtoMessage()  ***REMOVED******REMOVED***
