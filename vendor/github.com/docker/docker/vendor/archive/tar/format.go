// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tar

// Constants to identify various tar formats.
const (
	// The format is unknown.
	formatUnknown = (1 << iota) / 2 // Sequence of 0, 1, 2, 4, 8, etc...

	// The format of the original Unix V7 tar tool prior to standardization.
	formatV7

	// The old and new GNU formats, which are incompatible with USTAR.
	// This does cover the old GNU sparse extension.
	// This does not cover the GNU sparse extensions using PAX headers,
	// versions 0.0, 0.1, and 1.0; these fall under the PAX format.
	formatGNU

	// Schily's tar format, which is incompatible with USTAR.
	// This does not cover STAR extensions to the PAX format; these fall under
	// the PAX format.
	formatSTAR

	// USTAR is the former standardization of tar defined in POSIX.1-1988.
	// This is incompatible with the GNU and STAR formats.
	formatUSTAR

	// PAX is the latest standardization of tar defined in POSIX.1-2001.
	// This is an extension of USTAR and is "backwards compatible" with it.
	//
	// Some newer formats add their own extensions to PAX, such as GNU sparse
	// files and SCHILY extended attributes. Since they are backwards compatible
	// with PAX, they will be labelled as "PAX".
	formatPAX
)

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
// If the checksum fails, then formatUnknown is returned.
func (b *block) GetFormat() (format int) ***REMOVED***
	// Verify checksum.
	var p parser
	value := p.parseOctal(b.V7().Chksum())
	chksum1, chksum2 := b.ComputeChecksum()
	if p.err != nil || (value != chksum1 && value != chksum2) ***REMOVED***
		return formatUnknown
	***REMOVED***

	// Guess the magic values.
	magic := string(b.USTAR().Magic())
	version := string(b.USTAR().Version())
	trailer := string(b.STAR().Trailer())
	switch ***REMOVED***
	case magic == magicUSTAR && trailer == trailerSTAR:
		return formatSTAR
	case magic == magicUSTAR:
		return formatUSTAR
	case magic == magicGNU && version == versionGNU:
		return formatGNU
	default:
		return formatV7
	***REMOVED***
***REMOVED***

// SetFormat writes the magic values necessary for specified format
// and then updates the checksum accordingly.
func (b *block) SetFormat(format int) ***REMOVED***
	// Set the magic values.
	switch format ***REMOVED***
	case formatV7:
		// Do nothing.
	case formatGNU:
		copy(b.GNU().Magic(), magicGNU)
		copy(b.GNU().Version(), versionGNU)
	case formatSTAR:
		copy(b.STAR().Magic(), magicUSTAR)
		copy(b.STAR().Version(), versionUSTAR)
		copy(b.STAR().Trailer(), trailerSTAR)
	case formatUSTAR, formatPAX:
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

func (s sparseArray) Entry(i int) sparseNode ***REMOVED*** return (sparseNode)(s[i*24:]) ***REMOVED***
func (s sparseArray) IsExtended() []byte     ***REMOVED*** return s[24*s.MaxEntries():][:1] ***REMOVED***
func (s sparseArray) MaxEntries() int        ***REMOVED*** return len(s) / 24 ***REMOVED***

type sparseNode []byte

func (s sparseNode) Offset() []byte   ***REMOVED*** return s[00:][:12] ***REMOVED***
func (s sparseNode) NumBytes() []byte ***REMOVED*** return s[12:][:12] ***REMOVED***
