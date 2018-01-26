package system

import (
	"os"
	"time"
)

// StatT type contains status of a file. It contains metadata
// like permission, size, etc about a file.
type StatT struct ***REMOVED***
	mode os.FileMode
	size int64
	mtim time.Time
***REMOVED***

// Size returns file's size.
func (s StatT) Size() int64 ***REMOVED***
	return s.size
***REMOVED***

// Mode returns file's permission mode.
func (s StatT) Mode() os.FileMode ***REMOVED***
	return os.FileMode(s.mode)
***REMOVED***

// Mtim returns file's last modification time.
func (s StatT) Mtim() time.Time ***REMOVED***
	return time.Time(s.mtim)
***REMOVED***

// Stat takes a path to a file and returns
// a system.StatT type pertaining to that file.
//
// Throws an error if the file does not exist
func Stat(path string) (*StatT, error) ***REMOVED***
	fi, err := os.Stat(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return fromStatT(&fi)
***REMOVED***

// fromStatT converts a os.FileInfo type to a system.StatT type
func fromStatT(fi *os.FileInfo) (*StatT, error) ***REMOVED***
	return &StatT***REMOVED***
		size: (*fi).Size(),
		mode: (*fi).Mode(),
		mtim: (*fi).ModTime()***REMOVED***, nil
***REMOVED***
