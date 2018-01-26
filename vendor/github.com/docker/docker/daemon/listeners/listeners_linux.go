package listeners

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/coreos/go-systemd/activation"
	"github.com/docker/go-connections/sockets"
	"github.com/sirupsen/logrus"
)

// Init creates new listeners for the server.
// TODO: Clean up the fact that socketGroup and tlsConfig aren't always used.
func Init(proto, addr, socketGroup string, tlsConfig *tls.Config) ([]net.Listener, error) ***REMOVED***
	ls := []net.Listener***REMOVED******REMOVED***

	switch proto ***REMOVED***
	case "fd":
		fds, err := listenFD(addr, tlsConfig)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		ls = append(ls, fds...)
	case "tcp":
		l, err := sockets.NewTCPSocket(addr, tlsConfig)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		ls = append(ls, l)
	case "unix":
		gid, err := lookupGID(socketGroup)
		if err != nil ***REMOVED***
			if socketGroup != "" ***REMOVED***
				if socketGroup != defaultSocketGroup ***REMOVED***
					return nil, err
				***REMOVED***
				logrus.Warnf("could not change group %s to %s: %v", addr, defaultSocketGroup, err)
			***REMOVED***
			gid = os.Getgid()
		***REMOVED***
		l, err := sockets.NewUnixSocket(addr, gid)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("can't create unix socket %s: %v", addr, err)
		***REMOVED***
		ls = append(ls, l)
	default:
		return nil, fmt.Errorf("invalid protocol format: %q", proto)
	***REMOVED***

	return ls, nil
***REMOVED***

// listenFD returns the specified socket activated files as a slice of
// net.Listeners or all of the activated files if "*" is given.
func listenFD(addr string, tlsConfig *tls.Config) ([]net.Listener, error) ***REMOVED***
	var (
		err       error
		listeners []net.Listener
	)
	// socket activation
	if tlsConfig != nil ***REMOVED***
		listeners, err = activation.TLSListeners(false, tlsConfig)
	***REMOVED*** else ***REMOVED***
		listeners, err = activation.Listeners(false)
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if len(listeners) == 0 ***REMOVED***
		return nil, fmt.Errorf("no sockets found via socket activation: make sure the service was started by systemd")
	***REMOVED***

	// default to all fds just like unix:// and tcp://
	if addr == "" || addr == "*" ***REMOVED***
		return listeners, nil
	***REMOVED***

	fdNum, err := strconv.Atoi(addr)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to parse systemd fd address: should be a number: %v", addr)
	***REMOVED***
	fdOffset := fdNum - 3
	if len(listeners) < fdOffset+1 ***REMOVED***
		return nil, fmt.Errorf("too few socket activated files passed in by systemd")
	***REMOVED***
	if listeners[fdOffset] == nil ***REMOVED***
		return nil, fmt.Errorf("failed to listen on systemd activated file: fd %d", fdOffset+3)
	***REMOVED***
	for i, ls := range listeners ***REMOVED***
		if i == fdOffset || ls == nil ***REMOVED***
			continue
		***REMOVED***
		if err := ls.Close(); err != nil ***REMOVED***
			return nil, fmt.Errorf("failed to close systemd activated file: fd %d: %v", fdOffset+3, err)
		***REMOVED***
	***REMOVED***
	return []net.Listener***REMOVED***listeners[fdOffset]***REMOVED***, nil
***REMOVED***
