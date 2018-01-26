package tarsum

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

// Try to remove tarsum (in the BuilderContext) that do not exists, won't change a thing
func TestTarSumRemoveNonExistent(t *testing.T) ***REMOVED***
	filename := "testdata/46af0962ab5afeb5ce6740d4d91652e69206fc991fd5328c1a94d364ad00e457/layer.tar"
	reader, err := os.Open(filename)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer reader.Close()

	ts, err := NewTarSum(reader, false, Version0)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Read and discard bytes so that it populates sums
	_, err = io.Copy(ioutil.Discard, ts)
	if err != nil ***REMOVED***
		t.Errorf("failed to read from %s: %s", filename, err)
	***REMOVED***

	expected := len(ts.GetSums())

	ts.(BuilderContext).Remove("")
	ts.(BuilderContext).Remove("Anything")

	if len(ts.GetSums()) != expected ***REMOVED***
		t.Fatalf("Expected %v sums, go %v.", expected, ts.GetSums())
	***REMOVED***
***REMOVED***

// Remove a tarsum (in the BuilderContext)
func TestTarSumRemove(t *testing.T) ***REMOVED***
	filename := "testdata/46af0962ab5afeb5ce6740d4d91652e69206fc991fd5328c1a94d364ad00e457/layer.tar"
	reader, err := os.Open(filename)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer reader.Close()

	ts, err := NewTarSum(reader, false, Version0)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Read and discard bytes so that it populates sums
	_, err = io.Copy(ioutil.Discard, ts)
	if err != nil ***REMOVED***
		t.Errorf("failed to read from %s: %s", filename, err)
	***REMOVED***

	expected := len(ts.GetSums()) - 1

	ts.(BuilderContext).Remove("etc/sudoers")

	if len(ts.GetSums()) != expected ***REMOVED***
		t.Fatalf("Expected %v sums, go %v.", expected, len(ts.GetSums()))
	***REMOVED***
***REMOVED***
