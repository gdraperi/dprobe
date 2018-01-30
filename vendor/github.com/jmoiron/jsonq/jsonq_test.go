package jsonq

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

const TestData = `***REMOVED***
	"foo": 1,
	"bar": 2,
	"test": "Hello, world!",
	"baz": 123.1,
	"numstring": "42",
	"floatstring": "42.1",
	"array": [
		***REMOVED***"foo": 1***REMOVED***,
		***REMOVED***"bar": 2***REMOVED***,
		***REMOVED***"baz": 3***REMOVED***
	],
	"subobj": ***REMOVED***
		"foo": 1,
		"subarray": [1,2,3],
		"subsubobj": ***REMOVED***
			"bar": 2,
			"baz": 3,
			"array": ["hello", "world"]
		***REMOVED***
	***REMOVED***,
	"collections": ***REMOVED***
		"bools": [false, true, false],
		"strings": ["hello", "strings"],
		"numbers": [1,2,3,4],
		"arrays": [[1.0,2.0],[2.0,3.0],[4.0,3.0]],
		"objects": [
			***REMOVED***"obj1": 1***REMOVED***,
			***REMOVED***"obj2": 2***REMOVED***
		]
	***REMOVED***,
	"bool": true
***REMOVED***`

func tErr(t *testing.T, err error) ***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Error: %v\n", err)
	***REMOVED***
***REMOVED***

func TestQuery(t *testing.T) ***REMOVED***
	data := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	dec := json.NewDecoder(strings.NewReader(TestData))
	err := dec.Decode(&data)
	tErr(t, err)
	q := NewQuery(data)

	ival, err := q.Int("foo")
	if ival != 1 ***REMOVED***
		t.Errorf("Expecting 1, got %v\n", ival)
	***REMOVED***
	tErr(t, err)
	ival, err = q.Int("bar")
	if ival != 2 ***REMOVED***
		t.Errorf("Expecting 2, got %v\n", ival)
	***REMOVED***
	tErr(t, err)

	ival, err = q.Int("subobj", "foo")
	if ival != 1 ***REMOVED***
		t.Errorf("Expecting 1, got %v\n", ival)
	***REMOVED***
	tErr(t, err)

	// test that strings can get int-ed
	ival, err = q.Int("numstring")
	if ival != 42 ***REMOVED***
		t.Errorf("Expecting 42, got %v\n", ival)
	***REMOVED***
	tErr(t, err)

	for i := 0; i < 3; i++ ***REMOVED***
		ival, err := q.Int("subobj", "subarray", fmt.Sprintf("%d", i))
		if ival != i+1 ***REMOVED***
			t.Errorf("Expecting %d, got %v\n", i+1, ival)
		***REMOVED***
		tErr(t, err)
	***REMOVED***

	fval, err := q.Float("baz")
	if fval != 123.1 ***REMOVED***
		t.Errorf("Expecting 123.1, got %f\n", fval)
	***REMOVED***
	tErr(t, err)

	// test that strings can get float-ed
	fval, err = q.Float("floatstring")
	if fval != 42.1 ***REMOVED***
		t.Errorf("Expecting 42.1, got %v\n", fval)
	***REMOVED***
	tErr(t, err)

	sval, err := q.String("test")
	if sval != "Hello, world!" ***REMOVED***
		t.Errorf("Expecting \"Hello, World!\", got \"%v\"\n", sval)
	***REMOVED***

	sval, err = q.String("subobj", "subsubobj", "array", "0")
	if sval != "hello" ***REMOVED***
		t.Errorf("Expecting \"hello\", got \"%s\"\n", sval)
	***REMOVED***
	tErr(t, err)

	bval, err := q.Bool("bool")
	if !bval ***REMOVED***
		t.Errorf("Expecting true, got %v\n", bval)
	***REMOVED***
	tErr(t, err)

	obj, err := q.Object("subobj", "subsubobj")
	tErr(t, err)
	q2 := NewQuery(obj)
	sval, err = q2.String("array", "1")
	if sval != "world" ***REMOVED***
		t.Errorf("Expecting \"world\", got \"%s\"\n", sval)
	***REMOVED***
	tErr(t, err)

	aobj, err := q.Array("subobj", "subarray")
	tErr(t, err)
	if aobj[0].(float64) != 1 ***REMOVED***
		t.Errorf("Expecting 1, got %v\n", aobj[0])
	***REMOVED***

	iobj, err := q.Interface("numstring")
	tErr(t, err)
	if _, ok := iobj.(string); !ok ***REMOVED***
		t.Errorf("Expecting type string got: %s", iobj)
	***REMOVED***

	/*
		Test Extraction of typed slices
	*/

	//test array of strings
	astrings, err := q.ArrayOfStrings("collections", "strings")
	tErr(t, err)
	if astrings[0] != "hello" ***REMOVED***
		t.Errorf("Expecting hello, got %v\n", astrings[0])
	***REMOVED***

	//test array of ints
	aints, err := q.ArrayOfInts("collections", "numbers")
	tErr(t, err)
	if aints[0] != 1 ***REMOVED***
		t.Errorf("Expecting 1, got %v\n", aints[0])
	***REMOVED***

	//test array of floats
	afloats, err := q.ArrayOfFloats("collections", "numbers")
	tErr(t, err)
	if afloats[0] != 1.0 ***REMOVED***
		t.Errorf("Expecting 1.0, got %v\n", afloats[0])
	***REMOVED***

	//test array of bools
	abools, err := q.ArrayOfBools("collections", "bools")
	tErr(t, err)
	if abools[0] ***REMOVED***
		t.Errorf("Expecting true, got %v\n", abools[0])
	***REMOVED***

	//test array of arrays
	aa, err := q.ArrayOfArrays("collections", "arrays")
	tErr(t, err)
	if aa[0][0].(float64) != 1 ***REMOVED***
		t.Errorf("Expecting 1, got %v\n", aa[0][0])
	***REMOVED***

	//test array of objs
	aobjs, err := q.ArrayOfObjects("collections", "objects")
	tErr(t, err)
	if aobjs[0]["obj1"].(float64) != 1 ***REMOVED***
		t.Errorf("Expecting 1, got %v\n", aobjs[0]["obj1"])
	***REMOVED***

***REMOVED***
