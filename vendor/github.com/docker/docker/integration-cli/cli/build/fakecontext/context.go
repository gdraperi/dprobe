package fakecontext

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/archive"
)

type testingT interface ***REMOVED***
	Fatal(args ...interface***REMOVED******REMOVED***)
	Fatalf(string, ...interface***REMOVED******REMOVED***)
***REMOVED***

// New creates a fake build context
func New(t testingT, dir string, modifiers ...func(*Fake) error) *Fake ***REMOVED***
	fakeContext := &Fake***REMOVED***Dir: dir***REMOVED***
	if dir == "" ***REMOVED***
		if err := newDir(fakeContext); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***

	for _, modifier := range modifiers ***REMOVED***
		if err := modifier(fakeContext); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***

	return fakeContext
***REMOVED***

func newDir(fake *Fake) error ***REMOVED***
	tmp, err := ioutil.TempDir("", "fake-context")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := os.Chmod(tmp, 0755); err != nil ***REMOVED***
		return err
	***REMOVED***
	fake.Dir = tmp
	return nil
***REMOVED***

// WithFile adds the specified file (with content) in the build context
func WithFile(name, content string) func(*Fake) error ***REMOVED***
	return func(ctx *Fake) error ***REMOVED***
		return ctx.Add(name, content)
	***REMOVED***
***REMOVED***

// WithDockerfile adds the specified content as Dockerfile in the build context
func WithDockerfile(content string) func(*Fake) error ***REMOVED***
	return WithFile("Dockerfile", content)
***REMOVED***

// WithFiles adds the specified files in the build context, content is a string
func WithFiles(files map[string]string) func(*Fake) error ***REMOVED***
	return func(fakeContext *Fake) error ***REMOVED***
		for file, content := range files ***REMOVED***
			if err := fakeContext.Add(file, content); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// WithBinaryFiles adds the specified files in the build context, content is binary
func WithBinaryFiles(files map[string]*bytes.Buffer) func(*Fake) error ***REMOVED***
	return func(fakeContext *Fake) error ***REMOVED***
		for file, content := range files ***REMOVED***
			if err := fakeContext.Add(file, string(content.Bytes())); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// Fake creates directories that can be used as a build context
type Fake struct ***REMOVED***
	Dir string
***REMOVED***

// Add a file at a path, creating directories where necessary
func (f *Fake) Add(file, content string) error ***REMOVED***
	return f.addFile(file, []byte(content))
***REMOVED***

func (f *Fake) addFile(file string, content []byte) error ***REMOVED***
	fp := filepath.Join(f.Dir, filepath.FromSlash(file))
	dirpath := filepath.Dir(fp)
	if dirpath != "." ***REMOVED***
		if err := os.MkdirAll(dirpath, 0755); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return ioutil.WriteFile(fp, content, 0644)

***REMOVED***

// Delete a file at a path
func (f *Fake) Delete(file string) error ***REMOVED***
	fp := filepath.Join(f.Dir, filepath.FromSlash(file))
	return os.RemoveAll(fp)
***REMOVED***

// Close deletes the context
func (f *Fake) Close() error ***REMOVED***
	return os.RemoveAll(f.Dir)
***REMOVED***

// AsTarReader returns a ReadCloser with the contents of Dir as a tar archive.
func (f *Fake) AsTarReader(t testingT) io.ReadCloser ***REMOVED***
	reader, err := archive.TarWithOptions(f.Dir, &archive.TarOptions***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("Failed to create tar from %s: %s", f.Dir, err)
	***REMOVED***
	return reader
***REMOVED***
