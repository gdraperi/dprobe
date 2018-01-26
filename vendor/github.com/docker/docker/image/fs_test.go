package image

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/docker/internal/testutil"
	digest "github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
)

func defaultFSStoreBackend(t *testing.T) (StoreBackend, func()) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "images-fs-store")
	assert.NoError(t, err)

	fsBackend, err := NewFSStoreBackend(tmpdir)
	assert.NoError(t, err)

	return fsBackend, func() ***REMOVED*** os.RemoveAll(tmpdir) ***REMOVED***
***REMOVED***

func TestFSGetInvalidData(t *testing.T) ***REMOVED***
	store, cleanup := defaultFSStoreBackend(t)
	defer cleanup()

	id, err := store.Set([]byte("foobar"))
	assert.NoError(t, err)

	dgst := digest.Digest(id)

	err = ioutil.WriteFile(filepath.Join(store.(*fs).root, contentDirName, string(dgst.Algorithm()), dgst.Hex()), []byte("foobar2"), 0600)
	assert.NoError(t, err)

	_, err = store.Get(id)
	testutil.ErrorContains(t, err, "failed to verify")
***REMOVED***

func TestFSInvalidSet(t *testing.T) ***REMOVED***
	store, cleanup := defaultFSStoreBackend(t)
	defer cleanup()

	id := digest.FromBytes([]byte("foobar"))
	err := os.Mkdir(filepath.Join(store.(*fs).root, contentDirName, string(id.Algorithm()), id.Hex()), 0700)
	assert.NoError(t, err)

	_, err = store.Set([]byte("foobar"))
	testutil.ErrorContains(t, err, "failed to write digest data")
***REMOVED***

func TestFSInvalidRoot(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "images-fs-store")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	tcases := []struct ***REMOVED***
		root, invalidFile string
	***REMOVED******REMOVED***
		***REMOVED***"root", "root"***REMOVED***,
		***REMOVED***"root", "root/content"***REMOVED***,
		***REMOVED***"root", "root/metadata"***REMOVED***,
	***REMOVED***

	for _, tc := range tcases ***REMOVED***
		root := filepath.Join(tmpdir, tc.root)
		filePath := filepath.Join(tmpdir, tc.invalidFile)
		err := os.MkdirAll(filepath.Dir(filePath), 0700)
		assert.NoError(t, err)

		f, err := os.Create(filePath)
		assert.NoError(t, err)
		f.Close()

		_, err = NewFSStoreBackend(root)
		testutil.ErrorContains(t, err, "failed to create storage backend")

		os.RemoveAll(root)
	***REMOVED***

***REMOVED***

func TestFSMetadataGetSet(t *testing.T) ***REMOVED***
	store, cleanup := defaultFSStoreBackend(t)
	defer cleanup()

	id, err := store.Set([]byte("foo"))
	assert.NoError(t, err)

	id2, err := store.Set([]byte("bar"))
	assert.NoError(t, err)

	tcases := []struct ***REMOVED***
		id    digest.Digest
		key   string
		value []byte
	***REMOVED******REMOVED***
		***REMOVED***id, "tkey", []byte("tval1")***REMOVED***,
		***REMOVED***id, "tkey2", []byte("tval2")***REMOVED***,
		***REMOVED***id2, "tkey", []byte("tval3")***REMOVED***,
	***REMOVED***

	for _, tc := range tcases ***REMOVED***
		err = store.SetMetadata(tc.id, tc.key, tc.value)
		assert.NoError(t, err)

		actual, err := store.GetMetadata(tc.id, tc.key)
		assert.NoError(t, err)

		assert.Equal(t, tc.value, actual)
	***REMOVED***

	_, err = store.GetMetadata(id2, "tkey2")
	testutil.ErrorContains(t, err, "failed to read metadata")

	id3 := digest.FromBytes([]byte("baz"))
	err = store.SetMetadata(id3, "tkey", []byte("tval"))
	testutil.ErrorContains(t, err, "failed to get digest")

	_, err = store.GetMetadata(id3, "tkey")
	testutil.ErrorContains(t, err, "failed to get digest")
***REMOVED***

func TestFSInvalidWalker(t *testing.T) ***REMOVED***
	store, cleanup := defaultFSStoreBackend(t)
	defer cleanup()

	fooID, err := store.Set([]byte("foo"))
	assert.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(store.(*fs).root, contentDirName, "sha256/foobar"), []byte("foobar"), 0600)
	assert.NoError(t, err)

	n := 0
	err = store.Walk(func(id digest.Digest) error ***REMOVED***
		assert.Equal(t, fooID, id)
		n++
		return nil
	***REMOVED***)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
***REMOVED***

func TestFSGetSet(t *testing.T) ***REMOVED***
	store, cleanup := defaultFSStoreBackend(t)
	defer cleanup()

	type tcase struct ***REMOVED***
		input    []byte
		expected digest.Digest
	***REMOVED***
	tcases := []tcase***REMOVED***
		***REMOVED***[]byte("foobar"), digest.Digest("sha256:c3ab8ff13720e8ad9047dd39466b3c8974e592c2fa383d4a3960714caef0c4f2")***REMOVED***,
	***REMOVED***

	randomInput := make([]byte, 8*1024)
	_, err := rand.Read(randomInput)
	assert.NoError(t, err)

	// skipping use of digest pkg because it is used by the implementation
	h := sha256.New()
	_, err = h.Write(randomInput)
	assert.NoError(t, err)

	tcases = append(tcases, tcase***REMOVED***
		input:    randomInput,
		expected: digest.Digest("sha256:" + hex.EncodeToString(h.Sum(nil))),
	***REMOVED***)

	for _, tc := range tcases ***REMOVED***
		id, err := store.Set([]byte(tc.input))
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, id)
	***REMOVED***

	for _, tc := range tcases ***REMOVED***
		data, err := store.Get(tc.expected)
		assert.NoError(t, err)
		assert.Equal(t, tc.input, data)
	***REMOVED***
***REMOVED***

func TestFSGetUnsetKey(t *testing.T) ***REMOVED***
	store, cleanup := defaultFSStoreBackend(t)
	defer cleanup()

	for _, key := range []digest.Digest***REMOVED***"foobar:abc", "sha256:abc", "sha256:c3ab8ff13720e8ad9047dd39466b3c8974e592c2fa383d4a3960714caef0c4f2a"***REMOVED*** ***REMOVED***
		_, err := store.Get(key)
		testutil.ErrorContains(t, err, "failed to get digest")
	***REMOVED***
***REMOVED***

func TestFSGetEmptyData(t *testing.T) ***REMOVED***
	store, cleanup := defaultFSStoreBackend(t)
	defer cleanup()

	for _, emptyData := range [][]byte***REMOVED***nil, ***REMOVED******REMOVED******REMOVED*** ***REMOVED***
		_, err := store.Set(emptyData)
		testutil.ErrorContains(t, err, "invalid empty data")
	***REMOVED***
***REMOVED***

func TestFSDelete(t *testing.T) ***REMOVED***
	store, cleanup := defaultFSStoreBackend(t)
	defer cleanup()

	id, err := store.Set([]byte("foo"))
	assert.NoError(t, err)

	id2, err := store.Set([]byte("bar"))
	assert.NoError(t, err)

	err = store.Delete(id)
	assert.NoError(t, err)

	_, err = store.Get(id)
	testutil.ErrorContains(t, err, "failed to get digest")

	_, err = store.Get(id2)
	assert.NoError(t, err)

	err = store.Delete(id2)
	assert.NoError(t, err)

	_, err = store.Get(id2)
	testutil.ErrorContains(t, err, "failed to get digest")
***REMOVED***

func TestFSWalker(t *testing.T) ***REMOVED***
	store, cleanup := defaultFSStoreBackend(t)
	defer cleanup()

	id, err := store.Set([]byte("foo"))
	assert.NoError(t, err)

	id2, err := store.Set([]byte("bar"))
	assert.NoError(t, err)

	tcases := make(map[digest.Digest]struct***REMOVED******REMOVED***)
	tcases[id] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	tcases[id2] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	n := 0
	err = store.Walk(func(id digest.Digest) error ***REMOVED***
		delete(tcases, id)
		n++
		return nil
	***REMOVED***)
	assert.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.Len(t, tcases, 0)
***REMOVED***

func TestFSWalkerStopOnError(t *testing.T) ***REMOVED***
	store, cleanup := defaultFSStoreBackend(t)
	defer cleanup()

	id, err := store.Set([]byte("foo"))
	assert.NoError(t, err)

	tcases := make(map[digest.Digest]struct***REMOVED******REMOVED***)
	tcases[id] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	err = store.Walk(func(id digest.Digest) error ***REMOVED***
		return errors.New("what")
	***REMOVED***)
	testutil.ErrorContains(t, err, "what")
***REMOVED***
