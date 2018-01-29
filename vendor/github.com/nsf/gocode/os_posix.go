// +build !windows

package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
)

func create_sock_flag(name, desc string) *string ***REMOVED***
	return flag.String(name, "unix", desc)
***REMOVED***

// Full path of the current executable
func get_executable_filename() string ***REMOVED***
	// try readlink first
	path, err := os.Readlink("/proc/self/exe")
	if err == nil ***REMOVED***
		return path
	***REMOVED***
	// use argv[0]
	path = os.Args[0]
	if !filepath.IsAbs(path) ***REMOVED***
		cwd, _ := os.Getwd()
		path = filepath.Join(cwd, path)
	***REMOVED***
	if file_exists(path) ***REMOVED***
		return path
	***REMOVED***
	// Fallback : use "gocode" and assume we are in the PATH...
	path, err = exec.LookPath("gocode")
	if err == nil ***REMOVED***
		return path
	***REMOVED***
	return ""
***REMOVED***

// config location

func config_dir() string ***REMOVED***
	return filepath.Join(xdg_home_dir(), "gocode")
***REMOVED***

func config_file() string ***REMOVED***
	return filepath.Join(xdg_home_dir(), "gocode", "config.json")
***REMOVED***
