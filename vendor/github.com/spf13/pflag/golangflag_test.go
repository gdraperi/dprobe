// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	goflag "flag"
	"testing"
)

func TestGoflags(t *testing.T) ***REMOVED***
	goflag.String("stringFlag", "stringFlag", "stringFlag")
	goflag.Bool("boolFlag", false, "boolFlag")

	f := NewFlagSet("test", ContinueOnError)

	f.AddGoFlagSet(goflag.CommandLine)
	err := f.Parse([]string***REMOVED***"--stringFlag=bob", "--boolFlag"***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; get", err)
	***REMOVED***

	getString, err := f.GetString("stringFlag")
	if err != nil ***REMOVED***
		t.Fatal("expected no error; get", err)
	***REMOVED***
	if getString != "bob" ***REMOVED***
		t.Fatalf("expected getString=bob but got getString=%s", getString)
	***REMOVED***

	getBool, err := f.GetBool("boolFlag")
	if err != nil ***REMOVED***
		t.Fatal("expected no error; get", err)
	***REMOVED***
	if getBool != true ***REMOVED***
		t.Fatalf("expected getBool=true but got getBool=%v", getBool)
	***REMOVED***
***REMOVED***
