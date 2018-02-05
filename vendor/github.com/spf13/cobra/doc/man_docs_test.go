package doc

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func translate(in string) string ***REMOVED***
	return strings.Replace(in, "-", "\\-", -1)
***REMOVED***

func TestGenManDoc(t *testing.T) ***REMOVED***
	header := &GenManHeader***REMOVED***
		Title:   "Project",
		Section: "2",
	***REMOVED***

	// We generate on a subcommand so we have both subcommands and parents
	buf := new(bytes.Buffer)
	if err := GenMan(echoCmd, header, buf); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	output := buf.String()

	// Make sure parent has - in CommandPath() in SEE ALSO:
	parentPath := echoCmd.Parent().CommandPath()
	dashParentPath := strings.Replace(parentPath, " ", "-", -1)
	expected := translate(dashParentPath)
	expected = expected + "(" + header.Section + ")"
	checkStringContains(t, output, expected)

	checkStringContains(t, output, translate(echoCmd.Name()))
	checkStringContains(t, output, translate(echoCmd.Name()))
	checkStringContains(t, output, "boolone")
	checkStringContains(t, output, "rootflag")
	checkStringContains(t, output, translate(rootCmd.Name()))
	checkStringContains(t, output, translate(echoSubCmd.Name()))
	checkStringOmits(t, output, translate(deprecatedCmd.Name()))
	checkStringContains(t, output, translate("Auto generated"))
***REMOVED***

func TestGenManNoGenTag(t *testing.T) ***REMOVED***
	echoCmd.DisableAutoGenTag = true
	defer func() ***REMOVED*** echoCmd.DisableAutoGenTag = false ***REMOVED***()

	header := &GenManHeader***REMOVED***
		Title:   "Project",
		Section: "2",
	***REMOVED***

	// We generate on a subcommand so we have both subcommands and parents
	buf := new(bytes.Buffer)
	if err := GenMan(echoCmd, header, buf); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	output := buf.String()

	unexpected := translate("#HISTORY")
	checkStringOmits(t, output, unexpected)
***REMOVED***

func TestGenManSeeAlso(t *testing.T) ***REMOVED***
	rootCmd := &cobra.Command***REMOVED***Use: "root", Run: emptyRun***REMOVED***
	aCmd := &cobra.Command***REMOVED***Use: "aaa", Run: emptyRun, Hidden: true***REMOVED*** // #229
	bCmd := &cobra.Command***REMOVED***Use: "bbb", Run: emptyRun***REMOVED***
	cCmd := &cobra.Command***REMOVED***Use: "ccc", Run: emptyRun***REMOVED***
	rootCmd.AddCommand(aCmd, bCmd, cCmd)

	buf := new(bytes.Buffer)
	header := &GenManHeader***REMOVED******REMOVED***
	if err := GenMan(rootCmd, header, buf); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	scanner := bufio.NewScanner(buf)

	if err := assertLineFound(scanner, ".SH SEE ALSO"); err != nil ***REMOVED***
		t.Fatalf("Couldn't find SEE ALSO section header: %v", err)
	***REMOVED***
	if err := assertNextLineEquals(scanner, ".PP"); err != nil ***REMOVED***
		t.Fatalf("First line after SEE ALSO wasn't break-indent: %v", err)
	***REMOVED***
	if err := assertNextLineEquals(scanner, `\fBroot\-bbb(1)\fP, \fBroot\-ccc(1)\fP`); err != nil ***REMOVED***
		t.Fatalf("Second line after SEE ALSO wasn't correct: %v", err)
	***REMOVED***
***REMOVED***

func TestManPrintFlagsHidesShortDeperecated(t *testing.T) ***REMOVED***
	c := &cobra.Command***REMOVED******REMOVED***
	c.Flags().StringP("foo", "f", "default", "Foo flag")
	c.Flags().MarkShorthandDeprecated("foo", "don't use it no more")

	buf := new(bytes.Buffer)
	manPrintFlags(buf, c.Flags())

	got := buf.String()
	expected := "**--foo**=\"default\"\n\tFoo flag\n\n"
	if got != expected ***REMOVED***
		t.Errorf("Expected %v, got %v", expected, got)
	***REMOVED***
***REMOVED***

func TestGenManTree(t *testing.T) ***REMOVED***
	c := &cobra.Command***REMOVED***Use: "do [OPTIONS] arg1 arg2"***REMOVED***
	header := &GenManHeader***REMOVED***Section: "2"***REMOVED***
	tmpdir, err := ioutil.TempDir("", "test-gen-man-tree")
	if err != nil ***REMOVED***
		t.Fatalf("Failed to create tmpdir: %s", err.Error())
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	if err := GenManTree(c, header, tmpdir); err != nil ***REMOVED***
		t.Fatalf("GenManTree failed: %s", err.Error())
	***REMOVED***

	if _, err := os.Stat(filepath.Join(tmpdir, "do.2")); err != nil ***REMOVED***
		t.Fatalf("Expected file 'do.2' to exist")
	***REMOVED***

	if header.Title != "" ***REMOVED***
		t.Fatalf("Expected header.Title to be unmodified")
	***REMOVED***
***REMOVED***

func assertLineFound(scanner *bufio.Scanner, expectedLine string) error ***REMOVED***
	for scanner.Scan() ***REMOVED***
		line := scanner.Text()
		if line == expectedLine ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	if err := scanner.Err(); err != nil ***REMOVED***
		return fmt.Errorf("scan failed: %s", err)
	***REMOVED***

	return fmt.Errorf("hit EOF before finding %v", expectedLine)
***REMOVED***

func assertNextLineEquals(scanner *bufio.Scanner, expectedLine string) error ***REMOVED***
	if scanner.Scan() ***REMOVED***
		line := scanner.Text()
		if line == expectedLine ***REMOVED***
			return nil
		***REMOVED***
		return fmt.Errorf("got %v, not %v", line, expectedLine)
	***REMOVED***

	if err := scanner.Err(); err != nil ***REMOVED***
		return fmt.Errorf("scan failed: %v", err)
	***REMOVED***

	return fmt.Errorf("hit EOF before finding %v", expectedLine)
***REMOVED***

func BenchmarkGenManToFile(b *testing.B) ***REMOVED***
	file, err := ioutil.TempFile("", "")
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	defer os.Remove(file.Name())
	defer file.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if err := GenMan(rootCmd, nil, file); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***
