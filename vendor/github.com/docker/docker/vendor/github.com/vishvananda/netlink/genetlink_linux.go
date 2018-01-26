package netlink

import (
	"fmt"
	"syscall"

	"github.com/vishvananda/netlink/nl"
)

type GenlOp struct ***REMOVED***
	ID    uint32
	Flags uint32
***REMOVED***

type GenlMulticastGroup struct ***REMOVED***
	ID   uint32
	Name string
***REMOVED***

type GenlFamily struct ***REMOVED***
	ID      uint16
	HdrSize uint32
	Name    string
	Version uint32
	MaxAttr uint32
	Ops     []GenlOp
	Groups  []GenlMulticastGroup
***REMOVED***

func parseOps(b []byte) ([]GenlOp, error) ***REMOVED***
	attrs, err := nl.ParseRouteAttr(b)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ops := make([]GenlOp, 0, len(attrs))
	for _, a := range attrs ***REMOVED***
		nattrs, err := nl.ParseRouteAttr(a.Value)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		var op GenlOp
		for _, na := range nattrs ***REMOVED***
			switch na.Attr.Type ***REMOVED***
			case nl.GENL_CTRL_ATTR_OP_ID:
				op.ID = native.Uint32(na.Value)
			case nl.GENL_CTRL_ATTR_OP_FLAGS:
				op.Flags = native.Uint32(na.Value)
			***REMOVED***
		***REMOVED***
		ops = append(ops, op)
	***REMOVED***
	return ops, nil
***REMOVED***

func parseMulticastGroups(b []byte) ([]GenlMulticastGroup, error) ***REMOVED***
	attrs, err := nl.ParseRouteAttr(b)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	groups := make([]GenlMulticastGroup, 0, len(attrs))
	for _, a := range attrs ***REMOVED***
		nattrs, err := nl.ParseRouteAttr(a.Value)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		var g GenlMulticastGroup
		for _, na := range nattrs ***REMOVED***
			switch na.Attr.Type ***REMOVED***
			case nl.GENL_CTRL_ATTR_MCAST_GRP_NAME:
				g.Name = nl.BytesToString(na.Value)
			case nl.GENL_CTRL_ATTR_MCAST_GRP_ID:
				g.ID = native.Uint32(na.Value)
			***REMOVED***
		***REMOVED***
		groups = append(groups, g)
	***REMOVED***
	return groups, nil
***REMOVED***

func (f *GenlFamily) parseAttributes(attrs []syscall.NetlinkRouteAttr) error ***REMOVED***
	for _, a := range attrs ***REMOVED***
		switch a.Attr.Type ***REMOVED***
		case nl.GENL_CTRL_ATTR_FAMILY_NAME:
			f.Name = nl.BytesToString(a.Value)
		case nl.GENL_CTRL_ATTR_FAMILY_ID:
			f.ID = native.Uint16(a.Value)
		case nl.GENL_CTRL_ATTR_VERSION:
			f.Version = native.Uint32(a.Value)
		case nl.GENL_CTRL_ATTR_HDRSIZE:
			f.HdrSize = native.Uint32(a.Value)
		case nl.GENL_CTRL_ATTR_MAXATTR:
			f.MaxAttr = native.Uint32(a.Value)
		case nl.GENL_CTRL_ATTR_OPS:
			ops, err := parseOps(a.Value)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			f.Ops = ops
		case nl.GENL_CTRL_ATTR_MCAST_GROUPS:
			groups, err := parseMulticastGroups(a.Value)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			f.Groups = groups
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func parseFamilies(msgs [][]byte) ([]*GenlFamily, error) ***REMOVED***
	families := make([]*GenlFamily, 0, len(msgs))
	for _, m := range msgs ***REMOVED***
		attrs, err := nl.ParseRouteAttr(m[nl.SizeofGenlmsg:])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		family := &GenlFamily***REMOVED******REMOVED***
		if err := family.parseAttributes(attrs); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		families = append(families, family)
	***REMOVED***
	return families, nil
***REMOVED***

func (h *Handle) GenlFamilyList() ([]*GenlFamily, error) ***REMOVED***
	msg := &nl.Genlmsg***REMOVED***
		Command: nl.GENL_CTRL_CMD_GETFAMILY,
		Version: nl.GENL_CTRL_VERSION,
	***REMOVED***
	req := h.newNetlinkRequest(nl.GENL_ID_CTRL, syscall.NLM_F_DUMP)
	req.AddData(msg)
	msgs, err := req.Execute(syscall.NETLINK_GENERIC, 0)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return parseFamilies(msgs)
***REMOVED***

func GenlFamilyList() ([]*GenlFamily, error) ***REMOVED***
	return pkgHandle.GenlFamilyList()
***REMOVED***

func (h *Handle) GenlFamilyGet(name string) (*GenlFamily, error) ***REMOVED***
	msg := &nl.Genlmsg***REMOVED***
		Command: nl.GENL_CTRL_CMD_GETFAMILY,
		Version: nl.GENL_CTRL_VERSION,
	***REMOVED***
	req := h.newNetlinkRequest(nl.GENL_ID_CTRL, 0)
	req.AddData(msg)
	req.AddData(nl.NewRtAttr(nl.GENL_CTRL_ATTR_FAMILY_NAME, nl.ZeroTerminated(name)))
	msgs, err := req.Execute(syscall.NETLINK_GENERIC, 0)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	families, err := parseFamilies(msgs)
	if len(families) != 1 ***REMOVED***
		return nil, fmt.Errorf("invalid response for GENL_CTRL_CMD_GETFAMILY")
	***REMOVED***
	return families[0], nil
***REMOVED***

func GenlFamilyGet(name string) (*GenlFamily, error) ***REMOVED***
	return pkgHandle.GenlFamilyGet(name)
***REMOVED***
