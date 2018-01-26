// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tar

import "strings"

// Format represents the tar archive format.
//
// The original tar format was introduced in Unix V7.
// Since then, there have been multiple competing formats attempting to
// standardize or extend the V7 format to overcome its limitations.
// The most common formats are the USTAR, PAX, and GNU formats,
// each with their own advantages and limitations.
//
// The following table captures the capabilities of each format:
//
//	                  |  USTAR |       PAX |       GNU
//	------------------+--------+-----------+----------
//	Name              |   256B | unlimited | unlimited
//	Linkname          |   100B | unlimited | unlimited
//	Size              | uint33 | unlimited |    uint89
//	Mode              | uint21 |    uint21 |    uint57
//	Uid/Gid           | uint21 | unlimited |    uint57
//	Uname/Gname       |    32B | unlimited |       32B
//	ModTime           | uint33 | unlimited |     int89
//	AccessTime        |    n/a | unlimited |     int89
//	ChangeTime        |    n/a | unlimited |     int89
//	Devmajor/Devminor | uint21 |    uint21 |    uint57
//	------------------+--------+-----------+----------
//	string encoding   |  ASCII |     UTF-8 |    binary
//	sub-second times  |     no |       yes |        no
//	sparse files      |     no |       yes |       yes
//
// The table's upper portion shows the Header fields, where each format reports
// the maximum number of bytes allowed for each string field and
// the integer type used to store each numeric field
// (where timestamps are stored as the number of seconds since the Unix epoch).
//
// The table's lower portion shows specialized features of each format,
// such as supported string encodings, support for sub-second timestamps,
// or support for sparse files.
type Format int

// Constants to identify various tar formats.
const (
	// Deliberately hide the meaning of constants from public API.
	_ Format = (1 << iota) / 4 // Sequence of 0, 0, 1, 2, 4, 8, etc...

	// FormatUnknown indicates that the format is unknown.
	FormatUnknown

	// The format of the original Unix V7 tar tool prior to standardization.
	formatV7

	// FormatUSTAR represents the USTAR header format defined in POSIX.1-1988.
	//
	// While this format is compatible with most tar readers,
	// the format has several limitations making it unsuitable for some usages.
	// Most notably, it cannot support sparse files, files larger than 8GiB,
	// filenames larger than 256 characters, and non-ASCII filenames.
	//
	// Reference:
	//	http://pubs.opengroup.org/onlinepubs/9699919799/utilities/pax.html#tag_20_92_13_06
	FormatUSTAR

	// FormatPAX represents the PAX header format defined in POSIX.1-2001.
	//
	// PAX extends USTAR by writing a special file with Typeflag TypeXHeader
	// preceding the original header. This file contains a set of key-value
	// records, which are used to overcome USTAR's shortcomings, in addition to
	// providing the ability to have sub-second resolution for timestamps.
	//
	// Some newer formats add their own extensions to PAX by defining their
	// own keys and assigning certain semantic meaning to the associated values.
	// For example, sparse file support in PAX is implemented using keys
	// defined by the GNU manual (e.g., "GNU.sparse.map").
	//
	// Reference:
	//	http://pubs.opengroup.org/onlinepubs/009695399/utilities/pax.html
	FormatPAX

	// FormatGNU represents the GNU header format.
	//
	// The GNU header format is older than the USTAR and PAX standards and
	// is not compatible with them. The GNU format supports
	// arbitrary file sizes, filenames of arbitrary encoding and length,
	// sparse files, and other features.
	//
	// It is recommended that PAX be chosen over GNU unless the target
	// application can only parse GNU formatted archives.
	//
	// Reference:
	//	http://www.gnu.org/software/tar/manual/html_node/Standard.html
	FormatGNU

	// Schily's tar format, which is incompatible with USTAR.
	// This does not cover STAR extensions to the PAX format; these fall under
	// the PAX format.
	formatSTAR

	formatMax
)

func (f Format) has(f2 Format) bool   ***REMOVED*** return f&f2 != 0 ***REMOVED***
func (f *Format) mayBe(f2 Format)     ***REMOVED*** *f |= f2 ***REMOVED***
func (f *Format) mayOnlyBe(f2 Format) ***REMOVED*** *f &= f2 ***REMOVED***
func (f *Format) mustNotBe(f2 Format) ***REMOVED*** *f &^= f2 ***REMOVED***

var formatNames = map[Format]string***REMOVED***
	formatV7: "V7", FormatUSTAR: "USTAR", FormatPAX: "PAX", FormatGNU: "GNU", formatSTAR: "STAR",
***REMOVED***

func (f Format) String() string ***REMOVED***
	var ss []string
	for f2 := Format(1); f2 < formatMax; f2 <<= 1 ***REMOVED***
		if f.has(f2) ***REMOVED***
			ss = append(ss, formatNames[f2])
		***REMOVED***
	***REMOVED***
	switch len(ss) ***REMOVED***
	case 0:
		return "<unknown>"
	case 1:
		return ss[0]
	default:
		return "(" + strings.Join(ss, " | ") + ")"
	***REMOVED***
***REMOVED***

// Magics used to identify various formats.
const (
	magicGNU, versionGNU     = "ustar ", " \x00"
	magicUSTAR, versionUSTAR = "ustar\x00", "00"
	trailerSTAR              = "tar\x00"
)

// Size constants from various tar specifications.
const (
	blockSize  = 512 // Size of each block in a tar stream
	nameSize   = 100 // Max length of the name field in USTAR format
	prefixSize = 155 // Max length of the prefix field in USTAR format
)

// blockPadding computes the number of bytes needed to pad offset up to the
// nearest block edge where 0 <= n < blockSize.
func blockPadding(offset int64) (n int64) ***REMOVED***
	return -offset & (blockSize - 1)
***REMOVED***

var zeroBlock block

type block [blockSize]byte

// Convert block to any number of formats.
func (b *block) V7() *headerV7       ***REMOVED*** return (*headerV7)(b) ***REMOVED***
func (b *block) GNU() *headerGNU     ***REMOVED*** return (*headerGNU)(b) ***REMOVED***
func (b *block) STAR() *headerSTAR   ***REMOVED*** return (*headerSTAR)(b) ***REMOVED***
func (b *block) USTAR() *headerUSTAR ***REMOVED*** return (*headerUSTAR)(b) ***REMOVED***
func (b *block) Sparse() sparseArray ***REMOVED*** return (sparseArray)(b[:]) ***REMOVED***

// GetFormat checks that the block is a valid tar header based on the checksum.
// It then attempts to guess the specific format based on magic values.
// If the checksum fails, then FormatUnknown is returned.
func (b *block) GetFormat() Format ***REMOVED***
	// Verify checksum.
	var p parser
	value := p.parseOctal(b.V7().Chksum())
	chksum1, chksum2 := b.ComputeChecksum()
	if p.err != nil || (value != chksum1 && value != chksum2) ***REMOVED***
		return FormatUnknown
	***REMOVED***

	// Guess the magic values.
	magic := string(b.USTAR().Magic())
	version := string(b.USTAR().Version())
	trailer := string(b.STAR().Trailer())
	switch ***REMOVED***
	case magic == magicUSTAR && trailer == trailerSTAR:
		return formatSTAR
	case magic == magicUSTAR:
		return FormatUSTAR | FormatPAX
	case magic == magicGNU && version == versionGNU:
		return FormatGNU
	default:
		return formatV7
	***REMOVED***
***REMOVED***

// SetFormat writes the magic values necessary for specified format
// and then updates the checksum accordingly.
func (b *block) SetFormat(format Format) ***REMOVED***
	// Set the magic values.
	switch ***REMOVED***
	case format.has(formatV7):
		// Do nothing.
	case format.has(FormatGNU):
		copy(b.GNU().Magic(), magicGNU)
		copy(b.GNU().Version(), versionGNU)
	case format.has(formatSTAR):
		copy(b.STAR().Magic(), magicUSTAR)
		copy(b.STAR().Version(), versionUSTAR)
		copy(b.STAR().Trailer(), trailerSTAR)
	case format.has(FormatUSTAR | FormatPAX):
		copy(b.USTAR().Magic(), magicUSTAR)
		copy(b.USTAR().Version(), versionUSTAR)
	default:
		panic("invalid format")
	***REMOVED***

	// Update checksum.
	// This field is special in that it is terminated by a NULL then space.
	var f formatter
	field := b.V7().Chksum()
	chksum, _ := b.ComputeChecksum() // Possible values are 256..128776
	f.formatOctal(field[:7], chksum) // Never fails since 128776 < 262143
	field[7] = ' '
***REMOVED***

// ComputeChecksum computes the checksum for the header block.
// POSIX specifies a sum of the unsigned byte values, but the Sun tar used
// signed byte values.
// We compute and return both.
func (b *block) ComputeChecksum() (unsigned, signed int64) ***REMOVED***
	for i, c := range b ***REMOVED***
		if 148 <= i && i < 156 ***REMOVED***
			c = ' ' // Treat the checksum field itself as all spaces.
		***REMOVED***
		unsigned += int64(uint8(c))
		signed += int64(int8(c))
	***REMOVED***
	return unsigned, signed
***REMOVED***

// Reset clears the block with all zeros.
func (b *block) Reset() ***REMOVED***
	*b = block***REMOVED******REMOVED***
***REMOVED***

type headerV7 [blockSize]byte

func (h *headerV7) Name() []byte     ***REMOVED*** return h[000:][:100] ***REMOVED***
func (h *headerV7) Mode() []byte     ***REMOVED*** return h[100:][:8] ***REMOVED***
func (h *headerV7) UID() []byte      ***REMOVED*** return h[108:][:8] ***REMOVED***
func (h *headerV7) GID() []byte      ***REMOVED*** return h[116:][:8] ***REMOVED***
func (h *headerV7) Size() []byte     ***REMOVED*** return h[124:][:12] ***REMOVED***
func (h *headerV7) ModTime() []byte  ***REMOVED*** return h[136:][:12] ***REMOVED***
func (h *headerV7) Chksum() []byte   ***REMOVED*** return h[148:][:8] ***REMOVED***
func (h *headerV7) TypeFlag() []byte ***REMOVED*** return h[156:][:1] ***REMOVED***
func (h *headerV7) LinkName() []byte ***REMOVED*** return h[157:][:100] ***REMOVED***

type headerGNU [blockSize]byte

func (h *headerGNU) V7() *headerV7       ***REMOVED*** return (*headerV7)(h) ***REMOVED***
func (h *headerGNU) Magic() []byte       ***REMOVED*** return h[257:][:6] ***REMOVED***
func (h *headerGNU) Version() []byte     ***REMOVED*** return h[263:][:2] ***REMOVED***
func (h *headerGNU) UserName() []byte    ***REMOVED*** return h[265:][:32] ***REMOVED***
func (h *headerGNU) GroupName() []byte   ***REMOVED*** return h[297:][:32] ***REMOVED***
func (h *headerGNU) DevMajor() []byte    ***REMOVED*** return h[329:][:8] ***REMOVED***
func (h *headerGNU) DevMinor() []byte    ***REMOVED*** return h[337:][:8] ***REMOVED***
func (h *headerGNU) AccessTime() []byte  ***REMOVED*** return h[345:][:12] ***REMOVED***
func (h *headerGNU) ChangeTime() []byte  ***REMOVED*** return h[357:][:12] ***REMOVED***
func (h *headerGNU) Sparse() sparseArray ***REMOVED*** return (sparseArray)(h[386:][:24*4+1]) ***REMOVED***
func (h *headerGNU) RealSize() []byte    ***REMOVED*** return h[483:][:12] ***REMOVED***

type headerSTAR [blockSize]byte

func (h *headerSTAR) V7() *headerV7      ***REMOVED*** return (*headerV7)(h) ***REMOVED***
func (h *headerSTAR) Magic() []byte      ***REMOVED*** return h[257:][:6] ***REMOVED***
func (h *headerSTAR) Version() []byte    ***REMOVED*** return h[263:][:2] ***REMOVED***
func (h *headerSTAR) UserName() []byte   ***REMOVED*** return h[265:][:32] ***REMOVED***
func (h *headerSTAR) GroupName() []byte  ***REMOVED*** return h[297:][:32] ***REMOVED***
func (h *headerSTAR) DevMajor() []byte   ***REMOVED*** return h[329:][:8] ***REMOVED***
func (h *headerSTAR) DevMinor() []byte   ***REMOVED*** return h[337:][:8] ***REMOVED***
func (h *headerSTAR) Prefix() []byte     ***REMOVED*** return h[345:][:131] ***REMOVED***
func (h *headerSTAR) AccessTime() []byte ***REMOVED*** return h[476:][:12] ***REMOVED***
func (h *headerSTAR) ChangeTime() []byte ***REMOVED*** return h[488:][:12] ***REMOVED***
func (h *headerSTAR) Trailer() []byte    ***REMOVED*** return h[508:][:4] ***REMOVED***

type headerUSTAR [blockSize]byte

func (h *headerUSTAR) V7() *headerV7     ***REMOVED*** return (*headerV7)(h) ***REMOVED***
func (h *headerUSTAR) Magic() []byte     ***REMOVED*** return h[257:][:6] ***REMOVED***
func (h *headerUSTAR) Version() []byte   ***REMOVED*** return h[263:][:2] ***REMOVED***
func (h *headerUSTAR) UserName() []byte  ***REMOVED*** return h[265:][:32] ***REMOVED***
func (h *headerUSTAR) GroupName() []byte ***REMOVED*** return h[297:][:32] ***REMOVED***
func (h *headerUSTAR) DevMajor() []byte  ***REMOVED*** return h[329:][:8] ***REMOVED***
func (h *headerUSTAR) DevMinor() []byte  ***REMOVED*** return h[337:][:8] ***REMOVED***
func (h *headerUSTAR) Prefix() []byte    ***REMOVED*** return h[345:][:155] ***REMOVED***

type sparseArray []byte

func (s sparseArray) Entry(i int) sparseElem ***REMOVED*** return (sparseElem)(s[i*24:]) ***REMOVED***
func (s sparseArray) IsExtended() []byte     ***REMOVED*** return s[24*s.MaxEntries():][:1] ***REMOVED***
func (s sparseArray) MaxEntries() int        ***REMOVED*** return len(s) / 24 ***REMOVED***

type sparseElem []byte

func (s sparseElem) Offset() []byte ***REMOVED*** return s[00:][:12] ***REMOVED***
func (s sparseElem) Length() []byte ***REMOVED*** return s[12:][:12] ***REMOVED***
