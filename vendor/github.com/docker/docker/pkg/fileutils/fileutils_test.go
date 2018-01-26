package fileutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CopyFile with invalid src
func TestCopyFileWithInvalidSrc(t *testing.T) ***REMOVED***
	tempFolder, err := ioutil.TempDir("", "docker-fileutils-test")
	defer os.RemoveAll(tempFolder)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	bytes, err := CopyFile("/invalid/file/path", path.Join(tempFolder, "dest"))
	if err == nil ***REMOVED***
		t.Fatal("Should have fail to copy an invalid src file")
	***REMOVED***
	if bytes != 0 ***REMOVED***
		t.Fatal("Should have written 0 bytes")
	***REMOVED***

***REMOVED***

// CopyFile with invalid dest
func TestCopyFileWithInvalidDest(t *testing.T) ***REMOVED***
	tempFolder, err := ioutil.TempDir("", "docker-fileutils-test")
	defer os.RemoveAll(tempFolder)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	src := path.Join(tempFolder, "file")
	err = ioutil.WriteFile(src, []byte("content"), 0740)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	bytes, err := CopyFile(src, path.Join(tempFolder, "/invalid/dest/path"))
	if err == nil ***REMOVED***
		t.Fatal("Should have fail to copy an invalid src file")
	***REMOVED***
	if bytes != 0 ***REMOVED***
		t.Fatal("Should have written 0 bytes")
	***REMOVED***

***REMOVED***

// CopyFile with same src and dest
func TestCopyFileWithSameSrcAndDest(t *testing.T) ***REMOVED***
	tempFolder, err := ioutil.TempDir("", "docker-fileutils-test")
	defer os.RemoveAll(tempFolder)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	file := path.Join(tempFolder, "file")
	err = ioutil.WriteFile(file, []byte("content"), 0740)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	bytes, err := CopyFile(file, file)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if bytes != 0 ***REMOVED***
		t.Fatal("Should have written 0 bytes as it is the same file.")
	***REMOVED***
***REMOVED***

// CopyFile with same src and dest but path is different and not clean
func TestCopyFileWithSameSrcAndDestWithPathNameDifferent(t *testing.T) ***REMOVED***
	tempFolder, err := ioutil.TempDir("", "docker-fileutils-test")
	defer os.RemoveAll(tempFolder)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	testFolder := path.Join(tempFolder, "test")
	err = os.MkdirAll(testFolder, 0740)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	file := path.Join(testFolder, "file")
	sameFile := testFolder + "/../test/file"
	err = ioutil.WriteFile(file, []byte("content"), 0740)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	bytes, err := CopyFile(file, sameFile)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if bytes != 0 ***REMOVED***
		t.Fatal("Should have written 0 bytes as it is the same file.")
	***REMOVED***
***REMOVED***

func TestCopyFile(t *testing.T) ***REMOVED***
	tempFolder, err := ioutil.TempDir("", "docker-fileutils-test")
	defer os.RemoveAll(tempFolder)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	src := path.Join(tempFolder, "src")
	dest := path.Join(tempFolder, "dest")
	ioutil.WriteFile(src, []byte("content"), 0777)
	ioutil.WriteFile(dest, []byte("destContent"), 0777)
	bytes, err := CopyFile(src, dest)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if bytes != 7 ***REMOVED***
		t.Fatalf("Should have written %d bytes but wrote %d", 7, bytes)
	***REMOVED***
	actual, err := ioutil.ReadFile(dest)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if string(actual) != "content" ***REMOVED***
		t.Fatalf("Dest content was '%s', expected '%s'", string(actual), "content")
	***REMOVED***
***REMOVED***

// Reading a symlink to a directory must return the directory
func TestReadSymlinkedDirectoryExistingDirectory(t *testing.T) ***REMOVED***
	// TODO Windows: Port this test
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Needs porting to Windows")
	***REMOVED***
	var err error
	if err = os.Mkdir("/tmp/testReadSymlinkToExistingDirectory", 0777); err != nil ***REMOVED***
		t.Errorf("failed to create directory: %s", err)
	***REMOVED***

	if err = os.Symlink("/tmp/testReadSymlinkToExistingDirectory", "/tmp/dirLinkTest"); err != nil ***REMOVED***
		t.Errorf("failed to create symlink: %s", err)
	***REMOVED***

	var path string
	if path, err = ReadSymlinkedDirectory("/tmp/dirLinkTest"); err != nil ***REMOVED***
		t.Fatalf("failed to read symlink to directory: %s", err)
	***REMOVED***

	if path != "/tmp/testReadSymlinkToExistingDirectory" ***REMOVED***
		t.Fatalf("symlink returned unexpected directory: %s", path)
	***REMOVED***

	if err = os.Remove("/tmp/testReadSymlinkToExistingDirectory"); err != nil ***REMOVED***
		t.Errorf("failed to remove temporary directory: %s", err)
	***REMOVED***

	if err = os.Remove("/tmp/dirLinkTest"); err != nil ***REMOVED***
		t.Errorf("failed to remove symlink: %s", err)
	***REMOVED***
***REMOVED***

// Reading a non-existing symlink must fail
func TestReadSymlinkedDirectoryNonExistingSymlink(t *testing.T) ***REMOVED***
	var path string
	var err error
	if path, err = ReadSymlinkedDirectory("/tmp/test/foo/Non/ExistingPath"); err == nil ***REMOVED***
		t.Fatalf("error expected for non-existing symlink")
	***REMOVED***

	if path != "" ***REMOVED***
		t.Fatalf("expected empty path, but '%s' was returned", path)
	***REMOVED***
***REMOVED***

// Reading a symlink to a file must fail
func TestReadSymlinkedDirectoryToFile(t *testing.T) ***REMOVED***
	// TODO Windows: Port this test
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Needs porting to Windows")
	***REMOVED***
	var err error
	var file *os.File

	if file, err = os.Create("/tmp/testReadSymlinkToFile"); err != nil ***REMOVED***
		t.Fatalf("failed to create file: %s", err)
	***REMOVED***

	file.Close()

	if err = os.Symlink("/tmp/testReadSymlinkToFile", "/tmp/fileLinkTest"); err != nil ***REMOVED***
		t.Errorf("failed to create symlink: %s", err)
	***REMOVED***

	var path string
	if path, err = ReadSymlinkedDirectory("/tmp/fileLinkTest"); err == nil ***REMOVED***
		t.Fatalf("ReadSymlinkedDirectory on a symlink to a file should've failed")
	***REMOVED***

	if path != "" ***REMOVED***
		t.Fatalf("path should've been empty: %s", path)
	***REMOVED***

	if err = os.Remove("/tmp/testReadSymlinkToFile"); err != nil ***REMOVED***
		t.Errorf("failed to remove file: %s", err)
	***REMOVED***

	if err = os.Remove("/tmp/fileLinkTest"); err != nil ***REMOVED***
		t.Errorf("failed to remove symlink: %s", err)
	***REMOVED***
***REMOVED***

func TestWildcardMatches(t *testing.T) ***REMOVED***
	match, _ := Matches("fileutils.go", []string***REMOVED***"*"***REMOVED***)
	if !match ***REMOVED***
		t.Errorf("failed to get a wildcard match, got %v", match)
	***REMOVED***
***REMOVED***

// A simple pattern match should return true.
func TestPatternMatches(t *testing.T) ***REMOVED***
	match, _ := Matches("fileutils.go", []string***REMOVED***"*.go"***REMOVED***)
	if !match ***REMOVED***
		t.Errorf("failed to get a match, got %v", match)
	***REMOVED***
***REMOVED***

// An exclusion followed by an inclusion should return true.
func TestExclusionPatternMatchesPatternBefore(t *testing.T) ***REMOVED***
	match, _ := Matches("fileutils.go", []string***REMOVED***"!fileutils.go", "*.go"***REMOVED***)
	if !match ***REMOVED***
		t.Errorf("failed to get true match on exclusion pattern, got %v", match)
	***REMOVED***
***REMOVED***

// A folder pattern followed by an exception should return false.
func TestPatternMatchesFolderExclusions(t *testing.T) ***REMOVED***
	match, _ := Matches("docs/README.md", []string***REMOVED***"docs", "!docs/README.md"***REMOVED***)
	if match ***REMOVED***
		t.Errorf("failed to get a false match on exclusion pattern, got %v", match)
	***REMOVED***
***REMOVED***

// A folder pattern followed by an exception should return false.
func TestPatternMatchesFolderWithSlashExclusions(t *testing.T) ***REMOVED***
	match, _ := Matches("docs/README.md", []string***REMOVED***"docs/", "!docs/README.md"***REMOVED***)
	if match ***REMOVED***
		t.Errorf("failed to get a false match on exclusion pattern, got %v", match)
	***REMOVED***
***REMOVED***

// A folder pattern followed by an exception should return false.
func TestPatternMatchesFolderWildcardExclusions(t *testing.T) ***REMOVED***
	match, _ := Matches("docs/README.md", []string***REMOVED***"docs/*", "!docs/README.md"***REMOVED***)
	if match ***REMOVED***
		t.Errorf("failed to get a false match on exclusion pattern, got %v", match)
	***REMOVED***
***REMOVED***

// A pattern followed by an exclusion should return false.
func TestExclusionPatternMatchesPatternAfter(t *testing.T) ***REMOVED***
	match, _ := Matches("fileutils.go", []string***REMOVED***"*.go", "!fileutils.go"***REMOVED***)
	if match ***REMOVED***
		t.Errorf("failed to get false match on exclusion pattern, got %v", match)
	***REMOVED***
***REMOVED***

// A filename evaluating to . should return false.
func TestExclusionPatternMatchesWholeDirectory(t *testing.T) ***REMOVED***
	match, _ := Matches(".", []string***REMOVED***"*.go"***REMOVED***)
	if match ***REMOVED***
		t.Errorf("failed to get false match on ., got %v", match)
	***REMOVED***
***REMOVED***

// A single ! pattern should return an error.
func TestSingleExclamationError(t *testing.T) ***REMOVED***
	_, err := Matches("fileutils.go", []string***REMOVED***"!"***REMOVED***)
	if err == nil ***REMOVED***
		t.Errorf("failed to get an error for a single exclamation point, got %v", err)
	***REMOVED***
***REMOVED***

// Matches with no patterns
func TestMatchesWithNoPatterns(t *testing.T) ***REMOVED***
	matches, err := Matches("/any/path/there", []string***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if matches ***REMOVED***
		t.Fatalf("Should not have match anything")
	***REMOVED***
***REMOVED***

// Matches with malformed patterns
func TestMatchesWithMalformedPatterns(t *testing.T) ***REMOVED***
	matches, err := Matches("/any/path/there", []string***REMOVED***"["***REMOVED***)
	if err == nil ***REMOVED***
		t.Fatal("Should have failed because of a malformed syntax in the pattern")
	***REMOVED***
	if matches ***REMOVED***
		t.Fatalf("Should not have match anything")
	***REMOVED***
***REMOVED***

type matchesTestCase struct ***REMOVED***
	pattern string
	text    string
	pass    bool
***REMOVED***

func TestMatches(t *testing.T) ***REMOVED***
	tests := []matchesTestCase***REMOVED***
		***REMOVED***"**", "file", true***REMOVED***,
		***REMOVED***"**", "file/", true***REMOVED***,
		***REMOVED***"**/", "file", true***REMOVED***, // weird one
		***REMOVED***"**/", "file/", true***REMOVED***,
		***REMOVED***"**", "/", true***REMOVED***,
		***REMOVED***"**/", "/", true***REMOVED***,
		***REMOVED***"**", "dir/file", true***REMOVED***,
		***REMOVED***"**/", "dir/file", true***REMOVED***,
		***REMOVED***"**", "dir/file/", true***REMOVED***,
		***REMOVED***"**/", "dir/file/", true***REMOVED***,
		***REMOVED***"**/**", "dir/file", true***REMOVED***,
		***REMOVED***"**/**", "dir/file/", true***REMOVED***,
		***REMOVED***"dir/**", "dir/file", true***REMOVED***,
		***REMOVED***"dir/**", "dir/file/", true***REMOVED***,
		***REMOVED***"dir/**", "dir/dir2/file", true***REMOVED***,
		***REMOVED***"dir/**", "dir/dir2/file/", true***REMOVED***,
		***REMOVED***"**/dir2/*", "dir/dir2/file", true***REMOVED***,
		***REMOVED***"**/dir2/*", "dir/dir2/file/", true***REMOVED***,
		***REMOVED***"**/dir2/**", "dir/dir2/dir3/file", true***REMOVED***,
		***REMOVED***"**/dir2/**", "dir/dir2/dir3/file/", true***REMOVED***,
		***REMOVED***"**file", "file", true***REMOVED***,
		***REMOVED***"**file", "dir/file", true***REMOVED***,
		***REMOVED***"**/file", "dir/file", true***REMOVED***,
		***REMOVED***"**file", "dir/dir/file", true***REMOVED***,
		***REMOVED***"**/file", "dir/dir/file", true***REMOVED***,
		***REMOVED***"**/file*", "dir/dir/file", true***REMOVED***,
		***REMOVED***"**/file*", "dir/dir/file.txt", true***REMOVED***,
		***REMOVED***"**/file*txt", "dir/dir/file.txt", true***REMOVED***,
		***REMOVED***"**/file*.txt", "dir/dir/file.txt", true***REMOVED***,
		***REMOVED***"**/file*.txt*", "dir/dir/file.txt", true***REMOVED***,
		***REMOVED***"**/**/*.txt", "dir/dir/file.txt", true***REMOVED***,
		***REMOVED***"**/**/*.txt2", "dir/dir/file.txt", false***REMOVED***,
		***REMOVED***"**/*.txt", "file.txt", true***REMOVED***,
		***REMOVED***"**/**/*.txt", "file.txt", true***REMOVED***,
		***REMOVED***"a**/*.txt", "a/file.txt", true***REMOVED***,
		***REMOVED***"a**/*.txt", "a/dir/file.txt", true***REMOVED***,
		***REMOVED***"a**/*.txt", "a/dir/dir/file.txt", true***REMOVED***,
		***REMOVED***"a/*.txt", "a/dir/file.txt", false***REMOVED***,
		***REMOVED***"a/*.txt", "a/file.txt", true***REMOVED***,
		***REMOVED***"a/*.txt**", "a/file.txt", true***REMOVED***,
		***REMOVED***"a[b-d]e", "ae", false***REMOVED***,
		***REMOVED***"a[b-d]e", "ace", true***REMOVED***,
		***REMOVED***"a[b-d]e", "aae", false***REMOVED***,
		***REMOVED***"a[^b-d]e", "aze", true***REMOVED***,
		***REMOVED***".*", ".foo", true***REMOVED***,
		***REMOVED***".*", "foo", false***REMOVED***,
		***REMOVED***"abc.def", "abcdef", false***REMOVED***,
		***REMOVED***"abc.def", "abc.def", true***REMOVED***,
		***REMOVED***"abc.def", "abcZdef", false***REMOVED***,
		***REMOVED***"abc?def", "abcZdef", true***REMOVED***,
		***REMOVED***"abc?def", "abcdef", false***REMOVED***,
		***REMOVED***"a\\\\", "a\\", true***REMOVED***,
		***REMOVED***"**/foo/bar", "foo/bar", true***REMOVED***,
		***REMOVED***"**/foo/bar", "dir/foo/bar", true***REMOVED***,
		***REMOVED***"**/foo/bar", "dir/dir2/foo/bar", true***REMOVED***,
		***REMOVED***"abc/**", "abc", false***REMOVED***,
		***REMOVED***"abc/**", "abc/def", true***REMOVED***,
		***REMOVED***"abc/**", "abc/def/ghi", true***REMOVED***,
		***REMOVED***"**/.foo", ".foo", true***REMOVED***,
		***REMOVED***"**/.foo", "bar.foo", false***REMOVED***,
	***REMOVED***

	if runtime.GOOS != "windows" ***REMOVED***
		tests = append(tests, []matchesTestCase***REMOVED***
			***REMOVED***"a\\*b", "a*b", true***REMOVED***,
			***REMOVED***"a\\", "a", false***REMOVED***,
			***REMOVED***"a\\", "a\\", false***REMOVED***,
		***REMOVED***...)
	***REMOVED***

	for _, test := range tests ***REMOVED***
		desc := fmt.Sprintf("pattern=%q text=%q", test.pattern, test.text)
		pm, err := NewPatternMatcher([]string***REMOVED***test.pattern***REMOVED***)
		require.NoError(t, err, desc)
		res, _ := pm.Matches(test.text)
		assert.Equal(t, test.pass, res, desc)
	***REMOVED***
***REMOVED***

func TestCleanPatterns(t *testing.T) ***REMOVED***
	patterns := []string***REMOVED***"docs", "config"***REMOVED***
	pm, err := NewPatternMatcher(patterns)
	if err != nil ***REMOVED***
		t.Fatalf("invalid pattern %v", patterns)
	***REMOVED***
	cleaned := pm.Patterns()
	if len(cleaned) != 2 ***REMOVED***
		t.Errorf("expected 2 element slice, got %v", len(cleaned))
	***REMOVED***
***REMOVED***

func TestCleanPatternsStripEmptyPatterns(t *testing.T) ***REMOVED***
	patterns := []string***REMOVED***"docs", "config", ""***REMOVED***
	pm, err := NewPatternMatcher(patterns)
	if err != nil ***REMOVED***
		t.Fatalf("invalid pattern %v", patterns)
	***REMOVED***
	cleaned := pm.Patterns()
	if len(cleaned) != 2 ***REMOVED***
		t.Errorf("expected 2 element slice, got %v", len(cleaned))
	***REMOVED***
***REMOVED***

func TestCleanPatternsExceptionFlag(t *testing.T) ***REMOVED***
	patterns := []string***REMOVED***"docs", "!docs/README.md"***REMOVED***
	pm, err := NewPatternMatcher(patterns)
	if err != nil ***REMOVED***
		t.Fatalf("invalid pattern %v", patterns)
	***REMOVED***
	if !pm.Exclusions() ***REMOVED***
		t.Errorf("expected exceptions to be true, got %v", pm.Exclusions())
	***REMOVED***
***REMOVED***

func TestCleanPatternsLeadingSpaceTrimmed(t *testing.T) ***REMOVED***
	patterns := []string***REMOVED***"docs", "  !docs/README.md"***REMOVED***
	pm, err := NewPatternMatcher(patterns)
	if err != nil ***REMOVED***
		t.Fatalf("invalid pattern %v", patterns)
	***REMOVED***
	if !pm.Exclusions() ***REMOVED***
		t.Errorf("expected exceptions to be true, got %v", pm.Exclusions())
	***REMOVED***
***REMOVED***

func TestCleanPatternsTrailingSpaceTrimmed(t *testing.T) ***REMOVED***
	patterns := []string***REMOVED***"docs", "!docs/README.md  "***REMOVED***
	pm, err := NewPatternMatcher(patterns)
	if err != nil ***REMOVED***
		t.Fatalf("invalid pattern %v", patterns)
	***REMOVED***
	if !pm.Exclusions() ***REMOVED***
		t.Errorf("expected exceptions to be true, got %v", pm.Exclusions())
	***REMOVED***
***REMOVED***

func TestCleanPatternsErrorSingleException(t *testing.T) ***REMOVED***
	patterns := []string***REMOVED***"!"***REMOVED***
	_, err := NewPatternMatcher(patterns)
	if err == nil ***REMOVED***
		t.Errorf("expected error on single exclamation point, got %v", err)
	***REMOVED***
***REMOVED***

func TestCreateIfNotExistsDir(t *testing.T) ***REMOVED***
	tempFolder, err := ioutil.TempDir("", "docker-fileutils-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tempFolder)

	folderToCreate := filepath.Join(tempFolder, "tocreate")

	if err := CreateIfNotExists(folderToCreate, true); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	fileinfo, err := os.Stat(folderToCreate)
	if err != nil ***REMOVED***
		t.Fatalf("Should have create a folder, got %v", err)
	***REMOVED***

	if !fileinfo.IsDir() ***REMOVED***
		t.Fatalf("Should have been a dir, seems it's not")
	***REMOVED***
***REMOVED***

func TestCreateIfNotExistsFile(t *testing.T) ***REMOVED***
	tempFolder, err := ioutil.TempDir("", "docker-fileutils-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tempFolder)

	fileToCreate := filepath.Join(tempFolder, "file/to/create")

	if err := CreateIfNotExists(fileToCreate, false); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	fileinfo, err := os.Stat(fileToCreate)
	if err != nil ***REMOVED***
		t.Fatalf("Should have create a file, got %v", err)
	***REMOVED***

	if fileinfo.IsDir() ***REMOVED***
		t.Fatalf("Should have been a file, seems it's not")
	***REMOVED***
***REMOVED***

// These matchTests are stolen from go's filepath Match tests.
type matchTest struct ***REMOVED***
	pattern, s string
	match      bool
	err        error
***REMOVED***

var matchTests = []matchTest***REMOVED***
	***REMOVED***"abc", "abc", true, nil***REMOVED***,
	***REMOVED***"*", "abc", true, nil***REMOVED***,
	***REMOVED***"*c", "abc", true, nil***REMOVED***,
	***REMOVED***"a*", "a", true, nil***REMOVED***,
	***REMOVED***"a*", "abc", true, nil***REMOVED***,
	***REMOVED***"a*", "ab/c", true, nil***REMOVED***,
	***REMOVED***"a*/b", "abc/b", true, nil***REMOVED***,
	***REMOVED***"a*/b", "a/c/b", false, nil***REMOVED***,
	***REMOVED***"a*b*c*d*e*/f", "axbxcxdxe/f", true, nil***REMOVED***,
	***REMOVED***"a*b*c*d*e*/f", "axbxcxdxexxx/f", true, nil***REMOVED***,
	***REMOVED***"a*b*c*d*e*/f", "axbxcxdxe/xxx/f", false, nil***REMOVED***,
	***REMOVED***"a*b*c*d*e*/f", "axbxcxdxexxx/fff", false, nil***REMOVED***,
	***REMOVED***"a*b?c*x", "abxbbxdbxebxczzx", true, nil***REMOVED***,
	***REMOVED***"a*b?c*x", "abxbbxdbxebxczzy", false, nil***REMOVED***,
	***REMOVED***"ab[c]", "abc", true, nil***REMOVED***,
	***REMOVED***"ab[b-d]", "abc", true, nil***REMOVED***,
	***REMOVED***"ab[e-g]", "abc", false, nil***REMOVED***,
	***REMOVED***"ab[^c]", "abc", false, nil***REMOVED***,
	***REMOVED***"ab[^b-d]", "abc", false, nil***REMOVED***,
	***REMOVED***"ab[^e-g]", "abc", true, nil***REMOVED***,
	***REMOVED***"a\\*b", "a*b", true, nil***REMOVED***,
	***REMOVED***"a\\*b", "ab", false, nil***REMOVED***,
	***REMOVED***"a?b", "a☺b", true, nil***REMOVED***,
	***REMOVED***"a[^a]b", "a☺b", true, nil***REMOVED***,
	***REMOVED***"a???b", "a☺b", false, nil***REMOVED***,
	***REMOVED***"a[^a][^a][^a]b", "a☺b", false, nil***REMOVED***,
	***REMOVED***"[a-ζ]*", "α", true, nil***REMOVED***,
	***REMOVED***"*[a-ζ]", "A", false, nil***REMOVED***,
	***REMOVED***"a?b", "a/b", false, nil***REMOVED***,
	***REMOVED***"a*b", "a/b", false, nil***REMOVED***,
	***REMOVED***"[\\]a]", "]", true, nil***REMOVED***,
	***REMOVED***"[\\-]", "-", true, nil***REMOVED***,
	***REMOVED***"[x\\-]", "x", true, nil***REMOVED***,
	***REMOVED***"[x\\-]", "-", true, nil***REMOVED***,
	***REMOVED***"[x\\-]", "z", false, nil***REMOVED***,
	***REMOVED***"[\\-x]", "x", true, nil***REMOVED***,
	***REMOVED***"[\\-x]", "-", true, nil***REMOVED***,
	***REMOVED***"[\\-x]", "a", false, nil***REMOVED***,
	***REMOVED***"[]a]", "]", false, filepath.ErrBadPattern***REMOVED***,
	***REMOVED***"[-]", "-", false, filepath.ErrBadPattern***REMOVED***,
	***REMOVED***"[x-]", "x", false, filepath.ErrBadPattern***REMOVED***,
	***REMOVED***"[x-]", "-", false, filepath.ErrBadPattern***REMOVED***,
	***REMOVED***"[x-]", "z", false, filepath.ErrBadPattern***REMOVED***,
	***REMOVED***"[-x]", "x", false, filepath.ErrBadPattern***REMOVED***,
	***REMOVED***"[-x]", "-", false, filepath.ErrBadPattern***REMOVED***,
	***REMOVED***"[-x]", "a", false, filepath.ErrBadPattern***REMOVED***,
	***REMOVED***"\\", "a", false, filepath.ErrBadPattern***REMOVED***,
	***REMOVED***"[a-b-c]", "a", false, filepath.ErrBadPattern***REMOVED***,
	***REMOVED***"[", "a", false, filepath.ErrBadPattern***REMOVED***,
	***REMOVED***"[^", "a", false, filepath.ErrBadPattern***REMOVED***,
	***REMOVED***"[^bc", "a", false, filepath.ErrBadPattern***REMOVED***,
	***REMOVED***"a[", "a", false, filepath.ErrBadPattern***REMOVED***, // was nil but IMO its wrong
	***REMOVED***"a[", "ab", false, filepath.ErrBadPattern***REMOVED***,
	***REMOVED***"*x", "xxx", true, nil***REMOVED***,
***REMOVED***

func errp(e error) string ***REMOVED***
	if e == nil ***REMOVED***
		return "<nil>"
	***REMOVED***
	return e.Error()
***REMOVED***

// TestMatch test's our version of filepath.Match, called regexpMatch.
func TestMatch(t *testing.T) ***REMOVED***
	for _, tt := range matchTests ***REMOVED***
		pattern := tt.pattern
		s := tt.s
		if runtime.GOOS == "windows" ***REMOVED***
			if strings.Contains(pattern, "\\") ***REMOVED***
				// no escape allowed on windows.
				continue
			***REMOVED***
			pattern = filepath.Clean(pattern)
			s = filepath.Clean(s)
		***REMOVED***
		ok, err := Matches(s, []string***REMOVED***pattern***REMOVED***)
		if ok != tt.match || err != tt.err ***REMOVED***
			t.Fatalf("Match(%#q, %#q) = %v, %q want %v, %q", pattern, s, ok, errp(err), tt.match, errp(tt.err))
		***REMOVED***
	***REMOVED***
***REMOVED***
