package mapstructure

import (
	"encoding/json"
	"testing"
)

func Benchmark_Decode(b *testing.B) ***REMOVED***
	type Person struct ***REMOVED***
		Name   string
		Age    int
		Emails []string
		Extra  map[string]string
	***REMOVED***

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"name":   "Mitchell",
		"age":    91,
		"emails": []string***REMOVED***"one", "two", "three"***REMOVED***,
		"extra": map[string]string***REMOVED***
			"twitter": "mitchellh",
		***REMOVED***,
	***REMOVED***

	var result Person
	for i := 0; i < b.N; i++ ***REMOVED***
		Decode(input, &result)
	***REMOVED***
***REMOVED***

// decodeViaJSON takes the map data and passes it through encoding/json to convert it into the
// given Go native structure pointed to by v. v must be a pointer to a struct.
func decodeViaJSON(data interface***REMOVED******REMOVED***, v interface***REMOVED******REMOVED***) error ***REMOVED***
	// Perform the task by simply marshalling the input into JSON,
	// then unmarshalling it into target native Go struct.
	b, err := json.Marshal(data)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return json.Unmarshal(b, v)
***REMOVED***

func Benchmark_DecodeViaJSON(b *testing.B) ***REMOVED***
	type Person struct ***REMOVED***
		Name   string
		Age    int
		Emails []string
		Extra  map[string]string
	***REMOVED***

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"name":   "Mitchell",
		"age":    91,
		"emails": []string***REMOVED***"one", "two", "three"***REMOVED***,
		"extra": map[string]string***REMOVED***
			"twitter": "mitchellh",
		***REMOVED***,
	***REMOVED***

	var result Person
	for i := 0; i < b.N; i++ ***REMOVED***
		decodeViaJSON(input, &result)
	***REMOVED***
***REMOVED***

func Benchmark_DecodeBasic(b *testing.B) ***REMOVED***
	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vstring": "foo",
		"vint":    42,
		"Vuint":   42,
		"vbool":   true,
		"Vfloat":  42.42,
		"vsilent": true,
		"vdata":   42,
	***REMOVED***

	var result Basic
	for i := 0; i < b.N; i++ ***REMOVED***
		Decode(input, &result)
	***REMOVED***
***REMOVED***

func Benchmark_DecodeEmbedded(b *testing.B) ***REMOVED***
	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vstring": "foo",
		"Basic": map[string]interface***REMOVED******REMOVED******REMOVED***
			"vstring": "innerfoo",
		***REMOVED***,
		"vunique": "bar",
	***REMOVED***

	var result Embedded
	for i := 0; i < b.N; i++ ***REMOVED***
		Decode(input, &result)
	***REMOVED***
***REMOVED***

func Benchmark_DecodeTypeConversion(b *testing.B) ***REMOVED***
	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"IntToFloat":    42,
		"IntToUint":     42,
		"IntToBool":     1,
		"IntToString":   42,
		"UintToInt":     42,
		"UintToFloat":   42,
		"UintToBool":    42,
		"UintToString":  42,
		"BoolToInt":     true,
		"BoolToUint":    true,
		"BoolToFloat":   true,
		"BoolToString":  true,
		"FloatToInt":    42.42,
		"FloatToUint":   42.42,
		"FloatToBool":   42.42,
		"FloatToString": 42.42,
		"StringToInt":   "42",
		"StringToUint":  "42",
		"StringToBool":  "1",
		"StringToFloat": "42.42",
		"SliceToMap":    []interface***REMOVED******REMOVED******REMOVED******REMOVED***,
		"MapToSlice":    map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***,
	***REMOVED***

	var resultStrict TypeConversionResult
	for i := 0; i < b.N; i++ ***REMOVED***
		Decode(input, &resultStrict)
	***REMOVED***
***REMOVED***

func Benchmark_DecodeMap(b *testing.B) ***REMOVED***
	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vfoo": "foo",
		"vother": map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***
			"foo": "foo",
			"bar": "bar",
		***REMOVED***,
	***REMOVED***

	var result Map
	for i := 0; i < b.N; i++ ***REMOVED***
		Decode(input, &result)
	***REMOVED***
***REMOVED***

func Benchmark_DecodeMapOfStruct(b *testing.B) ***REMOVED***
	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"value": map[string]interface***REMOVED******REMOVED******REMOVED***
			"foo": map[string]string***REMOVED***"vstring": "one"***REMOVED***,
			"bar": map[string]string***REMOVED***"vstring": "two"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	var result MapOfStruct
	for i := 0; i < b.N; i++ ***REMOVED***
		Decode(input, &result)
	***REMOVED***
***REMOVED***

func Benchmark_DecodeSlice(b *testing.B) ***REMOVED***
	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vfoo": "foo",
		"vbar": []string***REMOVED***"foo", "bar", "baz"***REMOVED***,
	***REMOVED***

	var result Slice
	for i := 0; i < b.N; i++ ***REMOVED***
		Decode(input, &result)
	***REMOVED***
***REMOVED***

func Benchmark_DecodeSliceOfStruct(b *testing.B) ***REMOVED***
	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"value": []map[string]interface***REMOVED******REMOVED******REMOVED***
			***REMOVED***"vstring": "one"***REMOVED***,
			***REMOVED***"vstring": "two"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	var result SliceOfStruct
	for i := 0; i < b.N; i++ ***REMOVED***
		Decode(input, &result)
	***REMOVED***
***REMOVED***

func Benchmark_DecodeWeaklyTypedInput(b *testing.B) ***REMOVED***
	type Person struct ***REMOVED***
		Name   string
		Age    int
		Emails []string
	***REMOVED***

	// This input can come from anywhere, but typically comes from
	// something like decoding JSON, generated by a weakly typed language
	// such as PHP.
	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"name":   123,                      // number => string
		"age":    "42",                     // string => number
		"emails": map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***, // empty map => empty array
	***REMOVED***

	var result Person
	config := &DecoderConfig***REMOVED***
		WeaklyTypedInput: true,
		Result:           &result,
	***REMOVED***

	decoder, err := NewDecoder(config)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	for i := 0; i < b.N; i++ ***REMOVED***
		decoder.Decode(input)
	***REMOVED***
***REMOVED***

func Benchmark_DecodeMetadata(b *testing.B) ***REMOVED***
	type Person struct ***REMOVED***
		Name string
		Age  int
	***REMOVED***

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"name":  "Mitchell",
		"age":   91,
		"email": "foo@bar.com",
	***REMOVED***

	var md Metadata
	var result Person
	config := &DecoderConfig***REMOVED***
		Metadata: &md,
		Result:   &result,
	***REMOVED***

	decoder, err := NewDecoder(config)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	for i := 0; i < b.N; i++ ***REMOVED***
		decoder.Decode(input)
	***REMOVED***
***REMOVED***

func Benchmark_DecodeMetadataEmbedded(b *testing.B) ***REMOVED***
	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vstring": "foo",
		"vunique": "bar",
	***REMOVED***

	var md Metadata
	var result EmbeddedSquash
	config := &DecoderConfig***REMOVED***
		Metadata: &md,
		Result:   &result,
	***REMOVED***

	decoder, err := NewDecoder(config)
	if err != nil ***REMOVED***
		b.Fatalf("err: %s", err)
	***REMOVED***

	for i := 0; i < b.N; i++ ***REMOVED***
		decoder.Decode(input)
	***REMOVED***
***REMOVED***

func Benchmark_DecodeTagged(b *testing.B) ***REMOVED***
	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"foo": "bar",
		"bar": "value",
	***REMOVED***

	var result Tagged
	for i := 0; i < b.N; i++ ***REMOVED***
		Decode(input, &result)
	***REMOVED***
***REMOVED***
