// Derivative work from:
//	- https://golang.org/src/cmd/gofmt/gofmt.go
//	- https://github.com/fatih/hclfmt

package fmtcmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/hcl/printer"
)

var (
	ErrWriteStdin = errors.New("cannot use write option with standard input")
)

type Options struct ***REMOVED***
	List  bool // list files whose formatting differs
	Write bool // write result to (source) file instead of stdout
	Diff  bool // display diffs of formatting changes
***REMOVED***

func isValidFile(f os.FileInfo, extensions []string) bool ***REMOVED***
	if !f.IsDir() && !strings.HasPrefix(f.Name(), ".") ***REMOVED***
		for _, ext := range extensions ***REMOVED***
			if strings.HasSuffix(f.Name(), "."+ext) ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

// If in == nil, the source is the contents of the file with the given filename.
func processFile(filename string, in io.Reader, out io.Writer, stdin bool, opts Options) error ***REMOVED***
	if in == nil ***REMOVED***
		f, err := os.Open(filename)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		defer f.Close()
		in = f
	***REMOVED***

	src, err := ioutil.ReadAll(in)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	res, err := printer.Format(src)
	if err != nil ***REMOVED***
		return fmt.Errorf("In %s: %s", filename, err)
	***REMOVED***

	if !bytes.Equal(src, res) ***REMOVED***
		// formatting has changed
		if opts.List ***REMOVED***
			fmt.Fprintln(out, filename)
		***REMOVED***
		if opts.Write ***REMOVED***
			err = ioutil.WriteFile(filename, res, 0644)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if opts.Diff ***REMOVED***
			data, err := diff(src, res)
			if err != nil ***REMOVED***
				return fmt.Errorf("computing diff: %s", err)
			***REMOVED***
			fmt.Fprintf(out, "diff a/%s b/%s\n", filename, filename)
			out.Write(data)
		***REMOVED***
	***REMOVED***

	if !opts.List && !opts.Write && !opts.Diff ***REMOVED***
		_, err = out.Write(res)
	***REMOVED***

	return err
***REMOVED***

func walkDir(path string, extensions []string, stdout io.Writer, opts Options) error ***REMOVED***
	visitFile := func(path string, f os.FileInfo, err error) error ***REMOVED***
		if err == nil && isValidFile(f, extensions) ***REMOVED***
			err = processFile(path, nil, stdout, false, opts)
		***REMOVED***
		return err
	***REMOVED***

	return filepath.Walk(path, visitFile)
***REMOVED***

func Run(
	paths, extensions []string,
	stdin io.Reader,
	stdout io.Writer,
	opts Options,
) error ***REMOVED***
	if len(paths) == 0 ***REMOVED***
		if opts.Write ***REMOVED***
			return ErrWriteStdin
		***REMOVED***
		if err := processFile("<standard input>", stdin, stdout, true, opts); err != nil ***REMOVED***
			return err
		***REMOVED***
		return nil
	***REMOVED***

	for _, path := range paths ***REMOVED***
		switch dir, err := os.Stat(path); ***REMOVED***
		case err != nil:
			return err
		case dir.IsDir():
			if err := walkDir(path, extensions, stdout, opts); err != nil ***REMOVED***
				return err
			***REMOVED***
		default:
			if err := processFile(path, nil, stdout, false, opts); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func diff(b1, b2 []byte) (data []byte, err error) ***REMOVED***
	f1, err := ioutil.TempFile("", "")
	if err != nil ***REMOVED***
		return
	***REMOVED***
	defer os.Remove(f1.Name())
	defer f1.Close()

	f2, err := ioutil.TempFile("", "")
	if err != nil ***REMOVED***
		return
	***REMOVED***
	defer os.Remove(f2.Name())
	defer f2.Close()

	f1.Write(b1)
	f2.Write(b2)

	data, err = exec.Command("diff", "-u", f1.Name(), f2.Name()).CombinedOutput()
	if len(data) > 0 ***REMOVED***
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	***REMOVED***
	return
***REMOVED***
