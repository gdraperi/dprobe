// +build !windows,!linux,!freebsd freebsd,!cgo

package mount

import (
	"fmt"
	"runtime"
)

func parseMountTable() ([]*Info, error) ***REMOVED***
	return nil, fmt.Errorf("mount.parseMountTable is not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
***REMOVED***
