// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"bytes"
	"strconv"
	"testing"
)

// This value can be a boolean ("true", "false") or "maybe"
type triStateValue int

const (
	triStateFalse triStateValue = 0
	triStateTrue  triStateValue = 1
	triStateMaybe triStateValue = 2
)

const strTriStateMaybe = "maybe"

func (v *triStateValue) IsBoolFlag() bool ***REMOVED***
	return true
***REMOVED***

func (v *triStateValue) Get() interface***REMOVED******REMOVED*** ***REMOVED***
	return triStateValue(*v)
***REMOVED***

func (v *triStateValue) Set(s string) error ***REMOVED***
	if s == strTriStateMaybe ***REMOVED***
		*v = triStateMaybe
		return nil
	***REMOVED***
	boolVal, err := strconv.ParseBool(s)
	if boolVal ***REMOVED***
		*v = triStateTrue
	***REMOVED*** else ***REMOVED***
		*v = triStateFalse
	***REMOVED***
	return err
***REMOVED***

func (v *triStateValue) String() string ***REMOVED***
	if *v == triStateMaybe ***REMOVED***
		return strTriStateMaybe
	***REMOVED***
	return strconv.FormatBool(*v == triStateTrue)
***REMOVED***

// The type of the flag as required by the pflag.Value interface
func (v *triStateValue) Type() string ***REMOVED***
	return "version"
***REMOVED***

func setUpFlagSet(tristate *triStateValue) *FlagSet ***REMOVED***
	f := NewFlagSet("test", ContinueOnError)
	*tristate = triStateFalse
	flag := f.VarPF(tristate, "tristate", "t", "tristate value (true, maybe or false)")
	flag.NoOptDefVal = "true"
	return f
***REMOVED***

func TestExplicitTrue(t *testing.T) ***REMOVED***
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	err := f.Parse([]string***REMOVED***"--tristate=true"***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	if tristate != triStateTrue ***REMOVED***
		t.Fatal("expected", triStateTrue, "(triStateTrue) but got", tristate, "instead")
	***REMOVED***
***REMOVED***

func TestImplicitTrue(t *testing.T) ***REMOVED***
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	err := f.Parse([]string***REMOVED***"--tristate"***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	if tristate != triStateTrue ***REMOVED***
		t.Fatal("expected", triStateTrue, "(triStateTrue) but got", tristate, "instead")
	***REMOVED***
***REMOVED***

func TestShortFlag(t *testing.T) ***REMOVED***
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	err := f.Parse([]string***REMOVED***"-t"***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	if tristate != triStateTrue ***REMOVED***
		t.Fatal("expected", triStateTrue, "(triStateTrue) but got", tristate, "instead")
	***REMOVED***
***REMOVED***

func TestShortFlagExtraArgument(t *testing.T) ***REMOVED***
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	// The"maybe"turns into an arg, since short boolean options will only do true/false
	err := f.Parse([]string***REMOVED***"-t", "maybe"***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	if tristate != triStateTrue ***REMOVED***
		t.Fatal("expected", triStateTrue, "(triStateTrue) but got", tristate, "instead")
	***REMOVED***
	args := f.Args()
	if len(args) != 1 || args[0] != "maybe" ***REMOVED***
		t.Fatal("expected an extra 'maybe' argument to stick around")
	***REMOVED***
***REMOVED***

func TestExplicitMaybe(t *testing.T) ***REMOVED***
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	err := f.Parse([]string***REMOVED***"--tristate=maybe"***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	if tristate != triStateMaybe ***REMOVED***
		t.Fatal("expected", triStateMaybe, "(triStateMaybe) but got", tristate, "instead")
	***REMOVED***
***REMOVED***

func TestExplicitFalse(t *testing.T) ***REMOVED***
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	err := f.Parse([]string***REMOVED***"--tristate=false"***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	if tristate != triStateFalse ***REMOVED***
		t.Fatal("expected", triStateFalse, "(triStateFalse) but got", tristate, "instead")
	***REMOVED***
***REMOVED***

func TestImplicitFalse(t *testing.T) ***REMOVED***
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	err := f.Parse([]string***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	if tristate != triStateFalse ***REMOVED***
		t.Fatal("expected", triStateFalse, "(triStateFalse) but got", tristate, "instead")
	***REMOVED***
***REMOVED***

func TestInvalidValue(t *testing.T) ***REMOVED***
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	var buf bytes.Buffer
	f.SetOutput(&buf)
	err := f.Parse([]string***REMOVED***"--tristate=invalid"***REMOVED***)
	if err == nil ***REMOVED***
		t.Fatal("expected an error but did not get any, tristate has value", tristate)
	***REMOVED***
***REMOVED***

func TestBoolP(t *testing.T) ***REMOVED***
	b := BoolP("bool", "b", false, "bool value in CommandLine")
	c := BoolP("c", "c", false, "other bool value")
	args := []string***REMOVED***"--bool"***REMOVED***
	if err := CommandLine.Parse(args); err != nil ***REMOVED***
		t.Error("expected no error, got ", err)
	***REMOVED***
	if *b != true ***REMOVED***
		t.Errorf("expected b=true got b=%v", *b)
	***REMOVED***
	if *c != false ***REMOVED***
		t.Errorf("expect c=false got c=%v", *c)
	***REMOVED***
***REMOVED***
