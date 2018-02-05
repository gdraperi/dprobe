package doc

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestGenMdDoc(t *testing.T) ***REMOVED***
	// We generate on subcommand so we have both subcommands and parents.
	buf := new(bytes.Buffer)
	if err := GenMarkdown(echoCmd, buf); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	output := buf.String()

	checkStringContains(t, output, echoCmd.Long)
	checkStringContains(t, output, echoCmd.Example)
	checkStringContains(t, output, "boolone")
	checkStringContains(t, output, "rootflag")
	checkStringContains(t, output, rootCmd.Short)
	checkStringContains(t, output, echoSubCmd.Short)
	checkStringOmits(t, output, deprecatedCmd.Short)
***REMOVED***

func TestGenMdNoTag(t *testing.T) ***REMOVED***
	rootCmd.DisableAutoGenTag = true
	defer func() ***REMOVED*** rootCmd.DisableAutoGenTag = false ***REMOVED***()

	buf := new(bytes.Buffer)
	if err := GenMarkdown(rootCmd, buf); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	output := buf.String()

	checkStringOmits(t, output, "Auto generated")
***REMOVED***

func TestGenMdTree(t *testing.T) ***REMOVED***
	c := &cobra.Command***REMOVED***Use: "do [OPTIONS] arg1 arg2"***REMOVED***
	tmpdir, err := ioutil.TempDir("", "test-gen-md-tree")
	if err != nil ***REMOVED***
		t.Fatalf("Failed to create tmpdir: %v", err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	if err := GenMarkdownTree(c, tmpdir); err != nil ***REMOVED***
		t.Fatalf("GenMarkdownTree failed: %v", err)
	***REMOVED***

	if _, err := os.Stat(filepath.Join(tmpdir, "do.md")); err != nil ***REMOVED***
		t.Fatalf("Expected file 'do.md' to exist")
	***REMOVED***
***REMOVED***

func BenchmarkGenMarkdownToFile(b *testing.B) ***REMOVED***
	file, err := ioutil.TempFile("", "")
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	defer os.Remove(file.Name())
	defer file.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if err := GenMarkdown(rootCmd, file); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***
