// +build linux

package aufs

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
)

// Return all the directories
func loadIds(root string) ([]string, error) ***REMOVED***
	dirs, err := ioutil.ReadDir(root)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	out := []string***REMOVED******REMOVED***
	for _, d := range dirs ***REMOVED***
		if !d.IsDir() ***REMOVED***
			out = append(out, d.Name())
		***REMOVED***
	***REMOVED***
	return out, nil
***REMOVED***

// Read the layers file for the current id and return all the
// layers represented by new lines in the file
//
// If there are no lines in the file then the id has no parent
// and an empty slice is returned.
func getParentIDs(root, id string) ([]string, error) ***REMOVED***
	f, err := os.Open(path.Join(root, "layers", id))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()

	out := []string***REMOVED******REMOVED***
	s := bufio.NewScanner(f)

	for s.Scan() ***REMOVED***
		if t := s.Text(); t != "" ***REMOVED***
			out = append(out, s.Text())
		***REMOVED***
	***REMOVED***
	return out, s.Err()
***REMOVED***

func (a *Driver) getMountpoint(id string) string ***REMOVED***
	return path.Join(a.mntPath(), id)
***REMOVED***

func (a *Driver) mntPath() string ***REMOVED***
	return path.Join(a.rootPath(), "mnt")
***REMOVED***

func (a *Driver) getDiffPath(id string) string ***REMOVED***
	return path.Join(a.diffPath(), id)
***REMOVED***

func (a *Driver) diffPath() string ***REMOVED***
	return path.Join(a.rootPath(), "diff")
***REMOVED***
