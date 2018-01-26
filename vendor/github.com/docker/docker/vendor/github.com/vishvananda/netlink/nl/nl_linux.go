// Package nl has low level primitives for making Netlink calls.
package nl

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"

	"github.com/vishvananda/netns"
)

const (
	// Family type definitions
	FAMILY_ALL  = syscall.AF_UNSPEC
	FAMILY_V4   = syscall.AF_INET
	FAMILY_V6   = syscall.AF_INET6
	FAMILY_MPLS = AF_MPLS
)

// SupportedNlFamilies contains the list of netlink families this netlink package supports
var SupportedNlFamilies = []int***REMOVED***syscall.NETLINK_ROUTE, syscall.NETLINK_XFRM, syscall.NETLINK_NETFILTER***REMOVED***

var nextSeqNr uint32

// GetIPFamily returns the family type of a net.IP.
func GetIPFamily(ip net.IP) int ***REMOVED***
	if len(ip) <= net.IPv4len ***REMOVED***
		return FAMILY_V4
	***REMOVED***
	if ip.To4() != nil ***REMOVED***
		return FAMILY_V4
	***REMOVED***
	return FAMILY_V6
***REMOVED***

var nativeEndian binary.ByteOrder

// Get native endianness for the system
func NativeEndian() binary.ByteOrder ***REMOVED***
	if nativeEndian == nil ***REMOVED***
		var x uint32 = 0x01020304
		if *(*byte)(unsafe.Pointer(&x)) == 0x01 ***REMOVED***
			nativeEndian = binary.BigEndian
		***REMOVED*** else ***REMOVED***
			nativeEndian = binary.LittleEndian
		***REMOVED***
	***REMOVED***
	return nativeEndian
***REMOVED***

// Byte swap a 16 bit value if we aren't big endian
func Swap16(i uint16) uint16 ***REMOVED***
	if NativeEndian() == binary.BigEndian ***REMOVED***
		return i
	***REMOVED***
	return (i&0xff00)>>8 | (i&0xff)<<8
***REMOVED***

// Byte swap a 32 bit value if aren't big endian
func Swap32(i uint32) uint32 ***REMOVED***
	if NativeEndian() == binary.BigEndian ***REMOVED***
		return i
	***REMOVED***
	return (i&0xff000000)>>24 | (i&0xff0000)>>8 | (i&0xff00)<<8 | (i&0xff)<<24
***REMOVED***

type NetlinkRequestData interface ***REMOVED***
	Len() int
	Serialize() []byte
***REMOVED***

// IfInfomsg is related to links, but it is used for list requests as well
type IfInfomsg struct ***REMOVED***
	syscall.IfInfomsg
***REMOVED***

// Create an IfInfomsg with family specified
func NewIfInfomsg(family int) *IfInfomsg ***REMOVED***
	return &IfInfomsg***REMOVED***
		IfInfomsg: syscall.IfInfomsg***REMOVED***
			Family: uint8(family),
		***REMOVED***,
	***REMOVED***
***REMOVED***

func DeserializeIfInfomsg(b []byte) *IfInfomsg ***REMOVED***
	return (*IfInfomsg)(unsafe.Pointer(&b[0:syscall.SizeofIfInfomsg][0]))
***REMOVED***

func (msg *IfInfomsg) Serialize() []byte ***REMOVED***
	return (*(*[syscall.SizeofIfInfomsg]byte)(unsafe.Pointer(msg)))[:]
***REMOVED***

func (msg *IfInfomsg) Len() int ***REMOVED***
	return syscall.SizeofIfInfomsg
***REMOVED***

func (msg *IfInfomsg) EncapType() string ***REMOVED***
	switch msg.Type ***REMOVED***
	case 0:
		return "generic"
	case syscall.ARPHRD_ETHER:
		return "ether"
	case syscall.ARPHRD_EETHER:
		return "eether"
	case syscall.ARPHRD_AX25:
		return "ax25"
	case syscall.ARPHRD_PRONET:
		return "pronet"
	case syscall.ARPHRD_CHAOS:
		return "chaos"
	case syscall.ARPHRD_IEEE802:
		return "ieee802"
	case syscall.ARPHRD_ARCNET:
		return "arcnet"
	case syscall.ARPHRD_APPLETLK:
		return "atalk"
	case syscall.ARPHRD_DLCI:
		return "dlci"
	case syscall.ARPHRD_ATM:
		return "atm"
	case syscall.ARPHRD_METRICOM:
		return "metricom"
	case syscall.ARPHRD_IEEE1394:
		return "ieee1394"
	case syscall.ARPHRD_INFINIBAND:
		return "infiniband"
	case syscall.ARPHRD_SLIP:
		return "slip"
	case syscall.ARPHRD_CSLIP:
		return "cslip"
	case syscall.ARPHRD_SLIP6:
		return "slip6"
	case syscall.ARPHRD_CSLIP6:
		return "cslip6"
	case syscall.ARPHRD_RSRVD:
		return "rsrvd"
	case syscall.ARPHRD_ADAPT:
		return "adapt"
	case syscall.ARPHRD_ROSE:
		return "rose"
	case syscall.ARPHRD_X25:
		return "x25"
	case syscall.ARPHRD_HWX25:
		return "hwx25"
	case syscall.ARPHRD_PPP:
		return "ppp"
	case syscall.ARPHRD_HDLC:
		return "hdlc"
	case syscall.ARPHRD_LAPB:
		return "lapb"
	case syscall.ARPHRD_DDCMP:
		return "ddcmp"
	case syscall.ARPHRD_RAWHDLC:
		return "rawhdlc"
	case syscall.ARPHRD_TUNNEL:
		return "ipip"
	case syscall.ARPHRD_TUNNEL6:
		return "tunnel6"
	case syscall.ARPHRD_FRAD:
		return "frad"
	case syscall.ARPHRD_SKIP:
		return "skip"
	case syscall.ARPHRD_LOOPBACK:
		return "loopback"
	case syscall.ARPHRD_LOCALTLK:
		return "ltalk"
	case syscall.ARPHRD_FDDI:
		return "fddi"
	case syscall.ARPHRD_BIF:
		return "bif"
	case syscall.ARPHRD_SIT:
		return "sit"
	case syscall.ARPHRD_IPDDP:
		return "ip/ddp"
	case syscall.ARPHRD_IPGRE:
		return "gre"
	case syscall.ARPHRD_PIMREG:
		return "pimreg"
	case syscall.ARPHRD_HIPPI:
		return "hippi"
	case syscall.ARPHRD_ASH:
		return "ash"
	case syscall.ARPHRD_ECONET:
		return "econet"
	case syscall.ARPHRD_IRDA:
		return "irda"
	case syscall.ARPHRD_FCPP:
		return "fcpp"
	case syscall.ARPHRD_FCAL:
		return "fcal"
	case syscall.ARPHRD_FCPL:
		return "fcpl"
	case syscall.ARPHRD_FCFABRIC:
		return "fcfb0"
	case syscall.ARPHRD_FCFABRIC + 1:
		return "fcfb1"
	case syscall.ARPHRD_FCFABRIC + 2:
		return "fcfb2"
	case syscall.ARPHRD_FCFABRIC + 3:
		return "fcfb3"
	case syscall.ARPHRD_FCFABRIC + 4:
		return "fcfb4"
	case syscall.ARPHRD_FCFABRIC + 5:
		return "fcfb5"
	case syscall.ARPHRD_FCFABRIC + 6:
		return "fcfb6"
	case syscall.ARPHRD_FCFABRIC + 7:
		return "fcfb7"
	case syscall.ARPHRD_FCFABRIC + 8:
		return "fcfb8"
	case syscall.ARPHRD_FCFABRIC + 9:
		return "fcfb9"
	case syscall.ARPHRD_FCFABRIC + 10:
		return "fcfb10"
	case syscall.ARPHRD_FCFABRIC + 11:
		return "fcfb11"
	case syscall.ARPHRD_FCFABRIC + 12:
		return "fcfb12"
	case syscall.ARPHRD_IEEE802_TR:
		return "tr"
	case syscall.ARPHRD_IEEE80211:
		return "ieee802.11"
	case syscall.ARPHRD_IEEE80211_PRISM:
		return "ieee802.11/prism"
	case syscall.ARPHRD_IEEE80211_RADIOTAP:
		return "ieee802.11/radiotap"
	case syscall.ARPHRD_IEEE802154:
		return "ieee802.15.4"

	case 65534:
		return "none"
	case 65535:
		return "void"
	***REMOVED***
	return fmt.Sprintf("unknown%d", msg.Type)
***REMOVED***

func rtaAlignOf(attrlen int) int ***REMOVED***
	return (attrlen + syscall.RTA_ALIGNTO - 1) & ^(syscall.RTA_ALIGNTO - 1)
***REMOVED***

func NewIfInfomsgChild(parent *RtAttr, family int) *IfInfomsg ***REMOVED***
	msg := NewIfInfomsg(family)
	parent.children = append(parent.children, msg)
	return msg
***REMOVED***

// Extend RtAttr to handle data and children
type RtAttr struct ***REMOVED***
	syscall.RtAttr
	Data     []byte
	children []NetlinkRequestData
***REMOVED***

// Create a new Extended RtAttr object
func NewRtAttr(attrType int, data []byte) *RtAttr ***REMOVED***
	return &RtAttr***REMOVED***
		RtAttr: syscall.RtAttr***REMOVED***
			Type: uint16(attrType),
		***REMOVED***,
		children: []NetlinkRequestData***REMOVED******REMOVED***,
		Data:     data,
	***REMOVED***
***REMOVED***

// Create a new RtAttr obj anc add it as a child of an existing object
func NewRtAttrChild(parent *RtAttr, attrType int, data []byte) *RtAttr ***REMOVED***
	attr := NewRtAttr(attrType, data)
	parent.children = append(parent.children, attr)
	return attr
***REMOVED***

func (a *RtAttr) Len() int ***REMOVED***
	if len(a.children) == 0 ***REMOVED***
		return (syscall.SizeofRtAttr + len(a.Data))
	***REMOVED***

	l := 0
	for _, child := range a.children ***REMOVED***
		l += rtaAlignOf(child.Len())
	***REMOVED***
	l += syscall.SizeofRtAttr
	return rtaAlignOf(l + len(a.Data))
***REMOVED***

// Serialize the RtAttr into a byte array
// This can't just unsafe.cast because it must iterate through children.
func (a *RtAttr) Serialize() []byte ***REMOVED***
	native := NativeEndian()

	length := a.Len()
	buf := make([]byte, rtaAlignOf(length))

	next := 4
	if a.Data != nil ***REMOVED***
		copy(buf[next:], a.Data)
		next += rtaAlignOf(len(a.Data))
	***REMOVED***
	if len(a.children) > 0 ***REMOVED***
		for _, child := range a.children ***REMOVED***
			childBuf := child.Serialize()
			copy(buf[next:], childBuf)
			next += rtaAlignOf(len(childBuf))
		***REMOVED***
	***REMOVED***

	if l := uint16(length); l != 0 ***REMOVED***
		native.PutUint16(buf[0:2], l)
	***REMOVED***
	native.PutUint16(buf[2:4], a.Type)
	return buf
***REMOVED***

type NetlinkRequest struct ***REMOVED***
	syscall.NlMsghdr
	Data    []NetlinkRequestData
	RawData []byte
	Sockets map[int]*SocketHandle
***REMOVED***

// Serialize the Netlink Request into a byte array
func (req *NetlinkRequest) Serialize() []byte ***REMOVED***
	length := syscall.SizeofNlMsghdr
	dataBytes := make([][]byte, len(req.Data))
	for i, data := range req.Data ***REMOVED***
		dataBytes[i] = data.Serialize()
		length = length + len(dataBytes[i])
	***REMOVED***
	length += len(req.RawData)

	req.Len = uint32(length)
	b := make([]byte, length)
	hdr := (*(*[syscall.SizeofNlMsghdr]byte)(unsafe.Pointer(req)))[:]
	next := syscall.SizeofNlMsghdr
	copy(b[0:next], hdr)
	for _, data := range dataBytes ***REMOVED***
		for _, dataByte := range data ***REMOVED***
			b[next] = dataByte
			next = next + 1
		***REMOVED***
	***REMOVED***
	// Add the raw data if any
	if len(req.RawData) > 0 ***REMOVED***
		copy(b[next:length], req.RawData)
	***REMOVED***
	return b
***REMOVED***

func (req *NetlinkRequest) AddData(data NetlinkRequestData) ***REMOVED***
	if data != nil ***REMOVED***
		req.Data = append(req.Data, data)
	***REMOVED***
***REMOVED***

// AddRawData adds raw bytes to the end of the NetlinkRequest object during serialization
func (req *NetlinkRequest) AddRawData(data []byte) ***REMOVED***
	if data != nil ***REMOVED***
		req.RawData = append(req.RawData, data...)
	***REMOVED***
***REMOVED***

// Execute the request against a the given sockType.
// Returns a list of netlink messages in serialized format, optionally filtered
// by resType.
func (req *NetlinkRequest) Execute(sockType int, resType uint16) ([][]byte, error) ***REMOVED***
	var (
		s   *NetlinkSocket
		err error
	)

	if req.Sockets != nil ***REMOVED***
		if sh, ok := req.Sockets[sockType]; ok ***REMOVED***
			s = sh.Socket
			req.Seq = atomic.AddUint32(&sh.Seq, 1)
		***REMOVED***
	***REMOVED***
	sharedSocket := s != nil

	if s == nil ***REMOVED***
		s, err = getNetlinkSocket(sockType)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		defer s.Close()
	***REMOVED*** else ***REMOVED***
		s.Lock()
		defer s.Unlock()
	***REMOVED***

	if err := s.Send(req); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	pid, err := s.GetPid()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var res [][]byte

done:
	for ***REMOVED***
		msgs, err := s.Receive()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		for _, m := range msgs ***REMOVED***
			if m.Header.Seq != req.Seq ***REMOVED***
				if sharedSocket ***REMOVED***
					continue
				***REMOVED***
				return nil, fmt.Errorf("Wrong Seq nr %d, expected %d", m.Header.Seq, req.Seq)
			***REMOVED***
			if m.Header.Pid != pid ***REMOVED***
				return nil, fmt.Errorf("Wrong pid %d, expected %d", m.Header.Pid, pid)
			***REMOVED***
			if m.Header.Type == syscall.NLMSG_DONE ***REMOVED***
				break done
			***REMOVED***
			if m.Header.Type == syscall.NLMSG_ERROR ***REMOVED***
				native := NativeEndian()
				error := int32(native.Uint32(m.Data[0:4]))
				if error == 0 ***REMOVED***
					break done
				***REMOVED***
				return nil, syscall.Errno(-error)
			***REMOVED***
			if resType != 0 && m.Header.Type != resType ***REMOVED***
				continue
			***REMOVED***
			res = append(res, m.Data)
			if m.Header.Flags&syscall.NLM_F_MULTI == 0 ***REMOVED***
				break done
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return res, nil
***REMOVED***

// Create a new netlink request from proto and flags
// Note the Len value will be inaccurate once data is added until
// the message is serialized
func NewNetlinkRequest(proto, flags int) *NetlinkRequest ***REMOVED***
	return &NetlinkRequest***REMOVED***
		NlMsghdr: syscall.NlMsghdr***REMOVED***
			Len:   uint32(syscall.SizeofNlMsghdr),
			Type:  uint16(proto),
			Flags: syscall.NLM_F_REQUEST | uint16(flags),
			Seq:   atomic.AddUint32(&nextSeqNr, 1),
		***REMOVED***,
	***REMOVED***
***REMOVED***

type NetlinkSocket struct ***REMOVED***
	fd  int32
	lsa syscall.SockaddrNetlink
	sync.Mutex
***REMOVED***

func getNetlinkSocket(protocol int) (*NetlinkSocket, error) ***REMOVED***
	fd, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW|syscall.SOCK_CLOEXEC, protocol)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	s := &NetlinkSocket***REMOVED***
		fd: int32(fd),
	***REMOVED***
	s.lsa.Family = syscall.AF_NETLINK
	if err := syscall.Bind(fd, &s.lsa); err != nil ***REMOVED***
		syscall.Close(fd)
		return nil, err
	***REMOVED***

	return s, nil
***REMOVED***

// GetNetlinkSocketAt opens a netlink socket in the network namespace newNs
// and positions the thread back into the network namespace specified by curNs,
// when done. If curNs is close, the function derives the current namespace and
// moves back into it when done. If newNs is close, the socket will be opened
// in the current network namespace.
func GetNetlinkSocketAt(newNs, curNs netns.NsHandle, protocol int) (*NetlinkSocket, error) ***REMOVED***
	c, err := executeInNetns(newNs, curNs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer c()
	return getNetlinkSocket(protocol)
***REMOVED***

// executeInNetns sets execution of the code following this call to the
// network namespace newNs, then moves the thread back to curNs if open,
// otherwise to the current netns at the time the function was invoked
// In case of success, the caller is expected to execute the returned function
// at the end of the code that needs to be executed in the network namespace.
// Example:
// func jobAt(...) error ***REMOVED***
//      d, err := executeInNetns(...)
//      if err != nil ***REMOVED*** return err***REMOVED***
//      defer d()
//      < code which needs to be executed in specific netns>
//  ***REMOVED***
// TODO: his function probably belongs to netns pkg.
func executeInNetns(newNs, curNs netns.NsHandle) (func(), error) ***REMOVED***
	var (
		err       error
		moveBack  func(netns.NsHandle) error
		closeNs   func() error
		unlockThd func()
	)
	restore := func() ***REMOVED***
		// order matters
		if moveBack != nil ***REMOVED***
			moveBack(curNs)
		***REMOVED***
		if closeNs != nil ***REMOVED***
			closeNs()
		***REMOVED***
		if unlockThd != nil ***REMOVED***
			unlockThd()
		***REMOVED***
	***REMOVED***
	if newNs.IsOpen() ***REMOVED***
		runtime.LockOSThread()
		unlockThd = runtime.UnlockOSThread
		if !curNs.IsOpen() ***REMOVED***
			if curNs, err = netns.Get(); err != nil ***REMOVED***
				restore()
				return nil, fmt.Errorf("could not get current namespace while creating netlink socket: %v", err)
			***REMOVED***
			closeNs = curNs.Close
		***REMOVED***
		if err := netns.Set(newNs); err != nil ***REMOVED***
			restore()
			return nil, fmt.Errorf("failed to set into network namespace %d while creating netlink socket: %v", newNs, err)
		***REMOVED***
		moveBack = netns.Set
	***REMOVED***
	return restore, nil
***REMOVED***

// Create a netlink socket with a given protocol (e.g. NETLINK_ROUTE)
// and subscribe it to multicast groups passed in variable argument list.
// Returns the netlink socket on which Receive() method can be called
// to retrieve the messages from the kernel.
func Subscribe(protocol int, groups ...uint) (*NetlinkSocket, error) ***REMOVED***
	fd, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW, protocol)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	s := &NetlinkSocket***REMOVED***
		fd: int32(fd),
	***REMOVED***
	s.lsa.Family = syscall.AF_NETLINK

	for _, g := range groups ***REMOVED***
		s.lsa.Groups |= (1 << (g - 1))
	***REMOVED***

	if err := syscall.Bind(fd, &s.lsa); err != nil ***REMOVED***
		syscall.Close(fd)
		return nil, err
	***REMOVED***

	return s, nil
***REMOVED***

// SubscribeAt works like Subscribe plus let's the caller choose the network
// namespace in which the socket would be opened (newNs). Then control goes back
// to curNs if open, otherwise to the netns at the time this function was called.
func SubscribeAt(newNs, curNs netns.NsHandle, protocol int, groups ...uint) (*NetlinkSocket, error) ***REMOVED***
	c, err := executeInNetns(newNs, curNs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer c()
	return Subscribe(protocol, groups...)
***REMOVED***

func (s *NetlinkSocket) Close() ***REMOVED***
	fd := int(atomic.SwapInt32(&s.fd, -1))
	syscall.Close(fd)
***REMOVED***

func (s *NetlinkSocket) GetFd() int ***REMOVED***
	return int(atomic.LoadInt32(&s.fd))
***REMOVED***

func (s *NetlinkSocket) Send(request *NetlinkRequest) error ***REMOVED***
	fd := int(atomic.LoadInt32(&s.fd))
	if fd < 0 ***REMOVED***
		return fmt.Errorf("Send called on a closed socket")
	***REMOVED***
	if err := syscall.Sendto(fd, request.Serialize(), 0, &s.lsa); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func (s *NetlinkSocket) Receive() ([]syscall.NetlinkMessage, error) ***REMOVED***
	fd := int(atomic.LoadInt32(&s.fd))
	if fd < 0 ***REMOVED***
		return nil, fmt.Errorf("Receive called on a closed socket")
	***REMOVED***
	rb := make([]byte, syscall.Getpagesize())
	nr, _, err := syscall.Recvfrom(fd, rb, 0)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if nr < syscall.NLMSG_HDRLEN ***REMOVED***
		return nil, fmt.Errorf("Got short response from netlink")
	***REMOVED***
	rb = rb[:nr]
	return syscall.ParseNetlinkMessage(rb)
***REMOVED***

// SetSendTimeout allows to set a send timeout on the socket
func (s *NetlinkSocket) SetSendTimeout(timeout *syscall.Timeval) error ***REMOVED***
	// Set a send timeout of SOCKET_SEND_TIMEOUT, this will allow the Send to periodically unblock and avoid that a routine
	// remains stuck on a send on a closed fd
	return syscall.SetsockoptTimeval(int(s.fd), syscall.SOL_SOCKET, syscall.SO_SNDTIMEO, timeout)
***REMOVED***

// SetReceiveTimeout allows to set a receive timeout on the socket
func (s *NetlinkSocket) SetReceiveTimeout(timeout *syscall.Timeval) error ***REMOVED***
	// Set a read timeout of SOCKET_READ_TIMEOUT, this will allow the Read to periodically unblock and avoid that a routine
	// remains stuck on a recvmsg on a closed fd
	return syscall.SetsockoptTimeval(int(s.fd), syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, timeout)
***REMOVED***

func (s *NetlinkSocket) GetPid() (uint32, error) ***REMOVED***
	fd := int(atomic.LoadInt32(&s.fd))
	lsa, err := syscall.Getsockname(fd)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	switch v := lsa.(type) ***REMOVED***
	case *syscall.SockaddrNetlink:
		return v.Pid, nil
	***REMOVED***
	return 0, fmt.Errorf("Wrong socket type")
***REMOVED***

func ZeroTerminated(s string) []byte ***REMOVED***
	bytes := make([]byte, len(s)+1)
	for i := 0; i < len(s); i++ ***REMOVED***
		bytes[i] = s[i]
	***REMOVED***
	bytes[len(s)] = 0
	return bytes
***REMOVED***

func NonZeroTerminated(s string) []byte ***REMOVED***
	bytes := make([]byte, len(s))
	for i := 0; i < len(s); i++ ***REMOVED***
		bytes[i] = s[i]
	***REMOVED***
	return bytes
***REMOVED***

func BytesToString(b []byte) string ***REMOVED***
	n := bytes.Index(b, []byte***REMOVED***0***REMOVED***)
	return string(b[:n])
***REMOVED***

func Uint8Attr(v uint8) []byte ***REMOVED***
	return []byte***REMOVED***byte(v)***REMOVED***
***REMOVED***

func Uint16Attr(v uint16) []byte ***REMOVED***
	native := NativeEndian()
	bytes := make([]byte, 2)
	native.PutUint16(bytes, v)
	return bytes
***REMOVED***

func Uint32Attr(v uint32) []byte ***REMOVED***
	native := NativeEndian()
	bytes := make([]byte, 4)
	native.PutUint32(bytes, v)
	return bytes
***REMOVED***

func Uint64Attr(v uint64) []byte ***REMOVED***
	native := NativeEndian()
	bytes := make([]byte, 8)
	native.PutUint64(bytes, v)
	return bytes
***REMOVED***

func ParseRouteAttr(b []byte) ([]syscall.NetlinkRouteAttr, error) ***REMOVED***
	var attrs []syscall.NetlinkRouteAttr
	for len(b) >= syscall.SizeofRtAttr ***REMOVED***
		a, vbuf, alen, err := netlinkRouteAttrAndValue(b)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		ra := syscall.NetlinkRouteAttr***REMOVED***Attr: *a, Value: vbuf[:int(a.Len)-syscall.SizeofRtAttr]***REMOVED***
		attrs = append(attrs, ra)
		b = b[alen:]
	***REMOVED***
	return attrs, nil
***REMOVED***

func netlinkRouteAttrAndValue(b []byte) (*syscall.RtAttr, []byte, int, error) ***REMOVED***
	a := (*syscall.RtAttr)(unsafe.Pointer(&b[0]))
	if int(a.Len) < syscall.SizeofRtAttr || int(a.Len) > len(b) ***REMOVED***
		return nil, nil, 0, syscall.EINVAL
	***REMOVED***
	return a, b[syscall.SizeofRtAttr:], rtaAlignOf(int(a.Len)), nil
***REMOVED***

// SocketHandle contains the netlink socket and the associated
// sequence counter for a specific netlink family
type SocketHandle struct ***REMOVED***
	Seq    uint32
	Socket *NetlinkSocket
***REMOVED***

// Close closes the netlink socket
func (sh *SocketHandle) Close() ***REMOVED***
	if sh.Socket != nil ***REMOVED***
		sh.Socket.Close()
	***REMOVED***
***REMOVED***
