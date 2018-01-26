// Package registry contains client primitives to interact with a remote Docker registry.
package registry

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/distribution/registry/client/transport"
	"github.com/docker/go-connections/sockets"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/sirupsen/logrus"
)

var (
	// ErrAlreadyExists is an error returned if an image being pushed
	// already exists on the remote side
	ErrAlreadyExists = errors.New("Image already exists")
)

func newTLSConfig(hostname string, isSecure bool) (*tls.Config, error) ***REMOVED***
	// PreferredServerCipherSuites should have no effect
	tlsConfig := tlsconfig.ServerDefault()

	tlsConfig.InsecureSkipVerify = !isSecure

	if isSecure && CertsDir != "" ***REMOVED***
		hostDir := filepath.Join(CertsDir, cleanPath(hostname))
		logrus.Debugf("hostDir: %s", hostDir)
		if err := ReadCertsDirectory(tlsConfig, hostDir); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return tlsConfig, nil
***REMOVED***

func hasFile(files []os.FileInfo, name string) bool ***REMOVED***
	for _, f := range files ***REMOVED***
		if f.Name() == name ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// ReadCertsDirectory reads the directory for TLS certificates
// including roots and certificate pairs and updates the
// provided TLS configuration.
func ReadCertsDirectory(tlsConfig *tls.Config, directory string) error ***REMOVED***
	fs, err := ioutil.ReadDir(directory)
	if err != nil && !os.IsNotExist(err) ***REMOVED***
		return err
	***REMOVED***

	for _, f := range fs ***REMOVED***
		if strings.HasSuffix(f.Name(), ".crt") ***REMOVED***
			if tlsConfig.RootCAs == nil ***REMOVED***
				systemPool, err := tlsconfig.SystemCertPool()
				if err != nil ***REMOVED***
					return fmt.Errorf("unable to get system cert pool: %v", err)
				***REMOVED***
				tlsConfig.RootCAs = systemPool
			***REMOVED***
			logrus.Debugf("crt: %s", filepath.Join(directory, f.Name()))
			data, err := ioutil.ReadFile(filepath.Join(directory, f.Name()))
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			tlsConfig.RootCAs.AppendCertsFromPEM(data)
		***REMOVED***
		if strings.HasSuffix(f.Name(), ".cert") ***REMOVED***
			certName := f.Name()
			keyName := certName[:len(certName)-5] + ".key"
			logrus.Debugf("cert: %s", filepath.Join(directory, f.Name()))
			if !hasFile(fs, keyName) ***REMOVED***
				return fmt.Errorf("missing key %s for client certificate %s. Note that CA certificates should use the extension .crt", keyName, certName)
			***REMOVED***
			cert, err := tls.LoadX509KeyPair(filepath.Join(directory, certName), filepath.Join(directory, keyName))
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
		***REMOVED***
		if strings.HasSuffix(f.Name(), ".key") ***REMOVED***
			keyName := f.Name()
			certName := keyName[:len(keyName)-4] + ".cert"
			logrus.Debugf("key: %s", filepath.Join(directory, f.Name()))
			if !hasFile(fs, certName) ***REMOVED***
				return fmt.Errorf("Missing client certificate %s for key %s", certName, keyName)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// Headers returns request modifiers with a User-Agent and metaHeaders
func Headers(userAgent string, metaHeaders http.Header) []transport.RequestModifier ***REMOVED***
	modifiers := []transport.RequestModifier***REMOVED******REMOVED***
	if userAgent != "" ***REMOVED***
		modifiers = append(modifiers, transport.NewHeaderRequestModifier(http.Header***REMOVED***
			"User-Agent": []string***REMOVED***userAgent***REMOVED***,
		***REMOVED***))
	***REMOVED***
	if metaHeaders != nil ***REMOVED***
		modifiers = append(modifiers, transport.NewHeaderRequestModifier(metaHeaders))
	***REMOVED***
	return modifiers
***REMOVED***

// HTTPClient returns an HTTP client structure which uses the given transport
// and contains the necessary headers for redirected requests
func HTTPClient(transport http.RoundTripper) *http.Client ***REMOVED***
	return &http.Client***REMOVED***
		Transport:     transport,
		CheckRedirect: addRequiredHeadersToRedirectedRequests,
	***REMOVED***
***REMOVED***

func trustedLocation(req *http.Request) bool ***REMOVED***
	var (
		trusteds = []string***REMOVED***"docker.com", "docker.io"***REMOVED***
		hostname = strings.SplitN(req.Host, ":", 2)[0]
	)
	if req.URL.Scheme != "https" ***REMOVED***
		return false
	***REMOVED***

	for _, trusted := range trusteds ***REMOVED***
		if hostname == trusted || strings.HasSuffix(hostname, "."+trusted) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// addRequiredHeadersToRedirectedRequests adds the necessary redirection headers
// for redirected requests
func addRequiredHeadersToRedirectedRequests(req *http.Request, via []*http.Request) error ***REMOVED***
	if via != nil && via[0] != nil ***REMOVED***
		if trustedLocation(req) && trustedLocation(via[0]) ***REMOVED***
			req.Header = via[0].Header
			return nil
		***REMOVED***
		for k, v := range via[0].Header ***REMOVED***
			if k != "Authorization" ***REMOVED***
				for _, vv := range v ***REMOVED***
					req.Header.Add(k, vv)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// NewTransport returns a new HTTP transport. If tlsConfig is nil, it uses the
// default TLS configuration.
func NewTransport(tlsConfig *tls.Config) *http.Transport ***REMOVED***
	if tlsConfig == nil ***REMOVED***
		tlsConfig = tlsconfig.ServerDefault()
	***REMOVED***

	direct := &net.Dialer***REMOVED***
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	***REMOVED***

	base := &http.Transport***REMOVED***
		Proxy:               http.ProxyFromEnvironment,
		Dial:                direct.Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     tlsConfig,
		// TODO(dmcgowan): Call close idle connections when complete and use keep alive
		DisableKeepAlives: true,
	***REMOVED***

	proxyDialer, err := sockets.DialerFromEnvironment(direct)
	if err == nil ***REMOVED***
		base.Dial = proxyDialer.Dial
	***REMOVED***
	return base
***REMOVED***
