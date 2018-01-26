package netlink

import (
	"fmt"
	"net"
	"strings"
	"syscall"

	"github.com/vishvananda/netlink/nl"
)

type PDP struct ***REMOVED***
	Version     uint32
	TID         uint64
	PeerAddress net.IP
	MSAddress   net.IP
	Flow        uint16
	NetNSFD     uint32
	ITEI        uint32
	OTEI        uint32
***REMOVED***

func (pdp *PDP) String() string ***REMOVED***
	elems := []string***REMOVED******REMOVED***
	elems = append(elems, fmt.Sprintf("Version: %d", pdp.Version))
	if pdp.Version == 0 ***REMOVED***
		elems = append(elems, fmt.Sprintf("TID: %d", pdp.TID))
	***REMOVED*** else if pdp.Version == 1 ***REMOVED***
		elems = append(elems, fmt.Sprintf("TEI: %d/%d", pdp.ITEI, pdp.OTEI))
	***REMOVED***
	elems = append(elems, fmt.Sprintf("MS-Address: %s", pdp.MSAddress))
	elems = append(elems, fmt.Sprintf("Peer-Address: %s", pdp.PeerAddress))
	return fmt.Sprintf("***REMOVED***%s***REMOVED***", strings.Join(elems, " "))
***REMOVED***

func (p *PDP) parseAttributes(attrs []syscall.NetlinkRouteAttr) error ***REMOVED***
	for _, a := range attrs ***REMOVED***
		switch a.Attr.Type ***REMOVED***
		case nl.GENL_GTP_ATTR_VERSION:
			p.Version = native.Uint32(a.Value)
		case nl.GENL_GTP_ATTR_TID:
			p.TID = native.Uint64(a.Value)
		case nl.GENL_GTP_ATTR_PEER_ADDRESS:
			p.PeerAddress = net.IP(a.Value)
		case nl.GENL_GTP_ATTR_MS_ADDRESS:
			p.MSAddress = net.IP(a.Value)
		case nl.GENL_GTP_ATTR_FLOW:
			p.Flow = native.Uint16(a.Value)
		case nl.GENL_GTP_ATTR_NET_NS_FD:
			p.NetNSFD = native.Uint32(a.Value)
		case nl.GENL_GTP_ATTR_I_TEI:
			p.ITEI = native.Uint32(a.Value)
		case nl.GENL_GTP_ATTR_O_TEI:
			p.OTEI = native.Uint32(a.Value)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func parsePDP(msgs [][]byte) ([]*PDP, error) ***REMOVED***
	pdps := make([]*PDP, 0, len(msgs))
	for _, m := range msgs ***REMOVED***
		attrs, err := nl.ParseRouteAttr(m[nl.SizeofGenlmsg:])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		pdp := &PDP***REMOVED******REMOVED***
		if err := pdp.parseAttributes(attrs); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		pdps = append(pdps, pdp)
	***REMOVED***
	return pdps, nil
***REMOVED***

func (h *Handle) GTPPDPList() ([]*PDP, error) ***REMOVED***
	f, err := h.GenlFamilyGet(nl.GENL_GTP_NAME)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	msg := &nl.Genlmsg***REMOVED***
		Command: nl.GENL_GTP_CMD_GETPDP,
		Version: nl.GENL_GTP_VERSION,
	***REMOVED***
	req := h.newNetlinkRequest(int(f.ID), syscall.NLM_F_DUMP)
	req.AddData(msg)
	msgs, err := req.Execute(syscall.NETLINK_GENERIC, 0)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return parsePDP(msgs)
***REMOVED***

func GTPPDPList() ([]*PDP, error) ***REMOVED***
	return pkgHandle.GTPPDPList()
***REMOVED***

func gtpPDPGet(req *nl.NetlinkRequest) (*PDP, error) ***REMOVED***
	msgs, err := req.Execute(syscall.NETLINK_GENERIC, 0)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	pdps, err := parsePDP(msgs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(pdps) != 1 ***REMOVED***
		return nil, fmt.Errorf("invalid reqponse for GENL_GTP_CMD_GETPDP")
	***REMOVED***
	return pdps[0], nil
***REMOVED***

func (h *Handle) GTPPDPByTID(link Link, tid int) (*PDP, error) ***REMOVED***
	f, err := h.GenlFamilyGet(nl.GENL_GTP_NAME)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	msg := &nl.Genlmsg***REMOVED***
		Command: nl.GENL_GTP_CMD_GETPDP,
		Version: nl.GENL_GTP_VERSION,
	***REMOVED***
	req := h.newNetlinkRequest(int(f.ID), 0)
	req.AddData(msg)
	req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_VERSION, nl.Uint32Attr(0)))
	req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_LINK, nl.Uint32Attr(uint32(link.Attrs().Index))))
	req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_TID, nl.Uint64Attr(uint64(tid))))
	return gtpPDPGet(req)
***REMOVED***

func GTPPDPByTID(link Link, tid int) (*PDP, error) ***REMOVED***
	return pkgHandle.GTPPDPByTID(link, tid)
***REMOVED***

func (h *Handle) GTPPDPByITEI(link Link, itei int) (*PDP, error) ***REMOVED***
	f, err := h.GenlFamilyGet(nl.GENL_GTP_NAME)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	msg := &nl.Genlmsg***REMOVED***
		Command: nl.GENL_GTP_CMD_GETPDP,
		Version: nl.GENL_GTP_VERSION,
	***REMOVED***
	req := h.newNetlinkRequest(int(f.ID), 0)
	req.AddData(msg)
	req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_VERSION, nl.Uint32Attr(1)))
	req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_LINK, nl.Uint32Attr(uint32(link.Attrs().Index))))
	req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_I_TEI, nl.Uint32Attr(uint32(itei))))
	return gtpPDPGet(req)
***REMOVED***

func GTPPDPByITEI(link Link, itei int) (*PDP, error) ***REMOVED***
	return pkgHandle.GTPPDPByITEI(link, itei)
***REMOVED***

func (h *Handle) GTPPDPByMSAddress(link Link, addr net.IP) (*PDP, error) ***REMOVED***
	f, err := h.GenlFamilyGet(nl.GENL_GTP_NAME)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	msg := &nl.Genlmsg***REMOVED***
		Command: nl.GENL_GTP_CMD_GETPDP,
		Version: nl.GENL_GTP_VERSION,
	***REMOVED***
	req := h.newNetlinkRequest(int(f.ID), 0)
	req.AddData(msg)
	req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_VERSION, nl.Uint32Attr(0)))
	req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_LINK, nl.Uint32Attr(uint32(link.Attrs().Index))))
	req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_MS_ADDRESS, []byte(addr.To4())))
	return gtpPDPGet(req)
***REMOVED***

func GTPPDPByMSAddress(link Link, addr net.IP) (*PDP, error) ***REMOVED***
	return pkgHandle.GTPPDPByMSAddress(link, addr)
***REMOVED***

func (h *Handle) GTPPDPAdd(link Link, pdp *PDP) error ***REMOVED***
	f, err := h.GenlFamilyGet(nl.GENL_GTP_NAME)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	msg := &nl.Genlmsg***REMOVED***
		Command: nl.GENL_GTP_CMD_NEWPDP,
		Version: nl.GENL_GTP_VERSION,
	***REMOVED***
	req := h.newNetlinkRequest(int(f.ID), syscall.NLM_F_EXCL|syscall.NLM_F_ACK)
	req.AddData(msg)
	req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_VERSION, nl.Uint32Attr(pdp.Version)))
	req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_LINK, nl.Uint32Attr(uint32(link.Attrs().Index))))
	req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_PEER_ADDRESS, []byte(pdp.PeerAddress.To4())))
	req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_MS_ADDRESS, []byte(pdp.MSAddress.To4())))

	switch pdp.Version ***REMOVED***
	case 0:
		req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_TID, nl.Uint64Attr(pdp.TID)))
		req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_FLOW, nl.Uint16Attr(pdp.Flow)))
	case 1:
		req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_I_TEI, nl.Uint32Attr(pdp.ITEI)))
		req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_O_TEI, nl.Uint32Attr(pdp.OTEI)))
	default:
		return fmt.Errorf("unsupported GTP version: %d", pdp.Version)
	***REMOVED***
	_, err = req.Execute(syscall.NETLINK_GENERIC, 0)
	return err
***REMOVED***

func GTPPDPAdd(link Link, pdp *PDP) error ***REMOVED***
	return pkgHandle.GTPPDPAdd(link, pdp)
***REMOVED***

func (h *Handle) GTPPDPDel(link Link, pdp *PDP) error ***REMOVED***
	f, err := h.GenlFamilyGet(nl.GENL_GTP_NAME)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	msg := &nl.Genlmsg***REMOVED***
		Command: nl.GENL_GTP_CMD_DELPDP,
		Version: nl.GENL_GTP_VERSION,
	***REMOVED***
	req := h.newNetlinkRequest(int(f.ID), syscall.NLM_F_EXCL|syscall.NLM_F_ACK)
	req.AddData(msg)
	req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_VERSION, nl.Uint32Attr(pdp.Version)))
	req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_LINK, nl.Uint32Attr(uint32(link.Attrs().Index))))

	switch pdp.Version ***REMOVED***
	case 0:
		req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_TID, nl.Uint64Attr(pdp.TID)))
	case 1:
		req.AddData(nl.NewRtAttr(nl.GENL_GTP_ATTR_I_TEI, nl.Uint32Attr(pdp.ITEI)))
	default:
		return fmt.Errorf("unsupported GTP version: %d", pdp.Version)
	***REMOVED***
	_, err = req.Execute(syscall.NETLINK_GENERIC, 0)
	return err
***REMOVED***

func GTPPDPDel(link Link, pdp *PDP) error ***REMOVED***
	return pkgHandle.GTPPDPDel(link, pdp)
***REMOVED***
