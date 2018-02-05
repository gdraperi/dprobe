package pflag

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func setUpUISFlagSet(uisp *[]uint) *FlagSet ***REMOVED***
	f := NewFlagSet("test", ContinueOnError)
	f.UintSliceVar(uisp, "uis", []uint***REMOVED******REMOVED***, "Command separated list!")
	return f
***REMOVED***

func setUpUISFlagSetWithDefault(uisp *[]uint) *FlagSet ***REMOVED***
	f := NewFlagSet("test", ContinueOnError)
	f.UintSliceVar(uisp, "uis", []uint***REMOVED***0, 1***REMOVED***, "Command separated list!")
	return f
***REMOVED***

func TestEmptyUIS(t *testing.T) ***REMOVED***
	var uis []uint
	f := setUpUISFlagSet(&uis)
	err := f.Parse([]string***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	getUIS, err := f.GetUintSlice("uis")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetUintSlice():", err)
	***REMOVED***
	if len(getUIS) != 0 ***REMOVED***
		t.Fatalf("got is %v with len=%d but expected length=0", getUIS, len(getUIS))
	***REMOVED***
***REMOVED***

func TestUIS(t *testing.T) ***REMOVED***
	var uis []uint
	f := setUpUISFlagSet(&uis)

	vals := []string***REMOVED***"1", "2", "4", "3"***REMOVED***
	arg := fmt.Sprintf("--uis=%s", strings.Join(vals, ","))
	err := f.Parse([]string***REMOVED***arg***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	for i, v := range uis ***REMOVED***
		u, err := strconv.ParseUint(vals[i], 10, 0)
		if err != nil ***REMOVED***
			t.Fatalf("got error: %v", err)
		***REMOVED***
		if uint(u) != v ***REMOVED***
			t.Fatalf("expected uis[%d] to be %s but got %d", i, vals[i], v)
		***REMOVED***
	***REMOVED***
	getUIS, err := f.GetUintSlice("uis")
	if err != nil ***REMOVED***
		t.Fatalf("got error: %v", err)
	***REMOVED***
	for i, v := range getUIS ***REMOVED***
		u, err := strconv.ParseUint(vals[i], 10, 0)
		if err != nil ***REMOVED***
			t.Fatalf("got error: %v", err)
		***REMOVED***
		if uint(u) != v ***REMOVED***
			t.Fatalf("expected uis[%d] to be %s but got: %d from GetUintSlice", i, vals[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestUISDefault(t *testing.T) ***REMOVED***
	var uis []uint
	f := setUpUISFlagSetWithDefault(&uis)

	vals := []string***REMOVED***"0", "1"***REMOVED***

	err := f.Parse([]string***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	for i, v := range uis ***REMOVED***
		u, err := strconv.ParseUint(vals[i], 10, 0)
		if err != nil ***REMOVED***
			t.Fatalf("got error: %v", err)
		***REMOVED***
		if uint(u) != v ***REMOVED***
			t.Fatalf("expect uis[%d] to be %d but got: %d", i, u, v)
		***REMOVED***
	***REMOVED***

	getUIS, err := f.GetUintSlice("uis")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetUintSlice():", err)
	***REMOVED***
	for i, v := range getUIS ***REMOVED***
		u, err := strconv.ParseUint(vals[i], 10, 0)
		if err != nil ***REMOVED***
			t.Fatal("got an error from GetIntSlice():", err)
		***REMOVED***
		if uint(u) != v ***REMOVED***
			t.Fatalf("expected uis[%d] to be %d from GetUintSlice but got: %d", i, u, v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestUISWithDefault(t *testing.T) ***REMOVED***
	var uis []uint
	f := setUpUISFlagSetWithDefault(&uis)

	vals := []string***REMOVED***"1", "2"***REMOVED***
	arg := fmt.Sprintf("--uis=%s", strings.Join(vals, ","))
	err := f.Parse([]string***REMOVED***arg***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	for i, v := range uis ***REMOVED***
		u, err := strconv.ParseUint(vals[i], 10, 0)
		if err != nil ***REMOVED***
			t.Fatalf("got error: %v", err)
		***REMOVED***
		if uint(u) != v ***REMOVED***
			t.Fatalf("expected uis[%d] to be %d from GetUintSlice but got: %d", i, u, v)
		***REMOVED***
	***REMOVED***

	getUIS, err := f.GetUintSlice("uis")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetUintSlice():", err)
	***REMOVED***
	for i, v := range getUIS ***REMOVED***
		u, err := strconv.ParseUint(vals[i], 10, 0)
		if err != nil ***REMOVED***
			t.Fatalf("got error: %v", err)
		***REMOVED***
		if uint(u) != v ***REMOVED***
			t.Fatalf("expected uis[%d] to be %d from GetUintSlice but got: %d", i, u, v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestUISCalledTwice(t *testing.T) ***REMOVED***
	var uis []uint
	f := setUpUISFlagSet(&uis)

	in := []string***REMOVED***"1,2", "3"***REMOVED***
	expected := []int***REMOVED***1, 2, 3***REMOVED***
	argfmt := "--uis=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string***REMOVED***arg1, arg2***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	for i, v := range uis ***REMOVED***
		if uint(expected[i]) != v ***REMOVED***
			t.Fatalf("expected uis[%d] to be %d but got: %d", i, expected[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***
