package netlink

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/vishvananda/netlink/nl"
)

func writeStateAlgo(a *XfrmStateAlgo) []byte ***REMOVED***
	algo := nl.XfrmAlgo***REMOVED***
		AlgKeyLen: uint32(len(a.Key) * 8),
		AlgKey:    a.Key,
	***REMOVED***
	end := len(a.Name)
	if end > 64 ***REMOVED***
		end = 64
	***REMOVED***
	copy(algo.AlgName[:end], a.Name)
	return algo.Serialize()
***REMOVED***

func writeStateAlgoAuth(a *XfrmStateAlgo) []byte ***REMOVED***
	algo := nl.XfrmAlgoAuth***REMOVED***
		AlgKeyLen:   uint32(len(a.Key) * 8),
		AlgTruncLen: uint32(a.TruncateLen),
		AlgKey:      a.Key,
	***REMOVED***
	end := len(a.Name)
	if end > 64 ***REMOVED***
		end = 64
	***REMOVED***
	copy(algo.AlgName[:end], a.Name)
	return algo.Serialize()
***REMOVED***

func writeStateAlgoAead(a *XfrmStateAlgo) []byte ***REMOVED***
	algo := nl.XfrmAlgoAEAD***REMOVED***
		AlgKeyLen: uint32(len(a.Key) * 8),
		AlgICVLen: uint32(a.ICVLen),
		AlgKey:    a.Key,
	***REMOVED***
	end := len(a.Name)
	if end > 64 ***REMOVED***
		end = 64
	***REMOVED***
	copy(algo.AlgName[:end], a.Name)
	return algo.Serialize()
***REMOVED***

func writeMark(m *XfrmMark) []byte ***REMOVED***
	mark := &nl.XfrmMark***REMOVED***
		Value: m.Value,
		Mask:  m.Mask,
	***REMOVED***
	if mark.Mask == 0 ***REMOVED***
		mark.Mask = ^uint32(0)
	***REMOVED***
	return mark.Serialize()
***REMOVED***

func writeReplayEsn(replayWindow int) []byte ***REMOVED***
	replayEsn := &nl.XfrmReplayStateEsn***REMOVED***
		OSeq:         0,
		Seq:          0,
		OSeqHi:       0,
		SeqHi:        0,
		ReplayWindow: uint32(replayWindow),
	***REMOVED***

	// taken from iproute2/ip/xfrm_state.c:
	replayEsn.BmpLen = uint32((replayWindow + (4 * 8) - 1) / (4 * 8))

	return replayEsn.Serialize()
***REMOVED***

// XfrmStateAdd will add an xfrm state to the system.
// Equivalent to: `ip xfrm state add $state`
func XfrmStateAdd(state *XfrmState) error ***REMOVED***
	return pkgHandle.XfrmStateAdd(state)
***REMOVED***

// XfrmStateAdd will add an xfrm state to the system.
// Equivalent to: `ip xfrm state add $state`
func (h *Handle) XfrmStateAdd(state *XfrmState) error ***REMOVED***
	return h.xfrmStateAddOrUpdate(state, nl.XFRM_MSG_NEWSA)
***REMOVED***

// XfrmStateAllocSpi will allocate an xfrm state in the system.
// Equivalent to: `ip xfrm state allocspi`
func XfrmStateAllocSpi(state *XfrmState) (*XfrmState, error) ***REMOVED***
	return pkgHandle.xfrmStateAllocSpi(state)
***REMOVED***

// XfrmStateUpdate will update an xfrm state to the system.
// Equivalent to: `ip xfrm state update $state`
func XfrmStateUpdate(state *XfrmState) error ***REMOVED***
	return pkgHandle.XfrmStateUpdate(state)
***REMOVED***

// XfrmStateUpdate will update an xfrm state to the system.
// Equivalent to: `ip xfrm state update $state`
func (h *Handle) XfrmStateUpdate(state *XfrmState) error ***REMOVED***
	return h.xfrmStateAddOrUpdate(state, nl.XFRM_MSG_UPDSA)
***REMOVED***

func (h *Handle) xfrmStateAddOrUpdate(state *XfrmState, nlProto int) error ***REMOVED***

	// A state with spi 0 can't be deleted so don't allow it to be set
	if state.Spi == 0 ***REMOVED***
		return fmt.Errorf("Spi must be set when adding xfrm state.")
	***REMOVED***
	req := h.newNetlinkRequest(nlProto, syscall.NLM_F_CREATE|syscall.NLM_F_EXCL|syscall.NLM_F_ACK)

	msg := xfrmUsersaInfoFromXfrmState(state)

	if state.ESN ***REMOVED***
		if state.ReplayWindow == 0 ***REMOVED***
			return fmt.Errorf("ESN flag set without ReplayWindow")
		***REMOVED***
		msg.Flags |= nl.XFRM_STATE_ESN
		msg.ReplayWindow = 0
	***REMOVED***

	limitsToLft(state.Limits, &msg.Lft)
	req.AddData(msg)

	if state.Auth != nil ***REMOVED***
		out := nl.NewRtAttr(nl.XFRMA_ALG_AUTH_TRUNC, writeStateAlgoAuth(state.Auth))
		req.AddData(out)
	***REMOVED***
	if state.Crypt != nil ***REMOVED***
		out := nl.NewRtAttr(nl.XFRMA_ALG_CRYPT, writeStateAlgo(state.Crypt))
		req.AddData(out)
	***REMOVED***
	if state.Aead != nil ***REMOVED***
		out := nl.NewRtAttr(nl.XFRMA_ALG_AEAD, writeStateAlgoAead(state.Aead))
		req.AddData(out)
	***REMOVED***
	if state.Encap != nil ***REMOVED***
		encapData := make([]byte, nl.SizeofXfrmEncapTmpl)
		encap := nl.DeserializeXfrmEncapTmpl(encapData)
		encap.EncapType = uint16(state.Encap.Type)
		encap.EncapSport = nl.Swap16(uint16(state.Encap.SrcPort))
		encap.EncapDport = nl.Swap16(uint16(state.Encap.DstPort))
		encap.EncapOa.FromIP(state.Encap.OriginalAddress)
		out := nl.NewRtAttr(nl.XFRMA_ENCAP, encapData)
		req.AddData(out)
	***REMOVED***
	if state.Mark != nil ***REMOVED***
		out := nl.NewRtAttr(nl.XFRMA_MARK, writeMark(state.Mark))
		req.AddData(out)
	***REMOVED***
	if state.ESN ***REMOVED***
		out := nl.NewRtAttr(nl.XFRMA_REPLAY_ESN_VAL, writeReplayEsn(state.ReplayWindow))
		req.AddData(out)
	***REMOVED***

	_, err := req.Execute(syscall.NETLINK_XFRM, 0)
	return err
***REMOVED***

func (h *Handle) xfrmStateAllocSpi(state *XfrmState) (*XfrmState, error) ***REMOVED***
	req := h.newNetlinkRequest(nl.XFRM_MSG_ALLOCSPI,
		syscall.NLM_F_CREATE|syscall.NLM_F_EXCL|syscall.NLM_F_ACK)

	msg := &nl.XfrmUserSpiInfo***REMOVED******REMOVED***
	msg.XfrmUsersaInfo = *(xfrmUsersaInfoFromXfrmState(state))
	// 1-255 is reserved by IANA for future use
	msg.Min = 0x100
	msg.Max = 0xffffffff
	req.AddData(msg)

	if state.Mark != nil ***REMOVED***
		out := nl.NewRtAttr(nl.XFRMA_MARK, writeMark(state.Mark))
		req.AddData(out)
	***REMOVED***

	msgs, err := req.Execute(syscall.NETLINK_XFRM, 0)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	s, err := parseXfrmState(msgs[0], FAMILY_ALL)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return s, err
***REMOVED***

// XfrmStateDel will delete an xfrm state from the system. Note that
// the Algos are ignored when matching the state to delete.
// Equivalent to: `ip xfrm state del $state`
func XfrmStateDel(state *XfrmState) error ***REMOVED***
	return pkgHandle.XfrmStateDel(state)
***REMOVED***

// XfrmStateDel will delete an xfrm state from the system. Note that
// the Algos are ignored when matching the state to delete.
// Equivalent to: `ip xfrm state del $state`
func (h *Handle) XfrmStateDel(state *XfrmState) error ***REMOVED***
	_, err := h.xfrmStateGetOrDelete(state, nl.XFRM_MSG_DELSA)
	return err
***REMOVED***

// XfrmStateList gets a list of xfrm states in the system.
// Equivalent to: `ip [-4|-6] xfrm state show`.
// The list can be filtered by ip family.
func XfrmStateList(family int) ([]XfrmState, error) ***REMOVED***
	return pkgHandle.XfrmStateList(family)
***REMOVED***

// XfrmStateList gets a list of xfrm states in the system.
// Equivalent to: `ip xfrm state show`.
// The list can be filtered by ip family.
func (h *Handle) XfrmStateList(family int) ([]XfrmState, error) ***REMOVED***
	req := h.newNetlinkRequest(nl.XFRM_MSG_GETSA, syscall.NLM_F_DUMP)

	msgs, err := req.Execute(syscall.NETLINK_XFRM, nl.XFRM_MSG_NEWSA)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var res []XfrmState
	for _, m := range msgs ***REMOVED***
		if state, err := parseXfrmState(m, family); err == nil ***REMOVED***
			res = append(res, *state)
		***REMOVED*** else if err == familyError ***REMOVED***
			continue
		***REMOVED*** else ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return res, nil
***REMOVED***

// XfrmStateGet gets the xfrm state described by the ID, if found.
// Equivalent to: `ip xfrm state get ID [ mark MARK [ mask MASK ] ]`.
// Only the fields which constitue the SA ID must be filled in:
// ID := [ src ADDR ] [ dst ADDR ] [ proto XFRM-PROTO ] [ spi SPI ]
// mark is optional
func XfrmStateGet(state *XfrmState) (*XfrmState, error) ***REMOVED***
	return pkgHandle.XfrmStateGet(state)
***REMOVED***

// XfrmStateGet gets the xfrm state described by the ID, if found.
// Equivalent to: `ip xfrm state get ID [ mark MARK [ mask MASK ] ]`.
// Only the fields which constitue the SA ID must be filled in:
// ID := [ src ADDR ] [ dst ADDR ] [ proto XFRM-PROTO ] [ spi SPI ]
// mark is optional
func (h *Handle) XfrmStateGet(state *XfrmState) (*XfrmState, error) ***REMOVED***
	return h.xfrmStateGetOrDelete(state, nl.XFRM_MSG_GETSA)
***REMOVED***

func (h *Handle) xfrmStateGetOrDelete(state *XfrmState, nlProto int) (*XfrmState, error) ***REMOVED***
	req := h.newNetlinkRequest(nlProto, syscall.NLM_F_ACK)

	msg := &nl.XfrmUsersaId***REMOVED******REMOVED***
	msg.Family = uint16(nl.GetIPFamily(state.Dst))
	msg.Daddr.FromIP(state.Dst)
	msg.Proto = uint8(state.Proto)
	msg.Spi = nl.Swap32(uint32(state.Spi))
	req.AddData(msg)

	if state.Mark != nil ***REMOVED***
		out := nl.NewRtAttr(nl.XFRMA_MARK, writeMark(state.Mark))
		req.AddData(out)
	***REMOVED***
	if state.Src != nil ***REMOVED***
		out := nl.NewRtAttr(nl.XFRMA_SRCADDR, state.Src.To16())
		req.AddData(out)
	***REMOVED***

	resType := nl.XFRM_MSG_NEWSA
	if nlProto == nl.XFRM_MSG_DELSA ***REMOVED***
		resType = 0
	***REMOVED***

	msgs, err := req.Execute(syscall.NETLINK_XFRM, uint16(resType))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if nlProto == nl.XFRM_MSG_DELSA ***REMOVED***
		return nil, nil
	***REMOVED***

	s, err := parseXfrmState(msgs[0], FAMILY_ALL)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return s, nil
***REMOVED***

var familyError = fmt.Errorf("family error")

func xfrmStateFromXfrmUsersaInfo(msg *nl.XfrmUsersaInfo) *XfrmState ***REMOVED***
	var state XfrmState

	state.Dst = msg.Id.Daddr.ToIP()
	state.Src = msg.Saddr.ToIP()
	state.Proto = Proto(msg.Id.Proto)
	state.Mode = Mode(msg.Mode)
	state.Spi = int(nl.Swap32(msg.Id.Spi))
	state.Reqid = int(msg.Reqid)
	state.ReplayWindow = int(msg.ReplayWindow)
	lftToLimits(&msg.Lft, &state.Limits)

	return &state
***REMOVED***

func parseXfrmState(m []byte, family int) (*XfrmState, error) ***REMOVED***
	msg := nl.DeserializeXfrmUsersaInfo(m)

	// This is mainly for the state dump
	if family != FAMILY_ALL && family != int(msg.Family) ***REMOVED***
		return nil, familyError
	***REMOVED***

	state := xfrmStateFromXfrmUsersaInfo(msg)

	attrs, err := nl.ParseRouteAttr(m[nl.SizeofXfrmUsersaInfo:])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for _, attr := range attrs ***REMOVED***
		switch attr.Attr.Type ***REMOVED***
		case nl.XFRMA_ALG_AUTH, nl.XFRMA_ALG_CRYPT:
			var resAlgo *XfrmStateAlgo
			if attr.Attr.Type == nl.XFRMA_ALG_AUTH ***REMOVED***
				if state.Auth == nil ***REMOVED***
					state.Auth = new(XfrmStateAlgo)
				***REMOVED***
				resAlgo = state.Auth
			***REMOVED*** else ***REMOVED***
				state.Crypt = new(XfrmStateAlgo)
				resAlgo = state.Crypt
			***REMOVED***
			algo := nl.DeserializeXfrmAlgo(attr.Value[:])
			(*resAlgo).Name = nl.BytesToString(algo.AlgName[:])
			(*resAlgo).Key = algo.AlgKey
		case nl.XFRMA_ALG_AUTH_TRUNC:
			if state.Auth == nil ***REMOVED***
				state.Auth = new(XfrmStateAlgo)
			***REMOVED***
			algo := nl.DeserializeXfrmAlgoAuth(attr.Value[:])
			state.Auth.Name = nl.BytesToString(algo.AlgName[:])
			state.Auth.Key = algo.AlgKey
			state.Auth.TruncateLen = int(algo.AlgTruncLen)
		case nl.XFRMA_ALG_AEAD:
			state.Aead = new(XfrmStateAlgo)
			algo := nl.DeserializeXfrmAlgoAEAD(attr.Value[:])
			state.Aead.Name = nl.BytesToString(algo.AlgName[:])
			state.Aead.Key = algo.AlgKey
			state.Aead.ICVLen = int(algo.AlgICVLen)
		case nl.XFRMA_ENCAP:
			encap := nl.DeserializeXfrmEncapTmpl(attr.Value[:])
			state.Encap = new(XfrmStateEncap)
			state.Encap.Type = EncapType(encap.EncapType)
			state.Encap.SrcPort = int(nl.Swap16(encap.EncapSport))
			state.Encap.DstPort = int(nl.Swap16(encap.EncapDport))
			state.Encap.OriginalAddress = encap.EncapOa.ToIP()
		case nl.XFRMA_MARK:
			mark := nl.DeserializeXfrmMark(attr.Value[:])
			state.Mark = new(XfrmMark)
			state.Mark.Value = mark.Value
			state.Mark.Mask = mark.Mask
		***REMOVED***
	***REMOVED***

	return state, nil
***REMOVED***

// XfrmStateFlush will flush the xfrm state on the system.
// proto = 0 means any transformation protocols
// Equivalent to: `ip xfrm state flush [ proto XFRM-PROTO ]`
func XfrmStateFlush(proto Proto) error ***REMOVED***
	return pkgHandle.XfrmStateFlush(proto)
***REMOVED***

// XfrmStateFlush will flush the xfrm state on the system.
// proto = 0 means any transformation protocols
// Equivalent to: `ip xfrm state flush [ proto XFRM-PROTO ]`
func (h *Handle) XfrmStateFlush(proto Proto) error ***REMOVED***
	req := h.newNetlinkRequest(nl.XFRM_MSG_FLUSHSA, syscall.NLM_F_ACK)

	req.AddData(&nl.XfrmUsersaFlush***REMOVED***Proto: uint8(proto)***REMOVED***)

	_, err := req.Execute(syscall.NETLINK_XFRM, 0)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func limitsToLft(lmts XfrmStateLimits, lft *nl.XfrmLifetimeCfg) ***REMOVED***
	if lmts.ByteSoft != 0 ***REMOVED***
		lft.SoftByteLimit = lmts.ByteSoft
	***REMOVED*** else ***REMOVED***
		lft.SoftByteLimit = nl.XFRM_INF
	***REMOVED***
	if lmts.ByteHard != 0 ***REMOVED***
		lft.HardByteLimit = lmts.ByteHard
	***REMOVED*** else ***REMOVED***
		lft.HardByteLimit = nl.XFRM_INF
	***REMOVED***
	if lmts.PacketSoft != 0 ***REMOVED***
		lft.SoftPacketLimit = lmts.PacketSoft
	***REMOVED*** else ***REMOVED***
		lft.SoftPacketLimit = nl.XFRM_INF
	***REMOVED***
	if lmts.PacketHard != 0 ***REMOVED***
		lft.HardPacketLimit = lmts.PacketHard
	***REMOVED*** else ***REMOVED***
		lft.HardPacketLimit = nl.XFRM_INF
	***REMOVED***
	lft.SoftAddExpiresSeconds = lmts.TimeSoft
	lft.HardAddExpiresSeconds = lmts.TimeHard
	lft.SoftUseExpiresSeconds = lmts.TimeUseSoft
	lft.HardUseExpiresSeconds = lmts.TimeUseHard
***REMOVED***

func lftToLimits(lft *nl.XfrmLifetimeCfg, lmts *XfrmStateLimits) ***REMOVED***
	*lmts = *(*XfrmStateLimits)(unsafe.Pointer(lft))
***REMOVED***

func xfrmUsersaInfoFromXfrmState(state *XfrmState) *nl.XfrmUsersaInfo ***REMOVED***
	msg := &nl.XfrmUsersaInfo***REMOVED******REMOVED***
	msg.Family = uint16(nl.GetIPFamily(state.Dst))
	msg.Id.Daddr.FromIP(state.Dst)
	msg.Saddr.FromIP(state.Src)
	msg.Id.Proto = uint8(state.Proto)
	msg.Mode = uint8(state.Mode)
	msg.Id.Spi = nl.Swap32(uint32(state.Spi))
	msg.Reqid = uint32(state.Reqid)
	msg.ReplayWindow = uint8(state.ReplayWindow)

	return msg
***REMOVED***
