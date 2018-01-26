// +build !linux,!windows

package system

// ReadMemInfo is not supported on platforms other than linux and windows.
func ReadMemInfo() (*MemInfo, error) ***REMOVED***
	return nil, ErrNotSupportedPlatform
***REMOVED***
