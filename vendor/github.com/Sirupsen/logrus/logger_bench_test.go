package logrus

import (
	"os"
	"testing"
)

// smallFields is a small size data set for benchmarking
var loggerFields = Fields***REMOVED***
	"foo":   "bar",
	"baz":   "qux",
	"one":   "two",
	"three": "four",
***REMOVED***

func BenchmarkDummyLogger(b *testing.B) ***REMOVED***
	nullf, err := os.OpenFile("/dev/null", os.O_WRONLY, 0666)
	if err != nil ***REMOVED***
		b.Fatalf("%v", err)
	***REMOVED***
	defer nullf.Close()
	doLoggerBenchmark(b, nullf, &TextFormatter***REMOVED***DisableColors: true***REMOVED***, smallFields)
***REMOVED***

func BenchmarkDummyLoggerNoLock(b *testing.B) ***REMOVED***
	nullf, err := os.OpenFile("/dev/null", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil ***REMOVED***
		b.Fatalf("%v", err)
	***REMOVED***
	defer nullf.Close()
	doLoggerBenchmarkNoLock(b, nullf, &TextFormatter***REMOVED***DisableColors: true***REMOVED***, smallFields)
***REMOVED***

func doLoggerBenchmark(b *testing.B, out *os.File, formatter Formatter, fields Fields) ***REMOVED***
	logger := Logger***REMOVED***
		Out:       out,
		Level:     InfoLevel,
		Formatter: formatter,
	***REMOVED***
	entry := logger.WithFields(fields)
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		for pb.Next() ***REMOVED***
			entry.Info("aaa")
		***REMOVED***
	***REMOVED***)
***REMOVED***

func doLoggerBenchmarkNoLock(b *testing.B, out *os.File, formatter Formatter, fields Fields) ***REMOVED***
	logger := Logger***REMOVED***
		Out:       out,
		Level:     InfoLevel,
		Formatter: formatter,
	***REMOVED***
	logger.SetNoLock()
	entry := logger.WithFields(fields)
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		for pb.Next() ***REMOVED***
			entry.Info("aaa")
		***REMOVED***
	***REMOVED***)
***REMOVED***
