// +build linux,cgo

package loopback

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

func getLoopbackBackingFile(file *os.File) (uint64, uint64, error) ***REMOVED***
	loopInfo, err := ioctlLoopGetStatus64(file.Fd())
	if err != nil ***REMOVED***
		logrus.Errorf("Error get loopback backing file: %s", err)
		return 0, 0, ErrGetLoopbackBackingFile
	***REMOVED***
	return loopInfo.loDevice, loopInfo.loInode, nil
***REMOVED***

// SetCapacity reloads the size for the loopback device.
func SetCapacity(file *os.File) error ***REMOVED***
	if err := ioctlLoopSetCapacity(file.Fd(), 0); err != nil ***REMOVED***
		logrus.Errorf("Error loopbackSetCapacity: %s", err)
		return ErrSetCapacity
	***REMOVED***
	return nil
***REMOVED***

// FindLoopDeviceFor returns a loopback device file for the specified file which
// is backing file of a loop back device.
func FindLoopDeviceFor(file *os.File) *os.File ***REMOVED***
	var stat unix.Stat_t
	err := unix.Stat(file.Name(), &stat)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	targetInode := stat.Ino
	targetDevice := stat.Dev

	for i := 0; true; i++ ***REMOVED***
		path := fmt.Sprintf("/dev/loop%d", i)

		file, err := os.OpenFile(path, os.O_RDWR, 0)
		if err != nil ***REMOVED***
			if os.IsNotExist(err) ***REMOVED***
				return nil
			***REMOVED***

			// Ignore all errors until the first not-exist
			// we want to continue looking for the file
			continue
		***REMOVED***

		dev, inode, err := getLoopbackBackingFile(file)
		if err == nil && dev == targetDevice && inode == targetInode ***REMOVED***
			return file
		***REMOVED***
		file.Close()
	***REMOVED***

	return nil
***REMOVED***
