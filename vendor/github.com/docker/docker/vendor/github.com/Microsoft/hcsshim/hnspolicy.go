package hcsshim

// Type of Request Support in ModifySystem
type PolicyType string

// RequestType const
const (
	Nat                  PolicyType = "NAT"
	ACL                  PolicyType = "ACL"
	PA                   PolicyType = "PA"
	VLAN                 PolicyType = "VLAN"
	VSID                 PolicyType = "VSID"
	VNet                 PolicyType = "VNET"
	L2Driver             PolicyType = "L2Driver"
	Isolation            PolicyType = "Isolation"
	QOS                  PolicyType = "QOS"
	OutboundNat          PolicyType = "OutBoundNAT"
	ExternalLoadBalancer PolicyType = "ELB"
	Route                PolicyType = "ROUTE"
)

type NatPolicy struct ***REMOVED***
	Type         PolicyType `json:"Type"`
	Protocol     string
	InternalPort uint16
	ExternalPort uint16
***REMOVED***

type QosPolicy struct ***REMOVED***
	Type                            PolicyType `json:"Type"`
	MaximumOutgoingBandwidthInBytes uint64
***REMOVED***

type IsolationPolicy struct ***REMOVED***
	Type               PolicyType `json:"Type"`
	VLAN               uint
	VSID               uint
	InDefaultIsolation bool
***REMOVED***

type VlanPolicy struct ***REMOVED***
	Type PolicyType `json:"Type"`
	VLAN uint
***REMOVED***

type VsidPolicy struct ***REMOVED***
	Type PolicyType `json:"Type"`
	VSID uint
***REMOVED***

type PaPolicy struct ***REMOVED***
	Type PolicyType `json:"Type"`
	PA   string     `json:"PA"`
***REMOVED***

type OutboundNatPolicy struct ***REMOVED***
	Policy
	VIP        string   `json:"VIP,omitempty"`
	Exceptions []string `json:"ExceptionList,omitempty"`
***REMOVED***

type ActionType string
type DirectionType string
type RuleType string

const (
	Allow ActionType = "Allow"
	Block ActionType = "Block"

	In  DirectionType = "In"
	Out DirectionType = "Out"

	Host   RuleType = "Host"
	Switch RuleType = "Switch"
)

type ACLPolicy struct ***REMOVED***
	Type            PolicyType `json:"Type"`
	Protocol        uint16
	InternalPort    uint16
	Action          ActionType
	Direction       DirectionType
	LocalAddresses  string
	RemoteAddresses string
	LocalPort       uint16
	RemotePort      uint16
	RuleType        RuleType `json:"RuleType,omitempty"`
	Priority        uint16
	ServiceName     string
***REMOVED***

type Policy struct ***REMOVED***
	Type PolicyType `json:"Type"`
***REMOVED***
