// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build freebsd

package unix_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/sys/unix"
)

func TestSysctlUint64(t *testing.T) ***REMOVED***
	_, err := unix.SysctlUint64("vm.swap_total")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// FIXME: Infrastructure for launching tests in subprocesses stolen from openbsd_test.go - refactor?
// testCmd generates a proper command that, when executed, runs the test
// corresponding to the given key.

type testProc struct ***REMOVED***
	fn      func()                    // should always exit instead of returning
	arg     func(t *testing.T) string // generate argument for test
	cleanup func(arg string) error    // for instance, delete coredumps from testing pledge
	success bool                      // whether zero-exit means success or failure
***REMOVED***

var (
	testProcs = map[string]testProc***REMOVED******REMOVED***
	procName  = ""
	procArg   = ""
)

const (
	optName = "sys-unix-internal-procname"
	optArg  = "sys-unix-internal-arg"
)

func init() ***REMOVED***
	flag.StringVar(&procName, optName, "", "internal use only")
	flag.StringVar(&procArg, optArg, "", "internal use only")

***REMOVED***

func testCmd(procName string, procArg string) (*exec.Cmd, error) ***REMOVED***
	exe, err := filepath.Abs(os.Args[0])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	cmd := exec.Command(exe, "-"+optName+"="+procName, "-"+optArg+"="+procArg)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd, nil
***REMOVED***

// ExitsCorrectly is a comprehensive, one-line-of-use wrapper for testing
// a testProc with a key.
func ExitsCorrectly(t *testing.T, procName string) ***REMOVED***
	s := testProcs[procName]
	arg := "-"
	if s.arg != nil ***REMOVED***
		arg = s.arg(t)
	***REMOVED***
	c, err := testCmd(procName, arg)
	defer func(arg string) ***REMOVED***
		if err := s.cleanup(arg); err != nil ***REMOVED***
			t.Fatalf("Failed to run cleanup for %s %s %#v", procName, err, err)
		***REMOVED***
	***REMOVED***(arg)
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
		t := testProcs[procName]
		t.fn()
		os.Stderr.WriteString("test function did not exit\n")
		if t.success ***REMOVED***
			os.Exit(1)
		***REMOVED*** else ***REMOVED***
			os.Exit(0)
		***REMOVED***
	***REMOVED***
	os.Exit(m.Run())
***REMOVED***

// end of infrastructure

const testfile = "gocapmodetest"
const testfile2 = testfile + "2"

func CapEnterTest() ***REMOVED***
	_, err := os.OpenFile(path.Join(procArg, testfile), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("OpenFile: %s", err))
	***REMOVED***

	err = unix.CapEnter()
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("CapEnter: %s", err))
	***REMOVED***

	_, err = os.OpenFile(path.Join(procArg, testfile2), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil ***REMOVED***
		panic("OpenFile works!")
	***REMOVED***
	if err.(*os.PathError).Err != unix.ECAPMODE ***REMOVED***
		panic(fmt.Sprintf("OpenFile failed wrong: %s %#v", err, err))
	***REMOVED***
	os.Exit(0)
***REMOVED***

func makeTempDir(t *testing.T) string ***REMOVED***
	d, err := ioutil.TempDir("", "go_openat_test")
	if err != nil ***REMOVED***
		t.Fatalf("TempDir failed: %s", err)
	***REMOVED***
	return d
***REMOVED***

func removeTempDir(arg string) error ***REMOVED***
	err := os.RemoveAll(arg)
	if err != nil && err.(*os.PathError).Err == unix.ENOENT ***REMOVED***
		return nil
	***REMOVED***
	return err
***REMOVED***

func init() ***REMOVED***
	testProcs["cap_enter"] = testProc***REMOVED***
		CapEnterTest,
		makeTempDir,
		removeTempDir,
		true,
	***REMOVED***
***REMOVED***

func TestCapEnter(t *testing.T) ***REMOVED***
	if runtime.GOARCH != "amd64" ***REMOVED***
		t.Skipf("skipping test on %s", runtime.GOARCH)
	***REMOVED***
	ExitsCorrectly(t, "cap_enter")
***REMOVED***

func OpenatTest() ***REMOVED***
	f, err := os.Open(procArg)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	err = unix.CapEnter()
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("CapEnter: %s", err))
	***REMOVED***

	fxx, err := unix.Openat(int(f.Fd()), "xx", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	unix.Close(fxx)

	// The right to open BASE/xx is not ambient
	_, err = os.OpenFile(procArg+"/xx", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil ***REMOVED***
		panic("OpenFile succeeded")
	***REMOVED***
	if err.(*os.PathError).Err != unix.ECAPMODE ***REMOVED***
		panic(fmt.Sprintf("OpenFile failed wrong: %s %#v", err, err))
	***REMOVED***

	// Can't make a new directory either
	err = os.Mkdir(procArg+"2", 0777)
	if err == nil ***REMOVED***
		panic("MKdir succeeded")
	***REMOVED***
	if err.(*os.PathError).Err != unix.ECAPMODE ***REMOVED***
		panic(fmt.Sprintf("Mkdir failed wrong: %s %#v", err, err))
	***REMOVED***

	// Remove all caps except read and lookup.
	r, err := unix.CapRightsInit([]uint64***REMOVED***unix.CAP_READ, unix.CAP_LOOKUP***REMOVED***)
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("CapRightsInit failed: %s %#v", err, err))
	***REMOVED***
	err = unix.CapRightsLimit(f.Fd(), r)
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("CapRightsLimit failed: %s %#v", err, err))
	***REMOVED***

	// Check we can get the rights back again
	r, err = unix.CapRightsGet(f.Fd())
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("CapRightsGet failed: %s %#v", err, err))
	***REMOVED***
	b, err := unix.CapRightsIsSet(r, []uint64***REMOVED***unix.CAP_READ, unix.CAP_LOOKUP***REMOVED***)
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("CapRightsIsSet failed: %s %#v", err, err))
	***REMOVED***
	if !b ***REMOVED***
		panic(fmt.Sprintf("Unexpected rights"))
	***REMOVED***
	b, err = unix.CapRightsIsSet(r, []uint64***REMOVED***unix.CAP_READ, unix.CAP_LOOKUP, unix.CAP_WRITE***REMOVED***)
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("CapRightsIsSet failed: %s %#v", err, err))
	***REMOVED***
	if b ***REMOVED***
		panic(fmt.Sprintf("Unexpected rights (2)"))
	***REMOVED***

	// Can no longer create a file
	_, err = unix.Openat(int(f.Fd()), "xx2", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil ***REMOVED***
		panic("Openat succeeded")
	***REMOVED***
	if err != unix.ENOTCAPABLE ***REMOVED***
		panic(fmt.Sprintf("OpenFileAt failed wrong: %s %#v", err, err))
	***REMOVED***

	// But can read an existing one
	_, err = unix.Openat(int(f.Fd()), "xx", os.O_RDONLY, 0666)
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("Openat failed: %s %#v", err, err))
	***REMOVED***

	os.Exit(0)
***REMOVED***

func init() ***REMOVED***
	testProcs["openat"] = testProc***REMOVED***
		OpenatTest,
		makeTempDir,
		removeTempDir,
		true,
	***REMOVED***
***REMOVED***

func TestOpenat(t *testing.T) ***REMOVED***
	if runtime.GOARCH != "amd64" ***REMOVED***
		t.Skipf("skipping test on %s", runtime.GOARCH)
	***REMOVED***
	ExitsCorrectly(t, "openat")
***REMOVED***

func TestCapRightsSetAndClear(t *testing.T) ***REMOVED***
	r, err := unix.CapRightsInit([]uint64***REMOVED***unix.CAP_READ, unix.CAP_WRITE, unix.CAP_PDWAIT***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("CapRightsInit failed: %s", err)
	***REMOVED***

	err = unix.CapRightsSet(r, []uint64***REMOVED***unix.CAP_EVENT, unix.CAP_LISTEN***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("CapRightsSet failed: %s", err)
	***REMOVED***

	b, err := unix.CapRightsIsSet(r, []uint64***REMOVED***unix.CAP_READ, unix.CAP_WRITE, unix.CAP_PDWAIT, unix.CAP_EVENT, unix.CAP_LISTEN***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("CapRightsIsSet failed: %s", err)
	***REMOVED***
	if !b ***REMOVED***
		t.Fatalf("Wrong rights set")
	***REMOVED***

	err = unix.CapRightsClear(r, []uint64***REMOVED***unix.CAP_READ, unix.CAP_PDWAIT***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("CapRightsClear failed: %s", err)
	***REMOVED***

	b, err = unix.CapRightsIsSet(r, []uint64***REMOVED***unix.CAP_WRITE, unix.CAP_EVENT, unix.CAP_LISTEN***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("CapRightsIsSet failed: %s", err)
	***REMOVED***
	if !b ***REMOVED***
		t.Fatalf("Wrong rights set")
	***REMOVED***
***REMOVED***
