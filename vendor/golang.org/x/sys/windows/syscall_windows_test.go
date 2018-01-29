// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package windows_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
)

func TestWin32finddata(t *testing.T) ***REMOVED***
	dir, err := ioutil.TempDir("", "go-build")
	if err != nil ***REMOVED***
		t.Fatalf("failed to create temp directory: %v", err)
	***REMOVED***
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "long_name.and_extension")
	f, err := os.Create(path)
	if err != nil ***REMOVED***
		t.Fatalf("failed to create %v: %v", path, err)
	***REMOVED***
	f.Close()

	type X struct ***REMOVED***
		fd  windows.Win32finddata
		got byte
		pad [10]byte // to protect ourselves

	***REMOVED***
	var want byte = 2 // it is unlikely to have this character in the filename
	x := X***REMOVED***got: want***REMOVED***

	pathp, _ := windows.UTF16PtrFromString(path)
	h, err := windows.FindFirstFile(pathp, &(x.fd))
	if err != nil ***REMOVED***
		t.Fatalf("FindFirstFile failed: %v", err)
	***REMOVED***
	err = windows.FindClose(h)
	if err != nil ***REMOVED***
		t.Fatalf("FindClose failed: %v", err)
	***REMOVED***

	if x.got != want ***REMOVED***
		t.Fatalf("memory corruption: want=%d got=%d", want, x.got)
	***REMOVED***
***REMOVED***

func TestFormatMessage(t *testing.T) ***REMOVED***
	dll := windows.MustLoadDLL("pdh.dll")

	pdhOpenQuery := func(datasrc *uint16, userdata uint32, query *windows.Handle) (errno uintptr) ***REMOVED***
		r0, _, _ := syscall.Syscall(dll.MustFindProc("PdhOpenQueryW").Addr(), 3, uintptr(unsafe.Pointer(datasrc)), uintptr(userdata), uintptr(unsafe.Pointer(query)))
		return r0
	***REMOVED***

	pdhCloseQuery := func(query windows.Handle) (errno uintptr) ***REMOVED***
		r0, _, _ := syscall.Syscall(dll.MustFindProc("PdhCloseQuery").Addr(), 1, uintptr(query), 0, 0)
		return r0
	***REMOVED***

	var q windows.Handle
	name, err := windows.UTF16PtrFromString("no_such_source")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	errno := pdhOpenQuery(name, 0, &q)
	if errno == 0 ***REMOVED***
		pdhCloseQuery(q)
		t.Fatal("PdhOpenQuery succeeded, but expected to fail.")
	***REMOVED***

	const flags uint32 = syscall.FORMAT_MESSAGE_FROM_HMODULE | syscall.FORMAT_MESSAGE_ARGUMENT_ARRAY | syscall.FORMAT_MESSAGE_IGNORE_INSERTS
	buf := make([]uint16, 300)
	_, err = windows.FormatMessage(flags, uintptr(dll.Handle), uint32(errno), 0, buf, nil)
	if err != nil ***REMOVED***
		t.Fatalf("FormatMessage for handle=%x and errno=%x failed: %v", dll.Handle, errno, err)
	***REMOVED***
***REMOVED***

func abort(funcname string, err error) ***REMOVED***
	panic(funcname + " failed: " + err.Error())
***REMOVED***

func ExampleLoadLibrary() ***REMOVED***
	h, err := windows.LoadLibrary("kernel32.dll")
	if err != nil ***REMOVED***
		abort("LoadLibrary", err)
	***REMOVED***
	defer windows.FreeLibrary(h)
	proc, err := windows.GetProcAddress(h, "GetVersion")
	if err != nil ***REMOVED***
		abort("GetProcAddress", err)
	***REMOVED***
	r, _, _ := syscall.Syscall(uintptr(proc), 0, 0, 0, 0)
	major := byte(r)
	minor := uint8(r >> 8)
	build := uint16(r >> 16)
	print("windows version ", major, ".", minor, " (Build ", build, ")\n")
***REMOVED***
