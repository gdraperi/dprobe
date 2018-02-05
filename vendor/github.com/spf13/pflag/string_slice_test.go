// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"fmt"
	"strings"
	"testing"
)

func setUpSSFlagSet(ssp *[]string) *FlagSet ***REMOVED***
	f := NewFlagSet("test", ContinueOnError)
	f.StringSliceVar(ssp, "ss", []string***REMOVED******REMOVED***, "Command separated list!")
	return f
***REMOVED***

func setUpSSFlagSetWithDefault(ssp *[]string) *FlagSet ***REMOVED***
	f := NewFlagSet("test", ContinueOnError)
	f.StringSliceVar(ssp, "ss", []string***REMOVED***"default", "values"***REMOVED***, "Command separated list!")
	return f
***REMOVED***

func TestEmptySS(t *testing.T) ***REMOVED***
	var ss []string
	f := setUpSSFlagSet(&ss)
	err := f.Parse([]string***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	getSS, err := f.GetStringSlice("ss")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetStringSlice():", err)
	***REMOVED***
	if len(getSS) != 0 ***REMOVED***
		t.Fatalf("got ss %v with len=%d but expected length=0", getSS, len(getSS))
	***REMOVED***
***REMOVED***

func TestEmptySSValue(t *testing.T) ***REMOVED***
	var ss []string
	f := setUpSSFlagSet(&ss)
	err := f.Parse([]string***REMOVED***"--ss="***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	getSS, err := f.GetStringSlice("ss")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetStringSlice():", err)
	***REMOVED***
	if len(getSS) != 0 ***REMOVED***
		t.Fatalf("got ss %v with len=%d but expected length=0", getSS, len(getSS))
	***REMOVED***
***REMOVED***

func TestSS(t *testing.T) ***REMOVED***
	var ss []string
	f := setUpSSFlagSet(&ss)

	vals := []string***REMOVED***"one", "two", "4", "3"***REMOVED***
	arg := fmt.Sprintf("--ss=%s", strings.Join(vals, ","))
	err := f.Parse([]string***REMOVED***arg***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	for i, v := range ss ***REMOVED***
		if vals[i] != v ***REMOVED***
			t.Fatalf("expected ss[%d] to be %s but got: %s", i, vals[i], v)
		***REMOVED***
	***REMOVED***

	getSS, err := f.GetStringSlice("ss")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetStringSlice():", err)
	***REMOVED***
	for i, v := range getSS ***REMOVED***
		if vals[i] != v ***REMOVED***
			t.Fatalf("expected ss[%d] to be %s from GetStringSlice but got: %s", i, vals[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSSDefault(t *testing.T) ***REMOVED***
	var ss []string
	f := setUpSSFlagSetWithDefault(&ss)

	vals := []string***REMOVED***"default", "values"***REMOVED***

	err := f.Parse([]string***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	for i, v := range ss ***REMOVED***
		if vals[i] != v ***REMOVED***
			t.Fatalf("expected ss[%d] to be %s but got: %s", i, vals[i], v)
		***REMOVED***
	***REMOVED***

	getSS, err := f.GetStringSlice("ss")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetStringSlice():", err)
	***REMOVED***
	for i, v := range getSS ***REMOVED***
		if vals[i] != v ***REMOVED***
			t.Fatalf("expected ss[%d] to be %s from GetStringSlice but got: %s", i, vals[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSSWithDefault(t *testing.T) ***REMOVED***
	var ss []string
	f := setUpSSFlagSetWithDefault(&ss)

	vals := []string***REMOVED***"one", "two", "4", "3"***REMOVED***
	arg := fmt.Sprintf("--ss=%s", strings.Join(vals, ","))
	err := f.Parse([]string***REMOVED***arg***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	for i, v := range ss ***REMOVED***
		if vals[i] != v ***REMOVED***
			t.Fatalf("expected ss[%d] to be %s but got: %s", i, vals[i], v)
		***REMOVED***
	***REMOVED***

	getSS, err := f.GetStringSlice("ss")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetStringSlice():", err)
	***REMOVED***
	for i, v := range getSS ***REMOVED***
		if vals[i] != v ***REMOVED***
			t.Fatalf("expected ss[%d] to be %s from GetStringSlice but got: %s", i, vals[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSSCalledTwice(t *testing.T) ***REMOVED***
	var ss []string
	f := setUpSSFlagSet(&ss)

	in := []string***REMOVED***"one,two", "three"***REMOVED***
	expected := []string***REMOVED***"one", "two", "three"***REMOVED***
	argfmt := "--ss=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string***REMOVED***arg1, arg2***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	if len(expected) != len(ss) ***REMOVED***
		t.Fatalf("expected number of ss to be %d but got: %d", len(expected), len(ss))
	***REMOVED***
	for i, v := range ss ***REMOVED***
		if expected[i] != v ***REMOVED***
			t.Fatalf("expected ss[%d] to be %s but got: %s", i, expected[i], v)
		***REMOVED***
	***REMOVED***

	values, err := f.GetStringSlice("ss")
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	if len(expected) != len(values) ***REMOVED***
		t.Fatalf("expected number of values to be %d but got: %d", len(expected), len(ss))
	***REMOVED***
	for i, v := range values ***REMOVED***
		if expected[i] != v ***REMOVED***
			t.Fatalf("expected got ss[%d] to be %s but got: %s", i, expected[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSSWithComma(t *testing.T) ***REMOVED***
	var ss []string
	f := setUpSSFlagSet(&ss)

	in := []string***REMOVED***`"one,two"`, `"three"`, `"four,five",six`***REMOVED***
	expected := []string***REMOVED***"one,two", "three", "four,five", "six"***REMOVED***
	argfmt := "--ss=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	arg3 := fmt.Sprintf(argfmt, in[2])
	err := f.Parse([]string***REMOVED***arg1, arg2, arg3***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	if len(expected) != len(ss) ***REMOVED***
		t.Fatalf("expected number of ss to be %d but got: %d", len(expected), len(ss))
	***REMOVED***
	for i, v := range ss ***REMOVED***
		if expected[i] != v ***REMOVED***
			t.Fatalf("expected ss[%d] to be %s but got: %s", i, expected[i], v)
		***REMOVED***
	***REMOVED***

	values, err := f.GetStringSlice("ss")
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	if len(expected) != len(values) ***REMOVED***
		t.Fatalf("expected number of values to be %d but got: %d", len(expected), len(values))
	***REMOVED***
	for i, v := range values ***REMOVED***
		if expected[i] != v ***REMOVED***
			t.Fatalf("expected got ss[%d] to be %s but got: %s", i, expected[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSSWithSquareBrackets(t *testing.T) ***REMOVED***
	var ss []string
	f := setUpSSFlagSet(&ss)

	in := []string***REMOVED***`"[a-z]"`, `"[a-z]+"`***REMOVED***
	expected := []string***REMOVED***"[a-z]", "[a-z]+"***REMOVED***
	argfmt := "--ss=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string***REMOVED***arg1, arg2***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	if len(expected) != len(ss) ***REMOVED***
		t.Fatalf("expected number of ss to be %d but got: %d", len(expected), len(ss))
	***REMOVED***
	for i, v := range ss ***REMOVED***
		if expected[i] != v ***REMOVED***
			t.Fatalf("expected ss[%d] to be %s but got: %s", i, expected[i], v)
		***REMOVED***
	***REMOVED***

	values, err := f.GetStringSlice("ss")
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	if len(expected) != len(values) ***REMOVED***
		t.Fatalf("expected number of values to be %d but got: %d", len(expected), len(values))
	***REMOVED***
	for i, v := range values ***REMOVED***
		if expected[i] != v ***REMOVED***
			t.Fatalf("expected got ss[%d] to be %s but got: %s", i, expected[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***
