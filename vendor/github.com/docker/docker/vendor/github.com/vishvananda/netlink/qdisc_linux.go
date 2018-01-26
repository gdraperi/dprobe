package netlink

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"syscall"

	"github.com/vishvananda/netlink/nl"
)

// NOTE function is here because it uses other linux functions
func NewNetem(attrs QdiscAttrs, nattrs NetemQdiscAttrs) *Netem ***REMOVED***
	var limit uint32 = 1000
	var lossCorr, delayCorr, duplicateCorr uint32
	var reorderProb, reorderCorr uint32
	var corruptProb, corruptCorr uint32

	latency := nattrs.Latency
	loss := Percentage2u32(nattrs.Loss)
	gap := nattrs.Gap
	duplicate := Percentage2u32(nattrs.Duplicate)
	jitter := nattrs.Jitter

	// Correlation
	if latency > 0 && jitter > 0 ***REMOVED***
		delayCorr = Percentage2u32(nattrs.DelayCorr)
	***REMOVED***
	if loss > 0 ***REMOVED***
		lossCorr = Percentage2u32(nattrs.LossCorr)
	***REMOVED***
	if duplicate > 0 ***REMOVED***
		duplicateCorr = Percentage2u32(nattrs.DuplicateCorr)
	***REMOVED***
	// FIXME should validate values(like loss/duplicate are percentages...)
	latency = time2Tick(latency)

	if nattrs.Limit != 0 ***REMOVED***
		limit = nattrs.Limit
	***REMOVED***
	// Jitter is only value if latency is > 0
	if latency > 0 ***REMOVED***
		jitter = time2Tick(jitter)
	***REMOVED***

	reorderProb = Percentage2u32(nattrs.ReorderProb)
	reorderCorr = Percentage2u32(nattrs.ReorderCorr)

	if reorderProb > 0 ***REMOVED***
		// ERROR if lantency == 0
		if gap == 0 ***REMOVED***
			gap = 1
		***REMOVED***
	***REMOVED***

	corruptProb = Percentage2u32(nattrs.CorruptProb)
	corruptCorr = Percentage2u32(nattrs.CorruptCorr)

	return &Netem***REMOVED***
		QdiscAttrs:    attrs,
		Latency:       latency,
		DelayCorr:     delayCorr,
		Limit:         limit,
		Loss:          loss,
		LossCorr:      lossCorr,
		Gap:           gap,
		Duplicate:     duplicate,
		DuplicateCorr: duplicateCorr,
		Jitter:        jitter,
		ReorderProb:   reorderProb,
		ReorderCorr:   reorderCorr,
		CorruptProb:   corruptProb,
		CorruptCorr:   corruptCorr,
	***REMOVED***
***REMOVED***

// QdiscDel will delete a qdisc from the system.
// Equivalent to: `tc qdisc del $qdisc`
func QdiscDel(qdisc Qdisc) error ***REMOVED***
	return pkgHandle.QdiscDel(qdisc)
***REMOVED***

// QdiscDel will delete a qdisc from the system.
// Equivalent to: `tc qdisc del $qdisc`
func (h *Handle) QdiscDel(qdisc Qdisc) error ***REMOVED***
	return h.qdiscModify(syscall.RTM_DELQDISC, 0, qdisc)
***REMOVED***

// QdiscChange will change a qdisc in place
// Equivalent to: `tc qdisc change $qdisc`
// The parent and handle MUST NOT be changed.
func QdiscChange(qdisc Qdisc) error ***REMOVED***
	return pkgHandle.QdiscChange(qdisc)
***REMOVED***

// QdiscChange will change a qdisc in place
// Equivalent to: `tc qdisc change $qdisc`
// The parent and handle MUST NOT be changed.
func (h *Handle) QdiscChange(qdisc Qdisc) error ***REMOVED***
	return h.qdiscModify(syscall.RTM_NEWQDISC, 0, qdisc)
***REMOVED***

// QdiscReplace will replace a qdisc to the system.
// Equivalent to: `tc qdisc replace $qdisc`
// The handle MUST change.
func QdiscReplace(qdisc Qdisc) error ***REMOVED***
	return pkgHandle.QdiscReplace(qdisc)
***REMOVED***

// QdiscReplace will replace a qdisc to the system.
// Equivalent to: `tc qdisc replace $qdisc`
// The handle MUST change.
func (h *Handle) QdiscReplace(qdisc Qdisc) error ***REMOVED***
	return h.qdiscModify(
		syscall.RTM_NEWQDISC,
		syscall.NLM_F_CREATE|syscall.NLM_F_REPLACE,
		qdisc)
***REMOVED***

// QdiscAdd will add a qdisc to the system.
// Equivalent to: `tc qdisc add $qdisc`
func QdiscAdd(qdisc Qdisc) error ***REMOVED***
	return pkgHandle.QdiscAdd(qdisc)
***REMOVED***

// QdiscAdd will add a qdisc to the system.
// Equivalent to: `tc qdisc add $qdisc`
func (h *Handle) QdiscAdd(qdisc Qdisc) error ***REMOVED***
	return h.qdiscModify(
		syscall.RTM_NEWQDISC,
		syscall.NLM_F_CREATE|syscall.NLM_F_EXCL,
		qdisc)
***REMOVED***

func (h *Handle) qdiscModify(cmd, flags int, qdisc Qdisc) error ***REMOVED***
	req := h.newNetlinkRequest(cmd, flags|syscall.NLM_F_ACK)
	base := qdisc.Attrs()
	msg := &nl.TcMsg***REMOVED***
		Family:  nl.FAMILY_ALL,
		Ifindex: int32(base.LinkIndex),
		Handle:  base.Handle,
		Parent:  base.Parent,
	***REMOVED***
	req.AddData(msg)

	// When deleting don't bother building the rest of the netlink payload
	if cmd != syscall.RTM_DELQDISC ***REMOVED***
		if err := qdiscPayload(req, qdisc); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

func qdiscPayload(req *nl.NetlinkRequest, qdisc Qdisc) error ***REMOVED***

	req.AddData(nl.NewRtAttr(nl.TCA_KIND, nl.ZeroTerminated(qdisc.Type())))

	options := nl.NewRtAttr(nl.TCA_OPTIONS, nil)

	switch qdisc := qdisc.(type) ***REMOVED***
	case *Prio:
		tcmap := nl.TcPrioMap***REMOVED***
			Bands:   int32(qdisc.Bands),
			Priomap: qdisc.PriorityMap,
		***REMOVED***
		options = nl.NewRtAttr(nl.TCA_OPTIONS, tcmap.Serialize())
	case *Tbf:
		opt := nl.TcTbfQopt***REMOVED******REMOVED***
		opt.Rate.Rate = uint32(qdisc.Rate)
		opt.Peakrate.Rate = uint32(qdisc.Peakrate)
		opt.Limit = qdisc.Limit
		opt.Buffer = qdisc.Buffer
		nl.NewRtAttrChild(options, nl.TCA_TBF_PARMS, opt.Serialize())
		if qdisc.Rate >= uint64(1<<32) ***REMOVED***
			nl.NewRtAttrChild(options, nl.TCA_TBF_RATE64, nl.Uint64Attr(qdisc.Rate))
		***REMOVED***
		if qdisc.Peakrate >= uint64(1<<32) ***REMOVED***
			nl.NewRtAttrChild(options, nl.TCA_TBF_PRATE64, nl.Uint64Attr(qdisc.Peakrate))
		***REMOVED***
		if qdisc.Peakrate > 0 ***REMOVED***
			nl.NewRtAttrChild(options, nl.TCA_TBF_PBURST, nl.Uint32Attr(qdisc.Minburst))
		***REMOVED***
	case *Htb:
		opt := nl.TcHtbGlob***REMOVED******REMOVED***
		opt.Version = qdisc.Version
		opt.Rate2Quantum = qdisc.Rate2Quantum
		opt.Defcls = qdisc.Defcls
		// TODO: Handle Debug properly. For now default to 0
		opt.Debug = qdisc.Debug
		opt.DirectPkts = qdisc.DirectPkts
		nl.NewRtAttrChild(options, nl.TCA_HTB_INIT, opt.Serialize())
		// nl.NewRtAttrChild(options, nl.TCA_HTB_DIRECT_QLEN, opt.Serialize())
	case *Netem:
		opt := nl.TcNetemQopt***REMOVED******REMOVED***
		opt.Latency = qdisc.Latency
		opt.Limit = qdisc.Limit
		opt.Loss = qdisc.Loss
		opt.Gap = qdisc.Gap
		opt.Duplicate = qdisc.Duplicate
		opt.Jitter = qdisc.Jitter
		options = nl.NewRtAttr(nl.TCA_OPTIONS, opt.Serialize())
		// Correlation
		corr := nl.TcNetemCorr***REMOVED******REMOVED***
		corr.DelayCorr = qdisc.DelayCorr
		corr.LossCorr = qdisc.LossCorr
		corr.DupCorr = qdisc.DuplicateCorr

		if corr.DelayCorr > 0 || corr.LossCorr > 0 || corr.DupCorr > 0 ***REMOVED***
			nl.NewRtAttrChild(options, nl.TCA_NETEM_CORR, corr.Serialize())
		***REMOVED***
		// Corruption
		corruption := nl.TcNetemCorrupt***REMOVED******REMOVED***
		corruption.Probability = qdisc.CorruptProb
		corruption.Correlation = qdisc.CorruptCorr
		if corruption.Probability > 0 ***REMOVED***
			nl.NewRtAttrChild(options, nl.TCA_NETEM_CORRUPT, corruption.Serialize())
		***REMOVED***
		// Reorder
		reorder := nl.TcNetemReorder***REMOVED******REMOVED***
		reorder.Probability = qdisc.ReorderProb
		reorder.Correlation = qdisc.ReorderCorr
		if reorder.Probability > 0 ***REMOVED***
			nl.NewRtAttrChild(options, nl.TCA_NETEM_REORDER, reorder.Serialize())
		***REMOVED***
	case *Ingress:
		// ingress filters must use the proper handle
		if qdisc.Attrs().Parent != HANDLE_INGRESS ***REMOVED***
			return fmt.Errorf("Ingress filters must set Parent to HANDLE_INGRESS")
		***REMOVED***
	***REMOVED***

	req.AddData(options)
	return nil
***REMOVED***

// QdiscList gets a list of qdiscs in the system.
// Equivalent to: `tc qdisc show`.
// The list can be filtered by link.
func QdiscList(link Link) ([]Qdisc, error) ***REMOVED***
	return pkgHandle.QdiscList(link)
***REMOVED***

// QdiscList gets a list of qdiscs in the system.
// Equivalent to: `tc qdisc show`.
// The list can be filtered by link.
func (h *Handle) QdiscList(link Link) ([]Qdisc, error) ***REMOVED***
	req := h.newNetlinkRequest(syscall.RTM_GETQDISC, syscall.NLM_F_DUMP)
	index := int32(0)
	if link != nil ***REMOVED***
		base := link.Attrs()
		h.ensureIndex(base)
		index = int32(base.Index)
	***REMOVED***
	msg := &nl.TcMsg***REMOVED***
		Family:  nl.FAMILY_ALL,
		Ifindex: index,
	***REMOVED***
	req.AddData(msg)

	msgs, err := req.Execute(syscall.NETLINK_ROUTE, syscall.RTM_NEWQDISC)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var res []Qdisc
	for _, m := range msgs ***REMOVED***
		msg := nl.DeserializeTcMsg(m)

		attrs, err := nl.ParseRouteAttr(m[msg.Len():])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		// skip qdiscs from other interfaces
		if link != nil && msg.Ifindex != index ***REMOVED***
			continue
		***REMOVED***

		base := QdiscAttrs***REMOVED***
			LinkIndex: int(msg.Ifindex),
			Handle:    msg.Handle,
			Parent:    msg.Parent,
			Refcnt:    msg.Info,
		***REMOVED***
		var qdisc Qdisc
		qdiscType := ""
		for _, attr := range attrs ***REMOVED***
			switch attr.Attr.Type ***REMOVED***
			case nl.TCA_KIND:
				qdiscType = string(attr.Value[:len(attr.Value)-1])
				switch qdiscType ***REMOVED***
				case "pfifo_fast":
					qdisc = &PfifoFast***REMOVED******REMOVED***
				case "prio":
					qdisc = &Prio***REMOVED******REMOVED***
				case "tbf":
					qdisc = &Tbf***REMOVED******REMOVED***
				case "ingress":
					qdisc = &Ingress***REMOVED******REMOVED***
				case "htb":
					qdisc = &Htb***REMOVED******REMOVED***
				case "netem":
					qdisc = &Netem***REMOVED******REMOVED***
				default:
					qdisc = &GenericQdisc***REMOVED***QdiscType: qdiscType***REMOVED***
				***REMOVED***
			case nl.TCA_OPTIONS:
				switch qdiscType ***REMOVED***
				case "pfifo_fast":
					// pfifo returns TcPrioMap directly without wrapping it in rtattr
					if err := parsePfifoFastData(qdisc, attr.Value); err != nil ***REMOVED***
						return nil, err
					***REMOVED***
				case "prio":
					// prio returns TcPrioMap directly without wrapping it in rtattr
					if err := parsePrioData(qdisc, attr.Value); err != nil ***REMOVED***
						return nil, err
					***REMOVED***
				case "tbf":
					data, err := nl.ParseRouteAttr(attr.Value)
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***
					if err := parseTbfData(qdisc, data); err != nil ***REMOVED***
						return nil, err
					***REMOVED***
				case "htb":
					data, err := nl.ParseRouteAttr(attr.Value)
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***
					if err := parseHtbData(qdisc, data); err != nil ***REMOVED***
						return nil, err
					***REMOVED***
				case "netem":
					if err := parseNetemData(qdisc, attr.Value); err != nil ***REMOVED***
						return nil, err
					***REMOVED***

					// no options for ingress
				***REMOVED***
			***REMOVED***
		***REMOVED***
		*qdisc.Attrs() = base
		res = append(res, qdisc)
	***REMOVED***

	return res, nil
***REMOVED***

func parsePfifoFastData(qdisc Qdisc, value []byte) error ***REMOVED***
	pfifo := qdisc.(*PfifoFast)
	tcmap := nl.DeserializeTcPrioMap(value)
	pfifo.PriorityMap = tcmap.Priomap
	pfifo.Bands = uint8(tcmap.Bands)
	return nil
***REMOVED***

func parsePrioData(qdisc Qdisc, value []byte) error ***REMOVED***
	prio := qdisc.(*Prio)
	tcmap := nl.DeserializeTcPrioMap(value)
	prio.PriorityMap = tcmap.Priomap
	prio.Bands = uint8(tcmap.Bands)
	return nil
***REMOVED***

func parseHtbData(qdisc Qdisc, data []syscall.NetlinkRouteAttr) error ***REMOVED***
	native = nl.NativeEndian()
	htb := qdisc.(*Htb)
	for _, datum := range data ***REMOVED***
		switch datum.Attr.Type ***REMOVED***
		case nl.TCA_HTB_INIT:
			opt := nl.DeserializeTcHtbGlob(datum.Value)
			htb.Version = opt.Version
			htb.Rate2Quantum = opt.Rate2Quantum
			htb.Defcls = opt.Defcls
			htb.Debug = opt.Debug
			htb.DirectPkts = opt.DirectPkts
		case nl.TCA_HTB_DIRECT_QLEN:
			// TODO
			//htb.DirectQlen = native.uint32(datum.Value)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func parseNetemData(qdisc Qdisc, value []byte) error ***REMOVED***
	netem := qdisc.(*Netem)
	opt := nl.DeserializeTcNetemQopt(value)
	netem.Latency = opt.Latency
	netem.Limit = opt.Limit
	netem.Loss = opt.Loss
	netem.Gap = opt.Gap
	netem.Duplicate = opt.Duplicate
	netem.Jitter = opt.Jitter
	data, err := nl.ParseRouteAttr(value[nl.SizeofTcNetemQopt:])
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, datum := range data ***REMOVED***
		switch datum.Attr.Type ***REMOVED***
		case nl.TCA_NETEM_CORR:
			opt := nl.DeserializeTcNetemCorr(datum.Value)
			netem.DelayCorr = opt.DelayCorr
			netem.LossCorr = opt.LossCorr
			netem.DuplicateCorr = opt.DupCorr
		case nl.TCA_NETEM_CORRUPT:
			opt := nl.DeserializeTcNetemCorrupt(datum.Value)
			netem.CorruptProb = opt.Probability
			netem.CorruptCorr = opt.Correlation
		case nl.TCA_NETEM_REORDER:
			opt := nl.DeserializeTcNetemReorder(datum.Value)
			netem.ReorderProb = opt.Probability
			netem.ReorderCorr = opt.Correlation
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func parseTbfData(qdisc Qdisc, data []syscall.NetlinkRouteAttr) error ***REMOVED***
	native = nl.NativeEndian()
	tbf := qdisc.(*Tbf)
	for _, datum := range data ***REMOVED***
		switch datum.Attr.Type ***REMOVED***
		case nl.TCA_TBF_PARMS:
			opt := nl.DeserializeTcTbfQopt(datum.Value)
			tbf.Rate = uint64(opt.Rate.Rate)
			tbf.Peakrate = uint64(opt.Peakrate.Rate)
			tbf.Limit = opt.Limit
			tbf.Buffer = opt.Buffer
		case nl.TCA_TBF_RATE64:
			tbf.Rate = native.Uint64(datum.Value[0:8])
		case nl.TCA_TBF_PRATE64:
			tbf.Peakrate = native.Uint64(datum.Value[0:8])
		case nl.TCA_TBF_PBURST:
			tbf.Minburst = native.Uint32(datum.Value[0:4])
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

const (
	TIME_UNITS_PER_SEC = 1000000
)

var (
	tickInUsec  float64
	clockFactor float64
	hz          float64
)

func initClock() ***REMOVED***
	data, err := ioutil.ReadFile("/proc/net/psched")
	if err != nil ***REMOVED***
		return
	***REMOVED***
	parts := strings.Split(strings.TrimSpace(string(data)), " ")
	if len(parts) < 3 ***REMOVED***
		return
	***REMOVED***
	var vals [3]uint64
	for i := range vals ***REMOVED***
		val, err := strconv.ParseUint(parts[i], 16, 32)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		vals[i] = val
	***REMOVED***
	// compatibility
	if vals[2] == 1000000000 ***REMOVED***
		vals[0] = vals[1]
	***REMOVED***
	clockFactor = float64(vals[2]) / TIME_UNITS_PER_SEC
	tickInUsec = float64(vals[0]) / float64(vals[1]) * clockFactor
	hz = float64(vals[0])
***REMOVED***

func TickInUsec() float64 ***REMOVED***
	if tickInUsec == 0.0 ***REMOVED***
		initClock()
	***REMOVED***
	return tickInUsec
***REMOVED***

func ClockFactor() float64 ***REMOVED***
	if clockFactor == 0.0 ***REMOVED***
		initClock()
	***REMOVED***
	return clockFactor
***REMOVED***

func Hz() float64 ***REMOVED***
	if hz == 0.0 ***REMOVED***
		initClock()
	***REMOVED***
	return hz
***REMOVED***

func time2Tick(time uint32) uint32 ***REMOVED***
	return uint32(float64(time) * TickInUsec())
***REMOVED***

func tick2Time(tick uint32) uint32 ***REMOVED***
	return uint32(float64(tick) / TickInUsec())
***REMOVED***

func time2Ktime(time uint32) uint32 ***REMOVED***
	return uint32(float64(time) * ClockFactor())
***REMOVED***

func ktime2Time(ktime uint32) uint32 ***REMOVED***
	return uint32(float64(ktime) / ClockFactor())
***REMOVED***

func burst(rate uint64, buffer uint32) uint32 ***REMOVED***
	return uint32(float64(rate) * float64(tick2Time(buffer)) / TIME_UNITS_PER_SEC)
***REMOVED***

func latency(rate uint64, limit, buffer uint32) float64 ***REMOVED***
	return TIME_UNITS_PER_SEC*(float64(limit)/float64(rate)) - float64(tick2Time(buffer))
***REMOVED***

func Xmittime(rate uint64, size uint32) float64 ***REMOVED***
	return TickInUsec() * TIME_UNITS_PER_SEC * (float64(size) / float64(rate))
***REMOVED***
