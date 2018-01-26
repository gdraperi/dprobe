package volumedrivers

import (
	"testing"

	volumetestutils "github.com/docker/docker/volume/testutils"
)

func TestGetDriver(t *testing.T) ***REMOVED***
	_, err := GetDriver("missing")
	if err == nil ***REMOVED***
		t.Fatal("Expected error, was nil")
	***REMOVED***
	Register(volumetestutils.NewFakeDriver("fake"), "fake")

	d, err := GetDriver("fake")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if d.Name() != "fake" ***REMOVED***
		t.Fatalf("Expected fake driver, got %s\n", d.Name())
	***REMOVED***
***REMOVED***
