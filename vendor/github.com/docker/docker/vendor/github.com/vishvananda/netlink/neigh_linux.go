package netlink

import (
	"net"
	"syscall"
	"unsafe"

	"github.com/vishvananda/netlink/nl"
)

const (
	NDA_UNSPEC = iota
	NDA_DST
	NDA_LLADDR
	NDA_CACHEINFO
	NDA_PROBES
	NDA_VLAN
	NDA_PORT
	NDA_VNI
	NDA_IFINDEX
	NDA_MAX = NDA_IFINDEX
)

// Neighbor Cache Entry States.
const (
	NUD_NONE       = 0x00
	NUD_INCOMPLETE = 0x01
	NUD_REACHABLE  = 0x02
	NUD_STALE      = 0x04
	NUD_DELAY      = 0x08
	NUD_PROBE      = 0x10
	NUD_FAILED     = 0x20
	NUD_NOARP      = 0x40
	NUD_PERMANENT  = 0x80
)

// Neighbor Flags
const (
	NTF_USE    = 0x01
	NTF_SELF   = 0x02
	NTF_MASTER = 0x04
	NTF_PROXY  = 0x08
	NTF_ROUTER = 0x80
)

type Ndmsg struct ***REMOVED***
	Family uint8
	Index  uint32
	State  uint16
	Flags  uint8
	Type   uint8
***REMOVED***

func deserializeNdmsg(b []byte) *Ndmsg ***REMOVED***
	var dummy Ndmsg
	return (*Ndmsg)(unsafe.Pointer(&b[0:unsafe.Sizeof(dummy)][0]))
***REMOVED***

func (msg *Ndmsg) Serialize() []byte ***REMOVED***
	return (*(*[unsafe.Sizeof(*msg)]byte)(unsafe.Pointer(msg)))[:]
***REMOVED***

func (msg *Ndmsg) Len() int ***REMOVED***
	return int(unsafe.Sizeof(*msg))
***REMOVED***

// NeighAdd will add an IP to MAC mapping to the ARP table
// Equivalent to: `ip neigh add ....`
func NeighAdd(neigh *Neigh) error ***REMOVED***
	return pkgHandle.NeighAdd(neigh)
***REMOVED***

// NeighAdd will add an IP to MAC mapping to the ARP table
// Equivalent to: `ip neigh add ....`
func (h *Handle) NeighAdd(neigh *Neigh) error ***REMOVED***
	return h.neighAdd(neigh, syscall.NLM_F_CREATE|syscall.NLM_F_EXCL)
***REMOVED***

// NeighSet will add or replace an IP to MAC mapping to the ARP table
// Equivalent to: `ip neigh replace....`
func NeighSet(neigh *Neigh) error ***REMOVED***
	return pkgHandle.NeighSet(neigh)
***REMOVED***

// NeighSet will add or replace an IP to MAC mapping to the ARP table
// Equivalent to: `ip neigh replace....`
func (h *Handle) NeighSet(neigh *Neigh) error ***REMOVED***
	return h.neighAdd(neigh, syscall.NLM_F_CREATE|syscall.NLM_F_REPLACE)
***REMOVED***

// NeighAppend will append an entry to FDB
// Equivalent to: `bridge fdb append...`
func NeighAppend(neigh *Neigh) error ***REMOVED***
	return pkgHandle.NeighAppend(neigh)
***REMOVED***

// NeighAppend will append an entry to FDB
// Equivalent to: `bridge fdb append...`
func (h *Handle) NeighAppend(neigh *Neigh) error ***REMOVED***
	return h.neighAdd(neigh, syscall.NLM_F_CREATE|syscall.NLM_F_APPEND)
***REMOVED***

// NeighAppend will append an entry to FDB
// Equivalent to: `bridge fdb append...`
func neighAdd(neigh *Neigh, mode int) error ***REMOVED***
	return pkgHandle.neighAdd(neigh, mode)
***REMOVED***

// NeighAppend will append an entry to FDB
// Equivalent to: `bridge fdb append...`
func (h *Handle) neighAdd(neigh *Neigh, mode int) error ***REMOVED***
	req := h.newNetlinkRequest(syscall.RTM_NEWNEIGH, mode|syscall.NLM_F_ACK)
	return neighHandle(neigh, req)
***REMOVED***

// NeighDel will delete an IP address from a link device.
// Equivalent to: `ip addr del $addr dev $link`
func NeighDel(neigh *Neigh) error ***REMOVED***
	return pkgHandle.NeighDel(neigh)
***REMOVED***

// NeighDel will delete an IP address from a link device.
// Equivalent to: `ip addr del $addr dev $link`
func (h *Handle) NeighDel(neigh *Neigh) error ***REMOVED***
	req := h.newNetlinkRequest(syscall.RTM_DELNEIGH, syscall.NLM_F_ACK)
	return neighHandle(neigh, req)
***REMOVED***

func neighHandle(neigh *Neigh, req *nl.NetlinkRequest) error ***REMOVED***
	var family int

	if neigh.Family > 0 ***REMOVED***
		family = neigh.Family
	***REMOVED*** else ***REMOVED***
		family = nl.GetIPFamily(neigh.IP)
	***REMOVED***

	msg := Ndmsg***REMOVED***
		Family: uint8(family),
		Index:  uint32(neigh.LinkIndex),
		State:  uint16(neigh.State),
		Type:   uint8(neigh.Type),
		Flags:  uint8(neigh.Flags),
	***REMOVED***
	req.AddData(&msg)

	ipData := neigh.IP.To4()
	if ipData == nil ***REMOVED***
		ipData = neigh.IP.To16()
	***REMOVED***

	dstData := nl.NewRtAttr(NDA_DST, ipData)
	req.AddData(dstData)

	if neigh.LLIPAddr != nil ***REMOVED***
		llIPData := nl.NewRtAttr(NDA_LLADDR, neigh.LLIPAddr.To4())
		req.AddData(llIPData)
	***REMOVED*** else if neigh.Flags != NTF_PROXY || neigh.HardwareAddr != nil ***REMOVED***
		hwData := nl.NewRtAttr(NDA_LLADDR, []byte(neigh.HardwareAddr))
		req.AddData(hwData)
	***REMOVED***

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

// NeighList gets a list of IP-MAC mappings in the system (ARP table).
// Equivalent to: `ip neighbor show`.
// The list can be filtered by link and ip family.
func NeighList(linkIndex, family int) ([]Neigh, error) ***REMOVED***
	return pkgHandle.NeighList(linkIndex, family)
***REMOVED***

// NeighProxyList gets a list of neighbor proxies in the system.
// Equivalent to: `ip neighbor show proxy`.
// The list can be filtered by link and ip family.
func NeighProxyList(linkIndex, family int) ([]Neigh, error) ***REMOVED***
	return pkgHandle.NeighProxyList(linkIndex, family)
***REMOVED***

// NeighList gets a list of IP-MAC mappings in the system (ARP table).
// Equivalent to: `ip neighbor show`.
// The list can be filtered by link and ip family.
func (h *Handle) NeighList(linkIndex, family int) ([]Neigh, error) ***REMOVED***
	return h.neighList(linkIndex, family, 0)
***REMOVED***

// NeighProxyList gets a list of neighbor proxies in the system.
// Equivalent to: `ip neighbor show proxy`.
// The list can be filtered by link, ip family.
func (h *Handle) NeighProxyList(linkIndex, family int) ([]Neigh, error) ***REMOVED***
	return h.neighList(linkIndex, family, NTF_PROXY)
***REMOVED***

func (h *Handle) neighList(linkIndex, family, flags int) ([]Neigh, error) ***REMOVED***
	req := h.newNetlinkRequest(syscall.RTM_GETNEIGH, syscall.NLM_F_DUMP)
	msg := Ndmsg***REMOVED***
		Family: uint8(family),
		Index:  uint32(linkIndex),
		Flags:  uint8(flags),
	***REMOVED***
	req.AddData(&msg)

	msgs, err := req.Execute(syscall.NETLINK_ROUTE, syscall.RTM_NEWNEIGH)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var res []Neigh
	for _, m := range msgs ***REMOVED***
		ndm := deserializeNdmsg(m)
		if linkIndex != 0 && int(ndm.Index) != linkIndex ***REMOVED***
			// Ignore messages from other interfaces
			continue
		***REMOVED***

		neigh, err := NeighDeserialize(m)
		if err != nil ***REMOVED***
			continue
		***REMOVED***

		res = append(res, *neigh)
	***REMOVED***

	return res, nil
***REMOVED***

func NeighDeserialize(m []byte) (*Neigh, error) ***REMOVED***
	msg := deserializeNdmsg(m)

	neigh := Neigh***REMOVED***
		LinkIndex: int(msg.Index),
		Family:    int(msg.Family),
		State:     int(msg.State),
		Type:      int(msg.Type),
		Flags:     int(msg.Flags),
	***REMOVED***

	attrs, err := nl.ParseRouteAttr(m[msg.Len():])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// This should be cached for perfomance
	// once per table dump
	link, err := LinkByIndex(neigh.LinkIndex)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	encapType := link.Attrs().EncapType

	for _, attr := range attrs ***REMOVED***
		switch attr.Attr.Type ***REMOVED***
		case NDA_DST:
			neigh.IP = net.IP(attr.Value)
		case NDA_LLADDR:
			// BUG: Is this a bug in the netlink library?
			// #define RTA_LENGTH(len) (RTA_ALIGN(sizeof(struct rtattr)) + (len))
			// #define RTA_PAYLOAD(rta) ((int)((rta)->rta_len) - RTA_LENGTH(0))
			attrLen := attr.Attr.Len - syscall.SizeofRtAttr
			if attrLen == 4 && (encapType == "ipip" ||
				encapType == "sit" ||
				encapType == "gre") ***REMOVED***
				neigh.LLIPAddr = net.IP(attr.Value)
			***REMOVED*** else if attrLen == 16 &&
				encapType == "tunnel6" ***REMOVED***
				neigh.IP = net.IP(attr.Value)
			***REMOVED*** else ***REMOVED***
				neigh.HardwareAddr = net.HardwareAddr(attr.Value)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return &neigh, nil
***REMOVED***
