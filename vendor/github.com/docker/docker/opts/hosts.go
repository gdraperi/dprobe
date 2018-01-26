package opts

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

var (
	// DefaultHTTPPort Default HTTP Port used if only the protocol is provided to -H flag e.g. dockerd -H tcp://
	// These are the IANA registered port numbers for use with Docker
	// see http://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.xhtml?search=docker
	DefaultHTTPPort = 2375 // Default HTTP Port
	// DefaultTLSHTTPPort Default HTTP Port used when TLS enabled
	DefaultTLSHTTPPort = 2376 // Default TLS encrypted HTTP Port
	// DefaultUnixSocket Path for the unix socket.
	// Docker daemon by default always listens on the default unix socket
	DefaultUnixSocket = "/var/run/docker.sock"
	// DefaultTCPHost constant defines the default host string used by docker on Windows
	DefaultTCPHost = fmt.Sprintf("tcp://%s:%d", DefaultHTTPHost, DefaultHTTPPort)
	// DefaultTLSHost constant defines the default host string used by docker for TLS sockets
	DefaultTLSHost = fmt.Sprintf("tcp://%s:%d", DefaultHTTPHost, DefaultTLSHTTPPort)
	// DefaultNamedPipe defines the default named pipe used by docker on Windows
	DefaultNamedPipe = `//./pipe/docker_engine`
)

// ValidateHost validates that the specified string is a valid host and returns it.
func ValidateHost(val string) (string, error) ***REMOVED***
	host := strings.TrimSpace(val)
	// The empty string means default and is not handled by parseDaemonHost
	if host != "" ***REMOVED***
		_, err := parseDaemonHost(host)
		if err != nil ***REMOVED***
			return val, err
		***REMOVED***
	***REMOVED***
	// Note: unlike most flag validators, we don't return the mutated value here
	//       we need to know what the user entered later (using ParseHost) to adjust for TLS
	return val, nil
***REMOVED***

// ParseHost and set defaults for a Daemon host string
func ParseHost(defaultToTLS bool, val string) (string, error) ***REMOVED***
	host := strings.TrimSpace(val)
	if host == "" ***REMOVED***
		if defaultToTLS ***REMOVED***
			host = DefaultTLSHost
		***REMOVED*** else ***REMOVED***
			host = DefaultHost
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		var err error
		host, err = parseDaemonHost(host)
		if err != nil ***REMOVED***
			return val, err
		***REMOVED***
	***REMOVED***
	return host, nil
***REMOVED***

// parseDaemonHost parses the specified address and returns an address that will be used as the host.
// Depending of the address specified, this may return one of the global Default* strings defined in hosts.go.
func parseDaemonHost(addr string) (string, error) ***REMOVED***
	addrParts := strings.SplitN(addr, "://", 2)
	if len(addrParts) == 1 && addrParts[0] != "" ***REMOVED***
		addrParts = []string***REMOVED***"tcp", addrParts[0]***REMOVED***
	***REMOVED***

	switch addrParts[0] ***REMOVED***
	case "tcp":
		return ParseTCPAddr(addrParts[1], DefaultTCPHost)
	case "unix":
		return parseSimpleProtoAddr("unix", addrParts[1], DefaultUnixSocket)
	case "npipe":
		return parseSimpleProtoAddr("npipe", addrParts[1], DefaultNamedPipe)
	case "fd":
		return addr, nil
	default:
		return "", fmt.Errorf("Invalid bind address format: %s", addr)
	***REMOVED***
***REMOVED***

// parseSimpleProtoAddr parses and validates that the specified address is a valid
// socket address for simple protocols like unix and npipe. It returns a formatted
// socket address, either using the address parsed from addr, or the contents of
// defaultAddr if addr is a blank string.
func parseSimpleProtoAddr(proto, addr, defaultAddr string) (string, error) ***REMOVED***
	addr = strings.TrimPrefix(addr, proto+"://")
	if strings.Contains(addr, "://") ***REMOVED***
		return "", fmt.Errorf("Invalid proto, expected %s: %s", proto, addr)
	***REMOVED***
	if addr == "" ***REMOVED***
		addr = defaultAddr
	***REMOVED***
	return fmt.Sprintf("%s://%s", proto, addr), nil
***REMOVED***

// ParseTCPAddr parses and validates that the specified address is a valid TCP
// address. It returns a formatted TCP address, either using the address parsed
// from tryAddr, or the contents of defaultAddr if tryAddr is a blank string.
// tryAddr is expected to have already been Trim()'d
// defaultAddr must be in the full `tcp://host:port` form
func ParseTCPAddr(tryAddr string, defaultAddr string) (string, error) ***REMOVED***
	if tryAddr == "" || tryAddr == "tcp://" ***REMOVED***
		return defaultAddr, nil
	***REMOVED***
	addr := strings.TrimPrefix(tryAddr, "tcp://")
	if strings.Contains(addr, "://") || addr == "" ***REMOVED***
		return "", fmt.Errorf("Invalid proto, expected tcp: %s", tryAddr)
	***REMOVED***

	defaultAddr = strings.TrimPrefix(defaultAddr, "tcp://")
	defaultHost, defaultPort, err := net.SplitHostPort(defaultAddr)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	// url.Parse fails for trailing colon on IPv6 brackets on Go 1.5, but
	// not 1.4. See https://github.com/golang/go/issues/12200 and
	// https://github.com/golang/go/issues/6530.
	if strings.HasSuffix(addr, "]:") ***REMOVED***
		addr += defaultPort
	***REMOVED***

	u, err := url.Parse("tcp://" + addr)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil ***REMOVED***
		// try port addition once
		host, port, err = net.SplitHostPort(net.JoinHostPort(u.Host, defaultPort))
	***REMOVED***
	if err != nil ***REMOVED***
		return "", fmt.Errorf("Invalid bind address format: %s", tryAddr)
	***REMOVED***

	if host == "" ***REMOVED***
		host = defaultHost
	***REMOVED***
	if port == "" ***REMOVED***
		port = defaultPort
	***REMOVED***
	p, err := strconv.Atoi(port)
	if err != nil && p == 0 ***REMOVED***
		return "", fmt.Errorf("Invalid bind address format: %s", tryAddr)
	***REMOVED***

	return fmt.Sprintf("tcp://%s%s", net.JoinHostPort(host, port), u.Path), nil
***REMOVED***

// ValidateExtraHost validates that the specified string is a valid extrahost and returns it.
// ExtraHost is in the form of name:ip where the ip has to be a valid ip (IPv4 or IPv6).
func ValidateExtraHost(val string) (string, error) ***REMOVED***
	// allow for IPv6 addresses in extra hosts by only splitting on first ":"
	arr := strings.SplitN(val, ":", 2)
	if len(arr) != 2 || len(arr[0]) == 0 ***REMOVED***
		return "", fmt.Errorf("bad format for add-host: %q", val)
	***REMOVED***
	if _, err := ValidateIPAddress(arr[1]); err != nil ***REMOVED***
		return "", fmt.Errorf("invalid IP address in add-host: %q", arr[1])
	***REMOVED***
	return val, nil
***REMOVED***
