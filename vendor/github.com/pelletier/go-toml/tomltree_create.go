package toml

import (
	"fmt"
	"reflect"
	"time"
)

var kindToType = [reflect.String + 1]reflect.Type***REMOVED***
	reflect.Bool:    reflect.TypeOf(true),
	reflect.String:  reflect.TypeOf(""),
	reflect.Float32: reflect.TypeOf(float64(1)),
	reflect.Float64: reflect.TypeOf(float64(1)),
	reflect.Int:     reflect.TypeOf(int64(1)),
	reflect.Int8:    reflect.TypeOf(int64(1)),
	reflect.Int16:   reflect.TypeOf(int64(1)),
	reflect.Int32:   reflect.TypeOf(int64(1)),
	reflect.Int64:   reflect.TypeOf(int64(1)),
	reflect.Uint:    reflect.TypeOf(uint64(1)),
	reflect.Uint8:   reflect.TypeOf(uint64(1)),
	reflect.Uint16:  reflect.TypeOf(uint64(1)),
	reflect.Uint32:  reflect.TypeOf(uint64(1)),
	reflect.Uint64:  reflect.TypeOf(uint64(1)),
***REMOVED***

// typeFor returns a reflect.Type for a reflect.Kind, or nil if none is found.
// supported values:
// string, bool, int64, uint64, float64, time.Time, int, int8, int16, int32, uint, uint8, uint16, uint32, float32
func typeFor(k reflect.Kind) reflect.Type ***REMOVED***
	if k > 0 && int(k) < len(kindToType) ***REMOVED***
		return kindToType[k]
	***REMOVED***
	return nil
***REMOVED***

func simpleValueCoercion(object interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	switch original := object.(type) ***REMOVED***
	case string, bool, int64, uint64, float64, time.Time:
		return original, nil
	case int:
		return int64(original), nil
	case int8:
		return int64(original), nil
	case int16:
		return int64(original), nil
	case int32:
		return int64(original), nil
	case uint:
		return uint64(original), nil
	case uint8:
		return uint64(original), nil
	case uint16:
		return uint64(original), nil
	case uint32:
		return uint64(original), nil
	case float32:
		return float64(original), nil
	case fmt.Stringer:
		return original.String(), nil
	default:
		return nil, fmt.Errorf("cannot convert type %T to Tree", object)
	***REMOVED***
***REMOVED***

func sliceToTree(object interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	// arrays are a bit tricky, since they can represent either a
	// collection of simple values, which is represented by one
	// *tomlValue, or an array of tables, which is represented by an
	// array of *Tree.

	// holding the assumption that this function is called from toTree only when value.Kind() is Array or Slice
	value := reflect.ValueOf(object)
	insideType := value.Type().Elem()
	length := value.Len()
	if length > 0 ***REMOVED***
		insideType = reflect.ValueOf(value.Index(0).Interface()).Type()
	***REMOVED***
	if insideType.Kind() == reflect.Map ***REMOVED***
		// this is considered as an array of tables
		tablesArray := make([]*Tree, 0, length)
		for i := 0; i < length; i++ ***REMOVED***
			table := value.Index(i)
			tree, err := toTree(table.Interface())
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			tablesArray = append(tablesArray, tree.(*Tree))
		***REMOVED***
		return tablesArray, nil
	***REMOVED***

	sliceType := typeFor(insideType.Kind())
	if sliceType == nil ***REMOVED***
		sliceType = insideType
	***REMOVED***

	arrayValue := reflect.MakeSlice(reflect.SliceOf(sliceType), 0, length)

	for i := 0; i < length; i++ ***REMOVED***
		val := value.Index(i).Interface()
		simpleValue, err := simpleValueCoercion(val)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		arrayValue = reflect.Append(arrayValue, reflect.ValueOf(simpleValue))
	***REMOVED***
	return &tomlValue***REMOVED***value: arrayValue.Interface(), position: Position***REMOVED******REMOVED******REMOVED***, nil
***REMOVED***

func toTree(object interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	value := reflect.ValueOf(object)

	if value.Kind() == reflect.Map ***REMOVED***
		values := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
		keys := value.MapKeys()
		for _, key := range keys ***REMOVED***
			if key.Kind() != reflect.String ***REMOVED***
				if _, ok := key.Interface().(string); !ok ***REMOVED***
					return nil, fmt.Errorf("map key needs to be a string, not %T (%v)", key.Interface(), key.Kind())
				***REMOVED***
			***REMOVED***

			v := value.MapIndex(key)
			newValue, err := toTree(v.Interface())
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			values[key.String()] = newValue
		***REMOVED***
		return &Tree***REMOVED***values: values, position: Position***REMOVED******REMOVED******REMOVED***, nil
	***REMOVED***

	if value.Kind() == reflect.Array || value.Kind() == reflect.Slice ***REMOVED***
		return sliceToTree(object)
	***REMOVED***

	simpleValue, err := simpleValueCoercion(object)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &tomlValue***REMOVED***value: simpleValue, position: Position***REMOVED******REMOVED******REMOVED***, nil
***REMOVED***
