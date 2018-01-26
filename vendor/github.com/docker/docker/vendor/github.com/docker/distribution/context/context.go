package context

import (
	"sync"

	"github.com/docker/distribution/uuid"
	"golang.org/x/net/context"
)

// Context is a copy of Context from the golang.org/x/net/context package.
type Context interface ***REMOVED***
	context.Context
***REMOVED***

// instanceContext is a context that provides only an instance id. It is
// provided as the main background context.
type instanceContext struct ***REMOVED***
	Context
	id   string    // id of context, logged as "instance.id"
	once sync.Once // once protect generation of the id
***REMOVED***

func (ic *instanceContext) Value(key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	if key == "instance.id" ***REMOVED***
		ic.once.Do(func() ***REMOVED***
			// We want to lazy initialize the UUID such that we don't
			// call a random generator from the package initialization
			// code. For various reasons random could not be available
			// https://github.com/docker/distribution/issues/782
			ic.id = uuid.Generate().String()
		***REMOVED***)
		return ic.id
	***REMOVED***

	return ic.Context.Value(key)
***REMOVED***

var background = &instanceContext***REMOVED***
	Context: context.Background(),
***REMOVED***

// Background returns a non-nil, empty Context. The background context
// provides a single key, "instance.id" that is globally unique to the
// process.
func Background() Context ***REMOVED***
	return background
***REMOVED***

// WithValue returns a copy of parent in which the value associated with key is
// val. Use context Values only for request-scoped data that transits processes
// and APIs, not for passing optional parameters to functions.
func WithValue(parent Context, key, val interface***REMOVED******REMOVED***) Context ***REMOVED***
	return context.WithValue(parent, key, val)
***REMOVED***

// stringMapContext is a simple context implementation that checks a map for a
// key, falling back to a parent if not present.
type stringMapContext struct ***REMOVED***
	context.Context
	m map[string]interface***REMOVED******REMOVED***
***REMOVED***

// WithValues returns a context that proxies lookups through a map. Only
// supports string keys.
func WithValues(ctx context.Context, m map[string]interface***REMOVED******REMOVED***) context.Context ***REMOVED***
	mo := make(map[string]interface***REMOVED******REMOVED***, len(m)) // make our own copy.
	for k, v := range m ***REMOVED***
		mo[k] = v
	***REMOVED***

	return stringMapContext***REMOVED***
		Context: ctx,
		m:       mo,
	***REMOVED***
***REMOVED***

func (smc stringMapContext) Value(key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	if ks, ok := key.(string); ok ***REMOVED***
		if v, ok := smc.m[ks]; ok ***REMOVED***
			return v
		***REMOVED***
	***REMOVED***

	return smc.Context.Value(key)
***REMOVED***
