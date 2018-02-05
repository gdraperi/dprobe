package yaml_test

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	. "gopkg.in/check.v1"
	"gopkg.in/yaml.v2"
	"net"
	"os"
)

var marshalIntTest = 123

var marshalTests = []struct ***REMOVED***
	value interface***REMOVED******REMOVED***
	data  string
***REMOVED******REMOVED***
	***REMOVED***
		nil,
		"null\n",
	***REMOVED***, ***REMOVED***
		&struct***REMOVED******REMOVED******REMOVED******REMOVED***,
		"***REMOVED******REMOVED***\n",
	***REMOVED***, ***REMOVED***
		map[string]string***REMOVED***"v": "hi"***REMOVED***,
		"v: hi\n",
	***REMOVED***, ***REMOVED***
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": "hi"***REMOVED***,
		"v: hi\n",
	***REMOVED***, ***REMOVED***
		map[string]string***REMOVED***"v": "true"***REMOVED***,
		"v: \"true\"\n",
	***REMOVED***, ***REMOVED***
		map[string]string***REMOVED***"v": "false"***REMOVED***,
		"v: \"false\"\n",
	***REMOVED***, ***REMOVED***
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": true***REMOVED***,
		"v: true\n",
	***REMOVED***, ***REMOVED***
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": false***REMOVED***,
		"v: false\n",
	***REMOVED***, ***REMOVED***
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": 10***REMOVED***,
		"v: 10\n",
	***REMOVED***, ***REMOVED***
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": -10***REMOVED***,
		"v: -10\n",
	***REMOVED***, ***REMOVED***
		map[string]uint***REMOVED***"v": 42***REMOVED***,
		"v: 42\n",
	***REMOVED***, ***REMOVED***
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": int64(4294967296)***REMOVED***,
		"v: 4294967296\n",
	***REMOVED***, ***REMOVED***
		map[string]int64***REMOVED***"v": int64(4294967296)***REMOVED***,
		"v: 4294967296\n",
	***REMOVED***, ***REMOVED***
		map[string]uint64***REMOVED***"v": 4294967296***REMOVED***,
		"v: 4294967296\n",
	***REMOVED***, ***REMOVED***
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": "10"***REMOVED***,
		"v: \"10\"\n",
	***REMOVED***, ***REMOVED***
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": 0.1***REMOVED***,
		"v: 0.1\n",
	***REMOVED***, ***REMOVED***
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": float64(0.1)***REMOVED***,
		"v: 0.1\n",
	***REMOVED***, ***REMOVED***
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": -0.1***REMOVED***,
		"v: -0.1\n",
	***REMOVED***, ***REMOVED***
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": math.Inf(+1)***REMOVED***,
		"v: .inf\n",
	***REMOVED***, ***REMOVED***
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": math.Inf(-1)***REMOVED***,
		"v: -.inf\n",
	***REMOVED***, ***REMOVED***
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": math.NaN()***REMOVED***,
		"v: .nan\n",
	***REMOVED***, ***REMOVED***
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": nil***REMOVED***,
		"v: null\n",
	***REMOVED***, ***REMOVED***
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": ""***REMOVED***,
		"v: \"\"\n",
	***REMOVED***, ***REMOVED***
		map[string][]string***REMOVED***"v": []string***REMOVED***"A", "B"***REMOVED******REMOVED***,
		"v:\n- A\n- B\n",
	***REMOVED***, ***REMOVED***
		map[string][]string***REMOVED***"v": []string***REMOVED***"A", "B\nC"***REMOVED******REMOVED***,
		"v:\n- A\n- |-\n  B\n  C\n",
	***REMOVED***, ***REMOVED***
		map[string][]interface***REMOVED******REMOVED******REMOVED***"v": []interface***REMOVED******REMOVED******REMOVED***"A", 1, map[string][]int***REMOVED***"B": []int***REMOVED***2, 3***REMOVED******REMOVED******REMOVED******REMOVED***,
		"v:\n- A\n- 1\n- B:\n  - 2\n  - 3\n",
	***REMOVED***, ***REMOVED***
		map[string]interface***REMOVED******REMOVED******REMOVED***"a": map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***"b": "c"***REMOVED******REMOVED***,
		"a:\n  b: c\n",
	***REMOVED***, ***REMOVED***
		map[string]interface***REMOVED******REMOVED******REMOVED***"a": "-"***REMOVED***,
		"a: '-'\n",
	***REMOVED***,

	// Simple values.
	***REMOVED***
		&marshalIntTest,
		"123\n",
	***REMOVED***,

	// Structures
	***REMOVED***
		&struct***REMOVED*** Hello string ***REMOVED******REMOVED***"world"***REMOVED***,
		"hello: world\n",
	***REMOVED***, ***REMOVED***
		&struct ***REMOVED***
			A struct ***REMOVED***
				B string
			***REMOVED***
		***REMOVED******REMOVED***struct***REMOVED*** B string ***REMOVED******REMOVED***"c"***REMOVED******REMOVED***,
		"a:\n  b: c\n",
	***REMOVED***, ***REMOVED***
		&struct ***REMOVED***
			A *struct ***REMOVED***
				B string
			***REMOVED***
		***REMOVED******REMOVED***&struct***REMOVED*** B string ***REMOVED******REMOVED***"c"***REMOVED******REMOVED***,
		"a:\n  b: c\n",
	***REMOVED***, ***REMOVED***
		&struct ***REMOVED***
			A *struct ***REMOVED***
				B string
			***REMOVED***
		***REMOVED******REMOVED******REMOVED***,
		"a: null\n",
	***REMOVED***, ***REMOVED***
		&struct***REMOVED*** A int ***REMOVED******REMOVED***1***REMOVED***,
		"a: 1\n",
	***REMOVED***, ***REMOVED***
		&struct***REMOVED*** A []int ***REMOVED******REMOVED***[]int***REMOVED***1, 2***REMOVED******REMOVED***,
		"a:\n- 1\n- 2\n",
	***REMOVED***, ***REMOVED***
		&struct ***REMOVED***
			B int "a"
		***REMOVED******REMOVED***1***REMOVED***,
		"a: 1\n",
	***REMOVED***, ***REMOVED***
		&struct***REMOVED*** A bool ***REMOVED******REMOVED***true***REMOVED***,
		"a: true\n",
	***REMOVED***,

	// Conditional flag
	***REMOVED***
		&struct ***REMOVED***
			A int "a,omitempty"
			B int "b,omitempty"
		***REMOVED******REMOVED***1, 0***REMOVED***,
		"a: 1\n",
	***REMOVED***, ***REMOVED***
		&struct ***REMOVED***
			A int "a,omitempty"
			B int "b,omitempty"
		***REMOVED******REMOVED***0, 0***REMOVED***,
		"***REMOVED******REMOVED***\n",
	***REMOVED***, ***REMOVED***
		&struct ***REMOVED***
			A *struct***REMOVED*** X, y int ***REMOVED*** "a,omitempty,flow"
		***REMOVED******REMOVED***&struct***REMOVED*** X, y int ***REMOVED******REMOVED***1, 2***REMOVED******REMOVED***,
		"a: ***REMOVED***x: 1***REMOVED***\n",
	***REMOVED***, ***REMOVED***
		&struct ***REMOVED***
			A *struct***REMOVED*** X, y int ***REMOVED*** "a,omitempty,flow"
		***REMOVED******REMOVED***nil***REMOVED***,
		"***REMOVED******REMOVED***\n",
	***REMOVED***, ***REMOVED***
		&struct ***REMOVED***
			A *struct***REMOVED*** X, y int ***REMOVED*** "a,omitempty,flow"
		***REMOVED******REMOVED***&struct***REMOVED*** X, y int ***REMOVED******REMOVED******REMOVED******REMOVED***,
		"a: ***REMOVED***x: 0***REMOVED***\n",
	***REMOVED***, ***REMOVED***
		&struct ***REMOVED***
			A struct***REMOVED*** X, y int ***REMOVED*** "a,omitempty,flow"
		***REMOVED******REMOVED***struct***REMOVED*** X, y int ***REMOVED******REMOVED***1, 2***REMOVED******REMOVED***,
		"a: ***REMOVED***x: 1***REMOVED***\n",
	***REMOVED***, ***REMOVED***
		&struct ***REMOVED***
			A struct***REMOVED*** X, y int ***REMOVED*** "a,omitempty,flow"
		***REMOVED******REMOVED***struct***REMOVED*** X, y int ***REMOVED******REMOVED***0, 1***REMOVED******REMOVED***,
		"***REMOVED******REMOVED***\n",
	***REMOVED***, ***REMOVED***
		&struct ***REMOVED***
			A float64 "a,omitempty"
			B float64 "b,omitempty"
		***REMOVED******REMOVED***1, 0***REMOVED***,
		"a: 1\n",
	***REMOVED***,

	// Flow flag
	***REMOVED***
		&struct ***REMOVED***
			A []int "a,flow"
		***REMOVED******REMOVED***[]int***REMOVED***1, 2***REMOVED******REMOVED***,
		"a: [1, 2]\n",
	***REMOVED***, ***REMOVED***
		&struct ***REMOVED***
			A map[string]string "a,flow"
		***REMOVED******REMOVED***map[string]string***REMOVED***"b": "c", "d": "e"***REMOVED******REMOVED***,
		"a: ***REMOVED***b: c, d: e***REMOVED***\n",
	***REMOVED***, ***REMOVED***
		&struct ***REMOVED***
			A struct ***REMOVED***
				B, D string
			***REMOVED*** "a,flow"
		***REMOVED******REMOVED***struct***REMOVED*** B, D string ***REMOVED******REMOVED***"c", "e"***REMOVED******REMOVED***,
		"a: ***REMOVED***b: c, d: e***REMOVED***\n",
	***REMOVED***,

	// Unexported field
	***REMOVED***
		&struct ***REMOVED***
			u int
			A int
		***REMOVED******REMOVED***0, 1***REMOVED***,
		"a: 1\n",
	***REMOVED***,

	// Ignored field
	***REMOVED***
		&struct ***REMOVED***
			A int
			B int "-"
		***REMOVED******REMOVED***1, 2***REMOVED***,
		"a: 1\n",
	***REMOVED***,

	// Struct inlining
	***REMOVED***
		&struct ***REMOVED***
			A int
			C inlineB `yaml:",inline"`
		***REMOVED******REMOVED***1, inlineB***REMOVED***2, inlineC***REMOVED***3***REMOVED******REMOVED******REMOVED***,
		"a: 1\nb: 2\nc: 3\n",
	***REMOVED***,

	// Map inlining
	***REMOVED***
		&struct ***REMOVED***
			A int
			C map[string]int `yaml:",inline"`
		***REMOVED******REMOVED***1, map[string]int***REMOVED***"b": 2, "c": 3***REMOVED******REMOVED***,
		"a: 1\nb: 2\nc: 3\n",
	***REMOVED***,

	// Duration
	***REMOVED***
		map[string]time.Duration***REMOVED***"a": 3 * time.Second***REMOVED***,
		"a: 3s\n",
	***REMOVED***,

	// Issue #24: bug in map merging logic.
	***REMOVED***
		map[string]string***REMOVED***"a": "<foo>"***REMOVED***,
		"a: <foo>\n",
	***REMOVED***,

	// Issue #34: marshal unsupported base 60 floats quoted for compatibility
	// with old YAML 1.1 parsers.
	***REMOVED***
		map[string]string***REMOVED***"a": "1:1"***REMOVED***,
		"a: \"1:1\"\n",
	***REMOVED***,

	// Binary data.
	***REMOVED***
		map[string]string***REMOVED***"a": "\x00"***REMOVED***,
		"a: \"\\0\"\n",
	***REMOVED***, ***REMOVED***
		map[string]string***REMOVED***"a": "\x80\x81\x82"***REMOVED***,
		"a: !!binary gIGC\n",
	***REMOVED***, ***REMOVED***
		map[string]string***REMOVED***"a": strings.Repeat("\x90", 54)***REMOVED***,
		"a: !!binary |\n  " + strings.Repeat("kJCQ", 17) + "kJ\n  CQ\n",
	***REMOVED***,

	// Ordered maps.
	***REMOVED***
		&yaml.MapSlice***REMOVED******REMOVED***"b", 2***REMOVED***, ***REMOVED***"a", 1***REMOVED***, ***REMOVED***"d", 4***REMOVED***, ***REMOVED***"c", 3***REMOVED***, ***REMOVED***"sub", yaml.MapSlice***REMOVED******REMOVED***"e", 5***REMOVED******REMOVED******REMOVED******REMOVED***,
		"b: 2\na: 1\nd: 4\nc: 3\nsub:\n  e: 5\n",
	***REMOVED***,

	// Encode unicode as utf-8 rather than in escaped form.
	***REMOVED***
		map[string]string***REMOVED***"a": "你好"***REMOVED***,
		"a: 你好\n",
	***REMOVED***,

	// Support encoding.TextMarshaler.
	***REMOVED***
		map[string]net.IP***REMOVED***"a": net.IPv4(1, 2, 3, 4)***REMOVED***,
		"a: 1.2.3.4\n",
	***REMOVED***,
	***REMOVED***
		map[string]time.Time***REMOVED***"a": time.Unix(1424801979, 0)***REMOVED***,
		"a: 2015-02-24T18:19:39Z\n",
	***REMOVED***,

	// Ensure strings containing ": " are quoted (reported as PR #43, but not reproducible).
	***REMOVED***
		map[string]string***REMOVED***"a": "b: c"***REMOVED***,
		"a: 'b: c'\n",
	***REMOVED***,

	// Containing hash mark ('#') in string should be quoted
	***REMOVED***
		map[string]string***REMOVED***"a": "Hello #comment"***REMOVED***,
		"a: 'Hello #comment'\n",
	***REMOVED***,
	***REMOVED***
		map[string]string***REMOVED***"a": "你好 #comment"***REMOVED***,
		"a: '你好 #comment'\n",
	***REMOVED***,
***REMOVED***

func (s *S) TestMarshal(c *C) ***REMOVED***
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")
	for _, item := range marshalTests ***REMOVED***
		data, err := yaml.Marshal(item.value)
		c.Assert(err, IsNil)
		c.Assert(string(data), Equals, item.data)
	***REMOVED***
***REMOVED***

var marshalErrorTests = []struct ***REMOVED***
	value interface***REMOVED******REMOVED***
	error string
	panic string
***REMOVED******REMOVED******REMOVED***
	value: &struct ***REMOVED***
		B       int
		inlineB ",inline"
	***REMOVED******REMOVED***1, inlineB***REMOVED***2, inlineC***REMOVED***3***REMOVED******REMOVED******REMOVED***,
	panic: `Duplicated key 'b' in struct struct \***REMOVED*** B int; .*`,
***REMOVED***, ***REMOVED***
	value: &struct ***REMOVED***
		A int
		B map[string]int ",inline"
	***REMOVED******REMOVED***1, map[string]int***REMOVED***"a": 2***REMOVED******REMOVED***,
	panic: `Can't have key "a" in inlined map; conflicts with struct field`,
***REMOVED******REMOVED***

func (s *S) TestMarshalErrors(c *C) ***REMOVED***
	for _, item := range marshalErrorTests ***REMOVED***
		if item.panic != "" ***REMOVED***
			c.Assert(func() ***REMOVED*** yaml.Marshal(item.value) ***REMOVED***, PanicMatches, item.panic)
		***REMOVED*** else ***REMOVED***
			_, err := yaml.Marshal(item.value)
			c.Assert(err, ErrorMatches, item.error)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *S) TestMarshalTypeCache(c *C) ***REMOVED***
	var data []byte
	var err error
	func() ***REMOVED***
		type T struct***REMOVED*** A int ***REMOVED***
		data, err = yaml.Marshal(&T***REMOVED******REMOVED***)
		c.Assert(err, IsNil)
	***REMOVED***()
	func() ***REMOVED***
		type T struct***REMOVED*** B int ***REMOVED***
		data, err = yaml.Marshal(&T***REMOVED******REMOVED***)
		c.Assert(err, IsNil)
	***REMOVED***()
	c.Assert(string(data), Equals, "b: 0\n")
***REMOVED***

var marshalerTests = []struct ***REMOVED***
	data  string
	value interface***REMOVED******REMOVED***
***REMOVED******REMOVED***
	***REMOVED***"_:\n  hi: there\n", map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***"hi": "there"***REMOVED******REMOVED***,
	***REMOVED***"_:\n- 1\n- A\n", []interface***REMOVED******REMOVED******REMOVED***1, "A"***REMOVED******REMOVED***,
	***REMOVED***"_: 10\n", 10***REMOVED***,
	***REMOVED***"_: null\n", nil***REMOVED***,
	***REMOVED***"_: BAR!\n", "BAR!"***REMOVED***,
***REMOVED***

type marshalerType struct ***REMOVED***
	value interface***REMOVED******REMOVED***
***REMOVED***

func (o marshalerType) MarshalText() ([]byte, error) ***REMOVED***
	panic("MarshalText called on type with MarshalYAML")
***REMOVED***

func (o marshalerType) MarshalYAML() (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return o.value, nil
***REMOVED***

type marshalerValue struct ***REMOVED***
	Field marshalerType "_"
***REMOVED***

func (s *S) TestMarshaler(c *C) ***REMOVED***
	for _, item := range marshalerTests ***REMOVED***
		obj := &marshalerValue***REMOVED******REMOVED***
		obj.Field.value = item.value
		data, err := yaml.Marshal(obj)
		c.Assert(err, IsNil)
		c.Assert(string(data), Equals, string(item.data))
	***REMOVED***
***REMOVED***

func (s *S) TestMarshalerWholeDocument(c *C) ***REMOVED***
	obj := &marshalerType***REMOVED******REMOVED***
	obj.value = map[string]string***REMOVED***"hello": "world!"***REMOVED***
	data, err := yaml.Marshal(obj)
	c.Assert(err, IsNil)
	c.Assert(string(data), Equals, "hello: world!\n")
***REMOVED***

type failingMarshaler struct***REMOVED******REMOVED***

func (ft *failingMarshaler) MarshalYAML() (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return nil, failingErr
***REMOVED***

func (s *S) TestMarshalerError(c *C) ***REMOVED***
	_, err := yaml.Marshal(&failingMarshaler***REMOVED******REMOVED***)
	c.Assert(err, Equals, failingErr)
***REMOVED***

func (s *S) TestSortedOutput(c *C) ***REMOVED***
	order := []interface***REMOVED******REMOVED******REMOVED***
		false,
		true,
		1,
		uint(1),
		1.0,
		1.1,
		1.2,
		2,
		uint(2),
		2.0,
		2.1,
		"",
		".1",
		".2",
		".a",
		"1",
		"2",
		"a!10",
		"a/2",
		"a/10",
		"a~10",
		"ab/1",
		"b/1",
		"b/01",
		"b/2",
		"b/02",
		"b/3",
		"b/03",
		"b1",
		"b01",
		"b3",
		"c2.10",
		"c10.2",
		"d1",
		"d12",
		"d12a",
	***REMOVED***
	m := make(map[interface***REMOVED******REMOVED***]int)
	for _, k := range order ***REMOVED***
		m[k] = 1
	***REMOVED***
	data, err := yaml.Marshal(m)
	c.Assert(err, IsNil)
	out := "\n" + string(data)
	last := 0
	for i, k := range order ***REMOVED***
		repr := fmt.Sprint(k)
		if s, ok := k.(string); ok ***REMOVED***
			if _, err = strconv.ParseFloat(repr, 32); s == "" || err == nil ***REMOVED***
				repr = `"` + repr + `"`
			***REMOVED***
		***REMOVED***
		index := strings.Index(out, "\n"+repr+":")
		if index == -1 ***REMOVED***
			c.Fatalf("%#v is not in the output: %#v", k, out)
		***REMOVED***
		if index < last ***REMOVED***
			c.Fatalf("%#v was generated before %#v: %q", k, order[i-1], out)
		***REMOVED***
		last = index
	***REMOVED***
***REMOVED***
