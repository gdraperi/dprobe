// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tar

// TODO(dsymonds):
// - catch more errors (no first header, etc.)

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	ErrWriteTooLong    = errors.New("archive/tar: write too long")
	ErrFieldTooLong    = errors.New("archive/tar: header field too long")
	ErrWriteAfterClose = errors.New("archive/tar: write after close")
	errInvalidHeader   = errors.New("archive/tar: header field too long or contains invalid values")
)

// A Writer provides sequential writing of a tar archive in POSIX.1 format.
// A tar archive consists of a sequence of files.
// Call WriteHeader to begin a new file, and then call Write to supply that file's data,
// writing at most hdr.Size bytes in total.
type Writer struct ***REMOVED***
	w          io.Writer
	err        error
	nb         int64 // number of unwritten bytes for current file entry
	pad        int64 // amount of padding to write after current file entry
	closed     bool
	usedBinary bool            // whether the binary numeric field extension was used
	preferPax  bool            // use pax header instead of binary numeric header
	hdrBuff    [blockSize]byte // buffer to use in writeHeader when writing a regular header
	paxHdrBuff [blockSize]byte // buffer to use in writeHeader when writing a pax header
***REMOVED***

type formatter struct ***REMOVED***
	err error // Last error seen
***REMOVED***

// NewWriter creates a new Writer writing to w.
func NewWriter(w io.Writer) *Writer ***REMOVED*** return &Writer***REMOVED***w: w***REMOVED*** ***REMOVED***

// Flush finishes writing the current file (optional).
func (tw *Writer) Flush() error ***REMOVED***
	if tw.nb > 0 ***REMOVED***
		tw.err = fmt.Errorf("archive/tar: missed writing %d bytes", tw.nb)
		return tw.err
	***REMOVED***

	n := tw.nb + tw.pad
	for n > 0 && tw.err == nil ***REMOVED***
		nr := n
		if nr > blockSize ***REMOVED***
			nr = blockSize
		***REMOVED***
		var nw int
		nw, tw.err = tw.w.Write(zeroBlock[0:nr])
		n -= int64(nw)
	***REMOVED***
	tw.nb = 0
	tw.pad = 0
	return tw.err
***REMOVED***

// Write s into b, terminating it with a NUL if there is room.
func (f *formatter) formatString(b []byte, s string) ***REMOVED***
	if len(s) > len(b) ***REMOVED***
		f.err = ErrFieldTooLong
		return
	***REMOVED***
	ascii := toASCII(s)
	copy(b, ascii)
	if len(ascii) < len(b) ***REMOVED***
		b[len(ascii)] = 0
	***REMOVED***
***REMOVED***

// Encode x as an octal ASCII string and write it into b with leading zeros.
func (f *formatter) formatOctal(b []byte, x int64) ***REMOVED***
	s := strconv.FormatInt(x, 8)
	// leading zeros, but leave room for a NUL.
	for len(s)+1 < len(b) ***REMOVED***
		s = "0" + s
	***REMOVED***
	f.formatString(b, s)
***REMOVED***

// fitsInBase256 reports whether x can be encoded into n bytes using base-256
// encoding. Unlike octal encoding, base-256 encoding does not require that the
// string ends with a NUL character. Thus, all n bytes are available for output.
//
// If operating in binary mode, this assumes strict GNU binary mode; which means
// that the first byte can only be either 0x80 or 0xff. Thus, the first byte is
// equivalent to the sign bit in two's complement form.
func fitsInBase256(n int, x int64) bool ***REMOVED***
	var binBits = uint(n-1) * 8
	return n >= 9 || (x >= -1<<binBits && x < 1<<binBits)
***REMOVED***

// Write x into b, as binary (GNUtar/star extension).
func (f *formatter) formatNumeric(b []byte, x int64) ***REMOVED***
	if fitsInBase256(len(b), x) ***REMOVED***
		for i := len(b) - 1; i >= 0; i-- ***REMOVED***
			b[i] = byte(x)
			x >>= 8
		***REMOVED***
		b[0] |= 0x80 // Highest bit indicates binary format
		return
	***REMOVED***

	f.formatOctal(b, 0) // Last resort, just write zero
	f.err = ErrFieldTooLong
***REMOVED***

var (
	minTime = time.Unix(0, 0)
	// There is room for 11 octal digits (33 bits) of mtime.
	maxTime = minTime.Add((1<<33 - 1) * time.Second)
)

// WriteHeader writes hdr and prepares to accept the file's contents.
// WriteHeader calls Flush if it is not the first header.
// Calling after a Close will return ErrWriteAfterClose.
func (tw *Writer) WriteHeader(hdr *Header) error ***REMOVED***
	return tw.writeHeader(hdr, true)
***REMOVED***

// WriteHeader writes hdr and prepares to accept the file's contents.
// WriteHeader calls Flush if it is not the first header.
// Calling after a Close will return ErrWriteAfterClose.
// As this method is called internally by writePax header to allow it to
// suppress writing the pax header.
func (tw *Writer) writeHeader(hdr *Header, allowPax bool) error ***REMOVED***
	if tw.closed ***REMOVED***
		return ErrWriteAfterClose
	***REMOVED***
	if tw.err == nil ***REMOVED***
		tw.Flush()
	***REMOVED***
	if tw.err != nil ***REMOVED***
		return tw.err
	***REMOVED***

	// a map to hold pax header records, if any are needed
	paxHeaders := make(map[string]string)

	// TODO(shanemhansen): we might want to use PAX headers for
	// subsecond time resolution, but for now let's just capture
	// too long fields or non ascii characters

	var f formatter
	var header []byte

	// We need to select which scratch buffer to use carefully,
	// since this method is called recursively to write PAX headers.
	// If allowPax is true, this is the non-recursive call, and we will use hdrBuff.
	// If allowPax is false, we are being called by writePAXHeader, and hdrBuff is
	// already being used by the non-recursive call, so we must use paxHdrBuff.
	header = tw.hdrBuff[:]
	if !allowPax ***REMOVED***
		header = tw.paxHdrBuff[:]
	***REMOVED***
	copy(header, zeroBlock)
	s := slicer(header)

	// Wrappers around formatter that automatically sets paxHeaders if the
	// argument extends beyond the capacity of the input byte slice.
	var formatString = func(b []byte, s string, paxKeyword string) ***REMOVED***
		needsPaxHeader := paxKeyword != paxNone && len(s) > len(b) || !isASCII(s)
		if needsPaxHeader ***REMOVED***
			paxHeaders[paxKeyword] = s
			return
		***REMOVED***
		f.formatString(b, s)
	***REMOVED***
	var formatNumeric = func(b []byte, x int64, paxKeyword string) ***REMOVED***
		// Try octal first.
		s := strconv.FormatInt(x, 8)
		if len(s) < len(b) ***REMOVED***
			f.formatOctal(b, x)
			return
		***REMOVED***

		// If it is too long for octal, and PAX is preferred, use a PAX header.
		if paxKeyword != paxNone && tw.preferPax ***REMOVED***
			f.formatOctal(b, 0)
			s := strconv.FormatInt(x, 10)
			paxHeaders[paxKeyword] = s
			return
		***REMOVED***

		tw.usedBinary = true
		f.formatNumeric(b, x)
	***REMOVED***

	// keep a reference to the filename to allow to overwrite it later if we detect that we can use ustar longnames instead of pax
	pathHeaderBytes := s.next(fileNameSize)

	formatString(pathHeaderBytes, hdr.Name, paxPath)

	// Handle out of range ModTime carefully.
	var modTime int64
	if !hdr.ModTime.Before(minTime) && !hdr.ModTime.After(maxTime) ***REMOVED***
		modTime = hdr.ModTime.Unix()
	***REMOVED***

	f.formatOctal(s.next(8), hdr.Mode)               // 100:108
	formatNumeric(s.next(8), int64(hdr.Uid), paxUid) // 108:116
	formatNumeric(s.next(8), int64(hdr.Gid), paxGid) // 116:124
	formatNumeric(s.next(12), hdr.Size, paxSize)     // 124:136
	formatNumeric(s.next(12), modTime, paxNone)      // 136:148 --- consider using pax for finer granularity
	s.next(8)                                        // chksum (148:156)
	s.next(1)[0] = hdr.Typeflag                      // 156:157

	formatString(s.next(100), hdr.Linkname, paxLinkpath)

	copy(s.next(8), []byte("ustar\x0000"))          // 257:265
	formatString(s.next(32), hdr.Uname, paxUname)   // 265:297
	formatString(s.next(32), hdr.Gname, paxGname)   // 297:329
	formatNumeric(s.next(8), hdr.Devmajor, paxNone) // 329:337
	formatNumeric(s.next(8), hdr.Devminor, paxNone) // 337:345

	// keep a reference to the prefix to allow to overwrite it later if we detect that we can use ustar longnames instead of pax
	prefixHeaderBytes := s.next(155)
	formatString(prefixHeaderBytes, "", paxNone) // 345:500  prefix

	// Use the GNU magic instead of POSIX magic if we used any GNU extensions.
	if tw.usedBinary ***REMOVED***
		copy(header[257:265], []byte("ustar  \x00"))
	***REMOVED***

	_, paxPathUsed := paxHeaders[paxPath]
	// try to use a ustar header when only the name is too long
	if !tw.preferPax && len(paxHeaders) == 1 && paxPathUsed ***REMOVED***
		prefix, suffix, ok := splitUSTARPath(hdr.Name)
		if ok ***REMOVED***
			// Since we can encode in USTAR format, disable PAX header.
			delete(paxHeaders, paxPath)

			// Update the path fields
			formatString(pathHeaderBytes, suffix, paxNone)
			formatString(prefixHeaderBytes, prefix, paxNone)
		***REMOVED***
	***REMOVED***

	// The chksum field is terminated by a NUL and a space.
	// This is different from the other octal fields.
	chksum, _ := checksum(header)
	f.formatOctal(header[148:155], chksum) // Never fails
	header[155] = ' '

	// Check if there were any formatting errors.
	if f.err != nil ***REMOVED***
		tw.err = f.err
		return tw.err
	***REMOVED***

	if allowPax ***REMOVED***
		for k, v := range hdr.Xattrs ***REMOVED***
			paxHeaders[paxXattr+k] = v
		***REMOVED***
	***REMOVED***

	if len(paxHeaders) > 0 ***REMOVED***
		if !allowPax ***REMOVED***
			return errInvalidHeader
		***REMOVED***
		if err := tw.writePAXHeader(hdr, paxHeaders); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	tw.nb = int64(hdr.Size)
	tw.pad = (blockSize - (tw.nb % blockSize)) % blockSize

	_, tw.err = tw.w.Write(header)
	return tw.err
***REMOVED***

// splitUSTARPath splits a path according to USTAR prefix and suffix rules.
// If the path is not splittable, then it will return ("", "", false).
func splitUSTARPath(name string) (prefix, suffix string, ok bool) ***REMOVED***
	length := len(name)
	if length <= fileNameSize || !isASCII(name) ***REMOVED***
		return "", "", false
	***REMOVED*** else if length > fileNamePrefixSize+1 ***REMOVED***
		length = fileNamePrefixSize + 1
	***REMOVED*** else if name[length-1] == '/' ***REMOVED***
		length--
	***REMOVED***

	i := strings.LastIndex(name[:length], "/")
	nlen := len(name) - i - 1 // nlen is length of suffix
	plen := i                 // plen is length of prefix
	if i <= 0 || nlen > fileNameSize || nlen == 0 || plen > fileNamePrefixSize ***REMOVED***
		return "", "", false
	***REMOVED***
	return name[:i], name[i+1:], true
***REMOVED***

// writePaxHeader writes an extended pax header to the
// archive.
func (tw *Writer) writePAXHeader(hdr *Header, paxHeaders map[string]string) error ***REMOVED***
	// Prepare extended header
	ext := new(Header)
	ext.Typeflag = TypeXHeader
	// Setting ModTime is required for reader parsing to
	// succeed, and seems harmless enough.
	ext.ModTime = hdr.ModTime
	// The spec asks that we namespace our pseudo files
	// with the current pid.  However, this results in differing outputs
	// for identical inputs.  As such, the constant 0 is now used instead.
	// golang.org/issue/12358
	dir, file := path.Split(hdr.Name)
	fullName := path.Join(dir, "PaxHeaders.0", file)

	ascii := toASCII(fullName)
	if len(ascii) > 100 ***REMOVED***
		ascii = ascii[:100]
	***REMOVED***
	ext.Name = ascii
	// Construct the body
	var buf bytes.Buffer

	// Keys are sorted before writing to body to allow deterministic output.
	var keys []string
	for k := range paxHeaders ***REMOVED***
		keys = append(keys, k)
	***REMOVED***
	sort.Strings(keys)

	for _, k := range keys ***REMOVED***
		fmt.Fprint(&buf, formatPAXRecord(k, paxHeaders[k]))
	***REMOVED***

	ext.Size = int64(len(buf.Bytes()))
	if err := tw.writeHeader(ext, false); err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err := tw.Write(buf.Bytes()); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := tw.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// formatPAXRecord formats a single PAX record, prefixing it with the
// appropriate length.
func formatPAXRecord(k, v string) string ***REMOVED***
	const padding = 3 // Extra padding for ' ', '=', and '\n'
	size := len(k) + len(v) + padding
	size += len(strconv.Itoa(size))
	record := fmt.Sprintf("%d %s=%s\n", size, k, v)

	// Final adjustment if adding size field increased the record size.
	if len(record) != size ***REMOVED***
		size = len(record)
		record = fmt.Sprintf("%d %s=%s\n", size, k, v)
	***REMOVED***
	return record
***REMOVED***

// Write writes to the current entry in the tar archive.
// Write returns the error ErrWriteTooLong if more than
// hdr.Size bytes are written after WriteHeader.
func (tw *Writer) Write(b []byte) (n int, err error) ***REMOVED***
	if tw.closed ***REMOVED***
		err = ErrWriteAfterClose
		return
	***REMOVED***
	overwrite := false
	if int64(len(b)) > tw.nb ***REMOVED***
		b = b[0:tw.nb]
		overwrite = true
	***REMOVED***
	n, err = tw.w.Write(b)
	tw.nb -= int64(n)
	if err == nil && overwrite ***REMOVED***
		err = ErrWriteTooLong
		return
	***REMOVED***
	tw.err = err
	return
***REMOVED***

// Close closes the tar archive, flushing any unwritten
// data to the underlying writer.
func (tw *Writer) Close() error ***REMOVED***
	if tw.err != nil || tw.closed ***REMOVED***
		return tw.err
	***REMOVED***
	tw.Flush()
	tw.closed = true
	if tw.err != nil ***REMOVED***
		return tw.err
	***REMOVED***

	// trailer: two zero blocks
	for i := 0; i < 2; i++ ***REMOVED***
		_, tw.err = tw.w.Write(zeroBlock)
		if tw.err != nil ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return tw.err
***REMOVED***
