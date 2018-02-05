// Copyright ©2015 Steve Francia <spf@spf13.com>
// Portions Copyright ©2015 The Hugo Authors
//
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
//

package afero

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

var testFS = new(MemMapFs)

func TestDirExists(t *testing.T) ***REMOVED***
	type test struct ***REMOVED***
		input    string
		expected bool
	***REMOVED***

	// First create a couple directories so there is something in the filesystem
	//testFS := new(MemMapFs)
	testFS.MkdirAll("/foo/bar", 0777)

	data := []test***REMOVED***
		***REMOVED***".", true***REMOVED***,
		***REMOVED***"./", true***REMOVED***,
		***REMOVED***"..", true***REMOVED***,
		***REMOVED***"../", true***REMOVED***,
		***REMOVED***"./..", true***REMOVED***,
		***REMOVED***"./../", true***REMOVED***,
		***REMOVED***"/foo/", true***REMOVED***,
		***REMOVED***"/foo", true***REMOVED***,
		***REMOVED***"/foo/bar", true***REMOVED***,
		***REMOVED***"/foo/bar/", true***REMOVED***,
		***REMOVED***"/", true***REMOVED***,
		***REMOVED***"/some-really-random-directory-name", false***REMOVED***,
		***REMOVED***"/some/really/random/directory/name", false***REMOVED***,
		***REMOVED***"./some-really-random-local-directory-name", false***REMOVED***,
		***REMOVED***"./some/really/random/local/directory/name", false***REMOVED***,
	***REMOVED***

	for i, d := range data ***REMOVED***
		exists, _ := DirExists(testFS, filepath.FromSlash(d.input))
		if d.expected != exists ***REMOVED***
			t.Errorf("Test %d %q failed. Expected %t got %t", i, d.input, d.expected, exists)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestIsDir(t *testing.T) ***REMOVED***
	testFS = new(MemMapFs)

	type test struct ***REMOVED***
		input    string
		expected bool
	***REMOVED***
	data := []test***REMOVED***
		***REMOVED***"./", true***REMOVED***,
		***REMOVED***"/", true***REMOVED***,
		***REMOVED***"./this-directory-does-not-existi", false***REMOVED***,
		***REMOVED***"/this-absolute-directory/does-not-exist", false***REMOVED***,
	***REMOVED***

	for i, d := range data ***REMOVED***

		exists, _ := IsDir(testFS, d.input)
		if d.expected != exists ***REMOVED***
			t.Errorf("Test %d failed. Expected %t got %t", i, d.expected, exists)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestIsEmpty(t *testing.T) ***REMOVED***
	testFS = new(MemMapFs)

	zeroSizedFile, _ := createZeroSizedFileInTempDir()
	defer deleteFileInTempDir(zeroSizedFile)
	nonZeroSizedFile, _ := createNonZeroSizedFileInTempDir()
	defer deleteFileInTempDir(nonZeroSizedFile)
	emptyDirectory, _ := createEmptyTempDir()
	defer deleteTempDir(emptyDirectory)
	nonEmptyZeroLengthFilesDirectory, _ := createTempDirWithZeroLengthFiles()
	defer deleteTempDir(nonEmptyZeroLengthFilesDirectory)
	nonEmptyNonZeroLengthFilesDirectory, _ := createTempDirWithNonZeroLengthFiles()
	defer deleteTempDir(nonEmptyNonZeroLengthFilesDirectory)
	nonExistentFile := os.TempDir() + "/this-file-does-not-exist.txt"
	nonExistentDir := os.TempDir() + "/this/direcotry/does/not/exist/"

	fileDoesNotExist := fmt.Errorf("%q path does not exist", nonExistentFile)
	dirDoesNotExist := fmt.Errorf("%q path does not exist", nonExistentDir)

	type test struct ***REMOVED***
		input          string
		expectedResult bool
		expectedErr    error
	***REMOVED***

	data := []test***REMOVED***
		***REMOVED***zeroSizedFile.Name(), true, nil***REMOVED***,
		***REMOVED***nonZeroSizedFile.Name(), false, nil***REMOVED***,
		***REMOVED***emptyDirectory, true, nil***REMOVED***,
		***REMOVED***nonEmptyZeroLengthFilesDirectory, false, nil***REMOVED***,
		***REMOVED***nonEmptyNonZeroLengthFilesDirectory, false, nil***REMOVED***,
		***REMOVED***nonExistentFile, false, fileDoesNotExist***REMOVED***,
		***REMOVED***nonExistentDir, false, dirDoesNotExist***REMOVED***,
	***REMOVED***
	for i, d := range data ***REMOVED***
		exists, err := IsEmpty(testFS, d.input)
		if d.expectedResult != exists ***REMOVED***
			t.Errorf("Test %d %q failed exists. Expected result %t got %t", i, d.input, d.expectedResult, exists)
		***REMOVED***
		if d.expectedErr != nil ***REMOVED***
			if d.expectedErr.Error() != err.Error() ***REMOVED***
				t.Errorf("Test %d failed with err. Expected %q(%#v) got %q(%#v)", i, d.expectedErr, d.expectedErr, err, err)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if d.expectedErr != err ***REMOVED***
				t.Errorf("Test %d failed. Expected error %q(%#v) got %q(%#v)", i, d.expectedErr, d.expectedErr, err, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestReaderContains(t *testing.T) ***REMOVED***
	for i, this := range []struct ***REMOVED***
		v1     string
		v2     [][]byte
		expect bool
	***REMOVED******REMOVED***
		***REMOVED***"abc", [][]byte***REMOVED***[]byte("a")***REMOVED***, true***REMOVED***,
		***REMOVED***"abc", [][]byte***REMOVED***[]byte("b")***REMOVED***, true***REMOVED***,
		***REMOVED***"abcdefg", [][]byte***REMOVED***[]byte("efg")***REMOVED***, true***REMOVED***,
		***REMOVED***"abc", [][]byte***REMOVED***[]byte("d")***REMOVED***, false***REMOVED***,
		***REMOVED***"abc", [][]byte***REMOVED***[]byte("d"), []byte("e")***REMOVED***, false***REMOVED***,
		***REMOVED***"abc", [][]byte***REMOVED***[]byte("d"), []byte("a")***REMOVED***, true***REMOVED***,
		***REMOVED***"abc", [][]byte***REMOVED***[]byte("b"), []byte("e")***REMOVED***, true***REMOVED***,
		***REMOVED***"", nil, false***REMOVED***,
		***REMOVED***"", [][]byte***REMOVED***[]byte("a")***REMOVED***, false***REMOVED***,
		***REMOVED***"a", [][]byte***REMOVED***[]byte("")***REMOVED***, false***REMOVED***,
		***REMOVED***"", [][]byte***REMOVED***[]byte("")***REMOVED***, false***REMOVED******REMOVED*** ***REMOVED***
		result := readerContainsAny(strings.NewReader(this.v1), this.v2...)
		if result != this.expect ***REMOVED***
			t.Errorf("[%d] readerContains: got %t but expected %t", i, result, this.expect)
		***REMOVED***
	***REMOVED***

	if readerContainsAny(nil, []byte("a")) ***REMOVED***
		t.Error("readerContains with nil reader")
	***REMOVED***

	if readerContainsAny(nil, nil) ***REMOVED***
		t.Error("readerContains with nil arguments")
	***REMOVED***
***REMOVED***

func createZeroSizedFileInTempDir() (File, error) ***REMOVED***
	filePrefix := "_path_test_"
	f, e := TempFile(testFS, "", filePrefix) // dir is os.TempDir()
	if e != nil ***REMOVED***
		// if there was an error no file was created.
		// => no requirement to delete the file
		return nil, e
	***REMOVED***
	return f, nil
***REMOVED***

func createNonZeroSizedFileInTempDir() (File, error) ***REMOVED***
	f, err := createZeroSizedFileInTempDir()
	if err != nil ***REMOVED***
		// no file ??
	***REMOVED***
	byteString := []byte("byteString")
	err = WriteFile(testFS, f.Name(), byteString, 0644)
	if err != nil ***REMOVED***
		// delete the file
		deleteFileInTempDir(f)
		return nil, err
	***REMOVED***
	return f, nil
***REMOVED***

func deleteFileInTempDir(f File) ***REMOVED***
	err := testFS.Remove(f.Name())
	if err != nil ***REMOVED***
		// now what?
	***REMOVED***
***REMOVED***

func createEmptyTempDir() (string, error) ***REMOVED***
	dirPrefix := "_dir_prefix_"
	d, e := TempDir(testFS, "", dirPrefix) // will be in os.TempDir()
	if e != nil ***REMOVED***
		// no directory to delete - it was never created
		return "", e
	***REMOVED***
	return d, nil
***REMOVED***

func createTempDirWithZeroLengthFiles() (string, error) ***REMOVED***
	d, dirErr := createEmptyTempDir()
	if dirErr != nil ***REMOVED***
		//now what?
	***REMOVED***
	filePrefix := "_path_test_"
	_, fileErr := TempFile(testFS, d, filePrefix) // dir is os.TempDir()
	if fileErr != nil ***REMOVED***
		// if there was an error no file was created.
		// but we need to remove the directory to clean-up
		deleteTempDir(d)
		return "", fileErr
	***REMOVED***
	// the dir now has one, zero length file in it
	return d, nil

***REMOVED***

func createTempDirWithNonZeroLengthFiles() (string, error) ***REMOVED***
	d, dirErr := createEmptyTempDir()
	if dirErr != nil ***REMOVED***
		//now what?
	***REMOVED***
	filePrefix := "_path_test_"
	f, fileErr := TempFile(testFS, d, filePrefix) // dir is os.TempDir()
	if fileErr != nil ***REMOVED***
		// if there was an error no file was created.
		// but we need to remove the directory to clean-up
		deleteTempDir(d)
		return "", fileErr
	***REMOVED***
	byteString := []byte("byteString")
	fileErr = WriteFile(testFS, f.Name(), byteString, 0644)
	if fileErr != nil ***REMOVED***
		// delete the file
		deleteFileInTempDir(f)
		// also delete the directory
		deleteTempDir(d)
		return "", fileErr
	***REMOVED***

	// the dir now has one, zero length file in it
	return d, nil

***REMOVED***

func TestExists(t *testing.T) ***REMOVED***
	zeroSizedFile, _ := createZeroSizedFileInTempDir()
	defer deleteFileInTempDir(zeroSizedFile)
	nonZeroSizedFile, _ := createNonZeroSizedFileInTempDir()
	defer deleteFileInTempDir(nonZeroSizedFile)
	emptyDirectory, _ := createEmptyTempDir()
	defer deleteTempDir(emptyDirectory)
	nonExistentFile := os.TempDir() + "/this-file-does-not-exist.txt"
	nonExistentDir := os.TempDir() + "/this/direcotry/does/not/exist/"

	type test struct ***REMOVED***
		input          string
		expectedResult bool
		expectedErr    error
	***REMOVED***

	data := []test***REMOVED***
		***REMOVED***zeroSizedFile.Name(), true, nil***REMOVED***,
		***REMOVED***nonZeroSizedFile.Name(), true, nil***REMOVED***,
		***REMOVED***emptyDirectory, true, nil***REMOVED***,
		***REMOVED***nonExistentFile, false, nil***REMOVED***,
		***REMOVED***nonExistentDir, false, nil***REMOVED***,
	***REMOVED***
	for i, d := range data ***REMOVED***
		exists, err := Exists(testFS, d.input)
		if d.expectedResult != exists ***REMOVED***
			t.Errorf("Test %d failed. Expected result %t got %t", i, d.expectedResult, exists)
		***REMOVED***
		if d.expectedErr != err ***REMOVED***
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expectedErr, err)
		***REMOVED***
	***REMOVED***

***REMOVED***

func TestSafeWriteToDisk(t *testing.T) ***REMOVED***
	emptyFile, _ := createZeroSizedFileInTempDir()
	defer deleteFileInTempDir(emptyFile)
	tmpDir, _ := createEmptyTempDir()
	defer deleteTempDir(tmpDir)

	randomString := "This is a random string!"
	reader := strings.NewReader(randomString)

	fileExists := fmt.Errorf("%v already exists", emptyFile.Name())

	type test struct ***REMOVED***
		filename    string
		expectedErr error
	***REMOVED***

	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	data := []test***REMOVED***
		***REMOVED***emptyFile.Name(), fileExists***REMOVED***,
		***REMOVED***tmpDir + "/" + nowStr, nil***REMOVED***,
	***REMOVED***

	for i, d := range data ***REMOVED***
		e := SafeWriteReader(testFS, d.filename, reader)
		if d.expectedErr != nil ***REMOVED***
			if d.expectedErr.Error() != e.Error() ***REMOVED***
				t.Errorf("Test %d failed. Expected error %q but got %q", i, d.expectedErr.Error(), e.Error())
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if d.expectedErr != e ***REMOVED***
				t.Errorf("Test %d failed. Expected %q but got %q", i, d.expectedErr, e)
			***REMOVED***
			contents, _ := ReadFile(testFS, d.filename)
			if randomString != string(contents) ***REMOVED***
				t.Errorf("Test %d failed. Expected contents %q but got %q", i, randomString, string(contents))
			***REMOVED***
		***REMOVED***
		reader.Seek(0, 0)
	***REMOVED***
***REMOVED***

func TestWriteToDisk(t *testing.T) ***REMOVED***
	emptyFile, _ := createZeroSizedFileInTempDir()
	defer deleteFileInTempDir(emptyFile)
	tmpDir, _ := createEmptyTempDir()
	defer deleteTempDir(tmpDir)

	randomString := "This is a random string!"
	reader := strings.NewReader(randomString)

	type test struct ***REMOVED***
		filename    string
		expectedErr error
	***REMOVED***

	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	data := []test***REMOVED***
		***REMOVED***emptyFile.Name(), nil***REMOVED***,
		***REMOVED***tmpDir + "/" + nowStr, nil***REMOVED***,
	***REMOVED***

	for i, d := range data ***REMOVED***
		e := WriteReader(testFS, d.filename, reader)
		if d.expectedErr != e ***REMOVED***
			t.Errorf("Test %d failed. WriteToDisk Error Expected %q but got %q", i, d.expectedErr, e)
		***REMOVED***
		contents, e := ReadFile(testFS, d.filename)
		if e != nil ***REMOVED***
			t.Errorf("Test %d failed. Could not read file %s. Reason: %s\n", i, d.filename, e)
		***REMOVED***
		if randomString != string(contents) ***REMOVED***
			t.Errorf("Test %d failed. Expected contents %q but got %q", i, randomString, string(contents))
		***REMOVED***
		reader.Seek(0, 0)
	***REMOVED***
***REMOVED***

func TestGetTempDir(t *testing.T) ***REMOVED***
	dir := os.TempDir()
	if FilePathSeparator != dir[len(dir)-1:] ***REMOVED***
		dir = dir + FilePathSeparator
	***REMOVED***
	testDir := "hugoTestFolder" + FilePathSeparator
	tests := []struct ***REMOVED***
		input    string
		expected string
	***REMOVED******REMOVED***
		***REMOVED***"", dir***REMOVED***,
		***REMOVED***testDir + "  Foo bar  ", dir + testDir + "  Foo bar  " + FilePathSeparator***REMOVED***,
		***REMOVED***testDir + "Foo.Bar/foo_Bar-Foo", dir + testDir + "Foo.Bar/foo_Bar-Foo" + FilePathSeparator***REMOVED***,
		***REMOVED***testDir + "fOO,bar:foo%bAR", dir + testDir + "fOObarfoo%bAR" + FilePathSeparator***REMOVED***,
		***REMOVED***testDir + "FOo/BaR.html", dir + testDir + "FOo/BaR.html" + FilePathSeparator***REMOVED***,
		***REMOVED***testDir + "трям/трям", dir + testDir + "трям/трям" + FilePathSeparator***REMOVED***,
		***REMOVED***testDir + "은행", dir + testDir + "은행" + FilePathSeparator***REMOVED***,
		***REMOVED***testDir + "Банковский кассир", dir + testDir + "Банковский кассир" + FilePathSeparator***REMOVED***,
	***REMOVED***

	for _, test := range tests ***REMOVED***
		output := GetTempDir(new(MemMapFs), test.input)
		if output != test.expected ***REMOVED***
			t.Errorf("Expected %#v, got %#v\n", test.expected, output)
		***REMOVED***
	***REMOVED***
***REMOVED***

// This function is very dangerous. Don't use it.
func deleteTempDir(d string) ***REMOVED***
	err := os.RemoveAll(d)
	if err != nil ***REMOVED***
		// now what?
	***REMOVED***
***REMOVED***

func TestFullBaseFsPath(t *testing.T) ***REMOVED***
	type dirSpec struct ***REMOVED***
		Dir1, Dir2, Dir3 string
	***REMOVED***
	dirSpecs := []dirSpec***REMOVED***
		dirSpec***REMOVED***Dir1: "/", Dir2: "/", Dir3: "/"***REMOVED***,
		dirSpec***REMOVED***Dir1: "/", Dir2: "/path2", Dir3: "/"***REMOVED***,
		dirSpec***REMOVED***Dir1: "/path1/dir", Dir2: "/path2/dir/", Dir3: "/path3/dir"***REMOVED***,
		dirSpec***REMOVED***Dir1: "C:/path1", Dir2: "path2/dir", Dir3: "/path3/dir/"***REMOVED***,
	***REMOVED***

	for _, ds := range dirSpecs ***REMOVED***
		memFs := NewMemMapFs()
		level1Fs := NewBasePathFs(memFs, ds.Dir1)
		level2Fs := NewBasePathFs(level1Fs, ds.Dir2)
		level3Fs := NewBasePathFs(level2Fs, ds.Dir3)

		type spec struct ***REMOVED***
			BaseFs       Fs
			FileName     string
			ExpectedPath string
		***REMOVED***
		specs := []spec***REMOVED***
			spec***REMOVED***BaseFs: level3Fs, FileName: "f.txt", ExpectedPath: filepath.Join(ds.Dir1, ds.Dir2, ds.Dir3, "f.txt")***REMOVED***,
			spec***REMOVED***BaseFs: level3Fs, FileName: "", ExpectedPath: filepath.Join(ds.Dir1, ds.Dir2, ds.Dir3, "")***REMOVED***,
			spec***REMOVED***BaseFs: level2Fs, FileName: "f.txt", ExpectedPath: filepath.Join(ds.Dir1, ds.Dir2, "f.txt")***REMOVED***,
			spec***REMOVED***BaseFs: level2Fs, FileName: "", ExpectedPath: filepath.Join(ds.Dir1, ds.Dir2, "")***REMOVED***,
			spec***REMOVED***BaseFs: level1Fs, FileName: "f.txt", ExpectedPath: filepath.Join(ds.Dir1, "f.txt")***REMOVED***,
			spec***REMOVED***BaseFs: level1Fs, FileName: "", ExpectedPath: filepath.Join(ds.Dir1, "")***REMOVED***,
		***REMOVED***

		for _, s := range specs ***REMOVED***
			if actualPath := FullBaseFsPath(s.BaseFs.(*BasePathFs), s.FileName); actualPath != s.ExpectedPath ***REMOVED***
				t.Errorf("Expected \n%s got \n%s", s.ExpectedPath, actualPath)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
