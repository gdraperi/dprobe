// Copyright © 2014 Steve Francia <spf@spf13.com>.
// Copyright 2009 The Go Authors. All rights reserved.
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

package afero

import (
	"fmt"
	"os"
	"testing"
)

func TestWalk(t *testing.T) ***REMOVED***
	defer removeAllTestFiles(t)
	var testDir string
	for i, fs := range Fss ***REMOVED***
		if i == 0 ***REMOVED***
			testDir = setupTestDirRoot(t, fs)
		***REMOVED*** else ***REMOVED***
			setupTestDirReusePath(t, fs, testDir)
		***REMOVED***
	***REMOVED***

	outputs := make([]string, len(Fss))
	for i, fs := range Fss ***REMOVED***
		walkFn := func(path string, info os.FileInfo, err error) error ***REMOVED***
			if err != nil ***REMOVED***
				t.Error("walkFn err:", err)
			***REMOVED***
			var size int64
			if !info.IsDir() ***REMOVED***
				size = info.Size()
			***REMOVED***
			outputs[i] += fmt.Sprintln(path, info.Name(), size, info.IsDir(), err)
			return nil
		***REMOVED***
		err := Walk(fs, testDir, walkFn)
		if err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***
	***REMOVED***
	fail := false
	for i, o := range outputs ***REMOVED***
		if i == 0 ***REMOVED***
			continue
		***REMOVED***
		if o != outputs[i-1] ***REMOVED***
			fail = true
			break
		***REMOVED***
	***REMOVED***
	if fail ***REMOVED***
		t.Log("Walk outputs not equal!")
		for i, o := range outputs ***REMOVED***
			t.Log(Fss[i].Name() + "\n" + o)
		***REMOVED***
		t.Fail()
	***REMOVED***
***REMOVED***
