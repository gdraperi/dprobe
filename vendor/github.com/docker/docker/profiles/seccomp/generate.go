// +build ignore

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker/docker/profiles/seccomp"
)

// saves the default seccomp profile as a json file so people can use it as a
// base for their own custom profiles
func main() ***REMOVED***
	wd, err := os.Getwd()
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	f := filepath.Join(wd, "default.json")

	// write the default profile to the file
	b, err := json.MarshalIndent(seccomp.DefaultProfile(), "", "\t")
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	if err := ioutil.WriteFile(f, b, 0644); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***
