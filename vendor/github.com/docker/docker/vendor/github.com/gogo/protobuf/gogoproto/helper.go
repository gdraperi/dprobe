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

package gogoproto

import google_protobuf "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
import proto "github.com/gogo/protobuf/proto"

func IsEmbed(field *google_protobuf.FieldDescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(field.Options, E_Embed, false)
***REMOVED***

func IsNullable(field *google_protobuf.FieldDescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(field.Options, E_Nullable, true)
***REMOVED***

func IsStdTime(field *google_protobuf.FieldDescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(field.Options, E_Stdtime, false)
***REMOVED***

func IsStdDuration(field *google_protobuf.FieldDescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(field.Options, E_Stdduration, false)
***REMOVED***

func NeedsNilCheck(proto3 bool, field *google_protobuf.FieldDescriptorProto) bool ***REMOVED***
	nullable := IsNullable(field)
	if field.IsMessage() || IsCustomType(field) ***REMOVED***
		return nullable
	***REMOVED***
	if proto3 ***REMOVED***
		return false
	***REMOVED***
	return nullable || *field.Type == google_protobuf.FieldDescriptorProto_TYPE_BYTES
***REMOVED***

func IsCustomType(field *google_protobuf.FieldDescriptorProto) bool ***REMOVED***
	typ := GetCustomType(field)
	if len(typ) > 0 ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func IsCastType(field *google_protobuf.FieldDescriptorProto) bool ***REMOVED***
	typ := GetCastType(field)
	if len(typ) > 0 ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func IsCastKey(field *google_protobuf.FieldDescriptorProto) bool ***REMOVED***
	typ := GetCastKey(field)
	if len(typ) > 0 ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func IsCastValue(field *google_protobuf.FieldDescriptorProto) bool ***REMOVED***
	typ := GetCastValue(field)
	if len(typ) > 0 ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func HasEnumDecl(file *google_protobuf.FileDescriptorProto, enum *google_protobuf.EnumDescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(enum.Options, E_Enumdecl, proto.GetBoolExtension(file.Options, E_EnumdeclAll, true))
***REMOVED***

func HasTypeDecl(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_Typedecl, proto.GetBoolExtension(file.Options, E_TypedeclAll, true))
***REMOVED***

func GetCustomType(field *google_protobuf.FieldDescriptorProto) string ***REMOVED***
	if field == nil ***REMOVED***
		return ""
	***REMOVED***
	if field.Options != nil ***REMOVED***
		v, err := proto.GetExtension(field.Options, E_Customtype)
		if err == nil && v.(*string) != nil ***REMOVED***
			return *(v.(*string))
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

func GetCastType(field *google_protobuf.FieldDescriptorProto) string ***REMOVED***
	if field == nil ***REMOVED***
		return ""
	***REMOVED***
	if field.Options != nil ***REMOVED***
		v, err := proto.GetExtension(field.Options, E_Casttype)
		if err == nil && v.(*string) != nil ***REMOVED***
			return *(v.(*string))
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

func GetCastKey(field *google_protobuf.FieldDescriptorProto) string ***REMOVED***
	if field == nil ***REMOVED***
		return ""
	***REMOVED***
	if field.Options != nil ***REMOVED***
		v, err := proto.GetExtension(field.Options, E_Castkey)
		if err == nil && v.(*string) != nil ***REMOVED***
			return *(v.(*string))
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

func GetCastValue(field *google_protobuf.FieldDescriptorProto) string ***REMOVED***
	if field == nil ***REMOVED***
		return ""
	***REMOVED***
	if field.Options != nil ***REMOVED***
		v, err := proto.GetExtension(field.Options, E_Castvalue)
		if err == nil && v.(*string) != nil ***REMOVED***
			return *(v.(*string))
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

func IsCustomName(field *google_protobuf.FieldDescriptorProto) bool ***REMOVED***
	name := GetCustomName(field)
	if len(name) > 0 ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func IsEnumCustomName(field *google_protobuf.EnumDescriptorProto) bool ***REMOVED***
	name := GetEnumCustomName(field)
	if len(name) > 0 ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func IsEnumValueCustomName(field *google_protobuf.EnumValueDescriptorProto) bool ***REMOVED***
	name := GetEnumValueCustomName(field)
	if len(name) > 0 ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func GetCustomName(field *google_protobuf.FieldDescriptorProto) string ***REMOVED***
	if field == nil ***REMOVED***
		return ""
	***REMOVED***
	if field.Options != nil ***REMOVED***
		v, err := proto.GetExtension(field.Options, E_Customname)
		if err == nil && v.(*string) != nil ***REMOVED***
			return *(v.(*string))
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

func GetEnumCustomName(field *google_protobuf.EnumDescriptorProto) string ***REMOVED***
	if field == nil ***REMOVED***
		return ""
	***REMOVED***
	if field.Options != nil ***REMOVED***
		v, err := proto.GetExtension(field.Options, E_EnumCustomname)
		if err == nil && v.(*string) != nil ***REMOVED***
			return *(v.(*string))
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

func GetEnumValueCustomName(field *google_protobuf.EnumValueDescriptorProto) string ***REMOVED***
	if field == nil ***REMOVED***
		return ""
	***REMOVED***
	if field.Options != nil ***REMOVED***
		v, err := proto.GetExtension(field.Options, E_EnumvalueCustomname)
		if err == nil && v.(*string) != nil ***REMOVED***
			return *(v.(*string))
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

func GetJsonTag(field *google_protobuf.FieldDescriptorProto) *string ***REMOVED***
	if field == nil ***REMOVED***
		return nil
	***REMOVED***
	if field.Options != nil ***REMOVED***
		v, err := proto.GetExtension(field.Options, E_Jsontag)
		if err == nil && v.(*string) != nil ***REMOVED***
			return (v.(*string))
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func GetMoreTags(field *google_protobuf.FieldDescriptorProto) *string ***REMOVED***
	if field == nil ***REMOVED***
		return nil
	***REMOVED***
	if field.Options != nil ***REMOVED***
		v, err := proto.GetExtension(field.Options, E_Moretags)
		if err == nil && v.(*string) != nil ***REMOVED***
			return (v.(*string))
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type EnableFunc func(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool

func EnabledGoEnumPrefix(file *google_protobuf.FileDescriptorProto, enum *google_protobuf.EnumDescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(enum.Options, E_GoprotoEnumPrefix, proto.GetBoolExtension(file.Options, E_GoprotoEnumPrefixAll, true))
***REMOVED***

func EnabledGoStringer(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_GoprotoStringer, proto.GetBoolExtension(file.Options, E_GoprotoStringerAll, true))
***REMOVED***

func HasGoGetters(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_GoprotoGetters, proto.GetBoolExtension(file.Options, E_GoprotoGettersAll, true))
***REMOVED***

func IsUnion(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_Onlyone, proto.GetBoolExtension(file.Options, E_OnlyoneAll, false))
***REMOVED***

func HasGoString(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_Gostring, proto.GetBoolExtension(file.Options, E_GostringAll, false))
***REMOVED***

func HasEqual(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_Equal, proto.GetBoolExtension(file.Options, E_EqualAll, false))
***REMOVED***

func HasVerboseEqual(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_VerboseEqual, proto.GetBoolExtension(file.Options, E_VerboseEqualAll, false))
***REMOVED***

func IsStringer(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_Stringer, proto.GetBoolExtension(file.Options, E_StringerAll, false))
***REMOVED***

func IsFace(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_Face, proto.GetBoolExtension(file.Options, E_FaceAll, false))
***REMOVED***

func HasDescription(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_Description, proto.GetBoolExtension(file.Options, E_DescriptionAll, false))
***REMOVED***

func HasPopulate(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_Populate, proto.GetBoolExtension(file.Options, E_PopulateAll, false))
***REMOVED***

func HasTestGen(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_Testgen, proto.GetBoolExtension(file.Options, E_TestgenAll, false))
***REMOVED***

func HasBenchGen(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_Benchgen, proto.GetBoolExtension(file.Options, E_BenchgenAll, false))
***REMOVED***

func IsMarshaler(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_Marshaler, proto.GetBoolExtension(file.Options, E_MarshalerAll, false))
***REMOVED***

func IsUnmarshaler(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_Unmarshaler, proto.GetBoolExtension(file.Options, E_UnmarshalerAll, false))
***REMOVED***

func IsStableMarshaler(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_StableMarshaler, proto.GetBoolExtension(file.Options, E_StableMarshalerAll, false))
***REMOVED***

func IsSizer(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_Sizer, proto.GetBoolExtension(file.Options, E_SizerAll, false))
***REMOVED***

func IsProtoSizer(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_Protosizer, proto.GetBoolExtension(file.Options, E_ProtosizerAll, false))
***REMOVED***

func IsGoEnumStringer(file *google_protobuf.FileDescriptorProto, enum *google_protobuf.EnumDescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(enum.Options, E_GoprotoEnumStringer, proto.GetBoolExtension(file.Options, E_GoprotoEnumStringerAll, true))
***REMOVED***

func IsEnumStringer(file *google_protobuf.FileDescriptorProto, enum *google_protobuf.EnumDescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(enum.Options, E_EnumStringer, proto.GetBoolExtension(file.Options, E_EnumStringerAll, false))
***REMOVED***

func IsUnsafeMarshaler(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_UnsafeMarshaler, proto.GetBoolExtension(file.Options, E_UnsafeMarshalerAll, false))
***REMOVED***

func IsUnsafeUnmarshaler(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_UnsafeUnmarshaler, proto.GetBoolExtension(file.Options, E_UnsafeUnmarshalerAll, false))
***REMOVED***

func HasExtensionsMap(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_GoprotoExtensionsMap, proto.GetBoolExtension(file.Options, E_GoprotoExtensionsMapAll, true))
***REMOVED***

func HasUnrecognized(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	if IsProto3(file) ***REMOVED***
		return false
	***REMOVED***
	return proto.GetBoolExtension(message.Options, E_GoprotoUnrecognized, proto.GetBoolExtension(file.Options, E_GoprotoUnrecognizedAll, true))
***REMOVED***

func IsProto3(file *google_protobuf.FileDescriptorProto) bool ***REMOVED***
	return file.GetSyntax() == "proto3"
***REMOVED***

func ImportsGoGoProto(file *google_protobuf.FileDescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(file.Options, E_GogoprotoImport, true)
***REMOVED***

func HasCompare(file *google_protobuf.FileDescriptorProto, message *google_protobuf.DescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(message.Options, E_Compare, proto.GetBoolExtension(file.Options, E_CompareAll, false))
***REMOVED***

func RegistersGolangProto(file *google_protobuf.FileDescriptorProto) bool ***REMOVED***
	return proto.GetBoolExtension(file.Options, E_GoprotoRegistration, false)
***REMOVED***
