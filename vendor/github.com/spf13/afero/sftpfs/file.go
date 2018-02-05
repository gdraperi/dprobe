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
	"github.com/pkg/sftp"
	"os"
)

type File struct ***REMOVED***
	fd *sftp.File
***REMOVED***

func FileOpen(s *sftp.Client, name string) (*File, error) ***REMOVED***
	fd, err := s.Open(name)
	if err != nil ***REMOVED***
		return &File***REMOVED******REMOVED***, err
	***REMOVED***
	return &File***REMOVED***fd: fd***REMOVED***, nil
***REMOVED***

func FileCreate(s *sftp.Client, name string) (*File, error) ***REMOVED***
	fd, err := s.Create(name)
	if err != nil ***REMOVED***
		return &File***REMOVED******REMOVED***, err
	***REMOVED***
	return &File***REMOVED***fd: fd***REMOVED***, nil
***REMOVED***

func (f *File) Close() error ***REMOVED***
	return f.fd.Close()
***REMOVED***

func (f *File) Name() string ***REMOVED***
	return f.fd.Name()
***REMOVED***

func (f *File) Stat() (os.FileInfo, error) ***REMOVED***
	return f.fd.Stat()
***REMOVED***

func (f *File) Sync() error ***REMOVED***
	return nil
***REMOVED***

func (f *File) Truncate(size int64) error ***REMOVED***
	return f.fd.Truncate(size)
***REMOVED***

func (f *File) Read(b []byte) (n int, err error) ***REMOVED***
	return f.fd.Read(b)
***REMOVED***

// TODO
func (f *File) ReadAt(b []byte, off int64) (n int, err error) ***REMOVED***
	return 0, nil
***REMOVED***

// TODO
func (f *File) Readdir(count int) (res []os.FileInfo, err error) ***REMOVED***
	return nil, nil
***REMOVED***

// TODO
func (f *File) Readdirnames(n int) (names []string, err error) ***REMOVED***
	return nil, nil
***REMOVED***

func (f *File) Seek(offset int64, whence int) (int64, error) ***REMOVED***
	return f.fd.Seek(offset, whence)
***REMOVED***

func (f *File) Write(b []byte) (n int, err error) ***REMOVED***
	return f.fd.Write(b)
***REMOVED***

// TODO
func (f *File) WriteAt(b []byte, off int64) (n int, err error) ***REMOVED***
	return 0, nil
***REMOVED***

func (f *File) WriteString(s string) (ret int, err error) ***REMOVED***
	return f.fd.Write([]byte(s))
***REMOVED***
