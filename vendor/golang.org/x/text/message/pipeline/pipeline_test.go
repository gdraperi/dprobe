// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipeline

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/text/language"
)

var genFiles = flag.Bool("gen", false, "generate output files instead of comparing")

// setHelper is testing.T.Helper on Go 1.9+, overridden by go19_test.go.
var setHelper = func(t *testing.T) ***REMOVED******REMOVED***

func TestFullCycle(t *testing.T) ***REMOVED***
	const path = "./testdata"
	dirs, err := ioutil.ReadDir(path)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	for _, f := range dirs ***REMOVED***
		t.Run(f.Name(), func(t *testing.T) ***REMOVED***
			chk := func(t *testing.T, err error) ***REMOVED***
				setHelper(t)
				if err != nil ***REMOVED***
					t.Fatal(err)
				***REMOVED***
			***REMOVED***
			dir := filepath.Join(path, f.Name())
			pkgPath := fmt.Sprintf("%s/%s", path, f.Name())
			config := Config***REMOVED***
				SourceLanguage: language.AmericanEnglish,
				Packages:       []string***REMOVED***pkgPath***REMOVED***,
				Dir:            filepath.Join(dir, "locales"),
				GenFile:        "catalog_gen.go",
				GenPackage:     pkgPath,
			***REMOVED***
			// TODO: load config if available.
			s, err := Extract(&config)
			chk(t, err)
			chk(t, s.Import())
			chk(t, s.Merge())
			// TODO:
			//  for range s.Config.Actions ***REMOVED***
			//  	//  TODO: do the actions.
			//  ***REMOVED***
			chk(t, s.Export())
			chk(t, s.Generate())

			writeJSON(t, filepath.Join(dir, "extracted.gotext.json"), s.Extracted)
			checkOutput(t, dir)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func checkOutput(t *testing.T, p string) ***REMOVED***
	filepath.Walk(p, func(p string, f os.FileInfo, err error) error ***REMOVED***
		if f.IsDir() ***REMOVED***
			return nil
		***REMOVED***
		if filepath.Ext(p) != ".want" ***REMOVED***
			return nil
		***REMOVED***
		gotFile := p[:len(p)-len(".want")]
		got, err := ioutil.ReadFile(gotFile)
		if err != nil ***REMOVED***
			t.Errorf("failed to read %q", p)
			return nil
		***REMOVED***
		if *genFiles ***REMOVED***
			if err := ioutil.WriteFile(p, got, 0644); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
		***REMOVED***
		want, err := ioutil.ReadFile(p)
		if err != nil ***REMOVED***
			t.Errorf("failed to read %q", p)
		***REMOVED*** else ***REMOVED***
			scanGot := bufio.NewScanner(bytes.NewReader(got))
			scanWant := bufio.NewScanner(bytes.NewReader(want))
			line := 0
			clean := func(s string) string ***REMOVED***
				if i := strings.LastIndex(s, "//"); i != -1 ***REMOVED***
					s = s[:i]
				***REMOVED***
				return path.Clean(filepath.ToSlash(s))
			***REMOVED***
			for scanGot.Scan() && scanWant.Scan() ***REMOVED***
				got := clean(scanGot.Text())
				want := clean(scanWant.Text())
				if got != want ***REMOVED***
					t.Errorf("file %q differs from .want file at line %d:\n\t%s\n\t%s", gotFile, line, got, want)
					break
				***REMOVED***
				line++
			***REMOVED***
			if scanGot.Scan() || scanWant.Scan() ***REMOVED***
				t.Errorf("file %q differs from .want file at line %d.", gotFile, line)
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)
***REMOVED***

func writeJSON(t *testing.T, path string, x interface***REMOVED******REMOVED***) ***REMOVED***
	data, err := json.MarshalIndent(x, "", "    ")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := ioutil.WriteFile(path, data, 0644); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
