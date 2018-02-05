package doc

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestGenYamlDoc(t *testing.T) ***REMOVED***
	// We generate on s subcommand so we have both subcommands and parents
	buf := new(bytes.Buffer)
	if err := GenYaml(echoCmd, buf); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	output := buf.String()

	checkStringContains(t, output, echoCmd.Long)
	checkStringContains(t, output, echoCmd.Example)
	checkStringContains(t, output, "boolone")
	checkStringContains(t, output, "rootflag")
	checkStringContains(t, output, rootCmd.Short)
	checkStringContains(t, output, echoSubCmd.Short)
***REMOVED***

func TestGenYamlNoTag(t *testing.T) ***REMOVED***
	rootCmd.DisableAutoGenTag = true
	defer func() ***REMOVED*** rootCmd.DisableAutoGenTag = false ***REMOVED***()

	buf := new(bytes.Buffer)
	if err := GenYaml(rootCmd, buf); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	output := buf.String()

	checkStringOmits(t, output, "Auto generated")
***REMOVED***

func TestGenYamlTree(t *testing.T) ***REMOVED***
	c := &cobra.Command***REMOVED***Use: "do [OPTIONS] arg1 arg2"***REMOVED***

	tmpdir, err := ioutil.TempDir("", "test-gen-yaml-tree")
	if err != nil ***REMOVED***
		t.Fatalf("Failed to create tmpdir: %s", err.Error())
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	if err := GenYamlTree(c, tmpdir); err != nil ***REMOVED***
		t.Fatalf("GenYamlTree failed: %s", err.Error())
	***REMOVED***

	if _, err := os.Stat(filepath.Join(tmpdir, "do.yaml")); err != nil ***REMOVED***
		t.Fatalf("Expected file 'do.yaml' to exist")
	***REMOVED***
***REMOVED***

func BenchmarkGenYamlToFile(b *testing.B) ***REMOVED***
	file, err := ioutil.TempFile("", "")
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	defer os.Remove(file.Name())
	defer file.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if err := GenYaml(rootCmd, file); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***
