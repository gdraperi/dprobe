package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httputil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
)

const logReqMsg = `DEBUG: Request %s/%s Details:
---[ REQUEST POST-SIGN ]-----------------------------
%s
-----------------------------------------------------`

const logReqErrMsg = `DEBUG ERROR: Request %s/%s:
---[ REQUEST DUMP ERROR ]-----------------------------
%s
------------------------------------------------------`

type logWriter struct ***REMOVED***
	// Logger is what we will use to log the payload of a response.
	Logger aws.Logger
	// buf stores the contents of what has been read
	buf *bytes.Buffer
***REMOVED***

func (logger *logWriter) Write(b []byte) (int, error) ***REMOVED***
	return logger.buf.Write(b)
***REMOVED***

type teeReaderCloser struct ***REMOVED***
	// io.Reader will be a tee reader that is used during logging.
	// This structure will read from a body and write the contents to a logger.
	io.Reader
	// Source is used just to close when we are done reading.
	Source io.ReadCloser
***REMOVED***

func (reader *teeReaderCloser) Close() error ***REMOVED***
	return reader.Source.Close()
***REMOVED***

func logRequest(r *request.Request) ***REMOVED***
	logBody := r.Config.LogLevel.Matches(aws.LogDebugWithHTTPBody)
	dumpedBody, err := httputil.DumpRequestOut(r.HTTPRequest, logBody)
	if err != nil ***REMOVED***
		r.Config.Logger.Log(fmt.Sprintf(logReqErrMsg, r.ClientInfo.ServiceName, r.Operation.Name, err))
		return
	***REMOVED***

	if logBody ***REMOVED***
		// Reset the request body because dumpRequest will re-wrap the r.HTTPRequest's
		// Body as a NoOpCloser and will not be reset after read by the HTTP
		// client reader.
		r.ResetBody()
	***REMOVED***

	r.Config.Logger.Log(fmt.Sprintf(logReqMsg, r.ClientInfo.ServiceName, r.Operation.Name, string(dumpedBody)))
***REMOVED***

const logRespMsg = `DEBUG: Response %s/%s Details:
---[ RESPONSE ]--------------------------------------
%s
-----------------------------------------------------`

const logRespErrMsg = `DEBUG ERROR: Response %s/%s:
---[ RESPONSE DUMP ERROR ]-----------------------------
%s
-----------------------------------------------------`

func logResponse(r *request.Request) ***REMOVED***
	lw := &logWriter***REMOVED***r.Config.Logger, bytes.NewBuffer(nil)***REMOVED***
	r.HTTPResponse.Body = &teeReaderCloser***REMOVED***
		Reader: io.TeeReader(r.HTTPResponse.Body, lw),
		Source: r.HTTPResponse.Body,
	***REMOVED***

	handlerFn := func(req *request.Request) ***REMOVED***
		body, err := httputil.DumpResponse(req.HTTPResponse, false)
		if err != nil ***REMOVED***
			lw.Logger.Log(fmt.Sprintf(logRespErrMsg, req.ClientInfo.ServiceName, req.Operation.Name, err))
			return
		***REMOVED***

		b, err := ioutil.ReadAll(lw.buf)
		if err != nil ***REMOVED***
			lw.Logger.Log(fmt.Sprintf(logRespErrMsg, req.ClientInfo.ServiceName, req.Operation.Name, err))
			return
		***REMOVED***
		lw.Logger.Log(fmt.Sprintf(logRespMsg, req.ClientInfo.ServiceName, req.Operation.Name, string(body)))
		if req.Config.LogLevel.Matches(aws.LogDebugWithHTTPBody) ***REMOVED***
			lw.Logger.Log(string(b))
		***REMOVED***
	***REMOVED***

	const handlerName = "awsdk.client.LogResponse.ResponseBody"

	r.Handlers.Unmarshal.SetBackNamed(request.NamedHandler***REMOVED***
		Name: handlerName, Fn: handlerFn,
	***REMOVED***)
	r.Handlers.UnmarshalError.SetBackNamed(request.NamedHandler***REMOVED***
		Name: handlerName, Fn: handlerFn,
	***REMOVED***)
***REMOVED***
