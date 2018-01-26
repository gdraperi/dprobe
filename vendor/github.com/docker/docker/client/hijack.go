package client

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/sockets"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// tlsClientCon holds tls information and a dialed connection.
type tlsClientCon struct ***REMOVED***
	*tls.Conn
	rawConn net.Conn
***REMOVED***

func (c *tlsClientCon) CloseWrite() error ***REMOVED***
	// Go standard tls.Conn doesn't provide the CloseWrite() method so we do it
	// on its underlying connection.
	if conn, ok := c.rawConn.(types.CloseWriter); ok ***REMOVED***
		return conn.CloseWrite()
	***REMOVED***
	return nil
***REMOVED***

// postHijacked sends a POST request and hijacks the connection.
func (cli *Client) postHijacked(ctx context.Context, path string, query url.Values, body interface***REMOVED******REMOVED***, headers map[string][]string) (types.HijackedResponse, error) ***REMOVED***
	bodyEncoded, err := encodeData(body)
	if err != nil ***REMOVED***
		return types.HijackedResponse***REMOVED******REMOVED***, err
	***REMOVED***

	apiPath := cli.getAPIPath(path, query)
	req, err := http.NewRequest("POST", apiPath, bodyEncoded)
	if err != nil ***REMOVED***
		return types.HijackedResponse***REMOVED******REMOVED***, err
	***REMOVED***
	req = cli.addHeaders(req, headers)

	conn, err := cli.setupHijackConn(req, "tcp")
	if err != nil ***REMOVED***
		return types.HijackedResponse***REMOVED******REMOVED***, err
	***REMOVED***

	return types.HijackedResponse***REMOVED***Conn: conn, Reader: bufio.NewReader(conn)***REMOVED***, err
***REMOVED***

func tlsDial(network, addr string, config *tls.Config) (net.Conn, error) ***REMOVED***
	return tlsDialWithDialer(new(net.Dialer), network, addr, config)
***REMOVED***

// We need to copy Go's implementation of tls.Dial (pkg/cryptor/tls/tls.go) in
// order to return our custom tlsClientCon struct which holds both the tls.Conn
// object _and_ its underlying raw connection. The rationale for this is that
// we need to be able to close the write end of the connection when attaching,
// which tls.Conn does not provide.
func tlsDialWithDialer(dialer *net.Dialer, network, addr string, config *tls.Config) (net.Conn, error) ***REMOVED***
	// We want the Timeout and Deadline values from dialer to cover the
	// whole process: TCP connection and TLS handshake. This means that we
	// also need to start our own timers now.
	timeout := dialer.Timeout

	if !dialer.Deadline.IsZero() ***REMOVED***
		deadlineTimeout := time.Until(dialer.Deadline)
		if timeout == 0 || deadlineTimeout < timeout ***REMOVED***
			timeout = deadlineTimeout
		***REMOVED***
	***REMOVED***

	var errChannel chan error

	if timeout != 0 ***REMOVED***
		errChannel = make(chan error, 2)
		time.AfterFunc(timeout, func() ***REMOVED***
			errChannel <- errors.New("")
		***REMOVED***)
	***REMOVED***

	proxyDialer, err := sockets.DialerFromEnvironment(dialer)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	rawConn, err := proxyDialer.Dial(network, addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// When we set up a TCP connection for hijack, there could be long periods
	// of inactivity (a long running command with no output) that in certain
	// network setups may cause ECONNTIMEOUT, leaving the client in an unknown
	// state. Setting TCP KeepAlive on the socket connection will prohibit
	// ECONNTIMEOUT unless the socket connection truly is broken
	if tcpConn, ok := rawConn.(*net.TCPConn); ok ***REMOVED***
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)
	***REMOVED***

	colonPos := strings.LastIndex(addr, ":")
	if colonPos == -1 ***REMOVED***
		colonPos = len(addr)
	***REMOVED***
	hostname := addr[:colonPos]

	// If no ServerName is set, infer the ServerName
	// from the hostname we're connecting to.
	if config.ServerName == "" ***REMOVED***
		// Make a copy to avoid polluting argument or default.
		config = tlsConfigClone(config)
		config.ServerName = hostname
	***REMOVED***

	conn := tls.Client(rawConn, config)

	if timeout == 0 ***REMOVED***
		err = conn.Handshake()
	***REMOVED*** else ***REMOVED***
		go func() ***REMOVED***
			errChannel <- conn.Handshake()
		***REMOVED***()

		err = <-errChannel
	***REMOVED***

	if err != nil ***REMOVED***
		rawConn.Close()
		return nil, err
	***REMOVED***

	// This is Docker difference with standard's crypto/tls package: returned a
	// wrapper which holds both the TLS and raw connections.
	return &tlsClientCon***REMOVED***conn, rawConn***REMOVED***, nil
***REMOVED***

func dial(proto, addr string, tlsConfig *tls.Config) (net.Conn, error) ***REMOVED***
	if tlsConfig != nil && proto != "unix" && proto != "npipe" ***REMOVED***
		// Notice this isn't Go standard's tls.Dial function
		return tlsDial(proto, addr, tlsConfig)
	***REMOVED***
	if proto == "npipe" ***REMOVED***
		return sockets.DialPipe(addr, 32*time.Second)
	***REMOVED***
	return net.Dial(proto, addr)
***REMOVED***

func (cli *Client) setupHijackConn(req *http.Request, proto string) (net.Conn, error) ***REMOVED***
	req.Host = cli.addr
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", proto)

	conn, err := dial(cli.proto, cli.addr, resolveTLSConfig(cli.client.Transport))
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "cannot connect to the Docker daemon. Is 'docker daemon' running on this host?")
	***REMOVED***

	// When we set up a TCP connection for hijack, there could be long periods
	// of inactivity (a long running command with no output) that in certain
	// network setups may cause ECONNTIMEOUT, leaving the client in an unknown
	// state. Setting TCP KeepAlive on the socket connection will prohibit
	// ECONNTIMEOUT unless the socket connection truly is broken
	if tcpConn, ok := conn.(*net.TCPConn); ok ***REMOVED***
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)
	***REMOVED***

	clientconn := httputil.NewClientConn(conn, nil)
	defer clientconn.Close()

	// Server hijacks the connection, error 'connection closed' expected
	resp, err := clientconn.Do(req)
	if err != httputil.ErrPersistEOF ***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if resp.StatusCode != http.StatusSwitchingProtocols ***REMOVED***
			resp.Body.Close()
			return nil, fmt.Errorf("unable to upgrade to %s, received %d", proto, resp.StatusCode)
		***REMOVED***
	***REMOVED***

	c, br := clientconn.Hijack()
	if br.Buffered() > 0 ***REMOVED***
		// If there is buffered content, wrap the connection
		c = &hijackedConn***REMOVED***c, br***REMOVED***
	***REMOVED*** else ***REMOVED***
		br.Reset(nil)
	***REMOVED***

	return c, nil
***REMOVED***

type hijackedConn struct ***REMOVED***
	net.Conn
	r *bufio.Reader
***REMOVED***

func (c *hijackedConn) Read(b []byte) (int, error) ***REMOVED***
	return c.r.Read(b)
***REMOVED***
