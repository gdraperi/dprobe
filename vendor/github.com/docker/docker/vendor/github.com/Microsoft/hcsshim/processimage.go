package hcsshim

import "os"

// ProcessBaseLayer post-processes a base layer that has had its files extracted.
// The files should have been extracted to <path>\Files.
func ProcessBaseLayer(path string) error ***REMOVED***
	err := processBaseImage(path)
	if err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "ProcessBaseLayer", Path: path, Err: err***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// ProcessUtilityVMImage post-processes a utility VM image that has had its files extracted.
// The files should have been extracted to <path>\Files.
func ProcessUtilityVMImage(path string) error ***REMOVED***
	err := processUtilityImage(path)
	if err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "ProcessUtilityVMImage", Path: path, Err: err***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
