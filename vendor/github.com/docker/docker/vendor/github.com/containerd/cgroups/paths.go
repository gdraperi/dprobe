package cgroups

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
)

type Path func(subsystem Name) (string, error)

func RootPath(subsysem Name) (string, error) ***REMOVED***
	return "/", nil
***REMOVED***

// StaticPath returns a static path to use for all cgroups
func StaticPath(path string) Path ***REMOVED***
	return func(_ Name) (string, error) ***REMOVED***
		return path, nil
	***REMOVED***
***REMOVED***

// NestedPath will nest the cgroups based on the calling processes cgroup
// placing its child processes inside its own path
func NestedPath(suffix string) Path ***REMOVED***
	paths, err := parseCgroupFile("/proc/self/cgroup")
	if err != nil ***REMOVED***
		return errorPath(err)
	***REMOVED***
	return existingPath(paths, suffix)
***REMOVED***

// PidPath will return the correct cgroup paths for an existing process running inside a cgroup
// This is commonly used for the Load function to restore an existing container
func PidPath(pid int) Path ***REMOVED***
	p := fmt.Sprintf("/proc/%d/cgroup", pid)
	paths, err := parseCgroupFile(p)
	if err != nil ***REMOVED***
		return errorPath(errors.Wrapf(err, "parse cgroup file %s", p))
	***REMOVED***
	return existingPath(paths, "")
***REMOVED***

func existingPath(paths map[string]string, suffix string) Path ***REMOVED***
	// localize the paths based on the root mount dest for nested cgroups
	for n, p := range paths ***REMOVED***
		dest, err := getCgroupDestination(string(n))
		if err != nil ***REMOVED***
			return errorPath(err)
		***REMOVED***
		rel, err := filepath.Rel(dest, p)
		if err != nil ***REMOVED***
			return errorPath(err)
		***REMOVED***
		if rel == "." ***REMOVED***
			rel = dest
		***REMOVED***
		paths[n] = filepath.Join("/", rel)
	***REMOVED***
	return func(name Name) (string, error) ***REMOVED***
		root, ok := paths[string(name)]
		if !ok ***REMOVED***
			if root, ok = paths[fmt.Sprintf("name=%s", name)]; !ok ***REMOVED***
				return "", fmt.Errorf("unable to find %q in controller set", name)
			***REMOVED***
		***REMOVED***
		if suffix != "" ***REMOVED***
			return filepath.Join(root, suffix), nil
		***REMOVED***
		return root, nil
	***REMOVED***
***REMOVED***

func subPath(path Path, subName string) Path ***REMOVED***
	return func(name Name) (string, error) ***REMOVED***
		p, err := path(name)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		return filepath.Join(p, subName), nil
	***REMOVED***
***REMOVED***

func errorPath(err error) Path ***REMOVED***
	return func(_ Name) (string, error) ***REMOVED***
		return "", err
	***REMOVED***
***REMOVED***
