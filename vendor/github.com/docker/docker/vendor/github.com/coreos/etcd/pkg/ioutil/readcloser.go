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

package ioutil

import (
	"fmt"
	"io"
)

// ReaderAndCloser implements io.ReadCloser interface by combining
// reader and closer together.
type ReaderAndCloser struct ***REMOVED***
	io.Reader
	io.Closer
***REMOVED***

var (
	ErrShortRead = fmt.Errorf("ioutil: short read")
	ErrExpectEOF = fmt.Errorf("ioutil: expect EOF")
)

// NewExactReadCloser returns a ReadCloser that returns errors if the underlying
// reader does not read back exactly the requested number of bytes.
func NewExactReadCloser(rc io.ReadCloser, totalBytes int64) io.ReadCloser ***REMOVED***
	return &exactReadCloser***REMOVED***rc: rc, totalBytes: totalBytes***REMOVED***
***REMOVED***

type exactReadCloser struct ***REMOVED***
	rc         io.ReadCloser
	br         int64
	totalBytes int64
***REMOVED***

func (e *exactReadCloser) Read(p []byte) (int, error) ***REMOVED***
	n, err := e.rc.Read(p)
	e.br += int64(n)
	if e.br > e.totalBytes ***REMOVED***
		return 0, ErrExpectEOF
	***REMOVED***
	if e.br < e.totalBytes && n == 0 ***REMOVED***
		return 0, ErrShortRead
	***REMOVED***
	return n, err
***REMOVED***

func (e *exactReadCloser) Close() error ***REMOVED***
	if err := e.rc.Close(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if e.br < e.totalBytes ***REMOVED***
		return ErrShortRead
	***REMOVED***
	return nil
***REMOVED***
