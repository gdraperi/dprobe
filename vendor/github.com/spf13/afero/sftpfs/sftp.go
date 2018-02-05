// Copyright © 2015 Jerry Jacobs <jerry.jacobs@xor-gate.org>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sftpfs

import (
	"os"
	"time"

	"github.com/pkg/sftp"
	"github.com/spf13/afero"
)

// Fs is a afero.Fs implementation that uses functions provided by the sftp package.
//
// For details in any method, check the documentation of the sftp package
// (github.com/pkg/sftp).
type Fs struct ***REMOVED***
	client *sftp.Client
***REMOVED***

func New(client *sftp.Client) afero.Fs ***REMOVED***
	return &Fs***REMOVED***client: client***REMOVED***
***REMOVED***

func (s Fs) Name() string ***REMOVED*** return "sftpfs" ***REMOVED***

func (s Fs) Create(name string) (afero.File, error) ***REMOVED***
	return FileCreate(s.client, name)
***REMOVED***

func (s Fs) Mkdir(name string, perm os.FileMode) error ***REMOVED***
	err := s.client.Mkdir(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.client.Chmod(name, perm)
***REMOVED***

func (s Fs) MkdirAll(path string, perm os.FileMode) error ***REMOVED***
	// Fast path: if we can tell whether path is a directory or file, stop with success or error.
	dir, err := s.Stat(path)
	if err == nil ***REMOVED***
		if dir.IsDir() ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***

	// Slow path: make sure parent exists and then call Mkdir for path.
	i := len(path)
	for i > 0 && os.IsPathSeparator(path[i-1]) ***REMOVED*** // Skip trailing path separator.
		i--
	***REMOVED***

	j := i
	for j > 0 && !os.IsPathSeparator(path[j-1]) ***REMOVED*** // Scan backward over element.
		j--
	***REMOVED***

	if j > 1 ***REMOVED***
		// Create parent
		err = s.MkdirAll(path[0:j-1], perm)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Parent now exists; invoke Mkdir and use its result.
	err = s.Mkdir(path, perm)
	if err != nil ***REMOVED***
		// Handle arguments like "foo/." by
		// double-checking that directory doesn't exist.
		dir, err1 := s.Lstat(path)
		if err1 == nil && dir.IsDir() ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func (s Fs) Open(name string) (afero.File, error) ***REMOVED***
	return FileOpen(s.client, name)
***REMOVED***

func (s Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) ***REMOVED***
	return nil, nil
***REMOVED***

func (s Fs) Remove(name string) error ***REMOVED***
	return s.client.Remove(name)
***REMOVED***

func (s Fs) RemoveAll(path string) error ***REMOVED***
	// TODO have a look at os.RemoveAll
	// https://github.com/golang/go/blob/master/src/os/path.go#L66
	return nil
***REMOVED***

func (s Fs) Rename(oldname, newname string) error ***REMOVED***
	return s.client.Rename(oldname, newname)
***REMOVED***

func (s Fs) Stat(name string) (os.FileInfo, error) ***REMOVED***
	return s.client.Stat(name)
***REMOVED***

func (s Fs) Lstat(p string) (os.FileInfo, error) ***REMOVED***
	return s.client.Lstat(p)
***REMOVED***

func (s Fs) Chmod(name string, mode os.FileMode) error ***REMOVED***
	return s.client.Chmod(name, mode)
***REMOVED***

func (s Fs) Chtimes(name string, atime time.Time, mtime time.Time) error ***REMOVED***
	return s.client.Chtimes(name, atime, mtime)
***REMOVED***
