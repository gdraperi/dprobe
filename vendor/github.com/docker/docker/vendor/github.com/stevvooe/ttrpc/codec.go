package ttrpc

import (
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

type codec struct***REMOVED******REMOVED***

func (c codec) Marshal(msg interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	switch v := msg.(type) ***REMOVED***
	case proto.Message:
		return proto.Marshal(v)
	default:
		return nil, errors.Errorf("ttrpc: cannot marshal unknown type: %T", msg)
	***REMOVED***
***REMOVED***

func (c codec) Unmarshal(p []byte, msg interface***REMOVED******REMOVED***) error ***REMOVED***
	switch v := msg.(type) ***REMOVED***
	case proto.Message:
		return proto.Unmarshal(p, v)
	default:
		return errors.Errorf("ttrpc: cannot unmarshal into unknown type: %T", msg)
	***REMOVED***
***REMOVED***
