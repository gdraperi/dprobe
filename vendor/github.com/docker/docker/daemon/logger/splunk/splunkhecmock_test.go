package splunk

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"testing"
)

func (message *splunkMessage) EventAsString() (string, error) ***REMOVED***
	if val, ok := message.Event.(string); ok ***REMOVED***
		return val, nil
	***REMOVED***
	return "", fmt.Errorf("Cannot cast Event %v to string", message.Event)
***REMOVED***

func (message *splunkMessage) EventAsMap() (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	if val, ok := message.Event.(map[string]interface***REMOVED******REMOVED***); ok ***REMOVED***
		return val, nil
	***REMOVED***
	return nil, fmt.Errorf("Cannot cast Event %v to map", message.Event)
***REMOVED***

type HTTPEventCollectorMock struct ***REMOVED***
	tcpAddr     *net.TCPAddr
	tcpListener *net.TCPListener

	mu                  sync.Mutex
	token               string
	simulateServerError bool
	blockingCtx         context.Context

	test *testing.T

	connectionVerified bool
	gzipEnabled        *bool
	messages           []*splunkMessage
	numOfRequests      int
***REMOVED***

func NewHTTPEventCollectorMock(t *testing.T) *HTTPEventCollectorMock ***REMOVED***
	tcpAddr := &net.TCPAddr***REMOVED***IP: []byte***REMOVED***127, 0, 0, 1***REMOVED***, Port: 0, Zone: ""***REMOVED***
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	return &HTTPEventCollectorMock***REMOVED***
		tcpAddr:             tcpAddr,
		tcpListener:         tcpListener,
		token:               "4642492F-D8BD-47F1-A005-0C08AE4657DF",
		simulateServerError: false,
		test:                t,
		connectionVerified:  false***REMOVED***
***REMOVED***

func (hec *HTTPEventCollectorMock) simulateErr(b bool) ***REMOVED***
	hec.mu.Lock()
	hec.simulateServerError = b
	hec.mu.Unlock()
***REMOVED***

func (hec *HTTPEventCollectorMock) withBlock(ctx context.Context) ***REMOVED***
	hec.mu.Lock()
	hec.blockingCtx = ctx
	hec.mu.Unlock()
***REMOVED***

func (hec *HTTPEventCollectorMock) URL() string ***REMOVED***
	return "http://" + hec.tcpListener.Addr().String()
***REMOVED***

func (hec *HTTPEventCollectorMock) Serve() error ***REMOVED***
	return http.Serve(hec.tcpListener, hec)
***REMOVED***

func (hec *HTTPEventCollectorMock) Close() error ***REMOVED***
	return hec.tcpListener.Close()
***REMOVED***

func (hec *HTTPEventCollectorMock) ServeHTTP(writer http.ResponseWriter, request *http.Request) ***REMOVED***
	var err error

	hec.numOfRequests++

	hec.mu.Lock()
	simErr := hec.simulateServerError
	ctx := hec.blockingCtx
	hec.mu.Unlock()

	if ctx != nil ***REMOVED***
		<-hec.blockingCtx.Done()
	***REMOVED***

	if simErr ***REMOVED***
		if request.Body != nil ***REMOVED***
			defer request.Body.Close()
		***REMOVED***
		writer.WriteHeader(http.StatusInternalServerError)
		return
	***REMOVED***

	switch request.Method ***REMOVED***
	case http.MethodOptions:
		// Verify that options method is getting called only once
		if hec.connectionVerified ***REMOVED***
			hec.test.Errorf("Connection should not be verified more than once. Got second request with %s method.", request.Method)
		***REMOVED***
		hec.connectionVerified = true
		writer.WriteHeader(http.StatusOK)
	case http.MethodPost:
		// Always verify that Driver is using correct path to HEC
		if request.URL.String() != "/services/collector/event/1.0" ***REMOVED***
			hec.test.Errorf("Unexpected path %v", request.URL)
		***REMOVED***
		defer request.Body.Close()

		if authorization, ok := request.Header["Authorization"]; !ok || authorization[0] != ("Splunk "+hec.token) ***REMOVED***
			hec.test.Error("Authorization header is invalid.")
		***REMOVED***

		gzipEnabled := false
		if contentEncoding, ok := request.Header["Content-Encoding"]; ok && contentEncoding[0] == "gzip" ***REMOVED***
			gzipEnabled = true
		***REMOVED***

		if hec.gzipEnabled == nil ***REMOVED***
			hec.gzipEnabled = &gzipEnabled
		***REMOVED*** else if gzipEnabled != *hec.gzipEnabled ***REMOVED***
			// Nothing wrong with that, but we just know that Splunk Logging Driver does not do that
			hec.test.Error("Driver should not change Content Encoding.")
		***REMOVED***

		var gzipReader *gzip.Reader
		var reader io.Reader
		if gzipEnabled ***REMOVED***
			gzipReader, err = gzip.NewReader(request.Body)
			if err != nil ***REMOVED***
				hec.test.Fatal(err)
			***REMOVED***
			reader = gzipReader
		***REMOVED*** else ***REMOVED***
			reader = request.Body
		***REMOVED***

		// Read body
		var body []byte
		body, err = ioutil.ReadAll(reader)
		if err != nil ***REMOVED***
			hec.test.Fatal(err)
		***REMOVED***

		// Parse message
		messageStart := 0
		for i := 0; i < len(body); i++ ***REMOVED***
			if i == len(body)-1 || (body[i] == '***REMOVED***' && body[i+1] == '***REMOVED***') ***REMOVED***
				var message splunkMessage
				err = json.Unmarshal(body[messageStart:i+1], &message)
				if err != nil ***REMOVED***
					hec.test.Log(string(body[messageStart : i+1]))
					hec.test.Fatal(err)
				***REMOVED***
				hec.messages = append(hec.messages, &message)
				messageStart = i + 1
			***REMOVED***
		***REMOVED***

		if gzipEnabled ***REMOVED***
			gzipReader.Close()
		***REMOVED***

		writer.WriteHeader(http.StatusOK)
	default:
		hec.test.Errorf("Unexpected HTTP method %s", http.MethodOptions)
		writer.WriteHeader(http.StatusBadRequest)
	***REMOVED***
***REMOVED***
