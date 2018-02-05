// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"fmt"
	"testing"
)

func setUpSAFlagSet(sap *[]string) *FlagSet ***REMOVED***
	f := NewFlagSet("test", ContinueOnError)
	f.StringArrayVar(sap, "sa", []string***REMOVED******REMOVED***, "Command separated list!")
	return f
***REMOVED***

func setUpSAFlagSetWithDefault(sap *[]string) *FlagSet ***REMOVED***
	f := NewFlagSet("test", ContinueOnError)
	f.StringArrayVar(sap, "sa", []string***REMOVED***"default", "values"***REMOVED***, "Command separated list!")
	return f
***REMOVED***

func TestEmptySA(t *testing.T) ***REMOVED***
	var sa []string
	f := setUpSAFlagSet(&sa)
	err := f.Parse([]string***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	getSA, err := f.GetStringArray("sa")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetStringArray():", err)
	***REMOVED***
	if len(getSA) != 0 ***REMOVED***
		t.Fatalf("got sa %v with len=%d but expected length=0", getSA, len(getSA))
	***REMOVED***
***REMOVED***

func TestEmptySAValue(t *testing.T) ***REMOVED***
	var sa []string
	f := setUpSAFlagSet(&sa)
	err := f.Parse([]string***REMOVED***"--sa="***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	getSA, err := f.GetStringArray("sa")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetStringArray():", err)
	***REMOVED***
	if len(getSA) != 0 ***REMOVED***
		t.Fatalf("got sa %v with len=%d but expected length=0", getSA, len(getSA))
	***REMOVED***
***REMOVED***

func TestSADefault(t *testing.T) ***REMOVED***
	var sa []string
	f := setUpSAFlagSetWithDefault(&sa)

	vals := []string***REMOVED***"default", "values"***REMOVED***

	err := f.Parse([]string***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	for i, v := range sa ***REMOVED***
		if vals[i] != v ***REMOVED***
			t.Fatalf("expected sa[%d] to be %s but got: %s", i, vals[i], v)
		***REMOVED***
	***REMOVED***

	getSA, err := f.GetStringArray("sa")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetStringArray():", err)
	***REMOVED***
	for i, v := range getSA ***REMOVED***
		if vals[i] != v ***REMOVED***
			t.Fatalf("expected sa[%d] to be %s from GetStringArray but got: %s", i, vals[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSAWithDefault(t *testing.T) ***REMOVED***
	var sa []string
	f := setUpSAFlagSetWithDefault(&sa)

	val := "one"
	arg := fmt.Sprintf("--sa=%s", val)
	err := f.Parse([]string***REMOVED***arg***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	if len(sa) != 1 ***REMOVED***
		t.Fatalf("expected number of values to be %d but %d", 1, len(sa))
	***REMOVED***

	if sa[0] != val ***REMOVED***
		t.Fatalf("expected value to be %s but got: %s", sa[0], val)
	***REMOVED***

	getSA, err := f.GetStringArray("sa")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetStringArray():", err)
	***REMOVED***

	if len(getSA) != 1 ***REMOVED***
		t.Fatalf("expected number of values to be %d but %d", 1, len(getSA))
	***REMOVED***

	if getSA[0] != val ***REMOVED***
		t.Fatalf("expected value to be %s but got: %s", getSA[0], val)
	***REMOVED***
***REMOVED***

func TestSACalledTwice(t *testing.T) ***REMOVED***
	var sa []string
	f := setUpSAFlagSet(&sa)

	in := []string***REMOVED***"one", "two"***REMOVED***
	expected := []string***REMOVED***"one", "two"***REMOVED***
	argfmt := "--sa=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string***REMOVED***arg1, arg2***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	if len(expected) != len(sa) ***REMOVED***
		t.Fatalf("expected number of sa to be %d but got: %d", len(expected), len(sa))
	***REMOVED***
	for i, v := range sa ***REMOVED***
		if expected[i] != v ***REMOVED***
			t.Fatalf("expected sa[%d] to be %s but got: %s", i, expected[i], v)
		***REMOVED***
	***REMOVED***

	values, err := f.GetStringArray("sa")
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	if len(expected) != len(values) ***REMOVED***
		t.Fatalf("expected number of values to be %d but got: %d", len(expected), len(sa))
	***REMOVED***
	for i, v := range values ***REMOVED***
		if expected[i] != v ***REMOVED***
			t.Fatalf("expected got sa[%d] to be %s but got: %s", i, expected[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSAWithSpecialChar(t *testing.T) ***REMOVED***
	var sa []string
	f := setUpSAFlagSet(&sa)

	in := []string***REMOVED***"one,two", `"three"`, `"four,five",six`, "seven eight"***REMOVED***
	expected := []string***REMOVED***"one,two", `"three"`, `"four,five",six`, "seven eight"***REMOVED***
	argfmt := "--sa=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	arg3 := fmt.Sprintf(argfmt, in[2])
	arg4 := fmt.Sprintf(argfmt, in[3])
	err := f.Parse([]string***REMOVED***arg1, arg2, arg3, arg4***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	if len(expected) != len(sa) ***REMOVED***
		t.Fatalf("expected number of sa to be %d but got: %d", len(expected), len(sa))
	***REMOVED***
	for i, v := range sa ***REMOVED***
		if expected[i] != v ***REMOVED***
			t.Fatalf("expected sa[%d] to be %s but got: %s", i, expected[i], v)
		***REMOVED***
	***REMOVED***

	values, err := f.GetStringArray("sa")
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	if len(expected) != len(values) ***REMOVED***
		t.Fatalf("expected number of values to be %d but got: %d", len(expected), len(values))
	***REMOVED***
	for i, v := range values ***REMOVED***
		if expected[i] != v ***REMOVED***
			t.Fatalf("expected got sa[%d] to be %s but got: %s", i, expected[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSAWithSquareBrackets(t *testing.T) ***REMOVED***
	var sa []string
	f := setUpSAFlagSet(&sa)

	in := []string***REMOVED***"][]-[", "[a-z]", "[a-z]+"***REMOVED***
	expected := []string***REMOVED***"][]-[", "[a-z]", "[a-z]+"***REMOVED***
	argfmt := "--sa=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	arg3 := fmt.Sprintf(argfmt, in[2])
	err := f.Parse([]string***REMOVED***arg1, arg2, arg3***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	if len(expected) != len(sa) ***REMOVED***
		t.Fatalf("expected number of sa to be %d but got: %d", len(expected), len(sa))
	***REMOVED***
	for i, v := range sa ***REMOVED***
		if expected[i] != v ***REMOVED***
			t.Fatalf("expected sa[%d] to be %s but got: %s", i, expected[i], v)
		***REMOVED***
	***REMOVED***

	values, err := f.GetStringArray("sa")
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	if len(expected) != len(values) ***REMOVED***
		t.Fatalf("expected number of values to be %d but got: %d", len(expected), len(values))
	***REMOVED***
	for i, v := range values ***REMOVED***
		if expected[i] != v ***REMOVED***
			t.Fatalf("expected got sa[%d] to be %s but got: %s", i, expected[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***
