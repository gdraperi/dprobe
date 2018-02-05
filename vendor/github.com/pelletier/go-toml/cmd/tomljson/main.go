// Tomljson reads TOML and converts to JSON.
//
// Usage:
//   cat file.toml | tomljson > file.json
//   tomljson file1.toml > file.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/pelletier/go-toml"
)

func main() ***REMOVED***
	flag.Usage = func() ***REMOVED***
		fmt.Fprintln(os.Stderr, `tomljson can be used in two ways:
Writing to STDIN and reading from STDOUT:
  cat file.toml | tomljson > file.json

Reading from a file name:
  tomljson file.toml
`)
	***REMOVED***
	flag.Parse()
	os.Exit(processMain(flag.Args(), os.Stdin, os.Stdout, os.Stderr))
***REMOVED***

func processMain(files []string, defaultInput io.Reader, output io.Writer, errorOutput io.Writer) int ***REMOVED***
	// read from stdin and print to stdout
	inputReader := defaultInput

	if len(files) > 0 ***REMOVED***
		var err error
		inputReader, err = os.Open(files[0])
		if err != nil ***REMOVED***
			printError(err, errorOutput)
			return -1
		***REMOVED***
	***REMOVED***
	s, err := reader(inputReader)
	if err != nil ***REMOVED***
		printError(err, errorOutput)
		return -1
	***REMOVED***
	io.WriteString(output, s+"\n")
	return 0
***REMOVED***

func printError(err error, output io.Writer) ***REMOVED***
	io.WriteString(output, err.Error()+"\n")
***REMOVED***

func reader(r io.Reader) (string, error) ***REMOVED***
	tree, err := toml.LoadReader(r)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return mapToJSON(tree)
***REMOVED***

func mapToJSON(tree *toml.Tree) (string, error) ***REMOVED***
	treeMap := tree.ToMap()
	bytes, err := json.MarshalIndent(treeMap, "", "  ")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return string(bytes[:]), nil
***REMOVED***
