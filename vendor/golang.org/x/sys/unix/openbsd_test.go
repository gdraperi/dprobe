// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build openbsd

// This, on the face of it, bizarre testing mechanism is necessary because
// the only reliable way to gauge whether or not a pledge(2) call has succeeded
// is that the program has been killed as a result of breaking its pledge.

package unix_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"golang.org/x/sys/unix"
)

type testProc struct ***REMOVED***
	fn      func()       // should always exit instead of returning
	cleanup func() error // for instance, delete coredumps from testing pledge
	success bool         // whether zero-exit means success or failure
***REMOVED***

var (
	testProcs = map[string]testProc***REMOVED******REMOVED***
	procName  = ""
)

const (
	optName = "sys-unix-internal-procname"
)

func init() ***REMOVED***
	flag.StringVar(&procName, optName, "", "internal use only")
***REMOVED***

// testCmd generates a proper command that, when executed, runs the test
// corresponding to the given key.
func testCmd(procName string) (*exec.Cmd, error) ***REMOVED***
	exe, err := filepath.Abs(os.Args[0])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	cmd := exec.Command(exe, "-"+optName+"="+procName)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd, nil
***REMOVED***

// ExitsCorrectly is a comprehensive, one-line-of-use wrapper for testing
// a testProc with a key.
func ExitsCorrectly(procName string, t *testing.T) ***REMOVED***
	s := testProcs[procName]
	c, err := testCmd(procName)
	defer func() ***REMOVED***
		if s.cleanup() != nil ***REMOVED***
			t.Fatalf("Failed to run cleanup for %s", procName)
		***REMOVED***
	***REMOVED***()
	if err != nil ***REMOVED***
		t.Fatalf("Failed to construct command for %s", procName)
	***REMOVED***
	if (c.Run() == nil) != s.success ***REMOVED***
		result := "succeed"
		if !s.success ***REMOVED***
			result = "fail"
		***REMOVED***
		t.Fatalf("Process did not %s when it was supposed to", result)
	***REMOVED***
***REMOVED***

func TestMain(m *testing.M) ***REMOVED***
	flag.Parse()
	if procName != "" ***REMOVED***
		testProcs[procName].fn()
	***REMOVED***
	os.Exit(m.Run())
***REMOVED***

// For example, add a test for pledge.
func init() ***REMOVED***
	testProcs["pledge"] = testProc***REMOVED***
		func() ***REMOVED***
			fmt.Println(unix.Pledge("", nil))
			os.Exit(0)
		***REMOVED***,
		func() error ***REMOVED***
			files, err := ioutil.ReadDir(".")
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			for _, file := range files ***REMOVED***
				if filepath.Ext(file.Name()) == ".core" ***REMOVED***
					if err := os.Remove(file.Name()); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
			***REMOVED***
			return nil
		***REMOVED***,
		false,
	***REMOVED***
***REMOVED***

func TestPledge(t *testing.T) ***REMOVED***
	ExitsCorrectly("pledge", t)
***REMOVED***
