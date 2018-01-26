package session

import "context"

type contextKeyT string

var contextKey = contextKeyT("buildkit/session-id")

func NewContext(ctx context.Context, id string) context.Context ***REMOVED***
	if id != "" ***REMOVED***
		return context.WithValue(ctx, contextKey, id)
	***REMOVED***
	return ctx
***REMOVED***

func FromContext(ctx context.Context) string ***REMOVED***
	v := ctx.Value(contextKey)
	if v == nil ***REMOVED***
		return ""
	***REMOVED***
	return v.(string)
***REMOVED***
