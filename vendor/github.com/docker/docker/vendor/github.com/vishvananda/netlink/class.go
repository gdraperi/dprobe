package netlink

import (
	"fmt"
)

type Class interface ***REMOVED***
	Attrs() *ClassAttrs
	Type() string
***REMOVED***

// ClassAttrs represents a netlink class. A filter is associated with a link,
// has a handle and a parent. The root filter of a device should have a
// parent == HANDLE_ROOT.
type ClassAttrs struct ***REMOVED***
	LinkIndex int
	Handle    uint32
	Parent    uint32
	Leaf      uint32
***REMOVED***

func (q ClassAttrs) String() string ***REMOVED***
	return fmt.Sprintf("***REMOVED***LinkIndex: %d, Handle: %s, Parent: %s, Leaf: %d***REMOVED***", q.LinkIndex, HandleStr(q.Handle), HandleStr(q.Parent), q.Leaf)
***REMOVED***

type HtbClassAttrs struct ***REMOVED***
	// TODO handle all attributes
	Rate    uint64
	Ceil    uint64
	Buffer  uint32
	Cbuffer uint32
	Quantum uint32
	Level   uint32
	Prio    uint32
***REMOVED***

func (q HtbClassAttrs) String() string ***REMOVED***
	return fmt.Sprintf("***REMOVED***Rate: %d, Ceil: %d, Buffer: %d, Cbuffer: %d***REMOVED***", q.Rate, q.Ceil, q.Buffer, q.Cbuffer)
***REMOVED***

// HtbClass represents an Htb class
type HtbClass struct ***REMOVED***
	ClassAttrs
	Rate    uint64
	Ceil    uint64
	Buffer  uint32
	Cbuffer uint32
	Quantum uint32
	Level   uint32
	Prio    uint32
***REMOVED***

func (q HtbClass) String() string ***REMOVED***
	return fmt.Sprintf("***REMOVED***Rate: %d, Ceil: %d, Buffer: %d, Cbuffer: %d***REMOVED***", q.Rate, q.Ceil, q.Buffer, q.Cbuffer)
***REMOVED***

func (q *HtbClass) Attrs() *ClassAttrs ***REMOVED***
	return &q.ClassAttrs
***REMOVED***

func (q *HtbClass) Type() string ***REMOVED***
	return "htb"
***REMOVED***

// GenericClass classes represent types that are not currently understood
// by this netlink library.
type GenericClass struct ***REMOVED***
	ClassAttrs
	ClassType string
***REMOVED***

func (class *GenericClass) Attrs() *ClassAttrs ***REMOVED***
	return &class.ClassAttrs
***REMOVED***

func (class *GenericClass) Type() string ***REMOVED***
	return class.ClassType
***REMOVED***
