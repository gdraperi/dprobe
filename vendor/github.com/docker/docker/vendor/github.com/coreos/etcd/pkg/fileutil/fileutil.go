// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package fileutil implements utility functions related to files and paths.
package fileutil

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/coreos/pkg/capnslog"
)

const (
	// PrivateFileMode grants owner to read/write a file.
	PrivateFileMode = 0600
	// PrivateDirMode grants owner to make/remove files inside the directory.
	PrivateDirMode = 0700
)

var (
	plog = capnslog.NewPackageLogger("github.com/coreos/etcd", "pkg/fileutil")
)

// IsDirWriteable checks if dir is writable by writing and removing a file
// to dir. It returns nil if dir is writable.
func IsDirWriteable(dir string) error ***REMOVED***
	f := filepath.Join(dir, ".touch")
	if err := ioutil.WriteFile(f, []byte(""), PrivateFileMode); err != nil ***REMOVED***
		return err
	***REMOVED***
	return os.Remove(f)
***REMOVED***

// ReadDir returns the filenames in the given directory in sorted order.
func ReadDir(dirpath string) ([]string, error) ***REMOVED***
	dir, err := os.Open(dirpath)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	sort.Strings(names)
	return names, nil
***REMOVED***

// TouchDirAll is similar to os.MkdirAll. It creates directories with 0700 permission if any directory
// does not exists. TouchDirAll also ensures the given directory is writable.
func TouchDirAll(dir string) error ***REMOVED***
	// If path is already a directory, MkdirAll does nothing
	// and returns nil.
	err := os.MkdirAll(dir, PrivateDirMode)
	if err != nil ***REMOVED***
		// if mkdirAll("a/text") and "text" is not
		// a directory, this will return syscall.ENOTDIR
		return err
	***REMOVED***
	return IsDirWriteable(dir)
***REMOVED***

// CreateDirAll is similar to TouchDirAll but returns error
// if the deepest directory was not empty.
func CreateDirAll(dir string) error ***REMOVED***
	err := TouchDirAll(dir)
	if err == nil ***REMOVED***
		var ns []string
		ns, err = ReadDir(dir)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if len(ns) != 0 ***REMOVED***
			err = fmt.Errorf("expected %q to be empty, got %q", dir, ns)
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***

func Exist(name string) bool ***REMOVED***
	_, err := os.Stat(name)
	return err == nil
***REMOVED***

// ZeroToEnd zeros a file starting from SEEK_CUR to its SEEK_END. May temporarily
// shorten the length of the file.
func ZeroToEnd(f *os.File) error ***REMOVED***
	// TODO: support FALLOC_FL_ZERO_RANGE
	off, err := f.Seek(0, io.SeekCurrent)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	lenf, lerr := f.Seek(0, io.SeekEnd)
	if lerr != nil ***REMOVED***
		return lerr
	***REMOVED***
	if err = f.Truncate(off); err != nil ***REMOVED***
		return err
	***REMOVED***
	// make sure blocks remain allocated
	if err = Preallocate(f, lenf, true); err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = f.Seek(off, io.SeekStart)
	return err
***REMOVED***
