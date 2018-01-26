package netlink

import (
	"strings"
)

// Protinfo represents bridge flags from netlink.
type Protinfo struct ***REMOVED***
	Hairpin      bool
	Guard        bool
	FastLeave    bool
	RootBlock    bool
	Learning     bool
	Flood        bool
	ProxyArp     bool
	ProxyArpWiFi bool
***REMOVED***

// String returns a list of enabled flags
func (prot *Protinfo) String() string ***REMOVED***
	var boolStrings []string
	if prot.Hairpin ***REMOVED***
		boolStrings = append(boolStrings, "Hairpin")
	***REMOVED***
	if prot.Guard ***REMOVED***
		boolStrings = append(boolStrings, "Guard")
	***REMOVED***
	if prot.FastLeave ***REMOVED***
		boolStrings = append(boolStrings, "FastLeave")
	***REMOVED***
	if prot.RootBlock ***REMOVED***
		boolStrings = append(boolStrings, "RootBlock")
	***REMOVED***
	if prot.Learning ***REMOVED***
		boolStrings = append(boolStrings, "Learning")
	***REMOVED***
	if prot.Flood ***REMOVED***
		boolStrings = append(boolStrings, "Flood")
	***REMOVED***
	if prot.ProxyArp ***REMOVED***
		boolStrings = append(boolStrings, "ProxyArp")
	***REMOVED***
	if prot.ProxyArpWiFi ***REMOVED***
		boolStrings = append(boolStrings, "ProxyArpWiFi")
	***REMOVED***
	return strings.Join(boolStrings, " ")
***REMOVED***

func boolToByte(x bool) []byte ***REMOVED***
	if x ***REMOVED***
		return []byte***REMOVED***1***REMOVED***
	***REMOVED***
	return []byte***REMOVED***0***REMOVED***
***REMOVED***

func byteToBool(x byte) bool ***REMOVED***
	return uint8(x) != 0
***REMOVED***
