package tailfile

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestTailFile(t *testing.T) ***REMOVED***
	f, err := ioutil.TempFile("", "tail-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer f.Close()
	defer os.RemoveAll(f.Name())
	testFile := []byte(`first line
second line
third line
fourth line
fifth line
next first line
next second line
next third line
next fourth line
next fifth line
last first line
next first line
next second line
next third line
next fourth line
next fifth line
next first line
next second line
next third line
next fourth line
next fifth line
last second line
last third line
last fourth line
last fifth line
truncated line`)
	if _, err := f.Write(testFile); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := f.Seek(0, os.SEEK_SET); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expected := []string***REMOVED***"last fourth line", "last fifth line"***REMOVED***
	res, err := TailFile(f, 2)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	for i, l := range res ***REMOVED***
		t.Logf("%s", l)
		if expected[i] != string(l) ***REMOVED***
			t.Fatalf("Expected line %s, got %s", expected[i], l)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTailFileManyLines(t *testing.T) ***REMOVED***
	f, err := ioutil.TempFile("", "tail-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer f.Close()
	defer os.RemoveAll(f.Name())
	testFile := []byte(`first line
second line
truncated line`)
	if _, err := f.Write(testFile); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := f.Seek(0, os.SEEK_SET); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expected := []string***REMOVED***"first line", "second line"***REMOVED***
	res, err := TailFile(f, 10000)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	for i, l := range res ***REMOVED***
		t.Logf("%s", l)
		if expected[i] != string(l) ***REMOVED***
			t.Fatalf("Expected line %s, got %s", expected[i], l)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTailEmptyFile(t *testing.T) ***REMOVED***
	f, err := ioutil.TempFile("", "tail-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer f.Close()
	defer os.RemoveAll(f.Name())
	res, err := TailFile(f, 10000)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if len(res) != 0 ***REMOVED***
		t.Fatal("Must be empty slice from empty file")
	***REMOVED***
***REMOVED***

func TestTailNegativeN(t *testing.T) ***REMOVED***
	f, err := ioutil.TempFile("", "tail-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer f.Close()
	defer os.RemoveAll(f.Name())
	testFile := []byte(`first line
second line
truncated line`)
	if _, err := f.Write(testFile); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := f.Seek(0, os.SEEK_SET); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := TailFile(f, -1); err != ErrNonPositiveLinesNumber ***REMOVED***
		t.Fatalf("Expected ErrNonPositiveLinesNumber, got %s", err)
	***REMOVED***
	if _, err := TailFile(f, 0); err != ErrNonPositiveLinesNumber ***REMOVED***
		t.Fatalf("Expected ErrNonPositiveLinesNumber, got %s", err)
	***REMOVED***
***REMOVED***

func BenchmarkTail(b *testing.B) ***REMOVED***
	f, err := ioutil.TempFile("", "tail-test")
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	defer f.Close()
	defer os.RemoveAll(f.Name())
	for i := 0; i < 10000; i++ ***REMOVED***
		if _, err := f.Write([]byte("tailfile pretty interesting line\n")); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if _, err := TailFile(f, 1000); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***
