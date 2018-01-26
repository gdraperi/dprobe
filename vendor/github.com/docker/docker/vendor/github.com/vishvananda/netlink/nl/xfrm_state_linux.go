package nl

import (
	"unsafe"
)

const (
	SizeofXfrmUsersaId       = 0x18
	SizeofXfrmStats          = 0x0c
	SizeofXfrmUsersaInfo     = 0xe0
	SizeofXfrmUserSpiInfo    = 0xe8
	SizeofXfrmAlgo           = 0x44
	SizeofXfrmAlgoAuth       = 0x48
	SizeofXfrmAlgoAEAD       = 0x48
	SizeofXfrmEncapTmpl      = 0x18
	SizeofXfrmUsersaFlush    = 0x8
	SizeofXfrmReplayStateEsn = 0x18
)

const (
	XFRM_STATE_NOECN      = 1
	XFRM_STATE_DECAP_DSCP = 2
	XFRM_STATE_NOPMTUDISC = 4
	XFRM_STATE_WILDRECV   = 8
	XFRM_STATE_ICMP       = 16
	XFRM_STATE_AF_UNSPEC  = 32
	XFRM_STATE_ALIGN4     = 64
	XFRM_STATE_ESN        = 128
)

// struct xfrm_usersa_id ***REMOVED***
//   xfrm_address_t      daddr;
//   __be32        spi;
//   __u16       family;
//   __u8        proto;
// ***REMOVED***;

type XfrmUsersaId struct ***REMOVED***
	Daddr  XfrmAddress
	Spi    uint32 // big endian
	Family uint16
	Proto  uint8
	Pad    byte
***REMOVED***

func (msg *XfrmUsersaId) Len() int ***REMOVED***
	return SizeofXfrmUsersaId
***REMOVED***

func DeserializeXfrmUsersaId(b []byte) *XfrmUsersaId ***REMOVED***
	return (*XfrmUsersaId)(unsafe.Pointer(&b[0:SizeofXfrmUsersaId][0]))
***REMOVED***

func (msg *XfrmUsersaId) Serialize() []byte ***REMOVED***
	return (*(*[SizeofXfrmUsersaId]byte)(unsafe.Pointer(msg)))[:]
***REMOVED***

// struct xfrm_stats ***REMOVED***
//   __u32 replay_window;
//   __u32 replay;
//   __u32 integrity_failed;
// ***REMOVED***;

type XfrmStats struct ***REMOVED***
	ReplayWindow    uint32
	Replay          uint32
	IntegrityFailed uint32
***REMOVED***

func (msg *XfrmStats) Len() int ***REMOVED***
	return SizeofXfrmStats
***REMOVED***

func DeserializeXfrmStats(b []byte) *XfrmStats ***REMOVED***
	return (*XfrmStats)(unsafe.Pointer(&b[0:SizeofXfrmStats][0]))
***REMOVED***

func (msg *XfrmStats) Serialize() []byte ***REMOVED***
	return (*(*[SizeofXfrmStats]byte)(unsafe.Pointer(msg)))[:]
***REMOVED***

// struct xfrm_usersa_info ***REMOVED***
//   struct xfrm_selector    sel;
//   struct xfrm_id      id;
//   xfrm_address_t      saddr;
//   struct xfrm_lifetime_cfg  lft;
//   struct xfrm_lifetime_cur  curlft;
//   struct xfrm_stats   stats;
//   __u32       seq;
//   __u32       reqid;
//   __u16       family;
//   __u8        mode;   /* XFRM_MODE_xxx */
//   __u8        replay_window;
//   __u8        flags;
// #define XFRM_STATE_NOECN  1
// #define XFRM_STATE_DECAP_DSCP 2
// #define XFRM_STATE_NOPMTUDISC 4
// #define XFRM_STATE_WILDRECV 8
// #define XFRM_STATE_ICMP   16
// #define XFRM_STATE_AF_UNSPEC  32
// #define XFRM_STATE_ALIGN4 64
// #define XFRM_STATE_ESN    128
// ***REMOVED***;
//
// #define XFRM_SA_XFLAG_DONT_ENCAP_DSCP 1
//

type XfrmUsersaInfo struct ***REMOVED***
	Sel          XfrmSelector
	Id           XfrmId
	Saddr        XfrmAddress
	Lft          XfrmLifetimeCfg
	Curlft       XfrmLifetimeCur
	Stats        XfrmStats
	Seq          uint32
	Reqid        uint32
	Family       uint16
	Mode         uint8
	ReplayWindow uint8
	Flags        uint8
	Pad          [7]byte
***REMOVED***

func (msg *XfrmUsersaInfo) Len() int ***REMOVED***
	return SizeofXfrmUsersaInfo
***REMOVED***

func DeserializeXfrmUsersaInfo(b []byte) *XfrmUsersaInfo ***REMOVED***
	return (*XfrmUsersaInfo)(unsafe.Pointer(&b[0:SizeofXfrmUsersaInfo][0]))
***REMOVED***

func (msg *XfrmUsersaInfo) Serialize() []byte ***REMOVED***
	return (*(*[SizeofXfrmUsersaInfo]byte)(unsafe.Pointer(msg)))[:]
***REMOVED***

// struct xfrm_userspi_info ***REMOVED***
// 	struct xfrm_usersa_info		info;
// 	__u32				min;
// 	__u32				max;
// ***REMOVED***;

type XfrmUserSpiInfo struct ***REMOVED***
	XfrmUsersaInfo XfrmUsersaInfo
	Min            uint32
	Max            uint32
***REMOVED***

func (msg *XfrmUserSpiInfo) Len() int ***REMOVED***
	return SizeofXfrmUserSpiInfo
***REMOVED***

func DeserializeXfrmUserSpiInfo(b []byte) *XfrmUserSpiInfo ***REMOVED***
	return (*XfrmUserSpiInfo)(unsafe.Pointer(&b[0:SizeofXfrmUserSpiInfo][0]))
***REMOVED***

func (msg *XfrmUserSpiInfo) Serialize() []byte ***REMOVED***
	return (*(*[SizeofXfrmUserSpiInfo]byte)(unsafe.Pointer(msg)))[:]
***REMOVED***

// struct xfrm_algo ***REMOVED***
//   char    alg_name[64];
//   unsigned int  alg_key_len;    /* in bits */
//   char    alg_key[0];
// ***REMOVED***;

type XfrmAlgo struct ***REMOVED***
	AlgName   [64]byte
	AlgKeyLen uint32
	AlgKey    []byte
***REMOVED***

func (msg *XfrmAlgo) Len() int ***REMOVED***
	return SizeofXfrmAlgo + int(msg.AlgKeyLen/8)
***REMOVED***

func DeserializeXfrmAlgo(b []byte) *XfrmAlgo ***REMOVED***
	ret := XfrmAlgo***REMOVED******REMOVED***
	copy(ret.AlgName[:], b[0:64])
	ret.AlgKeyLen = *(*uint32)(unsafe.Pointer(&b[64]))
	ret.AlgKey = b[68:ret.Len()]
	return &ret
***REMOVED***

func (msg *XfrmAlgo) Serialize() []byte ***REMOVED***
	b := make([]byte, msg.Len())
	copy(b[0:64], msg.AlgName[:])
	copy(b[64:68], (*(*[4]byte)(unsafe.Pointer(&msg.AlgKeyLen)))[:])
	copy(b[68:msg.Len()], msg.AlgKey[:])
	return b
***REMOVED***

// struct xfrm_algo_auth ***REMOVED***
//   char    alg_name[64];
//   unsigned int  alg_key_len;    /* in bits */
//   unsigned int  alg_trunc_len;  /* in bits */
//   char    alg_key[0];
// ***REMOVED***;

type XfrmAlgoAuth struct ***REMOVED***
	AlgName     [64]byte
	AlgKeyLen   uint32
	AlgTruncLen uint32
	AlgKey      []byte
***REMOVED***

func (msg *XfrmAlgoAuth) Len() int ***REMOVED***
	return SizeofXfrmAlgoAuth + int(msg.AlgKeyLen/8)
***REMOVED***

func DeserializeXfrmAlgoAuth(b []byte) *XfrmAlgoAuth ***REMOVED***
	ret := XfrmAlgoAuth***REMOVED******REMOVED***
	copy(ret.AlgName[:], b[0:64])
	ret.AlgKeyLen = *(*uint32)(unsafe.Pointer(&b[64]))
	ret.AlgTruncLen = *(*uint32)(unsafe.Pointer(&b[68]))
	ret.AlgKey = b[72:ret.Len()]
	return &ret
***REMOVED***

func (msg *XfrmAlgoAuth) Serialize() []byte ***REMOVED***
	b := make([]byte, msg.Len())
	copy(b[0:64], msg.AlgName[:])
	copy(b[64:68], (*(*[4]byte)(unsafe.Pointer(&msg.AlgKeyLen)))[:])
	copy(b[68:72], (*(*[4]byte)(unsafe.Pointer(&msg.AlgTruncLen)))[:])
	copy(b[72:msg.Len()], msg.AlgKey[:])
	return b
***REMOVED***

// struct xfrm_algo_aead ***REMOVED***
//   char    alg_name[64];
//   unsigned int  alg_key_len;  /* in bits */
//   unsigned int  alg_icv_len;  /* in bits */
//   char    alg_key[0];
// ***REMOVED***

type XfrmAlgoAEAD struct ***REMOVED***
	AlgName   [64]byte
	AlgKeyLen uint32
	AlgICVLen uint32
	AlgKey    []byte
***REMOVED***

func (msg *XfrmAlgoAEAD) Len() int ***REMOVED***
	return SizeofXfrmAlgoAEAD + int(msg.AlgKeyLen/8)
***REMOVED***

func DeserializeXfrmAlgoAEAD(b []byte) *XfrmAlgoAEAD ***REMOVED***
	ret := XfrmAlgoAEAD***REMOVED******REMOVED***
	copy(ret.AlgName[:], b[0:64])
	ret.AlgKeyLen = *(*uint32)(unsafe.Pointer(&b[64]))
	ret.AlgICVLen = *(*uint32)(unsafe.Pointer(&b[68]))
	ret.AlgKey = b[72:ret.Len()]
	return &ret
***REMOVED***

func (msg *XfrmAlgoAEAD) Serialize() []byte ***REMOVED***
	b := make([]byte, msg.Len())
	copy(b[0:64], msg.AlgName[:])
	copy(b[64:68], (*(*[4]byte)(unsafe.Pointer(&msg.AlgKeyLen)))[:])
	copy(b[68:72], (*(*[4]byte)(unsafe.Pointer(&msg.AlgICVLen)))[:])
	copy(b[72:msg.Len()], msg.AlgKey[:])
	return b
***REMOVED***

// struct xfrm_encap_tmpl ***REMOVED***
//   __u16   encap_type;
//   __be16    encap_sport;
//   __be16    encap_dport;
//   xfrm_address_t  encap_oa;
// ***REMOVED***;

type XfrmEncapTmpl struct ***REMOVED***
	EncapType  uint16
	EncapSport uint16 // big endian
	EncapDport uint16 // big endian
	Pad        [2]byte
	EncapOa    XfrmAddress
***REMOVED***

func (msg *XfrmEncapTmpl) Len() int ***REMOVED***
	return SizeofXfrmEncapTmpl
***REMOVED***

func DeserializeXfrmEncapTmpl(b []byte) *XfrmEncapTmpl ***REMOVED***
	return (*XfrmEncapTmpl)(unsafe.Pointer(&b[0:SizeofXfrmEncapTmpl][0]))
***REMOVED***

func (msg *XfrmEncapTmpl) Serialize() []byte ***REMOVED***
	return (*(*[SizeofXfrmEncapTmpl]byte)(unsafe.Pointer(msg)))[:]
***REMOVED***

// struct xfrm_usersa_flush ***REMOVED***
//    __u8 proto;
// ***REMOVED***;

type XfrmUsersaFlush struct ***REMOVED***
	Proto uint8
***REMOVED***

func (msg *XfrmUsersaFlush) Len() int ***REMOVED***
	return SizeofXfrmUsersaFlush
***REMOVED***

func DeserializeXfrmUsersaFlush(b []byte) *XfrmUsersaFlush ***REMOVED***
	return (*XfrmUsersaFlush)(unsafe.Pointer(&b[0:SizeofXfrmUsersaFlush][0]))
***REMOVED***

func (msg *XfrmUsersaFlush) Serialize() []byte ***REMOVED***
	return (*(*[SizeofXfrmUsersaFlush]byte)(unsafe.Pointer(msg)))[:]
***REMOVED***

// struct xfrm_replay_state_esn ***REMOVED***
//     unsigned int    bmp_len;
//     __u32           oseq;
//     __u32           seq;
//     __u32           oseq_hi;
//     __u32           seq_hi;
//     __u32           replay_window;
//     __u32           bmp[0];
// ***REMOVED***;

type XfrmReplayStateEsn struct ***REMOVED***
	BmpLen       uint32
	OSeq         uint32
	Seq          uint32
	OSeqHi       uint32
	SeqHi        uint32
	ReplayWindow uint32
	Bmp          []uint32
***REMOVED***

func (msg *XfrmReplayStateEsn) Serialize() []byte ***REMOVED***
	// We deliberately do not pass Bmp, as it gets set by the kernel.
	return (*(*[SizeofXfrmReplayStateEsn]byte)(unsafe.Pointer(msg)))[:]
***REMOVED***
