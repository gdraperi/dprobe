// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package grpc_prometheus

import (
	"time"

	"google.golang.org/grpc/codes"

	prom "github.com/prometheus/client_golang/prometheus"
)

var (
	clientStartedCounter = prom.NewCounterVec(
		prom.CounterOpts***REMOVED***
			Namespace: "grpc",
			Subsystem: "client",
			Name:      "started_total",
			Help:      "Total number of RPCs started on the client.",
		***REMOVED***, []string***REMOVED***"grpc_type", "grpc_service", "grpc_method"***REMOVED***)

	clientHandledCounter = prom.NewCounterVec(
		prom.CounterOpts***REMOVED***
			Namespace: "grpc",
			Subsystem: "client",
			Name:      "handled_total",
			Help:      "Total number of RPCs completed by the client, regardless of success or failure.",
		***REMOVED***, []string***REMOVED***"grpc_type", "grpc_service", "grpc_method", "grpc_code"***REMOVED***)

	clientStreamMsgReceived = prom.NewCounterVec(
		prom.CounterOpts***REMOVED***
			Namespace: "grpc",
			Subsystem: "client",
			Name:      "msg_received_total",
			Help:      "Total number of RPC stream messages received by the client.",
		***REMOVED***, []string***REMOVED***"grpc_type", "grpc_service", "grpc_method"***REMOVED***)

	clientStreamMsgSent = prom.NewCounterVec(
		prom.CounterOpts***REMOVED***
			Namespace: "grpc",
			Subsystem: "client",
			Name:      "msg_sent_total",
			Help:      "Total number of gRPC stream messages sent by the client.",
		***REMOVED***, []string***REMOVED***"grpc_type", "grpc_service", "grpc_method"***REMOVED***)

	clientHandledHistogramEnabled = false
	clientHandledHistogramOpts    = prom.HistogramOpts***REMOVED***
		Namespace: "grpc",
		Subsystem: "client",
		Name:      "handling_seconds",
		Help:      "Histogram of response latency (seconds) of the gRPC until it is finished by the application.",
		Buckets:   prom.DefBuckets,
	***REMOVED***
	clientHandledHistogram *prom.HistogramVec
)

func init() ***REMOVED***
	prom.MustRegister(clientStartedCounter)
	prom.MustRegister(clientHandledCounter)
	prom.MustRegister(clientStreamMsgReceived)
	prom.MustRegister(clientStreamMsgSent)
***REMOVED***

// EnableClientHandlingTimeHistogram turns on recording of handling time of RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func EnableClientHandlingTimeHistogram(opts ...HistogramOption) ***REMOVED***
	for _, o := range opts ***REMOVED***
		o(&clientHandledHistogramOpts)
	***REMOVED***
	if !clientHandledHistogramEnabled ***REMOVED***
		clientHandledHistogram = prom.NewHistogramVec(
			clientHandledHistogramOpts,
			[]string***REMOVED***"grpc_type", "grpc_service", "grpc_method"***REMOVED***,
		)
		prom.Register(clientHandledHistogram)
	***REMOVED***
	clientHandledHistogramEnabled = true
***REMOVED***

type clientReporter struct ***REMOVED***
	rpcType     grpcType
	serviceName string
	methodName  string
	startTime   time.Time
***REMOVED***

func newClientReporter(rpcType grpcType, fullMethod string) *clientReporter ***REMOVED***
	r := &clientReporter***REMOVED***rpcType: rpcType***REMOVED***
	if clientHandledHistogramEnabled ***REMOVED***
		r.startTime = time.Now()
	***REMOVED***
	r.serviceName, r.methodName = splitMethodName(fullMethod)
	clientStartedCounter.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName).Inc()
	return r
***REMOVED***

func (r *clientReporter) ReceivedMessage() ***REMOVED***
	clientStreamMsgReceived.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName).Inc()
***REMOVED***

func (r *clientReporter) SentMessage() ***REMOVED***
	clientStreamMsgSent.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName).Inc()
***REMOVED***

func (r *clientReporter) Handled(code codes.Code) ***REMOVED***
	clientHandledCounter.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, code.String()).Inc()
	if clientHandledHistogramEnabled ***REMOVED***
		clientHandledHistogram.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName).Observe(time.Since(r.startTime).Seconds())
	***REMOVED***
***REMOVED***
