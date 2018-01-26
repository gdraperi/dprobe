package jsonutil

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/private/protocol"
)

// UnmarshalJSON reads a stream and unmarshals the results in object v.
func UnmarshalJSON(v interface***REMOVED******REMOVED***, stream io.Reader) error ***REMOVED***
	var out interface***REMOVED******REMOVED***

	b, err := ioutil.ReadAll(stream)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if len(b) == 0 ***REMOVED***
		return nil
	***REMOVED***

	if err := json.Unmarshal(b, &out); err != nil ***REMOVED***
		return err
	***REMOVED***

	return unmarshalAny(reflect.ValueOf(v), out, "")
***REMOVED***

func unmarshalAny(value reflect.Value, data interface***REMOVED******REMOVED***, tag reflect.StructTag) error ***REMOVED***
	vtype := value.Type()
	if vtype.Kind() == reflect.Ptr ***REMOVED***
		vtype = vtype.Elem() // check kind of actual element type
	***REMOVED***

	t := tag.Get("type")
	if t == "" ***REMOVED***
		switch vtype.Kind() ***REMOVED***
		case reflect.Struct:
			// also it can't be a time object
			if _, ok := value.Interface().(*time.Time); !ok ***REMOVED***
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
		return unmarshalStruct(value, data, tag)
	case "list":
		return unmarshalList(value, data, tag)
	case "map":
		return unmarshalMap(value, data, tag)
	default:
		return unmarshalScalar(value, data, tag)
	***REMOVED***
***REMOVED***

func unmarshalStruct(value reflect.Value, data interface***REMOVED******REMOVED***, tag reflect.StructTag) error ***REMOVED***
	if data == nil ***REMOVED***
		return nil
	***REMOVED***
	mapData, ok := data.(map[string]interface***REMOVED******REMOVED***)
	if !ok ***REMOVED***
		return fmt.Errorf("JSON value is not a structure (%#v)", data)
	***REMOVED***

	t := value.Type()
	if value.Kind() == reflect.Ptr ***REMOVED***
		if value.IsNil() ***REMOVED*** // create the structure if it's nil
			s := reflect.New(value.Type().Elem())
			value.Set(s)
			value = s
		***REMOVED***

		value = value.Elem()
		t = t.Elem()
	***REMOVED***

	// unwrap any payloads
	if payload := tag.Get("payload"); payload != "" ***REMOVED***
		field, _ := t.FieldByName(payload)
		return unmarshalAny(value.FieldByName(payload), data, field.Tag)
	***REMOVED***

	for i := 0; i < t.NumField(); i++ ***REMOVED***
		field := t.Field(i)
		if field.PkgPath != "" ***REMOVED***
			continue // ignore unexported fields
		***REMOVED***

		// figure out what this field is called
		name := field.Name
		if locName := field.Tag.Get("locationName"); locName != "" ***REMOVED***
			name = locName
		***REMOVED***

		member := value.FieldByIndex(field.Index)
		err := unmarshalAny(member, mapData[name], field.Tag)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func unmarshalList(value reflect.Value, data interface***REMOVED******REMOVED***, tag reflect.StructTag) error ***REMOVED***
	if data == nil ***REMOVED***
		return nil
	***REMOVED***
	listData, ok := data.([]interface***REMOVED******REMOVED***)
	if !ok ***REMOVED***
		return fmt.Errorf("JSON value is not a list (%#v)", data)
	***REMOVED***

	if value.IsNil() ***REMOVED***
		l := len(listData)
		value.Set(reflect.MakeSlice(value.Type(), l, l))
	***REMOVED***

	for i, c := range listData ***REMOVED***
		err := unmarshalAny(value.Index(i), c, "")
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func unmarshalMap(value reflect.Value, data interface***REMOVED******REMOVED***, tag reflect.StructTag) error ***REMOVED***
	if data == nil ***REMOVED***
		return nil
	***REMOVED***
	mapData, ok := data.(map[string]interface***REMOVED******REMOVED***)
	if !ok ***REMOVED***
		return fmt.Errorf("JSON value is not a map (%#v)", data)
	***REMOVED***

	if value.IsNil() ***REMOVED***
		value.Set(reflect.MakeMap(value.Type()))
	***REMOVED***

	for k, v := range mapData ***REMOVED***
		kvalue := reflect.ValueOf(k)
		vvalue := reflect.New(value.Type().Elem()).Elem()

		unmarshalAny(vvalue, v, "")
		value.SetMapIndex(kvalue, vvalue)
	***REMOVED***

	return nil
***REMOVED***

func unmarshalScalar(value reflect.Value, data interface***REMOVED******REMOVED***, tag reflect.StructTag) error ***REMOVED***
	errf := func() error ***REMOVED***
		return fmt.Errorf("unsupported value: %v (%s)", value.Interface(), value.Type())
	***REMOVED***

	switch d := data.(type) ***REMOVED***
	case nil:
		return nil // nothing to do here
	case string:
		switch value.Interface().(type) ***REMOVED***
		case *string:
			value.Set(reflect.ValueOf(&d))
		case []byte:
			b, err := base64.StdEncoding.DecodeString(d)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			value.Set(reflect.ValueOf(b))
		case aws.JSONValue:
			// No need to use escaping as the value is a non-quoted string.
			v, err := protocol.DecodeJSONValue(d, protocol.NoEscape)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			value.Set(reflect.ValueOf(v))
		default:
			return errf()
		***REMOVED***
	case float64:
		switch value.Interface().(type) ***REMOVED***
		case *int64:
			di := int64(d)
			value.Set(reflect.ValueOf(&di))
		case *float64:
			value.Set(reflect.ValueOf(&d))
		case *time.Time:
			t := time.Unix(int64(d), 0).UTC()
			value.Set(reflect.ValueOf(&t))
		default:
			return errf()
		***REMOVED***
	case bool:
		switch value.Interface().(type) ***REMOVED***
		case *bool:
			value.Set(reflect.ValueOf(&d))
		default:
			return errf()
		***REMOVED***
	default:
		return fmt.Errorf("unsupported JSON value (%v)", data)
	***REMOVED***
	return nil
***REMOVED***
