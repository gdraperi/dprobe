package netlink

import (
	"fmt"
	"syscall"

	"github.com/vishvananda/netlink/nl"
)

// BridgeVlanList gets a map of device id to bridge vlan infos.
// Equivalent to: `bridge vlan show`
func BridgeVlanList() (map[int32][]*nl.BridgeVlanInfo, error) ***REMOVED***
	return pkgHandle.BridgeVlanList()
***REMOVED***

// BridgeVlanList gets a map of device id to bridge vlan infos.
// Equivalent to: `bridge vlan show`
func (h *Handle) BridgeVlanList() (map[int32][]*nl.BridgeVlanInfo, error) ***REMOVED***
	req := h.newNetlinkRequest(syscall.RTM_GETLINK, syscall.NLM_F_DUMP)
	msg := nl.NewIfInfomsg(syscall.AF_BRIDGE)
	req.AddData(msg)
	req.AddData(nl.NewRtAttr(nl.IFLA_EXT_MASK, nl.Uint32Attr(uint32(nl.RTEXT_FILTER_BRVLAN))))

	msgs, err := req.Execute(syscall.NETLINK_ROUTE, syscall.RTM_NEWLINK)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ret := make(map[int32][]*nl.BridgeVlanInfo)
	for _, m := range msgs ***REMOVED***
		msg := nl.DeserializeIfInfomsg(m)

		attrs, err := nl.ParseRouteAttr(m[msg.Len():])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		for _, attr := range attrs ***REMOVED***
			switch attr.Attr.Type ***REMOVED***
			case nl.IFLA_AF_SPEC:
				//nested attr
				nestAttrs, err := nl.ParseRouteAttr(attr.Value)
				if err != nil ***REMOVED***
					return nil, fmt.Errorf("failed to parse nested attr %v", err)
				***REMOVED***
				for _, nestAttr := range nestAttrs ***REMOVED***
					switch nestAttr.Attr.Type ***REMOVED***
					case nl.IFLA_BRIDGE_VLAN_INFO:
						vlanInfo := nl.DeserializeBridgeVlanInfo(nestAttr.Value)
						ret[msg.Index] = append(ret[msg.Index], vlanInfo)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return ret, nil
***REMOVED***

// BridgeVlanAdd adds a new vlan filter entry
// Equivalent to: `bridge vlan add dev DEV vid VID [ pvid ] [ untagged ] [ self ] [ master ]`
func BridgeVlanAdd(link Link, vid uint16, pvid, untagged, self, master bool) error ***REMOVED***
	return pkgHandle.BridgeVlanAdd(link, vid, pvid, untagged, self, master)
***REMOVED***

// BridgeVlanAdd adds a new vlan filter entry
// Equivalent to: `bridge vlan add dev DEV vid VID [ pvid ] [ untagged ] [ self ] [ master ]`
func (h *Handle) BridgeVlanAdd(link Link, vid uint16, pvid, untagged, self, master bool) error ***REMOVED***
	return h.bridgeVlanModify(syscall.RTM_SETLINK, link, vid, pvid, untagged, self, master)
***REMOVED***

// BridgeVlanDel adds a new vlan filter entry
// Equivalent to: `bridge vlan del dev DEV vid VID [ pvid ] [ untagged ] [ self ] [ master ]`
func BridgeVlanDel(link Link, vid uint16, pvid, untagged, self, master bool) error ***REMOVED***
	return pkgHandle.BridgeVlanDel(link, vid, pvid, untagged, self, master)
***REMOVED***

// BridgeVlanDel adds a new vlan filter entry
// Equivalent to: `bridge vlan del dev DEV vid VID [ pvid ] [ untagged ] [ self ] [ master ]`
func (h *Handle) BridgeVlanDel(link Link, vid uint16, pvid, untagged, self, master bool) error ***REMOVED***
	return h.bridgeVlanModify(syscall.RTM_DELLINK, link, vid, pvid, untagged, self, master)
***REMOVED***

func (h *Handle) bridgeVlanModify(cmd int, link Link, vid uint16, pvid, untagged, self, master bool) error ***REMOVED***
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(cmd, syscall.NLM_F_ACK)

	msg := nl.NewIfInfomsg(syscall.AF_BRIDGE)
	msg.Index = int32(base.Index)
	req.AddData(msg)

	br := nl.NewRtAttr(nl.IFLA_AF_SPEC, nil)
	var flags uint16
	if self ***REMOVED***
		flags |= nl.BRIDGE_FLAGS_SELF
	***REMOVED***
	if master ***REMOVED***
		flags |= nl.BRIDGE_FLAGS_MASTER
	***REMOVED***
	if flags > 0 ***REMOVED***
		nl.NewRtAttrChild(br, nl.IFLA_BRIDGE_FLAGS, nl.Uint16Attr(flags))
	***REMOVED***
	vlanInfo := &nl.BridgeVlanInfo***REMOVED***Vid: vid***REMOVED***
	if pvid ***REMOVED***
		vlanInfo.Flags |= nl.BRIDGE_VLAN_INFO_PVID
	***REMOVED***
	if untagged ***REMOVED***
		vlanInfo.Flags |= nl.BRIDGE_VLAN_INFO_UNTAGGED
	***REMOVED***
	nl.NewRtAttrChild(br, nl.IFLA_BRIDGE_VLAN_INFO, vlanInfo.Serialize())
	req.AddData(br)
	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***
