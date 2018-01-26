package netlink

import (
	"fmt"
	"net"
	"strings"
	"syscall"

	"github.com/vishvananda/netlink/nl"
	"github.com/vishvananda/netns"
)

// IFA_FLAGS is a u32 attribute.
const IFA_FLAGS = 0x8

// AddrAdd will add an IP address to a link device.
// Equivalent to: `ip addr add $addr dev $link`
func AddrAdd(link Link, addr *Addr) error ***REMOVED***
	return pkgHandle.AddrAdd(link, addr)
***REMOVED***

// AddrAdd will add an IP address to a link device.
// Equivalent to: `ip addr add $addr dev $link`
func (h *Handle) AddrAdd(link Link, addr *Addr) error ***REMOVED***
	req := h.newNetlinkRequest(syscall.RTM_NEWADDR, syscall.NLM_F_CREATE|syscall.NLM_F_EXCL|syscall.NLM_F_ACK)
	return h.addrHandle(link, addr, req)
***REMOVED***

// AddrReplace will replace (or, if not present, add) an IP address on a link device.
// Equivalent to: `ip addr replace $addr dev $link`
func AddrReplace(link Link, addr *Addr) error ***REMOVED***
	return pkgHandle.AddrReplace(link, addr)
***REMOVED***

// AddrReplace will replace (or, if not present, add) an IP address on a link device.
// Equivalent to: `ip addr replace $addr dev $link`
func (h *Handle) AddrReplace(link Link, addr *Addr) error ***REMOVED***
	req := h.newNetlinkRequest(syscall.RTM_NEWADDR, syscall.NLM_F_CREATE|syscall.NLM_F_REPLACE|syscall.NLM_F_ACK)
	return h.addrHandle(link, addr, req)
***REMOVED***

// AddrDel will delete an IP address from a link device.
// Equivalent to: `ip addr del $addr dev $link`
func AddrDel(link Link, addr *Addr) error ***REMOVED***
	return pkgHandle.AddrDel(link, addr)
***REMOVED***

// AddrDel will delete an IP address from a link device.
// Equivalent to: `ip addr del $addr dev $link`
func (h *Handle) AddrDel(link Link, addr *Addr) error ***REMOVED***
	req := h.newNetlinkRequest(syscall.RTM_DELADDR, syscall.NLM_F_ACK)
	return h.addrHandle(link, addr, req)
***REMOVED***

func (h *Handle) addrHandle(link Link, addr *Addr, req *nl.NetlinkRequest) error ***REMOVED***
	base := link.Attrs()
	if addr.Label != "" && !strings.HasPrefix(addr.Label, base.Name) ***REMOVED***
		return fmt.Errorf("label must begin with interface name")
	***REMOVED***
	h.ensureIndex(base)

	family := nl.GetIPFamily(addr.IP)

	msg := nl.NewIfAddrmsg(family)
	msg.Index = uint32(base.Index)
	msg.Scope = uint8(addr.Scope)
	prefixlen, masklen := addr.Mask.Size()
	msg.Prefixlen = uint8(prefixlen)
	req.AddData(msg)

	var localAddrData []byte
	if family == FAMILY_V4 ***REMOVED***
		localAddrData = addr.IP.To4()
	***REMOVED*** else ***REMOVED***
		localAddrData = addr.IP.To16()
	***REMOVED***

	localData := nl.NewRtAttr(syscall.IFA_LOCAL, localAddrData)
	req.AddData(localData)
	var peerAddrData []byte
	if addr.Peer != nil ***REMOVED***
		if family == FAMILY_V4 ***REMOVED***
			peerAddrData = addr.Peer.IP.To4()
		***REMOVED*** else ***REMOVED***
			peerAddrData = addr.Peer.IP.To16()
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		peerAddrData = localAddrData
	***REMOVED***

	addressData := nl.NewRtAttr(syscall.IFA_ADDRESS, peerAddrData)
	req.AddData(addressData)

	if addr.Flags != 0 ***REMOVED***
		if addr.Flags <= 0xff ***REMOVED***
			msg.IfAddrmsg.Flags = uint8(addr.Flags)
		***REMOVED*** else ***REMOVED***
			b := make([]byte, 4)
			native.PutUint32(b, uint32(addr.Flags))
			flagsData := nl.NewRtAttr(IFA_FLAGS, b)
			req.AddData(flagsData)
		***REMOVED***
	***REMOVED***

	if addr.Broadcast == nil ***REMOVED***
		calcBroadcast := make(net.IP, masklen/8)
		for i := range localAddrData ***REMOVED***
			calcBroadcast[i] = localAddrData[i] | ^addr.Mask[i]
		***REMOVED***
		addr.Broadcast = calcBroadcast
	***REMOVED***
	req.AddData(nl.NewRtAttr(syscall.IFA_BROADCAST, addr.Broadcast))

	if addr.Label != "" ***REMOVED***
		labelData := nl.NewRtAttr(syscall.IFA_LABEL, nl.ZeroTerminated(addr.Label))
		req.AddData(labelData)
	***REMOVED***

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

// AddrList gets a list of IP addresses in the system.
// Equivalent to: `ip addr show`.
// The list can be filtered by link and ip family.
func AddrList(link Link, family int) ([]Addr, error) ***REMOVED***
	return pkgHandle.AddrList(link, family)
***REMOVED***

// AddrList gets a list of IP addresses in the system.
// Equivalent to: `ip addr show`.
// The list can be filtered by link and ip family.
func (h *Handle) AddrList(link Link, family int) ([]Addr, error) ***REMOVED***
	req := h.newNetlinkRequest(syscall.RTM_GETADDR, syscall.NLM_F_DUMP)
	msg := nl.NewIfInfomsg(family)
	req.AddData(msg)

	msgs, err := req.Execute(syscall.NETLINK_ROUTE, syscall.RTM_NEWADDR)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	indexFilter := 0
	if link != nil ***REMOVED***
		base := link.Attrs()
		h.ensureIndex(base)
		indexFilter = base.Index
	***REMOVED***

	var res []Addr
	for _, m := range msgs ***REMOVED***
		addr, msgFamily, ifindex, err := parseAddr(m)
		if err != nil ***REMOVED***
			return res, err
		***REMOVED***

		if link != nil && ifindex != indexFilter ***REMOVED***
			// Ignore messages from other interfaces
			continue
		***REMOVED***

		if family != FAMILY_ALL && msgFamily != family ***REMOVED***
			continue
		***REMOVED***

		res = append(res, addr)
	***REMOVED***

	return res, nil
***REMOVED***

func parseAddr(m []byte) (addr Addr, family, index int, err error) ***REMOVED***
	msg := nl.DeserializeIfAddrmsg(m)

	family = -1
	index = -1

	attrs, err1 := nl.ParseRouteAttr(m[msg.Len():])
	if err1 != nil ***REMOVED***
		err = err1
		return
	***REMOVED***

	family = int(msg.Family)
	index = int(msg.Index)

	var local, dst *net.IPNet
	for _, attr := range attrs ***REMOVED***
		switch attr.Attr.Type ***REMOVED***
		case syscall.IFA_ADDRESS:
			dst = &net.IPNet***REMOVED***
				IP:   attr.Value,
				Mask: net.CIDRMask(int(msg.Prefixlen), 8*len(attr.Value)),
			***REMOVED***
			addr.Peer = dst
		case syscall.IFA_LOCAL:
			local = &net.IPNet***REMOVED***
				IP:   attr.Value,
				Mask: net.CIDRMask(int(msg.Prefixlen), 8*len(attr.Value)),
			***REMOVED***
			addr.IPNet = local
		case syscall.IFA_BROADCAST:
			addr.Broadcast = attr.Value
		case syscall.IFA_LABEL:
			addr.Label = string(attr.Value[:len(attr.Value)-1])
		case IFA_FLAGS:
			addr.Flags = int(native.Uint32(attr.Value[0:4]))
		case nl.IFA_CACHEINFO:
			ci := nl.DeserializeIfaCacheInfo(attr.Value)
			addr.PreferedLft = int(ci.IfaPrefered)
			addr.ValidLft = int(ci.IfaValid)
		***REMOVED***
	***REMOVED***

	// IFA_LOCAL should be there but if not, fall back to IFA_ADDRESS
	if local != nil ***REMOVED***
		addr.IPNet = local
	***REMOVED*** else ***REMOVED***
		addr.IPNet = dst
	***REMOVED***
	addr.Scope = int(msg.Scope)

	return
***REMOVED***

type AddrUpdate struct ***REMOVED***
	LinkAddress net.IPNet
	LinkIndex   int
	Flags       int
	Scope       int
	PreferedLft int
	ValidLft    int
	NewAddr     bool // true=added false=deleted
***REMOVED***

// AddrSubscribe takes a chan down which notifications will be sent
// when addresses change.  Close the 'done' chan to stop subscription.
func AddrSubscribe(ch chan<- AddrUpdate, done <-chan struct***REMOVED******REMOVED***) error ***REMOVED***
	return addrSubscribeAt(netns.None(), netns.None(), ch, done, nil)
***REMOVED***

// AddrSubscribeAt works like AddrSubscribe plus it allows the caller
// to choose the network namespace in which to subscribe (ns).
func AddrSubscribeAt(ns netns.NsHandle, ch chan<- AddrUpdate, done <-chan struct***REMOVED******REMOVED***) error ***REMOVED***
	return addrSubscribeAt(ns, netns.None(), ch, done, nil)
***REMOVED***

// AddrSubscribeOptions contains a set of options to use with
// AddrSubscribeWithOptions.
type AddrSubscribeOptions struct ***REMOVED***
	Namespace     *netns.NsHandle
	ErrorCallback func(error)
***REMOVED***

// AddrSubscribeWithOptions work like AddrSubscribe but enable to
// provide additional options to modify the behavior. Currently, the
// namespace can be provided as well as an error callback.
func AddrSubscribeWithOptions(ch chan<- AddrUpdate, done <-chan struct***REMOVED******REMOVED***, options AddrSubscribeOptions) error ***REMOVED***
	if options.Namespace == nil ***REMOVED***
		none := netns.None()
		options.Namespace = &none
	***REMOVED***
	return addrSubscribeAt(*options.Namespace, netns.None(), ch, done, options.ErrorCallback)
***REMOVED***

func addrSubscribeAt(newNs, curNs netns.NsHandle, ch chan<- AddrUpdate, done <-chan struct***REMOVED******REMOVED***, cberr func(error)) error ***REMOVED***
	s, err := nl.SubscribeAt(newNs, curNs, syscall.NETLINK_ROUTE, syscall.RTNLGRP_IPV4_IFADDR, syscall.RTNLGRP_IPV6_IFADDR)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if done != nil ***REMOVED***
		go func() ***REMOVED***
			<-done
			s.Close()
		***REMOVED***()
	***REMOVED***
	go func() ***REMOVED***
		defer close(ch)
		for ***REMOVED***
			msgs, err := s.Receive()
			if err != nil ***REMOVED***
				if cberr != nil ***REMOVED***
					cberr(err)
				***REMOVED***
				return
			***REMOVED***
			for _, m := range msgs ***REMOVED***
				msgType := m.Header.Type
				if msgType != syscall.RTM_NEWADDR && msgType != syscall.RTM_DELADDR ***REMOVED***
					if cberr != nil ***REMOVED***
						cberr(fmt.Errorf("bad message type: %d", msgType))
					***REMOVED***
					return
				***REMOVED***

				addr, _, ifindex, err := parseAddr(m.Data)
				if err != nil ***REMOVED***
					if cberr != nil ***REMOVED***
						cberr(fmt.Errorf("could not parse address: %v", err))
					***REMOVED***
					return
				***REMOVED***

				ch <- AddrUpdate***REMOVED***LinkAddress: *addr.IPNet,
					LinkIndex:   ifindex,
					NewAddr:     msgType == syscall.RTM_NEWADDR,
					Flags:       addr.Flags,
					Scope:       addr.Scope,
					PreferedLft: addr.PreferedLft,
					ValidLft:    addr.ValidLft***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return nil
***REMOVED***
