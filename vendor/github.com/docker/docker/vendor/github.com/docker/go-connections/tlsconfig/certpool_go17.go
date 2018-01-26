// +build go1.7

package tlsconfig

import (
	"crypto/x509"
	"runtime"
)

// SystemCertPool returns a copy of the system cert pool,
// returns an error if failed to load or empty pool on windows.
func SystemCertPool() (*x509.CertPool, error) ***REMOVED***
	certpool, err := x509.SystemCertPool()
	if err != nil && runtime.GOOS == "windows" ***REMOVED***
		return x509.NewCertPool(), nil
	***REMOVED***
	return certpool, err
***REMOVED***
