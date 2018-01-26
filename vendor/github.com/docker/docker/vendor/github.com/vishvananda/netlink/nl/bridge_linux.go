package nl

import (
	"fmt"
	"unsafe"
)

const (
	SizeofBridgeVlanInfo = 0x04
)

/* Bridge Flags */
const (
	BRIDGE_FLAGS_MASTER = iota /* Bridge command to/from master */
	BRIDGE_FLAGS_SELF          /* Bridge command to/from lowerdev */
)

/* Bridge management nested attributes
 * [IFLA_AF_SPEC] = ***REMOVED***
 *     [IFLA_BRIDGE_FLAGS]
 *     [IFLA_BRIDGE_MODE]
 *     [IFLA_BRIDGE_VLAN_INFO]
 * ***REMOVED***
 */
const (
	IFLA_BRIDGE_FLAGS = iota
	IFLA_BRIDGE_MODE
	IFLA_BRIDGE_VLAN_INFO
)

const (
	BRIDGE_VLAN_INFO_MASTER = 1 << iota
	BRIDGE_VLAN_INFO_PVID
	BRIDGE_VLAN_INFO_UNTAGGED
	BRIDGE_VLAN_INFO_RANGE_BEGIN
	BRIDGE_VLAN_INFO_RANGE_END
)

// struct bridge_vlan_info ***REMOVED***
//   __u16 flags;
//   __u16 vid;
// ***REMOVED***;

type BridgeVlanInfo struct ***REMOVED***
	Flags uint16
	Vid   uint16
***REMOVED***

func (b *BridgeVlanInfo) Serialize() []byte ***REMOVED***
	return (*(*[SizeofBridgeVlanInfo]byte)(unsafe.Pointer(b)))[:]
***REMOVED***

func DeserializeBridgeVlanInfo(b []byte) *BridgeVlanInfo ***REMOVED***
	return (*BridgeVlanInfo)(unsafe.Pointer(&b[0:SizeofBridgeVlanInfo][0]))
***REMOVED***

func (b *BridgeVlanInfo) PortVID() bool ***REMOVED***
	return b.Flags&BRIDGE_VLAN_INFO_PVID > 0
***REMOVED***

func (b *BridgeVlanInfo) EngressUntag() bool ***REMOVED***
	return b.Flags&BRIDGE_VLAN_INFO_UNTAGGED > 0
***REMOVED***

func (b *BridgeVlanInfo) String() string ***REMOVED***
	return fmt.Sprintf("%+v", *b)
***REMOVED***

/* New extended info filters for IFLA_EXT_MASK */
const (
	RTEXT_FILTER_VF = 1 << iota
	RTEXT_FILTER_BRVLAN
	RTEXT_FILTER_BRVLAN_COMPRESSED
)
