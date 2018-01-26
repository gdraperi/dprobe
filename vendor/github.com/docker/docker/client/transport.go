package client

import (
	"crypto/tls"
	"net/http"
)

// resolveTLSConfig attempts to resolve the TLS configuration from the
// RoundTripper.
func resolveTLSConfig(transport http.RoundTripper) *tls.Config ***REMOVED***
	switch tr := transport.(type) ***REMOVED***
	case *http.Transport:
		return tr.TLSClientConfig
	default:
		return nil
	***REMOVED***
***REMOVED***
