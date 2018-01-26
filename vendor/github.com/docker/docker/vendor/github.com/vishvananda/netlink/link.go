package netlink

import (
	"fmt"
	"net"
)

// Link represents a link device from netlink. Shared link attributes
// like name may be retrieved using the Attrs() method. Unique data
// can be retrieved by casting the object to the proper type.
type Link interface ***REMOVED***
	Attrs() *LinkAttrs
	Type() string
***REMOVED***

type (
	NsPid int
	NsFd  int
)

// LinkAttrs represents data shared by most link types
type LinkAttrs struct ***REMOVED***
	Index        int
	MTU          int
	TxQLen       int // Transmit Queue Length
	Name         string
	HardwareAddr net.HardwareAddr
	Flags        net.Flags
	RawFlags     uint32
	ParentIndex  int         // index of the parent link device
	MasterIndex  int         // must be the index of a bridge
	Namespace    interface***REMOVED******REMOVED*** // nil | NsPid | NsFd
	Alias        string
	Statistics   *LinkStatistics
	Promisc      int
	Xdp          *LinkXdp
	EncapType    string
	Protinfo     *Protinfo
	OperState    LinkOperState
	NetNsID      int
***REMOVED***

// LinkOperState represents the values of the IFLA_OPERSTATE link
// attribute, which contains the RFC2863 state of the interface.
type LinkOperState uint8

const (
	OperUnknown        = iota // Status can't be determined.
	OperNotPresent            // Some component is missing.
	OperDown                  // Down.
	OperLowerLayerDown        // Down due to state of lower layer.
	OperTesting               // In some test mode.
	OperDormant               // Not up but pending an external event.
	OperUp                    // Up, ready to send packets.
)

func (s LinkOperState) String() string ***REMOVED***
	switch s ***REMOVED***
	case OperNotPresent:
		return "not-present"
	case OperDown:
		return "down"
	case OperLowerLayerDown:
		return "lower-layer-down"
	case OperTesting:
		return "testing"
	case OperDormant:
		return "dormant"
	case OperUp:
		return "up"
	default:
		return "unknown"
	***REMOVED***
***REMOVED***

// NewLinkAttrs returns LinkAttrs structure filled with default values
func NewLinkAttrs() LinkAttrs ***REMOVED***
	return LinkAttrs***REMOVED***
		TxQLen: -1,
	***REMOVED***
***REMOVED***

type LinkStatistics LinkStatistics64

/*
Ref: struct rtnl_link_stats ***REMOVED***...***REMOVED***
*/
type LinkStatistics32 struct ***REMOVED***
	RxPackets         uint32
	TxPackets         uint32
	RxBytes           uint32
	TxBytes           uint32
	RxErrors          uint32
	TxErrors          uint32
	RxDropped         uint32
	TxDropped         uint32
	Multicast         uint32
	Collisions        uint32
	RxLengthErrors    uint32
	RxOverErrors      uint32
	RxCrcErrors       uint32
	RxFrameErrors     uint32
	RxFifoErrors      uint32
	RxMissedErrors    uint32
	TxAbortedErrors   uint32
	TxCarrierErrors   uint32
	TxFifoErrors      uint32
	TxHeartbeatErrors uint32
	TxWindowErrors    uint32
	RxCompressed      uint32
	TxCompressed      uint32
***REMOVED***

func (s32 LinkStatistics32) to64() *LinkStatistics64 ***REMOVED***
	return &LinkStatistics64***REMOVED***
		RxPackets:         uint64(s32.RxPackets),
		TxPackets:         uint64(s32.TxPackets),
		RxBytes:           uint64(s32.RxBytes),
		TxBytes:           uint64(s32.TxBytes),
		RxErrors:          uint64(s32.RxErrors),
		TxErrors:          uint64(s32.TxErrors),
		RxDropped:         uint64(s32.RxDropped),
		TxDropped:         uint64(s32.TxDropped),
		Multicast:         uint64(s32.Multicast),
		Collisions:        uint64(s32.Collisions),
		RxLengthErrors:    uint64(s32.RxLengthErrors),
		RxOverErrors:      uint64(s32.RxOverErrors),
		RxCrcErrors:       uint64(s32.RxCrcErrors),
		RxFrameErrors:     uint64(s32.RxFrameErrors),
		RxFifoErrors:      uint64(s32.RxFifoErrors),
		RxMissedErrors:    uint64(s32.RxMissedErrors),
		TxAbortedErrors:   uint64(s32.TxAbortedErrors),
		TxCarrierErrors:   uint64(s32.TxCarrierErrors),
		TxFifoErrors:      uint64(s32.TxFifoErrors),
		TxHeartbeatErrors: uint64(s32.TxHeartbeatErrors),
		TxWindowErrors:    uint64(s32.TxWindowErrors),
		RxCompressed:      uint64(s32.RxCompressed),
		TxCompressed:      uint64(s32.TxCompressed),
	***REMOVED***
***REMOVED***

/*
Ref: struct rtnl_link_stats64 ***REMOVED***...***REMOVED***
*/
type LinkStatistics64 struct ***REMOVED***
	RxPackets         uint64
	TxPackets         uint64
	RxBytes           uint64
	TxBytes           uint64
	RxErrors          uint64
	TxErrors          uint64
	RxDropped         uint64
	TxDropped         uint64
	Multicast         uint64
	Collisions        uint64
	RxLengthErrors    uint64
	RxOverErrors      uint64
	RxCrcErrors       uint64
	RxFrameErrors     uint64
	RxFifoErrors      uint64
	RxMissedErrors    uint64
	TxAbortedErrors   uint64
	TxCarrierErrors   uint64
	TxFifoErrors      uint64
	TxHeartbeatErrors uint64
	TxWindowErrors    uint64
	RxCompressed      uint64
	TxCompressed      uint64
***REMOVED***

type LinkXdp struct ***REMOVED***
	Fd       int
	Attached bool
	Flags    uint32
	ProgId   uint32
***REMOVED***

// Device links cannot be created via netlink. These links
// are links created by udev like 'lo' and 'etho0'
type Device struct ***REMOVED***
	LinkAttrs
***REMOVED***

func (device *Device) Attrs() *LinkAttrs ***REMOVED***
	return &device.LinkAttrs
***REMOVED***

func (device *Device) Type() string ***REMOVED***
	return "device"
***REMOVED***

// Dummy links are dummy ethernet devices
type Dummy struct ***REMOVED***
	LinkAttrs
***REMOVED***

func (dummy *Dummy) Attrs() *LinkAttrs ***REMOVED***
	return &dummy.LinkAttrs
***REMOVED***

func (dummy *Dummy) Type() string ***REMOVED***
	return "dummy"
***REMOVED***

// Ifb links are advanced dummy devices for packet filtering
type Ifb struct ***REMOVED***
	LinkAttrs
***REMOVED***

func (ifb *Ifb) Attrs() *LinkAttrs ***REMOVED***
	return &ifb.LinkAttrs
***REMOVED***

func (ifb *Ifb) Type() string ***REMOVED***
	return "ifb"
***REMOVED***

// Bridge links are simple linux bridges
type Bridge struct ***REMOVED***
	LinkAttrs
	MulticastSnooping *bool
	HelloTime         *uint32
***REMOVED***

func (bridge *Bridge) Attrs() *LinkAttrs ***REMOVED***
	return &bridge.LinkAttrs
***REMOVED***

func (bridge *Bridge) Type() string ***REMOVED***
	return "bridge"
***REMOVED***

// Vlan links have ParentIndex set in their Attrs()
type Vlan struct ***REMOVED***
	LinkAttrs
	VlanId int
***REMOVED***

func (vlan *Vlan) Attrs() *LinkAttrs ***REMOVED***
	return &vlan.LinkAttrs
***REMOVED***

func (vlan *Vlan) Type() string ***REMOVED***
	return "vlan"
***REMOVED***

type MacvlanMode uint16

const (
	MACVLAN_MODE_DEFAULT MacvlanMode = iota
	MACVLAN_MODE_PRIVATE
	MACVLAN_MODE_VEPA
	MACVLAN_MODE_BRIDGE
	MACVLAN_MODE_PASSTHRU
	MACVLAN_MODE_SOURCE
)

// Macvlan links have ParentIndex set in their Attrs()
type Macvlan struct ***REMOVED***
	LinkAttrs
	Mode MacvlanMode
***REMOVED***

func (macvlan *Macvlan) Attrs() *LinkAttrs ***REMOVED***
	return &macvlan.LinkAttrs
***REMOVED***

func (macvlan *Macvlan) Type() string ***REMOVED***
	return "macvlan"
***REMOVED***

// Macvtap - macvtap is a virtual interfaces based on macvlan
type Macvtap struct ***REMOVED***
	Macvlan
***REMOVED***

func (macvtap Macvtap) Type() string ***REMOVED***
	return "macvtap"
***REMOVED***

type TuntapMode uint16
type TuntapFlag uint16

// Tuntap links created via /dev/tun/tap, but can be destroyed via netlink
type Tuntap struct ***REMOVED***
	LinkAttrs
	Mode  TuntapMode
	Flags TuntapFlag
***REMOVED***

func (tuntap *Tuntap) Attrs() *LinkAttrs ***REMOVED***
	return &tuntap.LinkAttrs
***REMOVED***

func (tuntap *Tuntap) Type() string ***REMOVED***
	return "tuntap"
***REMOVED***

// Veth devices must specify PeerName on create
type Veth struct ***REMOVED***
	LinkAttrs
	PeerName string // veth on create only
***REMOVED***

func (veth *Veth) Attrs() *LinkAttrs ***REMOVED***
	return &veth.LinkAttrs
***REMOVED***

func (veth *Veth) Type() string ***REMOVED***
	return "veth"
***REMOVED***

// GenericLink links represent types that are not currently understood
// by this netlink library.
type GenericLink struct ***REMOVED***
	LinkAttrs
	LinkType string
***REMOVED***

func (generic *GenericLink) Attrs() *LinkAttrs ***REMOVED***
	return &generic.LinkAttrs
***REMOVED***

func (generic *GenericLink) Type() string ***REMOVED***
	return generic.LinkType
***REMOVED***

type Vxlan struct ***REMOVED***
	LinkAttrs
	VxlanId      int
	VtepDevIndex int
	SrcAddr      net.IP
	Group        net.IP
	TTL          int
	TOS          int
	Learning     bool
	Proxy        bool
	RSC          bool
	L2miss       bool
	L3miss       bool
	UDPCSum      bool
	NoAge        bool
	GBP          bool
	FlowBased    bool
	Age          int
	Limit        int
	Port         int
	PortLow      int
	PortHigh     int
***REMOVED***

func (vxlan *Vxlan) Attrs() *LinkAttrs ***REMOVED***
	return &vxlan.LinkAttrs
***REMOVED***

func (vxlan *Vxlan) Type() string ***REMOVED***
	return "vxlan"
***REMOVED***

type IPVlanMode uint16

const (
	IPVLAN_MODE_L2 IPVlanMode = iota
	IPVLAN_MODE_L3
	IPVLAN_MODE_L3S
	IPVLAN_MODE_MAX
)

type IPVlan struct ***REMOVED***
	LinkAttrs
	Mode IPVlanMode
***REMOVED***

func (ipvlan *IPVlan) Attrs() *LinkAttrs ***REMOVED***
	return &ipvlan.LinkAttrs
***REMOVED***

func (ipvlan *IPVlan) Type() string ***REMOVED***
	return "ipvlan"
***REMOVED***

// BondMode type
type BondMode int

func (b BondMode) String() string ***REMOVED***
	s, ok := bondModeToString[b]
	if !ok ***REMOVED***
		return fmt.Sprintf("BondMode(%d)", b)
	***REMOVED***
	return s
***REMOVED***

// StringToBondMode returns bond mode, or uknonw is the s is invalid.
func StringToBondMode(s string) BondMode ***REMOVED***
	mode, ok := StringToBondModeMap[s]
	if !ok ***REMOVED***
		return BOND_MODE_UNKNOWN
	***REMOVED***
	return mode
***REMOVED***

// Possible BondMode
const (
	BOND_MODE_BALANCE_RR BondMode = iota
	BOND_MODE_ACTIVE_BACKUP
	BOND_MODE_BALANCE_XOR
	BOND_MODE_BROADCAST
	BOND_MODE_802_3AD
	BOND_MODE_BALANCE_TLB
	BOND_MODE_BALANCE_ALB
	BOND_MODE_UNKNOWN
)

var bondModeToString = map[BondMode]string***REMOVED***
	BOND_MODE_BALANCE_RR:    "balance-rr",
	BOND_MODE_ACTIVE_BACKUP: "active-backup",
	BOND_MODE_BALANCE_XOR:   "balance-xor",
	BOND_MODE_BROADCAST:     "broadcast",
	BOND_MODE_802_3AD:       "802.3ad",
	BOND_MODE_BALANCE_TLB:   "balance-tlb",
	BOND_MODE_BALANCE_ALB:   "balance-alb",
***REMOVED***
var StringToBondModeMap = map[string]BondMode***REMOVED***
	"balance-rr":    BOND_MODE_BALANCE_RR,
	"active-backup": BOND_MODE_ACTIVE_BACKUP,
	"balance-xor":   BOND_MODE_BALANCE_XOR,
	"broadcast":     BOND_MODE_BROADCAST,
	"802.3ad":       BOND_MODE_802_3AD,
	"balance-tlb":   BOND_MODE_BALANCE_TLB,
	"balance-alb":   BOND_MODE_BALANCE_ALB,
***REMOVED***

// BondArpValidate type
type BondArpValidate int

// Possible BondArpValidate value
const (
	BOND_ARP_VALIDATE_NONE BondArpValidate = iota
	BOND_ARP_VALIDATE_ACTIVE
	BOND_ARP_VALIDATE_BACKUP
	BOND_ARP_VALIDATE_ALL
)

// BondPrimaryReselect type
type BondPrimaryReselect int

// Possible BondPrimaryReselect value
const (
	BOND_PRIMARY_RESELECT_ALWAYS BondPrimaryReselect = iota
	BOND_PRIMARY_RESELECT_BETTER
	BOND_PRIMARY_RESELECT_FAILURE
)

// BondArpAllTargets type
type BondArpAllTargets int

// Possible BondArpAllTargets value
const (
	BOND_ARP_ALL_TARGETS_ANY BondArpAllTargets = iota
	BOND_ARP_ALL_TARGETS_ALL
)

// BondFailOverMac type
type BondFailOverMac int

// Possible BondFailOverMac value
const (
	BOND_FAIL_OVER_MAC_NONE BondFailOverMac = iota
	BOND_FAIL_OVER_MAC_ACTIVE
	BOND_FAIL_OVER_MAC_FOLLOW
)

// BondXmitHashPolicy type
type BondXmitHashPolicy int

func (b BondXmitHashPolicy) String() string ***REMOVED***
	s, ok := bondXmitHashPolicyToString[b]
	if !ok ***REMOVED***
		return fmt.Sprintf("XmitHashPolicy(%d)", b)
	***REMOVED***
	return s
***REMOVED***

// StringToBondXmitHashPolicy returns bond lacp arte, or uknonw is the s is invalid.
func StringToBondXmitHashPolicy(s string) BondXmitHashPolicy ***REMOVED***
	lacp, ok := StringToBondXmitHashPolicyMap[s]
	if !ok ***REMOVED***
		return BOND_XMIT_HASH_POLICY_UNKNOWN
	***REMOVED***
	return lacp
***REMOVED***

// Possible BondXmitHashPolicy value
const (
	BOND_XMIT_HASH_POLICY_LAYER2 BondXmitHashPolicy = iota
	BOND_XMIT_HASH_POLICY_LAYER3_4
	BOND_XMIT_HASH_POLICY_LAYER2_3
	BOND_XMIT_HASH_POLICY_ENCAP2_3
	BOND_XMIT_HASH_POLICY_ENCAP3_4
	BOND_XMIT_HASH_POLICY_UNKNOWN
)

var bondXmitHashPolicyToString = map[BondXmitHashPolicy]string***REMOVED***
	BOND_XMIT_HASH_POLICY_LAYER2:   "layer2",
	BOND_XMIT_HASH_POLICY_LAYER3_4: "layer3+4",
	BOND_XMIT_HASH_POLICY_LAYER2_3: "layer2+3",
	BOND_XMIT_HASH_POLICY_ENCAP2_3: "encap2+3",
	BOND_XMIT_HASH_POLICY_ENCAP3_4: "encap3+4",
***REMOVED***
var StringToBondXmitHashPolicyMap = map[string]BondXmitHashPolicy***REMOVED***
	"layer2":   BOND_XMIT_HASH_POLICY_LAYER2,
	"layer3+4": BOND_XMIT_HASH_POLICY_LAYER3_4,
	"layer2+3": BOND_XMIT_HASH_POLICY_LAYER2_3,
	"encap2+3": BOND_XMIT_HASH_POLICY_ENCAP2_3,
	"encap3+4": BOND_XMIT_HASH_POLICY_ENCAP3_4,
***REMOVED***

// BondLacpRate type
type BondLacpRate int

func (b BondLacpRate) String() string ***REMOVED***
	s, ok := bondLacpRateToString[b]
	if !ok ***REMOVED***
		return fmt.Sprintf("LacpRate(%d)", b)
	***REMOVED***
	return s
***REMOVED***

// StringToBondLacpRate returns bond lacp arte, or uknonw is the s is invalid.
func StringToBondLacpRate(s string) BondLacpRate ***REMOVED***
	lacp, ok := StringToBondLacpRateMap[s]
	if !ok ***REMOVED***
		return BOND_LACP_RATE_UNKNOWN
	***REMOVED***
	return lacp
***REMOVED***

// Possible BondLacpRate value
const (
	BOND_LACP_RATE_SLOW BondLacpRate = iota
	BOND_LACP_RATE_FAST
	BOND_LACP_RATE_UNKNOWN
)

var bondLacpRateToString = map[BondLacpRate]string***REMOVED***
	BOND_LACP_RATE_SLOW: "slow",
	BOND_LACP_RATE_FAST: "fast",
***REMOVED***
var StringToBondLacpRateMap = map[string]BondLacpRate***REMOVED***
	"slow": BOND_LACP_RATE_SLOW,
	"fast": BOND_LACP_RATE_FAST,
***REMOVED***

// BondAdSelect type
type BondAdSelect int

// Possible BondAdSelect value
const (
	BOND_AD_SELECT_STABLE BondAdSelect = iota
	BOND_AD_SELECT_BANDWIDTH
	BOND_AD_SELECT_COUNT
)

// BondAdInfo represents ad info for bond
type BondAdInfo struct ***REMOVED***
	AggregatorId int
	NumPorts     int
	ActorKey     int
	PartnerKey   int
	PartnerMac   net.HardwareAddr
***REMOVED***

// Bond representation
type Bond struct ***REMOVED***
	LinkAttrs
	Mode            BondMode
	ActiveSlave     int
	Miimon          int
	UpDelay         int
	DownDelay       int
	UseCarrier      int
	ArpInterval     int
	ArpIpTargets    []net.IP
	ArpValidate     BondArpValidate
	ArpAllTargets   BondArpAllTargets
	Primary         int
	PrimaryReselect BondPrimaryReselect
	FailOverMac     BondFailOverMac
	XmitHashPolicy  BondXmitHashPolicy
	ResendIgmp      int
	NumPeerNotif    int
	AllSlavesActive int
	MinLinks        int
	LpInterval      int
	PackersPerSlave int
	LacpRate        BondLacpRate
	AdSelect        BondAdSelect
	// looking at iproute tool AdInfo can only be retrived. It can't be set.
	AdInfo         *BondAdInfo
	AdActorSysPrio int
	AdUserPortKey  int
	AdActorSystem  net.HardwareAddr
	TlbDynamicLb   int
***REMOVED***

func NewLinkBond(atr LinkAttrs) *Bond ***REMOVED***
	return &Bond***REMOVED***
		LinkAttrs:       atr,
		Mode:            -1,
		ActiveSlave:     -1,
		Miimon:          -1,
		UpDelay:         -1,
		DownDelay:       -1,
		UseCarrier:      -1,
		ArpInterval:     -1,
		ArpIpTargets:    nil,
		ArpValidate:     -1,
		ArpAllTargets:   -1,
		Primary:         -1,
		PrimaryReselect: -1,
		FailOverMac:     -1,
		XmitHashPolicy:  -1,
		ResendIgmp:      -1,
		NumPeerNotif:    -1,
		AllSlavesActive: -1,
		MinLinks:        -1,
		LpInterval:      -1,
		PackersPerSlave: -1,
		LacpRate:        -1,
		AdSelect:        -1,
		AdActorSysPrio:  -1,
		AdUserPortKey:   -1,
		AdActorSystem:   nil,
		TlbDynamicLb:    -1,
	***REMOVED***
***REMOVED***

// Flag mask for bond options. Bond.Flagmask must be set to on for option to work.
const (
	BOND_MODE_MASK uint64 = 1 << (1 + iota)
	BOND_ACTIVE_SLAVE_MASK
	BOND_MIIMON_MASK
	BOND_UPDELAY_MASK
	BOND_DOWNDELAY_MASK
	BOND_USE_CARRIER_MASK
	BOND_ARP_INTERVAL_MASK
	BOND_ARP_VALIDATE_MASK
	BOND_ARP_ALL_TARGETS_MASK
	BOND_PRIMARY_MASK
	BOND_PRIMARY_RESELECT_MASK
	BOND_FAIL_OVER_MAC_MASK
	BOND_XMIT_HASH_POLICY_MASK
	BOND_RESEND_IGMP_MASK
	BOND_NUM_PEER_NOTIF_MASK
	BOND_ALL_SLAVES_ACTIVE_MASK
	BOND_MIN_LINKS_MASK
	BOND_LP_INTERVAL_MASK
	BOND_PACKETS_PER_SLAVE_MASK
	BOND_LACP_RATE_MASK
	BOND_AD_SELECT_MASK
)

// Attrs implementation.
func (bond *Bond) Attrs() *LinkAttrs ***REMOVED***
	return &bond.LinkAttrs
***REMOVED***

// Type implementation fro Vxlan.
func (bond *Bond) Type() string ***REMOVED***
	return "bond"
***REMOVED***

// Gretap devices must specify LocalIP and RemoteIP on create
type Gretap struct ***REMOVED***
	LinkAttrs
	IKey       uint32
	OKey       uint32
	EncapSport uint16
	EncapDport uint16
	Local      net.IP
	Remote     net.IP
	IFlags     uint16
	OFlags     uint16
	PMtuDisc   uint8
	Ttl        uint8
	Tos        uint8
	EncapType  uint16
	EncapFlags uint16
	Link       uint32
	FlowBased  bool
***REMOVED***

func (gretap *Gretap) Attrs() *LinkAttrs ***REMOVED***
	return &gretap.LinkAttrs
***REMOVED***

func (gretap *Gretap) Type() string ***REMOVED***
	return "gretap"
***REMOVED***

type Iptun struct ***REMOVED***
	LinkAttrs
	Ttl      uint8
	Tos      uint8
	PMtuDisc uint8
	Link     uint32
	Local    net.IP
	Remote   net.IP
***REMOVED***

func (iptun *Iptun) Attrs() *LinkAttrs ***REMOVED***
	return &iptun.LinkAttrs
***REMOVED***

func (iptun *Iptun) Type() string ***REMOVED***
	return "ipip"
***REMOVED***

type Vti struct ***REMOVED***
	LinkAttrs
	IKey   uint32
	OKey   uint32
	Link   uint32
	Local  net.IP
	Remote net.IP
***REMOVED***

func (vti *Vti) Attrs() *LinkAttrs ***REMOVED***
	return &vti.LinkAttrs
***REMOVED***

func (iptun *Vti) Type() string ***REMOVED***
	return "vti"
***REMOVED***

type Gretun struct ***REMOVED***
	LinkAttrs
	Link     uint32
	IFlags   uint16
	OFlags   uint16
	IKey     uint32
	OKey     uint32
	Local    net.IP
	Remote   net.IP
	Ttl      uint8
	Tos      uint8
	PMtuDisc uint8
***REMOVED***

func (gretun *Gretun) Attrs() *LinkAttrs ***REMOVED***
	return &gretun.LinkAttrs
***REMOVED***

func (gretun *Gretun) Type() string ***REMOVED***
	return "gre"
***REMOVED***

type Vrf struct ***REMOVED***
	LinkAttrs
	Table uint32
***REMOVED***

func (vrf *Vrf) Attrs() *LinkAttrs ***REMOVED***
	return &vrf.LinkAttrs
***REMOVED***

func (vrf *Vrf) Type() string ***REMOVED***
	return "vrf"
***REMOVED***

type GTP struct ***REMOVED***
	LinkAttrs
	FD0         int
	FD1         int
	Role        int
	PDPHashsize int
***REMOVED***

func (gtp *GTP) Attrs() *LinkAttrs ***REMOVED***
	return &gtp.LinkAttrs
***REMOVED***

func (gtp *GTP) Type() string ***REMOVED***
	return "gtp"
***REMOVED***

// iproute2 supported devices;
// vlan | veth | vcan | dummy | ifb | macvlan | macvtap |
// bridge | bond | ipoib | ip6tnl | ipip | sit | vxlan |
// gre | gretap | ip6gre | ip6gretap | vti | nlmon |
// bond_slave | ipvlan

// LinkNotFoundError wraps the various not found errors when
// getting/reading links. This is intended for better error
// handling by dependent code so that "not found error" can
// be distinguished from other errors
type LinkNotFoundError struct ***REMOVED***
	error
***REMOVED***
