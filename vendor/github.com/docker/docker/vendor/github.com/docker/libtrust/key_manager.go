package libtrust

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"sync"
)

// ClientKeyManager manages client keys on the filesystem
type ClientKeyManager struct ***REMOVED***
	key        PrivateKey
	clientFile string
	clientDir  string

	clientLock sync.RWMutex
	clients    []PublicKey

	configLock sync.Mutex
	configs    []*tls.Config
***REMOVED***

// NewClientKeyManager loads a new manager from a set of key files
// and managed by the given private key.
func NewClientKeyManager(trustKey PrivateKey, clientFile, clientDir string) (*ClientKeyManager, error) ***REMOVED***
	m := &ClientKeyManager***REMOVED***
		key:        trustKey,
		clientFile: clientFile,
		clientDir:  clientDir,
	***REMOVED***
	if err := m.loadKeys(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// TODO Start watching file and directory

	return m, nil
***REMOVED***

func (c *ClientKeyManager) loadKeys() (err error) ***REMOVED***
	// Load authorized keys file
	var clients []PublicKey
	if c.clientFile != "" ***REMOVED***
		clients, err = LoadKeySetFile(c.clientFile)
		if err != nil ***REMOVED***
			return fmt.Errorf("unable to load authorized keys: %s", err)
		***REMOVED***
	***REMOVED***

	// Add clients from authorized keys directory
	files, err := ioutil.ReadDir(c.clientDir)
	if err != nil && !os.IsNotExist(err) ***REMOVED***
		return fmt.Errorf("unable to open authorized keys directory: %s", err)
	***REMOVED***
	for _, f := range files ***REMOVED***
		if !f.IsDir() ***REMOVED***
			publicKey, err := LoadPublicKeyFile(path.Join(c.clientDir, f.Name()))
			if err != nil ***REMOVED***
				return fmt.Errorf("unable to load authorized key file: %s", err)
			***REMOVED***
			clients = append(clients, publicKey)
		***REMOVED***
	***REMOVED***

	c.clientLock.Lock()
	c.clients = clients
	c.clientLock.Unlock()

	return nil
***REMOVED***

// RegisterTLSConfig registers a tls configuration to manager
// such that any changes to the keys may be reflected in
// the tls client CA pool
func (c *ClientKeyManager) RegisterTLSConfig(tlsConfig *tls.Config) error ***REMOVED***
	c.clientLock.RLock()
	certPool, err := GenerateCACertPool(c.key, c.clients)
	if err != nil ***REMOVED***
		return fmt.Errorf("CA pool generation error: %s", err)
	***REMOVED***
	c.clientLock.RUnlock()

	tlsConfig.ClientCAs = certPool

	c.configLock.Lock()
	c.configs = append(c.configs, tlsConfig)
	c.configLock.Unlock()

	return nil
***REMOVED***

// NewIdentityAuthTLSConfig creates a tls.Config for the server to use for
// libtrust identity authentication for the domain specified
func NewIdentityAuthTLSConfig(trustKey PrivateKey, clients *ClientKeyManager, addr string, domain string) (*tls.Config, error) ***REMOVED***
	tlsConfig := newTLSConfig()

	tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	if err := clients.RegisterTLSConfig(tlsConfig); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Generate cert
	ips, domains, err := parseAddr(addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// add domain that it expects clients to use
	domains = append(domains, domain)
	x509Cert, err := GenerateSelfSignedServerCert(trustKey, domains, ips)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("certificate generation error: %s", err)
	***REMOVED***
	tlsConfig.Certificates = []tls.Certificate***REMOVED******REMOVED***
		Certificate: [][]byte***REMOVED***x509Cert.Raw***REMOVED***,
		PrivateKey:  trustKey.CryptoPrivateKey(),
		Leaf:        x509Cert,
	***REMOVED******REMOVED***

	return tlsConfig, nil
***REMOVED***

// NewCertAuthTLSConfig creates a tls.Config for the server to use for
// certificate authentication
func NewCertAuthTLSConfig(caPath, certPath, keyPath string) (*tls.Config, error) ***REMOVED***
	tlsConfig := newTLSConfig()

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("Couldn't load X509 key pair (%s, %s): %s. Key encrypted?", certPath, keyPath, err)
	***REMOVED***
	tlsConfig.Certificates = []tls.Certificate***REMOVED***cert***REMOVED***

	// Verify client certificates against a CA?
	if caPath != "" ***REMOVED***
		certPool := x509.NewCertPool()
		file, err := ioutil.ReadFile(caPath)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("Couldn't read CA certificate: %s", err)
		***REMOVED***
		certPool.AppendCertsFromPEM(file)

		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		tlsConfig.ClientCAs = certPool
	***REMOVED***

	return tlsConfig, nil
***REMOVED***

func newTLSConfig() *tls.Config ***REMOVED***
	return &tls.Config***REMOVED***
		NextProtos: []string***REMOVED***"http/1.1"***REMOVED***,
		// Avoid fallback on insecure SSL protocols
		MinVersion: tls.VersionTLS10,
	***REMOVED***
***REMOVED***

// parseAddr parses an address into an array of IPs and domains
func parseAddr(addr string) ([]net.IP, []string, error) ***REMOVED***
	host, _, err := net.SplitHostPort(addr)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	var domains []string
	var ips []net.IP
	ip := net.ParseIP(host)
	if ip != nil ***REMOVED***
		ips = []net.IP***REMOVED***ip***REMOVED***
	***REMOVED*** else ***REMOVED***
		domains = []string***REMOVED***host***REMOVED***
	***REMOVED***
	return ips, domains, nil
***REMOVED***
