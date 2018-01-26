// +build windows

package system

import (
	"testing"

	"github.com/containerd/continuity/pathdriver"
)

// TestCheckSystemDriveAndRemoveDriveLetter tests CheckSystemDriveAndRemoveDriveLetter
func TestCheckSystemDriveAndRemoveDriveLetter(t *testing.T) ***REMOVED***
	// Fails if not C drive.
	_, err := CheckSystemDriveAndRemoveDriveLetter(`d:\`, pathdriver.LocalPathDriver)
	if err == nil || (err != nil && err.Error() != "The specified path is not on the system drive (C:)") ***REMOVED***
		t.Fatalf("Expected error for d:")
	***REMOVED***

	// Single character is unchanged
	var path string
	if path, err = CheckSystemDriveAndRemoveDriveLetter("z", pathdriver.LocalPathDriver); err != nil ***REMOVED***
		t.Fatalf("Single character should pass")
	***REMOVED***
	if path != "z" ***REMOVED***
		t.Fatalf("Single character should be unchanged")
	***REMOVED***

	// Two characters without colon is unchanged
	if path, err = CheckSystemDriveAndRemoveDriveLetter("AB", pathdriver.LocalPathDriver); err != nil ***REMOVED***
		t.Fatalf("2 characters without colon should pass")
	***REMOVED***
	if path != "AB" ***REMOVED***
		t.Fatalf("2 characters without colon should be unchanged")
	***REMOVED***

	// Abs path without drive letter
	if path, err = CheckSystemDriveAndRemoveDriveLetter(`\l`, pathdriver.LocalPathDriver); err != nil ***REMOVED***
		t.Fatalf("abs path no drive letter should pass")
	***REMOVED***
	if path != `\l` ***REMOVED***
		t.Fatalf("abs path without drive letter should be unchanged")
	***REMOVED***

	// Abs path without drive letter, linux style
	if path, err = CheckSystemDriveAndRemoveDriveLetter(`/l`, pathdriver.LocalPathDriver); err != nil ***REMOVED***
		t.Fatalf("abs path no drive letter linux style should pass")
	***REMOVED***
	if path != `\l` ***REMOVED***
		t.Fatalf("abs path without drive letter linux failed %s", path)
	***REMOVED***

	// Drive-colon should be stripped
	if path, err = CheckSystemDriveAndRemoveDriveLetter(`c:\`, pathdriver.LocalPathDriver); err != nil ***REMOVED***
		t.Fatalf("An absolute path should pass")
	***REMOVED***
	if path != `\` ***REMOVED***
		t.Fatalf(`An absolute path should have been shortened to \ %s`, path)
	***REMOVED***

	// Verify with a linux-style path
	if path, err = CheckSystemDriveAndRemoveDriveLetter(`c:/`, pathdriver.LocalPathDriver); err != nil ***REMOVED***
		t.Fatalf("An absolute path should pass")
	***REMOVED***
	if path != `\` ***REMOVED***
		t.Fatalf(`A linux style absolute path should have been shortened to \ %s`, path)
	***REMOVED***

	// Failure on c:
	if path, err = CheckSystemDriveAndRemoveDriveLetter(`c:`, pathdriver.LocalPathDriver); err == nil ***REMOVED***
		t.Fatalf("c: should fail")
	***REMOVED***
	if err.Error() != `No relative path specified in "c:"` ***REMOVED***
		t.Fatalf(path, err)
	***REMOVED***

	// Failure on d:
	if path, err = CheckSystemDriveAndRemoveDriveLetter(`d:`, pathdriver.LocalPathDriver); err == nil ***REMOVED***
		t.Fatalf("c: should fail")
	***REMOVED***
	if err.Error() != `No relative path specified in "d:"` ***REMOVED***
		t.Fatalf(path, err)
	***REMOVED***
***REMOVED***
