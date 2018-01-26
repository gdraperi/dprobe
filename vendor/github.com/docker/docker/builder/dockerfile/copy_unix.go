// +build !windows

package dockerfile

import (
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/idtools"
)

func fixPermissions(source, destination string, rootIDs idtools.IDPair, overrideSkip bool) error ***REMOVED***
	var (
		skipChownRoot bool
		err           error
	)
	if !overrideSkip ***REMOVED***
		destEndpoint := &copyEndpoint***REMOVED***driver: containerfs.NewLocalDriver(), path: destination***REMOVED***
		skipChownRoot, err = isExistingDirectory(destEndpoint)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// We Walk on the source rather than on the destination because we don't
	// want to change permissions on things we haven't created or modified.
	return filepath.Walk(source, func(fullpath string, info os.FileInfo, err error) error ***REMOVED***
		// Do not alter the walk root iff. it existed before, as it doesn't fall under
		// the domain of "things we should chown".
		if skipChownRoot && source == fullpath ***REMOVED***
			return nil
		***REMOVED***

		// Path is prefixed by source: substitute with destination instead.
		cleaned, err := filepath.Rel(source, fullpath)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		fullpath = filepath.Join(destination, cleaned)
		return os.Lchown(fullpath, rootIDs.UID, rootIDs.GID)
	***REMOVED***)
***REMOVED***

func validateCopySourcePath(imageSource *imageMount, origPath, platform string) error ***REMOVED***
	return nil
***REMOVED***
