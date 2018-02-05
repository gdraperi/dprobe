package afero

import "testing"

func TestCopyOnWrite(t *testing.T) ***REMOVED***
	var fs Fs
	var err error
	base := NewOsFs()
	roBase := NewReadOnlyFs(base)
	ufs := NewCopyOnWriteFs(roBase, NewMemMapFs())
	fs = ufs
	err = fs.MkdirAll("nonexistent/directory/", 0744)
	if err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***
	_, err = fs.Create("nonexistent/directory/newfile")
	if err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***

***REMOVED***

func TestCopyOnWriteFileInMemMapBase(t *testing.T) ***REMOVED***
	base := &MemMapFs***REMOVED******REMOVED***
	layer := &MemMapFs***REMOVED******REMOVED***

	if err := WriteFile(base, "base.txt", []byte("base"), 0755); err != nil ***REMOVED***
		t.Fatalf("Failed to write file: %s", err)
	***REMOVED***

	ufs := NewCopyOnWriteFs(base, layer)

	_, err := ufs.Stat("base.txt")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
