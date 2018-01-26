package system

import (
	"os"
	"syscall"
	"time"

	"github.com/docker/docker/pkg/mount"
	"github.com/pkg/errors"
)

// EnsureRemoveAll wraps `os.RemoveAll` to check for specific errors that can
// often be remedied.
// Only use `EnsureRemoveAll` if you really want to make every effort to remove
// a directory.
//
// Because of the way `os.Remove` (and by extension `os.RemoveAll`) works, there
// can be a race between reading directory entries and then actually attempting
// to remove everything in the directory.
// These types of errors do not need to be returned since it's ok for the dir to
// be gone we can just retry the remove operation.
//
// This should not return a `os.ErrNotExist` kind of error under any circumstances
func EnsureRemoveAll(dir string) error ***REMOVED***
	notExistErr := make(map[string]bool)

	// track retries
	exitOnErr := make(map[string]int)
	maxRetry := 50

	// Attempt to unmount anything beneath this dir first
	mount.RecursiveUnmount(dir)

	for ***REMOVED***
		err := os.RemoveAll(dir)
		if err == nil ***REMOVED***
			return err
		***REMOVED***

		pe, ok := err.(*os.PathError)
		if !ok ***REMOVED***
			return err
		***REMOVED***

		if os.IsNotExist(err) ***REMOVED***
			if notExistErr[pe.Path] ***REMOVED***
				return err
			***REMOVED***
			notExistErr[pe.Path] = true

			// There is a race where some subdir can be removed but after the parent
			//   dir entries have been read.
			// So the path could be from `os.Remove(subdir)`
			// If the reported non-existent path is not the passed in `dir` we
			// should just retry, but otherwise return with no error.
			if pe.Path == dir ***REMOVED***
				return nil
			***REMOVED***
			continue
		***REMOVED***

		if pe.Err != syscall.EBUSY ***REMOVED***
			return err
		***REMOVED***

		if mounted, _ := mount.Mounted(pe.Path); mounted ***REMOVED***
			if e := mount.Unmount(pe.Path); e != nil ***REMOVED***
				if mounted, _ := mount.Mounted(pe.Path); mounted ***REMOVED***
					return errors.Wrapf(e, "error while removing %s", dir)
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if exitOnErr[pe.Path] == maxRetry ***REMOVED***
			return err
		***REMOVED***
		exitOnErr[pe.Path]++
		time.Sleep(100 * time.Millisecond)
	***REMOVED***
***REMOVED***
