package reference

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/distribution/reference"
	digest "github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	saveLoadTestCases = map[string]digest.Digest***REMOVED***
		"registry:5000/foobar:HEAD":                                                        "sha256:470022b8af682154f57a2163d030eb369549549cba00edc69e1b99b46bb924d6",
		"registry:5000/foobar:alternate":                                                   "sha256:ae300ebc4a4f00693702cfb0a5e0b7bc527b353828dc86ad09fb95c8a681b793",
		"registry:5000/foobar:latest":                                                      "sha256:6153498b9ac00968d71b66cca4eac37e990b5f9eb50c26877eb8799c8847451b",
		"registry:5000/foobar:master":                                                      "sha256:6c9917af4c4e05001b346421959d7ea81b6dc9d25718466a37a6add865dfd7fc",
		"jess/hollywood:latest":                                                            "sha256:ae7a5519a0a55a2d4ef20ddcbd5d0ca0888a1f7ab806acc8e2a27baf46f529fe",
		"registry@sha256:367eb40fd0330a7e464777121e39d2f5b3e8e23a1e159342e53ab05c9e4d94e6": "sha256:24126a56805beb9711be5f4590cc2eb55ab8d4a85ebd618eed72bb19fc50631c",
		"busybox:latest": "sha256:91e54dfb11794fad694460162bf0cb0a4fa710cfa3f60979c177d920813e267c",
	***REMOVED***

	marshalledSaveLoadTestCases = []byte(`***REMOVED***"Repositories":***REMOVED***"busybox":***REMOVED***"busybox:latest":"sha256:91e54dfb11794fad694460162bf0cb0a4fa710cfa3f60979c177d920813e267c"***REMOVED***,"jess/hollywood":***REMOVED***"jess/hollywood:latest":"sha256:ae7a5519a0a55a2d4ef20ddcbd5d0ca0888a1f7ab806acc8e2a27baf46f529fe"***REMOVED***,"registry":***REMOVED***"registry@sha256:367eb40fd0330a7e464777121e39d2f5b3e8e23a1e159342e53ab05c9e4d94e6":"sha256:24126a56805beb9711be5f4590cc2eb55ab8d4a85ebd618eed72bb19fc50631c"***REMOVED***,"registry:5000/foobar":***REMOVED***"registry:5000/foobar:HEAD":"sha256:470022b8af682154f57a2163d030eb369549549cba00edc69e1b99b46bb924d6","registry:5000/foobar:alternate":"sha256:ae300ebc4a4f00693702cfb0a5e0b7bc527b353828dc86ad09fb95c8a681b793","registry:5000/foobar:latest":"sha256:6153498b9ac00968d71b66cca4eac37e990b5f9eb50c26877eb8799c8847451b","registry:5000/foobar:master":"sha256:6c9917af4c4e05001b346421959d7ea81b6dc9d25718466a37a6add865dfd7fc"***REMOVED******REMOVED******REMOVED***`)
)

func TestLoad(t *testing.T) ***REMOVED***
	jsonFile, err := ioutil.TempFile("", "tag-store-test")
	if err != nil ***REMOVED***
		t.Fatalf("error creating temp file: %v", err)
	***REMOVED***
	defer os.RemoveAll(jsonFile.Name())

	// Write canned json to the temp file
	_, err = jsonFile.Write(marshalledSaveLoadTestCases)
	if err != nil ***REMOVED***
		t.Fatalf("error writing to temp file: %v", err)
	***REMOVED***
	jsonFile.Close()

	store, err := NewReferenceStore(jsonFile.Name())
	if err != nil ***REMOVED***
		t.Fatalf("error creating tag store: %v", err)
	***REMOVED***

	for refStr, expectedID := range saveLoadTestCases ***REMOVED***
		ref, err := reference.ParseNormalizedNamed(refStr)
		if err != nil ***REMOVED***
			t.Fatalf("failed to parse reference: %v", err)
		***REMOVED***
		id, err := store.Get(ref)
		if err != nil ***REMOVED***
			t.Fatalf("could not find reference %s: %v", refStr, err)
		***REMOVED***
		if id != expectedID ***REMOVED***
			t.Fatalf("expected %s - got %s", expectedID, id)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSave(t *testing.T) ***REMOVED***
	jsonFile, err := ioutil.TempFile("", "tag-store-test")
	require.NoError(t, err)

	_, err = jsonFile.Write([]byte(`***REMOVED******REMOVED***`))
	require.NoError(t, err)
	jsonFile.Close()
	defer os.RemoveAll(jsonFile.Name())

	store, err := NewReferenceStore(jsonFile.Name())
	if err != nil ***REMOVED***
		t.Fatalf("error creating tag store: %v", err)
	***REMOVED***

	for refStr, id := range saveLoadTestCases ***REMOVED***
		ref, err := reference.ParseNormalizedNamed(refStr)
		if err != nil ***REMOVED***
			t.Fatalf("failed to parse reference: %v", err)
		***REMOVED***
		if canonical, ok := ref.(reference.Canonical); ok ***REMOVED***
			err = store.AddDigest(canonical, id, false)
			if err != nil ***REMOVED***
				t.Fatalf("could not add digest reference %s: %v", refStr, err)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			err = store.AddTag(ref, id, false)
			if err != nil ***REMOVED***
				t.Fatalf("could not add reference %s: %v", refStr, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	jsonBytes, err := ioutil.ReadFile(jsonFile.Name())
	if err != nil ***REMOVED***
		t.Fatalf("could not read json file: %v", err)
	***REMOVED***

	if !bytes.Equal(jsonBytes, marshalledSaveLoadTestCases) ***REMOVED***
		t.Fatalf("save output did not match expectations\nexpected:\n%s\ngot:\n%s", marshalledSaveLoadTestCases, jsonBytes)
	***REMOVED***
***REMOVED***

func TestAddDeleteGet(t *testing.T) ***REMOVED***
	jsonFile, err := ioutil.TempFile("", "tag-store-test")
	if err != nil ***REMOVED***
		t.Fatalf("error creating temp file: %v", err)
	***REMOVED***
	_, err = jsonFile.Write([]byte(`***REMOVED******REMOVED***`))
	jsonFile.Close()
	defer os.RemoveAll(jsonFile.Name())

	store, err := NewReferenceStore(jsonFile.Name())
	if err != nil ***REMOVED***
		t.Fatalf("error creating tag store: %v", err)
	***REMOVED***

	testImageID1 := digest.Digest("sha256:9655aef5fd742a1b4e1b7b163aa9f1c76c186304bf39102283d80927c916ca9c")
	testImageID2 := digest.Digest("sha256:9655aef5fd742a1b4e1b7b163aa9f1c76c186304bf39102283d80927c916ca9d")
	testImageID3 := digest.Digest("sha256:9655aef5fd742a1b4e1b7b163aa9f1c76c186304bf39102283d80927c916ca9e")

	// Try adding a reference with no tag or digest
	nameOnly, err := reference.ParseNormalizedNamed("username/repo")
	if err != nil ***REMOVED***
		t.Fatalf("could not parse reference: %v", err)
	***REMOVED***
	if err = store.AddTag(nameOnly, testImageID1, false); err != nil ***REMOVED***
		t.Fatalf("error adding to store: %v", err)
	***REMOVED***

	// Add a few references
	ref1, err := reference.ParseNormalizedNamed("username/repo1:latest")
	if err != nil ***REMOVED***
		t.Fatalf("could not parse reference: %v", err)
	***REMOVED***
	if err = store.AddTag(ref1, testImageID1, false); err != nil ***REMOVED***
		t.Fatalf("error adding to store: %v", err)
	***REMOVED***

	ref2, err := reference.ParseNormalizedNamed("username/repo1:old")
	if err != nil ***REMOVED***
		t.Fatalf("could not parse reference: %v", err)
	***REMOVED***
	if err = store.AddTag(ref2, testImageID2, false); err != nil ***REMOVED***
		t.Fatalf("error adding to store: %v", err)
	***REMOVED***

	ref3, err := reference.ParseNormalizedNamed("username/repo1:alias")
	if err != nil ***REMOVED***
		t.Fatalf("could not parse reference: %v", err)
	***REMOVED***
	if err = store.AddTag(ref3, testImageID1, false); err != nil ***REMOVED***
		t.Fatalf("error adding to store: %v", err)
	***REMOVED***

	ref4, err := reference.ParseNormalizedNamed("username/repo2:latest")
	if err != nil ***REMOVED***
		t.Fatalf("could not parse reference: %v", err)
	***REMOVED***
	if err = store.AddTag(ref4, testImageID2, false); err != nil ***REMOVED***
		t.Fatalf("error adding to store: %v", err)
	***REMOVED***

	ref5, err := reference.ParseNormalizedNamed("username/repo3@sha256:58153dfb11794fad694460162bf0cb0a4fa710cfa3f60979c177d920813e267c")
	if err != nil ***REMOVED***
		t.Fatalf("could not parse reference: %v", err)
	***REMOVED***
	if err = store.AddDigest(ref5.(reference.Canonical), testImageID2, false); err != nil ***REMOVED***
		t.Fatalf("error adding to store: %v", err)
	***REMOVED***

	// Attempt to overwrite with force == false
	if err = store.AddTag(ref4, testImageID3, false); err == nil || !strings.HasPrefix(err.Error(), "Conflict:") ***REMOVED***
		t.Fatalf("did not get expected error on overwrite attempt - got %v", err)
	***REMOVED***
	// Repeat to overwrite with force == true
	if err = store.AddTag(ref4, testImageID3, true); err != nil ***REMOVED***
		t.Fatalf("failed to force tag overwrite: %v", err)
	***REMOVED***

	// Check references so far
	id, err := store.Get(nameOnly)
	if err != nil ***REMOVED***
		t.Fatalf("Get returned error: %v", err)
	***REMOVED***
	if id != testImageID1 ***REMOVED***
		t.Fatalf("id mismatch: got %s instead of %s", id.String(), testImageID1.String())
	***REMOVED***

	id, err = store.Get(ref1)
	if err != nil ***REMOVED***
		t.Fatalf("Get returned error: %v", err)
	***REMOVED***
	if id != testImageID1 ***REMOVED***
		t.Fatalf("id mismatch: got %s instead of %s", id.String(), testImageID1.String())
	***REMOVED***

	id, err = store.Get(ref2)
	if err != nil ***REMOVED***
		t.Fatalf("Get returned error: %v", err)
	***REMOVED***
	if id != testImageID2 ***REMOVED***
		t.Fatalf("id mismatch: got %s instead of %s", id.String(), testImageID2.String())
	***REMOVED***

	id, err = store.Get(ref3)
	if err != nil ***REMOVED***
		t.Fatalf("Get returned error: %v", err)
	***REMOVED***
	if id != testImageID1 ***REMOVED***
		t.Fatalf("id mismatch: got %s instead of %s", id.String(), testImageID1.String())
	***REMOVED***

	id, err = store.Get(ref4)
	if err != nil ***REMOVED***
		t.Fatalf("Get returned error: %v", err)
	***REMOVED***
	if id != testImageID3 ***REMOVED***
		t.Fatalf("id mismatch: got %s instead of %s", id.String(), testImageID3.String())
	***REMOVED***

	id, err = store.Get(ref5)
	if err != nil ***REMOVED***
		t.Fatalf("Get returned error: %v", err)
	***REMOVED***
	if id != testImageID2 ***REMOVED***
		t.Fatalf("id mismatch: got %s instead of %s", id.String(), testImageID3.String())
	***REMOVED***

	// Get should return ErrDoesNotExist for a nonexistent repo
	nonExistRepo, err := reference.ParseNormalizedNamed("username/nonexistrepo:latest")
	if err != nil ***REMOVED***
		t.Fatalf("could not parse reference: %v", err)
	***REMOVED***
	if _, err = store.Get(nonExistRepo); err != ErrDoesNotExist ***REMOVED***
		t.Fatal("Expected ErrDoesNotExist from Get")
	***REMOVED***

	// Get should return ErrDoesNotExist for a nonexistent tag
	nonExistTag, err := reference.ParseNormalizedNamed("username/repo1:nonexist")
	if err != nil ***REMOVED***
		t.Fatalf("could not parse reference: %v", err)
	***REMOVED***
	if _, err = store.Get(nonExistTag); err != ErrDoesNotExist ***REMOVED***
		t.Fatal("Expected ErrDoesNotExist from Get")
	***REMOVED***

	// Check References
	refs := store.References(testImageID1)
	if len(refs) != 3 ***REMOVED***
		t.Fatal("unexpected number of references")
	***REMOVED***
	// Looking for the references in this order verifies that they are
	// returned lexically sorted.
	if refs[0].String() != ref3.String() ***REMOVED***
		t.Fatalf("unexpected reference: %v", refs[0].String())
	***REMOVED***
	if refs[1].String() != ref1.String() ***REMOVED***
		t.Fatalf("unexpected reference: %v", refs[1].String())
	***REMOVED***
	if refs[2].String() != nameOnly.String()+":latest" ***REMOVED***
		t.Fatalf("unexpected reference: %v", refs[2].String())
	***REMOVED***

	// Check ReferencesByName
	repoName, err := reference.ParseNormalizedNamed("username/repo1")
	if err != nil ***REMOVED***
		t.Fatalf("could not parse reference: %v", err)
	***REMOVED***
	associations := store.ReferencesByName(repoName)
	if len(associations) != 3 ***REMOVED***
		t.Fatal("unexpected number of associations")
	***REMOVED***
	// Looking for the associations in this order verifies that they are
	// returned lexically sorted.
	if associations[0].Ref.String() != ref3.String() ***REMOVED***
		t.Fatalf("unexpected reference: %v", associations[0].Ref.String())
	***REMOVED***
	if associations[0].ID != testImageID1 ***REMOVED***
		t.Fatalf("unexpected reference: %v", associations[0].Ref.String())
	***REMOVED***
	if associations[1].Ref.String() != ref1.String() ***REMOVED***
		t.Fatalf("unexpected reference: %v", associations[1].Ref.String())
	***REMOVED***
	if associations[1].ID != testImageID1 ***REMOVED***
		t.Fatalf("unexpected reference: %v", associations[1].Ref.String())
	***REMOVED***
	if associations[2].Ref.String() != ref2.String() ***REMOVED***
		t.Fatalf("unexpected reference: %v", associations[2].Ref.String())
	***REMOVED***
	if associations[2].ID != testImageID2 ***REMOVED***
		t.Fatalf("unexpected reference: %v", associations[2].Ref.String())
	***REMOVED***

	// Delete should return ErrDoesNotExist for a nonexistent repo
	if _, err = store.Delete(nonExistRepo); err != ErrDoesNotExist ***REMOVED***
		t.Fatal("Expected ErrDoesNotExist from Delete")
	***REMOVED***

	// Delete should return ErrDoesNotExist for a nonexistent tag
	if _, err = store.Delete(nonExistTag); err != ErrDoesNotExist ***REMOVED***
		t.Fatal("Expected ErrDoesNotExist from Delete")
	***REMOVED***

	// Delete a few references
	if deleted, err := store.Delete(ref1); err != nil || !deleted ***REMOVED***
		t.Fatal("Delete failed")
	***REMOVED***
	if _, err := store.Get(ref1); err != ErrDoesNotExist ***REMOVED***
		t.Fatal("Expected ErrDoesNotExist from Get")
	***REMOVED***
	if deleted, err := store.Delete(ref5); err != nil || !deleted ***REMOVED***
		t.Fatal("Delete failed")
	***REMOVED***
	if _, err := store.Get(ref5); err != ErrDoesNotExist ***REMOVED***
		t.Fatal("Expected ErrDoesNotExist from Get")
	***REMOVED***
	if deleted, err := store.Delete(nameOnly); err != nil || !deleted ***REMOVED***
		t.Fatal("Delete failed")
	***REMOVED***
	if _, err := store.Get(nameOnly); err != ErrDoesNotExist ***REMOVED***
		t.Fatal("Expected ErrDoesNotExist from Get")
	***REMOVED***
***REMOVED***

func TestInvalidTags(t *testing.T) ***REMOVED***
	tmpDir, err := ioutil.TempDir("", "tag-store-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	store, err := NewReferenceStore(filepath.Join(tmpDir, "repositories.json"))
	require.NoError(t, err)
	id := digest.Digest("sha256:470022b8af682154f57a2163d030eb369549549cba00edc69e1b99b46bb924d6")

	// sha256 as repo name
	ref, err := reference.ParseNormalizedNamed("sha256:abc")
	require.NoError(t, err)
	err = store.AddTag(ref, id, true)
	assert.Error(t, err)

	// setting digest as a tag
	ref, err = reference.ParseNormalizedNamed("registry@sha256:367eb40fd0330a7e464777121e39d2f5b3e8e23a1e159342e53ab05c9e4d94e6")
	require.NoError(t, err)

	err = store.AddTag(ref, id, true)
	assert.Error(t, err)
***REMOVED***
