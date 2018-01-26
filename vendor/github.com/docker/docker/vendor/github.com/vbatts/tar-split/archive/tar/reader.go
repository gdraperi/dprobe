// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tar

// TODO(dsymonds):
//   - pax extensions

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ErrHeader = errors.New("archive/tar: invalid tar header")
)

const maxNanoSecondIntSize = 9

// A Reader provides sequential access to the contents of a tar archive.
// A tar archive consists of a sequence of files.
// The Next method advances to the next file in the archive (including the first),
// and then it can be treated as an io.Reader to access the file's data.
type Reader struct ***REMOVED***
	r       io.Reader
	err     error
	pad     int64           // amount of padding (ignored) after current file entry
	curr    numBytesReader  // reader for current file entry
	hdrBuff [blockSize]byte // buffer to use in readHeader

	RawAccounting bool          // Whether to enable the access needed to reassemble the tar from raw bytes. Some performance/memory hit for this.
	rawBytes      *bytes.Buffer // last raw bits
***REMOVED***

type parser struct ***REMOVED***
	err error // Last error seen
***REMOVED***

// RawBytes accesses the raw bytes of the archive, apart from the file payload itself.
// This includes the header and padding.
//
// This call resets the current rawbytes buffer
//
// Only when RawAccounting is enabled, otherwise this returns nil
func (tr *Reader) RawBytes() []byte ***REMOVED***
	if !tr.RawAccounting ***REMOVED***
		return nil
	***REMOVED***
	if tr.rawBytes == nil ***REMOVED***
		tr.rawBytes = bytes.NewBuffer(nil)
	***REMOVED***
	// if we've read them, then flush them.
	defer tr.rawBytes.Reset()
	return tr.rawBytes.Bytes()
***REMOVED***

// A numBytesReader is an io.Reader with a numBytes method, returning the number
// of bytes remaining in the underlying encoded data.
type numBytesReader interface ***REMOVED***
	io.Reader
	numBytes() int64
***REMOVED***

// A regFileReader is a numBytesReader for reading file data from a tar archive.
type regFileReader struct ***REMOVED***
	r  io.Reader // underlying reader
	nb int64     // number of unread bytes for current file entry
***REMOVED***

// A sparseFileReader is a numBytesReader for reading sparse file data from a
// tar archive.
type sparseFileReader struct ***REMOVED***
	rfr   numBytesReader // Reads the sparse-encoded file data
	sp    []sparseEntry  // The sparse map for the file
	pos   int64          // Keeps track of file position
	total int64          // Total size of the file
***REMOVED***

// A sparseEntry holds a single entry in a sparse file's sparse map.
//
// Sparse files are represented using a series of sparseEntrys.
// Despite the name, a sparseEntry represents an actual data fragment that
// references data found in the underlying archive stream. All regions not
// covered by a sparseEntry are logically filled with zeros.
//
// For example, if the underlying raw file contains the 10-byte data:
//	var compactData = "abcdefgh"
//
// And the sparse map has the following entries:
//	var sp = []sparseEntry***REMOVED***
//		***REMOVED***offset: 2,  numBytes: 5***REMOVED*** // Data fragment for [2..7]
//		***REMOVED***offset: 18, numBytes: 3***REMOVED*** // Data fragment for [18..21]
//	***REMOVED***
//
// Then the content of the resulting sparse file with a "real" size of 25 is:
//	var sparseData = "\x00"*2 + "abcde" + "\x00"*11 + "fgh" + "\x00"*4
type sparseEntry struct ***REMOVED***
	offset   int64 // Starting position of the fragment
	numBytes int64 // Length of the fragment
***REMOVED***

// Keywords for GNU sparse files in a PAX extended header
const (
	paxGNUSparseNumBlocks = "GNU.sparse.numblocks"
	paxGNUSparseOffset    = "GNU.sparse.offset"
	paxGNUSparseNumBytes  = "GNU.sparse.numbytes"
	paxGNUSparseMap       = "GNU.sparse.map"
	paxGNUSparseName      = "GNU.sparse.name"
	paxGNUSparseMajor     = "GNU.sparse.major"
	paxGNUSparseMinor     = "GNU.sparse.minor"
	paxGNUSparseSize      = "GNU.sparse.size"
	paxGNUSparseRealSize  = "GNU.sparse.realsize"
)

// Keywords for old GNU sparse headers
const (
	oldGNUSparseMainHeaderOffset               = 386
	oldGNUSparseMainHeaderIsExtendedOffset     = 482
	oldGNUSparseMainHeaderNumEntries           = 4
	oldGNUSparseExtendedHeaderIsExtendedOffset = 504
	oldGNUSparseExtendedHeaderNumEntries       = 21
	oldGNUSparseOffsetSize                     = 12
	oldGNUSparseNumBytesSize                   = 12
)

// NewReader creates a new Reader reading from r.
func NewReader(r io.Reader) *Reader ***REMOVED*** return &Reader***REMOVED***r: r***REMOVED*** ***REMOVED***

// Next advances to the next entry in the tar archive.
//
// io.EOF is returned at the end of the input.
func (tr *Reader) Next() (*Header, error) ***REMOVED***
	if tr.RawAccounting ***REMOVED***
		if tr.rawBytes == nil ***REMOVED***
			tr.rawBytes = bytes.NewBuffer(nil)
		***REMOVED*** else ***REMOVED***
			tr.rawBytes.Reset()
		***REMOVED***
	***REMOVED***

	if tr.err != nil ***REMOVED***
		return nil, tr.err
	***REMOVED***

	var hdr *Header
	var extHdrs map[string]string

	// Externally, Next iterates through the tar archive as if it is a series of
	// files. Internally, the tar format often uses fake "files" to add meta
	// data that describes the next file. These meta data "files" should not
	// normally be visible to the outside. As such, this loop iterates through
	// one or more "header files" until it finds a "normal file".
loop:
	for ***REMOVED***
		tr.err = tr.skipUnread()
		if tr.err != nil ***REMOVED***
			return nil, tr.err
		***REMOVED***

		hdr = tr.readHeader()
		if tr.err != nil ***REMOVED***
			return nil, tr.err
		***REMOVED***
		// Check for PAX/GNU special headers and files.
		switch hdr.Typeflag ***REMOVED***
		case TypeXHeader:
			extHdrs, tr.err = parsePAX(tr)
			if tr.err != nil ***REMOVED***
				return nil, tr.err
			***REMOVED***
			continue loop // This is a meta header affecting the next header
		case TypeGNULongName, TypeGNULongLink:
			var realname []byte
			realname, tr.err = ioutil.ReadAll(tr)
			if tr.err != nil ***REMOVED***
				return nil, tr.err
			***REMOVED***

			if tr.RawAccounting ***REMOVED***
				if _, tr.err = tr.rawBytes.Write(realname); tr.err != nil ***REMOVED***
					return nil, tr.err
				***REMOVED***
			***REMOVED***

			// Convert GNU extensions to use PAX headers.
			if extHdrs == nil ***REMOVED***
				extHdrs = make(map[string]string)
			***REMOVED***
			var p parser
			switch hdr.Typeflag ***REMOVED***
			case TypeGNULongName:
				extHdrs[paxPath] = p.parseString(realname)
			case TypeGNULongLink:
				extHdrs[paxLinkpath] = p.parseString(realname)
			***REMOVED***
			if p.err != nil ***REMOVED***
				tr.err = p.err
				return nil, tr.err
			***REMOVED***
			continue loop // This is a meta header affecting the next header
		default:
			mergePAX(hdr, extHdrs)

			// Check for a PAX format sparse file
			sp, err := tr.checkForGNUSparsePAXHeaders(hdr, extHdrs)
			if err != nil ***REMOVED***
				tr.err = err
				return nil, err
			***REMOVED***
			if sp != nil ***REMOVED***
				// Current file is a PAX format GNU sparse file.
				// Set the current file reader to a sparse file reader.
				tr.curr, tr.err = newSparseFileReader(tr.curr, sp, hdr.Size)
				if tr.err != nil ***REMOVED***
					return nil, tr.err
				***REMOVED***
			***REMOVED***
			break loop // This is a file, so stop
		***REMOVED***
	***REMOVED***
	return hdr, nil
***REMOVED***

// checkForGNUSparsePAXHeaders checks the PAX headers for GNU sparse headers. If they are found, then
// this function reads the sparse map and returns it. Unknown sparse formats are ignored, causing the file to
// be treated as a regular file.
func (tr *Reader) checkForGNUSparsePAXHeaders(hdr *Header, headers map[string]string) ([]sparseEntry, error) ***REMOVED***
	var sparseFormat string

	// Check for sparse format indicators
	major, majorOk := headers[paxGNUSparseMajor]
	minor, minorOk := headers[paxGNUSparseMinor]
	sparseName, sparseNameOk := headers[paxGNUSparseName]
	_, sparseMapOk := headers[paxGNUSparseMap]
	sparseSize, sparseSizeOk := headers[paxGNUSparseSize]
	sparseRealSize, sparseRealSizeOk := headers[paxGNUSparseRealSize]

	// Identify which, if any, sparse format applies from which PAX headers are set
	if majorOk && minorOk ***REMOVED***
		sparseFormat = major + "." + minor
	***REMOVED*** else if sparseNameOk && sparseMapOk ***REMOVED***
		sparseFormat = "0.1"
	***REMOVED*** else if sparseSizeOk ***REMOVED***
		sparseFormat = "0.0"
	***REMOVED*** else ***REMOVED***
		// Not a PAX format GNU sparse file.
		return nil, nil
	***REMOVED***

	// Check for unknown sparse format
	if sparseFormat != "0.0" && sparseFormat != "0.1" && sparseFormat != "1.0" ***REMOVED***
		return nil, nil
	***REMOVED***

	// Update hdr from GNU sparse PAX headers
	if sparseNameOk ***REMOVED***
		hdr.Name = sparseName
	***REMOVED***
	if sparseSizeOk ***REMOVED***
		realSize, err := strconv.ParseInt(sparseSize, 10, 0)
		if err != nil ***REMOVED***
			return nil, ErrHeader
		***REMOVED***
		hdr.Size = realSize
	***REMOVED*** else if sparseRealSizeOk ***REMOVED***
		realSize, err := strconv.ParseInt(sparseRealSize, 10, 0)
		if err != nil ***REMOVED***
			return nil, ErrHeader
		***REMOVED***
		hdr.Size = realSize
	***REMOVED***

	// Set up the sparse map, according to the particular sparse format in use
	var sp []sparseEntry
	var err error
	switch sparseFormat ***REMOVED***
	case "0.0", "0.1":
		sp, err = readGNUSparseMap0x1(headers)
	case "1.0":
		sp, err = readGNUSparseMap1x0(tr.curr)
	***REMOVED***
	return sp, err
***REMOVED***

// mergePAX merges well known headers according to PAX standard.
// In general headers with the same name as those found
// in the header struct overwrite those found in the header
// struct with higher precision or longer values. Esp. useful
// for name and linkname fields.
func mergePAX(hdr *Header, headers map[string]string) error ***REMOVED***
	for k, v := range headers ***REMOVED***
		switch k ***REMOVED***
		case paxPath:
			hdr.Name = v
		case paxLinkpath:
			hdr.Linkname = v
		case paxGname:
			hdr.Gname = v
		case paxUname:
			hdr.Uname = v
		case paxUid:
			uid, err := strconv.ParseInt(v, 10, 0)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			hdr.Uid = int(uid)
		case paxGid:
			gid, err := strconv.ParseInt(v, 10, 0)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			hdr.Gid = int(gid)
		case paxAtime:
			t, err := parsePAXTime(v)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			hdr.AccessTime = t
		case paxMtime:
			t, err := parsePAXTime(v)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			hdr.ModTime = t
		case paxCtime:
			t, err := parsePAXTime(v)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			hdr.ChangeTime = t
		case paxSize:
			size, err := strconv.ParseInt(v, 10, 0)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			hdr.Size = int64(size)
		default:
			if strings.HasPrefix(k, paxXattr) ***REMOVED***
				if hdr.Xattrs == nil ***REMOVED***
					hdr.Xattrs = make(map[string]string)
				***REMOVED***
				hdr.Xattrs[k[len(paxXattr):]] = v
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// parsePAXTime takes a string of the form %d.%d as described in
// the PAX specification.
func parsePAXTime(t string) (time.Time, error) ***REMOVED***
	buf := []byte(t)
	pos := bytes.IndexByte(buf, '.')
	var seconds, nanoseconds int64
	var err error
	if pos == -1 ***REMOVED***
		seconds, err = strconv.ParseInt(t, 10, 0)
		if err != nil ***REMOVED***
			return time.Time***REMOVED******REMOVED***, err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		seconds, err = strconv.ParseInt(string(buf[:pos]), 10, 0)
		if err != nil ***REMOVED***
			return time.Time***REMOVED******REMOVED***, err
		***REMOVED***
		nano_buf := string(buf[pos+1:])
		// Pad as needed before converting to a decimal.
		// For example .030 -> .030000000 -> 30000000 nanoseconds
		if len(nano_buf) < maxNanoSecondIntSize ***REMOVED***
			// Right pad
			nano_buf += strings.Repeat("0", maxNanoSecondIntSize-len(nano_buf))
		***REMOVED*** else if len(nano_buf) > maxNanoSecondIntSize ***REMOVED***
			// Right truncate
			nano_buf = nano_buf[:maxNanoSecondIntSize]
		***REMOVED***
		nanoseconds, err = strconv.ParseInt(string(nano_buf), 10, 0)
		if err != nil ***REMOVED***
			return time.Time***REMOVED******REMOVED***, err
		***REMOVED***
	***REMOVED***
	ts := time.Unix(seconds, nanoseconds)
	return ts, nil
***REMOVED***

// parsePAX parses PAX headers.
// If an extended header (type 'x') is invalid, ErrHeader is returned
func parsePAX(r io.Reader) (map[string]string, error) ***REMOVED***
	buf, err := ioutil.ReadAll(r)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// leaving this function for io.Reader makes it more testable
	if tr, ok := r.(*Reader); ok && tr.RawAccounting ***REMOVED***
		if _, err = tr.rawBytes.Write(buf); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	sbuf := string(buf)

	// For GNU PAX sparse format 0.0 support.
	// This function transforms the sparse format 0.0 headers into sparse format 0.1 headers.
	var sparseMap bytes.Buffer

	headers := make(map[string]string)
	// Each record is constructed as
	//     "%d %s=%s\n", length, keyword, value
	for len(sbuf) > 0 ***REMOVED***
		key, value, residual, err := parsePAXRecord(sbuf)
		if err != nil ***REMOVED***
			return nil, ErrHeader
		***REMOVED***
		sbuf = residual

		keyStr := string(key)
		if keyStr == paxGNUSparseOffset || keyStr == paxGNUSparseNumBytes ***REMOVED***
			// GNU sparse format 0.0 special key. Write to sparseMap instead of using the headers map.
			sparseMap.WriteString(value)
			sparseMap.Write([]byte***REMOVED***','***REMOVED***)
		***REMOVED*** else ***REMOVED***
			// Normal key. Set the value in the headers map.
			headers[keyStr] = string(value)
		***REMOVED***
	***REMOVED***
	if sparseMap.Len() != 0 ***REMOVED***
		// Add sparse info to headers, chopping off the extra comma
		sparseMap.Truncate(sparseMap.Len() - 1)
		headers[paxGNUSparseMap] = sparseMap.String()
	***REMOVED***
	return headers, nil
***REMOVED***

// parsePAXRecord parses the input PAX record string into a key-value pair.
// If parsing is successful, it will slice off the currently read record and
// return the remainder as r.
//
// A PAX record is of the following form:
//	"%d %s=%s\n" % (size, key, value)
func parsePAXRecord(s string) (k, v, r string, err error) ***REMOVED***
	// The size field ends at the first space.
	sp := strings.IndexByte(s, ' ')
	if sp == -1 ***REMOVED***
		return "", "", s, ErrHeader
	***REMOVED***

	// Parse the first token as a decimal integer.
	n, perr := strconv.ParseInt(s[:sp], 10, 0) // Intentionally parse as native int
	if perr != nil || n < 5 || int64(len(s)) < n ***REMOVED***
		return "", "", s, ErrHeader
	***REMOVED***

	// Extract everything between the space and the final newline.
	rec, nl, rem := s[sp+1:n-1], s[n-1:n], s[n:]
	if nl != "\n" ***REMOVED***
		return "", "", s, ErrHeader
	***REMOVED***

	// The first equals separates the key from the value.
	eq := strings.IndexByte(rec, '=')
	if eq == -1 ***REMOVED***
		return "", "", s, ErrHeader
	***REMOVED***
	return rec[:eq], rec[eq+1:], rem, nil
***REMOVED***

// parseString parses bytes as a NUL-terminated C-style string.
// If a NUL byte is not found then the whole slice is returned as a string.
func (*parser) parseString(b []byte) string ***REMOVED***
	n := 0
	for n < len(b) && b[n] != 0 ***REMOVED***
		n++
	***REMOVED***
	return string(b[0:n])
***REMOVED***

// parseNumeric parses the input as being encoded in either base-256 or octal.
// This function may return negative numbers.
// If parsing fails or an integer overflow occurs, err will be set.
func (p *parser) parseNumeric(b []byte) int64 ***REMOVED***
	// Check for base-256 (binary) format first.
	// If the first bit is set, then all following bits constitute a two's
	// complement encoded number in big-endian byte order.
	if len(b) > 0 && b[0]&0x80 != 0 ***REMOVED***
		// Handling negative numbers relies on the following identity:
		//	-a-1 == ^a
		//
		// If the number is negative, we use an inversion mask to invert the
		// data bytes and treat the value as an unsigned number.
		var inv byte // 0x00 if positive or zero, 0xff if negative
		if b[0]&0x40 != 0 ***REMOVED***
			inv = 0xff
		***REMOVED***

		var x uint64
		for i, c := range b ***REMOVED***
			c ^= inv // Inverts c only if inv is 0xff, otherwise does nothing
			if i == 0 ***REMOVED***
				c &= 0x7f // Ignore signal bit in first byte
			***REMOVED***
			if (x >> 56) > 0 ***REMOVED***
				p.err = ErrHeader // Integer overflow
				return 0
			***REMOVED***
			x = x<<8 | uint64(c)
		***REMOVED***
		if (x >> 63) > 0 ***REMOVED***
			p.err = ErrHeader // Integer overflow
			return 0
		***REMOVED***
		if inv == 0xff ***REMOVED***
			return ^int64(x)
		***REMOVED***
		return int64(x)
	***REMOVED***

	// Normal case is base-8 (octal) format.
	return p.parseOctal(b)
***REMOVED***

func (p *parser) parseOctal(b []byte) int64 ***REMOVED***
	// Because unused fields are filled with NULs, we need
	// to skip leading NULs. Fields may also be padded with
	// spaces or NULs.
	// So we remove leading and trailing NULs and spaces to
	// be sure.
	b = bytes.Trim(b, " \x00")

	if len(b) == 0 ***REMOVED***
		return 0
	***REMOVED***
	x, perr := strconv.ParseUint(p.parseString(b), 8, 64)
	if perr != nil ***REMOVED***
		p.err = ErrHeader
	***REMOVED***
	return int64(x)
***REMOVED***

// skipUnread skips any unread bytes in the existing file entry, as well as any
// alignment padding. It returns io.ErrUnexpectedEOF if any io.EOF is
// encountered in the data portion; it is okay to hit io.EOF in the padding.
//
// Note that this function still works properly even when sparse files are being
// used since numBytes returns the bytes remaining in the underlying io.Reader.
func (tr *Reader) skipUnread() error ***REMOVED***
	dataSkip := tr.numBytes()      // Number of data bytes to skip
	totalSkip := dataSkip + tr.pad // Total number of bytes to skip
	tr.curr, tr.pad = nil, 0
	if tr.RawAccounting ***REMOVED***
		_, tr.err = io.CopyN(tr.rawBytes, tr.r, totalSkip)
		return tr.err
	***REMOVED***
	// If possible, Seek to the last byte before the end of the data section.
	// Do this because Seek is often lazy about reporting errors; this will mask
	// the fact that the tar stream may be truncated. We can rely on the
	// io.CopyN done shortly afterwards to trigger any IO errors.
	var seekSkipped int64 // Number of bytes skipped via Seek
	if sr, ok := tr.r.(io.Seeker); ok && dataSkip > 1 ***REMOVED***
		// Not all io.Seeker can actually Seek. For example, os.Stdin implements
		// io.Seeker, but calling Seek always returns an error and performs
		// no action. Thus, we try an innocent seek to the current position
		// to see if Seek is really supported.
		pos1, err := sr.Seek(0, os.SEEK_CUR)
		if err == nil ***REMOVED***
			// Seek seems supported, so perform the real Seek.
			pos2, err := sr.Seek(dataSkip-1, os.SEEK_CUR)
			if err != nil ***REMOVED***
				tr.err = err
				return tr.err
			***REMOVED***
			seekSkipped = pos2 - pos1
		***REMOVED***
	***REMOVED***

	var copySkipped int64 // Number of bytes skipped via CopyN
	copySkipped, tr.err = io.CopyN(ioutil.Discard, tr.r, totalSkip-seekSkipped)
	if tr.err == io.EOF && seekSkipped+copySkipped < dataSkip ***REMOVED***
		tr.err = io.ErrUnexpectedEOF
	***REMOVED***
	return tr.err
***REMOVED***

func (tr *Reader) verifyChecksum(header []byte) bool ***REMOVED***
	if tr.err != nil ***REMOVED***
		return false
	***REMOVED***

	var p parser
	given := p.parseOctal(header[148:156])
	unsigned, signed := checksum(header)
	return p.err == nil && (given == unsigned || given == signed)
***REMOVED***

// readHeader reads the next block header and assumes that the underlying reader
// is already aligned to a block boundary.
//
// The err will be set to io.EOF only when one of the following occurs:
//	* Exactly 0 bytes are read and EOF is hit.
//	* Exactly 1 block of zeros is read and EOF is hit.
//	* At least 2 blocks of zeros are read.
func (tr *Reader) readHeader() *Header ***REMOVED***
	header := tr.hdrBuff[:]
	copy(header, zeroBlock)

	if n, err := io.ReadFull(tr.r, header); err != nil ***REMOVED***
		tr.err = err
		// because it could read some of the block, but reach EOF first
		if tr.err == io.EOF && tr.RawAccounting ***REMOVED***
			if _, err := tr.rawBytes.Write(header[:n]); err != nil ***REMOVED***
				tr.err = err
			***REMOVED***
		***REMOVED***
		return nil // io.EOF is okay here
	***REMOVED***
	if tr.RawAccounting ***REMOVED***
		if _, tr.err = tr.rawBytes.Write(header); tr.err != nil ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	// Two blocks of zero bytes marks the end of the archive.
	if bytes.Equal(header, zeroBlock[0:blockSize]) ***REMOVED***
		if n, err := io.ReadFull(tr.r, header); err != nil ***REMOVED***
			tr.err = err
			// because it could read some of the block, but reach EOF first
			if tr.err == io.EOF && tr.RawAccounting ***REMOVED***
				if _, err := tr.rawBytes.Write(header[:n]); err != nil ***REMOVED***
					tr.err = err
				***REMOVED***
			***REMOVED***
			return nil // io.EOF is okay here
		***REMOVED***
		if tr.RawAccounting ***REMOVED***
			if _, tr.err = tr.rawBytes.Write(header); tr.err != nil ***REMOVED***
				return nil
			***REMOVED***
		***REMOVED***
		if bytes.Equal(header, zeroBlock[0:blockSize]) ***REMOVED***
			tr.err = io.EOF
		***REMOVED*** else ***REMOVED***
			tr.err = ErrHeader // zero block and then non-zero block
		***REMOVED***
		return nil
	***REMOVED***

	if !tr.verifyChecksum(header) ***REMOVED***
		tr.err = ErrHeader
		return nil
	***REMOVED***

	// Unpack
	var p parser
	hdr := new(Header)
	s := slicer(header)

	hdr.Name = p.parseString(s.next(100))
	hdr.Mode = p.parseNumeric(s.next(8))
	hdr.Uid = int(p.parseNumeric(s.next(8)))
	hdr.Gid = int(p.parseNumeric(s.next(8)))
	hdr.Size = p.parseNumeric(s.next(12))
	hdr.ModTime = time.Unix(p.parseNumeric(s.next(12)), 0)
	s.next(8) // chksum
	hdr.Typeflag = s.next(1)[0]
	hdr.Linkname = p.parseString(s.next(100))

	// The remainder of the header depends on the value of magic.
	// The original (v7) version of tar had no explicit magic field,
	// so its magic bytes, like the rest of the block, are NULs.
	magic := string(s.next(8)) // contains version field as well.
	var format string
	switch ***REMOVED***
	case magic[:6] == "ustar\x00": // POSIX tar (1003.1-1988)
		if string(header[508:512]) == "tar\x00" ***REMOVED***
			format = "star"
		***REMOVED*** else ***REMOVED***
			format = "posix"
		***REMOVED***
	case magic == "ustar  \x00": // old GNU tar
		format = "gnu"
	***REMOVED***

	switch format ***REMOVED***
	case "posix", "gnu", "star":
		hdr.Uname = p.parseString(s.next(32))
		hdr.Gname = p.parseString(s.next(32))
		devmajor := s.next(8)
		devminor := s.next(8)
		if hdr.Typeflag == TypeChar || hdr.Typeflag == TypeBlock ***REMOVED***
			hdr.Devmajor = p.parseNumeric(devmajor)
			hdr.Devminor = p.parseNumeric(devminor)
		***REMOVED***
		var prefix string
		switch format ***REMOVED***
		case "posix", "gnu":
			prefix = p.parseString(s.next(155))
		case "star":
			prefix = p.parseString(s.next(131))
			hdr.AccessTime = time.Unix(p.parseNumeric(s.next(12)), 0)
			hdr.ChangeTime = time.Unix(p.parseNumeric(s.next(12)), 0)
		***REMOVED***
		if len(prefix) > 0 ***REMOVED***
			hdr.Name = prefix + "/" + hdr.Name
		***REMOVED***
	***REMOVED***

	if p.err != nil ***REMOVED***
		tr.err = p.err
		return nil
	***REMOVED***

	nb := hdr.Size
	if isHeaderOnlyType(hdr.Typeflag) ***REMOVED***
		nb = 0
	***REMOVED***
	if nb < 0 ***REMOVED***
		tr.err = ErrHeader
		return nil
	***REMOVED***

	// Set the current file reader.
	tr.pad = -nb & (blockSize - 1) // blockSize is a power of two
	tr.curr = &regFileReader***REMOVED***r: tr.r, nb: nb***REMOVED***

	// Check for old GNU sparse format entry.
	if hdr.Typeflag == TypeGNUSparse ***REMOVED***
		// Get the real size of the file.
		hdr.Size = p.parseNumeric(header[483:495])
		if p.err != nil ***REMOVED***
			tr.err = p.err
			return nil
		***REMOVED***

		// Read the sparse map.
		sp := tr.readOldGNUSparseMap(header)
		if tr.err != nil ***REMOVED***
			return nil
		***REMOVED***

		// Current file is a GNU sparse file. Update the current file reader.
		tr.curr, tr.err = newSparseFileReader(tr.curr, sp, hdr.Size)
		if tr.err != nil ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	return hdr
***REMOVED***

// readOldGNUSparseMap reads the sparse map as stored in the old GNU sparse format.
// The sparse map is stored in the tar header if it's small enough. If it's larger than four entries,
// then one or more extension headers are used to store the rest of the sparse map.
func (tr *Reader) readOldGNUSparseMap(header []byte) []sparseEntry ***REMOVED***
	var p parser
	isExtended := header[oldGNUSparseMainHeaderIsExtendedOffset] != 0
	spCap := oldGNUSparseMainHeaderNumEntries
	if isExtended ***REMOVED***
		spCap += oldGNUSparseExtendedHeaderNumEntries
	***REMOVED***
	sp := make([]sparseEntry, 0, spCap)
	s := slicer(header[oldGNUSparseMainHeaderOffset:])

	// Read the four entries from the main tar header
	for i := 0; i < oldGNUSparseMainHeaderNumEntries; i++ ***REMOVED***
		offset := p.parseNumeric(s.next(oldGNUSparseOffsetSize))
		numBytes := p.parseNumeric(s.next(oldGNUSparseNumBytesSize))
		if p.err != nil ***REMOVED***
			tr.err = p.err
			return nil
		***REMOVED***
		if offset == 0 && numBytes == 0 ***REMOVED***
			break
		***REMOVED***
		sp = append(sp, sparseEntry***REMOVED***offset: offset, numBytes: numBytes***REMOVED***)
	***REMOVED***

	for isExtended ***REMOVED***
		// There are more entries. Read an extension header and parse its entries.
		sparseHeader := make([]byte, blockSize)
		if _, tr.err = io.ReadFull(tr.r, sparseHeader); tr.err != nil ***REMOVED***
			return nil
		***REMOVED***
		if tr.RawAccounting ***REMOVED***
			if _, tr.err = tr.rawBytes.Write(sparseHeader); tr.err != nil ***REMOVED***
				return nil
			***REMOVED***
		***REMOVED***

		isExtended = sparseHeader[oldGNUSparseExtendedHeaderIsExtendedOffset] != 0
		s = slicer(sparseHeader)
		for i := 0; i < oldGNUSparseExtendedHeaderNumEntries; i++ ***REMOVED***
			offset := p.parseNumeric(s.next(oldGNUSparseOffsetSize))
			numBytes := p.parseNumeric(s.next(oldGNUSparseNumBytesSize))
			if p.err != nil ***REMOVED***
				tr.err = p.err
				return nil
			***REMOVED***
			if offset == 0 && numBytes == 0 ***REMOVED***
				break
			***REMOVED***
			sp = append(sp, sparseEntry***REMOVED***offset: offset, numBytes: numBytes***REMOVED***)
		***REMOVED***
	***REMOVED***
	return sp
***REMOVED***

// readGNUSparseMap1x0 reads the sparse map as stored in GNU's PAX sparse format
// version 1.0. The format of the sparse map consists of a series of
// newline-terminated numeric fields. The first field is the number of entries
// and is always present. Following this are the entries, consisting of two
// fields (offset, numBytes). This function must stop reading at the end
// boundary of the block containing the last newline.
//
// Note that the GNU manual says that numeric values should be encoded in octal
// format. However, the GNU tar utility itself outputs these values in decimal.
// As such, this library treats values as being encoded in decimal.
func readGNUSparseMap1x0(r io.Reader) ([]sparseEntry, error) ***REMOVED***
	var cntNewline int64
	var buf bytes.Buffer
	var blk = make([]byte, blockSize)

	// feedTokens copies data in numBlock chunks from r into buf until there are
	// at least cnt newlines in buf. It will not read more blocks than needed.
	var feedTokens = func(cnt int64) error ***REMOVED***
		for cntNewline < cnt ***REMOVED***
			if _, err := io.ReadFull(r, blk); err != nil ***REMOVED***
				if err == io.EOF ***REMOVED***
					err = io.ErrUnexpectedEOF
				***REMOVED***
				return err
			***REMOVED***
			buf.Write(blk)
			for _, c := range blk ***REMOVED***
				if c == '\n' ***REMOVED***
					cntNewline++
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***

	// nextToken gets the next token delimited by a newline. This assumes that
	// at least one newline exists in the buffer.
	var nextToken = func() string ***REMOVED***
		cntNewline--
		tok, _ := buf.ReadString('\n')
		return tok[:len(tok)-1] // Cut off newline
	***REMOVED***

	// Parse for the number of entries.
	// Use integer overflow resistant math to check this.
	if err := feedTokens(1); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	numEntries, err := strconv.ParseInt(nextToken(), 10, 0) // Intentionally parse as native int
	if err != nil || numEntries < 0 || int(2*numEntries) < int(numEntries) ***REMOVED***
		return nil, ErrHeader
	***REMOVED***

	// Parse for all member entries.
	// numEntries is trusted after this since a potential attacker must have
	// committed resources proportional to what this library used.
	if err := feedTokens(2 * numEntries); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	sp := make([]sparseEntry, 0, numEntries)
	for i := int64(0); i < numEntries; i++ ***REMOVED***
		offset, err := strconv.ParseInt(nextToken(), 10, 64)
		if err != nil ***REMOVED***
			return nil, ErrHeader
		***REMOVED***
		numBytes, err := strconv.ParseInt(nextToken(), 10, 64)
		if err != nil ***REMOVED***
			return nil, ErrHeader
		***REMOVED***
		sp = append(sp, sparseEntry***REMOVED***offset: offset, numBytes: numBytes***REMOVED***)
	***REMOVED***
	return sp, nil
***REMOVED***

// readGNUSparseMap0x1 reads the sparse map as stored in GNU's PAX sparse format
// version 0.1. The sparse map is stored in the PAX headers.
func readGNUSparseMap0x1(extHdrs map[string]string) ([]sparseEntry, error) ***REMOVED***
	// Get number of entries.
	// Use integer overflow resistant math to check this.
	numEntriesStr := extHdrs[paxGNUSparseNumBlocks]
	numEntries, err := strconv.ParseInt(numEntriesStr, 10, 0) // Intentionally parse as native int
	if err != nil || numEntries < 0 || int(2*numEntries) < int(numEntries) ***REMOVED***
		return nil, ErrHeader
	***REMOVED***

	// There should be two numbers in sparseMap for each entry.
	sparseMap := strings.Split(extHdrs[paxGNUSparseMap], ",")
	if int64(len(sparseMap)) != 2*numEntries ***REMOVED***
		return nil, ErrHeader
	***REMOVED***

	// Loop through the entries in the sparse map.
	// numEntries is trusted now.
	sp := make([]sparseEntry, 0, numEntries)
	for i := int64(0); i < numEntries; i++ ***REMOVED***
		offset, err := strconv.ParseInt(sparseMap[2*i], 10, 64)
		if err != nil ***REMOVED***
			return nil, ErrHeader
		***REMOVED***
		numBytes, err := strconv.ParseInt(sparseMap[2*i+1], 10, 64)
		if err != nil ***REMOVED***
			return nil, ErrHeader
		***REMOVED***
		sp = append(sp, sparseEntry***REMOVED***offset: offset, numBytes: numBytes***REMOVED***)
	***REMOVED***
	return sp, nil
***REMOVED***

// numBytes returns the number of bytes left to read in the current file's entry
// in the tar archive, or 0 if there is no current file.
func (tr *Reader) numBytes() int64 ***REMOVED***
	if tr.curr == nil ***REMOVED***
		// No current file, so no bytes
		return 0
	***REMOVED***
	return tr.curr.numBytes()
***REMOVED***

// Read reads from the current entry in the tar archive.
// It returns 0, io.EOF when it reaches the end of that entry,
// until Next is called to advance to the next entry.
//
// Calling Read on special types like TypeLink, TypeSymLink, TypeChar,
// TypeBlock, TypeDir, and TypeFifo returns 0, io.EOF regardless of what
// the Header.Size claims.
func (tr *Reader) Read(b []byte) (n int, err error) ***REMOVED***
	if tr.err != nil ***REMOVED***
		return 0, tr.err
	***REMOVED***
	if tr.curr == nil ***REMOVED***
		return 0, io.EOF
	***REMOVED***

	n, err = tr.curr.Read(b)
	if err != nil && err != io.EOF ***REMOVED***
		tr.err = err
	***REMOVED***
	return
***REMOVED***

func (rfr *regFileReader) Read(b []byte) (n int, err error) ***REMOVED***
	if rfr.nb == 0 ***REMOVED***
		// file consumed
		return 0, io.EOF
	***REMOVED***
	if int64(len(b)) > rfr.nb ***REMOVED***
		b = b[0:rfr.nb]
	***REMOVED***
	n, err = rfr.r.Read(b)
	rfr.nb -= int64(n)

	if err == io.EOF && rfr.nb > 0 ***REMOVED***
		err = io.ErrUnexpectedEOF
	***REMOVED***
	return
***REMOVED***

// numBytes returns the number of bytes left to read in the file's data in the tar archive.
func (rfr *regFileReader) numBytes() int64 ***REMOVED***
	return rfr.nb
***REMOVED***

// newSparseFileReader creates a new sparseFileReader, but validates all of the
// sparse entries before doing so.
func newSparseFileReader(rfr numBytesReader, sp []sparseEntry, total int64) (*sparseFileReader, error) ***REMOVED***
	if total < 0 ***REMOVED***
		return nil, ErrHeader // Total size cannot be negative
	***REMOVED***

	// Validate all sparse entries. These are the same checks as performed by
	// the BSD tar utility.
	for i, s := range sp ***REMOVED***
		switch ***REMOVED***
		case s.offset < 0 || s.numBytes < 0:
			return nil, ErrHeader // Negative values are never okay
		case s.offset > math.MaxInt64-s.numBytes:
			return nil, ErrHeader // Integer overflow with large length
		case s.offset+s.numBytes > total:
			return nil, ErrHeader // Region extends beyond the "real" size
		case i > 0 && sp[i-1].offset+sp[i-1].numBytes > s.offset:
			return nil, ErrHeader // Regions can't overlap and must be in order
		***REMOVED***
	***REMOVED***
	return &sparseFileReader***REMOVED***rfr: rfr, sp: sp, total: total***REMOVED***, nil
***REMOVED***

// readHole reads a sparse hole ending at endOffset.
func (sfr *sparseFileReader) readHole(b []byte, endOffset int64) int ***REMOVED***
	n64 := endOffset - sfr.pos
	if n64 > int64(len(b)) ***REMOVED***
		n64 = int64(len(b))
	***REMOVED***
	n := int(n64)
	for i := 0; i < n; i++ ***REMOVED***
		b[i] = 0
	***REMOVED***
	sfr.pos += n64
	return n
***REMOVED***

// Read reads the sparse file data in expanded form.
func (sfr *sparseFileReader) Read(b []byte) (n int, err error) ***REMOVED***
	// Skip past all empty fragments.
	for len(sfr.sp) > 0 && sfr.sp[0].numBytes == 0 ***REMOVED***
		sfr.sp = sfr.sp[1:]
	***REMOVED***

	// If there are no more fragments, then it is possible that there
	// is one last sparse hole.
	if len(sfr.sp) == 0 ***REMOVED***
		// This behavior matches the BSD tar utility.
		// However, GNU tar stops returning data even if sfr.total is unmet.
		if sfr.pos < sfr.total ***REMOVED***
			return sfr.readHole(b, sfr.total), nil
		***REMOVED***
		return 0, io.EOF
	***REMOVED***

	// In front of a data fragment, so read a hole.
	if sfr.pos < sfr.sp[0].offset ***REMOVED***
		return sfr.readHole(b, sfr.sp[0].offset), nil
	***REMOVED***

	// In a data fragment, so read from it.
	// This math is overflow free since we verify that offset and numBytes can
	// be safely added when creating the sparseFileReader.
	endPos := sfr.sp[0].offset + sfr.sp[0].numBytes // End offset of fragment
	bytesLeft := endPos - sfr.pos                   // Bytes left in fragment
	if int64(len(b)) > bytesLeft ***REMOVED***
		b = b[:bytesLeft]
	***REMOVED***

	n, err = sfr.rfr.Read(b)
	sfr.pos += int64(n)
	if err == io.EOF ***REMOVED***
		if sfr.pos < endPos ***REMOVED***
			err = io.ErrUnexpectedEOF // There was supposed to be more data
		***REMOVED*** else if sfr.pos < sfr.total ***REMOVED***
			err = nil // There is still an implicit sparse hole at the end
		***REMOVED***
	***REMOVED***

	if sfr.pos == endPos ***REMOVED***
		sfr.sp = sfr.sp[1:] // We are done with this fragment, so pop it
	***REMOVED***
	return n, err
***REMOVED***

// numBytes returns the number of bytes left to read in the sparse file's
// sparse-encoded data in the tar archive.
func (sfr *sparseFileReader) numBytes() int64 ***REMOVED***
	return sfr.rfr.numBytes()
***REMOVED***
