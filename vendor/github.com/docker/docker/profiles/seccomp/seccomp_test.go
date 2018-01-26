// +build linux

package seccomp

import (
	"io/ioutil"
	"testing"

	"github.com/docker/docker/oci"
)

func TestLoadProfile(t *testing.T) ***REMOVED***
	f, err := ioutil.ReadFile("fixtures/example.json")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	rs := oci.DefaultSpec()
	if _, err := LoadProfile(string(f), &rs); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestLoadDefaultProfile(t *testing.T) ***REMOVED***
	f, err := ioutil.ReadFile("default.json")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	rs := oci.DefaultSpec()
	if _, err := LoadProfile(string(f), &rs); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
