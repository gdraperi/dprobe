// +build linux,!btrfs_noversion

package btrfs

import (
	"testing"
)

func TestLibVersion(t *testing.T) ***REMOVED***
	if btrfsLibVersion() <= 0 ***REMOVED***
		t.Error("expected output from btrfs lib version > 0")
	***REMOVED***
***REMOVED***
