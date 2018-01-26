// +build !windows

package ioutils

import "io/ioutil"

// TempDir on Unix systems is equivalent to ioutil.TempDir.
func TempDir(dir, prefix string) (string, error) ***REMOVED***
	return ioutil.TempDir(dir, prefix)
***REMOVED***
