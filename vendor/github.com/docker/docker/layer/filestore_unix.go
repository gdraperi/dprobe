// +build !windows

package layer

import "runtime"

// setOS writes the "os" file to the layer filestore
func (fm *fileMetadataTransaction) setOS(os string) error ***REMOVED***
	return nil
***REMOVED***

// getOS reads the "os" file from the layer filestore
func (fms *fileMetadataStore) getOS(layer ChainID) (string, error) ***REMOVED***
	return runtime.GOOS, nil
***REMOVED***
