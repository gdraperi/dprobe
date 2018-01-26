package typeurl

import (
	"encoding/json"
	"path"
	"reflect"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
)

var (
	mu       sync.Mutex
	registry = make(map[reflect.Type]string)
)

var ErrNotFound = errors.New("not found")

// Register a type with the base url of the type
func Register(v interface***REMOVED******REMOVED***, args ...string) ***REMOVED***
	var (
		t = tryDereference(v)
		p = path.Join(args...)
	)
	mu.Lock()
	defer mu.Unlock()
	if et, ok := registry[t]; ok ***REMOVED***
		if et != p ***REMOVED***
			panic(errors.Errorf("type registred with alternate path %q != %q", et, p))
		***REMOVED***
		return
	***REMOVED***
	registry[t] = p
***REMOVED***

// TypeURL returns the type url for a registred type
func TypeURL(v interface***REMOVED******REMOVED***) (string, error) ***REMOVED***
	mu.Lock()
	u, ok := registry[tryDereference(v)]
	mu.Unlock()
	if !ok ***REMOVED***
		// fallback to the proto registry if it is a proto message
		pb, ok := v.(proto.Message)
		if !ok ***REMOVED***
			return "", errors.Wrapf(ErrNotFound, "type %s", reflect.TypeOf(v))
		***REMOVED***
		return proto.MessageName(pb), nil
	***REMOVED***
	return u, nil
***REMOVED***

// Is returns true if the type of the Any is the same as v
func Is(any *types.Any, v interface***REMOVED******REMOVED***) bool ***REMOVED***
	// call to check that v is a pointer
	tryDereference(v)
	url, err := TypeURL(v)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	return any.TypeUrl == url
***REMOVED***

// MarshalAny marshals the value v into an any with the correct TypeUrl
func MarshalAny(v interface***REMOVED******REMOVED***) (*types.Any, error) ***REMOVED***
	var marshal func(v interface***REMOVED******REMOVED***) ([]byte, error)
	switch t := v.(type) ***REMOVED***
	case *types.Any:
		// avoid reserializing the type if we have an any.
		return t, nil
	case proto.Message:
		marshal = func(v interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
			return proto.Marshal(t)
		***REMOVED***
	default:
		marshal = json.Marshal
	***REMOVED***

	url, err := TypeURL(v)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	data, err := marshal(v)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &types.Any***REMOVED***
		TypeUrl: url,
		Value:   data,
	***REMOVED***, nil
***REMOVED***

// UnmarshalAny unmarshals the any type into a concrete type
func UnmarshalAny(any *types.Any) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	t, err := getTypeByUrl(any.TypeUrl)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	v := reflect.New(t.t).Interface()
	if t.isProto ***REMOVED***
		err = proto.Unmarshal(any.Value, v.(proto.Message))
	***REMOVED*** else ***REMOVED***
		err = json.Unmarshal(any.Value, v)
	***REMOVED***
	return v, err
***REMOVED***

type urlType struct ***REMOVED***
	t       reflect.Type
	isProto bool
***REMOVED***

func getTypeByUrl(url string) (urlType, error) ***REMOVED***
	for t, u := range registry ***REMOVED***
		if u == url ***REMOVED***
			return urlType***REMOVED***
				t: t,
			***REMOVED***, nil
		***REMOVED***
	***REMOVED***
	// fallback to proto registry
	t := proto.MessageType(url)
	if t != nil ***REMOVED***
		return urlType***REMOVED***
			// get the underlying Elem because proto returns a pointer to the type
			t:       t.Elem(),
			isProto: true,
		***REMOVED***, nil
	***REMOVED***
	return urlType***REMOVED******REMOVED***, errors.Wrapf(ErrNotFound, "type with url %s", url)
***REMOVED***

func tryDereference(v interface***REMOVED******REMOVED***) reflect.Type ***REMOVED***
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr ***REMOVED***
		// require check of pointer but dereference to register
		return t.Elem()
	***REMOVED***
	panic("v is not a pointer to a type")
***REMOVED***
