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

var timeType = reflect.TypeOf((*time.Time)(nil)).Elem()

type timestamp struct ***REMOVED***
	Seconds int64 `protobuf:"varint,1,opt,name=seconds,proto3" json:"seconds,omitempty"`
	Nanos   int32 `protobuf:"varint,2,opt,name=nanos,proto3" json:"nanos,omitempty"`
***REMOVED***

func (m *timestamp) Reset()       ***REMOVED*** *m = timestamp***REMOVED******REMOVED*** ***REMOVED***
func (*timestamp) ProtoMessage()  ***REMOVED******REMOVED***
func (*timestamp) String() string ***REMOVED*** return "timestamp<string>" ***REMOVED***

func init() ***REMOVED***
	RegisterType((*timestamp)(nil), "gogo.protobuf.proto.timestamp")
***REMOVED***

func (o *Buffer) decTimestamp() (time.Time, error) ***REMOVED***
	b, err := o.DecodeRawBytes(true)
	if err != nil ***REMOVED***
		return time.Time***REMOVED******REMOVED***, err
	***REMOVED***
	tproto := &timestamp***REMOVED******REMOVED***
	if err := Unmarshal(b, tproto); err != nil ***REMOVED***
		return time.Time***REMOVED******REMOVED***, err
	***REMOVED***
	return timestampFromProto(tproto)
***REMOVED***

func (o *Buffer) dec_time(p *Properties, base structPointer) error ***REMOVED***
	t, err := o.decTimestamp()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	setPtrCustomType(base, p.field, &t)
	return nil
***REMOVED***

func (o *Buffer) dec_ref_time(p *Properties, base structPointer) error ***REMOVED***
	t, err := o.decTimestamp()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	setCustomType(base, p.field, &t)
	return nil
***REMOVED***

func (o *Buffer) dec_slice_time(p *Properties, base structPointer) error ***REMOVED***
	t, err := o.decTimestamp()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	newBas := appendStructPointer(base, p.field, reflect.SliceOf(reflect.PtrTo(timeType)))
	var zero field
	setPtrCustomType(newBas, zero, &t)
	return nil
***REMOVED***

func (o *Buffer) dec_slice_ref_time(p *Properties, base structPointer) error ***REMOVED***
	t, err := o.decTimestamp()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	newBas := appendStructPointer(base, p.field, reflect.SliceOf(timeType))
	var zero field
	setCustomType(newBas, zero, &t)
	return nil
***REMOVED***

func size_time(p *Properties, base structPointer) (n int) ***REMOVED***
	structp := structPointer_GetStructPointer(base, p.field)
	if structPointer_IsNil(structp) ***REMOVED***
		return 0
	***REMOVED***
	tim := structPointer_Interface(structp, timeType).(*time.Time)
	t, err := timestampProto(*tim)
	if err != nil ***REMOVED***
		return 0
	***REMOVED***
	size := Size(t)
	return size + sizeVarint(uint64(size)) + len(p.tagcode)
***REMOVED***

func (o *Buffer) enc_time(p *Properties, base structPointer) error ***REMOVED***
	structp := structPointer_GetStructPointer(base, p.field)
	if structPointer_IsNil(structp) ***REMOVED***
		return ErrNil
	***REMOVED***
	tim := structPointer_Interface(structp, timeType).(*time.Time)
	t, err := timestampProto(*tim)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	data, err := Marshal(t)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	o.buf = append(o.buf, p.tagcode...)
	o.EncodeRawBytes(data)
	return nil
***REMOVED***

func size_ref_time(p *Properties, base structPointer) (n int) ***REMOVED***
	tim := structPointer_InterfaceAt(base, p.field, timeType).(*time.Time)
	t, err := timestampProto(*tim)
	if err != nil ***REMOVED***
		return 0
	***REMOVED***
	size := Size(t)
	return size + sizeVarint(uint64(size)) + len(p.tagcode)
***REMOVED***

func (o *Buffer) enc_ref_time(p *Properties, base structPointer) error ***REMOVED***
	tim := structPointer_InterfaceAt(base, p.field, timeType).(*time.Time)
	t, err := timestampProto(*tim)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	data, err := Marshal(t)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	o.buf = append(o.buf, p.tagcode...)
	o.EncodeRawBytes(data)
	return nil
***REMOVED***

func size_slice_time(p *Properties, base structPointer) (n int) ***REMOVED***
	ptims := structPointer_InterfaceAt(base, p.field, reflect.SliceOf(reflect.PtrTo(timeType))).(*[]*time.Time)
	tims := *ptims
	for i := 0; i < len(tims); i++ ***REMOVED***
		if tims[i] == nil ***REMOVED***
			return 0
		***REMOVED***
		tproto, err := timestampProto(*tims[i])
		if err != nil ***REMOVED***
			return 0
		***REMOVED***
		size := Size(tproto)
		n += len(p.tagcode) + size + sizeVarint(uint64(size))
	***REMOVED***
	return n
***REMOVED***

func (o *Buffer) enc_slice_time(p *Properties, base structPointer) error ***REMOVED***
	ptims := structPointer_InterfaceAt(base, p.field, reflect.SliceOf(reflect.PtrTo(timeType))).(*[]*time.Time)
	tims := *ptims
	for i := 0; i < len(tims); i++ ***REMOVED***
		if tims[i] == nil ***REMOVED***
			return errRepeatedHasNil
		***REMOVED***
		tproto, err := timestampProto(*tims[i])
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		data, err := Marshal(tproto)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		o.buf = append(o.buf, p.tagcode...)
		o.EncodeRawBytes(data)
	***REMOVED***
	return nil
***REMOVED***

func size_slice_ref_time(p *Properties, base structPointer) (n int) ***REMOVED***
	ptims := structPointer_InterfaceAt(base, p.field, reflect.SliceOf(timeType)).(*[]time.Time)
	tims := *ptims
	for i := 0; i < len(tims); i++ ***REMOVED***
		tproto, err := timestampProto(tims[i])
		if err != nil ***REMOVED***
			return 0
		***REMOVED***
		size := Size(tproto)
		n += len(p.tagcode) + size + sizeVarint(uint64(size))
	***REMOVED***
	return n
***REMOVED***

func (o *Buffer) enc_slice_ref_time(p *Properties, base structPointer) error ***REMOVED***
	ptims := structPointer_InterfaceAt(base, p.field, reflect.SliceOf(timeType)).(*[]time.Time)
	tims := *ptims
	for i := 0; i < len(tims); i++ ***REMOVED***
		tproto, err := timestampProto(tims[i])
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		data, err := Marshal(tproto)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		o.buf = append(o.buf, p.tagcode...)
		o.EncodeRawBytes(data)
	***REMOVED***
	return nil
***REMOVED***
