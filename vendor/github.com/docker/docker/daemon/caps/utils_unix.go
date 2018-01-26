// +build !windows

package caps

import (
	"fmt"
	"strings"

	"github.com/syndtr/gocapability/capability"
)

var capabilityList Capabilities

func init() ***REMOVED***
	last := capability.CAP_LAST_CAP
	// hack for RHEL6 which has no /proc/sys/kernel/cap_last_cap
	if last == capability.Cap(63) ***REMOVED***
		last = capability.CAP_BLOCK_SUSPEND
	***REMOVED***
	for _, cap := range capability.List() ***REMOVED***
		if cap > last ***REMOVED***
			continue
		***REMOVED***
		capabilityList = append(capabilityList,
			&CapabilityMapping***REMOVED***
				Key:   "CAP_" + strings.ToUpper(cap.String()),
				Value: cap,
			***REMOVED***,
		)
	***REMOVED***
***REMOVED***

type (
	// CapabilityMapping maps linux capability name to its value of capability.Cap type
	// Capabilities is one of the security systems in Linux Security Module (LSM)
	// framework provided by the kernel.
	// For more details on capabilities, see http://man7.org/linux/man-pages/man7/capabilities.7.html
	CapabilityMapping struct ***REMOVED***
		Key   string         `json:"key,omitempty"`
		Value capability.Cap `json:"value,omitempty"`
	***REMOVED***
	// Capabilities contains all CapabilityMapping
	Capabilities []*CapabilityMapping
)

// String returns <key> of CapabilityMapping
func (c *CapabilityMapping) String() string ***REMOVED***
	return c.Key
***REMOVED***

// GetCapability returns CapabilityMapping which contains specific key
func GetCapability(key string) *CapabilityMapping ***REMOVED***
	for _, capp := range capabilityList ***REMOVED***
		if capp.Key == key ***REMOVED***
			cpy := *capp
			return &cpy
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// GetAllCapabilities returns all of the capabilities
func GetAllCapabilities() []string ***REMOVED***
	output := make([]string, len(capabilityList))
	for i, capability := range capabilityList ***REMOVED***
		output[i] = capability.String()
	***REMOVED***
	return output
***REMOVED***

// inSlice tests whether a string is contained in a slice of strings or not.
// Comparison is case insensitive
func inSlice(slice []string, s string) bool ***REMOVED***
	for _, ss := range slice ***REMOVED***
		if strings.ToLower(s) == strings.ToLower(ss) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// TweakCapabilities can tweak capabilities by adding or dropping capabilities
// based on the basics capabilities.
func TweakCapabilities(basics, adds, drops []string) ([]string, error) ***REMOVED***
	var (
		newCaps []string
		allCaps = GetAllCapabilities()
	)

	// FIXME(tonistiigi): docker format is without CAP_ prefix, oci is with prefix
	// Currently they are mixed in here. We should do conversion in one place.

	// look for invalid cap in the drop list
	for _, cap := range drops ***REMOVED***
		if strings.ToLower(cap) == "all" ***REMOVED***
			continue
		***REMOVED***

		if !inSlice(allCaps, "CAP_"+cap) ***REMOVED***
			return nil, fmt.Errorf("Unknown capability drop: %q", cap)
		***REMOVED***
	***REMOVED***

	// handle --cap-add=all
	if inSlice(adds, "all") ***REMOVED***
		basics = allCaps
	***REMOVED***

	if !inSlice(drops, "all") ***REMOVED***
		for _, cap := range basics ***REMOVED***
			// skip `all` already handled above
			if strings.ToLower(cap) == "all" ***REMOVED***
				continue
			***REMOVED***

			// if we don't drop `all`, add back all the non-dropped caps
			if !inSlice(drops, cap[4:]) ***REMOVED***
				newCaps = append(newCaps, strings.ToUpper(cap))
			***REMOVED***
		***REMOVED***
	***REMOVED***

	for _, cap := range adds ***REMOVED***
		// skip `all` already handled above
		if strings.ToLower(cap) == "all" ***REMOVED***
			continue
		***REMOVED***

		cap = "CAP_" + cap

		if !inSlice(allCaps, cap) ***REMOVED***
			return nil, fmt.Errorf("Unknown capability to add: %q", cap)
		***REMOVED***

		// add cap if not already in the list
		if !inSlice(newCaps, cap) ***REMOVED***
			newCaps = append(newCaps, strings.ToUpper(cap))
		***REMOVED***
	***REMOVED***
	return newCaps, nil
***REMOVED***
