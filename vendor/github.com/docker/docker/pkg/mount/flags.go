package mount

import (
	"fmt"
	"strings"
)

var flags = map[string]struct ***REMOVED***
	clear bool
	flag  int
***REMOVED******REMOVED***
	"defaults":      ***REMOVED***false, 0***REMOVED***,
	"ro":            ***REMOVED***false, RDONLY***REMOVED***,
	"rw":            ***REMOVED***true, RDONLY***REMOVED***,
	"suid":          ***REMOVED***true, NOSUID***REMOVED***,
	"nosuid":        ***REMOVED***false, NOSUID***REMOVED***,
	"dev":           ***REMOVED***true, NODEV***REMOVED***,
	"nodev":         ***REMOVED***false, NODEV***REMOVED***,
	"exec":          ***REMOVED***true, NOEXEC***REMOVED***,
	"noexec":        ***REMOVED***false, NOEXEC***REMOVED***,
	"sync":          ***REMOVED***false, SYNCHRONOUS***REMOVED***,
	"async":         ***REMOVED***true, SYNCHRONOUS***REMOVED***,
	"dirsync":       ***REMOVED***false, DIRSYNC***REMOVED***,
	"remount":       ***REMOVED***false, REMOUNT***REMOVED***,
	"mand":          ***REMOVED***false, MANDLOCK***REMOVED***,
	"nomand":        ***REMOVED***true, MANDLOCK***REMOVED***,
	"atime":         ***REMOVED***true, NOATIME***REMOVED***,
	"noatime":       ***REMOVED***false, NOATIME***REMOVED***,
	"diratime":      ***REMOVED***true, NODIRATIME***REMOVED***,
	"nodiratime":    ***REMOVED***false, NODIRATIME***REMOVED***,
	"bind":          ***REMOVED***false, BIND***REMOVED***,
	"rbind":         ***REMOVED***false, RBIND***REMOVED***,
	"unbindable":    ***REMOVED***false, UNBINDABLE***REMOVED***,
	"runbindable":   ***REMOVED***false, RUNBINDABLE***REMOVED***,
	"private":       ***REMOVED***false, PRIVATE***REMOVED***,
	"rprivate":      ***REMOVED***false, RPRIVATE***REMOVED***,
	"shared":        ***REMOVED***false, SHARED***REMOVED***,
	"rshared":       ***REMOVED***false, RSHARED***REMOVED***,
	"slave":         ***REMOVED***false, SLAVE***REMOVED***,
	"rslave":        ***REMOVED***false, RSLAVE***REMOVED***,
	"relatime":      ***REMOVED***false, RELATIME***REMOVED***,
	"norelatime":    ***REMOVED***true, RELATIME***REMOVED***,
	"strictatime":   ***REMOVED***false, STRICTATIME***REMOVED***,
	"nostrictatime": ***REMOVED***true, STRICTATIME***REMOVED***,
***REMOVED***

var validFlags = map[string]bool***REMOVED***
	"":          true,
	"size":      true,
	"mode":      true,
	"uid":       true,
	"gid":       true,
	"nr_inodes": true,
	"nr_blocks": true,
	"mpol":      true,
***REMOVED***

var propagationFlags = map[string]bool***REMOVED***
	"bind":        true,
	"rbind":       true,
	"unbindable":  true,
	"runbindable": true,
	"private":     true,
	"rprivate":    true,
	"shared":      true,
	"rshared":     true,
	"slave":       true,
	"rslave":      true,
***REMOVED***

// MergeTmpfsOptions merge mount options to make sure there is no duplicate.
func MergeTmpfsOptions(options []string) ([]string, error) ***REMOVED***
	// We use collisions maps to remove duplicates.
	// For flag, the key is the flag value (the key for propagation flag is -1)
	// For data=value, the key is the data
	flagCollisions := map[int]bool***REMOVED******REMOVED***
	dataCollisions := map[string]bool***REMOVED******REMOVED***

	var newOptions []string
	// We process in reverse order
	for i := len(options) - 1; i >= 0; i-- ***REMOVED***
		option := options[i]
		if option == "defaults" ***REMOVED***
			continue
		***REMOVED***
		if f, ok := flags[option]; ok && f.flag != 0 ***REMOVED***
			// There is only one propagation mode
			key := f.flag
			if propagationFlags[option] ***REMOVED***
				key = -1
			***REMOVED***
			// Check to see if there is collision for flag
			if !flagCollisions[key] ***REMOVED***
				// We prepend the option and add to collision map
				newOptions = append([]string***REMOVED***option***REMOVED***, newOptions...)
				flagCollisions[key] = true
			***REMOVED***
			continue
		***REMOVED***
		opt := strings.SplitN(option, "=", 2)
		if len(opt) != 2 || !validFlags[opt[0]] ***REMOVED***
			return nil, fmt.Errorf("Invalid tmpfs option %q", opt)
		***REMOVED***
		if !dataCollisions[opt[0]] ***REMOVED***
			// We prepend the option and add to collision map
			newOptions = append([]string***REMOVED***option***REMOVED***, newOptions...)
			dataCollisions[opt[0]] = true
		***REMOVED***
	***REMOVED***

	return newOptions, nil
***REMOVED***

// Parse fstab type mount options into mount() flags
// and device specific data
func parseOptions(options string) (int, string) ***REMOVED***
	var (
		flag int
		data []string
	)

	for _, o := range strings.Split(options, ",") ***REMOVED***
		// If the option does not exist in the flags table or the flag
		// is not supported on the platform,
		// then it is a data value for a specific fs type
		if f, exists := flags[o]; exists && f.flag != 0 ***REMOVED***
			if f.clear ***REMOVED***
				flag &= ^f.flag
			***REMOVED*** else ***REMOVED***
				flag |= f.flag
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			data = append(data, o)
		***REMOVED***
	***REMOVED***
	return flag, strings.Join(data, ",")
***REMOVED***

// ParseTmpfsOptions parse fstab type mount options into flags and data
func ParseTmpfsOptions(options string) (int, string, error) ***REMOVED***
	flags, data := parseOptions(options)
	for _, o := range strings.Split(data, ",") ***REMOVED***
		opt := strings.SplitN(o, "=", 2)
		if !validFlags[opt[0]] ***REMOVED***
			return 0, "", fmt.Errorf("Invalid tmpfs option %q", opt)
		***REMOVED***
	***REMOVED***
	return flags, data, nil
***REMOVED***
