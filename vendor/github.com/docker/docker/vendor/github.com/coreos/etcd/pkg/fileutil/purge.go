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

package fileutil

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func PurgeFile(dirname string, suffix string, max uint, interval time.Duration, stop <-chan struct***REMOVED******REMOVED***) <-chan error ***REMOVED***
	return purgeFile(dirname, suffix, max, interval, stop, nil)
***REMOVED***

// purgeFile is the internal implementation for PurgeFile which can post purged files to purgec if non-nil.
func purgeFile(dirname string, suffix string, max uint, interval time.Duration, stop <-chan struct***REMOVED******REMOVED***, purgec chan<- string) <-chan error ***REMOVED***
	errC := make(chan error, 1)
	go func() ***REMOVED***
		for ***REMOVED***
			fnames, err := ReadDir(dirname)
			if err != nil ***REMOVED***
				errC <- err
				return
			***REMOVED***
			newfnames := make([]string, 0)
			for _, fname := range fnames ***REMOVED***
				if strings.HasSuffix(fname, suffix) ***REMOVED***
					newfnames = append(newfnames, fname)
				***REMOVED***
			***REMOVED***
			sort.Strings(newfnames)
			fnames = newfnames
			for len(newfnames) > int(max) ***REMOVED***
				f := filepath.Join(dirname, newfnames[0])
				l, err := TryLockFile(f, os.O_WRONLY, PrivateFileMode)
				if err != nil ***REMOVED***
					break
				***REMOVED***
				if err = os.Remove(f); err != nil ***REMOVED***
					errC <- err
					return
				***REMOVED***
				if err = l.Close(); err != nil ***REMOVED***
					plog.Errorf("error unlocking %s when purging file (%v)", l.Name(), err)
					errC <- err
					return
				***REMOVED***
				plog.Infof("purged file %s successfully", f)
				newfnames = newfnames[1:]
			***REMOVED***
			if purgec != nil ***REMOVED***
				for i := 0; i < len(fnames)-len(newfnames); i++ ***REMOVED***
					purgec <- fnames[i]
				***REMOVED***
			***REMOVED***
			select ***REMOVED***
			case <-time.After(interval):
			case <-stop:
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	return errC
***REMOVED***
