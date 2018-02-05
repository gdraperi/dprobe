package mem

import (
	"testing"
	"time"
)

func TestFileDataNameRace(t *testing.T) ***REMOVED***
	t.Parallel()
	const someName = "someName"
	const someOtherName = "someOtherName"
	d := FileData***REMOVED***
		name: someName,
	***REMOVED***

	if d.Name() != someName ***REMOVED***
		t.Errorf("Failed to read correct Name, was %v", d.Name())
	***REMOVED***

	ChangeFileName(&d, someOtherName)
	if d.Name() != someOtherName ***REMOVED***
		t.Errorf("Failed to set Name, was %v", d.Name())
	***REMOVED***

	go func() ***REMOVED***
		ChangeFileName(&d, someName)
	***REMOVED***()

	if d.Name() != someName && d.Name() != someOtherName ***REMOVED***
		t.Errorf("Failed to read either Name, was %v", d.Name())
	***REMOVED***
***REMOVED***

func TestFileDataModTimeRace(t *testing.T) ***REMOVED***
	t.Parallel()
	someTime := time.Now()
	someOtherTime := someTime.Add(1 * time.Minute)

	d := FileData***REMOVED***
		modtime: someTime,
	***REMOVED***

	s := FileInfo***REMOVED***
		FileData: &d,
	***REMOVED***

	if s.ModTime() != someTime ***REMOVED***
		t.Errorf("Failed to read correct value, was %v", s.ModTime())
	***REMOVED***

	SetModTime(&d, someOtherTime)
	if s.ModTime() != someOtherTime ***REMOVED***
		t.Errorf("Failed to set ModTime, was %v", s.ModTime())
	***REMOVED***

	go func() ***REMOVED***
		SetModTime(&d, someTime)
	***REMOVED***()

	if s.ModTime() != someTime && s.ModTime() != someOtherTime ***REMOVED***
		t.Errorf("Failed to read either modtime, was %v", s.ModTime())
	***REMOVED***
***REMOVED***

func TestFileDataModeRace(t *testing.T) ***REMOVED***
	t.Parallel()
	const someMode = 0777
	const someOtherMode = 0660

	d := FileData***REMOVED***
		mode: someMode,
	***REMOVED***

	s := FileInfo***REMOVED***
		FileData: &d,
	***REMOVED***

	if s.Mode() != someMode ***REMOVED***
		t.Errorf("Failed to read correct value, was %v", s.Mode())
	***REMOVED***

	SetMode(&d, someOtherMode)
	if s.Mode() != someOtherMode ***REMOVED***
		t.Errorf("Failed to set Mode, was %v", s.Mode())
	***REMOVED***

	go func() ***REMOVED***
		SetMode(&d, someMode)
	***REMOVED***()

	if s.Mode() != someMode && s.Mode() != someOtherMode ***REMOVED***
		t.Errorf("Failed to read either mode, was %v", s.Mode())
	***REMOVED***
***REMOVED***

func TestFileDataIsDirRace(t *testing.T) ***REMOVED***
	t.Parallel()

	d := FileData***REMOVED***
		dir: true,
	***REMOVED***

	s := FileInfo***REMOVED***
		FileData: &d,
	***REMOVED***

	if s.IsDir() != true ***REMOVED***
		t.Errorf("Failed to read correct value, was %v", s.IsDir())
	***REMOVED***

	go func() ***REMOVED***
		s.Lock()
		d.dir = false
		s.Unlock()
	***REMOVED***()

	//just logging the value to trigger a read:
	t.Logf("Value is %v", s.IsDir())
***REMOVED***

func TestFileDataSizeRace(t *testing.T) ***REMOVED***
	t.Parallel()

	const someData = "Hello"
	const someOtherDataSize = "Hello World"

	d := FileData***REMOVED***
		data: []byte(someData),
		dir:  false,
	***REMOVED***

	s := FileInfo***REMOVED***
		FileData: &d,
	***REMOVED***

	if s.Size() != int64(len(someData)) ***REMOVED***
		t.Errorf("Failed to read correct value, was %v", s.Size())
	***REMOVED***

	go func() ***REMOVED***
		s.Lock()
		d.data = []byte(someOtherDataSize)
		s.Unlock()
	***REMOVED***()

	//just logging the value to trigger a read:
	t.Logf("Value is %v", s.Size())

	//Testing the Dir size case
	d.dir = true
	if s.Size() != int64(42) ***REMOVED***
		t.Errorf("Failed to read correct value for dir, was %v", s.Size())
	***REMOVED***
***REMOVED***
