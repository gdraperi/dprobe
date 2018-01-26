package daemon

import "github.com/opencontainers/selinux/go-selinux"

func selinuxSetDisabled() ***REMOVED***
	selinux.SetDisabled()
***REMOVED***

func selinuxFreeLxcContexts(label string) ***REMOVED***
	selinux.ReleaseLabel(label)
***REMOVED***

func selinuxEnabled() bool ***REMOVED***
	return selinux.GetEnabled()
***REMOVED***
