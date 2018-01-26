package opts

import (
	"testing"

	"github.com/docker/go-units"
)

func TestUlimitOpt(t *testing.T) ***REMOVED***
	ulimitMap := map[string]*units.Ulimit***REMOVED***
		"nofile": ***REMOVED***"nofile", 1024, 512***REMOVED***,
	***REMOVED***

	ulimitOpt := NewUlimitOpt(&ulimitMap)

	expected := "[nofile=512:1024]"
	if ulimitOpt.String() != expected ***REMOVED***
		t.Fatalf("Expected %v, got %v", expected, ulimitOpt)
	***REMOVED***

	// Valid ulimit append to opts
	if err := ulimitOpt.Set("core=1024:1024"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Invalid ulimit type returns an error and do not append to opts
	if err := ulimitOpt.Set("notavalidtype=1024:1024"); err == nil ***REMOVED***
		t.Fatalf("Expected error on invalid ulimit type")
	***REMOVED***
	expected = "[nofile=512:1024 core=1024:1024]"
	expected2 := "[core=1024:1024 nofile=512:1024]"
	result := ulimitOpt.String()
	if result != expected && result != expected2 ***REMOVED***
		t.Fatalf("Expected %v or %v, got %v", expected, expected2, ulimitOpt)
	***REMOVED***

	// And test GetList
	ulimits := ulimitOpt.GetList()
	if len(ulimits) != 2 ***REMOVED***
		t.Fatalf("Expected a ulimit list of 2, got %v", ulimits)
	***REMOVED***
***REMOVED***
