package driver

import (
	"io"
	"io/ioutil"
	"os"
	"sort"
)

// ReadFile works the same as ioutil.ReadFile with the Driver abstraction
func ReadFile(r Driver, filename string) ([]byte, error) ***REMOVED***
	f, err := r.Open(filename)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return data, nil
***REMOVED***

// WriteFile works the same as ioutil.WriteFile with the Driver abstraction
func WriteFile(r Driver, filename string, data []byte, perm os.FileMode) error ***REMOVED***
	f, err := r.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer f.Close()

	n, err := f.Write(data)
	if err != nil ***REMOVED***
		return err
	***REMOVED*** else if n != len(data) ***REMOVED***
		return io.ErrShortWrite
	***REMOVED***

	return nil
***REMOVED***

// ReadDir works the same as ioutil.ReadDir with the Driver abstraction
func ReadDir(r Driver, dirname string) ([]os.FileInfo, error) ***REMOVED***
	f, err := r.Open(dirname)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()

	dirs, err := f.Readdir(-1)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	sort.Sort(fileInfos(dirs))
	return dirs, nil
***REMOVED***

// Simple implementation of the sort.Interface for os.FileInfo
type fileInfos []os.FileInfo

func (fis fileInfos) Len() int ***REMOVED***
	return len(fis)
***REMOVED***

func (fis fileInfos) Less(i, j int) bool ***REMOVED***
	return fis[i].Name() < fis[j].Name()
***REMOVED***

func (fis fileInfos) Swap(i, j int) ***REMOVED***
	fis[i], fis[j] = fis[j], fis[i]
***REMOVED***
