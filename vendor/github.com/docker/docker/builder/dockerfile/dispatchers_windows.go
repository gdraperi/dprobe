package dockerfile

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/docker/docker/pkg/system"
)

var pattern = regexp.MustCompile(`^[a-zA-Z]:\.$`)

// normalizeWorkdir normalizes a user requested working directory in a
// platform semantically consistent way.
func normalizeWorkdir(platform string, current string, requested string) (string, error) ***REMOVED***
	if platform == "" ***REMOVED***
		platform = "windows"
	***REMOVED***
	if platform == "windows" ***REMOVED***
		return normalizeWorkdirWindows(current, requested)
	***REMOVED***
	return normalizeWorkdirUnix(current, requested)
***REMOVED***

// normalizeWorkdirUnix normalizes a user requested working directory in a
// platform semantically consistent way.
func normalizeWorkdirUnix(current string, requested string) (string, error) ***REMOVED***
	if requested == "" ***REMOVED***
		return "", errors.New("cannot normalize nothing")
	***REMOVED***
	current = strings.Replace(current, string(os.PathSeparator), "/", -1)
	requested = strings.Replace(requested, string(os.PathSeparator), "/", -1)
	if !path.IsAbs(requested) ***REMOVED***
		return path.Join(`/`, current, requested), nil
	***REMOVED***
	return requested, nil
***REMOVED***

// normalizeWorkdirWindows normalizes a user requested working directory in a
// platform semantically consistent way.
func normalizeWorkdirWindows(current string, requested string) (string, error) ***REMOVED***
	if requested == "" ***REMOVED***
		return "", errors.New("cannot normalize nothing")
	***REMOVED***

	// `filepath.Clean` will replace "" with "." so skip in that case
	if current != "" ***REMOVED***
		current = filepath.Clean(current)
	***REMOVED***
	if requested != "" ***REMOVED***
		requested = filepath.Clean(requested)
	***REMOVED***

	// If either current or requested in Windows is:
	// C:
	// C:.
	// then an error will be thrown as the definition for the above
	// refers to `current directory on drive C:`
	// Since filepath.Clean() will automatically normalize the above
	// to `C:.`, we only need to check the last format
	if pattern.MatchString(current) ***REMOVED***
		return "", fmt.Errorf("%s is not a directory. If you are specifying a drive letter, please add a trailing '\\'", current)
	***REMOVED***
	if pattern.MatchString(requested) ***REMOVED***
		return "", fmt.Errorf("%s is not a directory. If you are specifying a drive letter, please add a trailing '\\'", requested)
	***REMOVED***

	// Target semantics is C:\somefolder, specifically in the format:
	// UPPERCASEDriveLetter-Colon-Backslash-FolderName. We are already
	// guaranteed that `current`, if set, is consistent. This allows us to
	// cope correctly with any of the following in a Dockerfile:
	//	WORKDIR a                       --> C:\a
	//	WORKDIR c:\\foo                 --> C:\foo
	//	WORKDIR \\foo                   --> C:\foo
	//	WORKDIR /foo                    --> C:\foo
	//	WORKDIR c:\\foo \ WORKDIR bar   --> C:\foo --> C:\foo\bar
	//	WORKDIR C:/foo \ WORKDIR bar    --> C:\foo --> C:\foo\bar
	//	WORKDIR C:/foo \ WORKDIR \\bar  --> C:\foo --> C:\bar
	//	WORKDIR /foo \ WORKDIR c:/bar   --> C:\foo --> C:\bar
	if len(current) == 0 || system.IsAbs(requested) ***REMOVED***
		if (requested[0] == os.PathSeparator) ||
			(len(requested) > 1 && string(requested[1]) != ":") ||
			(len(requested) == 1) ***REMOVED***
			requested = filepath.Join(`C:\`, requested)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		requested = filepath.Join(current, requested)
	***REMOVED***
	// Upper-case drive letter
	return (strings.ToUpper(string(requested[0])) + requested[1:]), nil
***REMOVED***

// equalEnvKeys compare two strings and returns true if they are equal. On
// Windows this comparison is case insensitive.
func equalEnvKeys(from, to string) bool ***REMOVED***
	return strings.ToUpper(from) == strings.ToUpper(to)
***REMOVED***
