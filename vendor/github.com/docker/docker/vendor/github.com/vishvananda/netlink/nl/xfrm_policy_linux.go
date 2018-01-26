package nl

import (
	"unsafe"
)

const (
	SizeofXfrmUserpolicyId   = 0x40
	SizeofXfrmUserpolicyInfo = 0xa8
	SizeofXfrmUserTmpl       = 0x40
)

// struct xfrm_userpolicy_id ***REMOVED***
//   struct xfrm_selector    sel;
//   __u32       index;
//   __u8        dir;
// ***REMOVED***;
//

type XfrmUserpolicyId struct ***REMOVED***
	Sel   XfrmSelector
	Index uint32
	Dir   uint8
	Pad   [3]byte
***REMOVED***

func (msg *XfrmUserpolicyId) Len() int ***REMOVED***
	return SizeofXfrmUserpolicyId
***REMOVED***

func DeserializeXfrmUserpolicyId(b []byte) *XfrmUserpolicyId ***REMOVED***
	return (*XfrmUserpolicyId)(unsafe.Pointer(&b[0:SizeofXfrmUserpolicyId][0]))
***REMOVED***

func (msg *XfrmUserpolicyId) Serialize() []byte ***REMOVED***
	return (*(*[SizeofXfrmUserpolicyId]byte)(unsafe.Pointer(msg)))[:]
***REMOVED***

// struct xfrm_userpolicy_info ***REMOVED***
//   struct xfrm_selector    sel;
//   struct xfrm_lifetime_cfg  lft;
//   struct xfrm_lifetime_cur  curlft;
//   __u32       priority;
//   __u32       index;
//   __u8        dir;
//   __u8        action;
// #define XFRM_POLICY_ALLOW 0
// #define XFRM_POLICY_BLOCK 1
//   __u8        flags;
// #define XFRM_POLICY_LOCALOK 1 /* Allow user to override global policy */
//   /* Automatically expand selector to include matching ICMP payloads. */
// #define XFRM_POLICY_ICMP  2
//   __u8        share;
// ***REMOVED***;

type XfrmUserpolicyInfo struct ***REMOVED***
	Sel      XfrmSelector
	Lft      XfrmLifetimeCfg
	Curlft   XfrmLifetimeCur
	Priority uint32
	Index    uint32
	Dir      uint8
	Action   uint8
	Flags    uint8
	Share    uint8
	Pad      [4]byte
***REMOVED***

func (msg *XfrmUserpolicyInfo) Len() int ***REMOVED***
	return SizeofXfrmUserpolicyInfo
***REMOVED***

func DeserializeXfrmUserpolicyInfo(b []byte) *XfrmUserpolicyInfo ***REMOVED***
	return (*XfrmUserpolicyInfo)(unsafe.Pointer(&b[0:SizeofXfrmUserpolicyInfo][0]))
***REMOVED***

func (msg *XfrmUserpolicyInfo) Serialize() []byte ***REMOVED***
	return (*(*[SizeofXfrmUserpolicyInfo]byte)(unsafe.Pointer(msg)))[:]
***REMOVED***

// struct xfrm_user_tmpl ***REMOVED***
//   struct xfrm_id    id;
//   __u16     family;
//   xfrm_address_t    saddr;
//   __u32     reqid;
//   __u8      mode;
//   __u8      share;
//   __u8      optional;
//   __u32     aalgos;
//   __u32     ealgos;
//   __u32     calgos;
// ***REMOVED***

type XfrmUserTmpl struct ***REMOVED***
	XfrmId   XfrmId
	Family   uint16
	Pad1     [2]byte
	Saddr    XfrmAddress
	Reqid    uint32
	Mode     uint8
	Share    uint8
	Optional uint8
	Pad2     byte
	Aalgos   uint32
	Ealgos   uint32
	Calgos   uint32
***REMOVED***

func (msg *XfrmUserTmpl) Len() int ***REMOVED***
	return SizeofXfrmUserTmpl
***REMOVED***

func DeserializeXfrmUserTmpl(b []byte) *XfrmUserTmpl ***REMOVED***
	return (*XfrmUserTmpl)(unsafe.Pointer(&b[0:SizeofXfrmUserTmpl][0]))
***REMOVED***

func (msg *XfrmUserTmpl) Serialize() []byte ***REMOVED***
	return (*(*[SizeofXfrmUserTmpl]byte)(unsafe.Pointer(msg)))[:]
***REMOVED***
