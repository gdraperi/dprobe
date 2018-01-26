package netlink

import (
	"errors"
	"syscall"

	"github.com/vishvananda/netlink/nl"
)

// NOTE: function is in here because it uses other linux functions
func NewHtbClass(attrs ClassAttrs, cattrs HtbClassAttrs) *HtbClass ***REMOVED***
	mtu := 1600
	rate := cattrs.Rate / 8
	ceil := cattrs.Ceil / 8
	buffer := cattrs.Buffer
	cbuffer := cattrs.Cbuffer

	if ceil == 0 ***REMOVED***
		ceil = rate
	***REMOVED***

	if buffer == 0 ***REMOVED***
		buffer = uint32(float64(rate)/Hz() + float64(mtu))
	***REMOVED***
	buffer = uint32(Xmittime(rate, buffer))

	if cbuffer == 0 ***REMOVED***
		cbuffer = uint32(float64(ceil)/Hz() + float64(mtu))
	***REMOVED***
	cbuffer = uint32(Xmittime(ceil, cbuffer))

	return &HtbClass***REMOVED***
		ClassAttrs: attrs,
		Rate:       rate,
		Ceil:       ceil,
		Buffer:     buffer,
		Cbuffer:    cbuffer,
		Quantum:    10,
		Level:      0,
		Prio:       0,
	***REMOVED***
***REMOVED***

// ClassDel will delete a class from the system.
// Equivalent to: `tc class del $class`
func ClassDel(class Class) error ***REMOVED***
	return pkgHandle.ClassDel(class)
***REMOVED***

// ClassDel will delete a class from the system.
// Equivalent to: `tc class del $class`
func (h *Handle) ClassDel(class Class) error ***REMOVED***
	return h.classModify(syscall.RTM_DELTCLASS, 0, class)
***REMOVED***

// ClassChange will change a class in place
// Equivalent to: `tc class change $class`
// The parent and handle MUST NOT be changed.
func ClassChange(class Class) error ***REMOVED***
	return pkgHandle.ClassChange(class)
***REMOVED***

// ClassChange will change a class in place
// Equivalent to: `tc class change $class`
// The parent and handle MUST NOT be changed.
func (h *Handle) ClassChange(class Class) error ***REMOVED***
	return h.classModify(syscall.RTM_NEWTCLASS, 0, class)
***REMOVED***

// ClassReplace will replace a class to the system.
// quivalent to: `tc class replace $class`
// The handle MAY be changed.
// If a class already exist with this parent/handle pair, the class is changed.
// If a class does not already exist with this parent/handle, a new class is created.
func ClassReplace(class Class) error ***REMOVED***
	return pkgHandle.ClassReplace(class)
***REMOVED***

// ClassReplace will replace a class to the system.
// quivalent to: `tc class replace $class`
// The handle MAY be changed.
// If a class already exist with this parent/handle pair, the class is changed.
// If a class does not already exist with this parent/handle, a new class is created.
func (h *Handle) ClassReplace(class Class) error ***REMOVED***
	return h.classModify(syscall.RTM_NEWTCLASS, syscall.NLM_F_CREATE, class)
***REMOVED***

// ClassAdd will add a class to the system.
// Equivalent to: `tc class add $class`
func ClassAdd(class Class) error ***REMOVED***
	return pkgHandle.ClassAdd(class)
***REMOVED***

// ClassAdd will add a class to the system.
// Equivalent to: `tc class add $class`
func (h *Handle) ClassAdd(class Class) error ***REMOVED***
	return h.classModify(
		syscall.RTM_NEWTCLASS,
		syscall.NLM_F_CREATE|syscall.NLM_F_EXCL,
		class,
	)
***REMOVED***

func (h *Handle) classModify(cmd, flags int, class Class) error ***REMOVED***
	req := h.newNetlinkRequest(cmd, flags|syscall.NLM_F_ACK)
	base := class.Attrs()
	msg := &nl.TcMsg***REMOVED***
		Family:  nl.FAMILY_ALL,
		Ifindex: int32(base.LinkIndex),
		Handle:  base.Handle,
		Parent:  base.Parent,
	***REMOVED***
	req.AddData(msg)

	if cmd != syscall.RTM_DELTCLASS ***REMOVED***
		if err := classPayload(req, class); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

func classPayload(req *nl.NetlinkRequest, class Class) error ***REMOVED***
	req.AddData(nl.NewRtAttr(nl.TCA_KIND, nl.ZeroTerminated(class.Type())))

	options := nl.NewRtAttr(nl.TCA_OPTIONS, nil)
	if htb, ok := class.(*HtbClass); ok ***REMOVED***
		opt := nl.TcHtbCopt***REMOVED******REMOVED***
		opt.Buffer = htb.Buffer
		opt.Cbuffer = htb.Cbuffer
		opt.Quantum = htb.Quantum
		opt.Level = htb.Level
		opt.Prio = htb.Prio
		// TODO: Handle Debug properly. For now default to 0
		/* Calculate ***REMOVED***R,C***REMOVED***Tab and set Rate and Ceil */
		cellLog := -1
		ccellLog := -1
		linklayer := nl.LINKLAYER_ETHERNET
		mtu := 1600
		var rtab [256]uint32
		var ctab [256]uint32
		tcrate := nl.TcRateSpec***REMOVED***Rate: uint32(htb.Rate)***REMOVED***
		if CalcRtable(&tcrate, rtab, cellLog, uint32(mtu), linklayer) < 0 ***REMOVED***
			return errors.New("HTB: failed to calculate rate table")
		***REMOVED***
		opt.Rate = tcrate
		tcceil := nl.TcRateSpec***REMOVED***Rate: uint32(htb.Ceil)***REMOVED***
		if CalcRtable(&tcceil, ctab, ccellLog, uint32(mtu), linklayer) < 0 ***REMOVED***
			return errors.New("HTB: failed to calculate ceil rate table")
		***REMOVED***
		opt.Ceil = tcceil
		nl.NewRtAttrChild(options, nl.TCA_HTB_PARMS, opt.Serialize())
		nl.NewRtAttrChild(options, nl.TCA_HTB_RTAB, SerializeRtab(rtab))
		nl.NewRtAttrChild(options, nl.TCA_HTB_CTAB, SerializeRtab(ctab))
	***REMOVED***
	req.AddData(options)
	return nil
***REMOVED***

// ClassList gets a list of classes in the system.
// Equivalent to: `tc class show`.
// Generally returns nothing if link and parent are not specified.
func ClassList(link Link, parent uint32) ([]Class, error) ***REMOVED***
	return pkgHandle.ClassList(link, parent)
***REMOVED***

// ClassList gets a list of classes in the system.
// Equivalent to: `tc class show`.
// Generally returns nothing if link and parent are not specified.
func (h *Handle) ClassList(link Link, parent uint32) ([]Class, error) ***REMOVED***
	req := h.newNetlinkRequest(syscall.RTM_GETTCLASS, syscall.NLM_F_DUMP)
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

	msgs, err := req.Execute(syscall.NETLINK_ROUTE, syscall.RTM_NEWTCLASS)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var res []Class
	for _, m := range msgs ***REMOVED***
		msg := nl.DeserializeTcMsg(m)

		attrs, err := nl.ParseRouteAttr(m[msg.Len():])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		base := ClassAttrs***REMOVED***
			LinkIndex: int(msg.Ifindex),
			Handle:    msg.Handle,
			Parent:    msg.Parent,
		***REMOVED***

		var class Class
		classType := ""
		for _, attr := range attrs ***REMOVED***
			switch attr.Attr.Type ***REMOVED***
			case nl.TCA_KIND:
				classType = string(attr.Value[:len(attr.Value)-1])
				switch classType ***REMOVED***
				case "htb":
					class = &HtbClass***REMOVED******REMOVED***
				default:
					class = &GenericClass***REMOVED***ClassType: classType***REMOVED***
				***REMOVED***
			case nl.TCA_OPTIONS:
				switch classType ***REMOVED***
				case "htb":
					data, err := nl.ParseRouteAttr(attr.Value)
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***
					_, err = parseHtbClassData(class, data)
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		*class.Attrs() = base
		res = append(res, class)
	***REMOVED***

	return res, nil
***REMOVED***

func parseHtbClassData(class Class, data []syscall.NetlinkRouteAttr) (bool, error) ***REMOVED***
	htb := class.(*HtbClass)
	detailed := false
	for _, datum := range data ***REMOVED***
		switch datum.Attr.Type ***REMOVED***
		case nl.TCA_HTB_PARMS:
			opt := nl.DeserializeTcHtbCopt(datum.Value)
			htb.Rate = uint64(opt.Rate.Rate)
			htb.Ceil = uint64(opt.Ceil.Rate)
			htb.Buffer = opt.Buffer
			htb.Cbuffer = opt.Cbuffer
			htb.Quantum = opt.Quantum
			htb.Level = opt.Level
			htb.Prio = opt.Prio
		***REMOVED***
	***REMOVED***
	return detailed, nil
***REMOVED***
