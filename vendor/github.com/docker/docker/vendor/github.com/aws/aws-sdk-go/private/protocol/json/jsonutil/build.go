// Package jsonutil provides JSON serialization of AWS requests and responses.
package jsonutil

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/private/protocol"
)

var timeType = reflect.ValueOf(time.Time***REMOVED******REMOVED***).Type()
var byteSliceType = reflect.ValueOf([]byte***REMOVED******REMOVED***).Type()

// BuildJSON builds a JSON string for a given object v.
func BuildJSON(v interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	var buf bytes.Buffer

	err := buildAny(reflect.ValueOf(v), &buf, "")
	return buf.Bytes(), err
***REMOVED***

func buildAny(value reflect.Value, buf *bytes.Buffer, tag reflect.StructTag) error ***REMOVED***
	origVal := value
	value = reflect.Indirect(value)
	if !value.IsValid() ***REMOVED***
		return nil
	***REMOVED***

	vtype := value.Type()

	t := tag.Get("type")
	if t == "" ***REMOVED***
		switch vtype.Kind() ***REMOVED***
		case reflect.Struct:
			// also it can't be a time object
			if value.Type() != timeType ***REMOVED***
				t = "structure"
			***REMOVED***
		case reflect.Slice:
			// also it can't be a byte slice
			if _, ok := value.Interface().([]byte); !ok ***REMOVED***
				t = "list"
			***REMOVED***
		case reflect.Map:
			// cannot be a JSONValue map
			if _, ok := value.Interface().(aws.JSONValue); !ok ***REMOVED***
				t = "map"
			***REMOVED***
		***REMOVED***
	***REMOVED***

	switch t ***REMOVED***
	case "structure":
		if field, ok := vtype.FieldByName("_"); ok ***REMOVED***
			tag = field.Tag
		***REMOVED***
		return buildStruct(value, buf, tag)
	case "list":
		return buildList(value, buf, tag)
	case "map":
		return buildMap(value, buf, tag)
	default:
		return buildScalar(origVal, buf, tag)
	***REMOVED***
***REMOVED***

func buildStruct(value reflect.Value, buf *bytes.Buffer, tag reflect.StructTag) error ***REMOVED***
	if !value.IsValid() ***REMOVED***
		return nil
	***REMOVED***

	// unwrap payloads
	if payload := tag.Get("payload"); payload != "" ***REMOVED***
		field, _ := value.Type().FieldByName(payload)
		tag = field.Tag
		value = elemOf(value.FieldByName(payload))

		if !value.IsValid() ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	buf.WriteByte('***REMOVED***')

	t := value.Type()
	first := true
	for i := 0; i < t.NumField(); i++ ***REMOVED***
		member := value.Field(i)

		// This allocates the most memory.
		// Additionally, we cannot skip nil fields due to
		// idempotency auto filling.
		field := t.Field(i)

		if field.PkgPath != "" ***REMOVED***
			continue // ignore unexported fields
		***REMOVED***
		if field.Tag.Get("json") == "-" ***REMOVED***
			continue
		***REMOVED***
		if field.Tag.Get("location") != "" ***REMOVED***
			continue // ignore non-body elements
		***REMOVED***
		if field.Tag.Get("ignore") != "" ***REMOVED***
			continue
		***REMOVED***

		if protocol.CanSetIdempotencyToken(member, field) ***REMOVED***
			token := protocol.GetIdempotencyToken()
			member = reflect.ValueOf(&token)
		***REMOVED***

		if (member.Kind() == reflect.Ptr || member.Kind() == reflect.Slice || member.Kind() == reflect.Map) && member.IsNil() ***REMOVED***
			continue // ignore unset fields
		***REMOVED***

		if first ***REMOVED***
			first = false
		***REMOVED*** else ***REMOVED***
			buf.WriteByte(',')
		***REMOVED***

		// figure out what this field is called
		name := field.Name
		if locName := field.Tag.Get("locationName"); locName != "" ***REMOVED***
			name = locName
		***REMOVED***

		writeString(name, buf)
		buf.WriteString(`:`)

		err := buildAny(member, buf, field.Tag)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

	***REMOVED***

	buf.WriteString("***REMOVED***")

	return nil
***REMOVED***

func buildList(value reflect.Value, buf *bytes.Buffer, tag reflect.StructTag) error ***REMOVED***
	buf.WriteString("[")

	for i := 0; i < value.Len(); i++ ***REMOVED***
		buildAny(value.Index(i), buf, "")

		if i < value.Len()-1 ***REMOVED***
			buf.WriteString(",")
		***REMOVED***
	***REMOVED***

	buf.WriteString("]")

	return nil
***REMOVED***

type sortedValues []reflect.Value

func (sv sortedValues) Len() int           ***REMOVED*** return len(sv) ***REMOVED***
func (sv sortedValues) Swap(i, j int)      ***REMOVED*** sv[i], sv[j] = sv[j], sv[i] ***REMOVED***
func (sv sortedValues) Less(i, j int) bool ***REMOVED*** return sv[i].String() < sv[j].String() ***REMOVED***

func buildMap(value reflect.Value, buf *bytes.Buffer, tag reflect.StructTag) error ***REMOVED***
	buf.WriteString("***REMOVED***")

	sv := sortedValues(value.MapKeys())
	sort.Sort(sv)

	for i, k := range sv ***REMOVED***
		if i > 0 ***REMOVED***
			buf.WriteByte(',')
		***REMOVED***

		writeString(k.String(), buf)
		buf.WriteString(`:`)

		buildAny(value.MapIndex(k), buf, "")
	***REMOVED***

	buf.WriteString("***REMOVED***")

	return nil
***REMOVED***

func buildScalar(v reflect.Value, buf *bytes.Buffer, tag reflect.StructTag) error ***REMOVED***
	// prevents allocation on the heap.
	scratch := [64]byte***REMOVED******REMOVED***
	switch value := reflect.Indirect(v); value.Kind() ***REMOVED***
	case reflect.String:
		writeString(value.String(), buf)
	case reflect.Bool:
		if value.Bool() ***REMOVED***
			buf.WriteString("true")
		***REMOVED*** else ***REMOVED***
			buf.WriteString("false")
		***REMOVED***
	case reflect.Int64:
		buf.Write(strconv.AppendInt(scratch[:0], value.Int(), 10))
	case reflect.Float64:
		f := value.Float()
		if math.IsInf(f, 0) || math.IsNaN(f) ***REMOVED***
			return &json.UnsupportedValueError***REMOVED***Value: v, Str: strconv.FormatFloat(f, 'f', -1, 64)***REMOVED***
		***REMOVED***
		buf.Write(strconv.AppendFloat(scratch[:0], f, 'f', -1, 64))
	default:
		switch converted := value.Interface().(type) ***REMOVED***
		case time.Time:
			buf.Write(strconv.AppendInt(scratch[:0], converted.UTC().Unix(), 10))
		case []byte:
			if !value.IsNil() ***REMOVED***
				buf.WriteByte('"')
				if len(converted) < 1024 ***REMOVED***
					// for small buffers, using Encode directly is much faster.
					dst := make([]byte, base64.StdEncoding.EncodedLen(len(converted)))
					base64.StdEncoding.Encode(dst, converted)
					buf.Write(dst)
				***REMOVED*** else ***REMOVED***
					// for large buffers, avoid unnecessary extra temporary
					// buffer space.
					enc := base64.NewEncoder(base64.StdEncoding, buf)
					enc.Write(converted)
					enc.Close()
				***REMOVED***
				buf.WriteByte('"')
			***REMOVED***
		case aws.JSONValue:
			str, err := protocol.EncodeJSONValue(converted, protocol.QuotedEscape)
			if err != nil ***REMOVED***
				return fmt.Errorf("unable to encode JSONValue, %v", err)
			***REMOVED***
			buf.WriteString(str)
		default:
			return fmt.Errorf("unsupported JSON value %v (%s)", value.Interface(), value.Type())
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

var hex = "0123456789abcdef"

func writeString(s string, buf *bytes.Buffer) ***REMOVED***
	buf.WriteByte('"')
	for i := 0; i < len(s); i++ ***REMOVED***
		if s[i] == '"' ***REMOVED***
			buf.WriteString(`\"`)
		***REMOVED*** else if s[i] == '\\' ***REMOVED***
			buf.WriteString(`\\`)
		***REMOVED*** else if s[i] == '\b' ***REMOVED***
			buf.WriteString(`\b`)
		***REMOVED*** else if s[i] == '\f' ***REMOVED***
			buf.WriteString(`\f`)
		***REMOVED*** else if s[i] == '\r' ***REMOVED***
			buf.WriteString(`\r`)
		***REMOVED*** else if s[i] == '\t' ***REMOVED***
			buf.WriteString(`\t`)
		***REMOVED*** else if s[i] == '\n' ***REMOVED***
			buf.WriteString(`\n`)
		***REMOVED*** else if s[i] < 32 ***REMOVED***
			buf.WriteString("\\u00")
			buf.WriteByte(hex[s[i]>>4])
			buf.WriteByte(hex[s[i]&0xF])
		***REMOVED*** else ***REMOVED***
			buf.WriteByte(s[i])
		***REMOVED***
	***REMOVED***
	buf.WriteByte('"')
***REMOVED***

// Returns the reflection element of a value, if it is a pointer.
func elemOf(value reflect.Value) reflect.Value ***REMOVED***
	for value.Kind() == reflect.Ptr ***REMOVED***
		value = value.Elem()
	***REMOVED***
	return value
***REMOVED***
