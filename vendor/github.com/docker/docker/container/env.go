package container

import (
	"strings"
)

// ReplaceOrAppendEnvValues returns the defaults with the overrides either
// replaced by env key or appended to the list
func ReplaceOrAppendEnvValues(defaults, overrides []string) []string ***REMOVED***
	cache := make(map[string]int, len(defaults))
	for i, e := range defaults ***REMOVED***
		parts := strings.SplitN(e, "=", 2)
		cache[parts[0]] = i
	***REMOVED***

	for _, value := range overrides ***REMOVED***
		// Values w/o = means they want this env to be removed/unset.
		if !strings.Contains(value, "=") ***REMOVED***
			if i, exists := cache[value]; exists ***REMOVED***
				defaults[i] = "" // Used to indicate it should be removed
			***REMOVED***
			continue
		***REMOVED***

		// Just do a normal set/update
		parts := strings.SplitN(value, "=", 2)
		if i, exists := cache[parts[0]]; exists ***REMOVED***
			defaults[i] = value
		***REMOVED*** else ***REMOVED***
			defaults = append(defaults, value)
		***REMOVED***
	***REMOVED***

	// Now remove all entries that we want to "unset"
	for i := 0; i < len(defaults); i++ ***REMOVED***
		if defaults[i] == "" ***REMOVED***
			defaults = append(defaults[:i], defaults[i+1:]...)
			i--
		***REMOVED***
	***REMOVED***

	return defaults
***REMOVED***
