package netlink

import (
	"fmt"
	"syscall"
	"time"

	"github.com/vishvananda/netlink/nl"
	"github.com/vishvananda/netns"
)

// Empty handle used by the netlink package methods
var pkgHandle = &Handle***REMOVED******REMOVED***

// Handle is an handle for the netlink requests on a
// specific network namespace. All the requests on the
// same netlink family share the same netlink socket,
// which gets released when the handle is deleted.
type Handle struct ***REMOVED***
	sockets      map[int]*nl.SocketHandle
	lookupByDump bool
***REMOVED***

// SupportsNetlinkFamily reports whether the passed netlink family is supported by this Handle
func (h *Handle) SupportsNetlinkFamily(nlFamily int) bool ***REMOVED***
	_, ok := h.sockets[nlFamily]
	return ok
***REMOVED***

// NewHandle returns a netlink handle on the current network namespace.
// Caller may specify the netlink families the handle should support.
// If no families are specified, all the families the netlink package
// supports will be automatically added.
func NewHandle(nlFamilies ...int) (*Handle, error) ***REMOVED***
	return newHandle(netns.None(), netns.None(), nlFamilies...)
***REMOVED***

// SetSocketTimeout sets the send and receive timeout for each socket in the
// netlink handle. Although the socket timeout has granularity of one
// microsecond, the effective granularity is floored by the kernel timer tick,
// which default value is four milliseconds.
func (h *Handle) SetSocketTimeout(to time.Duration) error ***REMOVED***
	if to < time.Microsecond ***REMOVED***
		return fmt.Errorf("invalid timeout, minimul value is %s", time.Microsecond)
	***REMOVED***
	tv := syscall.NsecToTimeval(to.Nanoseconds())
	for _, sh := range h.sockets ***REMOVED***
		if err := sh.Socket.SetSendTimeout(&tv); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := sh.Socket.SetReceiveTimeout(&tv); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// SetSocketReceiveBufferSize sets the receive buffer size for each
// socket in the netlink handle. The maximum value is capped by
// /proc/sys/net/core/rmem_max.
func (h *Handle) SetSocketReceiveBufferSize(size int, force bool) error ***REMOVED***
	opt := syscall.SO_RCVBUF
	if force ***REMOVED***
		opt = syscall.SO_RCVBUFFORCE
	***REMOVED***
	for _, sh := range h.sockets ***REMOVED***
		fd := sh.Socket.GetFd()
		err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, opt, size)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// GetSocketReceiveBufferSize gets the receiver buffer size for each
// socket in the netlink handle. The retrieved value should be the
// double to the one set for SetSocketReceiveBufferSize.
func (h *Handle) GetSocketReceiveBufferSize() ([]int, error) ***REMOVED***
	results := make([]int, len(h.sockets))
	i := 0
	for _, sh := range h.sockets ***REMOVED***
		fd := sh.Socket.GetFd()
		size, err := syscall.GetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_RCVBUF)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		results[i] = size
		i++
	***REMOVED***
	return results, nil
***REMOVED***

// NewHandle returns a netlink handle on the network namespace
// specified by ns. If ns=netns.None(), current network namespace
// will be assumed
func NewHandleAt(ns netns.NsHandle, nlFamilies ...int) (*Handle, error) ***REMOVED***
	return newHandle(ns, netns.None(), nlFamilies...)
***REMOVED***

// NewHandleAtFrom works as NewHandle but allows client to specify the
// new and the origin netns Handle.
func NewHandleAtFrom(newNs, curNs netns.NsHandle) (*Handle, error) ***REMOVED***
	return newHandle(newNs, curNs)
***REMOVED***

func newHandle(newNs, curNs netns.NsHandle, nlFamilies ...int) (*Handle, error) ***REMOVED***
	h := &Handle***REMOVED***sockets: map[int]*nl.SocketHandle***REMOVED******REMOVED******REMOVED***
	fams := nl.SupportedNlFamilies
	if len(nlFamilies) != 0 ***REMOVED***
		fams = nlFamilies
	***REMOVED***
	for _, f := range fams ***REMOVED***
		s, err := nl.GetNetlinkSocketAt(newNs, curNs, f)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		h.sockets[f] = &nl.SocketHandle***REMOVED***Socket: s***REMOVED***
	***REMOVED***
	return h, nil
***REMOVED***

// Delete releases the resources allocated to this handle
func (h *Handle) Delete() ***REMOVED***
	for _, sh := range h.sockets ***REMOVED***
		sh.Close()
	***REMOVED***
	h.sockets = nil
***REMOVED***

func (h *Handle) newNetlinkRequest(proto, flags int) *nl.NetlinkRequest ***REMOVED***
	// Do this so that package API still use nl package variable nextSeqNr
	if h.sockets == nil ***REMOVED***
		return nl.NewNetlinkRequest(proto, flags)
	***REMOVED***
	return &nl.NetlinkRequest***REMOVED***
		NlMsghdr: syscall.NlMsghdr***REMOVED***
			Len:   uint32(syscall.SizeofNlMsghdr),
			Type:  uint16(proto),
			Flags: syscall.NLM_F_REQUEST | uint16(flags),
		***REMOVED***,
		Sockets: h.sockets,
	***REMOVED***
***REMOVED***
