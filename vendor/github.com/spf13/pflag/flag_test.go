// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	testBool                     = Bool("test_bool", false, "bool value")
	testInt                      = Int("test_int", 0, "int value")
	testInt64                    = Int64("test_int64", 0, "int64 value")
	testUint                     = Uint("test_uint", 0, "uint value")
	testUint64                   = Uint64("test_uint64", 0, "uint64 value")
	testString                   = String("test_string", "0", "string value")
	testFloat                    = Float64("test_float64", 0, "float64 value")
	testDuration                 = Duration("test_duration", 0, "time.Duration value")
	testOptionalInt              = Int("test_optional_int", 0, "optional int value")
	normalizeFlagNameInvocations = 0
)

func boolString(s string) string ***REMOVED***
	if s == "0" ***REMOVED***
		return "false"
	***REMOVED***
	return "true"
***REMOVED***

func TestEverything(t *testing.T) ***REMOVED***
	m := make(map[string]*Flag)
	desired := "0"
	visitor := func(f *Flag) ***REMOVED***
		if len(f.Name) > 5 && f.Name[0:5] == "test_" ***REMOVED***
			m[f.Name] = f
			ok := false
			switch ***REMOVED***
			case f.Value.String() == desired:
				ok = true
			case f.Name == "test_bool" && f.Value.String() == boolString(desired):
				ok = true
			case f.Name == "test_duration" && f.Value.String() == desired+"s":
				ok = true
			***REMOVED***
			if !ok ***REMOVED***
				t.Error("Visit: bad value", f.Value.String(), "for", f.Name)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	VisitAll(visitor)
	if len(m) != 9 ***REMOVED***
		t.Error("VisitAll misses some flags")
		for k, v := range m ***REMOVED***
			t.Log(k, *v)
		***REMOVED***
	***REMOVED***
	m = make(map[string]*Flag)
	Visit(visitor)
	if len(m) != 0 ***REMOVED***
		t.Errorf("Visit sees unset flags")
		for k, v := range m ***REMOVED***
			t.Log(k, *v)
		***REMOVED***
	***REMOVED***
	// Now set all flags
	Set("test_bool", "true")
	Set("test_int", "1")
	Set("test_int64", "1")
	Set("test_uint", "1")
	Set("test_uint64", "1")
	Set("test_string", "1")
	Set("test_float64", "1")
	Set("test_duration", "1s")
	Set("test_optional_int", "1")
	desired = "1"
	Visit(visitor)
	if len(m) != 9 ***REMOVED***
		t.Error("Visit fails after set")
		for k, v := range m ***REMOVED***
			t.Log(k, *v)
		***REMOVED***
	***REMOVED***
	// Now test they're visited in sort order.
	var flagNames []string
	Visit(func(f *Flag) ***REMOVED*** flagNames = append(flagNames, f.Name) ***REMOVED***)
	if !sort.StringsAreSorted(flagNames) ***REMOVED***
		t.Errorf("flag names not sorted: %v", flagNames)
	***REMOVED***
***REMOVED***

func TestUsage(t *testing.T) ***REMOVED***
	called := false
	ResetForTesting(func() ***REMOVED*** called = true ***REMOVED***)
	if GetCommandLine().Parse([]string***REMOVED***"--x"***REMOVED***) == nil ***REMOVED***
		t.Error("parse did not fail for unknown flag")
	***REMOVED***
	if called ***REMOVED***
		t.Error("did call Usage while using ContinueOnError")
	***REMOVED***
***REMOVED***

func TestAddFlagSet(t *testing.T) ***REMOVED***
	oldSet := NewFlagSet("old", ContinueOnError)
	newSet := NewFlagSet("new", ContinueOnError)

	oldSet.String("flag1", "flag1", "flag1")
	oldSet.String("flag2", "flag2", "flag2")

	newSet.String("flag2", "flag2", "flag2")
	newSet.String("flag3", "flag3", "flag3")

	oldSet.AddFlagSet(newSet)

	if len(oldSet.formal) != 3 ***REMOVED***
		t.Errorf("Unexpected result adding a FlagSet to a FlagSet %v", oldSet)
	***REMOVED***
***REMOVED***

func TestAnnotation(t *testing.T) ***REMOVED***
	f := NewFlagSet("shorthand", ContinueOnError)

	if err := f.SetAnnotation("missing-flag", "key", nil); err == nil ***REMOVED***
		t.Errorf("Expected error setting annotation on non-existent flag")
	***REMOVED***

	f.StringP("stringa", "a", "", "string value")
	if err := f.SetAnnotation("stringa", "key", nil); err != nil ***REMOVED***
		t.Errorf("Unexpected error setting new nil annotation: %v", err)
	***REMOVED***
	if annotation := f.Lookup("stringa").Annotations["key"]; annotation != nil ***REMOVED***
		t.Errorf("Unexpected annotation: %v", annotation)
	***REMOVED***

	f.StringP("stringb", "b", "", "string2 value")
	if err := f.SetAnnotation("stringb", "key", []string***REMOVED***"value1"***REMOVED***); err != nil ***REMOVED***
		t.Errorf("Unexpected error setting new annotation: %v", err)
	***REMOVED***
	if annotation := f.Lookup("stringb").Annotations["key"]; !reflect.DeepEqual(annotation, []string***REMOVED***"value1"***REMOVED***) ***REMOVED***
		t.Errorf("Unexpected annotation: %v", annotation)
	***REMOVED***

	if err := f.SetAnnotation("stringb", "key", []string***REMOVED***"value2"***REMOVED***); err != nil ***REMOVED***
		t.Errorf("Unexpected error updating annotation: %v", err)
	***REMOVED***
	if annotation := f.Lookup("stringb").Annotations["key"]; !reflect.DeepEqual(annotation, []string***REMOVED***"value2"***REMOVED***) ***REMOVED***
		t.Errorf("Unexpected annotation: %v", annotation)
	***REMOVED***
***REMOVED***

func testParse(f *FlagSet, t *testing.T) ***REMOVED***
	if f.Parsed() ***REMOVED***
		t.Error("f.Parse() = true before Parse")
	***REMOVED***
	boolFlag := f.Bool("bool", false, "bool value")
	bool2Flag := f.Bool("bool2", false, "bool2 value")
	bool3Flag := f.Bool("bool3", false, "bool3 value")
	intFlag := f.Int("int", 0, "int value")
	int8Flag := f.Int8("int8", 0, "int value")
	int16Flag := f.Int16("int16", 0, "int value")
	int32Flag := f.Int32("int32", 0, "int value")
	int64Flag := f.Int64("int64", 0, "int64 value")
	uintFlag := f.Uint("uint", 0, "uint value")
	uint8Flag := f.Uint8("uint8", 0, "uint value")
	uint16Flag := f.Uint16("uint16", 0, "uint value")
	uint32Flag := f.Uint32("uint32", 0, "uint value")
	uint64Flag := f.Uint64("uint64", 0, "uint64 value")
	stringFlag := f.String("string", "0", "string value")
	float32Flag := f.Float32("float32", 0, "float32 value")
	float64Flag := f.Float64("float64", 0, "float64 value")
	ipFlag := f.IP("ip", net.ParseIP("127.0.0.1"), "ip value")
	maskFlag := f.IPMask("mask", ParseIPv4Mask("0.0.0.0"), "mask value")
	durationFlag := f.Duration("duration", 5*time.Second, "time.Duration value")
	optionalIntNoValueFlag := f.Int("optional-int-no-value", 0, "int value")
	f.Lookup("optional-int-no-value").NoOptDefVal = "9"
	optionalIntWithValueFlag := f.Int("optional-int-with-value", 0, "int value")
	f.Lookup("optional-int-no-value").NoOptDefVal = "9"
	extra := "one-extra-argument"
	args := []string***REMOVED***
		"--bool",
		"--bool2=true",
		"--bool3=false",
		"--int=22",
		"--int8=-8",
		"--int16=-16",
		"--int32=-32",
		"--int64=0x23",
		"--uint", "24",
		"--uint8=8",
		"--uint16=16",
		"--uint32=32",
		"--uint64=25",
		"--string=hello",
		"--float32=-172e12",
		"--float64=2718e28",
		"--ip=10.11.12.13",
		"--mask=255.255.255.0",
		"--duration=2m",
		"--optional-int-no-value",
		"--optional-int-with-value=42",
		extra,
	***REMOVED***
	if err := f.Parse(args); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !f.Parsed() ***REMOVED***
		t.Error("f.Parse() = false after Parse")
	***REMOVED***
	if *boolFlag != true ***REMOVED***
		t.Error("bool flag should be true, is ", *boolFlag)
	***REMOVED***
	if v, err := f.GetBool("bool"); err != nil || v != *boolFlag ***REMOVED***
		t.Error("GetBool does not work.")
	***REMOVED***
	if *bool2Flag != true ***REMOVED***
		t.Error("bool2 flag should be true, is ", *bool2Flag)
	***REMOVED***
	if *bool3Flag != false ***REMOVED***
		t.Error("bool3 flag should be false, is ", *bool2Flag)
	***REMOVED***
	if *intFlag != 22 ***REMOVED***
		t.Error("int flag should be 22, is ", *intFlag)
	***REMOVED***
	if v, err := f.GetInt("int"); err != nil || v != *intFlag ***REMOVED***
		t.Error("GetInt does not work.")
	***REMOVED***
	if *int8Flag != -8 ***REMOVED***
		t.Error("int8 flag should be 0x23, is ", *int8Flag)
	***REMOVED***
	if *int16Flag != -16 ***REMOVED***
		t.Error("int16 flag should be -16, is ", *int16Flag)
	***REMOVED***
	if v, err := f.GetInt8("int8"); err != nil || v != *int8Flag ***REMOVED***
		t.Error("GetInt8 does not work.")
	***REMOVED***
	if v, err := f.GetInt16("int16"); err != nil || v != *int16Flag ***REMOVED***
		t.Error("GetInt16 does not work.")
	***REMOVED***
	if *int32Flag != -32 ***REMOVED***
		t.Error("int32 flag should be 0x23, is ", *int32Flag)
	***REMOVED***
	if v, err := f.GetInt32("int32"); err != nil || v != *int32Flag ***REMOVED***
		t.Error("GetInt32 does not work.")
	***REMOVED***
	if *int64Flag != 0x23 ***REMOVED***
		t.Error("int64 flag should be 0x23, is ", *int64Flag)
	***REMOVED***
	if v, err := f.GetInt64("int64"); err != nil || v != *int64Flag ***REMOVED***
		t.Error("GetInt64 does not work.")
	***REMOVED***
	if *uintFlag != 24 ***REMOVED***
		t.Error("uint flag should be 24, is ", *uintFlag)
	***REMOVED***
	if v, err := f.GetUint("uint"); err != nil || v != *uintFlag ***REMOVED***
		t.Error("GetUint does not work.")
	***REMOVED***
	if *uint8Flag != 8 ***REMOVED***
		t.Error("uint8 flag should be 8, is ", *uint8Flag)
	***REMOVED***
	if v, err := f.GetUint8("uint8"); err != nil || v != *uint8Flag ***REMOVED***
		t.Error("GetUint8 does not work.")
	***REMOVED***
	if *uint16Flag != 16 ***REMOVED***
		t.Error("uint16 flag should be 16, is ", *uint16Flag)
	***REMOVED***
	if v, err := f.GetUint16("uint16"); err != nil || v != *uint16Flag ***REMOVED***
		t.Error("GetUint16 does not work.")
	***REMOVED***
	if *uint32Flag != 32 ***REMOVED***
		t.Error("uint32 flag should be 32, is ", *uint32Flag)
	***REMOVED***
	if v, err := f.GetUint32("uint32"); err != nil || v != *uint32Flag ***REMOVED***
		t.Error("GetUint32 does not work.")
	***REMOVED***
	if *uint64Flag != 25 ***REMOVED***
		t.Error("uint64 flag should be 25, is ", *uint64Flag)
	***REMOVED***
	if v, err := f.GetUint64("uint64"); err != nil || v != *uint64Flag ***REMOVED***
		t.Error("GetUint64 does not work.")
	***REMOVED***
	if *stringFlag != "hello" ***REMOVED***
		t.Error("string flag should be `hello`, is ", *stringFlag)
	***REMOVED***
	if v, err := f.GetString("string"); err != nil || v != *stringFlag ***REMOVED***
		t.Error("GetString does not work.")
	***REMOVED***
	if *float32Flag != -172e12 ***REMOVED***
		t.Error("float32 flag should be -172e12, is ", *float32Flag)
	***REMOVED***
	if v, err := f.GetFloat32("float32"); err != nil || v != *float32Flag ***REMOVED***
		t.Errorf("GetFloat32 returned %v but float32Flag was %v", v, *float32Flag)
	***REMOVED***
	if *float64Flag != 2718e28 ***REMOVED***
		t.Error("float64 flag should be 2718e28, is ", *float64Flag)
	***REMOVED***
	if v, err := f.GetFloat64("float64"); err != nil || v != *float64Flag ***REMOVED***
		t.Errorf("GetFloat64 returned %v but float64Flag was %v", v, *float64Flag)
	***REMOVED***
	if !(*ipFlag).Equal(net.ParseIP("10.11.12.13")) ***REMOVED***
		t.Error("ip flag should be 10.11.12.13, is ", *ipFlag)
	***REMOVED***
	if v, err := f.GetIP("ip"); err != nil || !v.Equal(*ipFlag) ***REMOVED***
		t.Errorf("GetIP returned %v but ipFlag was %v", v, *ipFlag)
	***REMOVED***
	if (*maskFlag).String() != ParseIPv4Mask("255.255.255.0").String() ***REMOVED***
		t.Error("mask flag should be 255.255.255.0, is ", (*maskFlag).String())
	***REMOVED***
	if v, err := f.GetIPv4Mask("mask"); err != nil || v.String() != (*maskFlag).String() ***REMOVED***
		t.Errorf("GetIP returned %v maskFlag was %v error was %v", v, *maskFlag, err)
	***REMOVED***
	if *durationFlag != 2*time.Minute ***REMOVED***
		t.Error("duration flag should be 2m, is ", *durationFlag)
	***REMOVED***
	if v, err := f.GetDuration("duration"); err != nil || v != *durationFlag ***REMOVED***
		t.Error("GetDuration does not work.")
	***REMOVED***
	if _, err := f.GetInt("duration"); err == nil ***REMOVED***
		t.Error("GetInt parsed a time.Duration?!?!")
	***REMOVED***
	if *optionalIntNoValueFlag != 9 ***REMOVED***
		t.Error("optional int flag should be the default value, is ", *optionalIntNoValueFlag)
	***REMOVED***
	if *optionalIntWithValueFlag != 42 ***REMOVED***
		t.Error("optional int flag should be 42, is ", *optionalIntWithValueFlag)
	***REMOVED***
	if len(f.Args()) != 1 ***REMOVED***
		t.Error("expected one argument, got", len(f.Args()))
	***REMOVED*** else if f.Args()[0] != extra ***REMOVED***
		t.Errorf("expected argument %q got %q", extra, f.Args()[0])
	***REMOVED***
***REMOVED***

func testParseAll(f *FlagSet, t *testing.T) ***REMOVED***
	if f.Parsed() ***REMOVED***
		t.Error("f.Parse() = true before Parse")
	***REMOVED***
	f.BoolP("boola", "a", false, "bool value")
	f.BoolP("boolb", "b", false, "bool2 value")
	f.BoolP("boolc", "c", false, "bool3 value")
	f.BoolP("boold", "d", false, "bool4 value")
	f.StringP("stringa", "s", "0", "string value")
	f.StringP("stringz", "z", "0", "string value")
	f.StringP("stringx", "x", "0", "string value")
	f.StringP("stringy", "y", "0", "string value")
	f.Lookup("stringx").NoOptDefVal = "1"
	args := []string***REMOVED***
		"-ab",
		"-cs=xx",
		"--stringz=something",
		"-d=true",
		"-x",
		"-y",
		"ee",
	***REMOVED***
	want := []string***REMOVED***
		"boola", "true",
		"boolb", "true",
		"boolc", "true",
		"stringa", "xx",
		"stringz", "something",
		"boold", "true",
		"stringx", "1",
		"stringy", "ee",
	***REMOVED***
	got := []string***REMOVED******REMOVED***
	store := func(flag *Flag, value string) error ***REMOVED***
		got = append(got, flag.Name)
		if len(value) > 0 ***REMOVED***
			got = append(got, value)
		***REMOVED***
		return nil
	***REMOVED***
	if err := f.ParseAll(args, store); err != nil ***REMOVED***
		t.Errorf("expected no error, got %s", err)
	***REMOVED***
	if !f.Parsed() ***REMOVED***
		t.Errorf("f.Parse() = false after Parse")
	***REMOVED***
	if !reflect.DeepEqual(got, want) ***REMOVED***
		t.Errorf("f.ParseAll() fail to restore the args")
		t.Errorf("Got: %v", got)
		t.Errorf("Want: %v", want)
	***REMOVED***
***REMOVED***

func TestShorthand(t *testing.T) ***REMOVED***
	f := NewFlagSet("shorthand", ContinueOnError)
	if f.Parsed() ***REMOVED***
		t.Error("f.Parse() = true before Parse")
	***REMOVED***
	boolaFlag := f.BoolP("boola", "a", false, "bool value")
	boolbFlag := f.BoolP("boolb", "b", false, "bool2 value")
	boolcFlag := f.BoolP("boolc", "c", false, "bool3 value")
	booldFlag := f.BoolP("boold", "d", false, "bool4 value")
	stringaFlag := f.StringP("stringa", "s", "0", "string value")
	stringzFlag := f.StringP("stringz", "z", "0", "string value")
	extra := "interspersed-argument"
	notaflag := "--i-look-like-a-flag"
	args := []string***REMOVED***
		"-ab",
		extra,
		"-cs",
		"hello",
		"-z=something",
		"-d=true",
		"--",
		notaflag,
	***REMOVED***
	f.SetOutput(ioutil.Discard)
	if err := f.Parse(args); err != nil ***REMOVED***
		t.Error("expected no error, got ", err)
	***REMOVED***
	if !f.Parsed() ***REMOVED***
		t.Error("f.Parse() = false after Parse")
	***REMOVED***
	if *boolaFlag != true ***REMOVED***
		t.Error("boola flag should be true, is ", *boolaFlag)
	***REMOVED***
	if *boolbFlag != true ***REMOVED***
		t.Error("boolb flag should be true, is ", *boolbFlag)
	***REMOVED***
	if *boolcFlag != true ***REMOVED***
		t.Error("boolc flag should be true, is ", *boolcFlag)
	***REMOVED***
	if *booldFlag != true ***REMOVED***
		t.Error("boold flag should be true, is ", *booldFlag)
	***REMOVED***
	if *stringaFlag != "hello" ***REMOVED***
		t.Error("stringa flag should be `hello`, is ", *stringaFlag)
	***REMOVED***
	if *stringzFlag != "something" ***REMOVED***
		t.Error("stringz flag should be `something`, is ", *stringzFlag)
	***REMOVED***
	if len(f.Args()) != 2 ***REMOVED***
		t.Error("expected one argument, got", len(f.Args()))
	***REMOVED*** else if f.Args()[0] != extra ***REMOVED***
		t.Errorf("expected argument %q got %q", extra, f.Args()[0])
	***REMOVED*** else if f.Args()[1] != notaflag ***REMOVED***
		t.Errorf("expected argument %q got %q", notaflag, f.Args()[1])
	***REMOVED***
	if f.ArgsLenAtDash() != 1 ***REMOVED***
		t.Errorf("expected argsLenAtDash %d got %d", f.ArgsLenAtDash(), 1)
	***REMOVED***
***REMOVED***

func TestShorthandLookup(t *testing.T) ***REMOVED***
	f := NewFlagSet("shorthand", ContinueOnError)
	if f.Parsed() ***REMOVED***
		t.Error("f.Parse() = true before Parse")
	***REMOVED***
	f.BoolP("boola", "a", false, "bool value")
	f.BoolP("boolb", "b", false, "bool2 value")
	args := []string***REMOVED***
		"-ab",
	***REMOVED***
	f.SetOutput(ioutil.Discard)
	if err := f.Parse(args); err != nil ***REMOVED***
		t.Error("expected no error, got ", err)
	***REMOVED***
	if !f.Parsed() ***REMOVED***
		t.Error("f.Parse() = false after Parse")
	***REMOVED***
	flag := f.ShorthandLookup("a")
	if flag == nil ***REMOVED***
		t.Errorf("f.ShorthandLookup(\"a\") returned nil")
	***REMOVED***
	if flag.Name != "boola" ***REMOVED***
		t.Errorf("f.ShorthandLookup(\"a\") found %q instead of \"boola\"", flag.Name)
	***REMOVED***
	flag = f.ShorthandLookup("")
	if flag != nil ***REMOVED***
		t.Errorf("f.ShorthandLookup(\"\") did not return nil")
	***REMOVED***
	defer func() ***REMOVED***
		recover()
	***REMOVED***()
	flag = f.ShorthandLookup("ab")
	// should NEVER get here. lookup should panic. defer'd func should recover it.
	t.Errorf("f.ShorthandLookup(\"ab\") did not panic")
***REMOVED***

func TestParse(t *testing.T) ***REMOVED***
	ResetForTesting(func() ***REMOVED*** t.Error("bad parse") ***REMOVED***)
	testParse(GetCommandLine(), t)
***REMOVED***

func TestParseAll(t *testing.T) ***REMOVED***
	ResetForTesting(func() ***REMOVED*** t.Error("bad parse") ***REMOVED***)
	testParseAll(GetCommandLine(), t)
***REMOVED***

func TestFlagSetParse(t *testing.T) ***REMOVED***
	testParse(NewFlagSet("test", ContinueOnError), t)
***REMOVED***

func TestChangedHelper(t *testing.T) ***REMOVED***
	f := NewFlagSet("changedtest", ContinueOnError)
	f.Bool("changed", false, "changed bool")
	f.Bool("settrue", true, "true to true")
	f.Bool("setfalse", false, "false to false")
	f.Bool("unchanged", false, "unchanged bool")

	args := []string***REMOVED***"--changed", "--settrue", "--setfalse=false"***REMOVED***
	if err := f.Parse(args); err != nil ***REMOVED***
		t.Error("f.Parse() = false after Parse")
	***REMOVED***
	if !f.Changed("changed") ***REMOVED***
		t.Errorf("--changed wasn't changed!")
	***REMOVED***
	if !f.Changed("settrue") ***REMOVED***
		t.Errorf("--settrue wasn't changed!")
	***REMOVED***
	if !f.Changed("setfalse") ***REMOVED***
		t.Errorf("--setfalse wasn't changed!")
	***REMOVED***
	if f.Changed("unchanged") ***REMOVED***
		t.Errorf("--unchanged was changed!")
	***REMOVED***
	if f.Changed("invalid") ***REMOVED***
		t.Errorf("--invalid was changed!")
	***REMOVED***
	if f.ArgsLenAtDash() != -1 ***REMOVED***
		t.Errorf("Expected argsLenAtDash: %d but got %d", -1, f.ArgsLenAtDash())
	***REMOVED***
***REMOVED***

func replaceSeparators(name string, from []string, to string) string ***REMOVED***
	result := name
	for _, sep := range from ***REMOVED***
		result = strings.Replace(result, sep, to, -1)
	***REMOVED***
	// Type convert to indicate normalization has been done.
	return result
***REMOVED***

func wordSepNormalizeFunc(f *FlagSet, name string) NormalizedName ***REMOVED***
	seps := []string***REMOVED***"-", "_"***REMOVED***
	name = replaceSeparators(name, seps, ".")
	normalizeFlagNameInvocations++

	return NormalizedName(name)
***REMOVED***

func testWordSepNormalizedNames(args []string, t *testing.T) ***REMOVED***
	f := NewFlagSet("normalized", ContinueOnError)
	if f.Parsed() ***REMOVED***
		t.Error("f.Parse() = true before Parse")
	***REMOVED***
	withDashFlag := f.Bool("with-dash-flag", false, "bool value")
	// Set this after some flags have been added and before others.
	f.SetNormalizeFunc(wordSepNormalizeFunc)
	withUnderFlag := f.Bool("with_under_flag", false, "bool value")
	withBothFlag := f.Bool("with-both_flag", false, "bool value")
	if err := f.Parse(args); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !f.Parsed() ***REMOVED***
		t.Error("f.Parse() = false after Parse")
	***REMOVED***
	if *withDashFlag != true ***REMOVED***
		t.Error("withDashFlag flag should be true, is ", *withDashFlag)
	***REMOVED***
	if *withUnderFlag != true ***REMOVED***
		t.Error("withUnderFlag flag should be true, is ", *withUnderFlag)
	***REMOVED***
	if *withBothFlag != true ***REMOVED***
		t.Error("withBothFlag flag should be true, is ", *withBothFlag)
	***REMOVED***
***REMOVED***

func TestWordSepNormalizedNames(t *testing.T) ***REMOVED***
	args := []string***REMOVED***
		"--with-dash-flag",
		"--with-under-flag",
		"--with-both-flag",
	***REMOVED***
	testWordSepNormalizedNames(args, t)

	args = []string***REMOVED***
		"--with_dash_flag",
		"--with_under_flag",
		"--with_both_flag",
	***REMOVED***
	testWordSepNormalizedNames(args, t)

	args = []string***REMOVED***
		"--with-dash_flag",
		"--with-under_flag",
		"--with-both_flag",
	***REMOVED***
	testWordSepNormalizedNames(args, t)
***REMOVED***

func aliasAndWordSepFlagNames(f *FlagSet, name string) NormalizedName ***REMOVED***
	seps := []string***REMOVED***"-", "_"***REMOVED***

	oldName := replaceSeparators("old-valid_flag", seps, ".")
	newName := replaceSeparators("valid-flag", seps, ".")

	name = replaceSeparators(name, seps, ".")
	switch name ***REMOVED***
	case oldName:
		name = newName
	***REMOVED***

	return NormalizedName(name)
***REMOVED***

func TestCustomNormalizedNames(t *testing.T) ***REMOVED***
	f := NewFlagSet("normalized", ContinueOnError)
	if f.Parsed() ***REMOVED***
		t.Error("f.Parse() = true before Parse")
	***REMOVED***

	validFlag := f.Bool("valid-flag", false, "bool value")
	f.SetNormalizeFunc(aliasAndWordSepFlagNames)
	someOtherFlag := f.Bool("some-other-flag", false, "bool value")

	args := []string***REMOVED***"--old_valid_flag", "--some-other_flag"***REMOVED***
	if err := f.Parse(args); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if *validFlag != true ***REMOVED***
		t.Errorf("validFlag is %v even though we set the alias --old_valid_falg", *validFlag)
	***REMOVED***
	if *someOtherFlag != true ***REMOVED***
		t.Error("someOtherFlag should be true, is ", *someOtherFlag)
	***REMOVED***
***REMOVED***

// Every flag we add, the name (displayed also in usage) should normalized
func TestNormalizationFuncShouldChangeFlagName(t *testing.T) ***REMOVED***
	// Test normalization after addition
	f := NewFlagSet("normalized", ContinueOnError)

	f.Bool("valid_flag", false, "bool value")
	if f.Lookup("valid_flag").Name != "valid_flag" ***REMOVED***
		t.Error("The new flag should have the name 'valid_flag' instead of ", f.Lookup("valid_flag").Name)
	***REMOVED***

	f.SetNormalizeFunc(wordSepNormalizeFunc)
	if f.Lookup("valid_flag").Name != "valid.flag" ***REMOVED***
		t.Error("The new flag should have the name 'valid.flag' instead of ", f.Lookup("valid_flag").Name)
	***REMOVED***

	// Test normalization before addition
	f = NewFlagSet("normalized", ContinueOnError)
	f.SetNormalizeFunc(wordSepNormalizeFunc)

	f.Bool("valid_flag", false, "bool value")
	if f.Lookup("valid_flag").Name != "valid.flag" ***REMOVED***
		t.Error("The new flag should have the name 'valid.flag' instead of ", f.Lookup("valid_flag").Name)
	***REMOVED***
***REMOVED***

// Related to https://github.com/spf13/cobra/issues/521.
func TestNormalizationSharedFlags(t *testing.T) ***REMOVED***
	f := NewFlagSet("set f", ContinueOnError)
	g := NewFlagSet("set g", ContinueOnError)
	nfunc := wordSepNormalizeFunc
	testName := "valid_flag"
	normName := nfunc(nil, testName)
	if testName == string(normName) ***REMOVED***
		t.Error("TestNormalizationSharedFlags meaningless: the original and normalized flag names are identical:", testName)
	***REMOVED***

	f.Bool(testName, false, "bool value")
	g.AddFlagSet(f)

	f.SetNormalizeFunc(nfunc)
	g.SetNormalizeFunc(nfunc)

	if len(f.formal) != 1 ***REMOVED***
		t.Error("Normalizing flags should not result in duplications in the flag set:", f.formal)
	***REMOVED***
	if f.orderedFormal[0].Name != string(normName) ***REMOVED***
		t.Error("Flag name not normalized")
	***REMOVED***
	for k := range f.formal ***REMOVED***
		if k != "valid.flag" ***REMOVED***
			t.Errorf("The key in the flag map should have been normalized: wanted \"%s\", got \"%s\" instead", normName, k)
		***REMOVED***
	***REMOVED***

	if !reflect.DeepEqual(f.formal, g.formal) || !reflect.DeepEqual(f.orderedFormal, g.orderedFormal) ***REMOVED***
		t.Error("Two flag sets sharing the same flags should stay consistent after being normalized. Original set:", f.formal, "Duplicate set:", g.formal)
	***REMOVED***
***REMOVED***

func TestNormalizationSetFlags(t *testing.T) ***REMOVED***
	f := NewFlagSet("normalized", ContinueOnError)
	nfunc := wordSepNormalizeFunc
	testName := "valid_flag"
	normName := nfunc(nil, testName)
	if testName == string(normName) ***REMOVED***
		t.Error("TestNormalizationSetFlags meaningless: the original and normalized flag names are identical:", testName)
	***REMOVED***

	f.Bool(testName, false, "bool value")
	f.Set(testName, "true")
	f.SetNormalizeFunc(nfunc)

	if len(f.formal) != 1 ***REMOVED***
		t.Error("Normalizing flags should not result in duplications in the flag set:", f.formal)
	***REMOVED***
	if f.orderedFormal[0].Name != string(normName) ***REMOVED***
		t.Error("Flag name not normalized")
	***REMOVED***
	for k := range f.formal ***REMOVED***
		if k != "valid.flag" ***REMOVED***
			t.Errorf("The key in the flag map should have been normalized: wanted \"%s\", got \"%s\" instead", normName, k)
		***REMOVED***
	***REMOVED***

	if !reflect.DeepEqual(f.formal, f.actual) ***REMOVED***
		t.Error("The map of set flags should get normalized. Formal:", f.formal, "Actual:", f.actual)
	***REMOVED***
***REMOVED***

// Declare a user-defined flag type.
type flagVar []string

func (f *flagVar) String() string ***REMOVED***
	return fmt.Sprint([]string(*f))
***REMOVED***

func (f *flagVar) Set(value string) error ***REMOVED***
	*f = append(*f, value)
	return nil
***REMOVED***

func (f *flagVar) Type() string ***REMOVED***
	return "flagVar"
***REMOVED***

func TestUserDefined(t *testing.T) ***REMOVED***
	var flags FlagSet
	flags.Init("test", ContinueOnError)
	var v flagVar
	flags.VarP(&v, "v", "v", "usage")
	if err := flags.Parse([]string***REMOVED***"--v=1", "-v2", "-v", "3"***REMOVED***); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
	if len(v) != 3 ***REMOVED***
		t.Fatal("expected 3 args; got ", len(v))
	***REMOVED***
	expect := "[1 2 3]"
	if v.String() != expect ***REMOVED***
		t.Errorf("expected value %q got %q", expect, v.String())
	***REMOVED***
***REMOVED***

func TestSetOutput(t *testing.T) ***REMOVED***
	var flags FlagSet
	var buf bytes.Buffer
	flags.SetOutput(&buf)
	flags.Init("test", ContinueOnError)
	flags.Parse([]string***REMOVED***"--unknown"***REMOVED***)
	if out := buf.String(); !strings.Contains(out, "--unknown") ***REMOVED***
		t.Logf("expected output mentioning unknown; got %q", out)
	***REMOVED***
***REMOVED***

// This tests that one can reset the flags. This still works but not well, and is
// superseded by FlagSet.
func TestChangingArgs(t *testing.T) ***REMOVED***
	ResetForTesting(func() ***REMOVED*** t.Fatal("bad parse") ***REMOVED***)
	oldArgs := os.Args
	defer func() ***REMOVED*** os.Args = oldArgs ***REMOVED***()
	os.Args = []string***REMOVED***"cmd", "--before", "subcmd"***REMOVED***
	before := Bool("before", false, "")
	if err := GetCommandLine().Parse(os.Args[1:]); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	cmd := Arg(0)
	os.Args = []string***REMOVED***"subcmd", "--after", "args"***REMOVED***
	after := Bool("after", false, "")
	Parse()
	args := Args()

	if !*before || cmd != "subcmd" || !*after || len(args) != 1 || args[0] != "args" ***REMOVED***
		t.Fatalf("expected true subcmd true [args] got %v %v %v %v", *before, cmd, *after, args)
	***REMOVED***
***REMOVED***

// Test that -help invokes the usage message and returns ErrHelp.
func TestHelp(t *testing.T) ***REMOVED***
	var helpCalled = false
	fs := NewFlagSet("help test", ContinueOnError)
	fs.Usage = func() ***REMOVED*** helpCalled = true ***REMOVED***
	var flag bool
	fs.BoolVar(&flag, "flag", false, "regular flag")
	// Regular flag invocation should work
	err := fs.Parse([]string***REMOVED***"--flag=true"***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got ", err)
	***REMOVED***
	if !flag ***REMOVED***
		t.Error("flag was not set by --flag")
	***REMOVED***
	if helpCalled ***REMOVED***
		t.Error("help called for regular flag")
		helpCalled = false // reset for next test
	***REMOVED***
	// Help flag should work as expected.
	err = fs.Parse([]string***REMOVED***"--help"***REMOVED***)
	if err == nil ***REMOVED***
		t.Fatal("error expected")
	***REMOVED***
	if err != ErrHelp ***REMOVED***
		t.Fatal("expected ErrHelp; got ", err)
	***REMOVED***
	if !helpCalled ***REMOVED***
		t.Fatal("help was not called")
	***REMOVED***
	// If we define a help flag, that should override.
	var help bool
	fs.BoolVar(&help, "help", false, "help flag")
	helpCalled = false
	err = fs.Parse([]string***REMOVED***"--help"***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error for defined --help; got ", err)
	***REMOVED***
	if helpCalled ***REMOVED***
		t.Fatal("help was called; should not have been for defined help flag")
	***REMOVED***
***REMOVED***

func TestNoInterspersed(t *testing.T) ***REMOVED***
	f := NewFlagSet("test", ContinueOnError)
	f.SetInterspersed(false)
	f.Bool("true", true, "always true")
	f.Bool("false", false, "always false")
	err := f.Parse([]string***REMOVED***"--true", "break", "--false"***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got ", err)
	***REMOVED***
	args := f.Args()
	if len(args) != 2 || args[0] != "break" || args[1] != "--false" ***REMOVED***
		t.Fatal("expected interspersed options/non-options to fail")
	***REMOVED***
***REMOVED***

func TestTermination(t *testing.T) ***REMOVED***
	f := NewFlagSet("termination", ContinueOnError)
	boolFlag := f.BoolP("bool", "l", false, "bool value")
	if f.Parsed() ***REMOVED***
		t.Error("f.Parse() = true before Parse")
	***REMOVED***
	arg1 := "ls"
	arg2 := "-l"
	args := []string***REMOVED***
		"--",
		arg1,
		arg2,
	***REMOVED***
	f.SetOutput(ioutil.Discard)
	if err := f.Parse(args); err != nil ***REMOVED***
		t.Fatal("expected no error; got ", err)
	***REMOVED***
	if !f.Parsed() ***REMOVED***
		t.Error("f.Parse() = false after Parse")
	***REMOVED***
	if *boolFlag ***REMOVED***
		t.Error("expected boolFlag=false, got true")
	***REMOVED***
	if len(f.Args()) != 2 ***REMOVED***
		t.Errorf("expected 2 arguments, got %d: %v", len(f.Args()), f.Args())
	***REMOVED***
	if f.Args()[0] != arg1 ***REMOVED***
		t.Errorf("expected argument %q got %q", arg1, f.Args()[0])
	***REMOVED***
	if f.Args()[1] != arg2 ***REMOVED***
		t.Errorf("expected argument %q got %q", arg2, f.Args()[1])
	***REMOVED***
	if f.ArgsLenAtDash() != 0 ***REMOVED***
		t.Errorf("expected argsLenAtDash %d got %d", 0, f.ArgsLenAtDash())
	***REMOVED***
***REMOVED***

func TestDeprecatedFlagInDocs(t *testing.T) ***REMOVED***
	f := NewFlagSet("bob", ContinueOnError)
	f.Bool("badflag", true, "always true")
	f.MarkDeprecated("badflag", "use --good-flag instead")

	out := new(bytes.Buffer)
	f.SetOutput(out)
	f.PrintDefaults()

	if strings.Contains(out.String(), "badflag") ***REMOVED***
		t.Errorf("found deprecated flag in usage!")
	***REMOVED***
***REMOVED***

func TestDeprecatedFlagShorthandInDocs(t *testing.T) ***REMOVED***
	f := NewFlagSet("bob", ContinueOnError)
	name := "noshorthandflag"
	f.BoolP(name, "n", true, "always true")
	f.MarkShorthandDeprecated("noshorthandflag", fmt.Sprintf("use --%s instead", name))

	out := new(bytes.Buffer)
	f.SetOutput(out)
	f.PrintDefaults()

	if strings.Contains(out.String(), "-n,") ***REMOVED***
		t.Errorf("found deprecated flag shorthand in usage!")
	***REMOVED***
***REMOVED***

func parseReturnStderr(t *testing.T, f *FlagSet, args []string) (string, error) ***REMOVED***
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := f.Parse(args)

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() ***REMOVED***
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	***REMOVED***()

	w.Close()
	os.Stderr = oldStderr
	out := <-outC

	return out, err
***REMOVED***

func TestDeprecatedFlagUsage(t *testing.T) ***REMOVED***
	f := NewFlagSet("bob", ContinueOnError)
	f.Bool("badflag", true, "always true")
	usageMsg := "use --good-flag instead"
	f.MarkDeprecated("badflag", usageMsg)

	args := []string***REMOVED***"--badflag"***REMOVED***
	out, err := parseReturnStderr(t, f, args)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got ", err)
	***REMOVED***

	if !strings.Contains(out, usageMsg) ***REMOVED***
		t.Errorf("usageMsg not printed when using a deprecated flag!")
	***REMOVED***
***REMOVED***

func TestDeprecatedFlagShorthandUsage(t *testing.T) ***REMOVED***
	f := NewFlagSet("bob", ContinueOnError)
	name := "noshorthandflag"
	f.BoolP(name, "n", true, "always true")
	usageMsg := fmt.Sprintf("use --%s instead", name)
	f.MarkShorthandDeprecated(name, usageMsg)

	args := []string***REMOVED***"-n"***REMOVED***
	out, err := parseReturnStderr(t, f, args)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got ", err)
	***REMOVED***

	if !strings.Contains(out, usageMsg) ***REMOVED***
		t.Errorf("usageMsg not printed when using a deprecated flag!")
	***REMOVED***
***REMOVED***

func TestDeprecatedFlagUsageNormalized(t *testing.T) ***REMOVED***
	f := NewFlagSet("bob", ContinueOnError)
	f.Bool("bad-double_flag", true, "always true")
	f.SetNormalizeFunc(wordSepNormalizeFunc)
	usageMsg := "use --good-flag instead"
	f.MarkDeprecated("bad_double-flag", usageMsg)

	args := []string***REMOVED***"--bad_double_flag"***REMOVED***
	out, err := parseReturnStderr(t, f, args)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got ", err)
	***REMOVED***

	if !strings.Contains(out, usageMsg) ***REMOVED***
		t.Errorf("usageMsg not printed when using a deprecated flag!")
	***REMOVED***
***REMOVED***

// Name normalization function should be called only once on flag addition
func TestMultipleNormalizeFlagNameInvocations(t *testing.T) ***REMOVED***
	normalizeFlagNameInvocations = 0

	f := NewFlagSet("normalized", ContinueOnError)
	f.SetNormalizeFunc(wordSepNormalizeFunc)
	f.Bool("with_under_flag", false, "bool value")

	if normalizeFlagNameInvocations != 1 ***REMOVED***
		t.Fatal("Expected normalizeFlagNameInvocations to be 1; got ", normalizeFlagNameInvocations)
	***REMOVED***
***REMOVED***

//
func TestHiddenFlagInUsage(t *testing.T) ***REMOVED***
	f := NewFlagSet("bob", ContinueOnError)
	f.Bool("secretFlag", true, "shhh")
	f.MarkHidden("secretFlag")

	out := new(bytes.Buffer)
	f.SetOutput(out)
	f.PrintDefaults()

	if strings.Contains(out.String(), "secretFlag") ***REMOVED***
		t.Errorf("found hidden flag in usage!")
	***REMOVED***
***REMOVED***

//
func TestHiddenFlagUsage(t *testing.T) ***REMOVED***
	f := NewFlagSet("bob", ContinueOnError)
	f.Bool("secretFlag", true, "shhh")
	f.MarkHidden("secretFlag")

	args := []string***REMOVED***"--secretFlag"***REMOVED***
	out, err := parseReturnStderr(t, f, args)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got ", err)
	***REMOVED***

	if strings.Contains(out, "shhh") ***REMOVED***
		t.Errorf("usage message printed when using a hidden flag!")
	***REMOVED***
***REMOVED***

const defaultOutput = `      --A                         for bootstrapping, allow 'any' type
      --Alongflagname             disable bounds checking
  -C, --CCC                       a boolean defaulting to true (default true)
      --D path                    set relative path for local imports
  -E, --EEE num[=1234]            a num with NoOptDefVal (default 4321)
      --F number                  a non-zero number (default 2.7)
      --G float                   a float that defaults to zero
      --IP ip                     IP address with no default
      --IPMask ipMask             Netmask address with no default
      --IPNet ipNet               IP network with no default
      --Ints ints                 int slice with zero default
      --N int                     a non-zero int (default 27)
      --ND1 string[="bar"]        a string with NoOptDefVal (default "foo")
      --ND2 num[=4321]            a num with NoOptDefVal (default 1234)
      --StringArray stringArray   string array with zero default
      --StringSlice strings       string slice with zero default
      --Z int                     an int that defaults to zero
      --custom custom             custom Value implementation
      --customP custom            a VarP with default (default 10)
      --maxT timeout              set timeout for dial
  -v, --verbose count             verbosity
`

// Custom value that satisfies the Value interface.
type customValue int

func (cv *customValue) String() string ***REMOVED*** return fmt.Sprintf("%v", *cv) ***REMOVED***

func (cv *customValue) Set(s string) error ***REMOVED***
	v, err := strconv.ParseInt(s, 0, 64)
	*cv = customValue(v)
	return err
***REMOVED***

func (cv *customValue) Type() string ***REMOVED*** return "custom" ***REMOVED***

func TestPrintDefaults(t *testing.T) ***REMOVED***
	fs := NewFlagSet("print defaults test", ContinueOnError)
	var buf bytes.Buffer
	fs.SetOutput(&buf)
	fs.Bool("A", false, "for bootstrapping, allow 'any' type")
	fs.Bool("Alongflagname", false, "disable bounds checking")
	fs.BoolP("CCC", "C", true, "a boolean defaulting to true")
	fs.String("D", "", "set relative `path` for local imports")
	fs.Float64("F", 2.7, "a non-zero `number`")
	fs.Float64("G", 0, "a float that defaults to zero")
	fs.Int("N", 27, "a non-zero int")
	fs.IntSlice("Ints", []int***REMOVED******REMOVED***, "int slice with zero default")
	fs.IP("IP", nil, "IP address with no default")
	fs.IPMask("IPMask", nil, "Netmask address with no default")
	fs.IPNet("IPNet", net.IPNet***REMOVED******REMOVED***, "IP network with no default")
	fs.Int("Z", 0, "an int that defaults to zero")
	fs.Duration("maxT", 0, "set `timeout` for dial")
	fs.String("ND1", "foo", "a string with NoOptDefVal")
	fs.Lookup("ND1").NoOptDefVal = "bar"
	fs.Int("ND2", 1234, "a `num` with NoOptDefVal")
	fs.Lookup("ND2").NoOptDefVal = "4321"
	fs.IntP("EEE", "E", 4321, "a `num` with NoOptDefVal")
	fs.ShorthandLookup("E").NoOptDefVal = "1234"
	fs.StringSlice("StringSlice", []string***REMOVED******REMOVED***, "string slice with zero default")
	fs.StringArray("StringArray", []string***REMOVED******REMOVED***, "string array with zero default")
	fs.CountP("verbose", "v", "verbosity")

	var cv customValue
	fs.Var(&cv, "custom", "custom Value implementation")

	cv2 := customValue(10)
	fs.VarP(&cv2, "customP", "", "a VarP with default")

	fs.PrintDefaults()
	got := buf.String()
	if got != defaultOutput ***REMOVED***
		fmt.Println("\n" + got)
		fmt.Println("\n" + defaultOutput)
		t.Errorf("got %q want %q\n", got, defaultOutput)
	***REMOVED***
***REMOVED***

func TestVisitAllFlagOrder(t *testing.T) ***REMOVED***
	fs := NewFlagSet("TestVisitAllFlagOrder", ContinueOnError)
	fs.SortFlags = false
	// https://github.com/spf13/pflag/issues/120
	fs.SetNormalizeFunc(func(f *FlagSet, name string) NormalizedName ***REMOVED***
		return NormalizedName(name)
	***REMOVED***)

	names := []string***REMOVED***"C", "B", "A", "D"***REMOVED***
	for _, name := range names ***REMOVED***
		fs.Bool(name, false, "")
	***REMOVED***

	i := 0
	fs.VisitAll(func(f *Flag) ***REMOVED***
		if names[i] != f.Name ***REMOVED***
			t.Errorf("Incorrect order. Expected %v, got %v", names[i], f.Name)
		***REMOVED***
		i++
	***REMOVED***)
***REMOVED***

func TestVisitFlagOrder(t *testing.T) ***REMOVED***
	fs := NewFlagSet("TestVisitFlagOrder", ContinueOnError)
	fs.SortFlags = false
	names := []string***REMOVED***"C", "B", "A", "D"***REMOVED***
	for _, name := range names ***REMOVED***
		fs.Bool(name, false, "")
		fs.Set(name, "true")
	***REMOVED***

	i := 0
	fs.Visit(func(f *Flag) ***REMOVED***
		if names[i] != f.Name ***REMOVED***
			t.Errorf("Incorrect order. Expected %v, got %v", names[i], f.Name)
		***REMOVED***
		i++
	***REMOVED***)
***REMOVED***
