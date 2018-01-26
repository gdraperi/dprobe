// +build !linux,!windows

package configs

func (n *Namespace) Syscall() int ***REMOVED***
	panic("No namespace syscall support")
***REMOVED***

// CloneFlags parses the container's Namespaces options to set the correct
// flags on clone, unshare. This function returns flags only for new namespaces.
func (n *Namespaces) CloneFlags() uintptr ***REMOVED***
	panic("No namespace syscall support")
***REMOVED***
