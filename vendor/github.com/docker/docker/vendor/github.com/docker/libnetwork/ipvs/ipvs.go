// +build linux

package ipvs

import (
	"net"
	"syscall"
	"time"

	"fmt"

	"github.com/vishvananda/netlink/nl"
	"github.com/vishvananda/netns"
)

const (
	netlinkRecvSocketsTimeout = 3 * time.Second
	netlinkSendSocketTimeout  = 30 * time.Second
)

// Service defines an IPVS service in its entirety.
type Service struct ***REMOVED***
	// Virtual service address.
	Address  net.IP
	Protocol uint16
	Port     uint16
	FWMark   uint32 // Firewall mark of the service.

	// Virtual service options.
	SchedName     string
	Flags         uint32
	Timeout       uint32
	Netmask       uint32
	AddressFamily uint16
	PEName        string
	Stats         SvcStats
***REMOVED***

// SvcStats defines an IPVS service statistics
type SvcStats struct ***REMOVED***
	Connections uint32
	PacketsIn   uint32
	PacketsOut  uint32
	BytesIn     uint64
	BytesOut    uint64
	CPS         uint32
	BPSOut      uint32
	PPSIn       uint32
	PPSOut      uint32
	BPSIn       uint32
***REMOVED***

// Destination defines an IPVS destination (real server) in its
// entirety.
type Destination struct ***REMOVED***
	Address         net.IP
	Port            uint16
	Weight          int
	ConnectionFlags uint32
	AddressFamily   uint16
	UpperThreshold  uint32
	LowerThreshold  uint32
***REMOVED***

// Handle provides a namespace specific ipvs handle to program ipvs
// rules.
type Handle struct ***REMOVED***
	seq  uint32
	sock *nl.NetlinkSocket
***REMOVED***

// New provides a new ipvs handle in the namespace pointed to by the
// passed path. It will return a valid handle or an error in case an
// error occurred while creating the handle.
func New(path string) (*Handle, error) ***REMOVED***
	setup()

	n := netns.None()
	if path != "" ***REMOVED***
		var err error
		n, err = netns.GetFromPath(path)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	defer n.Close()

	sock, err := nl.GetNetlinkSocketAt(n, netns.None(), syscall.NETLINK_GENERIC)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// Add operation timeout to avoid deadlocks
	tv := syscall.NsecToTimeval(netlinkSendSocketTimeout.Nanoseconds())
	if err := sock.SetSendTimeout(&tv); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	tv = syscall.NsecToTimeval(netlinkRecvSocketsTimeout.Nanoseconds())
	if err := sock.SetReceiveTimeout(&tv); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Handle***REMOVED***sock: sock***REMOVED***, nil
***REMOVED***

// Close closes the ipvs handle. The handle is invalid after Close
// returns.
func (i *Handle) Close() ***REMOVED***
	if i.sock != nil ***REMOVED***
		i.sock.Close()
	***REMOVED***
***REMOVED***

// NewService creates a new ipvs service in the passed handle.
func (i *Handle) NewService(s *Service) error ***REMOVED***
	return i.doCmd(s, nil, ipvsCmdNewService)
***REMOVED***

// IsServicePresent queries for the ipvs service in the passed handle.
func (i *Handle) IsServicePresent(s *Service) bool ***REMOVED***
	return nil == i.doCmd(s, nil, ipvsCmdGetService)
***REMOVED***

// UpdateService updates an already existing service in the passed
// handle.
func (i *Handle) UpdateService(s *Service) error ***REMOVED***
	return i.doCmd(s, nil, ipvsCmdSetService)
***REMOVED***

// DelService deletes an already existing service in the passed
// handle.
func (i *Handle) DelService(s *Service) error ***REMOVED***
	return i.doCmd(s, nil, ipvsCmdDelService)
***REMOVED***

// Flush deletes all existing services in the passed
// handle.
func (i *Handle) Flush() error ***REMOVED***
	_, err := i.doCmdWithoutAttr(ipvsCmdFlush)
	return err
***REMOVED***

// NewDestination creates a new real server in the passed ipvs
// service which should already be existing in the passed handle.
func (i *Handle) NewDestination(s *Service, d *Destination) error ***REMOVED***
	return i.doCmd(s, d, ipvsCmdNewDest)
***REMOVED***

// UpdateDestination updates an already existing real server in the
// passed ipvs service in the passed handle.
func (i *Handle) UpdateDestination(s *Service, d *Destination) error ***REMOVED***
	return i.doCmd(s, d, ipvsCmdSetDest)
***REMOVED***

// DelDestination deletes an already existing real server in the
// passed ipvs service in the passed handle.
func (i *Handle) DelDestination(s *Service, d *Destination) error ***REMOVED***
	return i.doCmd(s, d, ipvsCmdDelDest)
***REMOVED***

// GetServices returns an array of services configured on the Node
func (i *Handle) GetServices() ([]*Service, error) ***REMOVED***
	return i.doGetServicesCmd(nil)
***REMOVED***

// GetDestinations returns an array of Destinations configured for this Service
func (i *Handle) GetDestinations(s *Service) ([]*Destination, error) ***REMOVED***
	return i.doGetDestinationsCmd(s, nil)
***REMOVED***

// GetService gets details of a specific IPVS services, useful in updating statisics etc.,
func (i *Handle) GetService(s *Service) (*Service, error) ***REMOVED***

	res, err := i.doGetServicesCmd(s)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// We are looking for exactly one service otherwise error out
	if len(res) != 1 ***REMOVED***
		return nil, fmt.Errorf("Expected only one service obtained=%d", len(res))
	***REMOVED***

	return res[0], nil
***REMOVED***
