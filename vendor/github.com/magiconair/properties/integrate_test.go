// Copyright 2017 Frank Schroeder. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package properties

import (
	"flag"
	"fmt"
	"testing"
)

// TestFlag verifies Properties.MustFlag without flag.FlagSet.Parse
func TestFlag(t *testing.T) ***REMOVED***
	f := flag.NewFlagSet("src", flag.PanicOnError)
	gotS := f.String("s", "?", "string flag")
	gotI := f.Int("i", -1, "int flag")

	p := NewProperties()
	p.MustSet("s", "t")
	p.MustSet("i", "9")
	p.MustFlag(f)

	if want := "t"; *gotS != want ***REMOVED***
		t.Errorf("Got string s=%q, want %q", *gotS, want)
	***REMOVED***
	if want := 9; *gotI != want ***REMOVED***
		t.Errorf("Got int i=%d, want %d", *gotI, want)
	***REMOVED***
***REMOVED***

// TestFlagOverride verifies Properties.MustFlag with flag.FlagSet.Parse.
func TestFlagOverride(t *testing.T) ***REMOVED***
	f := flag.NewFlagSet("src", flag.PanicOnError)
	gotA := f.Int("a", 1, "remain default")
	gotB := f.Int("b", 2, "customized")
	gotC := f.Int("c", 3, "overridden")

	if err := f.Parse([]string***REMOVED***"-c", "4"***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	p := NewProperties()
	p.MustSet("b", "5")
	p.MustSet("c", "6")
	p.MustFlag(f)

	if want := 1; *gotA != want ***REMOVED***
		t.Errorf("Got remain default a=%d, want %d", *gotA, want)
	***REMOVED***
	if want := 5; *gotB != want ***REMOVED***
		t.Errorf("Got customized b=%d, want %d", *gotB, want)
	***REMOVED***
	if want := 4; *gotC != want ***REMOVED***
		t.Errorf("Got overriden c=%d, want %d", *gotC, want)
	***REMOVED***
***REMOVED***

func ExampleProperties_MustFlag() ***REMOVED***
	x := flag.Int("x", 0, "demo customize")
	y := flag.Int("y", 0, "demo override")

	// Demo alternative for flag.Parse():
	flag.CommandLine.Parse([]string***REMOVED***"-y", "10"***REMOVED***)
	fmt.Printf("flagged as x=%d, y=%d\n", *x, *y)

	p := NewProperties()
	p.MustSet("x", "7")
	p.MustSet("y", "42") // note discard
	p.MustFlag(flag.CommandLine)
	fmt.Printf("configured to x=%d, y=%d\n", *x, *y)

	// Output:
	// flagged as x=0, y=10
	// configured to x=7, y=10
***REMOVED***
