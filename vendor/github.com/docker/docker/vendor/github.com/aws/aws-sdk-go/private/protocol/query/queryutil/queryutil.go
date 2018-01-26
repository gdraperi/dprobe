package queryutil

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/private/protocol"
)

// Parse parses an object i and fills a url.Values object. The isEC2 flag
// indicates if this is the EC2 Query sub-protocol.
func Parse(body url.Values, i interface***REMOVED******REMOVED***, isEC2 bool) error ***REMOVED***
	q := queryParser***REMOVED***isEC2: isEC2***REMOVED***
	return q.parseValue(body, reflect.ValueOf(i), "", "")
***REMOVED***

func elemOf(value reflect.Value) reflect.Value ***REMOVED***
	for value.Kind() == reflect.Ptr ***REMOVED***
		value = value.Elem()
	***REMOVED***
	return value
***REMOVED***

type queryParser struct ***REMOVED***
	isEC2 bool
***REMOVED***

func (q *queryParser) parseValue(v url.Values, value reflect.Value, prefix string, tag reflect.StructTag) error ***REMOVED***
	value = elemOf(value)

	// no need to handle zero values
	if !value.IsValid() ***REMOVED***
		return nil
	***REMOVED***

	t := tag.Get("type")
	if t == "" ***REMOVED***
		switch value.Kind() ***REMOVED***
		case reflect.Struct:
			t = "structure"
		case reflect.Slice:
			t = "list"
		case reflect.Map:
			t = "map"
		***REMOVED***
	***REMOVED***

	switch t ***REMOVED***
	case "structure":
		return q.parseStruct(v, value, prefix)
	case "list":
		return q.parseList(v, value, prefix, tag)
	case "map":
		return q.parseMap(v, value, prefix, tag)
	default:
		return q.parseScalar(v, value, prefix, tag)
	***REMOVED***
***REMOVED***

func (q *queryParser) parseStruct(v url.Values, value reflect.Value, prefix string) error ***REMOVED***
	if !value.IsValid() ***REMOVED***
		return nil
	***REMOVED***

	t := value.Type()
	for i := 0; i < value.NumField(); i++ ***REMOVED***
		elemValue := elemOf(value.Field(i))
		field := t.Field(i)

		if field.PkgPath != "" ***REMOVED***
			continue // ignore unexported fields
		***REMOVED***
		if field.Tag.Get("ignore") != "" ***REMOVED***
			continue
		***REMOVED***

		if protocol.CanSetIdempotencyToken(value.Field(i), field) ***REMOVED***
			token := protocol.GetIdempotencyToken()
			elemValue = reflect.ValueOf(token)
		***REMOVED***

		var name string
		if q.isEC2 ***REMOVED***
			name = field.Tag.Get("queryName")
		***REMOVED***
		if name == "" ***REMOVED***
			if field.Tag.Get("flattened") != "" && field.Tag.Get("locationNameList") != "" ***REMOVED***
				name = field.Tag.Get("locationNameList")
			***REMOVED*** else if locName := field.Tag.Get("locationName"); locName != "" ***REMOVED***
				name = locName
			***REMOVED***
			if name != "" && q.isEC2 ***REMOVED***
				name = strings.ToUpper(name[0:1]) + name[1:]
			***REMOVED***
		***REMOVED***
		if name == "" ***REMOVED***
			name = field.Name
		***REMOVED***

		if prefix != "" ***REMOVED***
			name = prefix + "." + name
		***REMOVED***

		if err := q.parseValue(v, elemValue, name, field.Tag); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (q *queryParser) parseList(v url.Values, value reflect.Value, prefix string, tag reflect.StructTag) error ***REMOVED***
	// If it's empty, generate an empty value
	if !value.IsNil() && value.Len() == 0 ***REMOVED***
		v.Set(prefix, "")
		return nil
	***REMOVED***

	if _, ok := value.Interface().([]byte); ok ***REMOVED***
		return q.parseScalar(v, value, prefix, tag)
	***REMOVED***

	// check for unflattened list member
	if !q.isEC2 && tag.Get("flattened") == "" ***REMOVED***
		if listName := tag.Get("locationNameList"); listName == "" ***REMOVED***
			prefix += ".member"
		***REMOVED*** else ***REMOVED***
			prefix += "." + listName
		***REMOVED***
	***REMOVED***

	for i := 0; i < value.Len(); i++ ***REMOVED***
		slicePrefix := prefix
		if slicePrefix == "" ***REMOVED***
			slicePrefix = strconv.Itoa(i + 1)
		***REMOVED*** else ***REMOVED***
			slicePrefix = slicePrefix + "." + strconv.Itoa(i+1)
		***REMOVED***
		if err := q.parseValue(v, value.Index(i), slicePrefix, ""); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (q *queryParser) parseMap(v url.Values, value reflect.Value, prefix string, tag reflect.StructTag) error ***REMOVED***
	// If it's empty, generate an empty value
	if !value.IsNil() && value.Len() == 0 ***REMOVED***
		v.Set(prefix, "")
		return nil
	***REMOVED***

	// check for unflattened list member
	if !q.isEC2 && tag.Get("flattened") == "" ***REMOVED***
		prefix += ".entry"
	***REMOVED***

	// sort keys for improved serialization consistency.
	// this is not strictly necessary for protocol support.
	mapKeyValues := value.MapKeys()
	mapKeys := map[string]reflect.Value***REMOVED******REMOVED***
	mapKeyNames := make([]string, len(mapKeyValues))
	for i, mapKey := range mapKeyValues ***REMOVED***
		name := mapKey.String()
		mapKeys[name] = mapKey
		mapKeyNames[i] = name
	***REMOVED***
	sort.Strings(mapKeyNames)

	for i, mapKeyName := range mapKeyNames ***REMOVED***
		mapKey := mapKeys[mapKeyName]
		mapValue := value.MapIndex(mapKey)

		kname := tag.Get("locationNameKey")
		if kname == "" ***REMOVED***
			kname = "key"
		***REMOVED***
		vname := tag.Get("locationNameValue")
		if vname == "" ***REMOVED***
			vname = "value"
		***REMOVED***

		// serialize key
		var keyName string
		if prefix == "" ***REMOVED***
			keyName = strconv.Itoa(i+1) + "." + kname
		***REMOVED*** else ***REMOVED***
			keyName = prefix + "." + strconv.Itoa(i+1) + "." + kname
		***REMOVED***

		if err := q.parseValue(v, mapKey, keyName, ""); err != nil ***REMOVED***
			return err
		***REMOVED***

		// serialize value
		var valueName string
		if prefix == "" ***REMOVED***
			valueName = strconv.Itoa(i+1) + "." + vname
		***REMOVED*** else ***REMOVED***
			valueName = prefix + "." + strconv.Itoa(i+1) + "." + vname
		***REMOVED***

		if err := q.parseValue(v, mapValue, valueName, ""); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (q *queryParser) parseScalar(v url.Values, r reflect.Value, name string, tag reflect.StructTag) error ***REMOVED***
	switch value := r.Interface().(type) ***REMOVED***
	case string:
		v.Set(name, value)
	case []byte:
		if !r.IsNil() ***REMOVED***
			v.Set(name, base64.StdEncoding.EncodeToString(value))
		***REMOVED***
	case bool:
		v.Set(name, strconv.FormatBool(value))
	case int64:
		v.Set(name, strconv.FormatInt(value, 10))
	case int:
		v.Set(name, strconv.Itoa(value))
	case float64:
		v.Set(name, strconv.FormatFloat(value, 'f', -1, 64))
	case float32:
		v.Set(name, strconv.FormatFloat(float64(value), 'f', -1, 32))
	case time.Time:
		const ISO8601UTC = "2006-01-02T15:04:05Z"
		v.Set(name, value.UTC().Format(ISO8601UTC))
	default:
		return fmt.Errorf("unsupported value for param %s: %v (%s)", name, r.Interface(), r.Type().Name())
	***REMOVED***
	return nil
***REMOVED***
