// Tomll is a linter for TOML
//
// Usage:
//   cat file.toml | tomll > file_linted.toml
//   tomll file1.toml file2.toml # lint the two files in place
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pelletier/go-toml"
)

func main() ***REMOVED***
	flag.Usage = func() ***REMOVED***
		fmt.Fprintln(os.Stderr, `tomll can be used in two ways:
Writing to STDIN and reading from STDOUT:
  cat file.toml | tomll > file.toml

Reading and updating a list of files:
  tomll a.toml b.toml c.toml

When given a list of files, tomll will modify all files in place without asking.
`)
	***REMOVED***
	flag.Parse()
	// read from stdin and print to stdout
	if flag.NArg() == 0 ***REMOVED***
		s, err := lintReader(os.Stdin)
		if err != nil ***REMOVED***
			io.WriteString(os.Stderr, err.Error())
			os.Exit(-1)
		***REMOVED***
		io.WriteString(os.Stdout, s)
	***REMOVED*** else ***REMOVED***
		// otherwise modify a list of files
		for _, filename := range flag.Args() ***REMOVED***
			s, err := lintFile(filename)
			if err != nil ***REMOVED***
				io.WriteString(os.Stderr, err.Error())
				os.Exit(-1)
			***REMOVED***
			ioutil.WriteFile(filename, []byte(s), 0644)
		***REMOVED***
	***REMOVED***
***REMOVED***

func lintFile(filename string) (string, error) ***REMOVED***
	tree, err := toml.LoadFile(filename)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return tree.String(), nil
***REMOVED***

func lintReader(r io.Reader) (string, error) ***REMOVED***
	tree, err := toml.LoadReader(r)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return tree.String(), nil
***REMOVED***
