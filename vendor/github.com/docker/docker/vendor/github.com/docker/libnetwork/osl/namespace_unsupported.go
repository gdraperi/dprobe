// +build !linux,!windows,!freebsd

package osl

// GC triggers garbage collection of namespace path right away
// and waits for it.
func GC() ***REMOVED***
***REMOVED***

// GetSandboxForExternalKey returns sandbox object for the supplied path
func GetSandboxForExternalKey(path string, key string) (Sandbox, error) ***REMOVED***
	return nil, nil
***REMOVED***

// SetBasePath sets the base url prefix for the ns path
func SetBasePath(path string) ***REMOVED***
***REMOVED***
