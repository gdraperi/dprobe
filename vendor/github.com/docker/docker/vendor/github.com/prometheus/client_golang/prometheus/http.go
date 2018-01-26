// Copyright 2014 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheus

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var instLabels = []string***REMOVED***"method", "code"***REMOVED***

type nower interface ***REMOVED***
	Now() time.Time
***REMOVED***

type nowFunc func() time.Time

func (n nowFunc) Now() time.Time ***REMOVED***
	return n()
***REMOVED***

var now nower = nowFunc(func() time.Time ***REMOVED***
	return time.Now()
***REMOVED***)

func nowSeries(t ...time.Time) nower ***REMOVED***
	return nowFunc(func() time.Time ***REMOVED***
		defer func() ***REMOVED***
			t = t[1:]
		***REMOVED***()

		return t[0]
	***REMOVED***)
***REMOVED***

// InstrumentHandler wraps the given HTTP handler for instrumentation. It
// registers four metric collectors (if not already done) and reports HTTP
// metrics to the (newly or already) registered collectors: http_requests_total
// (CounterVec), http_request_duration_microseconds (Summary),
// http_request_size_bytes (Summary), http_response_size_bytes (Summary). Each
// has a constant label named "handler" with the provided handlerName as
// value. http_requests_total is a metric vector partitioned by HTTP method
// (label name "method") and HTTP status code (label name "code").
//
// Note that InstrumentHandler has several issues:
//
// - It uses Summaries rather than Histograms. Summaries are not useful if
// aggregation across multiple instances is required.
//
// - It uses microseconds as unit, which is deprecated and should be replaced by
// seconds.
//
// - The size of the request is calculated in a separate goroutine. Since this
// calculator requires access to the request header, it creates a race with
// any writes to the header performed during request handling.
// httputil.ReverseProxy is a prominent example for a handler
// performing such writes.
//
// Upcoming versions of this package will provide ways of instrumenting HTTP
// handlers that are more flexible and have fewer issues. Consider this function
// DEPRECATED and prefer direct instrumentation in the meantime.
func InstrumentHandler(handlerName string, handler http.Handler) http.HandlerFunc ***REMOVED***
	return InstrumentHandlerFunc(handlerName, handler.ServeHTTP)
***REMOVED***

// InstrumentHandlerFunc wraps the given function for instrumentation. It
// otherwise works in the same way as InstrumentHandler (and shares the same
// issues).
func InstrumentHandlerFunc(handlerName string, handlerFunc func(http.ResponseWriter, *http.Request)) http.HandlerFunc ***REMOVED***
	return InstrumentHandlerFuncWithOpts(
		SummaryOpts***REMOVED***
			Subsystem:   "http",
			ConstLabels: Labels***REMOVED***"handler": handlerName***REMOVED***,
		***REMOVED***,
		handlerFunc,
	)
***REMOVED***

// InstrumentHandlerWithOpts works like InstrumentHandler (and shares the same
// issues) but provides more flexibility (at the cost of a more complex call
// syntax). As InstrumentHandler, this function registers four metric
// collectors, but it uses the provided SummaryOpts to create them. However, the
// fields "Name" and "Help" in the SummaryOpts are ignored. "Name" is replaced
// by "requests_total", "request_duration_microseconds", "request_size_bytes",
// and "response_size_bytes", respectively. "Help" is replaced by an appropriate
// help string. The names of the variable labels of the http_requests_total
// CounterVec are "method" (get, post, etc.), and "code" (HTTP status code).
//
// If InstrumentHandlerWithOpts is called as follows, it mimics exactly the
// behavior of InstrumentHandler:
//
//     prometheus.InstrumentHandlerWithOpts(
//         prometheus.SummaryOpts***REMOVED***
//              Subsystem:   "http",
//              ConstLabels: prometheus.Labels***REMOVED***"handler": handlerName***REMOVED***,
//     ***REMOVED***,
//         handler,
//     )
//
// Technical detail: "requests_total" is a CounterVec, not a SummaryVec, so it
// cannot use SummaryOpts. Instead, a CounterOpts struct is created internally,
// and all its fields are set to the equally named fields in the provided
// SummaryOpts.
func InstrumentHandlerWithOpts(opts SummaryOpts, handler http.Handler) http.HandlerFunc ***REMOVED***
	return InstrumentHandlerFuncWithOpts(opts, handler.ServeHTTP)
***REMOVED***

// InstrumentHandlerFuncWithOpts works like InstrumentHandlerFunc (and shares
// the same issues) but provides more flexibility (at the cost of a more complex
// call syntax). See InstrumentHandlerWithOpts for details how the provided
// SummaryOpts are used.
func InstrumentHandlerFuncWithOpts(opts SummaryOpts, handlerFunc func(http.ResponseWriter, *http.Request)) http.HandlerFunc ***REMOVED***
	reqCnt := NewCounterVec(
		CounterOpts***REMOVED***
			Namespace:   opts.Namespace,
			Subsystem:   opts.Subsystem,
			Name:        "requests_total",
			Help:        "Total number of HTTP requests made.",
			ConstLabels: opts.ConstLabels,
		***REMOVED***,
		instLabels,
	)

	opts.Name = "request_duration_microseconds"
	opts.Help = "The HTTP request latencies in microseconds."
	reqDur := NewSummary(opts)

	opts.Name = "request_size_bytes"
	opts.Help = "The HTTP request sizes in bytes."
	reqSz := NewSummary(opts)

	opts.Name = "response_size_bytes"
	opts.Help = "The HTTP response sizes in bytes."
	resSz := NewSummary(opts)

	regReqCnt := MustRegisterOrGet(reqCnt).(*CounterVec)
	regReqDur := MustRegisterOrGet(reqDur).(Summary)
	regReqSz := MustRegisterOrGet(reqSz).(Summary)
	regResSz := MustRegisterOrGet(resSz).(Summary)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		now := time.Now()

		delegate := &responseWriterDelegator***REMOVED***ResponseWriter: w***REMOVED***
		out := make(chan int)
		urlLen := 0
		if r.URL != nil ***REMOVED***
			urlLen = len(r.URL.String())
		***REMOVED***
		go computeApproximateRequestSize(r, out, urlLen)

		_, cn := w.(http.CloseNotifier)
		_, fl := w.(http.Flusher)
		_, hj := w.(http.Hijacker)
		_, rf := w.(io.ReaderFrom)
		var rw http.ResponseWriter
		if cn && fl && hj && rf ***REMOVED***
			rw = &fancyResponseWriterDelegator***REMOVED***delegate***REMOVED***
		***REMOVED*** else ***REMOVED***
			rw = delegate
		***REMOVED***
		handlerFunc(rw, r)

		elapsed := float64(time.Since(now)) / float64(time.Microsecond)

		method := sanitizeMethod(r.Method)
		code := sanitizeCode(delegate.status)
		regReqCnt.WithLabelValues(method, code).Inc()
		regReqDur.Observe(elapsed)
		regResSz.Observe(float64(delegate.written))
		regReqSz.Observe(float64(<-out))
	***REMOVED***)
***REMOVED***

func computeApproximateRequestSize(r *http.Request, out chan int, s int) ***REMOVED***
	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header ***REMOVED***
		s += len(name)
		for _, value := range values ***REMOVED***
			s += len(value)
		***REMOVED***
	***REMOVED***
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 ***REMOVED***
		s += int(r.ContentLength)
	***REMOVED***
	out <- s
***REMOVED***

type responseWriterDelegator struct ***REMOVED***
	http.ResponseWriter

	handler, method string
	status          int
	written         int64
	wroteHeader     bool
***REMOVED***

func (r *responseWriterDelegator) WriteHeader(code int) ***REMOVED***
	r.status = code
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(code)
***REMOVED***

func (r *responseWriterDelegator) Write(b []byte) (int, error) ***REMOVED***
	if !r.wroteHeader ***REMOVED***
		r.WriteHeader(http.StatusOK)
	***REMOVED***
	n, err := r.ResponseWriter.Write(b)
	r.written += int64(n)
	return n, err
***REMOVED***

type fancyResponseWriterDelegator struct ***REMOVED***
	*responseWriterDelegator
***REMOVED***

func (f *fancyResponseWriterDelegator) CloseNotify() <-chan bool ***REMOVED***
	return f.ResponseWriter.(http.CloseNotifier).CloseNotify()
***REMOVED***

func (f *fancyResponseWriterDelegator) Flush() ***REMOVED***
	f.ResponseWriter.(http.Flusher).Flush()
***REMOVED***

func (f *fancyResponseWriterDelegator) Hijack() (net.Conn, *bufio.ReadWriter, error) ***REMOVED***
	return f.ResponseWriter.(http.Hijacker).Hijack()
***REMOVED***

func (f *fancyResponseWriterDelegator) ReadFrom(r io.Reader) (int64, error) ***REMOVED***
	if !f.wroteHeader ***REMOVED***
		f.WriteHeader(http.StatusOK)
	***REMOVED***
	n, err := f.ResponseWriter.(io.ReaderFrom).ReadFrom(r)
	f.written += n
	return n, err
***REMOVED***

func sanitizeMethod(m string) string ***REMOVED***
	switch m ***REMOVED***
	case "GET", "get":
		return "get"
	case "PUT", "put":
		return "put"
	case "HEAD", "head":
		return "head"
	case "POST", "post":
		return "post"
	case "DELETE", "delete":
		return "delete"
	case "CONNECT", "connect":
		return "connect"
	case "OPTIONS", "options":
		return "options"
	case "NOTIFY", "notify":
		return "notify"
	default:
		return strings.ToLower(m)
	***REMOVED***
***REMOVED***

func sanitizeCode(s int) string ***REMOVED***
	switch s ***REMOVED***
	case 100:
		return "100"
	case 101:
		return "101"

	case 200:
		return "200"
	case 201:
		return "201"
	case 202:
		return "202"
	case 203:
		return "203"
	case 204:
		return "204"
	case 205:
		return "205"
	case 206:
		return "206"

	case 300:
		return "300"
	case 301:
		return "301"
	case 302:
		return "302"
	case 304:
		return "304"
	case 305:
		return "305"
	case 307:
		return "307"

	case 400:
		return "400"
	case 401:
		return "401"
	case 402:
		return "402"
	case 403:
		return "403"
	case 404:
		return "404"
	case 405:
		return "405"
	case 406:
		return "406"
	case 407:
		return "407"
	case 408:
		return "408"
	case 409:
		return "409"
	case 410:
		return "410"
	case 411:
		return "411"
	case 412:
		return "412"
	case 413:
		return "413"
	case 414:
		return "414"
	case 415:
		return "415"
	case 416:
		return "416"
	case 417:
		return "417"
	case 418:
		return "418"

	case 500:
		return "500"
	case 501:
		return "501"
	case 502:
		return "502"
	case 503:
		return "503"
	case 504:
		return "504"
	case 505:
		return "505"

	case 428:
		return "428"
	case 429:
		return "429"
	case 431:
		return "431"
	case 511:
		return "511"

	default:
		return strconv.Itoa(s)
	***REMOVED***
***REMOVED***
