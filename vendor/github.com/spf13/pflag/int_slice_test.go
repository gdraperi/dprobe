// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func setUpISFlagSet(isp *[]int) *FlagSet ***REMOVED***
	f := NewFlagSet("test", ContinueOnError)
	f.IntSliceVar(isp, "is", []int***REMOVED******REMOVED***, "Command separated list!")
	return f
***REMOVED***

func setUpISFlagSetWithDefault(isp *[]int) *FlagSet ***REMOVED***
	f := NewFlagSet("test", ContinueOnError)
	f.IntSliceVar(isp, "is", []int***REMOVED***0, 1***REMOVED***, "Command separated list!")
	return f
***REMOVED***

func TestEmptyIS(t *testing.T) ***REMOVED***
	var is []int
	f := setUpISFlagSet(&is)
	err := f.Parse([]string***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	getIS, err := f.GetIntSlice("is")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetIntSlice():", err)
	***REMOVED***
	if len(getIS) != 0 ***REMOVED***
		t.Fatalf("got is %v with len=%d but expected length=0", getIS, len(getIS))
	***REMOVED***
***REMOVED***

func TestIS(t *testing.T) ***REMOVED***
	var is []int
	f := setUpISFlagSet(&is)

	vals := []string***REMOVED***"1", "2", "4", "3"***REMOVED***
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string***REMOVED***arg***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	for i, v := range is ***REMOVED***
		d, err := strconv.Atoi(vals[i])
		if err != nil ***REMOVED***
			t.Fatalf("got error: %v", err)
		***REMOVED***
		if d != v ***REMOVED***
			t.Fatalf("expected is[%d] to be %s but got: %d", i, vals[i], v)
		***REMOVED***
	***REMOVED***
	getIS, err := f.GetIntSlice("is")
	if err != nil ***REMOVED***
		t.Fatalf("got error: %v", err)
	***REMOVED***
	for i, v := range getIS ***REMOVED***
		d, err := strconv.Atoi(vals[i])
		if err != nil ***REMOVED***
			t.Fatalf("got error: %v", err)
		***REMOVED***
		if d != v ***REMOVED***
			t.Fatalf("expected is[%d] to be %s but got: %d from GetIntSlice", i, vals[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestISDefault(t *testing.T) ***REMOVED***
	var is []int
	f := setUpISFlagSetWithDefault(&is)

	vals := []string***REMOVED***"0", "1"***REMOVED***

	err := f.Parse([]string***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	for i, v := range is ***REMOVED***
		d, err := strconv.Atoi(vals[i])
		if err != nil ***REMOVED***
			t.Fatalf("got error: %v", err)
		***REMOVED***
		if d != v ***REMOVED***
			t.Fatalf("expected is[%d] to be %d but got: %d", i, d, v)
		***REMOVED***
	***REMOVED***

	getIS, err := f.GetIntSlice("is")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetIntSlice():", err)
	***REMOVED***
	for i, v := range getIS ***REMOVED***
		d, err := strconv.Atoi(vals[i])
		if err != nil ***REMOVED***
			t.Fatal("got an error from GetIntSlice():", err)
		***REMOVED***
		if d != v ***REMOVED***
			t.Fatalf("expected is[%d] to be %d from GetIntSlice but got: %d", i, d, v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestISWithDefault(t *testing.T) ***REMOVED***
	var is []int
	f := setUpISFlagSetWithDefault(&is)

	vals := []string***REMOVED***"1", "2"***REMOVED***
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string***REMOVED***arg***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	for i, v := range is ***REMOVED***
		d, err := strconv.Atoi(vals[i])
		if err != nil ***REMOVED***
			t.Fatalf("got error: %v", err)
		***REMOVED***
		if d != v ***REMOVED***
			t.Fatalf("expected is[%d] to be %d but got: %d", i, d, v)
		***REMOVED***
	***REMOVED***

	getIS, err := f.GetIntSlice("is")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetIntSlice():", err)
	***REMOVED***
	for i, v := range getIS ***REMOVED***
		d, err := strconv.Atoi(vals[i])
		if err != nil ***REMOVED***
			t.Fatalf("got error: %v", err)
		***REMOVED***
		if d != v ***REMOVED***
			t.Fatalf("expected is[%d] to be %d from GetIntSlice but got: %d", i, d, v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestISCalledTwice(t *testing.T) ***REMOVED***
	var is []int
	f := setUpISFlagSet(&is)

	in := []string***REMOVED***"1,2", "3"***REMOVED***
	expected := []int***REMOVED***1, 2, 3***REMOVED***
	argfmt := "--is=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string***REMOVED***arg1, arg2***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	for i, v := range is ***REMOVED***
		if expected[i] != v ***REMOVED***
			t.Fatalf("expected is[%d] to be %d but got: %d", i, expected[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***
