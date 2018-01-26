package layer

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/opencontainers/go-digest"
)

func randomLayerID(seed int64) ChainID ***REMOVED***
	r := rand.New(rand.NewSource(seed))

	return ChainID(digest.FromBytes([]byte(fmt.Sprintf("%d", r.Int63()))))
***REMOVED***

func newFileMetadataStore(t *testing.T) (*fileMetadataStore, string, func()) ***REMOVED***
	td, err := ioutil.TempDir("", "layers-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	fms, err := NewFSMetadataStore(td)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	return fms.(*fileMetadataStore), td, func() ***REMOVED***
		if err := os.RemoveAll(td); err != nil ***REMOVED***
			t.Logf("Failed to cleanup %q: %s", td, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func assertNotDirectoryError(t *testing.T, err error) ***REMOVED***
	perr, ok := err.(*os.PathError)
	if !ok ***REMOVED***
		t.Fatalf("Unexpected error %#v, expected path error", err)
	***REMOVED***

	if perr.Err != syscall.ENOTDIR ***REMOVED***
		t.Fatalf("Unexpected error %s, expected %s", perr.Err, syscall.ENOTDIR)
	***REMOVED***
***REMOVED***

func TestCommitFailure(t *testing.T) ***REMOVED***
	fms, td, cleanup := newFileMetadataStore(t)
	defer cleanup()

	if err := ioutil.WriteFile(filepath.Join(td, "sha256"), []byte("was here first!"), 0644); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	tx, err := fms.StartTransaction()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := tx.SetSize(0); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = tx.Commit(randomLayerID(5))
	if err == nil ***REMOVED***
		t.Fatalf("Expected error committing with invalid layer parent directory")
	***REMOVED***
	assertNotDirectoryError(t, err)
***REMOVED***

func TestStartTransactionFailure(t *testing.T) ***REMOVED***
	fms, td, cleanup := newFileMetadataStore(t)
	defer cleanup()

	if err := ioutil.WriteFile(filepath.Join(td, "tmp"), []byte("was here first!"), 0644); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	_, err := fms.StartTransaction()
	if err == nil ***REMOVED***
		t.Fatalf("Expected error starting transaction with invalid layer parent directory")
	***REMOVED***
	assertNotDirectoryError(t, err)

	if err := os.Remove(filepath.Join(td, "tmp")); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	tx, err := fms.StartTransaction()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if expected := filepath.Join(td, "tmp"); strings.HasPrefix(expected, tx.String()) ***REMOVED***
		t.Fatalf("Unexpected transaction string %q, expected prefix %q", tx.String(), expected)
	***REMOVED***

	if err := tx.Cancel(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
