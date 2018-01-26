package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// PathOp is a function which accepts a Path to perform some operation
type PathOp func(path Path) error

// WithContent writes content to a file at Path
func WithContent(content string) PathOp ***REMOVED***
	return func(path Path) error ***REMOVED***
		return ioutil.WriteFile(path.Path(), []byte(content), 0644)
	***REMOVED***
***REMOVED***

// WithBytes write bytes to a file at Path
func WithBytes(raw []byte) PathOp ***REMOVED***
	return func(path Path) error ***REMOVED***
		return ioutil.WriteFile(path.Path(), raw, 0644)
	***REMOVED***
***REMOVED***

// AsUser changes ownership of the file system object at Path
func AsUser(uid, gid int) PathOp ***REMOVED***
	return func(path Path) error ***REMOVED***
		return os.Chown(path.Path(), uid, gid)
	***REMOVED***
***REMOVED***

// WithFile creates a file in the directory at path with content
func WithFile(filename, content string, ops ...PathOp) PathOp ***REMOVED***
	return func(path Path) error ***REMOVED***
		fullpath := filepath.Join(path.Path(), filepath.FromSlash(filename))
		if err := createFile(fullpath, content); err != nil ***REMOVED***
			return err
		***REMOVED***
		return applyPathOps(&File***REMOVED***path: fullpath***REMOVED***, ops)
	***REMOVED***
***REMOVED***

func createFile(fullpath string, content string) error ***REMOVED***
	return ioutil.WriteFile(fullpath, []byte(content), 0644)
***REMOVED***

// WithFiles creates all the files in the directory at path with their content
func WithFiles(files map[string]string) PathOp ***REMOVED***
	return func(path Path) error ***REMOVED***
		for filename, content := range files ***REMOVED***
			fullpath := filepath.Join(path.Path(), filepath.FromSlash(filename))
			if err := createFile(fullpath, content); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// FromDir copies the directory tree from the source path into the new Dir
func FromDir(source string) PathOp ***REMOVED***
	return func(path Path) error ***REMOVED***
		return copyDirectory(source, path.Path())
	***REMOVED***
***REMOVED***

// WithDir creates a subdirectory in the directory at path. Additional PathOp
// can be used to modify the subdirectory
func WithDir(name string, ops ...PathOp) PathOp ***REMOVED***
	return func(path Path) error ***REMOVED***
		fullpath := filepath.Join(path.Path(), filepath.FromSlash(name))
		err := os.MkdirAll(fullpath, 0755)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return applyPathOps(&Dir***REMOVED***path: fullpath***REMOVED***, ops)
	***REMOVED***
***REMOVED***

func applyPathOps(path Path, ops []PathOp) error ***REMOVED***
	for _, op := range ops ***REMOVED***
		if err := op(path); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// WithMode sets the file mode on the directory or file at path
func WithMode(mode os.FileMode) PathOp ***REMOVED***
	return func(path Path) error ***REMOVED***
		return os.Chmod(path.Path(), mode)
	***REMOVED***
***REMOVED***

func copyDirectory(source, dest string) error ***REMOVED***
	entries, err := ioutil.ReadDir(source)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, entry := range entries ***REMOVED***
		sourcePath := filepath.Join(source, entry.Name())
		destPath := filepath.Join(dest, entry.Name())
		if entry.IsDir() ***REMOVED***
			if err := os.Mkdir(destPath, 0755); err != nil ***REMOVED***
				return err
			***REMOVED***
			if err := copyDirectory(sourcePath, destPath); err != nil ***REMOVED***
				return err
			***REMOVED***
			continue
		***REMOVED***
		if err := copyFile(sourcePath, destPath); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func copyFile(source, dest string) error ***REMOVED***
	content, err := ioutil.ReadFile(source)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return ioutil.WriteFile(dest, content, 0644)
***REMOVED***
