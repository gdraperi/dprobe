package volume

import "strings"

// ***REMOVED***<copy mode>=isEnabled***REMOVED***
var copyModes = map[string]bool***REMOVED***
	"nocopy": false,
***REMOVED***

func copyModeExists(mode string) bool ***REMOVED***
	_, exists := copyModes[mode]
	return exists
***REMOVED***

// GetCopyMode gets the copy mode from the mode string for mounts
func getCopyMode(mode string, def bool) (bool, bool) ***REMOVED***
	for _, o := range strings.Split(mode, ",") ***REMOVED***
		if isEnabled, exists := copyModes[o]; exists ***REMOVED***
			return isEnabled, true
		***REMOVED***
	***REMOVED***
	return def, false
***REMOVED***
