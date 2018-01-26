// Package types contains types that are common across libnetwork project
package types

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// constants for the IP address type
const (
	IP = iota // IPv4 and IPv6
	IPv4
	IPv6
)

// EncryptionKey is the libnetwork representation of the key distributed by the lead
// manager.
type EncryptionKey struct ***REMOVED***
	Subsystem   string
	Algorithm   int32
	Key         []byte
	LamportTime uint64
***REMOVED***

// UUID represents a globally unique ID of various resources like network and endpoint
type UUID string

// QosPolicy represents a quality of service policy on an endpoint
type QosPolicy struct ***REMOVED***
	MaxEgressBandwidth uint64
***REMOVED***

// TransportPort represents a local Layer 4 endpoint
type TransportPort struct ***REMOVED***
	Proto Protocol
	Port  uint16
***REMOVED***

// Equal checks if this instance of Transportport is equal to the passed one
func (t *TransportPort) Equal(o *TransportPort) bool ***REMOVED***
	if t == o ***REMOVED***
		return true
	***REMOVED***

	if o == nil ***REMOVED***
		return false
	***REMOVED***

	if t.Proto != o.Proto || t.Port != o.Port ***REMOVED***
		return false
	***REMOVED***

	return true
***REMOVED***

// GetCopy returns a copy of this TransportPort structure instance
func (t *TransportPort) GetCopy() TransportPort ***REMOVED***
	return TransportPort***REMOVED***Proto: t.Proto, Port: t.Port***REMOVED***
***REMOVED***

// String returns the TransportPort structure in string form
func (t *TransportPort) String() string ***REMOVED***
	return fmt.Sprintf("%s/%d", t.Proto.String(), t.Port)
***REMOVED***

// FromString reads the TransportPort structure from string
func (t *TransportPort) FromString(s string) error ***REMOVED***
	ps := strings.Split(s, "/")
	if len(ps) == 2 ***REMOVED***
		t.Proto = ParseProtocol(ps[0])
		if p, err := strconv.ParseUint(ps[1], 10, 16); err == nil ***REMOVED***
			t.Port = uint16(p)
			return nil
		***REMOVED***
	***REMOVED***
	return BadRequestErrorf("invalid format for transport port: %s", s)
***REMOVED***

// PortBinding represents a port binding between the container and the host
type PortBinding struct ***REMOVED***
	Proto       Protocol
	IP          net.IP
	Port        uint16
	HostIP      net.IP
	HostPort    uint16
	HostPortEnd uint16
***REMOVED***

// HostAddr returns the host side transport address
func (p PortBinding) HostAddr() (net.Addr, error) ***REMOVED***
	switch p.Proto ***REMOVED***
	case UDP:
		return &net.UDPAddr***REMOVED***IP: p.HostIP, Port: int(p.HostPort)***REMOVED***, nil
	case TCP:
		return &net.TCPAddr***REMOVED***IP: p.HostIP, Port: int(p.HostPort)***REMOVED***, nil
	default:
		return nil, ErrInvalidProtocolBinding(p.Proto.String())
	***REMOVED***
***REMOVED***

// ContainerAddr returns the container side transport address
func (p PortBinding) ContainerAddr() (net.Addr, error) ***REMOVED***
	switch p.Proto ***REMOVED***
	case UDP:
		return &net.UDPAddr***REMOVED***IP: p.IP, Port: int(p.Port)***REMOVED***, nil
	case TCP:
		return &net.TCPAddr***REMOVED***IP: p.IP, Port: int(p.Port)***REMOVED***, nil
	default:
		return nil, ErrInvalidProtocolBinding(p.Proto.String())
	***REMOVED***
***REMOVED***

// GetCopy returns a copy of this PortBinding structure instance
func (p *PortBinding) GetCopy() PortBinding ***REMOVED***
	return PortBinding***REMOVED***
		Proto:       p.Proto,
		IP:          GetIPCopy(p.IP),
		Port:        p.Port,
		HostIP:      GetIPCopy(p.HostIP),
		HostPort:    p.HostPort,
		HostPortEnd: p.HostPortEnd,
	***REMOVED***
***REMOVED***

// String returns the PortBinding structure in string form
func (p *PortBinding) String() string ***REMOVED***
	ret := fmt.Sprintf("%s/", p.Proto)
	if p.IP != nil ***REMOVED***
		ret += p.IP.String()
	***REMOVED***
	ret = fmt.Sprintf("%s:%d/", ret, p.Port)
	if p.HostIP != nil ***REMOVED***
		ret += p.HostIP.String()
	***REMOVED***
	ret = fmt.Sprintf("%s:%d", ret, p.HostPort)
	return ret
***REMOVED***

// FromString reads the TransportPort structure from string
func (p *PortBinding) FromString(s string) error ***REMOVED***
	ps := strings.Split(s, "/")
	if len(ps) != 3 ***REMOVED***
		return BadRequestErrorf("invalid format for port binding: %s", s)
	***REMOVED***

	p.Proto = ParseProtocol(ps[0])

	var err error
	if p.IP, p.Port, err = parseIPPort(ps[1]); err != nil ***REMOVED***
		return BadRequestErrorf("failed to parse Container IP/Port in port binding: %s", err.Error())
	***REMOVED***

	if p.HostIP, p.HostPort, err = parseIPPort(ps[2]); err != nil ***REMOVED***
		return BadRequestErrorf("failed to parse Host IP/Port in port binding: %s", err.Error())
	***REMOVED***

	return nil
***REMOVED***

func parseIPPort(s string) (net.IP, uint16, error) ***REMOVED***
	pp := strings.Split(s, ":")
	if len(pp) != 2 ***REMOVED***
		return nil, 0, BadRequestErrorf("invalid format: %s", s)
	***REMOVED***

	var ip net.IP
	if pp[0] != "" ***REMOVED***
		if ip = net.ParseIP(pp[0]); ip == nil ***REMOVED***
			return nil, 0, BadRequestErrorf("invalid ip: %s", pp[0])
		***REMOVED***
	***REMOVED***

	port, err := strconv.ParseUint(pp[1], 10, 16)
	if err != nil ***REMOVED***
		return nil, 0, BadRequestErrorf("invalid port: %s", pp[1])
	***REMOVED***

	return ip, uint16(port), nil
***REMOVED***

// Equal checks if this instance of PortBinding is equal to the passed one
func (p *PortBinding) Equal(o *PortBinding) bool ***REMOVED***
	if p == o ***REMOVED***
		return true
	***REMOVED***

	if o == nil ***REMOVED***
		return false
	***REMOVED***

	if p.Proto != o.Proto || p.Port != o.Port ||
		p.HostPort != o.HostPort || p.HostPortEnd != o.HostPortEnd ***REMOVED***
		return false
	***REMOVED***

	if p.IP != nil ***REMOVED***
		if !p.IP.Equal(o.IP) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if o.IP != nil ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	if p.HostIP != nil ***REMOVED***
		if !p.HostIP.Equal(o.HostIP) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if o.HostIP != nil ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

// ErrInvalidProtocolBinding is returned when the port binding protocol is not valid.
type ErrInvalidProtocolBinding string

func (ipb ErrInvalidProtocolBinding) Error() string ***REMOVED***
	return fmt.Sprintf("invalid transport protocol: %s", string(ipb))
***REMOVED***

const (
	// ICMP is for the ICMP ip protocol
	ICMP = 1
	// TCP is for the TCP ip protocol
	TCP = 6
	// UDP is for the UDP ip protocol
	UDP = 17
)

// Protocol represents an IP protocol number
type Protocol uint8

func (p Protocol) String() string ***REMOVED***
	switch p ***REMOVED***
	case ICMP:
		return "icmp"
	case TCP:
		return "tcp"
	case UDP:
		return "udp"
	default:
		return fmt.Sprintf("%d", p)
	***REMOVED***
***REMOVED***

// ParseProtocol returns the respective Protocol type for the passed string
func ParseProtocol(s string) Protocol ***REMOVED***
	switch strings.ToLower(s) ***REMOVED***
	case "icmp":
		return ICMP
	case "udp":
		return UDP
	case "tcp":
		return TCP
	default:
		return 0
	***REMOVED***
***REMOVED***

// GetMacCopy returns a copy of the passed MAC address
func GetMacCopy(from net.HardwareAddr) net.HardwareAddr ***REMOVED***
	if from == nil ***REMOVED***
		return nil
	***REMOVED***
	to := make(net.HardwareAddr, len(from))
	copy(to, from)
	return to
***REMOVED***

// GetIPCopy returns a copy of the passed IP address
func GetIPCopy(from net.IP) net.IP ***REMOVED***
	if from == nil ***REMOVED***
		return nil
	***REMOVED***
	to := make(net.IP, len(from))
	copy(to, from)
	return to
***REMOVED***

// GetIPNetCopy returns a copy of the passed IP Network
func GetIPNetCopy(from *net.IPNet) *net.IPNet ***REMOVED***
	if from == nil ***REMOVED***
		return nil
	***REMOVED***
	bm := make(net.IPMask, len(from.Mask))
	copy(bm, from.Mask)
	return &net.IPNet***REMOVED***IP: GetIPCopy(from.IP), Mask: bm***REMOVED***
***REMOVED***

// GetIPNetCanonical returns the canonical form for the passed network
func GetIPNetCanonical(nw *net.IPNet) *net.IPNet ***REMOVED***
	if nw == nil ***REMOVED***
		return nil
	***REMOVED***
	c := GetIPNetCopy(nw)
	c.IP = c.IP.Mask(nw.Mask)
	return c
***REMOVED***

// CompareIPNet returns equal if the two IP Networks are equal
func CompareIPNet(a, b *net.IPNet) bool ***REMOVED***
	if a == b ***REMOVED***
		return true
	***REMOVED***
	if a == nil || b == nil ***REMOVED***
		return false
	***REMOVED***
	return a.IP.Equal(b.IP) && bytes.Equal(a.Mask, b.Mask)
***REMOVED***

// GetMinimalIP returns the address in its shortest form
func GetMinimalIP(ip net.IP) net.IP ***REMOVED***
	if ip != nil && ip.To4() != nil ***REMOVED***
		return ip.To4()
	***REMOVED***
	return ip
***REMOVED***

// GetMinimalIPNet returns a copy of the passed IP Network with congruent ip and mask notation
func GetMinimalIPNet(nw *net.IPNet) *net.IPNet ***REMOVED***
	if nw == nil ***REMOVED***
		return nil
	***REMOVED***
	if len(nw.IP) == 16 && nw.IP.To4() != nil ***REMOVED***
		m := nw.Mask
		if len(m) == 16 ***REMOVED***
			m = m[12:16]
		***REMOVED***
		return &net.IPNet***REMOVED***IP: nw.IP.To4(), Mask: m***REMOVED***
	***REMOVED***
	return nw
***REMOVED***

// IsIPNetValid returns true if the ipnet is a valid network/mask
// combination. Otherwise returns false.
func IsIPNetValid(nw *net.IPNet) bool ***REMOVED***
	return nw.String() != "0.0.0.0/0"
***REMOVED***

var v4inV6MaskPrefix = []byte***REMOVED***0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff***REMOVED***

// compareIPMask checks if the passed ip and mask are semantically compatible.
// It returns the byte indexes for the address and mask so that caller can
// do bitwise operations without modifying address representation.
func compareIPMask(ip net.IP, mask net.IPMask) (is int, ms int, err error) ***REMOVED***
	// Find the effective starting of address and mask
	if len(ip) == net.IPv6len && ip.To4() != nil ***REMOVED***
		is = 12
	***REMOVED***
	if len(ip[is:]) == net.IPv4len && len(mask) == net.IPv6len && bytes.Equal(mask[:12], v4inV6MaskPrefix) ***REMOVED***
		ms = 12
	***REMOVED***
	// Check if address and mask are semantically compatible
	if len(ip[is:]) != len(mask[ms:]) ***REMOVED***
		err = fmt.Errorf("ip and mask are not compatible: (%#v, %#v)", ip, mask)
	***REMOVED***
	return
***REMOVED***

// GetHostPartIP returns the host portion of the ip address identified by the mask.
// IP address representation is not modified. If address and mask are not compatible
// an error is returned.
func GetHostPartIP(ip net.IP, mask net.IPMask) (net.IP, error) ***REMOVED***
	// Find the effective starting of address and mask
	is, ms, err := compareIPMask(ip, mask)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("cannot compute host portion ip address because %s", err)
	***REMOVED***

	// Compute host portion
	out := GetIPCopy(ip)
	for i := 0; i < len(mask[ms:]); i++ ***REMOVED***
		out[is+i] &= ^mask[ms+i]
	***REMOVED***

	return out, nil
***REMOVED***

// GetBroadcastIP returns the broadcast ip address for the passed network (ip and mask).
// IP address representation is not modified. If address and mask are not compatible
// an error is returned.
func GetBroadcastIP(ip net.IP, mask net.IPMask) (net.IP, error) ***REMOVED***
	// Find the effective starting of address and mask
	is, ms, err := compareIPMask(ip, mask)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("cannot compute broadcast ip address because %s", err)
	***REMOVED***

	// Compute broadcast address
	out := GetIPCopy(ip)
	for i := 0; i < len(mask[ms:]); i++ ***REMOVED***
		out[is+i] |= ^mask[ms+i]
	***REMOVED***

	return out, nil
***REMOVED***

// ParseCIDR returns the *net.IPNet represented by the passed CIDR notation
func ParseCIDR(cidr string) (n *net.IPNet, e error) ***REMOVED***
	var i net.IP
	if i, n, e = net.ParseCIDR(cidr); e == nil ***REMOVED***
		n.IP = i
	***REMOVED***
	return
***REMOVED***

const (
	// NEXTHOP indicates a StaticRoute with an IP next hop.
	NEXTHOP = iota

	// CONNECTED indicates a StaticRoute with an interface for directly connected peers.
	CONNECTED
)

// StaticRoute is a statically-provisioned IP route.
type StaticRoute struct ***REMOVED***
	Destination *net.IPNet

	RouteType int // NEXT_HOP or CONNECTED

	// NextHop will be resolved by the kernel (i.e. as a loose hop).
	NextHop net.IP
***REMOVED***

// GetCopy returns a copy of this StaticRoute structure
func (r *StaticRoute) GetCopy() *StaticRoute ***REMOVED***
	d := GetIPNetCopy(r.Destination)
	nh := GetIPCopy(r.NextHop)
	return &StaticRoute***REMOVED***Destination: d,
		RouteType: r.RouteType,
		NextHop:   nh,
	***REMOVED***
***REMOVED***

// InterfaceStatistics represents the interface's statistics
type InterfaceStatistics struct ***REMOVED***
	RxBytes   uint64
	RxPackets uint64
	RxErrors  uint64
	RxDropped uint64
	TxBytes   uint64
	TxPackets uint64
	TxErrors  uint64
	TxDropped uint64
***REMOVED***

func (is *InterfaceStatistics) String() string ***REMOVED***
	return fmt.Sprintf("\nRxBytes: %d, RxPackets: %d, RxErrors: %d, RxDropped: %d, TxBytes: %d, TxPackets: %d, TxErrors: %d, TxDropped: %d",
		is.RxBytes, is.RxPackets, is.RxErrors, is.RxDropped, is.TxBytes, is.TxPackets, is.TxErrors, is.TxDropped)
***REMOVED***

/******************************
 * Well-known Error Interfaces
 ******************************/

// MaskableError is an interface for errors which can be ignored by caller
type MaskableError interface ***REMOVED***
	// Maskable makes implementer into MaskableError type
	Maskable()
***REMOVED***

// RetryError is an interface for errors which might get resolved through retry
type RetryError interface ***REMOVED***
	// Retry makes implementer into RetryError type
	Retry()
***REMOVED***

// BadRequestError is an interface for errors originated by a bad request
type BadRequestError interface ***REMOVED***
	// BadRequest makes implementer into BadRequestError type
	BadRequest()
***REMOVED***

// NotFoundError is an interface for errors raised because a needed resource is not available
type NotFoundError interface ***REMOVED***
	// NotFound makes implementer into NotFoundError type
	NotFound()
***REMOVED***

// ForbiddenError is an interface for errors which denote a valid request that cannot be honored
type ForbiddenError interface ***REMOVED***
	// Forbidden makes implementer into ForbiddenError type
	Forbidden()
***REMOVED***

// NoServiceError is an interface for errors returned when the required service is not available
type NoServiceError interface ***REMOVED***
	// NoService makes implementer into NoServiceError type
	NoService()
***REMOVED***

// TimeoutError is an interface for errors raised because of timeout
type TimeoutError interface ***REMOVED***
	// Timeout makes implementer into TimeoutError type
	Timeout()
***REMOVED***

// NotImplementedError is an interface for errors raised because of requested functionality is not yet implemented
type NotImplementedError interface ***REMOVED***
	// NotImplemented makes implementer into NotImplementedError type
	NotImplemented()
***REMOVED***

// InternalError is an interface for errors raised because of an internal error
type InternalError interface ***REMOVED***
	// Internal makes implementer into InternalError type
	Internal()
***REMOVED***

/******************************
 * Well-known Error Formatters
 ******************************/

// BadRequestErrorf creates an instance of BadRequestError
func BadRequestErrorf(format string, params ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return badRequest(fmt.Sprintf(format, params...))
***REMOVED***

// NotFoundErrorf creates an instance of NotFoundError
func NotFoundErrorf(format string, params ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return notFound(fmt.Sprintf(format, params...))
***REMOVED***

// ForbiddenErrorf creates an instance of ForbiddenError
func ForbiddenErrorf(format string, params ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return forbidden(fmt.Sprintf(format, params...))
***REMOVED***

// NoServiceErrorf creates an instance of NoServiceError
func NoServiceErrorf(format string, params ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return noService(fmt.Sprintf(format, params...))
***REMOVED***

// NotImplementedErrorf creates an instance of NotImplementedError
func NotImplementedErrorf(format string, params ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return notImpl(fmt.Sprintf(format, params...))
***REMOVED***

// TimeoutErrorf creates an instance of TimeoutError
func TimeoutErrorf(format string, params ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return timeout(fmt.Sprintf(format, params...))
***REMOVED***

// InternalErrorf creates an instance of InternalError
func InternalErrorf(format string, params ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return internal(fmt.Sprintf(format, params...))
***REMOVED***

// InternalMaskableErrorf creates an instance of InternalError and MaskableError
func InternalMaskableErrorf(format string, params ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return maskInternal(fmt.Sprintf(format, params...))
***REMOVED***

// RetryErrorf creates an instance of RetryError
func RetryErrorf(format string, params ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return retry(fmt.Sprintf(format, params...))
***REMOVED***

/***********************
 * Internal Error Types
 ***********************/
type badRequest string

func (br badRequest) Error() string ***REMOVED***
	return string(br)
***REMOVED***
func (br badRequest) BadRequest() ***REMOVED******REMOVED***

type maskBadRequest string

type notFound string

func (nf notFound) Error() string ***REMOVED***
	return string(nf)
***REMOVED***
func (nf notFound) NotFound() ***REMOVED******REMOVED***

type forbidden string

func (frb forbidden) Error() string ***REMOVED***
	return string(frb)
***REMOVED***
func (frb forbidden) Forbidden() ***REMOVED******REMOVED***

type noService string

func (ns noService) Error() string ***REMOVED***
	return string(ns)
***REMOVED***
func (ns noService) NoService() ***REMOVED******REMOVED***

type maskNoService string

type timeout string

func (to timeout) Error() string ***REMOVED***
	return string(to)
***REMOVED***
func (to timeout) Timeout() ***REMOVED******REMOVED***

type notImpl string

func (ni notImpl) Error() string ***REMOVED***
	return string(ni)
***REMOVED***
func (ni notImpl) NotImplemented() ***REMOVED******REMOVED***

type internal string

func (nt internal) Error() string ***REMOVED***
	return string(nt)
***REMOVED***
func (nt internal) Internal() ***REMOVED******REMOVED***

type maskInternal string

func (mnt maskInternal) Error() string ***REMOVED***
	return string(mnt)
***REMOVED***
func (mnt maskInternal) Internal() ***REMOVED******REMOVED***
func (mnt maskInternal) Maskable() ***REMOVED******REMOVED***

type retry string

func (r retry) Error() string ***REMOVED***
	return string(r)
***REMOVED***
func (r retry) Retry() ***REMOVED******REMOVED***
