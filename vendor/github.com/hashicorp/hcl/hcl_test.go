package hcl

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

// This is the directory where our test fixtures are.
const fixtureDir = "./test-fixtures"

func testReadFile(t *testing.T, n string) string ***REMOVED***
	d, err := ioutil.ReadFile(filepath.Join(fixtureDir, n))
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***

	return string(d)
***REMOVED***
