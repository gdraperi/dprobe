package tarsum

import (
	"runtime"
	"sort"
	"strings"
)

// FileInfoSumInterface provides an interface for accessing file checksum
// information within a tar file. This info is accessed through interface
// so the actual name and sum cannot be melded with.
type FileInfoSumInterface interface ***REMOVED***
	// File name
	Name() string
	// Checksum of this particular file and its headers
	Sum() string
	// Position of file in the tar
	Pos() int64
***REMOVED***

type fileInfoSum struct ***REMOVED***
	name string
	sum  string
	pos  int64
***REMOVED***

func (fis fileInfoSum) Name() string ***REMOVED***
	return fis.name
***REMOVED***
func (fis fileInfoSum) Sum() string ***REMOVED***
	return fis.sum
***REMOVED***
func (fis fileInfoSum) Pos() int64 ***REMOVED***
	return fis.pos
***REMOVED***

// FileInfoSums provides a list of FileInfoSumInterfaces.
type FileInfoSums []FileInfoSumInterface

// GetFile returns the first FileInfoSumInterface with a matching name.
func (fis FileInfoSums) GetFile(name string) FileInfoSumInterface ***REMOVED***
	// We do case insensitive matching on Windows as c:\APP and c:\app are
	// the same. See issue #33107.
	for i := range fis ***REMOVED***
		if (runtime.GOOS == "windows" && strings.EqualFold(fis[i].Name(), name)) ||
			(runtime.GOOS != "windows" && fis[i].Name() == name) ***REMOVED***
			return fis[i]
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// GetAllFile returns a FileInfoSums with all matching names.
func (fis FileInfoSums) GetAllFile(name string) FileInfoSums ***REMOVED***
	f := FileInfoSums***REMOVED******REMOVED***
	for i := range fis ***REMOVED***
		if fis[i].Name() == name ***REMOVED***
			f = append(f, fis[i])
		***REMOVED***
	***REMOVED***
	return f
***REMOVED***

// GetDuplicatePaths returns a FileInfoSums with all duplicated paths.
func (fis FileInfoSums) GetDuplicatePaths() (dups FileInfoSums) ***REMOVED***
	seen := make(map[string]int, len(fis)) // allocate earl. no need to grow this map.
	for i := range fis ***REMOVED***
		f := fis[i]
		if _, ok := seen[f.Name()]; ok ***REMOVED***
			dups = append(dups, f)
		***REMOVED*** else ***REMOVED***
			seen[f.Name()] = 0
		***REMOVED***
	***REMOVED***
	return dups
***REMOVED***

// Len returns the size of the FileInfoSums.
func (fis FileInfoSums) Len() int ***REMOVED*** return len(fis) ***REMOVED***

// Swap swaps two FileInfoSum values if a FileInfoSums list.
func (fis FileInfoSums) Swap(i, j int) ***REMOVED*** fis[i], fis[j] = fis[j], fis[i] ***REMOVED***

// SortByPos sorts FileInfoSums content by position.
func (fis FileInfoSums) SortByPos() ***REMOVED***
	sort.Sort(byPos***REMOVED***fis***REMOVED***)
***REMOVED***

// SortByNames sorts FileInfoSums content by name.
func (fis FileInfoSums) SortByNames() ***REMOVED***
	sort.Sort(byName***REMOVED***fis***REMOVED***)
***REMOVED***

// SortBySums sorts FileInfoSums content by sums.
func (fis FileInfoSums) SortBySums() ***REMOVED***
	dups := fis.GetDuplicatePaths()
	if len(dups) > 0 ***REMOVED***
		sort.Sort(bySum***REMOVED***fis, dups***REMOVED***)
	***REMOVED*** else ***REMOVED***
		sort.Sort(bySum***REMOVED***fis, nil***REMOVED***)
	***REMOVED***
***REMOVED***

// byName is a sort.Sort helper for sorting by file names.
// If names are the same, order them by their appearance in the tar archive
type byName struct***REMOVED*** FileInfoSums ***REMOVED***

func (bn byName) Less(i, j int) bool ***REMOVED***
	if bn.FileInfoSums[i].Name() == bn.FileInfoSums[j].Name() ***REMOVED***
		return bn.FileInfoSums[i].Pos() < bn.FileInfoSums[j].Pos()
	***REMOVED***
	return bn.FileInfoSums[i].Name() < bn.FileInfoSums[j].Name()
***REMOVED***

// bySum is a sort.Sort helper for sorting by the sums of all the fileinfos in the tar archive
type bySum struct ***REMOVED***
	FileInfoSums
	dups FileInfoSums
***REMOVED***

func (bs bySum) Less(i, j int) bool ***REMOVED***
	if bs.dups != nil && bs.FileInfoSums[i].Name() == bs.FileInfoSums[j].Name() ***REMOVED***
		return bs.FileInfoSums[i].Pos() < bs.FileInfoSums[j].Pos()
	***REMOVED***
	return bs.FileInfoSums[i].Sum() < bs.FileInfoSums[j].Sum()
***REMOVED***

// byPos is a sort.Sort helper for sorting by the sums of all the fileinfos by their original order
type byPos struct***REMOVED*** FileInfoSums ***REMOVED***

func (bp byPos) Less(i, j int) bool ***REMOVED***
	return bp.FileInfoSums[i].Pos() < bp.FileInfoSums[j].Pos()
***REMOVED***
