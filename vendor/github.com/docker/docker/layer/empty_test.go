package layer

import (
	"io"
	"testing"

	"github.com/opencontainers/go-digest"
)

func TestEmptyLayer(t *testing.T) ***REMOVED***
	if EmptyLayer.ChainID() != ChainID(DigestSHA256EmptyTar) ***REMOVED***
		t.Fatal("wrong ChainID for empty layer")
	***REMOVED***

	if EmptyLayer.DiffID() != DigestSHA256EmptyTar ***REMOVED***
		t.Fatal("wrong DiffID for empty layer")
	***REMOVED***

	if EmptyLayer.Parent() != nil ***REMOVED***
		t.Fatal("expected no parent for empty layer")
	***REMOVED***

	if size, err := EmptyLayer.Size(); err != nil || size != 0 ***REMOVED***
		t.Fatal("expected zero size for empty layer")
	***REMOVED***

	if diffSize, err := EmptyLayer.DiffSize(); err != nil || diffSize != 0 ***REMOVED***
		t.Fatal("expected zero diffsize for empty layer")
	***REMOVED***

	meta, err := EmptyLayer.Metadata()

	if len(meta) != 0 || err != nil ***REMOVED***
		t.Fatal("expected zero length metadata for empty layer")
	***REMOVED***

	tarStream, err := EmptyLayer.TarStream()
	if err != nil ***REMOVED***
		t.Fatalf("error streaming tar for empty layer: %v", err)
	***REMOVED***

	digester := digest.Canonical.Digester()
	_, err = io.Copy(digester.Hash(), tarStream)

	if err != nil ***REMOVED***
		t.Fatalf("error hashing empty tar layer: %v", err)
	***REMOVED***

	if digester.Digest() != digest.Digest(DigestSHA256EmptyTar) ***REMOVED***
		t.Fatal("empty layer tar stream hashes to wrong value")
	***REMOVED***
***REMOVED***
