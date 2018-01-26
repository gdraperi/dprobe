// +build !windows

package dockerfile

import (
	"errors"
	"os"
	"path/filepath"
)

// normalizeWorkdir normalizes a user requested working directory in a
// platform semantically consistent way.
func normalizeWorkdir(_ string, current string, requested string) (string, error) ***REMOVED***
	if requested == "" ***REMOVED***
		return "", errors.New("cannot normalize nothing")
	***REMOVED***
	current = filepath.FromSlash(current)
	requested = filepath.FromSlash(requested)
	if !filepath.IsAbs(requested) ***REMOVED***
		return filepath.Join(string(os.PathSeparator), current, requested), nil
	***REMOVED***
	return requested, nil
***REMOVED***

// equalEnvKeys compare two strings and returns true if they are equal. On
// Windows this comparison is case insensitive.
func equalEnvKeys(from, to string) bool ***REMOVED***
	return from == to
***REMOVED***
