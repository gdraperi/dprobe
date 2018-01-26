// +build linux freebsd

package system

import (
	"os"
	"testing"
)

// TestLstat tests Lstat for existing and non existing files
func TestLstat(t *testing.T) ***REMOVED***
	file, invalid, _, dir := prepareFiles(t)
	defer os.RemoveAll(dir)

	statFile, err := Lstat(file)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if statFile == nil ***REMOVED***
		t.Fatal("returned empty stat for existing file")
	***REMOVED***

	statInvalid, err := Lstat(invalid)
	if err == nil ***REMOVED***
		t.Fatal("did not return error for non-existing file")
	***REMOVED***
	if statInvalid != nil ***REMOVED***
		t.Fatal("returned non-nil stat for non-existing file")
	***REMOVED***
***REMOVED***
