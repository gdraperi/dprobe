package layer

import "github.com/docker/distribution"

var _ distribution.Describable = &roLayer***REMOVED******REMOVED***

func (rl *roLayer) Descriptor() distribution.Descriptor ***REMOVED***
	return rl.descriptor
***REMOVED***
