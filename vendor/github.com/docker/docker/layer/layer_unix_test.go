// +build !windows

package layer

import (
	"testing"
)

func graphDiffSize(ls Store, l Layer) (int64, error) ***REMOVED***
	cl := getCachedLayer(l)
	var parent string
	if cl.parent != nil ***REMOVED***
		parent = cl.parent.cacheID
	***REMOVED***
	return ls.(*layerStore).driver.DiffSize(cl.cacheID, parent)
***REMOVED***

// Unix as Windows graph driver does not support Changes which is indirectly
// invoked by calling DiffSize on the driver
func TestLayerSize(t *testing.T) ***REMOVED***
	ls, _, cleanup := newTestStore(t)
	defer cleanup()

	content1 := []byte("Base contents")
	content2 := []byte("Added contents")

	layer1, err := createLayer(ls, "", initWithFiles(newTestFile("file1", content1, 0644)))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	layer2, err := createLayer(ls, layer1.ChainID(), initWithFiles(newTestFile("file2", content2, 0644)))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	layer1DiffSize, err := graphDiffSize(ls, layer1)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if int(layer1DiffSize) != len(content1) ***REMOVED***
		t.Fatalf("Unexpected diff size %d, expected %d", layer1DiffSize, len(content1))
	***REMOVED***

	layer1Size, err := layer1.Size()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if expected := len(content1); int(layer1Size) != expected ***REMOVED***
		t.Fatalf("Unexpected size %d, expected %d", layer1Size, expected)
	***REMOVED***

	layer2DiffSize, err := graphDiffSize(ls, layer2)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if int(layer2DiffSize) != len(content2) ***REMOVED***
		t.Fatalf("Unexpected diff size %d, expected %d", layer2DiffSize, len(content2))
	***REMOVED***

	layer2Size, err := layer2.Size()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if expected := len(content1) + len(content2); int(layer2Size) != expected ***REMOVED***
		t.Fatalf("Unexpected size %d, expected %d", layer2Size, expected)
	***REMOVED***

***REMOVED***
