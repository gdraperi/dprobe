package netlink

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"syscall"
	"unsafe"

	"github.com/vishvananda/netlink/nl"
	"github.com/vishvananda/netns"
)

const (
	SizeofLinkStats32 = 0x5c
	SizeofLinkStats64 = 0xd8
	IFLA_STATS64      = 0x17 // syscall pkg does not contain this one
)

const (
	TUNTAP_MODE_TUN  TuntapMode = syscall.IFF_TUN
	TUNTAP_MODE_TAP  TuntapMode = syscall.IFF_TAP
	TUNTAP_DEFAULTS  TuntapFlag = syscall.IFF_TUN_EXCL | syscall.IFF_ONE_QUEUE
	TUNTAP_VNET_HDR  TuntapFlag = syscall.IFF_VNET_HDR
	TUNTAP_TUN_EXCL  TuntapFlag = syscall.IFF_TUN_EXCL
	TUNTAP_NO_PI     TuntapFlag = syscall.IFF_NO_PI
	TUNTAP_ONE_QUEUE TuntapFlag = syscall.IFF_ONE_QUEUE
)

var lookupByDump = false

var macvlanModes = [...]uint32***REMOVED***
	0,
	nl.MACVLAN_MODE_PRIVATE,
	nl.MACVLAN_MODE_VEPA,
	nl.MACVLAN_MODE_BRIDGE,
	nl.MACVLAN_MODE_PASSTHRU,
	nl.MACVLAN_MODE_SOURCE,
***REMOVED***

func ensureIndex(link *LinkAttrs) ***REMOVED***
	if link != nil && link.Index == 0 ***REMOVED***
		newlink, _ := LinkByName(link.Name)
		if newlink != nil ***REMOVED***
			link.Index = newlink.Attrs().Index
		***REMOVED***
	***REMOVED***
***REMOVED***

func (h *Handle) ensureIndex(link *LinkAttrs) ***REMOVED***
	if link != nil && link.Index == 0 ***REMOVED***
		newlink, _ := h.LinkByName(link.Name)
		if newlink != nil ***REMOVED***
			link.Index = newlink.Attrs().Index
		***REMOVED***
	***REMOVED***
***REMOVED***

func (h *Handle) LinkSetARPOff(link Link) error ***REMOVED***
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(syscall.RTM_SETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Change |= syscall.IFF_NOARP
	msg.Flags |= syscall.IFF_NOARP
	msg.Index = int32(base.Index)
	req.AddData(msg)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

func LinkSetARPOff(link Link) error ***REMOVED***
	return pkgHandle.LinkSetARPOff(link)
***REMOVED***

func (h *Handle) LinkSetARPOn(link Link) error ***REMOVED***
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(syscall.RTM_SETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Change |= syscall.IFF_NOARP
	msg.Flags &= ^uint32(syscall.IFF_NOARP)
	msg.Index = int32(base.Index)
	req.AddData(msg)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

func LinkSetARPOn(link Link) error ***REMOVED***
	return pkgHandle.LinkSetARPOn(link)
***REMOVED***

func (h *Handle) SetPromiscOn(link Link) error ***REMOVED***
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(syscall.RTM_SETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Change = syscall.IFF_PROMISC
	msg.Flags = syscall.IFF_PROMISC
	msg.Index = int32(base.Index)
	req.AddData(msg)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

func BridgeSetMcastSnoop(link Link, on bool) error ***REMOVED***
	return pkgHandle.BridgeSetMcastSnoop(link, on)
***REMOVED***

func (h *Handle) BridgeSetMcastSnoop(link Link, on bool) error ***REMOVED***
	bridge := link.(*Bridge)
	bridge.MulticastSnooping = &on
	return h.linkModify(bridge, syscall.NLM_F_ACK)
***REMOVED***

func SetPromiscOn(link Link) error ***REMOVED***
	return pkgHandle.SetPromiscOn(link)
***REMOVED***

func (h *Handle) SetPromiscOff(link Link) error ***REMOVED***
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(syscall.RTM_SETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Change = syscall.IFF_PROMISC
	msg.Flags = 0 & ^syscall.IFF_PROMISC
	msg.Index = int32(base.Index)
	req.AddData(msg)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

func SetPromiscOff(link Link) error ***REMOVED***
	return pkgHandle.SetPromiscOff(link)
***REMOVED***

// LinkSetUp enables the link device.
// Equivalent to: `ip link set $link up`
func LinkSetUp(link Link) error ***REMOVED***
	return pkgHandle.LinkSetUp(link)
***REMOVED***

// LinkSetUp enables the link device.
// Equivalent to: `ip link set $link up`
func (h *Handle) LinkSetUp(link Link) error ***REMOVED***
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(syscall.RTM_NEWLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Change = syscall.IFF_UP
	msg.Flags = syscall.IFF_UP
	msg.Index = int32(base.Index)
	req.AddData(msg)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

// LinkSetDown disables link device.
// Equivalent to: `ip link set $link down`
func LinkSetDown(link Link) error ***REMOVED***
	return pkgHandle.LinkSetDown(link)
***REMOVED***

// LinkSetDown disables link device.
// Equivalent to: `ip link set $link down`
func (h *Handle) LinkSetDown(link Link) error ***REMOVED***
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(syscall.RTM_NEWLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Change = syscall.IFF_UP
	msg.Flags = 0 & ^syscall.IFF_UP
	msg.Index = int32(base.Index)
	req.AddData(msg)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

// LinkSetMTU sets the mtu of the link device.
// Equivalent to: `ip link set $link mtu $mtu`
func LinkSetMTU(link Link, mtu int) error ***REMOVED***
	return pkgHandle.LinkSetMTU(link, mtu)
***REMOVED***

// LinkSetMTU sets the mtu of the link device.
// Equivalent to: `ip link set $link mtu $mtu`
func (h *Handle) LinkSetMTU(link Link, mtu int) error ***REMOVED***
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(syscall.RTM_SETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Index = int32(base.Index)
	req.AddData(msg)

	b := make([]byte, 4)
	native.PutUint32(b, uint32(mtu))

	data := nl.NewRtAttr(syscall.IFLA_MTU, b)
	req.AddData(data)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

// LinkSetName sets the name of the link device.
// Equivalent to: `ip link set $link name $name`
func LinkSetName(link Link, name string) error ***REMOVED***
	return pkgHandle.LinkSetName(link, name)
***REMOVED***

// LinkSetName sets the name of the link device.
// Equivalent to: `ip link set $link name $name`
func (h *Handle) LinkSetName(link Link, name string) error ***REMOVED***
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(syscall.RTM_SETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Index = int32(base.Index)
	req.AddData(msg)

	data := nl.NewRtAttr(syscall.IFLA_IFNAME, []byte(name))
	req.AddData(data)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

// LinkSetAlias sets the alias of the link device.
// Equivalent to: `ip link set dev $link alias $name`
func LinkSetAlias(link Link, name string) error ***REMOVED***
	return pkgHandle.LinkSetAlias(link, name)
***REMOVED***

// LinkSetAlias sets the alias of the link device.
// Equivalent to: `ip link set dev $link alias $name`
func (h *Handle) LinkSetAlias(link Link, name string) error ***REMOVED***
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(syscall.RTM_SETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Index = int32(base.Index)
	req.AddData(msg)

	data := nl.NewRtAttr(syscall.IFLA_IFALIAS, []byte(name))
	req.AddData(data)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

// LinkSetHardwareAddr sets the hardware address of the link device.
// Equivalent to: `ip link set $link address $hwaddr`
func LinkSetHardwareAddr(link Link, hwaddr net.HardwareAddr) error ***REMOVED***
	return pkgHandle.LinkSetHardwareAddr(link, hwaddr)
***REMOVED***

// LinkSetHardwareAddr sets the hardware address of the link device.
// Equivalent to: `ip link set $link address $hwaddr`
func (h *Handle) LinkSetHardwareAddr(link Link, hwaddr net.HardwareAddr) error ***REMOVED***
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(syscall.RTM_SETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Index = int32(base.Index)
	req.AddData(msg)

	data := nl.NewRtAttr(syscall.IFLA_ADDRESS, []byte(hwaddr))
	req.AddData(data)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

// LinkSetVfHardwareAddr sets the hardware address of a vf for the link.
// Equivalent to: `ip link set $link vf $vf mac $hwaddr`
func LinkSetVfHardwareAddr(link Link, vf int, hwaddr net.HardwareAddr) error ***REMOVED***
	return pkgHandle.LinkSetVfHardwareAddr(link, vf, hwaddr)
***REMOVED***

// LinkSetVfHardwareAddr sets the hardware address of a vf for the link.
// Equivalent to: `ip link set $link vf $vf mac $hwaddr`
func (h *Handle) LinkSetVfHardwareAddr(link Link, vf int, hwaddr net.HardwareAddr) error ***REMOVED***
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(syscall.RTM_SETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Index = int32(base.Index)
	req.AddData(msg)

	data := nl.NewRtAttr(nl.IFLA_VFINFO_LIST, nil)
	info := nl.NewRtAttrChild(data, nl.IFLA_VF_INFO, nil)
	vfmsg := nl.VfMac***REMOVED***
		Vf: uint32(vf),
	***REMOVED***
	copy(vfmsg.Mac[:], []byte(hwaddr))
	nl.NewRtAttrChild(info, nl.IFLA_VF_MAC, vfmsg.Serialize())
	req.AddData(data)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

// LinkSetVfVlan sets the vlan of a vf for the link.
// Equivalent to: `ip link set $link vf $vf vlan $vlan`
func LinkSetVfVlan(link Link, vf, vlan int) error ***REMOVED***
	return pkgHandle.LinkSetVfVlan(link, vf, vlan)
***REMOVED***

// LinkSetVfVlan sets the vlan of a vf for the link.
// Equivalent to: `ip link set $link vf $vf vlan $vlan`
func (h *Handle) LinkSetVfVlan(link Link, vf, vlan int) error ***REMOVED***
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(syscall.RTM_SETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Index = int32(base.Index)
	req.AddData(msg)

	data := nl.NewRtAttr(nl.IFLA_VFINFO_LIST, nil)
	info := nl.NewRtAttrChild(data, nl.IFLA_VF_INFO, nil)
	vfmsg := nl.VfVlan***REMOVED***
		Vf:   uint32(vf),
		Vlan: uint32(vlan),
	***REMOVED***
	nl.NewRtAttrChild(info, nl.IFLA_VF_VLAN, vfmsg.Serialize())
	req.AddData(data)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

// LinkSetVfTxRate sets the tx rate of a vf for the link.
// Equivalent to: `ip link set $link vf $vf rate $rate`
func LinkSetVfTxRate(link Link, vf, rate int) error ***REMOVED***
	return pkgHandle.LinkSetVfTxRate(link, vf, rate)
***REMOVED***

// LinkSetVfTxRate sets the tx rate of a vf for the link.
// Equivalent to: `ip link set $link vf $vf rate $rate`
func (h *Handle) LinkSetVfTxRate(link Link, vf, rate int) error ***REMOVED***
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(syscall.RTM_SETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Index = int32(base.Index)
	req.AddData(msg)

	data := nl.NewRtAttr(nl.IFLA_VFINFO_LIST, nil)
	info := nl.NewRtAttrChild(data, nl.IFLA_VF_INFO, nil)
	vfmsg := nl.VfTxRate***REMOVED***
		Vf:   uint32(vf),
		Rate: uint32(rate),
	***REMOVED***
	nl.NewRtAttrChild(info, nl.IFLA_VF_TX_RATE, vfmsg.Serialize())
	req.AddData(data)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

// LinkSetVfSpoofchk enables/disables spoof check on a vf for the link.
// Equivalent to: `ip link set $link vf $vf spoofchk $check`
func LinkSetVfSpoofchk(link Link, vf int, check bool) error ***REMOVED***
	return pkgHandle.LinkSetVfSpoofchk(link, vf, check)
***REMOVED***

// LinkSetVfSpookfchk enables/disables spoof check on a vf for the link.
// Equivalent to: `ip link set $link vf $vf spoofchk $check`
func (h *Handle) LinkSetVfSpoofchk(link Link, vf int, check bool) error ***REMOVED***
	var setting uint32
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(syscall.RTM_SETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Index = int32(base.Index)
	req.AddData(msg)

	data := nl.NewRtAttr(nl.IFLA_VFINFO_LIST, nil)
	info := nl.NewRtAttrChild(data, nl.IFLA_VF_INFO, nil)
	if check ***REMOVED***
		setting = 1
	***REMOVED***
	vfmsg := nl.VfSpoofchk***REMOVED***
		Vf:      uint32(vf),
		Setting: setting,
	***REMOVED***
	nl.NewRtAttrChild(info, nl.IFLA_VF_SPOOFCHK, vfmsg.Serialize())
	req.AddData(data)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

// LinkSetVfTrust enables/disables trust state on a vf for the link.
// Equivalent to: `ip link set $link vf $vf trust $state`
func LinkSetVfTrust(link Link, vf int, state bool) error ***REMOVED***
	return pkgHandle.LinkSetVfTrust(link, vf, state)
***REMOVED***

// LinkSetVfTrust enables/disables trust state on a vf for the link.
// Equivalent to: `ip link set $link vf $vf trust $state`
func (h *Handle) LinkSetVfTrust(link Link, vf int, state bool) error ***REMOVED***
	var setting uint32
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(syscall.RTM_SETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Index = int32(base.Index)
	req.AddData(msg)

	data := nl.NewRtAttr(nl.IFLA_VFINFO_LIST, nil)
	info := nl.NewRtAttrChild(data, nl.IFLA_VF_INFO, nil)
	if state ***REMOVED***
		setting = 1
	***REMOVED***
	vfmsg := nl.VfTrust***REMOVED***
		Vf:      uint32(vf),
		Setting: setting,
	***REMOVED***
	nl.NewRtAttrChild(info, nl.IFLA_VF_TRUST, vfmsg.Serialize())
	req.AddData(data)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

// LinkSetMaster sets the master of the link device.
// Equivalent to: `ip link set $link master $master`
func LinkSetMaster(link Link, master *Bridge) error ***REMOVED***
	return pkgHandle.LinkSetMaster(link, master)
***REMOVED***

// LinkSetMaster sets the master of the link device.
// Equivalent to: `ip link set $link master $master`
func (h *Handle) LinkSetMaster(link Link, master *Bridge) error ***REMOVED***
	index := 0
	if master != nil ***REMOVED***
		masterBase := master.Attrs()
		h.ensureIndex(masterBase)
		index = masterBase.Index
	***REMOVED***
	if index <= 0 ***REMOVED***
		return fmt.Errorf("Device does not exist")
	***REMOVED***
	return h.LinkSetMasterByIndex(link, index)
***REMOVED***

// LinkSetNoMaster removes the master of the link device.
// Equivalent to: `ip link set $link nomaster`
func LinkSetNoMaster(link Link) error ***REMOVED***
	return pkgHandle.LinkSetNoMaster(link)
***REMOVED***

// LinkSetNoMaster removes the master of the link device.
// Equivalent to: `ip link set $link nomaster`
func (h *Handle) LinkSetNoMaster(link Link) error ***REMOVED***
	return h.LinkSetMasterByIndex(link, 0)
***REMOVED***

// LinkSetMasterByIndex sets the master of the link device.
// Equivalent to: `ip link set $link master $master`
func LinkSetMasterByIndex(link Link, masterIndex int) error ***REMOVED***
	return pkgHandle.LinkSetMasterByIndex(link, masterIndex)
***REMOVED***

// LinkSetMasterByIndex sets the master of the link device.
// Equivalent to: `ip link set $link master $master`
func (h *Handle) LinkSetMasterByIndex(link Link, masterIndex int) error ***REMOVED***
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(syscall.RTM_SETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Index = int32(base.Index)
	req.AddData(msg)

	b := make([]byte, 4)
	native.PutUint32(b, uint32(masterIndex))

	data := nl.NewRtAttr(syscall.IFLA_MASTER, b)
	req.AddData(data)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

// LinkSetNsPid puts the device into a new network namespace. The
// pid must be a pid of a running process.
// Equivalent to: `ip link set $link netns $pid`
func LinkSetNsPid(link Link, nspid int) error ***REMOVED***
	return pkgHandle.LinkSetNsPid(link, nspid)
***REMOVED***

// LinkSetNsPid puts the device into a new network namespace. The
// pid must be a pid of a running process.
// Equivalent to: `ip link set $link netns $pid`
func (h *Handle) LinkSetNsPid(link Link, nspid int) error ***REMOVED***
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(syscall.RTM_SETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Index = int32(base.Index)
	req.AddData(msg)

	b := make([]byte, 4)
	native.PutUint32(b, uint32(nspid))

	data := nl.NewRtAttr(syscall.IFLA_NET_NS_PID, b)
	req.AddData(data)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

// LinkSetNsFd puts the device into a new network namespace. The
// fd must be an open file descriptor to a network namespace.
// Similar to: `ip link set $link netns $ns`
func LinkSetNsFd(link Link, fd int) error ***REMOVED***
	return pkgHandle.LinkSetNsFd(link, fd)
***REMOVED***

// LinkSetNsFd puts the device into a new network namespace. The
// fd must be an open file descriptor to a network namespace.
// Similar to: `ip link set $link netns $ns`
func (h *Handle) LinkSetNsFd(link Link, fd int) error ***REMOVED***
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(syscall.RTM_SETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Index = int32(base.Index)
	req.AddData(msg)

	b := make([]byte, 4)
	native.PutUint32(b, uint32(fd))

	data := nl.NewRtAttr(nl.IFLA_NET_NS_FD, b)
	req.AddData(data)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

// LinkSetXdpFd adds a bpf function to the driver. The fd must be a bpf
// program loaded with bpf(type=BPF_PROG_TYPE_XDP)
func LinkSetXdpFd(link Link, fd int) error ***REMOVED***
	return LinkSetXdpFdWithFlags(link, fd, 0)
***REMOVED***

// LinkSetXdpFdWithFlags adds a bpf function to the driver with the given
// options. The fd must be a bpf program loaded with bpf(type=BPF_PROG_TYPE_XDP)
func LinkSetXdpFdWithFlags(link Link, fd, flags int) error ***REMOVED***
	base := link.Attrs()
	ensureIndex(base)
	req := nl.NewNetlinkRequest(syscall.RTM_SETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Index = int32(base.Index)
	req.AddData(msg)

	addXdpAttrs(&LinkXdp***REMOVED***Fd: fd, Flags: uint32(flags)***REMOVED***, req)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

func boolAttr(val bool) []byte ***REMOVED***
	var v uint8
	if val ***REMOVED***
		v = 1
	***REMOVED***
	return nl.Uint8Attr(v)
***REMOVED***

type vxlanPortRange struct ***REMOVED***
	Lo, Hi uint16
***REMOVED***

func addVxlanAttrs(vxlan *Vxlan, linkInfo *nl.RtAttr) ***REMOVED***
	data := nl.NewRtAttrChild(linkInfo, nl.IFLA_INFO_DATA, nil)

	if vxlan.FlowBased ***REMOVED***
		vxlan.VxlanId = 0
	***REMOVED***

	nl.NewRtAttrChild(data, nl.IFLA_VXLAN_ID, nl.Uint32Attr(uint32(vxlan.VxlanId)))

	if vxlan.VtepDevIndex != 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_VXLAN_LINK, nl.Uint32Attr(uint32(vxlan.VtepDevIndex)))
	***REMOVED***
	if vxlan.SrcAddr != nil ***REMOVED***
		ip := vxlan.SrcAddr.To4()
		if ip != nil ***REMOVED***
			nl.NewRtAttrChild(data, nl.IFLA_VXLAN_LOCAL, []byte(ip))
		***REMOVED*** else ***REMOVED***
			ip = vxlan.SrcAddr.To16()
			if ip != nil ***REMOVED***
				nl.NewRtAttrChild(data, nl.IFLA_VXLAN_LOCAL6, []byte(ip))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if vxlan.Group != nil ***REMOVED***
		group := vxlan.Group.To4()
		if group != nil ***REMOVED***
			nl.NewRtAttrChild(data, nl.IFLA_VXLAN_GROUP, []byte(group))
		***REMOVED*** else ***REMOVED***
			group = vxlan.Group.To16()
			if group != nil ***REMOVED***
				nl.NewRtAttrChild(data, nl.IFLA_VXLAN_GROUP6, []byte(group))
			***REMOVED***
		***REMOVED***
	***REMOVED***

	nl.NewRtAttrChild(data, nl.IFLA_VXLAN_TTL, nl.Uint8Attr(uint8(vxlan.TTL)))
	nl.NewRtAttrChild(data, nl.IFLA_VXLAN_TOS, nl.Uint8Attr(uint8(vxlan.TOS)))
	nl.NewRtAttrChild(data, nl.IFLA_VXLAN_LEARNING, boolAttr(vxlan.Learning))
	nl.NewRtAttrChild(data, nl.IFLA_VXLAN_PROXY, boolAttr(vxlan.Proxy))
	nl.NewRtAttrChild(data, nl.IFLA_VXLAN_RSC, boolAttr(vxlan.RSC))
	nl.NewRtAttrChild(data, nl.IFLA_VXLAN_L2MISS, boolAttr(vxlan.L2miss))
	nl.NewRtAttrChild(data, nl.IFLA_VXLAN_L3MISS, boolAttr(vxlan.L3miss))

	if vxlan.UDPCSum ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_VXLAN_UDP_CSUM, boolAttr(vxlan.UDPCSum))
	***REMOVED***
	if vxlan.GBP ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_VXLAN_GBP, []byte***REMOVED******REMOVED***)
	***REMOVED***
	if vxlan.FlowBased ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_VXLAN_FLOWBASED, boolAttr(vxlan.FlowBased))
	***REMOVED***
	if vxlan.NoAge ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_VXLAN_AGEING, nl.Uint32Attr(0))
	***REMOVED*** else if vxlan.Age > 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_VXLAN_AGEING, nl.Uint32Attr(uint32(vxlan.Age)))
	***REMOVED***
	if vxlan.Limit > 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_VXLAN_LIMIT, nl.Uint32Attr(uint32(vxlan.Limit)))
	***REMOVED***
	if vxlan.Port > 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_VXLAN_PORT, htons(uint16(vxlan.Port)))
	***REMOVED***
	if vxlan.PortLow > 0 || vxlan.PortHigh > 0 ***REMOVED***
		pr := vxlanPortRange***REMOVED***uint16(vxlan.PortLow), uint16(vxlan.PortHigh)***REMOVED***

		buf := new(bytes.Buffer)
		binary.Write(buf, binary.BigEndian, &pr)

		nl.NewRtAttrChild(data, nl.IFLA_VXLAN_PORT_RANGE, buf.Bytes())
	***REMOVED***
***REMOVED***

func addBondAttrs(bond *Bond, linkInfo *nl.RtAttr) ***REMOVED***
	data := nl.NewRtAttrChild(linkInfo, nl.IFLA_INFO_DATA, nil)
	if bond.Mode >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_MODE, nl.Uint8Attr(uint8(bond.Mode)))
	***REMOVED***
	if bond.ActiveSlave >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_ACTIVE_SLAVE, nl.Uint32Attr(uint32(bond.ActiveSlave)))
	***REMOVED***
	if bond.Miimon >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_MIIMON, nl.Uint32Attr(uint32(bond.Miimon)))
	***REMOVED***
	if bond.UpDelay >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_UPDELAY, nl.Uint32Attr(uint32(bond.UpDelay)))
	***REMOVED***
	if bond.DownDelay >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_DOWNDELAY, nl.Uint32Attr(uint32(bond.DownDelay)))
	***REMOVED***
	if bond.UseCarrier >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_USE_CARRIER, nl.Uint8Attr(uint8(bond.UseCarrier)))
	***REMOVED***
	if bond.ArpInterval >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_ARP_INTERVAL, nl.Uint32Attr(uint32(bond.ArpInterval)))
	***REMOVED***
	if bond.ArpIpTargets != nil ***REMOVED***
		msg := nl.NewRtAttrChild(data, nl.IFLA_BOND_ARP_IP_TARGET, nil)
		for i := range bond.ArpIpTargets ***REMOVED***
			ip := bond.ArpIpTargets[i].To4()
			if ip != nil ***REMOVED***
				nl.NewRtAttrChild(msg, i, []byte(ip))
				continue
			***REMOVED***
			ip = bond.ArpIpTargets[i].To16()
			if ip != nil ***REMOVED***
				nl.NewRtAttrChild(msg, i, []byte(ip))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if bond.ArpValidate >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_ARP_VALIDATE, nl.Uint32Attr(uint32(bond.ArpValidate)))
	***REMOVED***
	if bond.ArpAllTargets >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_ARP_ALL_TARGETS, nl.Uint32Attr(uint32(bond.ArpAllTargets)))
	***REMOVED***
	if bond.Primary >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_PRIMARY, nl.Uint32Attr(uint32(bond.Primary)))
	***REMOVED***
	if bond.PrimaryReselect >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_PRIMARY_RESELECT, nl.Uint8Attr(uint8(bond.PrimaryReselect)))
	***REMOVED***
	if bond.FailOverMac >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_FAIL_OVER_MAC, nl.Uint8Attr(uint8(bond.FailOverMac)))
	***REMOVED***
	if bond.XmitHashPolicy >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_XMIT_HASH_POLICY, nl.Uint8Attr(uint8(bond.XmitHashPolicy)))
	***REMOVED***
	if bond.ResendIgmp >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_RESEND_IGMP, nl.Uint32Attr(uint32(bond.ResendIgmp)))
	***REMOVED***
	if bond.NumPeerNotif >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_NUM_PEER_NOTIF, nl.Uint8Attr(uint8(bond.NumPeerNotif)))
	***REMOVED***
	if bond.AllSlavesActive >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_ALL_SLAVES_ACTIVE, nl.Uint8Attr(uint8(bond.AllSlavesActive)))
	***REMOVED***
	if bond.MinLinks >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_MIN_LINKS, nl.Uint32Attr(uint32(bond.MinLinks)))
	***REMOVED***
	if bond.LpInterval >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_LP_INTERVAL, nl.Uint32Attr(uint32(bond.LpInterval)))
	***REMOVED***
	if bond.PackersPerSlave >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_PACKETS_PER_SLAVE, nl.Uint32Attr(uint32(bond.PackersPerSlave)))
	***REMOVED***
	if bond.LacpRate >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_AD_LACP_RATE, nl.Uint8Attr(uint8(bond.LacpRate)))
	***REMOVED***
	if bond.AdSelect >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_AD_SELECT, nl.Uint8Attr(uint8(bond.AdSelect)))
	***REMOVED***
	if bond.AdActorSysPrio >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_AD_ACTOR_SYS_PRIO, nl.Uint16Attr(uint16(bond.AdActorSysPrio)))
	***REMOVED***
	if bond.AdUserPortKey >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_AD_USER_PORT_KEY, nl.Uint16Attr(uint16(bond.AdUserPortKey)))
	***REMOVED***
	if bond.AdActorSystem != nil ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_AD_ACTOR_SYSTEM, []byte(bond.AdActorSystem))
	***REMOVED***
	if bond.TlbDynamicLb >= 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BOND_TLB_DYNAMIC_LB, nl.Uint8Attr(uint8(bond.TlbDynamicLb)))
	***REMOVED***
***REMOVED***

// LinkAdd adds a new link device. The type and features of the device
// are taken from the parameters in the link object.
// Equivalent to: `ip link add $link`
func LinkAdd(link Link) error ***REMOVED***
	return pkgHandle.LinkAdd(link)
***REMOVED***

// LinkAdd adds a new link device. The type and features of the device
// are taken fromt the parameters in the link object.
// Equivalent to: `ip link add $link`
func (h *Handle) LinkAdd(link Link) error ***REMOVED***
	return h.linkModify(link, syscall.NLM_F_CREATE|syscall.NLM_F_EXCL|syscall.NLM_F_ACK)
***REMOVED***

func (h *Handle) linkModify(link Link, flags int) error ***REMOVED***
	// TODO: support extra data for macvlan
	base := link.Attrs()

	if base.Name == "" ***REMOVED***
		return fmt.Errorf("LinkAttrs.Name cannot be empty!")
	***REMOVED***

	if tuntap, ok := link.(*Tuntap); ok ***REMOVED***
		// TODO: support user
		// TODO: support group
		// TODO: multi_queue
		// TODO: support non- persistent
		if tuntap.Mode < syscall.IFF_TUN || tuntap.Mode > syscall.IFF_TAP ***REMOVED***
			return fmt.Errorf("Tuntap.Mode %v unknown!", tuntap.Mode)
		***REMOVED***
		file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		defer file.Close()
		var req ifReq
		if tuntap.Flags == 0 ***REMOVED***
			req.Flags = uint16(TUNTAP_DEFAULTS)
		***REMOVED*** else ***REMOVED***
			req.Flags = uint16(tuntap.Flags)
		***REMOVED***
		req.Flags |= uint16(tuntap.Mode)
		copy(req.Name[:15], base.Name)
		_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req)))
		if errno != 0 ***REMOVED***
			return fmt.Errorf("Tuntap IOCTL TUNSETIFF failed, errno %v", errno)
		***REMOVED***
		_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(syscall.TUNSETPERSIST), 1)
		if errno != 0 ***REMOVED***
			return fmt.Errorf("Tuntap IOCTL TUNSETPERSIST failed, errno %v", errno)
		***REMOVED***
		h.ensureIndex(base)

		// can't set master during create, so set it afterwards
		if base.MasterIndex != 0 ***REMOVED***
			// TODO: verify MasterIndex is actually a bridge?
			return h.LinkSetMasterByIndex(link, base.MasterIndex)
		***REMOVED***
		return nil
	***REMOVED***

	req := h.newNetlinkRequest(syscall.RTM_NEWLINK, flags)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	// TODO: make it shorter
	if base.Flags&net.FlagUp != 0 ***REMOVED***
		msg.Change = syscall.IFF_UP
		msg.Flags = syscall.IFF_UP
	***REMOVED***
	if base.Flags&net.FlagBroadcast != 0 ***REMOVED***
		msg.Change |= syscall.IFF_BROADCAST
		msg.Flags |= syscall.IFF_BROADCAST
	***REMOVED***
	if base.Flags&net.FlagLoopback != 0 ***REMOVED***
		msg.Change |= syscall.IFF_LOOPBACK
		msg.Flags |= syscall.IFF_LOOPBACK
	***REMOVED***
	if base.Flags&net.FlagPointToPoint != 0 ***REMOVED***
		msg.Change |= syscall.IFF_POINTOPOINT
		msg.Flags |= syscall.IFF_POINTOPOINT
	***REMOVED***
	if base.Flags&net.FlagMulticast != 0 ***REMOVED***
		msg.Change |= syscall.IFF_MULTICAST
		msg.Flags |= syscall.IFF_MULTICAST
	***REMOVED***
	req.AddData(msg)

	if base.ParentIndex != 0 ***REMOVED***
		b := make([]byte, 4)
		native.PutUint32(b, uint32(base.ParentIndex))
		data := nl.NewRtAttr(syscall.IFLA_LINK, b)
		req.AddData(data)
	***REMOVED*** else if link.Type() == "ipvlan" ***REMOVED***
		return fmt.Errorf("Can't create ipvlan link without ParentIndex")
	***REMOVED***

	nameData := nl.NewRtAttr(syscall.IFLA_IFNAME, nl.ZeroTerminated(base.Name))
	req.AddData(nameData)

	if base.MTU > 0 ***REMOVED***
		mtu := nl.NewRtAttr(syscall.IFLA_MTU, nl.Uint32Attr(uint32(base.MTU)))
		req.AddData(mtu)
	***REMOVED***

	if base.TxQLen >= 0 ***REMOVED***
		qlen := nl.NewRtAttr(syscall.IFLA_TXQLEN, nl.Uint32Attr(uint32(base.TxQLen)))
		req.AddData(qlen)
	***REMOVED***

	if base.HardwareAddr != nil ***REMOVED***
		hwaddr := nl.NewRtAttr(syscall.IFLA_ADDRESS, []byte(base.HardwareAddr))
		req.AddData(hwaddr)
	***REMOVED***

	if base.Namespace != nil ***REMOVED***
		var attr *nl.RtAttr
		switch base.Namespace.(type) ***REMOVED***
		case NsPid:
			val := nl.Uint32Attr(uint32(base.Namespace.(NsPid)))
			attr = nl.NewRtAttr(syscall.IFLA_NET_NS_PID, val)
		case NsFd:
			val := nl.Uint32Attr(uint32(base.Namespace.(NsFd)))
			attr = nl.NewRtAttr(nl.IFLA_NET_NS_FD, val)
		***REMOVED***

		req.AddData(attr)
	***REMOVED***

	if base.Xdp != nil ***REMOVED***
		addXdpAttrs(base.Xdp, req)
	***REMOVED***

	linkInfo := nl.NewRtAttr(syscall.IFLA_LINKINFO, nil)
	nl.NewRtAttrChild(linkInfo, nl.IFLA_INFO_KIND, nl.NonZeroTerminated(link.Type()))

	switch link := link.(type) ***REMOVED***
	case *Vlan:
		b := make([]byte, 2)
		native.PutUint16(b, uint16(link.VlanId))
		data := nl.NewRtAttrChild(linkInfo, nl.IFLA_INFO_DATA, nil)
		nl.NewRtAttrChild(data, nl.IFLA_VLAN_ID, b)
	case *Veth:
		data := nl.NewRtAttrChild(linkInfo, nl.IFLA_INFO_DATA, nil)
		peer := nl.NewRtAttrChild(data, nl.VETH_INFO_PEER, nil)
		nl.NewIfInfomsgChild(peer, syscall.AF_UNSPEC)
		nl.NewRtAttrChild(peer, syscall.IFLA_IFNAME, nl.ZeroTerminated(link.PeerName))
		if base.TxQLen >= 0 ***REMOVED***
			nl.NewRtAttrChild(peer, syscall.IFLA_TXQLEN, nl.Uint32Attr(uint32(base.TxQLen)))
		***REMOVED***
		if base.MTU > 0 ***REMOVED***
			nl.NewRtAttrChild(peer, syscall.IFLA_MTU, nl.Uint32Attr(uint32(base.MTU)))
		***REMOVED***

	case *Vxlan:
		addVxlanAttrs(link, linkInfo)
	case *Bond:
		addBondAttrs(link, linkInfo)
	case *IPVlan:
		data := nl.NewRtAttrChild(linkInfo, nl.IFLA_INFO_DATA, nil)
		nl.NewRtAttrChild(data, nl.IFLA_IPVLAN_MODE, nl.Uint16Attr(uint16(link.Mode)))
	case *Macvlan:
		if link.Mode != MACVLAN_MODE_DEFAULT ***REMOVED***
			data := nl.NewRtAttrChild(linkInfo, nl.IFLA_INFO_DATA, nil)
			nl.NewRtAttrChild(data, nl.IFLA_MACVLAN_MODE, nl.Uint32Attr(macvlanModes[link.Mode]))
		***REMOVED***
	case *Macvtap:
		if link.Mode != MACVLAN_MODE_DEFAULT ***REMOVED***
			data := nl.NewRtAttrChild(linkInfo, nl.IFLA_INFO_DATA, nil)
			nl.NewRtAttrChild(data, nl.IFLA_MACVLAN_MODE, nl.Uint32Attr(macvlanModes[link.Mode]))
		***REMOVED***
	case *Gretap:
		addGretapAttrs(link, linkInfo)
	case *Iptun:
		addIptunAttrs(link, linkInfo)
	case *Gretun:
		addGretunAttrs(link, linkInfo)
	case *Vti:
		addVtiAttrs(link, linkInfo)
	case *Vrf:
		addVrfAttrs(link, linkInfo)
	case *Bridge:
		addBridgeAttrs(link, linkInfo)
	case *GTP:
		addGTPAttrs(link, linkInfo)
	***REMOVED***

	req.AddData(linkInfo)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	h.ensureIndex(base)

	// can't set master during create, so set it afterwards
	if base.MasterIndex != 0 ***REMOVED***
		// TODO: verify MasterIndex is actually a bridge?
		return h.LinkSetMasterByIndex(link, base.MasterIndex)
	***REMOVED***
	return nil
***REMOVED***

// LinkDel deletes link device. Either Index or Name must be set in
// the link object for it to be deleted. The other values are ignored.
// Equivalent to: `ip link del $link`
func LinkDel(link Link) error ***REMOVED***
	return pkgHandle.LinkDel(link)
***REMOVED***

// LinkDel deletes link device. Either Index or Name must be set in
// the link object for it to be deleted. The other values are ignored.
// Equivalent to: `ip link del $link`
func (h *Handle) LinkDel(link Link) error ***REMOVED***
	base := link.Attrs()

	h.ensureIndex(base)

	req := h.newNetlinkRequest(syscall.RTM_DELLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Index = int32(base.Index)
	req.AddData(msg)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

func (h *Handle) linkByNameDump(name string) (Link, error) ***REMOVED***
	links, err := h.LinkList()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for _, link := range links ***REMOVED***
		if link.Attrs().Name == name ***REMOVED***
			return link, nil
		***REMOVED***
	***REMOVED***
	return nil, LinkNotFoundError***REMOVED***fmt.Errorf("Link %s not found", name)***REMOVED***
***REMOVED***

func (h *Handle) linkByAliasDump(alias string) (Link, error) ***REMOVED***
	links, err := h.LinkList()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for _, link := range links ***REMOVED***
		if link.Attrs().Alias == alias ***REMOVED***
			return link, nil
		***REMOVED***
	***REMOVED***
	return nil, LinkNotFoundError***REMOVED***fmt.Errorf("Link alias %s not found", alias)***REMOVED***
***REMOVED***

// LinkByName finds a link by name and returns a pointer to the object.
func LinkByName(name string) (Link, error) ***REMOVED***
	return pkgHandle.LinkByName(name)
***REMOVED***

// LinkByName finds a link by name and returns a pointer to the object.
func (h *Handle) LinkByName(name string) (Link, error) ***REMOVED***
	if h.lookupByDump ***REMOVED***
		return h.linkByNameDump(name)
	***REMOVED***

	req := h.newNetlinkRequest(syscall.RTM_GETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	req.AddData(msg)

	nameData := nl.NewRtAttr(syscall.IFLA_IFNAME, nl.ZeroTerminated(name))
	req.AddData(nameData)

	link, err := execGetLink(req)
	if err == syscall.EINVAL ***REMOVED***
		// older kernels don't support looking up via IFLA_IFNAME
		// so fall back to dumping all links
		h.lookupByDump = true
		return h.linkByNameDump(name)
	***REMOVED***

	return link, err
***REMOVED***

// LinkByAlias finds a link by its alias and returns a pointer to the object.
// If there are multiple links with the alias it returns the first one
func LinkByAlias(alias string) (Link, error) ***REMOVED***
	return pkgHandle.LinkByAlias(alias)
***REMOVED***

// LinkByAlias finds a link by its alias and returns a pointer to the object.
// If there are multiple links with the alias it returns the first one
func (h *Handle) LinkByAlias(alias string) (Link, error) ***REMOVED***
	if h.lookupByDump ***REMOVED***
		return h.linkByAliasDump(alias)
	***REMOVED***

	req := h.newNetlinkRequest(syscall.RTM_GETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	req.AddData(msg)

	nameData := nl.NewRtAttr(syscall.IFLA_IFALIAS, nl.ZeroTerminated(alias))
	req.AddData(nameData)

	link, err := execGetLink(req)
	if err == syscall.EINVAL ***REMOVED***
		// older kernels don't support looking up via IFLA_IFALIAS
		// so fall back to dumping all links
		h.lookupByDump = true
		return h.linkByAliasDump(alias)
	***REMOVED***

	return link, err
***REMOVED***

// LinkByIndex finds a link by index and returns a pointer to the object.
func LinkByIndex(index int) (Link, error) ***REMOVED***
	return pkgHandle.LinkByIndex(index)
***REMOVED***

// LinkByIndex finds a link by index and returns a pointer to the object.
func (h *Handle) LinkByIndex(index int) (Link, error) ***REMOVED***
	req := h.newNetlinkRequest(syscall.RTM_GETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Index = int32(index)
	req.AddData(msg)

	return execGetLink(req)
***REMOVED***

func execGetLink(req *nl.NetlinkRequest) (Link, error) ***REMOVED***
	msgs, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	if err != nil ***REMOVED***
		if errno, ok := err.(syscall.Errno); ok ***REMOVED***
			if errno == syscall.ENODEV ***REMOVED***
				return nil, LinkNotFoundError***REMOVED***fmt.Errorf("Link not found")***REMOVED***
			***REMOVED***
		***REMOVED***
		return nil, err
	***REMOVED***

	switch ***REMOVED***
	case len(msgs) == 0:
		return nil, LinkNotFoundError***REMOVED***fmt.Errorf("Link not found")***REMOVED***

	case len(msgs) == 1:
		return LinkDeserialize(nil, msgs[0])

	default:
		return nil, fmt.Errorf("More than one link found")
	***REMOVED***
***REMOVED***

// linkDeserialize deserializes a raw message received from netlink into
// a link object.
func LinkDeserialize(hdr *syscall.NlMsghdr, m []byte) (Link, error) ***REMOVED***
	msg := nl.DeserializeIfInfomsg(m)

	attrs, err := nl.ParseRouteAttr(m[msg.Len():])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	base := LinkAttrs***REMOVED***Index: int(msg.Index), RawFlags: msg.Flags, Flags: linkFlags(msg.Flags), EncapType: msg.EncapType()***REMOVED***
	if msg.Flags&syscall.IFF_PROMISC != 0 ***REMOVED***
		base.Promisc = 1
	***REMOVED***
	var (
		link     Link
		stats32  []byte
		stats64  []byte
		linkType string
	)
	for _, attr := range attrs ***REMOVED***
		switch attr.Attr.Type ***REMOVED***
		case syscall.IFLA_LINKINFO:
			infos, err := nl.ParseRouteAttr(attr.Value)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			for _, info := range infos ***REMOVED***
				switch info.Attr.Type ***REMOVED***
				case nl.IFLA_INFO_KIND:
					linkType = string(info.Value[:len(info.Value)-1])
					switch linkType ***REMOVED***
					case "dummy":
						link = &Dummy***REMOVED******REMOVED***
					case "ifb":
						link = &Ifb***REMOVED******REMOVED***
					case "bridge":
						link = &Bridge***REMOVED******REMOVED***
					case "vlan":
						link = &Vlan***REMOVED******REMOVED***
					case "veth":
						link = &Veth***REMOVED******REMOVED***
					case "vxlan":
						link = &Vxlan***REMOVED******REMOVED***
					case "bond":
						link = &Bond***REMOVED******REMOVED***
					case "ipvlan":
						link = &IPVlan***REMOVED******REMOVED***
					case "macvlan":
						link = &Macvlan***REMOVED******REMOVED***
					case "macvtap":
						link = &Macvtap***REMOVED******REMOVED***
					case "gretap":
						link = &Gretap***REMOVED******REMOVED***
					case "ipip":
						link = &Iptun***REMOVED******REMOVED***
					case "gre":
						link = &Gretun***REMOVED******REMOVED***
					case "vti":
						link = &Vti***REMOVED******REMOVED***
					case "vrf":
						link = &Vrf***REMOVED******REMOVED***
					case "gtp":
						link = &GTP***REMOVED******REMOVED***
					default:
						link = &GenericLink***REMOVED***LinkType: linkType***REMOVED***
					***REMOVED***
				case nl.IFLA_INFO_DATA:
					data, err := nl.ParseRouteAttr(info.Value)
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***
					switch linkType ***REMOVED***
					case "vlan":
						parseVlanData(link, data)
					case "vxlan":
						parseVxlanData(link, data)
					case "bond":
						parseBondData(link, data)
					case "ipvlan":
						parseIPVlanData(link, data)
					case "macvlan":
						parseMacvlanData(link, data)
					case "macvtap":
						parseMacvtapData(link, data)
					case "gretap":
						parseGretapData(link, data)
					case "ipip":
						parseIptunData(link, data)
					case "gre":
						parseGretunData(link, data)
					case "vti":
						parseVtiData(link, data)
					case "vrf":
						parseVrfData(link, data)
					case "bridge":
						parseBridgeData(link, data)
					case "gtp":
						parseGTPData(link, data)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		case syscall.IFLA_ADDRESS:
			var nonzero bool
			for _, b := range attr.Value ***REMOVED***
				if b != 0 ***REMOVED***
					nonzero = true
				***REMOVED***
			***REMOVED***
			if nonzero ***REMOVED***
				base.HardwareAddr = attr.Value[:]
			***REMOVED***
		case syscall.IFLA_IFNAME:
			base.Name = string(attr.Value[:len(attr.Value)-1])
		case syscall.IFLA_MTU:
			base.MTU = int(native.Uint32(attr.Value[0:4]))
		case syscall.IFLA_LINK:
			base.ParentIndex = int(native.Uint32(attr.Value[0:4]))
		case syscall.IFLA_MASTER:
			base.MasterIndex = int(native.Uint32(attr.Value[0:4]))
		case syscall.IFLA_TXQLEN:
			base.TxQLen = int(native.Uint32(attr.Value[0:4]))
		case syscall.IFLA_IFALIAS:
			base.Alias = string(attr.Value[:len(attr.Value)-1])
		case syscall.IFLA_STATS:
			stats32 = attr.Value[:]
		case IFLA_STATS64:
			stats64 = attr.Value[:]
		case nl.IFLA_XDP:
			xdp, err := parseLinkXdp(attr.Value[:])
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			base.Xdp = xdp
		case syscall.IFLA_PROTINFO | syscall.NLA_F_NESTED:
			if hdr != nil && hdr.Type == syscall.RTM_NEWLINK &&
				msg.Family == syscall.AF_BRIDGE ***REMOVED***
				attrs, err := nl.ParseRouteAttr(attr.Value[:])
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				base.Protinfo = parseProtinfo(attrs)
			***REMOVED***
		case syscall.IFLA_OPERSTATE:
			base.OperState = LinkOperState(uint8(attr.Value[0]))
		case nl.IFLA_LINK_NETNSID:
			base.NetNsID = int(native.Uint32(attr.Value[0:4]))
		***REMOVED***
	***REMOVED***

	if stats64 != nil ***REMOVED***
		base.Statistics = parseLinkStats64(stats64)
	***REMOVED*** else if stats32 != nil ***REMOVED***
		base.Statistics = parseLinkStats32(stats32)
	***REMOVED***

	// Links that don't have IFLA_INFO_KIND are hardware devices
	if link == nil ***REMOVED***
		link = &Device***REMOVED******REMOVED***
	***REMOVED***
	*link.Attrs() = base

	return link, nil
***REMOVED***

// LinkList gets a list of link devices.
// Equivalent to: `ip link show`
func LinkList() ([]Link, error) ***REMOVED***
	return pkgHandle.LinkList()
***REMOVED***

// LinkList gets a list of link devices.
// Equivalent to: `ip link show`
func (h *Handle) LinkList() ([]Link, error) ***REMOVED***
	// NOTE(vish): This duplicates functionality in net/iface_linux.go, but we need
	//             to get the message ourselves to parse link type.
	req := h.newNetlinkRequest(syscall.RTM_GETLINK, syscall.NLM_F_DUMP)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	req.AddData(msg)

	msgs, err := req.Execute(syscall.NETLINK_ROUTE, syscall.RTM_NEWLINK)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var res []Link
	for _, m := range msgs ***REMOVED***
		link, err := LinkDeserialize(nil, m)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		res = append(res, link)
	***REMOVED***

	return res, nil
***REMOVED***

// LinkUpdate is used to pass information back from LinkSubscribe()
type LinkUpdate struct ***REMOVED***
	nl.IfInfomsg
	Header syscall.NlMsghdr
	Link
***REMOVED***

// LinkSubscribe takes a chan down which notifications will be sent
// when links change.  Close the 'done' chan to stop subscription.
func LinkSubscribe(ch chan<- LinkUpdate, done <-chan struct***REMOVED******REMOVED***) error ***REMOVED***
	return linkSubscribeAt(netns.None(), netns.None(), ch, done, nil)
***REMOVED***

// LinkSubscribeAt works like LinkSubscribe plus it allows the caller
// to choose the network namespace in which to subscribe (ns).
func LinkSubscribeAt(ns netns.NsHandle, ch chan<- LinkUpdate, done <-chan struct***REMOVED******REMOVED***) error ***REMOVED***
	return linkSubscribeAt(ns, netns.None(), ch, done, nil)
***REMOVED***

// LinkSubscribeOptions contains a set of options to use with
// LinkSubscribeWithOptions.
type LinkSubscribeOptions struct ***REMOVED***
	Namespace     *netns.NsHandle
	ErrorCallback func(error)
***REMOVED***

// LinkSubscribeWithOptions work like LinkSubscribe but enable to
// provide additional options to modify the behavior. Currently, the
// namespace can be provided as well as an error callback.
func LinkSubscribeWithOptions(ch chan<- LinkUpdate, done <-chan struct***REMOVED******REMOVED***, options LinkSubscribeOptions) error ***REMOVED***
	if options.Namespace == nil ***REMOVED***
		none := netns.None()
		options.Namespace = &none
	***REMOVED***
	return linkSubscribeAt(*options.Namespace, netns.None(), ch, done, options.ErrorCallback)
***REMOVED***

func linkSubscribeAt(newNs, curNs netns.NsHandle, ch chan<- LinkUpdate, done <-chan struct***REMOVED******REMOVED***, cberr func(error)) error ***REMOVED***
	s, err := nl.SubscribeAt(newNs, curNs, syscall.NETLINK_ROUTE, syscall.RTNLGRP_LINK)
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
				ifmsg := nl.DeserializeIfInfomsg(m.Data)
				link, err := LinkDeserialize(&m.Header, m.Data)
				if err != nil ***REMOVED***
					if cberr != nil ***REMOVED***
						cberr(err)
					***REMOVED***
					return
				***REMOVED***
				ch <- LinkUpdate***REMOVED***IfInfomsg: *ifmsg, Header: m.Header, Link: link***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return nil
***REMOVED***

func LinkSetHairpin(link Link, mode bool) error ***REMOVED***
	return pkgHandle.LinkSetHairpin(link, mode)
***REMOVED***

func (h *Handle) LinkSetHairpin(link Link, mode bool) error ***REMOVED***
	return h.setProtinfoAttr(link, mode, nl.IFLA_BRPORT_MODE)
***REMOVED***

func LinkSetGuard(link Link, mode bool) error ***REMOVED***
	return pkgHandle.LinkSetGuard(link, mode)
***REMOVED***

func (h *Handle) LinkSetGuard(link Link, mode bool) error ***REMOVED***
	return h.setProtinfoAttr(link, mode, nl.IFLA_BRPORT_GUARD)
***REMOVED***

func LinkSetFastLeave(link Link, mode bool) error ***REMOVED***
	return pkgHandle.LinkSetFastLeave(link, mode)
***REMOVED***

func (h *Handle) LinkSetFastLeave(link Link, mode bool) error ***REMOVED***
	return h.setProtinfoAttr(link, mode, nl.IFLA_BRPORT_FAST_LEAVE)
***REMOVED***

func LinkSetLearning(link Link, mode bool) error ***REMOVED***
	return pkgHandle.LinkSetLearning(link, mode)
***REMOVED***

func (h *Handle) LinkSetLearning(link Link, mode bool) error ***REMOVED***
	return h.setProtinfoAttr(link, mode, nl.IFLA_BRPORT_LEARNING)
***REMOVED***

func LinkSetRootBlock(link Link, mode bool) error ***REMOVED***
	return pkgHandle.LinkSetRootBlock(link, mode)
***REMOVED***

func (h *Handle) LinkSetRootBlock(link Link, mode bool) error ***REMOVED***
	return h.setProtinfoAttr(link, mode, nl.IFLA_BRPORT_PROTECT)
***REMOVED***

func LinkSetFlood(link Link, mode bool) error ***REMOVED***
	return pkgHandle.LinkSetFlood(link, mode)
***REMOVED***

func (h *Handle) LinkSetFlood(link Link, mode bool) error ***REMOVED***
	return h.setProtinfoAttr(link, mode, nl.IFLA_BRPORT_UNICAST_FLOOD)
***REMOVED***

func LinkSetBrProxyArp(link Link, mode bool) error ***REMOVED***
	return pkgHandle.LinkSetBrProxyArp(link, mode)
***REMOVED***

func (h *Handle) LinkSetBrProxyArp(link Link, mode bool) error ***REMOVED***
	return h.setProtinfoAttr(link, mode, nl.IFLA_BRPORT_PROXYARP)
***REMOVED***

func LinkSetBrProxyArpWiFi(link Link, mode bool) error ***REMOVED***
	return pkgHandle.LinkSetBrProxyArpWiFi(link, mode)
***REMOVED***

func (h *Handle) LinkSetBrProxyArpWiFi(link Link, mode bool) error ***REMOVED***
	return h.setProtinfoAttr(link, mode, nl.IFLA_BRPORT_PROXYARP_WIFI)
***REMOVED***

func (h *Handle) setProtinfoAttr(link Link, mode bool, attr int) error ***REMOVED***
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(syscall.RTM_SETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_BRIDGE)
	msg.Index = int32(base.Index)
	req.AddData(msg)

	br := nl.NewRtAttr(syscall.IFLA_PROTINFO|syscall.NLA_F_NESTED, nil)
	nl.NewRtAttrChild(br, attr, boolToByte(mode))
	req.AddData(br)
	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// LinkSetTxQLen sets the transaction queue length for the link.
// Equivalent to: `ip link set $link txqlen $qlen`
func LinkSetTxQLen(link Link, qlen int) error ***REMOVED***
	return pkgHandle.LinkSetTxQLen(link, qlen)
***REMOVED***

// LinkSetTxQLen sets the transaction queue length for the link.
// Equivalent to: `ip link set $link txqlen $qlen`
func (h *Handle) LinkSetTxQLen(link Link, qlen int) error ***REMOVED***
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(syscall.RTM_SETLINK, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_UNSPEC)
	msg.Index = int32(base.Index)
	req.AddData(msg)

	b := make([]byte, 4)
	native.PutUint32(b, uint32(qlen))

	data := nl.NewRtAttr(syscall.IFLA_TXQLEN, b)
	req.AddData(data)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

func parseVlanData(link Link, data []syscall.NetlinkRouteAttr) ***REMOVED***
	vlan := link.(*Vlan)
	for _, datum := range data ***REMOVED***
		switch datum.Attr.Type ***REMOVED***
		case nl.IFLA_VLAN_ID:
			vlan.VlanId = int(native.Uint16(datum.Value[0:2]))
		***REMOVED***
	***REMOVED***
***REMOVED***

func parseVxlanData(link Link, data []syscall.NetlinkRouteAttr) ***REMOVED***
	vxlan := link.(*Vxlan)
	for _, datum := range data ***REMOVED***
		switch datum.Attr.Type ***REMOVED***
		case nl.IFLA_VXLAN_ID:
			vxlan.VxlanId = int(native.Uint32(datum.Value[0:4]))
		case nl.IFLA_VXLAN_LINK:
			vxlan.VtepDevIndex = int(native.Uint32(datum.Value[0:4]))
		case nl.IFLA_VXLAN_LOCAL:
			vxlan.SrcAddr = net.IP(datum.Value[0:4])
		case nl.IFLA_VXLAN_LOCAL6:
			vxlan.SrcAddr = net.IP(datum.Value[0:16])
		case nl.IFLA_VXLAN_GROUP:
			vxlan.Group = net.IP(datum.Value[0:4])
		case nl.IFLA_VXLAN_GROUP6:
			vxlan.Group = net.IP(datum.Value[0:16])
		case nl.IFLA_VXLAN_TTL:
			vxlan.TTL = int(datum.Value[0])
		case nl.IFLA_VXLAN_TOS:
			vxlan.TOS = int(datum.Value[0])
		case nl.IFLA_VXLAN_LEARNING:
			vxlan.Learning = int8(datum.Value[0]) != 0
		case nl.IFLA_VXLAN_PROXY:
			vxlan.Proxy = int8(datum.Value[0]) != 0
		case nl.IFLA_VXLAN_RSC:
			vxlan.RSC = int8(datum.Value[0]) != 0
		case nl.IFLA_VXLAN_L2MISS:
			vxlan.L2miss = int8(datum.Value[0]) != 0
		case nl.IFLA_VXLAN_L3MISS:
			vxlan.L3miss = int8(datum.Value[0]) != 0
		case nl.IFLA_VXLAN_UDP_CSUM:
			vxlan.UDPCSum = int8(datum.Value[0]) != 0
		case nl.IFLA_VXLAN_GBP:
			vxlan.GBP = true
		case nl.IFLA_VXLAN_FLOWBASED:
			vxlan.FlowBased = int8(datum.Value[0]) != 0
		case nl.IFLA_VXLAN_AGEING:
			vxlan.Age = int(native.Uint32(datum.Value[0:4]))
			vxlan.NoAge = vxlan.Age == 0
		case nl.IFLA_VXLAN_LIMIT:
			vxlan.Limit = int(native.Uint32(datum.Value[0:4]))
		case nl.IFLA_VXLAN_PORT:
			vxlan.Port = int(ntohs(datum.Value[0:2]))
		case nl.IFLA_VXLAN_PORT_RANGE:
			buf := bytes.NewBuffer(datum.Value[0:4])
			var pr vxlanPortRange
			if binary.Read(buf, binary.BigEndian, &pr) != nil ***REMOVED***
				vxlan.PortLow = int(pr.Lo)
				vxlan.PortHigh = int(pr.Hi)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func parseBondData(link Link, data []syscall.NetlinkRouteAttr) ***REMOVED***
	bond := link.(*Bond)
	for i := range data ***REMOVED***
		switch data[i].Attr.Type ***REMOVED***
		case nl.IFLA_BOND_MODE:
			bond.Mode = BondMode(data[i].Value[0])
		case nl.IFLA_BOND_ACTIVE_SLAVE:
			bond.ActiveSlave = int(native.Uint32(data[i].Value[0:4]))
		case nl.IFLA_BOND_MIIMON:
			bond.Miimon = int(native.Uint32(data[i].Value[0:4]))
		case nl.IFLA_BOND_UPDELAY:
			bond.UpDelay = int(native.Uint32(data[i].Value[0:4]))
		case nl.IFLA_BOND_DOWNDELAY:
			bond.DownDelay = int(native.Uint32(data[i].Value[0:4]))
		case nl.IFLA_BOND_USE_CARRIER:
			bond.UseCarrier = int(data[i].Value[0])
		case nl.IFLA_BOND_ARP_INTERVAL:
			bond.ArpInterval = int(native.Uint32(data[i].Value[0:4]))
		case nl.IFLA_BOND_ARP_IP_TARGET:
			// TODO: implement
		case nl.IFLA_BOND_ARP_VALIDATE:
			bond.ArpValidate = BondArpValidate(native.Uint32(data[i].Value[0:4]))
		case nl.IFLA_BOND_ARP_ALL_TARGETS:
			bond.ArpAllTargets = BondArpAllTargets(native.Uint32(data[i].Value[0:4]))
		case nl.IFLA_BOND_PRIMARY:
			bond.Primary = int(native.Uint32(data[i].Value[0:4]))
		case nl.IFLA_BOND_PRIMARY_RESELECT:
			bond.PrimaryReselect = BondPrimaryReselect(data[i].Value[0])
		case nl.IFLA_BOND_FAIL_OVER_MAC:
			bond.FailOverMac = BondFailOverMac(data[i].Value[0])
		case nl.IFLA_BOND_XMIT_HASH_POLICY:
			bond.XmitHashPolicy = BondXmitHashPolicy(data[i].Value[0])
		case nl.IFLA_BOND_RESEND_IGMP:
			bond.ResendIgmp = int(native.Uint32(data[i].Value[0:4]))
		case nl.IFLA_BOND_NUM_PEER_NOTIF:
			bond.NumPeerNotif = int(data[i].Value[0])
		case nl.IFLA_BOND_ALL_SLAVES_ACTIVE:
			bond.AllSlavesActive = int(data[i].Value[0])
		case nl.IFLA_BOND_MIN_LINKS:
			bond.MinLinks = int(native.Uint32(data[i].Value[0:4]))
		case nl.IFLA_BOND_LP_INTERVAL:
			bond.LpInterval = int(native.Uint32(data[i].Value[0:4]))
		case nl.IFLA_BOND_PACKETS_PER_SLAVE:
			bond.PackersPerSlave = int(native.Uint32(data[i].Value[0:4]))
		case nl.IFLA_BOND_AD_LACP_RATE:
			bond.LacpRate = BondLacpRate(data[i].Value[0])
		case nl.IFLA_BOND_AD_SELECT:
			bond.AdSelect = BondAdSelect(data[i].Value[0])
		case nl.IFLA_BOND_AD_INFO:
			// TODO: implement
		case nl.IFLA_BOND_AD_ACTOR_SYS_PRIO:
			bond.AdActorSysPrio = int(native.Uint16(data[i].Value[0:2]))
		case nl.IFLA_BOND_AD_USER_PORT_KEY:
			bond.AdUserPortKey = int(native.Uint16(data[i].Value[0:2]))
		case nl.IFLA_BOND_AD_ACTOR_SYSTEM:
			bond.AdActorSystem = net.HardwareAddr(data[i].Value[0:6])
		case nl.IFLA_BOND_TLB_DYNAMIC_LB:
			bond.TlbDynamicLb = int(data[i].Value[0])
		***REMOVED***
	***REMOVED***
***REMOVED***

func parseIPVlanData(link Link, data []syscall.NetlinkRouteAttr) ***REMOVED***
	ipv := link.(*IPVlan)
	for _, datum := range data ***REMOVED***
		if datum.Attr.Type == nl.IFLA_IPVLAN_MODE ***REMOVED***
			ipv.Mode = IPVlanMode(native.Uint32(datum.Value[0:4]))
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func parseMacvtapData(link Link, data []syscall.NetlinkRouteAttr) ***REMOVED***
	macv := link.(*Macvtap)
	parseMacvlanData(&macv.Macvlan, data)
***REMOVED***

func parseMacvlanData(link Link, data []syscall.NetlinkRouteAttr) ***REMOVED***
	macv := link.(*Macvlan)
	for _, datum := range data ***REMOVED***
		if datum.Attr.Type == nl.IFLA_MACVLAN_MODE ***REMOVED***
			switch native.Uint32(datum.Value[0:4]) ***REMOVED***
			case nl.MACVLAN_MODE_PRIVATE:
				macv.Mode = MACVLAN_MODE_PRIVATE
			case nl.MACVLAN_MODE_VEPA:
				macv.Mode = MACVLAN_MODE_VEPA
			case nl.MACVLAN_MODE_BRIDGE:
				macv.Mode = MACVLAN_MODE_BRIDGE
			case nl.MACVLAN_MODE_PASSTHRU:
				macv.Mode = MACVLAN_MODE_PASSTHRU
			case nl.MACVLAN_MODE_SOURCE:
				macv.Mode = MACVLAN_MODE_SOURCE
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// copied from pkg/net_linux.go
func linkFlags(rawFlags uint32) net.Flags ***REMOVED***
	var f net.Flags
	if rawFlags&syscall.IFF_UP != 0 ***REMOVED***
		f |= net.FlagUp
	***REMOVED***
	if rawFlags&syscall.IFF_BROADCAST != 0 ***REMOVED***
		f |= net.FlagBroadcast
	***REMOVED***
	if rawFlags&syscall.IFF_LOOPBACK != 0 ***REMOVED***
		f |= net.FlagLoopback
	***REMOVED***
	if rawFlags&syscall.IFF_POINTOPOINT != 0 ***REMOVED***
		f |= net.FlagPointToPoint
	***REMOVED***
	if rawFlags&syscall.IFF_MULTICAST != 0 ***REMOVED***
		f |= net.FlagMulticast
	***REMOVED***
	return f
***REMOVED***

func addGretapAttrs(gretap *Gretap, linkInfo *nl.RtAttr) ***REMOVED***
	data := nl.NewRtAttrChild(linkInfo, nl.IFLA_INFO_DATA, nil)

	if gretap.FlowBased ***REMOVED***
		// In flow based mode, no other attributes need to be configured
		nl.NewRtAttrChild(data, nl.IFLA_GRE_COLLECT_METADATA, boolAttr(gretap.FlowBased))
		return
	***REMOVED***

	ip := gretap.Local.To4()
	if ip != nil ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_GRE_LOCAL, []byte(ip))
	***REMOVED***
	ip = gretap.Remote.To4()
	if ip != nil ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_GRE_REMOTE, []byte(ip))
	***REMOVED***

	if gretap.IKey != 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_GRE_IKEY, htonl(gretap.IKey))
		gretap.IFlags |= uint16(nl.GRE_KEY)
	***REMOVED***

	if gretap.OKey != 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_GRE_OKEY, htonl(gretap.OKey))
		gretap.OFlags |= uint16(nl.GRE_KEY)
	***REMOVED***

	nl.NewRtAttrChild(data, nl.IFLA_GRE_IFLAGS, htons(gretap.IFlags))
	nl.NewRtAttrChild(data, nl.IFLA_GRE_OFLAGS, htons(gretap.OFlags))

	if gretap.Link != 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_GRE_LINK, nl.Uint32Attr(gretap.Link))
	***REMOVED***

	nl.NewRtAttrChild(data, nl.IFLA_GRE_PMTUDISC, nl.Uint8Attr(gretap.PMtuDisc))
	nl.NewRtAttrChild(data, nl.IFLA_GRE_TTL, nl.Uint8Attr(gretap.Ttl))
	nl.NewRtAttrChild(data, nl.IFLA_GRE_TOS, nl.Uint8Attr(gretap.Tos))
	nl.NewRtAttrChild(data, nl.IFLA_GRE_ENCAP_TYPE, nl.Uint16Attr(gretap.EncapType))
	nl.NewRtAttrChild(data, nl.IFLA_GRE_ENCAP_FLAGS, nl.Uint16Attr(gretap.EncapFlags))
	nl.NewRtAttrChild(data, nl.IFLA_GRE_ENCAP_SPORT, htons(gretap.EncapSport))
	nl.NewRtAttrChild(data, nl.IFLA_GRE_ENCAP_DPORT, htons(gretap.EncapDport))
***REMOVED***

func parseGretapData(link Link, data []syscall.NetlinkRouteAttr) ***REMOVED***
	gre := link.(*Gretap)
	for _, datum := range data ***REMOVED***
		switch datum.Attr.Type ***REMOVED***
		case nl.IFLA_GRE_OKEY:
			gre.IKey = ntohl(datum.Value[0:4])
		case nl.IFLA_GRE_IKEY:
			gre.OKey = ntohl(datum.Value[0:4])
		case nl.IFLA_GRE_LOCAL:
			gre.Local = net.IP(datum.Value[0:4])
		case nl.IFLA_GRE_REMOTE:
			gre.Remote = net.IP(datum.Value[0:4])
		case nl.IFLA_GRE_ENCAP_SPORT:
			gre.EncapSport = ntohs(datum.Value[0:2])
		case nl.IFLA_GRE_ENCAP_DPORT:
			gre.EncapDport = ntohs(datum.Value[0:2])
		case nl.IFLA_GRE_IFLAGS:
			gre.IFlags = ntohs(datum.Value[0:2])
		case nl.IFLA_GRE_OFLAGS:
			gre.OFlags = ntohs(datum.Value[0:2])

		case nl.IFLA_GRE_TTL:
			gre.Ttl = uint8(datum.Value[0])
		case nl.IFLA_GRE_TOS:
			gre.Tos = uint8(datum.Value[0])
		case nl.IFLA_GRE_PMTUDISC:
			gre.PMtuDisc = uint8(datum.Value[0])
		case nl.IFLA_GRE_ENCAP_TYPE:
			gre.EncapType = native.Uint16(datum.Value[0:2])
		case nl.IFLA_GRE_ENCAP_FLAGS:
			gre.EncapFlags = native.Uint16(datum.Value[0:2])
		case nl.IFLA_GRE_COLLECT_METADATA:
			gre.FlowBased = int8(datum.Value[0]) != 0
		***REMOVED***
	***REMOVED***
***REMOVED***

func addGretunAttrs(gre *Gretun, linkInfo *nl.RtAttr) ***REMOVED***
	data := nl.NewRtAttrChild(linkInfo, nl.IFLA_INFO_DATA, nil)

	ip := gre.Local.To4()
	if ip != nil ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_GRE_LOCAL, []byte(ip))
	***REMOVED***
	ip = gre.Remote.To4()
	if ip != nil ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_GRE_REMOTE, []byte(ip))
	***REMOVED***

	if gre.IKey != 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_GRE_IKEY, htonl(gre.IKey))
		gre.IFlags |= uint16(nl.GRE_KEY)
	***REMOVED***

	if gre.OKey != 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_GRE_OKEY, htonl(gre.OKey))
		gre.OFlags |= uint16(nl.GRE_KEY)
	***REMOVED***

	nl.NewRtAttrChild(data, nl.IFLA_GRE_IFLAGS, htons(gre.IFlags))
	nl.NewRtAttrChild(data, nl.IFLA_GRE_OFLAGS, htons(gre.OFlags))

	if gre.Link != 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_GRE_LINK, nl.Uint32Attr(gre.Link))
	***REMOVED***

	nl.NewRtAttrChild(data, nl.IFLA_GRE_PMTUDISC, nl.Uint8Attr(gre.PMtuDisc))
	nl.NewRtAttrChild(data, nl.IFLA_GRE_TTL, nl.Uint8Attr(gre.Ttl))
	nl.NewRtAttrChild(data, nl.IFLA_GRE_TOS, nl.Uint8Attr(gre.Tos))
***REMOVED***

func parseGretunData(link Link, data []syscall.NetlinkRouteAttr) ***REMOVED***
	gre := link.(*Gretun)
	for _, datum := range data ***REMOVED***
		switch datum.Attr.Type ***REMOVED***
		case nl.IFLA_GRE_OKEY:
			gre.IKey = ntohl(datum.Value[0:4])
		case nl.IFLA_GRE_IKEY:
			gre.OKey = ntohl(datum.Value[0:4])
		case nl.IFLA_GRE_LOCAL:
			gre.Local = net.IP(datum.Value[0:4])
		case nl.IFLA_GRE_REMOTE:
			gre.Remote = net.IP(datum.Value[0:4])
		case nl.IFLA_GRE_IFLAGS:
			gre.IFlags = ntohs(datum.Value[0:2])
		case nl.IFLA_GRE_OFLAGS:
			gre.OFlags = ntohs(datum.Value[0:2])

		case nl.IFLA_GRE_TTL:
			gre.Ttl = uint8(datum.Value[0])
		case nl.IFLA_GRE_TOS:
			gre.Tos = uint8(datum.Value[0])
		case nl.IFLA_GRE_PMTUDISC:
			gre.PMtuDisc = uint8(datum.Value[0])
		***REMOVED***
	***REMOVED***
***REMOVED***

func parseLinkStats32(data []byte) *LinkStatistics ***REMOVED***
	return (*LinkStatistics)((*LinkStatistics32)(unsafe.Pointer(&data[0:SizeofLinkStats32][0])).to64())
***REMOVED***

func parseLinkStats64(data []byte) *LinkStatistics ***REMOVED***
	return (*LinkStatistics)((*LinkStatistics64)(unsafe.Pointer(&data[0:SizeofLinkStats64][0])))
***REMOVED***

func addXdpAttrs(xdp *LinkXdp, req *nl.NetlinkRequest) ***REMOVED***
	attrs := nl.NewRtAttr(nl.IFLA_XDP|syscall.NLA_F_NESTED, nil)
	b := make([]byte, 4)
	native.PutUint32(b, uint32(xdp.Fd))
	nl.NewRtAttrChild(attrs, nl.IFLA_XDP_FD, b)
	if xdp.Flags != 0 ***REMOVED***
		native.PutUint32(b, xdp.Flags)
		nl.NewRtAttrChild(attrs, nl.IFLA_XDP_FLAGS, b)
	***REMOVED***
	req.AddData(attrs)
***REMOVED***

func parseLinkXdp(data []byte) (*LinkXdp, error) ***REMOVED***
	attrs, err := nl.ParseRouteAttr(data)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	xdp := &LinkXdp***REMOVED******REMOVED***
	for _, attr := range attrs ***REMOVED***
		switch attr.Attr.Type ***REMOVED***
		case nl.IFLA_XDP_FD:
			xdp.Fd = int(native.Uint32(attr.Value[0:4]))
		case nl.IFLA_XDP_ATTACHED:
			xdp.Attached = attr.Value[0] != 0
		case nl.IFLA_XDP_FLAGS:
			xdp.Flags = native.Uint32(attr.Value[0:4])
		case nl.IFLA_XDP_PROG_ID:
			xdp.ProgId = native.Uint32(attr.Value[0:4])
		***REMOVED***
	***REMOVED***
	return xdp, nil
***REMOVED***

func addIptunAttrs(iptun *Iptun, linkInfo *nl.RtAttr) ***REMOVED***
	data := nl.NewRtAttrChild(linkInfo, nl.IFLA_INFO_DATA, nil)

	ip := iptun.Local.To4()
	if ip != nil ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_IPTUN_LOCAL, []byte(ip))
	***REMOVED***

	ip = iptun.Remote.To4()
	if ip != nil ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_IPTUN_REMOTE, []byte(ip))
	***REMOVED***

	if iptun.Link != 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_IPTUN_LINK, nl.Uint32Attr(iptun.Link))
	***REMOVED***
	nl.NewRtAttrChild(data, nl.IFLA_IPTUN_PMTUDISC, nl.Uint8Attr(iptun.PMtuDisc))
	nl.NewRtAttrChild(data, nl.IFLA_IPTUN_TTL, nl.Uint8Attr(iptun.Ttl))
	nl.NewRtAttrChild(data, nl.IFLA_IPTUN_TOS, nl.Uint8Attr(iptun.Tos))
***REMOVED***

func parseIptunData(link Link, data []syscall.NetlinkRouteAttr) ***REMOVED***
	iptun := link.(*Iptun)
	for _, datum := range data ***REMOVED***
		switch datum.Attr.Type ***REMOVED***
		case nl.IFLA_IPTUN_LOCAL:
			iptun.Local = net.IP(datum.Value[0:4])
		case nl.IFLA_IPTUN_REMOTE:
			iptun.Remote = net.IP(datum.Value[0:4])
		case nl.IFLA_IPTUN_TTL:
			iptun.Ttl = uint8(datum.Value[0])
		case nl.IFLA_IPTUN_TOS:
			iptun.Tos = uint8(datum.Value[0])
		case nl.IFLA_IPTUN_PMTUDISC:
			iptun.PMtuDisc = uint8(datum.Value[0])
		***REMOVED***
	***REMOVED***
***REMOVED***

func addVtiAttrs(vti *Vti, linkInfo *nl.RtAttr) ***REMOVED***
	data := nl.NewRtAttrChild(linkInfo, nl.IFLA_INFO_DATA, nil)

	ip := vti.Local.To4()
	if ip != nil ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_VTI_LOCAL, []byte(ip))
	***REMOVED***

	ip = vti.Remote.To4()
	if ip != nil ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_VTI_REMOTE, []byte(ip))
	***REMOVED***

	if vti.Link != 0 ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_VTI_LINK, nl.Uint32Attr(vti.Link))
	***REMOVED***

	nl.NewRtAttrChild(data, nl.IFLA_VTI_IKEY, htonl(vti.IKey))
	nl.NewRtAttrChild(data, nl.IFLA_VTI_OKEY, htonl(vti.OKey))
***REMOVED***

func parseVtiData(link Link, data []syscall.NetlinkRouteAttr) ***REMOVED***
	vti := link.(*Vti)
	for _, datum := range data ***REMOVED***
		switch datum.Attr.Type ***REMOVED***
		case nl.IFLA_VTI_LOCAL:
			vti.Local = net.IP(datum.Value[0:4])
		case nl.IFLA_VTI_REMOTE:
			vti.Remote = net.IP(datum.Value[0:4])
		case nl.IFLA_VTI_IKEY:
			vti.IKey = ntohl(datum.Value[0:4])
		case nl.IFLA_VTI_OKEY:
			vti.OKey = ntohl(datum.Value[0:4])
		***REMOVED***
	***REMOVED***
***REMOVED***

func addVrfAttrs(vrf *Vrf, linkInfo *nl.RtAttr) ***REMOVED***
	data := nl.NewRtAttrChild(linkInfo, nl.IFLA_INFO_DATA, nil)
	b := make([]byte, 4)
	native.PutUint32(b, uint32(vrf.Table))
	nl.NewRtAttrChild(data, nl.IFLA_VRF_TABLE, b)
***REMOVED***

func parseVrfData(link Link, data []syscall.NetlinkRouteAttr) ***REMOVED***
	vrf := link.(*Vrf)
	for _, datum := range data ***REMOVED***
		switch datum.Attr.Type ***REMOVED***
		case nl.IFLA_VRF_TABLE:
			vrf.Table = native.Uint32(datum.Value[0:4])
		***REMOVED***
	***REMOVED***
***REMOVED***

func addBridgeAttrs(bridge *Bridge, linkInfo *nl.RtAttr) ***REMOVED***
	data := nl.NewRtAttrChild(linkInfo, nl.IFLA_INFO_DATA, nil)
	if bridge.MulticastSnooping != nil ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BR_MCAST_SNOOPING, boolToByte(*bridge.MulticastSnooping))
	***REMOVED***
	if bridge.HelloTime != nil ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_BR_HELLO_TIME, nl.Uint32Attr(*bridge.HelloTime))
	***REMOVED***
***REMOVED***

func parseBridgeData(bridge Link, data []syscall.NetlinkRouteAttr) ***REMOVED***
	br := bridge.(*Bridge)
	for _, datum := range data ***REMOVED***
		switch datum.Attr.Type ***REMOVED***
		case nl.IFLA_BR_HELLO_TIME:
			helloTime := native.Uint32(datum.Value[0:4])
			br.HelloTime = &helloTime
		case nl.IFLA_BR_MCAST_SNOOPING:
			mcastSnooping := datum.Value[0] == 1
			br.MulticastSnooping = &mcastSnooping
		***REMOVED***
	***REMOVED***
***REMOVED***

func addGTPAttrs(gtp *GTP, linkInfo *nl.RtAttr) ***REMOVED***
	data := nl.NewRtAttrChild(linkInfo, nl.IFLA_INFO_DATA, nil)
	nl.NewRtAttrChild(data, nl.IFLA_GTP_FD0, nl.Uint32Attr(uint32(gtp.FD0)))
	nl.NewRtAttrChild(data, nl.IFLA_GTP_FD1, nl.Uint32Attr(uint32(gtp.FD1)))
	nl.NewRtAttrChild(data, nl.IFLA_GTP_PDP_HASHSIZE, nl.Uint32Attr(131072))
	if gtp.Role != nl.GTP_ROLE_GGSN ***REMOVED***
		nl.NewRtAttrChild(data, nl.IFLA_GTP_ROLE, nl.Uint32Attr(uint32(gtp.Role)))
	***REMOVED***
***REMOVED***

func parseGTPData(link Link, data []syscall.NetlinkRouteAttr) ***REMOVED***
	gtp := link.(*GTP)
	for _, datum := range data ***REMOVED***
		switch datum.Attr.Type ***REMOVED***
		case nl.IFLA_GTP_FD0:
			gtp.FD0 = int(native.Uint32(datum.Value))
		case nl.IFLA_GTP_FD1:
			gtp.FD1 = int(native.Uint32(datum.Value))
		case nl.IFLA_GTP_PDP_HASHSIZE:
			gtp.PDPHashsize = int(native.Uint32(datum.Value))
		case nl.IFLA_GTP_ROLE:
			gtp.Role = int(native.Uint32(datum.Value))
		***REMOVED***
	***REMOVED***
***REMOVED***
