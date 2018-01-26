// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package grpc_prometheus

import (
	"time"

	"google.golang.org/grpc/codes"

	prom "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

type grpcType string

const (
	Unary        grpcType = "unary"
	ClientStream grpcType = "client_stream"
	ServerStream grpcType = "server_stream"
	BidiStream   grpcType = "bidi_stream"
)

var (
	serverStartedCounter = prom.NewCounterVec(
		prom.CounterOpts***REMOVED***
			Namespace: "grpc",
			Subsystem: "server",
			Name:      "started_total",
			Help:      "Total number of RPCs started on the server.",
		***REMOVED***, []string***REMOVED***"grpc_type", "grpc_service", "grpc_method"***REMOVED***)

	serverHandledCounter = prom.NewCounterVec(
		prom.CounterOpts***REMOVED***
			Namespace: "grpc",
			Subsystem: "server",
			Name:      "handled_total",
			Help:      "Total number of RPCs completed on the server, regardless of success or failure.",
		***REMOVED***, []string***REMOVED***"grpc_type", "grpc_service", "grpc_method", "grpc_code"***REMOVED***)

	serverStreamMsgReceived = prom.NewCounterVec(
		prom.CounterOpts***REMOVED***
			Namespace: "grpc",
			Subsystem: "server",
			Name:      "msg_received_total",
			Help:      "Total number of RPC stream messages received on the server.",
		***REMOVED***, []string***REMOVED***"grpc_type", "grpc_service", "grpc_method"***REMOVED***)

	serverStreamMsgSent = prom.NewCounterVec(
		prom.CounterOpts***REMOVED***
			Namespace: "grpc",
			Subsystem: "server",
			Name:      "msg_sent_total",
			Help:      "Total number of gRPC stream messages sent by the server.",
		***REMOVED***, []string***REMOVED***"grpc_type", "grpc_service", "grpc_method"***REMOVED***)

	serverHandledHistogramEnabled = false
	serverHandledHistogramOpts    = prom.HistogramOpts***REMOVED***
		Namespace: "grpc",
		Subsystem: "server",
		Name:      "handling_seconds",
		Help:      "Histogram of response latency (seconds) of gRPC that had been application-level handled by the server.",
		Buckets:   prom.DefBuckets,
	***REMOVED***
	serverHandledHistogram *prom.HistogramVec
)

func init() ***REMOVED***
	prom.MustRegister(serverStartedCounter)
	prom.MustRegister(serverHandledCounter)
	prom.MustRegister(serverStreamMsgReceived)
	prom.MustRegister(serverStreamMsgSent)
***REMOVED***

type HistogramOption func(*prom.HistogramOpts)

// WithHistogramBuckets allows you to specify custom bucket ranges for histograms if EnableHandlingTimeHistogram is on.
func WithHistogramBuckets(buckets []float64) HistogramOption ***REMOVED***
	return func(o *prom.HistogramOpts) ***REMOVED*** o.Buckets = buckets ***REMOVED***
***REMOVED***

// EnableHandlingTimeHistogram turns on recording of handling time of RPCs for server-side interceptors.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func EnableHandlingTimeHistogram(opts ...HistogramOption) ***REMOVED***
	for _, o := range opts ***REMOVED***
		o(&serverHandledHistogramOpts)
	***REMOVED***
	if !serverHandledHistogramEnabled ***REMOVED***
		serverHandledHistogram = prom.NewHistogramVec(
			serverHandledHistogramOpts,
			[]string***REMOVED***"grpc_type", "grpc_service", "grpc_method"***REMOVED***,
		)
		prom.Register(serverHandledHistogram)
	***REMOVED***
	serverHandledHistogramEnabled = true
***REMOVED***

type serverReporter struct ***REMOVED***
	rpcType     grpcType
	serviceName string
	methodName  string
	startTime   time.Time
***REMOVED***

func newServerReporter(rpcType grpcType, fullMethod string) *serverReporter ***REMOVED***
	r := &serverReporter***REMOVED***rpcType: rpcType***REMOVED***
	if serverHandledHistogramEnabled ***REMOVED***
		r.startTime = time.Now()
	***REMOVED***
	r.serviceName, r.methodName = splitMethodName(fullMethod)
	serverStartedCounter.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName).Inc()
	return r
***REMOVED***

func (r *serverReporter) ReceivedMessage() ***REMOVED***
	serverStreamMsgReceived.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName).Inc()
***REMOVED***

func (r *serverReporter) SentMessage() ***REMOVED***
	serverStreamMsgSent.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName).Inc()
***REMOVED***

func (r *serverReporter) Handled(code codes.Code) ***REMOVED***
	serverHandledCounter.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, code.String()).Inc()
	if serverHandledHistogramEnabled ***REMOVED***
		serverHandledHistogram.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName).Observe(time.Since(r.startTime).Seconds())
	***REMOVED***
***REMOVED***

// preRegisterMethod is invoked on Register of a Server, allowing all gRPC services labels to be pre-populated.
func preRegisterMethod(serviceName string, mInfo *grpc.MethodInfo) ***REMOVED***
	methodName := mInfo.Name
	methodType := string(typeFromMethodInfo(mInfo))
	// These are just references (no increments), as just referencing will create the labels but not set values.
	serverStartedCounter.GetMetricWithLabelValues(methodType, serviceName, methodName)
	serverStreamMsgReceived.GetMetricWithLabelValues(methodType, serviceName, methodName)
	serverStreamMsgSent.GetMetricWithLabelValues(methodType, serviceName, methodName)
	if serverHandledHistogramEnabled ***REMOVED***
		serverHandledHistogram.GetMetricWithLabelValues(methodType, serviceName, methodName)
	***REMOVED***
	for _, code := range allCodes ***REMOVED***
		serverHandledCounter.GetMetricWithLabelValues(methodType, serviceName, methodName, code.String())
	***REMOVED***
***REMOVED***

func typeFromMethodInfo(mInfo *grpc.MethodInfo) grpcType ***REMOVED***
	if mInfo.IsClientStream == false && mInfo.IsServerStream == false ***REMOVED***
		return Unary
	***REMOVED***
	if mInfo.IsClientStream == true && mInfo.IsServerStream == false ***REMOVED***
		return ClientStream
	***REMOVED***
	if mInfo.IsClientStream == false && mInfo.IsServerStream == true ***REMOVED***
		return ServerStream
	***REMOVED***
	return BidiStream
***REMOVED***
