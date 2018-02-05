package toml

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
	"time"
)

type basicMarshalTestStruct struct ***REMOVED***
	String     string                      `toml:"string"`
	StringList []string                    `toml:"strlist"`
	Sub        basicMarshalTestSubStruct   `toml:"subdoc"`
	SubList    []basicMarshalTestSubStruct `toml:"sublist"`
***REMOVED***

type basicMarshalTestSubStruct struct ***REMOVED***
	String2 string
***REMOVED***

var basicTestData = basicMarshalTestStruct***REMOVED***
	String:     "Hello",
	StringList: []string***REMOVED***"Howdy", "Hey There"***REMOVED***,
	Sub:        basicMarshalTestSubStruct***REMOVED***"One"***REMOVED***,
	SubList:    []basicMarshalTestSubStruct***REMOVED******REMOVED***"Two"***REMOVED***, ***REMOVED***"Three"***REMOVED******REMOVED***,
***REMOVED***

var basicTestToml = []byte(`string = "Hello"
strlist = ["Howdy","Hey There"]

[subdoc]
  String2 = "One"

[[sublist]]
  String2 = "Two"

[[sublist]]
  String2 = "Three"
`)

func TestBasicMarshal(t *testing.T) ***REMOVED***
	result, err := Marshal(basicTestData)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expected := basicTestToml
	if !bytes.Equal(result, expected) ***REMOVED***
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	***REMOVED***
***REMOVED***

func TestBasicUnmarshal(t *testing.T) ***REMOVED***
	result := basicMarshalTestStruct***REMOVED******REMOVED***
	err := Unmarshal(basicTestToml, &result)
	expected := basicTestData
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !reflect.DeepEqual(result, expected) ***REMOVED***
		t.Errorf("Bad unmarshal: expected %v, got %v", expected, result)
	***REMOVED***
***REMOVED***

type testDoc struct ***REMOVED***
	Title       string            `toml:"title"`
	Basics      testDocBasics     `toml:"basic"`
	BasicLists  testDocBasicLists `toml:"basic_lists"`
	BasicMap    map[string]string `toml:"basic_map"`
	Subdocs     testDocSubs       `toml:"subdoc"`
	SubDocList  []testSubDoc      `toml:"subdoclist"`
	SubDocPtrs  []*testSubDoc     `toml:"subdocptrs"`
	err         int               `toml:"shouldntBeHere"`
	unexported  int               `toml:"shouldntBeHere"`
	Unexported2 int               `toml:"-"`
***REMOVED***

type testDocBasics struct ***REMOVED***
	Bool       bool      `toml:"bool"`
	Date       time.Time `toml:"date"`
	Float      float32   `toml:"float"`
	Int        int       `toml:"int"`
	Uint       uint      `toml:"uint"`
	String     *string   `toml:"string"`
	unexported int       `toml:"shouldntBeHere"`
***REMOVED***

type testDocBasicLists struct ***REMOVED***
	Bools   []bool      `toml:"bools"`
	Dates   []time.Time `toml:"dates"`
	Floats  []*float32  `toml:"floats"`
	Ints    []int       `toml:"ints"`
	Strings []string    `toml:"strings"`
	UInts   []uint      `toml:"uints"`
***REMOVED***

type testDocSubs struct ***REMOVED***
	First  testSubDoc  `toml:"first"`
	Second *testSubDoc `toml:"second"`
***REMOVED***

type testSubDoc struct ***REMOVED***
	Name       string `toml:"name"`
	unexported int    `toml:"shouldntBeHere"`
***REMOVED***

var biteMe = "Bite me"
var float1 float32 = 12.3
var float2 float32 = 45.6
var float3 float32 = 78.9
var subdoc = testSubDoc***REMOVED***"Second", 0***REMOVED***

var docData = testDoc***REMOVED***
	Title:       "TOML Marshal Testing",
	unexported:  0,
	Unexported2: 0,
	Basics: testDocBasics***REMOVED***
		Bool:       true,
		Date:       time.Date(1979, 5, 27, 7, 32, 0, 0, time.UTC),
		Float:      123.4,
		Int:        5000,
		Uint:       5001,
		String:     &biteMe,
		unexported: 0,
	***REMOVED***,
	BasicLists: testDocBasicLists***REMOVED***
		Bools: []bool***REMOVED***true, false, true***REMOVED***,
		Dates: []time.Time***REMOVED***
			time.Date(1979, 5, 27, 7, 32, 0, 0, time.UTC),
			time.Date(1980, 5, 27, 7, 32, 0, 0, time.UTC),
		***REMOVED***,
		Floats:  []*float32***REMOVED***&float1, &float2, &float3***REMOVED***,
		Ints:    []int***REMOVED***8001, 8001, 8002***REMOVED***,
		Strings: []string***REMOVED***"One", "Two", "Three"***REMOVED***,
		UInts:   []uint***REMOVED***5002, 5003***REMOVED***,
	***REMOVED***,
	BasicMap: map[string]string***REMOVED***
		"one": "one",
		"two": "two",
	***REMOVED***,
	Subdocs: testDocSubs***REMOVED***
		First:  testSubDoc***REMOVED***"First", 0***REMOVED***,
		Second: &subdoc,
	***REMOVED***,
	SubDocList: []testSubDoc***REMOVED***
		***REMOVED***"List.First", 0***REMOVED***,
		***REMOVED***"List.Second", 0***REMOVED***,
	***REMOVED***,
	SubDocPtrs: []*testSubDoc***REMOVED***&subdoc***REMOVED***,
***REMOVED***

func TestDocMarshal(t *testing.T) ***REMOVED***
	result, err := Marshal(docData)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expected, _ := ioutil.ReadFile("marshal_test.toml")
	if !bytes.Equal(result, expected) ***REMOVED***
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	***REMOVED***
***REMOVED***

func TestDocUnmarshal(t *testing.T) ***REMOVED***
	result := testDoc***REMOVED******REMOVED***
	tomlData, _ := ioutil.ReadFile("marshal_test.toml")
	err := Unmarshal(tomlData, &result)
	expected := docData
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !reflect.DeepEqual(result, expected) ***REMOVED***
		resStr, _ := json.MarshalIndent(result, "", "  ")
		expStr, _ := json.MarshalIndent(expected, "", "  ")
		t.Errorf("Bad unmarshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expStr, resStr)
	***REMOVED***
***REMOVED***

func TestDocPartialUnmarshal(t *testing.T) ***REMOVED***
	result := testDocSubs***REMOVED******REMOVED***

	tree, _ := LoadFile("marshal_test.toml")
	subTree := tree.Get("subdoc").(*Tree)
	err := subTree.Unmarshal(&result)
	expected := docData.Subdocs
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !reflect.DeepEqual(result, expected) ***REMOVED***
		resStr, _ := json.MarshalIndent(result, "", "  ")
		expStr, _ := json.MarshalIndent(expected, "", "  ")
		t.Errorf("Bad partial unmartial: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expStr, resStr)
	***REMOVED***
***REMOVED***

type tomlTypeCheckTest struct ***REMOVED***
	name string
	item interface***REMOVED******REMOVED***
	typ  int //0=primitive, 1=otherslice, 2=treeslice, 3=tree
***REMOVED***

func TestTypeChecks(t *testing.T) ***REMOVED***
	tests := []tomlTypeCheckTest***REMOVED***
		***REMOVED***"integer", 2, 0***REMOVED***,
		***REMOVED***"time", time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), 0***REMOVED***,
		***REMOVED***"stringlist", []string***REMOVED***"hello", "hi"***REMOVED***, 1***REMOVED***,
		***REMOVED***"timelist", []time.Time***REMOVED***time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)***REMOVED***, 1***REMOVED***,
		***REMOVED***"objectlist", []tomlTypeCheckTest***REMOVED******REMOVED***, 2***REMOVED***,
		***REMOVED***"object", tomlTypeCheckTest***REMOVED******REMOVED***, 3***REMOVED***,
	***REMOVED***

	for _, test := range tests ***REMOVED***
		expected := []bool***REMOVED***false, false, false, false***REMOVED***
		expected[test.typ] = true
		result := []bool***REMOVED***
			isPrimitive(reflect.TypeOf(test.item)),
			isOtherSlice(reflect.TypeOf(test.item)),
			isTreeSlice(reflect.TypeOf(test.item)),
			isTree(reflect.TypeOf(test.item)),
		***REMOVED***
		if !reflect.DeepEqual(expected, result) ***REMOVED***
			t.Errorf("Bad type check on %q: expected %v, got %v", test.name, expected, result)
		***REMOVED***
	***REMOVED***
***REMOVED***

type unexportedMarshalTestStruct struct ***REMOVED***
	String      string                      `toml:"string"`
	StringList  []string                    `toml:"strlist"`
	Sub         basicMarshalTestSubStruct   `toml:"subdoc"`
	SubList     []basicMarshalTestSubStruct `toml:"sublist"`
	unexported  int                         `toml:"shouldntBeHere"`
	Unexported2 int                         `toml:"-"`
***REMOVED***

var unexportedTestData = unexportedMarshalTestStruct***REMOVED***
	String:      "Hello",
	StringList:  []string***REMOVED***"Howdy", "Hey There"***REMOVED***,
	Sub:         basicMarshalTestSubStruct***REMOVED***"One"***REMOVED***,
	SubList:     []basicMarshalTestSubStruct***REMOVED******REMOVED***"Two"***REMOVED***, ***REMOVED***"Three"***REMOVED******REMOVED***,
	unexported:  0,
	Unexported2: 0,
***REMOVED***

var unexportedTestToml = []byte(`string = "Hello"
strlist = ["Howdy","Hey There"]
unexported = 1
shouldntBeHere = 2

[subdoc]
  String2 = "One"

[[sublist]]
  String2 = "Two"

[[sublist]]
  String2 = "Three"
`)

func TestUnexportedUnmarshal(t *testing.T) ***REMOVED***
	result := unexportedMarshalTestStruct***REMOVED******REMOVED***
	err := Unmarshal(unexportedTestToml, &result)
	expected := unexportedTestData
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !reflect.DeepEqual(result, expected) ***REMOVED***
		t.Errorf("Bad unexported unmarshal: expected %v, got %v", expected, result)
	***REMOVED***
***REMOVED***

type errStruct struct ***REMOVED***
	Bool   bool      `toml:"bool"`
	Date   time.Time `toml:"date"`
	Float  float64   `toml:"float"`
	Int    int16     `toml:"int"`
	String *string   `toml:"string"`
***REMOVED***

var errTomls = []string***REMOVED***
	"bool = truly\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = 5000\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:3200Z\nfloat = 123.4\nint = 5000\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:32:00Z\nfloat = 123a4\nint = 5000\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = j000\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = 5000\nstring = Bite me",
	"bool = true\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = 5000\nstring = Bite me",
	"bool = 1\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = 5000\nstring = \"Bite me\"",
	"bool = true\ndate = 1\nfloat = 123.4\nint = 5000\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:32:00Z\n\"sorry\"\nint = 5000\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = \"sorry\"\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = 5000\nstring = 1",
***REMOVED***

type mapErr struct ***REMOVED***
	Vals map[string]float64
***REMOVED***

type intErr struct ***REMOVED***
	Int1  int
	Int2  int8
	Int3  int16
	Int4  int32
	Int5  int64
	UInt1 uint
	UInt2 uint8
	UInt3 uint16
	UInt4 uint32
	UInt5 uint64
	Flt1  float32
	Flt2  float64
***REMOVED***

var intErrTomls = []string***REMOVED***
	"Int1 = []\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = []\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = []\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = []\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = []\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = []\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = []\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = []\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = []\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = []\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = []\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = []",
***REMOVED***

func TestErrUnmarshal(t *testing.T) ***REMOVED***
	for ind, toml := range errTomls ***REMOVED***
		result := errStruct***REMOVED******REMOVED***
		err := Unmarshal([]byte(toml), &result)
		if err == nil ***REMOVED***
			t.Errorf("Expected err from case %d\n", ind)
		***REMOVED***
	***REMOVED***
	result2 := mapErr***REMOVED******REMOVED***
	err := Unmarshal([]byte("[Vals]\nfred=\"1.2\""), &result2)
	if err == nil ***REMOVED***
		t.Errorf("Expected err from map")
	***REMOVED***
	for ind, toml := range intErrTomls ***REMOVED***
		result3 := intErr***REMOVED******REMOVED***
		err := Unmarshal([]byte(toml), &result3)
		if err == nil ***REMOVED***
			t.Errorf("Expected int err from case %d\n", ind)
		***REMOVED***
	***REMOVED***
***REMOVED***

type emptyMarshalTestStruct struct ***REMOVED***
	Title      string                  `toml:"title"`
	Bool       bool                    `toml:"bool"`
	Int        int                     `toml:"int"`
	String     string                  `toml:"string"`
	StringList []string                `toml:"stringlist"`
	Ptr        *basicMarshalTestStruct `toml:"ptr"`
	Map        map[string]string       `toml:"map"`
***REMOVED***

var emptyTestData = emptyMarshalTestStruct***REMOVED***
	Title:      "Placeholder",
	Bool:       false,
	Int:        0,
	String:     "",
	StringList: []string***REMOVED******REMOVED***,
	Ptr:        nil,
	Map:        map[string]string***REMOVED******REMOVED***,
***REMOVED***

var emptyTestToml = []byte(`bool = false
int = 0
string = ""
stringlist = []
title = "Placeholder"

[map]
`)

type emptyMarshalTestStruct2 struct ***REMOVED***
	Title      string                  `toml:"title"`
	Bool       bool                    `toml:"bool,omitempty"`
	Int        int                     `toml:"int, omitempty"`
	String     string                  `toml:"string,omitempty "`
	StringList []string                `toml:"stringlist,omitempty"`
	Ptr        *basicMarshalTestStruct `toml:"ptr,omitempty"`
	Map        map[string]string       `toml:"map,omitempty"`
***REMOVED***

var emptyTestData2 = emptyMarshalTestStruct2***REMOVED***
	Title:      "Placeholder",
	Bool:       false,
	Int:        0,
	String:     "",
	StringList: []string***REMOVED******REMOVED***,
	Ptr:        nil,
	Map:        map[string]string***REMOVED******REMOVED***,
***REMOVED***

var emptyTestToml2 = []byte(`title = "Placeholder"
`)

func TestEmptyMarshal(t *testing.T) ***REMOVED***
	result, err := Marshal(emptyTestData)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expected := emptyTestToml
	if !bytes.Equal(result, expected) ***REMOVED***
		t.Errorf("Bad empty marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	***REMOVED***
***REMOVED***

func TestEmptyMarshalOmit(t *testing.T) ***REMOVED***
	result, err := Marshal(emptyTestData2)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expected := emptyTestToml2
	if !bytes.Equal(result, expected) ***REMOVED***
		t.Errorf("Bad empty omit marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	***REMOVED***
***REMOVED***

func TestEmptyUnmarshal(t *testing.T) ***REMOVED***
	result := emptyMarshalTestStruct***REMOVED******REMOVED***
	err := Unmarshal(emptyTestToml, &result)
	expected := emptyTestData
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !reflect.DeepEqual(result, expected) ***REMOVED***
		t.Errorf("Bad empty unmarshal: expected %v, got %v", expected, result)
	***REMOVED***
***REMOVED***

func TestEmptyUnmarshalOmit(t *testing.T) ***REMOVED***
	result := emptyMarshalTestStruct2***REMOVED******REMOVED***
	err := Unmarshal(emptyTestToml, &result)
	expected := emptyTestData2
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !reflect.DeepEqual(result, expected) ***REMOVED***
		t.Errorf("Bad empty omit unmarshal: expected %v, got %v", expected, result)
	***REMOVED***
***REMOVED***

type pointerMarshalTestStruct struct ***REMOVED***
	Str       *string
	List      *[]string
	ListPtr   *[]*string
	Map       *map[string]string
	MapPtr    *map[string]*string
	EmptyStr  *string
	EmptyList *[]string
	EmptyMap  *map[string]string
	DblPtr    *[]*[]*string
***REMOVED***

var pointerStr = "Hello"
var pointerList = []string***REMOVED***"Hello back"***REMOVED***
var pointerListPtr = []*string***REMOVED***&pointerStr***REMOVED***
var pointerMap = map[string]string***REMOVED***"response": "Goodbye"***REMOVED***
var pointerMapPtr = map[string]*string***REMOVED***"alternate": &pointerStr***REMOVED***
var pointerTestData = pointerMarshalTestStruct***REMOVED***
	Str:       &pointerStr,
	List:      &pointerList,
	ListPtr:   &pointerListPtr,
	Map:       &pointerMap,
	MapPtr:    &pointerMapPtr,
	EmptyStr:  nil,
	EmptyList: nil,
	EmptyMap:  nil,
***REMOVED***

var pointerTestToml = []byte(`List = ["Hello back"]
ListPtr = ["Hello"]
Str = "Hello"

[Map]
  response = "Goodbye"

[MapPtr]
  alternate = "Hello"
`)

func TestPointerMarshal(t *testing.T) ***REMOVED***
	result, err := Marshal(pointerTestData)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expected := pointerTestToml
	if !bytes.Equal(result, expected) ***REMOVED***
		t.Errorf("Bad pointer marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	***REMOVED***
***REMOVED***

func TestPointerUnmarshal(t *testing.T) ***REMOVED***
	result := pointerMarshalTestStruct***REMOVED******REMOVED***
	err := Unmarshal(pointerTestToml, &result)
	expected := pointerTestData
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !reflect.DeepEqual(result, expected) ***REMOVED***
		t.Errorf("Bad pointer unmarshal: expected %v, got %v", expected, result)
	***REMOVED***
***REMOVED***

func TestUnmarshalTypeMismatch(t *testing.T) ***REMOVED***
	result := pointerMarshalTestStruct***REMOVED******REMOVED***
	err := Unmarshal([]byte("List = 123"), &result)
	if !strings.HasPrefix(err.Error(), "(1, 1): Can't convert 123(int64) to []string(slice)") ***REMOVED***
		t.Errorf("Type mismatch must be reported: got %v", err.Error())
	***REMOVED***
***REMOVED***

type nestedMarshalTestStruct struct ***REMOVED***
	String [][]string
	//Struct [][]basicMarshalTestSubStruct
	StringPtr *[]*[]*string
	// StructPtr *[]*[]*basicMarshalTestSubStruct
***REMOVED***

var str1 = "Three"
var str2 = "Four"
var strPtr = []*string***REMOVED***&str1, &str2***REMOVED***
var strPtr2 = []*[]*string***REMOVED***&strPtr***REMOVED***

var nestedTestData = nestedMarshalTestStruct***REMOVED***
	String:    [][]string***REMOVED******REMOVED***"Five", "Six"***REMOVED***, ***REMOVED***"One", "Two"***REMOVED******REMOVED***,
	StringPtr: &strPtr2,
***REMOVED***

var nestedTestToml = []byte(`String = [["Five","Six"],["One","Two"]]
StringPtr = [["Three","Four"]]
`)

func TestNestedMarshal(t *testing.T) ***REMOVED***
	result, err := Marshal(nestedTestData)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expected := nestedTestToml
	if !bytes.Equal(result, expected) ***REMOVED***
		t.Errorf("Bad nested marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	***REMOVED***
***REMOVED***

func TestNestedUnmarshal(t *testing.T) ***REMOVED***
	result := nestedMarshalTestStruct***REMOVED******REMOVED***
	err := Unmarshal(nestedTestToml, &result)
	expected := nestedTestData
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !reflect.DeepEqual(result, expected) ***REMOVED***
		t.Errorf("Bad nested unmarshal: expected %v, got %v", expected, result)
	***REMOVED***
***REMOVED***

type customMarshalerParent struct ***REMOVED***
	Self    customMarshaler   `toml:"me"`
	Friends []customMarshaler `toml:"friends"`
***REMOVED***

type customMarshaler struct ***REMOVED***
	FirsName string
	LastName string
***REMOVED***

func (c customMarshaler) MarshalTOML() ([]byte, error) ***REMOVED***
	fullName := fmt.Sprintf("%s %s", c.FirsName, c.LastName)
	return []byte(fullName), nil
***REMOVED***

var customMarshalerData = customMarshaler***REMOVED***FirsName: "Sally", LastName: "Fields"***REMOVED***
var customMarshalerToml = []byte(`Sally Fields`)
var nestedCustomMarshalerData = customMarshalerParent***REMOVED***
	Self:    customMarshaler***REMOVED***FirsName: "Maiku", LastName: "Suteda"***REMOVED***,
	Friends: []customMarshaler***REMOVED***customMarshalerData***REMOVED***,
***REMOVED***
var nestedCustomMarshalerToml = []byte(`friends = ["Sally Fields"]
me = "Maiku Suteda"
`)

func TestCustomMarshaler(t *testing.T) ***REMOVED***
	result, err := Marshal(customMarshalerData)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expected := customMarshalerToml
	if !bytes.Equal(result, expected) ***REMOVED***
		t.Errorf("Bad custom marshaler: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	***REMOVED***
***REMOVED***

func TestNestedCustomMarshaler(t *testing.T) ***REMOVED***
	result, err := Marshal(nestedCustomMarshalerData)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expected := nestedCustomMarshalerToml
	if !bytes.Equal(result, expected) ***REMOVED***
		t.Errorf("Bad nested custom marshaler: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	***REMOVED***
***REMOVED***

var commentTestToml = []byte(`
# it's a comment on type
[postgres]
  # isCommented = "dvalue"
  noComment = "cvalue"

  # A comment on AttrB with a
  # break line
  password = "bvalue"

  # A comment on AttrA
  user = "avalue"

  [[postgres.My]]

    # a comment on my on typeC
    My = "Foo"

  [[postgres.My]]

    # a comment on my on typeC
    My = "Baar"
`)

func TestMarshalComment(t *testing.T) ***REMOVED***
	type TypeC struct ***REMOVED***
		My string `comment:"a comment on my on typeC"`
	***REMOVED***
	type TypeB struct ***REMOVED***
		AttrA string `toml:"user" comment:"A comment on AttrA"`
		AttrB string `toml:"password" comment:"A comment on AttrB with a\n break line"`
		AttrC string `toml:"noComment"`
		AttrD string `toml:"isCommented" commented:"true"`
		My    []TypeC
	***REMOVED***
	type TypeA struct ***REMOVED***
		TypeB TypeB `toml:"postgres" comment:"it's a comment on type"`
	***REMOVED***

	ta := []TypeC***REMOVED******REMOVED***My: "Foo"***REMOVED***, ***REMOVED***My: "Baar"***REMOVED******REMOVED***
	config := TypeA***REMOVED***TypeB***REMOVED***AttrA: "avalue", AttrB: "bvalue", AttrC: "cvalue", AttrD: "dvalue", My: ta***REMOVED******REMOVED***
	result, err := Marshal(config)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expected := commentTestToml
	if !bytes.Equal(result, expected) ***REMOVED***
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	***REMOVED***
***REMOVED***

type mapsTestStruct struct ***REMOVED***
	Simple map[string]string
	Paths  map[string]string
	Other  map[string]float64
	X      struct ***REMOVED***
		Y struct ***REMOVED***
			Z map[string]bool
		***REMOVED***
	***REMOVED***
***REMOVED***

var mapsTestData = mapsTestStruct***REMOVED***
	Simple: map[string]string***REMOVED***
		"one plus one": "two",
		"next":         "three",
	***REMOVED***,
	Paths: map[string]string***REMOVED***
		"/this/is/a/path": "/this/is/also/a/path",
		"/heloo.txt":      "/tmp/lololo.txt",
	***REMOVED***,
	Other: map[string]float64***REMOVED***
		"testing": 3.9999,
	***REMOVED***,
	X: struct***REMOVED*** Y struct***REMOVED*** Z map[string]bool ***REMOVED*** ***REMOVED******REMOVED***
		Y: struct***REMOVED*** Z map[string]bool ***REMOVED******REMOVED***
			Z: map[string]bool***REMOVED***
				"is.Nested": true,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,
***REMOVED***
var mapsTestToml = []byte(`
[Other]
  "testing" = 3.9999

[Paths]
  "/heloo.txt" = "/tmp/lololo.txt"
  "/this/is/a/path" = "/this/is/also/a/path"

[Simple]
  "next" = "three"
  "one plus one" = "two"

[X]

  [X.Y]

    [X.Y.Z]
      "is.Nested" = true
`)

func TestEncodeQuotedMapKeys(t *testing.T) ***REMOVED***
	var buf bytes.Buffer
	if err := NewEncoder(&buf).QuoteMapKeys(true).Encode(mapsTestData); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	result := buf.Bytes()
	expected := mapsTestToml
	if !bytes.Equal(result, expected) ***REMOVED***
		t.Errorf("Bad maps marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	***REMOVED***
***REMOVED***

func TestDecodeQuotedMapKeys(t *testing.T) ***REMOVED***
	result := mapsTestStruct***REMOVED******REMOVED***
	err := NewDecoder(bytes.NewBuffer(mapsTestToml)).Decode(&result)
	expected := mapsTestData
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !reflect.DeepEqual(result, expected) ***REMOVED***
		t.Errorf("Bad maps unmarshal: expected %v, got %v", expected, result)
	***REMOVED***
***REMOVED***

type structArrayNoTag struct ***REMOVED***
	A struct ***REMOVED***
		B []int64
		C []int64
	***REMOVED***
***REMOVED***

func TestMarshalArray(t *testing.T) ***REMOVED***
	expected := []byte(`
[A]
  B = [1,2,3]
  C = [1]
`)

	m := structArrayNoTag***REMOVED***
		A: struct ***REMOVED***
			B []int64
			C []int64
		***REMOVED******REMOVED***
			B: []int64***REMOVED***1, 2, 3***REMOVED***,
			C: []int64***REMOVED***1***REMOVED***,
		***REMOVED***,
	***REMOVED***

	b, err := Marshal(m)

	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if !bytes.Equal(b, expected) ***REMOVED***
		t.Errorf("Bad arrays marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, b)
	***REMOVED***
***REMOVED***

func TestMarshalArrayOnePerLine(t *testing.T) ***REMOVED***
	expected := []byte(`
[A]
  B = [
    1,
    2,
    3
  ]
  C = [1]
`)

	m := structArrayNoTag***REMOVED***
		A: struct ***REMOVED***
			B []int64
			C []int64
		***REMOVED******REMOVED***
			B: []int64***REMOVED***1, 2, 3***REMOVED***,
			C: []int64***REMOVED***1***REMOVED***,
		***REMOVED***,
	***REMOVED***

	var buf bytes.Buffer
	encoder := NewEncoder(&buf).ArraysWithOneElementPerLine(true)
	err := encoder.Encode(m)

	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	b := buf.Bytes()

	if !bytes.Equal(b, expected) ***REMOVED***
		t.Errorf("Bad arrays marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, b)
	***REMOVED***
***REMOVED***
