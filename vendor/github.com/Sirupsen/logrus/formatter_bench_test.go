package logrus

import (
	"fmt"
	"testing"
	"time"
)

// smallFields is a small size data set for benchmarking
var smallFields = Fields***REMOVED***
	"foo":   "bar",
	"baz":   "qux",
	"one":   "two",
	"three": "four",
***REMOVED***

// largeFields is a large size data set for benchmarking
var largeFields = Fields***REMOVED***
	"foo":       "bar",
	"baz":       "qux",
	"one":       "two",
	"three":     "four",
	"five":      "six",
	"seven":     "eight",
	"nine":      "ten",
	"eleven":    "twelve",
	"thirteen":  "fourteen",
	"fifteen":   "sixteen",
	"seventeen": "eighteen",
	"nineteen":  "twenty",
	"a":         "b",
	"c":         "d",
	"e":         "f",
	"g":         "h",
	"i":         "j",
	"k":         "l",
	"m":         "n",
	"o":         "p",
	"q":         "r",
	"s":         "t",
	"u":         "v",
	"w":         "x",
	"y":         "z",
	"this":      "will",
	"make":      "thirty",
	"entries":   "yeah",
***REMOVED***

var errorFields = Fields***REMOVED***
	"foo": fmt.Errorf("bar"),
	"baz": fmt.Errorf("qux"),
***REMOVED***

func BenchmarkErrorTextFormatter(b *testing.B) ***REMOVED***
	doBenchmark(b, &TextFormatter***REMOVED***DisableColors: true***REMOVED***, errorFields)
***REMOVED***

func BenchmarkSmallTextFormatter(b *testing.B) ***REMOVED***
	doBenchmark(b, &TextFormatter***REMOVED***DisableColors: true***REMOVED***, smallFields)
***REMOVED***

func BenchmarkLargeTextFormatter(b *testing.B) ***REMOVED***
	doBenchmark(b, &TextFormatter***REMOVED***DisableColors: true***REMOVED***, largeFields)
***REMOVED***

func BenchmarkSmallColoredTextFormatter(b *testing.B) ***REMOVED***
	doBenchmark(b, &TextFormatter***REMOVED***ForceColors: true***REMOVED***, smallFields)
***REMOVED***

func BenchmarkLargeColoredTextFormatter(b *testing.B) ***REMOVED***
	doBenchmark(b, &TextFormatter***REMOVED***ForceColors: true***REMOVED***, largeFields)
***REMOVED***

func BenchmarkSmallJSONFormatter(b *testing.B) ***REMOVED***
	doBenchmark(b, &JSONFormatter***REMOVED******REMOVED***, smallFields)
***REMOVED***

func BenchmarkLargeJSONFormatter(b *testing.B) ***REMOVED***
	doBenchmark(b, &JSONFormatter***REMOVED******REMOVED***, largeFields)
***REMOVED***

func doBenchmark(b *testing.B, formatter Formatter, fields Fields) ***REMOVED***
	logger := New()

	entry := &Entry***REMOVED***
		Time:    time.Time***REMOVED******REMOVED***,
		Level:   InfoLevel,
		Message: "message",
		Data:    fields,
		Logger:  logger,
	***REMOVED***
	var d []byte
	var err error
	for i := 0; i < b.N; i++ ***REMOVED***
		d, err = formatter.Format(entry)
		if err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
		b.SetBytes(int64(len(d)))
	***REMOVED***
***REMOVED***
