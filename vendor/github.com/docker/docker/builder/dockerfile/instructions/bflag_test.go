package instructions

import (
	"testing"
)

func TestBuilderFlags(t *testing.T) ***REMOVED***
	var expected string
	var err error

	// ---

	bf := NewBFlags()
	bf.Args = []string***REMOVED******REMOVED***
	if err := bf.Parse(); err != nil ***REMOVED***
		t.Fatalf("Test1 of %q was supposed to work: %s", bf.Args, err)
	***REMOVED***

	// ---

	bf = NewBFlags()
	bf.Args = []string***REMOVED***"--"***REMOVED***
	if err := bf.Parse(); err != nil ***REMOVED***
		t.Fatalf("Test2 of %q was supposed to work: %s", bf.Args, err)
	***REMOVED***

	// ---

	bf = NewBFlags()
	flStr1 := bf.AddString("str1", "")
	flBool1 := bf.AddBool("bool1", false)
	bf.Args = []string***REMOVED******REMOVED***
	if err = bf.Parse(); err != nil ***REMOVED***
		t.Fatalf("Test3 of %q was supposed to work: %s", bf.Args, err)
	***REMOVED***

	if flStr1.IsUsed() ***REMOVED***
		t.Fatal("Test3 - str1 was not used!")
	***REMOVED***
	if flBool1.IsUsed() ***REMOVED***
		t.Fatal("Test3 - bool1 was not used!")
	***REMOVED***

	// ---

	bf = NewBFlags()
	flStr1 = bf.AddString("str1", "HI")
	flBool1 = bf.AddBool("bool1", false)
	bf.Args = []string***REMOVED******REMOVED***

	if err = bf.Parse(); err != nil ***REMOVED***
		t.Fatalf("Test4 of %q was supposed to work: %s", bf.Args, err)
	***REMOVED***

	if flStr1.Value != "HI" ***REMOVED***
		t.Fatal("Str1 was supposed to default to: HI")
	***REMOVED***
	if flBool1.IsTrue() ***REMOVED***
		t.Fatal("Bool1 was supposed to default to: false")
	***REMOVED***
	if flStr1.IsUsed() ***REMOVED***
		t.Fatal("Str1 was not used!")
	***REMOVED***
	if flBool1.IsUsed() ***REMOVED***
		t.Fatal("Bool1 was not used!")
	***REMOVED***

	// ---

	bf = NewBFlags()
	flStr1 = bf.AddString("str1", "HI")
	bf.Args = []string***REMOVED***"--str1"***REMOVED***

	if err = bf.Parse(); err == nil ***REMOVED***
		t.Fatalf("Test %q was supposed to fail", bf.Args)
	***REMOVED***

	// ---

	bf = NewBFlags()
	flStr1 = bf.AddString("str1", "HI")
	bf.Args = []string***REMOVED***"--str1="***REMOVED***

	if err = bf.Parse(); err != nil ***REMOVED***
		t.Fatalf("Test %q was supposed to work: %s", bf.Args, err)
	***REMOVED***

	expected = ""
	if flStr1.Value != expected ***REMOVED***
		t.Fatalf("Str1 (%q) should be: %q", flStr1.Value, expected)
	***REMOVED***

	// ---

	bf = NewBFlags()
	flStr1 = bf.AddString("str1", "HI")
	bf.Args = []string***REMOVED***"--str1=BYE"***REMOVED***

	if err = bf.Parse(); err != nil ***REMOVED***
		t.Fatalf("Test %q was supposed to work: %s", bf.Args, err)
	***REMOVED***

	expected = "BYE"
	if flStr1.Value != expected ***REMOVED***
		t.Fatalf("Str1 (%q) should be: %q", flStr1.Value, expected)
	***REMOVED***

	// ---

	bf = NewBFlags()
	flBool1 = bf.AddBool("bool1", false)
	bf.Args = []string***REMOVED***"--bool1"***REMOVED***

	if err = bf.Parse(); err != nil ***REMOVED***
		t.Fatalf("Test %q was supposed to work: %s", bf.Args, err)
	***REMOVED***

	if !flBool1.IsTrue() ***REMOVED***
		t.Fatal("Test-b1 Bool1 was supposed to be true")
	***REMOVED***

	// ---

	bf = NewBFlags()
	flBool1 = bf.AddBool("bool1", false)
	bf.Args = []string***REMOVED***"--bool1=true"***REMOVED***

	if err = bf.Parse(); err != nil ***REMOVED***
		t.Fatalf("Test %q was supposed to work: %s", bf.Args, err)
	***REMOVED***

	if !flBool1.IsTrue() ***REMOVED***
		t.Fatal("Test-b2 Bool1 was supposed to be true")
	***REMOVED***

	// ---

	bf = NewBFlags()
	flBool1 = bf.AddBool("bool1", false)
	bf.Args = []string***REMOVED***"--bool1=false"***REMOVED***

	if err = bf.Parse(); err != nil ***REMOVED***
		t.Fatalf("Test %q was supposed to work: %s", bf.Args, err)
	***REMOVED***

	if flBool1.IsTrue() ***REMOVED***
		t.Fatal("Test-b3 Bool1 was supposed to be false")
	***REMOVED***

	// ---

	bf = NewBFlags()
	flBool1 = bf.AddBool("bool1", false)
	bf.Args = []string***REMOVED***"--bool1=false1"***REMOVED***

	if err = bf.Parse(); err == nil ***REMOVED***
		t.Fatalf("Test %q was supposed to fail", bf.Args)
	***REMOVED***

	// ---

	bf = NewBFlags()
	flBool1 = bf.AddBool("bool1", false)
	bf.Args = []string***REMOVED***"--bool2"***REMOVED***

	if err = bf.Parse(); err == nil ***REMOVED***
		t.Fatalf("Test %q was supposed to fail", bf.Args)
	***REMOVED***

	// ---

	bf = NewBFlags()
	flStr1 = bf.AddString("str1", "HI")
	flBool1 = bf.AddBool("bool1", false)
	bf.Args = []string***REMOVED***"--bool1", "--str1=BYE"***REMOVED***

	if err = bf.Parse(); err != nil ***REMOVED***
		t.Fatalf("Test %q was supposed to work: %s", bf.Args, err)
	***REMOVED***

	if flStr1.Value != "BYE" ***REMOVED***
		t.Fatalf("Test %s, str1 should be BYE", bf.Args)
	***REMOVED***
	if !flBool1.IsTrue() ***REMOVED***
		t.Fatalf("Test %s, bool1 should be true", bf.Args)
	***REMOVED***
***REMOVED***
