// Protocol Buffers for Go with Gadgets
//
// Copyright (c) 2013, The GoGo Authors. All rights reserved.
// http://github.com/gogo/protobuf
//
// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2010 The Go Authors.  All rights reserved.
// http://github.com/golang/protobuf/
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
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
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

func NewRequiredNotSetError(field string) *RequiredNotSetError ***REMOVED***
	return &RequiredNotSetError***REMOVED***field***REMOVED***
***REMOVED***

type Sizer interface ***REMOVED***
	Size() int
***REMOVED***

func (o *Buffer) enc_ext_slice_byte(p *Properties, base structPointer) error ***REMOVED***
	s := *structPointer_Bytes(base, p.field)
	if s == nil ***REMOVED***
		return ErrNil
	***REMOVED***
	o.buf = append(o.buf, s...)
	return nil
***REMOVED***

func size_ext_slice_byte(p *Properties, base structPointer) (n int) ***REMOVED***
	s := *structPointer_Bytes(base, p.field)
	if s == nil ***REMOVED***
		return 0
	***REMOVED***
	n += len(s)
	return
***REMOVED***

// Encode a reference to bool pointer.
func (o *Buffer) enc_ref_bool(p *Properties, base structPointer) error ***REMOVED***
	v := *structPointer_BoolVal(base, p.field)
	x := 0
	if v ***REMOVED***
		x = 1
	***REMOVED***
	o.buf = append(o.buf, p.tagcode...)
	p.valEnc(o, uint64(x))
	return nil
***REMOVED***

func size_ref_bool(p *Properties, base structPointer) int ***REMOVED***
	return len(p.tagcode) + 1 // each bool takes exactly one byte
***REMOVED***

// Encode a reference to int32 pointer.
func (o *Buffer) enc_ref_int32(p *Properties, base structPointer) error ***REMOVED***
	v := structPointer_Word32Val(base, p.field)
	x := int32(word32Val_Get(v))
	o.buf = append(o.buf, p.tagcode...)
	p.valEnc(o, uint64(x))
	return nil
***REMOVED***

func size_ref_int32(p *Properties, base structPointer) (n int) ***REMOVED***
	v := structPointer_Word32Val(base, p.field)
	x := int32(word32Val_Get(v))
	n += len(p.tagcode)
	n += p.valSize(uint64(x))
	return
***REMOVED***

func (o *Buffer) enc_ref_uint32(p *Properties, base structPointer) error ***REMOVED***
	v := structPointer_Word32Val(base, p.field)
	x := word32Val_Get(v)
	o.buf = append(o.buf, p.tagcode...)
	p.valEnc(o, uint64(x))
	return nil
***REMOVED***

func size_ref_uint32(p *Properties, base structPointer) (n int) ***REMOVED***
	v := structPointer_Word32Val(base, p.field)
	x := word32Val_Get(v)
	n += len(p.tagcode)
	n += p.valSize(uint64(x))
	return
***REMOVED***

// Encode a reference to an int64 pointer.
func (o *Buffer) enc_ref_int64(p *Properties, base structPointer) error ***REMOVED***
	v := structPointer_Word64Val(base, p.field)
	x := word64Val_Get(v)
	o.buf = append(o.buf, p.tagcode...)
	p.valEnc(o, x)
	return nil
***REMOVED***

func size_ref_int64(p *Properties, base structPointer) (n int) ***REMOVED***
	v := structPointer_Word64Val(base, p.field)
	x := word64Val_Get(v)
	n += len(p.tagcode)
	n += p.valSize(x)
	return
***REMOVED***

// Encode a reference to a string pointer.
func (o *Buffer) enc_ref_string(p *Properties, base structPointer) error ***REMOVED***
	v := *structPointer_StringVal(base, p.field)
	o.buf = append(o.buf, p.tagcode...)
	o.EncodeStringBytes(v)
	return nil
***REMOVED***

func size_ref_string(p *Properties, base structPointer) (n int) ***REMOVED***
	v := *structPointer_StringVal(base, p.field)
	n += len(p.tagcode)
	n += sizeStringBytes(v)
	return
***REMOVED***

// Encode a reference to a message struct.
func (o *Buffer) enc_ref_struct_message(p *Properties, base structPointer) error ***REMOVED***
	var state errorState
	structp := structPointer_GetRefStructPointer(base, p.field)
	if structPointer_IsNil(structp) ***REMOVED***
		return ErrNil
	***REMOVED***

	// Can the object marshal itself?
	if p.isMarshaler ***REMOVED***
		m := structPointer_Interface(structp, p.stype).(Marshaler)
		data, err := m.Marshal()
		if err != nil && !state.shouldContinue(err, nil) ***REMOVED***
			return err
		***REMOVED***
		o.buf = append(o.buf, p.tagcode...)
		o.EncodeRawBytes(data)
		return nil
	***REMOVED***

	o.buf = append(o.buf, p.tagcode...)
	return o.enc_len_struct(p.sprop, structp, &state)
***REMOVED***

//TODO this is only copied, please fix this
func size_ref_struct_message(p *Properties, base structPointer) int ***REMOVED***
	structp := structPointer_GetRefStructPointer(base, p.field)
	if structPointer_IsNil(structp) ***REMOVED***
		return 0
	***REMOVED***

	// Can the object marshal itself?
	if p.isMarshaler ***REMOVED***
		m := structPointer_Interface(structp, p.stype).(Marshaler)
		data, _ := m.Marshal()
		n0 := len(p.tagcode)
		n1 := sizeRawBytes(data)
		return n0 + n1
	***REMOVED***

	n0 := len(p.tagcode)
	n1 := size_struct(p.sprop, structp)
	n2 := sizeVarint(uint64(n1)) // size of encoded length
	return n0 + n1 + n2
***REMOVED***

// Encode a slice of references to message struct pointers ([]struct).
func (o *Buffer) enc_slice_ref_struct_message(p *Properties, base structPointer) error ***REMOVED***
	var state errorState
	ss := structPointer_StructRefSlice(base, p.field, p.stype.Size())
	l := ss.Len()
	for i := 0; i < l; i++ ***REMOVED***
		structp := ss.Index(i)
		if structPointer_IsNil(structp) ***REMOVED***
			return errRepeatedHasNil
		***REMOVED***

		// Can the object marshal itself?
		if p.isMarshaler ***REMOVED***
			m := structPointer_Interface(structp, p.stype).(Marshaler)
			data, err := m.Marshal()
			if err != nil && !state.shouldContinue(err, nil) ***REMOVED***
				return err
			***REMOVED***
			o.buf = append(o.buf, p.tagcode...)
			o.EncodeRawBytes(data)
			continue
		***REMOVED***

		o.buf = append(o.buf, p.tagcode...)
		err := o.enc_len_struct(p.sprop, structp, &state)
		if err != nil && !state.shouldContinue(err, nil) ***REMOVED***
			if err == ErrNil ***REMOVED***
				return errRepeatedHasNil
			***REMOVED***
			return err
		***REMOVED***

	***REMOVED***
	return state.err
***REMOVED***

//TODO this is only copied, please fix this
func size_slice_ref_struct_message(p *Properties, base structPointer) (n int) ***REMOVED***
	ss := structPointer_StructRefSlice(base, p.field, p.stype.Size())
	l := ss.Len()
	n += l * len(p.tagcode)
	for i := 0; i < l; i++ ***REMOVED***
		structp := ss.Index(i)
		if structPointer_IsNil(structp) ***REMOVED***
			return // return the size up to this point
		***REMOVED***

		// Can the object marshal itself?
		if p.isMarshaler ***REMOVED***
			m := structPointer_Interface(structp, p.stype).(Marshaler)
			data, _ := m.Marshal()
			n += len(p.tagcode)
			n += sizeRawBytes(data)
			continue
		***REMOVED***

		n0 := size_struct(p.sprop, structp)
		n1 := sizeVarint(uint64(n0)) // size of encoded length
		n += n0 + n1
	***REMOVED***
	return
***REMOVED***

func (o *Buffer) enc_custom_bytes(p *Properties, base structPointer) error ***REMOVED***
	i := structPointer_InterfaceRef(base, p.field, p.ctype)
	if i == nil ***REMOVED***
		return ErrNil
	***REMOVED***
	custom := i.(Marshaler)
	data, err := custom.Marshal()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if data == nil ***REMOVED***
		return ErrNil
	***REMOVED***
	o.buf = append(o.buf, p.tagcode...)
	o.EncodeRawBytes(data)
	return nil
***REMOVED***

func size_custom_bytes(p *Properties, base structPointer) (n int) ***REMOVED***
	n += len(p.tagcode)
	i := structPointer_InterfaceRef(base, p.field, p.ctype)
	if i == nil ***REMOVED***
		return 0
	***REMOVED***
	custom := i.(Marshaler)
	data, _ := custom.Marshal()
	n += sizeRawBytes(data)
	return
***REMOVED***

func (o *Buffer) enc_custom_ref_bytes(p *Properties, base structPointer) error ***REMOVED***
	custom := structPointer_InterfaceAt(base, p.field, p.ctype).(Marshaler)
	data, err := custom.Marshal()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if data == nil ***REMOVED***
		return ErrNil
	***REMOVED***
	o.buf = append(o.buf, p.tagcode...)
	o.EncodeRawBytes(data)
	return nil
***REMOVED***

func size_custom_ref_bytes(p *Properties, base structPointer) (n int) ***REMOVED***
	n += len(p.tagcode)
	i := structPointer_InterfaceAt(base, p.field, p.ctype)
	if i == nil ***REMOVED***
		return 0
	***REMOVED***
	custom := i.(Marshaler)
	data, _ := custom.Marshal()
	n += sizeRawBytes(data)
	return
***REMOVED***

func (o *Buffer) enc_custom_slice_bytes(p *Properties, base structPointer) error ***REMOVED***
	inter := structPointer_InterfaceRef(base, p.field, p.ctype)
	if inter == nil ***REMOVED***
		return ErrNil
	***REMOVED***
	slice := reflect.ValueOf(inter)
	l := slice.Len()
	for i := 0; i < l; i++ ***REMOVED***
		v := slice.Index(i)
		custom := v.Interface().(Marshaler)
		data, err := custom.Marshal()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		o.buf = append(o.buf, p.tagcode...)
		o.EncodeRawBytes(data)
	***REMOVED***
	return nil
***REMOVED***

func size_custom_slice_bytes(p *Properties, base structPointer) (n int) ***REMOVED***
	inter := structPointer_InterfaceRef(base, p.field, p.ctype)
	if inter == nil ***REMOVED***
		return 0
	***REMOVED***
	slice := reflect.ValueOf(inter)
	l := slice.Len()
	n += l * len(p.tagcode)
	for i := 0; i < l; i++ ***REMOVED***
		v := slice.Index(i)
		custom := v.Interface().(Marshaler)
		data, _ := custom.Marshal()
		n += sizeRawBytes(data)
	***REMOVED***
	return
***REMOVED***
