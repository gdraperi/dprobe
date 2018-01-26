package srslog

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"os"
)

// This interface allows us to work with both local and network connections,
// and enables Solaris support (see syslog_unix.go).
type serverConn interface ***REMOVED***
	writeString(framer Framer, formatter Formatter, p Priority, hostname, tag, s string) error
	close() error
***REMOVED***

// New establishes a new connection to the system log daemon.  Each
// write to the returned Writer sends a log message with the given
// priority and prefix.
func New(priority Priority, tag string) (w *Writer, err error) ***REMOVED***
	return Dial("", "", priority, tag)
***REMOVED***

// Dial establishes a connection to a log daemon by connecting to
// address raddr on the specified network.  Each write to the returned
// Writer sends a log message with the given facility, severity and
// tag.
// If network is empty, Dial will connect to the local syslog server.
func Dial(network, raddr string, priority Priority, tag string) (*Writer, error) ***REMOVED***
	return DialWithTLSConfig(network, raddr, priority, tag, nil)
***REMOVED***

// DialWithTLSCertPath establishes a secure connection to a log daemon by connecting to
// address raddr on the specified network. It uses certPath to load TLS certificates and configure
// the secure connection.
func DialWithTLSCertPath(network, raddr string, priority Priority, tag, certPath string) (*Writer, error) ***REMOVED***
	serverCert, err := ioutil.ReadFile(certPath)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return DialWithTLSCert(network, raddr, priority, tag, serverCert)
***REMOVED***

// DialWIthTLSCert establishes a secure connection to a log daemon by connecting to
// address raddr on the specified network. It uses serverCert to load a TLS certificate
// and configure the secure connection.
func DialWithTLSCert(network, raddr string, priority Priority, tag string, serverCert []byte) (*Writer, error) ***REMOVED***
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(serverCert)
	config := tls.Config***REMOVED***
		RootCAs: pool,
	***REMOVED***

	return DialWithTLSConfig(network, raddr, priority, tag, &config)
***REMOVED***

// DialWithTLSConfig establishes a secure connection to a log daemon by connecting to
// address raddr on the specified network. It uses tlsConfig to configure the secure connection.
func DialWithTLSConfig(network, raddr string, priority Priority, tag string, tlsConfig *tls.Config) (*Writer, error) ***REMOVED***
	if err := validatePriority(priority); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if tag == "" ***REMOVED***
		tag = os.Args[0]
	***REMOVED***
	hostname, _ := os.Hostname()

	w := &Writer***REMOVED***
		priority:  priority,
		tag:       tag,
		hostname:  hostname,
		network:   network,
		raddr:     raddr,
		tlsConfig: tlsConfig,
	***REMOVED***

	_, err := w.connect()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return w, err
***REMOVED***

// NewLogger creates a log.Logger whose output is written to
// the system log service with the specified priority. The logFlag
// argument is the flag set passed through to log.New to create
// the Logger.
func NewLogger(p Priority, logFlag int) (*log.Logger, error) ***REMOVED***
	s, err := New(p, "")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return log.New(s, "", logFlag), nil
***REMOVED***
