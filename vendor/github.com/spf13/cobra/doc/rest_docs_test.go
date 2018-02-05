package doc

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestGenRSTDoc(t *testing.T) ***REMOVED***
	// We generate on a subcommand so we have both subcommands and parents
	buf := new(bytes.Buffer)
	if err := GenReST(echoCmd, buf); err != nil ***REMOVED***
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

func TestGenRSTNoTag(t *testing.T) ***REMOVED***
	rootCmd.DisableAutoGenTag = true
	defer func() ***REMOVED*** rootCmd.DisableAutoGenTag = false ***REMOVED***()

	buf := new(bytes.Buffer)
	if err := GenReST(rootCmd, buf); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	output := buf.String()

	unexpected := "Auto generated"
	checkStringOmits(t, output, unexpected)
***REMOVED***

func TestGenRSTTree(t *testing.T) ***REMOVED***
	c := &cobra.Command***REMOVED***Use: "do [OPTIONS] arg1 arg2"***REMOVED***

	tmpdir, err := ioutil.TempDir("", "test-gen-rst-tree")
	if err != nil ***REMOVED***
		t.Fatalf("Failed to create tmpdir: %s", err.Error())
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	if err := GenReSTTree(c, tmpdir); err != nil ***REMOVED***
		t.Fatalf("GenReSTTree failed: %s", err.Error())
	***REMOVED***

	if _, err := os.Stat(filepath.Join(tmpdir, "do.rst")); err != nil ***REMOVED***
		t.Fatalf("Expected file 'do.rst' to exist")
	***REMOVED***
***REMOVED***

func BenchmarkGenReSTToFile(b *testing.B) ***REMOVED***
	file, err := ioutil.TempFile("", "")
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	defer os.Remove(file.Name())
	defer file.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if err := GenReST(rootCmd, file); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***
