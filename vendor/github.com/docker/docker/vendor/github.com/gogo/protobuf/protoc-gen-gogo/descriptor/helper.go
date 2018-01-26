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

package descriptor

import (
	"strings"
)

func (msg *DescriptorProto) GetMapFields() (*FieldDescriptorProto, *FieldDescriptorProto) ***REMOVED***
	if !msg.GetOptions().GetMapEntry() ***REMOVED***
		return nil, nil
	***REMOVED***
	return msg.GetField()[0], msg.GetField()[1]
***REMOVED***

func dotToUnderscore(r rune) rune ***REMOVED***
	if r == '.' ***REMOVED***
		return '_'
	***REMOVED***
	return r
***REMOVED***

func (field *FieldDescriptorProto) WireType() (wire int) ***REMOVED***
	switch *field.Type ***REMOVED***
	case FieldDescriptorProto_TYPE_DOUBLE:
		return 1
	case FieldDescriptorProto_TYPE_FLOAT:
		return 5
	case FieldDescriptorProto_TYPE_INT64:
		return 0
	case FieldDescriptorProto_TYPE_UINT64:
		return 0
	case FieldDescriptorProto_TYPE_INT32:
		return 0
	case FieldDescriptorProto_TYPE_UINT32:
		return 0
	case FieldDescriptorProto_TYPE_FIXED64:
		return 1
	case FieldDescriptorProto_TYPE_FIXED32:
		return 5
	case FieldDescriptorProto_TYPE_BOOL:
		return 0
	case FieldDescriptorProto_TYPE_STRING:
		return 2
	case FieldDescriptorProto_TYPE_GROUP:
		return 2
	case FieldDescriptorProto_TYPE_MESSAGE:
		return 2
	case FieldDescriptorProto_TYPE_BYTES:
		return 2
	case FieldDescriptorProto_TYPE_ENUM:
		return 0
	case FieldDescriptorProto_TYPE_SFIXED32:
		return 5
	case FieldDescriptorProto_TYPE_SFIXED64:
		return 1
	case FieldDescriptorProto_TYPE_SINT32:
		return 0
	case FieldDescriptorProto_TYPE_SINT64:
		return 0
	***REMOVED***
	panic("unreachable")
***REMOVED***

func (field *FieldDescriptorProto) GetKeyUint64() (x uint64) ***REMOVED***
	packed := field.IsPacked()
	wireType := field.WireType()
	fieldNumber := field.GetNumber()
	if packed ***REMOVED***
		wireType = 2
	***REMOVED***
	x = uint64(uint32(fieldNumber)<<3 | uint32(wireType))
	return x
***REMOVED***

func (field *FieldDescriptorProto) GetKey3Uint64() (x uint64) ***REMOVED***
	packed := field.IsPacked3()
	wireType := field.WireType()
	fieldNumber := field.GetNumber()
	if packed ***REMOVED***
		wireType = 2
	***REMOVED***
	x = uint64(uint32(fieldNumber)<<3 | uint32(wireType))
	return x
***REMOVED***

func (field *FieldDescriptorProto) GetKey() []byte ***REMOVED***
	x := field.GetKeyUint64()
	i := 0
	keybuf := make([]byte, 0)
	for i = 0; x > 127; i++ ***REMOVED***
		keybuf = append(keybuf, 0x80|uint8(x&0x7F))
		x >>= 7
	***REMOVED***
	keybuf = append(keybuf, uint8(x))
	return keybuf
***REMOVED***

func (field *FieldDescriptorProto) GetKey3() []byte ***REMOVED***
	x := field.GetKey3Uint64()
	i := 0
	keybuf := make([]byte, 0)
	for i = 0; x > 127; i++ ***REMOVED***
		keybuf = append(keybuf, 0x80|uint8(x&0x7F))
		x >>= 7
	***REMOVED***
	keybuf = append(keybuf, uint8(x))
	return keybuf
***REMOVED***

func (desc *FileDescriptorSet) GetField(packageName, messageName, fieldName string) *FieldDescriptorProto ***REMOVED***
	msg := desc.GetMessage(packageName, messageName)
	if msg == nil ***REMOVED***
		return nil
	***REMOVED***
	for _, field := range msg.GetField() ***REMOVED***
		if field.GetName() == fieldName ***REMOVED***
			return field
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (file *FileDescriptorProto) GetMessage(typeName string) *DescriptorProto ***REMOVED***
	for _, msg := range file.GetMessageType() ***REMOVED***
		if msg.GetName() == typeName ***REMOVED***
			return msg
		***REMOVED***
		nes := file.GetNestedMessage(msg, strings.TrimPrefix(typeName, msg.GetName()+"."))
		if nes != nil ***REMOVED***
			return nes
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (file *FileDescriptorProto) GetNestedMessage(msg *DescriptorProto, typeName string) *DescriptorProto ***REMOVED***
	for _, nes := range msg.GetNestedType() ***REMOVED***
		if nes.GetName() == typeName ***REMOVED***
			return nes
		***REMOVED***
		res := file.GetNestedMessage(nes, strings.TrimPrefix(typeName, nes.GetName()+"."))
		if res != nil ***REMOVED***
			return res
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (desc *FileDescriptorSet) GetMessage(packageName string, typeName string) *DescriptorProto ***REMOVED***
	for _, file := range desc.GetFile() ***REMOVED***
		if strings.Map(dotToUnderscore, file.GetPackage()) != strings.Map(dotToUnderscore, packageName) ***REMOVED***
			continue
		***REMOVED***
		for _, msg := range file.GetMessageType() ***REMOVED***
			if msg.GetName() == typeName ***REMOVED***
				return msg
			***REMOVED***
		***REMOVED***
		for _, msg := range file.GetMessageType() ***REMOVED***
			for _, nes := range msg.GetNestedType() ***REMOVED***
				if nes.GetName() == typeName ***REMOVED***
					return nes
				***REMOVED***
				if msg.GetName()+"."+nes.GetName() == typeName ***REMOVED***
					return nes
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (desc *FileDescriptorSet) IsProto3(packageName string, typeName string) bool ***REMOVED***
	for _, file := range desc.GetFile() ***REMOVED***
		if strings.Map(dotToUnderscore, file.GetPackage()) != strings.Map(dotToUnderscore, packageName) ***REMOVED***
			continue
		***REMOVED***
		for _, msg := range file.GetMessageType() ***REMOVED***
			if msg.GetName() == typeName ***REMOVED***
				return file.GetSyntax() == "proto3"
			***REMOVED***
		***REMOVED***
		for _, msg := range file.GetMessageType() ***REMOVED***
			for _, nes := range msg.GetNestedType() ***REMOVED***
				if nes.GetName() == typeName ***REMOVED***
					return file.GetSyntax() == "proto3"
				***REMOVED***
				if msg.GetName()+"."+nes.GetName() == typeName ***REMOVED***
					return file.GetSyntax() == "proto3"
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (msg *DescriptorProto) IsExtendable() bool ***REMOVED***
	return len(msg.GetExtensionRange()) > 0
***REMOVED***

func (desc *FileDescriptorSet) FindExtension(packageName string, typeName string, fieldName string) (extPackageName string, field *FieldDescriptorProto) ***REMOVED***
	parent := desc.GetMessage(packageName, typeName)
	if parent == nil ***REMOVED***
		return "", nil
	***REMOVED***
	if !parent.IsExtendable() ***REMOVED***
		return "", nil
	***REMOVED***
	extendee := "." + packageName + "." + typeName
	for _, file := range desc.GetFile() ***REMOVED***
		for _, ext := range file.GetExtension() ***REMOVED***
			if strings.Map(dotToUnderscore, file.GetPackage()) == strings.Map(dotToUnderscore, packageName) ***REMOVED***
				if !(ext.GetExtendee() == typeName || ext.GetExtendee() == extendee) ***REMOVED***
					continue
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if ext.GetExtendee() != extendee ***REMOVED***
					continue
				***REMOVED***
			***REMOVED***
			if ext.GetName() == fieldName ***REMOVED***
				return file.GetPackage(), ext
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return "", nil
***REMOVED***

func (desc *FileDescriptorSet) FindExtensionByFieldNumber(packageName string, typeName string, fieldNum int32) (extPackageName string, field *FieldDescriptorProto) ***REMOVED***
	parent := desc.GetMessage(packageName, typeName)
	if parent == nil ***REMOVED***
		return "", nil
	***REMOVED***
	if !parent.IsExtendable() ***REMOVED***
		return "", nil
	***REMOVED***
	extendee := "." + packageName + "." + typeName
	for _, file := range desc.GetFile() ***REMOVED***
		for _, ext := range file.GetExtension() ***REMOVED***
			if strings.Map(dotToUnderscore, file.GetPackage()) == strings.Map(dotToUnderscore, packageName) ***REMOVED***
				if !(ext.GetExtendee() == typeName || ext.GetExtendee() == extendee) ***REMOVED***
					continue
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if ext.GetExtendee() != extendee ***REMOVED***
					continue
				***REMOVED***
			***REMOVED***
			if ext.GetNumber() == fieldNum ***REMOVED***
				return file.GetPackage(), ext
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return "", nil
***REMOVED***

func (desc *FileDescriptorSet) FindMessage(packageName string, typeName string, fieldName string) (msgPackageName string, msgName string) ***REMOVED***
	parent := desc.GetMessage(packageName, typeName)
	if parent == nil ***REMOVED***
		return "", ""
	***REMOVED***
	field := parent.GetFieldDescriptor(fieldName)
	if field == nil ***REMOVED***
		var extPackageName string
		extPackageName, field = desc.FindExtension(packageName, typeName, fieldName)
		if field == nil ***REMOVED***
			return "", ""
		***REMOVED***
		packageName = extPackageName
	***REMOVED***
	typeNames := strings.Split(field.GetTypeName(), ".")
	if len(typeNames) == 1 ***REMOVED***
		msg := desc.GetMessage(packageName, typeName)
		if msg == nil ***REMOVED***
			return "", ""
		***REMOVED***
		return packageName, msg.GetName()
	***REMOVED***
	if len(typeNames) > 2 ***REMOVED***
		for i := 1; i < len(typeNames)-1; i++ ***REMOVED***
			packageName = strings.Join(typeNames[1:len(typeNames)-i], ".")
			typeName = strings.Join(typeNames[len(typeNames)-i:], ".")
			msg := desc.GetMessage(packageName, typeName)
			if msg != nil ***REMOVED***
				typeNames := strings.Split(msg.GetName(), ".")
				if len(typeNames) == 1 ***REMOVED***
					return packageName, msg.GetName()
				***REMOVED***
				return strings.Join(typeNames[1:len(typeNames)-1], "."), typeNames[len(typeNames)-1]
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return "", ""
***REMOVED***

func (msg *DescriptorProto) GetFieldDescriptor(fieldName string) *FieldDescriptorProto ***REMOVED***
	for _, field := range msg.GetField() ***REMOVED***
		if field.GetName() == fieldName ***REMOVED***
			return field
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (desc *FileDescriptorSet) GetEnum(packageName string, typeName string) *EnumDescriptorProto ***REMOVED***
	for _, file := range desc.GetFile() ***REMOVED***
		if strings.Map(dotToUnderscore, file.GetPackage()) != strings.Map(dotToUnderscore, packageName) ***REMOVED***
			continue
		***REMOVED***
		for _, enum := range file.GetEnumType() ***REMOVED***
			if enum.GetName() == typeName ***REMOVED***
				return enum
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (f *FieldDescriptorProto) IsEnum() bool ***REMOVED***
	return *f.Type == FieldDescriptorProto_TYPE_ENUM
***REMOVED***

func (f *FieldDescriptorProto) IsMessage() bool ***REMOVED***
	return *f.Type == FieldDescriptorProto_TYPE_MESSAGE
***REMOVED***

func (f *FieldDescriptorProto) IsBytes() bool ***REMOVED***
	return *f.Type == FieldDescriptorProto_TYPE_BYTES
***REMOVED***

func (f *FieldDescriptorProto) IsRepeated() bool ***REMOVED***
	return f.Label != nil && *f.Label == FieldDescriptorProto_LABEL_REPEATED
***REMOVED***

func (f *FieldDescriptorProto) IsString() bool ***REMOVED***
	return *f.Type == FieldDescriptorProto_TYPE_STRING
***REMOVED***

func (f *FieldDescriptorProto) IsBool() bool ***REMOVED***
	return *f.Type == FieldDescriptorProto_TYPE_BOOL
***REMOVED***

func (f *FieldDescriptorProto) IsRequired() bool ***REMOVED***
	return f.Label != nil && *f.Label == FieldDescriptorProto_LABEL_REQUIRED
***REMOVED***

func (f *FieldDescriptorProto) IsPacked() bool ***REMOVED***
	return f.Options != nil && f.GetOptions().GetPacked()
***REMOVED***

func (f *FieldDescriptorProto) IsPacked3() bool ***REMOVED***
	if f.IsRepeated() && f.IsScalar() ***REMOVED***
		if f.Options == nil || f.GetOptions().Packed == nil ***REMOVED***
			return true
		***REMOVED***
		return f.Options != nil && f.GetOptions().GetPacked()
	***REMOVED***
	return false
***REMOVED***

func (m *DescriptorProto) HasExtension() bool ***REMOVED***
	return len(m.ExtensionRange) > 0
***REMOVED***
