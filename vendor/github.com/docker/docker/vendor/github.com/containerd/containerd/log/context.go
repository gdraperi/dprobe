package log

import (
	"context"
	"path"

	"github.com/sirupsen/logrus"
)

var (
	// G is an alias for GetLogger.
	//
	// We may want to define this locally to a package to get package tagged log
	// messages.
	G = GetLogger

	// L is an alias for the the standard logger.
	L = logrus.NewEntry(logrus.StandardLogger())
)

type (
	loggerKey struct***REMOVED******REMOVED***
	moduleKey struct***REMOVED******REMOVED***
)

// WithLogger returns a new context with the provided logger. Use in
// combination with logger.WithField(s) for great effect.
func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context ***REMOVED***
	return context.WithValue(ctx, loggerKey***REMOVED******REMOVED***, logger)
***REMOVED***

// GetLogger retrieves the current logger from the context. If no logger is
// available, the default logger is returned.
func GetLogger(ctx context.Context) *logrus.Entry ***REMOVED***
	logger := ctx.Value(loggerKey***REMOVED******REMOVED***)

	if logger == nil ***REMOVED***
		return L
	***REMOVED***

	return logger.(*logrus.Entry)
***REMOVED***

// WithModule adds the module to the context, appending it with a slash if a
// module already exists. A module is just an roughly correlated defined by the
// call tree for a given context.
//
// As an example, we might have a "node" module already part of a context. If
// this function is called with "tls", the new value of module will be
// "node/tls".
//
// Modules represent the call path. If the new module and last module are the
// same, a new module entry will not be created. If the new module and old
// older module are the same but separated by other modules, the cycle will be
// represented by the module path.
func WithModule(ctx context.Context, module string) context.Context ***REMOVED***
	parent := GetModulePath(ctx)

	if parent != "" ***REMOVED***
		// don't re-append module when module is the same.
		if path.Base(parent) == module ***REMOVED***
			return ctx
		***REMOVED***

		module = path.Join(parent, module)
	***REMOVED***

	ctx = WithLogger(ctx, GetLogger(ctx).WithField("module", module))
	return context.WithValue(ctx, moduleKey***REMOVED******REMOVED***, module)
***REMOVED***

// GetModulePath returns the module path for the provided context. If no module
// is set, an empty string is returned.
func GetModulePath(ctx context.Context) string ***REMOVED***
	module := ctx.Value(moduleKey***REMOVED******REMOVED***)
	if module == nil ***REMOVED***
		return ""
	***REMOVED***

	return module.(string)
***REMOVED***
