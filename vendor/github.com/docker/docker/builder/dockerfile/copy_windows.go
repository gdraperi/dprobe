package dockerfile

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/idtools"
)

var pathBlacklist = map[string]bool***REMOVED***
	"c:\\":        true,
	"c:\\windows": true,
***REMOVED***

func fixPermissions(source, destination string, rootIDs idtools.IDPair, overrideSkip bool) error ***REMOVED***
	// chown is not supported on Windows
	return nil
***REMOVED***

func validateCopySourcePath(imageSource *imageMount, origPath, platform string) error ***REMOVED***
	// validate windows paths from other images + LCOW
	if imageSource == nil || platform != "windows" ***REMOVED***
		return nil
	***REMOVED***

	origPath = filepath.FromSlash(origPath)
	p := strings.ToLower(filepath.Clean(origPath))
	if !filepath.IsAbs(p) ***REMOVED***
		if filepath.VolumeName(p) != "" ***REMOVED***
			if p[len(p)-2:] == ":." ***REMOVED*** // case where clean returns weird c:. paths
				p = p[:len(p)-1]
			***REMOVED***
			p += "\\"
		***REMOVED*** else ***REMOVED***
			p = filepath.Join("c:\\", p)
		***REMOVED***
	***REMOVED***
	if _, blacklisted := pathBlacklist[p]; blacklisted ***REMOVED***
		return errors.New("copy from c:\\ or c:\\windows is not allowed on windows")
	***REMOVED***
	return nil
***REMOVED***
