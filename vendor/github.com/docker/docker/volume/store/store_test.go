package store

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/docker/docker/volume"
	"github.com/docker/docker/volume/drivers"
	volumetestutils "github.com/docker/docker/volume/testutils"
)

func TestCreate(t *testing.T) ***REMOVED***
	volumedrivers.Register(volumetestutils.NewFakeDriver("fake"), "fake")
	defer volumedrivers.Unregister("fake")
	dir, err := ioutil.TempDir("", "test-create")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(dir)

	s, err := New(dir)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	v, err := s.Create("fake1", "fake", nil, nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if v.Name() != "fake1" ***REMOVED***
		t.Fatalf("Expected fake1 volume, got %v", v)
	***REMOVED***
	if l, _, _ := s.List(); len(l) != 1 ***REMOVED***
		t.Fatalf("Expected 1 volume in the store, got %v: %v", len(l), l)
	***REMOVED***

	if _, err := s.Create("none", "none", nil, nil); err == nil ***REMOVED***
		t.Fatalf("Expected unknown driver error, got nil")
	***REMOVED***

	_, err = s.Create("fakeerror", "fake", map[string]string***REMOVED***"error": "create error"***REMOVED***, nil)
	expected := &OpErr***REMOVED***Op: "create", Name: "fakeerror", Err: errors.New("create error")***REMOVED***
	if err != nil && err.Error() != expected.Error() ***REMOVED***
		t.Fatalf("Expected create fakeError: create error, got %v", err)
	***REMOVED***
***REMOVED***

func TestRemove(t *testing.T) ***REMOVED***
	volumedrivers.Register(volumetestutils.NewFakeDriver("fake"), "fake")
	volumedrivers.Register(volumetestutils.NewFakeDriver("noop"), "noop")
	defer volumedrivers.Unregister("fake")
	defer volumedrivers.Unregister("noop")
	dir, err := ioutil.TempDir("", "test-remove")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(dir)
	s, err := New(dir)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// doing string compare here since this error comes directly from the driver
	expected := "no such volume"
	if err := s.Remove(volumetestutils.NoopVolume***REMOVED******REMOVED***); err == nil || !strings.Contains(err.Error(), expected) ***REMOVED***
		t.Fatalf("Expected error %q, got %v", expected, err)
	***REMOVED***

	v, err := s.CreateWithRef("fake1", "fake", "fake", nil, nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := s.Remove(v); !IsInUse(err) ***REMOVED***
		t.Fatalf("Expected ErrVolumeInUse error, got %v", err)
	***REMOVED***
	s.Dereference(v, "fake")
	if err := s.Remove(v); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if l, _, _ := s.List(); len(l) != 0 ***REMOVED***
		t.Fatalf("Expected 0 volumes in the store, got %v, %v", len(l), l)
	***REMOVED***
***REMOVED***

func TestList(t *testing.T) ***REMOVED***
	volumedrivers.Register(volumetestutils.NewFakeDriver("fake"), "fake")
	volumedrivers.Register(volumetestutils.NewFakeDriver("fake2"), "fake2")
	defer volumedrivers.Unregister("fake")
	defer volumedrivers.Unregister("fake2")
	dir, err := ioutil.TempDir("", "test-list")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(dir)

	s, err := New(dir)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := s.Create("test", "fake", nil, nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := s.Create("test2", "fake2", nil, nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	ls, _, err := s.List()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if len(ls) != 2 ***REMOVED***
		t.Fatalf("expected 2 volumes, got: %d", len(ls))
	***REMOVED***
	if err := s.Shutdown(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// and again with a new store
	s, err = New(dir)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	ls, _, err = s.List()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if len(ls) != 2 ***REMOVED***
		t.Fatalf("expected 2 volumes, got: %d", len(ls))
	***REMOVED***
***REMOVED***

func TestFilterByDriver(t *testing.T) ***REMOVED***
	volumedrivers.Register(volumetestutils.NewFakeDriver("fake"), "fake")
	volumedrivers.Register(volumetestutils.NewFakeDriver("noop"), "noop")
	defer volumedrivers.Unregister("fake")
	defer volumedrivers.Unregister("noop")
	dir, err := ioutil.TempDir("", "test-filter-driver")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	s, err := New(dir)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := s.Create("fake1", "fake", nil, nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := s.Create("fake2", "fake", nil, nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := s.Create("fake3", "noop", nil, nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if l, _ := s.FilterByDriver("fake"); len(l) != 2 ***REMOVED***
		t.Fatalf("Expected 2 volumes, got %v, %v", len(l), l)
	***REMOVED***

	if l, _ := s.FilterByDriver("noop"); len(l) != 1 ***REMOVED***
		t.Fatalf("Expected 1 volume, got %v, %v", len(l), l)
	***REMOVED***
***REMOVED***

func TestFilterByUsed(t *testing.T) ***REMOVED***
	volumedrivers.Register(volumetestutils.NewFakeDriver("fake"), "fake")
	volumedrivers.Register(volumetestutils.NewFakeDriver("noop"), "noop")
	dir, err := ioutil.TempDir("", "test-filter-used")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	s, err := New(dir)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := s.CreateWithRef("fake1", "fake", "volReference", nil, nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := s.Create("fake2", "fake", nil, nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	vols, _, err := s.List()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	dangling := s.FilterByUsed(vols, false)
	if len(dangling) != 1 ***REMOVED***
		t.Fatalf("expected 1 dangling volume, got %v", len(dangling))
	***REMOVED***
	if dangling[0].Name() != "fake2" ***REMOVED***
		t.Fatalf("expected dangling volume fake2, got %s", dangling[0].Name())
	***REMOVED***

	used := s.FilterByUsed(vols, true)
	if len(used) != 1 ***REMOVED***
		t.Fatalf("expected 1 used volume, got %v", len(used))
	***REMOVED***
	if used[0].Name() != "fake1" ***REMOVED***
		t.Fatalf("expected used volume fake1, got %s", used[0].Name())
	***REMOVED***
***REMOVED***

func TestDerefMultipleOfSameRef(t *testing.T) ***REMOVED***
	volumedrivers.Register(volumetestutils.NewFakeDriver("fake"), "fake")
	dir, err := ioutil.TempDir("", "test-same-deref")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(dir)
	s, err := New(dir)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	v, err := s.CreateWithRef("fake1", "fake", "volReference", nil, nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := s.GetWithRef("fake1", "fake", "volReference"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	s.Dereference(v, "volReference")
	if err := s.Remove(v); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestCreateKeepOptsLabelsWhenExistsRemotely(t *testing.T) ***REMOVED***
	vd := volumetestutils.NewFakeDriver("fake")
	volumedrivers.Register(vd, "fake")
	dir, err := ioutil.TempDir("", "test-same-deref")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(dir)
	s, err := New(dir)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Create a volume in the driver directly
	if _, err := vd.Create("foo", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	v, err := s.Create("foo", "fake", nil, map[string]string***REMOVED***"hello": "world"***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	switch dv := v.(type) ***REMOVED***
	case volume.DetailedVolume:
		if dv.Labels()["hello"] != "world" ***REMOVED***
			t.Fatalf("labels don't match")
		***REMOVED***
	default:
		t.Fatalf("got unexpected type: %T", v)
	***REMOVED***
***REMOVED***

func TestDefererencePluginOnCreateError(t *testing.T) ***REMOVED***
	var (
		l   net.Listener
		err error
	)

	for i := 32768; l == nil && i < 40000; i++ ***REMOVED***
		l, err = net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", i))
	***REMOVED***
	if l == nil ***REMOVED***
		t.Fatalf("could not create listener: %v", err)
	***REMOVED***
	defer l.Close()

	d := volumetestutils.NewFakeDriver("TestDefererencePluginOnCreateError")
	p, err := volumetestutils.MakeFakePlugin(d, l)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	pg := volumetestutils.NewFakePluginGetter(p)
	volumedrivers.RegisterPluginGetter(pg)

	dir, err := ioutil.TempDir("", "test-plugin-deref-err")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(dir)

	s, err := New(dir)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// create a good volume so we have a plugin reference
	_, err = s.Create("fake1", d.Name(), nil, nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Now create another one expecting an error
	_, err = s.Create("fake2", d.Name(), map[string]string***REMOVED***"error": "some error"***REMOVED***, nil)
	if err == nil || !strings.Contains(err.Error(), "some error") ***REMOVED***
		t.Fatalf("expected an error on create: %v", err)
	***REMOVED***

	// There should be only 1 plugin reference
	if refs := volumetestutils.FakeRefs(p); refs != 1 ***REMOVED***
		t.Fatalf("expected 1 plugin reference, got: %d", refs)
	***REMOVED***
***REMOVED***
