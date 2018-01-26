// Copyright 2014 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package ini

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"
)

// NameMapper represents a ini tag name mapper.
type NameMapper func(string) string

// Built-in name getters.
var (
	// AllCapsUnderscore converts to format ALL_CAPS_UNDERSCORE.
	AllCapsUnderscore NameMapper = func(raw string) string ***REMOVED***
		newstr := make([]rune, 0, len(raw))
		for i, chr := range raw ***REMOVED***
			if isUpper := 'A' <= chr && chr <= 'Z'; isUpper ***REMOVED***
				if i > 0 ***REMOVED***
					newstr = append(newstr, '_')
				***REMOVED***
			***REMOVED***
			newstr = append(newstr, unicode.ToUpper(chr))
		***REMOVED***
		return string(newstr)
	***REMOVED***
	// TitleUnderscore converts to format title_underscore.
	TitleUnderscore NameMapper = func(raw string) string ***REMOVED***
		newstr := make([]rune, 0, len(raw))
		for i, chr := range raw ***REMOVED***
			if isUpper := 'A' <= chr && chr <= 'Z'; isUpper ***REMOVED***
				if i > 0 ***REMOVED***
					newstr = append(newstr, '_')
				***REMOVED***
				chr -= ('A' - 'a')
			***REMOVED***
			newstr = append(newstr, chr)
		***REMOVED***
		return string(newstr)
	***REMOVED***
)

func (s *Section) parseFieldName(raw, actual string) string ***REMOVED***
	if len(actual) > 0 ***REMOVED***
		return actual
	***REMOVED***
	if s.f.NameMapper != nil ***REMOVED***
		return s.f.NameMapper(raw)
	***REMOVED***
	return raw
***REMOVED***

func parseDelim(actual string) string ***REMOVED***
	if len(actual) > 0 ***REMOVED***
		return actual
	***REMOVED***
	return ","
***REMOVED***

var reflectTime = reflect.TypeOf(time.Now()).Kind()

// setSliceWithProperType sets proper values to slice based on its type.
func setSliceWithProperType(key *Key, field reflect.Value, delim string, allowShadow bool) error ***REMOVED***
	var strs []string
	if allowShadow ***REMOVED***
		strs = key.StringsWithShadows(delim)
	***REMOVED*** else ***REMOVED***
		strs = key.Strings(delim)
	***REMOVED***

	numVals := len(strs)
	if numVals == 0 ***REMOVED***
		return nil
	***REMOVED***

	var vals interface***REMOVED******REMOVED***

	sliceOf := field.Type().Elem().Kind()
	switch sliceOf ***REMOVED***
	case reflect.String:
		vals = strs
	case reflect.Int:
		vals, _ = key.parseInts(strs, true, false)
	case reflect.Int64:
		vals, _ = key.parseInt64s(strs, true, false)
	case reflect.Uint:
		vals = key.Uints(delim)
	case reflect.Uint64:
		vals = key.Uint64s(delim)
	case reflect.Float64:
		vals = key.Float64s(delim)
	case reflectTime:
		vals = key.Times(delim)
	default:
		return fmt.Errorf("unsupported type '[]%s'", sliceOf)
	***REMOVED***

	slice := reflect.MakeSlice(field.Type(), numVals, numVals)
	for i := 0; i < numVals; i++ ***REMOVED***
		switch sliceOf ***REMOVED***
		case reflect.String:
			slice.Index(i).Set(reflect.ValueOf(vals.([]string)[i]))
		case reflect.Int:
			slice.Index(i).Set(reflect.ValueOf(vals.([]int)[i]))
		case reflect.Int64:
			slice.Index(i).Set(reflect.ValueOf(vals.([]int64)[i]))
		case reflect.Uint:
			slice.Index(i).Set(reflect.ValueOf(vals.([]uint)[i]))
		case reflect.Uint64:
			slice.Index(i).Set(reflect.ValueOf(vals.([]uint64)[i]))
		case reflect.Float64:
			slice.Index(i).Set(reflect.ValueOf(vals.([]float64)[i]))
		case reflectTime:
			slice.Index(i).Set(reflect.ValueOf(vals.([]time.Time)[i]))
		***REMOVED***
	***REMOVED***
	field.Set(slice)
	return nil
***REMOVED***

// setWithProperType sets proper value to field based on its type,
// but it does not return error for failing parsing,
// because we want to use default value that is already assigned to strcut.
func setWithProperType(t reflect.Type, key *Key, field reflect.Value, delim string, allowShadow bool) error ***REMOVED***
	switch t.Kind() ***REMOVED***
	case reflect.String:
		if len(key.String()) == 0 ***REMOVED***
			return nil
		***REMOVED***
		field.SetString(key.String())
	case reflect.Bool:
		boolVal, err := key.Bool()
		if err != nil ***REMOVED***
			return nil
		***REMOVED***
		field.SetBool(boolVal)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		durationVal, err := key.Duration()
		// Skip zero value
		if err == nil && int(durationVal) > 0 ***REMOVED***
			field.Set(reflect.ValueOf(durationVal))
			return nil
		***REMOVED***

		intVal, err := key.Int64()
		if err != nil || intVal == 0 ***REMOVED***
			return nil
		***REMOVED***
		field.SetInt(intVal)
	//	byte is an alias for uint8, so supporting uint8 breaks support for byte
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		durationVal, err := key.Duration()
		// Skip zero value
		if err == nil && int(durationVal) > 0 ***REMOVED***
			field.Set(reflect.ValueOf(durationVal))
			return nil
		***REMOVED***

		uintVal, err := key.Uint64()
		if err != nil ***REMOVED***
			return nil
		***REMOVED***
		field.SetUint(uintVal)

	case reflect.Float32, reflect.Float64:
		floatVal, err := key.Float64()
		if err != nil ***REMOVED***
			return nil
		***REMOVED***
		field.SetFloat(floatVal)
	case reflectTime:
		timeVal, err := key.Time()
		if err != nil ***REMOVED***
			return nil
		***REMOVED***
		field.Set(reflect.ValueOf(timeVal))
	case reflect.Slice:
		return setSliceWithProperType(key, field, delim, allowShadow)
	default:
		return fmt.Errorf("unsupported type '%s'", t)
	***REMOVED***
	return nil
***REMOVED***

func parseTagOptions(tag string) (rawName string, omitEmpty bool, allowShadow bool) ***REMOVED***
	opts := strings.SplitN(tag, ",", 3)
	rawName = opts[0]
	if len(opts) > 1 ***REMOVED***
		omitEmpty = opts[1] == "omitempty"
	***REMOVED***
	if len(opts) > 2 ***REMOVED***
		allowShadow = opts[2] == "allowshadow"
	***REMOVED***
	return rawName, omitEmpty, allowShadow
***REMOVED***

func (s *Section) mapTo(val reflect.Value) error ***REMOVED***
	if val.Kind() == reflect.Ptr ***REMOVED***
		val = val.Elem()
	***REMOVED***
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ ***REMOVED***
		field := val.Field(i)
		tpField := typ.Field(i)

		tag := tpField.Tag.Get("ini")
		if tag == "-" ***REMOVED***
			continue
		***REMOVED***

		rawName, _, allowShadow := parseTagOptions(tag)
		fieldName := s.parseFieldName(tpField.Name, rawName)
		if len(fieldName) == 0 || !field.CanSet() ***REMOVED***
			continue
		***REMOVED***

		isAnonymous := tpField.Type.Kind() == reflect.Ptr && tpField.Anonymous
		isStruct := tpField.Type.Kind() == reflect.Struct
		if isAnonymous ***REMOVED***
			field.Set(reflect.New(tpField.Type.Elem()))
		***REMOVED***

		if isAnonymous || isStruct ***REMOVED***
			if sec, err := s.f.GetSection(fieldName); err == nil ***REMOVED***
				if err = sec.mapTo(field); err != nil ***REMOVED***
					return fmt.Errorf("error mapping field(%s): %v", fieldName, err)
				***REMOVED***
				continue
			***REMOVED***
		***REMOVED***

		if key, err := s.GetKey(fieldName); err == nil ***REMOVED***
			delim := parseDelim(tpField.Tag.Get("delim"))
			if err = setWithProperType(tpField.Type, key, field, delim, allowShadow); err != nil ***REMOVED***
				return fmt.Errorf("error mapping field(%s): %v", fieldName, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// MapTo maps section to given struct.
func (s *Section) MapTo(v interface***REMOVED******REMOVED***) error ***REMOVED***
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	if typ.Kind() == reflect.Ptr ***REMOVED***
		typ = typ.Elem()
		val = val.Elem()
	***REMOVED*** else ***REMOVED***
		return errors.New("cannot map to non-pointer struct")
	***REMOVED***

	return s.mapTo(val)
***REMOVED***

// MapTo maps file to given struct.
func (f *File) MapTo(v interface***REMOVED******REMOVED***) error ***REMOVED***
	return f.Section("").MapTo(v)
***REMOVED***

// MapTo maps data sources to given struct with name mapper.
func MapToWithMapper(v interface***REMOVED******REMOVED***, mapper NameMapper, source interface***REMOVED******REMOVED***, others ...interface***REMOVED******REMOVED***) error ***REMOVED***
	cfg, err := Load(source, others...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	cfg.NameMapper = mapper
	return cfg.MapTo(v)
***REMOVED***

// MapTo maps data sources to given struct.
func MapTo(v, source interface***REMOVED******REMOVED***, others ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return MapToWithMapper(v, nil, source, others...)
***REMOVED***

// reflectSliceWithProperType does the opposite thing as setSliceWithProperType.
func reflectSliceWithProperType(key *Key, field reflect.Value, delim string) error ***REMOVED***
	slice := field.Slice(0, field.Len())
	if field.Len() == 0 ***REMOVED***
		return nil
	***REMOVED***

	var buf bytes.Buffer
	sliceOf := field.Type().Elem().Kind()
	for i := 0; i < field.Len(); i++ ***REMOVED***
		switch sliceOf ***REMOVED***
		case reflect.String:
			buf.WriteString(slice.Index(i).String())
		case reflect.Int, reflect.Int64:
			buf.WriteString(fmt.Sprint(slice.Index(i).Int()))
		case reflect.Uint, reflect.Uint64:
			buf.WriteString(fmt.Sprint(slice.Index(i).Uint()))
		case reflect.Float64:
			buf.WriteString(fmt.Sprint(slice.Index(i).Float()))
		case reflectTime:
			buf.WriteString(slice.Index(i).Interface().(time.Time).Format(time.RFC3339))
		default:
			return fmt.Errorf("unsupported type '[]%s'", sliceOf)
		***REMOVED***
		buf.WriteString(delim)
	***REMOVED***
	key.SetValue(buf.String()[:buf.Len()-1])
	return nil
***REMOVED***

// reflectWithProperType does the opposite thing as setWithProperType.
func reflectWithProperType(t reflect.Type, key *Key, field reflect.Value, delim string) error ***REMOVED***
	switch t.Kind() ***REMOVED***
	case reflect.String:
		key.SetValue(field.String())
	case reflect.Bool:
		key.SetValue(fmt.Sprint(field.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		key.SetValue(fmt.Sprint(field.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		key.SetValue(fmt.Sprint(field.Uint()))
	case reflect.Float32, reflect.Float64:
		key.SetValue(fmt.Sprint(field.Float()))
	case reflectTime:
		key.SetValue(fmt.Sprint(field.Interface().(time.Time).Format(time.RFC3339)))
	case reflect.Slice:
		return reflectSliceWithProperType(key, field, delim)
	default:
		return fmt.Errorf("unsupported type '%s'", t)
	***REMOVED***
	return nil
***REMOVED***

// CR: copied from encoding/json/encode.go with modifications of time.Time support.
// TODO: add more test coverage.
func isEmptyValue(v reflect.Value) bool ***REMOVED***
	switch v.Kind() ***REMOVED***
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflectTime:
		return v.Interface().(time.Time).IsZero()
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	***REMOVED***
	return false
***REMOVED***

func (s *Section) reflectFrom(val reflect.Value) error ***REMOVED***
	if val.Kind() == reflect.Ptr ***REMOVED***
		val = val.Elem()
	***REMOVED***
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ ***REMOVED***
		field := val.Field(i)
		tpField := typ.Field(i)

		tag := tpField.Tag.Get("ini")
		if tag == "-" ***REMOVED***
			continue
		***REMOVED***

		opts := strings.SplitN(tag, ",", 2)
		if len(opts) == 2 && opts[1] == "omitempty" && isEmptyValue(field) ***REMOVED***
			continue
		***REMOVED***

		fieldName := s.parseFieldName(tpField.Name, opts[0])
		if len(fieldName) == 0 || !field.CanSet() ***REMOVED***
			continue
		***REMOVED***

		if (tpField.Type.Kind() == reflect.Ptr && tpField.Anonymous) ||
			(tpField.Type.Kind() == reflect.Struct && tpField.Type.Name() != "Time") ***REMOVED***
			// Note: The only error here is section doesn't exist.
			sec, err := s.f.GetSection(fieldName)
			if err != nil ***REMOVED***
				// Note: fieldName can never be empty here, ignore error.
				sec, _ = s.f.NewSection(fieldName)
			***REMOVED***
			if err = sec.reflectFrom(field); err != nil ***REMOVED***
				return fmt.Errorf("error reflecting field (%s): %v", fieldName, err)
			***REMOVED***
			continue
		***REMOVED***

		// Note: Same reason as secion.
		key, err := s.GetKey(fieldName)
		if err != nil ***REMOVED***
			key, _ = s.NewKey(fieldName, "")
		***REMOVED***
		if err = reflectWithProperType(tpField.Type, key, field, parseDelim(tpField.Tag.Get("delim"))); err != nil ***REMOVED***
			return fmt.Errorf("error reflecting field (%s): %v", fieldName, err)
		***REMOVED***

	***REMOVED***
	return nil
***REMOVED***

// ReflectFrom reflects secion from given struct.
func (s *Section) ReflectFrom(v interface***REMOVED******REMOVED***) error ***REMOVED***
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	if typ.Kind() == reflect.Ptr ***REMOVED***
		typ = typ.Elem()
		val = val.Elem()
	***REMOVED*** else ***REMOVED***
		return errors.New("cannot reflect from non-pointer struct")
	***REMOVED***

	return s.reflectFrom(val)
***REMOVED***

// ReflectFrom reflects file from given struct.
func (f *File) ReflectFrom(v interface***REMOVED******REMOVED***) error ***REMOVED***
	return f.Section("").ReflectFrom(v)
***REMOVED***

// ReflectFrom reflects data sources from given struct with name mapper.
func ReflectFromWithMapper(cfg *File, v interface***REMOVED******REMOVED***, mapper NameMapper) error ***REMOVED***
	cfg.NameMapper = mapper
	return cfg.ReflectFrom(v)
***REMOVED***

// ReflectFrom reflects data sources from given struct.
func ReflectFrom(cfg *File, v interface***REMOVED******REMOVED***) error ***REMOVED***
	return ReflectFromWithMapper(cfg, v, nil)
***REMOVED***
