package netlink

import (
	"syscall"

	"github.com/vishvananda/netlink/nl"
)

func selFromPolicy(sel *nl.XfrmSelector, policy *XfrmPolicy) ***REMOVED***
	sel.Family = uint16(nl.FAMILY_V4)
	if policy.Dst != nil ***REMOVED***
		sel.Family = uint16(nl.GetIPFamily(policy.Dst.IP))
		sel.Daddr.FromIP(policy.Dst.IP)
		prefixlenD, _ := policy.Dst.Mask.Size()
		sel.PrefixlenD = uint8(prefixlenD)
	***REMOVED***
	if policy.Src != nil ***REMOVED***
		sel.Saddr.FromIP(policy.Src.IP)
		prefixlenS, _ := policy.Src.Mask.Size()
		sel.PrefixlenS = uint8(prefixlenS)
	***REMOVED***
	sel.Proto = uint8(policy.Proto)
	sel.Dport = nl.Swap16(uint16(policy.DstPort))
	sel.Sport = nl.Swap16(uint16(policy.SrcPort))
	if sel.Dport != 0 ***REMOVED***
		sel.DportMask = ^uint16(0)
	***REMOVED***
	if sel.Sport != 0 ***REMOVED***
		sel.SportMask = ^uint16(0)
	***REMOVED***
***REMOVED***

// XfrmPolicyAdd will add an xfrm policy to the system.
// Equivalent to: `ip xfrm policy add $policy`
func XfrmPolicyAdd(policy *XfrmPolicy) error ***REMOVED***
	return pkgHandle.XfrmPolicyAdd(policy)
***REMOVED***

// XfrmPolicyAdd will add an xfrm policy to the system.
// Equivalent to: `ip xfrm policy add $policy`
func (h *Handle) XfrmPolicyAdd(policy *XfrmPolicy) error ***REMOVED***
	return h.xfrmPolicyAddOrUpdate(policy, nl.XFRM_MSG_NEWPOLICY)
***REMOVED***

// XfrmPolicyUpdate will update an xfrm policy to the system.
// Equivalent to: `ip xfrm policy update $policy`
func XfrmPolicyUpdate(policy *XfrmPolicy) error ***REMOVED***
	return pkgHandle.XfrmPolicyUpdate(policy)
***REMOVED***

// XfrmPolicyUpdate will update an xfrm policy to the system.
// Equivalent to: `ip xfrm policy update $policy`
func (h *Handle) XfrmPolicyUpdate(policy *XfrmPolicy) error ***REMOVED***
	return h.xfrmPolicyAddOrUpdate(policy, nl.XFRM_MSG_UPDPOLICY)
***REMOVED***

func (h *Handle) xfrmPolicyAddOrUpdate(policy *XfrmPolicy, nlProto int) error ***REMOVED***
	req := h.newNetlinkRequest(nlProto, syscall.NLM_F_CREATE|syscall.NLM_F_EXCL|syscall.NLM_F_ACK)

	msg := &nl.XfrmUserpolicyInfo***REMOVED******REMOVED***
	selFromPolicy(&msg.Sel, policy)
	msg.Priority = uint32(policy.Priority)
	msg.Index = uint32(policy.Index)
	msg.Dir = uint8(policy.Dir)
	msg.Lft.SoftByteLimit = nl.XFRM_INF
	msg.Lft.HardByteLimit = nl.XFRM_INF
	msg.Lft.SoftPacketLimit = nl.XFRM_INF
	msg.Lft.HardPacketLimit = nl.XFRM_INF
	req.AddData(msg)

	tmplData := make([]byte, nl.SizeofXfrmUserTmpl*len(policy.Tmpls))
	for i, tmpl := range policy.Tmpls ***REMOVED***
		start := i * nl.SizeofXfrmUserTmpl
		userTmpl := nl.DeserializeXfrmUserTmpl(tmplData[start : start+nl.SizeofXfrmUserTmpl])
		userTmpl.XfrmId.Daddr.FromIP(tmpl.Dst)
		userTmpl.Saddr.FromIP(tmpl.Src)
		userTmpl.XfrmId.Proto = uint8(tmpl.Proto)
		userTmpl.XfrmId.Spi = nl.Swap32(uint32(tmpl.Spi))
		userTmpl.Mode = uint8(tmpl.Mode)
		userTmpl.Reqid = uint32(tmpl.Reqid)
		userTmpl.Aalgos = ^uint32(0)
		userTmpl.Ealgos = ^uint32(0)
		userTmpl.Calgos = ^uint32(0)
	***REMOVED***
	if len(tmplData) > 0 ***REMOVED***
		tmpls := nl.NewRtAttr(nl.XFRMA_TMPL, tmplData)
		req.AddData(tmpls)
	***REMOVED***
	if policy.Mark != nil ***REMOVED***
		out := nl.NewRtAttr(nl.XFRMA_MARK, writeMark(policy.Mark))
		req.AddData(out)
	***REMOVED***

	_, err := req.Execute(syscall.NETLINK_XFRM, 0)
	return err
***REMOVED***

// XfrmPolicyDel will delete an xfrm policy from the system. Note that
// the Tmpls are ignored when matching the policy to delete.
// Equivalent to: `ip xfrm policy del $policy`
func XfrmPolicyDel(policy *XfrmPolicy) error ***REMOVED***
	return pkgHandle.XfrmPolicyDel(policy)
***REMOVED***

// XfrmPolicyDel will delete an xfrm policy from the system. Note that
// the Tmpls are ignored when matching the policy to delete.
// Equivalent to: `ip xfrm policy del $policy`
func (h *Handle) XfrmPolicyDel(policy *XfrmPolicy) error ***REMOVED***
	_, err := h.xfrmPolicyGetOrDelete(policy, nl.XFRM_MSG_DELPOLICY)
	return err
***REMOVED***

// XfrmPolicyList gets a list of xfrm policies in the system.
// Equivalent to: `ip xfrm policy show`.
// The list can be filtered by ip family.
func XfrmPolicyList(family int) ([]XfrmPolicy, error) ***REMOVED***
	return pkgHandle.XfrmPolicyList(family)
***REMOVED***

// XfrmPolicyList gets a list of xfrm policies in the system.
// Equivalent to: `ip xfrm policy show`.
// The list can be filtered by ip family.
func (h *Handle) XfrmPolicyList(family int) ([]XfrmPolicy, error) ***REMOVED***
	req := h.newNetlinkRequest(nl.XFRM_MSG_GETPOLICY, syscall.NLM_F_DUMP)

	msg := nl.NewIfInfomsg(family)
	req.AddData(msg)

	msgs, err := req.Execute(syscall.NETLINK_XFRM, nl.XFRM_MSG_NEWPOLICY)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var res []XfrmPolicy
	for _, m := range msgs ***REMOVED***
		if policy, err := parseXfrmPolicy(m, family); err == nil ***REMOVED***
			res = append(res, *policy)
		***REMOVED*** else if err == familyError ***REMOVED***
			continue
		***REMOVED*** else ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return res, nil
***REMOVED***

// XfrmPolicyGet gets a the policy described by the index or selector, if found.
// Equivalent to: `ip xfrm policy get ***REMOVED*** SELECTOR | index INDEX ***REMOVED*** dir DIR [ctx CTX ] [ mark MARK [ mask MASK ] ] [ ptype PTYPE ]`.
func XfrmPolicyGet(policy *XfrmPolicy) (*XfrmPolicy, error) ***REMOVED***
	return pkgHandle.XfrmPolicyGet(policy)
***REMOVED***

// XfrmPolicyGet gets a the policy described by the index or selector, if found.
// Equivalent to: `ip xfrm policy get ***REMOVED*** SELECTOR | index INDEX ***REMOVED*** dir DIR [ctx CTX ] [ mark MARK [ mask MASK ] ] [ ptype PTYPE ]`.
func (h *Handle) XfrmPolicyGet(policy *XfrmPolicy) (*XfrmPolicy, error) ***REMOVED***
	return h.xfrmPolicyGetOrDelete(policy, nl.XFRM_MSG_GETPOLICY)
***REMOVED***

// XfrmPolicyFlush will flush the policies on the system.
// Equivalent to: `ip xfrm policy flush`
func XfrmPolicyFlush() error ***REMOVED***
	return pkgHandle.XfrmPolicyFlush()
***REMOVED***

// XfrmPolicyFlush will flush the policies on the system.
// Equivalent to: `ip xfrm policy flush`
func (h *Handle) XfrmPolicyFlush() error ***REMOVED***
	req := h.newNetlinkRequest(nl.XFRM_MSG_FLUSHPOLICY, syscall.NLM_F_ACK)
	_, err := req.Execute(syscall.NETLINK_XFRM, 0)
	return err
***REMOVED***

func (h *Handle) xfrmPolicyGetOrDelete(policy *XfrmPolicy, nlProto int) (*XfrmPolicy, error) ***REMOVED***
	req := h.newNetlinkRequest(nlProto, syscall.NLM_F_ACK)

	msg := &nl.XfrmUserpolicyId***REMOVED******REMOVED***
	selFromPolicy(&msg.Sel, policy)
	msg.Index = uint32(policy.Index)
	msg.Dir = uint8(policy.Dir)
	req.AddData(msg)

	if policy.Mark != nil ***REMOVED***
		out := nl.NewRtAttr(nl.XFRMA_MARK, writeMark(policy.Mark))
		req.AddData(out)
	***REMOVED***

	resType := nl.XFRM_MSG_NEWPOLICY
	if nlProto == nl.XFRM_MSG_DELPOLICY ***REMOVED***
		resType = 0
	***REMOVED***

	msgs, err := req.Execute(syscall.NETLINK_XFRM, uint16(resType))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if nlProto == nl.XFRM_MSG_DELPOLICY ***REMOVED***
		return nil, err
	***REMOVED***

	p, err := parseXfrmPolicy(msgs[0], FAMILY_ALL)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return p, nil
***REMOVED***

func parseXfrmPolicy(m []byte, family int) (*XfrmPolicy, error) ***REMOVED***
	msg := nl.DeserializeXfrmUserpolicyInfo(m)

	// This is mainly for the policy dump
	if family != FAMILY_ALL && family != int(msg.Sel.Family) ***REMOVED***
		return nil, familyError
	***REMOVED***

	var policy XfrmPolicy

	policy.Dst = msg.Sel.Daddr.ToIPNet(msg.Sel.PrefixlenD)
	policy.Src = msg.Sel.Saddr.ToIPNet(msg.Sel.PrefixlenS)
	policy.Proto = Proto(msg.Sel.Proto)
	policy.DstPort = int(nl.Swap16(msg.Sel.Dport))
	policy.SrcPort = int(nl.Swap16(msg.Sel.Sport))
	policy.Priority = int(msg.Priority)
	policy.Index = int(msg.Index)
	policy.Dir = Dir(msg.Dir)

	attrs, err := nl.ParseRouteAttr(m[msg.Len():])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for _, attr := range attrs ***REMOVED***
		switch attr.Attr.Type ***REMOVED***
		case nl.XFRMA_TMPL:
			max := len(attr.Value)
			for i := 0; i < max; i += nl.SizeofXfrmUserTmpl ***REMOVED***
				var resTmpl XfrmPolicyTmpl
				tmpl := nl.DeserializeXfrmUserTmpl(attr.Value[i : i+nl.SizeofXfrmUserTmpl])
				resTmpl.Dst = tmpl.XfrmId.Daddr.ToIP()
				resTmpl.Src = tmpl.Saddr.ToIP()
				resTmpl.Proto = Proto(tmpl.XfrmId.Proto)
				resTmpl.Mode = Mode(tmpl.Mode)
				resTmpl.Spi = int(nl.Swap32(tmpl.XfrmId.Spi))
				resTmpl.Reqid = int(tmpl.Reqid)
				policy.Tmpls = append(policy.Tmpls, resTmpl)
			***REMOVED***
		case nl.XFRMA_MARK:
			mark := nl.DeserializeXfrmMark(attr.Value[:])
			policy.Mark = new(XfrmMark)
			policy.Mark.Value = mark.Value
			policy.Mark.Mask = mark.Mask
		***REMOVED***
	***REMOVED***

	return &policy, nil
***REMOVED***
