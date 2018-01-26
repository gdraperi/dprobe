// +build windows

package client

import (
	"fmt"
	"os"
	"path/filepath"
)

// LayerVhdDetails is a utility for getting a file name, size and indication of
// sandbox for a VHD(x) in a folder. A read-only layer will be layer.vhd. A
// read-write layer will be sandbox.vhdx.
func LayerVhdDetails(folder string) (string, int64, bool, error) ***REMOVED***
	var fileInfo os.FileInfo
	isSandbox := false
	filename := filepath.Join(folder, "layer.vhd")
	var err error

	if fileInfo, err = os.Stat(filename); err != nil ***REMOVED***
		filename = filepath.Join(folder, "sandbox.vhdx")
		if fileInfo, err = os.Stat(filename); err != nil ***REMOVED***
			if os.IsNotExist(err) ***REMOVED***
				return "", 0, isSandbox, fmt.Errorf("could not find layer or sandbox in %s", folder)
			***REMOVED***
			return "", 0, isSandbox, fmt.Errorf("error locating layer or sandbox in %s: %s", folder, err)
		***REMOVED***
		isSandbox = true
	***REMOVED***
	return filename, fileInfo.Size(), isSandbox, nil
***REMOVED***
