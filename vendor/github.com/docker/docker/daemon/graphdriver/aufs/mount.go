// +build linux

package aufs

import (
	"os/exec"

	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

// Unmount the target specified.
func Unmount(target string) error ***REMOVED***
	if err := exec.Command("auplink", target, "flush").Run(); err != nil ***REMOVED***
		logrus.Warnf("Couldn't run auplink before unmount %s: %s", target, err)
	***REMOVED***
	return unix.Unmount(target, 0)
***REMOVED***
