package registry

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
)

type handlerFunc func(w http.ResponseWriter, r *http.Request)

// Mock represent a registry mock
type Mock struct ***REMOVED***
	server   *httptest.Server
	hostport string
	handlers map[string]handlerFunc
	mu       sync.Mutex
***REMOVED***

// RegisterHandler register the specified handler for the registry mock
func (tr *Mock) RegisterHandler(path string, h handlerFunc) ***REMOVED***
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.handlers[path] = h
***REMOVED***

// NewMock creates a registry mock
func NewMock(t testingT) (*Mock, error) ***REMOVED***
	testReg := &Mock***REMOVED***handlers: make(map[string]handlerFunc)***REMOVED***

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		url := r.URL.String()

		var matched bool
		var err error
		for re, function := range testReg.handlers ***REMOVED***
			matched, err = regexp.MatchString(re, url)
			if err != nil ***REMOVED***
				t.Fatal("Error with handler regexp")
			***REMOVED***
			if matched ***REMOVED***
				function(w, r)
				break
			***REMOVED***
		***REMOVED***

		if !matched ***REMOVED***
			t.Fatalf("Unable to match %s with regexp", url)
		***REMOVED***
	***REMOVED***))

	testReg.server = ts
	testReg.hostport = strings.Replace(ts.URL, "http://", "", 1)
	return testReg, nil
***REMOVED***

// URL returns the url of the registry
func (tr *Mock) URL() string ***REMOVED***
	return tr.hostport
***REMOVED***

// Close closes mock and releases resources
func (tr *Mock) Close() ***REMOVED***
	tr.server.Close()
***REMOVED***
