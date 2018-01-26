package archive

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"testing"
)

func TestHardLinkOrder(t *testing.T) ***REMOVED***
	names := []string***REMOVED***"file1.txt", "file2.txt", "file3.txt"***REMOVED***
	msg := []byte("Hey y'all")

	// Create dir
	src, err := ioutil.TempDir("", "docker-hardlink-test-src-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(src)
	for _, name := range names ***REMOVED***
		func() ***REMOVED***
			fh, err := os.Create(path.Join(src, name))
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			defer fh.Close()
			if _, err = fh.Write(msg); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
		***REMOVED***()
	***REMOVED***
	// Create dest, with changes that includes hardlinks
	dest, err := ioutil.TempDir("", "docker-hardlink-test-dest-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	os.RemoveAll(dest) // we just want the name, at first
	if err := copyDir(src, dest); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(dest)
	for _, name := range names ***REMOVED***
		for i := 0; i < 5; i++ ***REMOVED***
			if err := os.Link(path.Join(dest, name), path.Join(dest, fmt.Sprintf("%s.link%d", name, i))); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// get changes
	changes, err := ChangesDirs(dest, src)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// sort
	sort.Sort(changesByPath(changes))

	// ExportChanges
	ar, err := ExportChanges(dest, changes, nil, nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	hdrs, err := walkHeaders(ar)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// reverse sort
	sort.Sort(sort.Reverse(changesByPath(changes)))
	// ExportChanges
	arRev, err := ExportChanges(dest, changes, nil, nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	hdrsRev, err := walkHeaders(arRev)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// line up the two sets
	sort.Sort(tarHeaders(hdrs))
	sort.Sort(tarHeaders(hdrsRev))

	// compare Size and LinkName
	for i := range hdrs ***REMOVED***
		if hdrs[i].Name != hdrsRev[i].Name ***REMOVED***
			t.Errorf("headers - expected name %q; but got %q", hdrs[i].Name, hdrsRev[i].Name)
		***REMOVED***
		if hdrs[i].Size != hdrsRev[i].Size ***REMOVED***
			t.Errorf("headers - %q expected size %d; but got %d", hdrs[i].Name, hdrs[i].Size, hdrsRev[i].Size)
		***REMOVED***
		if hdrs[i].Typeflag != hdrsRev[i].Typeflag ***REMOVED***
			t.Errorf("headers - %q expected type %d; but got %d", hdrs[i].Name, hdrs[i].Typeflag, hdrsRev[i].Typeflag)
		***REMOVED***
		if hdrs[i].Linkname != hdrsRev[i].Linkname ***REMOVED***
			t.Errorf("headers - %q expected linkname %q; but got %q", hdrs[i].Name, hdrs[i].Linkname, hdrsRev[i].Linkname)
		***REMOVED***
	***REMOVED***

***REMOVED***

type tarHeaders []tar.Header

func (th tarHeaders) Len() int           ***REMOVED*** return len(th) ***REMOVED***
func (th tarHeaders) Swap(i, j int)      ***REMOVED*** th[j], th[i] = th[i], th[j] ***REMOVED***
func (th tarHeaders) Less(i, j int) bool ***REMOVED*** return th[i].Name < th[j].Name ***REMOVED***

func walkHeaders(r io.Reader) ([]tar.Header, error) ***REMOVED***
	t := tar.NewReader(r)
	headers := []tar.Header***REMOVED******REMOVED***
	for ***REMOVED***
		hdr, err := t.Next()
		if err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				break
			***REMOVED***
			return headers, err
		***REMOVED***
		headers = append(headers, *hdr)
	***REMOVED***
	return headers, nil
***REMOVED***
