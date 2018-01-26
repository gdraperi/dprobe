package dockerignore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestReadAll(t *testing.T) ***REMOVED***
	tmpDir, err := ioutil.TempDir("", "dockerignore-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpDir)

	di, err := ReadAll(nil)
	if err != nil ***REMOVED***
		t.Fatalf("Expected not to have error, got %v", err)
	***REMOVED***

	if diLen := len(di); diLen != 0 ***REMOVED***
		t.Fatalf("Expected to have zero dockerignore entry, got %d", diLen)
	***REMOVED***

	diName := filepath.Join(tmpDir, ".dockerignore")
	content := fmt.Sprintf("test1\n/test2\n/a/file/here\n\nlastfile\n# this is a comment\n! /inverted/abs/path\n!\n! \n")
	err = ioutil.WriteFile(diName, []byte(content), 0777)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	diFd, err := os.Open(diName)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer diFd.Close()

	di, err = ReadAll(diFd)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if len(di) != 7 ***REMOVED***
		t.Fatalf("Expected 5 entries, got %v", len(di))
	***REMOVED***
	if di[0] != "test1" ***REMOVED***
		t.Fatal("First element is not test1")
	***REMOVED***
	if di[1] != "test2" ***REMOVED*** // according to https://docs.docker.com/engine/reference/builder/#dockerignore-file, /foo/bar should be treated as foo/bar
		t.Fatal("Second element is not test2")
	***REMOVED***
	if di[2] != "a/file/here" ***REMOVED*** // according to https://docs.docker.com/engine/reference/builder/#dockerignore-file, /foo/bar should be treated as foo/bar
		t.Fatal("Third element is not a/file/here")
	***REMOVED***
	if di[3] != "lastfile" ***REMOVED***
		t.Fatal("Fourth element is not lastfile")
	***REMOVED***
	if di[4] != "!inverted/abs/path" ***REMOVED***
		t.Fatal("Fifth element is not !inverted/abs/path")
	***REMOVED***
	if di[5] != "!" ***REMOVED***
		t.Fatalf("Sixth element is not !, but %s", di[5])
	***REMOVED***
	if di[6] != "!" ***REMOVED***
		t.Fatalf("Sixth element is not !, but %s", di[6])
	***REMOVED***
***REMOVED***
