package netlink

import (
	"fmt"
)

type Filter interface ***REMOVED***
	Attrs() *FilterAttrs
	Type() string
***REMOVED***

// FilterAttrs represents a netlink filter. A filter is associated with a link,
// has a handle and a parent. The root filter of a device should have a
// parent == HANDLE_ROOT.
type FilterAttrs struct ***REMOVED***
	LinkIndex int
	Handle    uint32
	Parent    uint32
	Priority  uint16 // lower is higher priority
	Protocol  uint16 // syscall.ETH_P_*
***REMOVED***

func (q FilterAttrs) String() string ***REMOVED***
	return fmt.Sprintf("***REMOVED***LinkIndex: %d, Handle: %s, Parent: %s, Priority: %d, Protocol: %d***REMOVED***", q.LinkIndex, HandleStr(q.Handle), HandleStr(q.Parent), q.Priority, q.Protocol)
***REMOVED***

type TcAct int32

const (
	TC_ACT_UNSPEC     TcAct = -1
	TC_ACT_OK         TcAct = 0
	TC_ACT_RECLASSIFY TcAct = 1
	TC_ACT_SHOT       TcAct = 2
	TC_ACT_PIPE       TcAct = 3
	TC_ACT_STOLEN     TcAct = 4
	TC_ACT_QUEUED     TcAct = 5
	TC_ACT_REPEAT     TcAct = 6
	TC_ACT_REDIRECT   TcAct = 7
	TC_ACT_JUMP       TcAct = 0x10000000
)

func (a TcAct) String() string ***REMOVED***
	switch a ***REMOVED***
	case TC_ACT_UNSPEC:
		return "unspec"
	case TC_ACT_OK:
		return "ok"
	case TC_ACT_RECLASSIFY:
		return "reclassify"
	case TC_ACT_SHOT:
		return "shot"
	case TC_ACT_PIPE:
		return "pipe"
	case TC_ACT_STOLEN:
		return "stolen"
	case TC_ACT_QUEUED:
		return "queued"
	case TC_ACT_REPEAT:
		return "repeat"
	case TC_ACT_REDIRECT:
		return "redirect"
	case TC_ACT_JUMP:
		return "jump"
	***REMOVED***
	return fmt.Sprintf("0x%x", int32(a))
***REMOVED***

type TcPolAct int32

const (
	TC_POLICE_UNSPEC     TcPolAct = TcPolAct(TC_ACT_UNSPEC)
	TC_POLICE_OK         TcPolAct = TcPolAct(TC_ACT_OK)
	TC_POLICE_RECLASSIFY TcPolAct = TcPolAct(TC_ACT_RECLASSIFY)
	TC_POLICE_SHOT       TcPolAct = TcPolAct(TC_ACT_SHOT)
	TC_POLICE_PIPE       TcPolAct = TcPolAct(TC_ACT_PIPE)
)

func (a TcPolAct) String() string ***REMOVED***
	switch a ***REMOVED***
	case TC_POLICE_UNSPEC:
		return "unspec"
	case TC_POLICE_OK:
		return "ok"
	case TC_POLICE_RECLASSIFY:
		return "reclassify"
	case TC_POLICE_SHOT:
		return "shot"
	case TC_POLICE_PIPE:
		return "pipe"
	***REMOVED***
	return fmt.Sprintf("0x%x", int32(a))
***REMOVED***

type ActionAttrs struct ***REMOVED***
	Index   int
	Capab   int
	Action  TcAct
	Refcnt  int
	Bindcnt int
***REMOVED***

func (q ActionAttrs) String() string ***REMOVED***
	return fmt.Sprintf("***REMOVED***Index: %d, Capab: %x, Action: %s, Refcnt: %d, Bindcnt: %d***REMOVED***", q.Index, q.Capab, q.Action.String(), q.Refcnt, q.Bindcnt)
***REMOVED***

// Action represents an action in any supported filter.
type Action interface ***REMOVED***
	Attrs() *ActionAttrs
	Type() string
***REMOVED***

type GenericAction struct ***REMOVED***
	ActionAttrs
***REMOVED***

func (action *GenericAction) Type() string ***REMOVED***
	return "generic"
***REMOVED***

func (action *GenericAction) Attrs() *ActionAttrs ***REMOVED***
	return &action.ActionAttrs
***REMOVED***

type BpfAction struct ***REMOVED***
	ActionAttrs
	Fd   int
	Name string
***REMOVED***

func (action *BpfAction) Type() string ***REMOVED***
	return "bpf"
***REMOVED***

func (action *BpfAction) Attrs() *ActionAttrs ***REMOVED***
	return &action.ActionAttrs
***REMOVED***

type MirredAct uint8

func (a MirredAct) String() string ***REMOVED***
	switch a ***REMOVED***
	case TCA_EGRESS_REDIR:
		return "egress redir"
	case TCA_EGRESS_MIRROR:
		return "egress mirror"
	case TCA_INGRESS_REDIR:
		return "ingress redir"
	case TCA_INGRESS_MIRROR:
		return "ingress mirror"
	***REMOVED***
	return "unknown"
***REMOVED***

const (
	TCA_EGRESS_REDIR   MirredAct = 1 /* packet redirect to EGRESS*/
	TCA_EGRESS_MIRROR  MirredAct = 2 /* mirror packet to EGRESS */
	TCA_INGRESS_REDIR  MirredAct = 3 /* packet redirect to INGRESS*/
	TCA_INGRESS_MIRROR MirredAct = 4 /* mirror packet to INGRESS */
)

type MirredAction struct ***REMOVED***
	ActionAttrs
	MirredAction MirredAct
	Ifindex      int
***REMOVED***

func (action *MirredAction) Type() string ***REMOVED***
	return "mirred"
***REMOVED***

func (action *MirredAction) Attrs() *ActionAttrs ***REMOVED***
	return &action.ActionAttrs
***REMOVED***

func NewMirredAction(redirIndex int) *MirredAction ***REMOVED***
	return &MirredAction***REMOVED***
		ActionAttrs: ActionAttrs***REMOVED***
			Action: TC_ACT_STOLEN,
		***REMOVED***,
		MirredAction: TCA_EGRESS_REDIR,
		Ifindex:      redirIndex,
	***REMOVED***
***REMOVED***

// Sel of the U32 filters that contains multiple TcU32Key. This is the copy
// and the frontend representation of nl.TcU32Sel. It is serialized into canonical
// nl.TcU32Sel with the appropriate endianness.
type TcU32Sel struct ***REMOVED***
	Flags    uint8
	Offshift uint8
	Nkeys    uint8
	Pad      uint8
	Offmask  uint16
	Off      uint16
	Offoff   int16
	Hoff     int16
	Hmask    uint32
	Keys     []TcU32Key
***REMOVED***

// TcU32Key contained of Sel in the U32 filters. This is the copy and the frontend
// representation of nl.TcU32Key. It is serialized into chanonical nl.TcU32Sel
// with the appropriate endianness.
type TcU32Key struct ***REMOVED***
	Mask    uint32
	Val     uint32
	Off     int32
	OffMask int32
***REMOVED***

// U32 filters on many packet related properties
type U32 struct ***REMOVED***
	FilterAttrs
	ClassId    uint32
	RedirIndex int
	Sel        *TcU32Sel
	Actions    []Action
***REMOVED***

func (filter *U32) Attrs() *FilterAttrs ***REMOVED***
	return &filter.FilterAttrs
***REMOVED***

func (filter *U32) Type() string ***REMOVED***
	return "u32"
***REMOVED***

type FilterFwAttrs struct ***REMOVED***
	ClassId   uint32
	InDev     string
	Mask      uint32
	Index     uint32
	Buffer    uint32
	Mtu       uint32
	Mpu       uint16
	Rate      uint32
	AvRate    uint32
	PeakRate  uint32
	Action    TcPolAct
	Overhead  uint16
	LinkLayer int
***REMOVED***

type BpfFilter struct ***REMOVED***
	FilterAttrs
	ClassId      uint32
	Fd           int
	Name         string
	DirectAction bool
***REMOVED***

func (filter *BpfFilter) Type() string ***REMOVED***
	return "bpf"
***REMOVED***

func (filter *BpfFilter) Attrs() *FilterAttrs ***REMOVED***
	return &filter.FilterAttrs
***REMOVED***

// GenericFilter filters represent types that are not currently understood
// by this netlink library.
type GenericFilter struct ***REMOVED***
	FilterAttrs
	FilterType string
***REMOVED***

func (filter *GenericFilter) Attrs() *FilterAttrs ***REMOVED***
	return &filter.FilterAttrs
***REMOVED***

func (filter *GenericFilter) Type() string ***REMOVED***
	return filter.FilterType
***REMOVED***
