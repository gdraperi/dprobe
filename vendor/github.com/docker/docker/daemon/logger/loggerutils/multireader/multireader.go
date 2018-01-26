package multireader

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type pos struct ***REMOVED***
	idx    int
	offset int64
***REMOVED***

type multiReadSeeker struct ***REMOVED***
	readers []io.ReadSeeker
	pos     *pos
	posIdx  map[io.ReadSeeker]int
***REMOVED***

func (r *multiReadSeeker) Seek(offset int64, whence int) (int64, error) ***REMOVED***
	var tmpOffset int64
	switch whence ***REMOVED***
	case os.SEEK_SET:
		for i, rdr := range r.readers ***REMOVED***
			// get size of the current reader
			s, err := rdr.Seek(0, os.SEEK_END)
			if err != nil ***REMOVED***
				return -1, err
			***REMOVED***

			if offset > tmpOffset+s ***REMOVED***
				if i == len(r.readers)-1 ***REMOVED***
					rdrOffset := s + (offset - tmpOffset)
					if _, err := rdr.Seek(rdrOffset, os.SEEK_SET); err != nil ***REMOVED***
						return -1, err
					***REMOVED***
					r.pos = &pos***REMOVED***i, rdrOffset***REMOVED***
					return offset, nil
				***REMOVED***

				tmpOffset += s
				continue
			***REMOVED***

			rdrOffset := offset - tmpOffset
			idx := i

			if _, err := rdr.Seek(rdrOffset, os.SEEK_SET); err != nil ***REMOVED***
				return -1, err
			***REMOVED***
			// make sure all following readers are at 0
			for _, rdr := range r.readers[i+1:] ***REMOVED***
				rdr.Seek(0, os.SEEK_SET)
			***REMOVED***

			if rdrOffset == s && i != len(r.readers)-1 ***REMOVED***
				idx++
				rdrOffset = 0
			***REMOVED***
			r.pos = &pos***REMOVED***idx, rdrOffset***REMOVED***
			return offset, nil
		***REMOVED***
	case os.SEEK_END:
		for _, rdr := range r.readers ***REMOVED***
			s, err := rdr.Seek(0, os.SEEK_END)
			if err != nil ***REMOVED***
				return -1, err
			***REMOVED***
			tmpOffset += s
		***REMOVED***
		if _, err := r.Seek(tmpOffset+offset, os.SEEK_SET); err != nil ***REMOVED***
			return -1, err
		***REMOVED***
		return tmpOffset + offset, nil
	case os.SEEK_CUR:
		if r.pos == nil ***REMOVED***
			return r.Seek(offset, os.SEEK_SET)
		***REMOVED***
		// Just return the current offset
		if offset == 0 ***REMOVED***
			return r.getCurOffset()
		***REMOVED***

		curOffset, err := r.getCurOffset()
		if err != nil ***REMOVED***
			return -1, err
		***REMOVED***
		rdr, rdrOffset, err := r.getReaderForOffset(curOffset + offset)
		if err != nil ***REMOVED***
			return -1, err
		***REMOVED***

		r.pos = &pos***REMOVED***r.posIdx[rdr], rdrOffset***REMOVED***
		return curOffset + offset, nil
	default:
		return -1, fmt.Errorf("Invalid whence: %d", whence)
	***REMOVED***

	return -1, fmt.Errorf("Error seeking for whence: %d, offset: %d", whence, offset)
***REMOVED***

func (r *multiReadSeeker) getReaderForOffset(offset int64) (io.ReadSeeker, int64, error) ***REMOVED***

	var offsetTo int64

	for _, rdr := range r.readers ***REMOVED***
		size, err := getReadSeekerSize(rdr)
		if err != nil ***REMOVED***
			return nil, -1, err
		***REMOVED***
		if offsetTo+size > offset ***REMOVED***
			return rdr, offset - offsetTo, nil
		***REMOVED***
		if rdr == r.readers[len(r.readers)-1] ***REMOVED***
			return rdr, offsetTo + offset, nil
		***REMOVED***
		offsetTo += size
	***REMOVED***

	return nil, 0, nil
***REMOVED***

func (r *multiReadSeeker) getCurOffset() (int64, error) ***REMOVED***
	var totalSize int64
	for _, rdr := range r.readers[:r.pos.idx+1] ***REMOVED***
		if r.posIdx[rdr] == r.pos.idx ***REMOVED***
			totalSize += r.pos.offset
			break
		***REMOVED***

		size, err := getReadSeekerSize(rdr)
		if err != nil ***REMOVED***
			return -1, fmt.Errorf("error getting seeker size: %v", err)
		***REMOVED***
		totalSize += size
	***REMOVED***
	return totalSize, nil
***REMOVED***

func (r *multiReadSeeker) getOffsetToReader(rdr io.ReadSeeker) (int64, error) ***REMOVED***
	var offset int64
	for _, r := range r.readers ***REMOVED***
		if r == rdr ***REMOVED***
			break
		***REMOVED***

		size, err := getReadSeekerSize(rdr)
		if err != nil ***REMOVED***
			return -1, err
		***REMOVED***
		offset += size
	***REMOVED***
	return offset, nil
***REMOVED***

func (r *multiReadSeeker) Read(b []byte) (int, error) ***REMOVED***
	if r.pos == nil ***REMOVED***
		// make sure all readers are at 0
		r.Seek(0, os.SEEK_SET)
	***REMOVED***

	bLen := int64(len(b))
	buf := bytes.NewBuffer(nil)
	var rdr io.ReadSeeker

	for _, rdr = range r.readers[r.pos.idx:] ***REMOVED***
		readBytes, err := io.CopyN(buf, rdr, bLen)
		if err != nil && err != io.EOF ***REMOVED***
			return -1, err
		***REMOVED***
		bLen -= readBytes

		if bLen == 0 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	rdrPos, err := rdr.Seek(0, os.SEEK_CUR)
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***
	r.pos = &pos***REMOVED***r.posIdx[rdr], rdrPos***REMOVED***
	return buf.Read(b)
***REMOVED***

func getReadSeekerSize(rdr io.ReadSeeker) (int64, error) ***REMOVED***
	// save the current position
	pos, err := rdr.Seek(0, os.SEEK_CUR)
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***

	// get the size
	size, err := rdr.Seek(0, os.SEEK_END)
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***

	// reset the position
	if _, err := rdr.Seek(pos, os.SEEK_SET); err != nil ***REMOVED***
		return -1, err
	***REMOVED***
	return size, nil
***REMOVED***

// MultiReadSeeker returns a ReadSeeker that's the logical concatenation of the provided
// input readseekers. After calling this method the initial position is set to the
// beginning of the first ReadSeeker. At the end of a ReadSeeker, Read always advances
// to the beginning of the next ReadSeeker and returns EOF at the end of the last ReadSeeker.
// Seek can be used over the sum of lengths of all readseekers.
//
// When a MultiReadSeeker is used, no Read and Seek operations should be made on
// its ReadSeeker components. Also, users should make no assumption on the state
// of individual readseekers while the MultiReadSeeker is used.
func MultiReadSeeker(readers ...io.ReadSeeker) io.ReadSeeker ***REMOVED***
	if len(readers) == 1 ***REMOVED***
		return readers[0]
	***REMOVED***
	idx := make(map[io.ReadSeeker]int)
	for i, rdr := range readers ***REMOVED***
		idx[rdr] = i
	***REMOVED***
	return &multiReadSeeker***REMOVED***
		readers: readers,
		posIdx:  idx,
	***REMOVED***
***REMOVED***
