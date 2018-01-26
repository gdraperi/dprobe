package remotecontext

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/docker/builder"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/reexec"
	"github.com/pkg/errors"
)

const (
	filename = "test"
	contents = "contents test"
)

func init() ***REMOVED***
	reexec.Init()
***REMOVED***

func TestCloseRootDirectory(t *testing.T) ***REMOVED***
	contextDir, err := ioutil.TempDir("", "builder-tarsum-test")
	defer os.RemoveAll(contextDir)
	if err != nil ***REMOVED***
		t.Fatalf("Error with creating temporary directory: %s", err)
	***REMOVED***

	src := makeTestArchiveContext(t, contextDir)
	err = src.Close()

	if err != nil ***REMOVED***
		t.Fatalf("Error while executing Close: %s", err)
	***REMOVED***

	_, err = os.Stat(src.Root().Path())

	if !os.IsNotExist(err) ***REMOVED***
		t.Fatal("Directory should not exist at this point")
	***REMOVED***
***REMOVED***

func TestHashFile(t *testing.T) ***REMOVED***
	contextDir, cleanup := createTestTempDir(t, "", "builder-tarsum-test")
	defer cleanup()

	createTestTempFile(t, contextDir, filename, contents, 0755)

	tarSum := makeTestArchiveContext(t, contextDir)

	sum, err := tarSum.Hash(filename)

	if err != nil ***REMOVED***
		t.Fatalf("Error when executing Stat: %s", err)
	***REMOVED***

	if len(sum) == 0 ***REMOVED***
		t.Fatalf("Hash returned empty sum")
	***REMOVED***

	expected := "1149ab94af7be6cc1da1335e398f24ee1cf4926b720044d229969dfc248ae7ec"

	if actual := sum; expected != actual ***REMOVED***
		t.Fatalf("invalid checksum. expected %s, got %s", expected, actual)
	***REMOVED***
***REMOVED***

func TestHashSubdir(t *testing.T) ***REMOVED***
	contextDir, cleanup := createTestTempDir(t, "", "builder-tarsum-test")
	defer cleanup()

	contextSubdir := filepath.Join(contextDir, "builder-tarsum-test-subdir")
	err := os.Mkdir(contextSubdir, 0755)
	if err != nil ***REMOVED***
		t.Fatalf("Failed to make directory: %s", contextSubdir)
	***REMOVED***

	testFilename := createTestTempFile(t, contextSubdir, filename, contents, 0755)

	tarSum := makeTestArchiveContext(t, contextDir)

	relativePath, err := filepath.Rel(contextDir, testFilename)

	if err != nil ***REMOVED***
		t.Fatalf("Error when getting relative path: %s", err)
	***REMOVED***

	sum, err := tarSum.Hash(relativePath)

	if err != nil ***REMOVED***
		t.Fatalf("Error when executing Stat: %s", err)
	***REMOVED***

	if len(sum) == 0 ***REMOVED***
		t.Fatalf("Hash returned empty sum")
	***REMOVED***

	expected := "d7f8d6353dee4816f9134f4156bf6a9d470fdadfb5d89213721f7e86744a4e69"

	if actual := sum; expected != actual ***REMOVED***
		t.Fatalf("invalid checksum. expected %s, got %s", expected, actual)
	***REMOVED***
***REMOVED***

func TestRemoveDirectory(t *testing.T) ***REMOVED***
	contextDir, cleanup := createTestTempDir(t, "", "builder-tarsum-test")
	defer cleanup()

	contextSubdir := createTestTempSubdir(t, contextDir, "builder-tarsum-test-subdir")

	relativePath, err := filepath.Rel(contextDir, contextSubdir)

	if err != nil ***REMOVED***
		t.Fatalf("Error when getting relative path: %s", err)
	***REMOVED***

	src := makeTestArchiveContext(t, contextDir)

	_, err = src.Root().Stat(src.Root().Join(src.Root().Path(), relativePath))
	if err != nil ***REMOVED***
		t.Fatalf("Statting %s shouldn't fail: %+v", relativePath, err)
	***REMOVED***

	tarSum := src.(modifiableContext)
	err = tarSum.Remove(relativePath)
	if err != nil ***REMOVED***
		t.Fatalf("Error when executing Remove: %s", err)
	***REMOVED***

	_, err = src.Root().Stat(src.Root().Join(src.Root().Path(), relativePath))
	if !os.IsNotExist(errors.Cause(err)) ***REMOVED***
		t.Fatalf("Directory should not exist at this point: %+v ", err)
	***REMOVED***
***REMOVED***

func makeTestArchiveContext(t *testing.T, dir string) builder.Source ***REMOVED***
	tarStream, err := archive.Tar(dir, archive.Uncompressed)
	if err != nil ***REMOVED***
		t.Fatalf("error: %s", err)
	***REMOVED***
	defer tarStream.Close()
	tarSum, err := FromArchive(tarStream)
	if err != nil ***REMOVED***
		t.Fatalf("Error when executing FromArchive: %s", err)
	***REMOVED***
	return tarSum
***REMOVED***
