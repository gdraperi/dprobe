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

package proto

import (
	"reflect"
)

// Decode a reference to a struct pointer.
func (o *Buffer) dec_ref_struct_message(p *Properties, base structPointer) (err error) ***REMOVED***
	raw, e := o.DecodeRawBytes(false)
	if e != nil ***REMOVED***
		return e
	***REMOVED***

	// If the object can unmarshal itself, let it.
	if p.isUnmarshaler ***REMOVED***
		panic("not supported, since this is a pointer receiver")
	***REMOVED***

	obuf := o.buf
	oi := o.index
	o.buf = raw
	o.index = 0

	bas := structPointer_FieldPointer(base, p.field)

	err = o.unmarshalType(p.stype, p.sprop, false, bas)
	o.buf = obuf
	o.index = oi

	return err
***REMOVED***

// Decode a slice of references to struct pointers ([]struct).
func (o *Buffer) dec_slice_ref_struct(p *Properties, is_group bool, base structPointer) error ***REMOVED***
	newBas := appendStructPointer(base, p.field, p.sstype)

	if is_group ***REMOVED***
		panic("not supported, maybe in future, if requested.")
	***REMOVED***

	raw, err := o.DecodeRawBytes(false)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// If the object can unmarshal itself, let it.
	if p.isUnmarshaler ***REMOVED***
		panic("not supported, since this is not a pointer receiver.")
	***REMOVED***

	obuf := o.buf
	oi := o.index
	o.buf = raw
	o.index = 0

	err = o.unmarshalType(p.stype, p.sprop, is_group, newBas)

	o.buf = obuf
	o.index = oi

	return err
***REMOVED***

// Decode a slice of references to struct pointers.
func (o *Buffer) dec_slice_ref_struct_message(p *Properties, base structPointer) error ***REMOVED***
	return o.dec_slice_ref_struct(p, false, base)
***REMOVED***

func setPtrCustomType(base structPointer, f field, v interface***REMOVED******REMOVED***) ***REMOVED***
	if v == nil ***REMOVED***
		return
	***REMOVED***
	structPointer_SetStructPointer(base, f, toStructPointer(reflect.ValueOf(v)))
***REMOVED***

func setCustomType(base structPointer, f field, value interface***REMOVED******REMOVED***) ***REMOVED***
	if value == nil ***REMOVED***
		return
	***REMOVED***
	v := reflect.ValueOf(value).Elem()
	t := reflect.TypeOf(value).Elem()
	kind := t.Kind()
	switch kind ***REMOVED***
	case reflect.Slice:
		slice := reflect.MakeSlice(t, v.Len(), v.Cap())
		reflect.Copy(slice, v)
		oldHeader := structPointer_GetSliceHeader(base, f)
		oldHeader.Data = slice.Pointer()
		oldHeader.Len = v.Len()
		oldHeader.Cap = v.Cap()
	default:
		size := reflect.TypeOf(value).Elem().Size()
		structPointer_Copy(toStructPointer(reflect.ValueOf(value)), structPointer_Add(base, f), int(size))
	***REMOVED***
***REMOVED***

func (o *Buffer) dec_custom_bytes(p *Properties, base structPointer) error ***REMOVED***
	b, err := o.DecodeRawBytes(true)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	i := reflect.New(p.ctype.Elem()).Interface()
	custom := (i).(Unmarshaler)
	if err := custom.Unmarshal(b); err != nil ***REMOVED***
		return err
	***REMOVED***
	setPtrCustomType(base, p.field, custom)
	return nil
***REMOVED***

func (o *Buffer) dec_custom_ref_bytes(p *Properties, base structPointer) error ***REMOVED***
	b, err := o.DecodeRawBytes(true)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	i := reflect.New(p.ctype).Interface()
	custom := (i).(Unmarshaler)
	if err := custom.Unmarshal(b); err != nil ***REMOVED***
		return err
	***REMOVED***
	if custom != nil ***REMOVED***
		setCustomType(base, p.field, custom)
	***REMOVED***
	return nil
***REMOVED***

// Decode a slice of bytes ([]byte) into a slice of custom types.
func (o *Buffer) dec_custom_slice_bytes(p *Properties, base structPointer) error ***REMOVED***
	b, err := o.DecodeRawBytes(true)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	i := reflect.New(p.ctype.Elem()).Interface()
	custom := (i).(Unmarshaler)
	if err := custom.Unmarshal(b); err != nil ***REMOVED***
		return err
	***REMOVED***
	newBas := appendStructPointer(base, p.field, p.ctype)

	var zero field
	setCustomType(newBas, zero, custom)

	return nil
***REMOVED***
