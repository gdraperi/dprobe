// ©2015 The Go Authors
// Copyright ©2015 Steve Francia <spf@spf13.com>
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

package afero

import "testing"

func checkSizePath(t *testing.T, path string, size int64) ***REMOVED***
	dir, err := testFS.Stat(path)
	if err != nil ***REMOVED***
		t.Fatalf("Stat %q (looking for size %d): %s", path, size, err)
	***REMOVED***
	if dir.Size() != size ***REMOVED***
		t.Errorf("Stat %q: size %d want %d", path, dir.Size(), size)
	***REMOVED***
***REMOVED***

func TestReadFile(t *testing.T) ***REMOVED***
	testFS = &MemMapFs***REMOVED******REMOVED***
	fsutil := &Afero***REMOVED***Fs: testFS***REMOVED***

	testFS.Create("this_exists.go")
	filename := "rumpelstilzchen"
	contents, err := fsutil.ReadFile(filename)
	if err == nil ***REMOVED***
		t.Fatalf("ReadFile %s: error expected, none found", filename)
	***REMOVED***

	filename = "this_exists.go"
	contents, err = fsutil.ReadFile(filename)
	if err != nil ***REMOVED***
		t.Fatalf("ReadFile %s: %v", filename, err)
	***REMOVED***

	checkSizePath(t, filename, int64(len(contents)))
***REMOVED***

func TestWriteFile(t *testing.T) ***REMOVED***
	testFS = &MemMapFs***REMOVED******REMOVED***
	fsutil := &Afero***REMOVED***Fs: testFS***REMOVED***
	f, err := fsutil.TempFile("", "ioutil-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	filename := f.Name()
	data := "Programming today is a race between software engineers striving to " +
		"build bigger and better idiot-proof programs, and the Universe trying " +
		"to produce bigger and better idiots. So far, the Universe is winning."

	if err := fsutil.WriteFile(filename, []byte(data), 0644); err != nil ***REMOVED***
		t.Fatalf("WriteFile %s: %v", filename, err)
	***REMOVED***

	contents, err := fsutil.ReadFile(filename)
	if err != nil ***REMOVED***
		t.Fatalf("ReadFile %s: %v", filename, err)
	***REMOVED***

	if string(contents) != data ***REMOVED***
		t.Fatalf("contents = %q\nexpected = %q", string(contents), data)
	***REMOVED***

	// cleanup
	f.Close()
	testFS.Remove(filename) // ignore error
***REMOVED***

func TestReadDir(t *testing.T) ***REMOVED***
	testFS = &MemMapFs***REMOVED******REMOVED***
	testFS.Mkdir("/i-am-a-dir", 0777)
	testFS.Create("/this_exists.go")
	dirname := "rumpelstilzchen"
	_, err := ReadDir(testFS, dirname)
	if err == nil ***REMOVED***
		t.Fatalf("ReadDir %s: error expected, none found", dirname)
	***REMOVED***

	dirname = ".."
	list, err := ReadDir(testFS, dirname)
	if err != nil ***REMOVED***
		t.Fatalf("ReadDir %s: %v", dirname, err)
	***REMOVED***

	foundFile := false
	foundSubDir := false
	for _, dir := range list ***REMOVED***
		switch ***REMOVED***
		case !dir.IsDir() && dir.Name() == "this_exists.go":
			foundFile = true
		case dir.IsDir() && dir.Name() == "i-am-a-dir":
			foundSubDir = true
		***REMOVED***
	***REMOVED***
	if !foundFile ***REMOVED***
		t.Fatalf("ReadDir %s: this_exists.go file not found", dirname)
	***REMOVED***
	if !foundSubDir ***REMOVED***
		t.Fatalf("ReadDir %s: i-am-a-dir directory not found", dirname)
	***REMOVED***
***REMOVED***
