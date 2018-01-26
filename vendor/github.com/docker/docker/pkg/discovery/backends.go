package discovery

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	// Backends is a global map of discovery backends indexed by their
	// associated scheme.
	backends = make(map[string]Backend)
)

// Register makes a discovery backend available by the provided scheme.
// If Register is called twice with the same scheme an error is returned.
func Register(scheme string, d Backend) error ***REMOVED***
	if _, exists := backends[scheme]; exists ***REMOVED***
		return fmt.Errorf("scheme already registered %s", scheme)
	***REMOVED***
	logrus.WithField("name", scheme).Debugf("Registering discovery service")
	backends[scheme] = d
	return nil
***REMOVED***

func parse(rawurl string) (string, string) ***REMOVED***
	parts := strings.SplitN(rawurl, "://", 2)

	// nodes:port,node2:port => nodes://node1:port,node2:port
	if len(parts) == 1 ***REMOVED***
		return "nodes", parts[0]
	***REMOVED***
	return parts[0], parts[1]
***REMOVED***

// ParseAdvertise parses the --cluster-advertise daemon config which accepts
// <ip-address>:<port> or <interface-name>:<port>
func ParseAdvertise(advertise string) (string, error) ***REMOVED***
	var (
		iface *net.Interface
		addrs []net.Addr
		err   error
	)

	addr, port, err := net.SplitHostPort(advertise)

	if err != nil ***REMOVED***
		return "", fmt.Errorf("invalid --cluster-advertise configuration: %s: %v", advertise, err)
	***REMOVED***

	ip := net.ParseIP(addr)
	// If it is a valid ip-address, use it as is
	if ip != nil ***REMOVED***
		return advertise, nil
	***REMOVED***

	// If advertise is a valid interface name, get the valid IPv4 address and use it to advertise
	ifaceName := addr
	iface, err = net.InterfaceByName(ifaceName)
	if err != nil ***REMOVED***
		return "", fmt.Errorf("invalid cluster advertise IP address or interface name (%s) : %v", advertise, err)
	***REMOVED***

	addrs, err = iface.Addrs()
	if err != nil ***REMOVED***
		return "", fmt.Errorf("unable to get advertise IP address from interface (%s) : %v", advertise, err)
	***REMOVED***

	if len(addrs) == 0 ***REMOVED***
		return "", fmt.Errorf("no available advertise IP address in interface (%s)", advertise)
	***REMOVED***

	addr = ""
	for _, a := range addrs ***REMOVED***
		ip, _, err := net.ParseCIDR(a.String())
		if err != nil ***REMOVED***
			return "", fmt.Errorf("error deriving advertise ip-address in interface (%s) : %v", advertise, err)
		***REMOVED***
		if ip.To4() == nil || ip.IsLoopback() ***REMOVED***
			continue
		***REMOVED***
		addr = ip.String()
		break
	***REMOVED***
	if addr == "" ***REMOVED***
		return "", fmt.Errorf("could not find a valid ip-address in interface %s", advertise)
	***REMOVED***

	addr = net.JoinHostPort(addr, port)
	return addr, nil
***REMOVED***

// New returns a new Discovery given a URL, heartbeat and ttl settings.
// Returns an error if the URL scheme is not supported.
func New(rawurl string, heartbeat time.Duration, ttl time.Duration, clusterOpts map[string]string) (Backend, error) ***REMOVED***
	scheme, uri := parse(rawurl)
	if backend, exists := backends[scheme]; exists ***REMOVED***
		logrus.WithFields(logrus.Fields***REMOVED***"name": scheme, "uri": uri***REMOVED***).Debugf("Initializing discovery service")
		err := backend.Initialize(uri, heartbeat, ttl, clusterOpts)
		return backend, err
	***REMOVED***

	return nil, ErrNotSupported
***REMOVED***
