package mapstructure

import (
	"encoding/json"
	"io"
	"reflect"
	"sort"
	"strings"
	"testing"
)

type Basic struct ***REMOVED***
	Vstring     string
	Vint        int
	Vuint       uint
	Vbool       bool
	Vfloat      float64
	Vextra      string
	vsilent     bool
	Vdata       interface***REMOVED******REMOVED***
	VjsonInt    int
	VjsonFloat  float64
	VjsonNumber json.Number
***REMOVED***

type BasicSquash struct ***REMOVED***
	Test Basic `mapstructure:",squash"`
***REMOVED***

type Embedded struct ***REMOVED***
	Basic
	Vunique string
***REMOVED***

type EmbeddedPointer struct ***REMOVED***
	*Basic
	Vunique string
***REMOVED***

type EmbeddedSquash struct ***REMOVED***
	Basic   `mapstructure:",squash"`
	Vunique string
***REMOVED***

type SliceAlias []string

type EmbeddedSlice struct ***REMOVED***
	SliceAlias `mapstructure:"slice_alias"`
	Vunique    string
***REMOVED***

type ArrayAlias [2]string

type EmbeddedArray struct ***REMOVED***
	ArrayAlias `mapstructure:"array_alias"`
	Vunique    string
***REMOVED***

type SquashOnNonStructType struct ***REMOVED***
	InvalidSquashType int `mapstructure:",squash"`
***REMOVED***

type Map struct ***REMOVED***
	Vfoo   string
	Vother map[string]string
***REMOVED***

type MapOfStruct struct ***REMOVED***
	Value map[string]Basic
***REMOVED***

type Nested struct ***REMOVED***
	Vfoo string
	Vbar Basic
***REMOVED***

type NestedPointer struct ***REMOVED***
	Vfoo string
	Vbar *Basic
***REMOVED***

type NilInterface struct ***REMOVED***
	W io.Writer
***REMOVED***

type Slice struct ***REMOVED***
	Vfoo string
	Vbar []string
***REMOVED***

type SliceOfStruct struct ***REMOVED***
	Value []Basic
***REMOVED***

type Array struct ***REMOVED***
	Vfoo string
	Vbar [2]string
***REMOVED***

type ArrayOfStruct struct ***REMOVED***
	Value [2]Basic
***REMOVED***

type Func struct ***REMOVED***
	Foo func() string
***REMOVED***

type Tagged struct ***REMOVED***
	Extra string `mapstructure:"bar,what,what"`
	Value string `mapstructure:"foo"`
***REMOVED***

type TypeConversionResult struct ***REMOVED***
	IntToFloat         float32
	IntToUint          uint
	IntToBool          bool
	IntToString        string
	UintToInt          int
	UintToFloat        float32
	UintToBool         bool
	UintToString       string
	BoolToInt          int
	BoolToUint         uint
	BoolToFloat        float32
	BoolToString       string
	FloatToInt         int
	FloatToUint        uint
	FloatToBool        bool
	FloatToString      string
	SliceUint8ToString string
	StringToSliceUint8 []byte
	ArrayUint8ToString string
	StringToInt        int
	StringToUint       uint
	StringToBool       bool
	StringToFloat      float32
	StringToStrSlice   []string
	StringToIntSlice   []int
	StringToStrArray   [1]string
	StringToIntArray   [1]int
	SliceToMap         map[string]interface***REMOVED******REMOVED***
	MapToSlice         []interface***REMOVED******REMOVED***
	ArrayToMap         map[string]interface***REMOVED******REMOVED***
	MapToArray         [1]interface***REMOVED******REMOVED***
***REMOVED***

func TestBasicTypes(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vstring":     "foo",
		"vint":        42,
		"Vuint":       42,
		"vbool":       true,
		"Vfloat":      42.42,
		"vsilent":     true,
		"vdata":       42,
		"vjsonInt":    json.Number("1234"),
		"vjsonFloat":  json.Number("1234.5"),
		"vjsonNumber": json.Number("1234.5"),
	***REMOVED***

	var result Basic
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	***REMOVED***

	if result.Vstring != "foo" ***REMOVED***
		t.Errorf("vstring value should be 'foo': %#v", result.Vstring)
	***REMOVED***

	if result.Vint != 42 ***REMOVED***
		t.Errorf("vint value should be 42: %#v", result.Vint)
	***REMOVED***

	if result.Vuint != 42 ***REMOVED***
		t.Errorf("vuint value should be 42: %#v", result.Vuint)
	***REMOVED***

	if result.Vbool != true ***REMOVED***
		t.Errorf("vbool value should be true: %#v", result.Vbool)
	***REMOVED***

	if result.Vfloat != 42.42 ***REMOVED***
		t.Errorf("vfloat value should be 42.42: %#v", result.Vfloat)
	***REMOVED***

	if result.Vextra != "" ***REMOVED***
		t.Errorf("vextra value should be empty: %#v", result.Vextra)
	***REMOVED***

	if result.vsilent != false ***REMOVED***
		t.Error("vsilent should not be set, it is unexported")
	***REMOVED***

	if result.Vdata != 42 ***REMOVED***
		t.Error("vdata should be valid")
	***REMOVED***

	if result.VjsonInt != 1234 ***REMOVED***
		t.Errorf("vjsonint value should be 1234: %#v", result.VjsonInt)
	***REMOVED***

	if result.VjsonFloat != 1234.5 ***REMOVED***
		t.Errorf("vjsonfloat value should be 1234.5: %#v", result.VjsonFloat)
	***REMOVED***

	if !reflect.DeepEqual(result.VjsonNumber, json.Number("1234.5")) ***REMOVED***
		t.Errorf("vjsonnumber value should be '1234.5': %T, %#v", result.VjsonNumber, result.VjsonNumber)
	***REMOVED***
***REMOVED***

func TestBasic_IntWithFloat(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vint": float64(42),
	***REMOVED***

	var result Basic
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got an err: %s", err)
	***REMOVED***
***REMOVED***

func TestBasic_Merge(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vint": 42,
	***REMOVED***

	var result Basic
	result.Vuint = 100
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got an err: %s", err)
	***REMOVED***

	expected := Basic***REMOVED***
		Vint:  42,
		Vuint: 100,
	***REMOVED***
	if !reflect.DeepEqual(result, expected) ***REMOVED***
		t.Fatalf("bad: %#v", result)
	***REMOVED***
***REMOVED***

// Test for issue #46.
func TestBasic_Struct(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vdata": map[string]interface***REMOVED******REMOVED******REMOVED***
			"vstring": "foo",
		***REMOVED***,
	***REMOVED***

	var result, inner Basic
	result.Vdata = &inner
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got an err: %s", err)
	***REMOVED***
	expected := Basic***REMOVED***
		Vdata: &Basic***REMOVED***
			Vstring: "foo",
		***REMOVED***,
	***REMOVED***
	if !reflect.DeepEqual(result, expected) ***REMOVED***
		t.Fatalf("bad: %#v", result)
	***REMOVED***
***REMOVED***

func TestDecode_BasicSquash(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vstring": "foo",
	***REMOVED***

	var result BasicSquash
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got an err: %s", err.Error())
	***REMOVED***

	if result.Test.Vstring != "foo" ***REMOVED***
		t.Errorf("vstring value should be 'foo': %#v", result.Test.Vstring)
	***REMOVED***
***REMOVED***

func TestDecode_Embedded(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vstring": "foo",
		"Basic": map[string]interface***REMOVED******REMOVED******REMOVED***
			"vstring": "innerfoo",
		***REMOVED***,
		"vunique": "bar",
	***REMOVED***

	var result Embedded
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got an err: %s", err.Error())
	***REMOVED***

	if result.Vstring != "innerfoo" ***REMOVED***
		t.Errorf("vstring value should be 'innerfoo': %#v", result.Vstring)
	***REMOVED***

	if result.Vunique != "bar" ***REMOVED***
		t.Errorf("vunique value should be 'bar': %#v", result.Vunique)
	***REMOVED***
***REMOVED***

func TestDecode_EmbeddedPointer(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vstring": "foo",
		"Basic": map[string]interface***REMOVED******REMOVED******REMOVED***
			"vstring": "innerfoo",
		***REMOVED***,
		"vunique": "bar",
	***REMOVED***

	var result EmbeddedPointer
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***

	expected := EmbeddedPointer***REMOVED***
		Basic: &Basic***REMOVED***
			Vstring: "innerfoo",
		***REMOVED***,
		Vunique: "bar",
	***REMOVED***
	if !reflect.DeepEqual(result, expected) ***REMOVED***
		t.Fatalf("bad: %#v", result)
	***REMOVED***
***REMOVED***

func TestDecode_EmbeddedSlice(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"slice_alias": []string***REMOVED***"foo", "bar"***REMOVED***,
		"vunique":     "bar",
	***REMOVED***

	var result EmbeddedSlice
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got an err: %s", err.Error())
	***REMOVED***

	if !reflect.DeepEqual(result.SliceAlias, SliceAlias([]string***REMOVED***"foo", "bar"***REMOVED***)) ***REMOVED***
		t.Errorf("slice value: %#v", result.SliceAlias)
	***REMOVED***

	if result.Vunique != "bar" ***REMOVED***
		t.Errorf("vunique value should be 'bar': %#v", result.Vunique)
	***REMOVED***
***REMOVED***

func TestDecode_EmbeddedArray(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"array_alias": [2]string***REMOVED***"foo", "bar"***REMOVED***,
		"vunique":     "bar",
	***REMOVED***

	var result EmbeddedArray
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got an err: %s", err.Error())
	***REMOVED***

	if !reflect.DeepEqual(result.ArrayAlias, ArrayAlias([2]string***REMOVED***"foo", "bar"***REMOVED***)) ***REMOVED***
		t.Errorf("array value: %#v", result.ArrayAlias)
	***REMOVED***

	if result.Vunique != "bar" ***REMOVED***
		t.Errorf("vunique value should be 'bar': %#v", result.Vunique)
	***REMOVED***
***REMOVED***

func TestDecode_EmbeddedSquash(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vstring": "foo",
		"vunique": "bar",
	***REMOVED***

	var result EmbeddedSquash
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got an err: %s", err.Error())
	***REMOVED***

	if result.Vstring != "foo" ***REMOVED***
		t.Errorf("vstring value should be 'foo': %#v", result.Vstring)
	***REMOVED***

	if result.Vunique != "bar" ***REMOVED***
		t.Errorf("vunique value should be 'bar': %#v", result.Vunique)
	***REMOVED***
***REMOVED***

func TestDecode_SquashOnNonStructType(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"InvalidSquashType": 42,
	***REMOVED***

	var result SquashOnNonStructType
	err := Decode(input, &result)
	if err == nil ***REMOVED***
		t.Fatal("unexpected success decoding invalid squash field type")
	***REMOVED*** else if !strings.Contains(err.Error(), "unsupported type for squash") ***REMOVED***
		t.Fatalf("unexpected error message for invalid squash field type: %s", err)
	***REMOVED***
***REMOVED***

func TestDecode_DecodeHook(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vint": "WHAT",
	***REMOVED***

	decodeHook := func(from reflect.Kind, to reflect.Kind, v interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
		if from == reflect.String && to != reflect.String ***REMOVED***
			return 5, nil
		***REMOVED***

		return v, nil
	***REMOVED***

	var result Basic
	config := &DecoderConfig***REMOVED***
		DecodeHook: decodeHook,
		Result:     &result,
	***REMOVED***

	decoder, err := NewDecoder(config)
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***

	err = decoder.Decode(input)
	if err != nil ***REMOVED***
		t.Fatalf("got an err: %s", err)
	***REMOVED***

	if result.Vint != 5 ***REMOVED***
		t.Errorf("vint should be 5: %#v", result.Vint)
	***REMOVED***
***REMOVED***

func TestDecode_DecodeHookType(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vint": "WHAT",
	***REMOVED***

	decodeHook := func(from reflect.Type, to reflect.Type, v interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
		if from.Kind() == reflect.String &&
			to.Kind() != reflect.String ***REMOVED***
			return 5, nil
		***REMOVED***

		return v, nil
	***REMOVED***

	var result Basic
	config := &DecoderConfig***REMOVED***
		DecodeHook: decodeHook,
		Result:     &result,
	***REMOVED***

	decoder, err := NewDecoder(config)
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***

	err = decoder.Decode(input)
	if err != nil ***REMOVED***
		t.Fatalf("got an err: %s", err)
	***REMOVED***

	if result.Vint != 5 ***REMOVED***
		t.Errorf("vint should be 5: %#v", result.Vint)
	***REMOVED***
***REMOVED***

func TestDecode_Nil(t *testing.T) ***REMOVED***
	t.Parallel()

	var input interface***REMOVED******REMOVED***
	result := Basic***REMOVED***
		Vstring: "foo",
	***REMOVED***

	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***

	if result.Vstring != "foo" ***REMOVED***
		t.Fatalf("bad: %#v", result.Vstring)
	***REMOVED***
***REMOVED***

func TestDecode_NilInterfaceHook(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"w": "",
	***REMOVED***

	decodeHook := func(f, t reflect.Type, v interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
		if t.String() == "io.Writer" ***REMOVED***
			return nil, nil
		***REMOVED***

		return v, nil
	***REMOVED***

	var result NilInterface
	config := &DecoderConfig***REMOVED***
		DecodeHook: decodeHook,
		Result:     &result,
	***REMOVED***

	decoder, err := NewDecoder(config)
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***

	err = decoder.Decode(input)
	if err != nil ***REMOVED***
		t.Fatalf("got an err: %s", err)
	***REMOVED***

	if result.W != nil ***REMOVED***
		t.Errorf("W should be nil: %#v", result.W)
	***REMOVED***
***REMOVED***

func TestDecode_FuncHook(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"foo": "baz",
	***REMOVED***

	decodeHook := func(f, t reflect.Type, v interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
		if t.Kind() != reflect.Func ***REMOVED***
			return v, nil
		***REMOVED***
		val := v.(string)
		return func() string ***REMOVED*** return val ***REMOVED***, nil
	***REMOVED***

	var result Func
	config := &DecoderConfig***REMOVED***
		DecodeHook: decodeHook,
		Result:     &result,
	***REMOVED***

	decoder, err := NewDecoder(config)
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***

	err = decoder.Decode(input)
	if err != nil ***REMOVED***
		t.Fatalf("got an err: %s", err)
	***REMOVED***

	if result.Foo() != "baz" ***REMOVED***
		t.Errorf("Foo call result should be 'baz': %s", result.Foo())
	***REMOVED***
***REMOVED***

func TestDecode_NonStruct(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"foo": "bar",
		"bar": "baz",
	***REMOVED***

	var result map[string]string
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***

	if result["foo"] != "bar" ***REMOVED***
		t.Fatal("foo is not bar")
	***REMOVED***
***REMOVED***

func TestDecode_StructMatch(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vbar": Basic***REMOVED***
			Vstring: "foo",
		***REMOVED***,
	***REMOVED***

	var result Nested
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got an err: %s", err.Error())
	***REMOVED***

	if result.Vbar.Vstring != "foo" ***REMOVED***
		t.Errorf("bad: %#v", result)
	***REMOVED***
***REMOVED***

func TestDecode_TypeConversion(t *testing.T) ***REMOVED***
	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"IntToFloat":         42,
		"IntToUint":          42,
		"IntToBool":          1,
		"IntToString":        42,
		"UintToInt":          42,
		"UintToFloat":        42,
		"UintToBool":         42,
		"UintToString":       42,
		"BoolToInt":          true,
		"BoolToUint":         true,
		"BoolToFloat":        true,
		"BoolToString":       true,
		"FloatToInt":         42.42,
		"FloatToUint":        42.42,
		"FloatToBool":        42.42,
		"FloatToString":      42.42,
		"SliceUint8ToString": []uint8("foo"),
		"StringToSliceUint8": "foo",
		"ArrayUint8ToString": [3]uint8***REMOVED***'f', 'o', 'o'***REMOVED***,
		"StringToInt":        "42",
		"StringToUint":       "42",
		"StringToBool":       "1",
		"StringToFloat":      "42.42",
		"StringToStrSlice":   "A",
		"StringToIntSlice":   "42",
		"StringToStrArray":   "A",
		"StringToIntArray":   "42",
		"SliceToMap":         []interface***REMOVED******REMOVED******REMOVED******REMOVED***,
		"MapToSlice":         map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***,
		"ArrayToMap":         []interface***REMOVED******REMOVED******REMOVED******REMOVED***,
		"MapToArray":         map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***,
	***REMOVED***

	expectedResultStrict := TypeConversionResult***REMOVED***
		IntToFloat:  42.0,
		IntToUint:   42,
		UintToInt:   42,
		UintToFloat: 42,
		BoolToInt:   0,
		BoolToUint:  0,
		BoolToFloat: 0,
		FloatToInt:  42,
		FloatToUint: 42,
	***REMOVED***

	expectedResultWeak := TypeConversionResult***REMOVED***
		IntToFloat:         42.0,
		IntToUint:          42,
		IntToBool:          true,
		IntToString:        "42",
		UintToInt:          42,
		UintToFloat:        42,
		UintToBool:         true,
		UintToString:       "42",
		BoolToInt:          1,
		BoolToUint:         1,
		BoolToFloat:        1,
		BoolToString:       "1",
		FloatToInt:         42,
		FloatToUint:        42,
		FloatToBool:        true,
		FloatToString:      "42.42",
		SliceUint8ToString: "foo",
		StringToSliceUint8: []byte("foo"),
		ArrayUint8ToString: "foo",
		StringToInt:        42,
		StringToUint:       42,
		StringToBool:       true,
		StringToFloat:      42.42,
		StringToStrSlice:   []string***REMOVED***"A"***REMOVED***,
		StringToIntSlice:   []int***REMOVED***42***REMOVED***,
		StringToStrArray:   [1]string***REMOVED***"A"***REMOVED***,
		StringToIntArray:   [1]int***REMOVED***42***REMOVED***,
		SliceToMap:         map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***,
		MapToSlice:         []interface***REMOVED******REMOVED******REMOVED******REMOVED***,
		ArrayToMap:         map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***,
		MapToArray:         [1]interface***REMOVED******REMOVED******REMOVED******REMOVED***,
	***REMOVED***

	// Test strict type conversion
	var resultStrict TypeConversionResult
	err := Decode(input, &resultStrict)
	if err == nil ***REMOVED***
		t.Errorf("should return an error")
	***REMOVED***
	if !reflect.DeepEqual(resultStrict, expectedResultStrict) ***REMOVED***
		t.Errorf("expected %v, got: %v", expectedResultStrict, resultStrict)
	***REMOVED***

	// Test weak type conversion
	var decoder *Decoder
	var resultWeak TypeConversionResult

	config := &DecoderConfig***REMOVED***
		WeaklyTypedInput: true,
		Result:           &resultWeak,
	***REMOVED***

	decoder, err = NewDecoder(config)
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***

	err = decoder.Decode(input)
	if err != nil ***REMOVED***
		t.Fatalf("got an err: %s", err)
	***REMOVED***

	if !reflect.DeepEqual(resultWeak, expectedResultWeak) ***REMOVED***
		t.Errorf("expected \n%#v, got: \n%#v", expectedResultWeak, resultWeak)
	***REMOVED***
***REMOVED***

func TestDecoder_ErrorUnused(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vstring": "hello",
		"foo":     "bar",
	***REMOVED***

	var result Basic
	config := &DecoderConfig***REMOVED***
		ErrorUnused: true,
		Result:      &result,
	***REMOVED***

	decoder, err := NewDecoder(config)
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***

	err = decoder.Decode(input)
	if err == nil ***REMOVED***
		t.Fatal("expected error")
	***REMOVED***
***REMOVED***

func TestMap(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vfoo": "foo",
		"vother": map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***
			"foo": "foo",
			"bar": "bar",
		***REMOVED***,
	***REMOVED***

	var result Map
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got an error: %s", err)
	***REMOVED***

	if result.Vfoo != "foo" ***REMOVED***
		t.Errorf("vfoo value should be 'foo': %#v", result.Vfoo)
	***REMOVED***

	if result.Vother == nil ***REMOVED***
		t.Fatal("vother should not be nil")
	***REMOVED***

	if len(result.Vother) != 2 ***REMOVED***
		t.Error("vother should have two items")
	***REMOVED***

	if result.Vother["foo"] != "foo" ***REMOVED***
		t.Errorf("'foo' key should be foo, got: %#v", result.Vother["foo"])
	***REMOVED***

	if result.Vother["bar"] != "bar" ***REMOVED***
		t.Errorf("'bar' key should be bar, got: %#v", result.Vother["bar"])
	***REMOVED***
***REMOVED***

func TestMapMerge(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vfoo": "foo",
		"vother": map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***
			"foo": "foo",
			"bar": "bar",
		***REMOVED***,
	***REMOVED***

	var result Map
	result.Vother = map[string]string***REMOVED***"hello": "world"***REMOVED***
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got an error: %s", err)
	***REMOVED***

	if result.Vfoo != "foo" ***REMOVED***
		t.Errorf("vfoo value should be 'foo': %#v", result.Vfoo)
	***REMOVED***

	expected := map[string]string***REMOVED***
		"foo":   "foo",
		"bar":   "bar",
		"hello": "world",
	***REMOVED***
	if !reflect.DeepEqual(result.Vother, expected) ***REMOVED***
		t.Errorf("bad: %#v", result.Vother)
	***REMOVED***
***REMOVED***

func TestMapOfStruct(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"value": map[string]interface***REMOVED******REMOVED******REMOVED***
			"foo": map[string]string***REMOVED***"vstring": "one"***REMOVED***,
			"bar": map[string]string***REMOVED***"vstring": "two"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	var result MapOfStruct
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got an err: %s", err)
	***REMOVED***

	if result.Value == nil ***REMOVED***
		t.Fatal("value should not be nil")
	***REMOVED***

	if len(result.Value) != 2 ***REMOVED***
		t.Error("value should have two items")
	***REMOVED***

	if result.Value["foo"].Vstring != "one" ***REMOVED***
		t.Errorf("foo value should be 'one', got: %s", result.Value["foo"].Vstring)
	***REMOVED***

	if result.Value["bar"].Vstring != "two" ***REMOVED***
		t.Errorf("bar value should be 'two', got: %s", result.Value["bar"].Vstring)
	***REMOVED***
***REMOVED***

func TestNestedType(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vfoo": "foo",
		"vbar": map[string]interface***REMOVED******REMOVED******REMOVED***
			"vstring": "foo",
			"vint":    42,
			"vbool":   true,
		***REMOVED***,
	***REMOVED***

	var result Nested
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got an err: %s", err.Error())
	***REMOVED***

	if result.Vfoo != "foo" ***REMOVED***
		t.Errorf("vfoo value should be 'foo': %#v", result.Vfoo)
	***REMOVED***

	if result.Vbar.Vstring != "foo" ***REMOVED***
		t.Errorf("vstring value should be 'foo': %#v", result.Vbar.Vstring)
	***REMOVED***

	if result.Vbar.Vint != 42 ***REMOVED***
		t.Errorf("vint value should be 42: %#v", result.Vbar.Vint)
	***REMOVED***

	if result.Vbar.Vbool != true ***REMOVED***
		t.Errorf("vbool value should be true: %#v", result.Vbar.Vbool)
	***REMOVED***

	if result.Vbar.Vextra != "" ***REMOVED***
		t.Errorf("vextra value should be empty: %#v", result.Vbar.Vextra)
	***REMOVED***
***REMOVED***

func TestNestedTypePointer(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vfoo": "foo",
		"vbar": &map[string]interface***REMOVED******REMOVED******REMOVED***
			"vstring": "foo",
			"vint":    42,
			"vbool":   true,
		***REMOVED***,
	***REMOVED***

	var result NestedPointer
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got an err: %s", err.Error())
	***REMOVED***

	if result.Vfoo != "foo" ***REMOVED***
		t.Errorf("vfoo value should be 'foo': %#v", result.Vfoo)
	***REMOVED***

	if result.Vbar.Vstring != "foo" ***REMOVED***
		t.Errorf("vstring value should be 'foo': %#v", result.Vbar.Vstring)
	***REMOVED***

	if result.Vbar.Vint != 42 ***REMOVED***
		t.Errorf("vint value should be 42: %#v", result.Vbar.Vint)
	***REMOVED***

	if result.Vbar.Vbool != true ***REMOVED***
		t.Errorf("vbool value should be true: %#v", result.Vbar.Vbool)
	***REMOVED***

	if result.Vbar.Vextra != "" ***REMOVED***
		t.Errorf("vextra value should be empty: %#v", result.Vbar.Vextra)
	***REMOVED***
***REMOVED***

// Test for issue #46.
func TestNestedTypeInterface(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vfoo": "foo",
		"vbar": &map[string]interface***REMOVED******REMOVED******REMOVED***
			"vstring": "foo",
			"vint":    42,
			"vbool":   true,

			"vdata": map[string]interface***REMOVED******REMOVED******REMOVED***
				"vstring": "bar",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	var result NestedPointer
	result.Vbar = new(Basic)
	result.Vbar.Vdata = new(Basic)
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got an err: %s", err.Error())
	***REMOVED***

	if result.Vfoo != "foo" ***REMOVED***
		t.Errorf("vfoo value should be 'foo': %#v", result.Vfoo)
	***REMOVED***

	if result.Vbar.Vstring != "foo" ***REMOVED***
		t.Errorf("vstring value should be 'foo': %#v", result.Vbar.Vstring)
	***REMOVED***

	if result.Vbar.Vint != 42 ***REMOVED***
		t.Errorf("vint value should be 42: %#v", result.Vbar.Vint)
	***REMOVED***

	if result.Vbar.Vbool != true ***REMOVED***
		t.Errorf("vbool value should be true: %#v", result.Vbar.Vbool)
	***REMOVED***

	if result.Vbar.Vextra != "" ***REMOVED***
		t.Errorf("vextra value should be empty: %#v", result.Vbar.Vextra)
	***REMOVED***

	if result.Vbar.Vdata.(*Basic).Vstring != "bar" ***REMOVED***
		t.Errorf("vstring value should be 'bar': %#v", result.Vbar.Vdata.(*Basic).Vstring)
	***REMOVED***
***REMOVED***

func TestSlice(t *testing.T) ***REMOVED***
	t.Parallel()

	inputStringSlice := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vfoo": "foo",
		"vbar": []string***REMOVED***"foo", "bar", "baz"***REMOVED***,
	***REMOVED***

	inputStringSlicePointer := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vfoo": "foo",
		"vbar": &[]string***REMOVED***"foo", "bar", "baz"***REMOVED***,
	***REMOVED***

	outputStringSlice := &Slice***REMOVED***
		"foo",
		[]string***REMOVED***"foo", "bar", "baz"***REMOVED***,
	***REMOVED***

	testSliceInput(t, inputStringSlice, outputStringSlice)
	testSliceInput(t, inputStringSlicePointer, outputStringSlice)
***REMOVED***

func TestInvalidSlice(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vfoo": "foo",
		"vbar": 42,
	***REMOVED***

	result := Slice***REMOVED******REMOVED***
	err := Decode(input, &result)
	if err == nil ***REMOVED***
		t.Errorf("expected failure")
	***REMOVED***
***REMOVED***

func TestSliceOfStruct(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"value": []map[string]interface***REMOVED******REMOVED******REMOVED***
			***REMOVED***"vstring": "one"***REMOVED***,
			***REMOVED***"vstring": "two"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	var result SliceOfStruct
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got unexpected error: %s", err)
	***REMOVED***

	if len(result.Value) != 2 ***REMOVED***
		t.Fatalf("expected two values, got %d", len(result.Value))
	***REMOVED***

	if result.Value[0].Vstring != "one" ***REMOVED***
		t.Errorf("first value should be 'one', got: %s", result.Value[0].Vstring)
	***REMOVED***

	if result.Value[1].Vstring != "two" ***REMOVED***
		t.Errorf("second value should be 'two', got: %s", result.Value[1].Vstring)
	***REMOVED***
***REMOVED***

func TestSliceToMap(t *testing.T) ***REMOVED***
	t.Parallel()

	input := []map[string]interface***REMOVED******REMOVED******REMOVED***
		***REMOVED***
			"foo": "bar",
		***REMOVED***,
		***REMOVED***
			"bar": "baz",
		***REMOVED***,
	***REMOVED***

	var result map[string]interface***REMOVED******REMOVED***
	err := WeakDecode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got an error: %s", err)
	***REMOVED***

	expected := map[string]interface***REMOVED******REMOVED******REMOVED***
		"foo": "bar",
		"bar": "baz",
	***REMOVED***
	if !reflect.DeepEqual(result, expected) ***REMOVED***
		t.Errorf("bad: %#v", result)
	***REMOVED***
***REMOVED***

func TestArray(t *testing.T) ***REMOVED***
	t.Parallel()

	inputStringArray := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vfoo": "foo",
		"vbar": [2]string***REMOVED***"foo", "bar"***REMOVED***,
	***REMOVED***

	inputStringArrayPointer := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vfoo": "foo",
		"vbar": &[2]string***REMOVED***"foo", "bar"***REMOVED***,
	***REMOVED***

	outputStringArray := &Array***REMOVED***
		"foo",
		[2]string***REMOVED***"foo", "bar"***REMOVED***,
	***REMOVED***

	testArrayInput(t, inputStringArray, outputStringArray)
	testArrayInput(t, inputStringArrayPointer, outputStringArray)
***REMOVED***

func TestInvalidArray(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vfoo": "foo",
		"vbar": 42,
	***REMOVED***

	result := Array***REMOVED******REMOVED***
	err := Decode(input, &result)
	if err == nil ***REMOVED***
		t.Errorf("expected failure")
	***REMOVED***
***REMOVED***

func TestArrayOfStruct(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"value": []map[string]interface***REMOVED******REMOVED******REMOVED***
			***REMOVED***"vstring": "one"***REMOVED***,
			***REMOVED***"vstring": "two"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	var result ArrayOfStruct
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got unexpected error: %s", err)
	***REMOVED***

	if len(result.Value) != 2 ***REMOVED***
		t.Fatalf("expected two values, got %d", len(result.Value))
	***REMOVED***

	if result.Value[0].Vstring != "one" ***REMOVED***
		t.Errorf("first value should be 'one', got: %s", result.Value[0].Vstring)
	***REMOVED***

	if result.Value[1].Vstring != "two" ***REMOVED***
		t.Errorf("second value should be 'two', got: %s", result.Value[1].Vstring)
	***REMOVED***
***REMOVED***

func TestArrayToMap(t *testing.T) ***REMOVED***
	t.Parallel()

	input := []map[string]interface***REMOVED******REMOVED******REMOVED***
		***REMOVED***
			"foo": "bar",
		***REMOVED***,
		***REMOVED***
			"bar": "baz",
		***REMOVED***,
	***REMOVED***

	var result map[string]interface***REMOVED******REMOVED***
	err := WeakDecode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got an error: %s", err)
	***REMOVED***

	expected := map[string]interface***REMOVED******REMOVED******REMOVED***
		"foo": "bar",
		"bar": "baz",
	***REMOVED***
	if !reflect.DeepEqual(result, expected) ***REMOVED***
		t.Errorf("bad: %#v", result)
	***REMOVED***
***REMOVED***

func TestMapOutputForStructuredInputs(t *testing.T) ***REMOVED***
	t.Parallel()

	tests := []struct ***REMOVED***
		name    string
		in      interface***REMOVED******REMOVED***
		target  interface***REMOVED******REMOVED***
		out     interface***REMOVED******REMOVED***
		wantErr bool
	***REMOVED******REMOVED***
		***REMOVED***
			"basic struct input",
			&Basic***REMOVED***
				Vstring: "vstring",
				Vint:    2,
				Vuint:   3,
				Vbool:   true,
				Vfloat:  4.56,
				Vextra:  "vextra",
				vsilent: true,
				Vdata:   []byte("data"),
			***REMOVED***,
			&map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***,
			&map[string]interface***REMOVED******REMOVED******REMOVED***
				"Vstring":     "vstring",
				"Vint":        2,
				"Vuint":       uint(3),
				"Vbool":       true,
				"Vfloat":      4.56,
				"Vextra":      "vextra",
				"Vdata":       []byte("data"),
				"VjsonInt":    0,
				"VjsonFloat":  0.0,
				"VjsonNumber": json.Number(""),
			***REMOVED***,
			false,
		***REMOVED***,
		***REMOVED***
			"embedded struct input",
			&Embedded***REMOVED***
				Vunique: "vunique",
				Basic: Basic***REMOVED***
					Vstring: "vstring",
					Vint:    2,
					Vuint:   3,
					Vbool:   true,
					Vfloat:  4.56,
					Vextra:  "vextra",
					vsilent: true,
					Vdata:   []byte("data"),
				***REMOVED***,
			***REMOVED***,
			&map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***,
			&map[string]interface***REMOVED******REMOVED******REMOVED***
				"Vunique": "vunique",
				"Basic": map[string]interface***REMOVED******REMOVED******REMOVED***
					"Vstring":     "vstring",
					"Vint":        2,
					"Vuint":       uint(3),
					"Vbool":       true,
					"Vfloat":      4.56,
					"Vextra":      "vextra",
					"Vdata":       []byte("data"),
					"VjsonInt":    0,
					"VjsonFloat":  0.0,
					"VjsonNumber": json.Number(""),
				***REMOVED***,
			***REMOVED***,
			false,
		***REMOVED***,
		***REMOVED***
			"slice input - should error",
			[]string***REMOVED***"foo", "bar"***REMOVED***,
			&map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***,
			&map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***,
			true,
		***REMOVED***,
		***REMOVED***
			"struct with slice property",
			&Slice***REMOVED***
				Vfoo: "vfoo",
				Vbar: []string***REMOVED***"foo", "bar"***REMOVED***,
			***REMOVED***,
			&map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***,
			&map[string]interface***REMOVED******REMOVED******REMOVED***
				"Vfoo": "vfoo",
				"Vbar": []string***REMOVED***"foo", "bar"***REMOVED***,
			***REMOVED***,
			false,
		***REMOVED***,
		***REMOVED***
			"struct with slice of struct property",
			&SliceOfStruct***REMOVED***
				Value: []Basic***REMOVED***
					Basic***REMOVED***
						Vstring: "vstring",
						Vint:    2,
						Vuint:   3,
						Vbool:   true,
						Vfloat:  4.56,
						Vextra:  "vextra",
						vsilent: true,
						Vdata:   []byte("data"),
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			&map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***,
			&map[string]interface***REMOVED******REMOVED******REMOVED***
				"Value": []Basic***REMOVED***
					Basic***REMOVED***
						Vstring: "vstring",
						Vint:    2,
						Vuint:   3,
						Vbool:   true,
						Vfloat:  4.56,
						Vextra:  "vextra",
						vsilent: true,
						Vdata:   []byte("data"),
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			false,
		***REMOVED***,
		***REMOVED***
			"struct with map property",
			&Map***REMOVED***
				Vfoo:   "vfoo",
				Vother: map[string]string***REMOVED***"vother": "vother"***REMOVED***,
			***REMOVED***,
			&map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***,
			&map[string]interface***REMOVED******REMOVED******REMOVED***
				"Vfoo": "vfoo",
				"Vother": map[string]string***REMOVED***
					"vother": "vother",
				***REMOVED******REMOVED***,
			false,
		***REMOVED***,
		***REMOVED***
			"tagged struct",
			&Tagged***REMOVED***
				Extra: "extra",
				Value: "value",
			***REMOVED***,
			&map[string]string***REMOVED******REMOVED***,
			&map[string]string***REMOVED***
				"bar": "extra",
				"foo": "value",
			***REMOVED***,
			false,
		***REMOVED***,
		***REMOVED***
			"omit tag struct",
			&struct ***REMOVED***
				Value string `mapstructure:"value"`
				Omit  string `mapstructure:"-"`
			***REMOVED******REMOVED***
				Value: "value",
				Omit:  "omit",
			***REMOVED***,
			&map[string]string***REMOVED******REMOVED***,
			&map[string]string***REMOVED***
				"value": "value",
			***REMOVED***,
			false,
		***REMOVED***,
		***REMOVED***
			"decode to wrong map type",
			&struct ***REMOVED***
				Value string
			***REMOVED******REMOVED***
				Value: "string",
			***REMOVED***,
			&map[string]int***REMOVED******REMOVED***,
			&map[string]int***REMOVED******REMOVED***,
			true,
		***REMOVED***,
	***REMOVED***

	for _, tt := range tests ***REMOVED***
		t.Run(tt.name, func(t *testing.T) ***REMOVED***
			if err := Decode(tt.in, tt.target); (err != nil) != tt.wantErr ***REMOVED***
				t.Fatalf("%q: TestMapOutputForStructuredInputs() unexpected error: %s", tt.name, err)
			***REMOVED***

			if !reflect.DeepEqual(tt.out, tt.target) ***REMOVED***
				t.Fatalf("%q: TestMapOutputForStructuredInputs() expected: %#v, got: %#v", tt.name, tt.out, tt.target)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestInvalidType(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vstring": 42,
	***REMOVED***

	var result Basic
	err := Decode(input, &result)
	if err == nil ***REMOVED***
		t.Fatal("error should exist")
	***REMOVED***

	derr, ok := err.(*Error)
	if !ok ***REMOVED***
		t.Fatalf("error should be kind of Error, instead: %#v", err)
	***REMOVED***

	if derr.Errors[0] != "'Vstring' expected type 'string', got unconvertible type 'int'" ***REMOVED***
		t.Errorf("got unexpected error: %s", err)
	***REMOVED***

	inputNegIntUint := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vuint": -42,
	***REMOVED***

	err = Decode(inputNegIntUint, &result)
	if err == nil ***REMOVED***
		t.Fatal("error should exist")
	***REMOVED***

	derr, ok = err.(*Error)
	if !ok ***REMOVED***
		t.Fatalf("error should be kind of Error, instead: %#v", err)
	***REMOVED***

	if derr.Errors[0] != "cannot parse 'Vuint', -42 overflows uint" ***REMOVED***
		t.Errorf("got unexpected error: %s", err)
	***REMOVED***

	inputNegFloatUint := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vuint": -42.0,
	***REMOVED***

	err = Decode(inputNegFloatUint, &result)
	if err == nil ***REMOVED***
		t.Fatal("error should exist")
	***REMOVED***

	derr, ok = err.(*Error)
	if !ok ***REMOVED***
		t.Fatalf("error should be kind of Error, instead: %#v", err)
	***REMOVED***

	if derr.Errors[0] != "cannot parse 'Vuint', -42.000000 overflows uint" ***REMOVED***
		t.Errorf("got unexpected error: %s", err)
	***REMOVED***
***REMOVED***

func TestDecodeMetadata(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vfoo": "foo",
		"vbar": map[string]interface***REMOVED******REMOVED******REMOVED***
			"vstring": "foo",
			"Vuint":   42,
			"foo":     "bar",
		***REMOVED***,
		"bar": "nil",
	***REMOVED***

	var md Metadata
	var result Nested

	err := DecodeMetadata(input, &result, &md)
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err.Error())
	***REMOVED***

	expectedKeys := []string***REMOVED***"Vbar", "Vbar.Vstring", "Vbar.Vuint", "Vfoo"***REMOVED***
	sort.Strings(md.Keys)
	if !reflect.DeepEqual(md.Keys, expectedKeys) ***REMOVED***
		t.Fatalf("bad keys: %#v", md.Keys)
	***REMOVED***

	expectedUnused := []string***REMOVED***"Vbar.foo", "bar"***REMOVED***
	if !reflect.DeepEqual(md.Unused, expectedUnused) ***REMOVED***
		t.Fatalf("bad unused: %#v", md.Unused)
	***REMOVED***
***REMOVED***

func TestMetadata(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vfoo": "foo",
		"vbar": map[string]interface***REMOVED******REMOVED******REMOVED***
			"vstring": "foo",
			"Vuint":   42,
			"foo":     "bar",
		***REMOVED***,
		"bar": "nil",
	***REMOVED***

	var md Metadata
	var result Nested
	config := &DecoderConfig***REMOVED***
		Metadata: &md,
		Result:   &result,
	***REMOVED***

	decoder, err := NewDecoder(config)
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***

	err = decoder.Decode(input)
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err.Error())
	***REMOVED***

	expectedKeys := []string***REMOVED***"Vbar", "Vbar.Vstring", "Vbar.Vuint", "Vfoo"***REMOVED***
	sort.Strings(md.Keys)
	if !reflect.DeepEqual(md.Keys, expectedKeys) ***REMOVED***
		t.Fatalf("bad keys: %#v", md.Keys)
	***REMOVED***

	expectedUnused := []string***REMOVED***"Vbar.foo", "bar"***REMOVED***
	if !reflect.DeepEqual(md.Unused, expectedUnused) ***REMOVED***
		t.Fatalf("bad unused: %#v", md.Unused)
	***REMOVED***
***REMOVED***

func TestMetadata_Embedded(t *testing.T) ***REMOVED***
	t.Parallel()

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
		t.Fatalf("err: %s", err)
	***REMOVED***

	err = decoder.Decode(input)
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err.Error())
	***REMOVED***

	expectedKeys := []string***REMOVED***"Vstring", "Vunique"***REMOVED***

	sort.Strings(md.Keys)
	if !reflect.DeepEqual(md.Keys, expectedKeys) ***REMOVED***
		t.Fatalf("bad keys: %#v", md.Keys)
	***REMOVED***

	expectedUnused := []string***REMOVED******REMOVED***
	if !reflect.DeepEqual(md.Unused, expectedUnused) ***REMOVED***
		t.Fatalf("bad unused: %#v", md.Unused)
	***REMOVED***
***REMOVED***

func TestNonPtrValue(t *testing.T) ***REMOVED***
	t.Parallel()

	err := Decode(map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***, Basic***REMOVED******REMOVED***)
	if err == nil ***REMOVED***
		t.Fatal("error should exist")
	***REMOVED***

	if err.Error() != "result must be a pointer" ***REMOVED***
		t.Errorf("got unexpected error: %s", err)
	***REMOVED***
***REMOVED***

func TestTagged(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"foo": "bar",
		"bar": "value",
	***REMOVED***

	var result Tagged
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error: %s", err)
	***REMOVED***

	if result.Value != "bar" ***REMOVED***
		t.Errorf("value should be 'bar', got: %#v", result.Value)
	***REMOVED***

	if result.Extra != "value" ***REMOVED***
		t.Errorf("extra should be 'value', got: %#v", result.Extra)
	***REMOVED***
***REMOVED***

func TestWeakDecode(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"foo": "4",
		"bar": "value",
	***REMOVED***

	var result struct ***REMOVED***
		Foo int
		Bar string
	***REMOVED***

	if err := WeakDecode(input, &result); err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***
	if result.Foo != 4 ***REMOVED***
		t.Fatalf("bad: %#v", result)
	***REMOVED***
	if result.Bar != "value" ***REMOVED***
		t.Fatalf("bad: %#v", result)
	***REMOVED***
***REMOVED***

func TestWeakDecodeMetadata(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"foo":    "4",
		"bar":    "value",
		"unused": "value",
	***REMOVED***

	var md Metadata
	var result struct ***REMOVED***
		Foo int
		Bar string
	***REMOVED***

	if err := WeakDecodeMetadata(input, &result, &md); err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***
	if result.Foo != 4 ***REMOVED***
		t.Fatalf("bad: %#v", result)
	***REMOVED***
	if result.Bar != "value" ***REMOVED***
		t.Fatalf("bad: %#v", result)
	***REMOVED***

	expectedKeys := []string***REMOVED***"Bar", "Foo"***REMOVED***
	sort.Strings(md.Keys)
	if !reflect.DeepEqual(md.Keys, expectedKeys) ***REMOVED***
		t.Fatalf("bad keys: %#v", md.Keys)
	***REMOVED***

	expectedUnused := []string***REMOVED***"unused"***REMOVED***
	if !reflect.DeepEqual(md.Unused, expectedUnused) ***REMOVED***
		t.Fatalf("bad unused: %#v", md.Unused)
	***REMOVED***
***REMOVED***

func testSliceInput(t *testing.T, input map[string]interface***REMOVED******REMOVED***, expected *Slice) ***REMOVED***
	var result Slice
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got error: %s", err)
	***REMOVED***

	if result.Vfoo != expected.Vfoo ***REMOVED***
		t.Errorf("Vfoo expected '%s', got '%s'", expected.Vfoo, result.Vfoo)
	***REMOVED***

	if result.Vbar == nil ***REMOVED***
		t.Fatalf("Vbar a slice, got '%#v'", result.Vbar)
	***REMOVED***

	if len(result.Vbar) != len(expected.Vbar) ***REMOVED***
		t.Errorf("Vbar length should be %d, got %d", len(expected.Vbar), len(result.Vbar))
	***REMOVED***

	for i, v := range result.Vbar ***REMOVED***
		if v != expected.Vbar[i] ***REMOVED***
			t.Errorf(
				"Vbar[%d] should be '%#v', got '%#v'",
				i, expected.Vbar[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func testArrayInput(t *testing.T, input map[string]interface***REMOVED******REMOVED***, expected *Array) ***REMOVED***
	var result Array
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got error: %s", err)
	***REMOVED***

	if result.Vfoo != expected.Vfoo ***REMOVED***
		t.Errorf("Vfoo expected '%s', got '%s'", expected.Vfoo, result.Vfoo)
	***REMOVED***

	if result.Vbar == [2]string***REMOVED******REMOVED*** ***REMOVED***
		t.Fatalf("Vbar a slice, got '%#v'", result.Vbar)
	***REMOVED***

	if len(result.Vbar) != len(expected.Vbar) ***REMOVED***
		t.Errorf("Vbar length should be %d, got %d", len(expected.Vbar), len(result.Vbar))
	***REMOVED***

	for i, v := range result.Vbar ***REMOVED***
		if v != expected.Vbar[i] ***REMOVED***
			t.Errorf(
				"Vbar[%d] should be '%#v', got '%#v'",
				i, expected.Vbar[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***
