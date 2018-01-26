package reexec

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var registeredInitializers = make(map[string]func())

// Register adds an initialization func under the specified name
func Register(name string, initializer func()) ***REMOVED***
	if _, exists := registeredInitializers[name]; exists ***REMOVED***
		panic(fmt.Sprintf("reexec func already registered under name %q", name))
	***REMOVED***

	registeredInitializers[name] = initializer
***REMOVED***

// Init is called as the first part of the exec process and returns true if an
// initialization function was called.
func Init() bool ***REMOVED***
	initializer, exists := registeredInitializers[os.Args[0]]
	if exists ***REMOVED***
		initializer()

		return true
	***REMOVED***
	return false
***REMOVED***

func naiveSelf() string ***REMOVED***
	name := os.Args[0]
	if filepath.Base(name) == name ***REMOVED***
		if lp, err := exec.LookPath(name); err == nil ***REMOVED***
			return lp
		***REMOVED***
	***REMOVED***
	// handle conversion of relative paths to absolute
	if absName, err := filepath.Abs(name); err == nil ***REMOVED***
		return absName
	***REMOVED***
	// if we couldn't get absolute name, return original
	// (NOTE: Go only errors on Abs() if os.Getwd fails)
	return name
***REMOVED***
