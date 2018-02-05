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
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// contains returns true if vector contains the string s.
func contains(vector []string, s string) bool ***REMOVED***
	for _, elem := range vector ***REMOVED***
		if elem == s ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func setupGlobDirRoot(t *testing.T, fs Fs) string ***REMOVED***
	path := testDir(fs)
	setupGlobFiles(t, fs, path)
	return path
***REMOVED***

func setupGlobDirReusePath(t *testing.T, fs Fs, path string) string ***REMOVED***
	testRegistry[fs] = append(testRegistry[fs], path)
	return setupGlobFiles(t, fs, path)
***REMOVED***

func setupGlobFiles(t *testing.T, fs Fs, path string) string ***REMOVED***
	testSubDir := filepath.Join(path, "globs", "bobs")
	err := fs.MkdirAll(testSubDir, 0700)
	if err != nil && !os.IsExist(err) ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	f, err := fs.Create(filepath.Join(testSubDir, "/matcher"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	f.WriteString("Testfile 1 content")
	f.Close()

	f, err = fs.Create(filepath.Join(testSubDir, "/../submatcher"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	f.WriteString("Testfile 2 content")
	f.Close()

	f, err = fs.Create(filepath.Join(testSubDir, "/../../match"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	f.WriteString("Testfile 3 content")
	f.Close()

	return testSubDir
***REMOVED***

func TestGlob(t *testing.T) ***REMOVED***
	defer removeAllTestFiles(t)
	var testDir string
	for i, fs := range Fss ***REMOVED***
		if i == 0 ***REMOVED***
			testDir = setupGlobDirRoot(t, fs)
		***REMOVED*** else ***REMOVED***
			setupGlobDirReusePath(t, fs, testDir)
		***REMOVED***
	***REMOVED***

	var globTests = []struct ***REMOVED***
		pattern, result string
	***REMOVED******REMOVED***
		***REMOVED***testDir + "/globs/bobs/matcher", testDir + "/globs/bobs/matcher"***REMOVED***,
		***REMOVED***testDir + "/globs/*/mat?her", testDir + "/globs/bobs/matcher"***REMOVED***,
		***REMOVED***testDir + "/globs/bobs/../*", testDir + "/globs/submatcher"***REMOVED***,
		***REMOVED***testDir + "/match", testDir + "/match"***REMOVED***,
	***REMOVED***

	for _, fs := range Fss ***REMOVED***

		for _, tt := range globTests ***REMOVED***
			pattern := tt.pattern
			result := tt.result
			if runtime.GOOS == "windows" ***REMOVED***
				pattern = filepath.Clean(pattern)
				result = filepath.Clean(result)
			***REMOVED***
			matches, err := Glob(fs, pattern)
			if err != nil ***REMOVED***
				t.Errorf("Glob error for %q: %s", pattern, err)
				continue
			***REMOVED***
			if !contains(matches, result) ***REMOVED***
				t.Errorf("Glob(%#q) = %#v want %v", pattern, matches, result)
			***REMOVED***
		***REMOVED***
		for _, pattern := range []string***REMOVED***"no_match", "../*/no_match"***REMOVED*** ***REMOVED***
			matches, err := Glob(fs, pattern)
			if err != nil ***REMOVED***
				t.Errorf("Glob error for %q: %s", pattern, err)
				continue
			***REMOVED***
			if len(matches) != 0 ***REMOVED***
				t.Errorf("Glob(%#q) = %#v want []", pattern, matches)
			***REMOVED***
		***REMOVED***

	***REMOVED***
***REMOVED***

func TestGlobSymlink(t *testing.T) ***REMOVED***
	defer removeAllTestFiles(t)

	fs := &OsFs***REMOVED******REMOVED***
	testDir := setupGlobDirRoot(t, fs)

	err := os.Symlink("target", filepath.Join(testDir, "symlink"))
	if err != nil ***REMOVED***
		t.Skipf("skipping on %s", runtime.GOOS)
	***REMOVED***

	var globSymlinkTests = []struct ***REMOVED***
		path, dest string
		brokenLink bool
	***REMOVED******REMOVED***
		***REMOVED***"test1", "link1", false***REMOVED***,
		***REMOVED***"test2", "link2", true***REMOVED***,
	***REMOVED***

	for _, tt := range globSymlinkTests ***REMOVED***
		path := filepath.Join(testDir, tt.path)
		dest := filepath.Join(testDir, tt.dest)
		f, err := fs.Create(path)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if err := f.Close(); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		err = os.Symlink(path, dest)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if tt.brokenLink ***REMOVED***
			// Break the symlink.
			fs.Remove(path)
		***REMOVED***
		matches, err := Glob(fs, dest)
		if err != nil ***REMOVED***
			t.Errorf("GlobSymlink error for %q: %s", dest, err)
		***REMOVED***
		if !contains(matches, dest) ***REMOVED***
			t.Errorf("Glob(%#q) = %#v want %v", dest, matches, dest)
		***REMOVED***
	***REMOVED***
***REMOVED***


func TestGlobError(t *testing.T) ***REMOVED***
	for _, fs := range Fss ***REMOVED***
		_, err := Glob(fs, "[7]")
		if err != nil ***REMOVED***
			t.Error("expected error for bad pattern; got none")
		***REMOVED***
	***REMOVED***
***REMOVED***
