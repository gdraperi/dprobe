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
	"io"
	"os"

	"github.com/coreos/etcd/pkg/fileutil"
)

// WriteAndSyncFile behaves just like ioutil.WriteFile in the standard library,
// but calls Sync before closing the file. WriteAndSyncFile guarantees the data
// is synced if there is no error returned.
func WriteAndSyncFile(filename string, data []byte, perm os.FileMode) error ***REMOVED***
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	n, err := f.Write(data)
	if err == nil && n < len(data) ***REMOVED***
		err = io.ErrShortWrite
	***REMOVED***
	if err == nil ***REMOVED***
		err = fileutil.Fsync(f)
	***REMOVED***
	if err1 := f.Close(); err == nil ***REMOVED***
		err = err1
	***REMOVED***
	return err
***REMOVED***
