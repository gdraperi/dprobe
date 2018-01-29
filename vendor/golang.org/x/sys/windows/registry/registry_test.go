// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package registry_test

import (
	"bytes"
	"crypto/rand"
	"os"
	"syscall"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

func randKeyName(prefix string) string ***REMOVED***
	const numbers = "0123456789"
	buf := make([]byte, 10)
	rand.Read(buf)
	for i, b := range buf ***REMOVED***
		buf[i] = numbers[b%byte(len(numbers))]
	***REMOVED***
	return prefix + string(buf)
***REMOVED***

func TestReadSubKeyNames(t *testing.T) ***REMOVED***
	k, err := registry.OpenKey(registry.CLASSES_ROOT, "TypeLib", registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer k.Close()

	names, err := k.ReadSubKeyNames(-1)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	var foundStdOle bool
	for _, name := range names ***REMOVED***
		// Every PC has "stdole 2.0 OLE Automation" library installed.
		if name == "***REMOVED***00020430-0000-0000-C000-000000000046***REMOVED***" ***REMOVED***
			foundStdOle = true
		***REMOVED***
	***REMOVED***
	if !foundStdOle ***REMOVED***
		t.Fatal("could not find stdole 2.0 OLE Automation")
	***REMOVED***
***REMOVED***

func TestCreateOpenDeleteKey(t *testing.T) ***REMOVED***
	k, err := registry.OpenKey(registry.CURRENT_USER, "Software", registry.QUERY_VALUE)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer k.Close()

	testKName := randKeyName("TestCreateOpenDeleteKey_")

	testK, exist, err := registry.CreateKey(k, testKName, registry.CREATE_SUB_KEY)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer testK.Close()

	if exist ***REMOVED***
		t.Fatalf("key %q already exists", testKName)
	***REMOVED***

	testKAgain, exist, err := registry.CreateKey(k, testKName, registry.CREATE_SUB_KEY)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer testKAgain.Close()

	if !exist ***REMOVED***
		t.Fatalf("key %q should already exist", testKName)
	***REMOVED***

	testKOpened, err := registry.OpenKey(k, testKName, registry.ENUMERATE_SUB_KEYS)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer testKOpened.Close()

	err = registry.DeleteKey(k, testKName)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	testKOpenedAgain, err := registry.OpenKey(k, testKName, registry.ENUMERATE_SUB_KEYS)
	if err == nil ***REMOVED***
		defer testKOpenedAgain.Close()
		t.Fatalf("key %q should already been deleted", testKName)
	***REMOVED***
	if err != registry.ErrNotExist ***REMOVED***
		t.Fatalf(`unexpected error ("not exist" expected): %v`, err)
	***REMOVED***
***REMOVED***

func equalStringSlice(a, b []string) bool ***REMOVED***
	if len(a) != len(b) ***REMOVED***
		return false
	***REMOVED***
	if a == nil ***REMOVED***
		return true
	***REMOVED***
	for i := range a ***REMOVED***
		if a[i] != b[i] ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

type ValueTest struct ***REMOVED***
	Type     uint32
	Name     string
	Value    interface***REMOVED******REMOVED***
	WillFail bool
***REMOVED***

var ValueTests = []ValueTest***REMOVED***
	***REMOVED***Type: registry.SZ, Name: "String1", Value: ""***REMOVED***,
	***REMOVED***Type: registry.SZ, Name: "String2", Value: "\000", WillFail: true***REMOVED***,
	***REMOVED***Type: registry.SZ, Name: "String3", Value: "Hello World"***REMOVED***,
	***REMOVED***Type: registry.SZ, Name: "String4", Value: "Hello World\000", WillFail: true***REMOVED***,
	***REMOVED***Type: registry.EXPAND_SZ, Name: "ExpString1", Value: ""***REMOVED***,
	***REMOVED***Type: registry.EXPAND_SZ, Name: "ExpString2", Value: "\000", WillFail: true***REMOVED***,
	***REMOVED***Type: registry.EXPAND_SZ, Name: "ExpString3", Value: "Hello World"***REMOVED***,
	***REMOVED***Type: registry.EXPAND_SZ, Name: "ExpString4", Value: "Hello\000World", WillFail: true***REMOVED***,
	***REMOVED***Type: registry.EXPAND_SZ, Name: "ExpString5", Value: "%PATH%"***REMOVED***,
	***REMOVED***Type: registry.EXPAND_SZ, Name: "ExpString6", Value: "%NO_SUCH_VARIABLE%"***REMOVED***,
	***REMOVED***Type: registry.EXPAND_SZ, Name: "ExpString7", Value: "%PATH%;."***REMOVED***,
	***REMOVED***Type: registry.BINARY, Name: "Binary1", Value: []byte***REMOVED******REMOVED******REMOVED***,
	***REMOVED***Type: registry.BINARY, Name: "Binary2", Value: []byte***REMOVED***1, 2, 3***REMOVED******REMOVED***,
	***REMOVED***Type: registry.BINARY, Name: "Binary3", Value: []byte***REMOVED***3, 2, 1, 0, 1, 2, 3***REMOVED******REMOVED***,
	***REMOVED***Type: registry.DWORD, Name: "Dword1", Value: uint64(0)***REMOVED***,
	***REMOVED***Type: registry.DWORD, Name: "Dword2", Value: uint64(1)***REMOVED***,
	***REMOVED***Type: registry.DWORD, Name: "Dword3", Value: uint64(0xff)***REMOVED***,
	***REMOVED***Type: registry.DWORD, Name: "Dword4", Value: uint64(0xffff)***REMOVED***,
	***REMOVED***Type: registry.QWORD, Name: "Qword1", Value: uint64(0)***REMOVED***,
	***REMOVED***Type: registry.QWORD, Name: "Qword2", Value: uint64(1)***REMOVED***,
	***REMOVED***Type: registry.QWORD, Name: "Qword3", Value: uint64(0xff)***REMOVED***,
	***REMOVED***Type: registry.QWORD, Name: "Qword4", Value: uint64(0xffff)***REMOVED***,
	***REMOVED***Type: registry.QWORD, Name: "Qword5", Value: uint64(0xffffff)***REMOVED***,
	***REMOVED***Type: registry.QWORD, Name: "Qword6", Value: uint64(0xffffffff)***REMOVED***,
	***REMOVED***Type: registry.MULTI_SZ, Name: "MultiString1", Value: []string***REMOVED***"a", "b", "c"***REMOVED******REMOVED***,
	***REMOVED***Type: registry.MULTI_SZ, Name: "MultiString2", Value: []string***REMOVED***"abc", "", "cba"***REMOVED******REMOVED***,
	***REMOVED***Type: registry.MULTI_SZ, Name: "MultiString3", Value: []string***REMOVED***""***REMOVED******REMOVED***,
	***REMOVED***Type: registry.MULTI_SZ, Name: "MultiString4", Value: []string***REMOVED***"abcdef"***REMOVED******REMOVED***,
	***REMOVED***Type: registry.MULTI_SZ, Name: "MultiString5", Value: []string***REMOVED***"\000"***REMOVED***, WillFail: true***REMOVED***,
	***REMOVED***Type: registry.MULTI_SZ, Name: "MultiString6", Value: []string***REMOVED***"a\000b"***REMOVED***, WillFail: true***REMOVED***,
	***REMOVED***Type: registry.MULTI_SZ, Name: "MultiString7", Value: []string***REMOVED***"ab", "\000", "cd"***REMOVED***, WillFail: true***REMOVED***,
	***REMOVED***Type: registry.MULTI_SZ, Name: "MultiString8", Value: []string***REMOVED***"\000", "cd"***REMOVED***, WillFail: true***REMOVED***,
	***REMOVED***Type: registry.MULTI_SZ, Name: "MultiString9", Value: []string***REMOVED***"ab", "\000"***REMOVED***, WillFail: true***REMOVED***,
***REMOVED***

func setValues(t *testing.T, k registry.Key) ***REMOVED***
	for _, test := range ValueTests ***REMOVED***
		var err error
		switch test.Type ***REMOVED***
		case registry.SZ:
			err = k.SetStringValue(test.Name, test.Value.(string))
		case registry.EXPAND_SZ:
			err = k.SetExpandStringValue(test.Name, test.Value.(string))
		case registry.MULTI_SZ:
			err = k.SetStringsValue(test.Name, test.Value.([]string))
		case registry.BINARY:
			err = k.SetBinaryValue(test.Name, test.Value.([]byte))
		case registry.DWORD:
			err = k.SetDWordValue(test.Name, uint32(test.Value.(uint64)))
		case registry.QWORD:
			err = k.SetQWordValue(test.Name, test.Value.(uint64))
		default:
			t.Fatalf("unsupported type %d for %s value", test.Type, test.Name)
		***REMOVED***
		if test.WillFail ***REMOVED***
			if err == nil ***REMOVED***
				t.Fatalf("setting %s value %q should fail, but succeeded", test.Name, test.Value)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func enumerateValues(t *testing.T, k registry.Key) ***REMOVED***
	names, err := k.ReadValueNames(-1)
	if err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***
	haveNames := make(map[string]bool)
	for _, n := range names ***REMOVED***
		haveNames[n] = false
	***REMOVED***
	for _, test := range ValueTests ***REMOVED***
		wantFound := !test.WillFail
		_, haveFound := haveNames[test.Name]
		if wantFound && !haveFound ***REMOVED***
			t.Errorf("value %s is not found while enumerating", test.Name)
		***REMOVED***
		if haveFound && !wantFound ***REMOVED***
			t.Errorf("value %s is found while enumerating, but expected to fail", test.Name)
		***REMOVED***
		if haveFound ***REMOVED***
			delete(haveNames, test.Name)
		***REMOVED***
	***REMOVED***
	for n, v := range haveNames ***REMOVED***
		t.Errorf("value %s (%v) is found while enumerating, but has not been cretaed", n, v)
	***REMOVED***
***REMOVED***

func testErrNotExist(t *testing.T, name string, err error) ***REMOVED***
	if err == nil ***REMOVED***
		t.Errorf("%s value should not exist", name)
		return
	***REMOVED***
	if err != registry.ErrNotExist ***REMOVED***
		t.Errorf("reading %s value should return 'not exist' error, but got: %s", name, err)
		return
	***REMOVED***
***REMOVED***

func testErrUnexpectedType(t *testing.T, test ValueTest, gottype uint32, err error) ***REMOVED***
	if err == nil ***REMOVED***
		t.Errorf("GetXValue(%q) should not succeed", test.Name)
		return
	***REMOVED***
	if err != registry.ErrUnexpectedType ***REMOVED***
		t.Errorf("reading %s value should return 'unexpected key value type' error, but got: %s", test.Name, err)
		return
	***REMOVED***
	if gottype != test.Type ***REMOVED***
		t.Errorf("want %s value type %v, got %v", test.Name, test.Type, gottype)
		return
	***REMOVED***
***REMOVED***

func testGetStringValue(t *testing.T, k registry.Key, test ValueTest) ***REMOVED***
	got, gottype, err := k.GetStringValue(test.Name)
	if err != nil ***REMOVED***
		t.Errorf("GetStringValue(%s) failed: %v", test.Name, err)
		return
	***REMOVED***
	if got != test.Value ***REMOVED***
		t.Errorf("want %s value %q, got %q", test.Name, test.Value, got)
		return
	***REMOVED***
	if gottype != test.Type ***REMOVED***
		t.Errorf("want %s value type %v, got %v", test.Name, test.Type, gottype)
		return
	***REMOVED***
	if gottype == registry.EXPAND_SZ ***REMOVED***
		_, err = registry.ExpandString(got)
		if err != nil ***REMOVED***
			t.Errorf("ExpandString(%s) failed: %v", got, err)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func testGetIntegerValue(t *testing.T, k registry.Key, test ValueTest) ***REMOVED***
	got, gottype, err := k.GetIntegerValue(test.Name)
	if err != nil ***REMOVED***
		t.Errorf("GetIntegerValue(%s) failed: %v", test.Name, err)
		return
	***REMOVED***
	if got != test.Value.(uint64) ***REMOVED***
		t.Errorf("want %s value %v, got %v", test.Name, test.Value, got)
		return
	***REMOVED***
	if gottype != test.Type ***REMOVED***
		t.Errorf("want %s value type %v, got %v", test.Name, test.Type, gottype)
		return
	***REMOVED***
***REMOVED***

func testGetBinaryValue(t *testing.T, k registry.Key, test ValueTest) ***REMOVED***
	got, gottype, err := k.GetBinaryValue(test.Name)
	if err != nil ***REMOVED***
		t.Errorf("GetBinaryValue(%s) failed: %v", test.Name, err)
		return
	***REMOVED***
	if !bytes.Equal(got, test.Value.([]byte)) ***REMOVED***
		t.Errorf("want %s value %v, got %v", test.Name, test.Value, got)
		return
	***REMOVED***
	if gottype != test.Type ***REMOVED***
		t.Errorf("want %s value type %v, got %v", test.Name, test.Type, gottype)
		return
	***REMOVED***
***REMOVED***

func testGetStringsValue(t *testing.T, k registry.Key, test ValueTest) ***REMOVED***
	got, gottype, err := k.GetStringsValue(test.Name)
	if err != nil ***REMOVED***
		t.Errorf("GetStringsValue(%s) failed: %v", test.Name, err)
		return
	***REMOVED***
	if !equalStringSlice(got, test.Value.([]string)) ***REMOVED***
		t.Errorf("want %s value %#v, got %#v", test.Name, test.Value, got)
		return
	***REMOVED***
	if gottype != test.Type ***REMOVED***
		t.Errorf("want %s value type %v, got %v", test.Name, test.Type, gottype)
		return
	***REMOVED***
***REMOVED***

func testGetValue(t *testing.T, k registry.Key, test ValueTest, size int) ***REMOVED***
	if size <= 0 ***REMOVED***
		return
	***REMOVED***
	// read data with no buffer
	gotsize, gottype, err := k.GetValue(test.Name, nil)
	if err != nil ***REMOVED***
		t.Errorf("GetValue(%s, [%d]byte) failed: %v", test.Name, size, err)
		return
	***REMOVED***
	if gotsize != size ***REMOVED***
		t.Errorf("want %s value size of %d, got %v", test.Name, size, gotsize)
		return
	***REMOVED***
	if gottype != test.Type ***REMOVED***
		t.Errorf("want %s value type %v, got %v", test.Name, test.Type, gottype)
		return
	***REMOVED***
	// read data with short buffer
	gotsize, gottype, err = k.GetValue(test.Name, make([]byte, size-1))
	if err == nil ***REMOVED***
		t.Errorf("GetValue(%s, [%d]byte) should fail, but succeeded", test.Name, size-1)
		return
	***REMOVED***
	if err != registry.ErrShortBuffer ***REMOVED***
		t.Errorf("reading %s value should return 'short buffer' error, but got: %s", test.Name, err)
		return
	***REMOVED***
	if gotsize != size ***REMOVED***
		t.Errorf("want %s value size of %d, got %v", test.Name, size, gotsize)
		return
	***REMOVED***
	if gottype != test.Type ***REMOVED***
		t.Errorf("want %s value type %v, got %v", test.Name, test.Type, gottype)
		return
	***REMOVED***
	// read full data
	gotsize, gottype, err = k.GetValue(test.Name, make([]byte, size))
	if err != nil ***REMOVED***
		t.Errorf("GetValue(%s, [%d]byte) failed: %v", test.Name, size, err)
		return
	***REMOVED***
	if gotsize != size ***REMOVED***
		t.Errorf("want %s value size of %d, got %v", test.Name, size, gotsize)
		return
	***REMOVED***
	if gottype != test.Type ***REMOVED***
		t.Errorf("want %s value type %v, got %v", test.Name, test.Type, gottype)
		return
	***REMOVED***
	// check GetValue returns ErrNotExist as required
	_, _, err = k.GetValue(test.Name+"_not_there", make([]byte, size))
	if err == nil ***REMOVED***
		t.Errorf("GetValue(%q) should not succeed", test.Name)
		return
	***REMOVED***
	if err != registry.ErrNotExist ***REMOVED***
		t.Errorf("GetValue(%q) should return 'not exist' error, but got: %s", test.Name, err)
		return
	***REMOVED***
***REMOVED***

func testValues(t *testing.T, k registry.Key) ***REMOVED***
	for _, test := range ValueTests ***REMOVED***
		switch test.Type ***REMOVED***
		case registry.SZ, registry.EXPAND_SZ:
			if test.WillFail ***REMOVED***
				_, _, err := k.GetStringValue(test.Name)
				testErrNotExist(t, test.Name, err)
			***REMOVED*** else ***REMOVED***
				testGetStringValue(t, k, test)
				_, gottype, err := k.GetIntegerValue(test.Name)
				testErrUnexpectedType(t, test, gottype, err)
				// Size of utf16 string in bytes is not perfect,
				// but correct for current test values.
				// Size also includes terminating 0.
				testGetValue(t, k, test, (len(test.Value.(string))+1)*2)
			***REMOVED***
			_, _, err := k.GetStringValue(test.Name + "_string_not_created")
			testErrNotExist(t, test.Name+"_string_not_created", err)
		case registry.DWORD, registry.QWORD:
			testGetIntegerValue(t, k, test)
			_, gottype, err := k.GetBinaryValue(test.Name)
			testErrUnexpectedType(t, test, gottype, err)
			_, _, err = k.GetIntegerValue(test.Name + "_int_not_created")
			testErrNotExist(t, test.Name+"_int_not_created", err)
			size := 8
			if test.Type == registry.DWORD ***REMOVED***
				size = 4
			***REMOVED***
			testGetValue(t, k, test, size)
		case registry.BINARY:
			testGetBinaryValue(t, k, test)
			_, gottype, err := k.GetStringsValue(test.Name)
			testErrUnexpectedType(t, test, gottype, err)
			_, _, err = k.GetBinaryValue(test.Name + "_byte_not_created")
			testErrNotExist(t, test.Name+"_byte_not_created", err)
			testGetValue(t, k, test, len(test.Value.([]byte)))
		case registry.MULTI_SZ:
			if test.WillFail ***REMOVED***
				_, _, err := k.GetStringsValue(test.Name)
				testErrNotExist(t, test.Name, err)
			***REMOVED*** else ***REMOVED***
				testGetStringsValue(t, k, test)
				_, gottype, err := k.GetStringValue(test.Name)
				testErrUnexpectedType(t, test, gottype, err)
				size := 0
				for _, s := range test.Value.([]string) ***REMOVED***
					size += len(s) + 1 // nil terminated
				***REMOVED***
				size += 1 // extra nil at the end
				size *= 2 // count bytes, not uint16
				testGetValue(t, k, test, size)
			***REMOVED***
			_, _, err := k.GetStringsValue(test.Name + "_strings_not_created")
			testErrNotExist(t, test.Name+"_strings_not_created", err)
		default:
			t.Errorf("unsupported type %d for %s value", test.Type, test.Name)
			continue
		***REMOVED***
	***REMOVED***
***REMOVED***

func testStat(t *testing.T, k registry.Key) ***REMOVED***
	subk, _, err := registry.CreateKey(k, "subkey", registry.CREATE_SUB_KEY)
	if err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***
	defer subk.Close()

	defer registry.DeleteKey(k, "subkey")

	ki, err := k.Stat()
	if err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***
	if ki.SubKeyCount != 1 ***REMOVED***
		t.Error("key must have 1 subkey")
	***REMOVED***
	if ki.MaxSubKeyLen != 6 ***REMOVED***
		t.Error("key max subkey name length must be 6")
	***REMOVED***
	if ki.ValueCount != 24 ***REMOVED***
		t.Errorf("key must have 24 values, but is %d", ki.ValueCount)
	***REMOVED***
	if ki.MaxValueNameLen != 12 ***REMOVED***
		t.Errorf("key max value name length must be 10, but is %d", ki.MaxValueNameLen)
	***REMOVED***
	if ki.MaxValueLen != 38 ***REMOVED***
		t.Errorf("key max value length must be 38, but is %d", ki.MaxValueLen)
	***REMOVED***
	if mt, ct := ki.ModTime(), time.Now(); ct.Sub(mt) > 100*time.Millisecond ***REMOVED***
		t.Errorf("key mod time is not close to current time: mtime=%v current=%v delta=%v", mt, ct, ct.Sub(mt))
	***REMOVED***
***REMOVED***

func deleteValues(t *testing.T, k registry.Key) ***REMOVED***
	for _, test := range ValueTests ***REMOVED***
		if test.WillFail ***REMOVED***
			continue
		***REMOVED***
		err := k.DeleteValue(test.Name)
		if err != nil ***REMOVED***
			t.Error(err)
			continue
		***REMOVED***
	***REMOVED***
	names, err := k.ReadValueNames(-1)
	if err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***
	if len(names) != 0 ***REMOVED***
		t.Errorf("some values remain after deletion: %v", names)
	***REMOVED***
***REMOVED***

func TestValues(t *testing.T) ***REMOVED***
	softwareK, err := registry.OpenKey(registry.CURRENT_USER, "Software", registry.QUERY_VALUE)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer softwareK.Close()

	testKName := randKeyName("TestValues_")

	k, exist, err := registry.CreateKey(softwareK, testKName, registry.CREATE_SUB_KEY|registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer k.Close()

	if exist ***REMOVED***
		t.Fatalf("key %q already exists", testKName)
	***REMOVED***

	defer registry.DeleteKey(softwareK, testKName)

	setValues(t, k)

	enumerateValues(t, k)

	testValues(t, k)

	testStat(t, k)

	deleteValues(t, k)
***REMOVED***

func walkKey(t *testing.T, k registry.Key, kname string) ***REMOVED***
	names, err := k.ReadValueNames(-1)
	if err != nil ***REMOVED***
		t.Fatalf("reading value names of %s failed: %v", kname, err)
	***REMOVED***
	for _, name := range names ***REMOVED***
		_, valtype, err := k.GetValue(name, nil)
		if err != nil ***REMOVED***
			t.Fatalf("reading value type of %s of %s failed: %v", name, kname, err)
		***REMOVED***
		switch valtype ***REMOVED***
		case registry.NONE:
		case registry.SZ:
			_, _, err := k.GetStringValue(name)
			if err != nil ***REMOVED***
				t.Error(err)
			***REMOVED***
		case registry.EXPAND_SZ:
			s, _, err := k.GetStringValue(name)
			if err != nil ***REMOVED***
				t.Error(err)
			***REMOVED***
			_, err = registry.ExpandString(s)
			if err != nil ***REMOVED***
				t.Error(err)
			***REMOVED***
		case registry.DWORD, registry.QWORD:
			_, _, err := k.GetIntegerValue(name)
			if err != nil ***REMOVED***
				t.Error(err)
			***REMOVED***
		case registry.BINARY:
			_, _, err := k.GetBinaryValue(name)
			if err != nil ***REMOVED***
				t.Error(err)
			***REMOVED***
		case registry.MULTI_SZ:
			_, _, err := k.GetStringsValue(name)
			if err != nil ***REMOVED***
				t.Error(err)
			***REMOVED***
		case registry.FULL_RESOURCE_DESCRIPTOR, registry.RESOURCE_LIST, registry.RESOURCE_REQUIREMENTS_LIST:
			// TODO: not implemented
		default:
			t.Fatalf("value type %d of %s of %s failed: %v", valtype, name, kname, err)
		***REMOVED***
	***REMOVED***

	names, err = k.ReadSubKeyNames(-1)
	if err != nil ***REMOVED***
		t.Fatalf("reading sub-keys of %s failed: %v", kname, err)
	***REMOVED***
	for _, name := range names ***REMOVED***
		func() ***REMOVED***
			subk, err := registry.OpenKey(k, name, registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
			if err != nil ***REMOVED***
				if err == syscall.ERROR_ACCESS_DENIED ***REMOVED***
					// ignore error, if we are not allowed to access this key
					return
				***REMOVED***
				t.Fatalf("opening sub-keys %s of %s failed: %v", name, kname, err)
			***REMOVED***
			defer subk.Close()

			walkKey(t, subk, kname+`\`+name)
		***REMOVED***()
	***REMOVED***
***REMOVED***

func TestWalkFullRegistry(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip("skipping long running test in short mode")
	***REMOVED***
	walkKey(t, registry.CLASSES_ROOT, "CLASSES_ROOT")
	walkKey(t, registry.CURRENT_USER, "CURRENT_USER")
	walkKey(t, registry.LOCAL_MACHINE, "LOCAL_MACHINE")
	walkKey(t, registry.USERS, "USERS")
	walkKey(t, registry.CURRENT_CONFIG, "CURRENT_CONFIG")
***REMOVED***

func TestExpandString(t *testing.T) ***REMOVED***
	got, err := registry.ExpandString("%PATH%")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	want := os.Getenv("PATH")
	if got != want ***REMOVED***
		t.Errorf("want %q string expanded, got %q", want, got)
	***REMOVED***
***REMOVED***

func TestInvalidValues(t *testing.T) ***REMOVED***
	softwareK, err := registry.OpenKey(registry.CURRENT_USER, "Software", registry.QUERY_VALUE)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer softwareK.Close()

	testKName := randKeyName("TestInvalidValues_")

	k, exist, err := registry.CreateKey(softwareK, testKName, registry.CREATE_SUB_KEY|registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer k.Close()

	if exist ***REMOVED***
		t.Fatalf("key %q already exists", testKName)
	***REMOVED***

	defer registry.DeleteKey(softwareK, testKName)

	var tests = []struct ***REMOVED***
		Type uint32
		Name string
		Data []byte
	***REMOVED******REMOVED***
		***REMOVED***registry.DWORD, "Dword1", nil***REMOVED***,
		***REMOVED***registry.DWORD, "Dword2", []byte***REMOVED***1, 2, 3***REMOVED******REMOVED***,
		***REMOVED***registry.QWORD, "Qword1", nil***REMOVED***,
		***REMOVED***registry.QWORD, "Qword2", []byte***REMOVED***1, 2, 3***REMOVED******REMOVED***,
		***REMOVED***registry.QWORD, "Qword3", []byte***REMOVED***1, 2, 3, 4, 5, 6, 7***REMOVED******REMOVED***,
		***REMOVED***registry.MULTI_SZ, "MultiString1", nil***REMOVED***,
		***REMOVED***registry.MULTI_SZ, "MultiString2", []byte***REMOVED***0***REMOVED******REMOVED***,
		***REMOVED***registry.MULTI_SZ, "MultiString3", []byte***REMOVED***'a', 'b', 0***REMOVED******REMOVED***,
		***REMOVED***registry.MULTI_SZ, "MultiString4", []byte***REMOVED***'a', 0, 0, 'b', 0***REMOVED******REMOVED***,
		***REMOVED***registry.MULTI_SZ, "MultiString5", []byte***REMOVED***'a', 0, 0***REMOVED******REMOVED***,
	***REMOVED***

	for _, test := range tests ***REMOVED***
		err := k.SetValue(test.Name, test.Type, test.Data)
		if err != nil ***REMOVED***
			t.Fatalf("SetValue for %q failed: %v", test.Name, err)
		***REMOVED***
	***REMOVED***

	for _, test := range tests ***REMOVED***
		switch test.Type ***REMOVED***
		case registry.DWORD, registry.QWORD:
			value, valType, err := k.GetIntegerValue(test.Name)
			if err == nil ***REMOVED***
				t.Errorf("GetIntegerValue(%q) succeeded. Returns type=%d value=%v", test.Name, valType, value)
			***REMOVED***
		case registry.MULTI_SZ:
			value, valType, err := k.GetStringsValue(test.Name)
			if err == nil ***REMOVED***
				if len(value) != 0 ***REMOVED***
					t.Errorf("GetStringsValue(%q) succeeded. Returns type=%d value=%v", test.Name, valType, value)
				***REMOVED***
			***REMOVED***
		default:
			t.Errorf("unsupported type %d for %s value", test.Type, test.Name)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestGetMUIStringValue(t *testing.T) ***REMOVED***
	if err := registry.LoadRegLoadMUIString(); err != nil ***REMOVED***
		t.Skip("regLoadMUIString not supported; skipping")
	***REMOVED***
	if err := procGetDynamicTimeZoneInformation.Find(); err != nil ***REMOVED***
		t.Skipf("%s not supported; skipping", procGetDynamicTimeZoneInformation.Name)
	***REMOVED***
	var dtzi DynamicTimezoneinformation
	if _, err := GetDynamicTimeZoneInformation(&dtzi); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	tzKeyName := syscall.UTF16ToString(dtzi.TimeZoneKeyName[:])
	timezoneK, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion\Time Zones\`+tzKeyName, registry.READ)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer timezoneK.Close()

	type testType struct ***REMOVED***
		name string
		want string
	***REMOVED***
	var tests = []testType***REMOVED***
		***REMOVED***"MUI_Std", syscall.UTF16ToString(dtzi.StandardName[:])***REMOVED***,
	***REMOVED***
	if dtzi.DynamicDaylightTimeDisabled == 0 ***REMOVED***
		tests = append(tests, testType***REMOVED***"MUI_Dlt", syscall.UTF16ToString(dtzi.DaylightName[:])***REMOVED***)
	***REMOVED***

	for _, test := range tests ***REMOVED***
		got, err := timezoneK.GetMUIStringValue(test.name)
		if err != nil ***REMOVED***
			t.Error("GetMUIStringValue:", err)
		***REMOVED***

		if got != test.want ***REMOVED***
			t.Errorf("GetMUIStringValue: %s: Got %q, want %q", test.name, got, test.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

type DynamicTimezoneinformation struct ***REMOVED***
	Bias                        int32
	StandardName                [32]uint16
	StandardDate                syscall.Systemtime
	StandardBias                int32
	DaylightName                [32]uint16
	DaylightDate                syscall.Systemtime
	DaylightBias                int32
	TimeZoneKeyName             [128]uint16
	DynamicDaylightTimeDisabled uint8
***REMOVED***

var (
	kernel32DLL = syscall.NewLazyDLL("kernel32")

	procGetDynamicTimeZoneInformation = kernel32DLL.NewProc("GetDynamicTimeZoneInformation")
)

func GetDynamicTimeZoneInformation(dtzi *DynamicTimezoneinformation) (rc uint32, err error) ***REMOVED***
	r0, _, e1 := syscall.Syscall(procGetDynamicTimeZoneInformation.Addr(), 1, uintptr(unsafe.Pointer(dtzi)), 0, 0)
	rc = uint32(r0)
	if rc == 0xffffffff ***REMOVED***
		if e1 != 0 ***REMOVED***
			err = error(e1)
		***REMOVED*** else ***REMOVED***
			err = syscall.EINVAL
		***REMOVED***
	***REMOVED***
	return
***REMOVED***
