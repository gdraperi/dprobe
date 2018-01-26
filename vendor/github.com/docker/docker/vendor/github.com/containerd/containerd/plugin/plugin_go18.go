// +build go1.8,!windows,amd64,!static_build

package plugin

import (
	"fmt"
	"path/filepath"
	"plugin"
	"runtime"
)

// loadPlugins loads all plugins for the OS and Arch
// that containerd is built for inside the provided path
func loadPlugins(path string) error ***REMOVED***
	abs, err := filepath.Abs(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	pattern := filepath.Join(abs, fmt.Sprintf(
		"*-%s-%s.%s",
		runtime.GOOS,
		runtime.GOARCH,
		getLibExt(),
	))
	libs, err := filepath.Glob(pattern)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, lib := range libs ***REMOVED***
		if _, err := plugin.Open(lib); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// getLibExt returns a platform specific lib extension for
// the platform that containerd is running on
func getLibExt() string ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "windows":
		return "dll"
	default:
		return "so"
	***REMOVED***
***REMOVED***
