package tarsum

import (
	"testing"
)

func TestVersionLabelForChecksum(t *testing.T) ***REMOVED***
	version := VersionLabelForChecksum("tarsum+sha256:deadbeef")
	if version != "tarsum" ***REMOVED***
		t.Fatalf("Version should have been 'tarsum', was %v", version)
	***REMOVED***
	version = VersionLabelForChecksum("tarsum.v1+sha256:deadbeef")
	if version != "tarsum.v1" ***REMOVED***
		t.Fatalf("Version should have been 'tarsum.v1', was %v", version)
	***REMOVED***
	version = VersionLabelForChecksum("something+somethingelse")
	if version != "something" ***REMOVED***
		t.Fatalf("Version should have been 'something', was %v", version)
	***REMOVED***
	version = VersionLabelForChecksum("invalidChecksum")
	if version != "" ***REMOVED***
		t.Fatalf("Version should have been empty, was %v", version)
	***REMOVED***
***REMOVED***

func TestVersion(t *testing.T) ***REMOVED***
	expected := "tarsum"
	var v Version
	if v.String() != expected ***REMOVED***
		t.Errorf("expected %q, got %q", expected, v.String())
	***REMOVED***

	expected = "tarsum.v1"
	v = 1
	if v.String() != expected ***REMOVED***
		t.Errorf("expected %q, got %q", expected, v.String())
	***REMOVED***

	expected = "tarsum.dev"
	v = 2
	if v.String() != expected ***REMOVED***
		t.Errorf("expected %q, got %q", expected, v.String())
	***REMOVED***
***REMOVED***

func TestGetVersion(t *testing.T) ***REMOVED***
	testSet := []struct ***REMOVED***
		Str      string
		Expected Version
	***REMOVED******REMOVED***
		***REMOVED***"tarsum+sha256:e58fcf7418d4390dec8e8fb69d88c06ec07039d651fedd3aa72af9972e7d046b", Version0***REMOVED***,
		***REMOVED***"tarsum+sha256", Version0***REMOVED***,
		***REMOVED***"tarsum", Version0***REMOVED***,
		***REMOVED***"tarsum.dev", VersionDev***REMOVED***,
		***REMOVED***"tarsum.dev+sha256:deadbeef", VersionDev***REMOVED***,
	***REMOVED***

	for _, ts := range testSet ***REMOVED***
		v, err := GetVersionFromTarsum(ts.Str)
		if err != nil ***REMOVED***
			t.Fatalf("%q : %s", err, ts.Str)
		***REMOVED***
		if v != ts.Expected ***REMOVED***
			t.Errorf("expected %d (%q), got %d (%q)", ts.Expected, ts.Expected, v, v)
		***REMOVED***
	***REMOVED***

	// test one that does not exist, to ensure it errors
	str := "weak+md5:abcdeabcde"
	_, err := GetVersionFromTarsum(str)
	if err != ErrNotVersion ***REMOVED***
		t.Fatalf("%q : %s", err, str)
	***REMOVED***
***REMOVED***

func TestGetVersions(t *testing.T) ***REMOVED***
	expected := []Version***REMOVED***
		Version0,
		Version1,
		VersionDev,
	***REMOVED***
	versions := GetVersions()
	if len(versions) != len(expected) ***REMOVED***
		t.Fatalf("Expected %v versions, got %v", len(expected), len(versions))
	***REMOVED***
	if !containsVersion(versions, expected[0]) || !containsVersion(versions, expected[1]) || !containsVersion(versions, expected[2]) ***REMOVED***
		t.Fatalf("Expected [%v], got [%v]", expected, versions)
	***REMOVED***
***REMOVED***

func containsVersion(versions []Version, version Version) bool ***REMOVED***
	for _, v := range versions ***REMOVED***
		if v == version ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
