package metadata

import (
	"encoding/hex"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"testing"

	"github.com/docker/docker/layer"
	"github.com/opencontainers/go-digest"
)

func TestV2MetadataService(t *testing.T) ***REMOVED***
	tmpDir, err := ioutil.TempDir("", "blobsum-storage-service-test")
	if err != nil ***REMOVED***
		t.Fatalf("could not create temp dir: %v", err)
	***REMOVED***
	defer os.RemoveAll(tmpDir)

	metadataStore, err := NewFSMetadataStore(tmpDir)
	if err != nil ***REMOVED***
		t.Fatalf("could not create metadata store: %v", err)
	***REMOVED***
	V2MetadataService := NewV2MetadataService(metadataStore)

	tooManyBlobSums := make([]V2Metadata, 100)
	for i := range tooManyBlobSums ***REMOVED***
		randDigest := randomDigest()
		tooManyBlobSums[i] = V2Metadata***REMOVED***Digest: randDigest***REMOVED***
	***REMOVED***

	testVectors := []struct ***REMOVED***
		diffID   layer.DiffID
		metadata []V2Metadata
	***REMOVED******REMOVED***
		***REMOVED***
			diffID: layer.DiffID("sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4"),
			metadata: []V2Metadata***REMOVED***
				***REMOVED***Digest: digest.Digest("sha256:f0cd5ca10b07f35512fc2f1cbf9a6cefbdb5cba70ac6b0c9e5988f4497f71937")***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			diffID: layer.DiffID("sha256:86e0e091d0da6bde2456dbb48306f3956bbeb2eae1b5b9a43045843f69fe4aaa"),
			metadata: []V2Metadata***REMOVED***
				***REMOVED***Digest: digest.Digest("sha256:f0cd5ca10b07f35512fc2f1cbf9a6cefbdb5cba70ac6b0c9e5988f4497f71937")***REMOVED***,
				***REMOVED***Digest: digest.Digest("sha256:9e3447ca24cb96d86ebd5960cb34d1299b07e0a0e03801d90b9969a2c187dd6e")***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			diffID:   layer.DiffID("sha256:03f4658f8b782e12230c1783426bd3bacce651ce582a4ffb6fbbfa2079428ecb"),
			metadata: tooManyBlobSums,
		***REMOVED***,
	***REMOVED***

	// Set some associations
	for _, vec := range testVectors ***REMOVED***
		for _, blobsum := range vec.metadata ***REMOVED***
			err := V2MetadataService.Add(vec.diffID, blobsum)
			if err != nil ***REMOVED***
				t.Fatalf("error calling Set: %v", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Check the correct values are read back
	for _, vec := range testVectors ***REMOVED***
		metadata, err := V2MetadataService.GetMetadata(vec.diffID)
		if err != nil ***REMOVED***
			t.Fatalf("error calling Get: %v", err)
		***REMOVED***
		expectedMetadataEntries := len(vec.metadata)
		if expectedMetadataEntries > 50 ***REMOVED***
			expectedMetadataEntries = 50
		***REMOVED***
		if !reflect.DeepEqual(metadata, vec.metadata[len(vec.metadata)-expectedMetadataEntries:len(vec.metadata)]) ***REMOVED***
			t.Fatal("Get returned incorrect layer ID")
		***REMOVED***
	***REMOVED***

	// Test GetMetadata on a nonexistent entry
	_, err = V2MetadataService.GetMetadata(layer.DiffID("sha256:82379823067823853223359023576437723560923756b03560378f4497753917"))
	if err == nil ***REMOVED***
		t.Fatal("expected error looking up nonexistent entry")
	***REMOVED***

	// Test GetDiffID on a nonexistent entry
	_, err = V2MetadataService.GetDiffID(digest.Digest("sha256:82379823067823853223359023576437723560923756b03560378f4497753917"))
	if err == nil ***REMOVED***
		t.Fatal("expected error looking up nonexistent entry")
	***REMOVED***

	// Overwrite one of the entries and read it back
	err = V2MetadataService.Add(testVectors[1].diffID, testVectors[0].metadata[0])
	if err != nil ***REMOVED***
		t.Fatalf("error calling Add: %v", err)
	***REMOVED***
	diffID, err := V2MetadataService.GetDiffID(testVectors[0].metadata[0].Digest)
	if err != nil ***REMOVED***
		t.Fatalf("error calling GetDiffID: %v", err)
	***REMOVED***
	if diffID != testVectors[1].diffID ***REMOVED***
		t.Fatal("GetDiffID returned incorrect diffID")
	***REMOVED***
***REMOVED***

func randomDigest() digest.Digest ***REMOVED***
	b := [32]byte***REMOVED******REMOVED***
	for i := 0; i < len(b); i++ ***REMOVED***
		b[i] = byte(rand.Intn(256))
	***REMOVED***
	d := hex.EncodeToString(b[:])
	return digest.Digest("sha256:" + d)
***REMOVED***
