package archive

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var testUntarFns = map[string]func(string, io.Reader) error***REMOVED***
	"untar": func(dest string, r io.Reader) error ***REMOVED***
		return Untar(r, dest, nil)
	***REMOVED***,
	"applylayer": func(dest string, r io.Reader) error ***REMOVED***
		_, err := ApplyLayer(dest, r)
		return err
	***REMOVED***,
***REMOVED***

// testBreakout is a helper function that, within the provided `tmpdir` directory,
// creates a `victim` folder with a generated `hello` file in it.
// `untar` extracts to a directory named `dest`, the tar file created from `headers`.
//
// Here are the tested scenarios:
// - removed `victim` folder				(write)
// - removed files from `victim` folder			(write)
// - new files in `victim` folder			(write)
// - modified files in `victim` folder			(write)
// - file in `dest` with same content as `victim/hello` (read)
//
// When using testBreakout make sure you cover one of the scenarios listed above.
func testBreakout(untarFn string, tmpdir string, headers []*tar.Header) error ***REMOVED***
	tmpdir, err := ioutil.TempDir("", tmpdir)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	dest := filepath.Join(tmpdir, "dest")
	if err := os.Mkdir(dest, 0755); err != nil ***REMOVED***
		return err
	***REMOVED***

	victim := filepath.Join(tmpdir, "victim")
	if err := os.Mkdir(victim, 0755); err != nil ***REMOVED***
		return err
	***REMOVED***
	hello := filepath.Join(victim, "hello")
	helloData, err := time.Now().MarshalText()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := ioutil.WriteFile(hello, helloData, 0644); err != nil ***REMOVED***
		return err
	***REMOVED***
	helloStat, err := os.Stat(hello)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	reader, writer := io.Pipe()
	go func() ***REMOVED***
		t := tar.NewWriter(writer)
		for _, hdr := range headers ***REMOVED***
			t.WriteHeader(hdr)
		***REMOVED***
		t.Close()
	***REMOVED***()

	untar := testUntarFns[untarFn]
	if untar == nil ***REMOVED***
		return fmt.Errorf("could not find untar function %q in testUntarFns", untarFn)
	***REMOVED***
	if err := untar(dest, reader); err != nil ***REMOVED***
		if _, ok := err.(breakoutError); !ok ***REMOVED***
			// If untar returns an error unrelated to an archive breakout,
			// then consider this an unexpected error and abort.
			return err
		***REMOVED***
		// Here, untar detected the breakout.
		// Let's move on verifying that indeed there was no breakout.
		fmt.Printf("breakoutError: %v\n", err)
	***REMOVED***

	// Check victim folder
	f, err := os.Open(victim)
	if err != nil ***REMOVED***
		// codepath taken if victim folder was removed
		return fmt.Errorf("archive breakout: error reading %q: %v", victim, err)
	***REMOVED***
	defer f.Close()

	// Check contents of victim folder
	//
	// We are only interested in getting 2 files from the victim folder, because if all is well
	// we expect only one result, the `hello` file. If there is a second result, it cannot
	// hold the same name `hello` and we assume that a new file got created in the victim folder.
	// That is enough to detect an archive breakout.
	names, err := f.Readdirnames(2)
	if err != nil ***REMOVED***
		// codepath taken if victim is not a folder
		return fmt.Errorf("archive breakout: error reading directory content of %q: %v", victim, err)
	***REMOVED***
	for _, name := range names ***REMOVED***
		if name != "hello" ***REMOVED***
			// codepath taken if new file was created in victim folder
			return fmt.Errorf("archive breakout: new file %q", name)
		***REMOVED***
	***REMOVED***

	// Check victim/hello
	f, err = os.Open(hello)
	if err != nil ***REMOVED***
		// codepath taken if read permissions were removed
		return fmt.Errorf("archive breakout: could not lstat %q: %v", hello, err)
	***REMOVED***
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	fi, err := f.Stat()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if helloStat.IsDir() != fi.IsDir() ||
		// TODO: cannot check for fi.ModTime() change
		helloStat.Mode() != fi.Mode() ||
		helloStat.Size() != fi.Size() ||
		!bytes.Equal(helloData, b) ***REMOVED***
		// codepath taken if hello has been modified
		return fmt.Errorf("archive breakout: file %q has been modified. Contents: expected=%q, got=%q. FileInfo: expected=%#v, got=%#v", hello, helloData, b, helloStat, fi)
	***REMOVED***

	// Check that nothing in dest/ has the same content as victim/hello.
	// Since victim/hello was generated with time.Now(), it is safe to assume
	// that any file whose content matches exactly victim/hello, managed somehow
	// to access victim/hello.
	return filepath.Walk(dest, func(path string, info os.FileInfo, err error) error ***REMOVED***
		if info.IsDir() ***REMOVED***
			if err != nil ***REMOVED***
				// skip directory if error
				return filepath.SkipDir
			***REMOVED***
			// enter directory
			return nil
		***REMOVED***
		if err != nil ***REMOVED***
			// skip file if error
			return nil
		***REMOVED***
		b, err := ioutil.ReadFile(path)
		if err != nil ***REMOVED***
			// Houston, we have a problem. Aborting (space)walk.
			return err
		***REMOVED***
		if bytes.Equal(helloData, b) ***REMOVED***
			return fmt.Errorf("archive breakout: file %q has been accessed via %q", hello, path)
		***REMOVED***
		return nil
	***REMOVED***)
***REMOVED***
