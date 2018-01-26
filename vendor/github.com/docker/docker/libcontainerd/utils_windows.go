package libcontainerd

import (
	"strings"

	"syscall"

	opengcs "github.com/Microsoft/opengcs/client"
)

// setupEnvironmentVariables converts a string array of environment variables
// into a map as required by the HCS. Source array is in format [v1=k1] [v2=k2] etc.
func setupEnvironmentVariables(a []string) map[string]string ***REMOVED***
	r := make(map[string]string)
	for _, s := range a ***REMOVED***
		arr := strings.SplitN(s, "=", 2)
		if len(arr) == 2 ***REMOVED***
			r[arr[0]] = arr[1]
		***REMOVED***
	***REMOVED***
	return r
***REMOVED***

// Apply for the LCOW option is a no-op.
func (s *LCOWOption) Apply(interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

// debugGCS is a dirty hack for debugging for Linux Utility VMs. It simply
// runs a bunch of commands inside the UVM, but seriously aides in advanced debugging.
func (c *container) debugGCS() ***REMOVED***
	if c == nil || c.isWindows || c.hcsContainer == nil ***REMOVED***
		return
	***REMOVED***
	cfg := opengcs.Config***REMOVED***
		Uvm:               c.hcsContainer,
		UvmTimeoutSeconds: 600,
	***REMOVED***
	cfg.DebugGCS()
***REMOVED***

// containerdSysProcAttr returns the SysProcAttr to use when exec'ing
// containerd
func containerdSysProcAttr() *syscall.SysProcAttr ***REMOVED***
	return nil
***REMOVED***
