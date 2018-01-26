package netlink

import (
	"fmt"
	"math"
)

const (
	HANDLE_NONE      = 0
	HANDLE_INGRESS   = 0xFFFFFFF1
	HANDLE_CLSACT    = HANDLE_INGRESS
	HANDLE_ROOT      = 0xFFFFFFFF
	PRIORITY_MAP_LEN = 16
)
const (
	HANDLE_MIN_INGRESS = 0xFFFFFFF2
	HANDLE_MIN_EGRESS  = 0xFFFFFFF3
)

type Qdisc interface ***REMOVED***
	Attrs() *QdiscAttrs
	Type() string
***REMOVED***

// QdiscAttrs represents a netlink qdisc. A qdisc is associated with a link,
// has a handle, a parent and a refcnt. The root qdisc of a device should
// have parent == HANDLE_ROOT.
type QdiscAttrs struct ***REMOVED***
	LinkIndex int
	Handle    uint32
	Parent    uint32
	Refcnt    uint32 // read only
***REMOVED***

func (q QdiscAttrs) String() string ***REMOVED***
	return fmt.Sprintf("***REMOVED***LinkIndex: %d, Handle: %s, Parent: %s, Refcnt: %d***REMOVED***", q.LinkIndex, HandleStr(q.Handle), HandleStr(q.Parent), q.Refcnt)
***REMOVED***

func MakeHandle(major, minor uint16) uint32 ***REMOVED***
	return (uint32(major) << 16) | uint32(minor)
***REMOVED***

func MajorMinor(handle uint32) (uint16, uint16) ***REMOVED***
	return uint16((handle & 0xFFFF0000) >> 16), uint16(handle & 0x0000FFFFF)
***REMOVED***

func HandleStr(handle uint32) string ***REMOVED***
	switch handle ***REMOVED***
	case HANDLE_NONE:
		return "none"
	case HANDLE_INGRESS:
		return "ingress"
	case HANDLE_ROOT:
		return "root"
	default:
		major, minor := MajorMinor(handle)
		return fmt.Sprintf("%x:%x", major, minor)
	***REMOVED***
***REMOVED***

func Percentage2u32(percentage float32) uint32 ***REMOVED***
	// FIXME this is most likely not the best way to convert from % to uint32
	if percentage == 100 ***REMOVED***
		return math.MaxUint32
	***REMOVED***
	return uint32(math.MaxUint32 * (percentage / 100))
***REMOVED***

// PfifoFast is the default qdisc created by the kernel if one has not
// been defined for the interface
type PfifoFast struct ***REMOVED***
	QdiscAttrs
	Bands       uint8
	PriorityMap [PRIORITY_MAP_LEN]uint8
***REMOVED***

func (qdisc *PfifoFast) Attrs() *QdiscAttrs ***REMOVED***
	return &qdisc.QdiscAttrs
***REMOVED***

func (qdisc *PfifoFast) Type() string ***REMOVED***
	return "pfifo_fast"
***REMOVED***

// Prio is a basic qdisc that works just like PfifoFast
type Prio struct ***REMOVED***
	QdiscAttrs
	Bands       uint8
	PriorityMap [PRIORITY_MAP_LEN]uint8
***REMOVED***

func NewPrio(attrs QdiscAttrs) *Prio ***REMOVED***
	return &Prio***REMOVED***
		QdiscAttrs:  attrs,
		Bands:       3,
		PriorityMap: [PRIORITY_MAP_LEN]uint8***REMOVED***1, 2, 2, 2, 1, 2, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1***REMOVED***,
	***REMOVED***
***REMOVED***

func (qdisc *Prio) Attrs() *QdiscAttrs ***REMOVED***
	return &qdisc.QdiscAttrs
***REMOVED***

func (qdisc *Prio) Type() string ***REMOVED***
	return "prio"
***REMOVED***

// Htb is a classful qdisc that rate limits based on tokens
type Htb struct ***REMOVED***
	QdiscAttrs
	Version      uint32
	Rate2Quantum uint32
	Defcls       uint32
	Debug        uint32
	DirectPkts   uint32
***REMOVED***

func NewHtb(attrs QdiscAttrs) *Htb ***REMOVED***
	return &Htb***REMOVED***
		QdiscAttrs:   attrs,
		Version:      3,
		Defcls:       0,
		Rate2Quantum: 10,
		Debug:        0,
		DirectPkts:   0,
	***REMOVED***
***REMOVED***

func (qdisc *Htb) Attrs() *QdiscAttrs ***REMOVED***
	return &qdisc.QdiscAttrs
***REMOVED***

func (qdisc *Htb) Type() string ***REMOVED***
	return "htb"
***REMOVED***

// Netem is a classless qdisc that rate limits based on tokens

type NetemQdiscAttrs struct ***REMOVED***
	Latency       uint32  // in us
	DelayCorr     float32 // in %
	Limit         uint32
	Loss          float32 // in %
	LossCorr      float32 // in %
	Gap           uint32
	Duplicate     float32 // in %
	DuplicateCorr float32 // in %
	Jitter        uint32  // in us
	ReorderProb   float32 // in %
	ReorderCorr   float32 // in %
	CorruptProb   float32 // in %
	CorruptCorr   float32 // in %
***REMOVED***

func (q NetemQdiscAttrs) String() string ***REMOVED***
	return fmt.Sprintf(
		"***REMOVED***Latency: %d, Limit: %d, Loss: %f, Gap: %d, Duplicate: %f, Jitter: %d***REMOVED***",
		q.Latency, q.Limit, q.Loss, q.Gap, q.Duplicate, q.Jitter,
	)
***REMOVED***

type Netem struct ***REMOVED***
	QdiscAttrs
	Latency       uint32
	DelayCorr     uint32
	Limit         uint32
	Loss          uint32
	LossCorr      uint32
	Gap           uint32
	Duplicate     uint32
	DuplicateCorr uint32
	Jitter        uint32
	ReorderProb   uint32
	ReorderCorr   uint32
	CorruptProb   uint32
	CorruptCorr   uint32
***REMOVED***

func (qdisc *Netem) Attrs() *QdiscAttrs ***REMOVED***
	return &qdisc.QdiscAttrs
***REMOVED***

func (qdisc *Netem) Type() string ***REMOVED***
	return "netem"
***REMOVED***

// Tbf is a classless qdisc that rate limits based on tokens
type Tbf struct ***REMOVED***
	QdiscAttrs
	Rate     uint64
	Limit    uint32
	Buffer   uint32
	Peakrate uint64
	Minburst uint32
	// TODO: handle other settings
***REMOVED***

func (qdisc *Tbf) Attrs() *QdiscAttrs ***REMOVED***
	return &qdisc.QdiscAttrs
***REMOVED***

func (qdisc *Tbf) Type() string ***REMOVED***
	return "tbf"
***REMOVED***

// Ingress is a qdisc for adding ingress filters
type Ingress struct ***REMOVED***
	QdiscAttrs
***REMOVED***

func (qdisc *Ingress) Attrs() *QdiscAttrs ***REMOVED***
	return &qdisc.QdiscAttrs
***REMOVED***

func (qdisc *Ingress) Type() string ***REMOVED***
	return "ingress"
***REMOVED***

// GenericQdisc qdiscs represent types that are not currently understood
// by this netlink library.
type GenericQdisc struct ***REMOVED***
	QdiscAttrs
	QdiscType string
***REMOVED***

func (qdisc *GenericQdisc) Attrs() *QdiscAttrs ***REMOVED***
	return &qdisc.QdiscAttrs
***REMOVED***

func (qdisc *GenericQdisc) Type() string ***REMOVED***
	return qdisc.QdiscType
***REMOVED***
