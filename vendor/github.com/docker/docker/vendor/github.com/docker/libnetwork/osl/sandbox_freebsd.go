package osl

import "testing"

// GenerateKey generates a sandbox key based on the passed
// container id.
func GenerateKey(containerID string) string ***REMOVED***
	maxLen := 12
	if len(containerID) < maxLen ***REMOVED***
		maxLen = len(containerID)
	***REMOVED***

	return containerID[:maxLen]
***REMOVED***

// NewSandbox provides a new sandbox instance created in an os specific way
// provided a key which uniquely identifies the sandbox
func NewSandbox(key string, osCreate, isRestore bool) (Sandbox, error) ***REMOVED***
	return nil, nil
***REMOVED***

// GetSandboxForExternalKey returns sandbox object for the supplied path
func GetSandboxForExternalKey(path string, key string) (Sandbox, error) ***REMOVED***
	return nil, nil
***REMOVED***

// GC triggers garbage collection of namespace path right away
// and waits for it.
func GC() ***REMOVED***
***REMOVED***

// InitOSContext initializes OS context while configuring network resources
func InitOSContext() func() ***REMOVED***
	return func() ***REMOVED******REMOVED***
***REMOVED***

// SetupTestOSContext sets up a separate test  OS context in which tests will be executed.
func SetupTestOSContext(t *testing.T) func() ***REMOVED***
	return func() ***REMOVED******REMOVED***
***REMOVED***

// SetBasePath sets the base url prefix for the ns path
func SetBasePath(path string) ***REMOVED***
***REMOVED***
