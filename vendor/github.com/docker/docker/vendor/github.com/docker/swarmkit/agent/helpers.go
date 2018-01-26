package agent

import "golang.org/x/net/context"

// runctx blocks until the function exits, closed is closed, or the context is
// cancelled. Call as part of go statement.
func runctx(ctx context.Context, closed chan struct***REMOVED******REMOVED***, errs chan error, fn func(ctx context.Context) error) ***REMOVED***
	select ***REMOVED***
	case errs <- fn(ctx):
	case <-closed:
	case <-ctx.Done():
	***REMOVED***
***REMOVED***
