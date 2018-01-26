package transport

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPTransport(t *testing.T) ***REMOVED***
	var r io.Reader
	roundTripper := &http.Transport***REMOVED******REMOVED***
	newTransport := NewHTTPTransport(roundTripper, "http", "0.0.0.0")
	request, err := newTransport.NewRequest("", r)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	assert.Equal(t, "POST", request.Method)
***REMOVED***
