package toml

import (
	"strconv"
	"testing"
	"time"
)

type customString string

type stringer struct***REMOVED******REMOVED***

func (s stringer) String() string ***REMOVED***
	return "stringer"
***REMOVED***

func validate(t *testing.T, path string, object interface***REMOVED******REMOVED***) ***REMOVED***
	switch o := object.(type) ***REMOVED***
	case *Tree:
		for key, tree := range o.values ***REMOVED***
			validate(t, path+"."+key, tree)
		***REMOVED***
	case []*Tree:
		for index, tree := range o ***REMOVED***
			validate(t, path+"."+strconv.Itoa(index), tree)
		***REMOVED***
	case *tomlValue:
		switch o.value.(type) ***REMOVED***
		case int64, uint64, bool, string, float64, time.Time,
			[]int64, []uint64, []bool, []string, []float64, []time.Time:
		default:
			t.Fatalf("tomlValue at key %s containing incorrect type %T", path, o.value)
		***REMOVED***
	default:
		t.Fatalf("value at key %s is of incorrect type %T", path, object)
	***REMOVED***
	t.Logf("validation ok %s as %T", path, object)
***REMOVED***

func validateTree(t *testing.T, tree *Tree) ***REMOVED***
	validate(t, "", tree)
***REMOVED***

func TestTreeCreateToTree(t *testing.T) ***REMOVED***
	data := map[string]interface***REMOVED******REMOVED******REMOVED***
		"a_string": "bar",
		"an_int":   42,
		"time":     time.Now(),
		"int8":     int8(2),
		"int16":    int16(2),
		"int32":    int32(2),
		"uint8":    uint8(2),
		"uint16":   uint16(2),
		"uint32":   uint32(2),
		"float32":  float32(2),
		"a_bool":   false,
		"stringer": stringer***REMOVED******REMOVED***,
		"nested": map[string]interface***REMOVED******REMOVED******REMOVED***
			"foo": "bar",
		***REMOVED***,
		"array":                 []string***REMOVED***"a", "b", "c"***REMOVED***,
		"array_uint":            []uint***REMOVED***uint(1), uint(2)***REMOVED***,
		"array_table":           []map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***"sub_map": 52***REMOVED******REMOVED***,
		"array_times":           []time.Time***REMOVED***time.Now(), time.Now()***REMOVED***,
		"map_times":             map[string]time.Time***REMOVED***"now": time.Now()***REMOVED***,
		"custom_string_map_key": map[customString]interface***REMOVED******REMOVED******REMOVED***customString("custom"): "custom"***REMOVED***,
	***REMOVED***
	tree, err := TreeFromMap(data)
	if err != nil ***REMOVED***
		t.Fatal("unexpected error:", err)
	***REMOVED***
	validateTree(t, tree)
***REMOVED***

func TestTreeCreateToTreeInvalidLeafType(t *testing.T) ***REMOVED***
	_, err := TreeFromMap(map[string]interface***REMOVED******REMOVED******REMOVED***"foo": t***REMOVED***)
	expected := "cannot convert type *testing.T to Tree"
	if err.Error() != expected ***REMOVED***
		t.Fatalf("expected error %s, got %s", expected, err.Error())
	***REMOVED***
***REMOVED***

func TestTreeCreateToTreeInvalidMapKeyType(t *testing.T) ***REMOVED***
	_, err := TreeFromMap(map[string]interface***REMOVED******REMOVED******REMOVED***"foo": map[int]interface***REMOVED******REMOVED******REMOVED***2: 1***REMOVED******REMOVED***)
	expected := "map key needs to be a string, not int (int)"
	if err.Error() != expected ***REMOVED***
		t.Fatalf("expected error %s, got %s", expected, err.Error())
	***REMOVED***
***REMOVED***

func TestTreeCreateToTreeInvalidArrayMemberType(t *testing.T) ***REMOVED***
	_, err := TreeFromMap(map[string]interface***REMOVED******REMOVED******REMOVED***"foo": []*testing.T***REMOVED***t***REMOVED******REMOVED***)
	expected := "cannot convert type *testing.T to Tree"
	if err.Error() != expected ***REMOVED***
		t.Fatalf("expected error %s, got %s", expected, err.Error())
	***REMOVED***
***REMOVED***

func TestTreeCreateToTreeInvalidTableGroupType(t *testing.T) ***REMOVED***
	_, err := TreeFromMap(map[string]interface***REMOVED******REMOVED******REMOVED***"foo": []map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***"hello": t***REMOVED******REMOVED******REMOVED***)
	expected := "cannot convert type *testing.T to Tree"
	if err.Error() != expected ***REMOVED***
		t.Fatalf("expected error %s, got %s", expected, err.Error())
	***REMOVED***
***REMOVED***

func TestRoundTripArrayOfTables(t *testing.T) ***REMOVED***
	orig := "\n[[stuff]]\n  name = \"foo\"\n  things = [\"a\",\"b\"]\n"
	tree, err := Load(orig)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error: %s", err)
	***REMOVED***

	m := tree.ToMap()

	tree, err = TreeFromMap(m)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error: %s", err)
	***REMOVED***
	want := orig
	got := tree.String()

	if got != want ***REMOVED***
		t.Errorf("want:\n%s\ngot:\n%s", want, got)
	***REMOVED***
***REMOVED***
