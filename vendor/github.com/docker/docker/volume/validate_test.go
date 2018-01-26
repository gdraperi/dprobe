package volume

import (
	"errors"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/mount"
)

func TestValidateMount(t *testing.T) ***REMOVED***
	testDir, err := ioutil.TempDir("", "test-validate-mount")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(testDir)

	cases := []struct ***REMOVED***
		input    mount.Mount
		expected error
	***REMOVED******REMOVED***
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeVolume***REMOVED***, errMissingField("Target")***REMOVED***,
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeVolume, Target: testDestinationPath, Source: "hello"***REMOVED***, nil***REMOVED***,
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeVolume, Target: testDestinationPath***REMOVED***, nil***REMOVED***,
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeBind***REMOVED***, errMissingField("Target")***REMOVED***,
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeBind, Target: testDestinationPath***REMOVED***, errMissingField("Source")***REMOVED***,
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeBind, Target: testDestinationPath, Source: testSourcePath, VolumeOptions: &mount.VolumeOptions***REMOVED******REMOVED******REMOVED***, errExtraField("VolumeOptions")***REMOVED***,

		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeBind, Source: testDir, Target: testDestinationPath***REMOVED***, nil***REMOVED***,
		***REMOVED***mount.Mount***REMOVED***Type: "invalid", Target: testDestinationPath***REMOVED***, errors.New("mount type unknown")***REMOVED***,
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeBind, Source: testSourcePath, Target: testDestinationPath***REMOVED***, errBindNotExist***REMOVED***,
	***REMOVED***

	lcowCases := []struct ***REMOVED***
		input    mount.Mount
		expected error
	***REMOVED******REMOVED***
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeVolume***REMOVED***, errMissingField("Target")***REMOVED***,
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeVolume, Target: "/foo", Source: "hello"***REMOVED***, nil***REMOVED***,
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeVolume, Target: "/foo"***REMOVED***, nil***REMOVED***,
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeBind***REMOVED***, errMissingField("Target")***REMOVED***,
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeBind, Target: "/foo"***REMOVED***, errMissingField("Source")***REMOVED***,
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeBind, Target: "/foo", Source: "c:\\foo", VolumeOptions: &mount.VolumeOptions***REMOVED******REMOVED******REMOVED***, errExtraField("VolumeOptions")***REMOVED***,
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeBind, Source: "c:\\foo", Target: "/foo"***REMOVED***, errBindNotExist***REMOVED***,
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeBind, Source: testDir, Target: "/foo"***REMOVED***, nil***REMOVED***,
		***REMOVED***mount.Mount***REMOVED***Type: "invalid", Target: "/foo"***REMOVED***, errors.New("mount type unknown")***REMOVED***,
	***REMOVED***
	parser := NewParser(runtime.GOOS)
	for i, x := range cases ***REMOVED***
		err := parser.ValidateMountConfig(&x.input)
		if err == nil && x.expected == nil ***REMOVED***
			continue
		***REMOVED***
		if (err == nil && x.expected != nil) || (x.expected == nil && err != nil) || !strings.Contains(err.Error(), x.expected.Error()) ***REMOVED***
			t.Errorf("expected %q, got %q, case: %d", x.expected, err, i)
		***REMOVED***
	***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		parser = &lcowParser***REMOVED******REMOVED***
		for i, x := range lcowCases ***REMOVED***
			err := parser.ValidateMountConfig(&x.input)
			if err == nil && x.expected == nil ***REMOVED***
				continue
			***REMOVED***
			if (err == nil && x.expected != nil) || (x.expected == nil && err != nil) || !strings.Contains(err.Error(), x.expected.Error()) ***REMOVED***
				t.Errorf("expected %q, got %q, case: %d", x.expected, err, i)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
