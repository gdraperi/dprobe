package sockaddr

// ifAddrAttrMap is a map of the IfAddr type-specific attributes.
var ifAddrAttrMap map[AttrName]func(IfAddr) string
var ifAddrAttrs []AttrName

func init() ***REMOVED***
	ifAddrAttrInit()
***REMOVED***

// GetPrivateIP returns a string with a single IP address that is part of RFC
// 6890 and has a default route.  If the system can't determine its IP address
// or find an RFC 6890 IP address, an empty string will be returned instead.
// This function is the `eval` equivalent of:
//
// ```
// $ sockaddr eval -r '***REMOVED******REMOVED***GetPrivateInterfaces | attr "address"***REMOVED******REMOVED***'
/// ```
func GetPrivateIP() (string, error) ***REMOVED***
	privateIfs, err := GetPrivateInterfaces()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if len(privateIfs) < 1 ***REMOVED***
		return "", nil
	***REMOVED***

	ifAddr := privateIfs[0]
	ip := *ToIPAddr(ifAddr.SockAddr)
	return ip.NetIP().String(), nil
***REMOVED***

// GetPublicIP returns a string with a single IP address that is NOT part of RFC
// 6890 and has a default route.  If the system can't determine its IP address
// or find a non RFC 6890 IP address, an empty string will be returned instead.
// This function is the `eval` equivalent of:
//
// ```
// $ sockaddr eval -r '***REMOVED******REMOVED***GetPublicInterfaces | attr "address"***REMOVED******REMOVED***'
/// ```
func GetPublicIP() (string, error) ***REMOVED***
	publicIfs, err := GetPublicInterfaces()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED*** else if len(publicIfs) < 1 ***REMOVED***
		return "", nil
	***REMOVED***

	ifAddr := publicIfs[0]
	ip := *ToIPAddr(ifAddr.SockAddr)
	return ip.NetIP().String(), nil
***REMOVED***

// GetInterfaceIP returns a string with a single IP address sorted by the size
// of the network (i.e. IP addresses with a smaller netmask, larger network
// size, are sorted first).  This function is the `eval` equivalent of:
//
// ```
// $ sockaddr eval -r '***REMOVED******REMOVED***GetAllInterfaces | include "name" <<ARG>> | sort "type,size" | include "flag" "forwardable" | attr "address" ***REMOVED******REMOVED***'
/// ```
func GetInterfaceIP(namedIfRE string) (string, error) ***REMOVED***
	ifAddrs, err := GetAllInterfaces()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	ifAddrs, _, err = IfByName(namedIfRE, ifAddrs)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	ifAddrs, _, err = IfByFlag("forwardable", ifAddrs)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	ifAddrs, err = SortIfBy("+type,+size", ifAddrs)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if len(ifAddrs) == 0 ***REMOVED***
		return "", err
	***REMOVED***

	ip := ToIPAddr(ifAddrs[0].SockAddr)
	if ip == nil ***REMOVED***
		return "", err
	***REMOVED***

	return IPAddrAttr(*ip, "address"), nil
***REMOVED***

// IfAddrAttrs returns a list of attributes supported by the IfAddr type
func IfAddrAttrs() []AttrName ***REMOVED***
	return ifAddrAttrs
***REMOVED***

// IfAddrAttr returns a string representation of an attribute for the given
// IfAddr.
func IfAddrAttr(ifAddr IfAddr, attrName AttrName) string ***REMOVED***
	fn, found := ifAddrAttrMap[attrName]
	if !found ***REMOVED***
		return ""
	***REMOVED***

	return fn(ifAddr)
***REMOVED***

// ifAddrAttrInit is called once at init()
func ifAddrAttrInit() ***REMOVED***
	// Sorted for human readability
	ifAddrAttrs = []AttrName***REMOVED***
		"flags",
		"name",
	***REMOVED***

	ifAddrAttrMap = map[AttrName]func(ifAddr IfAddr) string***REMOVED***
		"flags": func(ifAddr IfAddr) string ***REMOVED***
			return ifAddr.Interface.Flags.String()
		***REMOVED***,
		"name": func(ifAddr IfAddr) string ***REMOVED***
			return ifAddr.Interface.Name
		***REMOVED***,
	***REMOVED***
***REMOVED***
