package context

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"runtime"
)

// Logger provides a leveled-logging interface.
type Logger interface ***REMOVED***
	// standard logger methods
	Print(args ...interface***REMOVED******REMOVED***)
	Printf(format string, args ...interface***REMOVED******REMOVED***)
	Println(args ...interface***REMOVED******REMOVED***)

	Fatal(args ...interface***REMOVED******REMOVED***)
	Fatalf(format string, args ...interface***REMOVED******REMOVED***)
	Fatalln(args ...interface***REMOVED******REMOVED***)

	Panic(args ...interface***REMOVED******REMOVED***)
	Panicf(format string, args ...interface***REMOVED******REMOVED***)
	Panicln(args ...interface***REMOVED******REMOVED***)

	// Leveled methods, from logrus
	Debug(args ...interface***REMOVED******REMOVED***)
	Debugf(format string, args ...interface***REMOVED******REMOVED***)
	Debugln(args ...interface***REMOVED******REMOVED***)

	Error(args ...interface***REMOVED******REMOVED***)
	Errorf(format string, args ...interface***REMOVED******REMOVED***)
	Errorln(args ...interface***REMOVED******REMOVED***)

	Info(args ...interface***REMOVED******REMOVED***)
	Infof(format string, args ...interface***REMOVED******REMOVED***)
	Infoln(args ...interface***REMOVED******REMOVED***)

	Warn(args ...interface***REMOVED******REMOVED***)
	Warnf(format string, args ...interface***REMOVED******REMOVED***)
	Warnln(args ...interface***REMOVED******REMOVED***)
***REMOVED***

// WithLogger creates a new context with provided logger.
func WithLogger(ctx Context, logger Logger) Context ***REMOVED***
	return WithValue(ctx, "logger", logger)
***REMOVED***

// GetLoggerWithField returns a logger instance with the specified field key
// and value without affecting the context. Extra specified keys will be
// resolved from the context.
func GetLoggerWithField(ctx Context, key, value interface***REMOVED******REMOVED***, keys ...interface***REMOVED******REMOVED***) Logger ***REMOVED***
	return getLogrusLogger(ctx, keys...).WithField(fmt.Sprint(key), value)
***REMOVED***

// GetLoggerWithFields returns a logger instance with the specified fields
// without affecting the context. Extra specified keys will be resolved from
// the context.
func GetLoggerWithFields(ctx Context, fields map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***, keys ...interface***REMOVED******REMOVED***) Logger ***REMOVED***
	// must convert from interface***REMOVED******REMOVED*** -> interface***REMOVED******REMOVED*** to string -> interface***REMOVED******REMOVED*** for logrus.
	lfields := make(logrus.Fields, len(fields))
	for key, value := range fields ***REMOVED***
		lfields[fmt.Sprint(key)] = value
	***REMOVED***

	return getLogrusLogger(ctx, keys...).WithFields(lfields)
***REMOVED***

// GetLogger returns the logger from the current context, if present. If one
// or more keys are provided, they will be resolved on the context and
// included in the logger. While context.Value takes an interface, any key
// argument passed to GetLogger will be passed to fmt.Sprint when expanded as
// a logging key field. If context keys are integer constants, for example,
// its recommended that a String method is implemented.
func GetLogger(ctx Context, keys ...interface***REMOVED******REMOVED***) Logger ***REMOVED***
	return getLogrusLogger(ctx, keys...)
***REMOVED***

// GetLogrusLogger returns the logrus logger for the context. If one more keys
// are provided, they will be resolved on the context and included in the
// logger. Only use this function if specific logrus functionality is
// required.
func getLogrusLogger(ctx Context, keys ...interface***REMOVED******REMOVED***) *logrus.Entry ***REMOVED***
	var logger *logrus.Entry

	// Get a logger, if it is present.
	loggerInterface := ctx.Value("logger")
	if loggerInterface != nil ***REMOVED***
		if lgr, ok := loggerInterface.(*logrus.Entry); ok ***REMOVED***
			logger = lgr
		***REMOVED***
	***REMOVED***

	if logger == nil ***REMOVED***
		fields := logrus.Fields***REMOVED******REMOVED***

		// Fill in the instance id, if we have it.
		instanceID := ctx.Value("instance.id")
		if instanceID != nil ***REMOVED***
			fields["instance.id"] = instanceID
		***REMOVED***

		fields["go.version"] = runtime.Version()
		// If no logger is found, just return the standard logger.
		logger = logrus.StandardLogger().WithFields(fields)
	***REMOVED***

	fields := logrus.Fields***REMOVED******REMOVED***
	for _, key := range keys ***REMOVED***
		v := ctx.Value(key)
		if v != nil ***REMOVED***
			fields[fmt.Sprint(key)] = v
		***REMOVED***
	***REMOVED***

	return logger.WithFields(fields)
***REMOVED***
