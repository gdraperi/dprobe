package fscache

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/moby/buildkit/session/filesync"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestFSCache(t *testing.T) ***REMOVED***
	tmpDir, err := ioutil.TempDir("", "fscache")
	assert.Nil(t, err)
	defer os.RemoveAll(tmpDir)

	backend := NewNaiveCacheBackend(filepath.Join(tmpDir, "backend"))

	opt := Opt***REMOVED***
		Root:     tmpDir,
		Backend:  backend,
		GCPolicy: GCPolicy***REMOVED***MaxSize: 15, MaxKeepDuration: time.Hour***REMOVED***,
	***REMOVED***

	fscache, err := NewFSCache(opt)
	assert.Nil(t, err)

	defer fscache.Close()

	err = fscache.RegisterTransport("test", &testTransport***REMOVED******REMOVED***)
	assert.Nil(t, err)

	src1, err := fscache.SyncFrom(context.TODO(), &testIdentifier***REMOVED***"foo", "data", "bar"***REMOVED***)
	assert.Nil(t, err)

	dt, err := ioutil.ReadFile(filepath.Join(src1.Root().Path(), "foo"))
	assert.Nil(t, err)
	assert.Equal(t, string(dt), "data")

	// same id doesn't recalculate anything
	src2, err := fscache.SyncFrom(context.TODO(), &testIdentifier***REMOVED***"foo", "data2", "bar"***REMOVED***)
	assert.Nil(t, err)
	assert.Equal(t, src1.Root().Path(), src2.Root().Path())

	dt, err = ioutil.ReadFile(filepath.Join(src1.Root().Path(), "foo"))
	assert.Nil(t, err)
	assert.Equal(t, string(dt), "data")
	assert.Nil(t, src2.Close())

	src3, err := fscache.SyncFrom(context.TODO(), &testIdentifier***REMOVED***"foo2", "data2", "bar"***REMOVED***)
	assert.Nil(t, err)
	assert.NotEqual(t, src1.Root().Path(), src3.Root().Path())

	dt, err = ioutil.ReadFile(filepath.Join(src3.Root().Path(), "foo2"))
	assert.Nil(t, err)
	assert.Equal(t, string(dt), "data2")

	s, err := fscache.DiskUsage()
	assert.Nil(t, err)
	assert.Equal(t, s, int64(0))

	assert.Nil(t, src3.Close())

	s, err = fscache.DiskUsage()
	assert.Nil(t, err)
	assert.Equal(t, s, int64(5))

	// new upload with the same shared key shoutl overwrite
	src4, err := fscache.SyncFrom(context.TODO(), &testIdentifier***REMOVED***"foo3", "data3", "bar"***REMOVED***)
	assert.Nil(t, err)
	assert.NotEqual(t, src1.Root().Path(), src3.Root().Path())

	dt, err = ioutil.ReadFile(filepath.Join(src3.Root().Path(), "foo3"))
	assert.Nil(t, err)
	assert.Equal(t, string(dt), "data3")
	assert.Equal(t, src4.Root().Path(), src3.Root().Path())
	assert.Nil(t, src4.Close())

	s, err = fscache.DiskUsage()
	assert.Nil(t, err)
	assert.Equal(t, s, int64(10))

	// this one goes over the GC limit
	src5, err := fscache.SyncFrom(context.TODO(), &testIdentifier***REMOVED***"foo4", "datadata", "baz"***REMOVED***)
	assert.Nil(t, err)
	assert.Nil(t, src5.Close())

	// GC happens async
	time.Sleep(100 * time.Millisecond)

	// only last insertion after GC
	s, err = fscache.DiskUsage()
	assert.Nil(t, err)
	assert.Equal(t, s, int64(8))

	// prune deletes everything
	released, err := fscache.Prune(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, released, uint64(8))

	s, err = fscache.DiskUsage()
	assert.Nil(t, err)
	assert.Equal(t, s, int64(0))
***REMOVED***

type testTransport struct ***REMOVED***
***REMOVED***

func (t *testTransport) Copy(ctx context.Context, id RemoteIdentifier, dest string, cs filesync.CacheUpdater) error ***REMOVED***
	testid := id.(*testIdentifier)
	return ioutil.WriteFile(filepath.Join(dest, testid.filename), []byte(testid.data), 0600)
***REMOVED***

type testIdentifier struct ***REMOVED***
	filename  string
	data      string
	sharedKey string
***REMOVED***

func (t *testIdentifier) Key() string ***REMOVED***
	return t.filename
***REMOVED***
func (t *testIdentifier) SharedKey() string ***REMOVED***
	return t.sharedKey
***REMOVED***
func (t *testIdentifier) Transport() string ***REMOVED***
	return "test"
***REMOVED***
