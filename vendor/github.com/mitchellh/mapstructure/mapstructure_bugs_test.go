package mapstructure

import "testing"

// GH-1
func TestDecode_NilValue(t *testing.T) ***REMOVED***
	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vfoo":   nil,
		"vother": nil,
	***REMOVED***

	var result Map
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("should not error: %s", err)
	***REMOVED***

	if result.Vfoo != "" ***REMOVED***
		t.Fatalf("value should be default: %s", result.Vfoo)
	***REMOVED***

	if result.Vother != nil ***REMOVED***
		t.Fatalf("Vother should be nil: %s", result.Vother)
	***REMOVED***
***REMOVED***

// GH-10
func TestDecode_mapInterfaceInterface(t *testing.T) ***REMOVED***
	input := map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***
		"vfoo":   nil,
		"vother": nil,
	***REMOVED***

	var result Map
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("should not error: %s", err)
	***REMOVED***

	if result.Vfoo != "" ***REMOVED***
		t.Fatalf("value should be default: %s", result.Vfoo)
	***REMOVED***

	if result.Vother != nil ***REMOVED***
		t.Fatalf("Vother should be nil: %s", result.Vother)
	***REMOVED***
***REMOVED***

// #48
func TestNestedTypePointerWithDefaults(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vfoo": "foo",
		"vbar": map[string]interface***REMOVED******REMOVED******REMOVED***
			"vstring": "foo",
			"vint":    42,
			"vbool":   true,
		***REMOVED***,
	***REMOVED***

	result := NestedPointer***REMOVED***
		Vbar: &Basic***REMOVED***
			Vuint: 42,
		***REMOVED***,
	***REMOVED***
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

	// this is the error
	if result.Vbar.Vuint != 42 ***REMOVED***
		t.Errorf("vuint value should be 42: %#v", result.Vbar.Vuint)
	***REMOVED***

***REMOVED***

type NestedSlice struct ***REMOVED***
	Vfoo   string
	Vbars  []Basic
	Vempty []Basic
***REMOVED***

// #48
func TestNestedTypeSliceWithDefaults(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vfoo": "foo",
		"vbars": []map[string]interface***REMOVED******REMOVED******REMOVED***
			***REMOVED***"vstring": "foo", "vint": 42, "vbool": true***REMOVED***,
			***REMOVED***"vint": 42, "vbool": true***REMOVED***,
		***REMOVED***,
		"vempty": []map[string]interface***REMOVED******REMOVED******REMOVED***
			***REMOVED***"vstring": "foo", "vint": 42, "vbool": true***REMOVED***,
			***REMOVED***"vint": 42, "vbool": true***REMOVED***,
		***REMOVED***,
	***REMOVED***

	result := NestedSlice***REMOVED***
		Vbars: []Basic***REMOVED***
			***REMOVED***Vuint: 42***REMOVED***,
			***REMOVED***Vstring: "foo"***REMOVED***,
		***REMOVED***,
	***REMOVED***
	err := Decode(input, &result)
	if err != nil ***REMOVED***
		t.Fatalf("got an err: %s", err.Error())
	***REMOVED***

	if result.Vfoo != "foo" ***REMOVED***
		t.Errorf("vfoo value should be 'foo': %#v", result.Vfoo)
	***REMOVED***

	if result.Vbars[0].Vstring != "foo" ***REMOVED***
		t.Errorf("vstring value should be 'foo': %#v", result.Vbars[0].Vstring)
	***REMOVED***
	// this is the error
	if result.Vbars[0].Vuint != 42 ***REMOVED***
		t.Errorf("vuint value should be 42: %#v", result.Vbars[0].Vuint)
	***REMOVED***
***REMOVED***

// #48 workaround
func TestNestedTypeWithDefaults(t *testing.T) ***REMOVED***
	t.Parallel()

	input := map[string]interface***REMOVED******REMOVED******REMOVED***
		"vfoo": "foo",
		"vbar": map[string]interface***REMOVED******REMOVED******REMOVED***
			"vstring": "foo",
			"vint":    42,
			"vbool":   true,
		***REMOVED***,
	***REMOVED***

	result := Nested***REMOVED***
		Vbar: Basic***REMOVED***
			Vuint: 42,
		***REMOVED***,
	***REMOVED***
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

	// this is the error
	if result.Vbar.Vuint != 42 ***REMOVED***
		t.Errorf("vuint value should be 42: %#v", result.Vbar.Vuint)
	***REMOVED***

***REMOVED***

// #67 panic() on extending slices (decodeSlice with disabled ZeroValues)
func TestDecodeSliceToEmptySliceWOZeroing(t *testing.T) ***REMOVED***
	t.Parallel()

	type TestStruct struct ***REMOVED***
		Vfoo []string
	***REMOVED***

	decode := func(m interface***REMOVED******REMOVED***, rawVal interface***REMOVED******REMOVED***) error ***REMOVED***
		config := &DecoderConfig***REMOVED***
			Metadata:   nil,
			Result:     rawVal,
			ZeroFields: false,
		***REMOVED***

		decoder, err := NewDecoder(config)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		return decoder.Decode(m)
	***REMOVED***

	***REMOVED***
		input := map[string]interface***REMOVED******REMOVED******REMOVED***
			"vfoo": []string***REMOVED***"1"***REMOVED***,
		***REMOVED***

		result := &TestStruct***REMOVED******REMOVED***

		err := decode(input, &result)
		if err != nil ***REMOVED***
			t.Fatalf("got an err: %s", err.Error())
		***REMOVED***
	***REMOVED***

	***REMOVED***
		input := map[string]interface***REMOVED******REMOVED******REMOVED***
			"vfoo": []string***REMOVED***"1"***REMOVED***,
		***REMOVED***

		result := &TestStruct***REMOVED***
			Vfoo: []string***REMOVED******REMOVED***,
		***REMOVED***

		err := decode(input, &result)
		if err != nil ***REMOVED***
			t.Fatalf("got an err: %s", err.Error())
		***REMOVED***
	***REMOVED***

	***REMOVED***
		input := map[string]interface***REMOVED******REMOVED******REMOVED***
			"vfoo": []string***REMOVED***"2", "3"***REMOVED***,
		***REMOVED***

		result := &TestStruct***REMOVED***
			Vfoo: []string***REMOVED***"1"***REMOVED***,
		***REMOVED***

		err := decode(input, &result)
		if err != nil ***REMOVED***
			t.Fatalf("got an err: %s", err.Error())
		***REMOVED***
	***REMOVED***
***REMOVED***

// #70
func TestNextSquashMapstructure(t *testing.T) ***REMOVED***
	data := &struct ***REMOVED***
		Level1 struct ***REMOVED***
			Level2 struct ***REMOVED***
				Foo string
			***REMOVED*** `mapstructure:",squash"`
		***REMOVED*** `mapstructure:",squash"`
	***REMOVED******REMOVED******REMOVED***
	err := Decode(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***"foo": "baz"***REMOVED***, &data)
	if err != nil ***REMOVED***
		t.Fatalf("should not error: %s", err)
	***REMOVED***
	if data.Level1.Level2.Foo != "baz" ***REMOVED***
		t.Fatal("value should be baz")
	***REMOVED***
***REMOVED***
