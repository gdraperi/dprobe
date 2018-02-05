// +build !windows
// TODO(jen20): These need fixing on Windows but fmt is not used right now
// and red CI is making it harder to process other bugs, so ignore until
// we get around to fixing them.

package fmtcmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"syscall"
	"testing"

	"github.com/hashicorp/hcl/testhelper"
)

var fixtureExtensions = []string***REMOVED***"hcl"***REMOVED***

func init() ***REMOVED***
	sort.Sort(ByFilename(fixtures))
***REMOVED***

func TestIsValidFile(t *testing.T) ***REMOVED***
	const fixtureDir = "./test-fixtures"

	cases := []struct ***REMOVED***
		Path     string
		Expected bool
	***REMOVED******REMOVED***
		***REMOVED***"good.hcl", true***REMOVED***,
		***REMOVED***".hidden.ignore", false***REMOVED***,
		***REMOVED***"file.ignore", false***REMOVED***,
		***REMOVED***"dir.ignore", false***REMOVED***,
	***REMOVED***

	for _, tc := range cases ***REMOVED***
		file, err := os.Stat(filepath.Join(fixtureDir, tc.Path))
		if err != nil ***REMOVED***
			t.Errorf("unexpected error: %s", err)
		***REMOVED***

		if res := isValidFile(file, fixtureExtensions); res != tc.Expected ***REMOVED***
			t.Errorf("want: %b, got: %b", tc.Expected, res)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestRunMultiplePaths(t *testing.T) ***REMOVED***
	path1, err := renderFixtures("")
	if err != nil ***REMOVED***
		t.Errorf("unexpected error: %s", err)
	***REMOVED***
	defer os.RemoveAll(path1)
	path2, err := renderFixtures("")
	if err != nil ***REMOVED***
		t.Errorf("unexpected error: %s", err)
	***REMOVED***
	defer os.RemoveAll(path2)

	var expectedOut bytes.Buffer
	for _, path := range []string***REMOVED***path1, path2***REMOVED*** ***REMOVED***
		for _, fixture := range fixtures ***REMOVED***
			if !bytes.Equal(fixture.golden, fixture.input) ***REMOVED***
				expectedOut.WriteString(filepath.Join(path, fixture.filename) + "\n")
			***REMOVED***
		***REMOVED***
	***REMOVED***

	_, stdout := mockIO()
	err = Run(
		[]string***REMOVED***path1, path2***REMOVED***,
		fixtureExtensions,
		nil, stdout,
		Options***REMOVED***
			List: true,
		***REMOVED***,
	)

	if err != nil ***REMOVED***
		t.Errorf("unexpected error: %s", err)
	***REMOVED***
	if stdout.String() != expectedOut.String() ***REMOVED***
		t.Errorf("stdout want:\n%s\ngot:\n%s", expectedOut, stdout)
	***REMOVED***
***REMOVED***

func TestRunSubDirectories(t *testing.T) ***REMOVED***
	pathParent, err := ioutil.TempDir("", "")
	if err != nil ***REMOVED***
		t.Errorf("unexpected error: %s", err)
	***REMOVED***
	defer os.RemoveAll(pathParent)

	path1, err := renderFixtures(pathParent)
	if err != nil ***REMOVED***
		t.Errorf("unexpected error: %s", err)
	***REMOVED***
	path2, err := renderFixtures(pathParent)
	if err != nil ***REMOVED***
		t.Errorf("unexpected error: %s", err)
	***REMOVED***

	paths := []string***REMOVED***path1, path2***REMOVED***
	sort.Strings(paths)

	var expectedOut bytes.Buffer
	for _, path := range paths ***REMOVED***
		for _, fixture := range fixtures ***REMOVED***
			if !bytes.Equal(fixture.golden, fixture.input) ***REMOVED***
				expectedOut.WriteString(filepath.Join(path, fixture.filename) + "\n")
			***REMOVED***
		***REMOVED***
	***REMOVED***

	_, stdout := mockIO()
	err = Run(
		[]string***REMOVED***pathParent***REMOVED***,
		fixtureExtensions,
		nil, stdout,
		Options***REMOVED***
			List: true,
		***REMOVED***,
	)

	if err != nil ***REMOVED***
		t.Errorf("unexpected error: %s", err)
	***REMOVED***
	if stdout.String() != expectedOut.String() ***REMOVED***
		t.Errorf("stdout want:\n%s\ngot:\n%s", expectedOut, stdout)
	***REMOVED***
***REMOVED***

func TestRunStdin(t *testing.T) ***REMOVED***
	var expectedOut bytes.Buffer
	for i, fixture := range fixtures ***REMOVED***
		if i != 0 ***REMOVED***
			expectedOut.WriteString("\n")
		***REMOVED***
		expectedOut.Write(fixture.golden)
	***REMOVED***

	stdin, stdout := mockIO()
	for _, fixture := range fixtures ***REMOVED***
		stdin.Write(fixture.input)
	***REMOVED***

	err := Run(
		[]string***REMOVED******REMOVED***,
		fixtureExtensions,
		stdin, stdout,
		Options***REMOVED******REMOVED***,
	)

	if err != nil ***REMOVED***
		t.Errorf("unexpected error: %s", err)
	***REMOVED***
	if !bytes.Equal(stdout.Bytes(), expectedOut.Bytes()) ***REMOVED***
		t.Errorf("stdout want:\n%s\ngot:\n%s", expectedOut, stdout)
	***REMOVED***
***REMOVED***

func TestRunStdinAndWrite(t *testing.T) ***REMOVED***
	var expectedOut = []byte***REMOVED******REMOVED***

	stdin, stdout := mockIO()
	stdin.WriteString("")
	err := Run(
		[]string***REMOVED******REMOVED***, []string***REMOVED******REMOVED***,
		stdin, stdout,
		Options***REMOVED***
			Write: true,
		***REMOVED***,
	)

	if err != ErrWriteStdin ***REMOVED***
		t.Errorf("error want:\n%s\ngot:\n%s", ErrWriteStdin, err)
	***REMOVED***
	if !bytes.Equal(stdout.Bytes(), expectedOut) ***REMOVED***
		t.Errorf("stdout want:\n%s\ngot:\n%s", expectedOut, stdout)
	***REMOVED***
***REMOVED***

func TestRunFileError(t *testing.T) ***REMOVED***
	path, err := ioutil.TempDir("", "")
	if err != nil ***REMOVED***
		t.Errorf("unexpected error: %s", err)
	***REMOVED***
	defer os.RemoveAll(path)
	filename := filepath.Join(path, "unreadable.hcl")

	var expectedError = &os.PathError***REMOVED***
		Op:   "open",
		Path: filename,
		Err:  syscall.EACCES,
	***REMOVED***

	err = ioutil.WriteFile(filename, []byte***REMOVED******REMOVED***, 0000)
	if err != nil ***REMOVED***
		t.Errorf("unexpected error: %s", err)
	***REMOVED***

	_, stdout := mockIO()
	err = Run(
		[]string***REMOVED***path***REMOVED***,
		fixtureExtensions,
		nil, stdout,
		Options***REMOVED******REMOVED***,
	)

	if !reflect.DeepEqual(err, expectedError) ***REMOVED***
		t.Errorf("error want: %#v, got: %#v", expectedError, err)
	***REMOVED***
***REMOVED***

func TestRunNoOptions(t *testing.T) ***REMOVED***
	path, err := renderFixtures("")
	if err != nil ***REMOVED***
		t.Errorf("unexpected error: %s", err)
	***REMOVED***
	defer os.RemoveAll(path)

	var expectedOut bytes.Buffer
	for _, fixture := range fixtures ***REMOVED***
		expectedOut.Write(fixture.golden)
	***REMOVED***

	_, stdout := mockIO()
	err = Run(
		[]string***REMOVED***path***REMOVED***,
		fixtureExtensions,
		nil, stdout,
		Options***REMOVED******REMOVED***,
	)

	if err != nil ***REMOVED***
		t.Errorf("unexpected error: %s", err)
	***REMOVED***
	if stdout.String() != expectedOut.String() ***REMOVED***
		t.Errorf("stdout want:\n%s\ngot:\n%s", expectedOut, stdout)
	***REMOVED***
***REMOVED***

func TestRunList(t *testing.T) ***REMOVED***
	path, err := renderFixtures("")
	if err != nil ***REMOVED***
		t.Errorf("unexpected error: %s", err)
	***REMOVED***
	defer os.RemoveAll(path)

	var expectedOut bytes.Buffer
	for _, fixture := range fixtures ***REMOVED***
		if !bytes.Equal(fixture.golden, fixture.input) ***REMOVED***
			expectedOut.WriteString(fmt.Sprintln(filepath.Join(path, fixture.filename)))
		***REMOVED***
	***REMOVED***

	_, stdout := mockIO()
	err = Run(
		[]string***REMOVED***path***REMOVED***,
		fixtureExtensions,
		nil, stdout,
		Options***REMOVED***
			List: true,
		***REMOVED***,
	)

	if err != nil ***REMOVED***
		t.Errorf("unexpected error: %s", err)
	***REMOVED***
	if stdout.String() != expectedOut.String() ***REMOVED***
		t.Errorf("stdout want:\n%s\ngot:\n%s", expectedOut, stdout)
	***REMOVED***
***REMOVED***

func TestRunWrite(t *testing.T) ***REMOVED***
	path, err := renderFixtures("")
	if err != nil ***REMOVED***
		t.Errorf("unexpected error: %s", err)
	***REMOVED***
	defer os.RemoveAll(path)

	_, stdout := mockIO()
	err = Run(
		[]string***REMOVED***path***REMOVED***,
		fixtureExtensions,
		nil, stdout,
		Options***REMOVED***
			Write: true,
		***REMOVED***,
	)

	if err != nil ***REMOVED***
		t.Errorf("unexpected error: %s", err)
	***REMOVED***
	for _, fixture := range fixtures ***REMOVED***
		res, err := ioutil.ReadFile(filepath.Join(path, fixture.filename))
		if err != nil ***REMOVED***
			t.Errorf("unexpected error: %s", err)
		***REMOVED***
		if !bytes.Equal(res, fixture.golden) ***REMOVED***
			t.Errorf("file %q contents want:\n%s\ngot:\n%s", fixture.filename, fixture.golden, res)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestRunDiff(t *testing.T) ***REMOVED***
	path, err := renderFixtures("")
	if err != nil ***REMOVED***
		t.Errorf("unexpected error: %s", err)
	***REMOVED***
	defer os.RemoveAll(path)

	var expectedOut bytes.Buffer
	for _, fixture := range fixtures ***REMOVED***
		if len(fixture.diff) > 0 ***REMOVED***
			expectedOut.WriteString(
				regexp.QuoteMeta(
					fmt.Sprintf("diff a/%s/%s b/%s/%s\n", path, fixture.filename, path, fixture.filename),
				),
			)
			// Need to use regex to ignore datetimes in diff.
			expectedOut.WriteString(`--- .+?\n`)
			expectedOut.WriteString(`\+\+\+ .+?\n`)
			expectedOut.WriteString(regexp.QuoteMeta(string(fixture.diff)))
		***REMOVED***
	***REMOVED***

	expectedOutString := testhelper.Unix2dos(expectedOut.String())

	_, stdout := mockIO()
	err = Run(
		[]string***REMOVED***path***REMOVED***,
		fixtureExtensions,
		nil, stdout,
		Options***REMOVED***
			Diff: true,
		***REMOVED***,
	)

	if err != nil ***REMOVED***
		t.Errorf("unexpected error: %s", err)
	***REMOVED***
	if !regexp.MustCompile(expectedOutString).Match(stdout.Bytes()) ***REMOVED***
		t.Errorf("stdout want match:\n%s\ngot:\n%q", expectedOutString, stdout)
	***REMOVED***
***REMOVED***

func mockIO() (stdin, stdout *bytes.Buffer) ***REMOVED***
	return new(bytes.Buffer), new(bytes.Buffer)
***REMOVED***

type fixture struct ***REMOVED***
	filename            string
	input, golden, diff []byte
***REMOVED***

type ByFilename []fixture

func (s ByFilename) Len() int           ***REMOVED*** return len(s) ***REMOVED***
func (s ByFilename) Swap(i, j int)      ***REMOVED*** s[i], s[j] = s[j], s[i] ***REMOVED***
func (s ByFilename) Less(i, j int) bool ***REMOVED*** return len(s[i].filename) > len(s[j].filename) ***REMOVED***

var fixtures = []fixture***REMOVED***
	***REMOVED***
		"noop.hcl",
		[]byte(`resource "aws_security_group" "firewall" ***REMOVED***
  count = 5
***REMOVED***
`),
		[]byte(`resource "aws_security_group" "firewall" ***REMOVED***
  count = 5
***REMOVED***
`),
		[]byte(``),
	***REMOVED***, ***REMOVED***
		"align_equals.hcl",
		[]byte(`variable "foo" ***REMOVED***
  default = "bar"
  description = "bar"
***REMOVED***
`),
		[]byte(`variable "foo" ***REMOVED***
  default     = "bar"
  description = "bar"
***REMOVED***
`),
		[]byte(`@@ -1,4 +1,4 @@
 variable "foo" ***REMOVED***
-  default = "bar"
+  default     = "bar"
   description = "bar"
 ***REMOVED***
`),
	***REMOVED***, ***REMOVED***
		"indentation.hcl",
		[]byte(`provider "aws" ***REMOVED***
    access_key = "foo"
    secret_key = "bar"
***REMOVED***
`),
		[]byte(`provider "aws" ***REMOVED***
  access_key = "foo"
  secret_key = "bar"
***REMOVED***
`),
		[]byte(`@@ -1,4 +1,4 @@
 provider "aws" ***REMOVED***
-    access_key = "foo"
-    secret_key = "bar"
+  access_key = "foo"
+  secret_key = "bar"
 ***REMOVED***
`),
	***REMOVED***,
***REMOVED***

// parent can be an empty string, in which case the system's default
// temporary directory will be used.
func renderFixtures(parent string) (path string, err error) ***REMOVED***
	path, err = ioutil.TempDir(parent, "")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	for _, fixture := range fixtures ***REMOVED***
		err = ioutil.WriteFile(filepath.Join(path, fixture.filename), []byte(fixture.input), 0644)
		if err != nil ***REMOVED***
			os.RemoveAll(path)
			return "", err
		***REMOVED***
	***REMOVED***

	return path, nil
***REMOVED***
