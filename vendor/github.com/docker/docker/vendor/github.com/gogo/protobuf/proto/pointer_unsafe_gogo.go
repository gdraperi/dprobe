// Protocol Buffers for Go with Gadgets
//
// Copyright (c) 2013, The GoGo Authors. All rights reserved.
// http://github.com/gogo/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// +build !appengine,!js

// This file contains the implementation of the proto field accesses using package unsafe.

package proto

import (
	"reflect"
	"unsafe"
)

func structPointer_InterfaceAt(p structPointer, f field, t reflect.Type) interface***REMOVED******REMOVED*** ***REMOVED***
	point := unsafe.Pointer(uintptr(p) + uintptr(f))
	r := reflect.NewAt(t, point)
	return r.Interface()
***REMOVED***

func structPointer_InterfaceRef(p structPointer, f field, t reflect.Type) interface***REMOVED******REMOVED*** ***REMOVED***
	point := unsafe.Pointer(uintptr(p) + uintptr(f))
	r := reflect.NewAt(t, point)
	if r.Elem().IsNil() ***REMOVED***
		return nil
	***REMOVED***
	return r.Elem().Interface()
***REMOVED***

func copyUintPtr(oldptr, newptr uintptr, size int) ***REMOVED***
	oldbytes := make([]byte, 0)
	oldslice := (*reflect.SliceHeader)(unsafe.Pointer(&oldbytes))
	oldslice.Data = oldptr
	oldslice.Len = size
	oldslice.Cap = size
	newbytes := make([]byte, 0)
	newslice := (*reflect.SliceHeader)(unsafe.Pointer(&newbytes))
	newslice.Data = newptr
	newslice.Len = size
	newslice.Cap = size
	copy(newbytes, oldbytes)
***REMOVED***

func structPointer_Copy(oldptr structPointer, newptr structPointer, size int) ***REMOVED***
	copyUintPtr(uintptr(oldptr), uintptr(newptr), size)
***REMOVED***

func appendStructPointer(base structPointer, f field, typ reflect.Type) structPointer ***REMOVED***
	size := typ.Elem().Size()

	oldHeader := structPointer_GetSliceHeader(base, f)
	oldSlice := reflect.NewAt(typ, unsafe.Pointer(oldHeader)).Elem()
	newLen := oldHeader.Len + 1
	newSlice := reflect.MakeSlice(typ, newLen, newLen)
	reflect.Copy(newSlice, oldSlice)
	bas := toStructPointer(newSlice)
	oldHeader.Data = uintptr(bas)
	oldHeader.Len = newLen
	oldHeader.Cap = newLen

	return structPointer(unsafe.Pointer(uintptr(unsafe.Pointer(bas)) + uintptr(uintptr(newLen-1)*size)))
***REMOVED***

func structPointer_FieldPointer(p structPointer, f field) structPointer ***REMOVED***
	return structPointer(unsafe.Pointer(uintptr(p) + uintptr(f)))
***REMOVED***

func structPointer_GetRefStructPointer(p structPointer, f field) structPointer ***REMOVED***
	return structPointer((*structPointer)(unsafe.Pointer(uintptr(p) + uintptr(f))))
***REMOVED***

func structPointer_GetSliceHeader(p structPointer, f field) *reflect.SliceHeader ***REMOVED***
	return (*reflect.SliceHeader)(unsafe.Pointer(uintptr(p) + uintptr(f)))
***REMOVED***

func structPointer_Add(p structPointer, size field) structPointer ***REMOVED***
	return structPointer(unsafe.Pointer(uintptr(p) + uintptr(size)))
***REMOVED***

func structPointer_Len(p structPointer, f field) int ***REMOVED***
	return len(*(*[]interface***REMOVED******REMOVED***)(unsafe.Pointer(structPointer_GetRefStructPointer(p, f))))
***REMOVED***

func structPointer_StructRefSlice(p structPointer, f field, size uintptr) *structRefSlice ***REMOVED***
	return &structRefSlice***REMOVED***p: p, f: f, size: size***REMOVED***
***REMOVED***

// A structRefSlice represents a slice of structs (themselves submessages or groups).
type structRefSlice struct ***REMOVED***
	p    structPointer
	f    field
	size uintptr
***REMOVED***

func (v *structRefSlice) Len() int ***REMOVED***
	return structPointer_Len(v.p, v.f)
***REMOVED***

func (v *structRefSlice) Index(i int) structPointer ***REMOVED***
	ss := structPointer_GetStructPointer(v.p, v.f)
	ss1 := structPointer_GetRefStructPointer(ss, 0)
	return structPointer_Add(ss1, field(uintptr(i)*v.size))
***REMOVED***
