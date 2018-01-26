package nl

import (
	"unsafe"
)

// LinkLayer
const (
	LINKLAYER_UNSPEC = iota
	LINKLAYER_ETHERNET
	LINKLAYER_ATM
)

// ATM
const (
	ATM_CELL_PAYLOAD = 48
	ATM_CELL_SIZE    = 53
)

const TC_LINKLAYER_MASK = 0x0F

// Police
const (
	TCA_POLICE_UNSPEC = iota
	TCA_POLICE_TBF
	TCA_POLICE_RATE
	TCA_POLICE_PEAKRATE
	TCA_POLICE_AVRATE
	TCA_POLICE_RESULT
	TCA_POLICE_MAX = TCA_POLICE_RESULT
)

// Message types
const (
	TCA_UNSPEC = iota
	TCA_KIND
	TCA_OPTIONS
	TCA_STATS
	TCA_XSTATS
	TCA_RATE
	TCA_FCNT
	TCA_STATS2
	TCA_STAB
	TCA_MAX = TCA_STAB
)

const (
	TCA_ACT_TAB = 1
	TCAA_MAX    = 1
)

const (
	TCA_ACT_UNSPEC = iota
	TCA_ACT_KIND
	TCA_ACT_OPTIONS
	TCA_ACT_INDEX
	TCA_ACT_STATS
	TCA_ACT_MAX
)

const (
	TCA_PRIO_UNSPEC = iota
	TCA_PRIO_MQ
	TCA_PRIO_MAX = TCA_PRIO_MQ
)

const (
	SizeofTcMsg          = 0x14
	SizeofTcActionMsg    = 0x04
	SizeofTcPrioMap      = 0x14
	SizeofTcRateSpec     = 0x0c
	SizeofTcNetemQopt    = 0x18
	SizeofTcNetemCorr    = 0x0c
	SizeofTcNetemReorder = 0x08
	SizeofTcNetemCorrupt = 0x08
	SizeofTcTbfQopt      = 2*SizeofTcRateSpec + 0x0c
	SizeofTcHtbCopt      = 2*SizeofTcRateSpec + 0x14
	SizeofTcHtbGlob      = 0x14
	SizeofTcU32Key       = 0x10
	SizeofTcU32Sel       = 0x10 // without keys
	SizeofTcGen          = 0x14
	SizeofTcMirred       = SizeofTcGen + 0x08
	SizeofTcPolice       = 2*SizeofTcRateSpec + 0x20
)

// struct tcmsg ***REMOVED***
//   unsigned char tcm_family;
//   unsigned char tcm__pad1;
//   unsigned short  tcm__pad2;
//   int   tcm_ifindex;
//   __u32   tcm_handle;
//   __u32   tcm_parent;
//   __u32   tcm_info;
// ***REMOVED***;

type TcMsg struct ***REMOVED***
	Family  uint8
	Pad     [3]byte
	Ifindex int32
	Handle  uint32
	Parent  uint32
	Info    uint32
***REMOVED***

func (msg *TcMsg) Len() int ***REMOVED***
	return SizeofTcMsg
***REMOVED***

func DeserializeTcMsg(b []byte) *TcMsg ***REMOVED***
	return (*TcMsg)(unsafe.Pointer(&b[0:SizeofTcMsg][0]))
***REMOVED***

func (x *TcMsg) Serialize() []byte ***REMOVED***
	return (*(*[SizeofTcMsg]byte)(unsafe.Pointer(x)))[:]
***REMOVED***

// struct tcamsg ***REMOVED***
//   unsigned char tca_family;
//   unsigned char tca__pad1;
//   unsigned short  tca__pad2;
// ***REMOVED***;

type TcActionMsg struct ***REMOVED***
	Family uint8
	Pad    [3]byte
***REMOVED***

func (msg *TcActionMsg) Len() int ***REMOVED***
	return SizeofTcActionMsg
***REMOVED***

func DeserializeTcActionMsg(b []byte) *TcActionMsg ***REMOVED***
	return (*TcActionMsg)(unsafe.Pointer(&b[0:SizeofTcActionMsg][0]))
***REMOVED***

func (x *TcActionMsg) Serialize() []byte ***REMOVED***
	return (*(*[SizeofTcActionMsg]byte)(unsafe.Pointer(x)))[:]
***REMOVED***

const (
	TC_PRIO_MAX = 15
)

// struct tc_prio_qopt ***REMOVED***
// 	int bands;      /* Number of bands */
// 	__u8  priomap[TC_PRIO_MAX+1]; /* Map: logical priority -> PRIO band */
// ***REMOVED***;

type TcPrioMap struct ***REMOVED***
	Bands   int32
	Priomap [TC_PRIO_MAX + 1]uint8
***REMOVED***

func (msg *TcPrioMap) Len() int ***REMOVED***
	return SizeofTcPrioMap
***REMOVED***

func DeserializeTcPrioMap(b []byte) *TcPrioMap ***REMOVED***
	return (*TcPrioMap)(unsafe.Pointer(&b[0:SizeofTcPrioMap][0]))
***REMOVED***

func (x *TcPrioMap) Serialize() []byte ***REMOVED***
	return (*(*[SizeofTcPrioMap]byte)(unsafe.Pointer(x)))[:]
***REMOVED***

const (
	TCA_TBF_UNSPEC = iota
	TCA_TBF_PARMS
	TCA_TBF_RTAB
	TCA_TBF_PTAB
	TCA_TBF_RATE64
	TCA_TBF_PRATE64
	TCA_TBF_BURST
	TCA_TBF_PBURST
	TCA_TBF_MAX = TCA_TBF_PBURST
)

// struct tc_ratespec ***REMOVED***
//   unsigned char cell_log;
//   __u8    linklayer; /* lower 4 bits */
//   unsigned short  overhead;
//   short   cell_align;
//   unsigned short  mpu;
//   __u32   rate;
// ***REMOVED***;

type TcRateSpec struct ***REMOVED***
	CellLog   uint8
	Linklayer uint8
	Overhead  uint16
	CellAlign int16
	Mpu       uint16
	Rate      uint32
***REMOVED***

func (msg *TcRateSpec) Len() int ***REMOVED***
	return SizeofTcRateSpec
***REMOVED***

func DeserializeTcRateSpec(b []byte) *TcRateSpec ***REMOVED***
	return (*TcRateSpec)(unsafe.Pointer(&b[0:SizeofTcRateSpec][0]))
***REMOVED***

func (x *TcRateSpec) Serialize() []byte ***REMOVED***
	return (*(*[SizeofTcRateSpec]byte)(unsafe.Pointer(x)))[:]
***REMOVED***

/**
* NETEM
 */

const (
	TCA_NETEM_UNSPEC = iota
	TCA_NETEM_CORR
	TCA_NETEM_DELAY_DIST
	TCA_NETEM_REORDER
	TCA_NETEM_CORRUPT
	TCA_NETEM_LOSS
	TCA_NETEM_RATE
	TCA_NETEM_ECN
	TCA_NETEM_RATE64
	TCA_NETEM_MAX = TCA_NETEM_RATE64
)

// struct tc_netem_qopt ***REMOVED***
//	__u32	latency;	/* added delay (us) */
//	__u32   limit;		/* fifo limit (packets) */
//	__u32	loss;		/* random packet loss (0=none ~0=100%) */
//	__u32	gap;		/* re-ordering gap (0 for none) */
//	__u32   duplicate;	/* random packet dup  (0=none ~0=100%) */
// 	__u32	jitter;		/* random jitter in latency (us) */
// ***REMOVED***;

type TcNetemQopt struct ***REMOVED***
	Latency   uint32
	Limit     uint32
	Loss      uint32
	Gap       uint32
	Duplicate uint32
	Jitter    uint32
***REMOVED***

func (msg *TcNetemQopt) Len() int ***REMOVED***
	return SizeofTcNetemQopt
***REMOVED***

func DeserializeTcNetemQopt(b []byte) *TcNetemQopt ***REMOVED***
	return (*TcNetemQopt)(unsafe.Pointer(&b[0:SizeofTcNetemQopt][0]))
***REMOVED***

func (x *TcNetemQopt) Serialize() []byte ***REMOVED***
	return (*(*[SizeofTcNetemQopt]byte)(unsafe.Pointer(x)))[:]
***REMOVED***

// struct tc_netem_corr ***REMOVED***
//  __u32   delay_corr; /* delay correlation */
//  __u32   loss_corr;  /* packet loss correlation */
//  __u32   dup_corr;   /* duplicate correlation  */
// ***REMOVED***;

type TcNetemCorr struct ***REMOVED***
	DelayCorr uint32
	LossCorr  uint32
	DupCorr   uint32
***REMOVED***

func (msg *TcNetemCorr) Len() int ***REMOVED***
	return SizeofTcNetemCorr
***REMOVED***

func DeserializeTcNetemCorr(b []byte) *TcNetemCorr ***REMOVED***
	return (*TcNetemCorr)(unsafe.Pointer(&b[0:SizeofTcNetemCorr][0]))
***REMOVED***

func (x *TcNetemCorr) Serialize() []byte ***REMOVED***
	return (*(*[SizeofTcNetemCorr]byte)(unsafe.Pointer(x)))[:]
***REMOVED***

// struct tc_netem_reorder ***REMOVED***
//  __u32   probability;
//  __u32   correlation;
// ***REMOVED***;

type TcNetemReorder struct ***REMOVED***
	Probability uint32
	Correlation uint32
***REMOVED***

func (msg *TcNetemReorder) Len() int ***REMOVED***
	return SizeofTcNetemReorder
***REMOVED***

func DeserializeTcNetemReorder(b []byte) *TcNetemReorder ***REMOVED***
	return (*TcNetemReorder)(unsafe.Pointer(&b[0:SizeofTcNetemReorder][0]))
***REMOVED***

func (x *TcNetemReorder) Serialize() []byte ***REMOVED***
	return (*(*[SizeofTcNetemReorder]byte)(unsafe.Pointer(x)))[:]
***REMOVED***

// struct tc_netem_corrupt ***REMOVED***
//  __u32   probability;
//  __u32   correlation;
// ***REMOVED***;

type TcNetemCorrupt struct ***REMOVED***
	Probability uint32
	Correlation uint32
***REMOVED***

func (msg *TcNetemCorrupt) Len() int ***REMOVED***
	return SizeofTcNetemCorrupt
***REMOVED***

func DeserializeTcNetemCorrupt(b []byte) *TcNetemCorrupt ***REMOVED***
	return (*TcNetemCorrupt)(unsafe.Pointer(&b[0:SizeofTcNetemCorrupt][0]))
***REMOVED***

func (x *TcNetemCorrupt) Serialize() []byte ***REMOVED***
	return (*(*[SizeofTcNetemCorrupt]byte)(unsafe.Pointer(x)))[:]
***REMOVED***

// struct tc_tbf_qopt ***REMOVED***
//   struct tc_ratespec rate;
//   struct tc_ratespec peakrate;
//   __u32   limit;
//   __u32   buffer;
//   __u32   mtu;
// ***REMOVED***;

type TcTbfQopt struct ***REMOVED***
	Rate     TcRateSpec
	Peakrate TcRateSpec
	Limit    uint32
	Buffer   uint32
	Mtu      uint32
***REMOVED***

func (msg *TcTbfQopt) Len() int ***REMOVED***
	return SizeofTcTbfQopt
***REMOVED***

func DeserializeTcTbfQopt(b []byte) *TcTbfQopt ***REMOVED***
	return (*TcTbfQopt)(unsafe.Pointer(&b[0:SizeofTcTbfQopt][0]))
***REMOVED***

func (x *TcTbfQopt) Serialize() []byte ***REMOVED***
	return (*(*[SizeofTcTbfQopt]byte)(unsafe.Pointer(x)))[:]
***REMOVED***

const (
	TCA_HTB_UNSPEC = iota
	TCA_HTB_PARMS
	TCA_HTB_INIT
	TCA_HTB_CTAB
	TCA_HTB_RTAB
	TCA_HTB_DIRECT_QLEN
	TCA_HTB_RATE64
	TCA_HTB_CEIL64
	TCA_HTB_MAX = TCA_HTB_CEIL64
)

//struct tc_htb_opt ***REMOVED***
//	struct tc_ratespec	rate;
//	struct tc_ratespec	ceil;
//	__u32	buffer;
//	__u32	cbuffer;
//	__u32	quantum;
//	__u32	level;		/* out only */
//	__u32	prio;
//***REMOVED***;

type TcHtbCopt struct ***REMOVED***
	Rate    TcRateSpec
	Ceil    TcRateSpec
	Buffer  uint32
	Cbuffer uint32
	Quantum uint32
	Level   uint32
	Prio    uint32
***REMOVED***

func (msg *TcHtbCopt) Len() int ***REMOVED***
	return SizeofTcHtbCopt
***REMOVED***

func DeserializeTcHtbCopt(b []byte) *TcHtbCopt ***REMOVED***
	return (*TcHtbCopt)(unsafe.Pointer(&b[0:SizeofTcHtbCopt][0]))
***REMOVED***

func (x *TcHtbCopt) Serialize() []byte ***REMOVED***
	return (*(*[SizeofTcHtbCopt]byte)(unsafe.Pointer(x)))[:]
***REMOVED***

type TcHtbGlob struct ***REMOVED***
	Version      uint32
	Rate2Quantum uint32
	Defcls       uint32
	Debug        uint32
	DirectPkts   uint32
***REMOVED***

func (msg *TcHtbGlob) Len() int ***REMOVED***
	return SizeofTcHtbGlob
***REMOVED***

func DeserializeTcHtbGlob(b []byte) *TcHtbGlob ***REMOVED***
	return (*TcHtbGlob)(unsafe.Pointer(&b[0:SizeofTcHtbGlob][0]))
***REMOVED***

func (x *TcHtbGlob) Serialize() []byte ***REMOVED***
	return (*(*[SizeofTcHtbGlob]byte)(unsafe.Pointer(x)))[:]
***REMOVED***

const (
	TCA_U32_UNSPEC = iota
	TCA_U32_CLASSID
	TCA_U32_HASH
	TCA_U32_LINK
	TCA_U32_DIVISOR
	TCA_U32_SEL
	TCA_U32_POLICE
	TCA_U32_ACT
	TCA_U32_INDEV
	TCA_U32_PCNT
	TCA_U32_MARK
	TCA_U32_MAX = TCA_U32_MARK
)

// struct tc_u32_key ***REMOVED***
//   __be32    mask;
//   __be32    val;
//   int   off;
//   int   offmask;
// ***REMOVED***;

type TcU32Key struct ***REMOVED***
	Mask    uint32 // big endian
	Val     uint32 // big endian
	Off     int32
	OffMask int32
***REMOVED***

func (msg *TcU32Key) Len() int ***REMOVED***
	return SizeofTcU32Key
***REMOVED***

func DeserializeTcU32Key(b []byte) *TcU32Key ***REMOVED***
	return (*TcU32Key)(unsafe.Pointer(&b[0:SizeofTcU32Key][0]))
***REMOVED***

func (x *TcU32Key) Serialize() []byte ***REMOVED***
	return (*(*[SizeofTcU32Key]byte)(unsafe.Pointer(x)))[:]
***REMOVED***

// struct tc_u32_sel ***REMOVED***
//   unsigned char   flags;
//   unsigned char   offshift;
//   unsigned char   nkeys;
//
//   __be16      offmask;
//   __u16     off;
//   short     offoff;
//
//   short     hoff;
//   __be32      hmask;
//   struct tc_u32_key keys[0];
// ***REMOVED***;

const (
	TC_U32_TERMINAL  = 1 << iota
	TC_U32_OFFSET    = 1 << iota
	TC_U32_VAROFFSET = 1 << iota
	TC_U32_EAT       = 1 << iota
)

type TcU32Sel struct ***REMOVED***
	Flags    uint8
	Offshift uint8
	Nkeys    uint8
	Pad      uint8
	Offmask  uint16 // big endian
	Off      uint16
	Offoff   int16
	Hoff     int16
	Hmask    uint32 // big endian
	Keys     []TcU32Key
***REMOVED***

func (msg *TcU32Sel) Len() int ***REMOVED***
	return SizeofTcU32Sel + int(msg.Nkeys)*SizeofTcU32Key
***REMOVED***

func DeserializeTcU32Sel(b []byte) *TcU32Sel ***REMOVED***
	x := &TcU32Sel***REMOVED******REMOVED***
	copy((*(*[SizeofTcU32Sel]byte)(unsafe.Pointer(x)))[:], b)
	next := SizeofTcU32Sel
	var i uint8
	for i = 0; i < x.Nkeys; i++ ***REMOVED***
		x.Keys = append(x.Keys, *DeserializeTcU32Key(b[next:]))
		next += SizeofTcU32Key
	***REMOVED***
	return x
***REMOVED***

func (x *TcU32Sel) Serialize() []byte ***REMOVED***
	// This can't just unsafe.cast because it must iterate through keys.
	buf := make([]byte, x.Len())
	copy(buf, (*(*[SizeofTcU32Sel]byte)(unsafe.Pointer(x)))[:])
	next := SizeofTcU32Sel
	for _, key := range x.Keys ***REMOVED***
		keyBuf := key.Serialize()
		copy(buf[next:], keyBuf)
		next += SizeofTcU32Key
	***REMOVED***
	return buf
***REMOVED***

type TcGen struct ***REMOVED***
	Index   uint32
	Capab   uint32
	Action  int32
	Refcnt  int32
	Bindcnt int32
***REMOVED***

func (msg *TcGen) Len() int ***REMOVED***
	return SizeofTcGen
***REMOVED***

func DeserializeTcGen(b []byte) *TcGen ***REMOVED***
	return (*TcGen)(unsafe.Pointer(&b[0:SizeofTcGen][0]))
***REMOVED***

func (x *TcGen) Serialize() []byte ***REMOVED***
	return (*(*[SizeofTcGen]byte)(unsafe.Pointer(x)))[:]
***REMOVED***

// #define tc_gen \
//   __u32                 index; \
//   __u32                 capab; \
//   int                   action; \
//   int                   refcnt; \
//   int                   bindcnt

const (
	TCA_ACT_GACT = 5
)

const (
	TCA_GACT_UNSPEC = iota
	TCA_GACT_TM
	TCA_GACT_PARMS
	TCA_GACT_PROB
	TCA_GACT_MAX = TCA_GACT_PROB
)

type TcGact TcGen

const (
	TCA_ACT_BPF = 13
)

const (
	TCA_ACT_BPF_UNSPEC = iota
	TCA_ACT_BPF_TM
	TCA_ACT_BPF_PARMS
	TCA_ACT_BPF_OPS_LEN
	TCA_ACT_BPF_OPS
	TCA_ACT_BPF_FD
	TCA_ACT_BPF_NAME
	TCA_ACT_BPF_MAX = TCA_ACT_BPF_NAME
)

const (
	TCA_BPF_FLAG_ACT_DIRECT uint32 = 1 << iota
)

const (
	TCA_BPF_UNSPEC = iota
	TCA_BPF_ACT
	TCA_BPF_POLICE
	TCA_BPF_CLASSID
	TCA_BPF_OPS_LEN
	TCA_BPF_OPS
	TCA_BPF_FD
	TCA_BPF_NAME
	TCA_BPF_FLAGS
	TCA_BPF_MAX = TCA_BPF_FLAGS
)

type TcBpf TcGen

const (
	TCA_ACT_MIRRED = 8
)

const (
	TCA_MIRRED_UNSPEC = iota
	TCA_MIRRED_TM
	TCA_MIRRED_PARMS
	TCA_MIRRED_MAX = TCA_MIRRED_PARMS
)

// struct tc_mirred ***REMOVED***
// 	tc_gen;
// 	int                     eaction;   /* one of IN/EGRESS_MIRROR/REDIR */
// 	__u32                   ifindex;  /* ifindex of egress port */
// ***REMOVED***;

type TcMirred struct ***REMOVED***
	TcGen
	Eaction int32
	Ifindex uint32
***REMOVED***

func (msg *TcMirred) Len() int ***REMOVED***
	return SizeofTcMirred
***REMOVED***

func DeserializeTcMirred(b []byte) *TcMirred ***REMOVED***
	return (*TcMirred)(unsafe.Pointer(&b[0:SizeofTcMirred][0]))
***REMOVED***

func (x *TcMirred) Serialize() []byte ***REMOVED***
	return (*(*[SizeofTcMirred]byte)(unsafe.Pointer(x)))[:]
***REMOVED***

// struct tc_police ***REMOVED***
// 	__u32			index;
// 	int			action;
// 	__u32			limit;
// 	__u32			burst;
// 	__u32			mtu;
// 	struct tc_ratespec	rate;
// 	struct tc_ratespec	peakrate;
// 	int				refcnt;
// 	int				bindcnt;
// 	__u32			capab;
// ***REMOVED***;

type TcPolice struct ***REMOVED***
	Index    uint32
	Action   int32
	Limit    uint32
	Burst    uint32
	Mtu      uint32
	Rate     TcRateSpec
	PeakRate TcRateSpec
	Refcnt   int32
	Bindcnt  int32
	Capab    uint32
***REMOVED***

func (msg *TcPolice) Len() int ***REMOVED***
	return SizeofTcPolice
***REMOVED***

func DeserializeTcPolice(b []byte) *TcPolice ***REMOVED***
	return (*TcPolice)(unsafe.Pointer(&b[0:SizeofTcPolice][0]))
***REMOVED***

func (x *TcPolice) Serialize() []byte ***REMOVED***
	return (*(*[SizeofTcPolice]byte)(unsafe.Pointer(x)))[:]
***REMOVED***

const (
	TCA_FW_UNSPEC = iota
	TCA_FW_CLASSID
	TCA_FW_POLICE
	TCA_FW_INDEV
	TCA_FW_ACT
	TCA_FW_MASK
	TCA_FW_MAX = TCA_FW_MASK
)
