package libnetwork

import "github.com/docker/docker/pkg/reexec"

type setKeyData struct ***REMOVED***
	ContainerID string
	Key         string
***REMOVED***

func init() ***REMOVED***
	reexec.Register("libnetwork-setkey", processSetKeyReexec)
***REMOVED***
