// +build linux,cgo

package loopback

import (
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

// Loopback related errors
var (
	ErrAttachLoopbackDevice   = errors.New("loopback attach failed")
	ErrGetLoopbackBackingFile = errors.New("Unable to get loopback backing file")
	ErrSetCapacity            = errors.New("Unable set loopback capacity")
)

func stringToLoopName(src string) [LoNameSize]uint8 ***REMOVED***
	var dst [LoNameSize]uint8
	copy(dst[:], src[:])
	return dst
***REMOVED***

func getNextFreeLoopbackIndex() (int, error) ***REMOVED***
	f, err := os.OpenFile("/dev/loop-control", os.O_RDONLY, 0644)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	defer f.Close()

	index, err := ioctlLoopCtlGetFree(f.Fd())
	if index < 0 ***REMOVED***
		index = 0
	***REMOVED***
	return index, err
***REMOVED***

func openNextAvailableLoopback(index int, sparseFile *os.File) (loopFile *os.File, err error) ***REMOVED***
	// Start looking for a free /dev/loop
	for ***REMOVED***
		target := fmt.Sprintf("/dev/loop%d", index)
		index++

		fi, err := os.Stat(target)
		if err != nil ***REMOVED***
			if os.IsNotExist(err) ***REMOVED***
				logrus.Error("There are no more loopback devices available.")
			***REMOVED***
			return nil, ErrAttachLoopbackDevice
		***REMOVED***

		if fi.Mode()&os.ModeDevice != os.ModeDevice ***REMOVED***
			logrus.Errorf("Loopback device %s is not a block device.", target)
			continue
		***REMOVED***

		// OpenFile adds O_CLOEXEC
		loopFile, err = os.OpenFile(target, os.O_RDWR, 0644)
		if err != nil ***REMOVED***
			logrus.Errorf("Error opening loopback device: %s", err)
			return nil, ErrAttachLoopbackDevice
		***REMOVED***

		// Try to attach to the loop file
		if err := ioctlLoopSetFd(loopFile.Fd(), sparseFile.Fd()); err != nil ***REMOVED***
			loopFile.Close()

			// If the error is EBUSY, then try the next loopback
			if err != unix.EBUSY ***REMOVED***
				logrus.Errorf("Cannot set up loopback device %s: %s", target, err)
				return nil, ErrAttachLoopbackDevice
			***REMOVED***

			// Otherwise, we keep going with the loop
			continue
		***REMOVED***
		// In case of success, we finished. Break the loop.
		break
	***REMOVED***

	// This can't happen, but let's be sure
	if loopFile == nil ***REMOVED***
		logrus.Errorf("Unreachable code reached! Error attaching %s to a loopback device.", sparseFile.Name())
		return nil, ErrAttachLoopbackDevice
	***REMOVED***

	return loopFile, nil
***REMOVED***

// AttachLoopDevice attaches the given sparse file to the next
// available loopback device. It returns an opened *os.File.
func AttachLoopDevice(sparseName string) (loop *os.File, err error) ***REMOVED***

	// Try to retrieve the next available loopback device via syscall.
	// If it fails, we discard error and start looping for a
	// loopback from index 0.
	startIndex, err := getNextFreeLoopbackIndex()
	if err != nil ***REMOVED***
		logrus.Debugf("Error retrieving the next available loopback: %s", err)
	***REMOVED***

	// OpenFile adds O_CLOEXEC
	sparseFile, err := os.OpenFile(sparseName, os.O_RDWR, 0644)
	if err != nil ***REMOVED***
		logrus.Errorf("Error opening sparse file %s: %s", sparseName, err)
		return nil, ErrAttachLoopbackDevice
	***REMOVED***
	defer sparseFile.Close()

	loopFile, err := openNextAvailableLoopback(startIndex, sparseFile)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Set the status of the loopback device
	loopInfo := &loopInfo64***REMOVED***
		loFileName: stringToLoopName(loopFile.Name()),
		loOffset:   0,
		loFlags:    LoFlagsAutoClear,
	***REMOVED***

	if err := ioctlLoopSetStatus64(loopFile.Fd(), loopInfo); err != nil ***REMOVED***
		logrus.Errorf("Cannot set up loopback device info: %s", err)

		// If the call failed, then free the loopback device
		if err := ioctlLoopClrFd(loopFile.Fd()); err != nil ***REMOVED***
			logrus.Error("Error while cleaning up the loopback device")
		***REMOVED***
		loopFile.Close()
		return nil, ErrAttachLoopbackDevice
	***REMOVED***

	return loopFile, nil
***REMOVED***
