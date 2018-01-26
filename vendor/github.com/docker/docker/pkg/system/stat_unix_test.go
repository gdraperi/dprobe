// +build linux freebsd

package system

import (
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestFromStatT tests fromStatT for a tempfile
func TestFromStatT(t *testing.T) ***REMOVED***
	file, _, _, dir := prepareFiles(t)
	defer os.RemoveAll(dir)

	stat := &syscall.Stat_t***REMOVED******REMOVED***
	err := syscall.Lstat(file, stat)
	require.NoError(t, err)

	s, err := fromStatT(stat)
	require.NoError(t, err)

	if stat.Mode != s.Mode() ***REMOVED***
		t.Fatal("got invalid mode")
	***REMOVED***
	if stat.Uid != s.UID() ***REMOVED***
		t.Fatal("got invalid uid")
	***REMOVED***
	if stat.Gid != s.GID() ***REMOVED***
		t.Fatal("got invalid gid")
	***REMOVED***
	if stat.Rdev != s.Rdev() ***REMOVED***
		t.Fatal("got invalid rdev")
	***REMOVED***
	if stat.Mtim != s.Mtim() ***REMOVED***
		t.Fatal("got invalid mtim")
	***REMOVED***
***REMOVED***
