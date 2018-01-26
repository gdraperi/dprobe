package pathdriver

import (
	"path/filepath"
)

// PathDriver provides all of the path manipulation functions in a common
// interface. The context should call these and never use the `filepath`
// package or any other package to manipulate paths.
type PathDriver interface ***REMOVED***
	Join(paths ...string) string
	IsAbs(path string) bool
	Rel(base, target string) (string, error)
	Base(path string) string
	Dir(path string) string
	Clean(path string) string
	Split(path string) (dir, file string)
	Separator() byte
	Abs(path string) (string, error)
	Walk(string, filepath.WalkFunc) error
	FromSlash(path string) string
	ToSlash(path string) string
	Match(pattern, name string) (matched bool, err error)
***REMOVED***

// pathDriver is a simple default implementation calls the filepath package.
type pathDriver struct***REMOVED******REMOVED***

// LocalPathDriver is the exported pathDriver struct for convenience.
var LocalPathDriver PathDriver = &pathDriver***REMOVED******REMOVED***

func (*pathDriver) Join(paths ...string) string ***REMOVED***
	return filepath.Join(paths...)
***REMOVED***

func (*pathDriver) IsAbs(path string) bool ***REMOVED***
	return filepath.IsAbs(path)
***REMOVED***

func (*pathDriver) Rel(base, target string) (string, error) ***REMOVED***
	return filepath.Rel(base, target)
***REMOVED***

func (*pathDriver) Base(path string) string ***REMOVED***
	return filepath.Base(path)
***REMOVED***

func (*pathDriver) Dir(path string) string ***REMOVED***
	return filepath.Dir(path)
***REMOVED***

func (*pathDriver) Clean(path string) string ***REMOVED***
	return filepath.Clean(path)
***REMOVED***

func (*pathDriver) Split(path string) (dir, file string) ***REMOVED***
	return filepath.Split(path)
***REMOVED***

func (*pathDriver) Separator() byte ***REMOVED***
	return filepath.Separator
***REMOVED***

func (*pathDriver) Abs(path string) (string, error) ***REMOVED***
	return filepath.Abs(path)
***REMOVED***

// Note that filepath.Walk calls os.Stat, so if the context wants to
// to call Driver.Stat() for Walk, they need to create a new struct that
// overrides this method.
func (*pathDriver) Walk(root string, walkFn filepath.WalkFunc) error ***REMOVED***
	return filepath.Walk(root, walkFn)
***REMOVED***

func (*pathDriver) FromSlash(path string) string ***REMOVED***
	return filepath.FromSlash(path)
***REMOVED***

func (*pathDriver) ToSlash(path string) string ***REMOVED***
	return filepath.ToSlash(path)
***REMOVED***

func (*pathDriver) Match(pattern, name string) (bool, error) ***REMOVED***
	return filepath.Match(pattern, name)
***REMOVED***
