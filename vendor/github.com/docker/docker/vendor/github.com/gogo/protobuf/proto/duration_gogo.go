// Protocol Buffers for Go with Gadgets
//
// Copyright (c) 2016, The GoGo Authors. All rights reserved.
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
	"time"
)

var durationType = reflect.TypeOf((*time.Duration)(nil)).Elem()

type duration struct ***REMOVED***
	Seconds int64 `protobuf:"varint,1,opt,name=seconds,proto3" json:"seconds,omitempty"`
	Nanos   int32 `protobuf:"varint,2,opt,name=nanos,proto3" json:"nanos,omitempty"`
***REMOVED***

func (m *duration) Reset()       ***REMOVED*** *m = duration***REMOVED******REMOVED*** ***REMOVED***
func (*duration) ProtoMessage()  ***REMOVED******REMOVED***
func (*duration) String() string ***REMOVED*** return "duration<string>" ***REMOVED***

func init() ***REMOVED***
	RegisterType((*duration)(nil), "gogo.protobuf.proto.duration")
***REMOVED***

func (o *Buffer) decDuration() (time.Duration, error) ***REMOVED***
	b, err := o.DecodeRawBytes(true)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	dproto := &duration***REMOVED******REMOVED***
	if err := Unmarshal(b, dproto); err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return durationFromProto(dproto)
***REMOVED***

func (o *Buffer) dec_duration(p *Properties, base structPointer) error ***REMOVED***
	d, err := o.decDuration()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	word64_Set(structPointer_Word64(base, p.field), o, uint64(d))
	return nil
***REMOVED***

func (o *Buffer) dec_ref_duration(p *Properties, base structPointer) error ***REMOVED***
	d, err := o.decDuration()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	word64Val_Set(structPointer_Word64Val(base, p.field), o, uint64(d))
	return nil
***REMOVED***

func (o *Buffer) dec_slice_duration(p *Properties, base structPointer) error ***REMOVED***
	d, err := o.decDuration()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	newBas := appendStructPointer(base, p.field, reflect.SliceOf(reflect.PtrTo(durationType)))
	var zero field
	setPtrCustomType(newBas, zero, &d)
	return nil
***REMOVED***

func (o *Buffer) dec_slice_ref_duration(p *Properties, base structPointer) error ***REMOVED***
	d, err := o.decDuration()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	structPointer_Word64Slice(base, p.field).Append(uint64(d))
	return nil
***REMOVED***

func size_duration(p *Properties, base structPointer) (n int) ***REMOVED***
	structp := structPointer_GetStructPointer(base, p.field)
	if structPointer_IsNil(structp) ***REMOVED***
		return 0
	***REMOVED***
	dur := structPointer_Interface(structp, durationType).(*time.Duration)
	d := durationProto(*dur)
	size := Size(d)
	return size + sizeVarint(uint64(size)) + len(p.tagcode)
***REMOVED***

func (o *Buffer) enc_duration(p *Properties, base structPointer) error ***REMOVED***
	structp := structPointer_GetStructPointer(base, p.field)
	if structPointer_IsNil(structp) ***REMOVED***
		return ErrNil
	***REMOVED***
	dur := structPointer_Interface(structp, durationType).(*time.Duration)
	d := durationProto(*dur)
	data, err := Marshal(d)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	o.buf = append(o.buf, p.tagcode...)
	o.EncodeRawBytes(data)
	return nil
***REMOVED***

func size_ref_duration(p *Properties, base structPointer) (n int) ***REMOVED***
	dur := structPointer_InterfaceAt(base, p.field, durationType).(*time.Duration)
	d := durationProto(*dur)
	size := Size(d)
	return size + sizeVarint(uint64(size)) + len(p.tagcode)
***REMOVED***

func (o *Buffer) enc_ref_duration(p *Properties, base structPointer) error ***REMOVED***
	dur := structPointer_InterfaceAt(base, p.field, durationType).(*time.Duration)
	d := durationProto(*dur)
	data, err := Marshal(d)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	o.buf = append(o.buf, p.tagcode...)
	o.EncodeRawBytes(data)
	return nil
***REMOVED***

func size_slice_duration(p *Properties, base structPointer) (n int) ***REMOVED***
	pdurs := structPointer_InterfaceAt(base, p.field, reflect.SliceOf(reflect.PtrTo(durationType))).(*[]*time.Duration)
	durs := *pdurs
	for i := 0; i < len(durs); i++ ***REMOVED***
		if durs[i] == nil ***REMOVED***
			return 0
		***REMOVED***
		dproto := durationProto(*durs[i])
		size := Size(dproto)
		n += len(p.tagcode) + size + sizeVarint(uint64(size))
	***REMOVED***
	return n
***REMOVED***

func (o *Buffer) enc_slice_duration(p *Properties, base structPointer) error ***REMOVED***
	pdurs := structPointer_InterfaceAt(base, p.field, reflect.SliceOf(reflect.PtrTo(durationType))).(*[]*time.Duration)
	durs := *pdurs
	for i := 0; i < len(durs); i++ ***REMOVED***
		if durs[i] == nil ***REMOVED***
			return errRepeatedHasNil
		***REMOVED***
		dproto := durationProto(*durs[i])
		data, err := Marshal(dproto)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		o.buf = append(o.buf, p.tagcode...)
		o.EncodeRawBytes(data)
	***REMOVED***
	return nil
***REMOVED***

func size_slice_ref_duration(p *Properties, base structPointer) (n int) ***REMOVED***
	pdurs := structPointer_InterfaceAt(base, p.field, reflect.SliceOf(durationType)).(*[]time.Duration)
	durs := *pdurs
	for i := 0; i < len(durs); i++ ***REMOVED***
		dproto := durationProto(durs[i])
		size := Size(dproto)
		n += len(p.tagcode) + size + sizeVarint(uint64(size))
	***REMOVED***
	return n
***REMOVED***

func (o *Buffer) enc_slice_ref_duration(p *Properties, base structPointer) error ***REMOVED***
	pdurs := structPointer_InterfaceAt(base, p.field, reflect.SliceOf(durationType)).(*[]time.Duration)
	durs := *pdurs
	for i := 0; i < len(durs); i++ ***REMOVED***
		dproto := durationProto(durs[i])
		data, err := Marshal(dproto)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		o.buf = append(o.buf, p.tagcode...)
		o.EncodeRawBytes(data)
	***REMOVED***
	return nil
***REMOVED***
