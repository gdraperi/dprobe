package sockaddr

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// IfAddrs is a slice of IfAddr
type IfAddrs []IfAddr

func (ifs IfAddrs) Len() int ***REMOVED*** return len(ifs) ***REMOVED***

// CmpIfFunc is the function signature that must be met to be used in the
// OrderedIfAddrBy multiIfAddrSorter
type CmpIfAddrFunc func(p1, p2 *IfAddr) int

// multiIfAddrSorter implements the Sort interface, sorting the IfAddrs within.
type multiIfAddrSorter struct ***REMOVED***
	ifAddrs IfAddrs
	cmp     []CmpIfAddrFunc
***REMOVED***

// Sort sorts the argument slice according to the Cmp functions passed to
// OrderedIfAddrBy.
func (ms *multiIfAddrSorter) Sort(ifAddrs IfAddrs) ***REMOVED***
	ms.ifAddrs = ifAddrs
	sort.Sort(ms)
***REMOVED***

// OrderedIfAddrBy sorts SockAddr by the list of sort function pointers.
func OrderedIfAddrBy(cmpFuncs ...CmpIfAddrFunc) *multiIfAddrSorter ***REMOVED***
	return &multiIfAddrSorter***REMOVED***
		cmp: cmpFuncs,
	***REMOVED***
***REMOVED***

// Len is part of sort.Interface.
func (ms *multiIfAddrSorter) Len() int ***REMOVED***
	return len(ms.ifAddrs)
***REMOVED***

// Less is part of sort.Interface. It is implemented by looping along the Cmp()
// functions until it finds a comparison that is either less than or greater
// than.  A return value of 0 defers sorting to the next function in the
// multisorter (which means the results of sorting may leave the resutls in a
// non-deterministic order).
func (ms *multiIfAddrSorter) Less(i, j int) bool ***REMOVED***
	p, q := &ms.ifAddrs[i], &ms.ifAddrs[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.cmp)-1; k++ ***REMOVED***
		cmp := ms.cmp[k]
		x := cmp(p, q)
		switch x ***REMOVED***
		case -1:
			// p < q, so we have a decision.
			return true
		case 1:
			// p > q, so we have a decision.
			return false
		***REMOVED***
		// p == q; try the next comparison.
	***REMOVED***
	// All comparisons to here said "equal", so just return whatever the
	// final comparison reports.
	switch ms.cmp[k](p, q) ***REMOVED***
	case -1:
		return true
	case 1:
		return false
	default:
		// Still a tie! Now what?
		return false
		panic("undefined sort order for remaining items in the list")
	***REMOVED***
***REMOVED***

// Swap is part of sort.Interface.
func (ms *multiIfAddrSorter) Swap(i, j int) ***REMOVED***
	ms.ifAddrs[i], ms.ifAddrs[j] = ms.ifAddrs[j], ms.ifAddrs[i]
***REMOVED***

// AscIfAddress is a sorting function to sort IfAddrs by their respective
// address type.  Non-equal types are deferred in the sort.
func AscIfAddress(p1Ptr, p2Ptr *IfAddr) int ***REMOVED***
	return AscAddress(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
***REMOVED***

// AscIfName is a sorting function to sort IfAddrs by their interface names.
func AscIfName(p1Ptr, p2Ptr *IfAddr) int ***REMOVED***
	return strings.Compare(p1Ptr.Name, p2Ptr.Name)
***REMOVED***

// AscIfNetworkSize is a sorting function to sort IfAddrs by their respective
// network mask size.
func AscIfNetworkSize(p1Ptr, p2Ptr *IfAddr) int ***REMOVED***
	return AscNetworkSize(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
***REMOVED***

// AscIfPort is a sorting function to sort IfAddrs by their respective
// port type.  Non-equal types are deferred in the sort.
func AscIfPort(p1Ptr, p2Ptr *IfAddr) int ***REMOVED***
	return AscPort(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
***REMOVED***

// AscIfPrivate is a sorting function to sort IfAddrs by "private" values before
// "public" values.  Both IPv4 and IPv6 are compared against RFC6890 (RFC6890
// includes, and is not limited to, RFC1918 and RFC6598 for IPv4, and IPv6
// includes RFC4193).
func AscIfPrivate(p1Ptr, p2Ptr *IfAddr) int ***REMOVED***
	return AscPrivate(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
***REMOVED***

// AscIfType is a sorting function to sort IfAddrs by their respective address
// type.  Non-equal types are deferred in the sort.
func AscIfType(p1Ptr, p2Ptr *IfAddr) int ***REMOVED***
	return AscType(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
***REMOVED***

// DescIfAddress is identical to AscIfAddress but reverse ordered.
func DescIfAddress(p1Ptr, p2Ptr *IfAddr) int ***REMOVED***
	return -1 * AscAddress(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
***REMOVED***

// DescIfName is identical to AscIfName but reverse ordered.
func DescIfName(p1Ptr, p2Ptr *IfAddr) int ***REMOVED***
	return -1 * strings.Compare(p1Ptr.Name, p2Ptr.Name)
***REMOVED***

// DescIfNetworkSize is identical to AscIfNetworkSize but reverse ordered.
func DescIfNetworkSize(p1Ptr, p2Ptr *IfAddr) int ***REMOVED***
	return -1 * AscNetworkSize(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
***REMOVED***

// DescIfPort is identical to AscIfPort but reverse ordered.
func DescIfPort(p1Ptr, p2Ptr *IfAddr) int ***REMOVED***
	return -1 * AscPort(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
***REMOVED***

// DescIfPrivate is identical to AscIfPrivate but reverse ordered.
func DescIfPrivate(p1Ptr, p2Ptr *IfAddr) int ***REMOVED***
	return -1 * AscPrivate(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
***REMOVED***

// DescIfType is identical to AscIfType but reverse ordered.
func DescIfType(p1Ptr, p2Ptr *IfAddr) int ***REMOVED***
	return -1 * AscType(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
***REMOVED***

// FilterIfByType filters IfAddrs and returns a list of the matching type
func FilterIfByType(ifAddrs IfAddrs, type_ SockAddrType) (matchedIfs, excludedIfs IfAddrs) ***REMOVED***
	excludedIfs = make(IfAddrs, 0, len(ifAddrs))
	matchedIfs = make(IfAddrs, 0, len(ifAddrs))

	for _, ifAddr := range ifAddrs ***REMOVED***
		if ifAddr.SockAddr.Type()&type_ != 0 ***REMOVED***
			matchedIfs = append(matchedIfs, ifAddr)
		***REMOVED*** else ***REMOVED***
			excludedIfs = append(excludedIfs, ifAddr)
		***REMOVED***
	***REMOVED***
	return matchedIfs, excludedIfs
***REMOVED***

// IfAttr forwards the selector to IfAttr.Attr() for resolution.  If there is
// more than one IfAddr, only the first IfAddr is used.
func IfAttr(selectorName string, ifAddrs IfAddrs) (string, error) ***REMOVED***
	if len(ifAddrs) == 0 ***REMOVED***
		return "", nil
	***REMOVED***

	attrName := AttrName(strings.ToLower(selectorName))
	attrVal, err := ifAddrs[0].Attr(attrName)
	return attrVal, err
***REMOVED***

// GetAllInterfaces iterates over all available network interfaces and finds all
// available IP addresses on each interface and converts them to
// sockaddr.IPAddrs, and returning the result as an array of IfAddr.
func GetAllInterfaces() (IfAddrs, error) ***REMOVED***
	ifs, err := net.Interfaces()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ifAddrs := make(IfAddrs, 0, len(ifs))
	for _, intf := range ifs ***REMOVED***
		addrs, err := intf.Addrs()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		for _, addr := range addrs ***REMOVED***
			var ipAddr IPAddr
			ipAddr, err = NewIPAddr(addr.String())
			if err != nil ***REMOVED***
				return IfAddrs***REMOVED******REMOVED***, fmt.Errorf("unable to create an IP address from %q", addr.String())
			***REMOVED***

			ifAddr := IfAddr***REMOVED***
				SockAddr:  ipAddr,
				Interface: intf,
			***REMOVED***
			ifAddrs = append(ifAddrs, ifAddr)
		***REMOVED***
	***REMOVED***

	return ifAddrs, nil
***REMOVED***

// GetDefaultInterfaces returns IfAddrs of the addresses attached to the default
// route.
func GetDefaultInterfaces() (IfAddrs, error) ***REMOVED***
	ri, err := NewRouteInfo()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	defaultIfName, err := ri.GetDefaultInterfaceName()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var defaultIfs, ifAddrs IfAddrs
	ifAddrs, err = GetAllInterfaces()
	for _, ifAddr := range ifAddrs ***REMOVED***
		if ifAddr.Name == defaultIfName ***REMOVED***
			defaultIfs = append(defaultIfs, ifAddr)
		***REMOVED***
	***REMOVED***

	return defaultIfs, nil
***REMOVED***

// GetPrivateInterfaces returns an IfAddrs that are part of RFC 6890 and have a
// default route.  If the system can't determine its IP address or find an RFC
// 6890 IP address, an empty IfAddrs will be returned instead.  This function is
// the `eval` equivalent of:
//
// ```
// $ sockaddr eval -r '***REMOVED******REMOVED***GetDefaultInterfaces | include "type" "ip" | include "flags" "forwardable|up" | sort "type,size" | include "RFC" "6890" ***REMOVED******REMOVED***'
/// ```
func GetPrivateInterfaces() (IfAddrs, error) ***REMOVED***
	privateIfs, err := GetDefaultInterfaces()
	if err != nil ***REMOVED***
		return IfAddrs***REMOVED******REMOVED***, err
	***REMOVED***
	if len(privateIfs) == 0 ***REMOVED***
		return IfAddrs***REMOVED******REMOVED***, nil
	***REMOVED***

	privateIfs, _ = FilterIfByType(privateIfs, TypeIP)
	if len(privateIfs) == 0 ***REMOVED***
		return IfAddrs***REMOVED******REMOVED***, nil
	***REMOVED***

	privateIfs, _, err = IfByFlag("forwardable|up", privateIfs)
	if err != nil ***REMOVED***
		return IfAddrs***REMOVED******REMOVED***, err
	***REMOVED***
	if len(privateIfs) == 0 ***REMOVED***
		return IfAddrs***REMOVED******REMOVED***, nil
	***REMOVED***

	OrderedIfAddrBy(AscIfType, AscIfNetworkSize).Sort(privateIfs)

	privateIfs, _, err = IfByRFC("6890", privateIfs)
	if err != nil ***REMOVED***
		return IfAddrs***REMOVED******REMOVED***, err
	***REMOVED*** else if len(privateIfs) == 0 ***REMOVED***
		return IfAddrs***REMOVED******REMOVED***, nil
	***REMOVED***

	return privateIfs, nil
***REMOVED***

// GetPublicInterfaces returns an IfAddrs that are NOT part of RFC 6890 and has a
// default route.  If the system can't determine its IP address or find a non
// RFC 6890 IP address, an empty IfAddrs will be returned instead.  This
// function is the `eval` equivalent of:
//
// ```
// $ sockaddr eval -r '***REMOVED******REMOVED***GetDefaultInterfaces | include "type" "ip" | include "flags" "forwardable|up" | sort "type,size" | exclude "RFC" "6890" ***REMOVED******REMOVED***'
/// ```
func GetPublicInterfaces() (IfAddrs, error) ***REMOVED***
	publicIfs, err := GetDefaultInterfaces()
	if err != nil ***REMOVED***
		return IfAddrs***REMOVED******REMOVED***, err
	***REMOVED***
	if len(publicIfs) == 0 ***REMOVED***
		return IfAddrs***REMOVED******REMOVED***, nil
	***REMOVED***

	publicIfs, _ = FilterIfByType(publicIfs, TypeIP)
	if len(publicIfs) == 0 ***REMOVED***
		return IfAddrs***REMOVED******REMOVED***, nil
	***REMOVED***

	publicIfs, _, err = IfByFlag("forwardable|up", publicIfs)
	if err != nil ***REMOVED***
		return IfAddrs***REMOVED******REMOVED***, err
	***REMOVED***
	if len(publicIfs) == 0 ***REMOVED***
		return IfAddrs***REMOVED******REMOVED***, nil
	***REMOVED***

	OrderedIfAddrBy(AscIfType, AscIfNetworkSize).Sort(publicIfs)

	_, publicIfs, err = IfByRFC("6890", publicIfs)
	if err != nil ***REMOVED***
		return IfAddrs***REMOVED******REMOVED***, err
	***REMOVED*** else if len(publicIfs) == 0 ***REMOVED***
		return IfAddrs***REMOVED******REMOVED***, nil
	***REMOVED***

	return publicIfs, nil
***REMOVED***

// IfByAddress returns a list of matched and non-matched IfAddrs, or an error if
// the regexp fails to compile.
func IfByAddress(inputRe string, ifAddrs IfAddrs) (matched, remainder IfAddrs, err error) ***REMOVED***
	re, err := regexp.Compile(inputRe)
	if err != nil ***REMOVED***
		return nil, nil, fmt.Errorf("Unable to compile address regexp %+q: %v", inputRe, err)
	***REMOVED***

	matchedAddrs := make(IfAddrs, 0, len(ifAddrs))
	excludedAddrs := make(IfAddrs, 0, len(ifAddrs))
	for _, addr := range ifAddrs ***REMOVED***
		if re.MatchString(addr.SockAddr.String()) ***REMOVED***
			matchedAddrs = append(matchedAddrs, addr)
		***REMOVED*** else ***REMOVED***
			excludedAddrs = append(excludedAddrs, addr)
		***REMOVED***
	***REMOVED***

	return matchedAddrs, excludedAddrs, nil
***REMOVED***

// IfByName returns a list of matched and non-matched IfAddrs, or an error if
// the regexp fails to compile.
func IfByName(inputRe string, ifAddrs IfAddrs) (matched, remainder IfAddrs, err error) ***REMOVED***
	re, err := regexp.Compile(inputRe)
	if err != nil ***REMOVED***
		return nil, nil, fmt.Errorf("Unable to compile name regexp %+q: %v", inputRe, err)
	***REMOVED***

	matchedAddrs := make(IfAddrs, 0, len(ifAddrs))
	excludedAddrs := make(IfAddrs, 0, len(ifAddrs))
	for _, addr := range ifAddrs ***REMOVED***
		if re.MatchString(addr.Name) ***REMOVED***
			matchedAddrs = append(matchedAddrs, addr)
		***REMOVED*** else ***REMOVED***
			excludedAddrs = append(excludedAddrs, addr)
		***REMOVED***
	***REMOVED***

	return matchedAddrs, excludedAddrs, nil
***REMOVED***

// IfByPort returns a list of matched and non-matched IfAddrs, or an error if
// the regexp fails to compile.
func IfByPort(inputRe string, ifAddrs IfAddrs) (matchedIfs, excludedIfs IfAddrs, err error) ***REMOVED***
	re, err := regexp.Compile(inputRe)
	if err != nil ***REMOVED***
		return nil, nil, fmt.Errorf("Unable to compile port regexp %+q: %v", inputRe, err)
	***REMOVED***

	ipIfs, nonIfs := FilterIfByType(ifAddrs, TypeIP)
	matchedIfs = make(IfAddrs, 0, len(ipIfs))
	excludedIfs = append(IfAddrs(nil), nonIfs...)
	for _, addr := range ipIfs ***REMOVED***
		ipAddr := ToIPAddr(addr.SockAddr)
		if ipAddr == nil ***REMOVED***
			continue
		***REMOVED***

		port := strconv.FormatInt(int64((*ipAddr).IPPort()), 10)
		if re.MatchString(port) ***REMOVED***
			matchedIfs = append(matchedIfs, addr)
		***REMOVED*** else ***REMOVED***
			excludedIfs = append(excludedIfs, addr)
		***REMOVED***
	***REMOVED***

	return matchedIfs, excludedIfs, nil
***REMOVED***

// IfByRFC returns a list of matched and non-matched IfAddrs that contain the
// relevant RFC-specified traits.
func IfByRFC(selectorParam string, ifAddrs IfAddrs) (matched, remainder IfAddrs, err error) ***REMOVED***
	inputRFC, err := strconv.ParseUint(selectorParam, 10, 64)
	if err != nil ***REMOVED***
		return IfAddrs***REMOVED******REMOVED***, IfAddrs***REMOVED******REMOVED***, fmt.Errorf("unable to parse RFC number %q: %v", selectorParam, err)
	***REMOVED***

	matchedIfAddrs := make(IfAddrs, 0, len(ifAddrs))
	remainingIfAddrs := make(IfAddrs, 0, len(ifAddrs))

	rfcNetMap := KnownRFCs()
	rfcNets, ok := rfcNetMap[uint(inputRFC)]
	if !ok ***REMOVED***
		return nil, nil, fmt.Errorf("unsupported RFC %d", inputRFC)
	***REMOVED***

	for _, ifAddr := range ifAddrs ***REMOVED***
		var contained bool
		for _, rfcNet := range rfcNets ***REMOVED***
			if rfcNet.Contains(ifAddr.SockAddr) ***REMOVED***
				matchedIfAddrs = append(matchedIfAddrs, ifAddr)
				contained = true
				break
			***REMOVED***
		***REMOVED***
		if !contained ***REMOVED***
			remainingIfAddrs = append(remainingIfAddrs, ifAddr)
		***REMOVED***
	***REMOVED***

	return matchedIfAddrs, remainingIfAddrs, nil
***REMOVED***

// IfByRFCs returns a list of matched and non-matched IfAddrs that contain the
// relevant RFC-specified traits.  Multiple RFCs can be specified and separated
// by the `|` symbol.  No protection is taken to ensure an IfAddr does not end
// up in both the included and excluded list.
func IfByRFCs(selectorParam string, ifAddrs IfAddrs) (matched, remainder IfAddrs, err error) ***REMOVED***
	var includedIfs, excludedIfs IfAddrs
	for _, rfcStr := range strings.Split(selectorParam, "|") ***REMOVED***
		includedRFCIfs, excludedRFCIfs, err := IfByRFC(rfcStr, ifAddrs)
		if err != nil ***REMOVED***
			return IfAddrs***REMOVED******REMOVED***, IfAddrs***REMOVED******REMOVED***, fmt.Errorf("unable to lookup RFC number %q: %v", rfcStr, err)
		***REMOVED***
		includedIfs = append(includedIfs, includedRFCIfs...)
		excludedIfs = append(excludedIfs, excludedRFCIfs...)
	***REMOVED***

	return includedIfs, excludedIfs, nil
***REMOVED***

// IfByMaskSize returns a list of matched and non-matched IfAddrs that have the
// matching mask size.
func IfByMaskSize(selectorParam string, ifAddrs IfAddrs) (matchedIfs, excludedIfs IfAddrs, err error) ***REMOVED***
	maskSize, err := strconv.ParseUint(selectorParam, 10, 64)
	if err != nil ***REMOVED***
		return IfAddrs***REMOVED******REMOVED***, IfAddrs***REMOVED******REMOVED***, fmt.Errorf("invalid exclude size argument (%q): %v", selectorParam, err)
	***REMOVED***

	ipIfs, nonIfs := FilterIfByType(ifAddrs, TypeIP)
	matchedIfs = make(IfAddrs, 0, len(ipIfs))
	excludedIfs = append(IfAddrs(nil), nonIfs...)
	for _, addr := range ipIfs ***REMOVED***
		ipAddr := ToIPAddr(addr.SockAddr)
		if ipAddr == nil ***REMOVED***
			return IfAddrs***REMOVED******REMOVED***, IfAddrs***REMOVED******REMOVED***, fmt.Errorf("unable to filter mask sizes on non-IP type %s: %v", addr.SockAddr.Type().String(), addr.SockAddr.String())
		***REMOVED***

		switch ***REMOVED***
		case (*ipAddr).Type()&TypeIPv4 != 0 && maskSize > 32:
			return IfAddrs***REMOVED******REMOVED***, IfAddrs***REMOVED******REMOVED***, fmt.Errorf("mask size out of bounds for IPv4 address: %d", maskSize)
		case (*ipAddr).Type()&TypeIPv6 != 0 && maskSize > 128:
			return IfAddrs***REMOVED******REMOVED***, IfAddrs***REMOVED******REMOVED***, fmt.Errorf("mask size out of bounds for IPv6 address: %d", maskSize)
		***REMOVED***

		if (*ipAddr).Maskbits() == int(maskSize) ***REMOVED***
			matchedIfs = append(matchedIfs, addr)
		***REMOVED*** else ***REMOVED***
			excludedIfs = append(excludedIfs, addr)
		***REMOVED***
	***REMOVED***

	return matchedIfs, excludedIfs, nil
***REMOVED***

// IfByType returns a list of matching and non-matching IfAddr that match the
// specified type.  For instance:
//
// include "type" "IPv4,IPv6"
//
// will include any IfAddrs that is either an IPv4 or IPv6 address.  Any
// addresses on those interfaces that don't match will be included in the
// remainder results.
func IfByType(inputTypes string, ifAddrs IfAddrs) (matched, remainder IfAddrs, err error) ***REMOVED***
	matchingIfAddrs := make(IfAddrs, 0, len(ifAddrs))
	remainingIfAddrs := make(IfAddrs, 0, len(ifAddrs))

	ifTypes := strings.Split(strings.ToLower(inputTypes), "|")
	for _, ifType := range ifTypes ***REMOVED***
		switch ifType ***REMOVED***
		case "ip", "ipv4", "ipv6", "unix":
			// Valid types
		default:
			return nil, nil, fmt.Errorf("unsupported type %q %q", ifType, inputTypes)
		***REMOVED***
	***REMOVED***

	for _, ifAddr := range ifAddrs ***REMOVED***
		for _, ifType := range ifTypes ***REMOVED***
			var matched bool
			switch ***REMOVED***
			case ifType == "ip" && ifAddr.SockAddr.Type()&TypeIP != 0:
				matched = true
			case ifType == "ipv4" && ifAddr.SockAddr.Type()&TypeIPv4 != 0:
				matched = true
			case ifType == "ipv6" && ifAddr.SockAddr.Type()&TypeIPv6 != 0:
				matched = true
			case ifType == "unix" && ifAddr.SockAddr.Type()&TypeUnix != 0:
				matched = true
			***REMOVED***

			if matched ***REMOVED***
				matchingIfAddrs = append(matchingIfAddrs, ifAddr)
			***REMOVED*** else ***REMOVED***
				remainingIfAddrs = append(remainingIfAddrs, ifAddr)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return matchingIfAddrs, remainingIfAddrs, nil
***REMOVED***

// IfByFlag returns a list of matching and non-matching IfAddrs that match the
// specified type.  For instance:
//
// include "flag" "up,broadcast"
//
// will include any IfAddrs that have both the "up" and "broadcast" flags set.
// Any addresses on those interfaces that don't match will be omitted from the
// results.
func IfByFlag(inputFlags string, ifAddrs IfAddrs) (matched, remainder IfAddrs, err error) ***REMOVED***
	matchedAddrs := make(IfAddrs, 0, len(ifAddrs))
	excludedAddrs := make(IfAddrs, 0, len(ifAddrs))

	var wantForwardable,
		wantGlobalUnicast,
		wantInterfaceLocalMulticast,
		wantLinkLocalMulticast,
		wantLinkLocalUnicast,
		wantLoopback,
		wantMulticast,
		wantUnspecified bool
	var ifFlags net.Flags
	var checkFlags, checkAttrs bool
	for _, flagName := range strings.Split(strings.ToLower(inputFlags), "|") ***REMOVED***
		switch flagName ***REMOVED***
		case "broadcast":
			checkFlags = true
			ifFlags = ifFlags | net.FlagBroadcast
		case "down":
			checkFlags = true
			ifFlags = (ifFlags &^ net.FlagUp)
		case "forwardable":
			checkAttrs = true
			wantForwardable = true
		case "global unicast":
			checkAttrs = true
			wantGlobalUnicast = true
		case "interface-local multicast":
			checkAttrs = true
			wantInterfaceLocalMulticast = true
		case "link-local multicast":
			checkAttrs = true
			wantLinkLocalMulticast = true
		case "link-local unicast":
			checkAttrs = true
			wantLinkLocalUnicast = true
		case "loopback":
			checkAttrs = true
			checkFlags = true
			ifFlags = ifFlags | net.FlagLoopback
			wantLoopback = true
		case "multicast":
			checkAttrs = true
			checkFlags = true
			ifFlags = ifFlags | net.FlagMulticast
			wantMulticast = true
		case "point-to-point":
			checkFlags = true
			ifFlags = ifFlags | net.FlagPointToPoint
		case "unspecified":
			checkAttrs = true
			wantUnspecified = true
		case "up":
			checkFlags = true
			ifFlags = ifFlags | net.FlagUp
		default:
			return nil, nil, fmt.Errorf("Unknown interface flag: %+q", flagName)
		***REMOVED***
	***REMOVED***

	for _, ifAddr := range ifAddrs ***REMOVED***
		var matched bool
		if checkFlags && ifAddr.Interface.Flags&ifFlags == ifFlags ***REMOVED***
			matched = true
		***REMOVED***
		if checkAttrs ***REMOVED***
			if ip := ToIPAddr(ifAddr.SockAddr); ip != nil ***REMOVED***
				netIP := (*ip).NetIP()
				switch ***REMOVED***
				case wantGlobalUnicast && netIP.IsGlobalUnicast():
					matched = true
				case wantInterfaceLocalMulticast && netIP.IsInterfaceLocalMulticast():
					matched = true
				case wantLinkLocalMulticast && netIP.IsLinkLocalMulticast():
					matched = true
				case wantLinkLocalUnicast && netIP.IsLinkLocalUnicast():
					matched = true
				case wantLoopback && netIP.IsLoopback():
					matched = true
				case wantMulticast && netIP.IsMulticast():
					matched = true
				case wantUnspecified && netIP.IsUnspecified():
					matched = true
				case wantForwardable && !IsRFC(ForwardingBlacklist, ifAddr.SockAddr):
					matched = true
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if matched ***REMOVED***
			matchedAddrs = append(matchedAddrs, ifAddr)
		***REMOVED*** else ***REMOVED***
			excludedAddrs = append(excludedAddrs, ifAddr)
		***REMOVED***
	***REMOVED***
	return matchedAddrs, excludedAddrs, nil
***REMOVED***

// IfByNetwork returns an IfAddrs that are equal to or included within the
// network passed in by selector.
func IfByNetwork(selectorParam string, inputIfAddrs IfAddrs) (IfAddrs, IfAddrs, error) ***REMOVED***
	var includedIfs, excludedIfs IfAddrs
	for _, netStr := range strings.Split(selectorParam, "|") ***REMOVED***
		netAddr, err := NewIPAddr(netStr)
		if err != nil ***REMOVED***
			return nil, nil, fmt.Errorf("unable to create an IP address from %+q: %v", netStr, err)
		***REMOVED***

		for _, ifAddr := range inputIfAddrs ***REMOVED***
			if netAddr.Contains(ifAddr.SockAddr) ***REMOVED***
				includedIfs = append(includedIfs, ifAddr)
			***REMOVED*** else ***REMOVED***
				excludedIfs = append(excludedIfs, ifAddr)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return includedIfs, excludedIfs, nil
***REMOVED***

// IncludeIfs returns an IfAddrs based on the passed in selector.
func IncludeIfs(selectorName, selectorParam string, inputIfAddrs IfAddrs) (IfAddrs, error) ***REMOVED***
	var includedIfs IfAddrs
	var err error

	switch strings.ToLower(selectorName) ***REMOVED***
	case "address":
		includedIfs, _, err = IfByAddress(selectorParam, inputIfAddrs)
	case "flag", "flags":
		includedIfs, _, err = IfByFlag(selectorParam, inputIfAddrs)
	case "name":
		includedIfs, _, err = IfByName(selectorParam, inputIfAddrs)
	case "network":
		includedIfs, _, err = IfByNetwork(selectorParam, inputIfAddrs)
	case "port":
		includedIfs, _, err = IfByPort(selectorParam, inputIfAddrs)
	case "rfc", "rfcs":
		includedIfs, _, err = IfByRFCs(selectorParam, inputIfAddrs)
	case "size":
		includedIfs, _, err = IfByMaskSize(selectorParam, inputIfAddrs)
	case "type":
		includedIfs, _, err = IfByType(selectorParam, inputIfAddrs)
	default:
		return IfAddrs***REMOVED******REMOVED***, fmt.Errorf("invalid include selector %q", selectorName)
	***REMOVED***

	if err != nil ***REMOVED***
		return IfAddrs***REMOVED******REMOVED***, err
	***REMOVED***

	return includedIfs, nil
***REMOVED***

// ExcludeIfs returns an IfAddrs based on the passed in selector.
func ExcludeIfs(selectorName, selectorParam string, inputIfAddrs IfAddrs) (IfAddrs, error) ***REMOVED***
	var excludedIfs IfAddrs
	var err error

	switch strings.ToLower(selectorName) ***REMOVED***
	case "address":
		_, excludedIfs, err = IfByAddress(selectorParam, inputIfAddrs)
	case "flag", "flags":
		_, excludedIfs, err = IfByFlag(selectorParam, inputIfAddrs)
	case "name":
		_, excludedIfs, err = IfByName(selectorParam, inputIfAddrs)
	case "network":
		_, excludedIfs, err = IfByNetwork(selectorParam, inputIfAddrs)
	case "port":
		_, excludedIfs, err = IfByPort(selectorParam, inputIfAddrs)
	case "rfc", "rfcs":
		_, excludedIfs, err = IfByRFCs(selectorParam, inputIfAddrs)
	case "size":
		_, excludedIfs, err = IfByMaskSize(selectorParam, inputIfAddrs)
	case "type":
		_, excludedIfs, err = IfByType(selectorParam, inputIfAddrs)
	default:
		return IfAddrs***REMOVED******REMOVED***, fmt.Errorf("invalid exclude selector %q", selectorName)
	***REMOVED***

	if err != nil ***REMOVED***
		return IfAddrs***REMOVED******REMOVED***, err
	***REMOVED***

	return excludedIfs, nil
***REMOVED***

// SortIfBy returns an IfAddrs sorted based on the passed in selector.  Multiple
// sort clauses can be passed in as a comma delimited list without whitespace.
func SortIfBy(selectorParam string, inputIfAddrs IfAddrs) (IfAddrs, error) ***REMOVED***
	sortedIfs := append(IfAddrs(nil), inputIfAddrs...)

	clauses := strings.Split(selectorParam, ",")
	sortFuncs := make([]CmpIfAddrFunc, len(clauses))

	for i, clause := range clauses ***REMOVED***
		switch strings.TrimSpace(strings.ToLower(clause)) ***REMOVED***
		case "+address", "address":
			// The "address" selector returns an array of IfAddrs
			// ordered by the network address.  IfAddrs that are not
			// comparable will be at the end of the list and in a
			// non-deterministic order.
			sortFuncs[i] = AscIfAddress
		case "-address":
			sortFuncs[i] = DescIfAddress
		case "+name", "name":
			// The "name" selector returns an array of IfAddrs
			// ordered by the interface name.
			sortFuncs[i] = AscIfName
		case "-name":
			sortFuncs[i] = DescIfName
		case "+port", "port":
			// The "port" selector returns an array of IfAddrs
			// ordered by the port, if included in the IfAddr.
			// IfAddrs that are not comparable will be at the end of
			// the list and in a non-deterministic order.
			sortFuncs[i] = AscIfPort
		case "-port":
			sortFuncs[i] = DescIfPort
		case "+private", "private":
			// The "private" selector returns an array of IfAddrs
			// ordered by private addresses first.  IfAddrs that are
			// not comparable will be at the end of the list and in
			// a non-deterministic order.
			sortFuncs[i] = AscIfPrivate
		case "-private":
			sortFuncs[i] = DescIfPrivate
		case "+size", "size":
			// The "size" selector returns an array of IfAddrs
			// ordered by the size of the network mask, smaller mask
			// (larger number of hosts per network) to largest
			// (e.g. a /24 sorts before a /32).
			sortFuncs[i] = AscIfNetworkSize
		case "-size":
			sortFuncs[i] = DescIfNetworkSize
		case "+type", "type":
			// The "type" selector returns an array of IfAddrs
			// ordered by the type of the IfAddr.  The sort order is
			// Unix, IPv4, then IPv6.
			sortFuncs[i] = AscIfType
		case "-type":
			sortFuncs[i] = DescIfType
		default:
			// Return an empty list for invalid sort types.
			return IfAddrs***REMOVED******REMOVED***, fmt.Errorf("unknown sort type: %q", clause)
		***REMOVED***
	***REMOVED***

	OrderedIfAddrBy(sortFuncs...).Sort(sortedIfs)

	return sortedIfs, nil
***REMOVED***

// UniqueIfAddrsBy creates a unique set of IfAddrs based on the matching
// selector.  UniqueIfAddrsBy assumes the input has already been sorted.
func UniqueIfAddrsBy(selectorName string, inputIfAddrs IfAddrs) (IfAddrs, error) ***REMOVED***
	attrName := strings.ToLower(selectorName)

	ifs := make(IfAddrs, 0, len(inputIfAddrs))
	var lastMatch string
	for _, ifAddr := range inputIfAddrs ***REMOVED***
		var out string
		switch attrName ***REMOVED***
		case "address":
			out = ifAddr.SockAddr.String()
		case "name":
			out = ifAddr.Name
		default:
			return nil, fmt.Errorf("unsupported unique constraint %+q", selectorName)
		***REMOVED***

		switch ***REMOVED***
		case lastMatch == "", lastMatch != out:
			lastMatch = out
			ifs = append(ifs, ifAddr)
		case lastMatch == out:
			continue
		***REMOVED***
	***REMOVED***

	return ifs, nil
***REMOVED***

// JoinIfAddrs joins an IfAddrs and returns a string
func JoinIfAddrs(selectorName string, joinStr string, inputIfAddrs IfAddrs) (string, error) ***REMOVED***
	outputs := make([]string, 0, len(inputIfAddrs))
	attrName := AttrName(strings.ToLower(selectorName))

	for _, ifAddr := range inputIfAddrs ***REMOVED***
		var attrVal string
		var err error
		attrVal, err = ifAddr.Attr(attrName)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		outputs = append(outputs, attrVal)
	***REMOVED***
	return strings.Join(outputs, joinStr), nil
***REMOVED***

// LimitIfAddrs returns a slice of IfAddrs based on the specified limit.
func LimitIfAddrs(lim uint, in IfAddrs) (IfAddrs, error) ***REMOVED***
	// Clamp the limit to the length of the array
	if int(lim) > len(in) ***REMOVED***
		lim = uint(len(in))
	***REMOVED***

	return in[0:lim], nil
***REMOVED***

// OffsetIfAddrs returns a slice of IfAddrs based on the specified offset.
func OffsetIfAddrs(off int, in IfAddrs) (IfAddrs, error) ***REMOVED***
	var end bool
	if off < 0 ***REMOVED***
		end = true
		off = off * -1
	***REMOVED***

	if off > len(in) ***REMOVED***
		return IfAddrs***REMOVED******REMOVED***, fmt.Errorf("unable to seek past the end of the interface array: offset (%d) exceeds the number of interfaces (%d)", off, len(in))
	***REMOVED***

	if end ***REMOVED***
		return in[len(in)-off:], nil
	***REMOVED***
	return in[off:], nil
***REMOVED***

func (ifAddr IfAddr) String() string ***REMOVED***
	return fmt.Sprintf("%s %v", ifAddr.SockAddr, ifAddr.Interface)
***REMOVED***

// parseDefaultIfNameFromRoute parses standard route(8)'s output for the *BSDs
// and Solaris.
func parseDefaultIfNameFromRoute(routeOut string) (string, error) ***REMOVED***
	lines := strings.Split(routeOut, "\n")
	for _, line := range lines ***REMOVED***
		kvs := strings.SplitN(line, ":", 2)
		if len(kvs) != 2 ***REMOVED***
			continue
		***REMOVED***

		if strings.TrimSpace(kvs[0]) == "interface" ***REMOVED***
			ifName := strings.TrimSpace(kvs[1])
			return ifName, nil
		***REMOVED***
	***REMOVED***

	return "", errors.New("No default interface found")
***REMOVED***

// parseDefaultIfNameFromIPCmd parses the default interface from ip(8) for
// Linux.
func parseDefaultIfNameFromIPCmd(routeOut string) (string, error) ***REMOVED***
	lines := strings.Split(routeOut, "\n")
	re := regexp.MustCompile(`[\s]+`)
	for _, line := range lines ***REMOVED***
		kvs := re.Split(line, -1)
		if len(kvs) < 5 ***REMOVED***
			continue
		***REMOVED***

		if kvs[0] == "default" &&
			kvs[1] == "via" &&
			kvs[3] == "dev" ***REMOVED***
			ifName := strings.TrimSpace(kvs[4])
			return ifName, nil
		***REMOVED***
	***REMOVED***

	return "", errors.New("No default interface found")
***REMOVED***

// parseDefaultIfNameWindows parses the default interface from `netstat -rn` and
// `ipconfig` on Windows.
func parseDefaultIfNameWindows(routeOut, ipconfigOut string) (string, error) ***REMOVED***
	defaultIPAddr, err := parseDefaultIPAddrWindowsRoute(routeOut)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	ifName, err := parseDefaultIfNameWindowsIPConfig(defaultIPAddr, ipconfigOut)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return ifName, nil
***REMOVED***

// parseDefaultIPAddrWindowsRoute parses the IP address on the default interface
// `netstat -rn`.
//
// NOTES(sean): Only IPv4 addresses are parsed at this time.  If you have an
// IPv6 connected host, submit an issue on github.com/hashicorp/go-sockaddr with
// the output from `netstat -rn`, `ipconfig`, and version of Windows to see IPv6
// support added.
func parseDefaultIPAddrWindowsRoute(routeOut string) (string, error) ***REMOVED***
	lines := strings.Split(routeOut, "\n")
	re := regexp.MustCompile(`[\s]+`)
	for _, line := range lines ***REMOVED***
		kvs := re.Split(strings.TrimSpace(line), -1)
		if len(kvs) < 3 ***REMOVED***
			continue
		***REMOVED***

		if kvs[0] == "0.0.0.0" && kvs[1] == "0.0.0.0" ***REMOVED***
			defaultIPAddr := strings.TrimSpace(kvs[3])
			return defaultIPAddr, nil
		***REMOVED***
	***REMOVED***

	return "", errors.New("No IP on default interface found")
***REMOVED***

// parseDefaultIfNameWindowsIPConfig parses the output of `ipconfig` to find the
// interface name forwarding traffic to the default gateway.
func parseDefaultIfNameWindowsIPConfig(defaultIPAddr, routeOut string) (string, error) ***REMOVED***
	lines := strings.Split(routeOut, "\n")
	ifNameRE := regexp.MustCompile(`^Ethernet adapter ([^\s:]+):`)
	ipAddrRE := regexp.MustCompile(`^   IPv[46] Address\. \. \. \. \. \. \. \. \. \. \. : ([^\s]+)`)
	var ifName string
	for _, line := range lines ***REMOVED***
		switch ifNameMatches := ifNameRE.FindStringSubmatch(line); ***REMOVED***
		case len(ifNameMatches) > 1:
			ifName = ifNameMatches[1]
			continue
		***REMOVED***

		switch ipAddrMatches := ipAddrRE.FindStringSubmatch(line); ***REMOVED***
		case len(ipAddrMatches) > 1 && ipAddrMatches[1] == defaultIPAddr:
			return ifName, nil
		***REMOVED***
	***REMOVED***

	return "", errors.New("No default interface found with matching IP")
***REMOVED***
