package fsutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

func locateDummyIfEmpty(path string) (string, error) ***REMOVED***
	children, err := ioutil.ReadDir(path)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if len(children) != 0 ***REMOVED***
		return "", nil
	***REMOVED***
	dummyFile, err := ioutil.TempFile(path, "fsutils-dummy")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	name := dummyFile.Name()
	err = dummyFile.Close()
	return name, err
***REMOVED***

// SupportsDType returns whether the filesystem mounted on path supports d_type
func SupportsDType(path string) (bool, error) ***REMOVED***
	// locate dummy so that we have at least one dirent
	dummy, err := locateDummyIfEmpty(path)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	if dummy != "" ***REMOVED***
		defer os.Remove(dummy)
	***REMOVED***

	visited := 0
	supportsDType := true
	fn := func(ent *unix.Dirent) bool ***REMOVED***
		visited++
		if ent.Type == unix.DT_UNKNOWN ***REMOVED***
			supportsDType = false
			// stop iteration
			return true
		***REMOVED***
		// continue iteration
		return false
	***REMOVED***
	if err = iterateReadDir(path, fn); err != nil ***REMOVED***
		return false, err
	***REMOVED***
	if visited == 0 ***REMOVED***
		return false, fmt.Errorf("did not hit any dirent during iteration %s", path)
	***REMOVED***
	return supportsDType, nil
***REMOVED***

func iterateReadDir(path string, fn func(*unix.Dirent) bool) error ***REMOVED***
	d, err := os.Open(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer d.Close()
	fd := int(d.Fd())
	buf := make([]byte, 4096)
	for ***REMOVED***
		nbytes, err := unix.ReadDirent(fd, buf)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if nbytes == 0 ***REMOVED***
			break
		***REMOVED***
		for off := 0; off < nbytes; ***REMOVED***
			ent := (*unix.Dirent)(unsafe.Pointer(&buf[off]))
			if stop := fn(ent); stop ***REMOVED***
				return nil
			***REMOVED***
			off += int(ent.Reclen)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
