package fs

import (
	"os"

	"golang.org/x/sys/windows"
)

func detectDirDiff(upper, lower string) *diffDirOptions ***REMOVED***
	return nil
***REMOVED***

func compareSysStat(s1, s2 interface***REMOVED******REMOVED***) (bool, error) ***REMOVED***
	f1, ok := s1.(windows.Win32FileAttributeData)
	if !ok ***REMOVED***
		return false, nil
	***REMOVED***
	f2, ok := s2.(windows.Win32FileAttributeData)
	if !ok ***REMOVED***
		return false, nil
	***REMOVED***
	return f1.FileAttributes == f2.FileAttributes, nil
***REMOVED***

func compareCapabilities(p1, p2 string) (bool, error) ***REMOVED***
	// TODO: Use windows equivalent
	return true, nil
***REMOVED***

func isLinked(os.FileInfo) bool ***REMOVED***
	return false
***REMOVED***
