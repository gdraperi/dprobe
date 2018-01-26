package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDir = "testfiles"
const negativeTestDir = "testfiles-negative"
const testFileLineInfo = "testfile-line/Dockerfile"

func getDirs(t *testing.T, dir string) []string ***REMOVED***
	f, err := os.Open(dir)
	require.NoError(t, err)
	defer f.Close()

	dirs, err := f.Readdirnames(0)
	require.NoError(t, err)
	return dirs
***REMOVED***

func TestParseErrorCases(t *testing.T) ***REMOVED***
	for _, dir := range getDirs(t, negativeTestDir) ***REMOVED***
		dockerfile := filepath.Join(negativeTestDir, dir, "Dockerfile")

		df, err := os.Open(dockerfile)
		require.NoError(t, err, dockerfile)
		defer df.Close()

		_, err = Parse(df)
		assert.Error(t, err, dockerfile)
	***REMOVED***
***REMOVED***

func TestParseCases(t *testing.T) ***REMOVED***
	for _, dir := range getDirs(t, testDir) ***REMOVED***
		dockerfile := filepath.Join(testDir, dir, "Dockerfile")
		resultfile := filepath.Join(testDir, dir, "result")

		df, err := os.Open(dockerfile)
		require.NoError(t, err, dockerfile)
		defer df.Close()

		result, err := Parse(df)
		require.NoError(t, err, dockerfile)

		content, err := ioutil.ReadFile(resultfile)
		require.NoError(t, err, resultfile)

		if runtime.GOOS == "windows" ***REMOVED***
			// CRLF --> CR to match Unix behavior
			content = bytes.Replace(content, []byte***REMOVED***'\x0d', '\x0a'***REMOVED***, []byte***REMOVED***'\x0a'***REMOVED***, -1)
		***REMOVED***
		assert.Equal(t, result.AST.Dump()+"\n", string(content), "In "+dockerfile)
	***REMOVED***
***REMOVED***

func TestParseWords(t *testing.T) ***REMOVED***
	tests := []map[string][]string***REMOVED***
		***REMOVED***
			"input":  ***REMOVED***"foo"***REMOVED***,
			"expect": ***REMOVED***"foo"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"input":  ***REMOVED***"foo bar"***REMOVED***,
			"expect": ***REMOVED***"foo", "bar"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"input":  ***REMOVED***"foo\\ bar"***REMOVED***,
			"expect": ***REMOVED***"foo\\ bar"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"input":  ***REMOVED***"foo=bar"***REMOVED***,
			"expect": ***REMOVED***"foo=bar"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"input":  ***REMOVED***"foo bar 'abc xyz'"***REMOVED***,
			"expect": ***REMOVED***"foo", "bar", "'abc xyz'"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"input":  ***REMOVED***`foo bar "abc xyz"`***REMOVED***,
			"expect": ***REMOVED***"foo", "bar", `"abc xyz"`***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"input":  ***REMOVED***"àöû"***REMOVED***,
			"expect": ***REMOVED***"àöû"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"input":  ***REMOVED***`föo bàr "âbc xÿz"`***REMOVED***,
			"expect": ***REMOVED***"föo", "bàr", `"âbc xÿz"`***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for _, test := range tests ***REMOVED***
		words := parseWords(test["input"][0], NewDefaultDirective())
		assert.Equal(t, test["expect"], words)
	***REMOVED***
***REMOVED***

func TestParseIncludesLineNumbers(t *testing.T) ***REMOVED***
	df, err := os.Open(testFileLineInfo)
	require.NoError(t, err)
	defer df.Close()

	result, err := Parse(df)
	require.NoError(t, err)

	ast := result.AST
	assert.Equal(t, 5, ast.StartLine)
	assert.Equal(t, 31, ast.endLine)
	assert.Len(t, ast.Children, 3)
	expected := [][]int***REMOVED***
		***REMOVED***5, 5***REMOVED***,
		***REMOVED***11, 12***REMOVED***,
		***REMOVED***17, 31***REMOVED***,
	***REMOVED***
	for i, child := range ast.Children ***REMOVED***
		msg := fmt.Sprintf("Child %d", i)
		assert.Equal(t, expected[i], []int***REMOVED***child.StartLine, child.endLine***REMOVED***, msg)
	***REMOVED***
***REMOVED***

func TestParseWarnsOnEmptyContinutationLine(t *testing.T) ***REMOVED***
	dockerfile := bytes.NewBufferString(`
FROM alpine:3.6

RUN something \

    following \

    more

RUN another \

    thing
RUN non-indented \
# this is a comment
   after-comment

RUN indented \
    # this is an indented comment
    comment
	`)

	result, err := Parse(dockerfile)
	require.NoError(t, err)
	warnings := result.Warnings
	assert.Len(t, warnings, 3)
	assert.Contains(t, warnings[0], "Empty continuation line found in")
	assert.Contains(t, warnings[0], "RUN something     following     more")
	assert.Contains(t, warnings[1], "RUN another     thing")
	assert.Contains(t, warnings[2], "will become errors in a future release")
***REMOVED***

func TestParseReturnsScannerErrors(t *testing.T) ***REMOVED***
	label := strings.Repeat("a", bufio.MaxScanTokenSize)

	dockerfile := strings.NewReader(fmt.Sprintf(`
		FROM image
		LABEL test=%s
`, label))
	_, err := Parse(dockerfile)
	assert.EqualError(t, err, "dockerfile line greater than max allowed size of 65535")
***REMOVED***
