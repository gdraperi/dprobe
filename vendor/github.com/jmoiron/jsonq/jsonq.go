package jsonq

import (
	"fmt"
	"strconv"
)

// JsonQuery is an object that enables querying of a Go map with a simple
// positional query language.
type JsonQuery struct ***REMOVED***
	blob map[string]interface***REMOVED******REMOVED***
***REMOVED***

// stringFromInterface converts an interface***REMOVED******REMOVED*** to a string and returns an error if types don't match.
func stringFromInterface(val interface***REMOVED******REMOVED***) (string, error) ***REMOVED***
	switch val.(type) ***REMOVED***
	case string:
		return val.(string), nil
	***REMOVED***
	return "", fmt.Errorf("Expected string value for String, got \"%v\"\n", val)
***REMOVED***

// boolFromInterface converts an interface***REMOVED******REMOVED*** to a bool and returns an error if types don't match.
func boolFromInterface(val interface***REMOVED******REMOVED***) (bool, error) ***REMOVED***
	switch val.(type) ***REMOVED***
	case bool:
		return val.(bool), nil
	***REMOVED***
	return false, fmt.Errorf("Expected boolean value for Bool, got \"%v\"\n", val)
***REMOVED***

// floatFromInterface converts an interface***REMOVED******REMOVED*** to a float64 and returns an error if types don't match.
func floatFromInterface(val interface***REMOVED******REMOVED***) (float64, error) ***REMOVED***
	switch val.(type) ***REMOVED***
	case float64:
		return val.(float64), nil
	case int:
		return float64(val.(int)), nil
	case string:
		fval, err := strconv.ParseFloat(val.(string), 64)
		if err == nil ***REMOVED***
			return fval, nil
		***REMOVED***
	***REMOVED***
	return 0.0, fmt.Errorf("Expected numeric value for Float, got \"%v\"\n", val)
***REMOVED***

// intFromInterface converts an interface***REMOVED******REMOVED*** to an int and returns an error if types don't match.
func intFromInterface(val interface***REMOVED******REMOVED***) (int, error) ***REMOVED***
	switch val.(type) ***REMOVED***
	case float64:
		return int(val.(float64)), nil
	case string:
		ival, err := strconv.ParseFloat(val.(string), 64)
		if err == nil ***REMOVED***
			return int(ival), nil
		***REMOVED***
	case int:
		return val.(int), nil
	***REMOVED***
	return 0, fmt.Errorf("Expected numeric value for Int, got \"%v\"\n", val)
***REMOVED***

// objectFromInterface converts an interface***REMOVED******REMOVED*** to a map[string]interface***REMOVED******REMOVED*** and returns an error if types don't match.
func objectFromInterface(val interface***REMOVED******REMOVED***) (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	switch val.(type) ***REMOVED***
	case map[string]interface***REMOVED******REMOVED***:
		return val.(map[string]interface***REMOVED******REMOVED***), nil
	***REMOVED***
	return map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***, fmt.Errorf("Expected json object for Object, got \"%v\"\n", val)
***REMOVED***

// arrayFromInterface converts an interface***REMOVED******REMOVED*** to an []interface***REMOVED******REMOVED*** and returns an error if types don't match.
func arrayFromInterface(val interface***REMOVED******REMOVED***) ([]interface***REMOVED******REMOVED***, error) ***REMOVED***
	switch val.(type) ***REMOVED***
	case []interface***REMOVED******REMOVED***:
		return val.([]interface***REMOVED******REMOVED***), nil
	***REMOVED***
	return []interface***REMOVED******REMOVED******REMOVED******REMOVED***, fmt.Errorf("Expected json array for Array, got \"%v\"\n", val)
***REMOVED***

// NewQuery creates a new JsonQuery obj from an interface***REMOVED******REMOVED***.
func NewQuery(data interface***REMOVED******REMOVED***) *JsonQuery ***REMOVED***
	j := new(JsonQuery)
	j.blob = data.(map[string]interface***REMOVED******REMOVED***)
	return j
***REMOVED***

// Bool extracts a bool the JsonQuery
func (j *JsonQuery) Bool(s ...string) (bool, error) ***REMOVED***
	val, err := rquery(j.blob, s...)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	return boolFromInterface(val)
***REMOVED***

// Float extracts a float from the JsonQuery
func (j *JsonQuery) Float(s ...string) (float64, error) ***REMOVED***
	val, err := rquery(j.blob, s...)
	if err != nil ***REMOVED***
		return 0.0, err
	***REMOVED***
	return floatFromInterface(val)
***REMOVED***

// Int extracts an int from the JsonQuery
func (j *JsonQuery) Int(s ...string) (int, error) ***REMOVED***
	val, err := rquery(j.blob, s...)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return intFromInterface(val)
***REMOVED***

// String extracts a string from the JsonQuery
func (j *JsonQuery) String(s ...string) (string, error) ***REMOVED***
	val, err := rquery(j.blob, s...)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return stringFromInterface(val)
***REMOVED***

// Object extracts a json object from the JsonQuery
func (j *JsonQuery) Object(s ...string) (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	val, err := rquery(j.blob, s...)
	if err != nil ***REMOVED***
		return map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***, err
	***REMOVED***
	return objectFromInterface(val)
***REMOVED***

// Array extracts a []interface***REMOVED******REMOVED*** from the JsonQuery
func (j *JsonQuery) Array(s ...string) ([]interface***REMOVED******REMOVED***, error) ***REMOVED***
	val, err := rquery(j.blob, s...)
	if err != nil ***REMOVED***
		return []interface***REMOVED******REMOVED******REMOVED******REMOVED***, err
	***REMOVED***
	return arrayFromInterface(val)
***REMOVED***

// Interface extracts an interface***REMOVED******REMOVED*** from the JsonQuery
func (j *JsonQuery) Interface(s ...string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	val, err := rquery(j.blob, s...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return val, nil
***REMOVED***

// ArrayOfStrings extracts an array of strings from some json
func (j *JsonQuery) ArrayOfStrings(s ...string) ([]string, error) ***REMOVED***
	array, err := j.Array(s...)
	if err != nil ***REMOVED***
		return []string***REMOVED******REMOVED***, err
	***REMOVED***
	toReturn := make([]string, len(array))
	for index, val := range array ***REMOVED***
		toReturn[index], err = stringFromInterface(val)
		if err != nil ***REMOVED***
			return toReturn, err
		***REMOVED***
	***REMOVED***
	return toReturn, nil
***REMOVED***

// ArrayOfInts extracts an array of ints from some json
func (j *JsonQuery) ArrayOfInts(s ...string) ([]int, error) ***REMOVED***
	array, err := j.Array(s...)
	if err != nil ***REMOVED***
		return []int***REMOVED******REMOVED***, err
	***REMOVED***
	toReturn := make([]int, len(array))
	for index, val := range array ***REMOVED***
		toReturn[index], err = intFromInterface(val)
		if err != nil ***REMOVED***
			return toReturn, err
		***REMOVED***
	***REMOVED***
	return toReturn, nil
***REMOVED***

// ArrayOfFloats extracts an array of float64s from some json
func (j *JsonQuery) ArrayOfFloats(s ...string) ([]float64, error) ***REMOVED***
	array, err := j.Array(s...)
	if err != nil ***REMOVED***
		return []float64***REMOVED******REMOVED***, err
	***REMOVED***
	toReturn := make([]float64, len(array))
	for index, val := range array ***REMOVED***
		toReturn[index], err = floatFromInterface(val)
		if err != nil ***REMOVED***
			return toReturn, err
		***REMOVED***
	***REMOVED***
	return toReturn, nil
***REMOVED***

// ArrayOfBools extracts an array of bools from some json
func (j *JsonQuery) ArrayOfBools(s ...string) ([]bool, error) ***REMOVED***
	array, err := j.Array(s...)
	if err != nil ***REMOVED***
		return []bool***REMOVED******REMOVED***, err
	***REMOVED***
	toReturn := make([]bool, len(array))
	for index, val := range array ***REMOVED***
		toReturn[index], err = boolFromInterface(val)
		if err != nil ***REMOVED***
			return toReturn, err
		***REMOVED***
	***REMOVED***
	return toReturn, nil
***REMOVED***

// ArrayOfObjects extracts an array of map[string]interface***REMOVED******REMOVED*** (objects) from some json
func (j *JsonQuery) ArrayOfObjects(s ...string) ([]map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	array, err := j.Array(s...)
	if err != nil ***REMOVED***
		return []map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***, err
	***REMOVED***
	toReturn := make([]map[string]interface***REMOVED******REMOVED***, len(array))
	for index, val := range array ***REMOVED***
		toReturn[index], err = objectFromInterface(val)
		if err != nil ***REMOVED***
			return toReturn, err
		***REMOVED***
	***REMOVED***
	return toReturn, nil
***REMOVED***

// ArrayOfArrays extracts an array of []interface***REMOVED******REMOVED*** (arrays) from some json
func (j *JsonQuery) ArrayOfArrays(s ...string) ([][]interface***REMOVED******REMOVED***, error) ***REMOVED***
	array, err := j.Array(s...)
	if err != nil ***REMOVED***
		return [][]interface***REMOVED******REMOVED******REMOVED******REMOVED***, err
	***REMOVED***
	toReturn := make([][]interface***REMOVED******REMOVED***, len(array))
	for index, val := range array ***REMOVED***
		toReturn[index], err = arrayFromInterface(val)
		if err != nil ***REMOVED***
			return toReturn, err
		***REMOVED***
	***REMOVED***
	return toReturn, nil
***REMOVED***

// Matrix2D is an alias for ArrayOfArrays
func (j *JsonQuery) Matrix2D(s ...string) ([][]interface***REMOVED******REMOVED***, error) ***REMOVED***
	return j.ArrayOfArrays(s...)
***REMOVED***

// Recursively query a decoded json blob
func rquery(blob interface***REMOVED******REMOVED***, s ...string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	var (
		val interface***REMOVED******REMOVED***
		err error
	)
	val = blob
	for _, q := range s ***REMOVED***
		val, err = query(val, q)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	switch val.(type) ***REMOVED***
	case nil:
		return nil, fmt.Errorf("Nil value found at %s\n", s[len(s)-1])
	***REMOVED***
	return val, nil
***REMOVED***

// query a json blob for a single field or index.  If query is a string, then
// the blob is treated as a json object (map[string]interface***REMOVED******REMOVED***).  If query is
// an integer, the blob is treated as a json array ([]interface***REMOVED******REMOVED***).  Any kind
// of key or index error will result in a nil return value with an error set.
func query(blob interface***REMOVED******REMOVED***, query string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	index, err := strconv.Atoi(query)
	// if it's an integer, then we treat the current interface as an array
	if err == nil ***REMOVED***
		switch blob.(type) ***REMOVED***
		case []interface***REMOVED******REMOVED***:
		default:
			return nil, fmt.Errorf("Array index on non-array %v\n", blob)
		***REMOVED***
		if len(blob.([]interface***REMOVED******REMOVED***)) > index ***REMOVED***
			return blob.([]interface***REMOVED******REMOVED***)[index], nil
		***REMOVED***
		return nil, fmt.Errorf("Array index %d on array %v out of bounds\n", index, blob)
	***REMOVED***

	// blob is likely an object, but verify first
	switch blob.(type) ***REMOVED***
	case map[string]interface***REMOVED******REMOVED***:
	default:
		return nil, fmt.Errorf("Object lookup \"%s\" on non-object %v\n", query, blob)
	***REMOVED***

	val, ok := blob.(map[string]interface***REMOVED******REMOVED***)[query]
	if !ok ***REMOVED***
		return nil, fmt.Errorf("Object %v does not contain field %s\n", blob, query)
	***REMOVED***
	return val, nil
***REMOVED***
