package cmd

import (
	"testing"
)

func TestFindExistingPackage(t *testing.T) ***REMOVED***
	path := findPackage("github.com/spf13/cobra")
	if path == "" ***REMOVED***
		t.Fatal("findPackage didn't find the existing package")
	***REMOVED***
	if !hasGoPathPrefix(path) ***REMOVED***
		t.Fatalf("%q is not in GOPATH, but must be", path)
	***REMOVED***
***REMOVED***

func hasGoPathPrefix(path string) bool ***REMOVED***
	for _, srcPath := range srcPaths ***REMOVED***
		if filepathHasPrefix(path, srcPath) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
