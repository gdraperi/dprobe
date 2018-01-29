// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package registry

import (
	"errors"
	"io"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

const (
	// Registry value types.
	NONE                       = 0
	SZ                         = 1
	EXPAND_SZ                  = 2
	BINARY                     = 3
	DWORD                      = 4
	DWORD_BIG_ENDIAN           = 5
	LINK                       = 6
	MULTI_SZ                   = 7
	RESOURCE_LIST              = 8
	FULL_RESOURCE_DESCRIPTOR   = 9
	RESOURCE_REQUIREMENTS_LIST = 10
	QWORD                      = 11
)

var (
	// ErrShortBuffer is returned when the buffer was too short for the operation.
	ErrShortBuffer = syscall.ERROR_MORE_DATA

	// ErrNotExist is returned when a registry key or value does not exist.
	ErrNotExist = syscall.ERROR_FILE_NOT_FOUND

	// ErrUnexpectedType is returned by Get*Value when the value's type was unexpected.
	ErrUnexpectedType = errors.New("unexpected key value type")
)

// GetValue retrieves the type and data for the specified value associated
// with an open key k. It fills up buffer buf and returns the retrieved
// byte count n. If buf is too small to fit the stored value it returns
// ErrShortBuffer error along with the required buffer size n.
// If no buffer is provided, it returns true and actual buffer size n.
// If no buffer is provided, GetValue returns the value's type only.
// If the value does not exist, the error returned is ErrNotExist.
//
// GetValue is a low level function. If value's type is known, use the appropriate
// Get*Value function instead.
func (k Key) GetValue(name string, buf []byte) (n int, valtype uint32, err error) ***REMOVED***
	pname, err := syscall.UTF16PtrFromString(name)
	if err != nil ***REMOVED***
		return 0, 0, err
	***REMOVED***
	var pbuf *byte
	if len(buf) > 0 ***REMOVED***
		pbuf = (*byte)(unsafe.Pointer(&buf[0]))
	***REMOVED***
	l := uint32(len(buf))
	err = syscall.RegQueryValueEx(syscall.Handle(k), pname, nil, &valtype, pbuf, &l)
	if err != nil ***REMOVED***
		return int(l), valtype, err
	***REMOVED***
	return int(l), valtype, nil
***REMOVED***

func (k Key) getValue(name string, buf []byte) (date []byte, valtype uint32, err error) ***REMOVED***
	p, err := syscall.UTF16PtrFromString(name)
	if err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***
	var t uint32
	n := uint32(len(buf))
	for ***REMOVED***
		err = syscall.RegQueryValueEx(syscall.Handle(k), p, nil, &t, (*byte)(unsafe.Pointer(&buf[0])), &n)
		if err == nil ***REMOVED***
			return buf[:n], t, nil
		***REMOVED***
		if err != syscall.ERROR_MORE_DATA ***REMOVED***
			return nil, 0, err
		***REMOVED***
		if n <= uint32(len(buf)) ***REMOVED***
			return nil, 0, err
		***REMOVED***
		buf = make([]byte, n)
	***REMOVED***
***REMOVED***

// GetStringValue retrieves the string value for the specified
// value name associated with an open key k. It also returns the value's type.
// If value does not exist, GetStringValue returns ErrNotExist.
// If value is not SZ or EXPAND_SZ, it will return the correct value
// type and ErrUnexpectedType.
func (k Key) GetStringValue(name string) (val string, valtype uint32, err error) ***REMOVED***
	data, typ, err2 := k.getValue(name, make([]byte, 64))
	if err2 != nil ***REMOVED***
		return "", typ, err2
	***REMOVED***
	switch typ ***REMOVED***
	case SZ, EXPAND_SZ:
	default:
		return "", typ, ErrUnexpectedType
	***REMOVED***
	if len(data) == 0 ***REMOVED***
		return "", typ, nil
	***REMOVED***
	u := (*[1 << 29]uint16)(unsafe.Pointer(&data[0]))[:]
	return syscall.UTF16ToString(u), typ, nil
***REMOVED***

// GetMUIStringValue retrieves the localized string value for
// the specified value name associated with an open key k.
// If the value name doesn't exist or the localized string value
// can't be resolved, GetMUIStringValue returns ErrNotExist.
// GetMUIStringValue panics if the system doesn't support
// regLoadMUIString; use LoadRegLoadMUIString to check if
// regLoadMUIString is supported before calling this function.
func (k Key) GetMUIStringValue(name string) (string, error) ***REMOVED***
	pname, err := syscall.UTF16PtrFromString(name)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	buf := make([]uint16, 1024)
	var buflen uint32
	var pdir *uint16

	err = regLoadMUIString(syscall.Handle(k), pname, &buf[0], uint32(len(buf)), &buflen, 0, pdir)
	if err == syscall.ERROR_FILE_NOT_FOUND ***REMOVED*** // Try fallback path

		// Try to resolve the string value using the system directory as
		// a DLL search path; this assumes the string value is of the form
		// @[path]\dllname,-strID but with no path given, e.g. @tzres.dll,-320.

		// This approach works with tzres.dll but may have to be revised
		// in the future to allow callers to provide custom search paths.

		var s string
		s, err = ExpandString("%SystemRoot%\\system32\\")
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		pdir, err = syscall.UTF16PtrFromString(s)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***

		err = regLoadMUIString(syscall.Handle(k), pname, &buf[0], uint32(len(buf)), &buflen, 0, pdir)
	***REMOVED***

	for err == syscall.ERROR_MORE_DATA ***REMOVED*** // Grow buffer if needed
		if buflen <= uint32(len(buf)) ***REMOVED***
			break // Buffer not growing, assume race; break
		***REMOVED***
		buf = make([]uint16, buflen)
		err = regLoadMUIString(syscall.Handle(k), pname, &buf[0], uint32(len(buf)), &buflen, 0, pdir)
	***REMOVED***

	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return syscall.UTF16ToString(buf), nil
***REMOVED***

// ExpandString expands environment-variable strings and replaces
// them with the values defined for the current user.
// Use ExpandString to expand EXPAND_SZ strings.
func ExpandString(value string) (string, error) ***REMOVED***
	if value == "" ***REMOVED***
		return "", nil
	***REMOVED***
	p, err := syscall.UTF16PtrFromString(value)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	r := make([]uint16, 100)
	for ***REMOVED***
		n, err := expandEnvironmentStrings(p, &r[0], uint32(len(r)))
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		if n <= uint32(len(r)) ***REMOVED***
			u := (*[1 << 29]uint16)(unsafe.Pointer(&r[0]))[:]
			return syscall.UTF16ToString(u), nil
		***REMOVED***
		r = make([]uint16, n)
	***REMOVED***
***REMOVED***

// GetStringsValue retrieves the []string value for the specified
// value name associated with an open key k. It also returns the value's type.
// If value does not exist, GetStringsValue returns ErrNotExist.
// If value is not MULTI_SZ, it will return the correct value
// type and ErrUnexpectedType.
func (k Key) GetStringsValue(name string) (val []string, valtype uint32, err error) ***REMOVED***
	data, typ, err2 := k.getValue(name, make([]byte, 64))
	if err2 != nil ***REMOVED***
		return nil, typ, err2
	***REMOVED***
	if typ != MULTI_SZ ***REMOVED***
		return nil, typ, ErrUnexpectedType
	***REMOVED***
	if len(data) == 0 ***REMOVED***
		return nil, typ, nil
	***REMOVED***
	p := (*[1 << 29]uint16)(unsafe.Pointer(&data[0]))[:len(data)/2]
	if len(p) == 0 ***REMOVED***
		return nil, typ, nil
	***REMOVED***
	if p[len(p)-1] == 0 ***REMOVED***
		p = p[:len(p)-1] // remove terminating null
	***REMOVED***
	val = make([]string, 0, 5)
	from := 0
	for i, c := range p ***REMOVED***
		if c == 0 ***REMOVED***
			val = append(val, string(utf16.Decode(p[from:i])))
			from = i + 1
		***REMOVED***
	***REMOVED***
	return val, typ, nil
***REMOVED***

// GetIntegerValue retrieves the integer value for the specified
// value name associated with an open key k. It also returns the value's type.
// If value does not exist, GetIntegerValue returns ErrNotExist.
// If value is not DWORD or QWORD, it will return the correct value
// type and ErrUnexpectedType.
func (k Key) GetIntegerValue(name string) (val uint64, valtype uint32, err error) ***REMOVED***
	data, typ, err2 := k.getValue(name, make([]byte, 8))
	if err2 != nil ***REMOVED***
		return 0, typ, err2
	***REMOVED***
	switch typ ***REMOVED***
	case DWORD:
		if len(data) != 4 ***REMOVED***
			return 0, typ, errors.New("DWORD value is not 4 bytes long")
		***REMOVED***
		return uint64(*(*uint32)(unsafe.Pointer(&data[0]))), DWORD, nil
	case QWORD:
		if len(data) != 8 ***REMOVED***
			return 0, typ, errors.New("QWORD value is not 8 bytes long")
		***REMOVED***
		return uint64(*(*uint64)(unsafe.Pointer(&data[0]))), QWORD, nil
	default:
		return 0, typ, ErrUnexpectedType
	***REMOVED***
***REMOVED***

// GetBinaryValue retrieves the binary value for the specified
// value name associated with an open key k. It also returns the value's type.
// If value does not exist, GetBinaryValue returns ErrNotExist.
// If value is not BINARY, it will return the correct value
// type and ErrUnexpectedType.
func (k Key) GetBinaryValue(name string) (val []byte, valtype uint32, err error) ***REMOVED***
	data, typ, err2 := k.getValue(name, make([]byte, 64))
	if err2 != nil ***REMOVED***
		return nil, typ, err2
	***REMOVED***
	if typ != BINARY ***REMOVED***
		return nil, typ, ErrUnexpectedType
	***REMOVED***
	return data, typ, nil
***REMOVED***

func (k Key) setValue(name string, valtype uint32, data []byte) error ***REMOVED***
	p, err := syscall.UTF16PtrFromString(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if len(data) == 0 ***REMOVED***
		return regSetValueEx(syscall.Handle(k), p, 0, valtype, nil, 0)
	***REMOVED***
	return regSetValueEx(syscall.Handle(k), p, 0, valtype, &data[0], uint32(len(data)))
***REMOVED***

// SetDWordValue sets the data and type of a name value
// under key k to value and DWORD.
func (k Key) SetDWordValue(name string, value uint32) error ***REMOVED***
	return k.setValue(name, DWORD, (*[4]byte)(unsafe.Pointer(&value))[:])
***REMOVED***

// SetQWordValue sets the data and type of a name value
// under key k to value and QWORD.
func (k Key) SetQWordValue(name string, value uint64) error ***REMOVED***
	return k.setValue(name, QWORD, (*[8]byte)(unsafe.Pointer(&value))[:])
***REMOVED***

func (k Key) setStringValue(name string, valtype uint32, value string) error ***REMOVED***
	v, err := syscall.UTF16FromString(value)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	buf := (*[1 << 29]byte)(unsafe.Pointer(&v[0]))[:len(v)*2]
	return k.setValue(name, valtype, buf)
***REMOVED***

// SetStringValue sets the data and type of a name value
// under key k to value and SZ. The value must not contain a zero byte.
func (k Key) SetStringValue(name, value string) error ***REMOVED***
	return k.setStringValue(name, SZ, value)
***REMOVED***

// SetExpandStringValue sets the data and type of a name value
// under key k to value and EXPAND_SZ. The value must not contain a zero byte.
func (k Key) SetExpandStringValue(name, value string) error ***REMOVED***
	return k.setStringValue(name, EXPAND_SZ, value)
***REMOVED***

// SetStringsValue sets the data and type of a name value
// under key k to value and MULTI_SZ. The value strings
// must not contain a zero byte.
func (k Key) SetStringsValue(name string, value []string) error ***REMOVED***
	ss := ""
	for _, s := range value ***REMOVED***
		for i := 0; i < len(s); i++ ***REMOVED***
			if s[i] == 0 ***REMOVED***
				return errors.New("string cannot have 0 inside")
			***REMOVED***
		***REMOVED***
		ss += s + "\x00"
	***REMOVED***
	v := utf16.Encode([]rune(ss + "\x00"))
	buf := (*[1 << 29]byte)(unsafe.Pointer(&v[0]))[:len(v)*2]
	return k.setValue(name, MULTI_SZ, buf)
***REMOVED***

// SetBinaryValue sets the data and type of a name value
// under key k to value and BINARY.
func (k Key) SetBinaryValue(name string, value []byte) error ***REMOVED***
	return k.setValue(name, BINARY, value)
***REMOVED***

// DeleteValue removes a named value from the key k.
func (k Key) DeleteValue(name string) error ***REMOVED***
	return regDeleteValue(syscall.Handle(k), syscall.StringToUTF16Ptr(name))
***REMOVED***

// ReadValueNames returns the value names of key k.
// The parameter n controls the number of returned names,
// analogous to the way os.File.Readdirnames works.
func (k Key) ReadValueNames(n int) ([]string, error) ***REMOVED***
	ki, err := k.Stat()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	names := make([]string, 0, ki.ValueCount)
	buf := make([]uint16, ki.MaxValueNameLen+1) // extra room for terminating null character
loopItems:
	for i := uint32(0); ; i++ ***REMOVED***
		if n > 0 ***REMOVED***
			if len(names) == n ***REMOVED***
				return names, nil
			***REMOVED***
		***REMOVED***
		l := uint32(len(buf))
		for ***REMOVED***
			err := regEnumValue(syscall.Handle(k), i, &buf[0], &l, nil, nil, nil, nil)
			if err == nil ***REMOVED***
				break
			***REMOVED***
			if err == syscall.ERROR_MORE_DATA ***REMOVED***
				// Double buffer size and try again.
				l = uint32(2 * len(buf))
				buf = make([]uint16, l)
				continue
			***REMOVED***
			if err == _ERROR_NO_MORE_ITEMS ***REMOVED***
				break loopItems
			***REMOVED***
			return names, err
		***REMOVED***
		names = append(names, syscall.UTF16ToString(buf[:l]))
	***REMOVED***
	if n > len(names) ***REMOVED***
		return names, io.EOF
	***REMOVED***
	return names, nil
***REMOVED***
