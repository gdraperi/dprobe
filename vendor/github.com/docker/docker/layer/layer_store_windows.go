package layer

import (
	"io"

	"github.com/docker/distribution"
)

func (ls *layerStore) RegisterWithDescriptor(ts io.Reader, parent ChainID, descriptor distribution.Descriptor) (Layer, error) ***REMOVED***
	return ls.registerWithDescriptor(ts, parent, descriptor)
***REMOVED***
