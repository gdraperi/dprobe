package bridge

import (
	"bytes"
	"errors"
	"fmt"
	"net"

	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

var (
	defaultBindingIP = net.IPv4(0, 0, 0, 0)
)

func (n *bridgeNetwork) allocatePorts(ep *bridgeEndpoint, reqDefBindIP net.IP, ulPxyEnabled bool) ([]types.PortBinding, error) ***REMOVED***
	if ep.extConnConfig == nil || ep.extConnConfig.PortBindings == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	defHostIP := defaultBindingIP
	if reqDefBindIP != nil ***REMOVED***
		defHostIP = reqDefBindIP
	***REMOVED***

	return n.allocatePortsInternal(ep.extConnConfig.PortBindings, ep.addr.IP, defHostIP, ulPxyEnabled)
***REMOVED***

func (n *bridgeNetwork) allocatePortsInternal(bindings []types.PortBinding, containerIP, defHostIP net.IP, ulPxyEnabled bool) ([]types.PortBinding, error) ***REMOVED***
	bs := make([]types.PortBinding, 0, len(bindings))
	for _, c := range bindings ***REMOVED***
		b := c.GetCopy()
		if err := n.allocatePort(&b, containerIP, defHostIP, ulPxyEnabled); err != nil ***REMOVED***
			// On allocation failure, release previously allocated ports. On cleanup error, just log a warning message
			if cuErr := n.releasePortsInternal(bs); cuErr != nil ***REMOVED***
				logrus.Warnf("Upon allocation failure for %v, failed to clear previously allocated port bindings: %v", b, cuErr)
			***REMOVED***
			return nil, err
		***REMOVED***
		bs = append(bs, b)
	***REMOVED***
	return bs, nil
***REMOVED***

func (n *bridgeNetwork) allocatePort(bnd *types.PortBinding, containerIP, defHostIP net.IP, ulPxyEnabled bool) error ***REMOVED***
	var (
		host net.Addr
		err  error
	)

	// Store the container interface address in the operational binding
	bnd.IP = containerIP

	// Adjust the host address in the operational binding
	if len(bnd.HostIP) == 0 ***REMOVED***
		bnd.HostIP = defHostIP
	***REMOVED***

	// Adjust HostPortEnd if this is not a range.
	if bnd.HostPortEnd == 0 ***REMOVED***
		bnd.HostPortEnd = bnd.HostPort
	***REMOVED***

	// Construct the container side transport address
	container, err := bnd.ContainerAddr()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Try up to maxAllocatePortAttempts times to get a port that's not already allocated.
	for i := 0; i < maxAllocatePortAttempts; i++ ***REMOVED***
		if host, err = n.portMapper.MapRange(container, bnd.HostIP, int(bnd.HostPort), int(bnd.HostPortEnd), ulPxyEnabled); err == nil ***REMOVED***
			break
		***REMOVED***
		// There is no point in immediately retrying to map an explicitly chosen port.
		if bnd.HostPort != 0 ***REMOVED***
			logrus.Warnf("Failed to allocate and map port %d-%d: %s", bnd.HostPort, bnd.HostPortEnd, err)
			break
		***REMOVED***
		logrus.Warnf("Failed to allocate and map port: %s, retry: %d", err, i+1)
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Save the host port (regardless it was or not specified in the binding)
	switch netAddr := host.(type) ***REMOVED***
	case *net.TCPAddr:
		bnd.HostPort = uint16(host.(*net.TCPAddr).Port)
		return nil
	case *net.UDPAddr:
		bnd.HostPort = uint16(host.(*net.UDPAddr).Port)
		return nil
	default:
		// For completeness
		return ErrUnsupportedAddressType(fmt.Sprintf("%T", netAddr))
	***REMOVED***
***REMOVED***

func (n *bridgeNetwork) releasePorts(ep *bridgeEndpoint) error ***REMOVED***
	return n.releasePortsInternal(ep.portMapping)
***REMOVED***

func (n *bridgeNetwork) releasePortsInternal(bindings []types.PortBinding) error ***REMOVED***
	var errorBuf bytes.Buffer

	// Attempt to release all port bindings, do not stop on failure
	for _, m := range bindings ***REMOVED***
		if err := n.releasePort(m); err != nil ***REMOVED***
			errorBuf.WriteString(fmt.Sprintf("\ncould not release %v because of %v", m, err))
		***REMOVED***
	***REMOVED***

	if errorBuf.Len() != 0 ***REMOVED***
		return errors.New(errorBuf.String())
	***REMOVED***
	return nil
***REMOVED***

func (n *bridgeNetwork) releasePort(bnd types.PortBinding) error ***REMOVED***
	// Construct the host side transport address
	host, err := bnd.HostAddr()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return n.portMapper.Unmap(host)
***REMOVED***
