package afero

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var tempDirs []string

func NewTempOsBaseFs(t *testing.T) Fs ***REMOVED***
	name, err := TempDir(NewOsFs(), "", "")
	if err != nil ***REMOVED***
		t.Error("error creating tempDir", err)
	***REMOVED***

	tempDirs = append(tempDirs, name)

	return NewBasePathFs(NewOsFs(), name)
***REMOVED***

func CleanupTempDirs(t *testing.T) ***REMOVED***
	osfs := NewOsFs()
	type ev struct ***REMOVED***
		path string
		e    error
	***REMOVED***

	errs := []ev***REMOVED******REMOVED***

	for _, x := range tempDirs ***REMOVED***
		err := osfs.RemoveAll(x)
		if err != nil ***REMOVED***
			errs = append(errs, ev***REMOVED***path: x, e: err***REMOVED***)
		***REMOVED***
	***REMOVED***

	for _, e := range errs ***REMOVED***
		fmt.Println("error removing tempDir", e.path, e.e)
	***REMOVED***

	if len(errs) > 0 ***REMOVED***
		t.Error("error cleaning up tempDirs")
	***REMOVED***
	tempDirs = []string***REMOVED******REMOVED***
***REMOVED***

func TestUnionCreateExisting(t *testing.T) ***REMOVED***
	base := &MemMapFs***REMOVED******REMOVED***
	roBase := &ReadOnlyFs***REMOVED***source: base***REMOVED***

	ufs := NewCopyOnWriteFs(roBase, &MemMapFs***REMOVED******REMOVED***)

	base.MkdirAll("/home/test", 0777)
	fh, _ := base.Create("/home/test/file.txt")
	fh.WriteString("This is a test")
	fh.Close()

	fh, err := ufs.OpenFile("/home/test/file.txt", os.O_RDWR, 0666)
	if err != nil ***REMOVED***
		t.Errorf("Failed to open file r/w: %s", err)
	***REMOVED***

	_, err = fh.Write([]byte("####"))
	if err != nil ***REMOVED***
		t.Errorf("Failed to write file: %s", err)
	***REMOVED***
	fh.Seek(0, 0)
	data, err := ioutil.ReadAll(fh)
	if err != nil ***REMOVED***
		t.Errorf("Failed to read file: %s", err)
	***REMOVED***
	if string(data) != "#### is a test" ***REMOVED***
		t.Errorf("Got wrong data")
	***REMOVED***
	fh.Close()

	fh, _ = base.Open("/home/test/file.txt")
	data, err = ioutil.ReadAll(fh)
	if string(data) != "This is a test" ***REMOVED***
		t.Errorf("Got wrong data in base file")
	***REMOVED***
	fh.Close()

	fh, err = ufs.Create("/home/test/file.txt")
	switch err ***REMOVED***
	case nil:
		if fi, _ := fh.Stat(); fi.Size() != 0 ***REMOVED***
			t.Errorf("Create did not truncate file")
		***REMOVED***
		fh.Close()
	default:
		t.Errorf("Create failed on existing file")
	***REMOVED***

***REMOVED***

func TestUnionMergeReaddir(t *testing.T) ***REMOVED***
	base := &MemMapFs***REMOVED******REMOVED***
	roBase := &ReadOnlyFs***REMOVED***source: base***REMOVED***

	ufs := &CopyOnWriteFs***REMOVED***base: roBase, layer: &MemMapFs***REMOVED******REMOVED******REMOVED***

	base.MkdirAll("/home/test", 0777)
	fh, _ := base.Create("/home/test/file.txt")
	fh.WriteString("This is a test")
	fh.Close()

	fh, _ = ufs.Create("/home/test/file2.txt")
	fh.WriteString("This is a test")
	fh.Close()

	fh, _ = ufs.Open("/home/test")
	files, err := fh.Readdirnames(-1)
	if err != nil ***REMOVED***
		t.Errorf("Readdirnames failed")
	***REMOVED***
	if len(files) != 2 ***REMOVED***
		t.Errorf("Got wrong number of files: %v", files)
	***REMOVED***
***REMOVED***

func TestExistingDirectoryCollisionReaddir(t *testing.T) ***REMOVED***
	base := &MemMapFs***REMOVED******REMOVED***
	roBase := &ReadOnlyFs***REMOVED***source: base***REMOVED***
	overlay := &MemMapFs***REMOVED******REMOVED***

	ufs := &CopyOnWriteFs***REMOVED***base: roBase, layer: overlay***REMOVED***

	base.MkdirAll("/home/test", 0777)
	fh, _ := base.Create("/home/test/file.txt")
	fh.WriteString("This is a test")
	fh.Close()

	overlay.MkdirAll("home/test", 0777)
	fh, _ = overlay.Create("/home/test/file2.txt")
	fh.WriteString("This is a test")
	fh.Close()

	fh, _ = ufs.Create("/home/test/file3.txt")
	fh.WriteString("This is a test")
	fh.Close()

	fh, _ = ufs.Open("/home/test")
	files, err := fh.Readdirnames(-1)
	if err != nil ***REMOVED***
		t.Errorf("Readdirnames failed")
	***REMOVED***
	if len(files) != 3 ***REMOVED***
		t.Errorf("Got wrong number of files in union: %v", files)
	***REMOVED***

	fh, _ = overlay.Open("/home/test")
	files, err = fh.Readdirnames(-1)
	if err != nil ***REMOVED***
		t.Errorf("Readdirnames failed")
	***REMOVED***
	if len(files) != 2 ***REMOVED***
		t.Errorf("Got wrong number of files in overlay: %v", files)
	***REMOVED***
***REMOVED***

func TestNestedDirBaseReaddir(t *testing.T) ***REMOVED***
	base := &MemMapFs***REMOVED******REMOVED***
	roBase := &ReadOnlyFs***REMOVED***source: base***REMOVED***
	overlay := &MemMapFs***REMOVED******REMOVED***

	ufs := &CopyOnWriteFs***REMOVED***base: roBase, layer: overlay***REMOVED***

	base.MkdirAll("/home/test/foo/bar", 0777)
	fh, _ := base.Create("/home/test/file.txt")
	fh.WriteString("This is a test")
	fh.Close()

	fh, _ = base.Create("/home/test/foo/file2.txt")
	fh.WriteString("This is a test")
	fh.Close()
	fh, _ = base.Create("/home/test/foo/bar/file3.txt")
	fh.WriteString("This is a test")
	fh.Close()

	overlay.MkdirAll("/", 0777)

	// Opening something only in the base
	fh, _ = ufs.Open("/home/test/foo")
	list, err := fh.Readdir(-1)
	if err != nil ***REMOVED***
		t.Errorf("Readdir failed %s", err)
	***REMOVED***
	if len(list) != 2 ***REMOVED***
		for _, x := range list ***REMOVED***
			fmt.Println(x.Name())
		***REMOVED***
		t.Errorf("Got wrong number of files in union: %v", len(list))
	***REMOVED***
***REMOVED***

func TestNestedDirOverlayReaddir(t *testing.T) ***REMOVED***
	base := &MemMapFs***REMOVED******REMOVED***
	roBase := &ReadOnlyFs***REMOVED***source: base***REMOVED***
	overlay := &MemMapFs***REMOVED******REMOVED***

	ufs := &CopyOnWriteFs***REMOVED***base: roBase, layer: overlay***REMOVED***

	base.MkdirAll("/", 0777)
	overlay.MkdirAll("/home/test/foo/bar", 0777)
	fh, _ := overlay.Create("/home/test/file.txt")
	fh.WriteString("This is a test")
	fh.Close()
	fh, _ = overlay.Create("/home/test/foo/file2.txt")
	fh.WriteString("This is a test")
	fh.Close()
	fh, _ = overlay.Create("/home/test/foo/bar/file3.txt")
	fh.WriteString("This is a test")
	fh.Close()

	// Opening nested dir only in the overlay
	fh, _ = ufs.Open("/home/test/foo")
	list, err := fh.Readdir(-1)
	if err != nil ***REMOVED***
		t.Errorf("Readdir failed %s", err)
	***REMOVED***
	if len(list) != 2 ***REMOVED***
		for _, x := range list ***REMOVED***
			fmt.Println(x.Name())
		***REMOVED***
		t.Errorf("Got wrong number of files in union: %v", len(list))
	***REMOVED***
***REMOVED***

func TestNestedDirOverlayOsFsReaddir(t *testing.T) ***REMOVED***
	defer CleanupTempDirs(t)
	base := NewTempOsBaseFs(t)
	roBase := &ReadOnlyFs***REMOVED***source: base***REMOVED***
	overlay := NewTempOsBaseFs(t)

	ufs := &CopyOnWriteFs***REMOVED***base: roBase, layer: overlay***REMOVED***

	base.MkdirAll("/", 0777)
	overlay.MkdirAll("/home/test/foo/bar", 0777)
	fh, _ := overlay.Create("/home/test/file.txt")
	fh.WriteString("This is a test")
	fh.Close()
	fh, _ = overlay.Create("/home/test/foo/file2.txt")
	fh.WriteString("This is a test")
	fh.Close()
	fh, _ = overlay.Create("/home/test/foo/bar/file3.txt")
	fh.WriteString("This is a test")
	fh.Close()

	// Opening nested dir only in the overlay
	fh, _ = ufs.Open("/home/test/foo")
	list, err := fh.Readdir(-1)
	fh.Close()
	if err != nil ***REMOVED***
		t.Errorf("Readdir failed %s", err)
	***REMOVED***
	if len(list) != 2 ***REMOVED***
		for _, x := range list ***REMOVED***
			fmt.Println(x.Name())
		***REMOVED***
		t.Errorf("Got wrong number of files in union: %v", len(list))
	***REMOVED***
***REMOVED***

func TestCopyOnWriteFsWithOsFs(t *testing.T) ***REMOVED***
	defer CleanupTempDirs(t)
	base := NewTempOsBaseFs(t)
	roBase := &ReadOnlyFs***REMOVED***source: base***REMOVED***
	overlay := NewTempOsBaseFs(t)

	ufs := &CopyOnWriteFs***REMOVED***base: roBase, layer: overlay***REMOVED***

	base.MkdirAll("/home/test", 0777)
	fh, _ := base.Create("/home/test/file.txt")
	fh.WriteString("This is a test")
	fh.Close()

	overlay.MkdirAll("home/test", 0777)
	fh, _ = overlay.Create("/home/test/file2.txt")
	fh.WriteString("This is a test")
	fh.Close()

	fh, _ = ufs.Create("/home/test/file3.txt")
	fh.WriteString("This is a test")
	fh.Close()

	fh, _ = ufs.Open("/home/test")
	files, err := fh.Readdirnames(-1)
	fh.Close()
	if err != nil ***REMOVED***
		t.Errorf("Readdirnames failed")
	***REMOVED***
	if len(files) != 3 ***REMOVED***
		t.Errorf("Got wrong number of files in union: %v", files)
	***REMOVED***

	fh, _ = overlay.Open("/home/test")
	files, err = fh.Readdirnames(-1)
	fh.Close()
	if err != nil ***REMOVED***
		t.Errorf("Readdirnames failed")
	***REMOVED***
	if len(files) != 2 ***REMOVED***
		t.Errorf("Got wrong number of files in overlay: %v", files)
	***REMOVED***
***REMOVED***

func TestUnionCacheWrite(t *testing.T) ***REMOVED***
	base := &MemMapFs***REMOVED******REMOVED***
	layer := &MemMapFs***REMOVED******REMOVED***

	ufs := NewCacheOnReadFs(base, layer, 0)

	base.Mkdir("/data", 0777)

	fh, err := ufs.Create("/data/file.txt")
	if err != nil ***REMOVED***
		t.Errorf("Failed to create file")
	***REMOVED***
	_, err = fh.Write([]byte("This is a test"))
	if err != nil ***REMOVED***
		t.Errorf("Failed to write file")
	***REMOVED***

	fh.Seek(0, os.SEEK_SET)
	buf := make([]byte, 4)
	_, err = fh.Read(buf)
	fh.Write([]byte(" IS A"))
	fh.Close()

	baseData, _ := ReadFile(base, "/data/file.txt")
	layerData, _ := ReadFile(layer, "/data/file.txt")
	if string(baseData) != string(layerData) ***REMOVED***
		t.Errorf("Different data: %s <=> %s", baseData, layerData)
	***REMOVED***
***REMOVED***

func TestUnionCacheExpire(t *testing.T) ***REMOVED***
	base := &MemMapFs***REMOVED******REMOVED***
	layer := &MemMapFs***REMOVED******REMOVED***
	ufs := &CacheOnReadFs***REMOVED***base: base, layer: layer, cacheTime: 1 * time.Second***REMOVED***

	base.Mkdir("/data", 0777)

	fh, err := ufs.Create("/data/file.txt")
	if err != nil ***REMOVED***
		t.Errorf("Failed to create file")
	***REMOVED***
	_, err = fh.Write([]byte("This is a test"))
	if err != nil ***REMOVED***
		t.Errorf("Failed to write file")
	***REMOVED***
	fh.Close()

	fh, _ = base.Create("/data/file.txt")
	// sleep some time, so we really get a different time.Now() on write...
	time.Sleep(2 * time.Second)
	fh.WriteString("Another test")
	fh.Close()

	data, _ := ReadFile(ufs, "/data/file.txt")
	if string(data) != "Another test" ***REMOVED***
		t.Errorf("cache time failed: <%s>", data)
	***REMOVED***
***REMOVED***

func TestCacheOnReadFsNotInLayer(t *testing.T) ***REMOVED***
	base := NewMemMapFs()
	layer := NewMemMapFs()
	fs := NewCacheOnReadFs(base, layer, 0)

	fh, err := base.Create("/file.txt")
	if err != nil ***REMOVED***
		t.Fatal("unable to create file: ", err)
	***REMOVED***

	txt := []byte("This is a test")
	fh.Write(txt)
	fh.Close()

	fh, err = fs.Open("/file.txt")
	if err != nil ***REMOVED***
		t.Fatal("could not open file: ", err)
	***REMOVED***

	b, err := ReadAll(fh)
	fh.Close()

	if err != nil ***REMOVED***
		t.Fatal("could not read file: ", err)
	***REMOVED*** else if !bytes.Equal(txt, b) ***REMOVED***
		t.Fatalf("wanted file text %q, got %q", txt, b)
	***REMOVED***

	fh, err = layer.Open("/file.txt")
	if err != nil ***REMOVED***
		t.Fatal("could not open file from layer: ", err)
	***REMOVED***
	fh.Close()
***REMOVED***
