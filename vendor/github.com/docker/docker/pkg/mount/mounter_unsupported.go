// +build !linux,!freebsd freebsd,!cgo

package mount

func mount(device, target, mType string, flag uintptr, data string) error ***REMOVED***
	panic("Not implemented")
***REMOVED***

func unmount(target string, flag int) error ***REMOVED***
	panic("Not implemented")
***REMOVED***
