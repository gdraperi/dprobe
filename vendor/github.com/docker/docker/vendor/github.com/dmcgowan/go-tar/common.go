// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tar implements access to tar archives.
//
// Tape archives (tar) are a file format for storing a sequence of files that
// can be read and written in a streaming manner.
// This package aims to cover most variations of the format,
// including those produced by GNU and BSD tar tools.
package tar

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// BUG: Use of the Uid and Gid fields in Header could overflow on 32-bit
// architectures. If a large value is encountered when decoding, the result
// stored in Header will be the truncated version.

var (
	ErrHeader          = errors.New("tar: invalid tar header")
	ErrWriteTooLong    = errors.New("tar: write too long")
	ErrFieldTooLong    = errors.New("tar: header field too long")
	ErrWriteAfterClose = errors.New("tar: write after close")
	errMissData        = errors.New("tar: sparse file references non-existent data")
	errUnrefData       = errors.New("tar: sparse file contains unreferenced data")
	errWriteHole       = errors.New("tar: write non-NUL byte in sparse hole")
)

type headerError []string

func (he headerError) Error() string ***REMOVED***
	const prefix = "tar: cannot encode header"
	var ss []string
	for _, s := range he ***REMOVED***
		if s != "" ***REMOVED***
			ss = append(ss, s)
		***REMOVED***
	***REMOVED***
	if len(ss) == 0 ***REMOVED***
		return prefix
	***REMOVED***
	return fmt.Sprintf("%s: %v", prefix, strings.Join(ss, "; and "))
***REMOVED***

// Type flags for Header.Typeflag.
const (
	// Type '0' indicates a regular file.
	TypeReg  = '0'
	TypeRegA = '\x00' // For legacy support; use TypeReg instead

	// Type '1' to '6' are header-only flags and may not have a data body.
	TypeLink    = '1' // Hard link
	TypeSymlink = '2' // Symbolic link
	TypeChar    = '3' // Character device node
	TypeBlock   = '4' // Block device node
	TypeDir     = '5' // Directory
	TypeFifo    = '6' // FIFO node

	// Type '7' is reserved.
	TypeCont = '7'

	// Type 'x' is used by the PAX format to store key-value records that
	// are only relevant to the next file.
	// This package transparently handles these types.
	TypeXHeader = 'x'

	// Type 'g' is used by the PAX format to store key-value records that
	// are relevant to all subsequent files.
	// This package only supports parsing and composing such headers,
	// but does not currently support persisting the global state across files.
	TypeXGlobalHeader = 'g'

	// Type 'S' indicates a sparse file in the GNU format.
	// Header.SparseHoles should be populated when using this type.
	TypeGNUSparse = 'S'

	// Types 'L' and 'K' are used by the GNU format for a meta file
	// used to store the path or link name for the next file.
	// This package transparently handles these types.
	TypeGNULongName = 'L'
	TypeGNULongLink = 'K'
)

// Keywords for PAX extended header records.
const (
	paxNone     = "" // Indicates that no PAX key is suitable
	paxPath     = "path"
	paxLinkpath = "linkpath"
	paxSize     = "size"
	paxUid      = "uid"
	paxGid      = "gid"
	paxUname    = "uname"
	paxGname    = "gname"
	paxMtime    = "mtime"
	paxAtime    = "atime"
	paxCtime    = "ctime"   // Removed from later revision of PAX spec, but was valid
	paxCharset  = "charset" // Currently unused
	paxComment  = "comment" // Currently unused

	paxSchilyXattr = "SCHILY.xattr."

	// Keywords for GNU sparse files in a PAX extended header.
	paxGNUSparse          = "GNU.sparse."
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

// basicKeys is a set of the PAX keys for which we have built-in support.
// This does not contain "charset" or "comment", which are both PAX-specific,
// so adding them as first-class features of Header is unlikely.
// Users can use the PAXRecords field to set it themselves.
var basicKeys = map[string]bool***REMOVED***
	paxPath: true, paxLinkpath: true, paxSize: true, paxUid: true, paxGid: true,
	paxUname: true, paxGname: true, paxMtime: true, paxAtime: true, paxCtime: true,
***REMOVED***

// A Header represents a single header in a tar archive.
// Some fields may not be populated.
//
// For forward compatibility, users that retrieve a Header from Reader.Next,
// mutate it in some ways, and then pass it back to Writer.WriteHeader
// should do so by creating a new Header and copying the fields
// that they are interested in preserving.
type Header struct ***REMOVED***
	Typeflag byte // Type of header entry (should be TypeReg for most files)

	Name     string // Name of file entry
	Linkname string // Target name of link (valid for TypeLink or TypeSymlink)

	Size  int64  // Logical file size in bytes
	Mode  int64  // Permission and mode bits
	Uid   int    // User ID of owner
	Gid   int    // Group ID of owner
	Uname string // User name of owner
	Gname string // Group name of owner

	// If the Format is unspecified, then Writer.WriteHeader rounds ModTime
	// to the nearest second and ignores the AccessTime and ChangeTime fields.
	//
	// To use AccessTime or ChangeTime, specify the Format as PAX or GNU.
	// To use sub-second resolution, specify the Format as PAX.
	ModTime    time.Time // Modification time
	AccessTime time.Time // Access time (requires either PAX or GNU support)
	ChangeTime time.Time // Change time (requires either PAX or GNU support)

	Devmajor int64 // Major device number (valid for TypeChar or TypeBlock)
	Devminor int64 // Minor device number (valid for TypeChar or TypeBlock)

	// SparseHoles represents a sequence of holes in a sparse file.
	//
	// A file is sparse if len(SparseHoles) > 0 or Typeflag is TypeGNUSparse.
	// If TypeGNUSparse is set, then the format is GNU, otherwise
	// the format is PAX (by using GNU-specific PAX records).
	//
	// A sparse file consists of fragments of data, intermixed with holes
	// (described by this field). A hole is semantically a block of NUL-bytes,
	// but does not actually exist within the tar file.
	// The holes must be sorted in ascending order,
	// not overlap with each other, and not extend past the specified Size.
	SparseHoles []SparseEntry

	// Xattrs stores extended attributes as PAX records under the
	// "SCHILY.xattr." namespace.
	//
	// The following are semantically equivalent:
	//  h.Xattrs[key] = value
	//  h.PAXRecords["SCHILY.xattr."+key] = value
	//
	// When Writer.WriteHeader is called, the contents of Xattrs will take
	// precedence over those in PAXRecords.
	//
	// Deprecated: Use PAXRecords instead.
	Xattrs map[string]string

	// PAXRecords is a map of PAX extended header records.
	//
	// User-defined records should have keys of the following form:
	//	VENDOR.keyword
	// Where VENDOR is some namespace in all uppercase, and keyword may
	// not contain the '=' character (e.g., "GOLANG.pkg.version").
	// The key and value should be non-empty UTF-8 strings.
	//
	// When Writer.WriteHeader is called, PAX records derived from the
	// the other fields in Header take precedence over PAXRecords.
	PAXRecords map[string]string

	// Format specifies the format of the tar header.
	//
	// This is set by Reader.Next as a best-effort guess at the format.
	// Since the Reader liberally reads some non-compliant files,
	// it is possible for this to be FormatUnknown.
	//
	// If the format is unspecified when Writer.WriteHeader is called,
	// then it uses the first format (in the order of USTAR, PAX, GNU)
	// capable of encoding this Header (see Format).
	Format Format
***REMOVED***

// SparseEntry represents a Length-sized fragment at Offset in the file.
type SparseEntry struct***REMOVED*** Offset, Length int64 ***REMOVED***

func (s SparseEntry) endOffset() int64 ***REMOVED*** return s.Offset + s.Length ***REMOVED***

// A sparse file can be represented as either a sparseDatas or a sparseHoles.
// As long as the total size is known, they are equivalent and one can be
// converted to the other form and back. The various tar formats with sparse
// file support represent sparse files in the sparseDatas form. That is, they
// specify the fragments in the file that has data, and treat everything else as
// having zero bytes. As such, the encoding and decoding logic in this package
// deals with sparseDatas.
//
// However, the external API uses sparseHoles instead of sparseDatas because the
// zero value of sparseHoles logically represents a normal file (i.e., there are
// no holes in it). On the other hand, the zero value of sparseDatas implies
// that the file has no data in it, which is rather odd.
//
// As an example, if the underlying raw file contains the 10-byte data:
//	var compactFile = "abcdefgh"
//
// And the sparse map has the following entries:
//	var spd sparseDatas = []sparseEntry***REMOVED***
//		***REMOVED***Offset: 2,  Length: 5***REMOVED***,  // Data fragment for 2..6
//		***REMOVED***Offset: 18, Length: 3***REMOVED***,  // Data fragment for 18..20
//	***REMOVED***
//	var sph sparseHoles = []SparseEntry***REMOVED***
//		***REMOVED***Offset: 0,  Length: 2***REMOVED***,  // Hole fragment for 0..1
//		***REMOVED***Offset: 7,  Length: 11***REMOVED***, // Hole fragment for 7..17
//		***REMOVED***Offset: 21, Length: 4***REMOVED***,  // Hole fragment for 21..24
//	***REMOVED***
//
// Then the content of the resulting sparse file with a Header.Size of 25 is:
//	var sparseFile = "\x00"*2 + "abcde" + "\x00"*11 + "fgh" + "\x00"*4
type (
	sparseDatas []SparseEntry
	sparseHoles []SparseEntry
)

// validateSparseEntries reports whether sp is a valid sparse map.
// It does not matter whether sp represents data fragments or hole fragments.
func validateSparseEntries(sp []SparseEntry, size int64) bool ***REMOVED***
	// Validate all sparse entries. These are the same checks as performed by
	// the BSD tar utility.
	if size < 0 ***REMOVED***
		return false
	***REMOVED***
	var pre SparseEntry
	for _, cur := range sp ***REMOVED***
		switch ***REMOVED***
		case cur.Offset < 0 || cur.Length < 0:
			return false // Negative values are never okay
		case cur.Offset > math.MaxInt64-cur.Length:
			return false // Integer overflow with large length
		case cur.endOffset() > size:
			return false // Region extends beyond the actual size
		case pre.endOffset() > cur.Offset:
			return false // Regions cannot overlap and must be in order
		***REMOVED***
		pre = cur
	***REMOVED***
	return true
***REMOVED***

// alignSparseEntries mutates src and returns dst where each fragment's
// starting offset is aligned up to the nearest block edge, and each
// ending offset is aligned down to the nearest block edge.
//
// Even though the Go tar Reader and the BSD tar utility can handle entries
// with arbitrary offsets and lengths, the GNU tar utility can only handle
// offsets and lengths that are multiples of blockSize.
func alignSparseEntries(src []SparseEntry, size int64) []SparseEntry ***REMOVED***
	dst := src[:0]
	for _, s := range src ***REMOVED***
		pos, end := s.Offset, s.endOffset()
		pos += blockPadding(+pos) // Round-up to nearest blockSize
		if end != size ***REMOVED***
			end -= blockPadding(-end) // Round-down to nearest blockSize
		***REMOVED***
		if pos < end ***REMOVED***
			dst = append(dst, SparseEntry***REMOVED***Offset: pos, Length: end - pos***REMOVED***)
		***REMOVED***
	***REMOVED***
	return dst
***REMOVED***

// invertSparseEntries converts a sparse map from one form to the other.
// If the input is sparseHoles, then it will output sparseDatas and vice-versa.
// The input must have been already validated.
//
// This function mutates src and returns a normalized map where:
//	* adjacent fragments are coalesced together
//	* only the last fragment may be empty
//	* the endOffset of the last fragment is the total size
func invertSparseEntries(src []SparseEntry, size int64) []SparseEntry ***REMOVED***
	dst := src[:0]
	var pre SparseEntry
	for _, cur := range src ***REMOVED***
		if cur.Length == 0 ***REMOVED***
			continue // Skip empty fragments
		***REMOVED***
		pre.Length = cur.Offset - pre.Offset
		if pre.Length > 0 ***REMOVED***
			dst = append(dst, pre) // Only add non-empty fragments
		***REMOVED***
		pre.Offset = cur.endOffset()
	***REMOVED***
	pre.Length = size - pre.Offset // Possibly the only empty fragment
	return append(dst, pre)
***REMOVED***

// fileState tracks the number of logical (includes sparse holes) and physical
// (actual in tar archive) bytes remaining for the current file.
//
// Invariant: LogicalRemaining >= PhysicalRemaining
type fileState interface ***REMOVED***
	LogicalRemaining() int64
	PhysicalRemaining() int64
***REMOVED***

// allowedFormats determines which formats can be used.
// The value returned is the logical OR of multiple possible formats.
// If the value is FormatUnknown, then the input Header cannot be encoded
// and an error is returned explaining why.
//
// As a by-product of checking the fields, this function returns paxHdrs, which
// contain all fields that could not be directly encoded.
// A value receiver ensures that this method does not mutate the source Header.
func (h Header) allowedFormats() (format Format, paxHdrs map[string]string, err error) ***REMOVED***
	format = FormatUSTAR | FormatPAX | FormatGNU
	paxHdrs = make(map[string]string)

	var whyNoUSTAR, whyNoPAX, whyNoGNU string
	var preferPAX bool // Prefer PAX over USTAR
	verifyString := func(s string, size int, name, paxKey string) ***REMOVED***
		// NUL-terminator is optional for path and linkpath.
		// Technically, it is required for uname and gname,
		// but neither GNU nor BSD tar checks for it.
		tooLong := len(s) > size
		allowLongGNU := paxKey == paxPath || paxKey == paxLinkpath
		if hasNUL(s) || (tooLong && !allowLongGNU) ***REMOVED***
			whyNoGNU = fmt.Sprintf("GNU cannot encode %s=%q", name, s)
			format.mustNotBe(FormatGNU)
		***REMOVED***
		if !isASCII(s) || tooLong ***REMOVED***
			canSplitUSTAR := paxKey == paxPath
			if _, _, ok := splitUSTARPath(s); !canSplitUSTAR || !ok ***REMOVED***
				whyNoUSTAR = fmt.Sprintf("USTAR cannot encode %s=%q", name, s)
				format.mustNotBe(FormatUSTAR)
			***REMOVED***
			if paxKey == paxNone ***REMOVED***
				whyNoPAX = fmt.Sprintf("PAX cannot encode %s=%q", name, s)
				format.mustNotBe(FormatPAX)
			***REMOVED*** else ***REMOVED***
				paxHdrs[paxKey] = s
			***REMOVED***
		***REMOVED***
		if v, ok := h.PAXRecords[paxKey]; ok && v == s ***REMOVED***
			paxHdrs[paxKey] = v
		***REMOVED***
	***REMOVED***
	verifyNumeric := func(n int64, size int, name, paxKey string) ***REMOVED***
		if !fitsInBase256(size, n) ***REMOVED***
			whyNoGNU = fmt.Sprintf("GNU cannot encode %s=%d", name, n)
			format.mustNotBe(FormatGNU)
		***REMOVED***
		if !fitsInOctal(size, n) ***REMOVED***
			whyNoUSTAR = fmt.Sprintf("USTAR cannot encode %s=%d", name, n)
			format.mustNotBe(FormatUSTAR)
			if paxKey == paxNone ***REMOVED***
				whyNoPAX = fmt.Sprintf("PAX cannot encode %s=%d", name, n)
				format.mustNotBe(FormatPAX)
			***REMOVED*** else ***REMOVED***
				paxHdrs[paxKey] = strconv.FormatInt(n, 10)
			***REMOVED***
		***REMOVED***
		if v, ok := h.PAXRecords[paxKey]; ok && v == strconv.FormatInt(n, 10) ***REMOVED***
			paxHdrs[paxKey] = v
		***REMOVED***
	***REMOVED***
	verifyTime := func(ts time.Time, size int, name, paxKey string) ***REMOVED***
		if ts.IsZero() ***REMOVED***
			return // Always okay
		***REMOVED***
		if !fitsInBase256(size, ts.Unix()) ***REMOVED***
			whyNoGNU = fmt.Sprintf("GNU cannot encode %s=%v", name, ts)
			format.mustNotBe(FormatGNU)
		***REMOVED***
		isMtime := paxKey == paxMtime
		fitsOctal := fitsInOctal(size, ts.Unix())
		if (isMtime && !fitsOctal) || !isMtime ***REMOVED***
			whyNoUSTAR = fmt.Sprintf("USTAR cannot encode %s=%v", name, ts)
			format.mustNotBe(FormatUSTAR)
		***REMOVED***
		needsNano := ts.Nanosecond() != 0
		if !isMtime || !fitsOctal || needsNano ***REMOVED***
			preferPAX = true // USTAR may truncate sub-second measurements
			if paxKey == paxNone ***REMOVED***
				whyNoPAX = fmt.Sprintf("PAX cannot encode %s=%v", name, ts)
				format.mustNotBe(FormatPAX)
			***REMOVED*** else ***REMOVED***
				paxHdrs[paxKey] = formatPAXTime(ts)
			***REMOVED***
		***REMOVED***
		if v, ok := h.PAXRecords[paxKey]; ok && v == formatPAXTime(ts) ***REMOVED***
			paxHdrs[paxKey] = v
		***REMOVED***
	***REMOVED***

	// Check basic fields.
	var blk block
	v7 := blk.V7()
	ustar := blk.USTAR()
	gnu := blk.GNU()
	verifyString(h.Name, len(v7.Name()), "Name", paxPath)
	verifyString(h.Linkname, len(v7.LinkName()), "Linkname", paxLinkpath)
	verifyString(h.Uname, len(ustar.UserName()), "Uname", paxUname)
	verifyString(h.Gname, len(ustar.GroupName()), "Gname", paxGname)
	verifyNumeric(h.Mode, len(v7.Mode()), "Mode", paxNone)
	verifyNumeric(int64(h.Uid), len(v7.UID()), "Uid", paxUid)
	verifyNumeric(int64(h.Gid), len(v7.GID()), "Gid", paxGid)
	verifyNumeric(h.Size, len(v7.Size()), "Size", paxSize)
	verifyNumeric(h.Devmajor, len(ustar.DevMajor()), "Devmajor", paxNone)
	verifyNumeric(h.Devminor, len(ustar.DevMinor()), "Devminor", paxNone)
	verifyTime(h.ModTime, len(v7.ModTime()), "ModTime", paxMtime)
	verifyTime(h.AccessTime, len(gnu.AccessTime()), "AccessTime", paxAtime)
	verifyTime(h.ChangeTime, len(gnu.ChangeTime()), "ChangeTime", paxCtime)

	// Check for header-only types.
	var whyOnlyPAX, whyOnlyGNU string
	switch h.Typeflag ***REMOVED***
	case TypeReg, TypeChar, TypeBlock, TypeFifo, TypeGNUSparse:
		// Exclude TypeLink and TypeSymlink, since they may reference directories.
		if strings.HasSuffix(h.Name, "/") ***REMOVED***
			return FormatUnknown, nil, headerError***REMOVED***"filename may not have trailing slash"***REMOVED***
		***REMOVED***
	case TypeXHeader, TypeGNULongName, TypeGNULongLink:
		return FormatUnknown, nil, headerError***REMOVED***"cannot manually encode TypeXHeader, TypeGNULongName, or TypeGNULongLink headers"***REMOVED***
	case TypeXGlobalHeader:
		if !reflect.DeepEqual(h, Header***REMOVED***Typeflag: h.Typeflag, Xattrs: h.Xattrs, PAXRecords: h.PAXRecords, Format: h.Format***REMOVED***) ***REMOVED***
			return FormatUnknown, nil, headerError***REMOVED***"only PAXRecords may be set for TypeXGlobalHeader"***REMOVED***
		***REMOVED***
		whyOnlyPAX = "only PAX supports TypeXGlobalHeader"
		format.mayOnlyBe(FormatPAX)
	***REMOVED***
	if !isHeaderOnlyType(h.Typeflag) && h.Size < 0 ***REMOVED***
		return FormatUnknown, nil, headerError***REMOVED***"negative size on header-only type"***REMOVED***
	***REMOVED***

	// Check PAX records.
	if len(h.Xattrs) > 0 ***REMOVED***
		for k, v := range h.Xattrs ***REMOVED***
			paxHdrs[paxSchilyXattr+k] = v
		***REMOVED***
		whyOnlyPAX = "only PAX supports Xattrs"
		format.mayOnlyBe(FormatPAX)
	***REMOVED***
	if len(h.PAXRecords) > 0 ***REMOVED***
		for k, v := range h.PAXRecords ***REMOVED***
			switch _, exists := paxHdrs[k]; ***REMOVED***
			case exists:
				continue // Do not overwrite existing records
			case h.Typeflag == TypeXGlobalHeader:
				paxHdrs[k] = v // Copy all records
			case !basicKeys[k] && !strings.HasPrefix(k, paxGNUSparse):
				paxHdrs[k] = v // Ignore local records that may conflict
			***REMOVED***
		***REMOVED***
		whyOnlyPAX = "only PAX supports PAXRecords"
		format.mayOnlyBe(FormatPAX)
	***REMOVED***
	for k, v := range paxHdrs ***REMOVED***
		if !validPAXRecord(k, v) ***REMOVED***
			return FormatUnknown, nil, headerError***REMOVED***fmt.Sprintf("invalid PAX record: %q", k+" = "+v)***REMOVED***
		***REMOVED***
	***REMOVED***

	// Check sparse files.
	if len(h.SparseHoles) > 0 || h.Typeflag == TypeGNUSparse ***REMOVED***
		if isHeaderOnlyType(h.Typeflag) ***REMOVED***
			return FormatUnknown, nil, headerError***REMOVED***"header-only type cannot be sparse"***REMOVED***
		***REMOVED***
		if !validateSparseEntries(h.SparseHoles, h.Size) ***REMOVED***
			return FormatUnknown, nil, headerError***REMOVED***"invalid sparse holes"***REMOVED***
		***REMOVED***
		if h.Typeflag == TypeGNUSparse ***REMOVED***
			whyOnlyGNU = "only GNU supports TypeGNUSparse"
			format.mayOnlyBe(FormatGNU)
		***REMOVED*** else ***REMOVED***
			whyNoGNU = "GNU supports sparse files only with TypeGNUSparse"
			format.mustNotBe(FormatGNU)
		***REMOVED***
		whyNoUSTAR = "USTAR does not support sparse files"
		format.mustNotBe(FormatUSTAR)
	***REMOVED***

	// Check desired format.
	if wantFormat := h.Format; wantFormat != FormatUnknown ***REMOVED***
		if wantFormat.has(FormatPAX) && !preferPAX ***REMOVED***
			wantFormat.mayBe(FormatUSTAR) // PAX implies USTAR allowed too
		***REMOVED***
		format.mayOnlyBe(wantFormat) // Set union of formats allowed and format wanted
	***REMOVED***
	if format == FormatUnknown ***REMOVED***
		switch h.Format ***REMOVED***
		case FormatUSTAR:
			err = headerError***REMOVED***"Format specifies USTAR", whyNoUSTAR, whyOnlyPAX, whyOnlyGNU***REMOVED***
		case FormatPAX:
			err = headerError***REMOVED***"Format specifies PAX", whyNoPAX, whyOnlyGNU***REMOVED***
		case FormatGNU:
			err = headerError***REMOVED***"Format specifies GNU", whyNoGNU, whyOnlyPAX***REMOVED***
		default:
			err = headerError***REMOVED***whyNoUSTAR, whyNoPAX, whyNoGNU, whyOnlyPAX, whyOnlyGNU***REMOVED***
		***REMOVED***
	***REMOVED***
	return format, paxHdrs, err
***REMOVED***

var sysSparseDetect func(f *os.File) (sparseHoles, error)
var sysSparsePunch func(f *os.File, sph sparseHoles) error

// DetectSparseHoles searches for holes within f to populate SparseHoles
// on supported operating systems and filesystems.
// The file offset is cleared to zero.
//
// When packing a sparse file, DetectSparseHoles should be called prior to
// serializing the header to the archive with Writer.WriteHeader.
func (h *Header) DetectSparseHoles(f *os.File) (err error) ***REMOVED***
	defer func() ***REMOVED***
		if _, serr := f.Seek(0, io.SeekStart); err == nil ***REMOVED***
			err = serr
		***REMOVED***
	***REMOVED***()

	h.SparseHoles = nil
	if sysSparseDetect != nil ***REMOVED***
		sph, err := sysSparseDetect(f)
		h.SparseHoles = sph
		return err
	***REMOVED***
	return nil
***REMOVED***

// PunchSparseHoles destroys the contents of f, and prepares a sparse file
// (on supported operating systems and filesystems)
// with holes punched according to SparseHoles.
// The file offset is cleared to zero.
//
// When extracting a sparse file, PunchSparseHoles should be called prior to
// populating the content of a file with Reader.WriteTo.
func (h *Header) PunchSparseHoles(f *os.File) (err error) ***REMOVED***
	defer func() ***REMOVED***
		if _, serr := f.Seek(0, io.SeekStart); err == nil ***REMOVED***
			err = serr
		***REMOVED***
	***REMOVED***()

	if err := f.Truncate(0); err != nil ***REMOVED***
		return err
	***REMOVED***

	var size int64
	if len(h.SparseHoles) > 0 ***REMOVED***
		size = h.SparseHoles[len(h.SparseHoles)-1].endOffset()
	***REMOVED***
	if !validateSparseEntries(h.SparseHoles, size) ***REMOVED***
		return errors.New("tar: invalid sparse holes")
	***REMOVED***

	if size == 0 ***REMOVED***
		return nil // For non-sparse files, do nothing (other than Truncate)
	***REMOVED***
	if sysSparsePunch != nil ***REMOVED***
		return sysSparsePunch(f, h.SparseHoles)
	***REMOVED***
	return f.Truncate(size)
***REMOVED***

// FileInfo returns an os.FileInfo for the Header.
func (h *Header) FileInfo() os.FileInfo ***REMOVED***
	return headerFileInfo***REMOVED***h***REMOVED***
***REMOVED***

// headerFileInfo implements os.FileInfo.
type headerFileInfo struct ***REMOVED***
	h *Header
***REMOVED***

func (fi headerFileInfo) Size() int64        ***REMOVED*** return fi.h.Size ***REMOVED***
func (fi headerFileInfo) IsDir() bool        ***REMOVED*** return fi.Mode().IsDir() ***REMOVED***
func (fi headerFileInfo) ModTime() time.Time ***REMOVED*** return fi.h.ModTime ***REMOVED***
func (fi headerFileInfo) Sys() interface***REMOVED******REMOVED***   ***REMOVED*** return fi.h ***REMOVED***

// Name returns the base name of the file.
func (fi headerFileInfo) Name() string ***REMOVED***
	if fi.IsDir() ***REMOVED***
		return path.Base(path.Clean(fi.h.Name))
	***REMOVED***
	return path.Base(fi.h.Name)
***REMOVED***

// Mode returns the permission and mode bits for the headerFileInfo.
func (fi headerFileInfo) Mode() (mode os.FileMode) ***REMOVED***
	// Set file permission bits.
	mode = os.FileMode(fi.h.Mode).Perm()

	// Set setuid, setgid and sticky bits.
	if fi.h.Mode&c_ISUID != 0 ***REMOVED***
		mode |= os.ModeSetuid
	***REMOVED***
	if fi.h.Mode&c_ISGID != 0 ***REMOVED***
		mode |= os.ModeSetgid
	***REMOVED***
	if fi.h.Mode&c_ISVTX != 0 ***REMOVED***
		mode |= os.ModeSticky
	***REMOVED***

	// Set file mode bits; clear perm, setuid, setgid, and sticky bits.
	switch m := os.FileMode(fi.h.Mode) &^ 07777; m ***REMOVED***
	case c_ISDIR:
		mode |= os.ModeDir
	case c_ISFIFO:
		mode |= os.ModeNamedPipe
	case c_ISLNK:
		mode |= os.ModeSymlink
	case c_ISBLK:
		mode |= os.ModeDevice
	case c_ISCHR:
		mode |= os.ModeDevice
		mode |= os.ModeCharDevice
	case c_ISSOCK:
		mode |= os.ModeSocket
	***REMOVED***

	switch fi.h.Typeflag ***REMOVED***
	case TypeSymlink:
		mode |= os.ModeSymlink
	case TypeChar:
		mode |= os.ModeDevice
		mode |= os.ModeCharDevice
	case TypeBlock:
		mode |= os.ModeDevice
	case TypeDir:
		mode |= os.ModeDir
	case TypeFifo:
		mode |= os.ModeNamedPipe
	***REMOVED***

	return mode
***REMOVED***

// sysStat, if non-nil, populates h from system-dependent fields of fi.
var sysStat func(fi os.FileInfo, h *Header) error

const (
	// Mode constants from the USTAR spec:
	// See http://pubs.opengroup.org/onlinepubs/9699919799/utilities/pax.html#tag_20_92_13_06
	c_ISUID = 04000 // Set uid
	c_ISGID = 02000 // Set gid
	c_ISVTX = 01000 // Save text (sticky bit)

	// Common Unix mode constants; these are not defined in any common tar standard.
	// Header.FileInfo understands these, but FileInfoHeader will never produce these.
	c_ISDIR  = 040000  // Directory
	c_ISFIFO = 010000  // FIFO
	c_ISREG  = 0100000 // Regular file
	c_ISLNK  = 0120000 // Symbolic link
	c_ISBLK  = 060000  // Block special file
	c_ISCHR  = 020000  // Character special file
	c_ISSOCK = 0140000 // Socket
)

// FileInfoHeader creates a partially-populated Header from fi.
// If fi describes a symlink, FileInfoHeader records link as the link target.
// If fi describes a directory, a slash is appended to the name.
//
// Since os.FileInfo's Name method only returns the base name of
// the file it describes, it may be necessary to modify Header.Name
// to provide the full path name of the file.
//
// This function does not populate Header.SparseHoles;
// for sparse file support, additionally call Header.DetectSparseHoles.
func FileInfoHeader(fi os.FileInfo, link string) (*Header, error) ***REMOVED***
	if fi == nil ***REMOVED***
		return nil, errors.New("tar: FileInfo is nil")
	***REMOVED***
	fm := fi.Mode()
	h := &Header***REMOVED***
		Name:    fi.Name(),
		ModTime: fi.ModTime(),
		Mode:    int64(fm.Perm()), // or'd with c_IS* constants later
	***REMOVED***
	switch ***REMOVED***
	case fm.IsRegular():
		h.Typeflag = TypeReg
		h.Size = fi.Size()
	case fi.IsDir():
		h.Typeflag = TypeDir
		h.Name += "/"
	case fm&os.ModeSymlink != 0:
		h.Typeflag = TypeSymlink
		h.Linkname = link
	case fm&os.ModeDevice != 0:
		if fm&os.ModeCharDevice != 0 ***REMOVED***
			h.Typeflag = TypeChar
		***REMOVED*** else ***REMOVED***
			h.Typeflag = TypeBlock
		***REMOVED***
	case fm&os.ModeNamedPipe != 0:
		h.Typeflag = TypeFifo
	case fm&os.ModeSocket != 0:
		return nil, fmt.Errorf("tar: sockets not supported")
	default:
		return nil, fmt.Errorf("tar: unknown file mode %v", fm)
	***REMOVED***
	if fm&os.ModeSetuid != 0 ***REMOVED***
		h.Mode |= c_ISUID
	***REMOVED***
	if fm&os.ModeSetgid != 0 ***REMOVED***
		h.Mode |= c_ISGID
	***REMOVED***
	if fm&os.ModeSticky != 0 ***REMOVED***
		h.Mode |= c_ISVTX
	***REMOVED***
	// If possible, populate additional fields from OS-specific
	// FileInfo fields.
	if sys, ok := fi.Sys().(*Header); ok ***REMOVED***
		// This FileInfo came from a Header (not the OS). Use the
		// original Header to populate all remaining fields.
		h.Uid = sys.Uid
		h.Gid = sys.Gid
		h.Uname = sys.Uname
		h.Gname = sys.Gname
		h.AccessTime = sys.AccessTime
		h.ChangeTime = sys.ChangeTime
		if sys.Xattrs != nil ***REMOVED***
			h.Xattrs = make(map[string]string)
			for k, v := range sys.Xattrs ***REMOVED***
				h.Xattrs[k] = v
			***REMOVED***
		***REMOVED***
		if sys.Typeflag == TypeLink ***REMOVED***
			// hard link
			h.Typeflag = TypeLink
			h.Size = 0
			h.Linkname = sys.Linkname
		***REMOVED***
		if sys.SparseHoles != nil ***REMOVED***
			h.SparseHoles = append([]SparseEntry***REMOVED******REMOVED***, sys.SparseHoles...)
		***REMOVED***
		if sys.PAXRecords != nil ***REMOVED***
			h.PAXRecords = make(map[string]string)
			for k, v := range sys.PAXRecords ***REMOVED***
				h.PAXRecords[k] = v
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if sysStat != nil ***REMOVED***
		return h, sysStat(fi, h)
	***REMOVED***
	return h, nil
***REMOVED***

// isHeaderOnlyType checks if the given type flag is of the type that has no
// data section even if a size is specified.
func isHeaderOnlyType(flag byte) bool ***REMOVED***
	switch flag ***REMOVED***
	case TypeLink, TypeSymlink, TypeChar, TypeBlock, TypeDir, TypeFifo:
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

func min(a, b int64) int64 ***REMOVED***
	if a < b ***REMOVED***
		return a
	***REMOVED***
	return b
***REMOVED***
