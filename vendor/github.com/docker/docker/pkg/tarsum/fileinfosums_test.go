package tarsum

import "testing"

func newFileInfoSums() FileInfoSums ***REMOVED***
	return FileInfoSums***REMOVED***
		fileInfoSum***REMOVED***name: "file3", sum: "2abcdef1234567890", pos: 2***REMOVED***,
		fileInfoSum***REMOVED***name: "dup1", sum: "deadbeef1", pos: 5***REMOVED***,
		fileInfoSum***REMOVED***name: "file1", sum: "0abcdef1234567890", pos: 0***REMOVED***,
		fileInfoSum***REMOVED***name: "file4", sum: "3abcdef1234567890", pos: 3***REMOVED***,
		fileInfoSum***REMOVED***name: "dup1", sum: "deadbeef0", pos: 4***REMOVED***,
		fileInfoSum***REMOVED***name: "file2", sum: "1abcdef1234567890", pos: 1***REMOVED***,
	***REMOVED***
***REMOVED***

func TestSortFileInfoSums(t *testing.T) ***REMOVED***
	dups := newFileInfoSums().GetAllFile("dup1")
	if len(dups) != 2 ***REMOVED***
		t.Errorf("expected length 2, got %d", len(dups))
	***REMOVED***
	dups.SortByNames()
	if dups[0].Pos() != 4 ***REMOVED***
		t.Errorf("sorted dups should be ordered by position. Expected 4, got %d", dups[0].Pos())
	***REMOVED***

	fis := newFileInfoSums()
	expected := "0abcdef1234567890"
	fis.SortBySums()
	got := fis[0].Sum()
	if got != expected ***REMOVED***
		t.Errorf("Expected %q, got %q", expected, got)
	***REMOVED***

	fis = newFileInfoSums()
	expected = "dup1"
	fis.SortByNames()
	gotFis := fis[0]
	if gotFis.Name() != expected ***REMOVED***
		t.Errorf("Expected %q, got %q", expected, gotFis.Name())
	***REMOVED***
	// since a duplicate is first, ensure it is ordered first by position too
	if gotFis.Pos() != 4 ***REMOVED***
		t.Errorf("Expected %d, got %d", 4, gotFis.Pos())
	***REMOVED***

	fis = newFileInfoSums()
	fis.SortByPos()
	if fis[0].Pos() != 0 ***REMOVED***
		t.Error("sorted fileInfoSums by Pos should order them by position.")
	***REMOVED***

	fis = newFileInfoSums()
	expected = "deadbeef1"
	gotFileInfoSum := fis.GetFile("dup1")
	if gotFileInfoSum.Sum() != expected ***REMOVED***
		t.Errorf("Expected %q, got %q", expected, gotFileInfoSum)
	***REMOVED***
	if fis.GetFile("noPresent") != nil ***REMOVED***
		t.Error("Should have return nil if name not found.")
	***REMOVED***

***REMOVED***
