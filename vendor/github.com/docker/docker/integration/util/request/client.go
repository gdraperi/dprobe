package request

import (
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/docker/docker/api"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/sockets"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/stretchr/testify/require"
)

// NewAPIClient returns a docker API client configured from environment variables
func NewAPIClient(t *testing.T) client.APIClient ***REMOVED***
	clt, err := client.NewEnvClient()
	require.NoError(t, err)
	return clt
***REMOVED***

// NewTLSAPIClient returns a docker API client configured with the
// provided TLS settings
func NewTLSAPIClient(t *testing.T, host, cacertPath, certPath, keyPath string) (client.APIClient, error) ***REMOVED***
	opts := tlsconfig.Options***REMOVED***
		CAFile:             cacertPath,
		CertFile:           certPath,
		KeyFile:            keyPath,
		ExclusiveRootPools: true,
	***REMOVED***
	config, err := tlsconfig.Client(opts)
	require.Nil(t, err)
	tr := &http.Transport***REMOVED***
		TLSClientConfig: config,
		DialContext: (&net.Dialer***REMOVED***
			KeepAlive: 30 * time.Second,
			Timeout:   30 * time.Second,
		***REMOVED***).DialContext,
	***REMOVED***
	proto, addr, _, err := client.ParseHost(host)
	require.Nil(t, err)

	sockets.ConfigureTransport(tr, proto, addr)

	httpClient := &http.Client***REMOVED***
		Transport:     tr,
		CheckRedirect: client.CheckRedirect,
	***REMOVED***
	verStr := api.DefaultVersion
	customHeaders := map[string]string***REMOVED******REMOVED***
	return client.NewClient(host, verStr, httpClient, customHeaders)
***REMOVED***
