// +build linux

package configs

import "golang.org/x/sys/unix"

func (n *Namespace) Syscall() int ***REMOVED***
	return namespaceInfo[n.Type]
***REMOVED***

var namespaceInfo = map[NamespaceType]int***REMOVED***
	NEWNET:  unix.CLONE_NEWNET,
	NEWNS:   unix.CLONE_NEWNS,
	NEWUSER: unix.CLONE_NEWUSER,
	NEWIPC:  unix.CLONE_NEWIPC,
	NEWUTS:  unix.CLONE_NEWUTS,
	NEWPID:  unix.CLONE_NEWPID,
***REMOVED***

// CloneFlags parses the container's Namespaces options to set the correct
// flags on clone, unshare. This function returns flags only for new namespaces.
func (n *Namespaces) CloneFlags() uintptr ***REMOVED***
	var flag int
	for _, v := range *n ***REMOVED***
		if v.Path != "" ***REMOVED***
			continue
		***REMOVED***
		flag |= namespaceInfo[v.Type]
	***REMOVED***
	return uintptr(flag)
***REMOVED***
