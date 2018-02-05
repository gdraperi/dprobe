package mapstructure

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestComposeDecodeHookFunc(t *testing.T) ***REMOVED***
	f1 := func(
		f reflect.Kind,
		t reflect.Kind,
		data interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
		return data.(string) + "foo", nil
	***REMOVED***

	f2 := func(
		f reflect.Kind,
		t reflect.Kind,
		data interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
		return data.(string) + "bar", nil
	***REMOVED***

	f := ComposeDecodeHookFunc(f1, f2)

	result, err := DecodeHookExec(
		f, reflect.TypeOf(""), reflect.TypeOf([]byte("")), "")
	if err != nil ***REMOVED***
		t.Fatalf("bad: %s", err)
	***REMOVED***
	if result.(string) != "foobar" ***REMOVED***
		t.Fatalf("bad: %#v", result)
	***REMOVED***
***REMOVED***

func TestComposeDecodeHookFunc_err(t *testing.T) ***REMOVED***
	f1 := func(reflect.Kind, reflect.Kind, interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
		return nil, errors.New("foo")
	***REMOVED***

	f2 := func(reflect.Kind, reflect.Kind, interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
		panic("NOPE")
	***REMOVED***

	f := ComposeDecodeHookFunc(f1, f2)

	_, err := DecodeHookExec(
		f, reflect.TypeOf(""), reflect.TypeOf([]byte("")), 42)
	if err.Error() != "foo" ***REMOVED***
		t.Fatalf("bad: %s", err)
	***REMOVED***
***REMOVED***

func TestComposeDecodeHookFunc_kinds(t *testing.T) ***REMOVED***
	var f2From reflect.Kind

	f1 := func(
		f reflect.Kind,
		t reflect.Kind,
		data interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
		return int(42), nil
	***REMOVED***

	f2 := func(
		f reflect.Kind,
		t reflect.Kind,
		data interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
		f2From = f
		return data, nil
	***REMOVED***

	f := ComposeDecodeHookFunc(f1, f2)

	_, err := DecodeHookExec(
		f, reflect.TypeOf(""), reflect.TypeOf([]byte("")), "")
	if err != nil ***REMOVED***
		t.Fatalf("bad: %s", err)
	***REMOVED***
	if f2From != reflect.Int ***REMOVED***
		t.Fatalf("bad: %#v", f2From)
	***REMOVED***
***REMOVED***

func TestStringToSliceHookFunc(t *testing.T) ***REMOVED***
	f := StringToSliceHookFunc(",")

	strType := reflect.TypeOf("")
	sliceType := reflect.TypeOf([]byte(""))
	cases := []struct ***REMOVED***
		f, t   reflect.Type
		data   interface***REMOVED******REMOVED***
		result interface***REMOVED******REMOVED***
		err    bool
	***REMOVED******REMOVED***
		***REMOVED***sliceType, sliceType, 42, 42, false***REMOVED***,
		***REMOVED***strType, strType, 42, 42, false***REMOVED***,
		***REMOVED***
			strType,
			sliceType,
			"foo,bar,baz",
			[]string***REMOVED***"foo", "bar", "baz"***REMOVED***,
			false,
		***REMOVED***,
		***REMOVED***
			strType,
			sliceType,
			"",
			[]string***REMOVED******REMOVED***,
			false,
		***REMOVED***,
	***REMOVED***

	for i, tc := range cases ***REMOVED***
		actual, err := DecodeHookExec(f, tc.f, tc.t, tc.data)
		if tc.err != (err != nil) ***REMOVED***
			t.Fatalf("case %d: expected err %#v", i, tc.err)
		***REMOVED***
		if !reflect.DeepEqual(actual, tc.result) ***REMOVED***
			t.Fatalf(
				"case %d: expected %#v, got %#v",
				i, tc.result, actual)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestStringToTimeDurationHookFunc(t *testing.T) ***REMOVED***
	f := StringToTimeDurationHookFunc()

	strType := reflect.TypeOf("")
	timeType := reflect.TypeOf(time.Duration(5))
	cases := []struct ***REMOVED***
		f, t   reflect.Type
		data   interface***REMOVED******REMOVED***
		result interface***REMOVED******REMOVED***
		err    bool
	***REMOVED******REMOVED***
		***REMOVED***strType, timeType, "5s", 5 * time.Second, false***REMOVED***,
		***REMOVED***strType, timeType, "5", time.Duration(0), true***REMOVED***,
		***REMOVED***strType, strType, "5", "5", false***REMOVED***,
	***REMOVED***

	for i, tc := range cases ***REMOVED***
		actual, err := DecodeHookExec(f, tc.f, tc.t, tc.data)
		if tc.err != (err != nil) ***REMOVED***
			t.Fatalf("case %d: expected err %#v", i, tc.err)
		***REMOVED***
		if !reflect.DeepEqual(actual, tc.result) ***REMOVED***
			t.Fatalf(
				"case %d: expected %#v, got %#v",
				i, tc.result, actual)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestStringToTimeHookFunc(t *testing.T) ***REMOVED***
	strType := reflect.TypeOf("")
	timeType := reflect.TypeOf(time.Time***REMOVED******REMOVED***)
	cases := []struct ***REMOVED***
		f, t   reflect.Type
		layout string
		data   interface***REMOVED******REMOVED***
		result interface***REMOVED******REMOVED***
		err    bool
	***REMOVED******REMOVED***
		***REMOVED***strType, timeType, time.RFC3339, "2006-01-02T15:04:05Z",
			time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC), false***REMOVED***,
		***REMOVED***strType, timeType, time.RFC3339, "5", time.Time***REMOVED******REMOVED***, true***REMOVED***,
		***REMOVED***strType, strType, time.RFC3339, "5", "5", false***REMOVED***,
	***REMOVED***

	for i, tc := range cases ***REMOVED***
		f := StringToTimeHookFunc(tc.layout)
		actual, err := DecodeHookExec(f, tc.f, tc.t, tc.data)
		if tc.err != (err != nil) ***REMOVED***
			t.Fatalf("case %d: expected err %#v", i, tc.err)
		***REMOVED***
		if !reflect.DeepEqual(actual, tc.result) ***REMOVED***
			t.Fatalf(
				"case %d: expected %#v, got %#v",
				i, tc.result, actual)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestWeaklyTypedHook(t *testing.T) ***REMOVED***
	var f DecodeHookFunc = WeaklyTypedHook

	boolType := reflect.TypeOf(true)
	strType := reflect.TypeOf("")
	sliceType := reflect.TypeOf([]byte(""))
	cases := []struct ***REMOVED***
		f, t   reflect.Type
		data   interface***REMOVED******REMOVED***
		result interface***REMOVED******REMOVED***
		err    bool
	***REMOVED******REMOVED***
		// TO STRING
		***REMOVED***
			boolType,
			strType,
			false,
			"0",
			false,
		***REMOVED***,

		***REMOVED***
			boolType,
			strType,
			true,
			"1",
			false,
		***REMOVED***,

		***REMOVED***
			reflect.TypeOf(float32(1)),
			strType,
			float32(7),
			"7",
			false,
		***REMOVED***,

		***REMOVED***
			reflect.TypeOf(int(1)),
			strType,
			int(7),
			"7",
			false,
		***REMOVED***,

		***REMOVED***
			sliceType,
			strType,
			[]uint8("foo"),
			"foo",
			false,
		***REMOVED***,

		***REMOVED***
			reflect.TypeOf(uint(1)),
			strType,
			uint(7),
			"7",
			false,
		***REMOVED***,
	***REMOVED***

	for i, tc := range cases ***REMOVED***
		actual, err := DecodeHookExec(f, tc.f, tc.t, tc.data)
		if tc.err != (err != nil) ***REMOVED***
			t.Fatalf("case %d: expected err %#v", i, tc.err)
		***REMOVED***
		if !reflect.DeepEqual(actual, tc.result) ***REMOVED***
			t.Fatalf(
				"case %d: expected %#v, got %#v",
				i, tc.result, actual)
		***REMOVED***
	***REMOVED***
***REMOVED***
