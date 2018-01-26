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
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
)

func GetBoolExtension(pb Message, extension *ExtensionDesc, ifnotset bool) bool ***REMOVED***
	if reflect.ValueOf(pb).IsNil() ***REMOVED***
		return ifnotset
	***REMOVED***
	value, err := GetExtension(pb, extension)
	if err != nil ***REMOVED***
		return ifnotset
	***REMOVED***
	if value == nil ***REMOVED***
		return ifnotset
	***REMOVED***
	if value.(*bool) == nil ***REMOVED***
		return ifnotset
	***REMOVED***
	return *(value.(*bool))
***REMOVED***

func (this *Extension) Equal(that *Extension) bool ***REMOVED***
	return bytes.Equal(this.enc, that.enc)
***REMOVED***

func (this *Extension) Compare(that *Extension) int ***REMOVED***
	return bytes.Compare(this.enc, that.enc)
***REMOVED***

func SizeOfInternalExtension(m extendableProto) (n int) ***REMOVED***
	return SizeOfExtensionMap(m.extensionsWrite())
***REMOVED***

func SizeOfExtensionMap(m map[int32]Extension) (n int) ***REMOVED***
	return extensionsMapSize(m)
***REMOVED***

type sortableMapElem struct ***REMOVED***
	field int32
	ext   Extension
***REMOVED***

func newSortableExtensionsFromMap(m map[int32]Extension) sortableExtensions ***REMOVED***
	s := make(sortableExtensions, 0, len(m))
	for k, v := range m ***REMOVED***
		s = append(s, &sortableMapElem***REMOVED***field: k, ext: v***REMOVED***)
	***REMOVED***
	return s
***REMOVED***

type sortableExtensions []*sortableMapElem

func (this sortableExtensions) Len() int ***REMOVED*** return len(this) ***REMOVED***

func (this sortableExtensions) Swap(i, j int) ***REMOVED*** this[i], this[j] = this[j], this[i] ***REMOVED***

func (this sortableExtensions) Less(i, j int) bool ***REMOVED*** return this[i].field < this[j].field ***REMOVED***

func (this sortableExtensions) String() string ***REMOVED***
	sort.Sort(this)
	ss := make([]string, len(this))
	for i := range this ***REMOVED***
		ss[i] = fmt.Sprintf("%d: %v", this[i].field, this[i].ext)
	***REMOVED***
	return "map[" + strings.Join(ss, ",") + "]"
***REMOVED***

func StringFromInternalExtension(m extendableProto) string ***REMOVED***
	return StringFromExtensionsMap(m.extensionsWrite())
***REMOVED***

func StringFromExtensionsMap(m map[int32]Extension) string ***REMOVED***
	return newSortableExtensionsFromMap(m).String()
***REMOVED***

func StringFromExtensionsBytes(ext []byte) string ***REMOVED***
	m, err := BytesToExtensionsMap(ext)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return StringFromExtensionsMap(m)
***REMOVED***

func EncodeInternalExtension(m extendableProto, data []byte) (n int, err error) ***REMOVED***
	return EncodeExtensionMap(m.extensionsWrite(), data)
***REMOVED***

func EncodeExtensionMap(m map[int32]Extension, data []byte) (n int, err error) ***REMOVED***
	if err := encodeExtensionsMap(m); err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	keys := make([]int, 0, len(m))
	for k := range m ***REMOVED***
		keys = append(keys, int(k))
	***REMOVED***
	sort.Ints(keys)
	for _, k := range keys ***REMOVED***
		n += copy(data[n:], m[int32(k)].enc)
	***REMOVED***
	return n, nil
***REMOVED***

func GetRawExtension(m map[int32]Extension, id int32) ([]byte, error) ***REMOVED***
	if m[id].value == nil || m[id].desc == nil ***REMOVED***
		return m[id].enc, nil
	***REMOVED***
	if err := encodeExtensionsMap(m); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return m[id].enc, nil
***REMOVED***

func size(buf []byte, wire int) (int, error) ***REMOVED***
	switch wire ***REMOVED***
	case WireVarint:
		_, n := DecodeVarint(buf)
		return n, nil
	case WireFixed64:
		return 8, nil
	case WireBytes:
		v, n := DecodeVarint(buf)
		return int(v) + n, nil
	case WireFixed32:
		return 4, nil
	case WireStartGroup:
		offset := 0
		for ***REMOVED***
			u, n := DecodeVarint(buf[offset:])
			fwire := int(u & 0x7)
			offset += n
			if fwire == WireEndGroup ***REMOVED***
				return offset, nil
			***REMOVED***
			s, err := size(buf[offset:], wire)
			if err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			offset += s
		***REMOVED***
	***REMOVED***
	return 0, fmt.Errorf("proto: can't get size for unknown wire type %d", wire)
***REMOVED***

func BytesToExtensionsMap(buf []byte) (map[int32]Extension, error) ***REMOVED***
	m := make(map[int32]Extension)
	i := 0
	for i < len(buf) ***REMOVED***
		tag, n := DecodeVarint(buf[i:])
		if n <= 0 ***REMOVED***
			return nil, fmt.Errorf("unable to decode varint")
		***REMOVED***
		fieldNum := int32(tag >> 3)
		wireType := int(tag & 0x7)
		l, err := size(buf[i+n:], wireType)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		end := i + int(l) + n
		m[int32(fieldNum)] = Extension***REMOVED***enc: buf[i:end]***REMOVED***
		i = end
	***REMOVED***
	return m, nil
***REMOVED***

func NewExtension(e []byte) Extension ***REMOVED***
	ee := Extension***REMOVED***enc: make([]byte, len(e))***REMOVED***
	copy(ee.enc, e)
	return ee
***REMOVED***

func AppendExtension(e Message, tag int32, buf []byte) ***REMOVED***
	if ee, eok := e.(extensionsBytes); eok ***REMOVED***
		ext := ee.GetExtensions()
		*ext = append(*ext, buf...)
		return
	***REMOVED***
	if ee, eok := e.(extendableProto); eok ***REMOVED***
		m := ee.extensionsWrite()
		ext := m[int32(tag)] // may be missing
		ext.enc = append(ext.enc, buf...)
		m[int32(tag)] = ext
	***REMOVED***
***REMOVED***

func encodeExtension(e *Extension) error ***REMOVED***
	if e.value == nil || e.desc == nil ***REMOVED***
		// Extension is only in its encoded form.
		return nil
	***REMOVED***
	// We don't skip extensions that have an encoded form set,
	// because the extension value may have been mutated after
	// the last time this function was called.

	et := reflect.TypeOf(e.desc.ExtensionType)
	props := extensionProperties(e.desc)

	p := NewBuffer(nil)
	// If e.value has type T, the encoder expects a *struct***REMOVED*** X T ***REMOVED***.
	// Pass a *T with a zero field and hope it all works out.
	x := reflect.New(et)
	x.Elem().Set(reflect.ValueOf(e.value))
	if err := props.enc(p, props, toStructPointer(x)); err != nil ***REMOVED***
		return err
	***REMOVED***
	e.enc = p.buf
	return nil
***REMOVED***

func (this Extension) GoString() string ***REMOVED***
	if this.enc == nil ***REMOVED***
		if err := encodeExtension(&this); err != nil ***REMOVED***
			panic(err)
		***REMOVED***
	***REMOVED***
	return fmt.Sprintf("proto.NewExtension(%#v)", this.enc)
***REMOVED***

func SetUnsafeExtension(pb Message, fieldNum int32, value interface***REMOVED******REMOVED***) error ***REMOVED***
	typ := reflect.TypeOf(pb).Elem()
	ext, ok := extensionMaps[typ]
	if !ok ***REMOVED***
		return fmt.Errorf("proto: bad extended type; %s is not extendable", typ.String())
	***REMOVED***
	desc, ok := ext[fieldNum]
	if !ok ***REMOVED***
		return errors.New("proto: bad extension number; not in declared ranges")
	***REMOVED***
	return SetExtension(pb, desc, value)
***REMOVED***

func GetUnsafeExtension(pb Message, fieldNum int32) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	typ := reflect.TypeOf(pb).Elem()
	ext, ok := extensionMaps[typ]
	if !ok ***REMOVED***
		return nil, fmt.Errorf("proto: bad extended type; %s is not extendable", typ.String())
	***REMOVED***
	desc, ok := ext[fieldNum]
	if !ok ***REMOVED***
		return nil, fmt.Errorf("unregistered field number %d", fieldNum)
	***REMOVED***
	return GetExtension(pb, desc)
***REMOVED***

func NewUnsafeXXX_InternalExtensions(m map[int32]Extension) XXX_InternalExtensions ***REMOVED***
	x := &XXX_InternalExtensions***REMOVED***
		p: new(struct ***REMOVED***
			mu           sync.Mutex
			extensionMap map[int32]Extension
		***REMOVED***),
	***REMOVED***
	x.p.extensionMap = m
	return *x
***REMOVED***

func GetUnsafeExtensionsMap(extendable Message) map[int32]Extension ***REMOVED***
	pb := extendable.(extendableProto)
	return pb.extensionsWrite()
***REMOVED***
