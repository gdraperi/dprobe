package netlink

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/vishvananda/netlink/nl"
)

// Constants used in TcU32Sel.Flags.
const (
	TC_U32_TERMINAL  = nl.TC_U32_TERMINAL
	TC_U32_OFFSET    = nl.TC_U32_OFFSET
	TC_U32_VAROFFSET = nl.TC_U32_VAROFFSET
	TC_U32_EAT       = nl.TC_U32_EAT
)

// Fw filter filters on firewall marks
// NOTE: this is in filter_linux because it refers to nl.TcPolice which
//       is defined in nl/tc_linux.go
type Fw struct ***REMOVED***
	FilterAttrs
	ClassId uint32
	// TODO remove nl type from interface
	Police nl.TcPolice
	InDev  string
	// TODO Action
	Mask   uint32
	AvRate uint32
	Rtab   [256]uint32
	Ptab   [256]uint32
***REMOVED***

func NewFw(attrs FilterAttrs, fattrs FilterFwAttrs) (*Fw, error) ***REMOVED***
	var rtab [256]uint32
	var ptab [256]uint32
	rcellLog := -1
	pcellLog := -1
	avrate := fattrs.AvRate / 8
	police := nl.TcPolice***REMOVED******REMOVED***
	police.Rate.Rate = fattrs.Rate / 8
	police.PeakRate.Rate = fattrs.PeakRate / 8
	buffer := fattrs.Buffer
	linklayer := nl.LINKLAYER_ETHERNET

	if fattrs.LinkLayer != nl.LINKLAYER_UNSPEC ***REMOVED***
		linklayer = fattrs.LinkLayer
	***REMOVED***

	police.Action = int32(fattrs.Action)
	if police.Rate.Rate != 0 ***REMOVED***
		police.Rate.Mpu = fattrs.Mpu
		police.Rate.Overhead = fattrs.Overhead
		if CalcRtable(&police.Rate, rtab, rcellLog, fattrs.Mtu, linklayer) < 0 ***REMOVED***
			return nil, errors.New("TBF: failed to calculate rate table")
		***REMOVED***
		police.Burst = uint32(Xmittime(uint64(police.Rate.Rate), uint32(buffer)))
	***REMOVED***
	police.Mtu = fattrs.Mtu
	if police.PeakRate.Rate != 0 ***REMOVED***
		police.PeakRate.Mpu = fattrs.Mpu
		police.PeakRate.Overhead = fattrs.Overhead
		if CalcRtable(&police.PeakRate, ptab, pcellLog, fattrs.Mtu, linklayer) < 0 ***REMOVED***
			return nil, errors.New("POLICE: failed to calculate peak rate table")
		***REMOVED***
	***REMOVED***

	return &Fw***REMOVED***
		FilterAttrs: attrs,
		ClassId:     fattrs.ClassId,
		InDev:       fattrs.InDev,
		Mask:        fattrs.Mask,
		Police:      police,
		AvRate:      avrate,
		Rtab:        rtab,
		Ptab:        ptab,
	***REMOVED***, nil
***REMOVED***

func (filter *Fw) Attrs() *FilterAttrs ***REMOVED***
	return &filter.FilterAttrs
***REMOVED***

func (filter *Fw) Type() string ***REMOVED***
	return "fw"
***REMOVED***

// FilterDel will delete a filter from the system.
// Equivalent to: `tc filter del $filter`
func FilterDel(filter Filter) error ***REMOVED***
	return pkgHandle.FilterDel(filter)
***REMOVED***

// FilterDel will delete a filter from the system.
// Equivalent to: `tc filter del $filter`
func (h *Handle) FilterDel(filter Filter) error ***REMOVED***
	req := h.newNetlinkRequest(syscall.RTM_DELTFILTER, syscall.NLM_F_ACK)
	base := filter.Attrs()
	msg := &nl.TcMsg***REMOVED***
		Family:  nl.FAMILY_ALL,
		Ifindex: int32(base.LinkIndex),
		Handle:  base.Handle,
		Parent:  base.Parent,
		Info:    MakeHandle(base.Priority, nl.Swap16(base.Protocol)),
	***REMOVED***
	req.AddData(msg)

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

// FilterAdd will add a filter to the system.
// Equivalent to: `tc filter add $filter`
func FilterAdd(filter Filter) error ***REMOVED***
	return pkgHandle.FilterAdd(filter)
***REMOVED***

// FilterAdd will add a filter to the system.
// Equivalent to: `tc filter add $filter`
func (h *Handle) FilterAdd(filter Filter) error ***REMOVED***
	native = nl.NativeEndian()
	req := h.newNetlinkRequest(syscall.RTM_NEWTFILTER, syscall.NLM_F_CREATE|syscall.NLM_F_EXCL|syscall.NLM_F_ACK)
	base := filter.Attrs()
	msg := &nl.TcMsg***REMOVED***
		Family:  nl.FAMILY_ALL,
		Ifindex: int32(base.LinkIndex),
		Handle:  base.Handle,
		Parent:  base.Parent,
		Info:    MakeHandle(base.Priority, nl.Swap16(base.Protocol)),
	***REMOVED***
	req.AddData(msg)
	req.AddData(nl.NewRtAttr(nl.TCA_KIND, nl.ZeroTerminated(filter.Type())))

	options := nl.NewRtAttr(nl.TCA_OPTIONS, nil)

	switch filter := filter.(type) ***REMOVED***
	case *U32:
		// Convert TcU32Sel into nl.TcU32Sel as it is without copy.
		sel := (*nl.TcU32Sel)(unsafe.Pointer(filter.Sel))
		if sel == nil ***REMOVED***
			// match all
			sel = &nl.TcU32Sel***REMOVED***
				Nkeys: 1,
				Flags: nl.TC_U32_TERMINAL,
			***REMOVED***
			sel.Keys = append(sel.Keys, nl.TcU32Key***REMOVED******REMOVED***)
		***REMOVED***

		if native != networkOrder ***REMOVED***
			// Copy TcU32Sel.
			cSel := *sel
			keys := make([]nl.TcU32Key, cap(sel.Keys))
			copy(keys, sel.Keys)
			cSel.Keys = keys
			sel = &cSel

			// Handle the endianness of attributes
			sel.Offmask = native.Uint16(htons(sel.Offmask))
			sel.Hmask = native.Uint32(htonl(sel.Hmask))
			for i, key := range sel.Keys ***REMOVED***
				sel.Keys[i].Mask = native.Uint32(htonl(key.Mask))
				sel.Keys[i].Val = native.Uint32(htonl(key.Val))
			***REMOVED***
		***REMOVED***
		sel.Nkeys = uint8(len(sel.Keys))
		nl.NewRtAttrChild(options, nl.TCA_U32_SEL, sel.Serialize())
		if filter.ClassId != 0 ***REMOVED***
			nl.NewRtAttrChild(options, nl.TCA_U32_CLASSID, nl.Uint32Attr(filter.ClassId))
		***REMOVED***
		actionsAttr := nl.NewRtAttrChild(options, nl.TCA_U32_ACT, nil)
		// backwards compatibility
		if filter.RedirIndex != 0 ***REMOVED***
			filter.Actions = append([]Action***REMOVED***NewMirredAction(filter.RedirIndex)***REMOVED***, filter.Actions...)
		***REMOVED***
		if err := EncodeActions(actionsAttr, filter.Actions); err != nil ***REMOVED***
			return err
		***REMOVED***
	case *Fw:
		if filter.Mask != 0 ***REMOVED***
			b := make([]byte, 4)
			native.PutUint32(b, filter.Mask)
			nl.NewRtAttrChild(options, nl.TCA_FW_MASK, b)
		***REMOVED***
		if filter.InDev != "" ***REMOVED***
			nl.NewRtAttrChild(options, nl.TCA_FW_INDEV, nl.ZeroTerminated(filter.InDev))
		***REMOVED***
		if (filter.Police != nl.TcPolice***REMOVED******REMOVED***) ***REMOVED***

			police := nl.NewRtAttrChild(options, nl.TCA_FW_POLICE, nil)
			nl.NewRtAttrChild(police, nl.TCA_POLICE_TBF, filter.Police.Serialize())
			if (filter.Police.Rate != nl.TcRateSpec***REMOVED******REMOVED***) ***REMOVED***
				payload := SerializeRtab(filter.Rtab)
				nl.NewRtAttrChild(police, nl.TCA_POLICE_RATE, payload)
			***REMOVED***
			if (filter.Police.PeakRate != nl.TcRateSpec***REMOVED******REMOVED***) ***REMOVED***
				payload := SerializeRtab(filter.Ptab)
				nl.NewRtAttrChild(police, nl.TCA_POLICE_PEAKRATE, payload)
			***REMOVED***
		***REMOVED***
		if filter.ClassId != 0 ***REMOVED***
			b := make([]byte, 4)
			native.PutUint32(b, filter.ClassId)
			nl.NewRtAttrChild(options, nl.TCA_FW_CLASSID, b)
		***REMOVED***
	case *BpfFilter:
		var bpfFlags uint32
		if filter.ClassId != 0 ***REMOVED***
			nl.NewRtAttrChild(options, nl.TCA_BPF_CLASSID, nl.Uint32Attr(filter.ClassId))
		***REMOVED***
		if filter.Fd >= 0 ***REMOVED***
			nl.NewRtAttrChild(options, nl.TCA_BPF_FD, nl.Uint32Attr((uint32(filter.Fd))))
		***REMOVED***
		if filter.Name != "" ***REMOVED***
			nl.NewRtAttrChild(options, nl.TCA_BPF_NAME, nl.ZeroTerminated(filter.Name))
		***REMOVED***
		if filter.DirectAction ***REMOVED***
			bpfFlags |= nl.TCA_BPF_FLAG_ACT_DIRECT
		***REMOVED***
		nl.NewRtAttrChild(options, nl.TCA_BPF_FLAGS, nl.Uint32Attr(bpfFlags))
	***REMOVED***

	req.AddData(options)
	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

// FilterList gets a list of filters in the system.
// Equivalent to: `tc filter show`.
// Generally returns nothing if link and parent are not specified.
func FilterList(link Link, parent uint32) ([]Filter, error) ***REMOVED***
	return pkgHandle.FilterList(link, parent)
***REMOVED***

// FilterList gets a list of filters in the system.
// Equivalent to: `tc filter show`.
// Generally returns nothing if link and parent are not specified.
func (h *Handle) FilterList(link Link, parent uint32) ([]Filter, error) ***REMOVED***
	req := h.newNetlinkRequest(syscall.RTM_GETTFILTER, syscall.NLM_F_DUMP)
	msg := &nl.TcMsg***REMOVED***
		Family: nl.FAMILY_ALL,
		Parent: parent,
	***REMOVED***
	if link != nil ***REMOVED***
		base := link.Attrs()
		h.ensureIndex(base)
		msg.Ifindex = int32(base.Index)
	***REMOVED***
	req.AddData(msg)

	msgs, err := req.Execute(syscall.NETLINK_ROUTE, syscall.RTM_NEWTFILTER)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var res []Filter
	for _, m := range msgs ***REMOVED***
		msg := nl.DeserializeTcMsg(m)

		attrs, err := nl.ParseRouteAttr(m[msg.Len():])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		base := FilterAttrs***REMOVED***
			LinkIndex: int(msg.Ifindex),
			Handle:    msg.Handle,
			Parent:    msg.Parent,
		***REMOVED***
		base.Priority, base.Protocol = MajorMinor(msg.Info)
		base.Protocol = nl.Swap16(base.Protocol)

		var filter Filter
		filterType := ""
		detailed := false
		for _, attr := range attrs ***REMOVED***
			switch attr.Attr.Type ***REMOVED***
			case nl.TCA_KIND:
				filterType = string(attr.Value[:len(attr.Value)-1])
				switch filterType ***REMOVED***
				case "u32":
					filter = &U32***REMOVED******REMOVED***
				case "fw":
					filter = &Fw***REMOVED******REMOVED***
				case "bpf":
					filter = &BpfFilter***REMOVED******REMOVED***
				default:
					filter = &GenericFilter***REMOVED***FilterType: filterType***REMOVED***
				***REMOVED***
			case nl.TCA_OPTIONS:
				data, err := nl.ParseRouteAttr(attr.Value)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				switch filterType ***REMOVED***
				case "u32":
					detailed, err = parseU32Data(filter, data)
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***
				case "fw":
					detailed, err = parseFwData(filter, data)
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***
				case "bpf":
					detailed, err = parseBpfData(filter, data)
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***
				default:
					detailed = true
				***REMOVED***
			***REMOVED***
		***REMOVED***
		// only return the detailed version of the filter
		if detailed ***REMOVED***
			*filter.Attrs() = base
			res = append(res, filter)
		***REMOVED***
	***REMOVED***

	return res, nil
***REMOVED***

func toTcGen(attrs *ActionAttrs, tcgen *nl.TcGen) ***REMOVED***
	tcgen.Index = uint32(attrs.Index)
	tcgen.Capab = uint32(attrs.Capab)
	tcgen.Action = int32(attrs.Action)
	tcgen.Refcnt = int32(attrs.Refcnt)
	tcgen.Bindcnt = int32(attrs.Bindcnt)
***REMOVED***

func toAttrs(tcgen *nl.TcGen, attrs *ActionAttrs) ***REMOVED***
	attrs.Index = int(tcgen.Index)
	attrs.Capab = int(tcgen.Capab)
	attrs.Action = TcAct(tcgen.Action)
	attrs.Refcnt = int(tcgen.Refcnt)
	attrs.Bindcnt = int(tcgen.Bindcnt)
***REMOVED***

func EncodeActions(attr *nl.RtAttr, actions []Action) error ***REMOVED***
	tabIndex := int(nl.TCA_ACT_TAB)

	for _, action := range actions ***REMOVED***
		switch action := action.(type) ***REMOVED***
		default:
			return fmt.Errorf("unknown action type %s", action.Type())
		case *MirredAction:
			table := nl.NewRtAttrChild(attr, tabIndex, nil)
			tabIndex++
			nl.NewRtAttrChild(table, nl.TCA_ACT_KIND, nl.ZeroTerminated("mirred"))
			aopts := nl.NewRtAttrChild(table, nl.TCA_ACT_OPTIONS, nil)
			mirred := nl.TcMirred***REMOVED***
				Eaction: int32(action.MirredAction),
				Ifindex: uint32(action.Ifindex),
			***REMOVED***
			toTcGen(action.Attrs(), &mirred.TcGen)
			nl.NewRtAttrChild(aopts, nl.TCA_MIRRED_PARMS, mirred.Serialize())
		case *BpfAction:
			table := nl.NewRtAttrChild(attr, tabIndex, nil)
			tabIndex++
			nl.NewRtAttrChild(table, nl.TCA_ACT_KIND, nl.ZeroTerminated("bpf"))
			aopts := nl.NewRtAttrChild(table, nl.TCA_ACT_OPTIONS, nil)
			gen := nl.TcGen***REMOVED******REMOVED***
			toTcGen(action.Attrs(), &gen)
			nl.NewRtAttrChild(aopts, nl.TCA_ACT_BPF_PARMS, gen.Serialize())
			nl.NewRtAttrChild(aopts, nl.TCA_ACT_BPF_FD, nl.Uint32Attr(uint32(action.Fd)))
			nl.NewRtAttrChild(aopts, nl.TCA_ACT_BPF_NAME, nl.ZeroTerminated(action.Name))
		case *GenericAction:
			table := nl.NewRtAttrChild(attr, tabIndex, nil)
			tabIndex++
			nl.NewRtAttrChild(table, nl.TCA_ACT_KIND, nl.ZeroTerminated("gact"))
			aopts := nl.NewRtAttrChild(table, nl.TCA_ACT_OPTIONS, nil)
			gen := nl.TcGen***REMOVED******REMOVED***
			toTcGen(action.Attrs(), &gen)
			nl.NewRtAttrChild(aopts, nl.TCA_GACT_PARMS, gen.Serialize())
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func parseActions(tables []syscall.NetlinkRouteAttr) ([]Action, error) ***REMOVED***
	var actions []Action
	for _, table := range tables ***REMOVED***
		var action Action
		var actionType string
		aattrs, err := nl.ParseRouteAttr(table.Value)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	nextattr:
		for _, aattr := range aattrs ***REMOVED***
			switch aattr.Attr.Type ***REMOVED***
			case nl.TCA_KIND:
				actionType = string(aattr.Value[:len(aattr.Value)-1])
				// only parse if the action is mirred or bpf
				switch actionType ***REMOVED***
				case "mirred":
					action = &MirredAction***REMOVED******REMOVED***
				case "bpf":
					action = &BpfAction***REMOVED******REMOVED***
				case "gact":
					action = &GenericAction***REMOVED******REMOVED***
				default:
					break nextattr
				***REMOVED***
			case nl.TCA_OPTIONS:
				adata, err := nl.ParseRouteAttr(aattr.Value)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				for _, adatum := range adata ***REMOVED***
					switch actionType ***REMOVED***
					case "mirred":
						switch adatum.Attr.Type ***REMOVED***
						case nl.TCA_MIRRED_PARMS:
							mirred := *nl.DeserializeTcMirred(adatum.Value)
							toAttrs(&mirred.TcGen, action.Attrs())
							action.(*MirredAction).ActionAttrs = ActionAttrs***REMOVED******REMOVED***
							action.(*MirredAction).Ifindex = int(mirred.Ifindex)
							action.(*MirredAction).MirredAction = MirredAct(mirred.Eaction)
						***REMOVED***
					case "bpf":
						switch adatum.Attr.Type ***REMOVED***
						case nl.TCA_ACT_BPF_PARMS:
							gen := *nl.DeserializeTcGen(adatum.Value)
							toAttrs(&gen, action.Attrs())
						case nl.TCA_ACT_BPF_FD:
							action.(*BpfAction).Fd = int(native.Uint32(adatum.Value[0:4]))
						case nl.TCA_ACT_BPF_NAME:
							action.(*BpfAction).Name = string(adatum.Value[:len(adatum.Value)-1])
						***REMOVED***
					case "gact":
						switch adatum.Attr.Type ***REMOVED***
						case nl.TCA_GACT_PARMS:
							gen := *nl.DeserializeTcGen(adatum.Value)
							toAttrs(&gen, action.Attrs())
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		actions = append(actions, action)
	***REMOVED***
	return actions, nil
***REMOVED***

func parseU32Data(filter Filter, data []syscall.NetlinkRouteAttr) (bool, error) ***REMOVED***
	native = nl.NativeEndian()
	u32 := filter.(*U32)
	detailed := false
	for _, datum := range data ***REMOVED***
		switch datum.Attr.Type ***REMOVED***
		case nl.TCA_U32_SEL:
			detailed = true
			sel := nl.DeserializeTcU32Sel(datum.Value)
			u32.Sel = (*TcU32Sel)(unsafe.Pointer(sel))
			if native != networkOrder ***REMOVED***
				// Handle the endianness of attributes
				u32.Sel.Offmask = native.Uint16(htons(sel.Offmask))
				u32.Sel.Hmask = native.Uint32(htonl(sel.Hmask))
				for i, key := range u32.Sel.Keys ***REMOVED***
					u32.Sel.Keys[i].Mask = native.Uint32(htonl(key.Mask))
					u32.Sel.Keys[i].Val = native.Uint32(htonl(key.Val))
				***REMOVED***
			***REMOVED***
		case nl.TCA_U32_ACT:
			tables, err := nl.ParseRouteAttr(datum.Value)
			if err != nil ***REMOVED***
				return detailed, err
			***REMOVED***
			u32.Actions, err = parseActions(tables)
			if err != nil ***REMOVED***
				return detailed, err
			***REMOVED***
			for _, action := range u32.Actions ***REMOVED***
				if action, ok := action.(*MirredAction); ok ***REMOVED***
					u32.RedirIndex = int(action.Ifindex)
				***REMOVED***
			***REMOVED***
		case nl.TCA_U32_CLASSID:
			u32.ClassId = native.Uint32(datum.Value)
		***REMOVED***
	***REMOVED***
	return detailed, nil
***REMOVED***

func parseFwData(filter Filter, data []syscall.NetlinkRouteAttr) (bool, error) ***REMOVED***
	native = nl.NativeEndian()
	fw := filter.(*Fw)
	detailed := true
	for _, datum := range data ***REMOVED***
		switch datum.Attr.Type ***REMOVED***
		case nl.TCA_FW_MASK:
			fw.Mask = native.Uint32(datum.Value[0:4])
		case nl.TCA_FW_CLASSID:
			fw.ClassId = native.Uint32(datum.Value[0:4])
		case nl.TCA_FW_INDEV:
			fw.InDev = string(datum.Value[:len(datum.Value)-1])
		case nl.TCA_FW_POLICE:
			adata, _ := nl.ParseRouteAttr(datum.Value)
			for _, aattr := range adata ***REMOVED***
				switch aattr.Attr.Type ***REMOVED***
				case nl.TCA_POLICE_TBF:
					fw.Police = *nl.DeserializeTcPolice(aattr.Value)
				case nl.TCA_POLICE_RATE:
					fw.Rtab = DeserializeRtab(aattr.Value)
				case nl.TCA_POLICE_PEAKRATE:
					fw.Ptab = DeserializeRtab(aattr.Value)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return detailed, nil
***REMOVED***

func parseBpfData(filter Filter, data []syscall.NetlinkRouteAttr) (bool, error) ***REMOVED***
	native = nl.NativeEndian()
	bpf := filter.(*BpfFilter)
	detailed := true
	for _, datum := range data ***REMOVED***
		switch datum.Attr.Type ***REMOVED***
		case nl.TCA_BPF_FD:
			bpf.Fd = int(native.Uint32(datum.Value[0:4]))
		case nl.TCA_BPF_NAME:
			bpf.Name = string(datum.Value[:len(datum.Value)-1])
		case nl.TCA_BPF_CLASSID:
			bpf.ClassId = native.Uint32(datum.Value[0:4])
		case nl.TCA_BPF_FLAGS:
			flags := native.Uint32(datum.Value[0:4])
			if (flags & nl.TCA_BPF_FLAG_ACT_DIRECT) != 0 ***REMOVED***
				bpf.DirectAction = true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return detailed, nil
***REMOVED***

func AlignToAtm(size uint) uint ***REMOVED***
	var linksize, cells int
	cells = int(size / nl.ATM_CELL_PAYLOAD)
	if (size % nl.ATM_CELL_PAYLOAD) > 0 ***REMOVED***
		cells++
	***REMOVED***
	linksize = cells * nl.ATM_CELL_SIZE
	return uint(linksize)
***REMOVED***

func AdjustSize(sz uint, mpu uint, linklayer int) uint ***REMOVED***
	if sz < mpu ***REMOVED***
		sz = mpu
	***REMOVED***
	switch linklayer ***REMOVED***
	case nl.LINKLAYER_ATM:
		return AlignToAtm(sz)
	default:
		return sz
	***REMOVED***
***REMOVED***

func CalcRtable(rate *nl.TcRateSpec, rtab [256]uint32, cellLog int, mtu uint32, linklayer int) int ***REMOVED***
	bps := rate.Rate
	mpu := rate.Mpu
	var sz uint
	if mtu == 0 ***REMOVED***
		mtu = 2047
	***REMOVED***
	if cellLog < 0 ***REMOVED***
		cellLog = 0
		for (mtu >> uint(cellLog)) > 255 ***REMOVED***
			cellLog++
		***REMOVED***
	***REMOVED***
	for i := 0; i < 256; i++ ***REMOVED***
		sz = AdjustSize(uint((i+1)<<uint32(cellLog)), uint(mpu), linklayer)
		rtab[i] = uint32(Xmittime(uint64(bps), uint32(sz)))
	***REMOVED***
	rate.CellAlign = -1
	rate.CellLog = uint8(cellLog)
	rate.Linklayer = uint8(linklayer & nl.TC_LINKLAYER_MASK)
	return cellLog
***REMOVED***

func DeserializeRtab(b []byte) [256]uint32 ***REMOVED***
	var rtab [256]uint32
	native := nl.NativeEndian()
	r := bytes.NewReader(b)
	_ = binary.Read(r, native, &rtab)
	return rtab
***REMOVED***

func SerializeRtab(rtab [256]uint32) []byte ***REMOVED***
	native := nl.NativeEndian()
	var w bytes.Buffer
	_ = binary.Write(&w, native, rtab)
	return w.Bytes()
***REMOVED***
