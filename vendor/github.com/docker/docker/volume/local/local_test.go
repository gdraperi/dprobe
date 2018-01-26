package local

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/mount"
)

func TestGetAddress(t *testing.T) ***REMOVED***
	cases := map[string]string***REMOVED***
		"addr=11.11.11.1":   "11.11.11.1",
		" ":                 "",
		"addr=":             "",
		"addr=2001:db8::68": "2001:db8::68",
	***REMOVED***
	for name, success := range cases ***REMOVED***
		v := getAddress(name)
		if v != success ***REMOVED***
			t.Errorf("Test case failed for %s actual: %s expected : %s", name, v, success)
		***REMOVED***
	***REMOVED***

***REMOVED***

func TestRemove(t *testing.T) ***REMOVED***
	// TODO Windows: Investigate why this test fails on Windows under CI
	//               but passes locally.
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Test failing on Windows CI")
	***REMOVED***
	rootDir, err := ioutil.TempDir("", "local-volume-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(rootDir)

	r, err := New(rootDir, idtools.IDPair***REMOVED***UID: 0, GID: 0***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	vol, err := r.Create("testing", nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := r.Remove(vol); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	vol, err = r.Create("testing2", nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := os.RemoveAll(vol.Path()); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := r.Remove(vol); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := os.Stat(vol.Path()); err != nil && !os.IsNotExist(err) ***REMOVED***
		t.Fatal("volume dir not removed")
	***REMOVED***

	if l, _ := r.List(); len(l) != 0 ***REMOVED***
		t.Fatal("expected there to be no volumes")
	***REMOVED***
***REMOVED***

func TestInitializeWithVolumes(t *testing.T) ***REMOVED***
	rootDir, err := ioutil.TempDir("", "local-volume-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(rootDir)

	r, err := New(rootDir, idtools.IDPair***REMOVED***UID: 0, GID: 0***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	vol, err := r.Create("testing", nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	r, err = New(rootDir, idtools.IDPair***REMOVED***UID: 0, GID: 0***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	v, err := r.Get(vol.Name())
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if v.Path() != vol.Path() ***REMOVED***
		t.Fatal("expected to re-initialize root with existing volumes")
	***REMOVED***
***REMOVED***

func TestCreate(t *testing.T) ***REMOVED***
	rootDir, err := ioutil.TempDir("", "local-volume-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(rootDir)

	r, err := New(rootDir, idtools.IDPair***REMOVED***UID: 0, GID: 0***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	cases := map[string]bool***REMOVED***
		"name":                  true,
		"name-with-dash":        true,
		"name_with_underscore":  true,
		"name/with/slash":       false,
		"name/with/../../slash": false,
		"./name":                false,
		"../name":               false,
		"./":                    false,
		"../":                   false,
		"~":                     false,
		".":                     false,
		"..":                    false,
		"...":                   false,
	***REMOVED***

	for name, success := range cases ***REMOVED***
		v, err := r.Create(name, nil)
		if success ***REMOVED***
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			if v.Name() != name ***REMOVED***
				t.Fatalf("Expected volume with name %s, got %s", name, v.Name())
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if err == nil ***REMOVED***
				t.Fatalf("Expected error creating volume with name %s, got nil", name)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	r, err = New(rootDir, idtools.IDPair***REMOVED***UID: 0, GID: 0***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestValidateName(t *testing.T) ***REMOVED***
	r := &Root***REMOVED******REMOVED***
	names := map[string]bool***REMOVED***
		"x":           false,
		"/testvol":    false,
		"thing.d":     true,
		"hello-world": true,
		"./hello":     false,
		".hello":      false,
	***REMOVED***

	for vol, expected := range names ***REMOVED***
		err := r.validateName(vol)
		if expected && err != nil ***REMOVED***
			t.Fatalf("expected %s to be valid got %v", vol, err)
		***REMOVED***
		if !expected && err == nil ***REMOVED***
			t.Fatalf("expected %s to be invalid", vol)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestCreateWithOpts(t *testing.T) ***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip()
	***REMOVED***
	rootDir, err := ioutil.TempDir("", "local-volume-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(rootDir)

	r, err := New(rootDir, idtools.IDPair***REMOVED***UID: 0, GID: 0***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := r.Create("test", map[string]string***REMOVED***"invalidopt": "notsupported"***REMOVED***); err == nil ***REMOVED***
		t.Fatal("expected invalid opt to cause error")
	***REMOVED***

	vol, err := r.Create("test", map[string]string***REMOVED***"device": "tmpfs", "type": "tmpfs", "o": "size=1m,uid=1000"***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	v := vol.(*localVolume)

	dir, err := v.Mount("1234")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer func() ***REMOVED***
		if err := v.Unmount("1234"); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***()

	mountInfos, err := mount.GetMounts()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	var found bool
	for _, info := range mountInfos ***REMOVED***
		if info.Mountpoint == dir ***REMOVED***
			found = true
			if info.Fstype != "tmpfs" ***REMOVED***
				t.Fatalf("expected tmpfs mount, got %q", info.Fstype)
			***REMOVED***
			if info.Source != "tmpfs" ***REMOVED***
				t.Fatalf("expected tmpfs mount, got %q", info.Source)
			***REMOVED***
			if !strings.Contains(info.VfsOpts, "uid=1000") ***REMOVED***
				t.Fatalf("expected mount info to have uid=1000: %q", info.VfsOpts)
			***REMOVED***
			if !strings.Contains(info.VfsOpts, "size=1024k") ***REMOVED***
				t.Fatalf("expected mount info to have size=1024k: %q", info.VfsOpts)
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	if !found ***REMOVED***
		t.Fatal("mount not found")
	***REMOVED***

	if v.active.count != 1 ***REMOVED***
		t.Fatalf("Expected active mount count to be 1, got %d", v.active.count)
	***REMOVED***

	// test double mount
	if _, err := v.Mount("1234"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if v.active.count != 2 ***REMOVED***
		t.Fatalf("Expected active mount count to be 2, got %d", v.active.count)
	***REMOVED***

	if err := v.Unmount("1234"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if v.active.count != 1 ***REMOVED***
		t.Fatalf("Expected active mount count to be 1, got %d", v.active.count)
	***REMOVED***

	mounted, err := mount.Mounted(v.path)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !mounted ***REMOVED***
		t.Fatal("expected mount to still be active")
	***REMOVED***

	r, err = New(rootDir, idtools.IDPair***REMOVED***UID: 0, GID: 0***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	v2, exists := r.volumes["test"]
	if !exists ***REMOVED***
		t.Fatal("missing volume on restart")
	***REMOVED***

	if !reflect.DeepEqual(v.opts, v2.opts) ***REMOVED***
		t.Fatal("missing volume options on restart")
	***REMOVED***
***REMOVED***

func TestRealodNoOpts(t *testing.T) ***REMOVED***
	rootDir, err := ioutil.TempDir("", "volume-test-reload-no-opts")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(rootDir)

	r, err := New(rootDir, idtools.IDPair***REMOVED***UID: 0, GID: 0***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := r.Create("test1", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := r.Create("test2", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	// make sure a file with `null` (.e.g. empty opts map from older daemon) is ok
	if err := ioutil.WriteFile(filepath.Join(rootDir, "test2"), []byte("null"), 600); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := r.Create("test3", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	// make sure an empty opts file doesn't break us too
	if err := ioutil.WriteFile(filepath.Join(rootDir, "test3"), nil, 600); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := r.Create("test4", map[string]string***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	r, err = New(rootDir, idtools.IDPair***REMOVED***UID: 0, GID: 0***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	for _, name := range []string***REMOVED***"test1", "test2", "test3", "test4"***REMOVED*** ***REMOVED***
		v, err := r.Get(name)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		lv, ok := v.(*localVolume)
		if !ok ***REMOVED***
			t.Fatalf("expected *localVolume got: %v", reflect.TypeOf(v))
		***REMOVED***
		if lv.opts != nil ***REMOVED***
			t.Fatalf("expected opts to be nil, got: %v", lv.opts)
		***REMOVED***
		if _, err := lv.Mount("1234"); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***
