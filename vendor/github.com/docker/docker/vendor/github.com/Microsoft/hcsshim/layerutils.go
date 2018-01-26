package hcsshim

// This file contains utility functions to support storage (graph) related
// functionality.

import (
	"path/filepath"
	"syscall"

	"github.com/sirupsen/logrus"
)

/* To pass into syscall, we need a struct matching the following:
enum GraphDriverType
***REMOVED***
    DiffDriver,
    FilterDriver
***REMOVED***;

struct DriverInfo ***REMOVED***
    GraphDriverType Flavour;
    LPCWSTR HomeDir;
***REMOVED***;
*/
type DriverInfo struct ***REMOVED***
	Flavour int
	HomeDir string
***REMOVED***

type driverInfo struct ***REMOVED***
	Flavour  int
	HomeDirp *uint16
***REMOVED***

func convertDriverInfo(info DriverInfo) (driverInfo, error) ***REMOVED***
	homedirp, err := syscall.UTF16PtrFromString(info.HomeDir)
	if err != nil ***REMOVED***
		logrus.Debugf("Failed conversion of home to pointer for driver info: %s", err.Error())
		return driverInfo***REMOVED******REMOVED***, err
	***REMOVED***

	return driverInfo***REMOVED***
		Flavour:  info.Flavour,
		HomeDirp: homedirp,
	***REMOVED***, nil
***REMOVED***

/* To pass into syscall, we need a struct matching the following:
typedef struct _WC_LAYER_DESCRIPTOR ***REMOVED***

    //
    // The ID of the layer
    //

    GUID LayerId;

    //
    // Additional flags
    //

    union ***REMOVED***
        struct ***REMOVED***
            ULONG Reserved : 31;
            ULONG Dirty : 1;    // Created from sandbox as a result of snapshot
    ***REMOVED***;
        ULONG Value;
***REMOVED*** Flags;

    //
    // Path to the layer root directory, null-terminated
    //

    PCWSTR Path;

***REMOVED*** WC_LAYER_DESCRIPTOR, *PWC_LAYER_DESCRIPTOR;
*/
type WC_LAYER_DESCRIPTOR struct ***REMOVED***
	LayerId GUID
	Flags   uint32
	Pathp   *uint16
***REMOVED***

func layerPathsToDescriptors(parentLayerPaths []string) ([]WC_LAYER_DESCRIPTOR, error) ***REMOVED***
	// Array of descriptors that gets constructed.
	var layers []WC_LAYER_DESCRIPTOR

	for i := 0; i < len(parentLayerPaths); i++ ***REMOVED***
		// Create a layer descriptor, using the folder name
		// as the source for a GUID LayerId
		_, folderName := filepath.Split(parentLayerPaths[i])
		g, err := NameToGuid(folderName)
		if err != nil ***REMOVED***
			logrus.Debugf("Failed to convert name to guid %s", err)
			return nil, err
		***REMOVED***

		p, err := syscall.UTF16PtrFromString(parentLayerPaths[i])
		if err != nil ***REMOVED***
			logrus.Debugf("Failed conversion of parentLayerPath to pointer %s", err)
			return nil, err
		***REMOVED***

		layers = append(layers, WC_LAYER_DESCRIPTOR***REMOVED***
			LayerId: g,
			Flags:   0,
			Pathp:   p,
		***REMOVED***)
	***REMOVED***

	return layers, nil
***REMOVED***
