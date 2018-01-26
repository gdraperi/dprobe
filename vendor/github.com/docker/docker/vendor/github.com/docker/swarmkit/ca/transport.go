package ca

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"net"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/credentials"
)

var (
	// alpnProtoStr is the specified application level protocols for gRPC.
	alpnProtoStr = []string***REMOVED***"h2"***REMOVED***
)

type timeoutError struct***REMOVED******REMOVED***

func (timeoutError) Error() string   ***REMOVED*** return "mutablecredentials: Dial timed out" ***REMOVED***
func (timeoutError) Timeout() bool   ***REMOVED*** return true ***REMOVED***
func (timeoutError) Temporary() bool ***REMOVED*** return true ***REMOVED***

// MutableTLSCreds is the credentials required for authenticating a connection using TLS.
type MutableTLSCreds struct ***REMOVED***
	// Mutex for the tls config
	sync.Mutex
	// TLS configuration
	config *tls.Config
	// TLS Credentials
	tlsCreds credentials.TransportCredentials
	// store the subject for easy access
	subject pkix.Name
***REMOVED***

// Info implements the credentials.TransportCredentials interface
func (c *MutableTLSCreds) Info() credentials.ProtocolInfo ***REMOVED***
	return credentials.ProtocolInfo***REMOVED***
		SecurityProtocol: "tls",
		SecurityVersion:  "1.2",
	***REMOVED***
***REMOVED***

// Clone returns new MutableTLSCreds created from underlying *tls.Config.
// It panics if validation of underlying config fails.
func (c *MutableTLSCreds) Clone() credentials.TransportCredentials ***REMOVED***
	c.Lock()
	newCfg, err := NewMutableTLS(c.config)
	if err != nil ***REMOVED***
		panic("validation error on Clone")
	***REMOVED***
	c.Unlock()
	return newCfg
***REMOVED***

// OverrideServerName overrides *tls.Config.ServerName.
func (c *MutableTLSCreds) OverrideServerName(name string) error ***REMOVED***
	c.Lock()
	c.config.ServerName = name
	c.Unlock()
	return nil
***REMOVED***

// GetRequestMetadata implements the credentials.TransportCredentials interface
func (c *MutableTLSCreds) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) ***REMOVED***
	return nil, nil
***REMOVED***

// RequireTransportSecurity implements the credentials.TransportCredentials interface
func (c *MutableTLSCreds) RequireTransportSecurity() bool ***REMOVED***
	return true
***REMOVED***

// ClientHandshake implements the credentials.TransportCredentials interface
func (c *MutableTLSCreds) ClientHandshake(ctx context.Context, addr string, rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) ***REMOVED***
	// borrow all the code from the original TLS credentials
	c.Lock()
	if c.config.ServerName == "" ***REMOVED***
		colonPos := strings.LastIndex(addr, ":")
		if colonPos == -1 ***REMOVED***
			colonPos = len(addr)
		***REMOVED***
		c.config.ServerName = addr[:colonPos]
	***REMOVED***

	conn := tls.Client(rawConn, c.config)
	// Need to allow conn.Handshake to have access to config,
	// would create a deadlock otherwise
	c.Unlock()
	var err error
	errChannel := make(chan error, 1)
	go func() ***REMOVED***
		errChannel <- conn.Handshake()
	***REMOVED***()
	select ***REMOVED***
	case err = <-errChannel:
	case <-ctx.Done():
		err = ctx.Err()
	***REMOVED***
	if err != nil ***REMOVED***
		rawConn.Close()
		return nil, nil, err
	***REMOVED***
	return conn, nil, nil
***REMOVED***

// ServerHandshake implements the credentials.TransportCredentials interface
func (c *MutableTLSCreds) ServerHandshake(rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) ***REMOVED***
	c.Lock()
	conn := tls.Server(rawConn, c.config)
	c.Unlock()
	if err := conn.Handshake(); err != nil ***REMOVED***
		rawConn.Close()
		return nil, nil, err
	***REMOVED***

	return conn, credentials.TLSInfo***REMOVED***State: conn.ConnectionState()***REMOVED***, nil
***REMOVED***

// loadNewTLSConfig replaces the currently loaded TLS config with a new one
func (c *MutableTLSCreds) loadNewTLSConfig(newConfig *tls.Config) error ***REMOVED***
	newSubject, err := GetAndValidateCertificateSubject(newConfig.Certificates)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.Lock()
	defer c.Unlock()
	c.subject = newSubject
	c.config = newConfig

	return nil
***REMOVED***

// Config returns the current underlying TLS config.
func (c *MutableTLSCreds) Config() *tls.Config ***REMOVED***
	c.Lock()
	defer c.Unlock()

	return c.config
***REMOVED***

// Role returns the OU for the certificate encapsulated in this TransportCredentials
func (c *MutableTLSCreds) Role() string ***REMOVED***
	c.Lock()
	defer c.Unlock()

	return c.subject.OrganizationalUnit[0]
***REMOVED***

// Organization returns the O for the certificate encapsulated in this TransportCredentials
func (c *MutableTLSCreds) Organization() string ***REMOVED***
	c.Lock()
	defer c.Unlock()

	return c.subject.Organization[0]
***REMOVED***

// NodeID returns the CN for the certificate encapsulated in this TransportCredentials
func (c *MutableTLSCreds) NodeID() string ***REMOVED***
	c.Lock()
	defer c.Unlock()

	return c.subject.CommonName
***REMOVED***

// NewMutableTLS uses c to construct a mutable TransportCredentials based on TLS.
func NewMutableTLS(c *tls.Config) (*MutableTLSCreds, error) ***REMOVED***
	originalTC := credentials.NewTLS(c)

	if len(c.Certificates) < 1 ***REMOVED***
		return nil, errors.New("invalid configuration: needs at least one certificate")
	***REMOVED***

	subject, err := GetAndValidateCertificateSubject(c.Certificates)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	tc := &MutableTLSCreds***REMOVED***config: c, tlsCreds: originalTC, subject: subject***REMOVED***
	tc.config.NextProtos = alpnProtoStr

	return tc, nil
***REMOVED***

// GetAndValidateCertificateSubject is a helper method to retrieve and validate the subject
// from the x509 certificate underlying a tls.Certificate
func GetAndValidateCertificateSubject(certs []tls.Certificate) (pkix.Name, error) ***REMOVED***
	for i := range certs ***REMOVED***
		cert := &certs[i]
		x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
		if err != nil ***REMOVED***
			continue
		***REMOVED***
		if len(x509Cert.Subject.OrganizationalUnit) < 1 ***REMOVED***
			return pkix.Name***REMOVED******REMOVED***, errors.New("no OU found in certificate subject")
		***REMOVED***

		if len(x509Cert.Subject.Organization) < 1 ***REMOVED***
			return pkix.Name***REMOVED******REMOVED***, errors.New("no organization found in certificate subject")
		***REMOVED***
		if x509Cert.Subject.CommonName == "" ***REMOVED***
			return pkix.Name***REMOVED******REMOVED***, errors.New("no valid subject names found for TLS configuration")
		***REMOVED***

		return x509Cert.Subject, nil
	***REMOVED***

	return pkix.Name***REMOVED******REMOVED***, errors.New("no valid certificates found for TLS configuration")
***REMOVED***
