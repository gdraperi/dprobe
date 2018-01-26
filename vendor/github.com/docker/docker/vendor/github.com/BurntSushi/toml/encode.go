package toml

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type tomlEncodeError struct***REMOVED*** error ***REMOVED***

var (
	errArrayMixedElementTypes = errors.New(
		"can't encode array with mixed element types")
	errArrayNilElement = errors.New(
		"can't encode array with nil element")
	errNonString = errors.New(
		"can't encode a map with non-string key type")
	errAnonNonStruct = errors.New(
		"can't encode an anonymous field that is not a struct")
	errArrayNoTable = errors.New(
		"TOML array element can't contain a table")
	errNoKey = errors.New(
		"top-level values must be a Go map or struct")
	errAnything = errors.New("") // used in testing
)

var quotedReplacer = strings.NewReplacer(
	"\t", "\\t",
	"\n", "\\n",
	"\r", "\\r",
	"\"", "\\\"",
	"\\", "\\\\",
)

// Encoder controls the encoding of Go values to a TOML document to some
// io.Writer.
//
// The indentation level can be controlled with the Indent field.
type Encoder struct ***REMOVED***
	// A single indentation level. By default it is two spaces.
	Indent string

	// hasWritten is whether we have written any output to w yet.
	hasWritten bool
	w          *bufio.Writer
***REMOVED***

// NewEncoder returns a TOML encoder that encodes Go values to the io.Writer
// given. By default, a single indentation level is 2 spaces.
func NewEncoder(w io.Writer) *Encoder ***REMOVED***
	return &Encoder***REMOVED***
		w:      bufio.NewWriter(w),
		Indent: "  ",
	***REMOVED***
***REMOVED***

// Encode writes a TOML representation of the Go value to the underlying
// io.Writer. If the value given cannot be encoded to a valid TOML document,
// then an error is returned.
//
// The mapping between Go values and TOML values should be precisely the same
// as for the Decode* functions. Similarly, the TextMarshaler interface is
// supported by encoding the resulting bytes as strings. (If you want to write
// arbitrary binary data then you will need to use something like base64 since
// TOML does not have any binary types.)
//
// When encoding TOML hashes (i.e., Go maps or structs), keys without any
// sub-hashes are encoded first.
//
// If a Go map is encoded, then its keys are sorted alphabetically for
// deterministic output. More control over this behavior may be provided if
// there is demand for it.
//
// Encoding Go values without a corresponding TOML representation---like map
// types with non-string keys---will cause an error to be returned. Similarly
// for mixed arrays/slices, arrays/slices with nil elements, embedded
// non-struct types and nested slices containing maps or structs.
// (e.g., [][]map[string]string is not allowed but []map[string]string is OK
// and so is []map[string][]string.)
func (enc *Encoder) Encode(v interface***REMOVED******REMOVED***) error ***REMOVED***
	rv := eindirect(reflect.ValueOf(v))
	if err := enc.safeEncode(Key([]string***REMOVED******REMOVED***), rv); err != nil ***REMOVED***
		return err
	***REMOVED***
	return enc.w.Flush()
***REMOVED***

func (enc *Encoder) safeEncode(key Key, rv reflect.Value) (err error) ***REMOVED***
	defer func() ***REMOVED***
		if r := recover(); r != nil ***REMOVED***
			if terr, ok := r.(tomlEncodeError); ok ***REMOVED***
				err = terr.error
				return
			***REMOVED***
			panic(r)
		***REMOVED***
	***REMOVED***()
	enc.encode(key, rv)
	return nil
***REMOVED***

func (enc *Encoder) encode(key Key, rv reflect.Value) ***REMOVED***
	// Special case. Time needs to be in ISO8601 format.
	// Special case. If we can marshal the type to text, then we used that.
	// Basically, this prevents the encoder for handling these types as
	// generic structs (or whatever the underlying type of a TextMarshaler is).
	switch rv.Interface().(type) ***REMOVED***
	case time.Time, TextMarshaler:
		enc.keyEqElement(key, rv)
		return
	***REMOVED***

	k := rv.Kind()
	switch k ***REMOVED***
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
		enc.keyEqElement(key, rv)
	case reflect.Array, reflect.Slice:
		if typeEqual(tomlArrayHash, tomlTypeOfGo(rv)) ***REMOVED***
			enc.eArrayOfTables(key, rv)
		***REMOVED*** else ***REMOVED***
			enc.keyEqElement(key, rv)
		***REMOVED***
	case reflect.Interface:
		if rv.IsNil() ***REMOVED***
			return
		***REMOVED***
		enc.encode(key, rv.Elem())
	case reflect.Map:
		if rv.IsNil() ***REMOVED***
			return
		***REMOVED***
		enc.eTable(key, rv)
	case reflect.Ptr:
		if rv.IsNil() ***REMOVED***
			return
		***REMOVED***
		enc.encode(key, rv.Elem())
	case reflect.Struct:
		enc.eTable(key, rv)
	default:
		panic(e("Unsupported type for key '%s': %s", key, k))
	***REMOVED***
***REMOVED***

// eElement encodes any value that can be an array element (primitives and
// arrays).
func (enc *Encoder) eElement(rv reflect.Value) ***REMOVED***
	switch v := rv.Interface().(type) ***REMOVED***
	case time.Time:
		// Special case time.Time as a primitive. Has to come before
		// TextMarshaler below because time.Time implements
		// encoding.TextMarshaler, but we need to always use UTC.
		enc.wf(v.In(time.FixedZone("UTC", 0)).Format("2006-01-02T15:04:05Z"))
		return
	case TextMarshaler:
		// Special case. Use text marshaler if it's available for this value.
		if s, err := v.MarshalText(); err != nil ***REMOVED***
			encPanic(err)
		***REMOVED*** else ***REMOVED***
			enc.writeQuoted(string(s))
		***REMOVED***
		return
	***REMOVED***
	switch rv.Kind() ***REMOVED***
	case reflect.Bool:
		enc.wf(strconv.FormatBool(rv.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64:
		enc.wf(strconv.FormatInt(rv.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		enc.wf(strconv.FormatUint(rv.Uint(), 10))
	case reflect.Float32:
		enc.wf(floatAddDecimal(strconv.FormatFloat(rv.Float(), 'f', -1, 32)))
	case reflect.Float64:
		enc.wf(floatAddDecimal(strconv.FormatFloat(rv.Float(), 'f', -1, 64)))
	case reflect.Array, reflect.Slice:
		enc.eArrayOrSliceElement(rv)
	case reflect.Interface:
		enc.eElement(rv.Elem())
	case reflect.String:
		enc.writeQuoted(rv.String())
	default:
		panic(e("Unexpected primitive type: %s", rv.Kind()))
	***REMOVED***
***REMOVED***

// By the TOML spec, all floats must have a decimal with at least one
// number on either side.
func floatAddDecimal(fstr string) string ***REMOVED***
	if !strings.Contains(fstr, ".") ***REMOVED***
		return fstr + ".0"
	***REMOVED***
	return fstr
***REMOVED***

func (enc *Encoder) writeQuoted(s string) ***REMOVED***
	enc.wf("\"%s\"", quotedReplacer.Replace(s))
***REMOVED***

func (enc *Encoder) eArrayOrSliceElement(rv reflect.Value) ***REMOVED***
	length := rv.Len()
	enc.wf("[")
	for i := 0; i < length; i++ ***REMOVED***
		elem := rv.Index(i)
		enc.eElement(elem)
		if i != length-1 ***REMOVED***
			enc.wf(", ")
		***REMOVED***
	***REMOVED***
	enc.wf("]")
***REMOVED***

func (enc *Encoder) eArrayOfTables(key Key, rv reflect.Value) ***REMOVED***
	if len(key) == 0 ***REMOVED***
		encPanic(errNoKey)
	***REMOVED***
	for i := 0; i < rv.Len(); i++ ***REMOVED***
		trv := rv.Index(i)
		if isNil(trv) ***REMOVED***
			continue
		***REMOVED***
		panicIfInvalidKey(key)
		enc.newline()
		enc.wf("%s[[%s]]", enc.indentStr(key), key.maybeQuotedAll())
		enc.newline()
		enc.eMapOrStruct(key, trv)
	***REMOVED***
***REMOVED***

func (enc *Encoder) eTable(key Key, rv reflect.Value) ***REMOVED***
	panicIfInvalidKey(key)
	if len(key) == 1 ***REMOVED***
		// Output an extra new line between top-level tables.
		// (The newline isn't written if nothing else has been written though.)
		enc.newline()
	***REMOVED***
	if len(key) > 0 ***REMOVED***
		enc.wf("%s[%s]", enc.indentStr(key), key.maybeQuotedAll())
		enc.newline()
	***REMOVED***
	enc.eMapOrStruct(key, rv)
***REMOVED***

func (enc *Encoder) eMapOrStruct(key Key, rv reflect.Value) ***REMOVED***
	switch rv := eindirect(rv); rv.Kind() ***REMOVED***
	case reflect.Map:
		enc.eMap(key, rv)
	case reflect.Struct:
		enc.eStruct(key, rv)
	default:
		panic("eTable: unhandled reflect.Value Kind: " + rv.Kind().String())
	***REMOVED***
***REMOVED***

func (enc *Encoder) eMap(key Key, rv reflect.Value) ***REMOVED***
	rt := rv.Type()
	if rt.Key().Kind() != reflect.String ***REMOVED***
		encPanic(errNonString)
	***REMOVED***

	// Sort keys so that we have deterministic output. And write keys directly
	// underneath this key first, before writing sub-structs or sub-maps.
	var mapKeysDirect, mapKeysSub []string
	for _, mapKey := range rv.MapKeys() ***REMOVED***
		k := mapKey.String()
		if typeIsHash(tomlTypeOfGo(rv.MapIndex(mapKey))) ***REMOVED***
			mapKeysSub = append(mapKeysSub, k)
		***REMOVED*** else ***REMOVED***
			mapKeysDirect = append(mapKeysDirect, k)
		***REMOVED***
	***REMOVED***

	var writeMapKeys = func(mapKeys []string) ***REMOVED***
		sort.Strings(mapKeys)
		for _, mapKey := range mapKeys ***REMOVED***
			mrv := rv.MapIndex(reflect.ValueOf(mapKey))
			if isNil(mrv) ***REMOVED***
				// Don't write anything for nil fields.
				continue
			***REMOVED***
			enc.encode(key.add(mapKey), mrv)
		***REMOVED***
	***REMOVED***
	writeMapKeys(mapKeysDirect)
	writeMapKeys(mapKeysSub)
***REMOVED***

func (enc *Encoder) eStruct(key Key, rv reflect.Value) ***REMOVED***
	// Write keys for fields directly under this key first, because if we write
	// a field that creates a new table, then all keys under it will be in that
	// table (not the one we're writing here).
	rt := rv.Type()
	var fieldsDirect, fieldsSub [][]int
	var addFields func(rt reflect.Type, rv reflect.Value, start []int)
	addFields = func(rt reflect.Type, rv reflect.Value, start []int) ***REMOVED***
		for i := 0; i < rt.NumField(); i++ ***REMOVED***
			f := rt.Field(i)
			// skip unexporded fields
			if f.PkgPath != "" ***REMOVED***
				continue
			***REMOVED***
			frv := rv.Field(i)
			if f.Anonymous ***REMOVED***
				frv := eindirect(frv)
				t := frv.Type()
				if t.Kind() != reflect.Struct ***REMOVED***
					encPanic(errAnonNonStruct)
				***REMOVED***
				addFields(t, frv, f.Index)
			***REMOVED*** else if typeIsHash(tomlTypeOfGo(frv)) ***REMOVED***
				fieldsSub = append(fieldsSub, append(start, f.Index...))
			***REMOVED*** else ***REMOVED***
				fieldsDirect = append(fieldsDirect, append(start, f.Index...))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	addFields(rt, rv, nil)

	var writeFields = func(fields [][]int) ***REMOVED***
		for _, fieldIndex := range fields ***REMOVED***
			sft := rt.FieldByIndex(fieldIndex)
			sf := rv.FieldByIndex(fieldIndex)
			if isNil(sf) ***REMOVED***
				// Don't write anything for nil fields.
				continue
			***REMOVED***

			keyName := sft.Tag.Get("toml")
			if keyName == "-" ***REMOVED***
				continue
			***REMOVED***
			if keyName == "" ***REMOVED***
				keyName = sft.Name
			***REMOVED***
			enc.encode(key.add(keyName), sf)
		***REMOVED***
	***REMOVED***
	writeFields(fieldsDirect)
	writeFields(fieldsSub)
***REMOVED***

// tomlTypeName returns the TOML type name of the Go value's type. It is
// used to determine whether the types of array elements are mixed (which is
// forbidden). If the Go value is nil, then it is illegal for it to be an array
// element, and valueIsNil is returned as true.

// Returns the TOML type of a Go value. The type may be `nil`, which means
// no concrete TOML type could be found.
func tomlTypeOfGo(rv reflect.Value) tomlType ***REMOVED***
	if isNil(rv) || !rv.IsValid() ***REMOVED***
		return nil
	***REMOVED***
	switch rv.Kind() ***REMOVED***
	case reflect.Bool:
		return tomlBool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64:
		return tomlInteger
	case reflect.Float32, reflect.Float64:
		return tomlFloat
	case reflect.Array, reflect.Slice:
		if typeEqual(tomlHash, tomlArrayType(rv)) ***REMOVED***
			return tomlArrayHash
		***REMOVED*** else ***REMOVED***
			return tomlArray
		***REMOVED***
	case reflect.Ptr, reflect.Interface:
		return tomlTypeOfGo(rv.Elem())
	case reflect.String:
		return tomlString
	case reflect.Map:
		return tomlHash
	case reflect.Struct:
		switch rv.Interface().(type) ***REMOVED***
		case time.Time:
			return tomlDatetime
		case TextMarshaler:
			return tomlString
		default:
			return tomlHash
		***REMOVED***
	default:
		panic("unexpected reflect.Kind: " + rv.Kind().String())
	***REMOVED***
***REMOVED***

// tomlArrayType returns the element type of a TOML array. The type returned
// may be nil if it cannot be determined (e.g., a nil slice or a zero length
// slize). This function may also panic if it finds a type that cannot be
// expressed in TOML (such as nil elements, heterogeneous arrays or directly
// nested arrays of tables).
func tomlArrayType(rv reflect.Value) tomlType ***REMOVED***
	if isNil(rv) || !rv.IsValid() || rv.Len() == 0 ***REMOVED***
		return nil
	***REMOVED***
	firstType := tomlTypeOfGo(rv.Index(0))
	if firstType == nil ***REMOVED***
		encPanic(errArrayNilElement)
	***REMOVED***

	rvlen := rv.Len()
	for i := 1; i < rvlen; i++ ***REMOVED***
		elem := rv.Index(i)
		switch elemType := tomlTypeOfGo(elem); ***REMOVED***
		case elemType == nil:
			encPanic(errArrayNilElement)
		case !typeEqual(firstType, elemType):
			encPanic(errArrayMixedElementTypes)
		***REMOVED***
	***REMOVED***
	// If we have a nested array, then we must make sure that the nested
	// array contains ONLY primitives.
	// This checks arbitrarily nested arrays.
	if typeEqual(firstType, tomlArray) || typeEqual(firstType, tomlArrayHash) ***REMOVED***
		nest := tomlArrayType(eindirect(rv.Index(0)))
		if typeEqual(nest, tomlHash) || typeEqual(nest, tomlArrayHash) ***REMOVED***
			encPanic(errArrayNoTable)
		***REMOVED***
	***REMOVED***
	return firstType
***REMOVED***

func (enc *Encoder) newline() ***REMOVED***
	if enc.hasWritten ***REMOVED***
		enc.wf("\n")
	***REMOVED***
***REMOVED***

func (enc *Encoder) keyEqElement(key Key, val reflect.Value) ***REMOVED***
	if len(key) == 0 ***REMOVED***
		encPanic(errNoKey)
	***REMOVED***
	panicIfInvalidKey(key)
	enc.wf("%s%s = ", enc.indentStr(key), key.maybeQuoted(len(key)-1))
	enc.eElement(val)
	enc.newline()
***REMOVED***

func (enc *Encoder) wf(format string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
	if _, err := fmt.Fprintf(enc.w, format, v...); err != nil ***REMOVED***
		encPanic(err)
	***REMOVED***
	enc.hasWritten = true
***REMOVED***

func (enc *Encoder) indentStr(key Key) string ***REMOVED***
	return strings.Repeat(enc.Indent, len(key)-1)
***REMOVED***

func encPanic(err error) ***REMOVED***
	panic(tomlEncodeError***REMOVED***err***REMOVED***)
***REMOVED***

func eindirect(v reflect.Value) reflect.Value ***REMOVED***
	switch v.Kind() ***REMOVED***
	case reflect.Ptr, reflect.Interface:
		return eindirect(v.Elem())
	default:
		return v
	***REMOVED***
***REMOVED***

func isNil(rv reflect.Value) bool ***REMOVED***
	switch rv.Kind() ***REMOVED***
	case reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	***REMOVED***
***REMOVED***

func panicIfInvalidKey(key Key) ***REMOVED***
	for _, k := range key ***REMOVED***
		if len(k) == 0 ***REMOVED***
			encPanic(e("Key '%s' is not a valid table name. Key names "+
				"cannot be empty.", key.maybeQuotedAll()))
		***REMOVED***
	***REMOVED***
***REMOVED***

func isValidKeyName(s string) bool ***REMOVED***
	return len(s) != 0
***REMOVED***
