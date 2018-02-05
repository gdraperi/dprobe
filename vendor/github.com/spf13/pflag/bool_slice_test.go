package pflag

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func setUpBSFlagSet(bsp *[]bool) *FlagSet ***REMOVED***
	f := NewFlagSet("test", ContinueOnError)
	f.BoolSliceVar(bsp, "bs", []bool***REMOVED******REMOVED***, "Command separated list!")
	return f
***REMOVED***

func setUpBSFlagSetWithDefault(bsp *[]bool) *FlagSet ***REMOVED***
	f := NewFlagSet("test", ContinueOnError)
	f.BoolSliceVar(bsp, "bs", []bool***REMOVED***false, true***REMOVED***, "Command separated list!")
	return f
***REMOVED***

func TestEmptyBS(t *testing.T) ***REMOVED***
	var bs []bool
	f := setUpBSFlagSet(&bs)
	err := f.Parse([]string***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	getBS, err := f.GetBoolSlice("bs")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetBoolSlice():", err)
	***REMOVED***
	if len(getBS) != 0 ***REMOVED***
		t.Fatalf("got bs %v with len=%d but expected length=0", getBS, len(getBS))
	***REMOVED***
***REMOVED***

func TestBS(t *testing.T) ***REMOVED***
	var bs []bool
	f := setUpBSFlagSet(&bs)

	vals := []string***REMOVED***"1", "F", "TRUE", "0"***REMOVED***
	arg := fmt.Sprintf("--bs=%s", strings.Join(vals, ","))
	err := f.Parse([]string***REMOVED***arg***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	for i, v := range bs ***REMOVED***
		b, err := strconv.ParseBool(vals[i])
		if err != nil ***REMOVED***
			t.Fatalf("got error: %v", err)
		***REMOVED***
		if b != v ***REMOVED***
			t.Fatalf("expected is[%d] to be %s but got: %t", i, vals[i], v)
		***REMOVED***
	***REMOVED***
	getBS, err := f.GetBoolSlice("bs")
	if err != nil ***REMOVED***
		t.Fatalf("got error: %v", err)
	***REMOVED***
	for i, v := range getBS ***REMOVED***
		b, err := strconv.ParseBool(vals[i])
		if err != nil ***REMOVED***
			t.Fatalf("got error: %v", err)
		***REMOVED***
		if b != v ***REMOVED***
			t.Fatalf("expected bs[%d] to be %s but got: %t from GetBoolSlice", i, vals[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestBSDefault(t *testing.T) ***REMOVED***
	var bs []bool
	f := setUpBSFlagSetWithDefault(&bs)

	vals := []string***REMOVED***"false", "T"***REMOVED***

	err := f.Parse([]string***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	for i, v := range bs ***REMOVED***
		b, err := strconv.ParseBool(vals[i])
		if err != nil ***REMOVED***
			t.Fatalf("got error: %v", err)
		***REMOVED***
		if b != v ***REMOVED***
			t.Fatalf("expected bs[%d] to be %t from GetBoolSlice but got: %t", i, b, v)
		***REMOVED***
	***REMOVED***

	getBS, err := f.GetBoolSlice("bs")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetBoolSlice():", err)
	***REMOVED***
	for i, v := range getBS ***REMOVED***
		b, err := strconv.ParseBool(vals[i])
		if err != nil ***REMOVED***
			t.Fatal("got an error from GetBoolSlice():", err)
		***REMOVED***
		if b != v ***REMOVED***
			t.Fatalf("expected bs[%d] to be %t from GetBoolSlice but got: %t", i, b, v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestBSWithDefault(t *testing.T) ***REMOVED***
	var bs []bool
	f := setUpBSFlagSetWithDefault(&bs)

	vals := []string***REMOVED***"FALSE", "1"***REMOVED***
	arg := fmt.Sprintf("--bs=%s", strings.Join(vals, ","))
	err := f.Parse([]string***REMOVED***arg***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	for i, v := range bs ***REMOVED***
		b, err := strconv.ParseBool(vals[i])
		if err != nil ***REMOVED***
			t.Fatalf("got error: %v", err)
		***REMOVED***
		if b != v ***REMOVED***
			t.Fatalf("expected bs[%d] to be %t but got: %t", i, b, v)
		***REMOVED***
	***REMOVED***

	getBS, err := f.GetBoolSlice("bs")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetBoolSlice():", err)
	***REMOVED***
	for i, v := range getBS ***REMOVED***
		b, err := strconv.ParseBool(vals[i])
		if err != nil ***REMOVED***
			t.Fatalf("got error: %v", err)
		***REMOVED***
		if b != v ***REMOVED***
			t.Fatalf("expected bs[%d] to be %t from GetBoolSlice but got: %t", i, b, v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestBSCalledTwice(t *testing.T) ***REMOVED***
	var bs []bool
	f := setUpBSFlagSet(&bs)

	in := []string***REMOVED***"T,F", "T"***REMOVED***
	expected := []bool***REMOVED***true, false, true***REMOVED***
	argfmt := "--bs=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string***REMOVED***arg1, arg2***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	for i, v := range bs ***REMOVED***
		if expected[i] != v ***REMOVED***
			t.Fatalf("expected bs[%d] to be %t but got %t", i, expected[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestBSBadQuoting(t *testing.T) ***REMOVED***

	tests := []struct ***REMOVED***
		Want    []bool
		FlagArg []string
	***REMOVED******REMOVED***
		***REMOVED***
			Want:    []bool***REMOVED***true, false, true***REMOVED***,
			FlagArg: []string***REMOVED***"1", "0", "true"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Want:    []bool***REMOVED***true, false***REMOVED***,
			FlagArg: []string***REMOVED***"True", "F"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Want:    []bool***REMOVED***true, false***REMOVED***,
			FlagArg: []string***REMOVED***"T", "0"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Want:    []bool***REMOVED***true, false***REMOVED***,
			FlagArg: []string***REMOVED***"1", "0"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Want:    []bool***REMOVED***true, false, false***REMOVED***,
			FlagArg: []string***REMOVED***"true,false", "false"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Want:    []bool***REMOVED***true, false, false, true, false, true, false***REMOVED***,
			FlagArg: []string***REMOVED***`"true,false,false,1,0,     T"`, " false "***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Want:    []bool***REMOVED***false, false, true, false, true, false, true***REMOVED***,
			FlagArg: []string***REMOVED***`"0, False,  T,false  , true,F"`, "true"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***

		var bs []bool
		f := setUpBSFlagSet(&bs)

		if err := f.Parse([]string***REMOVED***fmt.Sprintf("--bs=%s", strings.Join(test.FlagArg, ","))***REMOVED***); err != nil ***REMOVED***
			t.Fatalf("flag parsing failed with error: %s\nparsing:\t%#v\nwant:\t\t%#v",
				err, test.FlagArg, test.Want[i])
		***REMOVED***

		for j, b := range bs ***REMOVED***
			if b != test.Want[j] ***REMOVED***
				t.Fatalf("bad value parsed for test %d on bool %d:\nwant:\t%t\ngot:\t%t", i, j, test.Want[j], b)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
