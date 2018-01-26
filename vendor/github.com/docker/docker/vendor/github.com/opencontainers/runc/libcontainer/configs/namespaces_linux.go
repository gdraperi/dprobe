package configs

import (
	"fmt"
	"os"
	"sync"
)

const (
	NEWNET  NamespaceType = "NEWNET"
	NEWPID  NamespaceType = "NEWPID"
	NEWNS   NamespaceType = "NEWNS"
	NEWUTS  NamespaceType = "NEWUTS"
	NEWIPC  NamespaceType = "NEWIPC"
	NEWUSER NamespaceType = "NEWUSER"
)

var (
	nsLock              sync.Mutex
	supportedNamespaces = make(map[NamespaceType]bool)
)

// NsName converts the namespace type to its filename
func NsName(ns NamespaceType) string ***REMOVED***
	switch ns ***REMOVED***
	case NEWNET:
		return "net"
	case NEWNS:
		return "mnt"
	case NEWPID:
		return "pid"
	case NEWIPC:
		return "ipc"
	case NEWUSER:
		return "user"
	case NEWUTS:
		return "uts"
	***REMOVED***
	return ""
***REMOVED***

// IsNamespaceSupported returns whether a namespace is available or
// not
func IsNamespaceSupported(ns NamespaceType) bool ***REMOVED***
	nsLock.Lock()
	defer nsLock.Unlock()
	supported, ok := supportedNamespaces[ns]
	if ok ***REMOVED***
		return supported
	***REMOVED***
	nsFile := NsName(ns)
	// if the namespace type is unknown, just return false
	if nsFile == "" ***REMOVED***
		return false
	***REMOVED***
	_, err := os.Stat(fmt.Sprintf("/proc/self/ns/%s", nsFile))
	// a namespace is supported if it exists and we have permissions to read it
	supported = err == nil
	supportedNamespaces[ns] = supported
	return supported
***REMOVED***

func NamespaceTypes() []NamespaceType ***REMOVED***
	return []NamespaceType***REMOVED***
		NEWUSER, // Keep user NS always first, don't move it.
		NEWIPC,
		NEWUTS,
		NEWNET,
		NEWPID,
		NEWNS,
	***REMOVED***
***REMOVED***

// Namespace defines configuration for each namespace.  It specifies an
// alternate path that is able to be joined via setns.
type Namespace struct ***REMOVED***
	Type NamespaceType `json:"type"`
	Path string        `json:"path"`
***REMOVED***

func (n *Namespace) GetPath(pid int) string ***REMOVED***
	return fmt.Sprintf("/proc/%d/ns/%s", pid, NsName(n.Type))
***REMOVED***

func (n *Namespaces) Remove(t NamespaceType) bool ***REMOVED***
	i := n.index(t)
	if i == -1 ***REMOVED***
		return false
	***REMOVED***
	*n = append((*n)[:i], (*n)[i+1:]...)
	return true
***REMOVED***

func (n *Namespaces) Add(t NamespaceType, path string) ***REMOVED***
	i := n.index(t)
	if i == -1 ***REMOVED***
		*n = append(*n, Namespace***REMOVED***Type: t, Path: path***REMOVED***)
		return
	***REMOVED***
	(*n)[i].Path = path
***REMOVED***

func (n *Namespaces) index(t NamespaceType) int ***REMOVED***
	for i, ns := range *n ***REMOVED***
		if ns.Type == t ***REMOVED***
			return i
		***REMOVED***
	***REMOVED***
	return -1
***REMOVED***

func (n *Namespaces) Contains(t NamespaceType) bool ***REMOVED***
	return n.index(t) != -1
***REMOVED***

func (n *Namespaces) PathOf(t NamespaceType) string ***REMOVED***
	i := n.index(t)
	if i == -1 ***REMOVED***
		return ""
	***REMOVED***
	return (*n)[i].Path
***REMOVED***
