package logger

import (
	"errors"
	"io"
)

type client interface ***REMOVED***
	Call(string, interface***REMOVED******REMOVED***, interface***REMOVED******REMOVED***) error
	Stream(string, interface***REMOVED******REMOVED***) (io.ReadCloser, error)
***REMOVED***

type logPluginProxy struct ***REMOVED***
	client
***REMOVED***

type logPluginProxyStartLoggingRequest struct ***REMOVED***
	File string
	Info Info
***REMOVED***

type logPluginProxyStartLoggingResponse struct ***REMOVED***
	Err string
***REMOVED***

func (pp *logPluginProxy) StartLogging(file string, info Info) (err error) ***REMOVED***
	var (
		req logPluginProxyStartLoggingRequest
		ret logPluginProxyStartLoggingResponse
	)

	req.File = file
	req.Info = info
	if err = pp.Call("LogDriver.StartLogging", req, &ret); err != nil ***REMOVED***
		return
	***REMOVED***

	if ret.Err != "" ***REMOVED***
		err = errors.New(ret.Err)
	***REMOVED***

	return
***REMOVED***

type logPluginProxyStopLoggingRequest struct ***REMOVED***
	File string
***REMOVED***

type logPluginProxyStopLoggingResponse struct ***REMOVED***
	Err string
***REMOVED***

func (pp *logPluginProxy) StopLogging(file string) (err error) ***REMOVED***
	var (
		req logPluginProxyStopLoggingRequest
		ret logPluginProxyStopLoggingResponse
	)

	req.File = file
	if err = pp.Call("LogDriver.StopLogging", req, &ret); err != nil ***REMOVED***
		return
	***REMOVED***

	if ret.Err != "" ***REMOVED***
		err = errors.New(ret.Err)
	***REMOVED***

	return
***REMOVED***

type logPluginProxyCapabilitiesResponse struct ***REMOVED***
	Cap Capability
	Err string
***REMOVED***

func (pp *logPluginProxy) Capabilities() (cap Capability, err error) ***REMOVED***
	var (
		ret logPluginProxyCapabilitiesResponse
	)

	if err = pp.Call("LogDriver.Capabilities", nil, &ret); err != nil ***REMOVED***
		return
	***REMOVED***

	cap = ret.Cap

	if ret.Err != "" ***REMOVED***
		err = errors.New(ret.Err)
	***REMOVED***

	return
***REMOVED***

type logPluginProxyReadLogsRequest struct ***REMOVED***
	Info   Info
	Config ReadConfig
***REMOVED***

func (pp *logPluginProxy) ReadLogs(info Info, config ReadConfig) (stream io.ReadCloser, err error) ***REMOVED***
	var (
		req logPluginProxyReadLogsRequest
	)

	req.Info = info
	req.Config = config
	return pp.Stream("LogDriver.ReadLogs", req)
***REMOVED***
