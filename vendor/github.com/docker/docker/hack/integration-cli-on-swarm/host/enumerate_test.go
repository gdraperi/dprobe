package main

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func getRepoTopDir(t *testing.T) string ***REMOVED***
	wd, err := os.Getwd()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	wd = filepath.Clean(wd)
	suffix := "hack/integration-cli-on-swarm/host"
	if !strings.HasSuffix(wd, suffix) ***REMOVED***
		t.Skipf("cwd seems strange (needs to have suffix %s): %v", suffix, wd)
	***REMOVED***
	return filepath.Clean(filepath.Join(wd, "../../.."))
***REMOVED***

func TestEnumerateTests(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip("skipping in short mode")
	***REMOVED***
	tests, err := enumerateTests(getRepoTopDir(t))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	sort.Strings(tests)
	t.Logf("enumerated %d test filter strings:", len(tests))
	for _, s := range tests ***REMOVED***
		t.Logf("- %q", s)
	***REMOVED***
***REMOVED***

func TestEnumerateTestsForBytes(t *testing.T) ***REMOVED***
	b := []byte(`package main
import (
	"github.com/go-check/check"
)

func (s *FooSuite) TestA(c *check.C) ***REMOVED***
***REMOVED***

func (s *FooSuite) TestAAA(c *check.C) ***REMOVED***
***REMOVED***

func (s *BarSuite) TestBar(c *check.C) ***REMOVED***
***REMOVED***

func (x *FooSuite) TestC(c *check.C) ***REMOVED***
***REMOVED***

func (*FooSuite) TestD(c *check.C) ***REMOVED***
***REMOVED***

// should not be counted
func (s *FooSuite) testE(c *check.C) ***REMOVED***
***REMOVED***

// counted, although we don't support ungofmt file
  func   (s *FooSuite)    TestF  (c   *check.C)***REMOVED******REMOVED***
`)
	expected := []string***REMOVED***
		"FooSuite.TestA$",
		"FooSuite.TestAAA$",
		"BarSuite.TestBar$",
		"FooSuite.TestC$",
		"FooSuite.TestD$",
		"FooSuite.TestF$",
	***REMOVED***

	actual, err := enumerateTestsForBytes(b)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !reflect.DeepEqual(expected, actual) ***REMOVED***
		t.Fatalf("expected %q, got %q", expected, actual)
	***REMOVED***
***REMOVED***
