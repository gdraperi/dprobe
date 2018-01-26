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
	usedBinary bool  // whether the binary numeric field extension was used
	preferPax  bool  // use PAX header instead of binary numeric header
	hdrBuff    block // buffer to use in writeHeader when writing a regular header
	paxHdrBuff block // buffer to use in writeHeader when writing a PAX header
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

	// TODO(dsnet): we might want to use PAX headers for
	// subsecond time resolution, but for now let's just capture
	// too long fields or non ascii characters

	// We need to select which scratch buffer to use carefully,
	// since this method is called recursively to write PAX headers.
	// If allowPax is true, this is the non-recursive call, and we will use hdrBuff.
	// If allowPax is false, we are being called by writePAXHeader, and hdrBuff is
	// already being used by the non-recursive call, so we must use paxHdrBuff.
	header := &tw.hdrBuff
	if !allowPax ***REMOVED***
		header = &tw.paxHdrBuff
	***REMOVED***
	copy(header[:], zeroBlock[:])

	// Wrappers around formatter that automatically sets paxHeaders if the
	// argument extends beyond the capacity of the input byte slice.
	var f formatter
	var formatString = func(b []byte, s string, paxKeyword string) ***REMOVED***
		needsPaxHeader := paxKeyword != paxNone && len(s) > len(b) || !isASCII(s)
		if needsPaxHeader ***REMOVED***
			paxHeaders[paxKeyword] = s
		***REMOVED***

		// Write string in a best-effort manner to satisfy readers that expect
		// the field to be non-empty.
		s = toASCII(s)
		if len(s) > len(b) ***REMOVED***
			s = s[:len(b)]
		***REMOVED***
		f.formatString(b, s) // Should never error
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

	// Handle out of range ModTime carefully.
	var modTime int64
	if !hdr.ModTime.Before(minTime) && !hdr.ModTime.After(maxTime) ***REMOVED***
		modTime = hdr.ModTime.Unix()
	***REMOVED***

	v7 := header.V7()
	formatString(v7.Name(), hdr.Name, paxPath)
	// TODO(dsnet): The GNU format permits the mode field to be encoded in
	// base-256 format. Thus, we can use formatNumeric instead of formatOctal.
	f.formatOctal(v7.Mode(), hdr.Mode)
	formatNumeric(v7.UID(), int64(hdr.Uid), paxUid)
	formatNumeric(v7.GID(), int64(hdr.Gid), paxGid)
	formatNumeric(v7.Size(), hdr.Size, paxSize)
	// TODO(dsnet): Consider using PAX for finer time granularity.
	formatNumeric(v7.ModTime(), modTime, paxNone)
	v7.TypeFlag()[0] = hdr.Typeflag
	formatString(v7.LinkName(), hdr.Linkname, paxLinkpath)

	ustar := header.USTAR()
	formatString(ustar.UserName(), hdr.Uname, paxUname)
	formatString(ustar.GroupName(), hdr.Gname, paxGname)
	formatNumeric(ustar.DevMajor(), hdr.Devmajor, paxNone)
	formatNumeric(ustar.DevMinor(), hdr.Devminor, paxNone)

	// TODO(dsnet): The logic surrounding the prefix field is broken when trying
	// to encode the header as GNU format. The challenge with the current logic
	// is that we are unsure what format we are using at any given moment until
	// we have processed *all* of the fields. The problem is that by the time
	// all fields have been processed, some work has already been done to handle
	// each field under the assumption that it is for one given format or
	// another. In some situations, this causes the Writer to be confused and
	// encode a prefix field when the format being used is GNU. Thus, producing
	// an invalid tar file.
	//
	// As a short-term fix, we disable the logic to use the prefix field, which
	// will force the badly generated GNU files to become encoded as being
	// the PAX format.
	//
	// As an alternative fix, we could hard-code preferPax to be true. However,
	// this is problematic for the following reasons:
	//	* The preferPax functionality is not tested at all.
	//	* This can result in headers that try to use both the GNU and PAX
	//	features at the same time, which is also wrong.
	//
	// The proper fix for this is to use a two-pass method:
	//	* The first pass simply determines what set of formats can possibly
	//	encode the given header.
	//	* The second pass actually encodes the header as that given format
	//	without worrying about violating the format.
	//
	// See the following:
	//	https://golang.org/issue/12594
	//	https://golang.org/issue/17630
	//	https://golang.org/issue/9683
	const usePrefix = false

	// try to use a ustar header when only the name is too long
	_, paxPathUsed := paxHeaders[paxPath]
	if usePrefix && !tw.preferPax && len(paxHeaders) == 1 && paxPathUsed ***REMOVED***
		prefix, suffix, ok := splitUSTARPath(hdr.Name)
		if ok ***REMOVED***
			// Since we can encode in USTAR format, disable PAX header.
			delete(paxHeaders, paxPath)

			// Update the path fields
			formatString(v7.Name(), suffix, paxNone)
			formatString(ustar.Prefix(), prefix, paxNone)
		***REMOVED***
	***REMOVED***

	if tw.usedBinary ***REMOVED***
		header.SetFormat(formatGNU)
	***REMOVED*** else ***REMOVED***
		header.SetFormat(formatUSTAR)
	***REMOVED***

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
	tw.nb = hdr.Size
	tw.pad = (blockSize - (tw.nb % blockSize)) % blockSize

	_, tw.err = tw.w.Write(header[:])
	return tw.err
***REMOVED***

// splitUSTARPath splits a path according to USTAR prefix and suffix rules.
// If the path is not splittable, then it will return ("", "", false).
func splitUSTARPath(name string) (prefix, suffix string, ok bool) ***REMOVED***
	length := len(name)
	if length <= nameSize || !isASCII(name) ***REMOVED***
		return "", "", false
	***REMOVED*** else if length > prefixSize+1 ***REMOVED***
		length = prefixSize + 1
	***REMOVED*** else if name[length-1] == '/' ***REMOVED***
		length--
	***REMOVED***

	i := strings.LastIndex(name[:length], "/")
	nlen := len(name) - i - 1 // nlen is length of suffix
	plen := i                 // plen is length of prefix
	if i <= 0 || nlen > nameSize || nlen == 0 || plen > prefixSize ***REMOVED***
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
	// with the current pid. However, this results in differing outputs
	// for identical inputs. As such, the constant 0 is now used instead.
	// golang.org/issue/12358
	dir, file := path.Split(hdr.Name)
	fullName := path.Join(dir, "PaxHeaders.0", file)

	ascii := toASCII(fullName)
	if len(ascii) > nameSize ***REMOVED***
		ascii = ascii[:nameSize]
	***REMOVED***
	ext.Name = ascii
	// Construct the body
	var buf bytes.Buffer

	// Keys are sorted before writing to body to allow deterministic output.
	keys := make([]string, 0, len(paxHeaders))
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
		_, tw.err = tw.w.Write(zeroBlock[:])
		if tw.err != nil ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return tw.err
***REMOVED***
