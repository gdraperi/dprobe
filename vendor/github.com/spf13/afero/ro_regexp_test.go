package afero

import (
	"regexp"
	"testing"
)

func TestFilterReadOnly(t *testing.T) ***REMOVED***
	fs := &ReadOnlyFs***REMOVED***source: &MemMapFs***REMOVED******REMOVED******REMOVED***
	_, err := fs.Create("/file.txt")
	if err == nil ***REMOVED***
		t.Errorf("Did not fail to create file")
	***REMOVED***
	// t.Logf("ERR=%s", err)
***REMOVED***

func TestFilterReadonlyRemoveAndRead(t *testing.T) ***REMOVED***
	mfs := &MemMapFs***REMOVED******REMOVED***
	fh, err := mfs.Create("/file.txt")
	fh.Write([]byte("content here"))
	fh.Close()

	fs := NewReadOnlyFs(mfs)
	err = fs.Remove("/file.txt")
	if err == nil ***REMOVED***
		t.Errorf("Did not fail to remove file")
	***REMOVED***

	fh, err = fs.Open("/file.txt")
	if err != nil ***REMOVED***
		t.Errorf("Failed to open file: %s", err)
	***REMOVED***

	buf := make([]byte, len("content here"))
	_, err = fh.Read(buf)
	fh.Close()
	if string(buf) != "content here" ***REMOVED***
		t.Errorf("Failed to read file: %s", err)
	***REMOVED***

	err = mfs.Remove("/file.txt")
	if err != nil ***REMOVED***
		t.Errorf("Failed to remove file")
	***REMOVED***

	fh, err = fs.Open("/file.txt")
	if err == nil ***REMOVED***
		fh.Close()
		t.Errorf("File still present")
	***REMOVED***
***REMOVED***

func TestFilterRegexp(t *testing.T) ***REMOVED***
	fs := NewRegexpFs(&MemMapFs***REMOVED******REMOVED***, regexp.MustCompile(`\.txt$`))
	_, err := fs.Create("/file.html")
	if err == nil ***REMOVED***

		t.Errorf("Did not fail to create file")
	***REMOVED***
	// t.Logf("ERR=%s", err)
***REMOVED***

func TestFilterRORegexpChain(t *testing.T) ***REMOVED***
	rofs := &ReadOnlyFs***REMOVED***source: &MemMapFs***REMOVED******REMOVED******REMOVED***
	fs := &RegexpFs***REMOVED***re: regexp.MustCompile(`\.txt$`), source: rofs***REMOVED***
	_, err := fs.Create("/file.txt")
	if err == nil ***REMOVED***
		t.Errorf("Did not fail to create file")
	***REMOVED***
	// t.Logf("ERR=%s", err)
***REMOVED***

func TestFilterRegexReadDir(t *testing.T) ***REMOVED***
	mfs := &MemMapFs***REMOVED******REMOVED***
	fs1 := &RegexpFs***REMOVED***re: regexp.MustCompile(`\.txt$`), source: mfs***REMOVED***
	fs := &RegexpFs***REMOVED***re: regexp.MustCompile(`^a`), source: fs1***REMOVED***

	mfs.MkdirAll("/dir/sub", 0777)
	for _, name := range []string***REMOVED***"afile.txt", "afile.html", "bfile.txt"***REMOVED*** ***REMOVED***
		for _, dir := range []string***REMOVED***"/dir/", "/dir/sub/"***REMOVED*** ***REMOVED***
			fh, _ := mfs.Create(dir + name)
			fh.Close()
		***REMOVED***
	***REMOVED***

	files, _ := ReadDir(fs, "/dir")
	if len(files) != 2 ***REMOVED*** // afile.txt, sub
		t.Errorf("Got wrong number of files: %#v", files)
	***REMOVED***

	f, _ := fs.Open("/dir/sub")
	names, _ := f.Readdirnames(-1)
	if len(names) != 1 ***REMOVED***
		t.Errorf("Got wrong number of names: %v", names)
	***REMOVED***
***REMOVED***
