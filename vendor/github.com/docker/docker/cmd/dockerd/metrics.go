package main

import (
	"net"
	"net/http"

	metrics "github.com/docker/go-metrics"
	"github.com/sirupsen/logrus"
)

func startMetricsServer(addr string) error ***REMOVED***
	if err := allocateDaemonPort(addr); err != nil ***REMOVED***
		return err
	***REMOVED***
	l, err := net.Listen("tcp", addr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	mux := http.NewServeMux()
	mux.Handle("/metrics", metrics.Handler())
	go func() ***REMOVED***
		if err := http.Serve(l, mux); err != nil ***REMOVED***
			logrus.Errorf("serve metrics api: %s", err)
		***REMOVED***
	***REMOVED***()
	return nil
***REMOVED***
